package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/skill"
	"github.com/spf13/cobra"
)

func skillShow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <skill-name>",
		Short: "Show skill file content",
		Long:  `Display the full content of a skill file (SKILL.md)`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("requires skill name argument")
			}
			skillName := args[0]
			skillsDir, _ := cmd.Flags().GetString("skills-dir") //nolint:errcheck // String flag never errors

			content, err := skill.ReadSkillContent(skillsDir, skillName)
			if err != nil {
				return fmt.Errorf("failed to read skill: %w", err)
			}

			fmt.Println(content)
			return nil
		},
	}

	return cmd
}
