package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

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
ERROR findings block commit (exit code 1).
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
