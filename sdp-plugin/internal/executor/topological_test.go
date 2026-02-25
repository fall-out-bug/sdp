package executor

import (
	"strings"
	"testing"
)

// TestExecutor_ParseWorkstreamDependencies tests parsing dependencies from workstream files
func TestExecutor_ParseWorkstreamDependencies(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
	}, newTestRunner())

	wsID := "00-054-02"
	deps, err := exec.ParseDependencies(wsID)

	if err != nil {
		t.Fatalf("ParseDependencies() failed: %v", err)
	}

	if len(deps) == 0 {
		t.Error("Expected 00-054-02 to have dependencies")
	}

	// 00-054-02 should depend on 00-054-01
	hasDep01 := false
	for _, dep := range deps {
		if dep == "00-054-01" {
			hasDep01 = true
			break
		}
	}

	if !hasDep01 {
		t.Error("Expected 00-054-02 to depend on 00-054-01")
	}
}

// TestExecutor_TopologicalSort tests topological sorting of workstreams
func TestExecutor_TopologicalSort(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
	}, newTestRunner())

	// Create a test graph with dependencies
	workstreams := []string{"00-054-01", "00-054-02", "00-054-03"}
	dependencies := map[string][]string{
		"00-054-01": {},
		"00-054-02": {"00-054-01"},
		"00-054-03": {"00-054-02"},
	}

	sorted, err := exec.TopologicalSort(workstreams, dependencies)

	if err != nil {
		t.Fatalf("TopologicalSort() failed: %v", err)
	}

	// Verify order: 00-054-01 should come before 00-054-02, which should come before 00-054-03
	pos01 := indexOf(sorted, "00-054-01")
	pos02 := indexOf(sorted, "00-054-02")
	pos03 := indexOf(sorted, "00-054-03")

	if pos01 >= pos02 {
		t.Error("00-054-01 should come before 00-054-02")
	}

	if pos02 >= pos03 {
		t.Error("00-054-02 should come before 00-054-03")
	}
}

// TestExecutor_CyclicDependencies tests detection of cyclic dependencies
func TestExecutor_CyclicDependencies(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
	}, newTestRunner())

	// Create a cyclic dependency
	workstreams := []string{"00-054-01", "00-054-02"}
	dependencies := map[string][]string{
		"00-054-01": {"00-054-02"},
		"00-054-02": {"00-054-01"},
	}

	_, err := exec.TopologicalSort(workstreams, dependencies)

	if err == nil {
		t.Error("Expected error for cyclic dependencies")
	}

	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("Expected cycle error, got: %v", err)
	}
}

// Helper function
func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}
