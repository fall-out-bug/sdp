package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/parser"
	"github.com/spf13/cobra"
)

func validateRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("workstream file path required")
	}

	wsPath := args[0]

	wsPath = filepath.Clean(wsPath)
	if containsPathTraversal(wsPath) {
		return fmt.Errorf("invalid file path: path traversal detected")
	}

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

	for _, issue := range issues {
		if issue.Severity == "ERROR" {
			return fmt.Errorf("validation failed with %d error(s)", countErrors(issues))
		}
	}

	return nil
}

func findWorkstreamFile(wsID string) (string, error) {
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

func containsPathTraversal(path string) bool {
	traversalPatterns := []string{
		"../",
		"..\\",
		"~/.",
	}
	for _, pattern := range traversalPatterns {
		if contains(path, pattern) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
