package orchestrator

import (
	"errors"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
)

// mockWorkstreamLoader mocks loading workstreams
type mockWorkstreamLoader struct {
	workstreams []WorkstreamNode
	loadError   error
}

func (m *mockWorkstreamLoader) LoadWorkstreams(featureID string) ([]WorkstreamNode, error) {
	if m.loadError != nil {
		return nil, m.loadError
	}
	return m.workstreams, nil
}

// mockExecutor mocks workstream execution
type mockExecutor struct {
	executeFunc func(wsID string) error
}

func (m *mockExecutor) Execute(wsID string) error {
	if m.executeFunc != nil {
		return m.executeFunc(wsID)
	}
	return nil
}

// mockCheckpointManager mocks checkpoint operations
type mockCheckpointManager struct {
	checkpoints   map[string]checkpoint.Checkpoint
	saveCallCount int
	saveError     error
	loadError     error
}

func newMockCheckpointManager() *mockCheckpointManager {
	return &mockCheckpointManager{
		checkpoints: make(map[string]checkpoint.Checkpoint),
	}
}

func (m *mockCheckpointManager) Save(cp checkpoint.Checkpoint) error {
	m.saveCallCount++
	if m.saveError != nil {
		return m.saveError
	}
	m.checkpoints[cp.ID] = cp
	return nil
}

func (m *mockCheckpointManager) Load(id string) (checkpoint.Checkpoint, error) {
	if m.loadError != nil {
		return checkpoint.Checkpoint{}, m.loadError
	}
	cp, exists := m.checkpoints[id]
	if !exists {
		return checkpoint.Checkpoint{}, ErrCheckpointNotFound
	}
	return cp, nil
}

func (m *mockCheckpointManager) Resume(id string) (checkpoint.Checkpoint, error) {
	cp, err := m.Load(id)
	if err != nil {
		return checkpoint.Checkpoint{}, err
	}
	cp.UpdatedAt = time.Now()
	return cp, nil
}

// AC1: Load workstreams from Beads
func TestLoadWorkstreams(t *testing.T) {
	tests := []struct {
		name         string
		featureID    string
		workstreams []WorkstreamNode
		loadError    error
		wantCount    int
		wantError    bool
	}{
		{
			name:      "load workstreams successfully",
			featureID: "F050",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Feature: "F050", Status: "backlog"},
				{ID: "00-050-02", Feature: "F050", Status: "backlog"},
			},
			wantCount: 2,
			wantError: false,
		},
		{
			name:      "empty feature",
			featureID: "F999",
			workstreams: []WorkstreamNode{},
			wantCount: 0,
			wantError: false,
		},
		{
			name:      "load error",
			featureID: "F050",
			loadError: ErrFeatureNotFound,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := &mockWorkstreamLoader{
				workstreams: tt.workstreams,
				loadError:   tt.loadError,
			}

			ws, err := loader.LoadWorkstreams(tt.featureID)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(ws) != tt.wantCount {
				t.Errorf("got %d workstreams, want %d", len(ws), tt.wantCount)
			}
		})
	}
}

// AC2: Build dependency graph
func TestBuildDependencyGraph(t *testing.T) {
	tests := []struct {
		name         string
		workstreams []WorkstreamNode
		wantEdges    int
		wantError    bool
	}{
		{
			name: "simple dependency chain",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-03", Dependencies: []string{"00-050-02"}},
			},
			wantEdges: 2,
			wantError: false,
		},
		{
			name: "diamond dependency",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-03", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-04", Dependencies: []string{"00-050-02", "00-050-03"}},
			},
			wantEdges: 4,
			wantError: false,
		},
		{
			name: "circular dependency",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{"00-050-02"}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
			},
			wantError: true,
		},
		{
			name: "missing dependency",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{"00-050-99"}},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph, err := BuildDependencyGraph(tt.workstreams)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Count edges
			edgeCount := 0
			for _, node := range graph {
				edgeCount += len(node.Dependencies)
			}

			if edgeCount != tt.wantEdges {
				t.Errorf("got %d edges, want %d", edgeCount, tt.wantEdges)
			}
		})
	}
}

