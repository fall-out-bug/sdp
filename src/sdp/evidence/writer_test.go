package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriter_New(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	writer, err := NewWriter(logPath)
	if err != nil {
		t.Fatalf("NewWriter failed: %v", err)
	}

	if writer == nil {
		t.Fatal("writer should not be nil")
	}
}

func TestWriter_WriteDecision(t *testing.T) {
	tempDir := t.TempDir()
	writer, _ := NewWriter(filepath.Join(tempDir, "events.jsonl"))

	event := NewDecisionEvent("Question?", "Answer", "")
	if err := writer.WriteDecision(&event); err != nil {
		t.Fatalf("WriteDecision failed: %v", err)
	}

	events, _ := writer.log.ReadAll()
	if len(events) != 1 {
		t.Errorf("expected 1 event, got %d", len(events))
	}
}

func TestWriter_WritePlan(t *testing.T) {
	tempDir := t.TempDir()
	writer, _ := NewWriter(filepath.Join(tempDir, "events.jsonl"))

	event := NewPlanEvent("F054", "Evidence", []string{"WS1"}, "")
	if err := writer.WritePlan(&event); err != nil {
		t.Fatalf("WritePlan failed: %v", err)
	}
}

func TestWriter_WriteGeneration(t *testing.T) {
	tempDir := t.TempDir()
	writer, _ := NewWriter(filepath.Join(tempDir, "events.jsonl"))

	event := NewGenerationEvent("claude", "hash", "")
	if err := writer.WriteGeneration(&event); err != nil {
		t.Fatalf("WriteGeneration failed: %v", err)
	}
}

func TestWriter_WriteVerification(t *testing.T) {
	tempDir := t.TempDir()
	writer, _ := NewWriter(filepath.Join(tempDir, "events.jsonl"))

	event := NewVerificationEvent("go", "go test ./...", true, "")
	if err := writer.WriteVerification(&event); err != nil {
		t.Fatalf("WriteVerification failed: %v", err)
	}
}

func TestWriter_WriteAcceptance(t *testing.T) {
	tempDir := t.TempDir()
	writer, _ := NewWriter(filepath.Join(tempDir, "events.jsonl"))

	event := NewAcceptanceEvent("npm start", true, "OK", 30, "")
	if err := writer.WriteAcceptance(&event); err != nil {
		t.Fatalf("WriteAcceptance failed: %v", err)
	}
}

func TestWriter_WriteApproval(t *testing.T) {
	tempDir := t.TempDir()
	writer, _ := NewWriter(filepath.Join(tempDir, "events.jsonl"))

	event := NewApprovalEvent("user@test.com", "human", "")
	if err := writer.WriteApproval(&event); err != nil {
		t.Fatalf("WriteApproval failed: %v", err)
	}
}

func TestWriter_Chain(t *testing.T) {
	tempDir := t.TempDir()
	writer, _ := NewWriter(filepath.Join(tempDir, "events.jsonl"))

	// Write multiple events
	plan := NewPlanEvent("F054", "Evidence", []string{"WS1"}, "")
	writer.WritePlan(&plan)

	decision := NewDecisionEvent("Q?", "A", "")
	writer.WriteDecision(&decision)

	gen := NewGenerationEvent("claude", "hash", "")
	writer.WriteGeneration(&gen)

	ver := NewVerificationEvent("go", "go test", true, "")
	writer.WriteVerification(&ver)

	approval := NewApprovalEvent("user", "human", "")
	writer.WriteApproval(&approval)

	// Verify chain
	if err := writer.log.VerifyChain(); err != nil {
		t.Errorf("chain verification failed: %v", err)
	}

	events, _ := writer.log.ReadAll()
	if len(events) != 5 {
		t.Errorf("expected 5 events, got %d", len(events))
	}
}

func TestWriter_DefaultPath(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	writer, err := NewWriter(DefaultLogPath)
	if err != nil {
		t.Fatalf("NewWriter with default path failed: %v", err)
	}

	if writer.log.Path() != DefaultLogPath {
		t.Errorf("expected path %s, got %s", DefaultLogPath, writer.log.Path())
	}
}
