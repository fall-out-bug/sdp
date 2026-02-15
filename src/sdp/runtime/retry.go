package runtime

import (
	"context"
	"math/rand"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxAttempts  int           // Maximum number of attempts (0 = no retry)
	InitialDelay time.Duration // Initial delay between retries
	MaxDelay     time.Duration // Maximum delay cap
	Multiplier   float64       // Delay multiplier for exponential backoff
	Jitter       bool          // Add random jitter to prevent thundering herd
}

// DefaultRetryConfig provides sensible defaults
var DefaultRetryConfig = RetryConfig{
	MaxAttempts:  3,
	InitialDelay: 100 * time.Millisecond,
	MaxDelay:     10 * time.Second,
	Multiplier:   2.0,
	Jitter:       true,
}

// RetryableError marks an error as retryable
type RetryableError struct {
	Err error
}

func (e *RetryableError) Error() string { return e.Err.Error() }
func (e *RetryableError) Unwrap() error { return e.Err }

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	_, ok := err.(*RetryableError)
	return ok
}

// WithRetry executes fn with retry logic
func WithRetry(ctx context.Context, cfg RetryConfig, fn func() error) error {
	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if attempt > 0 {
			// Wait before retry
			actualDelay := delay
			if cfg.Jitter {
				actualDelay = addJitter(delay)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(actualDelay):
			}

			// Increase delay for next attempt
			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Only retry retryable errors
		if !IsRetryable(err) {
			return err
		}
	}

	return lastErr
}

// WithFallback executes primary function, then fallback on failure
func WithFallback(ctx context.Context, primary, fallback func() error, cfg RetryConfig) error {
	err := WithRetry(ctx, cfg, primary)
	if err != nil {
		if fallback != nil {
			return fallback()
		}
	}
	return err
}

// WithFallbackValue executes primary function returning a value, then fallback on failure
func WithFallbackValue[T any](ctx context.Context, primary func() (T, error), fallback func() (T, error), cfg RetryConfig) (T, error) {
	var zero T

	var lastErr error
	delay := cfg.InitialDelay

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if attempt > 0 {
			actualDelay := delay
			if cfg.Jitter {
				actualDelay = addJitter(delay)
			}

			select {
			case <-ctx.Done():
				return zero, ctx.Err()
			case <-time.After(actualDelay):
			}

			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}

		result, err := primary()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if !IsRetryable(err) {
			break
		}
	}

	// Try fallback
	if fallback != nil {
		return fallback()
	}

	return zero, lastErr
}

// addJitter adds random jitter to delay (±25%)
func addJitter(delay time.Duration) time.Duration {
	// Random value between 0.75 and 1.25
	factor := 0.75 + rand.Float64()*0.5
	return time.Duration(float64(delay) * factor)
}
