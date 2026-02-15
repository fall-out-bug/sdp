package hooks

import (
	"errors"
	"testing"
	"time"
)

func TestGatewayEvent(t *testing.T) {
	tests := []struct {
		name         string
		model        string
		promptHash   string
		inputTokens  int
		outputTokens int
		latency      time.Duration
		cost         float64
		err          error
	}{
		{
			name:         "successful response",
			model:        "claude-3-opus",
			promptHash:   "abc123",
			inputTokens:  1000,
			outputTokens: 500,
			latency:      2 * time.Second,
			cost:         0.15,
			err:          nil,
		},
		{
			name:         "api error",
			model:        "claude-3-sonnet",
			promptHash:   "def456",
			inputTokens:  500,
			outputTokens: 0,
			latency:      500 * time.Millisecond,
			cost:         0,
			err:          errors.New("rate limit exceeded"),
		},
		{
			name:         "empty event",
			model:        "",
			promptHash:   "",
			inputTokens:  0,
			outputTokens: 0,
			latency:      0,
			cost:         0,
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewGatewayEvent(tt.model, tt.promptHash, tt.inputTokens, tt.outputTokens, tt.latency, tt.cost, tt.err)

			if event.Model != tt.model {
				t.Errorf("expected model %s, got %s", tt.model, event.Model)
			}
			if event.InputTokens != tt.inputTokens {
				t.Errorf("expected input tokens %d, got %d", tt.inputTokens, event.InputTokens)
			}
			if event.Cost != tt.cost {
				t.Errorf("expected cost %f, got %f", tt.cost, event.Cost)
			}
		})
	}
}

func TestGatewayEvent_ToHookEvent(t *testing.T) {
	gatewayEvent := NewGatewayEvent(
		"claude-3-opus",
		"hash789",
		2000,
		1000,
		3*time.Second,
		0.25,
		nil,
	)

	hookEvent := gatewayEvent.ToHookEvent(EventTypeGatewayResponse)

	if hookEvent.Type != EventTypeGatewayResponse {
		t.Errorf("expected type %s, got %s", EventTypeGatewayResponse, hookEvent.Type)
	}
	if hookEvent.Payload["model"] != "claude-3-opus" {
		t.Errorf("expected model in payload")
	}
	if hookEvent.Payload["cost"] != 0.25 {
		t.Errorf("expected cost in payload")
	}
}

func TestRegisterGatewayHooks(t *testing.T) {
	registry := NewRegistry()
	RegisterGatewayHooks(registry)

	types := registry.GetEventTypes()

	expectedTypes := []string{
		EventTypeGatewayRequest,
		EventTypeGatewayResponse,
		EventTypeGatewayError,
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

func TestGatewayHook_RequestEvent(t *testing.T) {
	registry := NewRegistry()

	var capturedModel string

	registry.Subscribe(EventTypeGatewayRequest, func(event HookEvent) error {
		if model, ok := event.Payload["model"].(string); ok {
			capturedModel = model
		}
		return nil
	}, 0)

	gatewayEvent := NewGatewayEvent("claude-3-haiku", "hash", 100, 0, 0, 0, nil)
	registry.Publish(gatewayEvent.ToHookEvent(EventTypeGatewayRequest))

	if capturedModel != "claude-3-haiku" {
		t.Errorf("expected model claude-3-haiku, got %s", capturedModel)
	}
}

func TestGatewayHook_ResponseEvent(t *testing.T) {
	registry := NewRegistry()

	var totalTokens int

	registry.Subscribe(EventTypeGatewayResponse, func(event HookEvent) error {
		input := event.Payload["input_tokens"].(int)
		output := event.Payload["output_tokens"].(int)
		totalTokens = input + output
		return nil
	}, 0)

	gatewayEvent := NewGatewayEvent("claude-3-opus", "hash", 1000, 500, 2*time.Second, 0.10, nil)
	registry.Publish(gatewayEvent.ToHookEvent(EventTypeGatewayResponse))

	if totalTokens != 1500 {
		t.Errorf("expected 1500 tokens, got %d", totalTokens)
	}
}

func TestGatewayHook_ErrorEvent(t *testing.T) {
	registry := NewRegistry()

	var capturedError error

	registry.Subscribe(EventTypeGatewayError, func(event HookEvent) error {
		if err, ok := event.Payload["error"].(error); ok {
			capturedError = err
		}
		return nil
	}, 0)

	testErr := errors.New("context length exceeded")
	gatewayEvent := NewGatewayEvent("claude-3-sonnet", "hash", 100000, 0, 0, 0, testErr)
	registry.Publish(gatewayEvent.ToHookEvent(EventTypeGatewayError))

	if capturedError == nil {
		t.Error("expected error to be captured")
	}
	if capturedError.Error() != "context length exceeded" {
		t.Errorf("expected 'context length exceeded', got %s", capturedError.Error())
	}
}

func TestGatewayEvent_TotalTokens(t *testing.T) {
	event := NewGatewayEvent("model", "hash", 1000, 500, 0, 0, nil)

	if event.TotalTokens() != 1500 {
		t.Errorf("expected 1500 total tokens, got %d", event.TotalTokens())
	}
}

func TestGatewayEvent_IsSuccess(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		isSuccess bool
	}{
		{"success", nil, true},
		{"failure", errors.New("failed"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewGatewayEvent("model", "hash", 0, 0, 0, 0, tt.err)
			if event.IsSuccess() != tt.isSuccess {
				t.Errorf("expected IsSuccess=%v, got %v", tt.isSuccess, event.IsSuccess())
			}
		})
	}
}

func TestGatewayEvent_CostPerToken(t *testing.T) {
	tests := []struct {
		name         string
		cost         float64
		totalTokens  int
		expectedCost float64
	}{
		{"normal calculation", 0.15, 1000, 0.00015},
		{"zero tokens", 0.10, 0, 0},
		{"zero cost", 0, 1000, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewGatewayEvent("model", "hash", tt.totalTokens, 0, 0, tt.cost, nil)
			costPerToken := event.CostPerToken()

			if costPerToken != tt.expectedCost {
				t.Errorf("expected cost per token %f, got %f", tt.expectedCost, costPerToken)
			}
		})
	}
}

func TestGatewayEvent_IsRateLimited(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"rate limit error", errors.New("rate limit exceeded"), true},
		{"other error", errors.New("network error"), false},
		{"no error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewGatewayEvent("model", "hash", 0, 0, 0, 0, tt.err)
			if event.IsRateLimited() != tt.expected {
				t.Errorf("expected IsRateLimited=%v, got %v", tt.expected, event.IsRateLimited())
			}
		})
	}
}
