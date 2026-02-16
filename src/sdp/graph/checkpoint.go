package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CheckpointManager manages atomic checkpoint persistence
type CheckpointManager struct {
	checkpointDir string
	featureID     string
}

// NewCheckpointManager creates a new checkpoint manager for the given feature
func NewCheckpointManager(featureID string) *CheckpointManager {
	return &CheckpointManager{
		checkpointDir: filepath.Join(".sdp", "checkpoints"),
		featureID:     featureID,
	}
}

// SetCheckpointDir sets the checkpoint directory (for testing)
func (cm *CheckpointManager) SetCheckpointDir(dir string) {
	cm.checkpointDir = dir
}

// GetFeatureID returns the feature ID
func (cm *CheckpointManager) GetFeatureID() string {
	return cm.featureID
}

// GetCheckpointPath returns the path to the checkpoint file
func (cm *CheckpointManager) GetCheckpointPath() string {
	return filepath.Join(cm.checkpointDir, fmt.Sprintf("%s-checkpoint.json", cm.featureID))
}

// GetTempPath returns the path to the temporary checkpoint file
func (cm *CheckpointManager) GetTempPath() string {
	return cm.GetCheckpointPath() + ".tmp"
}

// Save writes the checkpoint to disk atomically
// Algorithm: write to temp file -> fsync -> atomic rename
func (cm *CheckpointManager) Save(checkpoint *Checkpoint) error {
	// Ensure checkpoint directory exists
	if err := os.MkdirAll(cm.checkpointDir, 0755); err != nil {
		return fmt.Errorf("failed to create checkpoint directory: %w", err)
	}

	// Step 1: Write to temporary file
	tmpPath := cm.GetTempPath()
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}

	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Step 2: Fsync to disk (ensure data persistence)
	f, err := os.Open(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to open temp file for fsync: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return fmt.Errorf("failed to fsync temp file: %w", err)
	}
	f.Close()

	// Step 3: Atomic rename
	finalPath := cm.GetCheckpointPath()
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// Load reads the checkpoint from disk
// Returns nil if checkpoint doesn't exist
// Returns error if checkpoint is corrupt
func (cm *CheckpointManager) Load() (*Checkpoint, error) {
	finalPath := cm.GetCheckpointPath()

	// Check if file exists
	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		// No checkpoint exists, return nil (not an error)
		return nil, nil
	}

	// Read file
	data, err := os.ReadFile(finalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint: %w", err)
	}

	// Unmarshal JSON
	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		// Corrupt checkpoint - move to .corrupt suffix
		corruptPath := finalPath + ".corrupt"
		os.Rename(finalPath, corruptPath)
		return nil, fmt.Errorf("corrupt checkpoint (moved to %s): %w", corruptPath, err)
	}

	return &checkpoint, nil
}

// Delete removes the checkpoint file
func (cm *CheckpointManager) Delete() error {
	finalPath := cm.GetCheckpointPath()
	tmpPath := cm.GetTempPath()

	// Remove final checkpoint if exists
	if _, err := os.Stat(finalPath); err == nil {
		if err := os.Remove(finalPath); err != nil {
			return fmt.Errorf("failed to delete checkpoint: %w", err)
		}
	}

	// Remove temp file if exists
	if _, err := os.Stat(tmpPath); err == nil {
		os.Remove(tmpPath) // Ignore error for temp file
	}

	return nil
}
