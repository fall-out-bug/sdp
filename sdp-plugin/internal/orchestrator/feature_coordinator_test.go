package orchestrator

import (
	"strings"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
)

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// MockWorkstreamLoader mocks the WorkstreamLoader interface
type MockWorkstreamLoader struct {
	workstreams map[string][]WorkstreamNode
	loadErr     map[string]error
}

func (m *MockWorkstreamLoader) LoadWorkstreams(featureID string) ([]WorkstreamNode, error) {
	if err, ok := m.loadErr[featureID]; ok {
		return nil, err
	}
	return m.workstreams[featureID], nil
}

// MockWorkstreamExecutor mocks the WorkstreamExecutor interface
type MockWorkstreamExecutor struct {
	executed []string
	execErr  map[string]error
}

func (m *MockWorkstreamExecutor) Execute(wsID string) error {
	m.executed = append(m.executed, wsID)
	if err, ok := m.execErr[wsID]; ok {
		return err
	}
	return nil
}

func (m *MockWorkstreamExecutor) GetExecuted() []string {
	return m.executed
}

// MockCheckpointSaver mocks the CheckpointSaver interface
type MockCheckpointSaver struct {
	checkpoints map[string]checkpoint.Checkpoint
	saveErr     error
	loadErr     error
}

func (m *MockCheckpointSaver) Save(cp checkpoint.Checkpoint) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.checkpoints == nil {
		m.checkpoints = make(map[string]checkpoint.Checkpoint)
	}
	m.checkpoints[cp.ID] = cp
	return nil
}

func (m *MockCheckpointSaver) Load(id string) (checkpoint.Checkpoint, error) {
	if m.loadErr != nil {
		return checkpoint.Checkpoint{}, m.loadErr
	}
	if cp, ok := m.checkpoints[id]; ok {
		return cp, nil
	}
	return checkpoint.Checkpoint{}, ErrCheckpointNotFound
}

func (m *MockCheckpointSaver) Resume(id string) (checkpoint.Checkpoint, error) {
	cp, err := m.Load(id)
	if err != nil {
		return checkpoint.Checkpoint{}, err
	}
	cp.UpdatedAt = time.Now()
	return cp, nil
}

func TestFeatureCoordinator_NewFeatureCoordinator(t *testing.T) {
	loader := &MockWorkstreamLoader{
		workstreams: make(map[string][]WorkstreamNode),
	}
	executor := &MockWorkstreamExecutor{}
	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	if coordinator == nil {
		t.Fatal("Expected non-nil coordinator")
	}

	// Verify coordinator was created successfully
	if coordinator.orchestrator == nil {
		t.Error("Expected orchestrator to be initialized")
	}
}

func TestFeatureCoordinator_ExecuteFeature_Success(t *testing.T) {
	// AC1: Orchestrator executes workstreams in dependency order
	workstreams := []WorkstreamNode{
		{
			ID:           "WS-001",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{},
		},
		{
			ID:           "WS-002",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{"WS-001"},
		},
		{
			ID:           "WS-003",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{"WS-001"},
		},
	}

	loader := &MockWorkstreamLoader{
		workstreams: map[string][]WorkstreamNode{
			"F001": workstreams,
		},
	}

	executor := &MockWorkstreamExecutor{
		execErr: make(map[string]error),
	}

	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	// Execute feature
	err := coordinator.ExecuteFeature("F001")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify execution order (topological sort)
	executed := executor.GetExecuted()
	if len(executed) != 3 {
		t.Fatalf("Expected 3 executed workstreams, got %d", len(executed))
	}

	// WS-001 must execute first (no dependencies)
	if executed[0] != "WS-001" {
		t.Errorf("Expected first workstream to be WS-001, got %s", executed[0])
	}

	// Verify checkpoint was saved
	cp, err := saver.Load("F001")
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	// AC2: Checkpoint includes status, completed_ws, execution_order
	if cp.Status != checkpoint.StatusCompleted {
		t.Errorf("Expected status completed, got %s", cp.Status)
	}

	if len(cp.CompletedWorkstreams) != 3 {
		t.Errorf("Expected 3 completed workstreams, got %d", len(cp.CompletedWorkstreams))
	}
}

func TestFeatureCoordinator_ExecuteFeature_EmptyWorkstreams(t *testing.T) {
	loader := &MockWorkstreamLoader{
		workstreams: map[string][]WorkstreamNode{
			"F001": {},
		},
	}

	executor := &MockWorkstreamExecutor{}
	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	err := coordinator.ExecuteFeature("F001")

	if err != ErrFeatureNotFound {
		t.Errorf("Expected ErrFeatureNotFound, got %v", err)
	}
}

