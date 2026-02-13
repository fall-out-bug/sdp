package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/telemetry"
	"github.com/spf13/cobra"
)

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

var telemetryUploadCmd = &cobra.Command{
	Use:   "upload --format json|archive",
	Short: "Package telemetry data for sharing",
	Long: `Package telemetry data for voluntary sharing.

This command packages your telemetry data into a file that you can:
  - Attach to a GitHub Issue
  - Send via email
  - Share for debugging

🔒 Privacy: Review the packaged file before sharing to ensure no sensitive data.

Formats:
  json    - Structured JSON with metadata (default)
  archive - tar.gz archive with raw telemetry.jsonl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format") //nolint:errcheck // String flag never errors

		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config dir: %w", err)
		}

		telemetryFile := filepath.Join(configDir, "sdp", "telemetry.jsonl")

		// Generate filename with timestamp
		timestamp := time.Now().Format("2006-01-02")
		var outputPath string
		if format == "archive" {
			outputPath = fmt.Sprintf("telemetry_upload_%s.tar.gz", timestamp)
		} else {
			outputPath = fmt.Sprintf("telemetry_upload_%s.json", timestamp)
		}

		// Package telemetry data
		result, err := telemetry.PackForUpload(telemetryFile, outputPath, format)
		if err != nil {
			return fmt.Errorf("failed to package telemetry: %w", err)
		}

		// Print summary
		fmt.Printf("✓ Collected %d events\n", result.EventCount)
		fmt.Printf("✓ Packaged into: %s\n", result.File)
		fmt.Printf("  Size: %.2f KB\n", float64(result.Size)/1024)

		// Privacy reminder
		fmt.Println("\n🔒 Privacy Reminder:")
		fmt.Println("  Review the file before sharing to ensure no sensitive data.")
		fmt.Println("\n  You can now:")
		fmt.Println("  - Attach to GitHub Issue")
		fmt.Println("  - Send via email")
		fmt.Println("  - Share for debugging")

		return nil
	},
}

func init() {
	telemetryUploadCmd.Flags().String("format", "json", "Output format: json or archive")
}
