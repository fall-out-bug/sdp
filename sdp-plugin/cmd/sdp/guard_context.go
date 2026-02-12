package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/context"
	"github.com/spf13/cobra"
)

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
		Example: `  sdp guard context check
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
				for _, e := range result.Errors {
					fmt.Printf("  - %s\n", e)
				}
			}

			return nil
		},
	}
}
