package main

import (
	"fmt"
	"os"

	metrics "github.com/fall-out-bug/sdp/internal/metrics"
	"github.com/spf13/cobra"
)

// metricsReportCmd implements "sdp metrics report" (AC1-AC7).
func metricsReportCmd() *cobra.Command {
	var format string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate benchmark report",
		Long:  `Generate AI Code Quality Benchmark report from metrics and taxonomy data.`,
		Example: `  # Generate markdown report (default)
  sdp metrics report

  # Generate HTML report
  sdp metrics report --format=html

  # Generate JSON report
  sdp metrics report --format=json

  # Generate to custom path
  sdp metrics report --output ./my-report.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default paths
			metricsPath := ".sdp/metrics/latest.json"
			taxonomyPath := ".sdp/metrics/taxonomy.json"
			if outputPath == "" {
				outputPath = ".sdp/metrics/benchmark-{{QUARTER}}.md"
			}

			// Check if metrics exist
			if _, err := os.Stat(metricsPath); os.IsNotExist(err) {
				return fmt.Errorf("metrics not found: %s\nRun 'sdp metrics collect' first", metricsPath)
			}

			// Create reporter
			reporter := metrics.NewReporter(metricsPath, taxonomyPath)

			// Generate report based on format
			var report string
			var err error
			switch format {
			case "html":
				report, err = reporter.GenerateHTML()
			case "json":
				report, err = reporter.GenerateJSON()
			default: // md
				report, err = reporter.GenerateMarkdown()
			}

			if err != nil {
				return fmt.Errorf("generate report: %w", err)
			}

			// Write report
			if err := os.WriteFile(outputPath, []byte(report), 0644); err != nil {
				return fmt.Errorf("write report: %w", err)
			}

			fmt.Printf("✓ Benchmark report generated:\n")
			fmt.Printf("  Format: %s\n", format)
			fmt.Printf("  Output: %s\n", outputPath)

			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "md", "Report format: md, html, json")
	cmd.Flags().StringVar(&outputPath, "output", "", "Output path for report")

	return cmd
}
