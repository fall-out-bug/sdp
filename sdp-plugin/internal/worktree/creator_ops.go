package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/session"
)

// Delete removes a worktree and its session file.
func (c *Creator) Delete(featureID string) error {
	worktreeName := fmt.Sprintf("sdp-%s", featureID)
	worktreePath := filepath.Join(c.WorktreesDir, worktreeName)

	// Delete session file first
	if err := session.Delete(worktreePath); err != nil {
		// Log but continue - the worktree removal is more important
		fmt.Fprintf(os.Stderr, "warning: failed to delete session: %v\n", err)
	}

	// Remove the worktree
	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	cmd.Dir = c.MainRepoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("remove worktree: %w\n%s", err, string(output))
	}

	return nil
}

// List returns all worktrees with their sessions.
func (c *Creator) List() ([]WorktreeInfo, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = c.MainRepoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("list worktrees: %w", err)
	}

	return c.parseWorktreeList(string(output))
}

// WorktreeInfo holds information about a worktree.
type WorktreeInfo struct {
	// Path is the worktree path
	Path string
	// Branch is the current branch
	Branch string
	// Session is the session info (may be nil if no session)
	Session *session.Session
}

// parseWorktreeList parses the output of git worktree list --porcelain.
func (c *Creator) parseWorktreeList(output string) ([]WorktreeInfo, error) {
	var worktrees []WorktreeInfo
	var current *WorktreeInfo

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if current != nil && current.Path != "" {
				// Try to load session for this worktree
				if s, err := session.Load(current.Path); err == nil {
					current.Session = s
				}
				worktrees = append(worktrees, *current)
			}
			current = nil
			continue
		}

		if current == nil {
			current = &WorktreeInfo{}
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "worktree":
			current.Path = parts[1]
		case "branch":
			current.Branch = strings.TrimPrefix(parts[1], "refs/heads/")
		}
	}

	// Don't forget the last one
	if current != nil && current.Path != "" {
		if s, err := session.Load(current.Path); err == nil {
			current.Session = s
		}
		worktrees = append(worktrees, *current)
	}

	return worktrees, nil
}
