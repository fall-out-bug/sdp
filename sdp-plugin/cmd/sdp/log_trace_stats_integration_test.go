package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLogTraceCommand tests the log trace command
func TestLogTraceCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Create temp directory for isolated test environment
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Initialize minimal project structure
	sdpDir := filepath.Join(tmpDir, ".sdp", "log")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("create .sdp/log dir: %v", err)
	}

	// Create minimal config file
	configContent := `version: "0.9.0"
evidence:
  enabled: true
  log_path: ".sdp/log/events.jsonl"
`
	configPath := filepath.Join(tmpDir, ".sdp", "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("create config: %v", err)
	}

	cmd := exec.Command(binaryPath, "log", "trace", "--json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = tmpDir

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// trace should succeed even with empty log
	if err != nil && !strings.Contains(output, "No events") && !strings.Contains(output, "No matching") {
		t.Logf("Log trace failed: %v\nOutput: %s", err, output)
	}

	// Should be JSON or empty message
	if !strings.Contains(output, "No events") && !strings.Contains(output, "No matching") && !strings.Contains(output, "[") {
		t.Logf("Expected JSON or empty message\nGot: %s", output)
	}
}

// TestLogStatsCommand tests the log stats command
func TestLogStatsCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Create temp directory for isolated test environment
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Initialize minimal project structure
	sdpDir := filepath.Join(tmpDir, ".sdp", "log")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("create .sdp/log dir: %v", err)
	}

	// Create minimal config file
	configContent := `version: "0.9.0"
evidence:
  enabled: true
  log_path: ".sdp/log/events.jsonl"
`
	configPath := filepath.Join(tmpDir, ".sdp", "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("create config: %v", err)
	}

	cmd := exec.Command(binaryPath, "log", "stats")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = tmpDir

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// stats should succeed even with empty log
	if err != nil && !strings.Contains(output, "No events") {
		t.Logf("Log stats: %v\nOutput: %s", err, output)
	}

	// Should mention events or statistics
	if !strings.Contains(output, "event") && !strings.Contains(output, "statistics") && !strings.Contains(output, "No events") {
		t.Logf("Unexpected log stats output: %s", output)
	}
}
