package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestExecutor_ExecuteAllReady tests AC1: execute all ready workstreams (no blockers)
func TestExecutor_ExecuteAllReady(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
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
	})

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
	// Mock executor that fails on first attempt
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
		DryRun:     false,
		RetryCount: 3,
	})

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

// TestExecutor_ProgressBar tests AC5: streaming progress with progress bar
func TestExecutor_ProgressBar(t *testing.T) {
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

// TestExecutor_JSONOutput tests AC6: JSON output format for machine-readable events
func TestExecutor_JSONOutput(t *testing.T) {
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
		Output:     "json",
	})

	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	outputStr := output.String()
	t.Logf("Full output:\n%s\n", outputStr)

	lines := strings.Split(strings.TrimSpace(outputStr), "\n")

	// Parse each JSON line
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try to parse as ProgressEvent first
		var event ProgressEvent
		err := json.Unmarshal([]byte(line), &event)

		// Check if ws_id is present (if not, this is likely the summary)
		if event.WSID == "" {
			// Try to parse as ExecutionSummary
			var summary ExecutionSummary
			if sumErr := json.Unmarshal([]byte(line), &summary); sumErr != nil {
				t.Logf("Line %d: failed to parse as ProgressEvent or ExecutionSummary: %v\nLine: %s", i, err, line)
			}
			// Skip validation for summary lines
			continue
		}

		// Validate required fields for ProgressEvent
		if event.Status == "" {
			t.Errorf("Line %d: missing status field", i)
		}
		if event.Progress < 0 || event.Progress > 100 {
			t.Errorf("Line %d: invalid progress value: %d", i, event.Progress)
		}
	}
}

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

// TestExecutor_EvidenceChain tests AC8: emits full evidence chain
func TestExecutor_EvidenceChain(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir:      "testdata/backlog",
		DryRun:          false,
		RetryCount:      1,
		EvidenceLogPath: "testdata/evidence.jsonl",
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

	if len(result.EvidenceEvents) == 0 {
		t.Error("Expected evidence events to be emitted")
	}

	// Check for required evidence types
	hasPlanEvent := false
	hasGenEvent := false
	hasVerifyEvent := false
	hasApprovalEvent := false

	for _, event := range result.EvidenceEvents {
		switch event.Type {
		case "plan":
			hasPlanEvent = true
		case "generation":
			hasGenEvent = true
		case "verification":
			hasVerifyEvent = true
		case "approval":
			hasApprovalEvent = true
		}
	}

	if !hasPlanEvent {
		t.Error("Expected 'plan' evidence event")
	}
	if !hasGenEvent {
		t.Error("Expected 'generation' evidence event")
	}
	if !hasVerifyEvent {
		t.Error("Expected 'verification' evidence event")
	}
	if !hasApprovalEvent {
		t.Error("Expected 'approval' evidence event")
	}
}

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

// TestProgressRenderer_RenderProgressBar tests progress bar rendering
func TestProgressRenderer_RenderProgressBar(t *testing.T) {
	renderer := NewProgressRenderer("human")

	tests := []struct {
		name     string
		progress int
		message  string
	}{
		{
			name:     "0%",
			progress: 0,
			message:  "starting",
		},
		{
			name:     "50%",
			progress: 50,
			message:  "running tests",
		},
		{
			name:     "100%",
			progress: 100,
			message:  "complete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderer.RenderProgressBar("00-054-01", tt.progress, tt.message)

			// Verify the format is correct
			if !strings.Contains(output, "[00-054-01]") {
				t.Errorf("Expected workstream ID in output: %s", output)
			}
			if !strings.Contains(output, fmt.Sprintf("%d%%", tt.progress)) {
				t.Errorf("Expected progress percentage in output: %s", output)
			}
			if !strings.Contains(output, tt.message) {
				t.Errorf("Expected message in output: %s", output)
			}

			// Verify progress bar characters
			if !strings.Contains(output, "█") && tt.progress > 0 {
				t.Errorf("Expected filled bar character for progress > 0: %s", output)
			}
			if !strings.Contains(output, "░") && tt.progress < 100 {
				t.Errorf("Expected empty bar character for progress < 100: %s", output)
			}
		})
	}
}

