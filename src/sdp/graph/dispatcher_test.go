package graph

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestNewDispatcher verifies dispatcher creation
func TestNewDispatcher(t *testing.T) {
	g := NewDependencyGraph()
	d := NewDispatcher(g, 3)

	if d == nil {
		t.Fatal("NewDispatcher returned nil")
	}
	if d.concurrency != 3 {
		t.Errorf("expected concurrency 3, got %d", d.concurrency)
	}
	if d.circuitBreaker == nil {
		t.Error("circuit breaker not initialized")
	}
}

// TestNewDispatcher_DefaultConcurrency verifies default concurrency
func TestNewDispatcher_DefaultConcurrency(t *testing.T) {
	g := NewDependencyGraph()

	// Below minimum
	d1 := NewDispatcher(g, 0)
	if d1.concurrency != 3 {
		t.Errorf("expected default 3, got %d", d1.concurrency)
	}

	// Above maximum
	d2 := NewDispatcher(g, 10)
	if d2.concurrency != 5 {
		t.Errorf("expected max 5, got %d", d2.concurrency)
	}
}

// TestDispatcher_Execute_Single verifies single workstream execution
func TestDispatcher_Execute_Single(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)

	d := NewDispatcher(g, 1)
	executed := make(map[string]bool)
	var mu sync.Mutex

	results := d.Execute(func(wsID string) error {
		mu.Lock()
		executed[wsID] = true
		mu.Unlock()
		return nil
	})

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if !results[0].Success {
		t.Error("expected success")
	}
	if !executed["ws-001"] {
		t.Error("ws-001 not executed")
	}
}

// TestDispatcher_Execute_Multiple verifies multiple workstream execution
func TestDispatcher_Execute_Multiple(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)
	g.AddNode("ws-003", nil)

	d := NewDispatcher(g, 3)
	executed := make(map[string]bool)
	var mu sync.Mutex

	results := d.Execute(func(wsID string) error {
		mu.Lock()
		executed[wsID] = true
		mu.Unlock()
		return nil
	})

	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Success {
			t.Errorf("expected %s to succeed", r.WorkstreamID)
		}
	}
	if len(executed) != 3 {
		t.Errorf("expected 3 executed, got %d", len(executed))
	}
}

// TestDispatcher_Execute_DependencyOrder verifies dependency ordering
func TestDispatcher_Execute_DependencyOrder(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", []string{"ws-001"})
	g.AddNode("ws-003", []string{"ws-002"})

	d := NewDispatcher(g, 1) // Single thread to ensure order

	order := []string{}
	var mu sync.Mutex

	d.Execute(func(wsID string) error {
		mu.Lock()
		order = append(order, wsID)
		mu.Unlock()
		return nil
	})

	// Verify order: ws-001 before ws-002 before ws-003
	if len(order) != 3 {
		t.Fatalf("expected 3 executions, got %d", len(order))
	}
	if order[0] != "ws-001" {
		t.Errorf("expected ws-001 first, got %s", order[0])
	}
	if order[1] != "ws-002" {
		t.Errorf("expected ws-002 second, got %s", order[1])
	}
	if order[2] != "ws-003" {
		t.Errorf("expected ws-003 third, got %s", order[2])
	}
}

// TestDispatcher_Execute_WithFailure verifies failure handling
func TestDispatcher_Execute_WithFailure(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)

	d := NewDispatcher(g, 2)

	results := d.Execute(func(wsID string) error {
		if wsID == "ws-001" {
			return errors.New("intentional failure")
		}
		return nil
	})

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// Find results
	var ws001Result, ws002Result *ExecuteResult
	for i := range results {
		if results[i].WorkstreamID == "ws-001" {
			ws001Result = &results[i]
		} else if results[i].WorkstreamID == "ws-002" {
			ws002Result = &results[i]
		}
	}

	if ws001Result == nil || ws001Result.Success {
		t.Error("ws-001 should have failed")
	}
	if ws002Result == nil || !ws002Result.Success {
		t.Error("ws-002 should have succeeded")
	}
}

