package main

import (
	"github.com/spf13/cobra"
	"github.com/ai-masters/sdp/internal/hooks"
)

func hooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hooks",
		Short: "Manage Git hooks for SDP",
		Long: `Install or uninstall Git hooks for SDP quality checks.

Hooks are scripts that run automatically during Git operations:
  - pre-commit: Runs before each commit
  - pre-push: Runs before each push

You can customize hooks in .git/hooks/ after installation.`,
	}

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install Git hooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return hooks.Install()
		},
	}

	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Git hooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return hooks.Uninstall()
		},
	}

	cmd.AddCommand(installCmd)
	cmd.AddCommand(uninstallCmd)

	return cmd
}
