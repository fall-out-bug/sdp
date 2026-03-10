package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/nextstep"
)

func runTextStatus(jsonOutput bool) error {
	data, err := collectControlTowerData()
	if err != nil {
		return err
	}

	if jsonOutput {
		return printStatusJSON(data.StatusView)
	}
	return printStatusText(data.StatusView)
}

func printStatusText(view *nextstep.StatusView) error {
	fmt.Println("SDP Project Status")
	fmt.Println("=================")
	fmt.Println()

	fmt.Println("Environment:")
	fmt.Printf("  Git:          %s\n", boolIcon(view.HasGit))
	fmt.Printf("  Claude:       %s\n", boolIcon(view.HasClaude))
	fmt.Printf("  SDP:          %s\n", boolIcon(view.HasSDP))
	fmt.Printf("  Beads:        %s\n", boolIcon(view.HasBeads))
	if view.Environment.GitBranch != "" {
		fmt.Printf("  Branch:       %s\n", view.Environment.GitBranch)
	}
	if view.Environment.Uncommitted {
		fmt.Printf("  Uncommitted:  %s\n", boolIcon(false))
	} else {
		fmt.Printf("  Uncommitted:  %s\n", boolIcon(true))
	}
	fmt.Println()

	fmt.Println("Workstreams:")
	fmt.Printf("  Total:        %d\n", view.Workstreams.Total)
	fmt.Printf("  Ready:        %d\n", len(view.Workstreams.Ready))
	fmt.Printf("  In Progress:  %d\n", len(view.Workstreams.InProgress))
	fmt.Printf("  Blocked:      %d\n", len(view.Workstreams.Blocked))
	fmt.Printf("  Failed:       %d\n", len(view.Workstreams.Failed))
	fmt.Printf("  Backlog:      %d\n", view.Workstreams.Backlog)
	fmt.Printf("  Completed:    %d\n", view.Workstreams.Completed)
	fmt.Println()

	if view.ActiveSession != nil {
		fmt.Println("Active Session:")
		if view.ActiveSession.WorkstreamID != "" {
			fmt.Printf("  Workstream:   %s\n", view.ActiveSession.WorkstreamID)
		}
		if view.ActiveSession.FeatureID != "" {
			fmt.Printf("  Feature:      %s\n", view.ActiveSession.FeatureID)
		}
		if view.ActiveSession.ExpectedBranch != "" {
			fmt.Printf("  Branch:       %s\n", view.ActiveSession.ExpectedBranch)
		}
		fmt.Println()
	}

	fmt.Println("Next Action:")
	fmt.Printf("  Command:      %s\n", view.NextAction)
	if view.NextStep != nil {
		fmt.Printf("  Reason:       %s\n", view.NextStep.Reason)
		fmt.Printf("  Confidence:   %.0f%%\n", view.NextStep.Confidence*100)
		fmt.Printf("  Category:     %s\n", view.NextStep.Category)
	}

	return nil
}

func printStatusJSON(view *nextstep.StatusView) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(view)
}
