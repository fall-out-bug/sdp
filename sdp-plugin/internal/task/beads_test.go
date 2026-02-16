package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBeadsIntegration_IsEnabled(t *testing.T) {
	b := NewBeadsIntegration()
	// Just verify it doesn't panic
	_ = b.IsEnabled()
}

func TestBeadsIntegration_GenerateBeadsID(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"Fix CI Go version", "sdp-fix-ci"},
		{"Add user authentication", "sdp-add-user"},
		{"Test", "sdp-test"},
		{"A", "sdp-a"},
		{"", "sdp-unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			result := generateBeadsID(tt.title)
			if !strings.HasPrefix(result, "sdp-") {
				t.Errorf("expected beads ID to start with 'sdp-', got %s", result)
			}
		})
	}
}

func TestBeadsIntegration_LinkWorkstream(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a workstream file without beads_id
	wsContent := `---
ws_id: 00-064-01
feature_id: F064
title: "Test"
project_id: sdp
---
## Goal
Test
`
	wsPath := filepath.Join(tmpDir, "00-064-01.md")
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	b := &BeadsIntegration{enabled: true}

	err := b.LinkWorkstreamToBeads(wsPath, "sdp-test123")
	if err != nil {
		t.Fatalf("LinkWorkstreamToBeads() error = %v", err)
	}

	// Verify beads_id was added
	content, err := os.ReadFile(wsPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "beads_id: sdp-test123") {
		t.Error("expected beads_id to be added to workstream file")
		t.Logf("Content:\n%s", string(content))
	}
}

func TestBeadsIntegration_LinkWorkstream_Disabled(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "00-064-01.md")
	if err := os.WriteFile(wsPath, []byte("---\nws_id: test\n---"), 0644); err != nil {
		t.Fatal(err)
	}

	b := &BeadsIntegration{enabled: false}

	// Should not fail when disabled
	err := b.LinkWorkstreamToBeads(wsPath, "sdp-test")
	if err != nil {
		t.Errorf("expected no error when beads disabled, got: %v", err)
	}
}

func TestReadBeadsMapping(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		mappingPath := filepath.Join(tmpDir, "mapping.jsonl")
		content := `{"sdp_id":"00-064-01","beads_id":"sdp-abc","updated_at":"2026-02-12"}
{"sdp_id":"00-064-02","beads_id":"sdp-xyz","updated_at":"2026-02-12"}
`
		if err := os.WriteFile(mappingPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		mapping, err := ReadBeadsMapping(mappingPath)
		if err != nil {
			t.Fatalf("ReadBeadsMapping() error = %v", err)
		}

		if mapping["00-064-01"] != "sdp-abc" {
			t.Errorf("expected mapping[00-064-01] = sdp-abc, got %s", mapping["00-064-01"])
		}
		if mapping["sdp-xyz"] != "00-064-02" {
			t.Errorf("expected reverse mapping")
		}
	})

	t.Run("file missing", func(t *testing.T) {
		mapping, err := ReadBeadsMapping("/nonexistent/path.jsonl")
		if err != nil {
			t.Fatalf("expected no error for missing file, got: %v", err)
		}
		if len(mapping) != 0 {
			t.Error("expected empty mapping for missing file")
		}
	})
}

func TestCreator_CreateWorkstreamWithBeads(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	c := NewCreator(CreatorConfig{
		WorkstreamDir: wsDir,
		ProjectID:     "00",
	})

	task := &Task{
		Type:      TypeBug,
		Title:     "Test Bug",
		Priority:  PriorityP1,
		FeatureID: "F064",
		Goal:      "Fix test",
	}

	ws, err := c.CreateWorkstreamWithBeads(task)
	if err != nil {
		t.Fatalf("CreateWorkstreamWithBeads() error = %v", err)
	}

	if ws.WSID == "" {
		t.Error("expected non-empty WSID")
	}
}

func TestInsertBeadsID(t *testing.T) {
	tests := []struct {
		name    string
		content string
		beadsID string
	}{
		{
			"simple frontmatter",
			"---\nws_id: 00-064-01\n---\n## Goal\nTest",
			"sdp-test",
		},
		{
			"with project_id",
			"---\nws_id: 00-064-01\nproject_id: sdp\n---\n## Goal\nTest",
			"sdp-test",
		},
		{
			"no frontmatter",
			"No frontmatter here",
			"sdp-test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := insertBeadsID(tt.content, tt.beadsID)
			if strings.Contains(tt.content, "---") && !strings.Contains(result, "beads_id: "+tt.beadsID) {
				t.Errorf("expected beads_id to be inserted")
				t.Logf("Result:\n%s", result)
			}
		})
	}
}

func TestUpdateBeadsIDLine(t *testing.T) {
	content := "---\nws_id: 00-064-01\nbeads_id: old-value\n---\n"
	result := updateBeadsIDLine(content, "sdp-new")
	if !strings.Contains(result, "beads_id: sdp-new") {
		t.Errorf("expected beads_id to be updated")
	}
	if strings.Contains(result, "old-value") {
		t.Error("expected old beads_id to be replaced")
	}
}
