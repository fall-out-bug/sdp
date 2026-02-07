package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/quality"
)

func runQualityCoverage(strict bool) error {
	projectPath, err := os.Getwd()
	if err != nil {
		projectPath = "." // Fall back to current directory
	}
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	checker.SetStrictMode(strict)

	result, err := checker.CheckCoverage()
	if err != nil {
		return fmt.Errorf("coverage check failed: %w", err)
	}

	fmt.Printf("Project Type: %s\n", result.ProjectType)
	fmt.Printf("Coverage: %.1f%%\n", result.Coverage)
	fmt.Printf("Threshold: %.1f%%\n", result.Threshold)
	fmt.Printf("Status: ")
	if result.Passed {
		fmt.Println("✓ PASSED")
	} else {
		fmt.Println("✗ FAILED")
	}

	if result.Report != "" {
		fmt.Printf("\n%s\n", result.Report)
	}

	if len(result.FilesBelow) > 0 {
		fmt.Println("\nFiles below threshold:")
		for _, f := range result.FilesBelow {
			fmt.Printf("  %s: %.1f%%\n", f.File, f.Coverage)
		}
	}

	if !result.Passed {
		return fmt.Errorf("quality check failed")
	}

	return nil
}

func runQualityComplexity(strict bool) error {
	projectPath, err := os.Getwd()
	if err != nil {
		projectPath = "."
	}
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	checker.SetStrictMode(strict)

	result, err := checker.CheckComplexity()
	if err != nil {
		return fmt.Errorf("complexity check failed: %w", err)
	}

	fmt.Printf("Average CC: %.1f\n", result.AverageCC)
	fmt.Printf("Max CC: %d\n", result.MaxCC)
	fmt.Printf("Threshold: %d\n", result.Threshold)
	fmt.Printf("Status: ")
	if result.Passed {
		fmt.Println("✓ PASSED")
	} else {
		fmt.Println("✗ FAILED")
	}

	if len(result.ComplexFiles) > 0 {
		fmt.Printf("\n%d files exceed threshold:\n", len(result.ComplexFiles))
		for _, f := range result.ComplexFiles {
			fmt.Printf("  %s: CC %.1f (max: %d)\n", f.File, f.AverageCC, f.MaxCC)
		}
	}

	if !result.Passed {
		return fmt.Errorf("quality check failed")
	}

	return nil
}

func runQualitySize(strict bool) error {
	projectPath, err := os.Getwd()
	if err != nil {
		projectPath = "."
	}
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	checker.SetStrictMode(strict)

	result, err := checker.CheckFileSize()
	if err != nil {
		return fmt.Errorf("file size check failed: %w", err)
	}

	fmt.Printf("Total Files: %d\n", result.TotalFiles)
	fmt.Printf("Average LOC: %d\n", result.AverageLOC)
	fmt.Printf("Threshold: %d LOC\n", result.Threshold)
	fmt.Printf("Mode: ")
	if result.Strict {
		fmt.Println("STRICT (violations = errors)")
	} else {
		fmt.Println("PRAGMATIC (violations = warnings)")
	}

	// Show warnings (pragmatic mode)
	if len(result.Warnings) > 0 {
		fmt.Printf("\n⚠️  WARNINGS (%d):\n", len(result.Warnings))
		for _, f := range result.Warnings {
			fmt.Printf("  - %s: %d LOC (exceeds %d LOC threshold)\n", f.File, f.LOC, result.Threshold)
		}
	}

	// Show errors (strict mode)
	if len(result.Violators) > 0 {
		fmt.Printf("\n✗ ERRORS (%d):\n", len(result.Violators))
		for _, f := range result.Violators {
			fmt.Printf("  - %s: %d LOC (exceeds %d LOC threshold)\n", f.File, f.LOC, result.Threshold)
		}
	}

	// Show status
	fmt.Printf("\nStatus: ")
	if result.Passed {
		fmt.Println("✓ PASSED")
	} else {
		fmt.Println("✗ FAILED")
	}

	if !result.Passed {
		return fmt.Errorf("quality check failed")
	}

	return nil
}
