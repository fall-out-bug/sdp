package guard

import (
	"time"
)

// Exit codes for guard check command
const (
	ExitCodePass         = 0 // All checks passed
	ExitCodeViolation    = 1 // Policy violation (ERROR findings)
	ExitCodeRuntimeError = 2 // Runtime or configuration error
)

// Severity level for findings
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
)

// GuardResult represents the result of a guard check
type GuardResult struct {
	Allowed    bool     `json:"allowed"`
	WSID       string   `json:"ws_id,omitempty"`
	Reason     string   `json:"reason,omitempty"`
	ScopeFiles []string `json:"scope_files,omitempty"`
	Timestamp  string   `json:"timestamp"`
}

// CheckResult represents the result of a staged/CI guard check
type CheckResult struct {
	Success  bool         `json:"success"`
	ExitCode int          `json:"exit_code"`
	Findings []Finding    `json:"findings"`
	Summary  CheckSummary `json:"summary"`
}

// Finding represents a single policy finding
type Finding struct {
	Severity Severity `json:"severity"`
	Rule     string   `json:"rule"`
	File     string   `json:"file"`
	Line     int      `json:"line,omitempty"`
	Column   int      `json:"column,omitempty"`
	Message  string   `json:"message"`
}

// CheckSummary provides a summary of findings
type CheckSummary struct {
	Total    int `json:"total"`
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
}

// CheckOptions configures the check behavior
type CheckOptions struct {
	Staged bool   // Check only staged changes (--staged)
	JSON   bool   // Output JSON format (--json)
	Base   string // Base ref for diff (CI mode)
	Head   string // Head ref for diff (CI mode)
}

// GuardState represents the active workstream state
type GuardState struct {
	ActiveWS    string   `json:"active_ws"`
	ActivatedAt string   `json:"activated_at"`
	ScopeFiles  []string `json:"scope_files"`
	Timestamp   string   `json:"timestamp"`
}

// IsExpired checks if the guard state is expired (older than 24 hours)
func (gs *GuardState) IsExpired() bool {
	if gs.ActiveWS == "" {
		return true
	}

	activatedAt, err := time.Parse(time.RFC3339, gs.ActivatedAt)
	if err != nil {
		return true
	}

	return time.Since(activatedAt) > 24*time.Hour
}
