package telemetry

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEventTypes(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		wantValid bool
	}{
		{"command start", EventTypeCommandStart, true},
		{"command complete", EventTypeCommandComplete, true},
		{"ws start", EventTypeWSStart, true},
		{"ws complete", EventTypeWSComplete, true},
		{"quality gate result", EventTypeQualityGateResult, true},
		{"invalid type", EventType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{
				Type:      tt.eventType,
				Timestamp: time.Now(), //nolint:unusedwrite // Test fixture
			}

			isValid := event.Type.IsValid()

			if tt.wantValid && !isValid {
				t.Errorf("Expected event type %s to be valid", tt.eventType)
			}
			if !tt.wantValid && isValid {
				t.Errorf("Expected event type %s to be invalid", tt.eventType)
			}
		})
	}
}

func TestCollectorRecord(t *testing.T) {
	// Create temp directory for testing
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Test recording a command start event
	event := Event{
		Type:      EventTypeCommandStart,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"command": "init",
			"args":    []string{},
		},
	}

	err = collector.Record(event)
	if err != nil {
		t.Fatalf("Failed to record event: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(telemetryFile); os.IsNotExist(err) {
		t.Error("Telemetry file was not created")
	}
}

func TestCollectorOptOut(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create collector with telemetry disabled
	collector, err := NewCollector(telemetryFile, false)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	event := Event{
		Type:      EventTypeCommandStart,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = collector.Record(event)
	if err != nil {
		t.Fatalf("Failed to record event: %v", err)
	}

	// Verify file was NOT created when telemetry is disabled
	if _, err := os.Stat(telemetryFile); !os.IsNotExist(err) {
		t.Error("Telemetry file should not be created when disabled")
	}
}

func TestCollectorExportJSON(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")
	exportFile := filepath.Join(tmpDir, "export.json")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Record some test events
	events := []Event{
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-time.Hour),
			Data:      map[string]interface{}{"command": "init"},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"command": "init", "success": true},
		},
	}

	for _, event := range events {
		if err := collector.Record(event); err != nil {
			t.Fatalf("Failed to record event: %v", err)
		}
	}

	// Export to JSON
	err = collector.ExportJSON(exportFile)
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}

	// Verify export file exists
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		t.Error("Export file was not created")
	}
}

func TestCollectorExportCSV(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")
	exportFile := filepath.Join(tmpDir, "export.csv")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Record test event
	event := Event{
		Type:      EventTypeCommandStart,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"command": "doctor"},
	}

	if err := collector.Record(event); err != nil {
		t.Fatalf("Failed to record event: %v", err)
	}

	// Export to CSV
	err = collector.ExportCSV(exportFile)
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}

	// Verify export file exists
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		t.Error("Export file was not created")
	}
}

func TestCollectorStatus(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Test enabled status
	collectorEnabled, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	status := collectorEnabled.Status()
	if !status.Enabled {
		t.Error("Expected status to show telemetry is enabled")
	}

	// Test disabled status
	collectorDisabled, err := NewCollector(telemetryFile, false)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	status = collectorDisabled.Status()
	if status.Enabled {
		t.Error("Expected status to show telemetry is disabled")
	}
}

func TestCollectorEventCount(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Record multiple events
	for i := 0; i < 5; i++ {
		event := Event{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"count": i},
		}
		if err := collector.Record(event); err != nil {
			t.Fatalf("Failed to record event: %v", err)
		}
	}

	status := collector.Status()
	if status.EventCount != 5 {
		t.Errorf("Expected event count 5, got %d", status.EventCount)
	}
}

func TestCollectorDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Disable telemetry
	collector.Disable()

	// Try to record event
	event := Event{
		Type:      EventTypeCommandStart,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = collector.Record(event)
	if err != nil {
		t.Fatalf("Failed to record event: %v", err)
	}

	// Verify status is disabled
	status := collector.Status()
	if status.Enabled {
		t.Error("Expected telemetry to be disabled after calling Disable()")
	}
}

func TestCollectorEnable(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, false)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Verify initially disabled
	status := collector.Status()
	if status.Enabled {
		t.Error("Expected telemetry to be initially disabled")
	}

	// Enable telemetry
	collector.Enable()

	// Verify status is enabled
	status = collector.Status()
	if !status.Enabled {
		t.Error("Expected telemetry to be enabled after calling Enable()")
	}

	// Record event to verify it works
	event := Event{
		Type:      EventTypeCommandStart,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = collector.Record(event)
	if err != nil {
		t.Fatalf("Failed to record event: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(telemetryFile); os.IsNotExist(err) {
		t.Error("Telemetry file should be created after enabling and recording")
	}
}

func TestNewCollectorInvalidPath(t *testing.T) {
	// Test with empty path
	_, err := NewCollector("", true)
	if err == nil {
		t.Error("Expected error when creating collector with empty path")
	}
}

func TestCollectorRecordInvalidEvent(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Try to record invalid event
	event := Event{
		Type:      EventType("invalid_type"),
		Timestamp: time.Now(),
		Data:      map[string]interface{}{},
	}

	err = collector.Record(event)
	if err == nil {
		t.Error("Expected error when recording invalid event type")
	}
}

func TestCollectorExportEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")
	exportJSON := filepath.Join(tmpDir, "export.json")
	exportCSV := filepath.Join(tmpDir, "export.csv")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Export without any events
	err = collector.ExportJSON(exportJSON)
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}

	// Verify export file exists
	if _, err := os.Stat(exportJSON); os.IsNotExist(err) {
		t.Error("Export JSON file should be created even with no events")
	}

	// Export to CSV
	err = collector.ExportCSV(exportCSV)
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}

	// Verify export file exists
	if _, err := os.Stat(exportCSV); os.IsNotExist(err) {
		t.Error("Export CSV file should be created even with no events")
	}
}

func TestCollectorReadExistingEvents(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Record some events
	events := []Event{
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"command": "test1"},
		},
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"command": "test2"},
		},
	}

	for _, event := range events {
		if err := collector.Record(event); err != nil {
			t.Fatalf("Failed to record event: %v", err)
		}
	}

	// Verify event count in original collector
	status := collector.Status()
	if status.EventCount != 2 {
		t.Errorf("Expected event count 2, got %d", status.EventCount)
	}

	// Verify events were persisted to file
	data, err := os.ReadFile(telemetryFile)
	if err != nil {
		t.Fatalf("Failed to read telemetry file: %v", err)
	}

	// Count lines in file (should be 2 events)
	lineCount := 0
	for _, b := range data {
		if b == '\n' {
			lineCount++
		}
	}

	if lineCount != 2 {
		t.Errorf("Expected 2 lines in telemetry file, got %d", lineCount)
	}
}
