package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Global tracker instance
var globalTracker *Tracker
var trackerOnce sync.Once

// Tracker manages automatic telemetry tracking for CLI commands
type Tracker struct {
	collector      *Collector
	currentCommand *CommandEvent
	mu             sync.Mutex
}

// CommandEvent represents a command execution
type CommandEvent struct {
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Duration  time.Duration     `json:"duration"`
	Success   bool              `json:"success"`
	Error     string            `json:"error,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// GetTracker returns the global telemetry tracker instance
func GetTracker() *Tracker {
	trackerOnce.Do(func() {
		configDir, err := os.UserConfigDir()
		if err != nil {
			// If we can't get config dir, return a disabled tracker
			globalTracker = &Tracker{collector: &Collector{}}
			return
		}

		telemetryFile := filepath.Join(configDir, "sdp", "telemetry.jsonl")

		// Check if telemetry is disabled
		configPath := filepath.Join(configDir, "sdp", "telemetry.json")
		enabled := true
		if data, err := os.ReadFile(configPath); err == nil {
			var config map[string]bool
			if err := json.Unmarshal(data, &config); err == nil {
				if disabled, ok := config["disabled"]; ok && disabled {
					enabled = false
				}
			}
		}

		collector, err := NewCollector(telemetryFile, enabled)
		if err != nil {
			// If collector creation fails, return a disabled tracker
			globalTracker = &Tracker{collector: &Collector{}}
			return
		}

		globalTracker = &Tracker{collector: collector}
	})

	return globalTracker
}

// TrackCommandStart tracks the start of a command
func (t *Tracker) TrackCommandStart(command string, args []string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.collector == nil {
		return nil
	}

	event := &CommandEvent{
		Command:   command,
		Args:      args,
		StartTime: time.Now(),
	}

	t.currentCommand = event

	// Record command_start event
	telemetryEvent := Event{
		Type:      EventTypeCommandStart,
		Timestamp: event.StartTime,
		Data: map[string]interface{}{
			"command": event.Command,
			"args":    event.Args,
		},
	}

	return t.collector.Record(telemetryEvent)
}

// TrackCommandComplete tracks the completion of a command
func (t *Tracker) TrackCommandComplete(success bool, errMsg string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.collector == nil || t.currentCommand == nil {
		return nil
	}

	t.currentCommand.EndTime = time.Now()
	t.currentCommand.Duration = t.currentCommand.EndTime.Sub(t.currentCommand.StartTime)
	t.currentCommand.Success = success
	t.currentCommand.Error = errMsg

	// Record command_complete event
	telemetryEvent := Event{
		Type:      EventTypeCommandComplete,
		Timestamp: t.currentCommand.EndTime,
		Data: map[string]interface{}{
			"command":  t.currentCommand.Command,
			"args":     t.currentCommand.Args,
			"duration": t.currentCommand.Duration.Milliseconds(),
			"success":  t.currentCommand.Success,
			"error":    t.currentCommand.Error,
		},
	}

	err := t.collector.Record(telemetryEvent)

	// Reset current command
	t.currentCommand = nil

	return err
}

// TrackWorkstreamStart tracks the start of a workstream
func (t *Tracker) TrackWorkstreamStart(wsID string) error {
	if t.collector == nil {
		return nil
	}

	event := Event{
		Type:      EventTypeWSStart,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"ws_id": wsID,
		},
	}

	return t.collector.Record(event)
}

// TrackWorkstreamComplete tracks the completion of a workstream
func (t *Tracker) TrackWorkstreamComplete(wsID string, success bool, duration time.Duration) error {
	if t.collector == nil {
		return nil
	}

	event := Event{
		Type:      EventTypeWSComplete,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"ws_id":    wsID,
			"success":  success,
			"duration": duration.Milliseconds(),
		},
	}

	return t.collector.Record(event)
}

// TrackQualityGateResult tracks a quality gate result
func (t *Tracker) TrackQualityGateResult(gateName string, passed bool, score float64) error {
	if t.collector == nil {
		return nil
	}

	event := Event{
		Type:      EventTypeQualityGateResult,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"gate":   gateName,
			"passed": passed,
			"score":  score,
		},
	}

	return t.collector.Record(event)
}

// IsEnabled returns whether telemetry is enabled
func (t *Tracker) IsEnabled() bool {
	if t.collector == nil {
		return false
	}

	status := t.collector.Status()
	return status.Enabled
}

// Disable disables telemetry tracking
func (t *Tracker) Disable() {
	if t.collector != nil {
		t.collector.Disable()
	}
}

// Enable enables telemetry tracking
func (t *Tracker) Enable() {
	if t.collector != nil {
		t.collector.Enable()
	}
}

// GetStatus returns the current telemetry status
func (t *Tracker) GetStatus() *Status {
	if t.collector == nil {
		return &Status{Enabled: false, EventCount: 0}
	}

	status := t.collector.Status()
	return &status
}

// Helper functions for easy access without GetTracker()

// TrackCommandStart is a convenience function that uses the global tracker
func TrackCommandStart(command string, args []string) error {
	return GetTracker().TrackCommandStart(command, args)
}

// TrackCommandComplete is a convenience function that uses the global tracker
func TrackCommandComplete(success bool, errMsg string) error {
	return GetTracker().TrackCommandComplete(success, errMsg)
}

// TrackWorkstreamStart is a convenience function that uses the global tracker
func TrackWorkstreamStart(wsID string) error {
	return GetTracker().TrackWorkstreamStart(wsID)
}

// TrackWorkstreamComplete is a convenience function that uses the global tracker
func TrackWorkstreamComplete(wsID string, success bool, duration time.Duration) error {
	return GetTracker().TrackWorkstreamComplete(wsID, success, duration)
}

// TrackQualityGateResult is a convenience function that uses the global tracker
func TrackQualityGateResult(gateName string, passed bool, score float64) error {
	return GetTracker().TrackQualityGateResult(gateName, passed, score)
}

// IsTelemetryEnabled is a convenience function that uses the global tracker
func IsTelemetryEnabled() bool {
	return GetTracker().IsEnabled()
}

// GetTelemetryStatus is a convenience function that uses the global tracker
func GetTelemetryStatus() *Status {
	return GetTracker().GetStatus()
}

// DisableTelemetry is a convenience function that uses the global tracker
func DisableTelemetry() {
	GetTracker().Disable()
}

// EnableTelemetry is a convenience function that uses the global tracker
func EnableTelemetry() {
	GetTracker().Enable()
}
