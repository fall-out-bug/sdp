package guard

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Skill implements the guard skill logic
type Skill struct {
	stateManager *StateManager
	activeWS     string
}

// NewSkill creates a new guard skill
func NewSkill(configDir string) *Skill {
	return &Skill{
		stateManager: NewStateManager(configDir),
		activeWS:     "",
	}
}

// Activate sets the active workstream
func (s *Skill) Activate(wsID string) error {
	s.activeWS = wsID

	// Load existing state to get scope files
	state, err := s.stateManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Create new state with current WS
	newState := GuardState{
		ActiveWS:    wsID,
		ScopeFiles:  state.ScopeFiles,
		ActivatedAt: time.Now().Format(time.RFC3339),
		Timestamp:   "",
	}

	// Save state
	if err := s.stateManager.Save(newState); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// CheckEdit checks if a file edit is allowed
func (s *Skill) CheckEdit(filePath string) (*GuardResult, error) {
	// Load current state
	state, err := s.stateManager.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	// Check if state is expired
	if state.IsExpired() {
		return &GuardResult{
			Allowed: false,
			Reason:  "No active WS (state expired). Run 'sdp guard activate <ws_id>' first.",
		}, nil
	}

	// No active WS
	if state.ActiveWS == "" {
		return &GuardResult{
			Allowed: false,
			Reason:  "No active WS. Run 'sdp guard activate <ws_id>' first.",
		}, nil
	}

	// No scope restrictions = all files allowed
	if len(state.ScopeFiles) == 0 {
		return &GuardResult{
			Allowed:    true,
			WSID:       state.ActiveWS,
			Reason:     "No scope restrictions",
			ScopeFiles: state.ScopeFiles,
		}, nil
	}

	// Check if file is in scope
	allowed := false
	for _, scopeFile := range state.ScopeFiles {
		if filePath == scopeFile {
			allowed = true
			break
		}
	}

	if !allowed {
		return &GuardResult{
			Allowed:    false,
			WSID:       state.ActiveWS,
			Reason:     fmt.Sprintf("File %s not in WS scope", filePath),
			ScopeFiles: state.ScopeFiles,
		}, nil
	}

	return &GuardResult{
		Allowed:    true,
		WSID:       state.ActiveWS,
		Reason:     "File in scope",
		ScopeFiles: state.ScopeFiles,
	}, nil
}

// GetActiveWS returns the currently active workstream ID
func (s *Skill) GetActiveWS() string {
	state, _ := s.stateManager.Load()
	return state.ActiveWS
}

// Deactivate deactivates the current workstream
func (s *Skill) Deactivate() error {
	if err := s.stateManager.Clear(); err != nil {
		return fmt.Errorf("failed to clear state: %w", err)
	}
	s.activeWS = ""
	return nil
}

// ResolvePath resolves a relative file path to absolute path
func ResolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	return filepath.Join(wd, path), nil
}
