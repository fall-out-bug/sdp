package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveCmd_Workstream(t *testing.T) {
	// Create temp directory with workstream
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	wsContent := `---
ws_id: 00-064-01
feature_id: F064
title: "Test Workstream"
status: backlog
---
## Goal
Test goal
`
	if err := os.WriteFile(filepath.Join(wsDir, "00-064-01.md"), []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp directory
	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	cmd := resolveCmd()
	cmd.SetArgs([]string{"00-064-01"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("resolveCmd.Execute() error = %v", err)
	}
}

func TestResolveCmd_JSON(t *testing.T) {
	// Create temp directory with workstream
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	wsContent := `---
ws_id: 00-064-01
title: "Test"
---
## Goal
Test
`
	if err := os.WriteFile(filepath.Join(wsDir, "00-064-01.md"), []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	cmd := resolveCmd()
	cmd.SetArgs([]string{"--json", "00-064-01"})

	if err := cmd.Execute(); err != nil {
		t.Errorf("resolveCmd.Execute() --json error = %v", err)
	}
}

func TestResolveCmd_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	cmd := resolveCmd()
	cmd.SetArgs([]string{"00-999-99"})

	if err := cmd.Execute(); err == nil {
		t.Error("resolveCmd.Execute() expected error for non-existent workstream")
	}
}

func TestResolveCmd_UnknownID(t *testing.T) {
	cmd := resolveCmd()
	cmd.SetArgs([]string{"invalid-id-format"})

	if err := cmd.Execute(); err == nil {
		t.Error("resolveCmd.Execute() expected error for unknown ID format")
	}
}
