package decision_test

import (
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestLogger_Log_WithPresetTimestamp(t *testing.T) {
	tempDir := t.TempDir()
	logger, err := decision.NewLogger(tempDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Set a specific timestamp
	presetTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	d := decision.Decision{
		Question:  "Test?",
		Decision:  "Yes",
		Timestamp: presetTime,
	}

	if err := logger.Log(d); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	decisions, _ := logger.LoadAll()
	// Timestamp should be preserved
	if len(decisions) > 0 && !decisions[0].Timestamp.Equal(presetTime) {
		t.Errorf("Timestamp = %v, want %v", decisions[0].Timestamp, presetTime)
	}
}

func TestLogger_LogBatch_Single(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	decisions := []decision.Decision{
		{Question: "Q1", Decision: "D1"},
	}

	if err := logger.LogBatch(decisions); err != nil {
		t.Fatalf("LogBatch failed: %v", err)
	}

	loaded, _ := logger.LoadAll()
	if len(loaded) != 1 {
		t.Errorf("Expected 1 decision, got %d", len(loaded))
	}
}

func TestLogger_LogBatch_WithAllFields(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	decisions := []decision.Decision{
		{
			Type:         decision.DecisionTypeTechnical,
			Question:     "Q1",
			Decision:     "D1",
			Rationale:    "R1",
			Alternatives: []string{"A1", "A2"},
			FeatureID:    "F001",
			WorkstreamID: "00-001-01",
		},
	}

	if err := logger.LogBatch(decisions); err != nil {
		t.Fatalf("LogBatch failed: %v", err)
	}

	loaded, _ := logger.LoadAll()
	if len(loaded) > 0 && loaded[0].FeatureID != "F001" {
		t.Errorf("FeatureID = %s, want F001", loaded[0].FeatureID)
	}
}

func TestLogger_Log_AfterLogBatch(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	// Log batch first
	logger.LogBatch([]decision.Decision{
		{Question: "Q1", Decision: "D1"},
	})

	// Then log single
	logger.Log(decision.Decision{Question: "Q2", Decision: "D2"})

	loaded, _ := logger.LoadAll()
	if len(loaded) != 2 {
		t.Errorf("Expected 2 decisions, got %d", len(loaded))
	}
}

func TestLogger_LogBatch_AppendsToExisting(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	// Log initial
	logger.Log(decision.Decision{Question: "Q1", Decision: "D1"})

	// Log batch - should append
	logger.LogBatch([]decision.Decision{
		{Question: "Q2", Decision: "D2"},
		{Question: "Q3", Decision: "D3"},
	})

	loaded, _ := logger.LoadAll()
	if len(loaded) != 3 {
		t.Errorf("Expected 3 decisions, got %d", len(loaded))
	}
}

func TestLogger_Log_ConcurrentSafe(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	done := make(chan bool)

	// Log from multiple goroutines
	for i := 0; i < 10; i++ {
		go func(id int) {
			d := decision.Decision{
				Question: "Concurrent question",
				Decision: "Concurrent decision",
			}
			logger.Log(d)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	loaded, _ := logger.LoadAll()
	if len(loaded) != 10 {
		t.Errorf("Expected 10 decisions, got %d", len(loaded))
	}
}

func TestLogger_Log_AllFields(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	// Log decision with all fields populated
	presetTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	d := decision.Decision{
		Type:          decision.DecisionTypeTechnical,
		Question:      "Should we use Go?",
		Decision:      "Yes, use Go",
		Rationale:     "Go provides good performance",
		Alternatives:  []string{"Python", "Rust"},
		Outcome:       "Chose Go for backend",
		DecisionMaker: "team",
		FeatureID:     "F001",
		WorkstreamID:  "00-001-01",
		Tags:          []string{"language", "backend"},
		Timestamp:     presetTime,
	}

	if err := logger.Log(d); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	loaded, _ := logger.LoadAll()
	if len(loaded) != 1 {
		t.Fatalf("Expected 1 decision, got %d", len(loaded))
	}

	// Verify all fields were persisted
	if loaded[0].Type != decision.DecisionTypeTechnical {
		t.Errorf("Type = %s", loaded[0].Type)
	}
	if loaded[0].FeatureID != "F001" {
		t.Errorf("FeatureID = %s", loaded[0].FeatureID)
	}
	if len(loaded[0].Alternatives) != 2 {
		t.Errorf("Alternatives count = %d", len(loaded[0].Alternatives))
	}
	if len(loaded[0].Tags) != 2 {
		t.Errorf("Tags count = %d", len(loaded[0].Tags))
	}
}

func TestLogger_LogBatch_MultipleWithAllFields(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	decisions := []decision.Decision{
		{
			Type:         decision.DecisionTypeTechnical,
			Question:     "Q1",
			Decision:     "D1",
			Rationale:    "R1",
			Alternatives: []string{"A1"},
			FeatureID:    "F001",
			WorkstreamID: "00-001-01",
			Tags:         []string{"tag1"},
		},
		{
			Type:         decision.DecisionTypeVision,
			Question:     "Q2",
			Decision:     "D2",
			Rationale:    "R2",
			Alternatives: []string{"A2", "A3"},
			FeatureID:    "F002",
			WorkstreamID: "00-002-01",
			Tags:         []string{"tag2", "tag3"},
		},
	}

	if err := logger.LogBatch(decisions); err != nil {
		t.Fatalf("LogBatch failed: %v", err)
	}

	loaded, _ := logger.LoadAll()
	if len(loaded) != 2 {
		t.Errorf("Expected 2 decisions, got %d", len(loaded))
	}
}
