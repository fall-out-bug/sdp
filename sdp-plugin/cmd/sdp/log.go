package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

const defaultLogPath = ".sdp/log/events.jsonl"
const defaultRecentN = 20

func logCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Evidence log: show and trace events",
		Long:  `Show recent events (last 20) or trace evidence chain by commit/workstream.`,
		RunE:  func(cmd *cobra.Command, args []string) error { return runLogShow("", "") },
	}
	cmd.AddCommand(logShowCmd())
	cmd.AddCommand(logTraceCmd())
	return cmd
}

func logShowCmd() *cobra.Command {
	var eventType, search string
	c := &cobra.Command{
		Use:   "show",
		Short: "Show recent events (last 20)",
		Long:  `Show events. Use --type=decision for decisions only; --search for full-text search.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogShow(eventType, search)
		},
	}
	c.Flags().StringVar(&eventType, "type", "", "Filter by event type (e.g. decision)")
	c.Flags().StringVar(&search, "search", "", "Full-text search in question/choice/rationale")
	return c
}

func runLogShow(eventType, search string) error {
	path, err := evidenceLogPath()
	if err != nil {
		return err
	}
	r := evidence.NewReader(path)
	events, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}
	events = evidence.FilterByType(events, eventType)
	events = evidence.FilterBySearch(events, search)
	events = evidence.LastN(events, defaultRecentN)
	if len(events) == 0 {
		fmt.Println("No events in evidence log.")
		return nil
	}
	fmt.Printf("Recent events (last %d):\n", len(events))
	fmt.Println(evidence.FormatHuman(events))
	return nil
}

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

func evidenceLogPath() (string, error) {
	root, err := config.FindProjectRoot()
	if err != nil {
		return "", fmt.Errorf("find project root: %w", err)
	}
	cfg, _ := config.Load(root)
	logPath := defaultLogPath
	if cfg != nil && cfg.Evidence.LogPath != "" {
		logPath = cfg.Evidence.LogPath
	}
	return filepath.Join(root, logPath), nil
}
