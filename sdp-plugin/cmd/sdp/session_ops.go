package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/session"
	"github.com/spf13/cobra"
)

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

			branch, err := runGitCmd("branch", "--show-current")
			if err != nil {
				return fmt.Errorf("get current branch: %w", err)
			}

			remote, err := runGitCmd("rev-parse", "--abbrev-ref", "@{u}")
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

			branch, err := runGitCmd("branch", "--show-current")
			if err != nil {
				return fmt.Errorf("get current branch: %w", err)
			}

			featureID := extractFeatureID(branch)
			if featureID == "" {
				return fmt.Errorf("could not extract feature ID from branch %s", branch)
			}

			remote, err := runGitCmd("rev-parse", "--abbrev-ref", "@{u}")
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

// runGitCmd executes a git command and returns the output.
func runGitCmd(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// extractFeatureID extracts the feature ID from a branch name.
func extractFeatureID(branch string) string {
	if strings.HasPrefix(branch, "feature/") {
		return strings.TrimPrefix(branch, "feature/")
	}
	if strings.HasPrefix(branch, "bugfix/") {
		return strings.TrimPrefix(branch, "bugfix/")
	}
	if strings.HasPrefix(branch, "hotfix/") {
		return strings.TrimPrefix(branch, "hotfix/")
	}
	return ""
}
