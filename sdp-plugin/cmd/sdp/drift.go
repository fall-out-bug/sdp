package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ai-masters/sdp/internal/drift"
	"github.com/spf13/cobra"
)

func driftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect documentation-code drift",
		Long: `Detect drift between workstream documentation and actual code.

This helps prevent "wrong_approach" friction by validating that
workstream descriptions match actual codebase reality.`,
	}

	cmd.AddCommand(driftDetectCmd())

	return cmd
}

func driftDetectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detect <ws-id>",
		Short: "Detect drift for a workstream",
		Long: `Detect drift between documentation and code for a specific workstream.

Parses the workstream markdown file and checks:
  - All files in scope exist (ERROR if missing)
  - Files contain expected entities (WARNING if empty)
  - Generates actionable recommendations

Example:
  sdp drift detect 00-050-01`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsID := args[0]

			// Find project root
			projectRoot, err := findDriftProjectRoot()
			if err != nil {
				return fmt.Errorf("failed to find project root: %w", err)
			}

			// Find workstream file
			wsPath, err := findDriftWorkstreamFile(projectRoot, wsID)
			if err != nil {
				return fmt.Errorf("failed to find workstream file: %w", err)
			}

			// Create detector
			detector := drift.NewDetector(projectRoot)

			// Detect drift
			report, err := detector.DetectDrift(wsPath)
			if err != nil {
				return fmt.Errorf("failed to detect drift: %w", err)
			}

			// Print report
			fmt.Println(report.String())

			// Exit with error if verdict is FAIL
			if report.Verdict == "FAIL" {
				return fmt.Errorf("drift detected - %d error(s), %d warning(s)",
					countDriftErrors(report), countDriftWarnings(report))
			}

			return nil
		},
	}

	return cmd
}

// findDriftProjectRoot finds the project root by looking for .beads or docs directory
func findDriftProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if we're in sdp-plugin directory
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		// We're in sdp-plugin, go up one level
		parent := filepath.Dir(cwd)
		if _, err := os.Stat(filepath.Join(parent, "docs")); err == nil {
			return parent, nil
		}
	}

	// Check if we're already in project root
	if _, err := os.Stat(filepath.Join(cwd, "docs")); err == nil {
		return cwd, nil
	}

	// Check for .beads directory
	if _, err := os.Stat(filepath.Join(cwd, ".beads")); err == nil {
		return cwd, nil
	}

	// Traverse up looking for project root
	current := cwd
	for {
		if _, err := os.Stat(filepath.Join(current, "docs")); err == nil {
			return current, nil
		}
		if _, err := os.Stat(filepath.Join(current, ".beads")); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached root, return current directory
			return cwd, nil
		}
		current = parent
	}
}

// findDriftWorkstreamFile finds a workstream file by ID
func findDriftWorkstreamFile(projectRoot, wsID string) (string, error) {
	// Try backlog first
	wsPath := filepath.Join(projectRoot, "docs", "workstreams", "backlog", wsID+".md")
	if _, err := os.Stat(wsPath); err == nil {
		return wsPath, nil
	}

	// Try in_progress
	wsPath = filepath.Join(projectRoot, "docs", "workstreams", "in_progress", wsID+".md")
	if _, err := os.Stat(wsPath); err == nil {
		return wsPath, nil
	}

	// Try completed
	wsPath = filepath.Join(projectRoot, "docs", "workstreams", "completed", wsID+".md")
	if _, err := os.Stat(wsPath); err == nil {
		return wsPath, nil
	}

	return "", fmt.Errorf("workstream file not found for %s", wsID)
}

// countDriftErrors counts ERROR status issues in report
func countDriftErrors(report *drift.DriftReport) int {
	count := 0
	for _, issue := range report.Issues {
		if issue.Status == "ERROR" {
			count++
		}
	}
	return count
}

// countDriftWarnings counts WARNING status issues in report
func countDriftWarnings(report *drift.DriftReport) int {
	count := 0
	for _, issue := range report.Issues {
		if issue.Status == "WARNING" {
			count++
		}
	}
	return count
}
