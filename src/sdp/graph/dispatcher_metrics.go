package graph

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
