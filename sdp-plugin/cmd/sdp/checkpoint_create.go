package main

import (
	"fmt"
	"time"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

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
			ID:                   id,
			FeatureID:            featureID,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
			Status:               checkpoint.StatusPending,
			CurrentWorkstream:    "",
			CompletedWorkstreams: []string{},
			Metadata:             map[string]interface{}{},
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
