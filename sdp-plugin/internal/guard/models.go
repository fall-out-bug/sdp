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
	Total             int `json:"total"`
	Errors            int `json:"errors"`
	Warnings          int `json:"warnings"`
	AppliedExceptions int `json:"applied_exceptions,omitempty"`
}

// Exception represents a guard rule exception with TTL.
// AC1: Exception entries support rule_id, path glob, reason, owner, expires_at
type Exception struct {
	RuleID    string `json:"rule_id"`
	PathGlob  string `json:"path_glob"`
	Reason    string `json:"reason"`
	Owner     string `json:"owner"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at,omitempty"`
}

// IsExpired checks if the exception has expired.
// AC2: Expired exceptions are ignored automatically
func (e *Exception) IsExpired() bool {
	if e.ExpiresAt == "" {
		return false
	}

	expiresAt, err := time.Parse(time.RFC3339, e.ExpiresAt)
	if err != nil {
		return true
	}

	return time.Now().After(expiresAt)
}

// MatchesFile checks if the exception path glob matches a file
func (e *Exception) MatchesFile(filePath string) bool {
	return matchGlob(e.PathGlob, filePath)
}

// AppliedExceptionInfo tracks which exception was applied for auditability.
// AC6: JSON output includes exception metadata
type AppliedExceptionInfo struct {
	RuleID string `json:"rule_id"`
	File   string `json:"file"`
	Reason string `json:"reason"`
	Owner  string `json:"owner"`
}

// CheckResult represents the result of a staged/CI guard check
type CheckResult struct {
	Success           bool                   `json:"success"`
	ExitCode          int                    `json:"exit_code"`
	Findings          []Finding              `json:"findings"`
	Summary           CheckSummary           `json:"summary"`
	AppliedExceptions []AppliedExceptionInfo `json:"applied_exceptions,omitempty"`
}

// CheckOptions configures the check behavior
type CheckOptions struct {
	Staged bool   // Check only staged changes (--staged)
	JSON   bool   // Output JSON format (--json)
	Base   string // Base ref for diff (CI mode)
	Head   string // Head ref for diff (CI mode)
}

// ReviewFinding represents a finding from a review agent
type ReviewFinding struct {
	ID         string `json:"id"`
	FeatureID  string `json:"feature_id"`
	ReviewArea string `json:"review_area"`
	Title      string `json:"title"`
	Priority   int    `json:"priority"`
	BeadsID    string `json:"beads_id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	ResolvedAt string `json:"resolved_at"`
	ResolvedBy string `json:"resolved_by"`
}

// GuardState represents the active workstream state
type GuardState struct {
	ActiveWS       string          `json:"active_ws"`
	ActivatedAt    string          `json:"activated_at"`
	ScopeFiles     []string        `json:"scope_files"`
	Timestamp      string          `json:"timestamp"`
	ReviewFindings []ReviewFinding `json:"review_findings,omitempty"`
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

// AddFinding adds a review finding to the state
func (gs *GuardState) AddFinding(f ReviewFinding) {
	if gs.ReviewFindings == nil {
		gs.ReviewFindings = []ReviewFinding{}
	}
	gs.ReviewFindings = append(gs.ReviewFindings, f)
}

// ResolveFinding marks a finding as resolved
func (gs *GuardState) ResolveFinding(id string, resolvedBy string) bool {
	for i := range gs.ReviewFindings {
		if gs.ReviewFindings[i].ID == id {
			gs.ReviewFindings[i].Status = "resolved"
			gs.ReviewFindings[i].ResolvedAt = time.Now().Format(time.RFC3339)
			gs.ReviewFindings[i].ResolvedBy = resolvedBy
			return true
		}
	}
	return false
}

// GetOpenFindings returns all unresolved findings
func (gs *GuardState) GetOpenFindings() []ReviewFinding {
	var open []ReviewFinding
	for _, f := range gs.ReviewFindings {
		if f.Status != "resolved" {
			open = append(open, f)
		}
	}
	return open
}

// GetBlockingFindings returns P0/P1 findings that should block progress
func (gs *GuardState) GetBlockingFindings() []ReviewFinding {
	var blocking []ReviewFinding
	for _, f := range gs.ReviewFindings {
		if f.Status != "resolved" && f.Priority <= 1 {
			blocking = append(blocking, f)
		}
	}
	return blocking
}

// HasBlockingFindings returns true if there are unresolved P0/P1 findings
func (gs *GuardState) HasBlockingFindings() bool {
	return len(gs.GetBlockingFindings()) > 0
}

// FindingCount returns counts by status
func (gs *GuardState) FindingCount() (open int, resolved int, blocking int) {
	for _, f := range gs.ReviewFindings {
		if f.Status == "resolved" {
			resolved++
		} else {
			open++
			if f.Priority <= 1 {
				blocking++
			}
		}
	}
	return
}
