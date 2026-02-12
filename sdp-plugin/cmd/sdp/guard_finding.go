package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "finding",
		Short:   "Manage review findings",
		Long:    `Register and manage findings from @review skill.`,
		Aliases: []string{"findings"},
	}
	cmd.AddCommand(guardFindingAddCmd())
	cmd.AddCommand(guardFindingListCmd())
	cmd.AddCommand(guardFindingResolveCmd())
	cmd.AddCommand(guardFindingClearCmd())
	return cmd
}

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

func guardFindingResolveCmd() *cobra.Command {
	var resolvedBy string

	cmd := &cobra.Command{
		Use:   "resolve <finding-id>",
		Short: "Mark a finding as resolved",
		Long: `Mark a review finding as resolved.

The resolvedBy should describe how the finding was addressed.`,
		Example: `  sdp guard finding resolve finding-123 --by="Fixed in commit abc123"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			sm := guard.NewStateManager(configDir)

			state, err := sm.Load()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			if !state.ResolveFinding(args[0], resolvedBy) {
				return fmt.Errorf("finding not found: %s", args[0])
			}

			if err := sm.Save(*state); err != nil {
				return fmt.Errorf("failed to save state: %w", err)
			}

			fmt.Printf("✓ Resolved finding: %s\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringVar(&resolvedBy, "by", "manual", "How the finding was resolved")
	cmd.MarkFlagRequired("by")

	return cmd
}

func guardFindingClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all resolved findings",
		Long:  `Remove all resolved findings from the state file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			sm := guard.NewStateManager(configDir)

			state, err := sm.Load()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			var open []guard.ReviewFinding
			for _, f := range state.ReviewFindings {
				if f.Status != "resolved" {
					open = append(open, f)
				}
			}

			cleared := len(state.ReviewFindings) - len(open)
			state.ReviewFindings = open

			if err := sm.Save(*state); err != nil {
				return fmt.Errorf("failed to save state: %w", err)
			}

			fmt.Printf("✓ Cleared %d resolved findings\n", cleared)
			return nil
		},
	}

	return cmd
}

func getConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir, _ = os.UserConfigDir()
	}
	return filepath.Join(configDir, "sdp")
}
