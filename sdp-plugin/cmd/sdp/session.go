package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/session"
	"github.com/spf13/cobra"
)

func sessionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage per-worktree session state",
		Long: `Manage session state for git safety.

Sessions track the expected identity of a worktree to prevent
branch confusion and cross-feature commits.`,
	}
	cmd.AddCommand(sessionInitCmd())
	cmd.AddCommand(sessionSyncCmd())
	cmd.AddCommand(sessionRepairCmd())
	cmd.AddCommand(sessionShowCmd())
	cmd.AddCommand(sessionDeleteCmd())
	return cmd
}

func sessionInitCmd() *cobra.Command {
	var featureID string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize session in current worktree",
		Long: `Create a new session file for the current worktree.

This establishes the expected identity for git operations,
preventing accidental commits to wrong branches.`,
		Example: `  # Initialize session for feature F065
  sdp session init --feature=F065`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if featureID == "" {
				return fmt.Errorf("--feature flag is required")
			}

			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get working directory: %w", err)
			}

			// Check if session already exists
			if session.Exists(root) {
				existing, err := session.Load(root)
				if err == nil && existing.IsValid() {
					fmt.Printf("Session already exists for feature %s\n", existing.FeatureID)
					fmt.Printf("Run 'sdp session repair' to rebuild, or 'sdp session delete' first\n")
					return nil
				}
			}

			s, err := session.Init(featureID, cwd, "sdp session init")
			if err != nil {
				return fmt.Errorf("init session: %w", err)
			}

			if err := s.Save(root); err != nil {
				return fmt.Errorf("save session: %w", err)
			}

			fmt.Printf("Session initialized for feature %s\n", featureID)
			fmt.Printf("  Worktree: %s\n", cwd)
			fmt.Printf("  Branch: %s\n", s.ExpectedBranch)
			fmt.Printf("  Remote: %s\n", s.ExpectedRemote)
			return nil
		},
	}

	cmd.Flags().StringVar(&featureID, "feature", "", "Feature ID for this worktree")
	cmd.MarkFlagRequired("feature")

	return cmd
}

func sessionSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync session with current git state",
		Long: `Update the session file to match the current git branch and remote.

Use this after manually switching branches to keep the session in sync.`,
		Example: `  # Sync session with current branch
  sdp session sync`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			s, err := session.Load(root)
			if err != nil {
				return fmt.Errorf("load session: %w", err)
			}

			// Get current branch from git
			branch, err := getCurrentBranch()
			if err != nil {
				return fmt.Errorf("get current branch: %w", err)
			}

			// Get current remote tracking
			remote, err := getCurrentRemote(branch)
			if err != nil {
				remote = fmt.Sprintf("origin/%s", branch)
			}

			s.Sync(branch, remote)

			if err := s.Save(root); err != nil {
				return fmt.Errorf("save session: %w", err)
			}

			fmt.Printf("Session synced\n")
			fmt.Printf("  Branch: %s\n", branch)
			fmt.Printf("  Remote: %s\n", remote)
			return nil
		},
	}
}

func sessionRepairCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair corrupted session",
		Long: `Rebuild the session file from scratch.

Use this when the session file is corrupted or has been tampered with.`,
		Example: `  # Repair session
  sdp session repair --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			// Get current branch from git
			branch, err := getCurrentBranch()
			if err != nil {
				return fmt.Errorf("get current branch: %w", err)
			}

			// Extract feature ID from branch name
			featureID := extractFeatureID(branch)
			if featureID == "" {
				return fmt.Errorf("could not extract feature ID from branch %s", branch)
			}

			// Get current remote tracking
			remote, err := getCurrentRemote(branch)
			if err != nil {
				remote = fmt.Sprintf("origin/%s", branch)
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get working directory: %w", err)
			}

			s, err := session.Repair(root, featureID, branch, remote)
			if err != nil {
				return fmt.Errorf("repair session: %w", err)
			}

			// Update worktree path to current and save (Save recalculates hash)
			s.WorktreePath = cwd
			if err := s.Save(root); err != nil {
				return fmt.Errorf("save session: %w", err)
			}

			fmt.Printf("Session repaired for feature %s\n", featureID)
			fmt.Printf("  Worktree: %s\n", cwd)
			fmt.Printf("  Branch: %s\n", branch)
			fmt.Printf("  Remote: %s\n", remote)
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force repair even if session is valid")

	return cmd
}

func sessionShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current session details",
		Long:  `Display the current session file contents.`,
		Example: `  # Show session
  sdp session show`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			s, err := session.Load(root)
			if err != nil {
				return fmt.Errorf("load session: %w", err)
			}

			fmt.Printf("Session Details:\n")
			fmt.Printf("  Version: %s\n", s.Version)
			fmt.Printf("  Worktree: %s\n", s.WorktreePath)
			fmt.Printf("  Feature: %s\n", s.FeatureID)
			fmt.Printf("  Expected Branch: %s\n", s.ExpectedBranch)
			fmt.Printf("  Expected Remote: %s\n", s.ExpectedRemote)
			fmt.Printf("  Created At: %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("  Created By: %s\n", s.CreatedBy)
			fmt.Printf("  Hash: %s\n", s.Hash[:16]+"...")
			fmt.Printf("  Valid: %v\n", s.IsValid())

			return nil
		},
	}
}

func sessionDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete session file",
		Long:  `Remove the session file from this worktree.`,
		Example: `  # Delete session
  sdp session delete`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			if err := session.Delete(root); err != nil {
				return fmt.Errorf("delete session: %w", err)
			}

			fmt.Println("Session deleted")
			return nil
		},
	}
}

// Helper functions

// getCurrentBranch returns the current git branch name.
func getCurrentBranch() (string, error) {
	// This is a simplified implementation
	// In production, this would use git commands
	// For now, we'll use a placeholder that can be enhanced
	return "", fmt.Errorf("not implemented - use git branch --show-current")
}

// getCurrentRemote returns the remote tracking branch for the given branch.
func getCurrentRemote(branch string) (string, error) {
	// This is a simplified implementation
	// In production, this would use git rev-parse --abbrev-ref @{u}
	return "", fmt.Errorf("not implemented - use git rev-parse --abbrev-ref @{u}")
}

// extractFeatureID extracts the feature ID from a branch name.
func extractFeatureID(branch string) string {
	// Try feature/F### pattern
	if len(branch) > 8 && branch[:8] == "feature/" {
		return branch[8:]
	}
	return ""
}
