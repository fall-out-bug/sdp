package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardFindingListCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List review findings",
		Long: `List findings from @review skill.

By default shows only open findings. Use --all to include resolved.`,
		Example: `  sdp guard finding list
  sdp guard finding list --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			sm := guard.NewStateManager(configDir)

			state, err := sm.Load()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			if len(state.ReviewFindings) == 0 {
				fmt.Println("No review findings")
				return nil
			}

			open, resolved, blocking := state.FindingCount()
			fmt.Printf("Review Findings: %d open (%d blocking), %d resolved\n\n", open, blocking, resolved)

			for _, f := range state.ReviewFindings {
				if !all && f.Status == "resolved" {
					continue
				}

				status := " "
				if f.Status == "resolved" {
					status = "✓"
				} else if f.Priority <= 1 {
					status = "⚠"
				}

				fmt.Printf("%s [%s] P%d %s: %s\n", status, f.ReviewArea, f.Priority, f.FeatureID, f.Title)
				if f.BeadsID != "" {
					fmt.Printf("  → Beads: %s\n", f.BeadsID)
				}
				if f.Status == "resolved" {
					fmt.Printf("  → Resolved: %s\n", f.ResolvedBy)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Show all findings including resolved")

	return cmd
}
