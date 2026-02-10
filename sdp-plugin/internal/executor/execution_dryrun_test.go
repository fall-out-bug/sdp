package executor

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestExecutor_DryRun tests AC4: dry-run shows execution plan without running
func TestExecutor_DryRun(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     true,
		RetryCount: 1,
	})

	ctx := context.Background()
	var output bytes.Buffer

	result, err := exec.Execute(ctx, &output, ExecuteOptions{
		All:        true,
		SpecificWS: "",
		Retry:      0,
		Output:     "human",
	})

	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if result.Executed > 0 {
		t.Errorf("Dry-run should not execute any workstreams, got %d executed", result.Executed)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "DRY RUN") {
		t.Error("Expected output to contain 'DRY RUN'")
	}

	if !strings.Contains(outputStr, "Would execute") {
		t.Error("Expected output to contain 'Would execute'")
	}
}
