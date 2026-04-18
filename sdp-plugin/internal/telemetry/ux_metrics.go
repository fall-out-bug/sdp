package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// UXMetricsCollector manages UX metrics collection
type UXMetricsCollector struct {
	eventsFile string
	mu         sync.Mutex
	sessionID  string
}

// NewUXMetricsCollector creates a new UX metrics collector
// UX metrics are stored in the user's config directory (~/.config/sdp/ux-metrics.jsonl)
// rather than in the project-local .sdp/ directory to avoid polluting assessed repositories.
func NewUXMetricsCollector(sdpDir string) (*UXMetricsCollector, error) {
	// Get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config directory: %w", err)
	}

	// Create SDP config directory
	sdpConfigDir := filepath.Join(configDir, "sdp")
	if err := os.MkdirAll(sdpConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create SDP config directory: %w", err)
	}

	eventsFile := filepath.Join(sdpConfigDir, "ux-metrics.jsonl")

	// Create or verify UX metrics file exists
	if _, err := os.OpenFile(eventsFile, os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		return nil, fmt.Errorf("failed to create UX metrics file: %w", err)
	}

	// Generate session ID
	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())

	return &UXMetricsCollector{
		eventsFile: eventsFile,
		sessionID:  sessionID,
	}, nil
}

// RecordMetric records a UX metric event
func (ux *UXMetricsCollector) RecordMetric(event UXMetricEvent) error {
	// Validate metric type
	if !event.MetricType.IsValid() {
		return fmt.Errorf("invalid UX metric type: %s", event.MetricType)
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Set session ID if not provided
	if event.SessionID == "" {
		event.SessionID = ux.sessionID
	}

	// Create telemetry event
	teleEvent := Event{
		Type:      EventTypeUXMetric,
		Timestamp: event.Timestamp,
		Data:      map[string]any{},
	}

	// Marshal UX metric event to JSON and store in data
	dataBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal UX metric event: %w", err)
	}

	// Parse back into map for Data field
	if err := json.Unmarshal(dataBytes, &teleEvent.Data); err != nil {
		return fmt.Errorf("failed to unmarshal UX metric data: %w", err)
	}

	// Append to events file (appendEvent handles locking)
	return ux.appendEvent(teleEvent)
}

// RecordAssessComplete records completion of assess phase
func (ux *UXMetricsCollector) RecordAssessComplete(projectType string, duration time.Duration) error {
	event := Event{
		Type:      EventTypeAssessComplete,
		Timestamp: time.Now(),
		Data: map[string]any{
			"project_type": projectType,
			"step_name":    "assess",
			"step_number":  1,
			"duration_ms":  duration.Milliseconds(),
			"value":        duration.Milliseconds(),
		},
	}
	return ux.appendEvent(event)
}

// RecordTryComplete records completion of try phase
func (ux *UXMetricsCollector) RecordTryComplete(projectType string, duration time.Duration) error {
	event := Event{
		Type:      EventTypeTryComplete,
		Timestamp: time.Now(),
		Data: map[string]any{
			"project_type": projectType,
			"step_name":    "try",
			"step_number":  2,
			"duration_ms":  duration.Milliseconds(),
			"value":        duration.Milliseconds(),
		},
	}
	return ux.appendEvent(event)
}

// RecordTryDiscard records abandonment during try phase
func (ux *UXMetricsCollector) RecordTryDiscard(projectType, reason string, stepNumber int) error {
	event := Event{
		Type:      EventTypeTryDiscard,
		Timestamp: time.Now(),
		Data: map[string]any{
			"project_type": projectType,
			"step_name":    "try",
			"step_number":  stepNumber,
			"exit_reason":  reason,
			"value":        false,
		},
	}
	return ux.appendEvent(event)
}

