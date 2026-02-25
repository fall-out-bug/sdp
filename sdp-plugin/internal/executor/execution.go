package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"
)

// Execute runs workstreams according to the provided options.
//
//nolint:gocognit,gocyclo // orchestration with many branches by design
func (e *Executor) Execute(ctx context.Context, output io.Writer, opts ExecuteOptions) (*ExecutionResult, error) {
	startTime := time.Now()

	// Set output format
	e.SetOutputFormat(opts.Output)

	result := &ExecutionResult{
		EvidenceEvents: []EvidenceEvent{},
	}

	var workstreams []string
	var err error
	if opts.SpecificWS != "" {
		workstreams = []string{opts.SpecificWS}
	} else {
		workstreams, err = e.findReadyWorkstreams()
		if err != nil {
			return nil, fmt.Errorf("find ready workstreams: %w", err)
		}
	}

	result.TotalWorkstreams = len(workstreams)

	if e.config.DryRun {
		return e.executeDryRun(ctx, output, workstreams, result)
	}

	// Parse dependencies for all workstreams
	dependencies := make(map[string][]string)
	for _, wsID := range workstreams {
		deps, err := e.ParseDependencies(wsID)
		if err != nil {
			if err := writeFmt(output, "Warning: failed to parse dependencies for %s: %v\n", wsID, err); err != nil {
				return nil, fmt.Errorf("write: %w", err)
			}
			// Safe fallback: treat parse failure as no dependencies (not skip)
			dependencies[wsID] = []string{}
			continue
		}
		dependencies[wsID] = deps
	}

	// Sort workstreams topologically by dependencies
	sorted, err := e.TopologicalSort(workstreams, dependencies)
	if err != nil {
		return nil, fmt.Errorf("topological sort failed: %w", err)
	}

	// Execute each workstream
	for _, wsID := range sorted {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		if opts.Output != "json" {
			if err := writeFmt(output, "\nExecuting: %s\n", wsID); err != nil {
				return nil, fmt.Errorf("write: %w", err)
			}
		}

		// Execute workstream with retry logic
		retryCount, err := e.executeWorkstreamWithRetry(ctx, output, wsID, opts.Retry)
		result.Executed++
		result.Retries += retryCount

		if err != nil {
			result.Failed++
			if err := writeLine(output, e.progress.RenderError(wsID, err)); err != nil {
				return nil, fmt.Errorf("write: %w", err)
			}
			// Propagate context cancellation so caller can stop
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return result, err
			}
		} else {
			result.Succeeded++
		}

		if retryCount > 0 && opts.Output != "json" {
			if err := writeFmt(output, "Retry: %s retried %d time(s)\n", wsID, retryCount); err != nil {
				return nil, fmt.Errorf("write: %w", err)
			}
		}

		// Emit evidence events
		events := e.generateEvidenceEvents(wsID)
		result.EvidenceEvents = append(result.EvidenceEvents, events...)
	}

	result.Duration = time.Since(startTime)

	// Render summary
	summary := ExecutionSummary{
		TotalWorkstreams: result.TotalWorkstreams,
		Executed:         result.Executed,
		Succeeded:        result.Succeeded,
		Failed:           result.Failed,
		Skipped:          result.Skipped,
		Retries:          result.Retries,
		Duration:         result.Duration.Seconds(),
	}
	if err := writeLine(output, e.progress.RenderSummary(summary)); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	return result, nil
}
