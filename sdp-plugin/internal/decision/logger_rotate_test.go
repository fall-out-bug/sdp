package decision

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLogger_Rotate_FileSizeExceedsLimit(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	filePath := filepath.Join(decisionsDir, "decisions.jsonl")

	// Create a file larger than MaxFileSize (10MB)
	// Write a large JSON object to exceed the limit
	largeData := strings.Repeat("x", MaxFileSize+1000)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Write large data
	_, err = file.WriteString(largeData + "\n")
	if err != nil {
		t.Fatalf("Failed to write large data: %v", err)
	}
	file.Close()

	// Trigger rotation
	err = logger.rotate()
	if err != nil {
		t.Errorf("rotate failed: %v", err)
	}

	// Check that original file no longer exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Original file should have been rotated")
	}

	// Check that rotated file exists
	matches, err := filepath.Glob(filePath + ".*")
	if err != nil {
		t.Fatalf("Failed to glob rotated files: %v", err)
	}

	if len(matches) != 1 {
		t.Errorf("Expected 1 rotated file, got %d", len(matches))
	}
}

func TestLogger_Rotate_FileSizeUnderLimit(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	filePath := filepath.Join(decisionsDir, "decisions.jsonl")

	// Write a small file (under 10MB)
	smallData := strings.Repeat("x", 1000)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(smallData + "\n")
	if err != nil {
		t.Fatalf("Failed to write small data: %v", err)
	}
	file.Close()

	// Trigger rotation (should not rotate)
	err = logger.rotate()
	if err != nil {
		t.Errorf("rotate failed: %v", err)
	}

	// Check that original file still exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Original file should not have been rotated")
	}

	// Check that no rotated files exist
	matches, err := filepath.Glob(filePath + ".*")
	if err != nil {
		t.Fatalf("Failed to glob rotated files: %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("Expected 0 rotated files, got %d", len(matches))
	}
}

func TestLogger_Rotate_NoFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Don't create a file, just try to rotate
	err = logger.rotate()
	if err != nil {
		t.Errorf("rotate failed: %v", err)
	}

	// Should not error, just return successfully
}

func TestLogger_Rotate_MultipleRotations(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	filePath := filepath.Join(decisionsDir, "decisions.jsonl")

	// Function to fill file
	fillFile := func() {
		largeData := strings.Repeat("x", MaxFileSize+1000)
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString(largeData + "\n")
		if err != nil {
			t.Fatalf("Failed to write large data: %v", err)
		}
		file.Close()
	}

	// First rotation
	fillFile()
	err = logger.rotate()
	if err != nil {
		t.Errorf("First rotate failed: %v", err)
	}

	// Wait a moment to ensure timestamp differs
	time.Sleep(time.Second)

	// Second rotation
	fillFile()
	err = logger.rotate()
	if err != nil {
		t.Errorf("Second rotate failed: %v", err)
	}

	// Check that we have at least 1 rotated file
	matches, err := filepath.Glob(filePath + ".*")
	if err != nil {
		t.Fatalf("Failed to glob rotated files: %v", err)
	}

	if len(matches) < 1 {
		t.Errorf("Expected at least 1 rotated file, got %d", len(matches))
	}
}

func TestLogger_Log_TimestampAutoSet(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Log decision without timestamp
	decision := Decision{
		Type:          DecisionTypeTechnical,
		FeatureID:     "F001",
		Question:      "Question",
		Decision:      "Decision",
		Rationale:     "Rationale",
		DecisionMaker: "user",
	}

	err = logger.Log(decision)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Load decision back
	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(decisions) != 1 {
		t.Fatalf("Expected 1 decision, got %d", len(decisions))
	}

	// Check timestamp was set
	if decisions[0].Timestamp.IsZero() {
		t.Error("Timestamp should have been auto-set")
	}

	// Check timestamp is recent (within last minute)
	if time.Since(decisions[0].Timestamp) > time.Minute {
		t.Error("Timestamp should be recent")
	}
}

func TestLogger_Log_WithExistingTimestamp(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Log decision with timestamp
	pastTime := time.Now().Add(-24 * time.Hour)
	decision := Decision{
		Timestamp:     pastTime,
		Type:          DecisionTypeTechnical,
		FeatureID:     "F001",
		Question:      "Question",
		Decision:      "Decision",
		Rationale:     "Rationale",
		DecisionMaker: "user",
	}

	err = logger.Log(decision)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Load decision back
	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(decisions) != 1 {
		t.Fatalf("Expected 1 decision, got %d", len(decisions))
	}

	// Check timestamp was NOT overwritten
	if !decisions[0].Timestamp.Equal(pastTime) {
		t.Errorf("Timestamp should not have been overwritten, expected %v, got %v", pastTime, decisions[0].Timestamp)
	}
}

func TestLogger_JsonSerialization(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	decision := Decision{
		Timestamp:     time.Now(),
		Type:          DecisionTypeTechnical,
		FeatureID:     "F001",
		WorkstreamID:  "00-001-01",
		Question:      "Should we use Go?",
		Decision:      "Yes, use Go",
		Rationale:     "Go provides good performance and concurrency",
		Alternatives:  []string{"Python", "Rust"},
		Outcome:       "Chose Go for backend",
		DecisionMaker: "user",
		Tags:          []string{"language", "backend"},
	}

	err = logger.Log(decision)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	// Read the file directly and verify JSON
	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	filePath := filepath.Join(decisionsDir, "decisions.jsonl")

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Parse the JSON line
	var loaded Decision
	err = json.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify all fields
	if loaded.Type != decision.Type {
		t.Errorf("Type mismatch: expected %s, got %s", decision.Type, loaded.Type)
	}

	if loaded.FeatureID != decision.FeatureID {
		t.Errorf("FeatureID mismatch: expected %s, got %s", decision.FeatureID, loaded.FeatureID)
	}

	if len(loaded.Alternatives) != 2 {
		t.Errorf("Expected 2 alternatives, got %d", len(loaded.Alternatives))
	}

	if len(loaded.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(loaded.Tags))
	}
}