// TestProgressRenderer_RenderJSONEvent tests JSON event rendering
func TestProgressRenderer_RenderJSONEvent(t *testing.T) {
	renderer := NewProgressRenderer("json")

	event := ProgressEvent{
		WSID:      "00-054-01",
		Status:    "running",
		Progress:  50,
		Message:   "running tests",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	output := renderer.RenderJSONEvent(event)

	var parsed ProgressEvent
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if parsed.WSID != event.WSID {
		t.Errorf("ws_id mismatch: got %s, want %s", parsed.WSID, event.WSID)
	}

	if parsed.Status != event.Status {
		t.Errorf("status mismatch: got %s, want %s", parsed.Status, event.Status)
	}

	if parsed.Progress != event.Progress {
		t.Errorf("progress mismatch: got %d, want %d", parsed.Progress, event.Progress)
	}
}

// TestExecutor_ParseWorkstreamDependencies tests parsing dependencies from workstream files
func TestExecutor_ParseWorkstreamDependencies(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
	})

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
	})

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
	})

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

// TestExecutor_SetEvidenceWriter tests setting evidence writer
func TestExecutor_SetEvidenceWriter(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
	})

	var buf bytes.Buffer
	exec.SetEvidenceWriter(&buf)

	if exec.evidenceWriter == nil {
		t.Error("Expected evidence writer to be set")
	}
}

// TestExecutor_EmitEvidenceEvent tests evidence event emission
func TestExecutor_EmitEvidenceEvent(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
	})

	var buf bytes.Buffer
	exec.SetEvidenceWriter(&buf)

	event := EvidenceEvent{
		Type:      "test",
		WSID:      "00-054-01",
		Timestamp: "2026-02-10T10:00:00Z",
		Data:      map[string]interface{}{"key": "value"},
	}

	err := exec.emitEvidenceEvent(event)
	if err != nil {
		t.Errorf("emitEvidenceEvent() failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("Expected output in evidence writer")
	}

	if !strings.Contains(output, "test") {
		t.Errorf("Expected event type in output: %s", output)
	}
}

// TestProgressRenderer_RenderEvidenceEvent tests evidence event rendering
func TestProgressRenderer_RenderEvidenceEvent(t *testing.T) {
	renderer := NewProgressRenderer("json")

	event := EvidenceEvent{
		Type:      "plan",
		WSID:      "00-054-01",
		Timestamp: "2026-02-10T10:00:00Z",
		Data:      map[string]interface{}{"action": "test"},
	}

	output := renderer.RenderEvidenceEvent(event)

	if output == "" {
		t.Error("Expected non-empty output")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}

	if parsed["type"] != "plan" {
		t.Errorf("Expected type 'plan', got %v", parsed["type"])
	}
}

// TestProgressRenderer_RenderError tests error rendering
func TestProgressRenderer_RenderError(t *testing.T) {
	renderer := NewProgressRenderer("human")

	err := fmt.Errorf("test error")
	output := renderer.RenderError("00-054-01", err)

	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected 'ERROR' in output: %s", output)
	}

	if !strings.Contains(output, "00-054-01") {
		t.Errorf("Expected workstream ID in output: %s", output)
	}

	if !strings.Contains(output, "test error") {
		t.Errorf("Expected error message in output: %s", output)
	}
}

// TestProgressRenderer_RenderErrorJSON tests error rendering in JSON format
func TestProgressRenderer_RenderErrorJSON(t *testing.T) {
	renderer := NewProgressRenderer("json")

	err := fmt.Errorf("test error")
	output := renderer.RenderError("00-054-01", err)

	var event ProgressEvent
	if parseErr := json.Unmarshal([]byte(output), &event); parseErr != nil {
		t.Errorf("Failed to parse JSON: %v", parseErr)
	}

	if event.Status != "error" {
		t.Errorf("Expected status 'error', got %s", event.Status)
	}

	if event.WSID != "00-054-01" {
		t.Errorf("Expected ws_id '00-054-01', got %s", event.WSID)
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
