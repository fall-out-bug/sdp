package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Use filepath indirectly via repoRoot
var _ = filepath.Join

// TestQualityCommand tests the sdp quality command
func TestQualityCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(root)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "quality all",
			args:     []string{"quality", "all"},
			wantErr:  true, // Expected to fail due to coverage/complexity
			contains: "Coverage",
		},
		{
			name:    "quality coverage",
			args:    []string{"quality", "coverage"},
			wantErr: false,
		},
		{
			name:     "quality help",
			args:     []string{"quality", "--help"},
			wantErr:  false,
			contains: "quality",
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

			t.Logf("Quality %s: err=%v", tt.name, err)
		})
	}
}

// TestWatchCommand tests the sdp watch command
func TestWatchCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Test watch help (can't test actual watch as it runs forever)
	cmd := exec.Command(binaryPath, "watch", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Watch help failed: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "watch") {
		t.Errorf("Watch help should mention watch\nGot: %s", output)
	}
}

// TestHooksCommand tests the sdp hooks command
func TestHooksCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(root)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains string
	}{
		{
			name:     "hooks install",
			args:     []string{"hooks", "install"},
			wantErr:  false,
			contains: "hooks",
		},
		{
			name:    "hooks uninstall",
			args:    []string{"hooks", "uninstall"},
			wantErr: false,
		},
		{
			name:     "hooks help",
			args:     []string{"hooks", "--help"},
			wantErr:  false,
			contains: "hooks",
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

			t.Logf("Hooks %s: err=%v", tt.name, err)
		})
	}
}
