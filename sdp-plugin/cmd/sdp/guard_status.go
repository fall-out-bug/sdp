package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardStatus() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show guard status",
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

			activeWS := skill.GetActiveWS()
			if activeWS == "" {
				fmt.Println("Guard Status: INACTIVE")
				fmt.Println("No active workstream")
				return nil
			}

			fmt.Printf("Guard Status: ACTIVE\n")
			fmt.Printf("Active WS: %s\n", activeWS)

			// Load state to show scope files
			state, err := guard.NewStateManager(sdpDir).Load()
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			if len(state.ScopeFiles) > 0 {
				fmt.Println("Scope files:")
				for _, f := range state.ScopeFiles {
					fmt.Printf("  - %s\n", f)
				}
			} else {
				fmt.Println("Scope: No restrictions")
			}

			// Show review findings
			if len(state.ReviewFindings) > 0 {
				open, resolved, blocking := state.FindingCount()
				fmt.Printf("\nReview Findings: %d open (%d blocking), %d resolved\n", open, blocking, resolved)

				// Show blocking findings first
				if blocking > 0 {
					fmt.Println("\n⚠️  BLOCKING FINDINGS (must resolve before merge):")
					for _, f := range state.GetBlockingFindings() {
						fmt.Printf("  [%s] P%d %s\n", f.ReviewArea, f.Priority, f.Title)
						if f.BeadsID != "" {
							fmt.Printf("    → Beads: %s\n", f.BeadsID)
						}
					}
				}
			}

			return nil
		},
	}
}
