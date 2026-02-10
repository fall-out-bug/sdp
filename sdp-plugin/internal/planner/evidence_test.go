package planner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

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

	if !containsString(err.Error(), "evidence writer") {
		t.Errorf("Error should mention evidence writer, got: %v", err)
	}
}
