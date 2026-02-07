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
)

// IsValid checks if the event type is valid
func (et EventType) IsValid() bool {
	switch et {
	case EventTypeCommandStart, EventTypeCommandComplete,
		EventTypeWSStart, EventTypeWSComplete, EventTypeQualityGateResult:
		return true
	default:
		return false
	}
}

// Event represents a telemetry event
type Event struct {
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Status represents the current status of telemetry
type Status struct {
	Enabled    bool   `json:"enabled"`
	EventCount int    `json:"event_count"`
	FilePath   string `json:"file_path"`
}
