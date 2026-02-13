package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/collision"
	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

func runContractVerify(cmd *cobra.Command, args []string) error {
	featureName, err := cmd.Flags().GetString("feature")
	if err != nil {
		return fmt.Errorf("failed to get feature flag: %w", err)
	}
	contractPath, err := cmd.Flags().GetString("contract")
	if err != nil {
		return fmt.Errorf("failed to get contract flag: %w", err)
	}

	// Default contract path from feature name
	if contractPath == "" && featureName != "" {
		contractPath = fmt.Sprintf(".contracts/%s.yaml", featureName)
	}

	if contractPath == "" {
		return fmt.Errorf("either --feature or --contract must be specified")
	}

	// Derive lock path
	lockPath := strings.TrimSuffix(contractPath, filepath.Ext(contractPath)) + ".lock"

	// Run internal verify function
	matched, err := runContractVerifyInternal(contractPath, lockPath)
	if err != nil {
		return err
	}

	if matched {
		// Exit code 0 for success
		return nil
	}

	// Exit code 1 for mismatch
	return fmt.Errorf("contract mismatch detected")
}

// runContractVerifyInternal implements the core verify logic
func runContractVerifyInternal(contractPath, lockPath string) (bool, error) {
	// Check if lock exists
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		return false, fmt.Errorf("lock file not found: %s", lockPath)
	}

	// Check if contract exists
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		return false, fmt.Errorf("contract file not found: %s", contractPath)
	}

	// Read lock file
	lockData, err := os.ReadFile(lockPath)
	if err != nil {
		return false, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock ContractLock
	if err := yaml.Unmarshal(lockData, &lock); err != nil {
		return false, fmt.Errorf("failed to parse lock file: %w", err)
	}

	// Read contract file
	contractData, err := os.ReadFile(contractPath)
	if err != nil {
		return false, fmt.Errorf("failed to read contract: %w", err)
	}

	// Calculate current contract hash
	hash := sha256.Sum256(contractData)
	currentHash := hex.EncodeToString(hash[:])

	// Compare hashes
	if currentHash == lock.ContractHash {
		// Match
		fmt.Printf("✓ Contract matches lock\n")
		fmt.Printf("✓ Locked SHA: %s\n", lock.GitSHA)
		fmt.Printf("✓ Locked at: %s\n", lock.LockedAt)
		return true, nil
	}

	// Mismatch
	fmt.Printf("✗ Contract mismatch detected!\n")
	expectedHashDisplay := lock.ContractHash
	if len(expectedHashDisplay) > 16 {
		expectedHashDisplay = expectedHashDisplay[:16] + "..."
	}
	actualHashDisplay := currentHash
	if len(actualHashDisplay) > 16 {
		actualHashDisplay = actualHashDisplay[:16] + "..."
	}
	fmt.Printf("  Expected hash: %s\n", expectedHashDisplay)
	fmt.Printf("  Actual hash:   %s\n", actualHashDisplay)
	fmt.Println()
	fmt.Println("Contract has been modified since lock.")
	fmt.Println("Please re-lock or restore original contract.")
	return false, nil
}
