package guard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// GuardStateFile is the filename for storing guard state
	GuardStateFile = ".guard_state.json"
)

// StateManager manages guard state persistence
type StateManager struct {
	stateFile string
	configDir string
}

// NewStateManager creates a new StateManager
func NewStateManager(configDir string) *StateManager {
	return &StateManager{
		stateFile: GuardStateFile,
		configDir: configDir,
	}
}

// Save saves the guard state to file
func (sm *StateManager) Save(state GuardState) error {
	// Create config directory if needed
	if sm.configDir != "" {
		if err := os.MkdirAll(sm.configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	statePath := filepath.Join(sm.configDir, sm.stateFile)

	// Set timestamp
	state.Timestamp = time.Now().Format(time.RFC3339)

	// Marshal to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write with secure permissions
	if err := os.WriteFile(statePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// Load loads the guard state from file
func (sm *StateManager) Load() (*GuardState, error) {
	statePath := filepath.Join(sm.configDir, sm.stateFile)

	// Check if file exists
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// No active state
		return &GuardState{}, nil
	}

	// Read file
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	// Unmarshal JSON
	var state GuardState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return &state, nil
}

// Clear removes the guard state file
func (sm *StateManager) Clear() error {
	statePath := filepath.Join(sm.configDir, sm.stateFile)

	// Check if file exists
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return nil // Already cleared
	}

	// Remove file
	if err := os.Remove(statePath); err != nil {
		return fmt.Errorf("failed to remove state file: %w", err)
	}

	return nil
}
