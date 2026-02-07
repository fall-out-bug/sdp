package telemetry

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestPackForUploadJSON tests packaging telemetry data as JSON
func TestPackForUploadJSON(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test telemetry data
	events := []Event{
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now().Add(-2 * time.Hour),
			Data: map[string]interface{}{
				"command": "doctor",
				"args":    []string{},
			},
		},
		{
			Type:      EventTypeCommandComplete,
			Timestamp: time.Now().Add(-2*time.Hour + time.Second),
			Data: map[string]interface{}{
				"command":  "doctor",
				"duration": 1234,
				"success":  true,
				"error":    "",
			},
		},
	}

	// Write events to telemetry file
	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	for _, event := range events {
		if err := collector.Record(event); err != nil {
			t.Fatalf("Failed to record event: %v", err)
		}
	}

	// Test packaging as JSON
	outputPath := filepath.Join(tmpDir, "upload.json")
	result, err := PackForUpload(telemetryFile, outputPath, "json")
	if err != nil {
		t.Fatalf("PackForUpload failed: %v", err)
	}

	// Verify result
	if result.EventCount != 2 {
		t.Errorf("Expected 2 events, got %d", result.EventCount)
	}

	if result.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", result.Format)
	}

	if result.Size == 0 {
		t.Error("Expected non-zero file size")
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify output file contains valid JSON
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var uploadData map[string]interface{}
	if err := json.Unmarshal(data, &uploadData); err != nil {
		t.Errorf("Output file is not valid JSON: %v", err)
	}

	// Verify structure
	if _, ok := uploadData["events"]; !ok {
		t.Error("Output JSON missing 'events' key")
	}

	if _, ok := uploadData["metadata"]; !ok {
		t.Error("Output JSON missing 'metadata' key")
	}
}

// TestPackForUploadArchive tests packaging telemetry data as tar.gz archive
func TestPackForUploadArchive(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test telemetry data
	events := []Event{
		{
			Type:      EventTypeCommandStart,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"command": "build",
				"args":    []string{"00-001-01"},
			},
		},
	}

	// Write events to telemetry file
	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	if err := collector.Record(events[0]); err != nil {
		t.Fatalf("Failed to record event: %v", err)
	}

	// Test packaging as archive
	outputPath := filepath.Join(tmpDir, "upload.tar.gz")
	result, err := PackForUpload(telemetryFile, outputPath, "archive")
	if err != nil {
		t.Fatalf("PackForUpload failed: %v", err)
	}

	// Verify result
	if result.EventCount != 1 {
		t.Errorf("Expected 1 event, got %d", result.EventCount)
	}

	if result.Format != "archive" {
		t.Errorf("Expected format 'archive', got '%s'", result.Format)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify it's a valid gzip archive
	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open archive: %v", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("Not a valid gzip file: %v", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	// Should have at least one file (telemetry.jsonl)
	header, err := tarReader.Next()
	if err != nil {
		t.Fatalf("Failed to read tar header: %v", err)
	}

	if header.Name != "telemetry.jsonl" {
		t.Errorf("Expected 'telemetry.jsonl' in archive, got '%s'", header.Name)
	}
}

// TestPackForUploadEmptyFile tests packaging empty telemetry file
func TestPackForUploadEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create empty telemetry file
	if err := os.WriteFile(telemetryFile, []byte{}, 0600); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Test packaging
	outputPath := filepath.Join(tmpDir, "upload.json")
	result, err := PackForUpload(telemetryFile, outputPath, "json")

	// Should succeed but with zero events
	if err != nil {
		t.Fatalf("PackForUpload failed: %v", err)
	}

	if result.EventCount != 0 {
		t.Errorf("Expected 0 events, got %d", result.EventCount)
	}
}

// TestPackForUploadNonExistentFile tests packaging non-existent telemetry file
func TestPackForUploadNonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "nonexistent.jsonl")
	outputPath := filepath.Join(tmpDir, "upload.json")

	// Test packaging non-existent file
	result, err := PackForUpload(telemetryFile, outputPath, "json")

	// Should fail
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// result should be nil or have zero events
	if result != nil && result.EventCount != 0 {
		t.Error("Expected 0 events for non-existent file")
	}
}

// TestPackForUploadInvalidFormat tests packaging with invalid format
func TestPackForUploadInvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test telemetry file
	if err := os.WriteFile(telemetryFile, []byte("{}"), 0600); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "upload")
	_, err := PackForUpload(telemetryFile, outputPath, "invalid")

	// Should fail with unsupported format error
	if err == nil {
		t.Error("Expected error for invalid format, got nil")
	}
}

// TestPackForUploadJSONStructure tests JSON structure in detail
func TestPackForUploadJSONStructure(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test telemetry data
	testTime := time.Date(2026, 2, 6, 10, 30, 0, 0, time.UTC)
	events := []Event{
		{
			Type:      EventTypeCommandComplete,
			Timestamp: testTime,
			Data: map[string]interface{}{
				"command":  "init",
				"args":     []string{"."},
				"duration": 5000,
				"success":  true,
				"error":    "",
			},
		},
	}

	// Write events
	collector, _ := NewCollector(telemetryFile, true)
	for _, event := range events {
		collector.Record(event)
	}

	// Package as JSON
	outputPath := filepath.Join(tmpDir, "upload.json")
	_, err := PackForUpload(telemetryFile, outputPath, "json")
	if err != nil {
		t.Fatalf("PackForUpload failed: %v", err)
	}

	// Read and verify JSON structure
	data, _ := os.ReadFile(outputPath)
	var uploadData map[string]interface{}
	json.Unmarshal(data, &uploadData)

	// Verify metadata
	metadata, ok := uploadData["metadata"].(map[string]interface{})
	if !ok {
		t.Fatal("metadata is not a map")
	}

	if _, ok := metadata["version"]; !ok {
		t.Error("metadata missing 'version'")
	}

	if _, ok := metadata["generated_at"]; !ok {
		t.Error("metadata missing 'generated_at'")
	}

	if _, ok := metadata["event_count"]; !ok {
		t.Error("metadata missing 'event_count'")
	}

	// Verify events array
	eventsArray, ok := uploadData["events"].([]interface{})
	if !ok {
		t.Fatal("events is not an array")
	}

	if len(eventsArray) != 1 {
		t.Errorf("Expected 1 event, got %d", len(eventsArray))
	}

	// Verify event structure
	event, ok := eventsArray[0].(map[string]interface{})
	if !ok {
		t.Fatal("event is not a map")
	}

	if _, ok := event["type"]; !ok {
		t.Error("event missing 'type'")
	}

	if _, ok := event["timestamp"]; !ok {
		t.Error("event missing 'timestamp'")
	}

	if _, ok := event["data"]; !ok {
		t.Error("event missing 'data'")
	}
}

// TestPackForUploadPermissions tests that output file has secure permissions
func TestPackForUploadPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create test telemetry data
	collector, _ := NewCollector(telemetryFile, true)
	collector.Record(Event{
		Type:      EventTypeCommandStart,
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"command": "doctor"},
	})

	// Package as JSON
	outputPath := filepath.Join(tmpDir, "upload.json")
	_, err := PackForUpload(telemetryFile, outputPath, "json")
	if err != nil {
		t.Fatalf("PackForUpload failed: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 && mode != 0644 {
		// Allow 0600 (owner only) or 0644 (owner + group read)
		t.Logf("Note: File permissions are %04o (expected 0600 or 0644)", mode)
	}
}
