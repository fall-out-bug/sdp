package main

import (
	"fmt"
	"time"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

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
