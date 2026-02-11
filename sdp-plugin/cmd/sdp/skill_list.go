package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/skill"
	"github.com/spf13/cobra"
)

func skillList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available skills",
		Long:  `List all skill directories found in .claude/skills/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			skillsDir, _ := cmd.Flags().GetString("skills-dir") //nolint:errcheck // String flag never errors

			skills, err := skill.ListSkills(skillsDir)
			if err != nil {
				return fmt.Errorf("failed to list skills: %w", err)
			}

			if len(skills) == 0 {
				fmt.Println("No skills found")
				return nil
			}

			fmt.Printf("Found %d skills:\n", len(skills))
			for _, s := range skills {
				fmt.Printf("  - %s\n", s)
			}

			return nil
		},
	}

	return cmd
}
