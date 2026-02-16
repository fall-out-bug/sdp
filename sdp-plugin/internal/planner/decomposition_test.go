package planner

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPlanExecution_AC1 tests basic decomposition output
func TestPlanExecution_AC1(t *testing.T) {
	// Create temp directory for testing
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	// Mock decomposition
	result := &DecompositionResult{
		Workstreams: []Workstream{
			{
				ID:          "00-057-01",
				Title:       "OAuth2 configuration",
				Description: "Setup OAuth2 provider credentials",
				Status:      "pending",
			},
			{
				ID:          "00-057-02",
				Title:       "OAuth2 callback handler",
				Description: "Implement OAuth2 callback endpoint",
				Status:      "pending",
			},
		},
		Dependencies: []Dependency{
			{From: "00-057-02", To: "00-057-01", Reason: "Needs config first"},
		},
	}

	// Verify result structure
	if len(result.Workstreams) != 2 {
		t.Errorf("Expected 2 workstreams, got %d", len(result.Workstreams))
	}

	if result.Workstreams[0].ID != "00-057-01" {
		t.Errorf("Expected ID 00-057-01, got %s", result.Workstreams[0].ID)
	}

	if len(result.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(result.Dependencies))
	}
}

// TestPlanExecution_AC2 tests interactive mode
func TestPlanExecution_AC2(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		Interactive: true,
		Questions:   []string{"Which OAuth2 provider?", "What scopes?"},
	}

	if !p.Interactive {
		t.Error("Interactive mode not set")
	}

	if len(p.Questions) != 2 {
		t.Errorf("Expected 2 questions, got %d", len(p.Questions))
	}
}

// TestPlanExecution_AC3 tests auto-apply mode
func TestPlanExecution_AC3(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		AutoApply:   true,
	}

	if !p.AutoApply {
		t.Error("AutoApply mode not set")
	}
}

// TestPlanExecution_AC7 tests dry-run mode
func TestPlanExecution_AC7(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		DryRun:      true,
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{
			{ID: "00-057-01", Title: "OAuth2 config", Description: "Setup", Status: "pending"},
		},
	}

	// Dry run should not create files
	err := p.CreateWorkstreamFiles(result)
	if err != nil {
		t.Fatalf("Dry run failed: %v", err)
	}

	// Verify no files created
	expectedPath := filepath.Join(backlogDir, "00-057-01-oauth2-config.md")
	if _, err := os.Stat(expectedPath); !os.IsNotExist(err) {
		t.Error("Dry run should not create files")
	}
}

// TestPlanExecution_AC8 tests error when no model configured
func TestPlanExecution_AC8(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		ModelAPI:    "", // No model configured
	}

	// Attempt decomposition without model
	result, err := p.Decompose()
	if err == nil {
		t.Error("Expected error when no model configured")
	}

	if result != nil {
		t.Error("Should return nil result when model not configured")
	}

	expectedMsg := "model API"
	if err != nil && !containsString(err.Error(), expectedMsg) {
		t.Errorf("Error message should mention model API, got: %v", err)
	}
}

// TestDecompose_WithModel tests successful decomposition with model
func TestDecompose_WithModel(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		ModelAPI:    "https://api.example.com/v1",
	}

	result, err := p.Decompose()
	if err != nil {
		t.Fatalf("Decompose failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if len(result.Workstreams) == 0 {
		t.Error("Should have workstreams")
	}

	if result.FeatureID == "" {
		t.Error("Should have FeatureID")
	}

	if result.Summary == "" {
		t.Error("Should have Summary")
	}

	if result.CreatedAt == "" {
		t.Error("Should have CreatedAt timestamp")
	}
}

// Helper function for string contains
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
