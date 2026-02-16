package metrics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCollector_CollectMetrics_EmptyLog_ReturnsZeroMetrics(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")
	collector := NewCollector(logPath, outputPath)

	// Act
	metrics, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if metrics == nil {
		t.Fatal("Expected metrics object, got nil")
	}
	if metrics.CatchRate != 0 {
		t.Errorf("Expected catch rate 0, got %f", metrics.CatchRate)
	}
	if metrics.TotalVerifications != 0 {
		t.Errorf("Expected total verifications 0, got %d", metrics.TotalVerifications)
	}
}

func TestCollector_CollectMetrics_WithEvents_ComputesCatchRate(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")

	// Create test evidence events
	events := []string{
		`{"id":"evt-1","type":"verification","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"passed":false}}`,
		`{"id":"evt-2","type":"verification","timestamp":"2026-02-11T10:01:00Z","ws_id":"00-001-01","data":{"passed":true}}`,
		`{"id":"evt-3","type":"verification","timestamp":"2026-02-11T10:02:00Z","ws_id":"00-001-02","data":{"passed":false}}`,
		`{"id":"evt-4","type":"verification","timestamp":"2026-02-11T10:03:00Z","ws_id":"00-001-02","data":{"passed":true}}`,
	}
	if err := os.WriteFile(logPath, []byte(string(events[0])+"\n"+string(events[1])+"\n"+string(events[2])+"\n"+string(events[3])+"\n"), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)

	// Act
	metrics, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if metrics.TotalVerifications != 4 {
		t.Errorf("Expected 4 verifications, got %d", metrics.TotalVerifications)
	}
	if metrics.FailedVerifications != 2 {
		t.Errorf("Expected 2 failed verifications, got %d", metrics.FailedVerifications)
	}
	expectedCatchRate := 0.5
	if metrics.CatchRate != expectedCatchRate {
		t.Errorf("Expected catch rate %f, got %f", expectedCatchRate, metrics.CatchRate)
	}
}

func TestCollector_CollectMetrics_ComputesIterationCount(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")

	// Create test events with generations for same workstream
	events := []string{
		`{"id":"evt-1","type":"generation","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"model_id":"claude-sonnet-4"}}`,
		`{"id":"evt-2","type":"verification","timestamp":"2026-02-11T10:01:00Z","ws_id":"00-001-01","data":{"passed":false}}`,
		`{"id":"evt-3","type":"generation","timestamp":"2026-02-11T10:02:00Z","ws_id":"00-001-01","data":{"model_id":"claude-sonnet-4"}}`,
		`{"id":"evt-4","type":"verification","timestamp":"2026-02-11T10:03:00Z","ws_id":"00-001-01","data":{"passed":true}}`,
	}
	eventsData := string(events[0]) + "\n" + string(events[1]) + "\n" + string(events[2]) + "\n" + string(events[3]) + "\n"
	if err := os.WriteFile(logPath, []byte(eventsData), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)

	// Act
	metrics, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	iterations, exists := metrics.IterationCount["00-001-01"]
	if !exists {
		t.Fatal("Expected iteration count for 00-001-01, got not found")
	}
	if iterations != 2 {
		t.Errorf("Expected 2 iterations for 00-001-01, got %d", iterations)
	}
}

func TestCollector_CollectMetrics_ComputesModelPassRate(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")

	// Create test events with generation events to establish model mapping,
	// then verification events to track pass rate
	events := []string{
		`{"id":"evt-1","type":"generation","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"model_id":"claude-sonnet-4"}}`,
		`{"id":"evt-2","type":"verification","timestamp":"2026-02-11T10:01:00Z","ws_id":"00-001-01","data":{"passed":true}}`,
		`{"id":"evt-3","type":"generation","timestamp":"2026-02-11T10:02:00Z","ws_id":"00-001-02","data":{"model_id":"claude-opus-4"}}`,
		`{"id":"evt-4","type":"verification","timestamp":"2026-02-11T10:03:00Z","ws_id":"00-001-02","data":{"passed":false}}`,
	}
	eventsData := string(events[0]) + "\n" + string(events[1]) + "\n" + string(events[2]) + "\n" + string(events[3]) + "\n"
	if err := os.WriteFile(logPath, []byte(eventsData), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)

	// Act
	metrics, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(metrics.ModelPassRate) == 0 {
		t.Fatal("Expected model pass rates, got empty map")
	}
	// claude-sonnet-4: 1/1 = 1.0 (1 passed out of 1 total)
	sonnetRate, exists := metrics.ModelPassRate["claude-sonnet-4"]
	if !exists {
		t.Fatal("Expected pass rate for claude-sonnet-4, got not found")
	}
	if sonnetRate != 1.0 {
		t.Errorf("Expected pass rate 1.0 for claude-sonnet-4, got %f", sonnetRate)
	}
	// claude-opus-4: 0/1 = 0.0 (0 passed out of 1 total)
	opusRate, exists := metrics.ModelPassRate["claude-opus-4"]
	if !exists {
		t.Fatal("Expected pass rate for claude-opus-4, got not found")
	}
	if opusRate != 0.0 {
		t.Errorf("Expected pass rate 0.0 for claude-opus-4, got %f", opusRate)
	}
}

