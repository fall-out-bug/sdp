package planner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPlanExecution_AC4 tests workstream file creation
func TestPlanExecution_AC4(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{
			{
				ID:          "00-057-01",
				Title:       "OAuth2 configuration",
				Description: "Setup OAuth2 provider credentials",
				Status:      "pending",
			},
		},
	}

	// Create workstream files
	err := p.CreateWorkstreamFiles(result)
	if err != nil {
		t.Fatalf("Failed to create workstream files: %v", err)
	}

	// Verify file created
	expectedPath := filepath.Join(backlogDir, "00-057-01-oauth2-configuration.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Workstream file not created at %s", expectedPath)
	}

	// Verify file content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read workstream file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "OAuth2 configuration") {
		t.Error("Workstream file missing title")
	}

	if !strings.Contains(contentStr, "00-057-01") {
		t.Error("Workstream file missing ID")
	}
}

// TestWorkstreamFilename tests filename generation
func TestWorkstreamFilename(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"OAuth2 Configuration", "00-001-01-oauth2-configuration.md"},
		{"Add User Authentication", "00-001-01-add-user-authentication.md"},
		{"Fix login bug", "00-001-01-fix-login-bug.md"},
		{"Multiple   Spaces", "00-001-01-multiple-spaces.md"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			ws := Workstream{ID: "00-001-01", Title: tt.title}
			filename := ws.Filename()
			if filename != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, filename)
			}
		})
	}
}

// TestCreateWorkstreamFiles_ErrorCases tests error handling
func TestCreateWorkstreamFiles_ErrorCases(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{
			{ID: "00-057-01", Title: "OAuth2 config", Description: "Setup", Status: "pending"},
		},
	}

	// Create file first
	err := p.CreateWorkstreamFiles(result)
	if err != nil {
		t.Fatalf("First creation failed: %v", err)
	}

	// Try to create again (should fail)
	err = p.CreateWorkstreamFiles(result)
	if err == nil {
		t.Error("Expected error when creating duplicate workstream file")
	}
}
