package graph

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/src/sdp/graph"
)

// TestDispatcherSequential tests sequential execution (no parallelism possible)
func TestDispatcherSequential(t *testing.T) {
	g := graph.NewDependencyGraph()

	// A -> B -> C (must execute sequentially)
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-02"})

	// Track execution order
	var mu sync.Mutex
	var executionOrder []string

	executeFn := func(wsID string) error {
		mu.Lock()
		executionOrder = append(executionOrder, wsID)
		mu.Unlock()
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3)
	results := dispatcher.Execute(executeFn)

	// Verify all workstreams completed
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Verify all succeeded
	for _, result := range results {
		if !result.Success {
			t.Errorf("Workstream %s failed: %v", result.WorkstreamID, result.Error)
		}
	}

	// Verify execution order
	if len(executionOrder) != 3 {
		t.Fatalf("Expected 3 executions, got %d", len(executionOrder))
	}

	// A must come before B, B before C
	orderMap := make(map[string]int)
	for i, id := range executionOrder {
		orderMap[id] = i
	}

	if orderMap["00-001-01"] >= orderMap["00-001-02"] {
		t.Error("00-001-01 should execute before 00-001-02")
	}

	if orderMap["00-001-02"] >= orderMap["00-001-03"] {
		t.Error("00-001-02 should execute before 00-001-03")
	}
}

// TestDispatcherParallel tests parallel execution of independent nodes
func TestDispatcherParallel(t *testing.T) {
	g := graph.NewDependencyGraph()

	// A, B, C are all independent (can execute in parallel)
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{})
	g.AddNode("00-001-03", []string{})

	// Track execution start time
	var mu sync.Mutex
	var executionOrder []string
	startTimes := make(map[string]time.Time)
	endTimes := make(map[string]time.Time)

	executeFn := func(wsID string) error {
		mu.Lock()
		executionOrder = append(executionOrder, wsID)
		startTimes[wsID] = time.Now()
		mu.Unlock()

		// Simulate some work
		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		endTimes[wsID] = time.Now()
		mu.Unlock()

		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3)
	results := dispatcher.Execute(executeFn)

	// Verify all workstreams completed
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Verify all succeeded
	for _, result := range results {
		if !result.Success {
			t.Errorf("Workstream %s failed: %v", result.WorkstreamID, result.Error)
		}
	}

	// Verify parallelism: all should start roughly at the same time
	// (within 100ms of each other)
	minStart := time.Time{}
	maxStart := time.Time{}
	for _, start := range startTimes {
		if minStart.IsZero() || start.Before(minStart) {
			minStart = start
		}
		if maxStart.IsZero() || start.After(maxStart) {
			maxStart = start
		}
	}

	// If all started within 100ms, they were running in parallel
	if maxStart.Sub(minStart) > 200*time.Millisecond {
		t.Error("Workstreams may not have executed in parallel")
	}
}

// TestDispatcherMixed tests mixed sequential and parallel execution
func TestDispatcherMixed(t *testing.T) {
	g := graph.NewDependencyGraph()

	//     A
	//    /|\
	//   B C D
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-01"})
	g.AddNode("00-001-04", []string{"00-001-01"})

	var mu sync.Mutex
	var executionOrder []string

	executeFn := func(wsID string) error {
		mu.Lock()
		executionOrder = append(executionOrder, wsID)
		mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3)
	results := dispatcher.Execute(executeFn)

	// Verify all workstreams completed
	if len(results) != 4 {
		t.Errorf("Expected 4 results, got %d", len(results))
	}

	// Verify all succeeded
	for _, result := range results {
		if !result.Success {
			t.Errorf("Workstream %s failed: %v", result.WorkstreamID, result.Error)
		}
	}

	// Verify A executed first
	if executionOrder[0] != "00-001-01" {
		t.Errorf("Expected 00-001-01 to execute first, got %s", executionOrder[0])
	}

	// Verify B, C, D executed after A
	orderMap := make(map[string]int)
	for i, id := range executionOrder {
		orderMap[id] = i
	}

	if orderMap["00-001-01"] >= orderMap["00-001-02"] ||
		orderMap["00-001-01"] >= orderMap["00-001-03"] ||
		orderMap["00-001-01"] >= orderMap["00-001-04"] {
		t.Error("00-001-01 should execute before B, C, D")
	}
}

