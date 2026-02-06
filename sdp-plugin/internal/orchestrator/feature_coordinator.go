package orchestrator

import (
	"fmt"
	"time"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
)

// ProgressUpdate represents a progress update during execution
type ProgressUpdate struct {
	Timestamp time.Time
	Message   string
	WorkstreamID string
	Status    string
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

// ExecuteFeature executes all workstreams for a feature
func (fc *FeatureCoordinator) ExecuteFeature(featureID string) error {
	startTime := time.Now()

	fc.sendProgress(
		fmt.Sprintf("[%s] Starting feature execution: %s", formatTime(startTime), featureID),
		"",
		"started",
	)

	// Load workstreams
	fc.sendProgress(
		fmt.Sprintf("[%s] Loading workstreams for %s...", formatTime(time.Now()), featureID),
		"",
		"loading",
	)

	workstreams, err := fc.orchestrator.loader.LoadWorkstreams(featureID)
	if err != nil {
		return fmt.Errorf("failed to load workstreams: %w", err)
	}

	if len(workstreams) == 0 {
		return ErrFeatureNotFound
	}

	// Build dependency graph
	fc.sendProgress(
		fmt.Sprintf("[%s] Building dependency graph...", formatTime(time.Now())),
		"",
		"building_graph",
	)

	graph, err := BuildDependencyGraph(workstreams)
	if err != nil {
		return err
	}

	// Get execution order
	order, err := TopologicalSort(graph)
	if err != nil {
		return err
	}

	fc.sendProgress(
		fmt.Sprintf("[%s] Execution order: %v", formatTime(time.Now()), order),
		"",
		"execution_order",
	)

	// Create initial checkpoint
	cp := checkpoint.Checkpoint{
		ID:                   featureID,
		FeatureID:            featureID,
		Status:               checkpoint.StatusPending,
		CompletedWorkstreams: []string{},
		CurrentWorkstream:    "",
		CreatedAt:            startTime,
		UpdatedAt:            startTime,
	}

	// Execute workstreams
	for i, wsID := range order {
		cp.CurrentWorkstream = wsID
		cp.Status = checkpoint.StatusInProgress
		cp.UpdatedAt = time.Now()

		// Save checkpoint before execution
		if err := fc.orchestrator.checkpoint.Save(cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}

		// Send progress update
		fc.sendProgress(
			fmt.Sprintf("[%s] Executing %s (%d/%d)...", formatTime(time.Now()), wsID, i+1, len(order)),
			wsID,
			"executing",
		)

		// Execute with retry
		wsStartTime := time.Now()
		err := fc.executeWithRetry(wsID)
		wsDuration := time.Since(wsStartTime)

		if err != nil {
			// Update checkpoint with failure
			cp.Status = checkpoint.StatusFailed
			cp.UpdatedAt = time.Now()
			fc.orchestrator.checkpoint.Save(cp)

			fc.sendProgress(
				fmt.Sprintf("[%s] %s failed after %v: %v", formatTime(time.Now()), wsID, wsDuration.Round(time.Second), err),
				wsID,
				"failed",
			)

			return fmt.Errorf("%w: %s: %v", ErrExecutionFailed, wsID, err)
		}

		// Mark as completed
		cp.CompletedWorkstreams = append(cp.CompletedWorkstreams, wsID)
		cp.UpdatedAt = time.Now()

		// Save checkpoint after successful execution
		if err := fc.orchestrator.checkpoint.Save(cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}

		// Send progress update
		fc.sendProgress(
			fmt.Sprintf("[%s] %s complete (%v)", formatTime(time.Now()), wsID, wsDuration.Round(time.Second)),
			wsID,
			"completed",
		)
	}

	// Mark checkpoint as completed
	cp.Status = checkpoint.StatusCompleted
	cp.CurrentWorkstream = ""
	cp.UpdatedAt = time.Now()

	if err := fc.orchestrator.checkpoint.Save(cp); err != nil {
		return fmt.Errorf("failed to save final checkpoint: %w", err)
	}

	totalDuration := time.Since(startTime)

	// Send completion summary
	fc.sendProgress(
		fmt.Sprintf("[%s] Feature execution complete: %d/%d workstreams, %v total",
			formatTime(time.Now()), len(order), len(order), totalDuration.Round(time.Minute)),
		"",
		"completed",
	)

	return nil
}

// executeWithRetry executes a workstream with retry logic
func (fc *FeatureCoordinator) executeWithRetry(wsID string) error {
	maxRetries := fc.orchestrator.maxRetries

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := fc.orchestrator.executor.Execute(wsID)
		if err == nil {
			return nil
		}

		// If this is not the last attempt, log retry
		if attempt < maxRetries {
			fc.sendProgress(
				fmt.Sprintf("[%s] %s failed (attempt %d/%d), retrying...",
					formatTime(time.Now()), wsID, attempt+1, maxRetries+1),
				wsID,
				"retrying",
			)
		}

		// Wait a bit before retry (could be configurable)
		if attempt < maxRetries {
			time.Sleep(time.Second * 2)
		}
	}

	return fmt.Errorf("failed after %d retries", maxRetries)
}

