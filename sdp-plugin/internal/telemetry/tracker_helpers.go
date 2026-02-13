package telemetry

import "time"

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
