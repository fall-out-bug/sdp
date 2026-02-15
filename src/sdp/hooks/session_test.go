package hooks

import (
	"testing"
	"time"
)

func TestSessionEvent(t *testing.T) {
	tests := []struct {
		name       string
		sessionID  string
		startTime  time.Time
		endTime    time.Time
		workDone   []string
		tokensUsed int
	}{
		{
			name:       "active session",
			sessionID:  "sess-123",
			startTime:  time.Now().Add(-1 * time.Hour),
			endTime:    time.Time{},
			workDone:   []string{"implemented feature", "fixed bug"},
			tokensUsed: 5000,
		},
		{
			name:       "ended session",
			sessionID:  "sess-456",
			startTime:  time.Now().Add(-2 * time.Hour),
			endTime:    time.Now().Add(-30 * time.Minute),
			workDone:   []string{"code review"},
			tokensUsed: 2000,
		},
		{
			name:       "empty session",
			sessionID:  "",
			startTime:  time.Time{},
			endTime:    time.Time{},
			workDone:   nil,
			tokensUsed: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewSessionEvent(tt.sessionID, tt.startTime, tt.endTime, tt.workDone, tt.tokensUsed)

			if event.SessionID != tt.sessionID {
				t.Errorf("expected session ID %s, got %s", tt.sessionID, event.SessionID)
			}
			if event.TokensUsed != tt.tokensUsed {
				t.Errorf("expected tokens %d, got %d", tt.tokensUsed, event.TokensUsed)
			}
		})
	}
}

func TestSessionEvent_ToHookEvent(t *testing.T) {
	startTime := time.Now().Add(-1 * time.Hour)
	workDone := []string{"task1", "task2"}

	sessionEvent := NewSessionEvent("sess-789", startTime, time.Time{}, workDone, 3000)

	hookEvent := sessionEvent.ToHookEvent(EventTypeSessionStart)

	if hookEvent.Type != EventTypeSessionStart {
		t.Errorf("expected type %s, got %s", EventTypeSessionStart, hookEvent.Type)
	}
	if hookEvent.Payload["session_id"] != "sess-789" {
		t.Errorf("expected session_id in payload")
	}
	if hookEvent.Payload["tokens_used"] != 3000 {
		t.Errorf("expected tokens_used in payload")
	}
}

func TestRegisterSessionHooks(t *testing.T) {
	registry := NewRegistry()
	RegisterSessionHooks(registry)

	types := registry.GetEventTypes()

	expectedTypes := []string{
		EventTypeSessionStart,
		EventTypeSessionEnd,
		EventTypeSessionCompact,
		EventTypeSessionResume,
	}

	for _, expected := range expectedTypes {
		found := false
		for _, actual := range types {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing event type: %s", expected)
		}
	}
}

func TestSessionHook_StartEvent(t *testing.T) {
	registry := NewRegistry()

	var capturedID string

	registry.Subscribe(EventTypeSessionStart, func(event HookEvent) error {
		if id, ok := event.Payload["session_id"].(string); ok {
			capturedID = id
		}
		return nil
	}, 0)

	sessionEvent := NewSessionEvent("sess-start", time.Now(), time.Time{}, nil, 0)
	registry.Publish(sessionEvent.ToHookEvent(EventTypeSessionStart))

	if capturedID != "sess-start" {
		t.Errorf("expected session_id sess-start, got %s", capturedID)
	}
}

func TestSessionHook_EndEvent(t *testing.T) {
	registry := NewRegistry()

	var workDone []string

	registry.Subscribe(EventTypeSessionEnd, func(event HookEvent) error {
		if work, ok := event.Payload["work_done"].([]string); ok {
			workDone = work
		}
		return nil
	}, 0)

	sessionEvent := NewSessionEvent("sess-end", time.Now(), time.Now(), []string{"built feature", "wrote tests"}, 5000)
	registry.Publish(sessionEvent.ToHookEvent(EventTypeSessionEnd))

	if len(workDone) != 2 {
		t.Errorf("expected 2 work items, got %d", len(workDone))
	}
}

func TestSessionHook_CompactEvent(t *testing.T) {
	registry := NewRegistry()

	var tokensBefore int

	registry.Subscribe(EventTypeSessionCompact, func(event HookEvent) error {
		if tokens, ok := event.Payload["tokens_used"].(int); ok {
			tokensBefore = tokens
		}
		return nil
	}, 0)

	sessionEvent := NewSessionEvent("sess-compact", time.Now(), time.Time{}, nil, 100000)
	registry.Publish(sessionEvent.ToHookEvent(EventTypeSessionCompact))

	if tokensBefore != 100000 {
		t.Errorf("expected 100000 tokens, got %d", tokensBefore)
	}
}

func TestSessionHook_ResumeEvent(t *testing.T) {
	registry := NewRegistry()

	var resumed bool

	registry.Subscribe(EventTypeSessionResume, func(event HookEvent) error {
		resumed = true
		return nil
	}, 0)

	sessionEvent := NewSessionEvent("sess-resume", time.Now(), time.Time{}, nil, 0)
	registry.Publish(sessionEvent.ToHookEvent(EventTypeSessionResume))

	if !resumed {
		t.Error("resume hook should have been called")
	}
}

func TestSessionEvent_Duration(t *testing.T) {
	start := time.Now().Add(-2 * time.Hour)
	end := time.Now()

	sessionEvent := NewSessionEvent("sess-duration", start, end, nil, 0)

	duration := sessionEvent.Duration()

	if duration < 1*time.Hour || duration > 3*time.Hour {
		t.Errorf("expected duration around 2 hours, got %v", duration)
	}
}

func TestSessionEvent_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		endTime  time.Time
		isActive bool
	}{
		{"active session", time.Time{}, true},
		{"ended session", time.Now(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionEvent := NewSessionEvent("sess", time.Now(), tt.endTime, nil, 0)
			if sessionEvent.IsActive() != tt.isActive {
				t.Errorf("expected IsActive=%v, got %v", tt.isActive, sessionEvent.IsActive())
			}
		})
	}
}

func TestSessionEvent_AddWork(t *testing.T) {
	sessionEvent := NewSessionEvent("sess", time.Now(), time.Time{}, nil, 0)

	sessionEvent.AddWork("task1")
	sessionEvent.AddWork("task2")

	if len(sessionEvent.WorkDone) != 2 {
		t.Errorf("expected 2 work items, got %d", len(sessionEvent.WorkDone))
	}
}
