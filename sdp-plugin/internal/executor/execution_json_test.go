package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// TestExecutor_JSONOutput tests AC6: JSON output format for machine-readable events
func TestExecutor_JSONOutput(t *testing.T) {
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
