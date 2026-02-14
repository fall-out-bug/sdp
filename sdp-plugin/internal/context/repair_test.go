package context

import (
	"testing"
)

func TestFormatCheckResult_ValidResult(t *testing.T) {
	result := &ContextCheckResult{
		Valid:          true,
		ExitCode:       ExitCodeOK,
		WorktreePath:   "/path/to/worktree",
		CurrentBranch:  "feature/F067",
		RemoteTracking: "origin/feature/F067",
		SessionValid:   true,
	}

	output := FormatCheckResult(result)

	if output == "" {
		t.Error("FormatCheckResult should not return empty string")
	}

	// Check for expected markers
	expectedMarkers := []string{"[OK]", "Worktree:", "Branch:", "Tracking:", "Session:"}
	for _, marker := range expectedMarkers {
		if !containsSubstring(output, marker) {
			t.Errorf("Output should contain %q, got: %s", marker, output)
		}
	}
}

func TestFormatCheckResult_InvalidResult(t *testing.T) {
	result := &ContextCheckResult{
		Valid:          false,
		ExitCode:       ExitCodeContextMismatch,
		WorktreePath:   "/path/to/worktree",
		CurrentBranch:  "dev",
		ExpectedBranch: "feature/F067",
		SessionValid:   true,
		Errors:         []string{"Branch mismatch: expected feature/F067, on dev"},
	}

	output := FormatCheckResult(result)

	if output == "" {
		t.Error("FormatCheckResult should not return empty string")
	}

	if !containsSubstring(output, "[FAIL]") {
		t.Error("Invalid result should contain [FAIL]")
	}

	if !containsSubstring(output, "Branch mismatch") {
		t.Error("Output should contain error message")
	}
}

func TestFormatCheckResult_MultipleErrors(t *testing.T) {
	result := &ContextCheckResult{
		Valid:        false,
		ExitCode:     ExitCodeContextMismatch,
		WorktreePath: "/path/to/worktree",
		Errors: []string{
			"Error 1",
			"Error 2",
			"Error 3",
		},
	}

	output := FormatCheckResult(result)

	// All errors should be present
	for _, err := range result.Errors {
		if !containsSubstring(output, err) {
			t.Errorf("Output should contain error %q", err)
		}
	}
}

func TestFormatCheckResult_EmptyResult(t *testing.T) {
	result := &ContextCheckResult{
		Valid:    false,
		ExitCode: ExitCodeNoSession,
	}

	output := FormatCheckResult(result)

	if output == "" {
		t.Error("FormatCheckResult should return something even for empty result")
	}
}

func TestFormatCheckResult_SessionInvalid(t *testing.T) {
	result := &ContextCheckResult{
		Valid:        false,
		ExitCode:     ExitCodeHashMismatch,
		SessionValid: false,
		Errors:       []string{"Session file corrupted"},
	}

	output := FormatCheckResult(result)

	if !containsSubstring(output, "Session Valid: false") {
		t.Error("Output should show Session Valid: false")
	}
}

func TestContextCheckResult_AllFields(t *testing.T) {
	result := &ContextCheckResult{
		Valid:          true,
		ExitCode:       ExitCodeOK,
		WorktreePath:   "/path/to/worktree",
		FeatureID:      "F067",
		CurrentBranch:  "feature/F067",
		ExpectedBranch: "feature/F067",
		RemoteTracking: "origin/feature/F067",
		SessionValid:   true,
		Errors:         nil,
	}

	// Verify all fields are accessible
	if result.WorktreePath != "/path/to/worktree" {
		t.Errorf("WorktreePath = %q, want /path/to/worktree", result.WorktreePath)
	}
	if result.FeatureID != "F067" {
		t.Errorf("FeatureID = %q, want F067", result.FeatureID)
	}
	if result.CurrentBranch != "feature/F067" {
		t.Errorf("CurrentBranch = %q, want feature/F067", result.CurrentBranch)
	}
	if result.ExpectedBranch != "feature/F067" {
		t.Errorf("ExpectedBranch = %q, want feature/F067", result.ExpectedBranch)
	}
	if result.RemoteTracking != "origin/feature/F067" {
		t.Errorf("RemoteTracking = %q, want origin/feature/F067", result.RemoteTracking)
	}
}

func TestExitCodeConstants(t *testing.T) {
	// Verify all exit codes are distinct
	codes := map[string]int{
		"ExitCodeOK":              ExitCodeOK,
		"ExitCodeContextMismatch": ExitCodeContextMismatch,
		"ExitCodeNoSession":       ExitCodeNoSession,
		"ExitCodeHashMismatch":    ExitCodeHashMismatch,
		"ExitCodeRuntimeError":    ExitCodeRuntimeError,
	}

	seen := make(map[int]string)
	for name, code := range codes {
		if existing, exists := seen[code]; exists {
			t.Errorf("Exit code %d is duplicated: %s and %s", code, existing, name)
		}
		seen[code] = name
	}

	// Verify expected values
	if ExitCodeOK != 0 {
		t.Errorf("ExitCodeOK = %d, want 0", ExitCodeOK)
	}
	if ExitCodeContextMismatch != 1 {
		t.Errorf("ExitCodeContextMismatch = %d, want 1", ExitCodeContextMismatch)
	}
}
