package notification

import (
	"sync"
	"time"

	"github.com/fall-out-bug/sdp/internal/safetylog"
)

// NotificationType defines types of notifications (AC1)
type NotificationType string

const (
	NotifyFeatureComplete NotificationType = "feature_complete"
	NotifyFeatureFailed   NotificationType = "feature_failed"
	NotifyDriftDetected   NotificationType = "drift_detected"
	NotifyReviewComplete  NotificationType = "review_complete"
	NotifyDeployComplete  NotificationType = "deploy_complete"
	NotifyAgentError      NotificationType = "agent_error"
)

// Severity defines notification severity (AC1)
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)

// Notification represents a notification (AC3)
type Notification struct {
	Type      NotificationType  `json:"type"`
	Severity  Severity          `json:"severity"`
	Message   string            `json:"message"`
	FeatureID string            `json:"feature_id,omitempty"`
	WSID      string            `json:"ws_id,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Variables map[string]string `json:"variables,omitempty"`
}

// Channel is the interface for notification channels (AC2)
type Channel interface {
	Name() string
	Send(n *Notification) error
	IsEnabled() bool
}

// Gateway manages notification delivery (AC4, AC5)
type Gateway struct {
	channels  []Channel
	history   []*Notification
	historyMu sync.RWMutex
	rateLimit *RateLimiter
}

// RateLimiter controls notification frequency (AC4)
type RateLimiter struct {
	maxPerMinute int
	counts       map[string]int
	windowStart  time.Time
	mu           sync.Mutex
}

// NewRateLimiter creates a rate limiter
func NewRateLimiter(maxPerMinute int) *RateLimiter {
	return &RateLimiter{
		maxPerMinute: maxPerMinute,
		counts:       make(map[string]int),
		windowStart:  time.Now(),
	}
}

// Allow checks if a notification is allowed
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Reset window if minute has passed
	if time.Since(r.windowStart) > time.Minute {
		r.counts = make(map[string]int)
		r.windowStart = time.Now()
	}

	if r.counts[key] >= r.maxPerMinute {
		return false
	}

	r.counts[key]++
	return true
}

// NewGateway creates a new notification gateway
func NewGateway() *Gateway {
	return &Gateway{
		channels:  []Channel{},
		history:   []*Notification{},
		rateLimit: NewRateLimiter(10), // Default: 10 per minute
	}
}

// NewGatewayWithRateLimit creates a gateway with custom rate limit
func NewGatewayWithRateLimit(maxPerMinute, debounceSeconds int) *Gateway {
	return &Gateway{
		channels:  []Channel{},
		history:   []*Notification{},
		rateLimit: NewRateLimiter(maxPerMinute),
	}
}

// AddChannel adds a notification channel
func (g *Gateway) AddChannel(channel Channel) {
	g.channels = append(g.channels, channel)
}

// Send delivers a notification to all enabled channels (AC2)
func (g *Gateway) Send(n *Notification) error {
	// Set timestamp if not set
	if n.Timestamp.IsZero() {
		n.Timestamp = time.Now()
	}

	// Store in history
	g.historyMu.Lock()
	g.history = append(g.history, n)
	g.historyMu.Unlock()

	// Log the notification
	safetylog.Info("notification: [%s] %s - %s", n.Severity, n.Type, n.Message)

	// Send to all enabled channels
	for _, channel := range g.channels {
		if channel.IsEnabled() {
			if err := channel.Send(n); err != nil {
				safetylog.Warn("notification: channel %s failed: %v", channel.Name(), err)
			}
		}
	}

	return nil
}

// AllowSend checks rate limit for a feature
func (g *Gateway) AllowSend(featureID string) bool {
	return g.rateLimit.Allow(featureID)
}

// AllowSendCritical bypasses rate limit for critical notifications
func (g *Gateway) AllowSendCritical() bool {
	return true // Critical notifications always allowed
}

// GetHistory returns notification history (AC5)
func (g *Gateway) GetHistory() []*Notification {
	g.historyMu.RLock()
	defer g.historyMu.RUnlock()

	result := make([]*Notification, len(g.history))
	copy(result, g.history)
	return result
}

// RenderTemplate renders a notification template with variables (AC3)
func RenderTemplate(template string, vars map[string]string) string {
	// Simple variable substitution
	result := template
	for k, v := range vars {
		result = replaceAll(result, "{"+k+"}", v)
	}
	return result
}

func replaceAll(s, old, new string) string {
	// Simple implementation
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
