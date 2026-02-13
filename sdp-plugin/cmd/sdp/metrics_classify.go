package main

import (
	"fmt"
	"os"

	metrics "github.com/fall-out-bug/sdp/internal/metrics"
	"github.com/spf13/cobra"
)

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
					metrics.FailureWrongLogic:       true,
					metrics.FailureMissingEdgeCase:  true,
					metrics.FailureHallucinatedAPI:  true,
					metrics.FailureTypeError:        true,
					metrics.FailureTestPassingWrong: true,
					metrics.FailureCompilationError: true,
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

// autoClassifyFromLog reads evidence log and classifies all failures.
func autoClassifyFromLog(taxonomy *metrics.Taxonomy, logPath string) (int, error) {
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
