// Package context provides CWD recovery and context validation for git safety.
package context

import (
	"fmt"
	"os"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/safetylog"
	"github.com/fall-out-bug/sdp/internal/session"
)

// Exit codes for programmatic use (AC8)
const (
	ExitCodeOK              = 0
	ExitCodeContextMismatch = 1
	ExitCodeNoSession       = 2
	ExitCodeHashMismatch    = 3
	ExitCodeRuntimeError    = 4
)

// ContextCheckResult holds the result of a context check.
type ContextCheckResult struct {
	Valid          bool
	ExitCode       int
	WorktreePath   string
	FeatureID      string
	CurrentBranch  string
	ExpectedBranch string
	RemoteTracking string
	SessionValid   bool
	Errors         []string
}

// Recovery handles CWD recovery from multiple sources.
type Recovery struct {
	ProjectRoot string
}

// NewRecovery creates a new recovery handler.
func NewRecovery(projectRoot string) *Recovery {
	return &Recovery{ProjectRoot: projectRoot}
}

// Check validates the current context (AC1).
func (r *Recovery) Check() (*ContextCheckResult, error) {
	start := time.Now()
	result := &ContextCheckResult{
		Valid:    true,
		ExitCode: ExitCodeOK,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}
	result.WorktreePath = cwd

	s, err := session.Load(cwd)
	if err != nil {
		result.Valid = false
		result.ExitCode = ExitCodeNoSession
		result.Errors = append(result.Errors, "No session file found")
		return result, nil
	}

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

	currentBranch, err := getCurrentBranchIn(cwd)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to get branch: %v", err))
		return result, nil
	}
	result.CurrentBranch = currentBranch

	if cwd != s.WorktreePath {
		result.Valid = false
		result.ExitCode = ExitCodeContextMismatch
		result.Errors = append(result.Errors,
			fmt.Sprintf("Worktree mismatch: expected %s, in %s", s.WorktreePath, cwd))
	}

	if currentBranch != s.ExpectedBranch {
		result.Valid = false
		result.ExitCode = ExitCodeContextMismatch
		result.Errors = append(result.Errors,
			fmt.Sprintf("Branch mismatch: expected %s, on %s", s.ExpectedBranch, currentBranch))
	}

	status := "valid"
	if !result.Valid {
		status = "invalid"
	}
	safetylog.Context("check", cwd)
	safetylog.Debug("context check: %s (%v)", status, time.Since(start))

	return result, nil
}

// Show returns detailed context information (AC2).
func (r *Recovery) Show() (*ContextCheckResult, error) {
	return r.Check()
}

// FindProjectRoot finds the project root from the current directory.
func FindProjectRoot() (string, error) {
	return config.FindProjectRoot()
}
