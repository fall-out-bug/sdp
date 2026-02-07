package security

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPathTraversalPrevention tests that path traversal attacks are blocked
func TestPathTraversalPrevention(t *testing.T) {
	// Create a temporary base directory for testing
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		userInput   string
		shouldBlock bool
	}{
		{
			name:        "safe_workstream_file",
			userInput:   "docs/workstreams/backlog/00-001-01.md",
			shouldBlock: false,
		},
		{
			name:        "traversal_attack_etc_passwd",
			userInput:   "../../../etc/passwd",
			shouldBlock: true,
		},
		{
			name:        "traversal_attack_with_context",
			userInput:   "workstreams/active/../../../etc/passwd",
			shouldBlock: true,
		},
		{
			name:        "absolute_path_attack",
			userInput:   "/etc/passwd",
			shouldBlock: true,
		},
		{
			name:        "safe_subdirectory",
			userInput:   "internal/security/validator.go",
			shouldBlock: false,
		},
		{
			name:        "complex_traversal",
			userInput:   "./docs/../../etc/passwd",
			shouldBlock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SafeJoinPath(tmpDir, tt.userInput)

			if tt.shouldBlock {
				if err == nil {
					t.Errorf("SafeJoinPath should block '%s' but got: %s", tt.userInput, result)

					// Additional safety check: verify file doesn't exist
					if _, statErr := os.Stat(result); statErr == nil {
						t.Errorf("Blocked path '%s' actually exists! Security breach!", result)
					}
				}
			} else {
				if err != nil {
					t.Errorf("SafeJoinPath should allow '%s' but got error: %v", tt.userInput, err)
				}
			}
		})
	}
}

// TestSafeFileOperations tests safe file operation patterns
func TestSafeFileOperations(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a safe subdirectory
	safeSubDir := filepath.Join(tmpDir, "docs", "workstreams")
	err := os.MkdirAll(safeSubDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(safeSubDir, "test.md")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("safe_file_read", func(t *testing.T) {
		// Safe: user provides relative path within base
		userPath := "docs/workstreams/test.md"
		safePath, err := SafeJoinPath(tmpDir, userPath)
		if err != nil {
			t.Fatalf("SafeJoinPath failed: %v", err)
		}

		// Verify file exists and is readable
		content, err := os.ReadFile(safePath)
		if err != nil {
			t.Errorf("Failed to read safe file: %v", err)
		}
		if string(content) != "test content" {
			t.Errorf("Unexpected content: %s", string(content))
		}
	})

	t.Run("blocked_file_read", func(t *testing.T) {
		// Unsafe: user tries to read system file
		userPath := "../../../etc/passwd"
		_, err := SafeJoinPath(tmpDir, userPath)
		if err == nil {
			t.Error("SafeJoinPath should block traversal attack")
		}
	})

	t.Run("safe_file_write", func(t *testing.T) {
		userPath := "docs/workstreams/new.md"
		safePath, err := SafeJoinPath(tmpDir, userPath)
		if err != nil {
			t.Fatalf("SafeJoinPath failed: %v", err)
		}

		// Write file safely
		content := []byte("new content")
		err = os.WriteFile(safePath, content, 0644)
		if err != nil {
			t.Errorf("Failed to write safe file: %v", err)
		}

		// Verify file was written
		readContent, err := os.ReadFile(safePath)
		if err != nil {
			t.Errorf("Failed to read written file: %v", err)
		}
		if string(readContent) != "new content" {
			t.Errorf("Unexpected content: %s", string(readContent))
		}
	})
}

// TestValidateUserProvidedPaths tests common user input scenarios
func TestValidateUserProvidedPaths(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldPass  bool
		description string
	}{
		{
			name:        "workstream_id",
			input:       "00-001-01",
			shouldPass:  true,
			description: "Workstream ID (not a path)",
		},
		{
			name:        "file_path",
			input:       "docs/workstreams/00-001-01.md",
			shouldPass:  true,
			description: "Relative file path",
		},
		{
			name:        "path_with_dots",
			input:       "./workstreams/active",
			shouldPass:  true,
			description: "Path with current directory reference",
		},
		{
			name:        "traversal_attack",
			input:       "../secrets.txt",
			shouldPass:  false,
			description: "Path traversal attempt",
		},
		{
			name:        "windows_traversal",
			input:       "..\\windows\\system32",
			shouldPass:  false,
			description: "Windows-style traversal",
		},
		{
			name:        "absolute_path",
			input:       "/usr/local/bin",
			shouldPass:  false,
			description: "Absolute path (not allowed)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SanitizePath(tt.input)

			if tt.shouldPass {
				if err != nil {
					t.Errorf("%s: should pass but got error: %v", tt.description, err)
				}
			} else {
				if err == nil {
					t.Errorf("%s: should fail but passed", tt.description)
				}
			}
		})
	}
}
