package hooks

import (
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"
)

// HookRegistry manages hook subscriptions and event publishing.
type HookRegistry struct {
	mu       sync.RWMutex
	handlers map[string][]handlerEntry
}

// NewRegistry creates a new hook registry.
func NewRegistry() *HookRegistry {
	return &HookRegistry{
		handlers: make(map[string][]handlerEntry),
	}
}

// Subscribe registers a handler for an event type with a priority.
// Lower priority values execute first. Returns a subscription ID.
func (r *HookRegistry) Subscribe(eventType string, handler HookHandler, priority int) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	entry := handlerEntry{
		id:       id,
		handler:  handler,
		priority: priority,
	}

	r.handlers[eventType] = append(r.handlers[eventType], entry)

	// Sort by priority (ascending - lower values first)
	sort.Slice(r.handlers[eventType], func(i, j int) bool {
		return r.handlers[eventType][i].priority < r.handlers[eventType][j].priority
	})

	return id
}

// Unsubscribe removes a handler by its subscription ID.
func (r *HookRegistry) Unsubscribe(eventType, id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	handlers := r.handlers[eventType]
	for i, h := range handlers {
		if h.id == id {
			r.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Publish sends an event to all subscribed handlers synchronously.
// Errors from handlers are collected but don't stop execution.
func (r *HookRegistry) Publish(event HookEvent) error {
	r.mu.RLock()
	handlers := make([]handlerEntry, len(r.handlers[event.Type]))
	copy(handlers, r.handlers[event.Type])
	r.mu.RUnlock()

	var lastErr error
	for _, entry := range handlers {
		if err := entry.handler(event); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// PublishAsync sends an event to all handlers asynchronously.
func (r *HookRegistry) PublishAsync(ctx context.Context, event HookEvent) error {
	r.mu.RLock()
	handlers := make([]handlerEntry, len(r.handlers[event.Type]))
	copy(handlers, r.handlers[event.Type])
	r.mu.RUnlock()

	go func() {
		for _, entry := range handlers {
			select {
			case <-ctx.Done():
				return
			default:
				entry.handler(event)
			}
		}
	}()

	return nil
}

// GetHandlers returns all handlers for an event type (for testing).
func (r *HookRegistry) GetHandlers(eventType string) []handlerEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.handlers[eventType]
}

// GetEventTypes returns all registered event types.
func (r *HookRegistry) GetEventTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.handlers))
	for t := range r.handlers {
		types = append(types, t)
	}
	return types
}
