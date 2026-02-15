package evidence

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestEventTypeValues(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		expected  string
	}{
		{"decision", EventTypeDecision, "decision"},
		{"plan", EventTypePlan, "plan"},
		{"generation", EventTypeGeneration, "generation"},
		{"verification", EventTypeVerification, "verification"},
		{"acceptance", EventTypeAcceptance, "acceptance"},
		{"approval", EventTypeApproval, "approval"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.eventType) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.eventType)
			}
		})
	}
}

func TestValidEventTypes(t *testing.T) {
	validTypes := ValidEventTypes()

	expected := []EventType{
		EventTypeDecision,
		EventTypePlan,
		EventTypeGeneration,
		EventTypeVerification,
		EventTypeAcceptance,
		EventTypeApproval,
	}

	if len(validTypes) != len(expected) {
		t.Errorf("expected %d valid types, got %d", len(expected), len(validTypes))
	}

	for _, exp := range expected {
		found := false
		for _, v := range validTypes {
			if v == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing valid type: %s", exp)
		}
	}
}

func TestEventType_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		valid     bool
	}{
		{"decision is valid", EventTypeDecision, true},
		{"plan is valid", EventTypePlan, true},
		{"invalid type", EventType("invalid"), false},
		{"empty type", EventType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.eventType.IsValid() != tt.valid {
				t.Errorf("expected IsValid=%v for %s", tt.valid, tt.eventType)
			}
		})
	}
}

func TestBaseEvent_New(t *testing.T) {
	event := NewBaseEvent(EventTypeDecision, "prev-hash-123")

	if event.ID == "" {
		t.Error("ID should not be empty")
	}
	if event.Type != EventTypeDecision {
		t.Errorf("expected type decision, got %s", event.Type)
	}
	if event.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	if event.PrevHash != "prev-hash-123" {
		t.Errorf("expected prev hash prev-hash-123, got %s", event.PrevHash)
	}
	if event.Hash == "" {
		t.Error("hash should be calculated")
	}
}

func TestBaseEvent_NewGenesis(t *testing.T) {
	// Genesis event has no previous hash
	event := NewBaseEvent(EventTypePlan, "")

	if event.PrevHash != "" {
		t.Error("genesis event should have empty prev hash")
	}
	if event.Hash == "" {
		t.Error("hash should still be calculated")
	}
}

func TestBaseEvent_CalculateHash(t *testing.T) {
	event1 := NewBaseEvent(EventTypeDecision, "prev-1")
	event2 := NewBaseEvent(EventTypeDecision, "prev-2")

	// Different prev hashes should produce different hashes
	if event1.Hash == event2.Hash {
		t.Error("events with different prev hashes should have different hashes")
	}

	// Same event should have consistent hash
	hash1 := event1.Hash
	event1.CalculateHash()
	if event1.Hash != hash1 {
		t.Error("hash should be consistent for same content")
	}
}

func TestBaseEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   BaseEvent
		wantErr bool
	}{
		{
			name:    "valid event",
			event:   NewBaseEvent(EventTypeDecision, "prev"),
			wantErr: false,
		},
		{
			name: "missing ID",
			event: BaseEvent{
				Type:      EventTypeDecision,
				Timestamp: time.Now(),
				Hash:      "hash",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			event: BaseEvent{
				ID:        "id",
				Timestamp: time.Now(),
				Hash:      "hash",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			event: BaseEvent{
				ID:        "id",
				Type:      EventType("invalid"),
				Timestamp: time.Now(),
				Hash:      "hash",
			},
			wantErr: true,
		},
		{
			name: "missing hash",
			event: BaseEvent{
				ID:        "id",
				Type:      EventTypeDecision,
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero timestamp",
			event: BaseEvent{
				ID:   "id",
				Type: EventTypeDecision,
				Hash: "hash",
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

func TestBaseEvent_JSON(t *testing.T) {
	original := NewBaseEvent(EventTypeGeneration, "prev-hash")

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed BaseEvent
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed.ID != original.ID {
		t.Errorf("ID mismatch: %s != %s", parsed.ID, original.ID)
	}
	if parsed.Type != original.Type {
		t.Errorf("Type mismatch: %s != %s", parsed.Type, original.Type)
	}
	if parsed.Hash != original.Hash {
		t.Errorf("Hash mismatch: %s != %s", parsed.Hash, original.Hash)
	}
}

func TestBaseEvent_JSONL(t *testing.T) {
	event := NewBaseEvent(EventTypeVerification, "prev")

	// JSONL is just JSON with newline
	line, err := event.ToJSONL()
	if err != nil {
		t.Fatalf("ToJSONL failed: %v", err)
	}

	if !strings.HasSuffix(line, "\n") {
		t.Error("JSONL should end with newline")
	}

	parsed, err := ParseJSONL(line)
	if err != nil {
		t.Fatalf("ParseJSONL failed: %v", err)
	}

	if parsed.ID != event.ID {
		t.Errorf("ID mismatch after JSONL roundtrip")
	}
}

func TestParseJSONL_Invalid(t *testing.T) {
	_, err := ParseJSONL("not valid json")
	if err == nil {
		t.Error("should fail on invalid JSON")
	}
}

func TestBaseEvent_HashConsistency(t *testing.T) {
	// Create two events with same content
	event1 := BaseEvent{
		ID:        "test-id",
		Type:      EventTypeDecision,
		Timestamp: time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC),
		PrevHash:  "prev-hash",
	}
	event1.CalculateHash()

	event2 := BaseEvent{
		ID:        "test-id",
		Type:      EventTypeDecision,
		Timestamp: time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC),
		PrevHash:  "prev-hash",
	}
	event2.CalculateHash()

	if event1.Hash != event2.Hash {
		t.Error("same content should produce same hash")
	}
}
