package executor

import (
	"io"
	"time"
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
	TotalWorkstreams int               `json:"total_workstreams"`
	Executed         int               `json:"executed"`
	Succeeded        int               `json:"succeeded"`
	Failed           int               `json:"failed"`
	Skipped          int               `json:"skipped"`
	Retries          int               `json:"retries"`
	Duration         time.Duration     `json:"duration"`
	EvidenceEvents   []EvidenceEvent   `json:"evidence_events"`
}

// ExecutionSummary is a simplified summary for output
type ExecutionSummary struct {
	TotalWorkstreams int    `json:"total_workstreams"`
	Executed         int    `json:"executed"`
	Succeeded        int    `json:"succeeded"`
	Failed           int    `json:"failed"`
	Skipped          int    `json:"skipped"`
	Retries          int    `json:"retries"`
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
