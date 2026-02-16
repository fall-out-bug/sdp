package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/nextstep"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

// nextCmd returns the next command
func nextCmd() *cobra.Command {
	var outputJSON bool
	var showAlternatives bool

	cmd := &cobra.Command{
		Use:   "next",
		Short: "Get next-step recommendation for current project",
		Long: `Analyze current project state and recommend the next action.

The recommendation engine considers:
  - Workstream status and dependencies
  - Current execution state (in-progress, blocked, failed)
  - Git repository state
  - SDP configuration

Output includes:
  - Recommended command with confidence level
  - Reason for the recommendation
  - Alternative actions when available`,
		Example: `  # Get next step recommendation
  sdp next

  # Get JSON output for scripting
  sdp next --json

  # Show all alternatives
  sdp next --alternatives`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get project root
			projectRoot, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			// Collect state
			collector := nextstep.NewStateCollector(projectRoot)
			state, err := collector.Collect()
			if err != nil {
				return fmt.Errorf("failed to collect project state: %w", err)
			}

			// Get recommendation
			resolver := nextstep.NewResolver()
			rec, err := resolver.Recommend(state)
			if err != nil {
				return fmt.Errorf("failed to generate recommendation: %w", err)
			}

			// Output
			if outputJSON {
				return outputRecommendationJSON(rec)
			}

			return outputRecommendationHuman(rec, showAlternatives)
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&showAlternatives, "alternatives", false, "Show alternative recommendations")

	return cmd
}

// outputRecommendationJSON outputs the recommendation as JSON.
func outputRecommendationJSON(rec *nextstep.Recommendation) error {
	data, err := rec.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize recommendation: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// outputRecommendationHuman outputs the recommendation in human-readable format.
func outputRecommendationHuman(rec *nextstep.Recommendation, showAlternatives bool) error {
	// Header
	fmt.Println()
	fmt.Println(ui.BoldText("Next Step Recommendation"))
	fmt.Println()

	// Confidence indicator
	confidence := rec.Confidence
	confidenceStr := fmt.Sprintf("%.0f%%", confidence*100)
	if confidence >= 0.8 {
		confidenceStr = ui.Success(confidenceStr)
	} else if confidence >= 0.5 {
		confidenceStr = ui.Warning(confidenceStr)
	} else {
		confidenceStr = ui.Dim(confidenceStr)
	}

	// Category badge
	categoryStr := string(rec.Category)
	switch rec.Category {
	case nextstep.CategoryExecution:
		categoryStr = ui.Info(categoryStr)
	case nextstep.CategoryRecovery:
		categoryStr = ui.Error(categoryStr)
	case nextstep.CategoryPlanning:
		categoryStr = ui.Info(categoryStr)
	case nextstep.CategorySetup:
		categoryStr = ui.Warning(categoryStr)
	default:
		categoryStr = ui.Dim(categoryStr)
	}

	// Main recommendation
	fmt.Printf("  %s %s\n", ui.BoldText("Command:"), ui.Success(rec.Command))
	fmt.Printf("  %s %s\n", ui.BoldText("Reason:"), rec.Reason)
	fmt.Printf("  %s %s  %s %s\n", ui.BoldText("Confidence:"), confidenceStr, ui.BoldText("Category:"), categoryStr)

	// Alternatives
	if showAlternatives && len(rec.Alternatives) > 0 {
		fmt.Println()
		fmt.Println(ui.Dim("  Alternatives:"))
		for i, alt := range rec.Alternatives {
			fmt.Printf("    %d. %s - %s\n", i+1, ui.Dim(alt.Command), ui.Dim(alt.Reason))
		}
	}

	fmt.Println()
	return nil
}
