package graph

import (
	"log"
	"time"
)

// setState transitions the circuit breaker to a new state
func (cb *CircuitBreaker) setState(newState CircuitState) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = time.Now()

	log.Printf("[Circuit Breaker] State transition: %s → %s (failures: %d, consecutive opens: %d)",
		oldState, newState, cb.failureCount, cb.consecutiveOpens)
}

// calculateBackoff calculates the exponential backoff duration
func (cb *CircuitBreaker) calculateBackoff() time.Duration {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	// Base timeout
	baseTimeout := cb.timeout

	// Exponential backoff: base * 2^(consecutiveOpens - 1)
	backoff := baseTimeout * time.Duration(1<<uint(cb.consecutiveOpens-1))

	// Cap at max backoff
	if backoff > cb.maxBackoff {
		backoff = cb.maxBackoff
	}

	return backoff
}

// Restore restores the circuit breaker state from a snapshot
func (cb *CircuitBreaker) Restore(snapshot *CircuitBreakerSnapshot) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = CircuitState(snapshot.State)
	cb.failureCount = snapshot.FailureCount
	cb.successCount = snapshot.SuccessCount
	cb.consecutiveOpens = snapshot.ConsecutiveOpens
	cb.lastFailureTime = snapshot.LastFailureTime

	log.Printf("[Circuit Breaker] Restored state: %d (state=%s, failures=%d, consecutive opens=%d)",
		snapshot.State, CircuitState(snapshot.State), snapshot.FailureCount, snapshot.ConsecutiveOpens)
}
