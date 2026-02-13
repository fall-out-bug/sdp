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

func TestLogger_Log_WithAllFields(t *testing.T) {
	tempDir := t.TempDir()
	logger, err := decision.NewLogger(tempDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	d := decision.Decision{
		Type:         decision.DecisionTypeTechnical,
		Question:     "Should we use PostgreSQL?",
		Decision:     "Yes, PostgreSQL for relational data",
		Rationale:    "Better ACID compliance and JSON support",
		Alternatives: []string{"MySQL", "SQLite"},
		FeatureID:    "F067",
		WorkstreamID: "00-067-15",
		Timestamp:    time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC),
	}

	if err := logger.Log(d); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(decisions) != 1 {
		t.Fatalf("Expected 1 decision, got %d", len(decisions))
	}

	// Verify all fields are preserved
	got := decisions[0]
	if got.Type != decision.DecisionTypeTechnical {
		t.Errorf("Type = %v, want %v", got.Type, decision.DecisionTypeTechnical)
	}
	if got.Question != "Should we use PostgreSQL?" {
		t.Errorf("Question = %v", got.Question)
	}
	if got.Decision != "Yes, PostgreSQL for relational data" {
		t.Errorf("Decision = %v", got.Decision)
	}
	if got.FeatureID != "F067" {
		t.Errorf("FeatureID = %v, want F067", got.FeatureID)
	}
	if got.WorkstreamID != "00-067-15" {
		t.Errorf("WorkstreamID = %v, want 00-067-15", got.WorkstreamID)
	}
	if len(got.Alternatives) != 2 {
		t.Errorf("Alternatives count = %d, want 2", len(got.Alternatives))
	}
	// Timestamp should be preserved (not overwritten)
	if !got.Timestamp.Equal(time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC)) {
		t.Errorf("Timestamp = %v, want 2026-02-13 12:00:00 UTC", got.Timestamp)
	}
}

func TestLogger_Log_MultipleDecisions(t *testing.T) {
	tempDir := t.TempDir()
	logger, err := decision.NewLogger(tempDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Log multiple decisions
	for i := 0; i < 5; i++ {
		d := decision.Decision{
			Type:     decision.DecisionTypeTechnical,
			Question: "Question",
			Decision: "Decision {i}",
		}
		if err := logger.Log(d); err != nil {
			t.Fatalf("Log %d failed: %v", i, err)
		}
	}

	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(decisions) != 5 {
		t.Errorf("Expected 5 decisions, got %d", len(decisions))
	}
}

func TestLogger_Log_DifferentTypes(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	types := []string{
		decision.DecisionTypeVision,
		decision.DecisionTypeTechnical,
		decision.DecisionTypeTradeoff,
		decision.DecisionTypeExplicit,
	}

	for _, dt := range types {
		d := decision.Decision{
			Type:     dt,
			Question: dt,
			Decision: "Decision",
		}
		if err := logger.Log(d); err != nil {
			t.Fatalf("Log failed for type %v: %v", dt, err)
		}
	}

	decisions, _ := logger.LoadAll()
	if len(decisions) != 4 {
		t.Errorf("Expected 4 decisions, got %d", len(decisions))
	}
}
