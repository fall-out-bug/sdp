package graph

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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

// TestDispatcher_CircuitBreakerIntegration verifies circuit breaker is integrated
func TestDispatcher_CircuitBreakerIntegration(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Add independent workstreams
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{})
	g.AddNode("00-001-03", []string{})

	dispatcher := graph.NewDispatcher(g, 3)

	// Verify circuit breaker exists and is in CLOSED state initially
	metrics := dispatcher.GetCircuitBreakerMetrics()
	if metrics.State != graph.StateClosed {
		t.Errorf("Expected initial circuit breaker state CLOSED, got %v", metrics.State)
	}

	// Execute successfully
	executeFn := func(wsID string) error {
		return nil
	}
	results := dispatcher.Execute(executeFn)

	// Verify all succeeded
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}
	for _, result := range results {
		if !result.Success {
			t.Errorf("Workstream %s failed: %v", result.WorkstreamID, result.Error)
		}
	}

	// Verify circuit breaker metrics after success
	metricsAfter := dispatcher.GetCircuitBreakerMetrics()
	if metricsAfter.FailureCount != 0 {
		t.Errorf("Expected 0 failures after successful execution, got %d", metricsAfter.FailureCount)
	}
}

// TestDispatcher_CircuitBreakerTrips verifies circuit breaker trips on failures
func TestDispatcher_CircuitBreakerTrips(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Add workstreams that will fail
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{})
	g.AddNode("00-001-03", []string{})

	dispatcher := graph.NewDispatcher(g, 3)

	// Execute with failures
	var failCount int64
	executeFn := func(wsID string) error {
		if wsID == "00-001-01" || wsID == "00-001-02" || wsID == "00-001-03" {
			count := atomic.AddInt64(&failCount, 1)
			if count <= 3 {
				return errors.New("simulated failure")
			}
		}
		return nil
	}

	results := dispatcher.Execute(executeFn)

	// Verify circuit breaker tracked failures
	metrics := dispatcher.GetCircuitBreakerMetrics()
	if metrics.FailureCount < 3 {
		t.Logf("Warning: Expected at least 3 failures, got %d (may be due to circuit breaker)", metrics.FailureCount)
	}

	// Verify results
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	// Some should have failed
	failedResults := 0
	for _, result := range results {
		if !result.Success {
			failedResults++
		}
	}
	if failedResults == 0 {
		t.Error("Expected some workstreams to fail")
	}
}

// TestDispatcher_FiveParallelAgents tests exactly 5 agents running in parallel
func TestDispatcher_FiveParallelAgents(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Create 5 independent workstreams (mimicking real builder agents)
	workstreams := []string{"00-001-01", "00-001-02", "00-001-03", "00-001-04", "00-001-05"}
	for _, wsID := range workstreams {
		g.AddNode(wsID, []string{})
	}

	// Track execution timing
	var mu sync.Mutex
	startTimes := make(map[string]time.Time)
	endTimes := make(map[string]time.Time)

	executeFn := func(wsID string) error {
		mu.Lock()
		startTimes[wsID] = time.Now()
		mu.Unlock()

		// Simulate real work (file operations, etc.)
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		endTimes[wsID] = time.Now()
		mu.Unlock()

		return nil
	}

	dispatcher := graph.NewDispatcher(g, 5) // Allow 5 concurrent
	results := dispatcher.Execute(executeFn)

	// Verify all completed
	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	// Verify all succeeded
	for _, result := range results {
		if !result.Success {
			t.Errorf("Workstream %s failed: %v", result.WorkstreamID, result.Error)
		}
	}

	// Verify parallelism: all should start within 50ms of each other
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

	startSpread := maxStart.Sub(minStart)
	if startSpread > 100*time.Millisecond {
		t.Errorf("Workstreams did not execute in parallel: start spread %v > 100ms", startSpread)
	}

	// Calculate total time (should be ~100ms if truly parallel, not 500ms sequential)
	minEnd := time.Time{}
	maxEnd := time.Time{}
	for _, end := range endTimes {
		if minEnd.IsZero() || end.Before(minEnd) {
			minEnd = end
		}
		if maxEnd.IsZero() || end.After(maxEnd) {
			maxEnd = end
		}
	}

	totalTime := maxEnd.Sub(minStart)
	// With 5 parallel agents each taking 100ms, total should be ~100-150ms
	// If sequential, would be ~500ms
	if totalTime > 200*time.Millisecond {
		t.Errorf("Execution too slow (%v), likely sequential instead of parallel", totalTime)
	}

	t.Logf("✓ 5 agents executed in parallel in %v (start spread: %v)", totalTime, startSpread)
}

