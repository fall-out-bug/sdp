// Package context provides CWD recovery and context validation for git safety.
package context

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/session"
)

// Clean removes stale session files (AC5).
func (r *Recovery) Clean() ([]string, error) {
	var cleaned []string

	worktrees, err := r.listWorktrees()
	if err != nil {
		return nil, fmt.Errorf("list worktrees: %w", err)
	}

	for _, wt := range worktrees {
		sessionPath := filepath.Join(wt, ".sdp", session.SessionFileName)
		if _, err := os.Stat(sessionPath); err == nil {
			s, err := session.Load(wt)
			if err != nil || !s.IsValid() {
				if err := session.Delete(wt); err == nil {
					cleaned = append(cleaned, sessionPath)
				}
			}
		}
	}

	return cleaned, nil
}

// Repair rebuilds the session from git state (AC6).
func (r *Recovery) Repair() error {
	branch, err := getCurrentBranchIn(r.ProjectRoot)
	if err != nil {
		return fmt.Errorf("get current branch: %w", err)
	}

	featureID := extractFeatureID(branch)
	if featureID == "" {
		return fmt.Errorf("could not extract feature ID from branch %s", branch)
	}

	remote, err := getRemoteTrackingIn(r.ProjectRoot)
	if err != nil {
		remote = fmt.Sprintf("origin/%s", branch)
	}

	_, err = session.Repair(r.ProjectRoot, featureID, branch, remote)
	if err != nil {
		return fmt.Errorf("repair session: %w", err)
	}

	return nil
}

// getCurrentBranch returns the current git branch in the specified directory.
func getCurrentBranchIn(dir string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getCurrentBranch returns the current git branch (deprecated: use getCurrentBranchIn).
func getCurrentBranch() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return getCurrentBranchIn(cwd)
}

// getRemoteTrackingIn returns the remote tracking branch in the specified directory.
func getRemoteTrackingIn(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{u}")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getRemoteTracking returns the remote tracking branch (deprecated: use getRemoteTrackingIn).
func getRemoteTracking() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return getRemoteTrackingIn(cwd)
}

// FormatCheckResult formats the check result for display.
func FormatCheckResult(result *ContextCheckResult) string {
	var sb strings.Builder

	if result.Valid {
		sb.WriteString(fmt.Sprintf("  [OK] Worktree: %s\n", result.WorktreePath))
		sb.WriteString(fmt.Sprintf("  [OK] Branch: %s\n", result.CurrentBranch))
		sb.WriteString(fmt.Sprintf("  [OK] Tracking: %s\n", result.RemoteTracking))
		sb.WriteString("  [OK] Session: valid\n")
	} else {
		for _, err := range result.Errors {
			sb.WriteString(fmt.Sprintf("  [FAIL] %s\n", err))
		}
		if result.WorktreePath != "" {
			sb.WriteString(fmt.Sprintf("  Worktree: %s\n", result.WorktreePath))
		}
		if result.CurrentBranch != "" {
			sb.WriteString(fmt.Sprintf("  Current Branch: %s\n", result.CurrentBranch))
		}
		if result.ExpectedBranch != "" {
			sb.WriteString(fmt.Sprintf("  Expected Branch: %s\n", result.ExpectedBranch))
		}
		sb.WriteString(fmt.Sprintf("  Session Valid: %v\n", result.SessionValid))
	}

	return sb.String()
}
