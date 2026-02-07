package graph

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint represents a serialized state of dispatcher execution
type Checkpoint struct {
	Version        string                   `json:"version"`
	FeatureID      string                   `json:"feature_id"`
	Timestamp      time.Time                `json:"timestamp"`
	Completed      []string                 `json:"completed"`
	Failed         []string                 `json:"failed"`
	Graph          *GraphSnapshot           `json:"graph"`
	CircuitBreaker *CircuitBreakerSnapshot `json:"circuit_breaker"`
}

// GraphSnapshot represents the state of the dependency graph
type GraphSnapshot struct {
	Nodes []NodeSnapshot      `json:"nodes"`
	Edges map[string][]string `json:"edges"`
}

// NodeSnapshot represents the state of a single workstream node
type NodeSnapshot struct {
	ID        string   `json:"id"`
	DependsOn []string `json:"depends_on"`
	Indegree  int      `json:"indegree"`
	Completed bool     `json:"completed"`
}

// CircuitBreakerSnapshot represents the state of the circuit breaker
type CircuitBreakerSnapshot struct {
	State            int       `json:"state"`
	FailureCount     int       `json:"failure_count"`
	SuccessCount     int       `json:"success_count"`
	ConsecutiveOpens int       `json:"consecutive_opens"`
	LastFailureTime  time.Time `json:"last_failure_time"`
}

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

// copyStringSlice creates a deep copy of a string slice
func copyStringSlice(slice []string) []string {
	if slice == nil {
		return nil
	}
	copied := make([]string, len(slice))
	copy(copied, slice)
	return copied
}

// CreateCheckpoint creates a checkpoint from the current dispatcher state
func (cm *CheckpointManager) CreateCheckpoint(graph *DependencyGraph, featureID string, completed []string, failed []string) *Checkpoint {
	// Snapshot graph nodes
	nodes := make([]NodeSnapshot, 0, len(graph.nodes))
	for _, node := range graph.nodes {
		nodes = append(nodes, NodeSnapshot{
			ID:        node.ID,
			DependsOn: copyStringSlice(node.DependsOn),
			Indegree:  node.Indegree,
			Completed: node.Completed,
		})
	}

	// Snapshot graph edges (copy to avoid mutation)
	edges := make(map[string][]string)
	for from, toList := range graph.edges {
		edges[from] = copyStringSlice(toList)
	}

	return &Checkpoint{
		Version:   "1.0",
		FeatureID: featureID,
		Timestamp: time.Now().UTC(),
		Completed: copyStringSlice(completed),
		Failed:    copyStringSlice(failed),
		Graph: &GraphSnapshot{
			Nodes: nodes,
			Edges: edges,
		},
		CircuitBreaker: nil, // Will be populated by dispatcher
	}
}

// RestoreGraph recreates a dependency graph from a checkpoint
func (cm *CheckpointManager) RestoreGraph(checkpoint *Checkpoint) *DependencyGraph {
	graph := NewDependencyGraph()

	// Restore nodes
	for _, nodeSnapshot := range checkpoint.Graph.Nodes {
		node := &WorkstreamNode{
			ID:        nodeSnapshot.ID,
			DependsOn: copyStringSlice(nodeSnapshot.DependsOn),
			Indegree:  nodeSnapshot.Indegree,
			Completed: nodeSnapshot.Completed,
		}
		graph.nodes[nodeSnapshot.ID] = node
	}

	// Restore edges
	for from, toList := range checkpoint.Graph.Edges {
		graph.edges[from] = copyStringSlice(toList)
	}

	return graph
}
