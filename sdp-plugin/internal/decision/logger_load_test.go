package decision

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogger_Load(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Add some test decisions
	decisions := []Decision{
		{
			Timestamp:     time.Now(),
			Type:          DecisionTypeTechnical,
			FeatureID:     "F001",
			WorkstreamID:  "00-001-01",
			Question:      "Question 1",
			Decision:      "Decision 1",
			Rationale:     "Rationale 1",
			DecisionMaker: "user",
		},
		{
			Timestamp:     time.Now(),
			Type:          DecisionTypeVision,
			FeatureID:     "F002",
			WorkstreamID:  "00-002-01",
			Question:      "Question 2",
			Decision:      "Decision 2",
			Rationale:     "Rationale 2",
			DecisionMaker: "claude",
		},
		{
			Timestamp:     time.Now(),
			Type:          DecisionTypeTradeoff,
			FeatureID:     "F003",
			WorkstreamID:  "00-003-01",
			Question:      "Question 3",
			Decision:      "Decision 3",
			Rationale:     "Rationale 3",
			DecisionMaker: "user",
		},
	}

	err = logger.LogBatch(decisions)
	if err != nil {
		t.Fatalf("LogBatch failed: %v", err)
	}

	tests := []struct {
		name         string
		opts         LoadOptions
		expectedLen  int
		expectedErr  bool
		firstID      string
	}{
		{
			name:        "Load all decisions",
			opts:        LoadOptions{Offset: 0, Limit: 0},
			expectedLen: 3,
			expectedErr: false,
		},
		{
			name:        "Load first decision only",
			opts:        LoadOptions{Offset: 0, Limit: 1},
			expectedLen: 1,
			expectedErr: false,
		},
		{
			name:        "Load with offset",
			opts:        LoadOptions{Offset: 1, Limit: 0},
			expectedLen: 2,
			expectedErr: false,
		},
		{
			name:        "Load with offset and limit",
			opts:        LoadOptions{Offset: 1, Limit: 1},
			expectedLen: 1,
			expectedErr: false,
		},
		{
			name:        "Load with offset exceeding total",
			opts:        LoadOptions{Offset: 10, Limit: 1},
			expectedLen: 0,
			expectedErr: false,
		},
		{
			name:        "Load with limit larger than available",
			opts:        LoadOptions{Offset: 0, Limit: 10},
			expectedLen: 3,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := logger.Load(tt.opts)

			if tt.expectedErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(result) != tt.expectedLen {
				t.Errorf("Expected %d decisions, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestLogger_Load_EmptyLog(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Don't log anything, try to load
	result, err := logger.Load(LoadOptions{Offset: 0, Limit: 10})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 decisions, got %d", len(result))
	}
}

func TestLogger_Load_NonExistentFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Delete the decisions file
	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	filePath := filepath.Join(decisionsDir, "decisions.jsonl")
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove file: %v", err)
	}

	// Try to load from non-existent file
	result, err := logger.Load(LoadOptions{Offset: 0, Limit: 10})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 decisions, got %d", len(result))
	}
}

func TestLogger_Load_CorruptedFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Write a valid decision
	decision := Decision{
		Timestamp:     time.Now(),
		Type:          DecisionTypeTechnical,
		FeatureID:     "F001",
		Question:      "Question 1",
		Decision:      "Decision 1",
		DecisionMaker: "user",
	}

	err = logger.Log(decision)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Append corrupted data
	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	filePath := filepath.Join(decisionsDir, "decisions.jsonl")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("invalid json\n")
	if err != nil {
		t.Fatalf("Failed to write corrupted data: %v", err)
	}

	// Load should return valid decision and stop at corruption
	result, err := logger.Load(LoadOptions{Offset: 0, Limit: 10})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should return 1 valid decision before hitting corruption
	if len(result) != 1 {
		t.Errorf("Expected 1 decision, got %d", len(result))
	}
}

func TestLoadOptions_DefaultValues(t *testing.T) {
	opts := LoadOptions{}

	if opts.Offset != 0 {
		t.Errorf("Expected default Offset 0, got %d", opts.Offset)
	}

	if opts.Limit != 0 {
		t.Errorf("Expected default Limit 0, got %d", opts.Limit)
	}
}
