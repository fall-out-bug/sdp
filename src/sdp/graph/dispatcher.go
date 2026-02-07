package graph

import (
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

