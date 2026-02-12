package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/context"
	"github.com/spf13/cobra"
)

// guardContextCmd returns the guard context command group
func guardContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Context validation and recovery commands",
		Long: `Manage worktree context for git safety.

These commands help verify and recover the correct worktree context
when the CWD may have reset after tool calls.`,
	}
	cmd.AddCommand(guardContextCheckCmd())
	cmd.AddCommand(guardContextShowCmd())
	cmd.AddCommand(guardContextFindCmd())
	cmd.AddCommand(guardContextGoCmd())
	cmd.AddCommand(guardContextCleanCmd())
	cmd.AddCommand(guardContextRepairCmd())
	return cmd
}

// guardContextCheckCmd validates current context
func guardContextCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Validate current worktree context",
		Long: `Check that the current context is valid:
- Worktree path matches session
- Branch matches session
- Session file is valid (hash check)

Exit codes:
  0 - All checks pass
  1 - Context mismatch
  2 - No session file
  3 - Hash mismatch`,
		Example: `  # Check current context
  sdp guard context check

  # Use exit code for scripting
  sdp guard context check && git commit -m "message"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			recovery := context.NewRecovery(root)
			result, err := recovery.Check()
			if err != nil {
				return fmt.Errorf("check failed: %w", err)
			}

			fmt.Print(context.FormatCheckResult(result))

			if !result.Valid {
				os.Exit(result.ExitCode)
			}

			return nil
		},
	}
}

// guardContextShowCmd shows detailed context information
func guardContextShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Display full context details",
		Long:  `Show detailed information about the current worktree context.`,
		Example: `  sdp guard context show`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			recovery := context.NewRecovery(root)
			result, err := recovery.Show()
			if err != nil {
				return fmt.Errorf("show failed: %w", err)
			}

			fmt.Println("Context Details:")
			fmt.Printf("  Worktree Path: %s\n", result.WorktreePath)
			if result.FeatureID != "" {
				fmt.Printf("  Feature ID: %s\n", result.FeatureID)
			}
			fmt.Printf("  Current Branch: %s\n", result.CurrentBranch)
			if result.ExpectedBranch != "" {
				fmt.Printf("  Expected Branch: %s\n", result.ExpectedBranch)
			}
			if result.RemoteTracking != "" {
				fmt.Printf("  Remote Tracking: %s\n", result.RemoteTracking)
			}
			fmt.Printf("  Session Valid: %v\n", result.SessionValid)

			if !result.Valid {
				fmt.Println("\nErrors:")
				for _, err := range result.Errors {
					fmt.Printf("  - %s\n", err)
				}
			}

			return nil
		},
	}
}

// guardContextFindCmd locates worktree for a feature
func guardContextFindCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "find <feature-id>",
		Short: "Locate worktree for a feature",
		Long: `Find the worktree path for a given feature ID.

Uses hybrid recovery strategy:
1. Search session files
2. Parse git worktree list
3. Check workstream metadata`,
		Example: `  # Find worktree for F065
  sdp guard context find F065

  # Use in shell
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

// guardContextGoCmd provides instructions to change to a worktree
func guardContextGoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "go <feature-id>",
		Short: "Print command to change to feature worktree",
		Long: `Print the path and command to change to a feature worktree.

NOTE: This command cannot actually change your shell's CWD.
It outputs the path and instructions for you to execute.`,
		Example: `  # Get instructions to go to F065 worktree
  sdp guard context go F065`,
		Args: cobra.ExactArgs(1),
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

// guardContextCleanCmd cleans up stale sessions
func guardContextCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean up stale session files",
		Long: `Remove invalid or stale session files from all worktrees.`,
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

// guardContextRepairCmd rebuilds session from git state
func guardContextRepairCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "repair",
		Short: "Rebuild session from git state",
		Long: `Repair a corrupted session file by rebuilding it from the current git state.

Extracts feature ID from the current branch name and creates a new session.`,
		Example: `  # Repair session in current directory
  sdp guard context repair`,
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
