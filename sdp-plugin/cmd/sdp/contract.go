package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/collision"
	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// contractCmd returns the contract management command
func contractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract",
		Short: "Manage API contracts for component validation",
		Long: `Manage API contracts for component validation.

Commands:
  synthesize - Generate contract from requirements
  lock       - Lock contract as source of truth
  validate   - Validate contracts against each other`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// synthesize subcommand
	synthesizeCmd := &cobra.Command{
		Use:   "synthesize",
		Short: "Generate contract from requirements",
		Long: `Generate OpenAPI 3.0 contract from feature requirements.

Multi-agent synthesis:
1. Architect analyzes requirements
2. Proposes initial contract
3. Agents review in parallel (frontend/backend/sdk)
4. Synthesizer resolves conflicts
5. Outputs locked contract`,
		RunE: runContractSynthesize,
	}

	var featureName string
	var requirementsPath string
	var outputPath string

	synthesizeCmd.Flags().StringVar(&featureName, "feature", "", "Feature name (required)")
	synthesizeCmd.Flags().StringVar(&requirementsPath, "requirements", "", "Path to requirements document")
	synthesizeCmd.Flags().StringVar(&outputPath, "output", "", "Output contract path")

	// Mark flag as required (ignore error - programming error if this fails)
	_ = synthesizeCmd.MarkFlagRequired("feature") //nolint:errcheck

	cmd.AddCommand(synthesizeCmd)

	// generate subcommand (cross-feature contract generation)
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate contracts from shared boundaries",
		Long: `Generate interface contracts from shared boundaries detected across features.

Creates .contracts/<type>.yaml files defining the agreed interface
that multiple features must respect.`,
		RunE: runContractGenerate,
	}

	var featuresFlag string
	generateCmd.Flags().StringVar(&featuresFlag, "features", "", "Comma-separated feature IDs (e.g., F054,F055)")

	cmd.AddCommand(generateCmd)

	// lock subcommand
	lockCmd := &cobra.Command{
		Use:   "lock",
		Short: "Lock contract as source of truth",
		Long: `Lock contract to prevent modifications during implementation.

Creates .lock file with SHA256 checksum. Prevents agents
from diverging from agreed contract.`,
		RunE: runContractLock,
	}

	var contractPath string
	var gitSHA string
	var forceLock bool

	lockCmd.Flags().StringVar(&contractPath, "contract", "", "Contract file path")
	lockCmd.Flags().StringVar(&gitSHA, "sha", "", "Git commit SHA")
	lockCmd.Flags().BoolVar(&forceLock, "force", false, "Force re-lock if lock exists")

	cmd.AddCommand(lockCmd)

	// validate subcommand
	var contractPaths []string
	var reportPath string

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate contracts against implementation",
		Long: `Validate contracts against implementation files.

Detects:
- Missing required fields
- Type mismatches
- Extra fields (warning in P1)

Flags:
  --impl-dir: Directory containing implementation files
  --contracts-dir: Directory containing contract files`,
		RunE: runContractValidate,
	}

	var implDir string
	var contractsDir string
	validateCmd.Flags().StringSliceVar(&contractPaths, "contracts", []string{}, "Contract files to validate (min 2)")
	validateCmd.Flags().StringVar(&reportPath, "output", "", "Validation report output")
	validateCmd.Flags().StringVar(&implDir, "impl-dir", "", "Implementation directory")
	validateCmd.Flags().StringVar(&contractsDir, "contracts-dir", ".contracts", "Contracts directory")

	cmd.AddCommand(validateCmd)

	// verify subcommand
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify contract matches lock",
		Long: `Verify that contract file matches the locked version.

Returns exit code 0 if match, 1 if mismatch.`,
		RunE: runContractVerify,
	}

	var verifyFeature string
	var verifyContract string

	verifyCmd.Flags().StringVar(&verifyFeature, "feature", "", "Feature name")
	verifyCmd.Flags().StringVar(&verifyContract, "contract", "", "Contract file path")

	cmd.AddCommand(verifyCmd)

	return cmd
}

func runContractSynthesize(cmd *cobra.Command, args []string) error {
	featureName, err := cmd.Flags().GetString("feature")
	if err != nil {
		return fmt.Errorf("failed to get feature flag: %w", err)
	}
	requirementsPath, err := cmd.Flags().GetString("requirements")
	if err != nil {
		return fmt.Errorf("failed to get requirements flag: %w", err)
	}
	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %w", err)
	}

	// Set default requirements path if not provided
	if requirementsPath == "" {
		requirementsPath = fmt.Sprintf("docs/drafts/%s-idea.md", featureName)
	}

	// Set default output path if not provided
	if outputPath == "" {
		outputPath = fmt.Sprintf(".contracts/%s.yaml", featureName)
	}

	fmt.Printf("✓ Generating contract for feature: %s\n", featureName)
	fmt.Printf("  Requirements: %s\n", requirementsPath)
	fmt.Printf("  Output: %s\n", outputPath)
	fmt.Printf("\n⚠️  Contract synthesis not yet implemented\n")
	fmt.Printf("   This will require integration with multi-agent synthesis system\n")

	return nil
}

func runContractGenerate(cmd *cobra.Command, args []string) error {
	featuresFlag, err := cmd.Flags().GetString("features")
	if err != nil {
		return fmt.Errorf("failed to get features flag: %w", err)
	}

	// Parse feature IDs
	var featureIDs []string
	if featuresFlag != "" {
		for _, f := range strings.Split(featuresFlag, ",") {
			featureIDs = append(featureIDs, strings.TrimSpace(f))
		}
	}

	root, err := config.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}

	// Load feature scopes
	featureScopes, err := loadFeatureScopes(root)
	if err != nil {
		return fmt.Errorf("load feature scopes: %w", err)
	}

	// Filter by specified features if provided
	if len(featureIDs) > 0 {
		filtered := make([]collision.FeatureScope, 0)
		for _, fs := range featureScopes {
			for _, fid := range featureIDs {
				if fs.FeatureID == fid {
					filtered = append(filtered, fs)
					break
				}
			}
		}
		featureScopes = filtered
	}

	// Detect boundaries using the collision package
	boundaries := collision.DetectBoundaries(featureScopes)

	if len(boundaries) == 0 {
		fmt.Println("No shared boundaries detected.")
		fmt.Println("  Run 'sdp collision detect' to find shared interfaces.")
		return nil
	}

	// Generate contracts
	contractsDir := filepath.Join(root, ".contracts")
	contracts, err := collision.GenerateContracts(boundaries, contractsDir)
	if err != nil {
		return fmt.Errorf("generate contracts: %w", err)
	}

	fmt.Printf("✓ Generated %d contract(s)\n", len(contracts))
	for _, c := range contracts {
		fmt.Printf("  - %s.yaml (required by: %v)\n", c.TypeName, c.RequiredBy)
	}
	fmt.Printf("\n  Output directory: %s\n", contractsDir)
	fmt.Printf("  Next step: sdp contract lock --contract .contracts/<type>.yaml\n")

	return nil
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
	fmt.Printf("  Expected hash: %s\n", lock.ContractHash[:16]+"...")
	fmt.Printf("  Actual hash:   %s\n", currentHash[:16]+"...")
	fmt.Println()
	fmt.Println("Contract has been modified since lock.")
	fmt.Println("Please re-lock or restore original contract.")
	return false, nil
}
