package evidence

import (
	"encoding/json"
	"testing"
)

func TestGenerationEvent_New(t *testing.T) {
	event := NewGenerationEvent("claude-3-opus", "hash123", "")

	if event.Model != "claude-3-opus" {
		t.Errorf("expected model, got %s", event.Model)
	}
	if event.PromptHash != "hash123" {
		t.Errorf("expected prompt hash, got %s", event.PromptHash)
	}
	if event.BaseEvent.Type != EventTypeGeneration {
		t.Errorf("expected type generation, got %s", event.BaseEvent.Type)
	}
}

func TestGenerationEvent_WithParameters(t *testing.T) {
	event := NewGenerationEvent("claude-3-opus", "hash", "")
	event.Parameters = map[string]interface{}{
		"temperature": 0.7,
		"max_tokens":  4096,
	}

	if event.Parameters["temperature"] != 0.7 {
		t.Error("temperature not set")
	}
}

func TestGenerationEvent_WithCodeHash(t *testing.T) {
	event := NewGenerationEvent("claude-3-opus", "hash", "")
	event.CodeHash = "code-hash-abc"

	if event.CodeHash != "code-hash-abc" {
		t.Error("code hash not set")
	}
}

func TestGenerationEvent_WithSpecRef(t *testing.T) {
	event := NewGenerationEvent("claude-3-opus", "hash", "")
	event.SpecRef = "docs/workstreams/backlog/00-054-05.md"

	if event.SpecRef == "" {
		t.Error("spec ref should be set")
	}
}

func TestGenerationEvent_JSON(t *testing.T) {
	original := NewGenerationEvent("claude-3-sonnet", "prompt-hash", "")
	original.ModelVersion = "20240229"
	original.CodeHash = "code-hash"
	original.Workstream = "00-054-05"

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed GenerationEvent
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed.Model != original.Model {
		t.Error("model mismatch")
	}
	if parsed.Workstream != original.Workstream {
		t.Error("workstream mismatch")
	}
}

func TestGenerationEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   GenerationEvent
		wantErr bool
	}{
		{
			name:    "valid",
			event:   NewGenerationEvent("claude-3-opus", "hash", ""),
			wantErr: false,
		},
		{
			name: "missing model",
			event: GenerationEvent{
				BaseEvent:  NewBaseEvent(EventTypeGeneration, ""),
				PromptHash: "hash",
			},
			wantErr: true,
		},
		{
			name: "missing prompt hash",
			event: GenerationEvent{
				BaseEvent: NewBaseEvent(EventTypeGeneration, ""),
				Model:     "claude-3-opus",
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
