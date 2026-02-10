package executor

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ProgressEvent represents a progress update event for machine-readable output
type ProgressEvent struct {
	WSID      string `json:"ws_id"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// EvidenceEvent represents an evidence chain event
type EvidenceEvent struct {
	Type      string `json:"type"`      // plan, generation, verification, approval
	WSID      string `json:"ws_id"`
	Timestamp string `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ProgressRenderer handles formatting of progress output
type ProgressRenderer struct {
	outputFormat string // "human" or "json"
}

// NewProgressRenderer creates a new progress renderer
func NewProgressRenderer(outputFormat string) *ProgressRenderer {
	return &ProgressRenderer{
		outputFormat: outputFormat,
	}
}

// RenderProgressBar renders a human-readable progress bar
// Format: [WS-ID] ████░░░░ 50% — message
func (r *ProgressRenderer) RenderProgressBar(wsID string, progress int, message string) string {
	// Clamp progress to 0-100
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	// Calculate bar segments (12 segments total)
	filled := (progress * 12) / 100
	empty := 12 - filled

	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	return fmt.Sprintf("[%s] %s %d%% — %s", wsID, bar, progress, message)
}

// RenderJSONEvent renders a progress event as JSON
func (r *ProgressRenderer) RenderJSONEvent(event ProgressEvent) string {
	// Set timestamp if not provided
	if event.Timestamp == "" {
		event.Timestamp = time.Now().Format(time.RFC3339)
	}

	data, err := json.Marshal(event)
	if err != nil {
		// Fallback to simple JSON
		return fmt.Sprintf(`{"ws_id":"%s","status":"error","progress":0,"message":"json marshal error: %v"}`,
			event.WSID, err)
	}

	return string(data)
}

// RenderEvidenceEvent renders an evidence event
func (r *ProgressRenderer) RenderEvidenceEvent(event EvidenceEvent) string {
	if event.Timestamp == "" {
		event.Timestamp = time.Now().Format(time.RFC3339)
	}

	data, err := json.Marshal(event)
	if err != nil {
		return ""
	}

	return string(data)
}

// Output returns the formatted output based on format type
func (r *ProgressRenderer) Output(wsID string, progress int, status, message string) string {
	if r.outputFormat == "json" {
		return r.RenderJSONEvent(ProgressEvent{
			WSID:      wsID,
			Status:    status,
			Progress:  progress,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}

	// Human-readable format
	return r.RenderProgressBar(wsID, progress, message)
}

// RenderError renders an error message
func (r *ProgressRenderer) RenderError(wsID string, err error) string {
	if r.outputFormat == "json" {
		event := ProgressEvent{
			WSID:      wsID,
			Status:    "error",
			Progress:  0,
			Message:   err.Error(),
			Timestamp: time.Now().Format(time.RFC3339),
		}
		return r.RenderJSONEvent(event)
	}

	return fmt.Sprintf("[%s] ERROR: %v", wsID, err)
}

// RenderSuccess renders a success message
func (r *ProgressRenderer) RenderSuccess(wsID string, message string) string {
	if r.outputFormat == "json" {
		event := ProgressEvent{
			WSID:      wsID,
			Status:    "success",
			Progress:  100,
			Message:   message,
			Timestamp: time.Now().Format(time.RFC3339),
		}
		return r.RenderJSONEvent(event)
	}

	return fmt.Sprintf("[%s] ✓ %s", wsID, message)
}

// RenderSummary renders a summary of execution results
func (r *ProgressRenderer) RenderSummary(summary ExecutionSummary) string {
	if r.outputFormat == "json" {
		data, _ := json.Marshal(summary)
		return string(data)
	}

	var sb strings.Builder
	sb.WriteString("\n=== Execution Summary ===\n")
	sb.WriteString(fmt.Sprintf("Total: %d\n", summary.TotalWorkstreams))
	sb.WriteString(fmt.Sprintf("Executed: %d\n", summary.Executed))
	sb.WriteString(fmt.Sprintf("Succeeded: %d\n", summary.Succeeded))
	sb.WriteString(fmt.Sprintf("Failed: %d\n", summary.Failed))
	sb.WriteString(fmt.Sprintf("Skipped: %d\n", summary.Skipped))
	sb.WriteString(fmt.Sprintf("Retries: %d\n", summary.Retries))
	if summary.Duration > 0 {
		sb.WriteString(fmt.Sprintf("Duration: %v\n", summary.Duration))
	}

	return sb.String()
}
