package hooks

import (
	"strings"
	"time"
)

// Event type constants for gateway lifecycle.
const (
	EventTypeGatewayRequest  = "gateway:request"
	EventTypeGatewayResponse = "gateway:response"
	EventTypeGatewayError    = "gateway:error"
)

// GatewayEvent represents an AI gateway API event.
type GatewayEvent struct {
	Model        string        // Model identifier (e.g., "claude-3-opus")
	PromptHash   string        // Hash of the prompt for tracking
	InputTokens  int           // Tokens in the request
	OutputTokens int           // Tokens in the response
	Latency      time.Duration // API call duration
	Cost         float64       // Calculated cost in dollars
	Error        error         // Error if call failed
}

// NewGatewayEvent creates a new gateway event.
func NewGatewayEvent(model, promptHash string, inputTokens, outputTokens int, latency time.Duration, cost float64, err error) GatewayEvent {
	return GatewayEvent{
		Model:        model,
		PromptHash:   promptHash,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		Latency:      latency,
		Cost:         cost,
		Error:        err,
	}
}

// ToHookEvent converts a GatewayEvent to a HookEvent.
func (e GatewayEvent) ToHookEvent(eventType string) HookEvent {
	return NewEvent(eventType, map[string]interface{}{
		"model":         e.Model,
		"prompt_hash":   e.PromptHash,
		"input_tokens":  e.InputTokens,
		"output_tokens": e.OutputTokens,
		"latency":       e.Latency,
		"cost":          e.Cost,
		"error":         e.Error,
	})
}

// TotalTokens returns the sum of input and output tokens.
func (e GatewayEvent) TotalTokens() int {
	return e.InputTokens + e.OutputTokens
}

// IsSuccess returns true if the API call succeeded.
func (e GatewayEvent) IsSuccess() bool {
	return e.Error == nil
}

// CostPerToken returns the cost per token.
// Returns 0 if no tokens were used.
func (e GatewayEvent) CostPerToken() float64 {
	total := e.TotalTokens()
	if total == 0 {
		return 0
	}
	return e.Cost / float64(total)
}

// IsRateLimited returns true if the error indicates rate limiting.
func (e GatewayEvent) IsRateLimited() bool {
	if e.Error == nil {
		return false
	}
	return strings.Contains(e.Error.Error(), "rate limit")
}

// RegisterGatewayHooks registers default handlers for gateway lifecycle events.
func RegisterGatewayHooks(registry *HookRegistry) {
	registry.Subscribe(EventTypeGatewayRequest, func(event HookEvent) error {
		return nil
	}, 1000)

	registry.Subscribe(EventTypeGatewayResponse, func(event HookEvent) error {
		return nil
	}, 1000)

	registry.Subscribe(EventTypeGatewayError, func(event HookEvent) error {
		return nil
	}, 1000)
}
