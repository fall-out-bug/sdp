package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/context"
	"github.com/spf13/cobra"
)

func guardContextFindCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find <feature-id>",
		Short: "Locate worktree for a feature",
		Long: `Find the worktree path for a given feature ID.

Uses hybrid recovery strategy:
1. Search session files
2. Parse git worktree list
3. Check workstream metadata`,
		Example: `  sdp guard context find F065
  cd $(sdp guard context find F065)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			featureID := args[0]
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			recovery := context.NewRecovery(root)
			path, err := recovery.FindWorktree(featureID)
			if err != nil {
				return err
			}

			fmt.Println(path)
			return nil
		},
	}
}

func guardContextGoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "go <feature-id>",
		Short: "Print command to change to feature worktree",
		Long: `Print the path and command to change to a feature worktree.

NOTE: This command cannot actually change your shell's CWD.
It outputs the path and instructions for you to execute.`,
		Example: `  sdp guard context go F065`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			featureID := args[0]
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			recovery := context.NewRecovery(root)
			path, err := recovery.GoToWorktree(featureID)
			if err != nil {
				return err
			}

			fmt.Printf("Worktree path: %s\n", path)
			fmt.Printf("\nTo change directory, run:\n")
			fmt.Printf("  cd %s\n", path)

			return nil
		},
	}
}

func guardContextCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "clean",
		Short:   "Clean up stale session files",
		Long:    `Remove invalid or stale session files from all worktrees.`,
		Example: `  sdp guard context clean`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			recovery := context.NewRecovery(root)
			cleaned, err := recovery.Clean()
			if err != nil {
				return fmt.Errorf("clean failed: %w", err)
			}

			if len(cleaned) == 0 {
				fmt.Println("No stale sessions found")
			} else {
				fmt.Printf("Cleaned %d stale session(s):\n", len(cleaned))
				for _, path := range cleaned {
					fmt.Printf("  - %s\n", path)
				}
			}

			return nil
		},
	}
}

func guardContextRepairCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "repair",
		Short: "Rebuild session from git state",
		Long: `Repair a corrupted session file by rebuilding it from the current git state.

Extracts feature ID from the current branch name and creates a new session.`,
		Example: `  sdp guard context repair`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			recovery := context.NewRecovery(root)
			if err := recovery.Repair(); err != nil {
				return fmt.Errorf("repair failed: %w", err)
			}

			fmt.Println("Session repaired successfully")
			return nil
		},
	}
}
