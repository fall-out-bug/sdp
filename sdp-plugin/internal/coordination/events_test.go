package coordination

import (
	"testing"
	"time"
)

func TestAgentEvent_Types(t *testing.T) {
	// AC1: Define agent event types
	expectedTypes := []string{
		EventTypeAgentStart,
		EventTypeAgentAction,
		EventTypeAgentComplete,
		EventTypeAgentError,
		EventTypeAgentHandoff,
	}

	for _, et := range expectedTypes {
		if !isValidEventType(et) {
			t.Errorf("Invalid event type: %s", et)
		}
	}
}

func TestAgentEvent_Schema(t *testing.T) {
	// AC2: Event schema with agent_id, role, task_id, timestamp, payload
	event := &AgentEvent{
		ID:        "event-123",
		Type:      EventTypeAgentStart,
		AgentID:   "agent-001",
		Role:      "implementer",
		TaskID:    "task-456",
		Timestamp: time.Now(),
		Payload:   map[string]interface{}{"ws_id": "00-051-01"},
		PrevHash:  "prev123",
		Hash:      "hash123",
	}

	if event.ID != "event-123" {
		t.Errorf("Expected ID event-123, got %s", event.ID)
	}
	if event.AgentID != "agent-001" {
		t.Errorf("Expected AgentID agent-001, got %s", event.AgentID)
	}
	if event.Role != "implementer" {
		t.Errorf("Expected Role implementer, got %s", event.Role)
	}
	if event.TaskID != "task-456" {
		t.Errorf("Expected TaskID task-456, got %s", event.TaskID)
	}
	if event.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	if len(event.Payload) == 0 {
		t.Error("Payload should not be empty")
	}
}

func TestAgentEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   *AgentEvent
		wantErr bool
	}{
		{
			name: "valid event",
			event: &AgentEvent{
				ID:        "evt-1",
				Type:      EventTypeAgentStart,
				AgentID:   "agent-1",
				Role:      "implementer",
				TaskID:    "task-1",
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			event: &AgentEvent{
				Type:      EventTypeAgentStart,
				AgentID:   "agent-1",
				Role:      "implementer",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing agent_id",
			event: &AgentEvent{
				ID:        "evt-1",
				Type:      EventTypeAgentStart,
				Role:      "implementer",
				Timestamp: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			event: &AgentEvent{
				ID:        "evt-1",
				Type:      "invalid_type",
				AgentID:   "agent-1",
				Role:      "implementer",
				Timestamp: time.Now(),
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

func TestAgentEvent_ComputeHash(t *testing.T) {
	event := &AgentEvent{
		ID:        "evt-1",
		Type:      EventTypeAgentStart,
		AgentID:   "agent-1",
		Role:      "implementer",
		TaskID:    "task-1",
		Timestamp: time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC),
		Payload:   map[string]interface{}{"key": "value"},
		PrevHash:  "prev-hash",
	}

	hash1 := event.ComputeHash()
	hash2 := event.ComputeHash()

	if hash1 != hash2 {
		t.Error("Same event should produce same hash")
	}

	if len(hash1) != 64 {
		t.Errorf("Expected SHA256 hash (64 chars), got %d", len(hash1))
	}

	// Different content should produce different hash
	event2 := &AgentEvent{
		ID:        "evt-2",
		Type:      EventTypeAgentStart,
		AgentID:   "agent-1",
		Role:      "implementer",
		Timestamp: time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC),
		PrevHash:  "prev-hash",
	}

	hash3 := event2.ComputeHash()
	if hash1 == hash3 {
		t.Error("Different events should produce different hashes")
	}
}

func TestNewAgentEvent(t *testing.T) {
	event := NewAgentEvent(EventTypeAgentStart, "agent-123", "implementer")

	if event.ID == "" {
		t.Error("ID should be auto-generated")
	}
	if event.Type != EventTypeAgentStart {
		t.Errorf("Expected type %s, got %s", EventTypeAgentStart, event.Type)
	}
	if event.AgentID != "agent-123" {
		t.Errorf("Expected AgentID agent-123, got %s", event.AgentID)
	}
	if event.Role != "implementer" {
		t.Errorf("Expected role implementer, got %s", event.Role)
	}
	if event.Timestamp.IsZero() {
		t.Error("Timestamp should be set")
	}
	if event.Payload == nil {
		t.Error("Payload should be initialized")
	}
}
