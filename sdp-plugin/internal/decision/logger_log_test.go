package decision_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestLogger_Log_Success(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	logger, err := decision.NewLogger(tempDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Execute
	d := decision.Decision{
		Type:      decision.DecisionTypeTechnical,
		Question:  "Test question?",
		Decision:  "Test decision",
		Rationale: "Test rationale",
	}

	if err := logger.Log(d); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Verify
	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(decisions) != 1 {
		t.Errorf("Expected 1 decision, got %d", len(decisions))
	}

	if decisions[0].Decision != "Test decision" {
		t.Errorf("Expected 'Test decision', got '%s'", decisions[0].Decision)
	}

	if decisions[0].Timestamp.IsZero() {
		t.Error("Timestamp should be auto-set")
	}
}

func TestLogger_Log_TimestampAutoSet(t *testing.T) {
	tempDir := t.TempDir()
	logger, err := decision.NewLogger(tempDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	d := decision.Decision{
		Question: "Test?",
		Decision: "Yes",
	}
	// Timestamp is zero

	before := time.Now()
	if err := logger.Log(d); err != nil {
		t.Fatalf("Log failed: %v", err)
	}
	after := time.Now()

	decisions, _ := logger.LoadAll()
	if decisions[0].Timestamp.Before(before) || decisions[0].Timestamp.After(after) {
		t.Error("Timestamp not auto-set correctly")
	}
}

func TestLogger_Log_FileCreated(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)
	d := decision.Decision{Question: "Q", Decision: "D"}

	logger.Log(d)

	logPath := filepath.Join(tempDir, "docs", "decisions", "decisions.jsonl")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file not created")
	}
}
