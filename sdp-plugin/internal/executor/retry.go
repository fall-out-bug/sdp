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

	if maxRetries <= 0 {
		maxRetries = e.config.RetryCount
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
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

		if err := writeLine(output, e.progress.Output(wsID, progress, "running", message)); err != nil {
			return retries, fmt.Errorf("write: %w", err)
		}

		err := e.runner.Run(ctx, wsID)
		if err == nil {
			if err := writeLine(output, e.progress.RenderSuccess(wsID, "completed successfully")); err != nil {
				return retries, fmt.Errorf("write: %w", err)
			}
			return retries, nil
		}

		lastErr = err

		if attempt < maxRetries {
			if errW := writeLine(output, e.progress.Output(wsID, progress, "retrying", fmt.Sprintf("failed: %v", err))); errW != nil {
				return retries, fmt.Errorf("write: %w", errW)
			}
			time.Sleep(100 * time.Millisecond) // Small delay before retry
		}
	}

	return retries, lastErr
}
