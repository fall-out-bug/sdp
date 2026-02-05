package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Status represents the current status of a feature execution
type Status string

const (
	StatusPending     Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted, StatusFailed:
		return true
	default:
		return false
	}
}

// Checkpoint represents a saved state of feature execution
type Checkpoint struct {
	ID                  string                 `json:"id"`
	FeatureID           string                 `json:"feature_id"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	Status              Status                 `json:"status"`
	CompletedWorkstreams []string              `json:"completed_workstreams"`
	CurrentWorkstream   string                 `json:"current_workstream,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// Manager handles checkpoint operations
type Manager struct {
	dir string
}

// NewManager creates a new checkpoint manager
func NewManager(dir string) *Manager {
	return &Manager{
		dir: dir,
	}
}

// Save saves a checkpoint to disk
func (m *Manager) Save(cp Checkpoint) error {
	// Ensure directory exists
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return fmt.Errorf("failed to create checkpoint directory: %w", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}

	// Write to file
	filename := filepath.Join(m.dir, fmt.Sprintf("%s.json", cp.ID))
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint file: %w", err)
	}

	return nil
}

// Load loads a checkpoint from disk
func (m *Manager) Load(id string) (Checkpoint, error) {
	filename := filepath.Join(m.dir, fmt.Sprintf("%s.json", id))

	data, err := os.ReadFile(filename)
	if err != nil {
		return Checkpoint{}, err
	}

	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return Checkpoint{}, fmt.Errorf("failed to unmarshal checkpoint: %w", err)
	}

	return cp, nil
}

// Resume loads a checkpoint for resuming execution
func (m *Manager) Resume(id string) (Checkpoint, error) {
	cp, err := m.Load(id)
	if err != nil {
		return Checkpoint{}, fmt.Errorf("failed to resume checkpoint: %w", err)
	}

	// Update updated_at timestamp
	cp.UpdatedAt = time.Now()

	return cp, nil
}

// List returns all checkpoints
func (m *Manager) List() ([]Checkpoint, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []Checkpoint{}, nil
		}
		return nil, fmt.Errorf("failed to read checkpoint directory: %w", err)
	}

	var checkpoints []Checkpoint
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Skip non-JSON files
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		// Extract ID from filename
		id := entry.Name()[:len(entry.Name())-5] // Remove .json extension

		cp, err := m.Load(id)
		if err != nil {
			// Skip invalid checkpoints
			continue
		}

		checkpoints = append(checkpoints, cp)
	}

	// Sort by created_at (newest first)
	sort.Slice(checkpoints, func(i, j int) bool {
		return checkpoints[i].CreatedAt.After(checkpoints[j].CreatedAt)
	})

	return checkpoints, nil
}

// Clean removes checkpoints older than the specified age
func (m *Manager) Clean(age time.Duration) (int, error) {
	checkpoints, err := m.List()
	if err != nil {
		return 0, fmt.Errorf("failed to list checkpoints: %w", err)
	}

	cutoff := time.Now().Add(-age)
	deleted := 0

	for _, cp := range checkpoints {
		// Only clean completed checkpoints that are old enough
		if cp.Status != StatusCompleted {
			continue
		}

		if cp.UpdatedAt.Before(cutoff) {
			filename := filepath.Join(m.dir, fmt.Sprintf("%s.json", cp.ID))
			if err := os.Remove(filename); err != nil {
				return deleted, fmt.Errorf("failed to remove checkpoint %s: %w", cp.ID, err)
			}
			deleted++
		}
	}

	return deleted, nil
}

// GetDefaultDir returns the default checkpoint directory (.sdp/checkpoints)
func GetDefaultDir() (string, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create .sdp/checkpoints path
	checkpointDir := filepath.Join(cwd, ".sdp", "checkpoints")

	return checkpointDir, nil
}
