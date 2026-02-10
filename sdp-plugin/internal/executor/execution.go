package executor

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Execute runs workstreams according to the provided options
func (e *Executor) Execute(ctx context.Context, output io.Writer, opts ExecuteOptions) (*ExecutionResult, error) {
	startTime := time.Now()

	// Set output format
	e.SetOutputFormat(opts.Output)

	result := &ExecutionResult{
		EvidenceEvents: []EvidenceEvent{},
	}

	// Find workstreams to execute
	var workstreams []string
	var err error

	if opts.SpecificWS != "" {
		// Execute specific workstream
		workstreams = []string{opts.SpecificWS}
	} else if opts.All {
		// Execute all ready workstreams
		workstreams, err = e.findReadyWorkstreams()
		if err != nil {
			return nil, fmt.Errorf("find ready workstreams: %w", err)
		}
	} else {
		// Default: all ready workstreams
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
			fmt.Fprintf(output, "Warning: failed to parse dependencies for %s: %v\n", wsID, err)
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

		// Show executing message (only in human mode)
		if opts.Output != "json" {
			fmt.Fprintf(output, "\nExecuting: %s\n", wsID)
		}

		// Execute workstream with retry logic
		retryCount, err := e.executeWorkstreamWithRetry(ctx, output, wsID, opts.Retry)
		result.Executed++
		result.Retries += retryCount

		if err != nil {
			result.Failed++
			fmt.Fprintln(output, e.progress.RenderError(wsID, err))
		} else {
			result.Succeeded++
		}

		// Show retry message if retries occurred
		if retryCount > 0 && opts.Output != "json" {
			fmt.Fprintf(output, "Retry: %s retried %d time(s)\n", wsID, retryCount)
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
	fmt.Fprintln(output, e.progress.RenderSummary(summary))

	return result, nil
}

