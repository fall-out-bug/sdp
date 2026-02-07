package decision_test

import (
	"sync"
	"testing"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestLogger_Log_Concurrent(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	var wg sync.WaitGroup
	numGoroutines := 10
	decisionsPerGoroutine := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < decisionsPerGoroutine; j++ {
				d := decision.Decision{
					Question:  "Test",
					Decision:  "Yes",
					Rationale: "Because",
				}
				if err := logger.Log(d); err != nil {
					t.Errorf("Goroutine %d: Log failed: %v", id, err)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify all decisions were logged
	decisions, _ := logger.LoadAll()
	expected := numGoroutines * decisionsPerGoroutine
	if len(decisions) != expected {
		t.Errorf("Expected %d decisions, got %d", expected, len(decisions))
	}
}

func TestLogger_LogBatch_Concurrent(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	var wg sync.WaitGroup
	numGoroutines := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			batch := []decision.Decision{
				{Question: "Q1", Decision: "D1"},
				{Question: "Q2", Decision: "D2"},
			}
			if err := logger.LogBatch(batch); err != nil {
				t.Errorf("Goroutine %d: LogBatch failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	decisions, _ := logger.LoadAll()
	expected := numGoroutines * 2
	if len(decisions) != expected {
		t.Errorf("Expected %d decisions, got %d", expected, len(decisions))
	}
}
