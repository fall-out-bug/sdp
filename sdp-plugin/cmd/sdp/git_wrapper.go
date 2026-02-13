package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/git"
	"github.com/spf13/cobra"
)

func gitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Git wrapper with session validation",
		Long: `Execute git commands with automatic session validation.

This wrapper ensures that git operations are performed in the correct
worktree and branch by validating the session before execution.

Commands are categorized as:
- Safe: status, log, diff, show (read-only)
- Write: add, commit, reset (modifies repository)
- Remote: push, fetch, pull (interacts with remotes)
- Branch: checkout, merge (changes branches)

All commands require session validation, but only write/branch commands
perform post-execution checks to ensure branch hasn't changed.`,
		Example: `  # Safe commands
  sdp git status
  sdp git log --oneline -10

  # Write commands (with post-check)
  sdp git add .
  sdp git commit -m "feat: add new feature"

  # Remote commands
  sdp git push origin feature/F065

  # Branch commands (with post-check)
  sdp git checkout feature/F065`,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("git command required\nUsage: sdp git <command> [args...]")
			}

			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			wrapper := git.NewWrapper(root)

			gitCmd := args[0]
			gitArgs := args[1:]

			// Special handling for --help
			if gitCmd == "--help" || gitCmd == "-h" {
				cmd.Help()
				return nil
			}

			// Check if session exists
			if !wrapper.HasSession() {
				return fmt.Errorf(`no session found. Initialize one first:

  sdp session init --feature=F###`)
			}

			// Execute the git command
			if err := wrapper.Execute(gitCmd, gitArgs...); err != nil {
				// Print error and exit with non-zero code
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	return cmd
}
