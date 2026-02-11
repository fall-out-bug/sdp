package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// TestContainsPathTraversal tests path traversal detection
func TestContainsPathTraversal(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "normal path",
			path:     "docs/workstreams/backlog/00-050-01.md",
			expected: false,
		},
		{
			name:     "parent directory traversal",
			path:     "../../../etc/passwd",
			expected: true,
		},
		{
			name:     "windows backslash traversal",
			path:     "..\\..\\windows\\system32",
			expected: true,
		},
		{
			name:     "home directory access",
			path:     "~/.ssh/id_rsa",
			expected: true,
		},
		{
			name:     "absolute path",
			path:     "/etc/passwd",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsPathTraversal(tt.path)
			if result != tt.expected {
				t.Errorf("containsPathTraversal(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// TestContains tests substring search
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "substring exists",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "substring not exists",
			s:        "hello world",
			substr:   "foo",
			expected: false,
		},
		{
			name:     "empty substring",
			s:        "hello",
			substr:   "",
			expected: true,
		},
		{
			name:     "exact match",
			s:        "test",
			substr:   "test",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.s, tt.substr)
			if result != tt.expected {
				t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, result, tt.expected)
			}
		})
	}
}

// TestCountErrors tests error counting in validation issues
func TestCountErrors(t *testing.T) {
	tests := []struct {
		name     string
		issues   []parser.ValidationIssue
		expected int
	}{
		{
			name:     "no issues",
			issues:   []parser.ValidationIssue{},
			expected: 0,
		},
		{
			name: "only warnings",
			issues: []parser.ValidationIssue{
				{Severity: "WARNING", Field: "test", Message: "warning message"},
				{Severity: "WARNING", Field: "test2", Message: "another warning"},
			},
			expected: 0,
		},
		{
			name: "mixed issues",
			issues: []parser.ValidationIssue{
				{Severity: "WARNING", Field: "test", Message: "warning"},
				{Severity: "ERROR", Field: "test2", Message: "error"},
				{Severity: "ERROR", Field: "test3", Message: "another error"},
			},
			expected: 2,
		},
		{
			name: "only errors",
			issues: []parser.ValidationIssue{
				{Severity: "ERROR", Field: "test", Message: "error"},
				{Severity: "ERROR", Field: "test2", Message: "error2"},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countErrors(tt.issues)
			if result != tt.expected {
				t.Errorf("countErrors() = %d, want %d", result, tt.expected)
			}
		})
	}
}

// TestFindWorkstreamFile tests workstream file discovery
func TestFindWorkstreamFile(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()
	backlogDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test workstream file
	wsID := "00-050-01"
	wsPath := filepath.Join(backlogDir, wsID+".md")
	wsContent := `---
ws_id: 00-050-01
feature: F050
status: completed
---
# Test Workstream
`
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Test finding existing file
	foundPath, err := findWorkstreamFile(wsID)
	if err != nil {
		t.Errorf("findWorkstreamFile(%q) failed: %v", wsID, err)
	}
	if !strings.Contains(foundPath, wsID+".md") {
		t.Errorf("findWorkstreamFile(%q) = %q, should contain %q", wsID, foundPath, wsID+".md")
	}

	// Test finding non-existent file
	_, err = findWorkstreamFile("99-999-99")
	if err == nil {
		t.Error("findWorkstreamFile(99-999-99) should fail")
	}
}

// TestParseRun tests the parse command logic
func TestParseRun(t *testing.T) {
	// Create temp directory with workstream file
	tmpDir := t.TempDir()
	backlogDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	wsID := "00-050-01"
	wsPath := filepath.Join(backlogDir, wsID+".md")
	wsContent := `---
ws_id: 00-050-01
feature: F050
status: completed
size: MEDIUM
goal: Test workstream
acceptance_criteria:
  - AC1: Test passes
---
# Test Workstream
`
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
			errorMsg:    "required",
		},
		{
			name:        "path traversal attack",
			args:        []string{"../../../etc/passwd"},
			expectError: true,
			errorMsg:    "not found", // filepath.Clean() normalizes before traversal check
		},
		{
			name:        "valid workstream",
			args:        []string{wsID},
			expectError: false,
		},
		{
			name:        "non-existent workstream",
			args:        []string{"99-999-99"},
			expectError: true,
			errorMsg:    "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			var out bytes.Buffer
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			os.Stderr = w

			// Create fake command
			cmd := parseCmd()

			err := parseRun(cmd, tt.args)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout
			out.ReadFrom(r)

			if tt.expectError && err == nil {
				t.Errorf("parseRun(%v) expected error but got none", tt.args)
			}
			if !tt.expectError && err != nil {
				t.Errorf("parseRun(%v) unexpected error: %v", tt.args, err)
			}
			if tt.expectError && tt.errorMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("parseRun(%v) error = %q, should contain %q", tt.args, err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// TestValidateRun tests validation command logic
func TestValidateRun(t *testing.T) {
	// Create temp directory with workstream file
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "test.md")
	wsContent := `---
ws_id: 00-050-01
feature: F050
status: completed
---
# Test Workstream
`
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "no args",
			args:        []string{},
			expectError: true,
			errorMsg:    "required",
		},
		{
			name:        "path traversal attack",
			args:        []string{"../../../etc/passwd"},
			expectError: true,
			errorMsg:    "path traversal", // validateRun has explicit traversal check
		},
		{
			name:        "valid file",
			args:        []string{wsPath},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := parseCmd()
			err := validateRun(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("validateRun(%v) expected error but got none", tt.args)
			}
			if !tt.expectError && err != nil {
				t.Errorf("validateRun(%v) unexpected error: %v", tt.args, err)
			}
			if tt.expectError && tt.errorMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateRun(%v) error = %q, should contain %q", tt.args, err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// TestDisplayWorkstream tests workstream display
func TestDisplayWorkstream(t *testing.T) {
	ws := &parser.Workstream{
		ID:      "00-050-01",
		Feature: "F050",
		Status:  "completed",
		Size:    "MEDIUM",
		Goal:    "Test goal",
		Acceptance: []string{
			"AC1: First criterion",
			"AC2: Second criterion",
		},
		Scope: parser.Scope{
			Implementation: []string{"file1.go", "file2.go"},
			Tests:          []string{"file1_test.go"},
		},
	}

	// Capture output
	var out bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	displayWorkstream(ws)

	w.Close()
	os.Stdout = oldStdout
	out.ReadFrom(r)

	output := out.String()

	// Check that output contains expected fields
	expectedStrings := []string{
		"00-050-01",
		"F050",
		"completed",
		"MEDIUM",
		"Test goal",
		"AC1: First criterion",
		"AC2: Second criterion",
		"file1.go",
		"file2.go",
		"file1_test.go",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("displayWorkstream() output should contain %q\nGot: %s", expected, output)
		}
	}
}
