package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/config"
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

	// In adoption mode: skip enforcement, still log evidence
	if adoptionModeSkip("types", result.Passed) {
		fmt.Printf("Status: SKIP (adoption mode)\n")
		return nil
	}

	fmt.Printf("Status: ")
	if result.Passed {
		fmt.Println("PASSED")
	} else {
		fmt.Println("FAILED")
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
		if err := evidence.EmitSync(evidence.VerificationEvent("", result.Passed, "types", 0)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: evidence emit: %v\n", err)
		}
	}
	if !result.Passed {
		return fmt.Errorf("type check failed")
	}
	return nil
}

func runQualityAll(ctx context.Context, strict bool) error {
	if ctx == nil {
		ctx = context.Background()
	}
	adoptionMode := config.IsAdoptionMode()
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
	if adoptionMode {
		fmt.Println("Mode: ADOPTION (gates skipped, evidence logged)")
	} else if strict {
		fmt.Println("Mode: STRICT (violations = errors)")
	} else {
		fmt.Println("Mode: PRAGMATIC (violations = warnings)")
	}
	fmt.Println()

	// Coverage
	fmt.Println("=== Coverage ===")
	covResult, err := checker.CheckCoverage(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: coverage check failed: %v\n", err)
		covResult = &quality.CoverageResult{Coverage: 0, Threshold: 80, Passed: false}
	}
	fmt.Printf("Coverage: %.1f%% (threshold: %.1f%%) ", covResult.Coverage, covResult.Threshold)
	if adoptionMode {
		fmt.Println("SKIP")
		if evidence.Enabled() {
			_ = evidence.EmitSync(evidence.VerificationEvent("", covResult.Passed, "coverage", covResult.Coverage))
		}
	} else if covResult.Passed {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}

	// Complexity
	fmt.Println("\n=== Complexity ===")
	ccResult, err := checker.CheckComplexity()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: complexity check failed: %v\n", err)
		ccResult = &quality.ComplexityResult{MaxCC: 0, Threshold: 40, Passed: false}
	}
	fmt.Printf("Max CC: %d (threshold: %d) ", ccResult.MaxCC, ccResult.Threshold)
	if adoptionMode {
		fmt.Println("SKIP")
		if evidence.Enabled() {
			_ = evidence.EmitSync(evidence.VerificationEvent("", ccResult.Passed, "complexity", 0))
		}
	} else if ccResult.Passed {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}

	// File Size
	fmt.Println("\n=== File Size ===")
	sizeResult, err := checker.CheckFileSize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: file size check failed: %v\n", err)
		sizeResult = &quality.FileSizeResult{Threshold: 200, Passed: false}
	}
	if adoptionMode {
		fmt.Printf("Threshold: %d LOC -- SKIP\n", sizeResult.Threshold)
		if evidence.Enabled() {
			_ = evidence.EmitSync(evidence.VerificationEvent("", sizeResult.Passed, "size", 0))
		}
	} else {
		if len(sizeResult.Warnings) > 0 {
			fmt.Printf("Warnings: %d (threshold: %d LOC)\n", len(sizeResult.Warnings), sizeResult.Threshold)
		}
		if len(sizeResult.Violators) > 0 {
			fmt.Printf("Errors: %d (threshold: %d LOC)\n", len(sizeResult.Violators), sizeResult.Threshold)
		}
		if len(sizeResult.Warnings) == 0 && len(sizeResult.Violators) == 0 {
			fmt.Println("No violations")
		}
	}

	// Types
	fmt.Println("\n=== Types ===")
	typeResult, err := checker.CheckTypes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: type check failed: %v\n", err)
		typeResult = &quality.TypeResult{Passed: false}
	}
	fmt.Printf("Errors: %d ", len(typeResult.Errors))
	if adoptionMode {
		fmt.Println("SKIP")
		if evidence.Enabled() {
			_ = evidence.EmitSync(evidence.VerificationEvent("", typeResult.Passed, "types", 0))
		}
	} else if typeResult.Passed {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}

	fmt.Println()
	if adoptionMode {
		fmt.Println("Overall: SKIP (adoption mode active)")
		return nil
	}

	allPassed := covResult.Passed && ccResult.Passed && sizeResult.Passed && typeResult.Passed
	if allPassed {
		fmt.Println("Overall: ALL CHECKS PASSED")
	} else {
		fmt.Println("Overall: SOME CHECKS FAILED")
		return fmt.Errorf("quality checks failed")
	}

	return nil
}
