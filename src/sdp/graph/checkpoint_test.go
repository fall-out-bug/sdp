package graph

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewCheckpointManager verifies checkpoint manager creation
func TestNewCheckpointManager(t *testing.T) {
	cm := NewCheckpointManager("F001")
	if cm == nil {
		t.Fatal("NewCheckpointManager returned nil")
	}
	if cm.featureID != "F001" {
		t.Errorf("expected featureID F001, got %s", cm.featureID)
	}
}

// TestCheckpointManager_GetCheckpointPath verifies path generation
func TestCheckpointManager_GetCheckpointPath(t *testing.T) {
	cm := NewCheckpointManager("F001")
	expected := filepath.Join(".sdp", "checkpoints", "F001-checkpoint.json")
	if cm.GetCheckpointPath() != expected {
		t.Errorf("expected %s, got %s", expected, cm.GetCheckpointPath())
	}
}

// TestCheckpointManager_GetTempPath verifies temp path generation
func TestCheckpointManager_GetTempPath(t *testing.T) {
	cm := NewCheckpointManager("F001")
	expected := filepath.Join(".sdp", "checkpoints", "F001-checkpoint.json.tmp")
	if cm.GetTempPath() != expected {
		t.Errorf("expected %s, got %s", expected, cm.GetTempPath())
	}
}

// TestCheckpointManager_SaveAndLoad verifies save and load cycle
func TestCheckpointManager_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cm := NewCheckpointManager("F001")
	cm.SetCheckpointDir(tmpDir)

	checkpoint := &Checkpoint{
		Version:   "1.0",
		FeatureID: "F001",
		Timestamp: time.Now().UTC(),
		Completed: []string{"ws-001", "ws-002"},
		Failed:    []string{"ws-003"},
		Graph: &GraphSnapshot{
			Nodes: []NodeSnapshot{
				{ID: "ws-001", Indegree: 0, Completed: true},
				{ID: "ws-002", Indegree: 1, Completed: true},
				{ID: "ws-003", Indegree: 1, Completed: false},
			},
			Edges: map[string][]string{
				"ws-001": {"ws-002"},
			},
		},
		CircuitBreaker: &CircuitBreakerSnapshot{
			State:        int(StateClosed),
			FailureCount: 0,
			SuccessCount: 2,
		},
	}

	// Save
	err := cm.Save(checkpoint)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(cm.GetCheckpointPath()); os.IsNotExist(err) {
		t.Error("checkpoint file not created")
	}

	// Load
	loaded, err := cm.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.FeatureID != checkpoint.FeatureID {
		t.Errorf("expected FeatureID %s, got %s", checkpoint.FeatureID, loaded.FeatureID)
	}
	if len(loaded.Completed) != 2 {
		t.Errorf("expected 2 completed, got %d", len(loaded.Completed))
	}
	if len(loaded.Failed) != 1 {
		t.Errorf("expected 1 failed, got %d", len(loaded.Failed))
	}
}

// TestCheckpointManager_Load_NonExistent verifies load of non-existent file
func TestCheckpointManager_Load_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	cm := NewCheckpointManager("F001")
	cm.SetCheckpointDir(tmpDir)

	loaded, err := cm.Load()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil for non-existent checkpoint")
	}
}

// TestCheckpointManager_Delete verifies checkpoint deletion
func TestCheckpointManager_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	cm := NewCheckpointManager("F001")
	cm.SetCheckpointDir(tmpDir)

	// Create checkpoint
	checkpoint := &Checkpoint{
		Version:   "1.0",
		FeatureID: "F001",
		Timestamp: time.Now(),
	}
	cm.Save(checkpoint)

	// Delete
	err := cm.Delete()
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify file is gone
	if _, err := os.Stat(cm.GetCheckpointPath()); !os.IsNotExist(err) {
		t.Error("checkpoint file still exists after delete")
	}
}

// TestCheckpointManager_Delete_NonExistent verifies delete of non-existent
func TestCheckpointManager_Delete_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	cm := NewCheckpointManager("F001")
	cm.SetCheckpointDir(tmpDir)

	err := cm.Delete()
	if err != nil {
		t.Errorf("delete of non-existent should not error: %v", err)
	}
}

// TestCheckpointManager_CreateCheckpoint verifies checkpoint creation
func TestCheckpointManager_CreateCheckpoint(t *testing.T) {
	cm := NewCheckpointManager("F001")

	graph := NewDependencyGraph()
	graph.AddNode("ws-001", nil)
	graph.AddNode("ws-002", []string{"ws-001"})

	completed := []string{"ws-001"}
	failed := []string{}

	checkpoint := cm.CreateCheckpoint(graph, "F001", completed, failed)

	if checkpoint.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", checkpoint.Version)
	}
	if checkpoint.FeatureID != "F001" {
		t.Errorf("expected featureID F001, got %s", checkpoint.FeatureID)
	}
	if len(checkpoint.Completed) != 1 {
		t.Errorf("expected 1 completed, got %d", len(checkpoint.Completed))
	}
	if len(checkpoint.Graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(checkpoint.Graph.Nodes))
	}
}

