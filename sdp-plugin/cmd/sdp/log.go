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
		Short: "Evidence log: show, export, and trace events",
		Long:  `Show recent events, export logs, or trace evidence chain by commit/workstream.`,
		RunE:  func(cmd *cobra.Command, args []string) error { return runLogShow("", "") },
	}
	cmd.AddCommand(logShowCmd())
	cmd.AddCommand(logExportCmd())
	cmd.AddCommand(logStatsCmd())
	cmd.AddCommand(logTraceCmd())
	return cmd
}

func logShowCmd() *cobra.Command {
	var eventType, model, since, wsID string
	var page int
	c := &cobra.Command{
		Use:   "show",
		Short: "Show paginated events with filters",
		Long:  `Show events with pagination and filters. Supports --type, --model, --since, --ws, and --page.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogShowFiltered(eventType, model, since, wsID, page)
		},
	}
	c.Flags().StringVar(&eventType, "type", "", "Filter by event type (e.g. generation)")
	c.Flags().StringVar(&model, "model", "", "Filter by model ID (e.g. claude-sonnet-4)")
	c.Flags().StringVar(&since, "since", "", "Filter by date (RFC3339 format, e.g. 2026-02-01T00:00:00Z)")
	c.Flags().StringVar(&wsID, "ws", "", "Filter by workstream ID (e.g. 00-054-03)")
	c.Flags().IntVarP(&page, "page", "p", 1, "Page number (1-indexed, 20 events per page)")
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

func runLogShowFiltered(eventType, model, since, wsID string, page int) error {
	path, err := evidenceLogPath()
	if err != nil {
		return err
	}
	r := evidence.NewReader(path)
	events, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}

	browser := evidence.NewBrowser(events)

	// Apply filters
	if eventType != "" {
		events = browser.FilterByType(eventType)
	}
	if model != "" {
		events = browser.FilterByModel(model)
	}
	if since != "" {
		events = browser.FilterBySince(since)
	}
	if wsID != "" {
		events = browser.FilterByWS(wsID)
	}

	if len(events) == 0 {
		fmt.Println("No matching events.")
		return nil
	}

	// Paginate (20 per page)
	pageSize := 20
	pagedEvents, total := evidence.NewBrowser(events).Page(page, pageSize)

	if len(pagedEvents) == 0 {
		fmt.Printf("Page %d is empty (showing page 1 of %d)\n", page, (total+pageSize-1)/pageSize)
		pagedEvents, _ = evidence.NewBrowser(events).Page(1, pageSize)
	}

	fmt.Printf("Events (page %d of %d, %d total):\n", page, (total+pageSize-1)/pageSize, total)
	fmt.Println(evidence.FormatHuman(pagedEvents))
	return nil
}

func logExportCmd() *cobra.Command {
	var format string
	c := &cobra.Command{
		Use:   "export",
		Short: "Export events as CSV or JSON",
		Long:  `Export events for analysis. Supports --format=csv or --format=json.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogExport(format)
		},
	}
	c.Flags().StringVar(&format, "format", "csv", "Export format: csv or json")
	return c
}

func logStatsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "stats",
		Short: "Show event statistics",
		Long:  `Show summary statistics (event counts by type, model distribution).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogStats()
		},
	}
	return c
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

func runLogExport(format string) error {
	path, err := evidenceLogPath()
	if err != nil {
		return err
	}
	r := evidence.NewReader(path)
	events, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}
	if len(events) == 0 {
		fmt.Println("No events to export.")
		return nil
	}

	exporter := evidence.NewExporter()
	var output string
	if format == "json" {
		output, err = exporter.ToJSON(events)
	} else {
		output, err = exporter.ToCSV(events)
	}
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}
	fmt.Println(output)
	return nil
}

func runLogStats() error {
	path, err := evidenceLogPath()
	if err != nil {
		return err
	}
	r := evidence.NewReader(path)
	events, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}
	if len(events) == 0 {
		fmt.Println("No events in evidence log.")
		return nil
	}

	exporter := evidence.NewExporter()
	stats := exporter.Stats(events)

	fmt.Printf("Event Statistics (%d total):\n\n", stats.Total)

	fmt.Println("By Type:")
	for _, eventType := range []string{"plan", "generation", "verification", "approval", "decision", "lesson"} {
		count := stats.CountByType[eventType]
		if count > 0 {
			fmt.Printf("  %s: %d\n", eventType, count)
		}
	}

	if len(stats.ModelDistribution) > 0 {
		fmt.Println("\nBy Model:")
		for model, count := range stats.ModelDistribution {
			fmt.Printf("  %s: %d\n", model, count)
		}
	}

	if len(stats.DateDistribution) > 0 {
		fmt.Println("\nBy Date:")
		for date, count := range stats.DateDistribution {
			fmt.Printf("  %s: %d\n", date, count)
		}
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
