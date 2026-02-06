package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDoctorCmd tests the doctor command
func TestDoctorCmd(t *testing.T) {
	// Create .claude directory for doctor checks
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude", "skills"), 0755); err != nil {
		t.Fatalf("Failed to create .claude dir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := doctorCmd()

	// Test command structure
	if cmd.Use != "doctor" {
		t.Errorf("doctorCmd() has wrong use: %s", cmd.Use)
	}

	// Test flag exists
	if cmd.Flags().Lookup("drift") == nil {
		t.Error("doctorCmd() missing --drift flag")
	}

	// Test that command runs without crashing
	err := cmd.RunE(cmd, []string{})
	// Should succeed (all required checks should pass with .claude present)
	if err != nil {
		t.Errorf("doctorCmd() failed: %v", err)
	}
}

// TestDoctorCmdWithDriftFlag tests the doctor command with drift check enabled
func TestDoctorCmdWithDriftFlag(t *testing.T) {
	// Create .claude directory for doctor checks
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, ".claude", "skills"), 0755); err != nil {
		t.Fatalf("Failed to create .claude dir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := doctorCmd()

	// Set drift flag
	if err := cmd.Flags().Set("drift", "true"); err != nil {
		t.Fatalf("Failed to set drift flag: %v", err)
	}

	// Test that command runs without crashing
	err := cmd.RunE(cmd, []string{})
	// Should succeed (all required checks should pass with .claude present)
	if err != nil {
		t.Errorf("doctorCmd() with drift failed: %v", err)
	}
}

