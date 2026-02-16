package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardFindingResolveCmd() *cobra.Command {
	var resolvedBy string

	cmd := &cobra.Command{
		Use:   "resolve <finding-id>",
		Short: "Mark a finding as resolved",
		Long: `Mark a review finding as resolved.

The resolvedBy should describe how the finding was addressed.`,
		Example: `  sdp guard finding resolve finding-123 --by="Fixed in commit abc123"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			sm := guard.NewStateManager(configDir)

			state, err := sm.Load()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			if !state.ResolveFinding(args[0], resolvedBy) {
				return fmt.Errorf("finding not found: %s", args[0])
			}

			if err := sm.Save(*state); err != nil {
				return fmt.Errorf("failed to save state: %w", err)
			}

			fmt.Printf("✓ Resolved finding: %s\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringVar(&resolvedBy, "by", "manual", "How the finding was resolved")
	cmd.MarkFlagRequired("by")

	return cmd
}

func guardFindingClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all resolved findings",
		Long:  `Remove all resolved findings from the state file.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := getConfigDir()
			sm := guard.NewStateManager(configDir)

			state, err := sm.Load()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			var open []guard.ReviewFinding
			for _, f := range state.ReviewFindings {
				if f.Status != "resolved" {
					open = append(open, f)
				}
			}

			cleared := len(state.ReviewFindings) - len(open)
			state.ReviewFindings = open

			if err := sm.Save(*state); err != nil {
				return fmt.Errorf("failed to save state: %w", err)
			}

			fmt.Printf("✓ Cleared %d resolved findings\n", cleared)
			return nil
		},
	}

	return cmd
}

func getConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir, _ = os.UserConfigDir()
	}
	return filepath.Join(configDir, "sdp")
}
