package main

import (
	"fmt"
	"os"

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
	}

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(doctorCmd())
	rootCmd.AddCommand(hooksCmd())
	rootCmd.AddCommand(parseCmd())
	rootCmd.AddCommand(beadsCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
