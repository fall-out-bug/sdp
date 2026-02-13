// Package worktree provides git worktree management for SDP.
// It handles creation, deletion, and session management for worktrees.
package worktree

import (
	"fmt"
	"os/exec"
	"path/filepath"

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
