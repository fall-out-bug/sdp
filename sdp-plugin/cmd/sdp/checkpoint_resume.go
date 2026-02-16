package main

import (
	"fmt"
	"time"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

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
