package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

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
			fmt.Printf("  Created:  %s\n", ui.Dim(cp.CreatedAt.Format("2006-01-02 15:04:05")))
			fmt.Printf("  Updated:  %s\n", ui.Dim(cp.UpdatedAt.Format("2006-01-02 15:04:05")))
			fmt.Println()
		}

		return nil
	},
}
