package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// metricsClassifyCmd implements "sdp metrics classify" (AC3, AC5).
func metricsClassifyCmd() *cobra.Command {
	var eventID string
	var failureType string
	var notes string
	var autoClassify bool

	cmd := &cobra.Command{
		Use:   "classify",
		Short: "Classify AI failures by type",
		Long:  `Auto-classify verification failures by failure type using heuristic patterns.`,
		Example: `  # Auto-classify all unclassified failures from evidence log
  sdp metrics classify

  # Override classification for specific event
  sdp metrics classify --id=evt-123 --type=wrong_logic --notes="Manual correction"

  # Classify single event manually
  sdp metrics classify --id=evt-123 --type=type_error --notes="Missing import"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			taxonomyPath := ".sdp/metrics/taxonomy.json"
			taxonomy := metrics.NewTaxonomy(taxonomyPath)

			// Load existing taxonomy
			if err := taxonomy.Load(); err != nil {
				return fmt.Errorf("load taxonomy: %w", err)
			}

			// Manual override mode
			if eventID != "" && failureType != "" {
				// Validate failure type
				validTypes := map[string]bool{
					metrics.FailureWrongLogic:        true,
					metrics.FailureMissingEdgeCase:  true,
					metrics.FailureHallucinatedAPI:  true,
					metrics.FailureTypeError:        true,
					metrics.FailureTestPassingWrong: true,
					metrics.FailureCompilationError:  true,
					metrics.FailureImportError:      true,
				}
				if !validTypes[failureType] {
					return fmt.Errorf("invalid failure type: %s\nValid types: %v", failureType, getValidTypes())
				}

				taxonomy.SetClassification(eventID, failureType, notes)
				if err := taxonomy.Save(); err != nil {
					return fmt.Errorf("save taxonomy: %w", err)
				}

				fmt.Printf("✓ Classification updated: %s → %s\n", eventID, failureType)
				if notes != "" {
					fmt.Printf("  Notes: %s\n", notes)
				}
				return nil
			}

			// Auto-classify mode - read evidence log and classify failures
			logPath := ".sdp/log/events.jsonl"
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				return fmt.Errorf("evidence log not found: %s\nRun evidence collection first with 'sdp log show'", logPath)
			}

			classified, err := autoClassifyFromLog(taxonomy, logPath)
			if err != nil {
				return fmt.Errorf("auto-classify: %w", err)
			}

			if err := taxonomy.Save(); err != nil {
				return fmt.Errorf("save taxonomy: %w", err)
			}

			fmt.Printf("✓ Classified %d failures\n", classified)
			fmt.Printf("  Taxonomy saved to: %s\n", taxonomyPath)

			return nil
		},
	}

	cmd.Flags().StringVar(&eventID, "id", "", "Event ID to classify")
	cmd.Flags().StringVar(&failureType, "type", "", "Failure type (wrong_logic, missing_edge_case, etc)")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes for manual classification")
	cmd.Flags().BoolVar(&autoClassify, "auto", false, "Auto-classify all failures (default)")

	return cmd
}

// getValidTypes returns list of valid failure types.
func getValidTypes() []string {
	return []string{
		"wrong_logic",
		"missing_edge_case",
		"hallucinated_api",
		"type_error",
		"test_passing_but_wrong",
		"compilation_error",
		"import_error",
	}
}

// autoClassifyFromLog reads evidence log and classifies all failures.
func autoClassifyFromLog(taxonomy *metrics.Taxonomy, logPath string) (int, error) {
	// Import evidence reader to read events
	// For now, implement simple JSONL reading
	// TODO: Use internal/evidence/reader in future iteration

	events, err := readEventsJSONL(logPath)
	if err != nil {
		return 0, fmt.Errorf("read events: %w", err)
	}

	classified := 0
	for _, ev := range events {
		// Skip non-verification events
		if ev.Type != "verification" {
			continue
		}

		// Skip already classified
		if _, exists := taxonomy.GetClassification(ev.ID); exists {
			continue
		}

		// Extract output from data
		output := ""
		if outputStr, ok := ev.Data["output"].(string); ok {
			output = outputStr
		} else {
			// Build output from data fields
			if passed, ok := ev.Data["passed"].(bool); ok && !passed {
				if errMsg, ok := ev.Data["error"].(string); ok {
					output = errMsg
				} else {
					output = "verification failed"
				}
			}
		}

		if output == "" {
			continue
		}

		// Classify
		taxonomy.ClassifyFromOutput(ev.ID, ev.WSID, "", "", output)
		classified++
	}

	return classified, nil
}

// evidenceEvent represents a simplified evidence event for classification.
type evidenceEvent struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	WSID    string                 `json:"ws_id"`
	Data map[string]interface{} `json:"data"`
}

// readEventsJSONL reads events from JSONL file.
func readEventsJSONL(path string) ([]evidenceEvent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []evidenceEvent{}, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var events []evidenceEvent
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev evidenceEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue // Skip invalid lines
		}
		if ev.Data == nil {
			ev.Data = make(map[string]interface{})
		}
		events = append(events, ev)
	}

	return events, nil
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
