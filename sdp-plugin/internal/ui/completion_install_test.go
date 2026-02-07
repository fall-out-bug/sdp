package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstallCompletion_Bash tests bash completion installation
func TestInstallCompletion_Bash(t *testing.T) {
	// Save original HOME
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Create temp directory for home
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	err := InstallCompletion(Bash)

	if err != nil {
		t.Errorf("InstallCompletion(Bash) error = %v", err)
	}

	// Check completion script was created
	scriptPath := filepath.Join(tmpDir, ".bash_completion.d", "sdp")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Errorf("Bash completion script not created at %s", scriptPath)
	}

	// Verify script content
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read completion script: %v", err)
	}
	script := string(content)

	if !strings.Contains(script, "_sdp_completion") {
		t.Error("Bash completion should contain _sdp_completion function")
	}
	if !strings.Contains(script, "complete -F _sdp_completion sdp") {
		t.Error("Bash completion should register completion for sdp command")
	}
}

// TestInstallCompletion_Zsh tests zsh completion installation
func TestInstallCompletion_Zsh(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	err := InstallCompletion(Zsh)

	if err != nil {
		t.Errorf("InstallCompletion(Zsh) error = %v", err)
	}

	// Check completion script was created (tries both dirs)
	scriptPath1 := filepath.Join(tmpDir, ".zsh", "completion", "_sdp")
	scriptPath2 := filepath.Join(tmpDir, ".zsh", "completions", "_sdp")

	scriptPath := scriptPath1
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		scriptPath = scriptPath2
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Errorf("Zsh completion script not created (tried %s and %s)", scriptPath1, scriptPath2)
	}

	// Verify script content
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read completion script: %v", err)
	}
	script := string(content)

	if !strings.Contains(script, "#compdef sdp") {
		t.Error("Zsh completion should contain #compdef sdp directive")
	}
}

// TestInstallCompletion_Fish tests fish completion installation
func TestInstallCompletion_Fish(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	err := InstallCompletion(Fish)

	if err != nil {
		t.Errorf("InstallCompletion(Fish) error = %v", err)
	}

	// Check completion script was created
	scriptPath := filepath.Join(tmpDir, ".config", "fish", "completions", "sdp.fish")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Errorf("Fish completion script not created at %s", scriptPath)
	}

	// Verify script content
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read completion script: %v", err)
	}
	script := string(content)

	if !strings.Contains(script, "complete -c sdp") {
		t.Error("Fish completion should register completion for sdp command")
	}
}

// TestInstallCompletion_NoHome tests error when HOME not set
func TestInstallCompletion_NoHome(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	// Unset HOME
	os.Unsetenv("HOME")

	err := InstallCompletion(Bash)

	if err == nil {
		t.Error("InstallCompletion() should return error when HOME not set")
	}
	if !strings.Contains(err.Error(), "HOME environment variable not set") {
		t.Errorf("Error message should mention HOME, got: %v", err)
	}
}

// TestInstallCompletion_UnsupportedShell tests error for unsupported shell
func TestInstallCompletion_UnsupportedShell(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Create invalid shell type
	invalidShell := CompletionType("invalid")

	err := InstallCompletion(invalidShell)

	if err == nil {
		t.Error("InstallCompletion() should return error for unsupported shell")
	}
	if !strings.Contains(err.Error(), "unsupported shell") {
		t.Errorf("Error message should mention unsupported shell, got: %v", err)
	}
}

// TestInstallCompletion_CreatesDirectory tests that directories are created if needed
func TestInstallCompletion_CreatesDirectory(t *testing.T) {
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Don't create the completion directory beforehand
	err := InstallCompletion(Bash)

	if err != nil {
		t.Errorf("InstallCompletion(Bash) should create directory, error = %v", err)
	}

	// Verify directory was created
	dir := filepath.Join(tmpDir, ".bash_completion.d")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Directory should be created at %s", dir)
	}
}
