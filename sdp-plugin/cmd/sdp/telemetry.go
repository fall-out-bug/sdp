package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ai-masters/sdp/internal/telemetry"
)

var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Manage telemetry collection",
	Long: `Manage telemetry collection for SDP.

Telemetry tracks anonymized usage metrics to help improve SDP:
  - Command invocations
  - Execution duration
  - Success/failure rates

🔒 Privacy Policy:
  - No PII (names, emails, usernames) collected
  - No data transmitted remotely (stored locally)
  - Opt-out available: sdp telemetry disable
  - Auto-cleanup after 90 days
  - See docs/PRIVACY.md for details

All data is stored locally in ~/.sdp/telemetry.jsonl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default to status if no subcommand
		return telemetryStatusCmd.RunE(cmd, args)
	},
}

var telemetryStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show telemetry status",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		telemetryFile := filepath.Join(configDir, "sdp", "telemetry.jsonl")
		collector, err := telemetry.NewCollector(telemetryFile, true)
		if err != nil {
			return fmt.Errorf("failed to create collector: %w", err)
		}

		// Check if telemetry is disabled via config file
		configPath := filepath.Join(configDir, "sdp", "telemetry.json")
		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}

			var config map[string]bool
			if err := json.Unmarshal(data, &config); err == nil {
				if disabled, ok := config["disabled"]; ok && disabled {
					collector.Disable()
				}
			}
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
		}

		return nil
	},
}

var telemetryExportCmd = &cobra.Command{
	Use:   "export [format]",
	Short: "Export telemetry data",
	Long: `Export telemetry data to JSON or CSV.

If no format is specified, defaults to JSON.
The export file is saved to the current directory.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format := "json"
		if len(args) > 0 {
			format = args[0]
		}

		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		telemetryFile := filepath.Join(configDir, "sdp", "telemetry.jsonl")
		collector, err := telemetry.NewCollector(telemetryFile, true)
		if err != nil {
			return fmt.Errorf("failed to create collector: %w", err)
		}

		// Determine export filename
		exportPath := fmt.Sprintf("telemetry_export.%s", format)

		// Export based on format
		switch format {
		case "json":
			if err := collector.ExportJSON(exportPath); err != nil {
				return fmt.Errorf("failed to export JSON: %w", err)
			}
		case "csv":
			if err := collector.ExportCSV(exportPath); err != nil {
				return fmt.Errorf("failed to export CSV: %w", err)
			}
		default:
			return fmt.Errorf("unsupported format: %s (use json or csv)", format)
		}

		fmt.Printf("Exported telemetry to %s\n", exportPath)
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

		// Create config directory
		configPath := filepath.Join(configDir, "sdp")
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Write config file
		configFile := filepath.Join(configPath, "telemetry.json")
		config := map[string]bool{"disabled": true}
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		// Write config file with secure permissions (owner read/write only)
		if err := os.WriteFile(configFile, data, 0600); err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		fmt.Println("Telemetry disabled")
		return nil
	},
}

var telemetryEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable telemetry collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		// Remove config file if it exists
		configFile := filepath.Join(configDir, "sdp", "telemetry.json")
		if err := os.Remove(configFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove config: %w", err)
		}

		fmt.Println("Telemetry enabled")
		return nil
	},
}

var telemetryAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze telemetry data for insights",
	Long: `Analyze telemetry data to generate insights.

Calculates:
  - Success rate by command
  - Average execution time by command
  - Top error categories
  - Overall usage statistics`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		telemetryFile := filepath.Join(configDir, "sdp", "telemetry.jsonl")
		analyzer, err := telemetry.NewAnalyzer(telemetryFile)
		if err != nil {
			return fmt.Errorf("failed to create analyzer: %w", err)
		}

		// Generate report
		report, err := analyzer.GenerateReport(nil, nil)
		if err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}

		// Print report
		fmt.Println("\n📊 Telemetry Analysis Report")
		fmt.Println("==========================")
		fmt.Printf("\nTotal Events: %d\n", report.TotalEvents)

		if len(report.CommandStats) > 0 {
			fmt.Println("\n📈 Command Statistics:")
			fmt.Println("----------------------")
			for cmd, stats := range report.CommandStats {
				fmt.Printf("\n  %s:\n", cmd)
				fmt.Printf("    Total Runs: %d\n", stats.TotalRuns)
				fmt.Printf("    Success Rate: %.1f%%\n", stats.SuccessRate*100)
				fmt.Printf("    Avg Duration: %dms\n", stats.AvgDuration)
			}
		}

		if len(report.TopErrors) > 0 {
			fmt.Println("\n❌ Top Errors:")
			fmt.Println("-------------")
			for i, err := range report.TopErrors {
				fmt.Printf("  %d. %s (%d occurrences)\n", i+1, err.Message, err.Count)
			}
		}

		fmt.Println()
		return nil
	},
}

func init() {
	telemetryCmd.AddCommand(telemetryStatusCmd)
	telemetryCmd.AddCommand(telemetryExportCmd)
	telemetryCmd.AddCommand(telemetryDisableCmd)
	telemetryCmd.AddCommand(telemetryEnableCmd)
	telemetryCmd.AddCommand(telemetryAnalyzeCmd)
}
