package main

import (
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

func init() {
	checkpointCmd.PersistentFlags().String("dir", "", "Checkpoint directory (default: .sdp/checkpoints)")
	checkpointCleanCmd.Flags().Int("age", 24, "Age in hours (default: 24)")

	checkpointCmd.AddCommand(checkpointCreateCmd)
	checkpointCmd.AddCommand(checkpointResumeCmd)
	checkpointCmd.AddCommand(checkpointListCmd)
	checkpointCmd.AddCommand(checkpointCleanCmd)
}
