// Package nextstep provides deterministic next-step recommendations for SDP workflows.
package nextstep

// WorkstreamStatus represents the status of a single workstream.
type WorkstreamStatus struct {
	// ID is the workstream identifier (PP-FFF-SS format).
	ID string `json:"id"`

	// Status is the current state of the workstream.
	Status WorkstreamState `json:"status"`

	// Priority is the workstream priority (0 = highest).
	Priority int `json:"priority"`

	// BlockedBy lists workstream IDs that this workstream depends on.
	BlockedBy []string `json:"blocked_by,omitempty"`

	// Feature is the parent feature ID.
	Feature string `json:"feature,omitempty"`

	// Size indicates the workstream size (SMALL, MEDIUM, LARGE).
	Size string `json:"size,omitempty"`

	// LastError contains the last error if the workstream failed.
	LastError string `json:"last_error,omitempty"`
}

// WorkstreamState represents the state of a workstream.
type WorkstreamState string

const (
	// StatusBacklog indicates the workstream is in backlog.
	StatusBacklog WorkstreamState = "backlog"
	// StatusReady indicates the workstream is ready to start.
	StatusReady WorkstreamState = "ready"
	// StatusInProgress indicates the workstream is being worked on.
	StatusInProgress WorkstreamState = "in_progress"
	// StatusBlocked indicates the workstream is blocked.
	StatusBlocked WorkstreamState = "blocked"
	// StatusCompleted indicates the workstream is done.
	StatusCompleted WorkstreamState = "completed"
	// StatusFailed indicates the workstream failed.
	StatusFailed WorkstreamState = "failed"
)

// ProjectState represents the complete state of a project for recommendation.
type ProjectState struct {
	// Workstreams contains all known workstream statuses.
	Workstreams []WorkstreamStatus `json:"workstreams"`

	// LastCommand is the previously executed SDP command.
	LastCommand string `json:"last_command,omitempty"`

	// LastCommandError is the error from the last command, if any.
	LastCommandError string `json:"last_command_error,omitempty"`

	// Mode is the current execution mode.
	Mode ExecutionMode `json:"mode"`

	// GitStatus contains git repository information.
	GitStatus GitStatusInfo `json:"git_status"`

	// Config contains SDP configuration information.
	Config ConfigInfo `json:"config"`

	// ActiveWorkstream is the currently active workstream, if any.
	ActiveWorkstream string `json:"active_workstream,omitempty"`

	// SessionID is the current session identifier.
	SessionID string `json:"session_id,omitempty"`
}

// ExecutionMode represents the current execution mode.
type ExecutionMode string

const (
	// ModeDrive is the guided/driven mode with confirmations.
	ModeDrive ExecutionMode = "drive"
	// ModeShip is the autonomous mode without confirmations.
	ModeShip ExecutionMode = "ship"
	// ModeInteractive is the step-by-step interactive mode.
	ModeInteractive ExecutionMode = "interactive"
	// ModeManual is the manual mode without guidance.
	ModeManual ExecutionMode = "manual"
)

// GitStatusInfo contains git repository status information.
type GitStatusInfo struct {
	// Branch is the current git branch.
	Branch string `json:"branch,omitempty"`

	// Uncommitted indicates uncommitted changes exist.
	Uncommitted bool `json:"uncommitted"`

	// UpstreamDiverg indicates local and upstream have diverged.
	UpstreamDiverg bool `json:"upstream_diverged"`

	// IsRepo indicates if currently in a git repository.
	IsRepo bool `json:"is_repo"`

	// MainBranch is the name of the main/master branch.
	MainBranch string `json:"main_branch,omitempty"`
}

// ConfigInfo contains SDP configuration information.
type ConfigInfo struct {
	// HasSDPConfig indicates if .sdp/config.yml exists.
	HasSDPConfig bool `json:"has_sdp_config"`

	// Version is the SDP protocol version.
	Version string `json:"version,omitempty"`

	// ProjectRoot is the root directory of the project.
	ProjectRoot string `json:"project_root,omitempty"`

	// EvidenceEnabled indicates if evidence logging is enabled.
	EvidenceEnabled bool `json:"evidence_enabled"`
}

// ComparePriority compares two workstreams for priority ordering.
// Returns:
//
//	<0 if ws1 should be executed before ws2
//	0 if equal priority
//	>0 if ws2 should be executed before ws1
//
// Tie-break rules (in order):
//  1. Ready status over other states
//  2. Lower priority value (0 is highest)
//  3. Unblocked over blocked
//  4. Lower ID (lexicographic)
func ComparePriority(ws1, ws2 WorkstreamStatus) int {
	// Rule 1: Ready status takes precedence
	ws1Ready := ws1.Status == StatusReady
	ws2Ready := ws2.Status == StatusReady
	if ws1Ready != ws2Ready {
		if ws1Ready {
			return -1
		}
		return 1
	}

	// Rule 2: Unblocked over blocked
	ws1Blocked := len(ws1.BlockedBy) > 0
	ws2Blocked := len(ws2.BlockedBy) > 0
	if ws1Blocked != ws2Blocked {
		if ws2Blocked {
			return -1
		}
		return 1
	}

	// Rule 3: Lower priority value (0 is highest)
	if ws1.Priority != ws2.Priority {
		return ws1.Priority - ws2.Priority
	}

	// Rule 4: Lower ID (lexicographic)
	if ws1.ID < ws2.ID {
		return -1
	} else if ws1.ID > ws2.ID {
		return 1
	}

	return 0
}
