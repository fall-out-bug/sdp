package main

import (
	"fmt"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/collision"
	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/spf13/cobra"
)

func runContractValidate(cmd *cobra.Command, args []string) error {
	implDir, err := cmd.Flags().GetString("impl-dir")
	if err != nil {
		return fmt.Errorf("failed to get impl-dir flag: %w", err)
	}
	contractsDir, err := cmd.Flags().GetString("contracts-dir")
	if err != nil {
		return fmt.Errorf("failed to get contracts-dir flag: %w", err)
	}

	// If impl-dir is specified, validate implementation against contracts
	if implDir != "" {
		return validateImplementation(contractsDir, implDir)
	}

	// Otherwise, use original contract-to-contract validation (placeholder)
	contractPaths, err := cmd.Flags().GetStringSlice("contracts")
	if err != nil {
		return fmt.Errorf("failed to get contracts flag: %w", err)
	}
	reportPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %w", err)
	}

	if len(contractPaths) < 2 {
		return fmt.Errorf("at least 2 contracts required for validation")
	}

	fmt.Printf("✓ Validating %d contracts...\n", len(contractPaths))
	fmt.Printf("  Report: %s\n", reportPath)
	fmt.Printf("\n⚠️  Contract-to-contract validation not yet implemented\n")
	fmt.Printf("   Use --impl-dir to validate implementation against contracts\n")

	return nil
}

// validateImplementation validates implementation files against contracts.
func validateImplementation(contractsDir, implDir string) error {
	root, err := config.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}

	// Resolve paths
	if !filepath.IsAbs(contractsDir) {
		contractsDir = filepath.Join(root, contractsDir)
	}
	if !filepath.IsAbs(implDir) {
		implDir = filepath.Join(root, implDir)
	}

	violations, err := collision.ValidateContractsInDir(contractsDir, implDir)
	if err != nil {
		return fmt.Errorf("validate contracts: %w", err)
	}

	if len(violations) == 0 {
		fmt.Println("✓ No contract violations found")
		fmt.Println("  All implementations match their contracts")
		return nil
	}

	fmt.Printf("⚠️  Found %d contract violation(s):\n", len(violations))
	fmt.Println()

	errorCount := 0
	warningCount := 0
	for _, v := range violations {
		if v.Severity == "error" {
			fmt.Printf("  ❌ %s: %s\n", v.Field, v.Message)
			errorCount++
		} else {
			fmt.Printf("  ⚠️  %s: %s\n", v.Field, v.Message)
			warningCount++
		}
	}

	fmt.Println()
	fmt.Printf("  Errors: %d, Warnings: %d\n", errorCount, warningCount)
	fmt.Println("  Note: Extra fields are warnings in P1 (enforcement in P2)")

	if errorCount > 0 {
		return fmt.Errorf("%d contract violation(s) found", errorCount)
	}

	return nil
}
