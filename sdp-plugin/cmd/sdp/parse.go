package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ai-masters/sdp/internal/parser"
	"github.com/spf13/cobra"
)

func parseCmd() *cobra.Command {
	var validateFlag bool

	cmd := &cobra.Command{
		Use:   "parse <ws-id>",
		Short: "Parse and display workstream information",
		Long: `Parse a workstream markdown file and display its contents.

Args:
  ws-id    Workstream ID (e.g., 00-050-01)

Examples:
  sdp parse 00-050-01
  sdp parse --validate docs/workstreams/backlog/00-050-01.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if validateFlag {
				return validateRun(cmd, args)
			}

			return parseRun(cmd, args)
		},
	}

	cmd.Flags().BoolVar(&validateFlag, "validate", false, "Validate workstream file")

	return cmd
}

func parseRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("workstream ID required")
	}

	wsID := args[0]

	// Sanitize wsID to prevent path traversal attacks
	wsID = filepath.Clean(wsID)
	if wsID != args[0] {
		return fmt.Errorf("invalid workstream ID: path traversal detected")
	}

	// Find workstream file
	wsPath, err := findWorkstreamFile(wsID)
	if err != nil {
		return err
	}

	// Parse workstream
	ws, err := parser.ParseWorkstream(wsPath)
	if err != nil {
		return fmt.Errorf("failed to parse workstream: %w", err)
	}

	// Display workstream
	displayWorkstream(ws)

	return nil
}

func validateRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("workstream file path required")
	}

	wsPath := args[0]

	// Sanitize path to prevent path traversal attacks
	wsPath = filepath.Clean(wsPath)
	// Check if path tries to escape expected directories
	if containsPathTraversal(wsPath) {
		return fmt.Errorf("invalid file path: path traversal detected")
	}

	// Validate file
	issues, err := parser.ValidateFile(wsPath)
	if err != nil {
		return fmt.Errorf("failed to validate workstream: %w", err)
	}

	if len(issues) == 0 {
		fmt.Println("✅ No validation issues found")
		return nil
	}

	fmt.Printf("Found %d validation issue(s):\n\n", len(issues))
	for _, issue := range issues {
		symbol := "⚠️"
		if issue.Severity == "ERROR" {
			symbol = "❌"
		}
		fmt.Printf("%s [%s] %s: %s\n", symbol, issue.Severity, issue.Field, issue.Message)
	}

	// Return error if there are ERROR severity issues
	for _, issue := range issues {
		if issue.Severity == "ERROR" {
			return fmt.Errorf("validation failed with %d error(s)", countErrors(issues))
		}
	}

	return nil
}

func findWorkstreamFile(wsID string) (string, error) {
	// Try common locations
	locations := []string{
		filepath.Join("docs/workstreams/backlog", wsID+".md"),
		filepath.Join("docs/workstreams/in_progress", wsID+".md"),
		filepath.Join("docs/workstreams/completed", wsID+".md"),
		filepath.Join("..", "docs", "workstreams", "backlog", wsID+".md"),
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc, nil
		}
	}

	return "", fmt.Errorf("workstream file not found: %s", wsID)
}

func displayWorkstream(ws *parser.Workstream) {
	fmt.Printf("📋 Workstream: %s\n", ws.ID)
	fmt.Printf("Feature: %s\n", ws.Feature)
	fmt.Printf("Status: %s\n", ws.Status)
	fmt.Printf("Size: %s\n", ws.Size)
	fmt.Printf("\n### Goal\n\n%s\n\n", ws.Goal)

	if len(ws.Acceptance) > 0 {
		fmt.Println("### Acceptance Criteria")
		for i, ac := range ws.Acceptance {
			fmt.Printf("  %d. %s\n", i+1, ac)
		}
		fmt.Println()
	}

	if len(ws.Scope.Implementation) > 0 || len(ws.Scope.Tests) > 0 {
		fmt.Println("### Scope Files")

		if len(ws.Scope.Implementation) > 0 {
			fmt.Println("\n  Implementation:")
			for _, f := range ws.Scope.Implementation {
				fmt.Printf("    - %s\n", f)
			}
		}

		if len(ws.Scope.Tests) > 0 {
			fmt.Println("\n  Tests:")
			for _, f := range ws.Scope.Tests {
				fmt.Printf("    - %s\n", f)
			}
		}
		fmt.Println()
	}
}

func countErrors(issues []parser.ValidationIssue) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == "ERROR" {
			count++
		}
	}
	return count
}

// containsPathTraversal checks if a path contains traversal attempts
func containsPathTraversal(path string) bool {
	// Check for obvious traversal patterns
	traversalPatterns := []string{
		"../",
		"..\\",
		"~/.", // Don't allow home directory access
	}
	for _, pattern := range traversalPatterns {
		if contains(path, pattern) {
			return true
		}
	}
	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

// findSubstring finds a substring in a string
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