func TestFeatureCoordinator_ExecuteFeature_WorkstreamRetry(t *testing.T) {
	// AC5: Auto-retry for HIGH/MEDIUM issues (max 2 retries)
	workstreams := []WorkstreamNode{
		{
			ID:           "WS-001",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{},
		},
	}

	loader := &MockWorkstreamLoader{
		workstreams: map[string][]WorkstreamNode{
			"F001": workstreams,
		},
	}

	executor := &MockWorkstreamExecutor{
		execErr: map[string]error{
			"WS-001": ErrExecutionFailed, // Fail first time
		},
	}

	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	err := coordinator.ExecuteFeature("F001")

	// Should fail after max retries
	if err == nil {
		t.Fatal("Expected error after retries, got nil")
	}

	// Verify retry attempts (should be maxRetries + 1 initial attempts)
	executed := executor.GetExecuted()
	expectedAttempts := 3 // 1 initial + 2 retries
	if len(executed) != expectedAttempts {
		t.Errorf("Expected %d retry attempts, got %d", expectedAttempts, len(executed))
	}
}

func TestFeatureCoordinator_ResumeFeature(t *testing.T) {
	// AC2: On resume, skip completed workstreams
	workstreams := []WorkstreamNode{
		{
			ID:           "WS-001",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{},
		},
		{
			ID:           "WS-002",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{"WS-001"},
		},
	}

	loader := &MockWorkstreamLoader{
		workstreams: map[string][]WorkstreamNode{
			"F001": workstreams,
		},
	}

	executor := &MockWorkstreamExecutor{}
	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	// Create a checkpoint showing WS-001 is complete
	existingCheckpoint := checkpoint.Checkpoint{
		ID:                   "F001",
		FeatureID:            "F001",
		Status:               checkpoint.StatusInProgress,
		CompletedWorkstreams: []string{"WS-001"},
		CurrentWorkstream:    "WS-002",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	saver.Save(existingCheckpoint)

	// Resume from checkpoint
	err := coordinator.ResumeFeature("F001")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify only remaining workstream was executed
	executed := executor.GetExecuted()
	if len(executed) != 1 {
		t.Fatalf("Expected 1 executed workstream (WS-002 only), got %d", len(executed))
	}

	if executed[0] != "WS-002" {
		t.Errorf("Expected WS-002 to be executed, got %s", executed[0])
	}

	// Verify checkpoint updated
	cp, _ := saver.Load("F001")
	if cp.Status != checkpoint.StatusCompleted {
		t.Errorf("Expected status completed after resume, got %s", cp.Status)
	}

	if len(cp.CompletedWorkstreams) != 2 {
		t.Errorf("Expected 2 completed workstreams, got %d", len(cp.CompletedWorkstreams))
	}
}

func TestFeatureCoordinator_ExecuteFeature_ProgressTracking(t *testing.T) {
	// AC4: Progress tracking with timestamps
	workstreams := []WorkstreamNode{
		{
			ID:           "WS-001",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{},
		},
	}

	loader := &MockWorkstreamLoader{
		workstreams: map[string][]WorkstreamNode{
			"F001": workstreams,
		},
	}

	executor := &MockWorkstreamExecutor{}
	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	// Execute with progress tracking
	var progressUpdates []ProgressUpdate
	coordinator.progressCallback = func(update ProgressUpdate) {
		progressUpdates = append(progressUpdates, update)
	}

	err := coordinator.ExecuteFeature("F001")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify progress updates were captured
	if len(progressUpdates) == 0 {
		t.Fatal("Expected progress updates, got none")
	}

	// Check for start timestamp
	if progressUpdates[0].Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp in progress update")
	}

	// Check for completion summary
	lastUpdate := progressUpdates[len(progressUpdates)-1]
	if lastUpdate.Message == "" {
		t.Error("Expected message in progress update")
	}
}

func TestFeatureCoordinator_ExecuteFeature_CircularDependency(t *testing.T) {
	// Test circular dependency detection
	workstreams := []WorkstreamNode{
		{
			ID:           "WS-001",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{"WS-002"},
		},
		{
			ID:           "WS-002",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{"WS-001"},
		},
	}

	loader := &MockWorkstreamLoader{
		workstreams: map[string][]WorkstreamNode{
			"F001": workstreams,
		},
	}

	executor := &MockWorkstreamExecutor{}
	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	err := coordinator.ExecuteFeature("F001")

	if err != ErrCircularDependency {
		t.Errorf("Expected ErrCircularDependency, got %v", err)
	}
}

func TestFeatureCoordinator_ExecuteFeature_MissingDependency(t *testing.T) {
	// Test missing dependency detection
	workstreams := []WorkstreamNode{
		{
			ID:           "WS-001",
			Feature:      "F001",
			Status:       "pending",
			Dependencies: []string{"WS-NONEXISTENT"},
		},
	}

	loader := &MockWorkstreamLoader{
		workstreams: map[string][]WorkstreamNode{
			"F001": workstreams,
		},
	}

	executor := &MockWorkstreamExecutor{}
	saver := &MockCheckpointSaver{}

	coordinator := NewFeatureCoordinator(loader, executor, saver, 2)

	err := coordinator.ExecuteFeature("F001")

	// Error should wrap ErrMissingDependency
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	// Check if error message contains "missing dependency"
	if !containsString(err.Error(), "missing dependency") {
		t.Errorf("Expected error to contain 'missing dependency', got %v", err)
	}
}