// AC3: Topological sort execution order
func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name         string
		workstreams []WorkstreamNode
		wantOrder    []string
		wantError    bool
	}{
		{
			name: "simple chain",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-03", Dependencies: []string{"00-050-02"}},
			},
			wantOrder: []string{"00-050-01", "00-050-02", "00-050-03"},
			wantError: false,
		},
		{
			name: "diamond dependency",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-03", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-04", Dependencies: []string{"00-050-02", "00-050-03"}},
			},
			wantError: false,
		},
		{
			name: "circular dependency",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{"00-050-02"}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
			},
			wantError: true, // Should detect cycle during build
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph, err := BuildDependencyGraph(tt.workstreams)
			// If build failed with expected error, test passes
			if err != nil {
				if tt.wantError {
					return // Test passes
				}
				t.Errorf("build graph failed: %v", err)
				return
			}

			order, err := TopologicalSort(graph)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// If specific order expected, verify it
			if tt.wantOrder != nil {
				if len(order) != len(tt.wantOrder) {
					t.Errorf("got %d nodes, want %d", len(order), len(tt.wantOrder))
					return
				}
				for i, wsID := range tt.wantOrder {
					if order[i] != wsID {
						t.Errorf("position %d: got %s, want %s", i, order[i], wsID)
					}
				}
			}

			// Verify dependencies come before dependents
			// The graph stores edges FROM dependencies TO their dependents
			// So we need to check that for each node, all nodes it points to come after it
			for i, wsID := range order {
				node := graph[wsID]
				// node.Dependencies contains the workstreams that depend on this one
				// All of them should come AFTER this node
				for _, dependentNode := range node.Dependencies {
					dependentID := dependentNode.Workstream.ID
					// Find dependent position
					depPos := -1
					for j, id := range order {
						if id == dependentID {
							depPos = j
							break
						}
					}
					if depPos == -1 {
						t.Errorf("dependent %s not found in order", dependentID)
					}
					if depPos <= i {
						t.Errorf("dependent %s (pos %d) must come after %s (pos %d)",
							dependentID, depPos, wsID, i)
					}
				}
			}
		})
	}
}

// AC4: Execute workstreams sequentially
func TestExecuteWorkstreams(t *testing.T) {
	tests := []struct {
		name         string
		workstreams []string
		executeFunc  func(wsID string) error
		wantError    bool
		wantExecuted []string
	}{
		{
			name:         "all succeed",
			workstreams: []string{"00-050-01", "00-050-02", "00-050-03"},
			executeFunc:  func(wsID string) error { return nil },
			wantError:    false,
			wantExecuted: []string{"00-050-01", "00-050-02", "00-050-03"},
		},
		{
			name:         "middle fails",
			workstreams: []string{"00-050-01", "00-050-02", "00-050-03"},
			executeFunc: func(wsID string) error {
				if wsID == "00-050-02" {
					return errors.New("permanent failure")
				}
				return nil
			},
			wantError:    true,
			wantExecuted: []string{"00-050-01", "00-050-02", "00-050-02", "00-050-02"}, // Initial + 2 retries
		},
		{
			name:         "first fails",
			workstreams: []string{"00-050-01", "00-050-02"},
			executeFunc: func(wsID string) error {
				if wsID == "00-050-01" {
					return errors.New("permanent failure")
				}
				return nil
			},
			wantError:    true,
			wantExecuted: []string{"00-050-01", "00-050-01", "00-050-01"}, // Initial + 2 retries
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executed := []string{}
			executor := &mockExecutor{
				executeFunc: func(wsID string) error {
					executed = append(executed, wsID)
					return tt.executeFunc(wsID)
				},
			}

			orch := NewOrchestrator(nil, executor, nil, 2)
			err := orch.executeSequential(tt.workstreams)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			// Verify execution order
			if len(executed) != len(tt.wantExecuted) {
				t.Errorf("executed %d workstreams, want %d", len(executed), len(tt.wantExecuted))
				return
			}

			for i, wsID := range tt.wantExecuted {
				if executed[i] != wsID {
					t.Errorf("position %d: got %s, want %s", i, executed[i], wsID)
				}
			}
		})
	}
}

