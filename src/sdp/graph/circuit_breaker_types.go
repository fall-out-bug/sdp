package graph

import "time"

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// String returns the string representation of the state
func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig contains configuration for the circuit breaker
type CircuitBreakerConfig struct {
	Threshold  int           // Number of failures to trip the breaker
	Window     int           // Number of requests to measure over
	Timeout    time.Duration // Time to wait before transitioning from OPEN to HALF_OPEN
	MaxBackoff time.Duration // Maximum backoff time
}

// CircuitBreakerMetrics represents the current metrics of the circuit breaker
type CircuitBreakerMetrics struct {
	State            CircuitState // Current state
	FailureCount     int          // Number of failures in current window
	SuccessCount     int          // Number of successes in current window
	ConsecutiveOpens int          // Number of times the breaker has opened consecutively
	LastFailureTime  time.Time    // Timestamp of the last failure
	LastStateChange  time.Time    // Timestamp of the last state change
}
