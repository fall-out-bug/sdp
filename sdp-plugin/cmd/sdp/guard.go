package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/collision"
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

Prevents editing files outside of active workstream's scope.
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

func guardCheck() *cobra.Command {
	var staged, jsonOutput bool

	cmd := &cobra.Command{
		Use:   "check [file]",
		Short: "Check if file edit is allowed or check staged files",
		Long: `Check file edit permissions or staged files for policy compliance.

Single file mode (legacy):
  sdp guard check <file>

Staged mode (new):
  sdp guard check --staged [--json]

Staged mode checks only staged files using git diff --cached.
ERROR findings block the commit (exit code 1).
WARNING findings are displayed but don't block (hybrid mode).

Uses environment variables for CI diff-range:
  CI_BASE_SHA: Base commit SHA
  CI_HEAD_SHA: Head commit SHA`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get config directory
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

			// Staged mode: check staged files (AC1, AC4, AC5)
			if staged {
				// Parse options (including CI env vars)
				opts := guard.ParseCheckOptions()
				opts.Staged = true
				opts.JSON = jsonOutput

				// Run staged check
				result, err := skill.StagedCheck(opts)
				if err != nil {
					// Runtime error (AC3: exit code 2)
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(guard.ExitCodeRuntimeError)
					return nil
				}

				// Output based on format (AC4: human-readable, AC5: JSON)
				if opts.JSON {
					data, err := json.MarshalIndent(result, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal JSON: %w", err)
					}
					fmt.Println(string(data))
				} else {
					printHumanReadableResult(result)
				}

				// Exit with appropriate code (AC3)
				if !result.Success {
					os.Exit(result.ExitCode)
				}

				return nil
			}

			// Legacy single file mode
			if len(args) == 0 {
				return fmt.Errorf("requires a <file> argument or --staged flag")
			}

			filePath := args[0]

			// Resolve to absolute path
			absPath, err := guard.ResolvePath(filePath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Check edit permission
			result, err := skill.CheckEdit(absPath)
			if err != nil {
				return fmt.Errorf("failed to check edit: %w", err)
			}

			// Display result
			if result.Allowed {
				fmt.Printf("ALLOWED: %s\n", result.Reason)
				fmt.Printf("   Active WS: %s\n", result.WSID)
				return nil
			}

			// Not allowed
			fmt.Printf("BLOCKED: %s\n", result.Reason)
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

	// Add flags for staged mode
	cmd.Flags().BoolVar(&staged, "staged", false, "Check staged files")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output JSON format for CI")

	return cmd
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
	fmt.Fprintln(os.Stderr, "Similar past decision(s) found:")
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

// warnContractViolations checks for contract violations and prints warnings.
// AC2: @build runs contract validation after code generation.
// AC3: Contract violation = build warning (not error in P1; enforcement in P2).
func warnContractViolations() {
	// Find project root
	root, err := config.FindProjectRoot()
	if err != nil {
		return // Not in project, skip validation
	}

	// Check for .contracts directory
	contractsDir := filepath.Join(root, ".contracts")
	if _, err := os.Stat(contractsDir); os.IsNotExist(err) {
		return // No contracts, skip validation
	}

	// Find implementation directory (default: internal/)
	implDir := filepath.Join(root, "internal")
	if _, err := os.Stat(implDir); os.IsNotExist(err) {
		return // No impl dir, skip validation
	}

	// Run validation
	violations, err := collision.ValidateContractsInDir(contractsDir, implDir)
	if err != nil {
		// Validation failed (not violations list), log warning
		fmt.Fprintf(os.Stderr, "Contract validation error: %v\n", err)
		return
	}

	// Print violations as warnings
	if len(violations) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "Contract violations detected (%d):\n", len(violations))
	for _, v := range violations {
		severity := "WARNING"
		if v.Severity == "error" {
			severity = "ERROR"
		}
		fmt.Fprintf(os.Stderr, "   [%s] %s: %s\n", severity, v.Type, v.Message)
	}

	// In P1, violations are warnings only (not blocking)
	fmt.Fprintf(os.Stderr, "   Review violations before proceeding\n")
}

// printHumanReadableResult prints check results in human-readable format (AC4)
func printHumanReadableResult(result *guard.CheckResult) {
	if result.Success && result.Summary.Total == 0 {
		fmt.Println("No issues found")
		return
	}

	// Print findings
	for _, f := range result.Findings {
		severity := string(f.Severity)
		icon := "WARNING"
		if f.Severity == guard.SeverityError {
			icon = "ERROR"
		}
		fmt.Printf("[%s] %s: %s\n", icon, severity, f.Message)
		if f.File != "" {
			fmt.Printf("   File: %s\n", f.File)
			if f.Line > 0 {
				fmt.Printf("   Line: %d\n", f.Line)
			}
		}
	}

	// Print summary
	fmt.Printf("\nSummary: %d total (%d errors, %d warnings)\n",
		result.Summary.Total, result.Summary.Errors, result.Summary.Warnings)

	// AC8: Hybrid mode message
	if !result.Success {
		fmt.Println("\nCheck failed: Errors must be fixed before committing")
	} else if result.Summary.Warnings > 0 {
		fmt.Println("\nCheck passed with warnings")
	}
}
