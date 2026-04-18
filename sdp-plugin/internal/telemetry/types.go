package telemetry

import (
	"time"
)

// EventType represents the type of telemetry event
type EventType string

const (
	EventTypeCommandStart      EventType = "command_start"
	EventTypeCommandComplete   EventType = "command_complete"
	EventTypeWSStart           EventType = "ws_start"
	EventTypeWSComplete        EventType = "ws_complete"
	EventTypeQualityGateResult EventType = "quality_gate_result"

	// UX metric events for measuring adoption journey
	EventTypeUXMetric       EventType = "ux_metric"
	EventTypeAssessComplete EventType = "assess_complete"
	EventTypeTryComplete    EventType = "try_complete"
	EventTypeTryDiscard     EventType = "try_discard"
	EventTypeAdoptComplete  EventType = "adopt_complete"
	EventTypeReset          EventType = "reset"
	EventTypeUninstall      EventType = "uninstall"
)

// IsValid checks if the event type is valid
func (et EventType) IsValid() bool {
	switch et {
	case EventTypeCommandStart, EventTypeCommandComplete,
		EventTypeWSStart, EventTypeWSComplete, EventTypeQualityGateResult,
		EventTypeUXMetric, EventTypeAssessComplete, EventTypeTryComplete,
		EventTypeTryDiscard, EventTypeAdoptComplete, EventTypeReset, EventTypeUninstall:
		return true
	default:
		return false
	}
}

// Event represents a telemetry event
type Event struct {
	Type      EventType      `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data"`
}

// Status represents the current status of telemetry
type Status struct {
	Enabled    bool   `json:"enabled"`
	EventCount int    `json:"event_count"`
	FilePath   string `json:"file_path"`
}

// UXMetricType represents the type of UX metric
type UXMetricType string

const (
	UXMetricTimeToFirstValue         UXMetricType = "time_to_first_value"
	UXMetricStepAbandonRate          UXMetricType = "step_abandon_rate"
	UXMetricResetUninstallFrequency  UXMetricType = "reset_uninstall_frequency"
	UXMetricBrownfieldInitCompletion UXMetricType = "brownfield_init_completion"
	UXMetricRecoverySuccessRate      UXMetricType = "recovery_success_rate"
	UXMetricSecondSessionReturn      UXMetricType = "second_session_return"
)

// UXMetricEvent represents a UX metric measurement
type UXMetricEvent struct {
	MetricType  UXMetricType   `json:"metric_type"`
	Timestamp   time.Time      `json:"timestamp"`
	Value       interface{}    `json:"value"` // Can be number, bool, or string
	SessionID   string         `json:"session_id,omitempty"`
	ProjectType string         `json:"project_type,omitempty"` // "greenfield", "brownfield", "unknown"
	StepName    string         `json:"step_name,omitempty"`    // e.g., "assess", "try", "adopt"
	StepNumber  int            `json:"step_number,omitempty"`
	Context     map[string]any `json:"context,omitempty"`
}

// IsValid checks if the UX metric type is valid
func (mt UXMetricType) IsValid() bool {
	switch mt {
	case UXMetricTimeToFirstValue, UXMetricStepAbandonRate,
		UXMetricResetUninstallFrequency, UXMetricBrownfieldInitCompletion,
		UXMetricRecoverySuccessRate, UXMetricSecondSessionReturn:
		return true
	default:
		return false
	}
}
