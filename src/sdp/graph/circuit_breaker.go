package graph

import (
	"errors"
	"log"
	"sync"
	"time"
)

// ErrCircuitBreakerOpen is returned when circuit breaker is in OPEN state
var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu               sync.RWMutex
	state            CircuitState
	failureCount     int
	successCount     int
	lastFailureTime  time.Time
	lastStateChange  time.Time
	threshold        int
	window           int
	timeout          time.Duration
	maxBackoff       time.Duration
	consecutiveOpens int
}

// NewCircuitBreaker creates a new circuit breaker with the given config
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	// Apply defaults
	if config.Threshold <= 0 {
		config.Threshold = 3
	}
	if config.Window <= 0 {
		config.Window = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 60 * time.Second
	}
	if config.MaxBackoff <= 0 {
		config.MaxBackoff = 5 * time.Minute
	}

	return &CircuitBreaker{
		state:     StateClosed,
		threshold: config.Threshold,
		window:    config.Window,
		timeout:   config.Timeout,
		maxBackoff: config.MaxBackoff,
	}
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Metrics returns the current metrics of the circuit breaker
func (cb *CircuitBreaker) Metrics() CircuitBreakerMetrics {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return CircuitBreakerMetrics{
		State:            cb.state,
		FailureCount:     cb.failureCount,
		SuccessCount:     cb.successCount,
		ConsecutiveOpens: cb.consecutiveOpens,
		LastFailureTime:  cb.lastFailureTime,
		LastStateChange:  cb.lastStateChange,
	}
}

// Execute runs the given function if the circuit breaker is not OPEN
// If the circuit breaker is OPEN, it returns ErrCircuitBreakerOpen without running the function
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	// Check if we're in OPEN state and need to wait for backoff
	if cb.state == StateOpen {
		backoff := cb.calculateBackoff()
		cb.mu.Unlock()

		// Wait for backoff
		log.Printf("[Circuit Breaker] OPEN - waiting %v before attempting recovery", backoff)
		time.Sleep(backoff)

		// Transition to HALF_OPEN
		cb.mu.Lock()
		cb.setState(StateHalfOpen)
		cb.mu.Unlock()
	}

	// Check if we can proceed
	if cb.state == StateOpen {
		cb.mu.Unlock()
		return ErrCircuitBreakerOpen
	}

	cb.mu.Unlock()

	// Execute the function
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// Check if we should trip the breaker
		if cb.failureCount >= cb.threshold {
			cb.consecutiveOpens++
			cb.setState(StateOpen)
		}
	} else {
		cb.successCount++
		// Reset failure count on success
		if cb.failureCount > 0 {
			cb.failureCount--
		}

		// If we're in HALF_OPEN, transition back to CLOSED on success
		if cb.state == StateHalfOpen {
			cb.consecutiveOpens = 0
			cb.setState(StateClosed)
		}
	}

	return err
}
