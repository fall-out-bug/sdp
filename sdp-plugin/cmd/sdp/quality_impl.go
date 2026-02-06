package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/quality"
)

func runQualityCoverage() error {
	projectPath, _ := os.Getwd()
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}

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
		return fmt.Errorf("coverage check failed: %.1f%% < %.1f%%", result.Coverage, result.Threshold)
	}

	return nil
}

func runQualityComplexity() error {
	projectPath, _ := os.Getwd()
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}

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
		return fmt.Errorf("complexity check failed: max CC %d > threshold %d", result.MaxCC, result.Threshold)
	}

	return nil
}

func runQualitySize() error {
	projectPath, _ := os.Getwd()
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}

	result, err := checker.CheckFileSize()
	if err != nil {
		return fmt.Errorf("file size check failed: %w", err)
	}

	fmt.Printf("Total Files: %d\n", result.TotalFiles)
	fmt.Printf("Average LOC: %d\n", result.AverageLOC)
	fmt.Printf("Threshold: %d LOC\n", result.Threshold)
	fmt.Printf("Status: ")
	if result.Passed {
		fmt.Println("✓ PASSED")
	} else {
		fmt.Println("✗ FAILED")
	}

	if len(result.Violators) > 0 {
		fmt.Printf("\n%d files exceed threshold:\n", len(result.Violators))
		for _, f := range result.Violators {
			fmt.Printf("  %s: %d LOC\n", f.File, f.LOC)
		}
	}

	if !result.Passed {
		return fmt.Errorf("file size check failed: %d files exceed threshold", len(result.Violators))
	}

	return nil
}
