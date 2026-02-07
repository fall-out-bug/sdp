package decision_test

import (
	"fmt"
	"testing"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestLogger_LoadAll_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(decisions) != 0 {
		t.Errorf("Expected 0 decisions, got %d", len(decisions))
	}
}

func TestLogger_LoadAll_AfterLog(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	d := decision.Decision{
		Type:      decision.DecisionTypeTechnical,
		Question:  "Q1",
		Decision:  "D1",
		Rationale: "R1",
	}
	logger.Log(d)

	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(decisions) != 1 {
		t.Fatalf("Expected 1 decision, got %d", len(decisions))
	}

	if decisions[0].Question != "Q1" {
		t.Errorf("Expected 'Q1', got '%s'", decisions[0].Question)
	}
}

func TestLogger_LoadAll_MultipleDecisions(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	for i := 1; i <= 3; i++ {
		d := decision.Decision{
			Question:  fmt.Sprintf("Q%d", i),
			Decision:  fmt.Sprintf("D%d", i),
			Rationale: fmt.Sprintf("R%d", i),
		}
		logger.Log(d)
	}

	decisions, _ := logger.LoadAll()
	if len(decisions) != 3 {
		t.Errorf("Expected 3 decisions, got %d", len(decisions))
	}
}
