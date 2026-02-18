package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/sdpinit"
	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var projectType string
	var skipBeads bool
	var autoMode bool
	var headless bool
	var guided bool
	var output string
	var force bool
	var dryRun bool
	var noEvidence bool
	var interactive bool
	var projectName string
	var skills []string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize project with SDP prompts",
		Long: `Initialize current project with SDP prompts and configuration.

Creates .claude/ directory structure:
  skills/     - Claude Code skills
  agents/     - Multi-agent prompts
  validators/ - AI-based quality validators

Modes:
  Interactive (default): Prompts for configuration options
  --auto:                Non-interactive with safe defaults
  --headless:            CI/CD mode with JSON output

Preflight checks detect project type and validate environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Backward-compatible aliases
			if guided {
				interactive = true
			}

			// Run preflight checks
			preflight := sdpinit.RunPreflight()

			// Build config
			cfg := sdpinit.Config{
				ProjectType: projectType,
				SkipBeads:   skipBeads,
				Headless:    headless,
				Output:      output,
				Force:       force,
				DryRun:      dryRun,
				NoEvidence:  noEvidence,
				Interactive: interactive,
				ProjectName: projectName,
				Skills:      skills,
			}

			// Determine project type
			if cfg.ProjectType == "" {
				cfg.ProjectType = preflight.ProjectType
			}

			// Headless mode - JSON output
			if headless {
				return runHeadlessInit(cfg)
			}

			// Auto mode - non-interactive with defaults
			if autoMode {
				return runAutoInit(cfg, preflight)
			}

			// Interactive mode - run wizard
			if interactive || !autoMode {
				return runInteractiveInit(cfg, preflight)
			}

			return nil
		},
	}

	// Basic flags
	cmd.Flags().StringVarP(&projectType, "project-type", "p", "", "Project type (python, go, node, mixed, unknown)")
	cmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	cmd.Flags().BoolVar(&skipBeads, "skip-beads", false, "Skip Beads integration")
	cmd.Flags().StringSliceVar(&skills, "skills", nil, "Skills to enable (comma-separated)")

	// Mode flags
	cmd.Flags().BoolVar(&autoMode, "auto", false, "Non-interactive mode with safe defaults")
	cmd.Flags().BoolVar(&headless, "headless", false, "CI/CD mode with JSON output")
	cmd.Flags().BoolVar(&guided, "guided", false, "Deprecated alias for --interactive")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Force interactive mode")
	cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text, json)")

	// Action flags
	cmd.Flags().BoolVar(&force, "force", false, "Force overwrite existing files")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	cmd.Flags().BoolVar(&noEvidence, "no-evidence", false, "Disable evidence logging")

	return cmd
}

// runHeadlessInit runs initialization in headless mode for CI/CD.
func runHeadlessInit(cfg sdpinit.Config) error {
	output, err := sdpinit.RunHeadless(cfg)

	// Always output JSON (for both success and error cases)
	if jsonErr := output.OutputJSON(); jsonErr != nil {
		fmt.Fprintf(os.Stderr, "Error outputting JSON: %v\n", jsonErr)
	}

	if err != nil {
		// Signal exit code through special error type
		return &headlessError{exitCode: output.GetExitCode(), err: err}
	}

	return nil
}

// headlessError wraps an error with an exit code for headless mode.
type headlessError struct {
	exitCode int
	err      error
}

func (e *headlessError) Error() string {
	return e.err.Error()
}

func (e *headlessError) ExitCode() int {
	return e.exitCode
}

// runAutoInit runs initialization with safe defaults.
func runAutoInit(cfg sdpinit.Config, preflight *sdpinit.PreflightResult) error {
	fmt.Println("SDP Auto Initialization")
	fmt.Println("=======================")
	fmt.Printf("Detected project type: %s\n", preflight.ProjectType)

	if preflight.HasClaude && !cfg.Force {
		fmt.Println("Warning: .claude/ already exists (use --force to overwrite)")
	}
	if preflight.HasSDP {
		fmt.Println("Info: .sdp/ already exists")
	}
	if !preflight.HasGit {
		fmt.Println("Warning: Not a git repository (version control recommended)")
	}

	for _, conflict := range preflight.Conflicts {
		fmt.Printf("Conflict: %s\n", conflict)
	}

	for _, warning := range preflight.Warnings {
		fmt.Printf("Note: %s\n", warning)
	}

	// Get defaults
	defaults := sdpinit.GetDefaults(cfg.ProjectType)
	fmt.Printf("\nUsing safe defaults for %s:\n", cfg.ProjectType)
	fmt.Printf("  Skills: %v\n", defaults.Skills)
	fmt.Printf("  Evidence: %v\n", defaults.EvidenceEnabled && !cfg.NoEvidence)
	fmt.Printf("  Beads: %v\n", !cfg.SkipBeads)

	if cfg.DryRun {
		fmt.Println("\n[DRY RUN] Would create:")
		fmt.Println("  .claude/")
		fmt.Println("  .claude/skills/")
		fmt.Println("  .claude/agents/")
		fmt.Println("  .claude/validators/")
		fmt.Println("  .claude/settings.json")
		return nil
	}

	fmt.Printf("\nInitializing with project type: %s\n", cfg.ProjectType)
	return sdpinit.Run(cfg)
}

// runInteractiveInit runs the interactive wizard.
func runInteractiveInit(cfg sdpinit.Config, preflight *sdpinit.PreflightResult) error {
	// Run the wizard
	answers, err := sdpinit.RunWizard(preflight)
	if err != nil {
		return fmt.Errorf("wizard: %w", err)
	}

	// Apply wizard answers to config
	cfg.ProjectName = answers.ProjectName
	cfg.ProjectType = answers.ProjectType
	cfg.Skills = answers.Skills
	cfg.NoEvidence = answers.NoEvidence
	cfg.SkipBeads = answers.SkipBeads

	// Run initialization
	return sdpinit.Run(cfg)
}
