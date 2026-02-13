package orchestrator

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
)

// executeSequential executes workstreams sequentially with retry
func (o *Orchestrator) executeSequential(workstreams []string) error {
	for _, wsID := range workstreams {
		err := o.executeWithRetry(wsID)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrExecutionFailed, wsID, err)
		}
	}
	return nil
}

// executeWithRetry executes a workstream with retry logic
func (o *Orchestrator) executeWithRetry(wsID string) error {
	var lastErr error
	for attempt := 0; attempt <= o.maxRetries; attempt++ {
		err := o.executor.Execute(wsID)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("failed after %d retries: %w", o.maxRetries, lastErr)
}

// executeWithCheckpoint executes workstreams with checkpointing
func (o *Orchestrator) executeWithCheckpoint(workstreams []string, cp *checkpoint.Checkpoint) error {
	for _, wsID := range workstreams {
		// Update current workstream
		cp.CurrentWorkstream = wsID
		cp.Status = checkpoint.StatusInProgress

		// Save checkpoint before execution
		if err := o.checkpoint.Save(*cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}

		// Execute workstream
		err := o.executeWithRetry(wsID)
		if err != nil {
			// Update checkpoint with failure
			cp.Status = checkpoint.StatusFailed
			o.checkpoint.Save(*cp)
			return err
		}

		// Mark as completed
		cp.CompletedWorkstreams = append(cp.CompletedWorkstreams, wsID)

		// Save checkpoint after successful execution
		if err := o.checkpoint.Save(*cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}
	}

	// Mark checkpoint as completed
	cp.Status = checkpoint.StatusCompleted
	cp.CurrentWorkstream = ""
	if err := o.checkpoint.Save(*cp); err != nil {
		return fmt.Errorf("failed to save final checkpoint: %w", err)
	}

	return nil
}

// getRemainingWorkstreams determines which workstreams still need to execute
func (o *Orchestrator) getRemainingWorkstreams(allWorkstreams []string, cp *checkpoint.Checkpoint) ([]string, error) {
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

	// Return remaining workstreams
	return allWorkstreams[startIndex:], nil
}

// Run executes all workstreams for a feature
func (o *Orchestrator) Run(featureID string) error {
	// Load workstreams
	workstreams, err := o.loader.LoadWorkstreams(featureID)
	if err != nil {
		return fmt.Errorf("failed to load workstreams: %w", err)
	}

	if len(workstreams) == 0 {
		return fmt.Errorf("%w: feature %s has no workstreams", ErrFeatureNotFound, featureID)
	}

	// Build dependency graph
	graph, err := BuildDependencyGraph(workstreams)
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Get execution order
	order, err := TopologicalSort(graph)
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	// Create initial checkpoint
	cp := &checkpoint.Checkpoint{
		ID:                   featureID,
		FeatureID:            featureID,
		Status:               checkpoint.StatusPending,
		CompletedWorkstreams: []string{},
		CurrentWorkstream:    "",
	}

	// Execute with checkpointing
	return o.executeWithCheckpoint(order, cp)
}

// Resume resumes execution from a checkpoint
func (o *Orchestrator) Resume(checkpointID string) error {
	// Load checkpoint
	cp, err := o.checkpoint.Resume(checkpointID)
	if err != nil {
		return fmt.Errorf("failed to resume checkpoint: %w", err)
	}

	// If already completed, nothing to do
	if cp.Status == checkpoint.StatusCompleted {
		return nil
	}

	// Load workstreams
	workstreams, err := o.loader.LoadWorkstreams(cp.FeatureID)
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
	remaining, err := o.getRemainingWorkstreams(fullOrder, &cp)
	if err != nil {
		return fmt.Errorf("failed to determine remaining workstreams: %w", err)
	}

	// Execute remaining
	return o.executeWithCheckpoint(remaining, &cp)
}
