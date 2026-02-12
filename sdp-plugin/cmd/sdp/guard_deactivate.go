package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardDeactivate() *cobra.Command {
	return &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate guard",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get config directory (respect XDG_CONFIG_HOME for testing)
			configDir := os.Getenv("XDG_CONFIG_HOME")
			if configDir == "" {
				var err error
				configDir, err = os.UserConfigDir()
				if err != nil {
					return fmt.Errorf("failed to get config dir: %w", err)
				}
			}

			sdpDir := filepath.Join(configDir, "sdp")
			skill := guard.NewSkill(sdpDir)

			// Deactivate
			if err := skill.Deactivate(); err != nil {
				return fmt.Errorf("failed to deactivate: %w", err)
			}

			fmt.Println("Guard deactivated")

			return nil
		},
	}
}
