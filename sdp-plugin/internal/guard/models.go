package guard

import (
	"time"
)

// GuardResult represents the result of a guard check
type GuardResult struct {
	Allowed    bool     `json:"allowed"`
	WSID       string   `json:"ws_id,omitempty"`
	Reason     string   `json:"reason,omitempty"`
	ScopeFiles []string `json:"scope_files,omitempty"`
	Timestamp  string   `json:"timestamp"`
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
