package coordination

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

// ComputeHash computes SHA256 hash of the event (excluding Hash field)
func (e *AgentEvent) ComputeHash() string {
	// Create a copy without the Hash field for hashing
	data := map[string]interface{}{
		"id":        e.ID,
		"type":      e.Type,
		"agent_id":  e.AgentID,
		"role":      e.Role,
		"task_id":   e.TaskID,
		"timestamp": e.Timestamp.Format(time.RFC3339Nano),
		"prev_hash": e.PrevHash,
	}

	// Add payload if present
	if e.Payload != nil {
		data["payload"] = e.Payload
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

// NewAgentEvent creates a new agent event with generated ID and timestamp
func NewAgentEvent(eventType, agentID, role string) *AgentEvent {
	return &AgentEvent{
		ID:        generateEventID(),
		Type:      eventType,
		AgentID:   agentID,
		Role:      role,
		Timestamp: time.Now(),
		Payload:   make(map[string]interface{}),
	}
}

// generateEventID generates a unique event ID
func generateEventID() string {
	ts := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(time.Now().String()))
	return hex.EncodeToString(hash[:8]) + "-" + time.Unix(0, ts).Format("20060102150405")
}
