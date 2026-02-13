package main

import (
	"fmt"
	"time"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardFindingAddCmd() *cobra.Command {
	var featureID, reviewArea, title, beadsID string
	var priority int

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Register a review finding",
		Long: `Register a finding from @review skill.

This is called automatically by review agents when findings are detected.`,
		Example: `  sdp guard finding add --feature=F051 --area=SRE --title="Missing logging" --priority=1 --beads=sdp-abc123`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			sm := guard.NewStateManager(configDir)

			state, err := sm.Load()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			finding := guard.ReviewFinding{
				ID:         fmt.Sprintf("finding-%d", time.Now().Unix()),
				FeatureID:  featureID,
				ReviewArea: reviewArea,
				Title:      title,
				Priority:   priority,
				BeadsID:    beadsID,
				Status:     "open",
				CreatedAt:  time.Now().Format(time.RFC3339),
			}

			state.AddFinding(finding)

			if err := sm.Save(*state); err != nil {
				return fmt.Errorf("failed to save state: %w", err)
			}

			fmt.Printf("✓ Registered finding: %s\n", finding.ID)
			fmt.Printf("  Feature: %s\n", featureID)
			fmt.Printf("  Area: %s\n", reviewArea)
			fmt.Printf("  Priority: P%d\n", priority)
			if beadsID != "" {
				fmt.Printf("  Beads: %s\n", beadsID)
			}

			if priority <= 1 {
				fmt.Println("\n⚠️  BLOCKING: P0/P1 finding requires resolution before merge")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&featureID, "feature", "", "Feature ID (e.g., F051)")
	cmd.Flags().StringVar(&reviewArea, "area", "", "Review area (QA, Security, DevOps, SRE, TechLead, Documentation)")
	cmd.Flags().StringVar(&title, "title", "", "Finding title")
	cmd.Flags().StringVar(&beadsID, "beads", "", "Beads issue ID if created")
	cmd.Flags().IntVar(&priority, "priority", 2, "Priority (0=P0, 1=P1, 2=P2, 3=P3)")
	cmd.MarkFlagRequired("feature")
	cmd.MarkFlagRequired("area")
	cmd.MarkFlagRequired("title")

	return cmd
}
