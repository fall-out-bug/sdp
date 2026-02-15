package evidence

import (
	"bytes"
	"testing"
)

func TestVerifyChain(t *testing.T) {
	tempDir := t.TempDir()
	logPath := tempDir + "/events.jsonl"

	log, _ := NewEvidenceLog(logPath)

	// Empty log should verify
	if err := log.VerifyChain(); err != nil {
		t.Errorf("empty log should verify: %v", err)
	}

	// Add events with proper chain
	log.AppendEvent(EventTypePlan, nil)
	log.AppendEvent(EventTypeDecision, nil)
	log.AppendEvent(EventTypeGeneration, nil)

	if err := log.VerifyChain(); err != nil {
		t.Errorf("valid chain should verify: %v", err)
	}
}

func TestVerifyHash(t *testing.T) {
	event := NewBaseEvent(EventTypeDecision, "")

	if !VerifyHash(event) {
		t.Error("valid event hash should verify")
	}

	// Corrupt the hash
	event.Hash = "corrupted"
	if VerifyHash(event) {
		t.Error("corrupted hash should not verify")
	}
}

func TestVerifyEvents(t *testing.T) {
	// Create valid chain
	event1 := NewBaseEvent(EventTypePlan, "")
	event2 := NewBaseEvent(EventTypeDecision, event1.Hash)

	events := []BaseEvent{event1, event2}
	if err := VerifyEvents(events); err != nil {
		t.Errorf("valid events should verify: %v", err)
	}

	// Broken chain
	event3 := NewBaseEvent(EventTypeGeneration, "wrong-hash")
	events = []BaseEvent{event1, event3}
	if err := VerifyEvents(events); err == nil {
		t.Error("broken chain should fail verification")
	}
}

func TestFindEventByID(t *testing.T) {
	event1 := NewBaseEvent(EventTypePlan, "")
	event2 := NewBaseEvent(EventTypeDecision, event1.Hash)
	event3 := NewBaseEvent(EventTypeGeneration, event2.Hash)

	events := []BaseEvent{event1, event2, event3}

	found, idx := FindEventByID(events, event2.ID)
	if idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
	if found.ID != event2.ID {
		t.Error("wrong event returned")
	}

	_, idx = FindEventByID(events, "nonexistent")
	if idx != -1 {
		t.Error("nonexistent event should return -1")
	}
}

func TestFilterByType(t *testing.T) {
	event1 := NewBaseEvent(EventTypePlan, "")
	event2 := NewBaseEvent(EventTypeDecision, event1.Hash)
	event3 := NewBaseEvent(EventTypeDecision, event2.Hash)
	event4 := NewBaseEvent(EventTypeGeneration, event3.Hash)

	events := []BaseEvent{event1, event2, event3, event4}

	decisions := FilterByType(events, EventTypeDecision)
	if len(decisions) != 2 {
		t.Errorf("expected 2 decisions, got %d", len(decisions))
	}

	plans := FilterByType(events, EventTypePlan)
	if len(plans) != 1 {
		t.Errorf("expected 1 plan, got %d", len(plans))
	}

	approvals := FilterByType(events, EventTypeApproval)
	if len(approvals) != 0 {
		t.Errorf("expected 0 approvals, got %d", len(approvals))
	}
}

func TestEventsToJSONL(t *testing.T) {
	event1 := NewBaseEvent(EventTypePlan, "")
	event2 := NewBaseEvent(EventTypeDecision, event1.Hash)

	events := []BaseEvent{event1, event2}

	jsonl, err := EventsToJSONL(events)
	if err != nil {
		t.Fatalf("EventsToJSONL failed: %v", err)
	}

	// Should have two lines
	lines := bytes.Count([]byte(jsonl), []byte("\n"))
	if lines != 2 {
		t.Errorf("expected 2 lines, got %d", lines)
	}
}

func TestCopyTo(t *testing.T) {
	tempDir := t.TempDir()
	logPath := tempDir + "/events.jsonl"

	log, _ := NewEvidenceLog(logPath)
	log.AppendEvent(EventTypePlan, nil)
	log.AppendEvent(EventTypeDecision, nil)

	var buf bytes.Buffer
	if err := log.CopyTo(&buf); err != nil {
		t.Fatalf("CopyTo failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("CopyTo should write data")
	}
}

func TestCopyTo_EmptyLog(t *testing.T) {
	tempDir := t.TempDir()
	logPath := tempDir + "/events.jsonl"

	log, _ := NewEvidenceLog(logPath)

	var buf bytes.Buffer
	if err := log.CopyTo(&buf); err != nil {
		t.Fatalf("CopyTo on empty log failed: %v", err)
	}

	if buf.Len() != 0 {
		t.Error("empty log should write nothing")
	}
}