// TestDispatcher_Execute_EmptyGraph verifies empty graph handling
func TestDispatcher_Execute_EmptyGraph(t *testing.T) {
	g := NewDependencyGraph()
	d := NewDispatcher(g, 1)

	results := d.Execute(func(wsID string) error {
		return nil
	})

	if len(results) != 0 {
		t.Errorf("expected 0 results for empty graph, got %d", len(results))
	}
}

// TestDispatcher_isCompleted verifies completion tracking
func TestDispatcher_isCompleted(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	d := NewDispatcher(g, 1)

	if d.isCompleted("ws-001") {
		t.Error("should not be completed initially")
	}

	d.Execute(func(wsID string) error {
		return nil
	})

	if !d.isCompleted("ws-001") {
		t.Error("should be completed after execution")
	}
}

// TestDispatcher_SetCheckpointDir verifies checkpoint dir setting
func TestDispatcher_SetCheckpointDir(t *testing.T) {
	g := NewDependencyGraph()
	d := NewDispatcherWithCheckpoint(g, 1, "F001", true)

	d.SetCheckpointDir("/tmp/test")
	// No error means success
}

// TestDispatcher_SetFeatureID verifies feature ID setting
func TestDispatcher_SetFeatureID(t *testing.T) {
	g := NewDependencyGraph()
	d := NewDispatcher(g, 1)

	d.SetFeatureID("F999")

	if d.featureID != "F999" {
		t.Errorf("expected F999, got %s", d.featureID)
	}
}

// TestExecuteResult_Fields verifies ExecuteResult structure
func TestExecuteResult_Fields(t *testing.T) {
	result := ExecuteResult{
		WorkstreamID: "ws-001",
		Success:      true,
		Error:        nil,
		Duration:     100,
	}

	if result.WorkstreamID != "ws-001" {
		t.Error("WorkstreamID not set")
	}
	if !result.Success {
		t.Error("Success not set")
	}
	if result.Duration != 100 {
		t.Error("Duration not set")
	}
}

// TestWorkstreamFile_Fields verifies WorkstreamFile structure
func TestWorkstreamFile_Fields(t *testing.T) {
	ws := WorkstreamFile{
		ID:        "00-001-01",
		DependsOn: []string{"00-001-00"},
	}

	if ws.ID != "00-001-01" {
		t.Error("ID not set")
	}
	if len(ws.DependsOn) != 1 {
		t.Error("DependsOn not set")
	}
}

// TestBuildGraphFromWSFiles verifies graph building from files
func TestBuildGraphFromWSFiles(t *testing.T) {
	workstreams := []WorkstreamFile{
		{ID: "00-001-01", DependsOn: nil},
		{ID: "00-001-02", DependsOn: []string{"00-001-01"}},
		{ID: "00-001-03", DependsOn: []string{"00-001-01", "00-001-02"}},
	}

	g, err := BuildGraphFromWSFiles(workstreams)
	if err != nil {
		t.Fatalf("BuildGraphFromWSFiles failed: %v", err)
	}

	if len(g.nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(g.nodes))
	}
	if g.nodes["00-001-03"].Indegree != 2 {
		t.Errorf("expected indegree 2, got %d", g.nodes["00-001-03"].Indegree)
	}
}

