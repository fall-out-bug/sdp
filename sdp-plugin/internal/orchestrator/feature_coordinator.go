package orchestrator

import (
	"time"
)

// ProgressUpdate represents a progress update during execution
type ProgressUpdate struct {
	Timestamp    time.Time
	Message      string
	WorkstreamID string
	Status       string
}

// ProgressCallback is called for progress updates
type ProgressCallback func(update ProgressUpdate)

// FeatureCoordinator coordinates feature execution with orchestrator
type FeatureCoordinator struct {
	orchestrator     *Orchestrator
	progressCallback ProgressCallback
}

// NewFeatureCoordinator creates a new feature coordinator
func NewFeatureCoordinator(
	loader WorkstreamLoader,
	executor WorkstreamExecutor,
	saver CheckpointSaver,
	maxRetries int,
) *FeatureCoordinator {
	orch := NewOrchestrator(loader, executor, saver, maxRetries)
	return &FeatureCoordinator{
		orchestrator:     orch,
		progressCallback: nil,
	}
}

// SetProgressCallback sets the progress callback
func (fc *FeatureCoordinator) SetProgressCallback(callback ProgressCallback) {
	fc.progressCallback = callback
}

// sendProgress sends a progress update if callback is set
func (fc *FeatureCoordinator) sendProgress(message, workstreamID, status string) {
	if fc.progressCallback != nil {
		fc.progressCallback(ProgressUpdate{
			Timestamp:    time.Now(),
			Message:      message,
			WorkstreamID: workstreamID,
			Status:       status,
		})
	}
}

// formatTime formats a time as HH:MM
func formatTime(t time.Time) string {
	return t.Format("15:04")
}
