package evidence

import (
	"encoding/json"
	"errors"
)

// GenerationEvent records AI code generation provenance.
type GenerationEvent struct {
	BaseEvent
	Model        string                 `json:"model"`
	ModelVersion string                 `json:"model_version,omitempty"`
	PromptHash   string                 `json:"prompt_hash"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	CodeHash     string                 `json:"code_hash,omitempty"`
	SpecRef      string                 `json:"spec_ref,omitempty"`
	Workstream   string                 `json:"workstream,omitempty"`
}

// NewGenerationEvent creates a new generation event.
func NewGenerationEvent(model, promptHash, prevHash string) GenerationEvent {
	return GenerationEvent{
		BaseEvent:  NewBaseEvent(EventTypeGeneration, prevHash),
		Model:      model,
		PromptHash: promptHash,
	}
}

// Validate checks required fields.
func (e GenerationEvent) Validate() error {
	if err := e.BaseEvent.Validate(); err != nil {
		return err
	}
	if e.Model == "" {
		return errors.New("model is required")
	}
	if e.PromptHash == "" {
		return errors.New("prompt_hash is required")
	}
	return nil
}

// ToJSONL returns the event as JSONL.
func (e GenerationEvent) ToJSONL() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ParseGenerationJSONL parses JSONL to GenerationEvent.
func ParseGenerationJSONL(line string) (GenerationEvent, error) {
	var e GenerationEvent
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return GenerationEvent{}, err
	}
	return e, nil
}
