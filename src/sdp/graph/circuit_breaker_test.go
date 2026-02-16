package graph

import (
	"errors"
	"testing"
	"time"
)

// TestNewCircuitBreaker verifies circuit breaker creation
func TestNewCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		Threshold:  5,
		Window:     10,
		Timeout:    30 * time.Second,
		MaxBackoff: 3 * time.Minute,
	}
	cb := NewCircuitBreaker(config)

	if cb == nil {
		t.Fatal("NewCircuitBreaker returned nil")
	}
	if cb.State() != StateClosed {
		t.Errorf("expected initial state CLOSED, got %v", cb.State())
	}
}

// TestNewCircuitBreaker_Defaults verifies default values
func TestNewCircuitBreaker_Defaults(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{})

	if cb.threshold != 3 {
		t.Errorf("expected default threshold 3, got %d", cb.threshold)
	}
	if cb.window != 5 {
		t.Errorf("expected default window 5, got %d", cb.window)
	}
	if cb.timeout != 60*time.Second {
		t.Errorf("expected default timeout 60s, got %v", cb.timeout)
	}
	if cb.maxBackoff != 5*time.Minute {
		t.Errorf("expected default maxBackoff 5m, got %v", cb.maxBackoff)
	}
}

// TestCircuitBreaker_Execute_Success verifies successful execution
func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{Threshold: 3})

	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
	if cb.State() != StateClosed {
		t.Errorf("expected state CLOSED, got %v", cb.State())
	}

	metrics := cb.Metrics()
	if metrics.SuccessCount != 1 {
		t.Errorf("expected success count 1, got %d", metrics.SuccessCount)
	}
}

// TestCircuitBreaker_Execute_Failure verifies failed execution
func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{Threshold: 3})
	testErr := errors.New("test error")

	err := cb.Execute(func() error {
		return testErr
	})

	if err == nil {
		t.Error("expected error")
	}
	if !errors.Is(err, testErr) {
		t.Errorf("expected test error, got %v", err)
	}

	metrics := cb.Metrics()
	if metrics.FailureCount != 1 {
		t.Errorf("expected failure count 1, got %d", metrics.FailureCount)
	}
}

// TestCircuitBreaker_TripsAfterThreshold verifies circuit trips after threshold
func TestCircuitBreaker_TripsAfterThreshold(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{Threshold: 2})

	// First failure
	cb.Execute(func() error { return errors.New("fail 1") })
	if cb.State() != StateClosed {
		t.Errorf("expected CLOSED after 1 failure, got %v", cb.State())
	}

	// Second failure - should trip
	cb.Execute(func() error { return errors.New("fail 2") })
	if cb.State() != StateOpen {
		t.Errorf("expected OPEN after 2 failures, got %v", cb.State())
	}
}

// TestCircuitBreaker_RejectsWhenOpen verifies rejection when open
func TestCircuitBreaker_RejectsWhenOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Threshold:  1,
		Timeout:    1 * time.Hour, // Long timeout so it stays open
		MaxBackoff: 1 * time.Hour,
	})

	// Trip the breaker
	cb.Execute(func() error { return errors.New("fail") })

	// Try to execute - should be rejected
	executed := false
	err := cb.Execute(func() error {
		executed = true
		return nil
	})

	if !errors.Is(err, ErrCircuitBreakerOpen) {
		t.Errorf("expected ErrCircuitBreakerOpen, got %v", err)
	}
	if executed {
		t.Error("function should not have been executed")
	}
}

