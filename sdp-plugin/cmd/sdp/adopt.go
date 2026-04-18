package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/sdpinit"
	"github.com/fall-out-bug/sdp/internal/telemetry"
	"github.com/spf13/cobra"
)

func adoptCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "adopt",
		Short: "Adopt current changes into SDP",
		Long: `Convert a successful 'sdp try' session into a full SDP setup:
	  - Creates .sdp/ directory structure (equivalent to 'sdp init')
	  - Creates .claude/settings.json with SDP skill configuration
	  - Commits both .sdp/ and .claude/ to git
	  - Preserves all code changes from the trial

	This is the next step after accepting a trial with 'sdp try --keep'.`,
		Example: `  # Adopt current changes
  sdp adopt

  # Force adopt even if .sdp exists
  sdp adopt --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()
			projectPath := "."

			// Convert to absolute path
			absPath, err := filepath.Abs(projectPath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Check if .sdp already exists and is initialized
			sdpPath := filepath.Join(absPath, ".sdp")
			claudeSettingsPath := filepath.Join(absPath, ".claude", "settings.json")

			sdpExists := false
			alreadyInitialized := false

			if _, err := os.Stat(sdpPath); err == nil {
				sdpExists = true
			}

			if _, err := os.Stat(claudeSettingsPath); err == nil {
				alreadyInitialized = true
			}

			if alreadyInitialized && !force {
				return fmt.Errorf("SDP already initialized. Use --force to reinitialize")
			}

			if sdpExists && !force {
				fmt.Println("⚠ .sdp directory exists but may not be fully initialized")
				fmt.Println("   Use --force to reinitialize completely")
			}

			// Initialize telemetry collector (after checks, UX metrics now go to user config dir)
			uxMetrics, err := telemetry.NewUXMetricsCollector("")
			if err != nil {
				// Don't fail the command if telemetry fails
				fmt.Fprintf(os.Stderr, "Warning: failed to initialize telemetry: %v\n", err)
			}

			// Run SDP init
			fmt.Println("Adopting project into SDP...")
			cfg := sdpinit.Config{
				ProjectType: "auto",
				Force:       force,
				Headless:    false,
			}
			if err := sdpinit.Run(cfg); err != nil {
				return fmt.Errorf("failed to initialize SDP: %w", err)
			}

			fmt.Println("✓ SDP structure created")

			// Commit the .sdp/ structure
			fmt.Println("\nCommitting .sdp/ structure...")
			commitSuccess := true
			if err := commitSDPStructure(); err != nil {
				commitSuccess = false
				fmt.Printf("⚠ Warning: failed to commit .sdp/ structure: %v\n", err)
				fmt.Println("  Please commit manually: git add .sdp/ && git commit -m 'Initialize SDP'")
			} else {
				fmt.Println("✓ SDP structure committed")
			}

			fmt.Println("\nNext steps:")
			fmt.Println("  1. Review the .sdp/ structure")
			fmt.Println("  2. Continue with SDP workflow: sdp plan 'your feature'")

			// Record telemetry
			if uxMetrics != nil && commitSuccess {
				duration := time.Since(startTime)
				if err := uxMetrics.RecordAdoptComplete("unknown", duration); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to record telemetry: %v\n", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Reinitialize even if .sdp exists")

	return cmd
}

// commitSDPStructure commits the .sdp/ structure to git
func commitSDPStructure() error {
	// Add .sdp/ directory
	addCmd := exec.Command("git", "add", ".sdp/", ".claude/")
	if output, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add .sdp/ to git: %s: %w", string(output), err)
	}

	// Check if there's anything to commit
	statusCmd := exec.Command("git", "diff", "--cached", "--quiet")
	if err := statusCmd.Run(); err == nil {
		// No changes to commit (exit status 0 means no differences)
		return fmt.Errorf("no changes to commit")
	}

	// Commit the changes
	commitCmd := exec.Command("git", "commit", "-m", "Initialize SDP structure")
	if output, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to commit: %s: %w", string(output), err)
	}

	return nil
}