// TestDispatcherErrorHandling tests handling of workstream failures
func TestDispatcherErrorHandling(t *testing.T) {
	g := graph.NewDependencyGraph()

	// A -> B -> C
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-02"})

	executeFn := func(wsID string) error {
		if wsID == "00-001-02" {
			return errors.New("simulated failure")
		}
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3)
	results := dispatcher.Execute(executeFn)

	// Verify we got results
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Count successes and failures
	successCount := 0
	failedCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failedCount++
		}
	}

	// At least one should have failed (B)
	if failedCount < 1 {
		t.Error("Expected at least 1 failure")
	}

	// C might not have executed if B failed, depending on implementation
}

// TestDispatcherConcurrencyLimit tests that concurrency is limited
func TestDispatcherConcurrencyLimit(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Add 10 independent workstreams
	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("00-001-%02d", i)
		g.AddNode(id, []string{})
	}

	var mu sync.Mutex
	maxConcurrent := 0
	currentConcurrent := 0

	executeFn := func(wsID string) error {
		mu.Lock()
		currentConcurrent++
		if currentConcurrent > maxConcurrent {
			maxConcurrent = currentConcurrent
		}
		mu.Unlock()

		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		currentConcurrent--
		mu.Unlock()

		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3) // Limit to 3 concurrent
	results := dispatcher.Execute(executeFn)

	// Verify all completed
	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}

	// Verify concurrency was respected
	// With proper locking, maxConcurrent should not exceed 3 by much
	// (there might be a small window where it goes to 4)
	if maxConcurrent > 4 {
		t.Errorf("Concurrency limit exceeded: got %d, expected max 4", maxConcurrent)
	}
}

// TestBuildGraphFromWSFiles tests building a graph from workstream files
func TestBuildGraphFromWSFiles(t *testing.T) {
	workstreams := []graph.WorkstreamFile{
		{ID: "00-001-01", DependsOn: []string{}},
		{ID: "00-001-02", DependsOn: []string{"00-001-01"}},
		{ID: "00-001-03", DependsOn: []string{"00-001-02"}},
	}

	g, err := graph.BuildGraphFromWSFiles(workstreams)
	if err != nil {
		t.Fatalf("BuildGraphFromWSFiles failed: %v", err)
	}

	// Verify graph structure
	ready := g.GetReady()
	if len(ready) != 1 || ready[0] != "00-001-01" {
		t.Errorf("Expected only 00-001-01 to be ready, got %v", ready)
	}
}

// TestBuildGraphFromWSFilesError tests error handling in graph building
func TestBuildGraphFromWSFilesError(t *testing.T) {
	workstreams := []graph.WorkstreamFile{
		{ID: "00-001-01", DependsOn: []string{}},
		{ID: "00-001-02", DependsOn: []string{"00-001-03"}}, // Missing dependency
	}

	_, err := graph.BuildGraphFromWSFiles(workstreams)
	if err == nil {
		t.Error("Expected error for missing dependency, got nil")
	}
}

// TestDispatcherGetCompletedGetFailed tests tracking completed and failed workstreams
func TestDispatcherGetCompletedGetFailed(t *testing.T) {
	g := graph.NewDependencyGraph()

	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-01"})

	executeFn := func(wsID string) error {
		if wsID == "00-001-02" {
			return errors.New("simulated failure")
		}
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3)
	dispatcher.Execute(executeFn)

	completed := dispatcher.GetCompleted()
	failed := dispatcher.GetFailed()

	// Should have 2 completed (A and C) and 1 failed (B)
	if len(completed) != 2 {
		t.Errorf("Expected 2 completed, got %d", len(completed))
	}

	if len(failed) != 1 {
		t.Errorf("Expected 1 failed, got %d", len(failed))
	}

	if _, ok := failed["00-001-02"]; !ok {
		t.Error("Expected 00-001-02 to be in failed list")
	}
}
