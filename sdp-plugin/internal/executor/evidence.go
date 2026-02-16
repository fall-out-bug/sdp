package executor

import (
	"fmt"
	"time"
)

// generateEvidenceEvents generates evidence chain events for a workstream
func (e *Executor) generateEvidenceEvents(wsID string) []EvidenceEvent {
	now := time.Now().Format(time.RFC3339)

	return []EvidenceEvent{
		{
			Type:      "plan",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "execution_plan"},
		},
		{
			Type:      "generation",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "code_generation"},
		},
		{
			Type:      "verification",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "test_verification"},
		},
		{
			Type:      "approval",
			WSID:      wsID,
			Timestamp: now,
			Data:      map[string]interface{}{"action": "auto_approval"},
		},
	}
}

// emitEvidenceEvent writes an evidence event to the evidence log
func (e *Executor) emitEvidenceEvent(event EvidenceEvent) error {
	if e.evidenceWriter == nil {
		return nil // No evidence writer configured
	}

	output := e.progress.RenderEvidenceEvent(event)
	if output == "" {
		return fmt.Errorf("failed to render evidence event")
	}

	_, err := fmt.Fprintln(e.evidenceWriter, output)
	return err
}
