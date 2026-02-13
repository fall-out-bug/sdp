package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ContractLock represents the lock file structure
type ContractLock struct {
	ContractFile string `yaml:"contract_file"`
	ContractHash string `yaml:"contract_hash"`
	GitSHA       string `yaml:"git_sha"`
	LockedAt     string `yaml:"locked_at"`
	Checksum     string `yaml:"checksum"`
	Metadata     struct {
		Feature   string `yaml:"feature"`
		Version   string `yaml:"version"`
		Endpoints int    `yaml:"endpoints"`
		Schemas   int    `yaml:"schemas"`
	} `yaml:"metadata,omitempty"`
}

func runContractLock(cmd *cobra.Command, args []string) error {
	contractPath, err := cmd.Flags().GetString("contract")
	if err != nil {
		return fmt.Errorf("failed to get contract flag: %w", err)
	}
	gitSHA, err := cmd.Flags().GetString("sha")
	if err != nil {
		return fmt.Errorf("failed to get sha flag: %w", err)
	}
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return fmt.Errorf("failed to get force flag: %w", err)
	}

	// Default contract path if not provided
	featureName := ""
	if contractPath == "" {
		// Try to get feature name from contract path or use default
		featureName = "feature"
		contractPath = fmt.Sprintf(".contracts/%s.yaml", featureName)
	} else {
		// Extract feature name from contract path
		base := filepath.Base(contractPath)
		featureName = strings.TrimSuffix(base, filepath.Ext(base))
	}

	// Default git SHA if not provided
	if gitSHA == "" {
		gitSHA = "unknown"
	}

	// Derive lock path from contract path
	lockPath := strings.TrimSuffix(contractPath, filepath.Ext(contractPath)) + ".lock"

	// Run internal lock function
	return runContractLockInternal(featureName, gitSHA, contractPath, lockPath, force)
}

// runContractLockInternal implements the core lock logic
func runContractLockInternal(featureName, gitSHA, contractPath, lockPath string, force bool) error {
	// Check if contract exists
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		return fmt.Errorf("contract file not found: %s", contractPath)
	}

	// Check if lock already exists
	if !force {
		if _, err := os.Stat(lockPath); err == nil {
			return fmt.Errorf("lock file already exists: %s (use --force to re-lock)", lockPath)
		}
	}

	// Read contract content
	contractContent, err := os.ReadFile(contractPath)
	if err != nil {
		return fmt.Errorf("failed to read contract: %w", err)
	}

	// Calculate SHA256 hash
	hash := sha256.Sum256(contractContent)
	contractHash := hex.EncodeToString(hash[:])

	// Get current timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Calculate checksum (hash of contract hash + git SHA + timestamp)
	checksumInput := contractHash + gitSHA + timestamp
	checksumHash := sha256.Sum256([]byte(checksumInput))
	checksum := "sha256:" + hex.EncodeToString(checksumHash[:])

	// Create lock structure
	lock := ContractLock{
		ContractFile: contractPath,
		ContractHash: contractHash,
		GitSHA:       gitSHA,
		LockedAt:     timestamp,
		Checksum:     checksum,
	}

	// Try to parse contract to extract metadata
	var contract map[string]interface{}
	if err := yaml.Unmarshal(contractContent, &contract); err == nil {
		// Extract metadata if available
		if _, ok := contract["info"].(map[string]interface{}); ok {
			lock.Metadata.Feature = featureName
			lock.Metadata.Version = "1.0.0"
		}

		// Count endpoints and schemas
		if paths, ok := contract["paths"].(map[string]interface{}); ok {
			lock.Metadata.Endpoints = len(paths)
		}
		if schemas, ok := contract["components"].(map[string]interface{}); ok {
			if s, ok := schemas["schemas"].(map[string]interface{}); ok {
				lock.Metadata.Schemas = len(s)
			}
		}
	}

	// Marshal lock to YAML
	lockData, err := yaml.Marshal(lock)
	if err != nil {
		return fmt.Errorf("failed to marshal lock: %w", err)
	}

	// Write lock file
	if err := os.WriteFile(lockPath, lockData, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	// Print success message
	fmt.Printf("✓ Contract read: %s\n", contractPath)
	fmt.Printf("✓ Contract hash: %s\n", contractHash[:16]+"...")
	fmt.Printf("✓ Lock file created: %s\n", lockPath)
	fmt.Printf("✓ Locked at: %s\n", timestamp)
	fmt.Printf("✓ Git SHA: %s\n", gitSHA)
	fmt.Println()
	fmt.Println("Contract is now immutable. Any changes will require re-lock.")

	return nil
}
