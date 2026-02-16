package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardActivate() *cobra.Command {
	return &cobra.Command{
		Use:   "activate <ws-id>",
		Short: "Activate workstream for editing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsID := args[0]

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

			// Activate workstream
			if err := skill.Activate(wsID); err != nil {
				return fmt.Errorf("failed to activate WS: %w", err)
			}

			activeWS := skill.GetActiveWS()
			fmt.Printf("Activated WS: %s\n", activeWS)

			if evidence.Enabled() {
				scopeFiles := scopeFilesForWS(wsID)
				if err := evidence.EmitSync(evidence.PlanEvent(wsID, scopeFiles)); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "warning: evidence emit: %v\n", err)
				}
			}

			// AC1: Check for scope overlap with other in-progress workstreams (warning only)
			warnCollisionIfAny(wsID)

			// AC3/AC4: Check for similar past failed decisions (warning only)
			warnSimilarFailures(wsID)

			// AC2: Contract validation after code generation (warning in P1)
			warnContractViolations()

			return nil
		},
	}
}
