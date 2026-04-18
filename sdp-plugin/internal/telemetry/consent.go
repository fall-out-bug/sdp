package telemetry

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
// Returns true if user consented, false otherwise
func AskForConsent() (bool, error) {
	// Write consent banner to stderr to avoid corrupting structured stdout (e.g., --json)
	w := os.Stderr
	fmt.Fprintln(w, "\n"+strings.Repeat("=", 60))
	fmt.Fprintln(w, "📊 Telemetry Consent")
	fmt.Fprintln(w, strings.Repeat("=", 60))
	fmt.Fprintln(w)
	fmt.Fprintln(w, "SDP can collect anonymous usage statistics")
	fmt.Fprintln(w, "to improve quality and reliability.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "🔒 What is collected:")
	fmt.Fprintln(w, "  • Commands (@build, @review, @oneshot)")
	fmt.Fprintln(w, "  • Command execution duration")
	fmt.Fprintln(w, "  • Success/failure of execution")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "❌ What is NOT collected:")
	fmt.Fprintln(w, "  • PII (names, email, usernames)")
	fmt.Fprintln(w, "  • Code content")
	fmt.Fprintln(w, "  • File paths")
	fmt.Fprintln(w, "  • Data stays local (not transmitted)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "📜 Privacy policy: docs/PRIVACY.md")
	fmt.Fprintln(w)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(w, "Help improve SDP? (y/n): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			// Non-interactive environment (e.g., script)
			fmt.Fprintln(w, "\n(non-interactive mode: telemetry disabled)")
			return false, nil
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			fmt.Fprintln(w, "\n✓ Thank you! Your contribution helps improve SDP.")
			fmt.Fprintln(w, "  You can disable anytime with:")
			fmt.Fprintln(w, "  sdp telemetry disable")
			return true, nil

		case "n", "no":
			fmt.Fprintln(w, "\n✓ Telemetry disabled.")
			fmt.Fprintln(w, "  You can enable later with:")
			fmt.Fprintln(w, "  sdp telemetry enable")
			return false, nil

		default:
			fmt.Fprintln(w, "Please enter 'y' or 'n'")
		}
	}
}

// IsFirstRun checks if this is the first run (no consent config exists)
func IsFirstRun(configPath string) bool {
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}
