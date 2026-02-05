package main

import (
	"fmt"
	"os"

	"github.com/ai-masters/sdp/internal/telemetry"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sdp",
		Short: "Spec-Driven Protocol - AI workflow tools",
		Long: `SDP provides convenience commands for Spec-Driven Protocol:

  init     Initialize project with SDP prompts
  doctor   Check environment (Git, Claude Code, .claude/)
  hooks    Manage Git hooks for SDP

These commands are optional convenience tools. The core SDP functionality
is provided by the Claude Plugin prompts in .claude/.`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(doctorCmd())
	rootCmd.AddCommand(hooksCmd())
	rootCmd.AddCommand(parseCmd())
	rootCmd.AddCommand(beadsCmd())
	rootCmd.AddCommand(tddCmd())
	rootCmd.AddCommand(driftCmd())
	rootCmd.AddCommand(qualityCmd())
	rootCmd.AddCommand(telemetryCmd)

	if err := rootCmd.Execute(); err != nil {
		// Track command failure
		telemetry.TrackCommandComplete(false, err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
