package evidence

import (
	"encoding/json"
	"errors"
)

// PlanEvent records a feature decomposition plan.
type PlanEvent struct {
	BaseEvent
	Feature      string   `json:"feature"`
	Description  string   `json:"description"`
	Workstreams  []string `json:"workstreams"`
	Dependencies []string `json:"dependencies,omitempty"`
	CostEstimate string   `json:"cost_estimate,omitempty"`
	SessionID    string   `json:"session_id,omitempty"`
}

// NewPlanEvent creates a new plan event.
func NewPlanEvent(feature, description string, workstreams []string, prevHash string) PlanEvent {
	return PlanEvent{
		BaseEvent:   NewBaseEvent(EventTypePlan, prevHash),
		Feature:     feature,
		Description: description,
		Workstreams: workstreams,
	}
}

// Validate checks required fields.
func (e PlanEvent) Validate() error {
	if err := e.BaseEvent.Validate(); err != nil {
		return err
	}
	if e.Feature == "" {
		return errors.New("feature is required")
	}
	if len(e.Workstreams) == 0 {
		return errors.New("at least one workstream is required")
	}
	return nil
}

// ToJSONL returns the event as JSONL.
func (e PlanEvent) ToJSONL() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ParsePlanJSONL parses JSONL to PlanEvent.
func ParsePlanJSONL(line string) (PlanEvent, error) {
	var e PlanEvent
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return PlanEvent{}, err
	}
	return e, nil
}
