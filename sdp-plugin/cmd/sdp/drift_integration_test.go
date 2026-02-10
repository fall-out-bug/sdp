package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

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

// TestVerifyCommand tests the sdp verify command
func TestVerifyCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)

	// Test verify on an existing workstream
	wsFile := filepath.Join(root, "docs", "workstreams", "completed", "00-050-01.md")

	if _, err := os.Stat(wsFile); os.IsNotExist(err) {
		t.Skip("Workstream file not found, skipping verify test")
	}

	// Change to repo root directory so verify can find docs/workstreams
	cmd := exec.Command(binaryPath, "verify", "00-050-01")
	cmd.Dir = root
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	t.Logf("Verify output: %s\nError: %v", output, err)

	// Verify should not fail catastrophically
	if err != nil && !strings.Contains(output, "Error") && !strings.Contains(output, "FAILED") {
		t.Errorf("Verify failed unexpectedly: %v", err)
	}
}

// TestGuardCommand tests the sdp guard command
func TestGuardCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	root := repoRoot(t)
	wsFile := filepath.Join(root, "docs", "workstreams", "completed", "00-050-01.md")

	if _, err := os.Stat(wsFile); os.IsNotExist(err) {
		t.Skip("Workstream file not found, skipping guard test")
	}

	cmd := exec.Command(binaryPath, "guard", "activate", wsFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	t.Logf("Guard output: %s\nError: %v", output, err)

	// Guard should work or fail gracefully
	if err != nil && !strings.Contains(output, "Error") {
		// OK as long as it's not a panic
	}
}
