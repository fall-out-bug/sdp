package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/skill"
	"github.com/spf13/cobra"
)

func skillCheckAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-all",
		Short: "Validate all skills in .claude/skills/",
		Long: `Validate all skill files in the .claude/skills/ directory
against SDP standards.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			skillsDir, _ := cmd.Flags().GetString("skills-dir") //nolint:errcheck // String flag never errors
			validator := skill.NewValidator()

			results, err := validator.ValidateAll(skillsDir)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			total := len(results)
			failed := 0

			for skillName, result := range results {
				if len(result.Errors) > 0 {
					fmt.Printf("❌ %s: %d errors\n", skillName, len(result.Errors))
					for _, e := range result.Errors {
						fmt.Printf("   - %s\n", e)
					}
					failed++
				}

				if len(result.Warnings) > 0 {
					fmt.Printf("⚠️  %s: %d warnings\n", skillName, len(result.Warnings))
					for _, w := range result.Warnings {
						fmt.Printf("   - %s\n", w)
					}
				}

				if result.IsValid {
					fmt.Printf("✅ %s: valid (%d lines)\n", skillName, result.LineCount)
				}
			}

			fmt.Printf("\nSummary: %d/%d skills valid\n", total-failed, total)
			if failed > 0 {
				return fmt.Errorf("skill validation failed")
			}

			return nil
		},
	}

	cmd.Flags().String("skills-dir", "", "Skills directory")
	return cmd
}