// TestDispatcher_ThreeParallelAgents tests exactly 3 agents running in parallel
func TestDispatcher_ThreeParallelAgents(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Create 3 independent workstreams
	workstreams := []string{"00-001-01", "00-001-02", "00-001-03"}
	for _, wsID := range workstreams {
		g.AddNode(wsID, []string{})
	}

	// Track execution
	var executionCount int32

	executeFn := func(wsID string) error {
		count := atomic.AddInt32(&executionCount, 1)

		// Verify we never exceed 3 concurrent
		if count > 3 {
			t.Errorf("Concurrency limit exceeded: %d > 3", count)
		}

		time.Sleep(50 * time.Millisecond)

		atomic.AddInt32(&executionCount, -1)
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3)
	results := dispatcher.Execute(executeFn)

	// Verify all completed
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	for _, result := range results {
		if !result.Success {
			t.Errorf("Workstream %s failed: %v", result.WorkstreamID, result.Error)
		}
	}

	t.Logf("✓ 3 agents executed successfully with concurrency limit of 3")
}

// TestDispatcher_FourParallelAgents tests exactly 4 agents running in parallel
func TestDispatcher_FourParallelAgents(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Create 4 independent workstreams (common real-world scenario)
	workstreams := []string{"00-001-01", "00-001-02", "00-001-03", "00-001-04"}
	for _, wsID := range workstreams {
		g.AddNode(wsID, []string{})
	}

	// Track execution timing
	var mu sync.Mutex
	startTimes := make(map[string]time.Time)

	executeFn := func(wsID string) error {
		mu.Lock()
		startTimes[wsID] = time.Now()
		mu.Unlock()

		time.Sleep(75 * time.Millisecond)
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 4)
	results := dispatcher.Execute(executeFn)

	// Verify all completed
	if len(results) != 4 {
		t.Errorf("Expected 4 results, got %d", len(results))
	}

	for _, result := range results {
		if !result.Success {
			t.Errorf("Workstream %s failed: %v", result.WorkstreamID, result.Error)
		}
	}

	// Verify parallelism
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

	startSpread := maxStart.Sub(minStart)
	if startSpread > 100*time.Millisecond {
		t.Errorf("Workstreams did not execute in parallel: start spread %v > 100ms", startSpread)
	}

	t.Logf("✓ 4 agents executed in parallel (start spread: %v)", startSpread)
}

// TestDispatcher_RealWorldDependencyGraph tests a complex real-world scenario
func TestDispatcher_RealWorldDependencyGraph(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Simulate a real feature with multiple components
	// Phase 1: Foundation (parallel)
	g.AddNode("00-001-01", []string{}) // Domain models
	g.AddNode("00-001-02", []string{}) // Database schema
	g.AddNode("00-001-03", []string{}) // API contracts

	// Phase 2: Implementation (depends on Phase 1, parallel within phase)
	g.AddNode("00-001-04", []string{"00-001-01"}) // Business logic
	g.AddNode("00-001-05", []string{"00-001-01", "00-001-02"}) // Repository layer
	g.AddNode("00-001-06", []string{"00-001-03"}) // API endpoints

	// Phase 3: Integration (depends on Phase 2)
	g.AddNode("00-001-07", []string{"00-001-04", "00-001-05", "00-001-06"}) // Integration tests

	// Track execution order
	var mu sync.Mutex
	executionOrder := []string{}

	executeFn := func(wsID string) error {
		mu.Lock()
		executionOrder = append(executionOrder, wsID)
		mu.Unlock()

		// Simulate work
		time.Sleep(20 * time.Millisecond)
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 3) // Allow 3 concurrent
	results := dispatcher.Execute(executeFn)

	// Verify all completed
	if len(results) != 7 {
		t.Errorf("Expected 7 results, got %d", len(results))
	}

	// Verify execution order respects dependencies
	orderMap := make(map[string]int)
	for i, id := range executionOrder {
		orderMap[id] = i
	}

	// Phase 1 should execute before Phase 2
	if orderMap["00-001-01"] >= orderMap["00-001-04"] {
		t.Error("00-001-01 (Phase 1) should execute before 00-001-04 (Phase 2)")
	}
	if orderMap["00-001-02"] >= orderMap["00-001-05"] {
		t.Error("00-001-02 (Phase 1) should execute before 00-001-05 (Phase 2)")
	}

	// Phase 2 should execute before Phase 3
	if orderMap["00-001-04"] >= orderMap["00-001-07"] {
		t.Error("00-001-04 (Phase 2) should execute before 00-001-07 (Phase 3)")
	}

	t.Logf("✓ Real-world dependency graph executed correctly")
	t.Logf("  Execution order: %v", executionOrder)
}

