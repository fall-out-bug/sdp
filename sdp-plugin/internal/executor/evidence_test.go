package executor

import (
	"bytes"
	"strings"
	"testing"
)

// TestExecutor_SetEvidenceWriter tests setting evidence writer
func TestExecutor_SetEvidenceWriter(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir: "testdata/backlog",
	}, newTestRunner())

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
	}, newTestRunner())

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