// TestCheckpointManager_RestoreGraph verifies graph restoration
func TestCheckpointManager_RestoreGraph(t *testing.T) {
	cm := NewCheckpointManager("F001")

	checkpoint := &Checkpoint{
		Version:   "1.0",
		FeatureID: "F001",
		Timestamp: time.Now(),
		Graph: &GraphSnapshot{
			Nodes: []NodeSnapshot{
				{ID: "ws-001", DependsOn: []string{}, Indegree: 0, Completed: true},
				{ID: "ws-002", DependsOn: []string{"ws-001"}, Indegree: 1, Completed: false},
			},
			Edges: map[string][]string{
				"ws-001": {"ws-002"},
			},
		},
	}

	graph := cm.RestoreGraph(checkpoint)

	if len(graph.nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(graph.nodes))
	}
	if graph.nodes["ws-001"].Completed != true {
		t.Error("ws-001 should be completed")
	}
	if graph.nodes["ws-002"].Indegree != 1 {
		t.Errorf("expected indegree 1, got %d", graph.nodes["ws-002"].Indegree)
	}
}

// TestCheckpointManager_GetFeatureID verifies GetFeatureID
func TestCheckpointManager_GetFeatureID(t *testing.T) {
	cm := NewCheckpointManager("F999")
	if cm.GetFeatureID() != "F999" {
		t.Errorf("expected F999, got %s", cm.GetFeatureID())
	}
}

// TestCopyStringSlice verifies slice copying
func TestCopyStringSlice(t *testing.T) {
	original := []string{"a", "b", "c"}
	copied := copyStringSlice(original)

	if len(copied) != len(original) {
		t.Error("length mismatch")
	}

	// Modify original - copy should not change
	original[0] = "x"
	if copied[0] == "x" {
		t.Error("copy was not independent")
	}
}

// TestCopyStringSlice_Nil verifies nil slice handling
func TestCopyStringSlice_Nil(t *testing.T) {
	result := copyStringSlice(nil)
	if result != nil {
		t.Error("expected nil for nil input")
	}
}

// TestCheckpoint_Types verifies checkpoint types
func TestCheckpoint_Types(t *testing.T) {
	checkpoint := Checkpoint{
		Version:   "1.0",
		FeatureID: "F001",
		Timestamp: time.Now(),
		Completed: []string{"ws-001"},
		Failed:    []string{"ws-002"},
		Graph: &GraphSnapshot{
			Nodes: []NodeSnapshot{},
			Edges: map[string][]string{},
		},
		CircuitBreaker: &CircuitBreakerSnapshot{
			State:        int(StateClosed),
			FailureCount: 0,
		},
	}

	if checkpoint.Version != "1.0" {
		t.Error("Version not set")
	}
	if len(checkpoint.Completed) != 1 {
		t.Error("Completed not set")
	}
}

// TestNodeSnapshot verifies node snapshot structure
func TestNodeSnapshot(t *testing.T) {
	snapshot := NodeSnapshot{
		ID:        "ws-001",
		DependsOn: []string{"ws-000"},
		Indegree:  1,
		Completed: true,
	}

	if snapshot.ID != "ws-001" {
		t.Error("ID not set")
	}
	if len(snapshot.DependsOn) != 1 {
		t.Error("DependsOn not set")
	}
}

// TestGraphSnapshot verifies graph snapshot structure
func TestGraphSnapshot(t *testing.T) {
	snapshot := GraphSnapshot{
		Nodes: []NodeSnapshot{
			{ID: "a"},
			{ID: "b"},
		},
		Edges: map[string][]string{
			"a": {"b"},
		},
	}

	if len(snapshot.Nodes) != 2 {
		t.Error("Nodes not set")
	}
	if len(snapshot.Edges) != 1 {
		t.Error("Edges not set")
	}
}

// TestCircuitBreakerSnapshot verifies circuit breaker snapshot structure
func TestCircuitBreakerSnapshot(t *testing.T) {
	now := time.Now()
	snapshot := CircuitBreakerSnapshot{
		State:            int(StateOpen),
		FailureCount:     3,
		SuccessCount:     5,
		ConsecutiveOpens: 2,
		LastFailureTime:  now,
		LastStateChange:  now,
	}

	if snapshot.State != int(StateOpen) {
		t.Error("State not set")
	}
	if snapshot.FailureCount != 3 {
		t.Error("FailureCount not set")
	}
}
