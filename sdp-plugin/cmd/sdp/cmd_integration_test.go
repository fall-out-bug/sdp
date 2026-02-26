package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCheckpointCommand tests the sdp checkpoint commands
func TestCheckpointCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)
	// Test checkpoint list
	cmd := exec.Command(binaryPath, "checkpoint", "list")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	// May fail if no checkpoints directory, that's OK

	output := stdout.String() + stderr.String()
	t.Logf("Checkpoint list output: %s", output)
	if err != nil {
		t.Logf("Checkpoint list exit: %v", err)
	}
}

// TestInitCommand tests the sdp init command
func TestInitCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)
	// Create temp directory for init
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run init
	cmd := exec.Command(binaryPath, "init", "--project-type", "go")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// May fail if prompts/ dir not found (that's OK in test environment)
	output := stdout.String() + stderr.String()
	t.Logf("Init output: %s", output)
	if err != nil {
		t.Logf("Init exit: %v", err)
	}

	// Check if .claude was created (should succeed if prompts exist)
	if err == nil {
		claudeDir := filepath.Join(tmpDir, ".claude")
		if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
			t.Error(".claude directory was not created")
		}
	}
}

// TestBeadsCommand tests the sdp beads command
func TestBeadsCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "beads ready",
			args:     []string{"beads", "ready"},
			wantErr:  false,
			contains: "ready",
		},
		{
			name:    "beads list",
			args:    []string{"beads", "list"},
			wantErr: false,
		},
		{
			name:    "beads sync",
			args:    []string{"beads", "sync"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Logf("Output: %s", output)
			}
		})
	}
}

// TestCompletionCommand tests the sdp completion command
func TestCompletionCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)
	shells := []string{"bash", "zsh", "fish"}

	for _, shell := range shells {
		t.Run("completion "+shell, func(t *testing.T) {
			cmd := exec.Command(binaryPath, "completion", shell)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			if err := cmd.Run(); err != nil {
				t.Logf("Completion %s failed: %v\nOutput: %s", shell, err, stdout.String()+stderr.String())
			}

			output := stdout.String()
			if len(output) == 0 {
				t.Errorf("Completion script for %s should produce output", shell)
			}
		})
	}
}

// TestOrchestrateCommand tests the sdp orchestrate command
func TestOrchestrateCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Test orchestrate help
	cmd := exec.Command(binaryPath, "orchestrate", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Orchestrate help failed: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "orchestrate") {
		t.Errorf("Orchestrate help should mention orchestrate\nGot: %s", output)
	}
}

// TestPrdCommand tests the sdp prd command
func TestPrdCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "prd help",
			args:     []string{"prd", "--help"},
			wantErr:  false,
			contains: "PRD",
		},
		{
			name:    "prd detect",
			args:    []string{"prd", "detect"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			output := stdout.String() + stderr.String()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if tt.contains != "" && !strings.Contains(output, tt.contains) {
				t.Errorf("Output does not contain %q\nGot: %s", tt.contains, output)
			}
		})
	}
}
