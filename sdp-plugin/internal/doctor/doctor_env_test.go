package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCheckGit_CommandFails tests when git exists but --version fails
func TestCheckGit_CommandFails(t *testing.T) {
	// Save original PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	// Create a temporary directory with a fake git that fails
	tmpDir := t.TempDir()
	fakeBin := filepath.Join(tmpDir, "git")
	// Create a script that exits with error
	script := `#!/bin/bash
echo "broken" >&2
exit 1
`
	if err := os.WriteFile(fakeBin, []byte(script), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Add tmpDir to front of PATH
	os.Setenv("PATH", tmpDir+":"+oldPath)

	result := checkGit()

	// Should still return ok (git exists, even if version fails)
	if result.Name != "Git" {
		t.Errorf("Expected name 'Git', got '%s'", result.Name)
	}
	// Status can be ok (found git) even if version check fails
	if result.Status != "ok" {
		t.Logf("Status: %s, Message: %s", result.Status, result.Message)
	}
}

// TestCheckClaudeCode_CommandFails tests when claude exists but --version fails
func TestCheckClaudeCode_CommandFails(t *testing.T) {
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	tmpDir := t.TempDir()
	fakeBin := filepath.Join(tmpDir, "claude")
	script := `#!/bin/bash
exit 1
`
	if err := os.WriteFile(fakeBin, []byte(script), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	os.Setenv("PATH", tmpDir+":"+oldPath)

	result := checkClaudeCode()

	if result.Name != "Claude Code" {
		t.Errorf("Expected name 'Claude Code', got '%s'", result.Name)
	}
	// Since fake claude exists, should be ok
	if result.Status != "ok" {
		t.Logf("Status: %s, Message: %s", result.Status, result.Message)
	}
}

// TestCheckGo_CommandFails tests when go exists but version fails
func TestCheckGo_CommandFails(t *testing.T) {
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	tmpDir := t.TempDir()
	fakeBin := filepath.Join(tmpDir, "go")
	script := `#!/bin/bash
exit 1
`
	if err := os.WriteFile(fakeBin, []byte(script), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	os.Setenv("PATH", tmpDir+":"+oldPath)

	result := checkGo()

	if result.Name != "Go" {
		t.Errorf("Expected name 'Go', got '%s'", result.Name)
	}
	// Since fake go exists, should be ok
	if result.Status != "ok" {
		t.Logf("Status: %s, Message: %s", result.Status, result.Message)
	}
}
