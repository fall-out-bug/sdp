package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/doctor"
	"github.com/spf13/cobra"
)

func doctorCmd() *cobra.Command {
	var driftCheck bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check SDP environment",
		Long: `Check that your environment is properly configured for SDP.

Verifies:
  - Git is installed
  - Claude Code CLI is available (optional)
  - Go compiler is available (for building binary)
  - .claude/ directory exists and is properly structured
  - Documentation-code drift (with --drift flag)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run checks with drift detection if flag is set
			opts := doctor.RunOptions{
				DriftCheck: driftCheck,
			}
			results := doctor.RunWithOptions(opts)

			// Print results
			fmt.Println("SDP Environment Check")
			fmt.Println("=====================")

			for _, r := range results {
				icon := "✓"
				color := ""
				if r.Status == "warning" {
					icon = "⚠"
					color = " (optional)"
				} else if r.Status == "error" {
					icon = "✗"
				}

				fmt.Printf("%s %s%s\n", icon, r.Name, color)
				fmt.Printf("    %s\n\n", r.Message)
			}

			// Exit code based on results
			hasErrors := false
			for _, r := range results {
				if r.Status == "error" {
					hasErrors = true
				}
			}

			if hasErrors {
				return fmt.Errorf("some required checks failed")
			}

			fmt.Println("All required checks passed!")
			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&driftCheck, "drift", false, "Check for documentation-code drift in recent workstreams")

	return cmd
}
