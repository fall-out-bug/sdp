package main

import (
	"encoding/json"
	"fmt"
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

			// Initialize telemetry collector
			uxMetrics, err := telemetry.NewUXMetricsCollector("")
			if err != nil {
				// Don't fail the command if telemetry fails
				fmt.Fprintf(os.Stderr, "Warning: failed to initialize telemetry: %v\n", err)
			}

			// Run assessment
			result, err := assess.Assess(absPath)
			if err != nil {
				return fmt.Errorf("assessment failed: %w", err)
			}

			// Output results
			if jsonOutput {
				if err := printAssessmentJSON(result); err != nil {
					return err
				}
			} else {
				if err := printAssessment(result, absPath); err != nil {
					return err
				}
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

func printAssessment(result *assess.Assessment, projectPath string) error {
	fmt.Println("SDP Project Assessment")
	fmt.Println("=====================")
	fmt.Printf("Project: %s\n\n", projectPath)

	// Language
	fmt.Printf("Language: %s\n", result.Language)

	// Frameworks
	if len(result.Framework) > 0 {
		fmt.Printf("Frameworks: %s\n", strings.Join(result.Framework, ", "))
	} else {
		fmt.Println("Frameworks: None detected")
	}

	// Structure
	if len(result.Structure) > 0 {
		fmt.Printf("Structure: %s\n", strings.Join(result.Structure, ", "))
	} else {
		fmt.Println("Structure: standard")
	}

	// Flags
	fmt.Printf("Monorepo: %v\n", result.IsMonorepo)
	fmt.Printf("Has Tests: %v\n", result.HasTests)
	fmt.Printf("Has CI: %v\n", result.HasCI)

	// Recommendations
	fmt.Println("\nRecommendations")
	fmt.Println("--------------")

	if len(result.Recommendations) == 0 {
		fmt.Println("No recommendations - project looks good!")
	} else {
		for _, rec := range result.Recommendations {
			priorityIcon := "ℹ"
			if rec.Priority == "high" {
				priorityIcon = "⚠"
			} else if rec.Priority == "medium" {
				priorityIcon = "→"
			}

			fmt.Printf("%s [%s] %s\n", priorityIcon, rec.Category, rec.Title)
			fmt.Printf("    %s\n\n", rec.Message)
		}
	}

	return nil
}

func printAssessmentJSON(result *assess.Assessment) error {
	// Define a JSON-serializable structure
	type JSONRecommendation struct {
		Category string `json:"category"`
		Title    string `json:"title"`
		Message  string `json:"message"`
		Priority string `json:"priority"`
	}

	type JSONAssessment struct {
		Language       string              `json:"language"`
		Frameworks     []string            `json:"frameworks"`
		Structure      []string            `json:"structure"`
		IsMonorepo     bool                `json:"is_monorepo"`
		HasTests       bool                `json:"has_tests"`
		HasCI          bool                `json:"has_ci"`
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
		Language:       result.Language,
		Frameworks:     result.Framework,
		Structure:      result.Structure,
		IsMonorepo:     result.IsMonorepo,
		HasTests:       result.HasTests,
		HasCI:          result.HasCI,
		Recommendations: jsonRecs,
	}

	// Marshal to JSON with proper escaping
	data, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal assessment to JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}
