package runtime

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_Closed(t *testing.T) {
	cb := NewCircuitBreaker(CircuitConfig{
		FailureThreshold: 3,
		Timeout:          100 * time.Millisecond,
	})

	// Should be closed initially
	if cb.State() != StateClosed {
		t.Errorf("Expected closed state, got %v", cb.State())
	}

	// Execute successful calls
	for i := 0; i < 5; i++ {
		err := cb.Execute(nil, func() error { return nil })
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	// Should still be closed
	if cb.State() != StateClosed {
		t.Errorf("Expected closed state after successes, got %v", cb.State())
	}
}

func TestCircuitBreaker_Open(t *testing.T) {
	cb := NewCircuitBreaker(CircuitConfig{
		FailureThreshold: 3,
		Timeout:          100 * time.Millisecond,
	})

	// Trigger failures
	for i := 0; i < 3; i++ {
		cb.Execute(nil, func() error { return errors.New("fail") })
	}

	if cb.State() != StateOpen {
		t.Errorf("Expected open state after failures, got %v", cb.State())
	}

	// Should reject calls when open
	err := cb.Execute(nil, func() error { return nil })
	if err != ErrCircuitOpen {
		t.Errorf("Expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitConfig{
		FailureThreshold: 2,
		SuccessThreshold: 2,
		Timeout:          50 * time.Millisecond,
	})

	// Trigger open
	for i := 0; i < 2; i++ {
		cb.Execute(nil, func() error { return errors.New("fail") })
	}

	if cb.State() != StateOpen {
		t.Fatalf("Expected open state, got %v", cb.State())
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Should transition to half-open on next request
	err := cb.Execute(nil, func() error { return nil })
	if err != nil {
		t.Errorf("Unexpected error in half-open: %v", err)
	}

	// Still half-open after 1 success
	if cb.State() != StateHalfOpen {
		t.Errorf("Expected half-open state, got %v", cb.State())
	}

	// Another success should close
	cb.Execute(nil, func() error { return nil })
	if cb.State() != StateClosed {
		t.Errorf("Expected closed state after success threshold, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenToOpen(t *testing.T) {
	cb := NewCircuitBreaker(CircuitConfig{
		FailureThreshold: 2,
		SuccessThreshold: 2,
		Timeout:          50 * time.Millisecond,
	})

	// Trigger open
	for i := 0; i < 2; i++ {
		cb.Execute(nil, func() error { return errors.New("fail") })
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	// Should transition to half-open
	cb.Execute(nil, func() error { return nil })

	// Failure in half-open should go back to open
	cb.Execute(nil, func() error { return errors.New("fail") })

	if cb.State() != StateOpen {
		t.Errorf("Expected open state after failure in half-open, got %v", cb.State())
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	cb := NewCircuitBreaker(CircuitConfig{
		FailureThreshold: 2,
		Timeout:          100 * time.Millisecond,
	})

	// Trigger open
	for i := 0; i < 2; i++ {
		cb.Execute(nil, func() error { return errors.New("fail") })
	}

	if cb.State() != StateOpen {
		t.Fatalf("Expected open state, got %v", cb.State())
	}

	// Reset
	cb.Reset()

	if cb.State() != StateClosed {
		t.Errorf("Expected closed state after reset, got %v", cb.State())
	}
}

func TestCircuitBreaker_Concurrent(t *testing.T) {
	cb := NewCircuitBreaker(CircuitConfig{
		FailureThreshold: 100,
		Timeout:          100 * time.Millisecond,
	})

	done := make(chan bool)

	// Launch concurrent operations
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				cb.Execute(nil, func() error { return nil })
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should still be closed
	if cb.State() != StateClosed {
		t.Errorf("Expected closed state after concurrent operations, got %v", cb.State())
	}
}

func TestCircuitState_String(t *testing.T) {
	tests := []struct {
		state    CircuitState
		expected string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{CircuitState(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expected {
			t.Errorf("State %d: expected %s, got %s", tt.state, tt.expected, got)
		}
	}
}
