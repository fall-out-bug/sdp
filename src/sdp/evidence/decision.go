package evidence

import (
	"encoding/json"
	"errors"
)

// DecisionEvent records a decision made during discovery.
type DecisionEvent struct {
	BaseEvent
	Question   string                 `json:"question"`
	Answer     string                 `json:"answer"`
	Rationale  string                 `json:"rationale,omitempty"`
	Context    map[string]interface{} `json:"context,omitempty"`
	SessionID  string                 `json:"session_id,omitempty"`
	Workstream string                 `json:"workstream,omitempty"`
}

// NewDecisionEvent creates a new decision event.
func NewDecisionEvent(question, answer, prevHash string) DecisionEvent {
	return DecisionEvent{
		BaseEvent: NewBaseEvent(EventTypeDecision, prevHash),
		Question:  question,
		Answer:    answer,
	}
}

// Validate checks required fields.
func (e DecisionEvent) Validate() error {
	if err := e.BaseEvent.Validate(); err != nil {
		return err
	}
	if e.Question == "" {
		return errors.New("question is required")
	}
	if e.Answer == "" {
		return errors.New("answer is required")
	}
	return nil
}

// ToJSONL returns the event as a JSONL line.
func (e DecisionEvent) ToJSONL() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ParseDecisionJSONL parses a JSONL line into a DecisionEvent.
func ParseDecisionJSONL(line string) (DecisionEvent, error) {
	var e DecisionEvent
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return DecisionEvent{}, err
	}
	return e, nil
}
