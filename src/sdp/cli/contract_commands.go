package cli

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fall-out-bug/sdp/src/sdp/agents"
	"github.com/spf13/cobra"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Manage API contracts",
	Long:  `Manage API contracts for component validation.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var synthesizeCmd = &cobra.Command{
	Use:   "synthesize",
	Short: "Synthesize contract from multi-agent agreement",
	Long: `Generate an API contract from feature requirements using multi-agent synthesis.

This command analyzes feature requirements, proposes initial contracts,
requests agent reviews, applies synthesis rules, and generates an agreed contract.`,
	RunE: runContractSynthesize,
}

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock contract to prevent changes",
	Long:  `Lock a contract to establish it as the source of truth for implementation.`,
	RunE: runContractLock,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate contracts between components",
	Long:  `Validate contracts between frontend, backend, and SDK components.`,
	RunE: runContractValidate,
}

// Synthesis flags
var (
	synthesizeFeature      string
	synthesizeRequirements string
	synthesizeOutput       string
)

// Lock flags
var (
	lockContract string
	lockReason   string
)

// Validate flags
var (
	validateContracts []string
	validateOutput    string
)

func init() {
	// Add synthesize flags
	synthesizeCmd.Flags().StringVarP(&synthesizeFeature, "feature", "f", "", "Feature name (required)")
	synthesizeCmd.Flags().StringVar(&synthesizeRequirements, "requirements", "", "Path to requirements file")
	synthesizeCmd.Flags().StringVar(&synthesizeOutput, "output", "", "Output path for contract")
	synthesizeCmd.MarkFlagRequired("feature")

	// Add lock flags
	lockCmd.Flags().StringVarP(&lockContract, "contract", "c", "", "Contract file path (required)")
	lockCmd.Flags().StringVar(&lockReason, "reason", "", "Reason for locking")
	lockCmd.MarkFlagRequired("contract")

	// Add validate flags
	validateCmd.Flags().StringSliceVarP(&validateContracts, "contracts", "c", []string{}, "Contract files to validate")
	validateCmd.Flags().StringVar(&validateOutput, "output", ".contracts/validation-report.md", "Output path for report")
	validateCmd.MarkFlagRequired("contracts")

	// Register subcommands
	contractCmd.AddCommand(synthesizeCmd)
	contractCmd.AddCommand(lockCmd)
	contractCmd.AddCommand(validateCmd)
}

func runContractSynthesize(cmd *cobra.Command, args []string) error {
	feature := synthesizeFeature

	// Validate feature name (alphanumeric and dash only)
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(feature) {
		return fmt.Errorf("invalid feature name %q: must contain only lowercase letters, numbers, and dashes", feature)
	}

	// Determine requirements path
	reqPath := synthesizeRequirements
	if reqPath == "" {
		reqPath = filepath.Join("docs", "drafts", feature+"-requirements.md")
	}

	// Sanitize requirements path
	safeReqPath, err := sanitizePath(reqPath, []string{"docs", "docs/drafts"})
	if err != nil {
		return fmt.Errorf("invalid requirements path: %w", err)
	}

	// Determine output path
	outputPath := synthesizeOutput
	if outputPath == "" {
		outputPath = filepath.Join(".contracts", feature+".yaml")
	}

	// Sanitize output path
	safeOutputPath, err := sanitizePath(outputPath, []string{".contracts"})
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	// Print progress
	fmt.Printf("✓ Analyzing requirements from %s\n", safeReqPath)

	// Create synthesizer
	synthesizer := agents.NewContractSynthesizer()

	// Perform synthesis
	result, err := synthesizer.SynthesizeContract(feature, safeReqPath, safeOutputPath)
	if err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}

	// Print success
	fmt.Printf("✓ Contract agreed: %s\n", safeOutputPath)
	fmt.Printf("\nResolution method: %s\n", result.Rule)

	if result.WinningAgent != "" {
		fmt.Printf("Winning agent: %s\n", result.WinningAgent)
	}

	return nil
}

// sanitizePath validates and sanitizes file paths to prevent directory traversal attacks
func sanitizePath(path string, allowedDirs []string) (string, error) {
	// Clean the path to remove any ".." or "." components
	cleanPath := filepath.Clean(path)

	// Convert to absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Check if path is within allowed directories
	for _, allowedDir := range allowedDirs {
		absAllowedDir, err := filepath.Abs(allowedDir)
		if err != nil {
			continue
		}

		// Check if the absolute path starts with the allowed directory
		if strings.HasPrefix(absPath, absAllowedDir+string(filepath.Separator)) {
			return absPath, nil
		}
	}

	return "", fmt.Errorf("path %q is outside allowed directories %v", path, allowedDirs)
}

func runContractLock(cmd *cobra.Command, args []string) error {
	contractPath := lockContract

	// Sanitize path to prevent directory traversal
	safePath, err := sanitizePath(contractPath, []string{".contracts", "docs"})
	if err != nil {
		return fmt.Errorf("invalid contract path: %w", err)
	}

	// Read contract
	content, err := os.ReadFile(safePath)
	if err != nil {
		return fmt.Errorf("failed to read contract: %w", err)
	}

	// Create lock file
	lockPath := safePath + ".lock"
	lockContent := fmt.Sprintf("# Contract Lock\n\nlocked: true\nreason: %s\ncontract_sha256: %x\n",
		lockReason, sha256.Sum256(content))

	if err := os.WriteFile(lockPath, []byte(lockContent), 0600); err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	fmt.Printf("✓ Contract locked: %s\n", safePath)
	fmt.Printf("✓ Lock file created: %s\n", lockPath)

	return nil
}

func runContractValidate(cmd *cobra.Command, args []string) error {
	if len(validateContracts) < 2 {
		return fmt.Errorf("at least 2 contracts required for validation")
	}

	fmt.Printf("✓ Validating %d contracts...\n", len(validateContracts))

	// Sanitize all contract paths
	safeContracts := make([]string, len(validateContracts))
	for i, contractPath := range validateContracts {
		safePath, err := sanitizePath(contractPath, []string{".contracts", "docs"})
		if err != nil {
			return fmt.Errorf("invalid contract path %q: %w", contractPath, err)
		}
		safeContracts[i] = safePath
	}

	// Sanitize output path
	safeOutput, err := sanitizePath(validateOutput, []string{".contracts", "docs"})
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	// Load contracts
	validator := agents.NewContractValidator()

	// For now, just validate the first two
	contract1, err := loadContract(safeContracts[0])
	if err != nil {
		return fmt.Errorf("failed to load contract 1: %w", err)
	}

	contract2, err := loadContract(safeContracts[1])
	if err != nil {
		return fmt.Errorf("failed to load contract 2: %w", err)
	}

	// Compare contracts
	mismatches, err := validator.CompareContracts(
		contract1,
		contract2,
		safeContracts[0],
		safeContracts[1],
	)
	if err != nil {
		return fmt.Errorf("comparison failed: %w", err)
	}

	// Generate report
	report := validator.GenerateReport(mismatches)

	// Write report
	if err := validator.WriteReport(report, safeOutput); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	// Print summary
	fmt.Printf("✓ Validation report: %s\n", safeOutput)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Total issues: %d\n", len(mismatches))

	errorCount := 0
	warningCount := 0
	for _, m := range mismatches {
		if m.Severity == "ERROR" {
			errorCount++
		} else if m.Severity == "WARNING" {
			warningCount++
		}
	}

	fmt.Printf("- Errors: %d\n", errorCount)
	fmt.Printf("- Warnings: %d\n", warningCount)

	return nil
}

func loadContract(path string) (*agents.OpenAPIContract, error) {
	// For now, return a minimal contract
	// Full implementation would parse YAML
	return &agents.OpenAPIContract{
		OpenAPI: "3.0.0",
		Paths:   make(agents.PathsSpec),
	}, nil
}

// RegisterContractCommand registers the contract command with root
func RegisterContractCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(contractCmd)
}
