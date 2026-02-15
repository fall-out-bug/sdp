package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/sdpinit"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var projectType string
	var skipBeads bool
	var autoMode bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize project with SDP prompts",
		Long: `Initialize current project with SDP prompts and configuration.

Creates .claude/ directory structure:
  skills/     - Claude Code skills
  agents/     - Multi-agent prompts
  validators/ - AI-based quality validators

Preflight checks detect project type and validate environment.
Use --auto for non-interactive initialization with safe defaults.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run preflight checks
			preflight := sdpinit.RunPreflight()

			// Display preflight results
			fmt.Println("SDP Initialization")
			fmt.Println("==================")
			fmt.Printf("Detected project type: %s\n", preflight.ProjectType)

			if preflight.HasClaude {
				fmt.Println("⚠ .claude/ already exists - will update")
			}
			if preflight.HasSDP {
				fmt.Println("⚠ .sdp/ already exists")
			}
			if !preflight.HasGit {
				fmt.Println("⚠ Not a git repository")
			}

			for _, conflict := range preflight.Conflicts {
				fmt.Printf("⚠ Conflict: %s\n", conflict)
			}

			for _, warning := range preflight.Warnings {
				fmt.Printf("  Note: %s\n", warning)
			}

			// Determine project type
			if projectType == "" {
				projectType = preflight.ProjectType
			}

			// In auto mode, proceed without prompts
			if autoMode {
				fmt.Println("\nAuto mode: Using safe defaults")
			}

			cfg := sdpinit.Config{
				ProjectType: projectType,
				SkipBeads:   skipBeads,
			}

			fmt.Printf("\nInitializing with project type: %s\n", projectType)
			return sdpinit.Run(cfg)
		},
	}

	cmd.Flags().StringVarP(&projectType, "project-type", "p", "", "Project type (python, java, go, node, mixed, unknown)")
	cmd.Flags().BoolVar(&skipBeads, "skip-beads", false, "Skip Beads integration")
	cmd.Flags().BoolVar(&autoMode, "auto", false, "Non-interactive mode with safe defaults")

	return cmd
}
