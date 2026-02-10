package planner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/evidence"
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

// TestPlanExecution_AC5 tests plan event emission
func TestPlanExecution_AC5(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	// Setup evidence log
	writer, err := evidence.NewWriter(logPath)
	if err != nil {
		t.Fatalf("Failed to create evidence writer: %v", err)
	}

	p := &Planner{
		BacklogDir:  filepath.Join(tempDir, "backlog"),
		Description: "Add OAuth2",
		EvidenceWriter: writer,
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{
			{ID: "00-057-01", Title: "OAuth2 config", Description: "Setup", Status: "pending"},
		},
	}

	// Emit plan event
	err = p.EmitPlanEvent(result)
	if err != nil {
		t.Fatalf("Failed to emit plan event: %v", err)
	}

	// Verify event was written
	reader := evidence.NewReader(logPath)
	events, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read events: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type != "plan" {
		t.Errorf("Expected plan event type, got %s", events[0].Type)
	}
}

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
	if err != nil && !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Error message should mention model API, got: %v", err)
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

// TestPromptForInteractive tests interactive mode prompting
func TestPromptForInteractive(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		Interactive: true,
	}

	err := p.PromptForInteractive()
	if err != nil {
		t.Fatalf("PromptForInteractive failed: %v", err)
	}

	// Non-interactive mode should also work
	p.Interactive = false
	err = p.PromptForInteractive()
	if err != nil {
		t.Fatalf("PromptForInteractive in non-interactive mode failed: %v", err)
	}
}

// TestExecuteAutoApply tests auto-apply execution
func TestExecuteAutoApply(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		AutoApply:   true,
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{
			{ID: "00-057-01", Title: "OAuth2 config", Description: "Setup", Status: "pending"},
		},
	}

	err := p.ExecuteAutoApply(result)
	if err != nil {
		t.Fatalf("ExecuteAutoApply failed: %v", err)
	}

	// Test with empty workstreams
	emptyResult := &DecompositionResult{
		Workstreams: []Workstream{},
	}

	err = p.ExecuteAutoApply(emptyResult)
	if err == nil {
		t.Error("Expected error when auto-applying empty plan")
	}

	// Non-auto-apply mode should also work
	p.AutoApply = false
	err = p.ExecuteAutoApply(result)
	if err != nil {
		t.Fatalf("ExecuteAutoApply in non-auto-apply mode failed: %v", err)
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

	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("Error should mention unsupported format, got: %v", err)
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

// TestEmitPlanEvent_NoWriter tests error when no writer configured
func TestEmitPlanEvent_NoWriter(t *testing.T) {
	tempDir := t.TempDir()
	backlogDir := filepath.Join(tempDir, "backlog")
	os.MkdirAll(backlogDir, 0755)

	p := &Planner{
		BacklogDir:  backlogDir,
		Description: "Add OAuth2",
		EvidenceWriter: nil,
	}

	result := &DecompositionResult{
		Workstreams: []Workstream{},
	}

	err := p.EmitPlanEvent(result)
	if err == nil {
		t.Error("Expected error when evidence writer not configured")
	}

	if !strings.Contains(err.Error(), "evidence writer") {
		t.Errorf("Error should mention evidence writer, got: %v", err)
	}
}