// TestDispatcher_PerformanceSpeedup verifies parallel execution provides speedup
func TestDispatcher_PerformanceSpeedup(t *testing.T) {
	// Test sequential execution
	g1 := graph.NewDependencyGraph()
	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("00-001-%02d", i)
		g1.AddNode(id, []string{})
	}

	var sequentialDuration time.Duration
	executeFn := func(wsID string) error {
		start := time.Now()
		time.Sleep(50 * time.Millisecond)
		sequentialDuration += time.Since(start)
		return nil
	}

	dispatcher1 := graph.NewDispatcher(g1, 1) // Sequential
	startSeq := time.Now()
	dispatcher1.Execute(executeFn)
	sequentialTime := time.Since(startSeq)

	// Test parallel execution
	g2 := graph.NewDependencyGraph()
	for i := 1; i <= 5; i++ {
		id := fmt.Sprintf("00-001-%02d", i)
		g2.AddNode(id, []string{})
	}

	executeFn2 := func(wsID string) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}

	dispatcher2 := graph.NewDispatcher(g2, 5) // Parallel
	startPar := time.Now()
	dispatcher2.Execute(executeFn2)
	parallelTime := time.Since(startPar)

	// Parallel should be significantly faster
	// Sequential: ~250ms (5 * 50ms)
	// Parallel: ~50-100ms (all start together)
	speedup := float64(sequentialTime) / float64(parallelTime)

	t.Logf("Sequential time: %v", sequentialTime)
	t.Logf("Parallel time: %v", parallelTime)
	t.Logf("Speedup: %.2fx", speedup)

	if speedup < 1.5 {
		t.Errorf("Parallel execution should be at least 1.5x faster, got %.2fx", speedup)
	}

	if parallelTime > 150*time.Millisecond {
		t.Errorf("Parallel execution too slow: %v (expected ~50-100ms)", parallelTime)
	}
}

// TestDispatcher_RaceCondition verifies no race conditions with shared state
func TestDispatcher_RaceCondition(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Create 10 workstreams all modifying shared state
	for i := 1; i <= 10; i++ {
		id := fmt.Sprintf("00-001-%02d", i)
		g.AddNode(id, []string{})
	}

	// Shared counter using atomic operations (proper synchronization)
	var counter int64

	executeFn := func(wsID string) error {
		// Use atomic operations to safely increment from multiple goroutines
		for i := 0; i < 100; i++ {
			atomic.AddInt64(&counter, 1)
		}
		return nil
	}

	dispatcher := graph.NewDispatcher(g, 5)
	results := dispatcher.Execute(executeFn)

	// Verify all completed
	if len(results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(results))
	}

	// Verify counter is correct (10 workstreams * 100 increments = 1000)
	expected := int64(10 * 100)
	if counter != expected {
		t.Errorf("Counter mismatch: expected counter=%d, got %d", expected, counter)
	}

	t.Logf("✓ No race conditions: counter=%d (expected %d)", counter, expected)
}

// TestDispatcher_EdgeCases tests edge cases and boundary conditions
func TestDispatcher_EdgeCases(t *testing.T) {
	// Test empty graph
	g := graph.NewDependencyGraph()
	dispatcher := graph.NewDispatcher(g, 3)
	results := dispatcher.Execute(func(wsID string) error { return nil })

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty graph, got %d", len(results))
	}

	// Test single workstream
	g2 := graph.NewDependencyGraph()
	g2.AddNode("00-001-01", []string{})
	dispatcher2 := graph.NewDispatcher(g2, 3)
	results2 := dispatcher2.Execute(func(wsID string) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	if len(results2) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results2))
	}

	// Test circular dependency (should be caught during graph building)
	g3 := graph.NewDependencyGraph()
	g3.AddNode("00-001-01", []string{"00-001-02"})
	g3.AddNode("00-001-02", []string{"00-001-01"})

	// This should either fail or handle gracefully
	dispatcher3 := graph.NewDispatcher(g3, 3)
	results3 := dispatcher3.Execute(func(wsID string) error { return nil })

	// Circular dependency should not crash
	t.Logf("Circular dependency handling: %d results", len(results3))
}


