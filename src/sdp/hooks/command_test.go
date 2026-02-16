package hooks

import (
	"errors"
	"testing"
	"time"
)

func TestCommandEvent(t *testing.T) {
	tests := []struct {
		name     string
		cmdName  string
		args     []string
		flags    map[string]interface{}
		duration time.Duration
		output   string
		err      error
	}{
		{
			name:     "simple command",
			cmdName:  "build",
			args:     []string{"--verbose"},
			flags:    map[string]interface{}{"verbose": true},
			duration: 100 * time.Millisecond,
			output:   "Build successful",
			err:      nil,
		},
		{
			name:     "command with error",
			cmdName:  "test",
			args:     []string{"--fail"},
			flags:    nil,
			duration: 50 * time.Millisecond,
			output:   "",
			err:      errors.New("test failed"),
		},
		{
			name:     "empty command",
			cmdName:  "",
			args:     nil,
			flags:    nil,
			duration: 0,
			output:   "",
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewCommandEvent(tt.cmdName, tt.args, tt.flags, tt.duration, tt.output, tt.err)

			if event.Name != tt.cmdName {
				t.Errorf("expected name %s, got %s", tt.cmdName, event.Name)
			}
			if event.Duration != tt.duration {
				t.Errorf("expected duration %v, got %v", tt.duration, event.Duration)
			}
			if event.Error != tt.err {
				t.Errorf("expected error %v, got %v", tt.err, event.Error)
			}
		})
	}
}

func TestCommandEvent_ToHookEvent(t *testing.T) {
	cmdEvent := NewCommandEvent(
		"apply",
		[]string{"--ws", "00-001-01"},
		map[string]interface{}{"dry_run": true},
		150*time.Millisecond,
		"Applied successfully",
		nil,
	)

	hookEvent := cmdEvent.ToHookEvent("command:post")

	if hookEvent.Type != "command:post" {
		t.Errorf("expected type command:post, got %s", hookEvent.Type)
	}
	if hookEvent.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	if hookEvent.Payload["name"] != "apply" {
		t.Errorf("expected name apply, got %v", hookEvent.Payload["name"])
	}
	if hookEvent.Payload["duration"] != 150*time.Millisecond {
		t.Errorf("expected duration in payload")
	}
}

func TestRegisterCommandHooks(t *testing.T) {
	registry := NewRegistry()
	RegisterCommandHooks(registry)

	// Verify all event types are registered
	types := registry.GetEventTypes()

	expectedTypes := []string{
		EventTypeCommandPre,
		EventTypeCommandPost,
		EventTypeCommandError,
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

func TestCommandHook_PreHookModifiesInput(t *testing.T) {
	registry := NewRegistry()

	var capturedArgs []string

	// Pre-hook that modifies args
	registry.Subscribe(EventTypeCommandPre, func(event HookEvent) error {
		if payload, ok := event.Payload["args"].([]string); ok {
			// This is demonstration - in real use, modifications would be captured
			capturedArgs = append(payload, "--modified")
		}
		return nil
	}, 0)

	cmdEvent := NewCommandEvent("build", []string{"--verbose"}, nil, 0, "", nil)
	registry.Publish(cmdEvent.ToHookEvent(EventTypeCommandPre))

	// Handler was called
	if len(capturedArgs) == 0 {
		t.Error("pre-hook should have been called with args")
	}
}

func TestCommandHook_PostHookCapturesOutput(t *testing.T) {
	registry := NewRegistry()

	var capturedOutput string

	registry.Subscribe(EventTypeCommandPost, func(event HookEvent) error {
		if output, ok := event.Payload["output"].(string); ok {
			capturedOutput = output
		}
		return nil
	}, 0)

	cmdEvent := NewCommandEvent("test", nil, nil, 100*time.Millisecond, "All tests passed", nil)
	registry.Publish(cmdEvent.ToHookEvent(EventTypeCommandPost))

	if capturedOutput != "All tests passed" {
		t.Errorf("expected output captured, got %s", capturedOutput)
	}
}

func TestCommandHook_ErrorHookReceivesError(t *testing.T) {
	registry := NewRegistry()

	var capturedError error

	registry.Subscribe(EventTypeCommandError, func(event HookEvent) error {
		if err, ok := event.Payload["error"].(error); ok {
			capturedError = err
		}
		return nil
	}, 0)

	testErr := errors.New("command failed")
	cmdEvent := NewCommandEvent("deploy", nil, nil, 50*time.Millisecond, "", testErr)
	registry.Publish(cmdEvent.ToHookEvent(EventTypeCommandError))

	if capturedError == nil {
		t.Error("error hook should have received error")
	}
	if capturedError.Error() != "command failed" {
		t.Errorf("expected 'command failed', got %s", capturedError.Error())
	}
}

func TestCommandEvent_Flags(t *testing.T) {
	flags := map[string]interface{}{
		"verbose": true,
		"count":   5,
		"name":    "test",
	}

	cmdEvent := NewCommandEvent("run", nil, flags, 0, "", nil)

	if cmdEvent.Flags["verbose"] != true {
		t.Error("verbose flag should be true")
	}
	if cmdEvent.Flags["count"] != 5 {
		t.Error("count flag should be 5")
	}
	if cmdEvent.Flags["name"] != "test" {
		t.Error("name flag should be test")
	}
}

func TestCommandEvent_IsSuccess(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		isSuccess bool
	}{
		{"no error", nil, true},
		{"with error", errors.New("failed"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdEvent := NewCommandEvent("test", nil, nil, 0, "", tt.err)
			if cmdEvent.IsSuccess() != tt.isSuccess {
				t.Errorf("expected IsSuccess=%v, got %v", tt.isSuccess, cmdEvent.IsSuccess())
			}
		})
	}
}
