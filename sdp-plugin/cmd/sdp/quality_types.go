package main

import (
	"fmt"
	"os"

	"github.com/ai-masters/sdp/internal/quality"
)

func runQualityTypes() error {
	projectPath, _ := os.Getwd()
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}

	result, err := checker.CheckTypes()
	if err != nil {
		return fmt.Errorf("type check failed: %w", err)
	}

	fmt.Printf("Project Type: %s\n", result.ProjectType)
	fmt.Printf("Status: ")
	if result.Passed {
		fmt.Println("✓ PASSED")
	} else {
		fmt.Println("✗ FAILED")
	}

	if len(result.Errors) > 0 {
		fmt.Printf("\n%d errors:\n", len(result.Errors))
		for _, e := range result.Errors {
			if e.Line > 0 {
				fmt.Printf("  %s:%d: %s\n", e.File, e.Line, e.Message)
			} else {
				fmt.Printf("  %s: %s\n", e.File, e.Message)
			}
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Printf("\n%d warnings:\n", len(result.Warnings))
		for _, w := range result.Warnings {
			fmt.Printf("  %s\n", w.Message)
		}
	}

	if !result.Passed {
		os.Exit(1)
	}

	return nil
}

func runQualityAll() error {
	projectPath, _ := os.Getwd()
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}

	fmt.Println("Running all quality checks...")
	fmt.Println()

	// Coverage
	fmt.Println("=== Coverage ===")
	covResult, _ := checker.CheckCoverage()
	fmt.Printf("Coverage: %.1f%% (threshold: %.1f%%) ", covResult.Coverage, covResult.Threshold)
	if covResult.Passed {
		fmt.Println("✓")
	} else {
		fmt.Println("✗")
	}

	// Complexity
	fmt.Println("\n=== Complexity ===")
	ccResult, _ := checker.CheckComplexity()
	fmt.Printf("Max CC: %d (threshold: %d) ", ccResult.MaxCC, ccResult.Threshold)
	if ccResult.Passed {
		fmt.Println("✓")
	} else {
		fmt.Println("✗")
	}

	// File Size
	fmt.Println("\n=== File Size ===")
	sizeResult, _ := checker.CheckFileSize()
	fmt.Printf("Violators: %d (threshold: %d LOC) ", len(sizeResult.Violators), sizeResult.Threshold)
	if sizeResult.Passed {
		fmt.Println("✓")
	} else {
		fmt.Println("✗")
	}

	// Types
	fmt.Println("\n=== Types ===")
	typeResult, _ := checker.CheckTypes()
	fmt.Printf("Errors: %d ", len(typeResult.Errors))
	if typeResult.Passed {
		fmt.Println("✓")
	} else {
		fmt.Println("✗")
	}

	fmt.Println()
	allPassed := covResult.Passed && ccResult.Passed && sizeResult.Passed && typeResult.Passed
	if allPassed {
		fmt.Println("Overall: ✓ ALL CHECKS PASSED")
	} else {
		fmt.Println("Overall: ✗ SOME CHECKS FAILED")
		os.Exit(1)
	}

	return nil
}
