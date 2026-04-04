package main

import (
	"context"

	"github.com/spf13/cobra"
)

var (
	runQualityCoverageCmd   = runQualityCoverage
	runQualityComplexityCmd = runQualityComplexity
	runQualitySizeCmd       = runQualitySize
	runQualityTypesCmd      = runQualityTypes
	runQualityAllCmd        = runQualityAll
)

func qualityCmd() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:   "quality",
		Short: "Run quality checks on the project",
		Long: `Run quality checks on the project.

Checks include:
  coverage   - Test coverage analysis (≥80% required)
  complexity - Cyclomatic complexity analysis (<10 required)
  size       - File size analysis (<200 LOC required)
  types      - Type checking (mypy, go vet, etc.)
  all        - Run all quality checks

Pragmatic Mode (default):
  File size violations are WARNINGS (build continues)

Strict Mode (--strict):
  File size violations are ERRORS (build fails)`,
	}

	// Add persistent flag for strict mode (applies to all subcommands)
	cmd.PersistentFlags().BoolVar(&strict, "strict", false, "Enable strict quality gates (file size violations = errors)")

	cmd.AddCommand(&cobra.Command{
		Use:   "coverage",
		Short: "Check test coverage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualityCoverageCmd(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "complexity",
		Short: "Check cyclomatic complexity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualityComplexityCmd(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "size",
		Short: "Check file sizes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualitySizeCmd(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "types",
		Short: "Check type completeness",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualityTypesCmd(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "Run all quality checks",
		RunE: func(c *cobra.Command, args []string) error {
			ctx := c.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			return runQualityAllCmd(ctx, strict)
		},
	})

	return cmd
}
