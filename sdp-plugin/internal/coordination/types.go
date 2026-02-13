package coordination

import (
	"time"
)

// Agent event types (AC1)
const (
	EventTypeAgentStart    = "agent_start"
	EventTypeAgentAction   = "agent_action"
	EventTypeAgentComplete = "agent_complete"
	EventTypeAgentError    = "agent_error"
	EventTypeAgentHandoff  = "agent_handoff"
)

// validEventTypes is the set of valid event types
var validEventTypes = map[string]bool{
	EventTypeAgentStart:    true,
	EventTypeAgentAction:   true,
	EventTypeAgentComplete: true,
	EventTypeAgentError:    true,
	EventTypeAgentHandoff:  true,
}

// isValidEventType checks if the event type is valid
func isValidEventType(t string) bool {
	return validEventTypes[t]
}

// AgentEvent represents an agent coordination event (AC2)
type AgentEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	AgentID   string                 `json:"agent_id"`
	Role      string                 `json:"role"`
	TaskID    string                 `json:"task_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	PrevHash  string                 `json:"prev_hash,omitempty"`
	Hash      string                 `json:"hash,omitempty"`
}

// Validate validates the agent event
func (e *AgentEvent) Validate() error {
	if e.ID == "" {
		return ErrMissingID
	}
	if e.AgentID == "" {
		return ErrMissingAgentID
	}
	if !isValidEventType(e.Type) {
		return ErrInvalidEventType
	}
	return nil
}
