package security

import (
	"context"
	"testing"
)

// TestMustSafeCommand tests MustSafeCommand panic behavior
func TestMustSafeCommand(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustSafeCommand should panic on invalid command")
		}
	}()

	// This should panic
	MustSafeCommand(context.Background(), "rm", []string{"-rf", "/"}...)
}

// TestMustSafeCommandValid tests MustSafeCommand with valid command
func TestMustSafeCommandValid(t *testing.T) {
	ctx := context.Background()

	// This should not panic
	cmd := MustSafeCommand(ctx, "git", []string{"--version"}...)
	if cmd == nil {
		t.Error("MustSafeCommand should return non-nil cmd")
	}
}

// TestSanitizePathEdgeCases tests edge cases for SanitizePath
func TestSanitizePathEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "empty_string",
			input:       "",
			shouldError: true,
		},
		{
			name:        "single_dot",
			input:       ".",
			shouldError: false,
		},
		{
			name:        "double_dot_only",
			input:       "..",
			shouldError: true,
		},
		{
			name:        "trailing_slash",
			input:       "docs/workstreams/",
			shouldError: false,
		},
		{
			name:        "double_slash",
			input:       "docs//workstreams",
			shouldError: false,
		},
		{
			name:        "triple_dot",
			input:       "...",
			shouldError: true, // Starts with ".." which is blocked
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SanitizePath(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("SanitizePath(%q) should error", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("SanitizePath(%q) should not error: %v", tt.input, err)
				}
			}
		})
	}
}

// TestValidatePathInDirectoryErrors tests error cases for ValidatePathInDirectory
func TestValidatePathInDirectoryErrors(t *testing.T) {
	tests := []struct {
		name        string
		baseDir     string
		targetPath  string
		shouldError bool
	}{
		{
			name:        "empty_base_dir",
			baseDir:     "",
			targetPath:  "/some/path",
			shouldError: true,
		},
		{
			name:        "empty_target_path",
			baseDir:     "/home/user",
			targetPath:  "",
			shouldError: true,
		},
		{
			name:        "both_empty",
			baseDir:     "",
			targetPath:  "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePathInDirectory(tt.baseDir, tt.targetPath)

			if tt.shouldError {
				if err == nil {
					t.Error("ValidatePathInDirectory should error")
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePathInDirectory should not error: %v", err)
				}
			}
		})
	}
}

// TestSafeJoinPathErrors tests error cases for SafeJoinPath
func TestSafeJoinPathErrors(t *testing.T) {
	tests := []struct {
		name        string
		base        string
		userPath    string
		shouldError bool
	}{
		{
			name:        "empty_base",
			base:        "",
			userPath:    "docs/file.md",
			shouldError: true,
		},
		{
			name:        "empty_user_path",
			base:        "/home/user",
			userPath:    "",
			shouldError: true,
		},
		{
			name:        "both_empty",
			base:        "",
			userPath:    "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SafeJoinPath(tt.base, tt.userPath)

			if tt.shouldError {
				if err == nil {
					t.Error("SafeJoinPath should error")
				}
			} else {
				if err != nil {
					t.Errorf("SafeJoinPath should not error: %v", err)
				}
			}
		})
	}
}

// TestValidateTestCommandEdgeCases tests edge cases for ValidateTestCommand
func TestValidateTestCommandEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		testCmd     string
		shouldError bool
	}{
		{
			name:        "empty_command",
			testCmd:     "",
			shouldError: true,
		},
		{
			name:        "whitespace_only",
			testCmd:     "   ",
			shouldError: true,
		},
		{
			name:        "command_only_no_args",
			testCmd:     "pytest",
			shouldError: false,
		},
		{
			name:        "command_with_single_arg",
			testCmd:     "go test",
			shouldError: false,
		},
		{
			name:        "command_with_multiple_args",
			testCmd:     "pytest tests/ -v --cov",
			shouldError: false,
		},
		{
			name:        "trailing_injection",
			testCmd:     "pytest tests/",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTestCommand(tt.testCmd)

			if tt.shouldError {
				if err == nil {
					t.Errorf("ValidateTestCommand(%q) should error", tt.testCmd)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateTestCommand(%q) should not error: %v", tt.testCmd, err)
				}
			}
		})
	}
}
