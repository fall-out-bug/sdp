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

// ReviewFinding represents a finding from a review agent
type ReviewFinding struct {
	ID          string `json:"id"`           // Beads issue ID (e.g., sdp-abc123)
	FeatureID   string `json:"feature_id"`   // Feature ID (e.g., F051)
	ReviewArea  string `json:"review_area"`  // QA, Security, DevOps, SRE, TechLead, Documentation
	Title       string `json:"title"`        // Issue title
	Priority    int    `json:"priority"`     // 0=P0, 1=P1, 2=P2, 3=P3
	BeadsID     string `json:"beads_id"`     // Beads issue ID if created
	Status      string `json:"status"`       // open, in_progress, resolved
	CreatedAt   string `json:"created_at"`   // ISO timestamp
	ResolvedAt  string `json:"resolved_at"`  // ISO timestamp when resolved
	ResolvedBy  string `json:"resolved_by"`  // How it was resolved
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
