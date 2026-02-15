package evidence

import (
	"encoding/json"
	"errors"
)

// ApprovalEvent records approval of work.
type ApprovalEvent struct {
	BaseEvent
	Approver   string `json:"approver"`
	Mode       string `json:"mode"` // "auto" or "human"
	Reasoning  string `json:"reasoning,omitempty"`
	Workstream string `json:"workstream,omitempty"`
}

// NewApprovalEvent creates a new approval event.
func NewApprovalEvent(approver, mode, prevHash string) ApprovalEvent {
	return ApprovalEvent{
		BaseEvent: NewBaseEvent(EventTypeApproval, prevHash),
		Approver:  approver,
		Mode:      mode,
	}
}

// Validate checks required fields.
func (e ApprovalEvent) Validate() error {
	if err := e.BaseEvent.Validate(); err != nil {
		return err
	}
	if e.Approver == "" {
		return errors.New("approver is required")
	}
	if e.Mode != "auto" && e.Mode != "human" {
		return errors.New("mode must be 'auto' or 'human'")
	}
	return nil
}

// ToJSONL returns the event as JSONL.
func (e ApprovalEvent) ToJSONL() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ParseApprovalJSONL parses JSONL to ApprovalEvent.
func ParseApprovalJSONL(line string) (ApprovalEvent, error) {
	var e ApprovalEvent
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return ApprovalEvent{}, err
	}
	return e, nil
}
