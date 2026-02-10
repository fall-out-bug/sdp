package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/skill"
	"github.com/spf13/cobra"
)

func skillValidate() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:   "validate <skill-file>",
		Short: "Validate a skill file against standards",
		Long: `Validate a skill file against SDP standards.

Checks:
- Line count ≤150 (warning if >100)
- Required sections present
- Frontmatter starts with ---
- References resolve`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires skill file argument")
			}
			skillPath := args[0]
			validator := skill.NewValidator()

			result, err := validator.ValidateFile(skillPath)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			if len(result.Errors) > 0 {
				fmt.Printf("❌ %s: %d errors\n", skillPath, len(result.Errors))
				for _, e := range result.Errors {
					fmt.Printf("   - %s\n", e)
				}
			}

			if len(result.Warnings) > 0 {
				fmt.Printf("⚠️  %s: %d warnings\n", skillPath, len(result.Warnings))
				for _, w := range result.Warnings {
					fmt.Printf("   - %s\n", w)
				}
			}

			if result.IsValid {
				fmt.Printf("✅ %s: valid (%d lines)\n", skillPath, result.LineCount)
			}

			if !result.IsValid || (strict && len(result.Warnings) > 0) {
				return fmt.Errorf("skill validation failed")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "Fail on warnings")

	return cmd
}
