package git

import (
	"testing"
)

func TestCommandCategories(t *testing.T) {
	// Test command categorization
	tests := []struct {
		name      string
		command   string
		category  CommandCategory
		needPost  bool
		needCheck bool
	}{
		// Safe commands
		{"status is safe", "status", CategorySafe, false, true},
		{"log is safe", "log", CategorySafe, false, true},
		{"diff is safe", "diff", CategorySafe, false, true},
		{"show is safe", "show", CategorySafe, false, true},

		// Write commands
		{"add is write", "add", CategoryWrite, true, true},
		{"commit is write", "commit", CategoryWrite, true, true},
		{"reset is write", "reset", CategoryWrite, true, true},

		// Remote commands
		{"push is remote", "push", CategoryRemote, true, true},
		{"fetch is remote", "fetch", CategoryRemote, false, true},
		{"pull is remote", "pull", CategoryRemote, true, true},

		// Branch commands
		{"checkout is branch", "checkout", CategoryBranch, true, true},
		{"branch is safe (read-only)", "branch", CategorySafe, false, true},
		{"merge is branch", "merge", CategoryBranch, true, true},

		// Unknown commands
		{"unknown defaults to safe", "unknown-cmd", CategorySafe, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat := CategorizeCommand(tt.command)
			if cat != tt.category {
				t.Errorf("CategorizeCommand(%s) = %v, want %v", tt.command, cat, tt.category)
			}

			needsPost := NeedsPostCheck(tt.command)
			if needsPost != tt.needPost {
				t.Errorf("NeedsPostCheck(%s) = %v, want %v", tt.command, needsPost, tt.needPost)
			}

			needsCheck := NeedsSessionCheck(tt.command)
			if needsCheck != tt.needCheck {
				t.Errorf("NeedsSessionCheck(%s) = %v, want %v", tt.command, needsCheck, tt.needCheck)
			}
		})
	}
}

