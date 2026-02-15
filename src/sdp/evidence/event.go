// Package evidence provides the evidence layer for SDP.
// It implements an append-only log with hash chain for provenance.
package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of evidence event.
type EventType string

// Event type constants.
const (
	EventTypeDecision     EventType = "decision"
	EventTypePlan         EventType = "plan"
	EventTypeGeneration   EventType = "generation"
	EventTypeVerification EventType = "verification"
	EventTypeAcceptance   EventType = "acceptance"
	EventTypeApproval     EventType = "approval"
)

// ValidEventTypes returns all valid event types.
func ValidEventTypes() []EventType {
	return []EventType{
		EventTypeDecision,
		EventTypePlan,
		EventTypeGeneration,
		EventTypeVerification,
		EventTypeAcceptance,
		EventTypeApproval,
	}
}

// IsValid returns true if the event type is valid.
func (et EventType) IsValid() bool {
	for _, t := range ValidEventTypes() {
		if et == t {
			return true
		}
	}
	return false
}

// BaseEvent is the common structure for all evidence events.
type BaseEvent struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	PrevHash  string    `json:"prev_hash"`
	Hash      string    `json:"hash"`
}

// NewBaseEvent creates a new base event with generated ID and calculated hash.
func NewBaseEvent(eventType EventType, prevHash string) BaseEvent {
	e := BaseEvent{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		PrevHash:  prevHash,
	}
	e.CalculateHash()
	return e
}

// CalculateHash computes the hash for this event.
func (e *BaseEvent) CalculateHash() {
	data := e.ID + string(e.Type) + e.Timestamp.Format(time.RFC3339Nano) + e.PrevHash
	hash := sha256.Sum256([]byte(data))
	e.Hash = hex.EncodeToString(hash[:])
}

// Validate checks if the event has all required fields.
func (e BaseEvent) Validate() error {
	if e.ID == "" {
		return errors.New("event ID is required")
	}
	if !e.Type.IsValid() {
		return errors.New("invalid event type")
	}
	if e.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if e.Hash == "" {
		return errors.New("hash is required")
	}
	return nil
}

// ToJSONL returns the event as a JSONL line (JSON + newline).
func (e BaseEvent) ToJSONL() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ParseJSONL parses a JSONL line into a BaseEvent.
func ParseJSONL(line string) (BaseEvent, error) {
	var e BaseEvent
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return BaseEvent{}, err
	}
	return e, nil
}
