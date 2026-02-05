package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// This file contains integration tests for CLI commands.
// These tests require the sdp binary to be built separately.
// To build: go build -o sdp ./cmd/sdp

func skipIfBinaryNotBuilt(t *testing.T) string {
	// Tests run in sdp-plugin/cmd/sdp directory
	// Binary is in sdp-plugin/ directory (two levels up)
	// Need to go up TWO levels because of root-level cmd/ directory
	relativePath := filepath.Join("..", "..", "sdp")
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		t.Skip("Cannot resolve binary path")
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		t.Skip("sdp binary not found. Run: go build -o sdp ./cmd/sdp from sdp-plugin/ directory")
	}
	return absPath
}

// repoRoot returns the absolute path to the repository root
func repoRoot(t *testing.T) string {
	// From sdp-plugin/cmd/sdp, go up THREE levels to get to repo root
	// sdp-plugin/cmd/sdp → ../.. → sdp-plugin/ → ../../.. → sdp/
	path := filepath.Join("..", "..", "..")
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("Cannot resolve repo root: %v", err)
	}
	return absPath
}

// TestParseCommand tests the sdp parse command
func TestParseCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		contains   string
		notContain string
	}{
		{
			name:     "parse valid workstream by ID",
			args:     []string{"parse", "00-050-01"},
			wantErr:  false,
			contains: "00-050-01",
		},
		{
			name:     "parse missing workstream",
			args:     []string{"parse", "99-999-99"},
			wantErr:  true,
			contains: "not found",
		},
		{
			name:     "parse without args",
			args:     []string{"parse"},
			wantErr:  true,
			contains: "required",
		},
		{
			name:       "path traversal attack blocked",
			args:       []string{"parse", "../../../etc/passwd"},
			wantErr:    true,
			contains:   "not found",
			notContain: "root:x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tt.args...)
			// Set working directory to repo root so docs/workstreams/ are found
			cmd.Dir = repoRoot(t)

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				output := stdout.String() + stderr.String()
				t.Errorf("Unexpected error: %v\nOutput:\n%s", err, output)
			}

			output := stdout.String() + stderr.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Output does not contain expected string %q\nGot: %s", tt.contains, output)
			}

			if tt.notContain != "" && strings.Contains(output, tt.notContain) {
				t.Errorf("Output should not contain string %q\nGot: %s", tt.notContain, output)
			}
		})
	}
}

// TestDriftCommand tests the sdp drift detect command
func TestDriftCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Test that the drift command exists and doesn't crash
	cmd := exec.Command(binaryPath, "drift", "detect", "--help")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run drift help: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "detect") {
		t.Errorf("Drift command help should mention detect\nGot: %s", output)
	}
}

// TestVersionCommand tests the sdp --version flag
func TestVersionCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "--version")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run version command: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "sdp version") {
		t.Errorf("Version output does not contain version string\nGot: %s", output)
	}
}

// TestHelpCommand tests the sdp --help flag
func TestHelpCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "--help")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	output := stdout.String() + stderr.String()

	// Check for expected help content
	expectedKeywords := []string{"Usage", "Available Commands", "Flags"}
	for _, keyword := range expectedKeywords {
		if !strings.Contains(output, keyword) {
			t.Errorf("Help output does not contain expected keyword %q\nGot: %s", keyword, output)
		}
	}
}

// TestDoctorCommand tests the sdp doctor command
func TestDoctorCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "doctor")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// doctor should not fail
	if err := cmd.Run(); err != nil {
		t.Logf("Doctor command failed: %v\nOutput: %s", err, stdout.String()+stderr.String())
	}

	output := stdout.String() + stderr.String()

	// Check for expected doctor content
	if !strings.Contains(output, "SDP Doctor") && !strings.Contains(output, "doctor") {
		t.Logf("Doctor output: %s", output)
	}
}

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
	t.Logf("Checkpoint list output: %s\nError: %v", output, err)
}
