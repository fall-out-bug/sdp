package hooks

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestHookEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		payload   map[string]interface{}
	}{
		{
			name:      "command pre event",
			eventType: "command:pre",
			payload:   map[string]interface{}{"command": "build", "args": []string{"--verbose"}},
		},
		{
			name:      "session start event",
			eventType: "session:start",
			payload:   map[string]interface{}{"session_id": "abc123"},
		},
		{
			name:      "gateway response event",
			eventType: "gateway:response",
			payload:   map[string]interface{}{"model": "claude", "tokens": 100},
		},
		{
			name:      "empty payload",
			eventType: "test:event",
			payload:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewEvent(tt.eventType, tt.payload)

			if event.Type != tt.eventType {
				t.Errorf("expected type %s, got %s", tt.eventType, event.Type)
			}
			if event.Timestamp.IsZero() {
				t.Error("timestamp should not be zero")
			}
			if tt.payload != nil && event.Payload == nil {
				t.Error("payload should not be nil")
			}
		})
	}
}

func TestHookRegistry_Subscribe(t *testing.T) {
	registry := NewRegistry()

	handler := func(event HookEvent) error { return nil }

	// Subscribe to event
	registry.Subscribe("command:pre", handler, 0)

	handlers := registry.GetHandlers("command:pre")
	if len(handlers) != 1 {
		t.Errorf("expected 1 handler, got %d", len(handlers))
	}

	// Subscribe multiple handlers
	registry.Subscribe("command:pre", handler, 0)
	handlers = registry.GetHandlers("command:pre")
	if len(handlers) != 2 {
		t.Errorf("expected 2 handlers, got %d", len(handlers))
	}
}

func TestHookRegistry_Publish(t *testing.T) {
	registry := NewRegistry()

	var called bool
	handler := func(event HookEvent) error {
		called = true
		return nil
	}

	registry.Subscribe("command:pre", handler, 0)

	event := NewEvent("command:pre", map[string]interface{}{"test": true})
	err := registry.Publish(event)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("handler was not called")
	}
}

func TestHookRegistry_Priority(t *testing.T) {
	registry := NewRegistry()

	var order []int

	// Subscribe with different priorities (lower = higher priority)
	registry.Subscribe("test", func(event HookEvent) error {
		order = append(order, 2)
		return nil
	}, 10) // Lower priority

	registry.Subscribe("test", func(event HookEvent) error {
		order = append(order, 1)
		return nil
	}, 0) // Higher priority

	registry.Subscribe("test", func(event HookEvent) error {
		order = append(order, 3)
		return nil
	}, 20) // Lowest priority

	event := NewEvent("test", nil)
	registry.Publish(event)

	// Handlers should execute in priority order
	expected := []int{1, 2, 3}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("position %d: expected %d, got %d (order: %v)", i, v, order[i], order)
		}
	}
}

func TestHookRegistry_ErrorIsolation(t *testing.T) {
	registry := NewRegistry()

	var callCount int

	// First handler errors
	registry.Subscribe("test", func(event HookEvent) error {
		callCount++
		return errors.New("handler error")
	}, 0)

	// Second handler should still execute
	registry.Subscribe("test", func(event HookEvent) error {
		callCount++
		return nil
	}, 1)

	event := NewEvent("test", nil)
	err := registry.Publish(event)

	// Should return error but both handlers called
	if err == nil {
		t.Error("expected error from first handler")
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestHookRegistry_AsyncPublish(t *testing.T) {
	registry := NewRegistry()

	var mu sync.Mutex
	var callCount int

	// Slow handler
	registry.Subscribe("test", func(event HookEvent) error {
		time.Sleep(50 * time.Millisecond)
		mu.Lock()
		callCount++
		mu.Unlock()
		return nil
	}, 0)

	event := NewEvent("test", nil)

	start := time.Now()
	err := registry.PublishAsync(context.Background(), event)
	elapsed := time.Since(start)

	// Should return quickly (async)
	if elapsed > 30*time.Millisecond {
		t.Error("PublishAsync should return quickly")
	}

	// Wait for async completion
	time.Sleep(100 * time.Millisecond)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mu.Lock()
	count := callCount
	mu.Unlock()

	if count != 1 {
		t.Errorf("expected 1 call, got %d", count)
	}
}

func TestHookRegistry_NoHandlers(t *testing.T) {
	registry := NewRegistry()

	event := NewEvent("nonexistent", nil)
	err := registry.Publish(event)

	// Should not error when no handlers
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHookRegistry_Unsubscribe(t *testing.T) {
	registry := NewRegistry()

	var callCount int
	handler := func(event HookEvent) error {
		callCount++
		return nil
	}

	id := registry.Subscribe("test", handler, 0)

	event := NewEvent("test", nil)
	registry.Publish(event)

	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Unsubscribe
	registry.Unsubscribe("test", id)

	callCount = 0
	registry.Publish(event)

	if callCount != 0 {
		t.Error("handler should not be called after unsubscribe")
	}
}

func TestHookRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewRegistry()

	var wg sync.WaitGroup

	// Concurrent subscriptions
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			registry.Subscribe("test", func(event HookEvent) error {
				return nil
			}, 0)
		}()
	}

	// Concurrent publishes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			event := NewEvent("test", nil)
			registry.Publish(event)
		}()
	}

	wg.Wait()

	// Should have 10 handlers
	handlers := registry.GetHandlers("test")
	if len(handlers) != 10 {
		t.Errorf("expected 10 handlers, got %d", len(handlers))
	}
}

func TestHookRegistry_GetEventTypes(t *testing.T) {
	registry := NewRegistry()

	registry.Subscribe("command:pre", func(event HookEvent) error { return nil }, 0)
	registry.Subscribe("command:post", func(event HookEvent) error { return nil }, 0)
	registry.Subscribe("session:start", func(event HookEvent) error { return nil }, 0)

	types := registry.GetEventTypes()

	if len(types) != 3 {
		t.Errorf("expected 3 event types, got %d", len(types))
	}

	// Check all types present
	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	expected := []string{"command:pre", "command:post", "session:start"}
	for _, e := range expected {
		if !typeMap[e] {
			t.Errorf("missing event type: %s", e)
		}
	}
}
