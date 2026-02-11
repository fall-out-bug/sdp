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
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 Telemetry Consent")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("SDP может собирать анонимную статистику использования")
	fmt.Println("для улучшения качества и надежности.")
	fmt.Println()
	fmt.Println("🔒 Что собирается:")
	fmt.Println("  • Команды (@build, @review, @oneshot)")
	fmt.Println("  • Время выполнения команд")
	fmt.Println("  • Успех/ошибки выполнения")
	fmt.Println()
	fmt.Println("❌ Что НЕ собирается:")
	fmt.Println("  • PII (имена, email, логины)")
	fmt.Println("  • Содержимое кода")
	fmt.Println("  • Пути к файлам")
	fmt.Println("  • Данные остаются локальными (не отправляются)")
	fmt.Println()
	fmt.Println("📜 Политика конфиденциальности: docs/PRIVACY.md")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Помочь улучшить SDP? (y/n): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			// Non-interactive environment (e.g., script)
			fmt.Println("\n(non-interactive mode: telemetry disabled)")
			return false, nil
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes", "да", "д":
			fmt.Println("\n✓ Спасибо! Ваш вклад помогает улучшить SDP.")
			fmt.Println("  Вы можете отключить в любой момент:")
			fmt.Println("  sdp telemetry disable")
			return true, nil

		case "n", "no", "нет", "н":
			fmt.Println("\n✓ Телеметрия отключена.")
			fmt.Println("  Вы можете включить позже:")
			fmt.Println("  sdp telemetry enable")
			return false, nil

		default:
			fmt.Println("Пожалуйста, введите 'y' или 'n'")
		}
	}
}

// IsFirstRun checks if this is the first run (no consent config exists)
func IsFirstRun(configPath string) bool {
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}
