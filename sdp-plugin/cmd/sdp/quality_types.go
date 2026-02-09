package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/quality"
)

func runQualityTypes(strict bool) error {
	projectPath, err := os.Getwd()
	if err != nil {
		projectPath = "." // Fall back to current directory
	}
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	checker.SetStrictMode(strict)

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

	if evidence.Enabled() {
		evidence.Emit(evidence.VerificationEvent("", result.Passed, "types", 0))
	}
	if !result.Passed {
		return fmt.Errorf("type check failed")
	}
	return nil
}

func runQualityAll(strict bool) error {
	projectPath, err := os.Getwd()
	if err != nil {
		projectPath = "." // Fall back to current directory
	}
	checker, err := quality.NewChecker(projectPath)
	if err != nil {
		return fmt.Errorf("failed to create checker: %w", err)
	}
	checker.SetStrictMode(strict)

	fmt.Println("Running all quality checks...")
	if strict {
		fmt.Println("Mode: STRICT (violations = errors)")
	} else {
		fmt.Println("Mode: PRAGMATIC (violations = warnings)")
	}
	fmt.Println()

	// Coverage
	fmt.Println("=== Coverage ===")
	covResult, _ := checker.CheckCoverage() //nolint:errcheck // UI display, error is non-critical
	fmt.Printf("Coverage: %.1f%% (threshold: %.1f%%) ", covResult.Coverage, covResult.Threshold)
	if covResult.Passed {
		fmt.Println("✓")
	} else {
		fmt.Println("✗")
	}

	// Complexity
	fmt.Println("\n=== Complexity ===")
	ccResult, _ := checker.CheckComplexity() //nolint:errcheck // UI display, error is non-critical
	fmt.Printf("Max CC: %d (threshold: %d) ", ccResult.MaxCC, ccResult.Threshold)
	if ccResult.Passed {
		fmt.Println("✓")
	} else {
		fmt.Println("✗")
	}

	// File Size
	fmt.Println("\n=== File Size ===")
	sizeResult, _ := checker.CheckFileSize() //nolint:errcheck // UI display, error is non-critical
	if len(sizeResult.Warnings) > 0 {
		fmt.Printf("Warnings: %d (threshold: %d LOC) ⚠️\n", len(sizeResult.Warnings), sizeResult.Threshold)
	}
	if len(sizeResult.Violators) > 0 {
		fmt.Printf("Errors: %d (threshold: %d LOC) ✗\n", len(sizeResult.Violators), sizeResult.Threshold)
	}
	if len(sizeResult.Warnings) == 0 && len(sizeResult.Violators) == 0 {
		fmt.Println("No violations ✓")
	}

	// Types
	fmt.Println("\n=== Types ===")
	typeResult, _ := checker.CheckTypes() //nolint:errcheck // UI display, error is non-critical
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
		return fmt.Errorf("quality checks failed")
	}

	return nil
}
