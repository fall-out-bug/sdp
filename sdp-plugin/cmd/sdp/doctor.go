package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/doctor"
	"github.com/spf13/cobra"
)

func doctorCmd() *cobra.Command {
	var driftCheck bool
	var repair bool
	var deep bool
	var migrate bool
	var dryRun bool
	var rollback string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check SDP environment",
		Long: `Check that your environment is properly configured for SDP.

Verifies:
  - Git is installed
  - Claude Code CLI is available (optional)
  - Go compiler is available (for building binary)
  - .claude/ directory exists and is properly structured
  - Documentation-code drift (with --drift flag)

Modes:
  --repair    Automatically fix detected issues
  --deep      Comprehensive environment analysis
  --migrate   Migrate config to latest version
  --rollback  Restore config from backup`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle rollback first
			if rollback != "" {
				fmt.Println("Rolling back config...")
				if err := doctor.RollbackMigration(rollback); err != nil {
					return fmt.Errorf("rollback failed: %w", err)
				}
				fmt.Printf("✓ Config restored from %s\n", rollback)
				return nil
			}

			// Handle migration
			if migrate {
				fmt.Println("Migrating config...")
				m, err := doctor.MigrateConfig(dryRun)
				if err != nil {
					return fmt.Errorf("migration failed: %w", err)
				}

				if dryRun {
					fmt.Printf("✓ %s (dry run)\n", m.Message)
				} else {
					fmt.Printf("✓ %s\n", m.Message)
					if m.BackupPath != "" {
						fmt.Printf("  Backup: %s\n", m.BackupPath)
					}
				}
				return nil
			}

			// Run standard checks
			opts := doctor.RunOptions{
				DriftCheck: driftCheck,
			}
			results := doctor.RunWithOptions(opts)

			// Print results
			fmt.Println("SDP Environment Check")
			fmt.Println("=====================")

			for _, r := range results {
				icon := "✓"
				color := ""
				if r.Status == "warning" {
					icon = "⚠"
					color = " (optional)"
				} else if r.Status == "error" {
					icon = "✗"
				}

				fmt.Printf("%s %s%s\n", icon, r.Name, color)
				fmt.Printf("    %s\n\n", r.Message)
			}

			// Run repair if requested
			if repair {
				fmt.Println("\nRepair Mode")
				fmt.Println("===========")
				actions := doctor.RunWithRepair()
				for _, a := range actions {
					icon := "✓"
					if a.Status == "failed" || a.Status == "manual" {
						icon = "✗"
					} else if a.Status == "skipped" {
						icon = "→"
					}
					fmt.Printf("%s %s [%s]\n", icon, a.Check, a.Status)
					fmt.Printf("    %s\n\n", a.Message)
				}

				if doctor.HasUnfixableErrors(actions) {
					return fmt.Errorf("some issues require manual intervention")
				}
				fmt.Println("All repairable issues fixed!")
			}

			// Run deep checks if requested
			if deep {
				fmt.Println("\nDeep Diagnostics")
				fmt.Println("================")
				deepResults := doctor.RunDeepChecks()
				for _, r := range deepResults {
					icon := "✓"
					if r.Status == "warning" {
						icon = "⚠"
					} else if r.Status == "error" {
						icon = "✗"
					}
					fmt.Printf("%s %s [%v]\n", icon, r.Check, r.Duration.Round(0))
					fmt.Printf("    %s\n\n", r.Message)
				}
			}

			// Exit code based on results
			hasErrors := false
			for _, r := range results {
				if r.Status == "error" {
					hasErrors = true
				}
			}

			if hasErrors {
				return fmt.Errorf("some required checks failed")
			}

			if !repair && !deep {
				fmt.Println("All required checks passed!")
			}
			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&driftCheck, "drift", false, "Check for documentation-code drift in recent workstreams")
	cmd.Flags().BoolVar(&repair, "repair", false, "Automatically fix detected issues where possible")
	cmd.Flags().BoolVar(&deep, "deep", false, "Run comprehensive environment analysis")
	cmd.Flags().BoolVar(&migrate, "migrate", false, "Migrate config to latest version")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview migration without making changes")
	cmd.Flags().StringVar(&rollback, "rollback", "", "Restore config from backup file")

	return cmd
}
