package runtime

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithRetry_Success(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: 1 * time.Millisecond}
	calls := 0

	err := WithRetry(context.Background(), cfg, func() error {
		calls++
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}
}

func TestWithRetry_RetryableError(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: 1 * time.Millisecond, Multiplier: 1.0}
	calls := 0

	err := WithRetry(context.Background(), cfg, func() error {
		calls++
		if calls < 3 {
			return &RetryableError{Err: errors.New("temporary")}
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error after retries, got %v", err)
	}
	if calls != 3 {
		t.Errorf("Expected 3 calls, got %d", calls)
	}
}

func TestWithRetry_MaxAttempts(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: 1 * time.Millisecond, Multiplier: 1.0}
	calls := 0

	err := WithRetry(context.Background(), cfg, func() error {
		calls++
		return &RetryableError{Err: errors.New("always fails")}
	})

	if err == nil {
		t.Error("Expected error after max attempts")
	}
	if calls != 3 {
		t.Errorf("Expected 3 calls, got %d", calls)
	}
}

func TestWithRetry_NonRetryableError(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: 1 * time.Millisecond}
	calls := 0

	err := WithRetry(context.Background(), cfg, func() error {
		calls++
		return errors.New("permanent error")
	})

	if err == nil {
		t.Error("Expected error")
	}
	if calls != 1 {
		t.Errorf("Expected 1 call (no retry for non-retryable), got %d", calls)
	}
}

func TestWithRetry_ContextCancellation(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 10, InitialDelay: 100 * time.Millisecond}
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := WithRetry(ctx, cfg, func() error {
		calls++
		return &RetryableError{Err: errors.New("fail")}
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestWithFallback_PrimarySucceeds(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, InitialDelay: 1 * time.Millisecond}
	fallbackCalled := false

	err := WithFallback(context.Background(),
		func() error { return nil },
		func() error { fallbackCalled = true; return nil },
		cfg,
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if fallbackCalled {
		t.Error("Fallback should not be called when primary succeeds")
	}
}

func TestWithFallback_FallbackOnFailure(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, InitialDelay: 1 * time.Millisecond}
	fallbackCalled := false

	err := WithFallback(context.Background(),
		func() error { return &RetryableError{Err: errors.New("fail")} },
		func() error { fallbackCalled = true; return nil },
		cfg,
	)

	if err != nil {
		t.Errorf("Expected no error from fallback, got %v", err)
	}
	if !fallbackCalled {
		t.Error("Fallback should be called when primary fails")
	}
}

func TestWithFallbackValue_Success(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, InitialDelay: 1 * time.Millisecond}

	result, err := WithFallbackValue(context.Background(),
		func() (string, error) { return "primary", nil },
		func() (string, error) { return "fallback", nil },
		cfg,
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "primary" {
		t.Errorf("Expected 'primary', got %s", result)
	}
}

func TestWithFallbackValue_Fallback(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 2, InitialDelay: 1 * time.Millisecond}

	result, err := WithFallbackValue(context.Background(),
		func() (string, error) { return "", &RetryableError{Err: errors.New("fail")} },
		func() (string, error) { return "fallback", nil },
		cfg,
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "fallback" {
		t.Errorf("Expected 'fallback', got %s", result)
	}
}

func TestIsRetryable(t *testing.T) {
	if !IsRetryable(&RetryableError{Err: errors.New("test")}) {
		t.Error("RetryableError should be retryable")
	}
	if IsRetryable(errors.New("test")) {
		t.Error("Regular error should not be retryable")
	}
}

func TestAddJitter(t *testing.T) {
	delay := 100 * time.Millisecond

	for i := 0; i < 100; i++ {
		result := addJitter(delay)
		// Jitter is ±25%, so 75ms to 125ms
		if result < 75*time.Millisecond || result > 125*time.Millisecond {
			t.Errorf("Jittered delay %v out of expected range", result)
		}
	}
}
