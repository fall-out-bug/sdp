package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunAcceptanceFromConfig_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()

	passed, skipped, err := runAcceptanceFromConfig(tmpDir)
	if err != nil {
		t.Errorf("Expected no error when config doesn't exist, got: %v", err)
	}
	if passed {
		t.Error("Expected passed=false when no config")
	}
	if !skipped {
		t.Error("Expected skipped=true when no config")
	}
}

func TestRunAcceptanceFromConfig_EmptyCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .sdp directory
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create config without acceptance command
	configContent := `version: 1
acceptance:
  command: ""
`
	if err := os.WriteFile(filepath.Join(sdpDir, "config.yml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	passed, skipped, err := runAcceptanceFromConfig(tmpDir)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if passed {
		t.Error("Expected passed=false with empty command")
	}
	if !skipped {
		t.Error("Expected skipped=true with empty command")
	}
}

func TestRunAcceptanceFromConfig_WithInvalidCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .sdp directory
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create config with invalid command
	configContent := `version: 1
acceptance:
  command: "nonexistent-command-that-should-fail"
  timeout: "5s"
`
	if err := os.WriteFile(filepath.Join(sdpDir, "config.yml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	passed, skipped, err := runAcceptanceFromConfig(tmpDir)
	// Should either skip or fail
	t.Logf("passed=%v, skipped=%v, err=%v", passed, skipped, err)
}
