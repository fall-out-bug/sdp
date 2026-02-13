package notification

import (
	"testing"
)

func TestLogChannel_Send(t *testing.T) {
	channel := NewLogChannel()

	if channel.Name() != "log" {
		t.Errorf("Expected name 'log', got %s", channel.Name())
	}

	if !channel.IsEnabled() {
		t.Error("Log channel should be enabled by default")
	}

	notification := &Notification{
		Type:    NotifyFeatureComplete,
		Message: "Test notification",
	}

	err := channel.Send(notification)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

func TestLogChannel_Disabled(t *testing.T) {
	channel := NewLogChannel()
	channel.Enabled = false

	if channel.IsEnabled() {
		t.Error("Channel should be disabled")
	}
}

func TestWebhookChannel_Send(t *testing.T) {
	channel := NewWebhookChannel("https://example.com/webhook")

	if channel.Name() != "webhook" {
		t.Errorf("Expected name 'webhook', got %s", channel.Name())
	}

	// Webhook without URL should fail gracefully
	channel = NewWebhookChannel("")
	notification := &Notification{
		Type:    NotifyFeatureComplete,
		Message: "Test notification",
	}

	// Should not panic with empty URL
	_ = channel.Send(notification)
}

func TestRenderTemplate(t *testing.T) {
	template := "Feature {feature} completed with {status}"
	vars := map[string]string{
		"feature": "F051",
		"status":  "success",
	}

	result := RenderTemplate(template, vars)

	if result != "Feature F051 completed with success" {
		t.Errorf("Template render failed, got: %s", result)
	}
}

func TestRenderTemplate_MissingVar(t *testing.T) {
	template := "Feature {feature} {missing}"
	vars := map[string]string{
		"feature": "F051",
	}

	result := RenderTemplate(template, vars)

	// Missing variable should remain as placeholder
	if result != "Feature F051 {missing}" {
		t.Errorf("Expected placeholder to remain, got: %s", result)
	}
}
