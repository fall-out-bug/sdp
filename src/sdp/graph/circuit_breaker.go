package graph

import (
	"errors"
	"sync"
	"time"
)

// CircuitState represents the three states of the circuit breaker
type CircuitState int

const (
	StateClosed   CircuitState = iota // Normal operation
	StateOpen                         // Failing, reject requests
	StateHalfOpen                     // Testing recovery
)

// ErrCircuitBreakerOpen is returned when circuit breaker is in OPEN state
var ErrCircuitBreakerOpen = errors.New("circuit breaker is open")

// CircuitBreakerConfig holds configuration for the circuit breaker
type CircuitBreakerConfig struct {
	Threshold  int           // Default: 3 failures
	Window     int           // Default: 5 requests
	Timeout    time.Duration // Default: 60s
	MaxBackoff time.Duration // Default: 5min
}

// CircuitBreakerMetrics holds current metrics from the circuit breaker
type CircuitBreakerMetrics struct {
	State            CircuitState
	FailureCount     int
	SuccessCount     int
	ConsecutiveOpens int
	LastFailureTime  time.Time
	LastStateChange  time.Time
}

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
		state:           StateClosed,
		threshold:       config.Threshold,
		window:          config.Window,
		timeout:         config.Timeout,
		maxBackoff:      config.MaxBackoff,
		lastStateChange: time.Now(),
	}
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Metrics returns current metrics from the circuit breaker
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

// Execute runs the given function according to circuit breaker state
func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()

	// Check if we should transition from OPEN to HALF_OPEN
	if cb.state == StateOpen {
		backoff := cb.calculateBackoff()
		elapsed := time.Since(cb.lastStateChange)

		if elapsed >= backoff {
			// Transition to HALF_OPEN
			cb.setState(StateHalfOpen)
		} else {
			// Still in timeout period, reject
			cb.mu.Unlock()
			return ErrCircuitBreakerOpen
		}
	}

	// If we're still in OPEN state (shouldn't happen after above check), reject
	if cb.state == StateOpen {
		cb.mu.Unlock()
		return ErrCircuitBreakerOpen
	}

	// Allow execution in CLOSED or HALF_OPEN state
	cb.mu.Unlock()

	// Execute the function
	err := fn()

	// Handle result
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// State-specific failure handling
		switch cb.state {
		case StateClosed:
			if cb.failureCount >= cb.threshold {
				cb.setState(StateOpen)
				cb.consecutiveOpens++
			}
		case StateHalfOpen:
			// Test request failed, back to OPEN
			cb.setState(StateOpen)
			cb.consecutiveOpens++
		}
	} else {
		cb.successCount++

		// State-specific success handling
		switch cb.state {
		case StateClosed:
			// Reset failure count on success
			cb.failureCount = 0
		case StateHalfOpen:
			// Test request succeeded, back to CLOSED
			cb.setState(StateClosed)
			cb.failureCount = 0
			cb.consecutiveOpens = 0
		}
	}

	return err
}

// setState changes the circuit breaker state and records the transition time
func (cb *CircuitBreaker) setState(newState CircuitState) {
	cb.state = newState
	cb.lastStateChange = time.Now()
}

// calculateBackoff computes the backoff duration with exponential increase
func (cb *CircuitBreaker) calculateBackoff() time.Duration {
	// Exponential backoff: baseTimeout * (2 ^ (consecutiveOpens - 1))
	// consecutiveOpens=1 (first open): baseTimeout * 1
	// consecutiveOpens=2 (second open): baseTimeout * 2
	// consecutiveOpens=3 (third open): baseTimeout * 4
	var backoff time.Duration
	if cb.consecutiveOpens <= 1 {
		backoff = cb.timeout
	} else {
		backoff = cb.timeout * (1 << uint(cb.consecutiveOpens-1))
	}

	// Cap at max backoff
	if backoff > cb.maxBackoff {
		backoff = cb.maxBackoff
	}

	return backoff
}
