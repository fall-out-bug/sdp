package context

import (
	"testing"
)

func TestExitCodes(t *testing.T) {
	// Test that exit codes are defined correctly
	if ExitCodeOK != 0 {
		t.Errorf("ExitCodeOK = %d, want 0", ExitCodeOK)
	}
	if ExitCodeContextMismatch != 1 {
		t.Errorf("ExitCodeContextMismatch = %d, want 1", ExitCodeContextMismatch)
	}
	if ExitCodeNoSession != 2 {
		t.Errorf("ExitCodeNoSession = %d, want 2", ExitCodeNoSession)
	}
	if ExitCodeHashMismatch != 3 {
		t.Errorf("ExitCodeHashMismatch = %d, want 3", ExitCodeHashMismatch)
	}
}

func TestNewRecovery(t *testing.T) {
	recovery := NewRecovery("/path/to/project")
	if recovery == nil {
		t.Fatal("NewRecovery returned nil")
	}
	if recovery.ProjectRoot != "/path/to/project" {
		t.Errorf("ProjectRoot = %v, want /path/to/project", recovery.ProjectRoot)
	}
}

func TestContextCheckResult(t *testing.T) {
	// Test ContextCheckResult structure
	result := &ContextCheckResult{
		Valid:          true,
		ExitCode:       ExitCodeOK,
		WorktreePath:   "/path/to/worktree",
		FeatureID:      "F065",
		CurrentBranch:  "feature/F065",
		ExpectedBranch: "feature/F065",
		RemoteTracking: "origin/feature/F065",
		SessionValid:   true,
		Errors:         nil,
	}

	if !result.Valid {
		t.Error("Valid should be true")
	}
	if result.ExitCode != ExitCodeOK {
		t.Errorf("ExitCode = %d, want %d", result.ExitCode, ExitCodeOK)
	}
	if result.FeatureID != "F065" {
		t.Errorf("FeatureID = %s, want F065", result.FeatureID)
	}
}

func TestExtractFeatureID(t *testing.T) {
	tests := []struct {
		branch   string
		expected string
	}{
		{"feature/F065", "F065"},
		{"feature/F123", "F123"},
		{"bugfix/sdp-1234", "sdp-1234"},
		{"hotfix/sdp-5678", "sdp-5678"},
		{"dev", ""},
		{"main", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.branch, func(t *testing.T) {
			result := extractFeatureID(tt.branch)
			if result != tt.expected {
				t.Errorf("extractFeatureID(%s) = %s, want %s", tt.branch, result, tt.expected)
			}
		})
	}
}

func TestFormatCheckResult(t *testing.T) {
	// Test valid result formatting
	validResult := &ContextCheckResult{
		Valid:          true,
		ExitCode:       ExitCodeOK,
		WorktreePath:   "/path/to/worktree",
		CurrentBranch:  "feature/F065",
		RemoteTracking: "origin/feature/F065",
		SessionValid:   true,
	}

	output := FormatCheckResult(validResult)
	if output == "" {
		t.Error("FormatCheckResult should not return empty string")
	}
	if !containsSubstring(output, "[OK]") {
		t.Error("Valid result should contain [OK] markers")
	}

	// Test invalid result formatting
	invalidResult := &ContextCheckResult{
		Valid:          false,
		ExitCode:       ExitCodeContextMismatch,
		WorktreePath:   "/path/to/worktree",
		CurrentBranch:  "dev",
		ExpectedBranch: "feature/F065",
		SessionValid:   true,
		Errors:         []string{"Branch mismatch"},
	}

	output = FormatCheckResult(invalidResult)
	if output == "" {
		t.Error("FormatCheckResult should not return empty string for invalid result")
	}
	if !containsSubstring(output, "[FAIL]") {
		t.Error("Invalid result should contain [FAIL] markers")
	}
}

func TestContextCheckResultWithErrors(t *testing.T) {
	// Test result with multiple errors
	result := &ContextCheckResult{
		Valid:   false,
		ExitCode: ExitCodeContextMismatch,
		Errors: []string{
			"Worktree mismatch",
			"Branch mismatch",
		},
	}

	if len(result.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result.Errors))
	}
}

func TestContextCheckResultInvalidSession(t *testing.T) {
	// Test result for invalid session
	result := &ContextCheckResult{
		Valid:        false,
		ExitCode:     ExitCodeHashMismatch,
		SessionValid: false,
		Errors:       []string{"Session file corrupted or tampered"},
	}

	if result.ExitCode != ExitCodeHashMismatch {
		t.Errorf("ExitCode = %d, want %d", result.ExitCode, ExitCodeHashMismatch)
	}
	if result.SessionValid {
		t.Error("SessionValid should be false for corrupted session")
	}
}

func TestContextCheckResultNoSession(t *testing.T) {
	// Test result when no session file exists
	result := &ContextCheckResult{
		Valid:    false,
		ExitCode: ExitCodeNoSession,
		Errors:   []string{"No session file found"},
	}

	if result.ExitCode != ExitCodeNoSession {
		t.Errorf("ExitCode = %d, want %d", result.ExitCode, ExitCodeNoSession)
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