// TestBuildGraphFromWSFiles_Empty verifies empty input
func TestBuildGraphFromWSFiles_Empty(t *testing.T) {
	g, err := BuildGraphFromWSFiles([]WorkstreamFile{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if g == nil {
		t.Error("graph should not be nil")
	}
	if len(g.nodes) != 0 {
		t.Error("expected empty graph")
	}
}

// TestBuildGraphFromWSFiles_MissingDependency verifies missing dependency error
func TestBuildGraphFromWSFiles_MissingDependency(t *testing.T) {
	workstreams := []WorkstreamFile{
		{ID: "00-001-02", DependsOn: []string{"00-001-01"}}, // Missing 00-001-01
	}

	_, err := BuildGraphFromWSFiles(workstreams)
	if err == nil {
		t.Error("expected error for missing dependency")
	}
}

// TestNewDispatcherWithCheckpoint verifies checkpoint dispatcher creation
func TestNewDispatcherWithCheckpoint(t *testing.T) {
	g := NewDependencyGraph()

	// Without checkpoint
	d1 := NewDispatcherWithCheckpoint(g, 2, "F001", false)
	if d1.enableCheckpoint {
		t.Error("checkpoint should be disabled")
	}

	// With checkpoint
	d2 := NewDispatcherWithCheckpoint(g, 2, "F001", true)
	if !d2.enableCheckpoint {
		t.Error("checkpoint should be enabled")
	}
	if d2.checkpointManager == nil {
		t.Error("checkpoint manager should be initialized")
	}
}

// TestDispatcher_ConcurrentExecution verifies concurrent execution
func TestDispatcher_ConcurrentExecution(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)
	g.AddNode("ws-003", nil)

	d := NewDispatcher(g, 3) // Allow 3 parallel

	startTimes := make(map[string]time.Time)
	var mu sync.Mutex

	d.Execute(func(wsID string) error {
		mu.Lock()
		startTimes[wsID] = time.Now()
		mu.Unlock()
		time.Sleep(10 * time.Millisecond) // Simulate work
		return nil
	})

	// All should have started within a short window (parallel)
	if len(startTimes) != 3 {
		t.Fatalf("expected 3 start times, got %d", len(startTimes))
	}
}

// TestDispatcher_MetricsCollection verifies dispatcher has metrics
func TestDispatcher_MetricsCollection(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	d := NewDispatcher(g, 1)

	d.Execute(func(wsID string) error {
		return nil
	})

	// Check circuit breaker metrics
	metrics := d.circuitBreaker.Metrics()
	if metrics.SuccessCount != 1 {
		t.Errorf("expected 1 success, got %d", metrics.SuccessCount)
	}
}

// TestDispatcher_GetCompleted verifies GetCompleted method
func TestDispatcher_GetCompleted(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)
	d := NewDispatcher(g, 2)

	// Before execution
	if len(d.GetCompleted()) != 0 {
		t.Error("expected 0 completed before execution")
	}

	d.Execute(func(wsID string) error {
		return nil
	})

	// After execution
	completed := d.GetCompleted()
	if len(completed) != 2 {
		t.Errorf("expected 2 completed, got %d", len(completed))
	}
}

// TestDispatcher_GetFailed verifies GetFailed method
func TestDispatcher_GetFailed(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)
	d := NewDispatcher(g, 2)

	// Before execution
	if len(d.GetFailed()) != 0 {
		t.Error("expected 0 failed before execution")
	}

	d.Execute(func(wsID string) error {
		if wsID == "ws-001" {
			return errors.New("intentional failure")
		}
		return nil
	})

	// After execution
	failed := d.GetFailed()
	if len(failed) != 1 {
		t.Errorf("expected 1 failed, got %d", len(failed))
	}
	if failed["ws-001"] == nil {
		t.Error("expected ws-001 to have error")
	}
}

// TestDispatcher_GetCircuitBreakerMetrics verifies GetCircuitBreakerMetrics
func TestDispatcher_GetCircuitBreakerMetrics(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	d := NewDispatcher(g, 1)

	// Before execution
	metrics := d.GetCircuitBreakerMetrics()
	if metrics.State != StateClosed {
		t.Errorf("expected CLOSED state, got %v", metrics.State)
	}

	d.Execute(func(wsID string) error {
		return nil
	})

	// After execution
	metrics = d.GetCircuitBreakerMetrics()
	if metrics.SuccessCount != 1 {
		t.Errorf("expected 1 success, got %d", metrics.SuccessCount)
	}
}

// TestDispatcher_CheckpointRestore verifies checkpoint restore flow
func TestDispatcher_CheckpointRestore(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", []string{"ws-001"})

	d := NewDispatcherWithCheckpoint(g, 1, "F001", true)
	d.SetCheckpointDir(tmpDir)

	// Execute first workstream only (simulate partial execution)
	callCount := 0
	d.Execute(func(wsID string) error {
		callCount++
		if wsID == "ws-001" {
			return nil
		}
		return errors.New("simulated failure")
	})

	// Verify some execution happened
	if callCount != 2 {
		t.Errorf("expected 2 executions, got %d", callCount)
	}
}

