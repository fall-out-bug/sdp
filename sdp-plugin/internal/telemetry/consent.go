package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConsentConfig represents the user's telemetry consent choice
type ConsentConfig struct {
	Enabled    bool   `json:"enabled"`
	AskedAt    string `json:"asked_at,omitempty"`    // When user was asked
	AnsweredAt string `json:"answered_at,omitempty"` // When user answered
	Version    string `json:"version,omitempty"`     // Privacy policy version
}

// CheckConsent checks if user has granted telemetry consent
// Returns (granted, error)
func CheckConsent(configPath string) (bool, error) {
	// If config doesn't exist, consent not granted
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false, nil
	}

	// Read config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false, fmt.Errorf("failed to read config: %w", err)
	}

	var config ConsentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return false, fmt.Errorf("failed to parse config: %w", err)
	}

	return config.Enabled, nil
}

// GrantConsent saves user's consent choice
// enabled=true means user consented to telemetry
// enabled=false means user declined
func GrantConsent(configPath string, enabled bool) error {
	// Create directory if needed
	dir := configPath[:len(configPath)-len("telemetry.json")]
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	// Load existing config if any
	var config ConsentConfig
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse existing config: %w", err)
		}
	}

	// Update config
	config.Enabled = enabled

	// Save with secure permissions (owner read/write only)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// AskForConsent prompts user for telemetry consent (interactive)
// This should be called on first run
func AskForConsent() (bool, error) {
	// For now, return false (consent not granted)
	// In real implementation, this would show an interactive prompt
	// and wait for user input

	// TODO: Implement interactive prompt:
	//
	// SDP collects anonymized usage telemetry to improve quality.
	//
	// What's collected:
	//   - Command usage (@build, @review, etc.)
	//   - Execution duration
	//   - Success/failure rates
	//
	// What's NOT collected:
	//   - No PII (names, emails, usernames)
	//   - No code content
	//   - No file paths
	//   - Data stays local (never transmitted)
	//
	// Help improve SDP? (y/n):

	return false, nil
}

// IsFirstRun checks if this is the first run (no consent config exists)
func IsFirstRun(configPath string) bool {
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}