// RecordAdoptComplete records completion of adopt phase
func (ux *UXMetricsCollector) RecordAdoptComplete(projectType string, duration time.Duration) error {
	event := Event{
		Type:      EventTypeAdoptComplete,
		Timestamp: time.Now(),
		Data: map[string]any{
			"project_type": projectType,
			"step_name":    "adopt",
			"step_number":  3,
			"duration_ms":  duration.Milliseconds(),
			"value":        duration.Milliseconds(),
		},
	}
	return ux.appendEvent(event)
}

// RecordReset records a reset event
func (ux *UXMetricsCollector) RecordReset(reason string) error {
	uxEvent := UXMetricEvent{
		MetricType: UXMetricResetUninstallFrequency,
		Timestamp:  time.Now(),
		Value:      1,
		Context: map[string]any{
			"action":      "reset",
			"exit_reason": reason,
		},
	}
	return ux.RecordMetric(uxEvent)
}

// RecordUninstall records an uninstall event
func (ux *UXMetricsCollector) RecordUninstall(reason string) error {
	uxEvent := UXMetricEvent{
		MetricType: UXMetricResetUninstallFrequency,
		Timestamp:  time.Now(),
		Value:      1,
		Context: map[string]any{
			"action":      "uninstall",
			"exit_reason": reason,
		},
	}
	return ux.RecordMetric(uxEvent)
}

// RecordTimeToFirstValue records time from init to first successful feature
func (ux *UXMetricsCollector) RecordTimeToFirstValue(duration time.Duration) error {
	uxEvent := UXMetricEvent{
		MetricType: UXMetricTimeToFirstValue,
		Timestamp:  time.Now(),
		Value:      duration.Milliseconds(),
		Context: map[string]any{
			"duration_ms": duration.Milliseconds(),
		},
	}
	return ux.RecordMetric(uxEvent)
}

// RecordBrownfieldInitCompletion records brownfield init completion status
func (ux *UXMetricsCollector) RecordBrownfieldInitCompletion(success bool, phase string, details map[string]any) error {
	uxEvent := UXMetricEvent{
		MetricType:  UXMetricBrownfieldInitCompletion,
		Timestamp:   time.Now(),
		Value:       success,
		ProjectType: "brownfield",
		Context:     make(map[string]any),
	}

	if phase != "" {
		uxEvent.Context["init_phase"] = phase
	}

	// Merge additional details
	for k, v := range details {
		uxEvent.Context[k] = v
	}

	return ux.RecordMetric(uxEvent)
}

// RecordRecoveryAttempt records a recovery attempt and its success
func (ux *UXMetricsCollector) RecordRecoveryAttempt(success bool, recoveryType string) error {
	uxEvent := UXMetricEvent{
		MetricType: UXMetricRecoverySuccessRate,
		Timestamp:  time.Now(),
		Value:      success,
		Context: map[string]any{
			"recovery_type": recoveryType,
			"success":       success,
		},
	}
	return ux.RecordMetric(uxEvent)
}

// RecordSecondSessionReturn records whether user returned for a second session
func (ux *UXMetricsCollector) RecordSecondSessionReturn(daysSinceFirst int) error {
	uxEvent := UXMetricEvent{
		MetricType: UXMetricSecondSessionReturn,
		Timestamp:  time.Now(),
		Value:      daysSinceFirst > 0,
		Context: map[string]any{
			"days_since_first_session": daysSinceFirst,
		},
	}
	return ux.RecordMetric(uxEvent)
}

// appendEvent appends an event to the events file
func (ux *UXMetricsCollector) appendEvent(event Event) error {
	// Marshal event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Append to file with secure permissions
	ux.mu.Lock()
	defer ux.mu.Unlock()

	file, err := os.OpenFile(ux.eventsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open events file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	return nil
}

// GetEventsFile returns the path to the events file
func (ux *UXMetricsCollector) GetEventsFile() string {
	return ux.eventsFile
}

// GetSessionID returns the current session ID
func (ux *UXMetricsCollector) GetSessionID() string {
	return ux.sessionID
}
