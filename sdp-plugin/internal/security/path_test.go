package security

import (
	"testing"
)

// TestSanitizePath tests path sanitization
func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name:        "safe_path",
			input:       "docs/workstreams/backlog/00-001-01.md",
			expected:    "docs/workstreams/backlog/00-001-01.md",
			shouldError: false,
		},
		{
			name:        "path_traversal_single",
			input:       "../../../etc/passwd",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "path_traversal_double",
			input:       "00-050-01/../../../../../../etc/passwd",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "path_traversal_mixed",
			input:       "docs/../../etc/passwd",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "absolute_path_blocked",
			input:       "/etc/passwd",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "safe_relative_path",
			input:       "workstreams/active/00-050-01.md",
			expected:    "workstreams/active/00-050-01.md",
			shouldError: false,
		},
		{
			name:        "path_with_current_dir",
			input:       "./docs/workstreams/00-001-01.md",
			expected:    "docs/workstreams/00-001-01.md",
			shouldError: false,
		},
		{
			name:        "path_with_extra_slashes",
			input:       "docs//workstreams///backlog/00-001-01.md",
			expected:    "docs/workstreams/backlog/00-001-01.md",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SanitizePath(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("SanitizePath(%q) should error but got: %s", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("SanitizePath(%q) should not error but got: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("SanitizePath(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

// TestValidatePathInDirectory tests path containment validation
func TestValidatePathInDirectory(t *testing.T) {
	baseDir := "/home/user/project"

	tests := []struct {
		name        string
		baseDir     string
		inputPath   string
		shouldError bool
	}{
		{
			name:        "safe_path_within_base",
			baseDir:     baseDir,
			inputPath:   "/home/user/project/docs/workstreams/00-001-01.md",
			shouldError: false,
		},
		{
			name:        "path_traversal_outside_base",
			baseDir:     baseDir,
			inputPath:   "/home/user/project/../etc/passwd",
			shouldError: true,
		},
		{
			name:        "path_completely_outside",
			baseDir:     baseDir,
			inputPath:   "/etc/passwd",
			shouldError: true,
		},
		{
			name:        "safe_subdirectory",
			baseDir:     baseDir,
			inputPath:   "/home/user/project/internal/security/validator.go",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePathInDirectory(tt.baseDir, tt.inputPath)

			if tt.shouldError {
				if err == nil {
					t.Errorf("ValidatePathInDirectory(%q, %q) should error", tt.baseDir, tt.inputPath)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePathInDirectory(%q, %q) should not error but got: %v", tt.baseDir, tt.inputPath, err)
				}
			}
		})
	}
}

// TestSafeJoinPath tests safe path joining
func TestSafeJoinPath(t *testing.T) {
	baseDir := "/home/user/project"

	tests := []struct {
		name        string
		base        string
		input       string
		expected    string
		shouldError bool
	}{
		{
			name:        "safe_join",
			base:        baseDir,
			input:       "docs/workstreams/00-001-01.md",
			expected:    "/home/user/project/docs/workstreams/00-001-01.md",
			shouldError: false,
		},
		{
			name:        "join_with_traversal",
			base:        baseDir,
			input:       "../../../etc/passwd",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "join_with_absolute_path",
			base:        baseDir,
			input:       "/etc/passwd",
			expected:    "",
			shouldError: true,
		},
		{
			name:        "join_with_current_dir",
			base:        baseDir,
			input:       "./internal/security",
			expected:    "/home/user/project/internal/security",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeJoinPath(tt.base, tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("SafeJoinPath(%q, %q) should error but got: %s", tt.base, tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("SafeJoinPath(%q, %q) should not error but got: %v", tt.base, tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("SafeJoinPath(%q, %q) = %q, want %q", tt.base, tt.input, result, tt.expected)
				}
			}
		})
	}
}
