package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/telemetry"
	"github.com/spf13/cobra"
)

var telemetryStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show telemetry status",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		configPath := filepath.Join(configDir, "sdp", "telemetry.json")
		telemetryFile := filepath.Join(configDir, "sdp", "telemetry.jsonl")

		// Check if user has granted consent
		enabled, err := telemetry.CheckConsent(configPath)
		if err != nil {
			return fmt.Errorf("failed to check consent: %w", err)
		}

		collector, err := telemetry.NewCollector(telemetryFile, enabled)
		if err != nil {
			return fmt.Errorf("failed to create collector: %w", err)
		}

		status := collector.Status()

		fmt.Println("Telemetry Status:")
		fmt.Printf("  Enabled: %v\n", status.Enabled)
		fmt.Printf("  Events: %d\n", status.EventCount)
		fmt.Printf("  File: %s\n", status.FilePath)

		if status.Enabled {
			fmt.Println("\n🔒 Privacy:")
			fmt.Println("  - No PII collected")
			fmt.Println("  - Data stays local")
			fmt.Println("  - Auto-cleanup after 90 days")
			fmt.Println("  - See: docs/PRIVACY.md")
			fmt.Println("\n  To disable: sdp telemetry disable")
		} else {
			fmt.Println("\n📊 Opt-in:")
			fmt.Println("  - Telemetry is currently disabled")
			fmt.Println("  - To help improve SDP: sdp telemetry enable")
			fmt.Println("  - See: docs/PRIVACY.md")
		}

		return nil
	},
}

var telemetryConsentCmd = &cobra.Command{
	Use:   "consent",
	Short: "Manage telemetry consent",
	Long: `Manage your telemetry consent preference.

Telemetry is opt-in by default. Use this command to:
  - Grant consent: sdp telemetry consent grant
  - Revoke consent: sdp telemetry consent revoke
  - Check status: sdp telemetry status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Telemetry Consent:")
		fmt.Println("==================")
		fmt.Println()
		fmt.Println("SDP collects anonymized usage telemetry to improve quality.")
		fmt.Println()
		fmt.Println("🔒 What's collected:")
		fmt.Println("  - Command usage (@build, @review, etc.)")
		fmt.Println("  - Execution duration")
		fmt.Println("  - Success/failure rates")
		fmt.Println()
		fmt.Println("❌ What's NOT collected:")
		fmt.Println("  - No PII (names, emails, usernames)")
		fmt.Println("  - No code content")
		fmt.Println("  - No file paths")
		fmt.Println("  - Data stays local (never transmitted)")
		fmt.Println()
		fmt.Println("To grant consent:  sdp telemetry enable")
		fmt.Println("To revoke consent: sdp telemetry disable")
		fmt.Println()
		fmt.Println("See: docs/PRIVACY.md for full privacy policy")
		return nil
	},
}

var telemetryDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable telemetry collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		configFile := filepath.Join(configDir, "sdp", "telemetry.json")

		// Revoke consent (disable telemetry)
		if err := telemetry.GrantConsent(configFile, false); err != nil {
			return fmt.Errorf("failed to disable telemetry: %w", err)
		}

		fmt.Println("✓ Telemetry disabled")
		fmt.Println("  Your data remains local and will not be collected.")
		return nil
	},
}

var telemetryEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable telemetry collection (opt-in)",
	Long: `Enable telemetry collection to help improve SDP.

This is an opt-in choice. SDP will collect:
  - Command usage patterns
  - Execution duration
  - Success/failure rates

NO PII is collected. Data stays local. See docs/PRIVACY.md for details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		configFile := filepath.Join(configDir, "sdp", "telemetry.json")

		// Grant consent (enable telemetry)
		if err := telemetry.GrantConsent(configFile, true); err != nil {
			return fmt.Errorf("failed to enable telemetry: %w", err)
		}

		fmt.Println("✓ Telemetry enabled")
		fmt.Println("  Thank you for helping improve SDP!")
		fmt.Println("  To disable: sdp telemetry disable")
		return nil
	},
}
