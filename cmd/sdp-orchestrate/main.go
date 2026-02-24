package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fall-out-bug/sdp/internal/ciloop"
	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	feature := flag.String("feature", "", "Feature ID (e.g. F016)")
	nextAction := flag.Bool("next-action", false, "Output next action as JSON")
	advance := flag.Bool("advance", false, "Advance to next phase after current action")
	result := flag.String("result", "", "Result for advance (e.g. commit hash for build phase)")
	resume := flag.Bool("resume", false, "Resume from existing checkpoint")
	checkpointDir := flag.String("checkpoint-dir", ".sdp/checkpoints", "Checkpoint directory")
	runsDir := flag.String("runs-dir", ".sdp/runs", "Runs directory")
	runtime := flag.String("runtime", "", "Runtime for LLM phases: opencode (invokes opencode run as subprocess)")
	hydrate := flag.Bool("hydrate", false, "Gather context and write .sdp/context-packet.json (before LLM invocation)")
	ws := flag.String("ws", "", "Workstream ID for --hydrate (default: current build ws from next-action)")
	flag.Parse()

	if *feature == "" {
		fmt.Fprintln(os.Stderr, "error: --feature is required")
		flag.Usage()
		os.Exit(1)
	}

	featureID := strings.ToUpper(*feature)
	if !strings.HasPrefix(featureID, "F") {
		featureID = "F" + featureID
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	projectRoot, err := orchestrate.FindProjectRoot(wd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	workstreams, err := orchestrate.DiscoverWorkstreams(projectRoot, featureID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cpPath := filepath.Join(projectRoot, *checkpointDir)
	runsPath := filepath.Join(projectRoot, *runsDir)

	// Remove orphan .tmp files from previous runs
	ciloop.RemoveOrphanTmpFiles(cpPath, runsPath, filepath.Join(projectRoot, ".sdp"))

	cp, err := orchestrate.LoadCheckpoint(cpPath, featureID)
	if err != nil {
		if *resume || !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		branch, err := orchestrate.CurrentBranch(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		cp = orchestrate.CreateInitialCheckpoint(featureID, branch, workstreams)
		cp.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		if err := os.MkdirAll(cpPath, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := orchestrate.SaveCheckpoint(cpPath, cp); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := orchestrate.EnsureRunFile(runsPath, featureID, cp.Branch); err != nil {
			fmt.Fprintf(os.Stderr, "error: ensure run file: %v\n", err)
			os.Exit(1)
		}
	}

	if *nextAction {
		runNextAction(cp, workstreams, projectRoot)
		return
	}
	if *hydrate {
		runHydrate(projectRoot, featureID, *ws, cp, workstreams)
		return
	}
	if *runtime == "opencode" {
		orchestrate.RunOpenCodeLoop(projectRoot, featureID, cpPath, runsPath, cp, workstreams)
		return
	}
	if *advance {
		runAdvance(projectRoot, featureID, cpPath, runsPath, *result, false, cp, workstreams)
		return
	}

	action, err := orchestrate.ComputeNextAction(cp, workstreams, projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	switch action.Action {
	case "build":
		fmt.Printf("INVOKE: @build %s\n", action.WSID)
	case "review":
		cpFilePath := filepath.Join(cpPath, featureID+".json")
		hookEnv := orchestrate.HookEnv{FeatureID: action.Feature, Phase: "review", CheckpointPath: cpFilePath}
		if err := orchestrate.RunHooks(ctx, projectRoot, "review", "pre", hookEnv, func(msg string) { fmt.Fprintln(os.Stderr, msg) }); err != nil {
			fmt.Fprintf(os.Stderr, "error: pre-review hook: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("INVOKE: @review %s\n", action.Feature)
	case "pr":
		fmt.Println("INVOKE: git push && gh pr create")
	case "ci-loop":
		fmt.Printf("INVOKE: sdp-ci-loop --pr %d --feature %s\n", action.PR, action.Feature)
	case "done":
		fmt.Println("CI GREEN - @oneshot complete")
	}
}
