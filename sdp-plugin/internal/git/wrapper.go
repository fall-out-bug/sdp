// Package git provides git operation wrappers with session validation.
// It implements safe git command execution with pre/post checks.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/session"
)

// CommandCategory represents the type of git command.
type CommandCategory int

const (
	// CategorySafe commands are read-only (status, log, diff, show)
	CategorySafe CommandCategory = iota
	// CategoryWrite commands modify the repository (add, commit, reset)
	CategoryWrite
	// CategoryRemote commands interact with remotes (push, fetch, pull)
	CategoryRemote
	// CategoryBranch commands change branches (checkout, branch, merge)
	CategoryBranch
)

// String returns the string representation of a command category.
func (c CommandCategory) String() string {
	switch c {
	case CategorySafe:
		return "safe"
	case CategoryWrite:
		return "write"
	case CategoryRemote:
		return "remote"
	case CategoryBranch:
		return "branch"
	default:
		return "unknown"
	}
}

// ValidationResult holds the result of session validation.
type ValidationResult struct {
	// Valid indicates whether all checks passed
	Valid bool
	// WorktreePath is the expected worktree path from session
	WorktreePath string
	// ActualPath is the current working directory
	ActualPath string
	// CurrentBranch is the current git branch
	CurrentBranch string
	// ExpectedBranch is the expected branch from session
	ExpectedBranch string
	// ExpectedRemote is the expected remote tracking branch
	ExpectedRemote string
	// ActualRemote is the actual remote tracking
	ActualRemote string
	// Error describes what validation failed
	Error string
	// Fix suggests how to resolve the error
	Fix string
}

