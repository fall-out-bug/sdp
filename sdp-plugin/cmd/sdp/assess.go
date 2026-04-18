package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/assess"
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

			// Run assessment
			result, err := assess.Assess(absPath)
			if err != nil {
				return fmt.Errorf("assessment failed: %w", err)
			}

			// Output results
			if jsonOutput {
				return printAssessmentJSON(result)
			}
			return printAssessment(result, absPath)
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
	fmt.Printf("Frameworks: %s\n", fmt.Sprintf("[%s]", fmt.Sprintf("%s", result.Framework)))

	// Structure
	if len(result.Structure) > 0 {
		fmt.Printf("Structure: %s\n", fmt.Sprintf("[%s]", fmt.Sprintf("%s", result.Structure)))
	} else {
		fmt.Println("Structure: [standard]")
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
	// Simple JSON output (for now - could use json.Marshal later)
	fmt.Printf(`{
  "language": "%s",
  "frameworks": [%s],
  "structure": [%s],
  "is_monorepo": %v,
  "has_tests": %v,
  "has_ci": %v,
  "recommendations": [
`, result.Language, result.Framework, result.Structure, result.IsMonorepo, result.HasTests, result.HasCI)

	for i, rec := range result.Recommendations {
		comma := ","
		if i == len(result.Recommendations)-1 {
			comma = ""
		}
		fmt.Printf(`    {
      "category": "%s",
      "title": "%s",
      "message": "%s",
      "priority": "%s"
    }%s
`, rec.Category, rec.Title, rec.Message, rec.Priority, comma)
	}

	fmt.Println("  ]")
	fmt.Println("}")

	return nil
}
