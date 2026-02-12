// Package context provides CWD recovery and context validation for git safety.
// It implements hybrid recovery from session files, git worktree list, and workstream metadata.
package context

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/session"
)

// Exit codes for programmatic use (AC8)
const (
	ExitCodeOK             = 0
	ExitCodeContextMismatch = 1
	ExitCodeNoSession      = 2
	ExitCodeHashMismatch   = 3
	ExitCodeRuntimeError   = 4
)

// ContextCheckResult holds the result of a context check.
type ContextCheckResult struct {
	// Valid indicates all checks passed
	Valid bool
	// ExitCode for programmatic use
	ExitCode int
	// WorktreePath is the current worktree path
	WorktreePath string
	// FeatureID is the feature for this worktree
	FeatureID string
	// CurrentBranch is the current git branch
	CurrentBranch string
	// ExpectedBranch is the expected branch from session
	ExpectedBranch string
	// RemoteTracking is the remote tracking branch
	RemoteTracking string
	// SessionValid indicates if the session file is valid
	SessionValid bool
	// Errors contains any validation errors
	Errors []string
}

// Recovery handles CWD recovery from multiple sources.
type Recovery struct {
	// ProjectRoot is the main repository root
	ProjectRoot string
}

// NewRecovery creates a new recovery handler.
func NewRecovery(projectRoot string) *Recovery {
	return &Recovery{ProjectRoot: projectRoot}
}

// Check validates the current context (AC1).
func (r *Recovery) Check() (*ContextCheckResult, error) {
	result := &ContextCheckResult{
		Valid:    true,
		ExitCode: ExitCodeOK,
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}
	result.WorktreePath = cwd

	// Try to load session
	s, err := session.Load(cwd)
	if err != nil {
		// No session file
		result.Valid = false
		result.ExitCode = ExitCodeNoSession
		result.Errors = append(result.Errors, "No session file found")
		return result, nil
	}

	// Verify session hash
	if !s.IsValid() {
		result.Valid = false
		result.ExitCode = ExitCodeHashMismatch
		result.SessionValid = false
		result.Errors = append(result.Errors, "Session file corrupted or tampered")
		return result, nil
	}

	result.SessionValid = true
	result.FeatureID = s.FeatureID
	result.ExpectedBranch = s.ExpectedBranch
	result.RemoteTracking = s.ExpectedRemote

	// Get current branch
	currentBranch, err := r.getCurrentBranch()
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to get current branch: %v", err))
		return result, nil
	}
	result.CurrentBranch = currentBranch

	// Verify worktree path
	if cwd != s.WorktreePath {
		result.Valid = false
		result.ExitCode = ExitCodeContextMismatch
		result.Errors = append(result.Errors,
			fmt.Sprintf("Worktree mismatch: expected %s, in %s", s.WorktreePath, cwd))
	}

	// Verify branch
	if currentBranch != s.ExpectedBranch {
		result.Valid = false
		result.ExitCode = ExitCodeContextMismatch
		result.Errors = append(result.Errors,
			fmt.Sprintf("Branch mismatch: expected %s, on %s", s.ExpectedBranch, currentBranch))
	}

	return result, nil
}

// Show returns detailed context information (AC2).
func (r *Recovery) Show() (*ContextCheckResult, error) {
	return r.Check()
}

// FindWorktree locates the worktree for a given feature (AC3, AC7).
// Uses hybrid recovery: session files → git worktree list → workstream metadata
func (r *Recovery) FindWorktree(featureID string) (string, error) {
	// Strategy 1: Look for existing session files with this feature ID
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

	// Strategy 2: Check for worktree with feature ID in path name
	expectedName := fmt.Sprintf("sdp-%s", featureID)
	for _, wt := range worktrees {
		if strings.Contains(wt, expectedName) {
			return wt, nil
		}
	}

	// Strategy 3: Check workstream files for feature context
	wsPath := filepath.Join(r.ProjectRoot, "docs", "workstreams", "backlog")
	if _, err := os.Stat(wsPath); err == nil {
		// Feature workstreams exist, suggest creating worktree
		return "", fmt.Errorf("no worktree found for feature %s - create one with: sdp worktree create %s", featureID, featureID)
	}

	return "", fmt.Errorf("feature %s not found", featureID)
}

// GoToWorktree returns instructions to change to a feature worktree (AC4).
// Note: This cannot actually change the shell's CWD, but returns the path
// and instructions for the caller.
func (r *Recovery) GoToWorktree(featureID string) (string, error) {
	path, err := r.FindWorktree(featureID)
	if err != nil {
		return "", err
	}
	return path, nil
}

// Clean removes stale session files (AC5).
func (r *Recovery) Clean() ([]string, error) {
	var cleaned []string

	worktrees, err := r.listWorktrees()
	if err != nil {
		return nil, fmt.Errorf("list worktrees: %w", err)
	}

	// Create a map of valid worktree paths
	validPaths := make(map[string]bool)
	for _, wt := range worktrees {
		validPaths[wt] = true
	}

	// Check all worktrees for orphaned session files
	for _, wt := range worktrees {
		sessionPath := filepath.Join(wt, ".sdp", session.SessionFileName)
		if _, err := os.Stat(sessionPath); err == nil {
			// Session file exists - verify it's valid
			s, err := session.Load(wt)
			if err != nil || !s.IsValid() {
				// Invalid session - remove it
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
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Get current branch
	branch, err := r.getCurrentBranch()
	if err != nil {
		return fmt.Errorf("get current branch: %w", err)
	}

	// Extract feature ID from branch
	featureID := extractFeatureID(branch)
	if featureID == "" {
		return fmt.Errorf("could not extract feature ID from branch %s", branch)
	}

	// Get remote tracking
	remote, err := r.getRemoteTracking(branch)
	if err != nil {
		remote = fmt.Sprintf("origin/%s", branch)
	}

	// Repair the session
	_, err = session.Repair(cwd, featureID, branch, remote)
	if err != nil {
		return fmt.Errorf("repair session: %w", err)
	}

	return nil
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

// getCurrentBranch returns the current git branch.
func (r *Recovery) getCurrentBranch() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getRemoteTracking returns the remote tracking branch.
func (r *Recovery) getRemoteTracking(branch string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{u}")
	cmd.Dir = cwd
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// extractFeatureID extracts feature ID from branch name.
func extractFeatureID(branch string) string {
	// Try feature/F### pattern
	if strings.HasPrefix(branch, "feature/") {
		return strings.TrimPrefix(branch, "feature/")
	}
	// Try bugfix/ pattern
	if strings.HasPrefix(branch, "bugfix/") {
		return strings.TrimPrefix(branch, "bugfix/")
	}
	// Try hotfix/ pattern
	if strings.HasPrefix(branch, "hotfix/") {
		return strings.TrimPrefix(branch, "hotfix/")
	}
	return ""
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

// FindProjectRoot finds the project root from the current directory.
func FindProjectRoot() (string, error) {
	return config.FindProjectRoot()
}
