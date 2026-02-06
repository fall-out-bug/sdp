package orchestrator

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// TestNewOrchestratorLogger tests logger creation
func TestNewOrchestratorLogger(t *testing.T) {
	ol := NewOrchestratorLogger("test-feature")

	if ol == nil {
		t.Fatal("NewOrchestratorLogger() returned nil")
	}

	if ol.logger == nil {
		t.Error("logger not initialized")
	}

	if ol.correlationID == "" {
		t.Error("correlation ID not generated")
	}
}

// TestCorrelationIDUnique tests that each logger gets unique correlation ID
func TestCorrelationIDUnique(t *testing.T) {
	ol1 := NewOrchestratorLogger("feature-1")
	ol2 := NewOrchestratorLogger("feature-2")

	if ol1.correlationID == ol2.correlationID {
		t.Error("Correlation IDs should be unique")
	}
}

// TestLogStart tests logging orchestration start
func TestLogStart(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{
		logger:         logger,
		correlationID:  "test-correlation-123",
	}

	ol.LogStart("test-feature")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "test-feature") {
		t.Error("LogStart() missing feature_id")
	}
	if !strings.Contains(logOutput, "test-correlation-123") {
		t.Error("LogStart() missing correlation_id")
	}
	if !strings.Contains(logOutput, "orchestration_start") {
		t.Error("LogStart() missing event type")
	}
}

// TestLogWSStart tests logging workstream start
func TestLogWSStart(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{
		logger:         logger,
		correlationID:  "test-correlation-123",
	}

	ol.LogWSStart("00-001-01")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "00-001-01") {
		t.Error("LogWSStart() missing ws_id")
	}
	if !strings.Contains(logOutput, "test-correlation-123") {
		t.Error("LogWSStart() missing correlation_id")
	}
	if !strings.Contains(logOutput, "workstream_start") {
		t.Error("LogWSStart() missing event type")
	}
}

// TestLogWSComplete tests logging workstream completion
func TestLogWSComplete(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{
		logger:         logger,
		correlationID:  "test-correlation-123",
	}

	ol.LogWSComplete("00-001-01", 5*time.Minute, 85.5)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "00-001-01") {
		t.Error("LogWSComplete() missing ws_id")
	}
	if !strings.Contains(logOutput, "duration") {
		t.Error("LogWSComplete() missing duration")
	}
	if !strings.Contains(logOutput, "85.5") {
		t.Error("LogWSComplete() missing coverage")
	}
	if !strings.Contains(logOutput, "workstream_complete") {
		t.Error("LogWSComplete() missing event type")
	}
}

// TestLogWSCheckpoint tests logging checkpoint save
func TestLogWSCheckpoint(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{
		logger:         logger,
		correlationID:  "test-correlation-123",
	}

	ol.LogWSCheckpoint("00-001-01", ".oneshot/test-checkpoint.json")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "00-001-01") {
		t.Error("LogWSCheckpoint() missing ws_id")
	}
	if !strings.Contains(logOutput, ".oneshot/test-checkpoint.json") {
		t.Error("LogWSCheckpoint() missing checkpoint file")
	}
	if !strings.Contains(logOutput, "checkpoint_saved") {
		t.Error("LogWSCheckpoint() missing event type")
	}
}

// TestLogWSError tests logging workstream errors
func TestLogWSError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{
		logger:         logger,
		correlationID:  "test-correlation-123",
	}

	ol.LogWSError("00-001-01", 1, errors.New("test error"))

	logOutput := buf.String()
	if !strings.Contains(logOutput, "00-001-01") {
		t.Error("LogWSError() missing ws_id")
	}
	if !strings.Contains(logOutput, "test error") {
		t.Error("LogWSError() missing error message")
	}
	if !strings.Contains(logOutput, "attempt") {
		t.Error("LogWSError() missing attempt count")
	}
	if !strings.Contains(logOutput, "workstream_error") {
		t.Error("LogWSError() missing event type")
	}
}

// TestLogOrchestrationComplete tests logging orchestration completion
func TestLogOrchestrationComplete(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{
		logger:         logger,
		correlationID:  "test-correlation-123",
	}

	ol.LogOrchestrationComplete("test-feature", 2*time.Hour, 18, true)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "test-feature") {
		t.Error("LogOrchestrationComplete() missing feature_id")
	}
	if !strings.Contains(logOutput, "duration") {
		t.Error("LogOrchestrationComplete() missing duration")
	}
	if !strings.Contains(logOutput, "total_workstreams") {
		t.Error("LogOrchestrationComplete() missing total workstreams")
	}
	if !strings.Contains(logOutput, "18") {
		t.Error("LogOrchestrationComplete() missing workstream count")
	}
	if !strings.Contains(logOutput, "true") {
		t.Error("LogOrchestrationComplete() missing success status")
	}
	if !strings.Contains(logOutput, "orchestration_complete") {
		t.Error("LogOrchestrationComplete() missing event type")
	}
}

// TestStructuredFields tests that all required fields are present
func TestStructuredFields(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	ol := &OrchestratorLogger{
		logger:         logger,
		correlationID:  "test-correlation-456",
	}

	ol.LogWSStart("00-001-01")

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log JSON: %v", err)
	}

	// Check for required fields
	requiredFields := []string{"level", "msg", "time", "correlation_id"}
	for _, field := range requiredFields {
		if _, exists := logEntry[field]; !exists {
			t.Errorf("Log missing required field: %s", field)
		}
	}

	if logEntry["correlation_id"] != "test-correlation-456" {
		t.Errorf("correlation_id mismatch: got %v, want test-correlation-456", logEntry["correlation_id"])
	}
}
