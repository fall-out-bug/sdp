package planner

import (
	"os"
	"path/filepath"
	"testing"
)

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
