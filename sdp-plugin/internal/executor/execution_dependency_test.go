package executor

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestExecutor_DependencyOrder tests AC7: respects dependency order
func TestExecutor_DependencyOrder(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	}, newTestRunner())

	ctx := context.Background()
	var output bytes.Buffer

	_, err := exec.Execute(ctx, &output, ExecuteOptions{
		All:        true,
		SpecificWS: "",
		Retry:      0,
		Output:     "human",
	})

	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	// Verify execution order respects dependencies
	// 00-054-01 depends on nothing
	// 00-054-02 depends on 00-054-01
	// So 00-054-01 should execute before 00-054-02
	outputStr := output.String()
	pos01 := strings.Index(outputStr, "00-054-01")
	pos02 := strings.Index(outputStr, "00-054-02")

	if pos01 == -1 || pos02 == -1 {
		t.Fatal("Expected both workstreams in output")
	}

	if pos01 > pos02 {
		t.Error("00-054-01 should execute before 00-054-02 (dependency order violation)")
	}
}

// TestExecutor_ParseDependenciesError_SafeFallback tests that when ParseDependencies
// fails (e.g. workstream file not found), we use safe fallback (empty deps) and continue
// execution rather than skipping or failing the whole run.
func TestExecutor_ParseDependenciesError_SafeFallback(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	}, newTestRunner())

	ctx := context.Background()
	var output bytes.Buffer

	// 00-054-99 has no file in testdata/backlog, so ParseDependencies will fail
	_, err := exec.Execute(ctx, &output, ExecuteOptions{
		All:        false,
		SpecificWS: "00-054-99",
		Retry:      0,
		Output:     "human",
	})

	if err != nil {
		t.Fatalf("Execute() should not fail on ParseDependencies error (safe fallback): %v", err)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Warning: failed to parse dependencies") {
		t.Error("Expected warning about failed dependency parse in output")
	}
	// Workstream should still be executed (with empty deps)
	if !strings.Contains(outputStr, "00-054-99") {
		t.Error("Expected workstream 00-054-99 to be executed despite parse failure")
	}
}
