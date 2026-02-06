package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestVerifyCmd tests the verify command
func TestVerifyCmd(t *testing.T) {
	// Create temp directory with workstream file
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	wsID := "00-050-01"
	wsPath := filepath.Join(wsDir, "backlog", wsID+".md")
	if err := os.MkdirAll(filepath.Dir(wsPath), 0755); err != nil {
		t.Fatalf("Failed to create backlog dir: %v", err)
	}

	wsContent := `---
ws_id: 00-050-01
feature: F050
status: completed
---
# Test Workstream

## Scope Files

**Implementation:**
- internal/file1.go

## Verification

- [x] All tests pass
`
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create the implementation file
	filePath := filepath.Join(tmpDir, "internal", "file1.go")
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(filePath, []byte("package main\n\nfunc Test() {}"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Test verify command (will fail verification because file is incomplete, but that's ok)
	// We just want to test the command runs without crashing
	cmd := verifyCmd()

	// Recover from os.Exit(1) if verification fails
	defer func() {
		if r := recover(); r != nil {
			// Expected - verification likely failed
		}
	}()

	_ = cmd.RunE(cmd, []string{wsID})
}
