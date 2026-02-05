package main

import (
	"fmt"

	"github.com/ai-masters/sdp/internal/doctor"
	"github.com/spf13/cobra"
)

func doctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check SDP environment",
		Long: `Check that your environment is properly configured for SDP.

Verifies:
  - Git is installed
  - Claude Code CLI is available (optional)
  - Go compiler is available (for building binary)
  - .claude/ directory exists and is properly structured`,
		RunE: func(cmd *cobra.Command, args []string) error {
			results := doctor.Run()

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

	return cmd
}
