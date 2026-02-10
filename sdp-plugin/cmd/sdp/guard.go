package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/decision"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func guardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guard",
		Short: "Pre-edit guard for quality gate enforcement",
		Long: `Guard commands for managing workstream editing scope.

Prevents editing files outside the active workstream's scope.
This is part of TDD discipline - one workstream at a time.

Examples:
  sdp guard activate 00-001-01
  sdp guard check internal/file.go
  sdp guard status`,
	}

	cmd.AddCommand(guardActivate())
	cmd.AddCommand(guardCheck())
	cmd.AddCommand(guardStatus())
	cmd.AddCommand(guardDeactivate())

	return cmd
}

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
			fmt.Printf("✅ Activated WS: %s\n", activeWS)

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

			return nil
		},
	}
}

func guardCheck() *cobra.Command {
	return &cobra.Command{
		Use:   "check <file>",
		Short: "Check if file edit is allowed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			// Resolve to absolute path
			absPath, err := guard.ResolvePath(filePath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

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

			// Check edit permission
			result, err := skill.CheckEdit(absPath)
			if err != nil {
				return fmt.Errorf("failed to check edit: %w", err)
			}

			// Display result
			if result.Allowed {
				fmt.Printf("✅ ALLOWED: %s\n", result.Reason)
				fmt.Printf("   Active WS: %s\n", result.WSID)
				return nil
			}

			// Not allowed
			fmt.Printf("❌ BLOCKED: %s\n", result.Reason)
			if result.WSID != "" {
				fmt.Printf("   Active WS: %s\n", result.WSID)
			}
			if len(result.ScopeFiles) > 0 {
				fmt.Printf("   Scope files:\n")
				for _, f := range result.ScopeFiles {
					fmt.Printf("     - %s\n", f)
				}
			}
			return fmt.Errorf("file edit not allowed: %s", result.Reason)
		},
	}
}

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

			return nil
		},
	}
}

// warnSimilarFailures loads past decisions and prints warning if similar failed decisions exist (AC3, AC4, AC8).
func warnSimilarFailures(wsID string) {
	root, err := config.FindProjectRoot()
	if err != nil {
		return
	}
	logger, err := decision.NewLogger(root)
	if err != nil {
		return
	}
	decisions, err := logger.LoadAll()
	if err != nil || len(decisions) == 0 {
		return
	}
	matches := evidence.FindSimilarDecisions("", nil, decisions)
	if len(matches) == 0 {
		return
	}
	fmt.Fprintln(os.Stderr, "⚠️  Similar past decision(s) found:")
	for _, m := range matches {
		fmt.Fprintf(os.Stderr, "   Decision: %q\n", m.Question)
		fmt.Fprintf(os.Stderr, "   Outcome: %s\n", m.Outcome)
		fmt.Fprintf(os.Stderr, "   Source: WS %s\n", m.WorkstreamID)
		if len(m.Tags) > 0 {
			fmt.Fprintf(os.Stderr, "   Tags: %v\n", m.Tags)
		}
		fmt.Fprintln(os.Stderr, "   Continue anyway? [y/N]")
	}
}

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

			fmt.Println("✓ Guard deactivated")

			return nil
		},
	}
}
