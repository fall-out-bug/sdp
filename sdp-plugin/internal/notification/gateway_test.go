package notification

import (
	"testing"
)

func TestNotification_Types(t *testing.T) {
	types := []NotificationType{
		NotifyFeatureComplete,
		NotifyFeatureFailed,
		NotifyDriftDetected,
		NotifyReviewComplete,
		NotifyDeployComplete,
		NotifyAgentError,
	}

	for _, nt := range types {
		if string(nt) == "" {
			t.Error("Notification type should not be empty")
		}
	}
}

func TestNotification_Severities(t *testing.T) {
	severities := []Severity{
		SeverityInfo,
		SeverityWarning,
		SeverityError,
		SeverityCritical,
	}

	for _, s := range severities {
		if string(s) == "" {
			t.Error("Severity should not be empty")
		}
	}
}

func TestGateway_Send(t *testing.T) {
	gateway := NewGateway()
	gateway.AddChannel(&MockChannel{name: "test", enabled: true})

	notification := &Notification{
		Type:    NotifyFeatureComplete,
		Message: "Feature F051 completed",
	}

	err := gateway.Send(notification)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

func TestGateway_SendDisabled(t *testing.T) {
	gateway := NewGateway()
	gateway.AddChannel(&MockChannel{name: "disabled", enabled: false})

	notification := &Notification{
		Type:    NotifyFeatureComplete,
		Message: "Feature F051 completed",
	}

	err := gateway.Send(notification)
	if err != nil {
		t.Fatalf("Send should not fail for disabled channel: %v", err)
	}
}

func TestGateway_RateLimit(t *testing.T) {
	gateway := NewGatewayWithRateLimit(2, 60) // 2 per minute

	// First two should succeed
	for i := 0; i < 2; i++ {
		if !gateway.AllowSend("test-feature") {
			t.Errorf("Notification %d should be allowed", i+1)
		}
	}

	// Third should be rate limited
	if gateway.AllowSend("test-feature") {
		t.Error("Third notification should be rate limited")
	}
}

func TestGateway_RateLimit_CriticalBypass(t *testing.T) {
	gateway := NewGatewayWithRateLimit(1, 60) // 1 per minute

	// Use up the limit
	gateway.AllowSend("test-feature")

	// Critical should still be allowed
	if !gateway.AllowSendCritical() {
		t.Error("Critical notification should bypass rate limit")
	}
}

func TestGateway_History(t *testing.T) {
	gateway := NewGateway()
	gateway.AddChannel(&MockChannel{name: "test", enabled: true})

	// Send multiple notifications
	for i := 0; i < 3; i++ {
		gateway.Send(&Notification{
			Type:    NotifyFeatureComplete,
			Message: "Test notification",
		})
	}

	history := gateway.GetHistory()
	if len(history) != 3 {
		t.Errorf("Expected 3 history entries, got %d", len(history))
	}
}

// MockChannel for testing
type MockChannel struct {
	name    string
	enabled bool
}

func (c *MockChannel) Name() string    { return c.name }
func (c *MockChannel) IsEnabled() bool { return c.enabled }
func (c *MockChannel) Send(n *Notification) error {
	return nil
}