// ResumeFeature resumes execution from a checkpoint
func (fc *FeatureCoordinator) ResumeFeature(checkpointID string) error {
	startTime := time.Now()

	fc.sendProgress(
		fmt.Sprintf("[%s] Resuming feature execution from checkpoint: %s", formatTime(startTime), checkpointID),
		"",
		"resuming",
	)

	// Load checkpoint
	cp, err := fc.orchestrator.checkpoint.Resume(checkpointID)
	if err != nil {
		return fmt.Errorf("failed to resume checkpoint: %w", err)
	}

	// If already completed, nothing to do
	if cp.Status == checkpoint.StatusCompleted {
		fc.sendProgress(
			fmt.Sprintf("[%s] Checkpoint already completed", formatTime(time.Now())),
			"",
			"already_completed",
		)
		return nil
	}

	// Load workstreams
	workstreams, err := fc.orchestrator.loader.LoadWorkstreams(cp.FeatureID)
	if err != nil {
		return fmt.Errorf("failed to load workstreams: %w", err)
	}

	// Build dependency graph
	graph, err := BuildDependencyGraph(workstreams)
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Get full execution order
	fullOrder, err := TopologicalSort(graph)
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	// Get remaining workstreams
	remaining, err := fc.getRemainingWorkstreams(fullOrder, &cp)
	if err != nil {
		return fmt.Errorf("failed to determine remaining workstreams: %w", err)
	}

	fc.sendProgress(
		fmt.Sprintf("[%s] Resuming from %s, %d workstreams remaining",
			formatTime(time.Now()), cp.CurrentWorkstream, len(remaining)),
		"",
		"resuming",
	)

	// Execute remaining workstreams
	for i, wsID := range remaining {
		cp.CurrentWorkstream = wsID
		cp.Status = checkpoint.StatusInProgress
		cp.UpdatedAt = time.Now()

		// Save checkpoint before execution
		if err := fc.orchestrator.checkpoint.Save(cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}

		// Send progress update
		fc.sendProgress(
			fmt.Sprintf("[%s] Executing %s (%d/%d remaining)...",
				formatTime(time.Now()), wsID, i+1, len(remaining)),
			wsID,
			"executing",
		)

		// Execute with retry
		wsStartTime := time.Now()
		err := fc.executeWithRetry(wsID)
		wsDuration := time.Since(wsStartTime)

		if err != nil {
			// Update checkpoint with failure
			cp.Status = checkpoint.StatusFailed
			cp.UpdatedAt = time.Now()
			fc.orchestrator.checkpoint.Save(cp)

			fc.sendProgress(
				fmt.Sprintf("[%s] %s failed after %v: %v",
					formatTime(time.Now()), wsID, wsDuration.Round(time.Second), err),
				wsID,
				"failed",
			)

			return fmt.Errorf("%w: %s: %v", ErrExecutionFailed, wsID, err)
		}

		// Mark as completed
		cp.CompletedWorkstreams = append(cp.CompletedWorkstreams, wsID)
		cp.UpdatedAt = time.Now()

		// Save checkpoint after successful execution
		if err := fc.orchestrator.checkpoint.Save(cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}

		// Send progress update
		fc.sendProgress(
			fmt.Sprintf("[%s] %s complete (%v)", formatTime(time.Now()), wsID, wsDuration.Round(time.Second)),
			wsID,
			"completed",
		)
	}

	// Mark checkpoint as completed
	cp.Status = checkpoint.StatusCompleted
	cp.CurrentWorkstream = ""
	cp.UpdatedAt = time.Now()

	if err := fc.orchestrator.checkpoint.Save(cp); err != nil {
		return fmt.Errorf("failed to save final checkpoint: %w", err)
	}

	totalDuration := time.Since(startTime)

	// Send completion summary
	fc.sendProgress(
		fmt.Sprintf("[%s] Feature resume complete: %d/%d workstreams, %v total",
			formatTime(time.Now()), len(cp.CompletedWorkstreams), len(fullOrder), totalDuration.Round(time.Minute)),
		"",
		"completed",
	)

	return nil
}

// getRemainingWorkstreams determines which workstreams still need to execute
func (fc *FeatureCoordinator) getRemainingWorkstreams(allWorkstreams []string, cp *checkpoint.Checkpoint) ([]string, error) {
	// If already completed, nothing to do
	if cp.Status == checkpoint.StatusCompleted {
		return []string{}, nil
	}

	// Find position of current workstream
	startIndex := 0
	if cp.CurrentWorkstream != "" {
		for i, wsID := range allWorkstreams {
			if wsID == cp.CurrentWorkstream {
				startIndex = i
				break
			}
		}
	}

	// Return remaining workstreams (including current if not completed)
	remaining := allWorkstreams[startIndex:]

	// Filter out already completed workstreams
	var result []string
	for _, wsID := range remaining {
		completed := false
		for _, completedID := range cp.CompletedWorkstreams {
			if wsID == completedID {
				completed = true
				break
			}
		}
		if !completed {
			result = append(result, wsID)
		}
	}

	return result, nil
}