func TestCollector_CollectMetrics_ComputesAcceptanceCatchRate(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")

	// Create test events with approval (acceptance) events
	events := []string{
		`{"id":"evt-1","type":"approval","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"approved":false}}`,
		`{"id":"evt-2","type":"approval","timestamp":"2026-02-11T10:01:00Z","ws_id":"00-001-02","data":{"approved":true}}`,
		`{"id":"evt-3","type":"approval","timestamp":"2026-02-11T10:02:00Z","ws_id":"00-001-03","data":{"approved":true}}`,
	}
	eventsData := string(events[0]) + "\n" + string(events[1]) + "\n" + string(events[2]) + "\n"
	if err := os.WriteFile(logPath, []byte(eventsData), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)

	// Act
	metrics, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if metrics.TotalApprovals != 3 {
		t.Errorf("Expected 3 total approvals, got %d", metrics.TotalApprovals)
	}
	if metrics.FailedApprovals != 1 {
		t.Errorf("Expected 1 failed approval, got %d", metrics.FailedApprovals)
	}
	// AC5: acceptance_catch_rate = acceptance failures / total builds
	expectedCatchRate := 1.0 / 3.0
	if metrics.AcceptanceCatchRate != expectedCatchRate {
		t.Errorf("Expected acceptance catch rate %f, got %f", expectedCatchRate, metrics.AcceptanceCatchRate)
	}
}

func TestCollector_CollectMetrics_WritesJSONOutput(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")

	// Create minimal test data
	events := `{"id":"evt-1","type":"verification","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"passed":true}}`
	if err := os.WriteFile(logPath, []byte(events+"\n"), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)

	// Act
	_, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Expected output file %s to exist", outputPath)
	}

	// Verify valid JSON
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var metrics Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		t.Fatalf("Output file is not valid JSON: %v", err)
	}
}

func TestCollector_CollectMetrics_IncrementalOnlyProcessesSinceWatermark(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")
	watermarkPath := filepath.Join(tempDir, "watermark.txt")

	// Create initial events
	events := `{"id":"evt-1","type":"verification","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"passed":true}}`
	if err := os.WriteFile(logPath, []byte(events+"\n"), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)
	collector.SetWatermarkPath(watermarkPath)

	// Act - First collection
	firstMetrics, err := collector.Collect()

	// Assert - First collection
	if err != nil {
		t.Fatalf("Expected no error on first collect, got %v", err)
	}
	if firstMetrics.TotalVerifications != 1 {
		t.Errorf("Expected 1 verification on first collect, got %d", firstMetrics.TotalVerifications)
	}

	// Add more events
	moreEvents := `{"id":"evt-2","type":"verification","timestamp":"2026-02-11T10:01:00Z","ws_id":"00-001-02","data":{"passed":false}}`
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	if _, err := f.Write([]byte(moreEvents + "\n")); err != nil {
		_ = f.Close()
		t.Fatalf("Failed to append events: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Failed to close log file: %v", err)
	}

	// Act - Second collection (incremental)
	collector2 := NewCollector(logPath, outputPath)
	collector2.SetWatermarkPath(watermarkPath)
	secondMetrics, err := collector2.Collect()

	// Assert - Second collection should only process new event
	if err != nil {
		t.Fatalf("Expected no error on second collect, got %v", err)
	}
	if secondMetrics.TotalVerifications != 1 {
		t.Errorf("Expected 1 new verification on incremental collect, got %d", secondMetrics.TotalVerifications)
	}
}

func TestGetLatestWatermark_NoFile_ReturnsEmpty(t *testing.T) {
	// Arrange & Act
	watermark, err := GetLatestWatermark("/nonexistent/path/watermark.json")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error for non-existent file, got %v", err)
	}
	if watermark != "" {
		t.Errorf("Expected empty watermark, got %s", watermark)
	}
}

