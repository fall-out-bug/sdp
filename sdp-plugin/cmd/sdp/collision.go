package main

import (
	"github.com/spf13/cobra"
)

func collisionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collision",
		Short: "Scope collision detection for parallel workstreams",
		Long:  `Detect when in-progress workstreams touch the same files or directories.`,
	}
	cmd.AddCommand(collisionCheckCmd())
	cmd.AddCommand(collisionDetectCmd())
	return cmd
}

func collisionCheckCmd() *cobra.Command {
	var deep bool
	cmd := &cobra.Command{
		Use:   "check",
		Short: "List scope overlaps across in-progress workstreams",
		RunE:  runCollisionCheck,
	}
	cmd.Flags().BoolVar(&deep, "deep", false, "Analyze interface boundaries (shared types/structs)")
	return cmd
}
