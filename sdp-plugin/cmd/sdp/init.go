package main

import (
	"os"

	"github.com/fall-out-bug/sdp/internal/sdpinit"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var projectType string
	var skipBeads bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize project with SDP prompts",
		Long: `Initialize current project with SDP prompts and configuration.

Creates .claude/ directory structure:
  skills/     - Claude Code skills
  agents/     - Multi-agent prompts
  validators/ - AI-based quality validators

Automatically detects project type (python, java, go) or use --project-type flag.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Detect project type if not specified
			if projectType == "" {
				projectType = detectProjectType()
			}

			cfg := sdpinit.Config{
				ProjectType: projectType,
				SkipBeads:   skipBeads,
			}

			return sdpinit.Run(cfg)
		},
	}

	cmd.Flags().StringVarP(&projectType, "project-type", "p", "", "Project type (python, java, go, agnostic)")
	cmd.Flags().BoolVar(&skipBeads, "skip-beads", false, "Skip Beads integration")

	return cmd
}

func detectProjectType() string {
	// Check for build files
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return "python"
	}
	if _, err := os.Stat("pom.xml"); err == nil {
		return "java"
	}
	if _, err := os.Stat("build.gradle"); err == nil {
		return "java"
	}
	if _, err := os.Stat("go.mod"); err == nil {
		return "go"
	}
	return "agnostic"
}
