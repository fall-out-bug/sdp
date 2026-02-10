package executor

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestExecutor_ContextCancellation tests context cancellation
func TestExecutor_ContextCancellation(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var output bytes.Buffer

	_, err := exec.Execute(ctx, &output, ExecuteOptions{
		All:        true,
		SpecificWS: "",
		Retry:      0,
		Output:     "human",
	})

	if err == nil {
		t.Error("Expected error when context is cancelled")
	}

	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context canceled error, got: %v", err)
	}
}
