package hooks

import (
	"time"
)

// Event type constants for command lifecycle.
const (
	EventTypeCommandPre   = "command:pre"
	EventTypeCommandPost  = "command:post"
	EventTypeCommandError = "command:error"
)

// CommandEvent represents a command execution event.
type CommandEvent struct {
	Name     string                 // Command name (e.g., "build", "apply")
	Args     []string               // Command arguments
	Flags    map[string]interface{} // Command flags/options
	Duration time.Duration          // Execution duration
	Output   string                 // Command output
	Error    error                  // Error if command failed
}

// NewCommandEvent creates a new command event.
func NewCommandEvent(name string, args []string, flags map[string]interface{}, duration time.Duration, output string, err error) CommandEvent {
	return CommandEvent{
		Name:     name,
		Args:     args,
		Flags:    flags,
		Duration: duration,
		Output:   output,
		Error:    err,
	}
}

// ToHookEvent converts a CommandEvent to a HookEvent.
func (e CommandEvent) ToHookEvent(eventType string) HookEvent {
	return NewEvent(eventType, map[string]interface{}{
		"name":     e.Name,
		"args":     e.Args,
		"flags":    e.Flags,
		"duration": e.Duration,
		"output":   e.Output,
		"error":    e.Error,
	})
}

// IsSuccess returns true if the command completed without error.
func (e CommandEvent) IsSuccess() bool {
	return e.Error == nil
}

// RegisterCommandHooks registers default handlers for command lifecycle events.
// This sets up the event types in the registry but doesn't add business logic.
func RegisterCommandHooks(registry *HookRegistry) {
	// Register placeholder handlers to ensure event types exist
	registry.Subscribe(EventTypeCommandPre, func(event HookEvent) error {
		return nil
	}, 1000)

	registry.Subscribe(EventTypeCommandPost, func(event HookEvent) error {
		return nil
	}, 1000)

	registry.Subscribe(EventTypeCommandError, func(event HookEvent) error {
		return nil
	}, 1000)
}
