package hooks

import (
	"time"
)

// Event type constants for session lifecycle.
const (
	EventTypeSessionStart   = "session:start"
	EventTypeSessionEnd     = "session:end"
	EventTypeSessionCompact = "session:compact"
	EventTypeSessionResume  = "session:resume"
)

// SessionEvent represents a session lifecycle event.
type SessionEvent struct {
	SessionID   string    // Unique session identifier
	StartTime   time.Time // When session started
	EndTime     time.Time // When session ended (zero if active)
	WorkDone    []string  // List of completed work items
	TokensUsed  int       // Total tokens consumed
}

// NewSessionEvent creates a new session event.
func NewSessionEvent(sessionID string, startTime, endTime time.Time, workDone []string, tokensUsed int) SessionEvent {
	return SessionEvent{
		SessionID:  sessionID,
		StartTime:  startTime,
		EndTime:    endTime,
		WorkDone:   workDone,
		TokensUsed: tokensUsed,
	}
}

// ToHookEvent converts a SessionEvent to a HookEvent.
func (e SessionEvent) ToHookEvent(eventType string) HookEvent {
	return NewEvent(eventType, map[string]interface{}{
		"session_id":  e.SessionID,
		"start_time":  e.StartTime,
		"end_time":    e.EndTime,
		"work_done":   e.WorkDone,
		"tokens_used": e.TokensUsed,
	})
}

// Duration returns the session duration.
// Returns 0 if session hasn't ended.
func (e SessionEvent) Duration() time.Duration {
	if e.EndTime.IsZero() {
		return 0
	}
	return e.EndTime.Sub(e.StartTime)
}

// IsActive returns true if the session is still active.
func (e SessionEvent) IsActive() bool {
	return e.EndTime.IsZero()
}

// AddWork appends a work item to the session.
func (e *SessionEvent) AddWork(item string) {
	e.WorkDone = append(e.WorkDone, item)
}

// RegisterSessionHooks registers default handlers for session lifecycle events.
func RegisterSessionHooks(registry *HookRegistry) {
	registry.Subscribe(EventTypeSessionStart, func(event HookEvent) error {
		return nil
	}, 1000)

	registry.Subscribe(EventTypeSessionEnd, func(event HookEvent) error {
		return nil
	}, 1000)

	registry.Subscribe(EventTypeSessionCompact, func(event HookEvent) error {
		return nil
	}, 1000)

	registry.Subscribe(EventTypeSessionResume, func(event HookEvent) error {
		return nil
	}, 1000)
}
