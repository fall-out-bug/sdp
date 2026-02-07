package graph

import (
	"fmt"
	"log"
)

// SetCheckpointDir sets the checkpoint directory (for testing)
func (d *Dispatcher) SetCheckpointDir(dir string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.checkpointManager != nil {
		d.checkpointManager.SetCheckpointDir(dir)
	}
}

// SetFeatureID sets the feature ID for checkpointing
func (d *Dispatcher) SetFeatureID(featureID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.featureID = featureID
	if d.enableCheckpoint && d.checkpointManager == nil {
		d.checkpointManager = NewCheckpointManager(featureID)
	}
}

// createCheckpoint creates a checkpoint from the current dispatcher state
func (d *Dispatcher) createCheckpoint() *Checkpoint {
	if d.checkpointManager == nil {
		return nil
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	// Convert completed map to slice
	completed := make([]string, 0, len(d.completed))
	for id := range d.completed {
		completed = append(completed, id)
	}

	// Convert failed map to slice
	failed := make([]string, 0, len(d.failed))
	for id := range d.failed {
		failed = append(failed, id)
	}

	checkpoint := d.checkpointManager.CreateCheckpoint(d.graph, d.featureID, completed, failed)

	// Add circuit breaker snapshot
	cbMetrics := d.circuitBreaker.Metrics()
	checkpoint.CircuitBreaker = &CircuitBreakerSnapshot{
		State:            int(cbMetrics.State),
		FailureCount:     cbMetrics.FailureCount,
		SuccessCount:     cbMetrics.SuccessCount,
		ConsecutiveOpens: cbMetrics.ConsecutiveOpens,
		LastFailureTime:  cbMetrics.LastFailureTime,
		LastStateChange:  cbMetrics.LastStateChange,
	}

	return checkpoint
}

// restoreFromCheckpoint restores the dispatcher state from a checkpoint
func (d *Dispatcher) restoreFromCheckpoint(checkpoint *Checkpoint) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Restore completed workstreams
	d.completed = make(map[string]bool)
	for _, id := range checkpoint.Completed {
		d.completed[id] = true
	}

	// Restore failed workstreams
	d.failed = make(map[string]error)
	for _, id := range checkpoint.Failed {
		d.failed[id] = fmt.Errorf("restored from failed state")
	}

	// Restore circuit breaker state
	cbState := -1 // Default to unknown
	if checkpoint.CircuitBreaker != nil {
		d.circuitBreaker.Restore(checkpoint.CircuitBreaker)
		cbState = checkpoint.CircuitBreaker.State
	}

	log.Printf("[Checkpoint] Restored state: %d completed, %d failed, CB state=%d",
		len(checkpoint.Completed), len(checkpoint.Failed), cbState)
}

// tryRestoreCheckpoint attempts to restore from checkpoint if enabled
func (d *Dispatcher) tryRestoreCheckpoint() {
	if !d.enableCheckpoint || d.checkpointManager == nil {
		return
	}

	checkpoint, err := d.checkpointManager.Load()
	if err != nil {
		log.Printf("[Checkpoint] Failed to load checkpoint: %v", err)
		return
	}

	if checkpoint != nil {
		// Verify feature ID matches
		if checkpoint.FeatureID != d.featureID {
			log.Printf("[Checkpoint] Feature ID mismatch: expected %s, got %s",
				d.featureID, checkpoint.FeatureID)
			return
		}

		// Restore state
		d.restoreFromCheckpoint(checkpoint)
		log.Printf("[Checkpoint] Successfully restored from checkpoint")
	}
}
