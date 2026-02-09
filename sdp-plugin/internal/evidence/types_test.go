package evidence

import (
	"encoding/json"
	"testing"
)

func TestEventTypes_Count(t *testing.T) {
	want := 6
	if got := len(EventTypes); got != want {
		t.Errorf("EventTypes: want %d, got %d", want, got)
	}
}

func TestEventTypes_Values(t *testing.T) {
	expected := map[string]bool{"plan": true, "generation": true, "verification": true, "approval": true, "decision": true, "lesson": true}
	for _, et := range EventTypes {
		if !expected[et] {
			t.Errorf("unexpected event type: %s", et)
		}
	}
}

func TestEvent_JSONRoundTrip(t *testing.T) {
	ev := Event{
		ID:        "evt-1",
		Type:      "decision",
		Timestamp: "2026-02-09T12:00:00Z",
		WSID:      "00-054-03",
		CommitSHA: "abc123",
		PrevHash:  "sha256-prev",
		Data: DecisionEventData{
			Question:  "How to store evidence?",
			Choice:    "JSONL",
			Rationale: "Simple",
			Tags:      []string{"storage"},
		},
	}
	b, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Event
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ID != ev.ID || out.Type != ev.Type || out.WSID != ev.WSID {
		t.Errorf("roundtrip: got %+v", out)
	}
}

func TestGenerationData_Fields(t *testing.T) {
	d := GenerationData{
		ModelID:      "claude-sonnet",
		ModelVersion: "20250514",
		PromptHash:   "sha256-abc",
		FilesChanged: []string{"internal/evidence/types.go"},
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out GenerationData
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.ModelID != d.ModelID || len(out.FilesChanged) != 1 {
		t.Errorf("roundtrip: got %+v", out)
	}
}

func TestLessonEventData_Outcome(t *testing.T) {
	d := LessonEventData{
		Category:   "architecture",
		Insight:    "JSONL scales",
		SourceWSID: "00-054-04",
		Outcome:    "worked",
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out LessonEventData
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Outcome != "worked" || out.SourceWSID != d.SourceWSID {
		t.Errorf("roundtrip: got %+v", out)
	}
}
