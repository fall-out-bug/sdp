package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/verify"
	"github.com/spf13/cobra"
)

func verifyCmd() *cobra.Command {
	var wsDir string

	cmd := &cobra.Command{
		Use:   "verify <ws-id>",
		Short: "Verify workstream completion with evidence",
		Long: `Verify workstream completion by checking:
  - All scope_files output exist
  - All Verification commands pass
  - Test coverage meets threshold

Usage:
  sdp verify 00-001-01`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsID := args[0]

			// Default workstream directory
			if wsDir == "" {
				wsDir = "docs/workstreams"
			}

			// Create verifier
			verifier := verify.NewVerifier(wsDir)

			// Run verification
			result := verifier.Verify(wsID)

			// AC1: Emit lesson event when workstream completes (auto-extracted)
			lesson := evidence.ExtractLesson(wsID, result)
			evidence.EmitLesson(lesson)

			// Print results
			if result.Passed {
				fmt.Printf("✅ Workstream %s verification PASSED\n", wsID)
			} else {
				fmt.Printf("❌ Workstream %s verification FAILED\n", wsID)
			}

			fmt.Printf("\nChecks run: %d\n", len(result.Checks))
			fmt.Printf("Duration: %v\n", result.Duration)

			for _, check := range result.Checks {
				status := "✅"
				if !check.Passed {
					status = "❌"
				}
				fmt.Printf("  %s %s: %s\n", status, check.Name, check.Message)
				if check.Evidence != "" {
					fmt.Printf("     Evidence: %s\n", check.Evidence)
				}
			}

			if result.CoverageActual > 0 {
				fmt.Printf("\nCoverage: %.1f%%\n", result.CoverageActual)
			}

			if len(result.MissingFiles) > 0 {
				fmt.Printf("\nMissing files (%d):\n", len(result.MissingFiles))
				for _, f := range result.MissingFiles {
					fmt.Printf("  - %s\n", f)
				}
			}

			if len(result.FailedCommands) > 0 {
				fmt.Printf("\nFailed commands (%d):\n", len(result.FailedCommands))
				for _, cmd := range result.FailedCommands {
					fmt.Printf("  - %s\n", cmd)
				}
			}

			// Exit with error code if failed
			if !result.Passed {
				os.Exit(1)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&wsDir, "ws-dir", "", "Workstream directory (default: docs/workstreams)")

	return cmd
}
