package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fall-out-bug/sdp/internal/watcher"
	"github.com/spf13/cobra"
)

func watchCmd() *cobra.Command {
	var (
		quiet     bool
		include   []string
		exclude   []string
		watchPath string
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch files for quality violations",
		Long: `Watch files for quality violations in real-time.

Monitors source files and runs quality checks on save. Detects:
  - Files exceeding 200 LOC
  - Cyclomatic complexity >= 10
  - Type errors (via go vet)
  - Coverage violations (via coverage check)

Press Ctrl+C to stop watching.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine watch path
			if len(args) > 0 {
				watchPath = args[0]
			}
			if watchPath == "" {
				watchPath = "."
			}

			// Create quality watcher
			config := &watcher.QualityWatcherConfig{
				IncludePatterns: include,
				ExcludePatterns: exclude,
				Quiet:           quiet,
			}

			qw, err := watcher.NewQualityWatcher(watchPath, config)
			if err != nil {
				return err
			}
			defer qw.Close()

			// Handle Ctrl+C gracefully
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

			// Start watching in background
			errChan := make(chan error, 1)
			go func() {
				errChan <- qw.Start()
			}()

			// Wait for interrupt
			select {
			case <-sigChan:
				cmd.Println("\nStopping watcher...")
				qw.Stop()

				// Print summary
				violations := qw.GetViolations()
				if len(violations) > 0 {
					cmd.Printf("\nFound %d violations:\n", len(violations))
					for _, v := range violations {
						cmd.Printf("  - %s: %s\n", v.File, v.Message)
					}
				} else {
					cmd.Println("\nNo violations detected!")
				}

				return nil
			case err := <-errChan:
				return err
			}
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
	cmd.Flags().StringSliceVarP(&include, "include", "i", []string{"*.go"}, "Include patterns (glob)")
	cmd.Flags().StringSliceVarP(&exclude, "exclude", "e", []string{"*_test.go", "mock_*.go"}, "Exclude patterns (glob)")

	return cmd
}
