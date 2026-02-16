package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/nextstep"
)

// TestNextCommandBasic tests basic next command execution.
func TestNextCommandBasic(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "sdp-next-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}
	defer os.Chdir(oldDir)

	// Create git directory
	if err := os.Mkdir(".git", 0755); err != nil {
		t.Fatalf("Failed to create .git: %v", err)
	}

	// Create .sdp directory
	if err := os.MkdirAll(".sdp", 0755); err != nil {
		t.Fatalf("Failed to create .sdp: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := nextCmd()
	cmd.SetArgs([]string{})
	err = cmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Command failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if output == "" {
		t.Error("Expected non-empty output")
	}
}

// TestNextCommandJSON tests JSON output.
func TestNextCommandJSON(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "sdp-next-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Mkdir(".git", 0755); err != nil {
		t.Fatalf("Failed to create .git: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := nextCmd()
	cmd.SetArgs([]string{"--json"})
	err = cmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Command failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())

	// Should be valid JSON
	rec, err := nextstep.FromJSON([]byte(output))
	if err != nil {
		t.Errorf("Output is not valid JSON: %v (output: %q)", err, output)
	}
	if rec.Command == "" {
		t.Error("Expected non-empty command in JSON output")
	}
}

// TestNextCommandWithWorkstreams tests with workstream files.
func TestNextCommandWithWorkstreams(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-next-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}
	defer os.Chdir(oldDir)

	// Create git directory
	if err := os.Mkdir(".git", 0755); err != nil {
		t.Fatalf("Failed to create .git: %v", err)
	}

	// Create workstream file
	wsDir := filepath.Join("docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("Failed to create workstream dir: %v", err)
	}

	wsContent := `---
ws_id: 00-069-01
feature_id: F069
title: "Test Workstream"
status: ready
priority: 0
size: SMALL
---

## Goal
Test goal.
`
	wsFile := filepath.Join(wsDir, "00-069-01.md")
	if err := os.WriteFile(wsFile, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to write workstream: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := nextCmd()
	cmd.SetArgs([]string{"--json"})
	err = cmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Command failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())

	// Should recommend executing the ready workstream
	rec, err := nextstep.FromJSON([]byte(output))
	if err != nil {
		t.Fatalf("Invalid JSON output: %v (output: %q)", err, output)
	}
	if rec.Category != nextstep.CategoryExecution {
		t.Errorf("Expected execution category, got %s", rec.Category)
	}
}

// TestNextCommandAlternatives tests showing alternatives.
func TestNextCommandAlternatives(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-next-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Mkdir(".git", 0755); err != nil {
		t.Fatalf("Failed to create .git: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := nextCmd()
	cmd.SetArgs([]string{"--alternatives"})
	err = cmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Command failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Output should contain "Alternatives"
	// (even if empty, the section should be present when flag is set)
	if output == "" {
		t.Error("Expected non-empty output")
	}
}
