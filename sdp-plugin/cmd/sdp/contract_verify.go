package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

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
