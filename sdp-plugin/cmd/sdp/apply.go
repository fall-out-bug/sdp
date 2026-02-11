package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/executor"
	"github.com/spf13/cobra"
)

func applyCmd() *cobra.Command {
	var dryRun bool
	var retryCount int
	var outputFormat string
	var specificWS string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Execute workstreams from the terminal",
		Long: `Execute workstreams with streaming progress reporting.

Modes:
  - Default: Execute all ready workstreams (no unresolved blockers)
  - --ws <id>: Execute specific workstream by ID
  - --dry-run: Show execution plan without running
  - --retry <n>: Retry failed workstreams up to N times
  - --output=json: Machine-readable JSON progress events

Progress output:
  - Human (default): [00-054-01] ██████░░░░░░ 50% — running tests
  - JSON: {"ws_id":"00-054-01","status":"running","progress":50,...}

Examples:
  sdp apply
  sdp apply --ws 00-054-01
  sdp apply --retry 3
  sdp apply --dry-run
  sdp apply --output=json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find project root
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			// Setup paths
			backlogDir := filepath.Join(root, defaultBacklogDir)
			logPath := filepath.Join(root, ".sdp/log/events.jsonl")

			// Validate backlog directory exists
			if _, err := os.Stat(backlogDir); os.IsNotExist(err) {
				return fmt.Errorf("backlog directory not found: %s\nRun 'sdp plan' first to create workstreams", backlogDir)
			}

			// Create executor (evidence logging is handled internally)
			exec := executor.NewExecutor(executor.ExecutorConfig{
				BacklogDir:      backlogDir,
				DryRun:          dryRun,
				RetryCount:      retryCount,
				EvidenceLogPath: logPath,
			})
			exec.SetOutputFormat(outputFormat)

			// Create context for cancellation
			ctx := context.Background()

			// Determine execution options
			opts := executor.ExecuteOptions{
				All:        specificWS == "",
				SpecificWS: specificWS,
				Retry:      retryCount,
				Output:     outputFormat,
			}

			// Execute workstreams
			fmt.Println("SDP Workstream Executor")
			fmt.Println("======================")
			if dryRun {
				fmt.Println("DRY RUN MODE - No changes will be made")
				fmt.Println()
			}

			result, err := exec.Execute(ctx, os.Stdout, opts)
			if err != nil {
				return fmt.Errorf("execution failed: %w", err)
			}

			// Show summary
			if dryRun {
				fmt.Println("\nDry run complete. To execute, run:")
				if specificWS != "" {
					fmt.Printf("  sdp apply --ws %s\n", specificWS)
				} else {
					fmt.Println("  sdp apply")
				}
			} else if result.Failed > 0 {
				fmt.Printf("\n⚠ %d workstream(s) failed. Check logs for details.\n", result.Failed)
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show execution plan without running")
	cmd.Flags().IntVar(&retryCount, "retry", 1, "Retry failed workstreams up to N times")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "human", "Output format: human or json")
	cmd.Flags().StringVarP(&specificWS, "ws", "w", "", "Execute specific workstream by ID")

	return cmd
}
