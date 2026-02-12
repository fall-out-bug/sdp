package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/collision"
	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/decision"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/guard"
)

// warnSimilarFailures loads past decisions and prints warning if similar failed decisions exist (AC3, AC4, AC8).
func warnSimilarFailures(wsID string) {
	root, err := config.FindProjectRoot()
	if err != nil {
		return
	}
	logger, err := decision.NewLogger(root)
	if err != nil {
		return
	}
	decisions, err := logger.LoadAll()
	if err != nil || len(decisions) == 0 {
		return
	}
	matches := evidence.FindSimilarDecisions("", nil, decisions)
	if len(matches) == 0 {
		return
	}
	fmt.Fprintln(os.Stderr, "Similar past decision(s) found:")
	for _, m := range matches {
		fmt.Fprintf(os.Stderr, "   Decision: %q\n", m.Question)
		fmt.Fprintf(os.Stderr, "   Outcome: %s\n", m.Outcome)
		fmt.Fprintf(os.Stderr, "   Source: WS %s\n", m.WorkstreamID)
		if len(m.Tags) > 0 {
			fmt.Fprintf(os.Stderr, "   Tags: %v\n", m.Tags)
		}
	}
	fmt.Fprintln(os.Stderr, "   Continue anyway? [y/N]")
}

// warnContractViolations checks for contract violations and prints warnings.
// AC2: @build runs contract validation after code generation.
// AC3: Contract violation = build warning (not error in P1; enforcement in P2).
func warnContractViolations() {
	// Find project root
	root, err := config.FindProjectRoot()
	if err != nil {
		return // Not in project, skip validation
	}

	// Check for .contracts directory
	contractsDir := filepath.Join(root, ".contracts")
	if _, err := os.Stat(contractsDir); os.IsNotExist(err) {
		return // No contracts, skip validation
	}

	// Find implementation directory (default: internal/)
	implDir := filepath.Join(root, "internal")
	if _, err := os.Stat(implDir); os.IsNotExist(err) {
		return // No impl dir, skip validation
	}

	// Run validation
	violations, err := collision.ValidateContractsInDir(contractsDir, implDir)
	if err != nil {
		// Validation failed (not violations list), log warning
		fmt.Fprintf(os.Stderr, "Contract validation error: %v\n", err)
		return
	}

	// Print violations as warnings
	if len(violations) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "Contract violations detected (%d):\n", len(violations))
	for _, v := range violations {
		severity := "WARNING"
		if v.Severity == "error" {
			severity = "ERROR"
		}
		fmt.Fprintf(os.Stderr, "   [%s] %s: %s\n", severity, v.Type, v.Message)
	}

	// In P1, violations are warnings only (not blocking)
	fmt.Fprintf(os.Stderr, "   Review violations before proceeding\n")
}

// printHumanReadableResult prints check results in human-readable format (AC4)
func printHumanReadableResult(result *guard.CheckResult) {
	if result.Success && result.Summary.Total == 0 {
		fmt.Println("No issues found")
		return
	}

	// Print findings
	for _, f := range result.Findings {
		severity := string(f.Severity)
		icon := "WARNING"
		if f.Severity == guard.SeverityError {
			icon = "ERROR"
		}
		fmt.Printf("[%s] %s: %s\n", icon, severity, f.Message)
		if f.File != "" {
			fmt.Printf("   File: %s\n", f.File)
			if f.Line > 0 {
				fmt.Printf("   Line: %d\n", f.Line)
			}
		}
	}

	// Print summary
	fmt.Printf("\nSummary: %d total (%d errors, %d warnings)\n",
		result.Summary.Total, result.Summary.Errors, result.Summary.Warnings)

	// AC8: Hybrid mode message
	if !result.Success {
		fmt.Println("\nCheck failed: Errors must be fixed before committing")
	} else if result.Summary.Warnings > 0 {
		fmt.Println("\nCheck passed with warnings")
	}
}
