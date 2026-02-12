package guard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/guard"
)

const (
	// StateFile is the guard state filename
	StateFile = "state.json"
)

// Skill implements the guard skill logic with exceptions support
type Skill struct {
	stateManager *StateManager
	exceptions    *Manager
	activeWS      string
}

// NewSkill creates a new guard skill
func NewSkill(configDir string) *Skill {
	return &Skill{
		stateManager: NewStateManager(configDir),
		exceptions:    NewManager(configDir),
		activeWS:      "",
	}
}

// Activate activates a workstream and records activation timestamp
func (s *Skill) Activate(wsID string) error {
	// Store active workstream
	s.activeWS = wsID

	// Save state with new active WS
	state := GuardState{
		ActiveWS:    wsID,
		ActivatedAt: time.Now().Format(time.RFC3339),
		ScopeFiles:  scopeFilesForWS(wsID),
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	if err := s.stateManager.Save(state); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// GetActiveWS returns the currently active workstream ID
func (s *Skill) GetActiveWS() string {
	return s.activeWS
}

// Deactivate deactivates the current workstream
func (s *Skill) Deactivate() error {
	s.activeWS = ""
	return s.stateManager.Clear()
}

// CheckEdit checks if editing a file is allowed
func (s *Skill) CheckEdit(filePath string) (*GuardResult, error) {
	// Get active state
	state, _ := s.stateManager.Load()

	// No active WS or expired = blocked
	if state.ActiveWS == "" || state.IsExpired() {
		return &GuardResult{
			Allowed: false,
			Reason:  "No active WS (state expired)",
			WSID:   state.ActiveWS,
		}, nil
	}

	// Check if file is in scope
	inScope := false
	for _, scopeFile := range state.ScopeFiles {
		if filePath == scopeFile {
			inScope = true
			break
		}
	}

	// Not in scope
	if !inScope {
		return &GuardResult{
			Allowed: false,
			Reason:  "File not in workstream scope",
			WSID:   state.ActiveWS,
			ScopeFiles: state.ScopeFiles,
		}, nil
	}

	// Check for exceptions (AC3)
	activeExceptions, err := s.exceptions.GetActive()
	if err == nil && len(activeExceptions) > 0 {
		// Apply exceptions
		for _, exc := range activeExceptions {
			if exc.MatchesPath(filePath) {
				// Exception applies, allow the edit
				return &GuardResult{
					Allowed:    true,
					Reason:     fmt.Sprintf("Exception %%s: %%s", exc.ID, exc.Reason),
					WSID:       state.ActiveWS,
				}, nil
			}
		}
	}

	// No exception matched, not in scope
	return &GuardResult{
		Allowed: false,
		Reason:  "File edit not allowed (no active exception for this path)",
		WSID:   state.ActiveWS,
	}, nil
}
