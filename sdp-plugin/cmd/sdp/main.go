package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/telemetry"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

var version = "dev"

var consentAsked = false // Track if we've asked for consent this session

func main() {
	var noColor bool

	var rootCmd = &cobra.Command{
		Use:   "sdp",
		Short: "Spec-Driven Protocol - AI workflow tools",
		Long: `SDP provides convenience commands for Spec-Driven Protocol:

  init       Initialize project with SDP prompts
  doctor     Check environment (Git, Claude Code, .claude/)
  hooks      Manage Git hooks for SDP
  watch      Watch files for quality violations
  checkpoint Manage checkpoints for long-running features
  completion Generate shell completion script

These commands are optional convenience tools. The core SDP functionality
is provided by the Claude Plugin prompts in .claude/.`,
		Example: `  # Initialize SDP in a project
  sdp init .

  # Check environment setup
  sdp doctor

  # Generate shell completion
  sdp completion bash > ~/.bash_completion.d/sdp
  sdp completion zsh > ~/.zsh/completion/_sdp

  # Create a checkpoint
  sdp checkpoint create my-feature F042

  # List checkpoints
  sdp checkpoint list`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set NoColor flag
			ui.NoColor = noColor

			// Check for first-run consent (only once per session)
			if !consentAsked && cmd.Name() != "telemetry" {
				configDir, err := os.UserConfigDir()
				if err == nil {
					configPath := filepath.Join(configDir, "sdp", "telemetry.json")
					if telemetry.IsFirstRun(configPath) {
						// Ask for consent on first run
						granted, err := telemetry.AskForConsent()
						if err == nil {
							// Save user's choice
							telemetry.GrantConsent(configPath, granted)
							consentAsked = true
						}
					}
				}
			}

			// Track command start (skip telemetry commands to avoid infinite loops)
			if cmd.Parent() == nil || cmd.Parent().Use != "telemetry" {
				telemetry.TrackCommandStart(cmd.Name(), args)
			}
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// Track command completion (skip telemetry commands)
			if cmd.Parent() == nil || cmd.Parent().Use != "telemetry" {
				telemetry.TrackCommandComplete(true, "")
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(doctorCmd())
	rootCmd.AddCommand(hooksCmd())
	rootCmd.AddCommand(guardCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(prdCmd())
	rootCmd.AddCommand(skillCmd())
	rootCmd.AddCommand(parseCmd())
	rootCmd.AddCommand(beadsCmd())
	rootCmd.AddCommand(tddCmd())
	rootCmd.AddCommand(driftCmd())
	rootCmd.AddCommand(qualityCmd())
	rootCmd.AddCommand(watchCmd())
	rootCmd.AddCommand(telemetryCmd)
	rootCmd.AddCommand(checkpointCmd)
	rootCmd.AddCommand(orchestrateCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(statusCmd())

	if err := rootCmd.Execute(); err != nil {
		// Track command failure
		telemetry.TrackCommandComplete(false, err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
