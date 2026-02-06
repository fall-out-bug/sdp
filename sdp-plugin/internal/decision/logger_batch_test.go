package decision_test

import (
	"fmt"
	"testing"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestLogger_LogBatch_Empty(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	err := logger.LogBatch([]decision.Decision{})
	if err != nil {
		t.Fatalf("LogBatch with empty slice failed: %v", err)
	}

	decisions, _ := logger.LoadAll()
	if len(decisions) != 0 {
		t.Errorf("Expected 0 decisions, got %d", len(decisions))
	}
}

func TestLogger_LogBatch_Multiple(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	decisions := []decision.Decision{
		{Question: "Q1", Decision: "D1", Rationale: "R1"},
		{Question: "Q2", Decision: "D2", Rationale: "R2"},
		{Question: "Q3", Decision: "D3", Rationale: "R3"},
	}

	if err := logger.LogBatch(decisions); err != nil {
		t.Fatalf("LogBatch failed: %v", err)
	}

	loaded, _ := logger.LoadAll()
	if len(loaded) != 3 {
		t.Errorf("Expected 3 decisions, got %d", len(loaded))
	}

	for i, d := range loaded {
		expectedQ := fmt.Sprintf("Q%d", i+1)
		if d.Question != expectedQ {
			t.Errorf("Decision %d: expected '%s', got '%s'", i, expectedQ, d.Question)
		}
	}
}

func TestLogger_LogBatch_ConcurrentWithLog(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	// Mix Log and LogBatch
	logger.Log(decision.Decision{Question: "Single", Decision: "S"})

	batch := []decision.Decision{
		{Question: "B1", Decision: "B1"},
		{Question: "B2", Decision: "B2"},
	}
	logger.LogBatch(batch)

	logger.Log(decision.Decision{Question: "Single2", Decision: "S2"})

	loaded, _ := logger.LoadAll()
	if len(loaded) != 4 {
		t.Errorf("Expected 4 decisions, got %d", len(loaded))
	}
}
