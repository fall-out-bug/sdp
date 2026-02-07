package graph

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Dispatcher coordinates parallel execution of workstreams
type Dispatcher struct {
	graph             *DependencyGraph
	concurrency       int
	completed         map[string]bool
	failed            map[string]error
	circuitBreaker    *CircuitBreaker
	checkpointManager *CheckpointManager
	featureID         string
	enableCheckpoint  bool
	mu                sync.RWMutex
}

// NewDispatcher creates a new dispatcher for parallel execution
func NewDispatcher(g *DependencyGraph, concurrency int) *Dispatcher {
	if concurrency < 1 {
		concurrency = 3 // Default to 3 parallel workers
	}
	if concurrency > 5 {
		concurrency = 5 // Max 5 parallel workers
	}

	// Create circuit breaker with default configuration
	cbConfig := CircuitBreakerConfig{
		Threshold:  3,                // Trip after 3 failures
		Window:     5,                // Within 5 requests
		Timeout:    60 * time.Second, // Wait 60s before retry
		MaxBackoff: 5 * time.Minute,  // Max backoff 5min
	}

	return &Dispatcher{
		graph:             g,
		concurrency:       concurrency,
		completed:         make(map[string]bool),
		failed:            make(map[string]error),
		circuitBreaker:    NewCircuitBreaker(cbConfig),
		checkpointManager: nil,
		featureID:         "",
		enableCheckpoint:  false,
	}
}

// ExecuteResult represents the result of executing a workstream
type ExecuteResult struct {
	WorkstreamID string
	Success      bool
	Error        error
	Duration     int64 // Duration in milliseconds
}

// ExecuteFunc is a function that executes a single workstream
type ExecuteFunc func(wsID string) error

// Execute runs all workstreams in parallel respecting dependencies
func (d *Dispatcher) Execute(executeFn ExecuteFunc) []ExecuteResult {
	// Try to restore from checkpoint if enabled
	d.tryRestoreCheckpoint()

	results := []ExecuteResult{}
	totalNodes := len(d.graph.nodes)

	// Continue until all nodes are processed
	for len(d.completed)+len(d.failed) < totalNodes {
		// Get ready nodes
		ready := d.graph.GetReady()

		// Filter out already completed nodes
		readyToRun := []string{}
		for _, id := range ready {
			if !d.isCompleted(id) {
				readyToRun = append(readyToRun, id)
			}
		}

		// If no nodes are ready but we haven't finished, we might have a problem
		if len(readyToRun) == 0 && len(d.completed)+len(d.failed) < totalNodes {
			// This shouldn't happen if the graph is valid
			// Check if we're just waiting on nodes already in flight
			continue
		}

		// Execute ready nodes in parallel
		batchSize := len(readyToRun)
		if batchSize > d.concurrency {
			batchSize = d.concurrency
		}

		// Process batch
		var wg sync.WaitGroup
		resultsChan := make(chan ExecuteResult, batchSize)
		for i := 0; i < batchSize && i < len(readyToRun); i++ {
			wg.Add(1)
			go func(wsID string) {
				defer wg.Done()
				// Wrap execution with circuit breaker
				err := d.circuitBreaker.Execute(func() error {
					return executeFn(wsID)
				})
				// Log circuit breaker state for observability
				metrics := d.circuitBreaker.Metrics()
				if err != nil && err == ErrCircuitBreakerOpen {
					log.Printf("[Circuit Breaker] Workstream %s rejected - circuit is OPEN (state=%v, failures=%d)",
						wsID, metrics.State, metrics.FailureCount)
				} else if err != nil {
					log.Printf("[Circuit Breaker] Workstream %s failed - circuit state=%v, failures=%d",
						wsID, metrics.State, metrics.FailureCount)
				}
				result := ExecuteResult{
					WorkstreamID: wsID,
					Success:      err == nil,
					Error:        err,
				}
				resultsChan <- result
				// Update graph state
				d.mu.Lock()
				if err != nil {
					d.failed[wsID] = err
					// Mark as complete in graph so dependents can run
					// (even though execution failed, we want to continue with others)
					d.graph.MarkComplete(wsID)
				} else {
					d.completed[wsID] = true
					d.graph.MarkComplete(wsID)
				}
				d.mu.Unlock()
			}(readyToRun[i])
		}
		// Wait for all goroutines in this batch
		wg.Wait()
		close(resultsChan)

		// Collect results
		for result := range resultsChan {
			results = append(results, result)
		}

		// Save checkpoint after each batch if enabled
		if d.enableCheckpoint && d.checkpointManager != nil {
			checkpoint := d.createCheckpoint()
			if err := d.checkpointManager.Save(checkpoint); err != nil {
				log.Printf("[Checkpoint] Failed to save checkpoint: %v", err)
			}
		}
	}

	// Delete checkpoint on successful completion if enabled
	if d.enableCheckpoint && d.checkpointManager != nil && len(d.failed) == 0 {
		if err := d.checkpointManager.Delete(); err != nil {
			log.Printf("[Checkpoint] Failed to delete checkpoint: %v", err)
		} else {
			log.Printf("[Checkpoint] Deleted checkpoint after successful completion")
		}
	}

	return results
}

// GetCompleted returns the list of completed workstream IDs
func (d *Dispatcher) GetCompleted() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	completed := []string{}
	for id := range d.completed {
		completed = append(completed, id)
	}
	return completed
}

// GetFailed returns the list of failed workstream IDs and their errors
func (d *Dispatcher) GetFailed() map[string]error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	failed := make(map[string]error)
	for id, err := range d.failed {
		failed[id] = err
	}
	return failed
}

// GetCircuitBreakerMetrics returns the current circuit breaker metrics
func (d *Dispatcher) GetCircuitBreakerMetrics() CircuitBreakerMetrics {
	return d.circuitBreaker.Metrics()
}

// isCompleted checks if a workstream is completed
func (d *Dispatcher) isCompleted(id string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.completed[id]
}

// BuildGraphFromWSFiles creates a dependency graph from workstream files
// This is a placeholder - the actual implementation would parse WS files
func BuildGraphFromWSFiles(workstreams []WorkstreamFile) (*DependencyGraph, error) {
	graph := NewDependencyGraph()

	// First pass: add all nodes
	for _, ws := range workstreams {
		err := graph.AddNode(ws.ID, ws.DependsOn)
		if err != nil {
			return nil, fmt.Errorf("failed to add workstream %s: %w", ws.ID, err)
		}
	}

	return graph, nil
}

// WorkstreamFile represents a workstream file
type WorkstreamFile struct {
	ID        string
	DependsOn []string
}

// NewDispatcherWithCheckpoint creates a new dispatcher with checkpoint support
func NewDispatcherWithCheckpoint(g *DependencyGraph, concurrency int, featureID string, enableCheckpoint bool) *Dispatcher {
	d := NewDispatcher(g, concurrency)
	if enableCheckpoint {
		d.featureID = featureID
		d.enableCheckpoint = true
		d.checkpointManager = NewCheckpointManager(featureID)
	}
	return d
}

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