// FormatError returns a formatted error message with remediation guidance.
func (r ValidationResult) FormatError() string {
	if r.Valid {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ERROR: %s\n", r.Error))

	if r.ActualPath != "" && r.WorktreePath != "" && r.ActualPath != r.WorktreePath {
		sb.WriteString(fmt.Sprintf("  Expected worktree: %s\n", r.WorktreePath))
		sb.WriteString(fmt.Sprintf("  Current directory: %s\n", r.ActualPath))
	}

	if r.CurrentBranch != "" && r.ExpectedBranch != "" && r.CurrentBranch != r.ExpectedBranch {
		sb.WriteString(fmt.Sprintf("  Expected branch: %s\n", r.ExpectedBranch))
		sb.WriteString(fmt.Sprintf("  Current branch: %s\n", r.CurrentBranch))
	}

	if r.Fix != "" {
		sb.WriteString(fmt.Sprintf("\nFIX: %s\n", r.Fix))
	}

	return sb.String()
}

// Wrapper provides git command execution with session validation.
type Wrapper struct {
	// ProjectRoot is the root directory of the project
	ProjectRoot string
}

// NewWrapper creates a new git wrapper for the given project root.
func NewWrapper(projectRoot string) *Wrapper {
	return &Wrapper{ProjectRoot: projectRoot}
}

// safeCommands lists commands that are read-only.
var safeCommands = map[string]bool{
	"status": true,
	"log":    true,
	"diff":   true,
	"show":   true,
	"ls-files": true,
	"rev-parse": true,
	"branch": true, // 'git branch' without args is safe
	"remote": true,
	"tag":    true, // 'git tag' without args is safe
}

// writeCommands lists commands that modify the repository.
var writeCommands = map[string]bool{
	"add":    true,
	"commit": true,
	"reset":  true,
	"rm":     true,
	"mv":     true,
	"stash":  true,
}

// remoteCommands lists commands that interact with remotes.
var remoteCommands = map[string]bool{
	"push":  true,
	"fetch": true,
	"pull":  true,
	"clone": true,
}

// branchChangeCommands lists commands that change branches.
var branchChangeCommands = map[string]bool{
	"checkout": true,
	"switch":   true,
	"merge":    true,
	"rebase":   true,
	"cherry-pick": true,
}

// CategorizeCommand determines the category of a git command.
func CategorizeCommand(cmd string) CommandCategory {
	cmd = strings.ToLower(cmd)

	if safeCommands[cmd] {
		return CategorySafe
	}
	if writeCommands[cmd] {
		return CategoryWrite
	}
	if remoteCommands[cmd] {
		return CategoryRemote
	}
	if branchChangeCommands[cmd] {
		return CategoryBranch
	}

	// Default to safe for unknown commands
	return CategorySafe
}

// NeedsSessionCheck returns true if the command requires session validation.
func NeedsSessionCheck(cmd string) bool {
	// All commands need session check to ensure we're in the right worktree
	return true
}

// NeedsPostCheck returns true if the command needs post-execution validation.
func NeedsPostCheck(cmd string) bool {
	cmd = strings.ToLower(cmd)

	// Write and branch commands need post-check to ensure branch didn't change
	return writeCommands[cmd] || branchChangeCommands[cmd] || cmd == "push" || cmd == "pull"
}

// ValidateSession checks if the current context matches the session.
func (w *Wrapper) ValidateSession() (*ValidationResult, error) {
	// Load session
	s, err := session.Load(w.ProjectRoot)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Error: "session file not found or invalid",
			Fix:   "sdp session init --feature=F###",
		}, err
	}

	// Verify session integrity
	if !s.IsValid() {
		return &ValidationResult{
			Valid: false,
			Error: "session file corrupted or tampered",
			Fix:   "sdp session repair --force",
		}, fmt.Errorf("session hash validation failed")
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	// Validate worktree path
	if cwd != s.WorktreePath {
		return &ValidationResult{
			Valid:         false,
			WorktreePath:  s.WorktreePath,
			ActualPath:    cwd,
			ExpectedBranch: s.ExpectedBranch,
			Error:         "wrong worktree",
			Fix:           fmt.Sprintf("sdp guard context go %s", s.FeatureID),
		}, nil
	}

	// Get current branch
	currentBranch, err := w.getCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("get current branch: %w", err)
	}

	// Validate branch
	if currentBranch != s.ExpectedBranch {
		return &ValidationResult{
			Valid:          false,
			WorktreePath:   s.WorktreePath,
			ActualPath:     cwd,
			CurrentBranch:  currentBranch,
			ExpectedBranch: s.ExpectedBranch,
			Error:          "wrong branch",
			Fix:            fmt.Sprintf("git checkout %s", s.ExpectedBranch),
		}, nil
	}

	return &ValidationResult{
		Valid:          true,
		WorktreePath:   s.WorktreePath,
		ActualPath:     cwd,
		CurrentBranch:  currentBranch,
		ExpectedBranch: s.ExpectedBranch,
		ExpectedRemote: s.ExpectedRemote,
	}, nil
}

// Execute runs a git command with session validation.
func (w *Wrapper) Execute(command string, args ...string) error {
	// Check if session validation is needed
	if NeedsSessionCheck(command) {
		result, err := w.ValidateSession()
		if err != nil {
			return fmt.Errorf("session validation failed: %w", err)
		}
		if !result.Valid {
			return fmt.Errorf("%s", result.FormatError())
		}
	}

	// Build git command
	cmd := exec.Command("git", append([]string{command}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = w.ProjectRoot

	// Execute command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s failed: %w", command, err)
	}

	// Post-check for write commands
	if NeedsPostCheck(command) {
		result, _ := w.ValidateSession()
		if result != nil && !result.Valid {
			return fmt.Errorf("CRITICAL: branch changed during command!\n%s", result.FormatError())
		}
	}

	return nil
}

// getCurrentBranch returns the current git branch name.
func (w *Wrapper) getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = w.ProjectRoot
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetWorktreePath returns the worktree path from session.
func (w *Wrapper) GetWorktreePath() (string, error) {
	s, err := session.Load(w.ProjectRoot)
	if err != nil {
		return "", err
	}
	return s.WorktreePath, nil
}

// HasSession returns true if a valid session exists.
func (w *Wrapper) HasSession() bool {
	return session.Exists(w.ProjectRoot)
}

// FindProjectRoot finds the project root from the current directory.
func FindProjectRoot() (string, error) {
	return config.FindProjectRoot()
}
