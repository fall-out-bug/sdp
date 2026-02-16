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

// Wrapper provides git command execution with session validation.
type Wrapper struct {
	ProjectRoot string
	validator   *Validator
}

// NewWrapper creates a new git wrapper for the given project root.
func NewWrapper(projectRoot string) *Wrapper {
	return &Wrapper{
		ProjectRoot: projectRoot,
		validator:   NewValidator(projectRoot),
	}
}

// Execute runs a git command with session validation.
func (w *Wrapper) Execute(command string, args ...string) error {
	if NeedsSessionCheck(command) {
		result, err := w.validator.ValidateSession()
		if err != nil {
			return fmt.Errorf("session validation failed: %w", err)
		}
		if !result.Valid {
			return fmt.Errorf("%s", result.FormatError())
		}
	}

	cmd := exec.Command("git", append([]string{command}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = w.ProjectRoot

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s failed: %w", command, err)
	}

	if NeedsPostCheck(command) {
		result, _ := w.validator.ValidateSession()
		if result != nil && !result.Valid {
			return fmt.Errorf("CRITICAL: branch changed during command!\n%s", result.FormatError())
		}
	}

	return nil
}

// GetCurrentBranch returns the current git branch name.
func GetCurrentBranch(projectRoot string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = projectRoot
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
