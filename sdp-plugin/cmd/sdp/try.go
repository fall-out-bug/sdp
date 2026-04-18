package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/telemetry"
	"github.com/fall-out-bug/sdp/internal/trial"
	"github.com/spf13/cobra"
)

func tryCmd() *cobra.Command {
	var discard bool
	var keep bool

	cmd := &cobra.Command{
		Use:   "try \"task description\"",
		Short: "Try a task on a temporary branch",
		Long: `Execute a bounded task on a temporary branch with zero residue:
  - Creates temporary branch (sdp-try-{timestamp})
  - Executes one bounded task based on description
  - Shows results for review
  - On accept: keeps branch, suggests 'sdp adopt'
  - On discard: deletes branch, returns to original state

This provides a zero-commitment first experience with SDP.`,
		Example: `  # Try a task
  sdp try "Add user authentication"

  # Try and discard if not satisfied
  sdp try "Refactor API" --discard

  # Try and keep for adoption
  sdp try "Add tests" --keep`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()
			taskDescription := args[0]
			projectPath := "."

			// Convert to absolute path
			absPath, err := filepath.Abs(projectPath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Create trial session
			t, err := trial.NewTrial(absPath, taskDescription)
			if err != nil {
				return fmt.Errorf("failed to create trial: %w", err)
			}

			// Verify clean state
			clean, err := t.VerifyClean()
			if err != nil {
				return fmt.Errorf("failed to verify clean state: %w", err)
			}
			if !clean {
				return fmt.Errorf("working directory not clean - commit or stash changes first")
			}

			// Initialize telemetry collector (after clean-state check, UX metrics now go to user config dir)
			uxMetrics, err := telemetry.NewUXMetricsCollector("")
			if err != nil {
				// Don't fail the command if telemetry fails
				fmt.Fprintf(os.Stderr, "Warning: failed to initialize telemetry: %v\n", err)
			}

			// Start trial
			fmt.Printf("Starting trial on branch: %s\n", t.BranchName)
			fmt.Printf("Task: %s\n\n", taskDescription)

			if err := t.Start(); err != nil {
				return fmt.Errorf("failed to start trial: %w", err)
			}

			fmt.Println("✓ Trial branch created")

			// Execute task
			fmt.Println("\nExecuting task...")
			result, err := t.Execute()
			if err != nil {
				// Record discard telemetry on execution failure
				if uxMetrics != nil {
					_ = uxMetrics.RecordTryDiscard("unknown", "execution_failure", 1)
				}
				return fmt.Errorf("execution failed: %w", err)
			}

			// Show results
			fmt.Printf("\nExecution completed in %v\n", result.Duration.Round(time.Second))
			fmt.Printf("Result: %s\n", result.Message)

			// Determine outcome
			var outcome string
			var stepNumber int

			// Handle flags
			if discard {
				fmt.Println("\nDiscarding trial...")
				outcome = "user_discarded"
				stepNumber = 2
				if err := t.Discard(); err != nil {
					return err
				}
			} else if keep {
				fmt.Println("\nKeeping trial...")
				outcome = "user_accepted"
				stepNumber = 2
				if err := t.Accept(); err != nil {
					return err
				}
			} else {
				// Interactive prompt
				fmt.Println("\nWhat would you like to do?")
				fmt.Println("  [1] Accept - Keep branch and adopt changes")
				fmt.Println("  [2] Discard - Delete branch and restore original state")
				fmt.Print("Choice: ")

				reader := bufio.NewReader(os.Stdin)
				choice, _ := reader.ReadString('\n')
				choice = strings.TrimSpace(choice)

				switch choice {
				case "1", "a", "accept":
					outcome = "user_accepted"
					stepNumber = 2
					if err := t.Accept(); err != nil {
						return err
					}
				case "2", "d", "discard":
					outcome = "user_discarded"
					stepNumber = 2
					if err := t.Discard(); err != nil {
						return err
					}
				default:
					fmt.Println("Invalid choice. Discarding trial.")
					outcome = "invalid_choice"
					stepNumber = 2
					if err := t.Discard(); err != nil {
						return err
					}
				}
			}

			// Record telemetry
			if uxMetrics != nil {
				duration := time.Since(startTime)
				if outcome == "user_accepted" && result.Success {
					// Record successful completion
					if err := uxMetrics.RecordTryComplete("unknown", duration); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to record telemetry: %v\n", err)
					}
				} else {
					// Record discard
					if err := uxMetrics.RecordTryDiscard("unknown", outcome, stepNumber); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to record telemetry: %v\n", err)
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&discard, "discard", false, "Discard trial after execution")
	cmd.Flags().BoolVar(&keep, "keep", false, "Keep trial after execution")

	return cmd
}
