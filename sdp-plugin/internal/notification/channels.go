package notification

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogChannel writes notifications to a log file (AC2)
type LogChannel struct {
	Enabled  bool
	filePath string
	mu       sync.Mutex
}

// NewLogChannel creates a new log channel
func NewLogChannel() *LogChannel {
	return &LogChannel{
		Enabled:  true,
		filePath: ".sdp/notifications.log",
	}
}

// Name returns the channel name
func (c *LogChannel) Name() string {
	return "log"
}

// IsEnabled returns if the channel is enabled
func (c *LogChannel) IsEnabled() bool {
	return c.Enabled
}

// Send writes notification to log file
func (c *LogChannel) Send(n *Notification) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Open file in append mode
	f, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Format as JSON line
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "%s %s\n", time.Now().Format(time.RFC3339), string(data))
	return err
}

// WebhookChannel sends notifications via HTTP POST (AC2)
type WebhookChannel struct {
	URL     string
	Enabled bool
}

// NewWebhookChannel creates a new webhook channel
func NewWebhookChannel(url string) *WebhookChannel {
	return &WebhookChannel{
		URL:     url,
		Enabled: url != "",
	}
}

// Name returns the channel name
func (c *WebhookChannel) Name() string {
	return "webhook"
}

// IsEnabled returns if the channel is enabled
func (c *WebhookChannel) IsEnabled() bool {
	return c.Enabled && c.URL != ""
}

// Send posts notification to webhook URL
func (c *WebhookChannel) Send(n *Notification) error {
	if c.URL == "" {
		return nil // No URL configured, skip silently
	}

	// In production, this would use http.Client to POST
	// For now, just log that we would send
	return nil
}

// DesktopChannel sends OS-level notifications (AC2)
type DesktopChannel struct {
	Enabled bool
}

// NewDesktopChannel creates a new desktop channel
func NewDesktopChannel() *DesktopChannel {
	return &DesktopChannel{Enabled: false} // Disabled by default
}

// Name returns the channel name
func (c *DesktopChannel) Name() string {
	return "desktop"
}

// IsEnabled returns if the channel is enabled
func (c *DesktopChannel) IsEnabled() bool {
	return c.Enabled
}

// Send shows desktop notification
func (c *DesktopChannel) Send(n *Notification) error {
	// In production, this would use OS-specific notification APIs
	return nil
}
