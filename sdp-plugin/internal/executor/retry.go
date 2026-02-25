package executor

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
)

// executeWorkstreamWithRetry executes a workstream with retry logic
func (e *Executor) executeWorkstreamWithRetry(ctx context.Context, output io.Writer, wsID string, maxRetries int) (int, error) {
	var lastErr error
	retries := 0

	if maxRetries <= 0 {
		maxRetries = e.config.RetryCount
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check context cancellation in loop body (not only at start)
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
		slog.Debug("workstream run failed", "ws_id", wsID, "attempt", attempt, "max_retries", maxRetries, "error", err)

		if attempt < maxRetries {
			if errW := writeLine(output, e.progress.Output(wsID, progress, "retrying", fmt.Sprintf("failed: %v", err))); errW != nil {
				return retries, fmt.Errorf("write: %w", errW)
			}
			delay := e.retryDelayFromConfigCached()
			select {
			case <-ctx.Done():
				return retries, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	slog.Info("workstream execution failed after retries", "ws_id", wsID, "retries", retries, "error", lastErr)
	return retries, lastErr
}

// retryDelayFromConfigCached returns retry delay from config, cached to avoid repeated config I/O per retry.
func (e *Executor) retryDelayFromConfigCached() time.Duration {
	e.cachedRetryDelayOnce.Do(func() {
		e.cachedRetryDelay = retryDelayFromConfig()
	})
	return e.cachedRetryDelay
}

// retryDelayFromConfig returns retry delay from config (or env, or default).
func retryDelayFromConfig() time.Duration {
	root, err := config.FindProjectRoot()
	if err != nil {
		return config.TimeoutFromEnv("SDP_TIMEOUT_RETRY_DELAY", 100*time.Millisecond)
	}
	cfg, err := config.Load(root)
	if err != nil || cfg == nil {
		return config.TimeoutFromEnv("SDP_TIMEOUT_RETRY_DELAY", 100*time.Millisecond)
	}
	return config.TimeoutFromConfigOrEnv(cfg.Timeouts.RetryDelay, "SDP_TIMEOUT_RETRY_DELAY", 100*time.Millisecond)
}
