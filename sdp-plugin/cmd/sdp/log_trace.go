package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

func logTraceCmd() *cobra.Command {
	var wsID string
	var jsonOut bool
	var verify bool
	c := &cobra.Command{
		Use:   "trace [commit-sha]",
		Short: "Trace evidence chain by commit or workstream",
		Long:  `Filter events by commit_sha or --ws <ws-id>. Use --verify to check chain integrity.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			commitSHA := ""
			if len(args) > 0 {
				commitSHA = args[0]
			}
			return runLogTrace(commitSHA, wsID, jsonOut, verify)
		},
	}
	c.Flags().StringVar(&wsID, "ws", "", "Filter by workstream ID")
	c.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")
	c.Flags().BoolVar(&verify, "verify", false, "Verify hash chain integrity")
	return c
}

func runLogTrace(commitSHA, wsID string, jsonOut, verify bool) error {
	path, err := evidenceLogPath()
	if err != nil {
		return err
	}
	r := evidence.NewReader(path)
	events, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}
	if verify {
		if err := r.Verify(); err != nil {
			fmt.Fprintf(os.Stderr, "Chain integrity: ✗ %v\n", err)
			return err
		}
	}
	events = evidence.FilterByCommit(events, commitSHA)
	events = evidence.FilterByWS(events, wsID)
	if len(events) == 0 {
		fmt.Println("No matching events.")
		return nil
	}
	if jsonOut {
		out, err := evidence.FormatJSON(events)
		if err != nil {
			return err
		}
		fmt.Println(out)
		return nil
	}
	if commitSHA != "" {
		fmt.Printf("Evidence trail for commit %s:\n\n", commitSHA)
	} else if wsID != "" {
		fmt.Printf("Evidence trail for WS %s:\n\n", wsID)
	} else {
		fmt.Printf("Evidence trail:\n\n")
	}
	fmt.Println(evidence.FormatHuman(events))
	if verify {
		fmt.Printf("Chain integrity: ✓ valid (%d events, no breaks)\n", len(events))
	} else {
		fmt.Printf("%d events\n", len(events))
	}
	return nil
}
