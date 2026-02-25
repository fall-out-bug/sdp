package executor

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// TestExecutor_ContextCancellation tests context cancellation before execution
func TestExecutor_ContextCancellation(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	}, newTestRunner())

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

// TestExecutor_ContextCancellationMidExecution tests that long-running workstreams
// respect ctx.Done() when cancelled mid-execution (runner checks ctx periodically).
func TestExecutor_ContextCancellationMidExecution(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	}, newBlockingRunner())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	var execErr error
	go func() {
		var output bytes.Buffer
		_, execErr = exec.Execute(ctx, &output, ExecuteOptions{
			All:        false,
			SpecificWS: "00-054-01",
			Retry:      0,
			Output:     "human",
		})
		close(done)
	}()

	// Cancel after runner has started
	time.Sleep(50 * time.Millisecond)
	cancel()

	<-done
	if execErr == nil {
		t.Error("Expected error when context is cancelled mid-execution")
	}
	if !strings.Contains(execErr.Error(), "context canceled") {
		t.Errorf("Expected context canceled error, got: %v", execErr)
	}
}
