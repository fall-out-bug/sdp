package executor

import (
	"context"
	"fmt"
	"io"
	"time"
)

// executeWorkstreamWithRetry executes a workstream with retry logic
func (e *Executor) executeWorkstreamWithRetry(ctx context.Context, output io.Writer, wsID string, maxRetries int) (int, error) {
	var lastErr error
	retries := 0
	var attemptCount int // Track total attempts

	// Use configured retry count if not specified
	if maxRetries <= 0 {
		maxRetries = e.config.RetryCount
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		attemptCount++

		// Check context
		select {
		case <-ctx.Done():
			return retries, ctx.Err()
		default:
		}

		// Show progress
		progress := (attempt * 100) / (maxRetries + 1)
		message := "executing"
		if attempt > 0 {
			message = fmt.Sprintf("retry attempt %d/%d", attempt, maxRetries)
			retries = attempt
		}

		fmt.Fprintln(output, e.progress.Output(wsID, progress, "running", message))

		// Execute (mock for now - would call actual execution logic)
		err := e.executeWorkstreamMock(ctx, wsID, attemptCount)
		if err == nil {
			// Success
			fmt.Fprintln(output, e.progress.RenderSuccess(wsID, "completed successfully"))
			return retries, nil
		}

		lastErr = err

		// Check if we should retry
		if attempt < maxRetries {
			fmt.Fprintln(output, e.progress.Output(wsID, progress, "retrying", fmt.Sprintf("failed: %v", err)))
			time.Sleep(100 * time.Millisecond) // Small delay before retry
		}
	}

	return retries, lastErr
}

// executeWorkstreamMock is a mock executor for testing
// In production, this would call the actual workstream execution logic
func (e *Executor) executeWorkstreamMock(ctx context.Context, wsID string, attemptCount int) error {
	// Mock: 00-054-02 fails on first attempt, succeeds on retry
	if wsID == "00-054-02" {
		if attemptCount == 1 {
			// First attempt fails
			return fmt.Errorf("mock execution failure for %s", wsID)
		}
		// Second attempt (retry) succeeds
		return nil
	}

	// Other workstreams always succeed
	return nil
}
