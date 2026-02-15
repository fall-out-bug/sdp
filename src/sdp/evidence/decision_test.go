package evidence

import (
	"encoding/json"
	"testing"
)

func TestDecisionEvent_New(t *testing.T) {
	prevHash := "prev-123"
	event := NewDecisionEvent("What framework?", "React", prevHash)

	if event.Question != "What framework?" {
		t.Errorf("expected question, got %s", event.Question)
	}
	if event.Answer != "React" {
		t.Errorf("expected answer, got %s", event.Answer)
	}
	if event.BaseEvent.Type != EventTypeDecision {
		t.Errorf("expected type decision, got %s", event.BaseEvent.Type)
	}
	if event.BaseEvent.PrevHash != prevHash {
		t.Errorf("expected prev hash, got %s", event.BaseEvent.PrevHash)
	}
}

func TestDecisionEvent_WithRationale(t *testing.T) {
	event := NewDecisionEvent("Auth method?", "JWT", "")
	event.Rationale = "Stateless and scalable"

	if event.Rationale != "Stateless and scalable" {
		t.Errorf("expected rationale, got %s", event.Rationale)
	}
}

func TestDecisionEvent_WithContext(t *testing.T) {
	event := NewDecisionEvent("Database?", "PostgreSQL", "")
	event.Context = map[string]interface{}{
		"budget": 1000,
		"scale":  "medium",
	}

	if event.Context["budget"] != 1000 {
		t.Error("context budget not set")
	}
}

func TestDecisionEvent_WithSession(t *testing.T) {
	event := NewDecisionEvent("Cache?", "Redis", "")
	event.SessionID = "sess-456"
	event.Workstream = "00-054-03"

	if event.SessionID != "sess-456" {
		t.Error("session ID not set")
	}
	if event.Workstream != "00-054-03" {
		t.Error("workstream not set")
	}
}

func TestDecisionEvent_JSON(t *testing.T) {
	original := NewDecisionEvent("Language?", "Go", "prev-hash")
	original.Rationale = "Performance and simplicity"
	original.SessionID = "sess-123"

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed DecisionEvent
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed.Question != original.Question {
		t.Error("question mismatch after roundtrip")
	}
	if parsed.Answer != original.Answer {
		t.Error("answer mismatch after roundtrip")
	}
	if parsed.Rationale != original.Rationale {
		t.Error("rationale mismatch after roundtrip")
	}
}

func TestDecisionEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   DecisionEvent
		wantErr bool
	}{
		{
			name:    "valid event",
			event:   NewDecisionEvent("Question?", "Answer", ""),
			wantErr: false,
		},
		{
			name: "missing question",
			event: DecisionEvent{
				BaseEvent: NewBaseEvent(EventTypeDecision, ""),
				Answer:    "Answer",
			},
			wantErr: true,
		},
		{
			name: "missing answer",
			event: DecisionEvent{
				BaseEvent: NewBaseEvent(EventTypeDecision, ""),
				Question:  "Question?",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecisionEvent_ToJSONL(t *testing.T) {
	event := NewDecisionEvent("Framework?", "Vue", "")
	event.Workstream = "00-054-03"

	line, err := event.ToJSONL()
	if err != nil {
		t.Fatalf("ToJSONL failed: %v", err)
	}

	parsed, err := ParseDecisionJSONL(line)
	if err != nil {
		t.Fatalf("ParseDecisionJSONL failed: %v", err)
	}

	if parsed.Question != event.Question {
		t.Error("question mismatch")
	}
}