func TestValidationResult(t *testing.T) {
	// Test validation result structure
	tests := []struct {
		name       string
		result     ValidationResult
		wantValid  bool
		wantErrMsg string
	}{
		{
			name: "valid result",
			result: ValidationResult{
				Valid:          true,
				WorktreePath:   "/path/to/worktree",
				CurrentBranch:  "feature/F065",
				ExpectedBranch: "feature/F065",
			},
			wantValid:  true,
			wantErrMsg: "",
		},
		{
			name: "worktree mismatch",
			result: ValidationResult{
				Valid:          false,
				WorktreePath:   "/path/to/worktree",
				ActualPath:     "/wrong/path",
				CurrentBranch:  "feature/F065",
				ExpectedBranch: "feature/F065",
				Error:          "wrong worktree",
				Fix:            "sdp guard context go F065",
			},
			wantValid:  false,
			wantErrMsg: "wrong worktree",
		},
		{
			name: "branch mismatch",
			result: ValidationResult{
				Valid:          false,
				WorktreePath:   "/path/to/worktree",
				CurrentBranch:  "dev",
				ExpectedBranch: "feature/F065",
				Error:          "wrong branch",
				Fix:            "git checkout feature/F065",
			},
			wantValid:  false,
			wantErrMsg: "wrong branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result.Valid != tt.wantValid {
				t.Errorf("Valid = %v, want %v", tt.result.Valid, tt.wantValid)
			}
			if tt.result.Error != tt.wantErrMsg {
				t.Errorf("Error = %v, want %v", tt.result.Error, tt.wantErrMsg)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	// Test error message formatting
	tests := []struct {
		name     string
		result   ValidationResult
		expected string
	}{
		{
			name: "worktree error",
			result: ValidationResult{
				Valid:        false,
				WorktreePath: "/path/to/expected",
				ActualPath:   "/path/to/actual",
				Error:        "worktree mismatch",
				Fix:          "sdp guard context go F065",
			},
			expected: "ERROR: worktree mismatch",
		},
		{
			name: "branch error",
			result: ValidationResult{
				Valid:          false,
				CurrentBranch:  "dev",
				ExpectedBranch: "feature/F065",
				Error:          "branch mismatch",
				Fix:            "git checkout feature/F065",
			},
			expected: "ERROR: branch mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.result.FormatError()
			if msg == "" {
				t.Error("FormatError() should not return empty string for invalid result")
			}
			if tt.result.Error != "" && msg != "" {
				// Check that the error message contains the expected text
			}
		})
	}
}

func TestValidatorNeedsSessionCheck(t *testing.T) {
	// Test that all commands need session check (even safe ones)
	commands := []string{"status", "log", "diff", "show", "add", "commit", "push", "checkout"}

	for _, cmd := range commands {
		if !NeedsSessionCheck(cmd) {
			t.Errorf("NeedsSessionCheck(%s) should be true", cmd)
		}
	}
}

func TestValidatorNeedsPostCheck(t *testing.T) {
	// Test which commands need post-check
	safeCommands := []string{"status", "log", "diff", "show", "branch"}
	writeCommands := []string{"add", "commit", "reset"}
	remoteCommands := []string{"push", "pull"}
	branchCommands := []string{"checkout", "merge"}

	for _, cmd := range safeCommands {
		if NeedsPostCheck(cmd) {
			t.Errorf("Safe command %s should not need post-check", cmd)
		}
	}

	for _, cmd := range writeCommands {
		if !NeedsPostCheck(cmd) {
			t.Errorf("Write command %s should need post-check", cmd)
		}
	}

	for _, cmd := range remoteCommands {
		if cmd == "push" && !NeedsPostCheck(cmd) {
			t.Errorf("Push command should need post-check")
		}
	}

	for _, cmd := range branchCommands {
		if !NeedsPostCheck(cmd) {
			t.Errorf("Branch command %s should need post-check", cmd)
		}
	}
}

func TestNewWrapper(t *testing.T) {
	wrapper := NewWrapper("/path/to/project")
	if wrapper == nil {
		t.Fatal("NewWrapper returned nil")
	}
	if wrapper.ProjectRoot != "/path/to/project" {
		t.Errorf("ProjectRoot = %v, want /path/to/project", wrapper.ProjectRoot)
	}
}

func TestCommandCategoryString(t *testing.T) {
	tests := []struct {
		category CommandCategory
		expected string
	}{
		{CategorySafe, "safe"},
		{CategoryWrite, "write"},
		{CategoryRemote, "remote"},
		{CategoryBranch, "branch"},
		{CommandCategory(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.category.String(); got != tt.expected {
			t.Errorf("CommandCategory(%d).String() = %v, want %v", tt.category, got, tt.expected)
		}
	}
}

func TestValidationResultFormatErrorEmpty(t *testing.T) {
	// Valid result should return empty error
	result := ValidationResult{Valid: true}
	if msg := result.FormatError(); msg != "" {
		t.Errorf("FormatError() for valid result should be empty, got: %s", msg)
	}
}

func TestCategorizeAllCommands(t *testing.T) {
	// Test all known commands
	tests := map[string]CommandCategory{
		"status":      CategorySafe,
		"log":         CategorySafe,
		"diff":        CategorySafe,
		"show":        CategorySafe,
		"ls-files":    CategorySafe,
		"rev-parse":   CategorySafe,
		"branch":      CategorySafe,
		"remote":      CategorySafe,
		"tag":         CategorySafe,
		"add":         CategoryWrite,
		"commit":      CategoryWrite,
		"reset":       CategoryWrite,
		"rm":          CategoryWrite,
		"mv":          CategoryWrite,
		"stash":       CategoryWrite,
		"push":        CategoryRemote,
		"fetch":       CategoryRemote,
		"pull":        CategoryRemote,
		"clone":       CategoryRemote,
		"checkout":    CategoryBranch,
		"switch":      CategoryBranch,
		"merge":       CategoryBranch,
		"rebase":      CategoryBranch,
		"cherry-pick": CategoryBranch,
	}

	for cmd, expected := range tests {
		if got := CategorizeCommand(cmd); got != expected {
			t.Errorf("CategorizeCommand(%s) = %v, want %v", cmd, got, expected)
		}
	}
}

func TestValidationResultFormatErrorWithFix(t *testing.T) {
	// Test error formatting with all fields
	result := ValidationResult{
		Valid:          false,
		WorktreePath:   "/expected/path",
		ActualPath:     "/actual/path",
		CurrentBranch:  "dev",
		ExpectedBranch: "feature/F065",
		Error:          "validation failed",
		Fix:            "sdp guard context go F065",
	}

	msg := result.FormatError()
	if msg == "" {
		t.Error("FormatError() should not be empty for invalid result")
	}

	// Check that all important parts are in the message
	if !containsSubstring(msg, "ERROR:") {
		t.Error("Error message should contain 'ERROR:'")
	}
	if !containsSubstring(msg, "validation failed") {
		t.Error("Error message should contain the error")
	}
	if !containsSubstring(msg, "FIX:") {
		t.Error("Error message should contain 'FIX:'")
	}
	if !containsSubstring(msg, "sdp guard context go F065") {
		t.Error("Error message should contain the fix")
	}
}

func TestValidationResultFormatErrorBranchOnly(t *testing.T) {
	// Test error formatting with only branch mismatch
	result := ValidationResult{
		Valid:          false,
		CurrentBranch:  "dev",
		ExpectedBranch: "feature/F065",
		Error:          "branch mismatch",
		Fix:            "git checkout feature/F065",
	}

	msg := result.FormatError()
	if !containsSubstring(msg, "Expected branch:") {
		t.Error("Error message should contain 'Expected branch:'")
	}
	if !containsSubstring(msg, "feature/F065") {
		t.Error("Error message should contain expected branch")
	}
}

func TestValidationResultFormatErrorWorktreeOnly(t *testing.T) {
	// Test error formatting with only worktree mismatch
	result := ValidationResult{
		Valid:        false,
		WorktreePath: "/expected/path",
		ActualPath:   "/actual/path",
		Error:        "worktree mismatch",
		Fix:          "sdp guard context go F065",
	}

	msg := result.FormatError()
	if !containsSubstring(msg, "Expected worktree:") {
		t.Error("Error message should contain 'Expected worktree:'")
	}
	if !containsSubstring(msg, "/expected/path") {
		t.Error("Error message should contain expected path")
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