// TestCircuitBreaker_TransitionsToHalfOpen verifies transition to half-open
func TestCircuitBreaker_TransitionsToHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Threshold:  1,
		Timeout:    10 * time.Millisecond,
		MaxBackoff: 10 * time.Millisecond,
	})

	// Trip the breaker
	cb.Execute(func() error { return errors.New("fail") })
	if cb.State() != StateOpen {
		t.Fatalf("expected OPEN, got %v", cb.State())
	}

	// Wait for timeout
	time.Sleep(20 * time.Millisecond)

	// Next execution should transition to half-open
	err := cb.Execute(func() error { return nil })
	if err != nil {
		t.Errorf("execution failed: %v", err)
	}
	// After success in half-open, should be closed
	if cb.State() != StateClosed {
		t.Errorf("expected CLOSED after half-open success, got %v", cb.State())
	}
}

// TestCircuitBreaker_SuccessResetsFailures verifies success resets failure count
func TestCircuitBreaker_SuccessResetsFailures(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{Threshold: 3})

	// Two failures
	cb.Execute(func() error { return errors.New("fail 1") })
	cb.Execute(func() error { return errors.New("fail 2") })

	if cb.Metrics().FailureCount != 2 {
		t.Fatalf("expected 2 failures, got %d", cb.Metrics().FailureCount)
	}

	// Success should reset
	cb.Execute(func() error { return nil })

	if cb.Metrics().FailureCount != 0 {
		t.Errorf("expected 0 failures after success, got %d", cb.Metrics().FailureCount)
	}
}

// TestCircuitBreaker_Restore verifies state restoration
func TestCircuitBreaker_Restore(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{Threshold: 3})

	// Create snapshot
	snapshot := &CircuitBreakerSnapshot{
		State:            int(StateOpen),
		FailureCount:     5,
		SuccessCount:     10,
		ConsecutiveOpens: 2,
		LastFailureTime:  time.Now(),
	}

	cb.Restore(snapshot)

	if cb.State() != StateOpen {
		t.Errorf("expected OPEN after restore, got %v", cb.State())
	}

	metrics := cb.Metrics()
	if metrics.FailureCount != 5 {
		t.Errorf("expected failure count 5, got %d", metrics.FailureCount)
	}
	if metrics.SuccessCount != 10 {
		t.Errorf("expected success count 10, got %d", metrics.SuccessCount)
	}
	if metrics.ConsecutiveOpens != 2 {
		t.Errorf("expected consecutive opens 2, got %d", metrics.ConsecutiveOpens)
	}
}

// TestCircuitState_String verifies state string conversion
func TestCircuitState_String(t *testing.T) {
	tests := []struct {
		state    CircuitState
		expected string
	}{
		{StateClosed, "CLOSED"},
		{StateOpen, "OPEN"},
		{StateHalfOpen, "HALF_OPEN"},
		{CircuitState(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expected {
			t.Errorf("state %d: expected %s, got %s", tt.state, tt.expected, got)
		}
	}
}

// TestCircuitBreakerMetrics verifies metrics structure
func TestCircuitBreakerMetrics(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{Threshold: 3})

	cb.Execute(func() error { return nil })
	cb.Execute(func() error { return errors.New("fail") })

	metrics := cb.Metrics()
	if metrics.SuccessCount != 1 {
		t.Errorf("expected success count 1, got %d", metrics.SuccessCount)
	}
	if metrics.FailureCount != 1 {
		t.Errorf("expected failure count 1, got %d", metrics.FailureCount)
	}
	if metrics.State != StateClosed {
		t.Errorf("expected state CLOSED, got %v", metrics.State)
	}
}

// TestCircuitBreakerConfig_Fields verifies config structure
func TestCircuitBreakerConfig_Fields(t *testing.T) {
	config := CircuitBreakerConfig{
		Threshold:  5,
		Window:     10,
		Timeout:    30 * time.Second,
		MaxBackoff: 3 * time.Minute,
	}

	if config.Threshold != 5 {
		t.Error("Threshold not set")
	}
	if config.Window != 10 {
		t.Error("Window not set")
	}
	if config.Timeout != 30*time.Second {
		t.Error("Timeout not set")
	}
	if config.MaxBackoff != 3*time.Minute {
		t.Error("MaxBackoff not set")
	}
}
