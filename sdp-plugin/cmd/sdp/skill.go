package main

import (
	"fmt"
	"strings"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/skill"
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

func skillRecord() *cobra.Command {
	var skillName, eventType, wsID string
	var dataPairs []string

	cmd := &cobra.Command{
		Use:   "record",
		Short: "Record skill execution in evidence log (F056)",
		Long:  `Emit a thin evidence event for a skill. Non-blocking; used by @vision, @reality, etc.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if skillName == "" || eventType == "" {
				return fmt.Errorf("--skill and --type required")
			}
			if wsID == "" {
				wsID = "00-000-00"
			}
			data := make(map[string]interface{})
			for _, p := range dataPairs {
				kv := strings.SplitN(p, "=", 2)
				if len(kv) == 2 {
					data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			}
			_ = evidence.EmitSync(evidence.SkillEvent(skillName, eventType, wsID, data))
			fmt.Printf("Recorded: %s %s\n", skillName, eventType)
			return nil
		},
	}
	cmd.Flags().StringVar(&skillName, "skill", "", "Skill name (vision, reality, oneshot, prototype, hotfix, bugfix, issue, debug)")
	cmd.Flags().StringVar(&eventType, "type", "", "Event type (plan, verification, generation, approval)")
	cmd.Flags().StringVar(&wsID, "ws-id", "00-000-00", "WS ID for event")
	cmd.Flags().StringArrayVar(&dataPairs, "data", nil, "key=value (repeatable)")
	return cmd
}

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

			// Output results
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

			// Exit with error if not valid or strict mode with warnings
			if !result.IsValid || (strict && len(result.Warnings) > 0) {
				return fmt.Errorf("skill validation failed")
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "Fail on warnings")

	return cmd
}

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