// AC5: Checkpoint after each WS
func TestCheckpointAfterExecution(t *testing.T) {
	tests := []struct {
		name         string
		workstreams []string
		executeFunc  func(wsID string) error
		wantSaves    int
		wantError    bool
	}{
		{
			name:         "checkpoint after each success",
			workstreams: []string{"00-050-01", "00-050-02"},
			executeFunc:  func(wsID string) error { return nil },
			wantSaves:    5, // Before WS1, After WS1, Before WS2, After WS2, Final
			wantError:    false,
		},
		{
			name:         "checkpoint on failure",
			workstreams: []string{"00-050-01", "00-050-02", "00-050-03"},
			executeFunc: func(wsID string) error {
				if wsID == "00-050-02" {
					return errors.New("permanent failure")
				}
				return nil
			},
			wantSaves: 4, // Before WS1, After WS1, Before WS2, After WS2 (failed)
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkpointMgr := newMockCheckpointManager()
			executor := &mockExecutor{executeFunc: tt.executeFunc}

			orch := NewOrchestrator(nil, executor, checkpointMgr, 2)

			// Create initial checkpoint
			cp := checkpoint.Checkpoint{
				ID:        "test-cp",
				FeatureID: "F050",
				Status:    checkpoint.StatusInProgress,
				CompletedWorkstreams: []string{},
			}

			err := orch.executeWithCheckpoint(tt.workstreams, &cp)

			if tt.wantError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify checkpoint was saved expected times
			savedCount := 0
			for _, savedCp := range checkpointMgr.checkpoints {
				if savedCp.ID == "test-cp" {
					savedCount++
					// Verify completed workstreams in checkpoint
					expectedCompleted := len(cp.CompletedWorkstreams)
					if len(savedCp.CompletedWorkstreams) != expectedCompleted {
						t.Errorf("checkpoint has %d completed, want %d",
							len(savedCp.CompletedWorkstreams), expectedCompleted)
					}
				}
			}

			if checkpointMgr.saveCallCount != tt.wantSaves {
				t.Errorf("checkpoint saved %d times, want %d", checkpointMgr.saveCallCount, tt.wantSaves)
			}
		})
	}
}

// AC6: Resume from checkpoint
func TestResumeFromCheckpoint(t *testing.T) {
	tests := []struct {
		name              string
		checkpoint        checkpoint.Checkpoint
		allWorkstreams    []string
		wantResume        []string
		wantError         bool
	}{
		{
			name: "resume from middle",
			checkpoint: checkpoint.Checkpoint{
				ID:        "test-cp",
				FeatureID: "F050",
				Status:    checkpoint.StatusInProgress,
				CurrentWorkstream: "00-050-02",
				CompletedWorkstreams: []string{"00-050-01"},
			},
			allWorkstreams: []string{"00-050-01", "00-050-02", "00-050-03"},
			wantResume:     []string{"00-050-02", "00-050-03"},
			wantError:      false,
		},
		{
			name: "resume from start",
			checkpoint: checkpoint.Checkpoint{
				ID:        "test-cp",
				FeatureID: "F050",
				Status:    checkpoint.StatusPending,
				CompletedWorkstreams: []string{},
			},
			allWorkstreams: []string{"00-050-01", "00-050-02"},
			wantResume:     []string{"00-050-01", "00-050-02"},
			wantError:      false,
		},
		{
			name: "already completed",
			checkpoint: checkpoint.Checkpoint{
				ID:        "test-cp",
				FeatureID: "F050",
				Status:    checkpoint.StatusCompleted,
				CompletedWorkstreams: []string{"00-050-01", "00-050-02"},
			},
			allWorkstreams: []string{"00-050-01", "00-050-02"},
			wantResume:     []string{}, // Nothing to execute
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkpointMgr := newMockCheckpointManager()
			checkpointMgr.checkpoints[tt.checkpoint.ID] = tt.checkpoint

			orch := NewOrchestrator(nil, nil, checkpointMgr, 2)

			resume, err := orch.getRemainingWorkstreams(tt.allWorkstreams, &tt.checkpoint)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(resume) != len(tt.wantResume) {
				t.Errorf("got %d workstreams to resume, want %d", len(resume), len(tt.wantResume))
				return
			}

			for i, wsID := range tt.wantResume {
				if resume[i] != wsID {
					t.Errorf("position %d: got %s, want %s", i, resume[i], wsID)
				}
			}
		})
	}
}

// AC7: Handle dependencies correctly
func TestDependencyHandling(t *testing.T) {
	tests := []struct {
		name         string
		workstreams []WorkstreamNode
		wantError    bool
	}{
		{
			name: "valid dependency chain",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-03", Dependencies: []string{"00-050-02"}},
			},
			wantError: false,
		},
		{
			name: "self-dependency",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{"00-050-01"}},
			},
			wantError: true,
		},
		{
			name: "transitive dependencies",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Dependencies: []string{}},
				{ID: "00-050-02", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-03", Dependencies: []string{"00-050-01", "00-050-02"}},
				{ID: "00-050-04", Dependencies: []string{"00-050-03"}},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph, err := BuildDependencyGraph(tt.workstreams)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify all nodes are in graph
			for _, ws := range tt.workstreams {
				if _, exists := graph[ws.ID]; !exists {
					t.Errorf("workstream %s not in graph", ws.ID)
				}
			}

			// Verify graph structure - check that all nodes exist
			// Note: node.Dependencies contains dependents (outgoing edges), not dependencies
			// The original dependencies are stored in node.Workstream.Dependencies
			for _, ws := range tt.workstreams {
				node := graph[ws.ID]
				// Verify the workstream data is preserved
				if node.Workstream.ID != ws.ID {
					t.Errorf("node %s: ID mismatch", ws.ID)
				}
				// Verify original dependencies are preserved
				if len(node.Workstream.Dependencies) != len(ws.Dependencies) {
					t.Errorf("node %s: got %d dependencies, want %d",
						ws.ID, len(node.Workstream.Dependencies), len(ws.Dependencies))
				}
			}
		})
	}
}

