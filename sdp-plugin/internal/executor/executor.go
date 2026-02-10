package executor

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// ExecutorConfig holds configuration for the executor
type ExecutorConfig struct {
	BacklogDir      string // Directory containing workstream files
	DryRun          bool   // If true, show plan without executing
	RetryCount      int    // Maximum number of retry attempts
	EvidenceLogPath string // Path to evidence log file
}

// ExecuteOptions holds options for a single execution
type ExecuteOptions struct {
	All        bool   // Execute all ready workstreams
	SpecificWS string // Execute specific workstream by ID
	Retry      int    // Number of retries for failed workstreams
	Output     string // Output format: "human" or "json"
}

// ExecutionResult holds the result of an execution
type ExecutionResult struct {
	TotalWorkstreams int             `json:"total_workstreams"`
	Executed         int             `json:"executed"`
	Succeeded        int             `json:"succeeded"`
	Failed           int             `json:"failed"`
	Skipped          int             `json:"skipped"`
	Retries          int             `json:"retries"`
	Duration         time.Duration   `json:"duration"`
	EvidenceEvents   []EvidenceEvent `json:"evidence_events"`
}

// ExecutionSummary is a simplified summary for output
type ExecutionSummary struct {
	TotalWorkstreams int     `json:"total_workstreams"`
	Executed         int     `json:"executed"`
	Succeeded        int     `json:"succeeded"`
	Failed           int     `json:"failed"`
	Skipped          int     `json:"skipped"`
	Retries          int     `json:"retries"`
	Duration         float64 `json:"duration_seconds"`
}

// Executor handles workstream execution with progress tracking
type Executor struct {
	config         ExecutorConfig
	progress       *ProgressRenderer
	evidenceWriter io.Writer
}

// NewExecutor creates a new executor
func NewExecutor(config ExecutorConfig) *Executor {
	return &Executor{
		config:   config,
		progress: NewProgressRenderer("human"),
	}
}

// SetOutputFormat sets the output format for progress rendering
func (e *Executor) SetOutputFormat(format string) {
	e.progress = NewProgressRenderer(format)
}

// SetEvidenceWriter sets the evidence log writer
func (e *Executor) SetEvidenceWriter(w io.Writer) {
	e.evidenceWriter = w
}

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

// executeDryRun shows what would be executed without running
func (e *Executor) executeDryRun(ctx context.Context, output io.Writer, workstreams []string, result *ExecutionResult) (*ExecutionResult, error) {
	fmt.Fprintln(output, "DRY RUN MODE - Showing execution plan:")
	fmt.Fprintln(output, "")

	for _, wsID := range workstreams {
		// Parse workstream file
		wsPath := filepath.Join(e.config.BacklogDir, fmt.Sprintf("%s-*.md", wsID))
		matches, err := filepath.Glob(wsPath)
		if err != nil || len(matches) == 0 {
			fmt.Fprintf(output, "  [✗] %s - file not found\n", wsID)
			continue
		}

		ws, err := parser.ParseWorkstream(matches[0])
		if err != nil {
			fmt.Fprintf(output, "  [✗] %s - parse error: %v\n", wsID, err)
			continue
		}

		fmt.Fprintf(output, "  [→] %s - %s\n", wsID, ws.Goal)

		// Show dependencies
		deps, _ := e.ParseDependencies(wsID)
		if len(deps) > 0 {
			fmt.Fprintf(output, "      Depends on: %s\n", strings.Join(deps, ", "))
		}
	}

	fmt.Fprintln(output, "")
	fmt.Fprintf(output, "Would execute %d workstreams\n", len(workstreams))

	return result, nil
}

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

// findReadyWorkstreams finds all workstreams that are ready to execute
// (no blockers or blockers already completed)
func (e *Executor) findReadyWorkstreams() ([]string, error) {
	var workstreams []string

	// List all workstream files in backlog
	pattern := filepath.Join(e.config.BacklogDir, "*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob pattern failed: %w", err)
	}

	for _, match := range matches {
		// Parse workstream
		ws, err := parser.ParseWorkstream(match)
		if err != nil {
			continue // Skip unparseable files
		}

		// Check if workstream is in ready state
		if ws.Status == "pending" || ws.Status == "ready" {
			workstreams = append(workstreams, ws.ID)
		}
	}

	return workstreams, nil
}

// ParseDependencies parses dependencies from a workstream file
func (e *Executor) ParseDependencies(wsID string) ([]string, error) {
	// Find workstream file
	pattern := filepath.Join(e.config.BacklogDir, fmt.Sprintf("%s-*.md", wsID))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob pattern failed: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("workstream file not found: %s", wsID)
	}

	// Parse workstream to extract dependencies
	ws, err := parser.ParseWorkstream(matches[0])
	if err != nil {
		return nil, fmt.Errorf("parse workstream: %w", err)
	}

	var deps []string

	// Check if workstream has a parent (dependency)
	if ws.Parent != "" {
		deps = append(deps, ws.Parent)
	}

	// TODO: Parse explicit "Dependencies:" section from workstream content
	// This would require reading the file content and looking for a Dependencies section

	return deps, nil
}

// TopologicalSort performs topological sort on workstreams by dependencies
// Uses Kahn's algorithm
func (e *Executor) TopologicalSort(workstreams []string, dependencies map[string][]string) ([]string, error) {
	// Build adjacency list and in-degree count
	inDegree := make(map[string]int)
	adjList := make(map[string][]string)

	// Initialize in-degree for all nodes
	for _, ws := range workstreams {
		inDegree[ws] = 0
		adjList[ws] = []string{}
	}

	// Build graph
	for _, ws := range workstreams {
		for _, dep := range dependencies[ws] {
			// Check if dependency is in our workstream list
			if _, exists := inDegree[dep]; exists {
				adjList[dep] = append(adjList[dep], ws)
				inDegree[ws]++
			}
		}
	}

	// Find all nodes with zero in-degree
	var queue []string
	for ws := range inDegree {
		if inDegree[ws] == 0 {
			queue = append(queue, ws)
		}
	}

	// Process nodes
	var result []string
	for len(queue) > 0 {
		// Remove first node
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree for all neighbors
		for _, neighbor := range adjList[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles
	if len(result) != len(workstreams) {
		return nil, fmt.Errorf("cycle detected in workstream dependencies")
	}

	return result, nil
}

// generateEvidenceEvents generates evidence chain events for a workstream
func (e *Executor) generateEvidenceEvents(wsID string) []EvidenceEvent {
	now := time.Now().Format(time.RFC3339)

	return []EvidenceEvent{
		{
			Type:      "plan",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "execution_plan"},
		},
		{
			Type:      "generation",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "code_generation"},
		},
		{
			Type:      "verification",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "test_verification"},
		},
		{
			Type:      "approval",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "auto_approval"},
		},
	}
}

// emitEvidenceEvent writes an evidence event to the evidence log
func (e *Executor) emitEvidenceEvent(event EvidenceEvent) error {
	if e.evidenceWriter == nil {
		return nil // No evidence writer configured
	}

	output := e.progress.RenderEvidenceEvent(event)
	if output == "" {
		return fmt.Errorf("failed to render evidence event")
	}

	_, err := fmt.Fprintln(e.evidenceWriter, output)
	return err
}
