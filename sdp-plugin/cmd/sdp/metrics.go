package main

import (
	"fmt"
	"os"
	"path/filepath"

	metrics "github.com/fall-out-bug/sdp/internal/metrics"
	"github.com/spf13/cobra"
)

// metricsCmd collects metrics from evidence log (AC1).
func metricsCmd() *cobra.Command {
	var output string
	var watermark bool

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Collect and report metrics from evidence log",
		Long: `Metrics collection and reporting for evidence events.

This command reads the evidence log (.sdp/log/events.jsonl) and computes:
- catch_rate: verification failures / total verifications
- iteration_count: red→green cycles per workstream
- model_pass_rate: pass rate per model ID
- acceptance_catch_rate: acceptance failures / total approvals

Output is written to .sdp/metrics/latest.json by default.`,
		Example: `  # Collect metrics from evidence log
  sdp metrics collect

  # Collect metrics with custom output path
  sdp metrics collect --output /path/to/metrics.json

  # Incremental collection (only process new events)
  sdp metrics collect --watermark`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires a subcommand (collect, classify, report)")
			}
			return nil
		},
	}

	// Add subcommands
	cmd.AddCommand(metricsCollectCmd())
	cmd.AddCommand(metricsClassifyCmd())
	cmd.AddCommand(metricsReportCmd())

	cmd.PersistentFlags().StringVar(&output, "output", "", "Output path for metrics")
	cmd.PersistentFlags().BoolVar(&watermark, "watermark", false, "Enable incremental collection using watermark")

	return cmd
}

// metricsCollectCmd implements "sdp metrics collect" (AC1).
func metricsCollectCmd() *cobra.Command {
	var outputPath string
	var enableWatermark bool

	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect metrics from evidence log",
		Long:  `Scan the evidence log and compute metrics (catch rate, iterations, model performance).`,
		Example: `  # Collect metrics to default location
  sdp metrics collect

  # Collect to custom path
  sdp metrics collect --output ./my-metrics.json

  # Enable incremental collection
  sdp metrics collect --watermark`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default output path: .sdp/metrics/latest.json
			if outputPath == "" {
				outputPath = ".sdp/metrics/latest.json"
			}

			// Evidence log path: .sdp/log/events.jsonl
			logPath := ".sdp/log/events.jsonl"

			// Check if evidence log exists
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				return fmt.Errorf("evidence log not found: %s\nRun evidence collection first with 'sdp log show'", logPath)
			}

			// Create collector
			collector := metrics.NewCollector(logPath, outputPath)

			// Set watermark path if enabled
			if enableWatermark {
				watermarkPath := filepath.Join(filepath.Dir(outputPath), ".watermark.json")
				collector.SetWatermarkPath(watermarkPath)
			}

			// Collect metrics
			collectedMetrics, err := collector.Collect()
			if err != nil {
				return fmt.Errorf("collect metrics: %w", err)
			}

			// Print summary
			fmt.Printf("✓ Metrics collected:\n")
			fmt.Printf("  Catch Rate: %.2f%%\n", collectedMetrics.CatchRate*100)
			fmt.Printf("  Total Verifications: %d\n", collectedMetrics.TotalVerifications)
			fmt.Printf("  Failed Verifications: %d\n", collectedMetrics.FailedVerifications)
			fmt.Printf("  Acceptance Catch Rate: %.2f%%\n", collectedMetrics.AcceptanceCatchRate*100)
			fmt.Printf("  Output: %s\n", outputPath)

			return nil
		},
	}

	cmd.Flags().StringVar(&outputPath, "output", "", "Output path for metrics JSON")
	cmd.Flags().BoolVar(&enableWatermark, "watermark", false, "Enable incremental collection using watermark")

	return cmd
}

// metricsClassifyCmd implements "sdp metrics classify" (placeholder for F061-02).
func metricsClassifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "classify",
		Short: "Classify AI failures by type (requires F061-02)",
		Long:  `Auto-classify verification failures by failure type. Requires F061-02 (AI Failure Taxonomy).`,
		Example: `  # Classify failures from evidence log
  sdp metrics classify

  # Override classification for specific event
  sdp metrics classify --id=evt-123 --type=wrong_logic`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("metrics classify requires F061-02 (AI Failure Taxonomy)\nStatus: https://github.com/fall-out-bug/sdp/issues/xxx")
		},
	}

	var eventID string
	var failureType string

	cmd.Flags().StringVar(&eventID, "id", "", "Event ID to classify")
	cmd.Flags().StringVar(&failureType, "type", "", "Failure type (wrong_logic, missing_edge_case, etc)")

	return cmd
}

// metricsReportCmd implements "sdp metrics report" (placeholder for F061-03).
func metricsReportCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate benchmark report (requires F061-03)",
		Long:  `Generate AI Code Quality Benchmark report from metrics. Requires F061-03 (Benchmark Report Generator).`,
		Example: `  # Generate markdown report
  sdp metrics report

  # Generate HTML report
  sdp metrics report --format=html

  # Generate JSON report
  sdp metrics report --format=json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("metrics report requires F061-03 (Benchmark Report Generator)\nStatus: https://github.com/fall-out-bug/sdp/issues/xxx")
		},
	}

	cmd.Flags().StringVar(&format, "format", "md", "Report format: md, html, json")

	return cmd
}

func init() {
	// Ensure metrics directory exists
	metricsDir := filepath.Join(".sdp", "metrics")
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to create metrics directory: %v\n", err)
	}
}
