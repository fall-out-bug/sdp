package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstall(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that hooks were created
	expectedHooks := []string{"pre-commit", "pre-push"}
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			t.Errorf("Hook %s was not created", hookName)
			continue
		}

		// Check content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", hookPath, err)
		}

		if !strings.Contains(string(content), "# SDP Git Hook") {
			t.Errorf("Hook %s has wrong content: %s", hookName, string(content))
		}
	}
}

func TestInstall_NoGitDir(t *testing.T) {
	// Create temp directory WITHOUT .git
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install - should fail
	err := Install()
	if err == nil {
		t.Fatal("Install() should fail when .git doesn't exist")
	}

	if !strings.Contains(err.Error(), ".git directory not found") {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestInstall_SkipExisting(t *testing.T) {
	// Create temp directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create existing hook
	existingHook := filepath.Join(gitDir, "pre-commit")
	if err := os.WriteFile(existingHook, []byte("# existing hook"), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Run install - should skip existing hook
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that existing hook wasn't overwritten
	content, err := os.ReadFile(existingHook)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(content) != "# existing hook" {
		t.Errorf("Hook was overwritten, got: %s", string(content))
	}
}

func TestUninstall(t *testing.T) {
	// Create temp directory with hooks
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create hooks
	for _, hookName := range []string{"pre-commit", "pre-push"} {
		hookPath := filepath.Join(gitDir, hookName)
		if err := os.WriteFile(hookPath, []byte("# test hook"), 0755); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	// Run uninstall
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	// Check that hooks were removed
	for _, hookName := range []string{"pre-commit", "pre-push"} {
		hookPath := filepath.Join(gitDir, hookName)
		if _, err := os.Stat(hookPath); !os.IsNotExist(err) {
			t.Errorf("Hook %s was not removed", hookName)
		}
	}
}

func TestUninstall_NotExists(t *testing.T) {
	// Create temp directory WITHOUT hooks
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run uninstall - should not fail
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() should not fail when hooks don't exist: %v", err)
	}
}
