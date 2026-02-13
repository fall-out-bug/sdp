package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/coordination"
	"github.com/spf13/cobra"
)

func coordinationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coordination",
		Short: "Agent coordination management",
		Long: `Manage agent coordination events.

The coordination system tracks agent activities via event sourcing
with hash-chain verification for integrity.

Examples:
  sdp coordination events              # List all events
  sdp coordination events --agent=a1   # Filter by agent
  sdp coordination stats               # Show aggregated stats
  sdp coordination verify              # Verify hash chain`,
	}

	cmd.AddCommand(coordinationEventsCmd())
	cmd.AddCommand(coordinationStatsCmd())
	cmd.AddCommand(coordinationVerifyCmd())

	return cmd
}

func coordinationEventsCmd() *cobra.Command {
	var agentID string
	var taskID string
	var dbPath string

	cmd := &cobra.Command{
		Use:   "events",
		Short: "List agent coordination events",
		Long:  `List all coordination events, optionally filtered by agent or task.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				dbPath = ".sdp/coordination.jsonl"
			}
			if !filepath.IsAbs(dbPath) {
				cwd, _ := os.Getwd()
				dbPath = filepath.Join(cwd, dbPath)
			}

			store, err := coordination.NewStore(dbPath)
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}
			defer store.Close()

			var events []*coordination.AgentEvent
			if agentID != "" {
				events, err = store.FilterByAgent(agentID)
			} else if taskID != "" {
				events, err = store.FilterByTask(taskID)
			} else {
				events, err = store.ReadAll()
			}
			if err != nil {
				return fmt.Errorf("failed to read events: %w", err)
			}

			if len(events) == 0 {
				fmt.Println("No events found")
				return nil
			}

			fmt.Printf("Found %d events:\n\n", len(events))
			for i, e := range events {
				fmt.Printf("%d. [%s] %s\n", i+1, e.Timestamp.Format("2006-01-02 15:04:05"), e.Type)
				fmt.Printf("   Agent: %s (Role: %s)\n", e.AgentID, e.Role)
				if e.TaskID != "" {
					fmt.Printf("   Task: %s\n", e.TaskID)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&agentID, "agent", "", "Filter by agent ID")
	cmd.Flags().StringVar(&taskID, "task", "", "Filter by task ID")
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: .sdp/coordination.jsonl)")

	return cmd
}

func coordinationStatsCmd() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show aggregated coordination statistics",
		Long:  `Display aggregated statistics about agent coordination events.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				dbPath = ".sdp/coordination.jsonl"
			}
			if !filepath.IsAbs(dbPath) {
				cwd, _ := os.Getwd()
				dbPath = filepath.Join(cwd, dbPath)
			}

			store, err := coordination.NewStore(dbPath)
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}
			defer store.Close()

			stats, err := store.GetAggregatedStats()
			if err != nil {
				return fmt.Errorf("failed to get stats: %w", err)
			}

			fmt.Println("Agent Coordination Statistics")
			fmt.Println("=============================")
			fmt.Printf("Total Events: %d\n\n", stats.TotalEvents)

			fmt.Println("By Type:")
			for t, c := range stats.ByType {
				fmt.Printf("  %s: %d\n", t, c)
			}

			if len(stats.ByAgent) > 0 {
				fmt.Println("\nBy Agent:")
				for a, c := range stats.ByAgent {
					fmt.Printf("  %s: %d\n", a, c)
				}
			}

			if len(stats.ByRole) > 0 {
				fmt.Println("\nBy Role:")
				for r, c := range stats.ByRole {
					fmt.Printf("  %s: %d\n", r, c)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: .sdp/coordination.jsonl)")

	return cmd
}

func coordinationVerifyCmd() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify hash chain integrity",
		Long:  `Verify the hash chain integrity of the coordination event log.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if dbPath == "" {
				dbPath = ".sdp/coordination.jsonl"
			}
			if !filepath.IsAbs(dbPath) {
				cwd, _ := os.Getwd()
				dbPath = filepath.Join(cwd, dbPath)
			}

			store, err := coordination.NewStore(dbPath)
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}
			defer store.Close()

			if err := store.VerifyHashChain(); err != nil {
				fmt.Printf("Hash chain verification FAILED: %v\n", err)
				return err
			}

			fmt.Println("Hash chain verification PASSED")
			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: .sdp/coordination.jsonl)")

	return cmd
}
