package executor

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestExecutor_ExecuteAllReady tests AC1: execute all ready workstreams (no blockers)
func TestExecutor_ExecuteAllReady(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	}, newTestRunner())

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

	if result.TotalWorkstreams == 0 {
		t.Error("Expected at least one workstream, got 0")
	}

	if result.Executed == 0 {
		t.Error("Expected at least one executed workstream, got 0")
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Executing") {
		t.Error("Expected output to contain 'Executing'")
	}
}

// TestExecutor_ExecuteSpecificWS tests AC2: execute specific workstream
func TestExecutor_ExecuteSpecificWS(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 1,
	}, newTestRunner())

	ctx := context.Background()
	var output bytes.Buffer

	result, err := exec.Execute(ctx, &output, ExecuteOptions{
		All:        false,
		SpecificWS: "00-054-01",
		Retry:      0,
		Output:     "human",
	})

	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if result.Executed != 1 {
		t.Errorf("Expected 1 executed workstream, got %d", result.Executed)
	}

	if result.TotalWorkstreams != 1 {
		t.Errorf("Expected 1 total workstream, got %d", result.TotalWorkstreams)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "00-054-01") {
		t.Error("Expected output to contain workstream ID '00-054-01'")
	}
}

// TestExecutor_RetryFailed tests AC3: retry failed workstream up to N times
func TestExecutor_RetryFailed(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 3,
	}, newTestRunner())

	ctx := context.Background()
	var output bytes.Buffer

	result, err := exec.Execute(ctx, &output, ExecuteOptions{
		All:        false,
		SpecificWS: "00-054-02", // This WS will fail in mock
		Retry:      3,
		Output:     "human",
	})

	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if result.Retries < 1 {
		t.Errorf("Expected at least 1 retry, got %d", result.Retries)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Retry") {
		t.Error("Expected output to contain 'Retry'")
	}
}

// TestExecutor_ProgressBar tests AC5: streaming progress with progress bar
func TestExecutor_ProgressBar(t *testing.T) {
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

	outputStr := output.String()
	// Check for progress bar characters (█ and ░)
	if !strings.Contains(outputStr, "█") || !strings.Contains(outputStr, "░") {
		t.Error("Expected progress bar with █ and ░ characters")
	}

	// Check for percentage
	if !strings.Contains(outputStr, "%") {
		t.Error("Expected progress percentage in output")
	}
}
