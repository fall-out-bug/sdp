package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAnalyzer_CalculateSuccessRate(t *testing.T) {
	// Create temporary telemetry file
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test events
	events := []Event{
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-2 * time.Hour),
			Data:      map[string]interface{}{"command": "test"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-2 * time.Hour),
			Data:      map[string]interface{}{"command": "test", "success": true},
		},
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-1 * time.Hour),
			Data:      map[string]interface{}{"command": "test"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-1 * time.Hour),
			Data:      map[string]interface{}{"command": "test", "success": false},
		},
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-30 * time.Minute),
			Data:      map[string]interface{}{"command": "build"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-30 * time.Minute),
			Data:      map[string]interface{}{"command": "build", "success": true},
		},
	}

	// Write events to file
	data, err := json.Marshal(events)
	if err != nil {
		t.Fatalf("Failed to marshal events: %v", err)
	}
	if err := os.WriteFile(telemetryFile, data, 0644); err != nil {
		t.Fatalf("Failed to write telemetry file: %v", err)
	}

	// Create analyzer
	analyzer, err := NewAnalyzer(telemetryFile)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Calculate success rate
	rates, err := analyzer.CalculateSuccessRate()
	if err != nil {
		t.Fatalf("CalculateSuccessRate failed: %v", err)
	}

	// Verify results
	if len(rates) != 2 {
		t.Errorf("Expected 2 command rates, got %d", len(rates))
	}

	testRate, ok := rates["test"]
	if !ok {
		t.Fatal("Expected 'test' command in rates")
	}
	if testRate != 0.5 {
		t.Errorf("Expected test success rate 0.5, got %f", testRate)
	}

	buildRate, ok := rates["build"]
	if !ok {
		t.Fatal("Expected 'build' command in rates")
	}
	if buildRate != 1.0 {
		t.Errorf("Expected build success rate 1.0, got %f", buildRate)
	}
}

func TestAnalyzer_CalculateAverageDuration(t *testing.T) {
	// Create temporary telemetry file
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test events
	events := []Event{
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-2 * time.Hour),
			Data:      map[string]interface{}{"command": "test"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-2 * time.Hour).Add(100 * time.Millisecond),
			Data:      map[string]interface{}{"command": "test", "success": true, "duration": 100},
		},
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-1 * time.Hour),
			Data:      map[string]interface{}{"command": "test"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-1 * time.Hour).Add(200 * time.Millisecond),
			Data:      map[string]interface{}{"command": "test", "success": true, "duration": 200},
		},
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-30 * time.Minute),
			Data:      map[string]interface{}{"command": "build"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-30 * time.Minute).Add(500 * time.Millisecond),
			Data:      map[string]interface{}{"command": "build", "success": true, "duration": 500},
		},
	}

	// Write events to file
	data, err := json.Marshal(events)
	if err != nil {
		t.Fatalf("Failed to marshal events: %v", err)
	}
	if err := os.WriteFile(telemetryFile, data, 0644); err != nil {
		t.Fatalf("Failed to write telemetry file: %v", err)
	}

	// Create analyzer
	analyzer, err := NewAnalyzer(telemetryFile)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Calculate average duration
	durations, err := analyzer.CalculateAverageDuration()
	if err != nil {
		t.Fatalf("CalculateAverageDuration failed: %v", err)
	}

	// Verify results
	if len(durations) != 2 {
		t.Errorf("Expected 2 command durations, got %d", len(durations))
	}

	testDuration, ok := durations["test"]
	if !ok {
		t.Fatal("Expected 'test' command in durations")
	}
	if testDuration != 150 {
		t.Errorf("Expected test average duration 150ms, got %d", testDuration)
	}

	buildDuration, ok := durations["build"]
	if !ok {
		t.Fatal("Expected 'build' command in durations")
	}
	if buildDuration != 500 {
		t.Errorf("Expected build average duration 500ms, got %d", buildDuration)
	}
}

func TestAnalyzer_TopErrorCategories(t *testing.T) {
	// Create temporary telemetry file
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test events with various errors
	events := []Event{
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-2 * time.Hour),
			Data:      map[string]interface{}{"command": "test", "success": false, "error": "file not found"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-1*time.Hour - 30*time.Minute),
			Data:      map[string]interface{}{"command": "test", "success": false, "error": "file not found"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-1 * time.Hour),
			Data:      map[string]interface{}{"command": "test", "success": false, "error": "permission denied"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-30 * time.Minute),
			Data:      map[string]interface{}{"command": "build", "success": false, "error": "syntax error"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-15 * time.Minute),
			Data:      map[string]interface{}{"command": "build", "success": false, "error": "timeout"},
		},
	}

	// Write events to file
	data, err := json.Marshal(events)
	if err != nil {
		t.Fatalf("Failed to marshal events: %v", err)
	}
	if err := os.WriteFile(telemetryFile, data, 0644); err != nil {
		t.Fatalf("Failed to write telemetry file: %v", err)
	}

	// Create analyzer
	analyzer, err := NewAnalyzer(telemetryFile)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Get top errors
	errors, err := analyzer.TopErrorCategories(5)
	if err != nil {
		t.Fatalf("TopErrorCategories failed: %v", err)
	}

	// Verify results
	if len(errors) != 4 {
		t.Errorf("Expected 4 error categories, got %d", len(errors))
	}

	// Top error should be "file not found" with count 2
	if errors[0].Message != "file not found" {
		t.Errorf("Expected top error 'file not found', got '%s'", errors[0].Message)
	}
	if errors[0].Count != 2 {
		t.Errorf("Expected top error count 2, got %d", errors[0].Count)
	}
}

