package main

import (
	"github.com/spf13/cobra"
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
			return runQualityCoverage(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "complexity",
		Short: "Check cyclomatic complexity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualityComplexity(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "size",
		Short: "Check file sizes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualitySize(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "types",
		Short: "Check type completeness",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualityTypes(strict)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "Run all quality checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQualityAll(strict)
		},
	})

	return cmd
}