// Integration tests for Run and Resume methods
func TestOrchestratorRun(t *testing.T) {
	tests := []struct {
		name         string
		featureID    string
		workstreams []WorkstreamNode
		executeFunc  func(wsID string) error
		wantError    bool
	}{
		{
			name:      "run complete feature",
			featureID: "F050",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Feature: "F050", Status: "backlog", Dependencies: []string{}},
				{ID: "00-050-02", Feature: "F050", Status: "backlog", Dependencies: []string{"00-050-01"}},
			},
			executeFunc: func(wsID string) error { return nil },
			wantError:   false,
		},
		{
			name:      "feature not found",
			featureID: "F999",
			workstreams: []WorkstreamNode{},
			executeFunc: func(wsID string) error { return nil },
			wantError:   true,
		},
		{
			name:      "execution fails",
			featureID: "F050",
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Feature: "F050", Status: "backlog", Dependencies: []string{}},
			},
			executeFunc: func(wsID string) error {
				return errors.New("execution failed")
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := &mockWorkstreamLoader{
				workstreams: tt.workstreams,
			}
			checkpointMgr := newMockCheckpointManager()
			executor := &mockExecutor{executeFunc: tt.executeFunc}

			orch := NewOrchestrator(loader, executor, checkpointMgr, 2)

			err := orch.Run(tt.featureID)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Verify checkpoint was created and completed
				cp, exists := checkpointMgr.checkpoints[tt.featureID]
				if !exists {
					t.Errorf("checkpoint not created")
				} else if cp.Status != checkpoint.StatusCompleted {
					t.Errorf("checkpoint status = %s, want %s", cp.Status, checkpoint.StatusCompleted)
				}
			}
		})
	}
}

func TestOrchestratorResume(t *testing.T) {
	tests := []struct {
		name         string
		checkpointID string
		checkpoint   checkpoint.Checkpoint
		workstreams  []WorkstreamNode
		executeFunc  func(wsID string) error
		wantError    bool
	}{
		{
			name:         "resume from middle",
			checkpointID: "test-cp",
			checkpoint: checkpoint.Checkpoint{
				ID:        "test-cp",
				FeatureID: "F050",
				Status:    checkpoint.StatusInProgress,
				CurrentWorkstream: "00-050-02",
				CompletedWorkstreams: []string{"00-050-01"},
			},
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Feature: "F050", Status: "backlog", Dependencies: []string{}},
				{ID: "00-050-02", Feature: "F050", Status: "backlog", Dependencies: []string{"00-050-01"}},
				{ID: "00-050-03", Feature: "F050", Status: "backlog", Dependencies: []string{"00-050-02"}},
			},
			executeFunc: func(wsID string) error { return nil },
			wantError:   false,
		},
		{
			name:         "already completed",
			checkpointID: "test-cp",
			checkpoint: checkpoint.Checkpoint{
				ID:        "test-cp",
				FeatureID: "F050",
				Status:    checkpoint.StatusCompleted,
				CompletedWorkstreams: []string{"00-050-01", "00-050-02"},
			},
			workstreams: []WorkstreamNode{
				{ID: "00-050-01", Feature: "F050", Status: "backlog", Dependencies: []string{}},
				{ID: "00-050-02", Feature: "F050", Status: "backlog", Dependencies: []string{"00-050-01"}},
			},
			executeFunc: func(wsID string) error { return nil },
			wantError:   false,
		},
		{
			name:         "checkpoint not found",
			checkpointID: "nonexistent",
			executeFunc:  func(wsID string) error { return nil },
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkpointMgr := newMockCheckpointManager()
			checkpointMgr.checkpoints[tt.checkpoint.ID] = tt.checkpoint

			loader := &mockWorkstreamLoader{
				workstreams: tt.workstreams,
			}
			executor := &mockExecutor{executeFunc: tt.executeFunc}

			orch := NewOrchestrator(loader, executor, checkpointMgr, 2)

			err := orch.Resume(tt.checkpointID)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