func TestAnalyzer_GenerateReport(t *testing.T) {
	// Create temporary telemetry file
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test events
	now := time.Now()
	events := []Event{
		{
			Type:      EventTypeCommandStart,
			Timestamp: now.Add(-2 * time.Hour),
			Data:      map[string]interface{}{"command": "test"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: now.Add(-2 * time.Hour).Add(100 * time.Millisecond),
			Data:      map[string]interface{}{"command": "test", "success": true, "duration": 100},
		},
		{
			Type:      EventTypeCommandStart,
			Timestamp: now.Add(-1 * time.Hour),
			Data:      map[string]interface{}{"command": "test"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: now.Add(-1 * time.Hour).Add(200 * time.Millisecond),
			Data:      map[string]interface{}{"command": "test", "success": false, "error": "timeout", "duration": 200},
		},
		{
			Type:      EventTypeCommandStart,
			Timestamp: now.Add(-30 * time.Minute),
			Data:      map[string]interface{}{"command": "build"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: now.Add(-30 * time.Minute).Add(500 * time.Millisecond),
			Data:      map[string]interface{}{"command": "build", "success": true, "duration": 500},
		},
	}

	// Write events to file
	data, err := json.Marshal(events)
	if err != nil {
		t.Fatalf("Failed to marshal events: %v", err)
	}
	if err := os.WriteFile(telemetryFile, data, 0644); err != nil {
		t.Fatalf("Failed to write telemetry file: %v", err)
	}

	// Create analyzer
	analyzer, err := NewAnalyzer(telemetryFile)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Generate report
	report, err := analyzer.GenerateReport(nil, nil)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	// Verify report structure
	if report.TotalEvents != 6 {
		t.Errorf("Expected 6 total events, got %d", report.TotalEvents)
	}

	if len(report.CommandStats) != 2 {
		t.Errorf("Expected 2 command stats, got %d", len(report.CommandStats))
	}

	// Verify test command stats
	testStats, ok := report.CommandStats["test"]
	if !ok {
		t.Fatal("Expected 'test' command in stats")
	}
	if testStats.TotalRuns != 2 {
		t.Errorf("Expected test total runs 2, got %d", testStats.TotalRuns)
	}
	if testStats.SuccessRate != 0.5 {
		t.Errorf("Expected test success rate 0.5, got %f", testStats.SuccessRate)
	}
	if testStats.AvgDuration != 150 {
		t.Errorf("Expected test avg duration 150ms, got %d", testStats.AvgDuration)
	}

	// Verify error counts
	if len(report.TopErrors) != 1 {
		t.Errorf("Expected 1 top error, got %d", len(report.TopErrors))
	}
	if report.TopErrors[0].Message != "timeout" {
		t.Errorf("Expected error 'timeout', got '%s'", report.TopErrors[0].Message)
	}
}

func TestAnalyzer_FilterByDateRange(t *testing.T) {
	// Create temporary telemetry file
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test events spanning multiple days
	now := time.Now()
	events := []Event{
		{
			Type:      EventTypeCommandComplete,
			Timestamp: now.Add(-72 * time.Hour), // 3 days ago
			Data:      map[string]interface{}{"command": "test", "success": true},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: now.Add(-48 * time.Hour), // 2 days ago
			Data:      map[string]interface{}{"command": "test", "success": true},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: now.Add(-24 * time.Hour), // 1 day ago
			Data:      map[string]interface{}{"command": "test", "success": true},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: now.Add(-12 * time.Hour), // 12 hours ago
			Data:      map[string]interface{}{"command": "test", "success": true},
		},
	}

	// Write events to file
	data, err := json.Marshal(events)
	if err != nil {
		t.Fatalf("Failed to marshal events: %v", err)
	}
	if err := os.WriteFile(telemetryFile, data, 0644); err != nil {
		t.Fatalf("Failed to write telemetry file: %v", err)
	}

	// Create analyzer
	analyzer, err := NewAnalyzer(telemetryFile)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Filter to last 48 hours
	startTime := now.Add(-48 * time.Hour)
	endTime := now

	report, err := analyzer.GenerateReport(&startTime, &endTime)
	if err != nil {
		t.Fatalf("GenerateReport with date range failed: %v", err)
	}

	// Should only include events from last 48 hours (48h, 24h, 12h)
	if report.TotalEvents != 3 {
		t.Errorf("Expected 3 events in date range, got %d", report.TotalEvents)
	}
}

func TestAnalyzer_EmptyTelemetryFile(t *testing.T) {
	// Create temporary telemetry file
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create empty file
	if err := os.WriteFile(telemetryFile, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create telemetry file: %v", err)
	}

	// Create analyzer
	analyzer, err := NewAnalyzer(telemetryFile)
	if err != nil {
		t.Fatalf("Failed to create analyzer: %v", err)
	}

	// Generate report
	report, err := analyzer.GenerateReport(nil, nil)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	// Should return empty report
	if report.TotalEvents != 0 {
		t.Errorf("Expected 0 total events, got %d", report.TotalEvents)
	}
	if len(report.CommandStats) != 0 {
		t.Errorf("Expected 0 command stats, got %d", len(report.CommandStats))
	}
}
