package executor

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

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
