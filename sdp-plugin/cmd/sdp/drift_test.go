package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/drift"
)

// TestDriftDetectCmd tests the drift detect command
func TestDriftDetectCmd(t *testing.T) {
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
---
# Test Workstream

## Scope Files

**Implementation:**
- internal/file1.go
- internal/file2.go

**Tests:**
- internal/file1_test.go
`
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create the implementation files to avoid drift errors
	for _, file := range []string{"internal/file1.go", "internal/file2.go", "internal/file1_test.go"} {
		fullPath := filepath.Join(tmpDir, file)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		content := []byte("package main\n\nfunc TestFunction() {}\n")
		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
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
			errorMsg:    "requires exactly 1 arg",
		},
		{
			name:        "valid workstream with all files",
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
			// Recover from panic caused by cobra.ExactArgs validation
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("driftDetectCmd() panicked with: %v", r)
					}
				}
			}()

			cmd := driftDetectCmd()
			err := cmd.RunE(cmd, tt.args)

			if tt.expectError && err == nil {
				t.Errorf("driftDetectCmd() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("driftDetectCmd() unexpected error: %v", err)
			}
			if tt.expectError && tt.errorMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("driftDetectCmd() error = %q, should contain %q", err.Error(), tt.errorMsg)
				}
			}
		})
	}
}

// TestDriftDetectWithMissingFiles tests drift detection with missing files
func TestDriftDetectWithMissingFiles(t *testing.T) {
	// Create temp directory with workstream file that references non-existent files
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
---
# Test Workstream

## Scope Files

**Implementation:**
- internal/missing.go

**Tests:**
- internal/missing_test.go
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

	cmd := driftDetectCmd()
	err := cmd.RunE(cmd, []string{wsID})

	// Should error due to missing files (drift detected)
	if err == nil {
		t.Errorf("driftDetectCmd() expected error for missing files but got none")
	} else if !strings.Contains(err.Error(), "drift detected") {
		t.Errorf("driftDetectCmd() error = %q, should contain 'drift detected'", err.Error())
	}
}

// TestDriftReportString tests the report string output
func TestDriftReportString(t *testing.T) {
	report := &drift.DriftReport{
		WorkstreamID: "00-050-01",
		Verdict:      "WARNING",
		Issues: []drift.DriftIssue{
			{
				Status:        drift.StatusWarning,
				File:          "internal/file.go",
				Expected:      "File contains implementation",
				Actual:        "Found TODO comments",
				Recommendation: "Remove TODO comments",
			},
		},
	}

	output := report.String()

	// Check that output contains expected fields
	expectedStrings := []string{
		"00-050-01",
		"WARNING",
		"internal/file.go",
		"TODO comments",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Report String() should contain %q\nGot: %s", expected, output)
		}
	}
}
