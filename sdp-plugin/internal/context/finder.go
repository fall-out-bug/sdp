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

// FindWorktree locates the worktree for a given feature (AC3, AC7).
func (r *Recovery) FindWorktree(featureID string) (string, error) {
	worktrees, err := r.listWorktrees()
	if err == nil {
		for _, wt := range worktrees {
			if session.Exists(wt) {
				s, err := session.Load(wt)
				if err == nil && s.FeatureID == featureID {
					return wt, nil
				}
			}
		}
	}

	expectedName := fmt.Sprintf("sdp-%s", featureID)
	for _, wt := range worktrees {
		if strings.Contains(wt, expectedName) {
			return wt, nil
		}
	}

	wsPath := filepath.Join(r.ProjectRoot, "docs", "workstreams", "backlog")
	if _, err := os.Stat(wsPath); err == nil {
		return "", fmt.Errorf("no worktree for %s - create: sdp worktree create %s", featureID, featureID)
	}

	return "", fmt.Errorf("feature %s not found", featureID)
}

// GoToWorktree returns instructions to change to a feature worktree (AC4).
func (r *Recovery) GoToWorktree(featureID string) (string, error) {
	return r.FindWorktree(featureID)
}

// listWorktrees returns all worktree paths.
func (r *Recovery) listWorktrees() ([]string, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = r.ProjectRoot
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var worktrees []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "worktree ") {
			worktrees = append(worktrees, strings.TrimPrefix(line, "worktree "))
		}
	}

	return worktrees, nil
}

// extractFeatureID extracts feature ID from branch name.
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