func TestGetLatestWatermark_WithFile_ReturnsLastID(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	watermarkPath := filepath.Join(tempDir, "watermark.json")
	ids := []string{"evt-1", "evt-2", "evt-3"}
	data, _ := json.Marshal(ids)
	if err := os.WriteFile(watermarkPath, data, 0644); err != nil {
		t.Fatalf("Failed to write watermark: %v", err)
	}

	// Act
	watermark, err := GetLatestWatermark(watermarkPath)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if watermark != "evt-3" {
		t.Errorf("Expected last ID evt-3, got %s", watermark)
	}
}

func TestGetLatestWatermark_EmptyFile_ReturnsEmpty(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	watermarkPath := filepath.Join(tempDir, "watermark.json")
	data, _ := json.Marshal([]string{})
	if err := os.WriteFile(watermarkPath, data, 0644); err != nil {
		t.Fatalf("Failed to write watermark: %v", err)
	}

	// Act
	watermark, err := GetLatestWatermark(watermarkPath)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if watermark != "" {
		t.Errorf("Expected empty watermark, got %s", watermark)
	}
}

func TestParseIntFromPath_ValidInt_ReturnsInt(t *testing.T) {
	// Arrange & Act
	result := ParseIntFromPath("123")

	// Assert
	if result != 123 {
		t.Errorf("Expected 123, got %d", result)
	}
}

func TestParseIntFromPath_InvalidInt_ReturnsZero(t *testing.T) {
	// Arrange & Act
	result := ParseIntFromPath("abc")

	// Assert
	if result != 0 {
		t.Errorf("Expected 0 for invalid int, got %d", result)
	}
}

func TestCollector_CollectMetrics_WithMultipleIterations_TracksCorrectly(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")

	// Create events with multiple generation/verification cycles
	events := []string{
		`{"id":"evt-1","type":"generation","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"model_id":"claude-sonnet-4"}}`,
		`{"id":"evt-2","type":"verification","timestamp":"2026-02-11T10:01:00Z","ws_id":"00-001-01","data":{"passed":false}}`,
		`{"id":"evt-3","type":"generation","timestamp":"2026-02-11T10:02:00Z","ws_id":"00-001-01","data":{"model_id":"claude-sonnet-4"}}`,
		`{"id":"evt-4","type":"verification","timestamp":"2026-02-11T10:03:00Z","ws_id":"00-001-01","data":{"passed":true}}`,
		`{"id":"evt-5","type":"generation","timestamp":"2026-02-11T10:04:00Z","ws_id":"00-001-01","data":{"model_id":"claude-sonnet-4"}}`,
		`{"id":"evt-6","type":"verification","timestamp":"2026-02-11T10:05:00Z","ws_id":"00-001-01","data":{"passed":true}}`,
	}
	eventsData := ""
	for _, e := range events {
		eventsData += e + "\n"
	}
	if err := os.WriteFile(logPath, []byte(eventsData), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)

	// Act
	metrics, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// 3 verification events for 00-001-01
	iterations, exists := metrics.IterationCount["00-001-01"]
	if !exists {
		t.Fatal("Expected iteration count for 00-001-01, got not found")
	}
	if iterations != 3 {
		t.Errorf("Expected 3 iterations for 00-001-01, got %d", iterations)
	}
}

func TestCollector_CollectMetrics_InvalidEvents_SkippedGracefully(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")
	outputPath := filepath.Join(tempDir, "metrics.json")

	// Create events with some invalid JSON
	events := `{"id":"evt-1","type":"verification","timestamp":"2026-02-11T10:00:00Z","ws_id":"00-001-01","data":{"passed":true}}
this is not valid json
{"id":"evt-2","type":"verification","timestamp":"2026-02-11T10:01:00Z","ws_id":"00-001-02","data":{"passed":false}}
`
	if err := os.WriteFile(logPath, []byte(events), 0644); err != nil {
		t.Fatalf("Failed to write test events: %v", err)
	}

	collector := NewCollector(logPath, outputPath)

	// Act - Should not error, just skip invalid lines
	metrics, err := collector.Collect()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error despite invalid JSON, got %v", err)
	}
	if metrics.TotalVerifications != 2 {
		t.Errorf("Expected 2 valid verifications, got %d", metrics.TotalVerifications)
	}
}
