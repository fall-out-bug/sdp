package main

import (
	"github.com/spf13/cobra"
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
