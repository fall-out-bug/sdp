package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func runAdvance(projectRoot, featureID, cpPath, runsPath, result string, skipGuard bool, cp *orchestrate.Checkpoint, workstreams []string) {
	advanceCtx, advanceStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer advanceStop()

	if cp.Phase == orchestrate.PhasePR {
		if err := orchestrate.AdvancePRPhase(advanceCtx, projectRoot, featureID, cpPath, cp); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if cp.Phase == orchestrate.PhaseCI {
		if err := orchestrate.AdvanceCIPhase(advanceCtx, projectRoot, featureID, cpPath, runsPath, cp); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if cp.Phase == orchestrate.PhaseBuild && result != "" && !skipGuard {
		wsID := orchestrate.CurrentBuildWS(cp)
		if wsID != "" {
			if err := orchestrate.RunGuardCheck(projectRoot, wsID); err != nil {
				var scopeErr *orchestrate.ScopeViolationError
				if errors.As(err, &scopeErr) {
					fmt.Fprintf(os.Stderr, "SCOPE VIOLATION: %s\n", err)
					if createErr := orchestrate.CreateScopeEscalationBead(scopeErr.WSID, scopeErr.Violations); createErr != nil {
						fmt.Fprintf(os.Stderr, "warning: bd create failed: %v\n", createErr)
					}
				}
				fmt.Fprintf(os.Stderr, "error: advance blocked by scope guard: %v\n", err)
				os.Exit(1)
			}
		}
	}
	cpFilePath := filepath.Join(cpPath, featureID+".json")
	hookEnv := orchestrate.HookEnv{
		WSID:           orchestrate.CurrentBuildWS(cp),
		FeatureID:      featureID,
		Phase:          cp.Phase,
		CheckpointPath: cpFilePath,
	}
	logHook := func(msg string) { fmt.Fprintln(os.Stderr, msg) }
	switch cp.Phase {
	case orchestrate.PhaseInit:
		if err := orchestrate.RunHooks(advanceCtx, projectRoot, "build", "pre", hookEnv, logHook); err != nil {
			fmt.Fprintf(os.Stderr, "error: pre-build hook: %v\n", err)
			os.Exit(1)
		}
	case orchestrate.PhaseBuild:
		if err := orchestrate.RunHooks(advanceCtx, projectRoot, "build", "post", hookEnv, logHook); err != nil {
			fmt.Fprintf(os.Stderr, "error: post-build hook: %v\n", err)
			os.Exit(1)
		}
	case orchestrate.PhaseReview:
		if err := orchestrate.RunHooks(advanceCtx, projectRoot, "review", "post", hookEnv, logHook); err != nil {
			fmt.Fprintf(os.Stderr, "error: post-review hook: %v\n", err)
			os.Exit(1)
		}
	}
	// Evaluate OPA policies at phase transition (before advancing).
	// Blocking mode halts; advisory mode logs and continues.
	changedFiles := orchestrate.GetChangedFiles(projectRoot)
	scopeViolations := 0
	policyInput := orchestrate.BuildPolicyInput(cp, scopeViolations, changedFiles)
	policyResult, policyErr := orchestrate.EvaluatePolicies(projectRoot, policyInput)
	if policyErr != nil {
		fmt.Fprintf(os.Stderr, "warning: policy evaluation error: %v\n", policyErr)
	} else {
		for _, w := range policyResult.Warnings {
			fmt.Fprintf(os.Stderr, "POLICY WARN: %s\n", w)
		}
		if len(policyResult.Denials) > 0 {
			for _, d := range policyResult.Denials {
				fmt.Fprintf(os.Stderr, "POLICY DENY [%s]: %s\n", policyResult.Level, d)
			}
			if policyResult.Level == "blocking" {
				fmt.Fprintf(os.Stderr, "error: advance blocked by %d policy denial(s)\n", len(policyResult.Denials))
				os.Exit(1)
			}
		}
	}

	// Validate FSM transition before advancing.
	if err := orchestrate.ValidateAdvance(cp, workstreams); err != nil {
		fmt.Fprintf(os.Stderr, "error: FSM conformance violation: %v\n", err)
		fmt.Fprintf(os.Stderr, "Halting to prevent protocol violation. Fix the issue and retry.\n")
		os.Exit(1)
	}

	prevPhase := cp.Phase
	if err := orchestrate.Advance(cp, workstreams, result); err != nil {
		fmt.Fprintf(os.Stderr, "error: advance: %v\n", err)
		os.Exit(1)
	}
	if err := orchestrate.SaveCheckpoint(cpPath, cp); err != nil {
		fmt.Fprintf(os.Stderr, "error: save checkpoint: %v\n", err)
		os.Exit(1)
	}

	// Generate in-toto attestation on key phase transitions.
	// Written to .sdp/evidence/FXXX.json — updated at each step.
	shouldAttest := prevPhase == orchestrate.PhaseBuild ||
		prevPhase == orchestrate.PhaseReview ||
		cp.Phase == orchestrate.PhaseDone
	if shouldAttest {
		if err := orchestrate.WriteOrchestratorAttestation(projectRoot, cp); err != nil {
			// Non-fatal: log warning but don't block
			fmt.Fprintf(os.Stderr, "warning: attestation generation failed: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "attestation updated: .sdp/evidence/%s.json\n", featureID)
		}
	}
}
