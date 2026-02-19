package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestDoctorHooksProvenanceSubcommandExists(t *testing.T) {
	cmd := doctorCmd()
	found := false
	for _, sub := range cmd.Commands() {
		if sub.Name() == "hooks-provenance" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("doctor command missing hooks-provenance subcommand")
	}
}

func TestDoctorHooksProvenanceRunE(t *testing.T) {
	tmpDir := t.TempDir()
	hooksDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("mkdir hooks: %v", err)
	}

	commitMsg := filepath.Join(hooksDir, "commit-msg")
	postCommit := filepath.Join(hooksDir, "post-commit")
	if err := os.WriteFile(commitMsg, []byte("#!/bin/sh\n# SDP-Agent\n# SDP-Model\n# SDP-Task\n"), 0755); err != nil {
		t.Fatalf("write commit-msg: %v", err)
	}
	if err := os.WriteFile(postCommit, []byte("#!/bin/sh\nsdp skill record\n# commit_sha\n# agent\n# model\n"), 0755); err != nil {
		t.Fatalf("write post-commit: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cmd := doctorHooksProvenanceCmd()
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Fatalf("doctor hooks-provenance failed: %v", err)
	}
}

func TestDoctorHooksProvenanceRunE_MissingHook(t *testing.T) {
	tmpDir := t.TempDir()
	hooksDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatalf("mkdir hooks: %v", err)
	}

	commitMsg := filepath.Join(hooksDir, "commit-msg")
	if err := os.WriteFile(commitMsg, []byte("#!/bin/sh\n# SDP-Agent\n# SDP-Model\n# SDP-Task\n"), 0755); err != nil {
		t.Fatalf("write commit-msg: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cmd := doctorHooksProvenanceCmd()
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Fatal("expected error when post-commit hook is missing")
	}
	if !strings.Contains(err.Error(), "failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}
