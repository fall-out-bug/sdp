package evidence

import (
	"encoding/json"
	"errors"
)

// AcceptanceEvent records smoke test / acceptance results.
type AcceptanceEvent struct {
	BaseEvent
	Command        string `json:"command"`
	Passed         bool   `json:"passed"`
	Output         string `json:"output"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	Workstream     string `json:"workstream,omitempty"`
}

// NewAcceptanceEvent creates a new acceptance event.
func NewAcceptanceEvent(command string, passed bool, output string, timeout int, prevHash string) AcceptanceEvent {
	return AcceptanceEvent{
		BaseEvent:      NewBaseEvent(EventTypeAcceptance, prevHash),
		Command:        command,
		Passed:         passed,
		Output:         output,
		TimeoutSeconds: timeout,
	}
}

// Validate checks required fields.
func (e AcceptanceEvent) Validate() error {
	if err := e.BaseEvent.Validate(); err != nil {
		return err
	}
	if e.Command == "" {
		return errors.New("command is required")
	}
	return nil
}

// ToJSONL returns the event as JSONL.
func (e AcceptanceEvent) ToJSONL() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ParseAcceptanceJSONL parses JSONL to AcceptanceEvent.
func ParseAcceptanceJSONL(line string) (AcceptanceEvent, error) {
	var e AcceptanceEvent
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return AcceptanceEvent{}, err
	}
	return e, nil
}