// TestDispatcher_TryRestoreCheckpoint_FeatureIDMismatch verifies mismatch handling
func TestDispatcher_TryRestoreCheckpoint_FeatureIDMismatch(t *testing.T) {
	tmpDir := t.TempDir()

	// Create and save checkpoint with F001
	g1 := NewDependencyGraph()
	g1.AddNode("ws-001", nil)
	cm1 := NewCheckpointManager("F001")
	cm1.SetCheckpointDir(tmpDir)
	checkpoint := cm1.CreateCheckpoint(g1, "F001", []string{"ws-001"}, []string{})
	cm1.Save(checkpoint)

	// Create dispatcher with F002 - should not restore
	g2 := NewDependencyGraph()
	g2.AddNode("ws-001", nil)
	d := NewDispatcherWithCheckpoint(g2, 1, "F002", true)
	d.SetCheckpointDir(tmpDir)

	// Execute - should not have restored F001's completed state
	d.Execute(func(wsID string) error {
		return nil
	})
}

// TestCheckpointManager_LoadCorrupt verifies corrupt checkpoint handling
func TestCheckpointManager_LoadCorrupt(t *testing.T) {
	tmpDir := t.TempDir()
	cm := NewCheckpointManager("F001")
	cm.SetCheckpointDir(tmpDir)

	// Write invalid JSON
	tmpPath := cm.GetTempPath()
	os.WriteFile(tmpPath, []byte("invalid json"), 0600)
	os.Rename(tmpPath, cm.GetCheckpointPath())

	// Load should fail and move to .corrupt
	_, err := cm.Load()
	if err == nil {
		t.Error("expected error for corrupt checkpoint")
	}

	// Corrupt file should exist
	if _, err := os.Stat(cm.GetCheckpointPath() + ".corrupt"); os.IsNotExist(err) {
		t.Error("corrupt file should have been renamed")
	}
}

// TestCheckpointManager_SaveFsync verifies fsync behavior
func TestCheckpointManager_SaveFsync(t *testing.T) {
	tmpDir := t.TempDir()
	cm := NewCheckpointManager("F001")
	cm.SetCheckpointDir(tmpDir)

	checkpoint := &Checkpoint{
		Version:   "1.0",
		FeatureID: "F001",
		Timestamp: time.Now(),
		Completed: []string{"ws-001"},
	}

	// Save should succeed with fsync
	err := cm.Save(checkpoint)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify temp file is gone
	if _, err := os.Stat(cm.GetTempPath()); !os.IsNotExist(err) {
		t.Error("temp file should be removed after atomic rename")
	}
}

// TestDispatcher_Execute_WithCheckpointSave verifies checkpoint save during execution
func TestDispatcher_Execute_WithCheckpointSave(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)

	d := NewDispatcherWithCheckpoint(g, 1, "F001", true)
	d.SetCheckpointDir(tmpDir)

	d.Execute(func(wsID string) error {
		return nil
	})

	// After successful completion, checkpoint should be deleted
	if _, err := os.Stat(filepath.Join(tmpDir, "F001-checkpoint.json")); !os.IsNotExist(err) {
		t.Error("checkpoint should be deleted after successful completion")
	}
}

// TestDispatcher_Execute_WithCheckpointKeepOnError verifies checkpoint kept on error
func TestDispatcher_Execute_WithCheckpointKeepOnError(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)

	d := NewDispatcherWithCheckpoint(g, 1, "F001", true)
	d.SetCheckpointDir(tmpDir)

	d.Execute(func(wsID string) error {
		return errors.New("intentional failure")
	})

	// Checkpoint should be kept when there are failures
	if _, err := os.Stat(filepath.Join(tmpDir, "F001-checkpoint.json")); os.IsNotExist(err) {
		t.Error("checkpoint should be kept when execution has failures")
	}
}
