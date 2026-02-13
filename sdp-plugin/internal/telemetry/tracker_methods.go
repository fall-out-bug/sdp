package telemetry

import (
	"time"
)

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
