// Package worktree provides git worktree management for SDP.
// It handles creation, deletion, and session management for worktrees.
package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/session"
)

// Creator handles git worktree creation with session initialization.
type Creator struct {
	// MainRepoPath is the path to the main git repository
	MainRepoPath string
	// WorktreesDir is the directory where worktrees are created (default: parent of main repo)
	WorktreesDir string
}

// NewCreator creates a new worktree creator.
func NewCreator(mainRepoPath string) *Creator {
	// Default worktrees to be siblings of the main repo
	worktreesDir := filepath.Dir(mainRepoPath)
	return &Creator{
		MainRepoPath: mainRepoPath,
		WorktreesDir: worktreesDir,
	}
}

// CreateOptions holds options for creating a worktree.
type CreateOptions struct {
	// FeatureID is the feature this worktree is for
	FeatureID string
	// BranchName is the git branch to create/use
	BranchName string
	// BaseBranch is the branch to create from (default: dev)
	BaseBranch string
	// CreateBranch indicates whether to create a new branch
	CreateBranch bool
}

// CreateResult holds the result of worktree creation.
type CreateResult struct {
	// WorktreePath is the path to the created worktree
	WorktreePath string
	// BranchName is the branch used
	BranchName string
	// SessionFile is the path to the created session file
	SessionFile string
}

// Create creates a new git worktree with session initialization (AC6).
func (c *Creator) Create(opts CreateOptions) (*CreateResult, error) {
	if opts.FeatureID == "" {
		return nil, fmt.Errorf("feature ID is required")
	}

	// Set defaults
	branchName := opts.BranchName
	if branchName == "" {
		branchName = fmt.Sprintf("feature/%s", opts.FeatureID)
	}
	baseBranch := opts.BaseBranch
	if baseBranch == "" {
		baseBranch = "dev"
	}

	// Determine worktree path
	worktreeName := fmt.Sprintf("sdp-%s", opts.FeatureID)
	worktreePath := filepath.Join(c.WorktreesDir, worktreeName)

	// Build git worktree command
	var cmd *exec.Cmd
	if opts.CreateBranch {
		// Create new branch from base
		cmd = exec.Command("git", "worktree", "add", "-b", branchName, worktreePath, baseBranch)
	} else {
		// Use existing branch or create without explicit base
		cmd = exec.Command("git", "worktree", "add", worktreePath, branchName)
	}
	cmd.Dir = c.MainRepoPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("create worktree: %w\n%s", err, string(output))
	}

	// Initialize session in the new worktree (AC6)
	s, err := session.Init(opts.FeatureID, worktreePath, "sdp worktree create")
	if err != nil {
		// Rollback: remove the worktree if session init fails
		_ = c.removeWorktree(worktreePath) //nolint:errcheck // rollback best effort
		return nil, fmt.Errorf("init session: %w", err)
	}

	// Save the session
	if err := s.Save(worktreePath); err != nil {
		// Rollback: remove the worktree if session save fails
		_ = c.removeWorktree(worktreePath) //nolint:errcheck // rollback best effort
		return nil, fmt.Errorf("save session: %w", err)
	}

	return &CreateResult{
		WorktreePath: worktreePath,
		BranchName:   branchName,
		SessionFile:  filepath.Join(worktreePath, ".sdp", "session.json"),
	}, nil
}

// removeWorktree removes a worktree (helper for rollback).
func (c *Creator) removeWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", "--force", worktreePath)
	cmd.Dir = c.MainRepoPath
	return cmd.Run()
}

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
