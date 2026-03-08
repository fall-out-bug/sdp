package nextstep

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCollectWorkstreams tests workstream collection from files.
func TestCollectWorkstreams(t *testing.T) {
	// Create temp directory with workstream files
	tmpDir, err := os.MkdirTemp("", "sdp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create workstream directory structure
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("Failed to create workstream dir: %v", err)
	}

	// Create a sample workstream file
	wsContent := `---
ws_id: 00-069-01
feature_id: F069
title: "Test Workstream"
status: ready
priority: 0
size: SMALL
depends_on: []
---

## Goal
Test goal.

## Acceptance Criteria
- [ ] AC1: Test criterion
`
	wsFile := filepath.Join(wsDir, "00-069-01.md")
	if err := os.WriteFile(wsFile, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to write workstream file: %v", err)
	}

	collector := NewStateCollector(tmpDir)
	workstreams := collector.collectWorkstreams()

	if len(workstreams) == 0 {
		t.Error("Expected to collect workstreams, got none")
	}

	found := false
	for _, ws := range workstreams {
		if ws.ID == "00-069-01" {
			found = true
			if ws.Status != StatusReady {
				t.Errorf("Expected status ready, got %s", ws.Status)
			}
			if ws.Feature != "F069" {
				t.Errorf("Expected feature F069, got %s", ws.Feature)
			}
			if ws.Priority != 0 {
				t.Errorf("Expected priority 0, got %d", ws.Priority)
			}
			break
		}
	}

	if !found {
		t.Error("Expected to find workstream 00-069-01")
	}
}

// TestMapStatus tests status string mapping.
func TestMapStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected WorkstreamState
	}{
		{"backlog", StatusBacklog},
		{"BACKLOG", StatusBacklog},
		{"ready", StatusReady},
		{"open", StatusReady},
		{"in_progress", StatusInProgress},
		{"in-progress", StatusInProgress},
		{"started", StatusInProgress},
		{"blocked", StatusBlocked},
		{"completed", StatusCompleted},
		{"done", StatusCompleted},
		{"failed", StatusFailed},
		{"error", StatusFailed},
		{"unknown", StatusBacklog},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapStatus(tt.input)
			if result != tt.expected {
				t.Errorf("mapStatus(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCollectGitStatus tests git status collection.
func TestCollectGitStatus(t *testing.T) {
	// Non-git directory
	tmpDir, err := os.MkdirTemp("", "sdp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	collector := NewStateCollector(tmpDir)
	status := collector.collectGitStatus()

	if status.IsRepo {
		t.Error("Expected IsRepo=false for non-git directory")
	}

	// Git directory
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	status = collector.collectGitStatus()
	if !status.IsRepo {
		t.Error("Expected IsRepo=true for git directory")
	}
}

// TestCollectConfig tests config collection.
func TestCollectConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	collector := NewStateCollector(tmpDir)
	config := collector.collectConfig()

	if config.HasSDPConfig {
		t.Error("Expected HasSDPConfig=false without config")
	}

	// Create config directory
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.Mkdir(sdpDir, 0755); err != nil {
		t.Fatalf("Failed to create .sdp dir: %v", err)
	}

	// Create config file
	configFile := filepath.Join(sdpDir, "config.yml")
	if err := os.WriteFile(configFile, []byte("version: 0.10.0"), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config = collector.collectConfig()
	if !config.HasSDPConfig {
		t.Error("Expected HasSDPConfig=true with config.yml")
	}
}

// TestCollect tests full state collection.
func TestCollect(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	collector := NewStateCollector(tmpDir)
	state, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// Verify basic fields are populated
	if state.Config.ProjectRoot != tmpDir {
		t.Errorf("Expected ProjectRoot=%s, got %s", tmpDir, state.Config.ProjectRoot)
	}
	if state.Mode == "" {
		t.Error("Expected non-empty default mode")
	}
}

func TestCollectSessionLegacyFormat(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, ".sdp"), 0o755); err != nil {
		t.Fatalf("mkdir .sdp: %v", err)
	}
	content := `{"workstream_id":"00-069-01","feature_id":"F069","branch":"feature/F069"}`
	if err := os.WriteFile(filepath.Join(tmpDir, ".sdp", "session.json"), []byte(content), 0o644); err != nil {
		t.Fatalf("write session: %v", err)
	}

	collector := NewStateCollector(tmpDir)
	state, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}
	if state.Session == nil {
		t.Fatal("expected session to be collected")
	}
	if state.Session.WorkstreamID != "00-069-01" {
		t.Fatalf("WorkstreamID = %q, want 00-069-01", state.Session.WorkstreamID)
	}
	if state.ActiveWorkstream != "00-069-01" {
		t.Fatalf("ActiveWorkstream = %q, want 00-069-01", state.ActiveWorkstream)
	}
}

func TestCollectWorkstreamsIgnoresCompletedDependencies(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatalf("mkdir workstreams: %v", err)
	}

	completed := `---
ws_id: 00-069-01
feature_id: F069
status: completed
priority: 0
size: S
depends_on: []
---

## Goal

Complete blocker.

## Acceptance Criteria

- [x] Done
`
	ready := `---
ws_id: 00-069-02
feature_id: F069
status: ready
priority: 1
size: M
depends_on: ["00-069-01"]
---

## Goal

Ready after dependency completes.

## Acceptance Criteria

- [ ] Execute
`
	if err := os.WriteFile(filepath.Join(wsDir, "00-069-01.md"), []byte(completed), 0o644); err != nil {
		t.Fatalf("write completed workstream: %v", err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, "00-069-02.md"), []byte(ready), 0o644); err != nil {
		t.Fatalf("write ready workstream: %v", err)
	}

	collector := NewStateCollector(tmpDir)
	workstreams := collector.collectWorkstreams()

	for _, ws := range workstreams {
		if ws.ID != "00-069-02" {
			continue
		}
		if len(ws.BlockedBy) != 0 {
			t.Fatalf("BlockedBy = %#v, want no unmet blockers", ws.BlockedBy)
		}
		return
	}
	t.Fatal("expected to find workstream 00-069-02")
}
