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
var telemetryConsentSkipCommands = map[string]struct{}{
	"completion": {},
	"demo":       {},
	"doctor":     {},
	"init":       {},
	"next":       {},
	"status":     {},
}

func shouldAskForTelemetryConsent(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}

	if _, skip := telemetryConsentSkipCommands[cmd.Name()]; skip {
		return false
	}

	for _, flagName := range []string{"auto", "headless"} {
		flag := cmd.Flags().Lookup(flagName)
		if flag != nil && flag.Value.String() == "true" {
			return false
		}
	}

	return true
}

func main() {
	var noColor bool

	var rootCmd = &cobra.Command{
		Use:   "sdp",
		Short: "Spec-Driven Protocol - AI workflow tools",
		Long: `SDP provides convenience commands for Spec-Driven Protocol:

	  init       Initialize project with SDP prompts
	  doctor     Check environment and setup health
	  plan       Decompose a feature into workstreams
	  apply      Execute ready workstreams from the terminal
	  status     Show current project state
	  next       Recommend the next action to take
	  log        Inspect the evidence log
	  demo       Run a guided first-success walkthrough
	  guard      Enforce workstream scope and context
	  completion Generate shell completion script

These commands are optional convenience tools. The core SDP functionality
is provided by the prompts installed into your supported IDE integration
directory (.claude/, .cursor/, .opencode/, or .codex/).`,
		Example: `  # Initialize SDP in a project
  sdp init .

	  # Plan and execute a feature from the CLI
	  sdp plan "Add auth"
	  sdp apply

	  # Inspect state and get the next action
	  sdp status --text
	  sdp next --json

	  # Check environment setup
  sdp doctor

	  # Run a guided demo
	  sdp demo

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
			if !consentAsked && cmd.Name() != "telemetry" && shouldAskForTelemetryConsent(cmd) {
				configDir, err := os.UserConfigDir()
				if err == nil {
					configPath := filepath.Join(configDir, "sdp", "telemetry.json")
					if telemetry.IsFirstRun(configPath) {
						// Ask for consent on first run
						granted, err := telemetry.AskForConsent()
						if err == nil {
							// Save user's choice
							func() {
								if cerr := telemetry.GrantConsent(configPath, granted); cerr != nil {
									fmt.Fprintf(os.Stderr, "warning: failed to save telemetry consent: %v\n", cerr)
								}
							}()
							consentAsked = true
						}
					}
				}
			}

			// Track command start (skip telemetry commands to avoid infinite loops)
			if cmd.Parent() == nil || cmd.Parent().Use != "telemetry" {
				func() {
					if cerr := telemetry.TrackCommandStart(cmd.Name(), args); cerr != nil {
						fmt.Fprintf(os.Stderr, "warning: failed to track command start: %v\n", cerr)
					}
				}()
			}
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// Track command completion (skip telemetry commands)
			if cmd.Parent() == nil || cmd.Parent().Use != "telemetry" {
				func() {
					if cerr := telemetry.TrackCommandComplete(true, ""); cerr != nil {
						fmt.Fprintf(os.Stderr, "warning: failed to track command completion: %v\n", cerr)
					}
				}()
			}
			return nil
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(doctorCmd())
	rootCmd.AddCommand(assessCmd())
	rootCmd.AddCommand(tryCmd())
	rootCmd.AddCommand(adoptCmd())
	rootCmd.AddCommand(hooksCmd())
	rootCmd.AddCommand(guardCmd())
	rootCmd.AddCommand(collisionCmd())
	rootCmd.AddCommand(acceptanceCmd())
	rootCmd.AddCommand(logCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(deployCmd())
	rootCmd.AddCommand(prdCmd())
	rootCmd.AddCommand(skillCmd())
	rootCmd.AddCommand(parseCmd())
	rootCmd.AddCommand(designCmd())
	rootCmd.AddCommand(ideaCmd())
	rootCmd.AddCommand(prototypeCmd())
	rootCmd.AddCommand(planCmd())
	rootCmd.AddCommand(applyCmd())
	rootCmd.AddCommand(beadsCmd())
	rootCmd.AddCommand(buildCmd())
	rootCmd.AddCommand(tddCmd())
	rootCmd.AddCommand(driftCmd())
	rootCmd.AddCommand(qualityCmd())
	rootCmd.AddCommand(watchCmd())
	rootCmd.AddCommand(telemetryCmd)
	rootCmd.AddCommand(checkpointCmd)
	rootCmd.AddCommand(contractCmd())
	rootCmd.AddCommand(orchestrateCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(decisionsCmd())
	rootCmd.AddCommand(metricsCmd())
	rootCmd.AddCommand(resolveCmd())
	rootCmd.AddCommand(taskCmd())
	rootCmd.AddCommand(sessionCmd())
	rootCmd.AddCommand(gitCmd())
	rootCmd.AddCommand(memoryCmd())
	rootCmd.AddCommand(coordinationCmd())
	rootCmd.AddCommand(healthCmd())
	rootCmd.AddCommand(diagnoseCmd)
	rootCmd.AddCommand(nextCmd())
	rootCmd.AddCommand(demoCmd())

	if err := rootCmd.Execute(); err != nil {
		// Track command failure
		func() {
			if cerr := telemetry.TrackCommandComplete(false, err.Error()); cerr != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to track command failure: %v\n", cerr)
			}
		}()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
