package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/doctor"
)

// TestDoctorCmd tests the doctor command
func TestDoctorCmd(t *testing.T) {
	originalRunWithOptions := doctorRunWithOptions
	doctorRunWithOptions = func(opts doctor.RunOptions) []doctor.CheckResult {
		if opts.DriftCheck {
			t.Fatal("doctor command unexpectedly enabled drift checks")
		}
		return []doctor.CheckResult{
			{Name: "Git", Status: "ok", Message: "Installed (git version test)"},
			{Name: "IDE integration", Status: "ok", Message: "Found: Codex"},
		}
	}
	t.Cleanup(func() {
		doctorRunWithOptions = originalRunWithOptions
	})

	cmd := doctorCmd()

	if cmd.Use != "doctor" {
		t.Errorf("doctorCmd() has wrong use: %s", cmd.Use)
	}
	if cmd.Flags().Lookup("drift") == nil {
		t.Error("doctorCmd() missing --drift flag")
	}

	output, err := captureCommandOutput(t, func() error {
		return cmd.RunE(cmd, []string{})
	})
	if err != nil {
		t.Errorf("doctorCmd() failed: %v", err)
	}
	for _, snippet := range []string{"SDP Environment Check", "✓ Git", "✓ IDE integration", "All required checks passed!"} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("doctor output missing %q:\n%s", snippet, output)
		}
	}
}

// TestDoctorCmdWithDriftFlag tests the doctor command with drift check enabled
func TestDoctorCmdWithDriftFlag(t *testing.T) {
	originalRunWithOptions := doctorRunWithOptions
	doctorRunWithOptions = func(opts doctor.RunOptions) []doctor.CheckResult {
		if !opts.DriftCheck {
			t.Fatal("doctor command did not pass drift flag to runner")
		}
		return []doctor.CheckResult{
			{Name: "Git", Status: "ok", Message: "Installed (git version test)"},
			{Name: "Drift Detection", Status: "ok", Message: "No recent workstreams to check"},
		}
	}
	t.Cleanup(func() {
		doctorRunWithOptions = originalRunWithOptions
	})

	cmd := doctorCmd()
	if err := cmd.Flags().Set("drift", "true"); err != nil {
		t.Fatalf("Failed to set drift flag: %v", err)
	}

	output, err := captureCommandOutput(t, func() error {
		return cmd.RunE(cmd, []string{})
	})
	if err != nil {
		t.Errorf("doctorCmd() with drift failed: %v", err)
	}
	if !strings.Contains(output, "✓ Drift Detection") {
		t.Fatalf("doctor drift output missing drift result:\n%s", output)
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
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatalf("mkdir hooks: %v", err)
	}

	commitMsg := filepath.Join(hooksDir, "commit-msg")
	postCommit := filepath.Join(hooksDir, "post-commit")
	if err := os.WriteFile(commitMsg, []byte("#!/bin/sh\n# SDP-Agent\n# SDP-Model\n# SDP-Task\n"), 0o755); err != nil {
		t.Fatalf("write commit-msg: %v", err)
	}
	if err := os.WriteFile(postCommit, []byte("#!/bin/sh\nsdp skill record\n# commit_sha\n# agent\n# model\n"), 0o755); err != nil {
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
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatalf("mkdir hooks: %v", err)
	}

	commitMsg := filepath.Join(hooksDir, "commit-msg")
	if err := os.WriteFile(commitMsg, []byte("#!/bin/sh\n# SDP-Agent\n# SDP-Model\n# SDP-Task\n"), 0o755); err != nil {
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

func captureCommandOutput(t *testing.T, run func() error) (string, error) {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	runErr := run()

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read output: %v", err)
	}
	return buf.String(), runErr
}
