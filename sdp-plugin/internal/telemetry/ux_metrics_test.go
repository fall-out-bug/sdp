package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestUXMetricsCollector(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	// Create UX metrics collector
	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	// Verify collector is initialized
	if collector.eventsFile == "" {
		t.Error("Events file path is empty")
	}
	if collector.sessionID == "" {
		t.Error("Session ID is empty")
	}

	// Verify events file was created
	if _, err := os.Stat(collector.eventsFile); os.IsNotExist(err) {
		t.Error("Events file was not created")
	}
}

func TestRecordTimeToFirstValue(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	duration := 5 * time.Minute
	err = collector.RecordTimeToFirstValue(duration)
	if err != nil {
		t.Fatalf("Failed to record time to first value: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if event.Type != EventTypeUXMetric {
		t.Errorf("Expected event type %s, got %s", EventTypeUXMetric, event.Type)
	}

	// Check metric type in data
	metricType, ok := event.Data["metric_type"].(string)
	if !ok || metricType != string(UXMetricTimeToFirstValue) {
		t.Errorf("Expected metric type %s, got %v", UXMetricTimeToFirstValue, metricType)
	}
}

func TestRecordAssessComplete(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	duration := 2 * time.Minute
	err = collector.RecordAssessComplete("greenfield", duration)
	if err != nil {
		t.Fatalf("Failed to record assess complete: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check step name
	stepName, ok := event.Data["step_name"].(string)
	if !ok || stepName != "assess" {
		t.Errorf("Expected step name 'assess', got %v", stepName)
	}

	// Check project type
	projectType, ok := event.Data["project_type"].(string)
	if !ok || projectType != "greenfield" {
		t.Errorf("Expected project type 'greenfield', got %v", projectType)
	}
}

func TestRecordTryComplete(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	duration := 10 * time.Minute
	err = collector.RecordTryComplete("brownfield", duration)
	if err != nil {
		t.Fatalf("Failed to record try complete: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check step name
	stepName, ok := event.Data["step_name"].(string)
	if !ok || stepName != "try" {
		t.Errorf("Expected step name 'try', got %v", stepName)
	}

	// Check step number
	stepNumber, ok := event.Data["step_number"].(float64)
	if !ok || int(stepNumber) != 2 {
		t.Errorf("Expected step number 2, got %v", stepNumber)
	}
}

func TestRecordTryDiscard(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	err = collector.RecordTryDiscard("brownfield", "user_exited", 2)
	if err != nil {
		t.Fatalf("Failed to record try discard: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check metric_type
	metricType, ok := event.Data["metric_type"].(string)
	if !ok || metricType != "step_abandon_rate" {
		t.Errorf("Expected metric_type 'step_abandon_rate', got %v", metricType)
	}

	// Check exit reason in context
	context, ok := event.Data["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected context to be a map, got %T", event.Data["context"])
	}
	exitReason, ok := context["exit_reason"].(string)
	if !ok || exitReason != "user_exited" {
		t.Errorf("Expected exit reason 'user_exited', got %v", exitReason)
	}

	// Check value (should be false)
	value, ok := event.Data["value"].(bool)
	if !ok || value != false {
		t.Errorf("Expected value false, got %v", value)
	}
}

func TestRecordAdoptComplete(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	duration := 30 * time.Minute
	err = collector.RecordAdoptComplete("greenfield", duration)
	if err != nil {
		t.Fatalf("Failed to record adopt complete: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check step name
	stepName, ok := event.Data["step_name"].(string)
	if !ok || stepName != "adopt" {
		t.Errorf("Expected step name 'adopt', got %v", stepName)
	}

	// Check step number
	stepNumber, ok := event.Data["step_number"].(float64)
	if !ok || int(stepNumber) != 3 {
		t.Errorf("Expected step number 3, got %v", stepNumber)
	}
}

func TestRecordReset(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	err = collector.RecordReset("configuration_error")
	if err != nil {
		t.Fatalf("Failed to record reset: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check metric type
	metricType, ok := event.Data["metric_type"].(string)
	if !ok || metricType != string(UXMetricResetUninstallFrequency) {
		t.Errorf("Expected metric type %s, got %v", UXMetricResetUninstallFrequency, metricType)
	}

	// Check action in context
	context, ok := event.Data["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected context to be a map, got %T", event.Data["context"])
	}
	action, ok := context["action"].(string)
	if !ok || action != "reset" {
		t.Errorf("Expected action 'reset', got %v", action)
	}
}

func TestRecordUninstall(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	err = collector.RecordUninstall("not_suitable")
	if err != nil {
		t.Fatalf("Failed to record uninstall: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check action in context
	context, ok := event.Data["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected context to be a map, got %T", event.Data["context"])
	}
	action, ok := context["action"].(string)
	if !ok || action != "uninstall" {
		t.Errorf("Expected action 'uninstall', got %v", action)
	}
}

func TestRecordBrownfieldInitCompletion(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	details := map[string]any{
		"project_size": "large",
		"language":     "go",
	}

	err = collector.RecordBrownfieldInitCompletion(true, "complete", details)
	if err != nil {
		t.Fatalf("Failed to record brownfield init completion: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check project type
	projectType, ok := event.Data["project_type"].(string)
	if !ok || projectType != "brownfield" {
		t.Errorf("Expected project type 'brownfield', got %v", projectType)
	}

	// Check init phase in context
	context, ok := event.Data["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected context to be a map, got %T", event.Data["context"])
	}
	initPhase, ok := context["init_phase"].(string)
	if !ok || initPhase != "complete" {
		t.Errorf("Expected init phase 'complete', got %v", initPhase)
	}

	// Check value (should be true)
	value, ok := event.Data["value"].(bool)
	if !ok || value != true {
		t.Errorf("Expected value true, got %v", value)
	}

	// Check additional details in context
	projectSize, ok := context["project_size"].(string)
	if !ok || projectSize != "large" {
		t.Errorf("Expected project size 'large', got %v", projectSize)
	}
}

func TestRecordRecoveryAttempt(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	err = collector.RecordRecoveryAttempt(true, "auto_fix")
	if err != nil {
		t.Fatalf("Failed to record recovery attempt: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check recovery type in context
	context, ok := event.Data["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected context to be a map, got %T", event.Data["context"])
	}
	recoveryType, ok := context["recovery_type"].(string)
	if !ok || recoveryType != "auto_fix" {
		t.Errorf("Expected recovery type 'auto_fix', got %v", recoveryType)
	}

	// Check success value
	success, ok := event.Data["value"].(bool)
	if !ok || success != true {
		t.Errorf("Expected success true, got %v", success)
	}
}

func TestRecordSecondSessionReturn(t *testing.T) {
	tempDir := t.TempDir()
	sdpDir := filepath.Join(tempDir, ".sdp")

	collector, err := NewUXMetricsCollector(sdpDir)
	if err != nil {
		t.Fatalf("Failed to create UX metrics collector: %v", err)
	}

	err = collector.RecordSecondSessionReturn(3)
	if err != nil {
		t.Fatalf("Failed to record second session return: %v", err)
	}

	// Verify event was written
	eventsFile := collector.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("Failed to read events file: %v", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Check days since first session in context
	context, ok := event.Data["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected context to be a map, got %T", event.Data["context"])
	}
	daysSinceFirst, ok := context["days_since_first_session"].(float64)
	if !ok || int(daysSinceFirst) != 3 {
		t.Errorf("Expected 3 days since first session, got %v", daysSinceFirst)
	}

	// Check value (should be true since days > 0)
	value, ok := event.Data["value"].(bool)
	if !ok || value != true {
		t.Errorf("Expected value true, got %v", value)
	}
}

func TestUXMetricTypeValidation(t *testing.T) {
	tests := []struct {
		name   string
		metric UXMetricType
		valid  bool
	}{
		{"time_to_first_value", UXMetricTimeToFirstValue, true},
		{"step_abandon_rate", UXMetricStepAbandonRate, true},
		{"reset_uninstall_frequency", UXMetricResetUninstallFrequency, true},
		{"brownfield_init_completion", UXMetricBrownfieldInitCompletion, true},
		{"recovery_success_rate", UXMetricRecoverySuccessRate, true},
		{"second_session_return", UXMetricSecondSessionReturn, true},
		{"invalid_metric", UXMetricType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.metric.IsValid()
			if got != tt.valid {
				t.Errorf("UXMetricType.IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestEventTypeValidationWithUXTypes(t *testing.T) {
	tests := []struct {
		name  string
		etype EventType
		valid bool
	}{
		{"ux_metric", EventTypeUXMetric, true},
		{"assess_complete", EventTypeAssessComplete, true},
		{"try_complete", EventTypeTryComplete, true},
		{"try_discard", EventTypeTryDiscard, true},
		{"adopt_complete", EventTypeAdoptComplete, true},
		{"reset", EventTypeReset, true},
		{"uninstall", EventTypeUninstall, true},
		{"invalid_type", EventType("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.etype.IsValid()
			if got != tt.valid {
				t.Errorf("EventType.IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestExportUXMetrics(t *testing.T) {
	tempDir := t.TempDir()

	// Create a collector with events file
	eventsFile := filepath.Join(tempDir, "events.jsonl")
	collector, err := NewCollector(eventsFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Record some UX metric events
	collector.Record(Event{
		Type:      EventTypeUXMetric,
		Timestamp: time.Now(),
		Data: map[string]any{
			"metric_type": string(UXMetricTimeToFirstValue),
			"value":       int64(300000),
		},
	})

	collector.Record(Event{
		Type:      EventTypeUXMetric,
		Timestamp: time.Now(),
		Data: map[string]any{
			"metric_type": string(UXMetricStepAbandonRate),
			"value":       15.5,
		},
	})

	// Export UX metrics
	exportPath := filepath.Join(tempDir, "ux_metrics_export.json")
	err = collector.ExportUXMetrics(exportPath)
	if err != nil {
		t.Fatalf("Failed to export UX metrics: %v", err)
	}

	// Verify export file exists
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		t.Error("Export file was not created")
	}

	// Read and verify export
	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	var exportedMetrics []map[string]any
	if err := json.Unmarshal(data, &exportedMetrics); err != nil {
		t.Fatalf("Failed to unmarshal exported metrics: %v", err)
	}

	if len(exportedMetrics) != 2 {
		t.Errorf("Expected 2 exported metrics, got %d", len(exportedMetrics))
	}
}

func TestGetUXMetrics(t *testing.T) {
	tempDir := t.TempDir()

	// Create a collector with events file
	eventsFile := filepath.Join(tempDir, "events.jsonl")
	collector, err := NewCollector(eventsFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Record mixed events (UX and non-UX)
	collector.Record(Event{
		Type:      EventTypeCommandStart,
		Timestamp: time.Now(),
		Data: map[string]any{
			"command": "test",
		},
	})

	collector.Record(Event{
		Type:      EventTypeUXMetric,
		Timestamp: time.Now(),
		Data: map[string]any{
			"metric_type": string(UXMetricTimeToFirstValue),
			"value":       int64(300000),
		},
	})

	// Get UX metrics
	uxMetrics, err := collector.GetUXMetrics()
	if err != nil {
		t.Fatalf("Failed to get UX metrics: %v", err)
	}

	if len(uxMetrics) != 1 {
		t.Errorf("Expected 1 UX metric, got %d", len(uxMetrics))
	}

	// Verify the metric
	metricType, ok := uxMetrics[0]["metric_type"].(string)
	if !ok || metricType != string(UXMetricTimeToFirstValue) {
		t.Errorf("Expected metric type %s, got %v", UXMetricTimeToFirstValue, metricType)
	}
}
