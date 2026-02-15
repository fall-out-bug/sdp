package evidence

import (
	"encoding/json"
	"errors"
)

// VerificationEvent records test/tool verification results.
type VerificationEvent struct {
	BaseEvent
	Tool       string  `json:"tool"`
	Command    string  `json:"command"`
	Output     string  `json:"output"`
	Passed     bool    `json:"passed"`
	Coverage   float64 `json:"coverage,omitempty"`
	Duration   int64   `json:"duration_ms,omitempty"`
	Workstream string  `json:"workstream,omitempty"`
}

// NewVerificationEvent creates a new verification event.
func NewVerificationEvent(tool, command string, passed bool, prevHash string) VerificationEvent {
	return VerificationEvent{
		BaseEvent: NewBaseEvent(EventTypeVerification, prevHash),
		Tool:      tool,
		Command:   command,
		Passed:    passed,
	}
}

// Validate checks required fields.
func (e VerificationEvent) Validate() error {
	if err := e.BaseEvent.Validate(); err != nil {
		return err
	}
	if e.Tool == "" {
		return errors.New("tool is required")
	}
	if e.Command == "" {
		return errors.New("command is required")
	}
	return nil
}

// ToJSONL returns the event as JSONL.
func (e VerificationEvent) ToJSONL() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ParseVerificationJSONL parses JSONL to VerificationEvent.
func ParseVerificationJSONL(line string) (VerificationEvent, error) {
	var e VerificationEvent
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return VerificationEvent{}, err
	}
	return e, nil
}
