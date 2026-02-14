package git

import (
	"strings"
	"testing"
)

func TestValidationResult_FormatError_Valid(t *testing.T) {
	r := &ValidationResult{Valid: true}
	if r.FormatError() != "" {
		t.Error("FormatError should return empty string for valid result")
	}
}

func TestValidationResult_FormatError_Invalid(t *testing.T) {
	r := &ValidationResult{
		Valid:          false,
		Error:          "wrong branch",
		WorktreePath:   "/path/to/worktree",
		ActualPath:     "/path/to/actual",
		ExpectedBranch: "feature/F001",
		CurrentBranch:  "main",
		Fix:            "git checkout feature/F001",
	}

	output := r.FormatError()
	if output == "" {
		t.Error("FormatError should return non-empty string for invalid result")
	}

	// Check that key elements are present
	if !strings.Contains(output, "ERROR:") {
		t.Error("FormatError output should contain ERROR:")
	}
	if !strings.Contains(output, "wrong branch") {
		t.Error("FormatError output should contain error message")
	}
	if !strings.Contains(output, "Expected worktree:") {
		t.Error("FormatError output should contain Expected worktree")
	}
	if !strings.Contains(output, "Expected branch:") {
		t.Error("FormatError output should contain Expected branch")
	}
	if !strings.Contains(output, "FIX:") {
		t.Error("FormatError output should contain FIX:")
	}
}

func TestValidationResult_FormatError_OnlyError(t *testing.T) {
	r := &ValidationResult{
		Valid: false,
		Error: "session not found",
		Fix:   "run sdp init",
	}

	output := r.FormatError()
	if output == "" {
		t.Error("FormatError should return non-empty string")
	}

	if !strings.Contains(output, "ERROR:") {
		t.Error("FormatError output should contain ERROR:")
	}
	if !strings.Contains(output, "session not found") {
		t.Error("FormatError output should contain error message")
	}
	if !strings.Contains(output, "FIX:") {
		t.Error("FormatError output should contain FIX:")
	}
}

func TestValidationResult_FormatError_NoFix(t *testing.T) {
	r := &ValidationResult{
		Valid: false,
		Error: "unknown error",
	}

	output := r.FormatError()
	if output == "" {
		t.Error("FormatError should return non-empty string")
	}

	if !strings.Contains(output, "ERROR:") {
		t.Error("FormatError output should contain ERROR:")
	}
}

func TestValidationResult_FormatError_MismatchedPaths(t *testing.T) {
	r := &ValidationResult{
		Valid:        false,
		Error:        "wrong worktree",
		WorktreePath: "/expected/path",
		ActualPath:   "/actual/path",
		Fix:          "cd /expected/path",
	}

	output := r.FormatError()

	if !strings.Contains(output, "Expected worktree:") {
		t.Error("FormatError should show expected worktree when paths differ")
	}
	if !strings.Contains(output, "Current directory:") {
		t.Error("FormatError should show current directory when paths differ")
	}
}

func TestValidationResult_FormatError_MismatchedBranches(t *testing.T) {
	r := &ValidationResult{
		Valid:          false,
		Error:          "wrong branch",
		CurrentBranch:  "main",
		ExpectedBranch: "feature/F001",
		Fix:            "git checkout feature/F001",
	}

	output := r.FormatError()

	if !strings.Contains(output, "Expected branch:") {
		t.Error("FormatError should show expected branch when branches differ")
	}
	if !strings.Contains(output, "Current branch:") {
		t.Error("FormatError should show current branch when branches differ")
	}
}
