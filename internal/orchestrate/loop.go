package orchestrate

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// RunOpenCodeLoop drives the full workflow using opencode as the inner loop.
func RunOpenCodeLoop(projectRoot, featureID, cpPath, runsPath string, cp *Checkpoint, workstreams []string) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-ctx.Done():
			_ = SaveCheckpoint(cpPath, cp) // best-effort so resume does not re-run last phase
			slog.Warn("shutdown", "error", ctx.Err())
			os.Exit(1)
		default:
		}

		action, err := ComputeNextAction(cp, workstreams, projectRoot)
		if err != nil {
			fatal("error: %v", err)
		}
		switch action.Action {
		case "build":
			cpFilePath := filepath.Join(cpPath, featureID+".json")
			hookEnv := HookEnv{WSID: action.WSID, FeatureID: featureID, Phase: "build", CheckpointPath: cpFilePath}
			if err := RunHooks(ctx, projectRoot, "build", "pre", hookEnv, func(msg string) { slog.Info("hook", "msg", msg) }); err != nil {
				fatal("error: pre-build hook: %v", err)
			}
			if _, err := Hydrate(projectRoot, featureID, action.WSID, cp); err != nil {
				slog.Error("hydration failed", "error", err, "ws", action.WSID)
				os.Exit(1)
			}
			phaseCtx, cancel := context.WithTimeout(ctx, buildPhaseTimeout)
			commit, err := RunBuildPhase(phaseCtx, projectRoot, action.Feature, action.WSID)
			cancel()
			if err != nil {
				slog.Error("opencode build failed", "error", err, "ws", action.WSID)
				os.Exit(1)
			}
			pending := 0
			for _, ws := range cp.Workstreams {
				if ws.Status != "done" {
					pending++
				}
			}
			if pending == 1 {
				if err := RunHooks(ctx, projectRoot, "build", "post", hookEnv, func(msg string) { slog.Info("hook", "msg", msg) }); err != nil {
					fatal("error: post-build hook: %v", err)
				}
			}
			if err := Advance(cp, workstreams, commit); err != nil {
				fatal("error: advance: %v", err)
			}
			if err := SaveCheckpoint(cpPath, cp); err != nil {
				fatal("error: save checkpoint: %v", err)
			}
		case "review":
			cpFilePath := filepath.Join(cpPath, featureID+".json")
			hookEnv := HookEnv{FeatureID: action.Feature, Phase: "review", CheckpointPath: cpFilePath}
			if err := RunHooks(ctx, projectRoot, "review", "pre", hookEnv, func(msg string) { slog.Info("hook", "msg", msg) }); err != nil {
				fatal("error: pre-review hook: %v", err)
			}
			if _, err := HydrateForReview(projectRoot, action.Feature, cp, workstreams); err != nil {
				slog.Error("hydration failed", "error", err, "feature", action.Feature)
				os.Exit(1)
			}
			phaseCtx, cancel := context.WithTimeout(ctx, reviewPhaseTimeout)
			approved, err := RunReviewPhase(phaseCtx, projectRoot, action.Feature)
			cancel()
			if err != nil || !approved {
				slog.Error("opencode review failed", "error", err, "approved", approved, "feature", action.Feature)
				os.Exit(1)
			}
			if err := RunHooks(ctx, projectRoot, "review", "post", hookEnv, func(msg string) { slog.Info("hook", "msg", msg) }); err != nil {
				fatal("error: post-review hook: %v", err)
			}
			if err := Advance(cp, workstreams, ""); err != nil {
				fatal("error: advance: %v", err)
			}
			if err := SaveCheckpoint(cpPath, cp); err != nil {
				fatal("error: save checkpoint: %v", err)
			}
		case "pr":
			if err := AdvancePRPhase(ctx, projectRoot, featureID, cpPath, cp); err != nil {
				fatal("error: %v", err)
			}
		case "ci-loop":
			if err := AdvanceCIPhase(ctx, projectRoot, featureID, cpPath, runsPath, cp); err != nil {
				fatal("error: %v", err)
			}
		case "done":
			slog.Info("oneshot complete", "feature", featureID)
			fmt.Println("CI GREEN - @oneshot complete")
			return
		}
	}
}
