package cli

import (
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
	synthesizeCmd.Flags().StringVarP(&synthesizeFeature, "feature", "f", "", "Feature name (required)")
	synthesizeCmd.Flags().StringVar(&synthesizeRequirements, "requirements", "", "Path to requirements file")
	synthesizeCmd.Flags().StringVar(&synthesizeOutput, "output", "", "Output path for contract")
	synthesizeCmd.MarkFlagRequired("feature")

	lockCmd.Flags().StringVarP(&lockContract, "contract", "c", "", "Contract file path (required)")
	lockCmd.Flags().StringVar(&lockReason, "reason", "", "Reason for locking")
	lockCmd.MarkFlagRequired("contract")

	validateCmd.Flags().StringSliceVarP(&validateContracts, "contracts", "c", []string{}, "Contract files to validate")
	validateCmd.Flags().StringVar(&validateOutput, "output", ".contracts/validation-report.md", "Output path for report")
	validateCmd.MarkFlagRequired("contracts")

	contractCmd.AddCommand(synthesizeCmd)
	contractCmd.AddCommand(lockCmd)
	contractCmd.AddCommand(validateCmd)
}

// RegisterContractCommand registers the contract command with root
func RegisterContractCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(contractCmd)
}
