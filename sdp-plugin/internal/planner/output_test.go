package planner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPlanExecution_AC6 tests JSON output
func TestPlanExecution_AC6(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		OutputFormat: "json",
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{
			{ID: "00-057-01", Title: "OAuth2 config", Description: "Setup", Status: "pending"},
		},
		Dependencies: []Dependency{
			{From: "00-057-02", To: "00-057-01", Reason: "Needs config"},
		},
	}

	// Generate JSON output
	output, err := p.FormatOutput(result)
	if err != nil {
		t.Fatalf("Failed to format output: %v", err)
	}

	// Verify valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	// Verify structure
	if _, ok := parsed["workstreams"]; !ok {
		t.Error("Missing 'workstreams' key in JSON output")
	}

	if _, ok := parsed["dependencies"]; !ok {
		t.Error("Missing 'dependencies' key in JSON output")
	}
}

// TestFormatHumanOutput tests human-readable output
func TestFormatHumanOutput(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		OutputFormat: "human",
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{
			{ID: "00-057-01", Title: "OAuth2 config", Description: "Setup", Status: "pending"},
		},
	}

	output, err := p.FormatOutput(result)
	if err != nil {
		t.Fatalf("Failed to format output: %v", err)
	}

	if !strings.Contains(output, "OAuth2 config") {
		t.Error("Human output missing workstream title")
	}

	if !strings.Contains(output, "00-057-01") {
		t.Error("Human output missing workstream ID")
	}
}

// TestFormatOutput_UnsupportedFormat tests unsupported format error
func TestFormatOutput_UnsupportedFormat(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		OutputFormat: "yaml",
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{},
	}

	_, err := p.FormatOutput(result)
	if err == nil {
		t.Error("Expected error for unsupported format")
	}

	if !containsString(err.Error(), "unsupported") {
		t.Errorf("Error should mention unsupported format, got: %v", err)
	}
}
