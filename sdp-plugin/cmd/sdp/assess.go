package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/assess"
	"github.com/fall-out-bug/sdp/internal/telemetry"
	"github.com/spf13/cobra"
)

func assessCmd() *cobra.Command {
	var outputPath string
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "assess [project-path]",
		Short: "Assess project without making changes",
		Long: `Perform a read-only scan of the repository to detect:
  - Programming language
  - Frameworks and libraries
  - Project structure
  - Testing setup
  - CI/CD configuration
  - Monorepo patterns

Outputs recommendations to stdout only. No files are created.`,
		Example: `  # Assess current directory
  sdp assess

  # Assess specific project
  sdp assess /path/to/project

  # Output JSON
  sdp assess --json`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()

			// Determine project path
			projectPath := "."
			if len(args) > 0 {
				projectPath = args[0]
			}

			// Convert to absolute path
			absPath, err := filepath.Abs(projectPath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Check if path exists
			if _, err := os.Stat(absPath); os.IsNotExist(err) {
				return fmt.Errorf("path does not exist: %s", absPath)
			}

			// Run assessment (read-only, no clean-state check needed)
			result, err := assess.Assess(absPath)
			if err != nil {
				return fmt.Errorf("assessment failed: %w", err)
			}

			// Output results
			var outputWriter = os.Stdout
			if outputPath != "" {
				f, err := os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer f.Close()
				outputWriter = f
			}

			if jsonOutput {
				if err := printAssessmentJSONTo(result, outputWriter); err != nil {
					return err
				}
			} else {
				if err := printAssessmentTo(result, absPath, outputWriter); err != nil {
					return err
				}
			}

			// Initialize telemetry collector (after assessment, to avoid creating files in assessed repos)
			uxMetrics, err := telemetry.NewUXMetricsCollector("")
			if err != nil {
				// Don't fail the command if telemetry fails
				fmt.Fprintf(os.Stderr, "Warning: failed to initialize telemetry: %v\n", err)
			}

			// Record telemetry
			if uxMetrics != nil {
				duration := time.Since(startTime)
				projectType := result.Language
				if projectType == "" {
					projectType = "unknown"
				}

				if err := uxMetrics.RecordAssessComplete(projectType, duration); err != nil {
					// Don't fail the command if telemetry fails
					fmt.Fprintf(os.Stderr, "Warning: failed to record telemetry: %v\n", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Write output to file")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output JSON format")

	return cmd
}

func printAssessmentTo(result *assess.Assessment, projectPath string, w io.Writer) error {
	fmt.Fprintln(w, "SDP Project Assessment")
	fmt.Fprintln(w, "=====================")
	fmt.Fprintf(w, "Project: %s\n\n", projectPath)

	// Language
	fmt.Fprintf(w, "Language: %s\n", result.Language)

	// Frameworks
	if len(result.Framework) > 0 {
		fmt.Fprintf(w, "Frameworks: %s\n", strings.Join(result.Framework, ", "))
	} else {
		fmt.Fprintln(w, "Frameworks: None detected")
	}

	// Structure
	if len(result.Structure) > 0 {
		fmt.Fprintf(w, "Structure: %s\n", strings.Join(result.Structure, ", "))
	} else {
		fmt.Fprintln(w, "Structure: standard")
	}

	// Flags
	fmt.Fprintf(w, "Monorepo: %v\n", result.IsMonorepo)
	fmt.Fprintf(w, "Has Tests: %v\n", result.HasTests)
	fmt.Fprintf(w, "Has CI: %v\n", result.HasCI)

	// Recommendations
	fmt.Fprintln(w, "\nRecommendations")
	fmt.Fprintln(w, "--------------")

	if len(result.Recommendations) == 0 {
		fmt.Fprintln(w, "No recommendations - project looks good!")
	} else {
		for _, rec := range result.Recommendations {
			priorityIcon := "ℹ"
			if rec.Priority == "high" {
				priorityIcon = "⚠"
			} else if rec.Priority == "medium" {
				priorityIcon = "→"
			}

			fmt.Fprintf(w, "%s [%s] %s\n", priorityIcon, rec.Category, rec.Title)
			fmt.Fprintf(w, "    %s\n\n", rec.Message)
		}
	}

	return nil
}

func printAssessmentJSONTo(result *assess.Assessment, w io.Writer) error {
	// Define a JSON-serializable structure
	type JSONRecommendation struct {
		Category string `json:"category"`
		Title    string `json:"title"`
		Message  string `json:"message"`
		Priority string `json:"priority"`
	}

	type JSONAssessment struct {
		Language        string               `json:"language"`
		Frameworks      []string             `json:"frameworks"`
		Structure       []string             `json:"structure"`
		IsMonorepo      bool                 `json:"is_monorepo"`
		HasTests        bool                 `json:"has_tests"`
		HasCI           bool                 `json:"has_ci"`
		Recommendations []JSONRecommendation `json:"recommendations"`
	}

	// Convert recommendations to JSON format
	jsonRecs := make([]JSONRecommendation, len(result.Recommendations))
	for i, rec := range result.Recommendations {
		jsonRecs[i] = JSONRecommendation{
			Category: rec.Category,
			Title:    rec.Title,
			Message:  rec.Message,
			Priority: rec.Priority,
		}
	}

	jsonResult := JSONAssessment{
		Language:        result.Language,
		Frameworks:      result.Framework,
		Structure:       result.Structure,
		IsMonorepo:      result.IsMonorepo,
		HasTests:        result.HasTests,
		HasCI:           result.HasCI,
		Recommendations: jsonRecs,
	}

	// Marshal to JSON with proper escaping
	data, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal assessment to JSON: %w", err)
	}

	fmt.Fprintln(w, string(data))
	return nil
}
