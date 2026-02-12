// Package git provides git operation wrappers with session validation.
package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/fall-out-bug/sdp/internal/session"
)

// ValidationResult holds the result of session validation.
type ValidationResult struct {
	Valid          bool
	WorktreePath   string
	ActualPath     string
	CurrentBranch  string
	ExpectedBranch string
	ExpectedRemote string
	ActualRemote   string
	Error          string
	Fix            string
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

// Validator provides session validation for git operations.
type Validator struct {
	ProjectRoot string
}

// NewValidator creates a new validator.
func NewValidator(projectRoot string) *Validator {
	return &Validator{ProjectRoot: projectRoot}
}

// ValidateSession checks if the current context matches the session.
func (v *Validator) ValidateSession() (*ValidationResult, error) {
	s, err := session.Load(v.ProjectRoot)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Error: "session file not found or invalid",
			Fix:   "sdp session init --feature=F###",
		}, err
	}

	if !s.IsValid() {
		return &ValidationResult{
			Valid: false,
			Error: "session file corrupted or tampered",
			Fix:   "sdp session repair --force",
		}, fmt.Errorf("session hash validation failed")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	if cwd != s.WorktreePath {
		return &ValidationResult{
			Valid:          false,
			WorktreePath:   s.WorktreePath,
			ActualPath:     cwd,
			ExpectedBranch: s.ExpectedBranch,
			Error:          "wrong worktree",
			Fix:            fmt.Sprintf("sdp guard context go %s", s.FeatureID),
		}, nil
	}

	currentBranch, err := getCurrentBranch(v.ProjectRoot)
	if err != nil {
		return nil, fmt.Errorf("get current branch: %w", err)
	}

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

// getCurrentBranch returns the current git branch name.
func getCurrentBranch(projectRoot string) (string, error) {
	return GetCurrentBranch(projectRoot)
}
