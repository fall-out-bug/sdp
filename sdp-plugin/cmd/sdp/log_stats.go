package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

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
