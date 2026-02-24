package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func runNextAction(cp *orchestrate.Checkpoint, workstreams []string, projectRoot string) {
	action, err := orchestrate.ComputeNextAction(cp, workstreams, projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(action); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
