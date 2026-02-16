package runtime

import (
	"context"
	"errors"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota // Normal operation
	StateOpen                       // Failing, reject all calls
	StateHalfOpen                   // Testing if recovered
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitConfig configures circuit breaker behavior
type CircuitConfig struct {
	FailureThreshold int           // Failures before opening
	SuccessThreshold int           // Successes before closing from half-open
	Timeout          time.Duration // Time before trying half-open
}

// DefaultCircuitConfig provides sensible defaults
var DefaultCircuitConfig = CircuitConfig{
	FailureThreshold: 5,
	SuccessThreshold: 3,
	Timeout:          30 * time.Second,
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config CircuitConfig

	mu           sync.RWMutex
	state        CircuitState
	failures     int
	successes    int
	lastFailTime time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config CircuitConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute runs fn through the circuit breaker
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if !cb.allowRequest() {
		return ErrCircuitOpen
	}

	err := fn()
	cb.recordResult(err)
	return err
}

// ErrCircuitOpen is returned when circuit is open
var ErrCircuitOpen = errors.New("circuit breaker is open")

// allowRequest checks if a request should be allowed
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if timeout has passed
		if time.Since(cb.lastFailTime) > cb.config.Timeout {
			cb.state = StateHalfOpen
			cb.successes = 0
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// recordResult records the result of a request
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err == nil {
		cb.onSuccess()
	} else {
		cb.onFailure()
	}
}

func (cb *CircuitBreaker) onSuccess() {
	cb.failures = 0

	switch cb.state {
	case StateClosed:
		// Already closed, nothing to do
	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.state = StateClosed
			cb.successes = 0
		}
	}
}

func (cb *CircuitBreaker) onFailure() {
	cb.failures++
	cb.lastFailTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.config.FailureThreshold {
			cb.state = StateOpen
		}
	case StateHalfOpen:
		cb.state = StateOpen
		cb.successes = 0
	}
}

// State returns current circuit state
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
}
