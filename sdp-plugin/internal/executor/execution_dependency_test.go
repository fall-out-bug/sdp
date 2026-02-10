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
	})

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
