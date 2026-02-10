package main

import (
	"github.com/spf13/cobra"
)

func skillCmd() *cobra.Command {
	var skillsDir string

	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Skill management commands",
		Long: `Skill management operations for validating and listing
Claude Code skills.

Subcommands:
  validate    - Validate a skill file against standards
  check-all   - Validate all skills in .claude/skills/
  list        - List all available skills
  show        - Show skill file content`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if skillsDir == "" {
				skillsDir = ".claude/skills"
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&skillsDir, "skills-dir", "", "Skills directory (default: .claude/skills)")

	cmd.AddCommand(skillValidate())
	cmd.AddCommand(skillCheckAll())
	cmd.AddCommand(skillList())
	cmd.AddCommand(skillShow())
	cmd.AddCommand(skillRecord())

	return cmd
}
