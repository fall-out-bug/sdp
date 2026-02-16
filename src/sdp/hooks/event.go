// Package hooks provides a runtime hook system for lifecycle events.
package hooks

import (
	"time"
)

// HookEvent represents a lifecycle event in the system.
type HookEvent struct {
	Type      string                 // Event type (e.g., "command:pre", "session:start")
	Timestamp time.Time              // When the event occurred
	Payload   map[string]interface{} // Event-specific data
}

// NewEvent creates a new hook event with the current timestamp.
func NewEvent(eventType string, payload map[string]interface{}) HookEvent {
	return HookEvent{
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

// HookHandler is a function that processes a hook event.
type HookHandler func(event HookEvent) error

// handlerEntry wraps a handler with its priority and ID.
type handlerEntry struct {
	id       string
	handler  HookHandler
	priority int
}
