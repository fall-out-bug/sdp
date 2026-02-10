package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

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
