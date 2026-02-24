package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func runHydrate(projectRoot, featureID, wsFlag string, cp *orchestrate.Checkpoint, workstreams []string) {
	action, err := orchestrate.ComputeNextAction(cp, workstreams, projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if action.Action == "review" {
		if _, err := orchestrate.HydrateForReview(projectRoot, featureID, cp, workstreams); err != nil {
			fmt.Fprintf(os.Stderr, "error: hydrate: %v\n", err)
			os.Exit(1)
		}
	} else {
		wsID := wsFlag
		if wsID == "" && action.Action == "build" {
			wsID = action.WSID
		}
		if wsID == "" {
			fmt.Fprintf(os.Stderr, "error: cannot hydrate: action=%s, specify --ws\n", action.Action)
			os.Exit(1)
		}
		if _, err := orchestrate.Hydrate(projectRoot, featureID, wsID, cp); err != nil {
			fmt.Fprintf(os.Stderr, "error: hydrate: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("Wrote .sdp/context-packet.json")
}
