package telemetry

import (
	"path/filepath"
	"testing"
	"time"
)

func TestTrackerTrackCommand(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	tracker := &Tracker{collector: collector}

	// Track command start
	err = tracker.TrackCommandStart("test-cmd", []string{"arg1", "arg2"})
	if err != nil {
		t.Fatalf("Failed to track command start: %v", err)
	}

	// Track command complete
	err = tracker.TrackCommandComplete(true, "")
	if err != nil {
		t.Fatalf("Failed to track command complete: %v", err)
	}

	// Verify events were recorded
	status := collector.Status()
	if status.EventCount != 2 {
		t.Errorf("Expected 2 events, got %d", status.EventCount)
	}
}

func TestTrackerTrackWorkstream(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	tracker := &Tracker{collector: collector}

	// Track workstream start
	err = tracker.TrackWorkstreamStart("00-001-01")
	if err != nil {
		t.Fatalf("Failed to track workstream start: %v", err)
	}

	// Track workstream complete
	duration := 5 * time.Minute
	err = tracker.TrackWorkstreamComplete("00-001-01", true, duration)
	if err != nil {
		t.Fatalf("Failed to track workstream complete: %v", err)
	}

	// Verify events were recorded
	status := collector.Status()
	if status.EventCount != 2 {
		t.Errorf("Expected 2 events, got %d", status.EventCount)
	}
}

func TestTrackerTrackQualityGate(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	tracker := &Tracker{collector: collector}

	// Track quality gate result
	err = tracker.TrackQualityGateResult("coverage", true, 85.5)
	if err != nil {
		t.Fatalf("Failed to track quality gate: %v", err)
	}

	// Verify event was recorded
	status := collector.Status()
	if status.EventCount != 1 {
		t.Errorf("Expected 1 event, got %d", status.EventCount)
	}
}

func TestTrackerDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, false)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	tracker := &Tracker{collector: collector}

	// Try to track events
	tracker.TrackCommandStart("test", []string{})
	tracker.TrackCommandComplete(true, "")

	// Verify no events were recorded
	status := collector.Status()
	if status.EventCount != 0 {
		t.Errorf("Expected 0 events when disabled, got %d", status.EventCount)
	}
}

func TestTrackerNilCollector(t *testing.T) {
	tracker := &Tracker{collector: nil}

	// These should not panic
	err := tracker.TrackCommandStart("test", []string{})
	if err != nil {
		t.Errorf("Expected nil error with nil collector, got %v", err)
	}

	err = tracker.TrackCommandComplete(true, "")
	if err != nil {
		t.Errorf("Expected nil error with nil collector, got %v", err)
	}

	err = tracker.TrackWorkstreamStart("00-001-01")
	if err != nil {
		t.Errorf("Expected nil error with nil collector, got %v", err)
	}

	err = tracker.TrackWorkstreamComplete("00-001-01", true, time.Minute)
	if err != nil {
		t.Errorf("Expected nil error with nil collector, got %v", err)
	}

	err = tracker.TrackQualityGateResult("coverage", true, 85.5)
	if err != nil {
		t.Errorf("Expected nil error with nil collector, got %v", err)
	}
}

func TestTrackerIsEnabled(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Test enabled
	collectorEnabled, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	trackerEnabled := &Tracker{collector: collectorEnabled}
	if !trackerEnabled.IsEnabled() {
		t.Error("Expected tracker to be enabled")
	}

	// Test disabled
	collectorDisabled, err := NewCollector(telemetryFile, false)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	trackerDisabled := &Tracker{collector: collectorDisabled}
	if trackerDisabled.IsEnabled() {
		t.Error("Expected tracker to be disabled")
	}

	// Test nil collector
	trackerNil := &Tracker{collector: nil}
	if trackerNil.IsEnabled() {
		t.Error("Expected tracker with nil collector to be disabled")
	}
}

func TestTrackerGetStatus(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	tracker := &Tracker{collector: collector}

	// Record an event
	tracker.TrackWorkstreamStart("00-001-01")

	// Get status
	status := tracker.GetStatus()
	if status == nil {
		t.Fatal("Expected status to not be nil")
	}

	if !status.Enabled {
		t.Error("Expected status to show enabled")
	}

	if status.EventCount != 1 {
		t.Errorf("Expected event count 1, got %d", status.EventCount)
	}
}

func TestTrackerCommandDuration(t *testing.T) {
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	collector, err := NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	tracker := &Tracker{collector: collector}

	// Track command start
	err = tracker.TrackCommandStart("test-cmd", []string{})
	if err != nil {
		t.Fatalf("Failed to track command start: %v", err)
	}

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Track command complete
	err = tracker.TrackCommandComplete(true, "")
	if err != nil {
		t.Fatalf("Failed to track command complete: %v", err)
	}

	// Read events to verify duration
	events, err := collector.readEvents()
	if err != nil {
		t.Fatalf("Failed to read events: %v", err)
	}

	// Find command_complete event
	var completeEvent *Event
	for _, event := range events {
		if event.Type == EventTypeCommandComplete {
			completeEvent = &event
			break
		}
	}

	if completeEvent == nil {
		t.Fatal("Command complete event not found")
	}

	duration, ok := completeEvent.Data["duration"].(float64)
	if !ok {
		t.Fatal("Duration not found in event data")
	}

	if duration < 10 {
		t.Errorf("Expected duration >= 10ms, got %v", duration)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Note: These tests use the global tracker, so we can't easily test them
	// in isolation. We'll just verify they don't panic.

	// These should not panic
	_ = TrackCommandStart("test", []string{})
	_ = TrackCommandComplete(true, "")
	_ = TrackWorkstreamStart("00-001-01")
	_ = TrackWorkstreamComplete("00-001-01", true, time.Minute)
	_ = TrackQualityGateResult("coverage", true, 85.5)
	_ = IsTelemetryEnabled()
	_ = GetTelemetryStatus()
	DisableTelemetry()
	EnableTelemetry()
}
