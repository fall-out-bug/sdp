package main

import (
	"fmt"
	"time"

	"github.com/ai-masters/sdp/internal/checkpoint"
	"github.com/ai-masters/sdp/internal/ui"
	"github.com/spf13/cobra"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Manage checkpoints for long-running features",
	Long: `Checkpoint system for saving and resuming feature execution.

Commands:
  create   Create a new checkpoint
  resume   Resume from an existing checkpoint
  list     List all checkpoints
  clean    Clean old checkpoints

Examples:
  # Create a checkpoint for feature F042
  sdp checkpoint create my-feature F042

  # List all checkpoints
  sdp checkpoint list

  # Resume from checkpoint
  sdp checkpoint resume my-feature

  # Clean checkpoints older than 48 hours
  sdp checkpoint clean --age 48`,
}

var checkpointCreateCmd = &cobra.Command{
	Use:   "create <id> <feature-id>",
	Short: "Create a new checkpoint",
	Example: `  sdp checkpoint create feature-01 F042
  sdp checkpoint create feature-01 F042 --dir /tmp/checkpoints`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		featureID := args[1]

		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return fmt.Errorf("failed to get dir flag: %w", err)
		}

		if dir == "" {
			dir, err = checkpoint.GetDefaultDir()
			if err != nil {
				return fmt.Errorf("failed to get default checkpoint directory: %w", err)
			}
		}

		manager := checkpoint.NewManager(dir)

		cp := checkpoint.Checkpoint{
			ID:                  id,
			FeatureID:           featureID,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			Status:              checkpoint.StatusPending,
			CurrentWorkstream:   "",
			CompletedWorkstreams: []string{},
			Metadata:            map[string]interface{}{},
		}

		if err := manager.Save(cp); err != nil {
			return fmt.Errorf("failed to create checkpoint: %w", err)
		}

		ui.SuccessLine("Checkpoint created: %s", ui.BoldText(id))
		fmt.Printf("   Feature:  %s\n", featureID)
		fmt.Printf("   Location: %s/%s.json\n", ui.Dim(dir), id)

		return nil
	},
}

var checkpointResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume from an existing checkpoint",
	Example: `  sdp checkpoint resume feature-01
  sdp checkpoint resume feature-01 --dir /tmp/checkpoints`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return fmt.Errorf("failed to get dir flag: %w", err)
		}

		if dir == "" {
			dir, err = checkpoint.GetDefaultDir()
			if err != nil {
				return fmt.Errorf("failed to get default checkpoint directory: %w", err)
			}
		}

		manager := checkpoint.NewManager(dir)

		cp, err := manager.Resume(id)
		if err != nil {
			return fmt.Errorf("failed to resume checkpoint: %w", err)
		}

		ui.SuccessLine("Resumed checkpoint: %s", ui.BoldText(cp.ID))
		fmt.Printf("   Feature:              %s\n", cp.FeatureID)
		fmt.Printf("   Status:               %s\n", ui.Info(string(cp.Status)))
		fmt.Printf("   Current Workstream:   %s\n", cp.CurrentWorkstream)
		fmt.Printf("   Completed Workstreams: %d\n", len(cp.CompletedWorkstreams))
		fmt.Printf("   Created:              %s\n", ui.Dim(cp.CreatedAt.Format(time.RFC3339)))
		fmt.Printf("   Updated:              %s\n", ui.Dim(cp.UpdatedAt.Format(time.RFC3339)))

		return nil
	},
}

var checkpointListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all checkpoints",
	Example: `  sdp checkpoint list
  sdp checkpoint list --dir /tmp/checkpoints`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return fmt.Errorf("failed to get dir flag: %w", err)
		}

		if dir == "" {
			dir, err = checkpoint.GetDefaultDir()
			if err != nil {
				return fmt.Errorf("failed to get default checkpoint directory: %w", err)
			}
		}

		manager := checkpoint.NewManager(dir)

		checkpoints, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list checkpoints: %w", err)
		}

		if len(checkpoints) == 0 {
			ui.InfoLine("No checkpoints found")
			return nil
		}

		ui.Header(fmt.Sprintf("Found %d checkpoint(s)", len(checkpoints)))
		for _, cp := range checkpoints {
			fmt.Printf("ID:        %s\n", ui.BoldText(cp.ID))
			fmt.Printf("  Feature:  %s\n", cp.FeatureID)
			fmt.Printf("  Status:   %s\n", ui.Info(string(cp.Status)))
			fmt.Printf("  Current:  %s\n", cp.CurrentWorkstream)
			fmt.Printf("  Progress: %d/%d workstreams\n",
				len(cp.CompletedWorkstreams),
				len(cp.CompletedWorkstreams)+1) // +1 for current
			fmt.Printf("  Created:  %s\n", ui.Dim(cp.CreatedAt.Format(time.RFC3339)))
			fmt.Printf("  Updated:  %s\n", ui.Dim(cp.UpdatedAt.Format(time.RFC3339)))
			fmt.Println()
		}

		return nil
	},
}

var checkpointCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean old checkpoints",
	Long: `Remove checkpoints older than the specified age.

This command permanently deletes checkpoint files that have not been modified
within the specified time period. Use with caution.`,
	Example: `  # Clean checkpoints older than 24 hours (default)
  sdp checkpoint clean

  # Clean checkpoints older than 48 hours
  sdp checkpoint clean --age 48

  # Clean checkpoints older than 7 days
  sdp checkpoint clean --age 168`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return fmt.Errorf("failed to get dir flag: %w", err)
		}

		if dir == "" {
			dir, err = checkpoint.GetDefaultDir()
			if err != nil {
				return fmt.Errorf("failed to get default checkpoint directory: %w", err)
			}
		}

		ageHours, err := cmd.Flags().GetInt("age")
		if err != nil {
			return fmt.Errorf("failed to get age flag: %w", err)
		}

		manager := checkpoint.NewManager(dir)

		age := time.Duration(ageHours) * time.Hour
		deleted, err := manager.Clean(age)
		if err != nil {
			return fmt.Errorf("failed to clean checkpoints: %w", err)
		}

		if deleted == 0 {
			ui.InfoLine("No old checkpoints to clean")
		} else {
			ui.SuccessLine("Cleaned %d old checkpoint(s)", deleted)
		}

		return nil
	},
}

func init() {
	checkpointCmd.PersistentFlags().String("dir", "", "Checkpoint directory (default: .sdp/checkpoints)")
	checkpointCleanCmd.Flags().Int("age", 24, "Age in hours (default: 24)")

	checkpointCmd.AddCommand(checkpointCreateCmd)
	checkpointCmd.AddCommand(checkpointResumeCmd)
	checkpointCmd.AddCommand(checkpointListCmd)
	checkpointCmd.AddCommand(checkpointCleanCmd)
}
