package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/sdpinit"
	"github.com/fall-out-bug/sdp/internal/telemetry"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func adoptCmd() *cobra.Command {
	var force bool
	var full bool

	cmd := &cobra.Command{
		Use:   "adopt",
		Short: "Adopt current changes into SDP",
		Long: `Convert a successful 'sdp try' session into a full SDP setup:
  - Creates .sdp/ directory structure (equivalent to 'sdp init')
  - Creates .claude/settings.json with SDP skill configuration
  - Commits both .sdp/ and .claude/ to git
  - Preserves all code changes from the trial

Adoption mode:
  By default, 'sdp adopt' enables adoption_mode which disables quality
  gates (file size, coverage, TDD) so legacy code passes @build without
  enforcement. Evidence logging stays enabled (lightweight, non-blocking).

  Use 'sdp adopt --full' when you are ready to enable all quality gates.

This is the next step after accepting a trial with 'sdp try --keep'.`,
		Example: `  # Adopt with quality gates disabled (adoption mode)
  sdp adopt

  # Force adopt even if .sdp exists
  sdp adopt --force

  # Enable all quality gates (graduation from adoption mode)
  sdp adopt --full`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			startTime := time.Now()
			projectPath := "."

			// Convert to absolute path
			absPath, err := filepath.Abs(projectPath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Handle --full: just flip adoption_mode to false and return
			if full {
				return setFullMode(absPath)
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
				fmt.Println("Warning: .sdp directory exists but may not be fully initialized")
				fmt.Println("   Use --force to reinitialize completely")
			}

			// Initialize telemetry collector (after checks, UX metrics now go to user config dir)
			uxMetrics, err := telemetry.NewUXMetricsCollector("")
			if err != nil {
				// Don't fail the command if telemetry fails
				fmt.Fprintf(os.Stderr, "Warning: failed to initialize telemetry: %v\n", err)
			}

			// Create .sdp/ directory structure before sdpinit (which only creates .claude/)
			if err := createSDPDirectory(absPath); err != nil {
				return fmt.Errorf("failed to create .sdp/ directory: %w", err)
			}
			fmt.Println("ok .sdp/ directory created")

			// Run SDP init (creates .claude/ with settings, skills, agents)
			// Note: sdpinit.Run() may overwrite config.yml, so set adoption mode AFTER.
			fmt.Println("Adopting project into SDP...")
			cfg := sdpinit.Config{
				ProjectType: "auto",
				Force:       force,
				Headless:    false,
			}
			if err := sdpinit.Run(cfg); err != nil {
				return fmt.Errorf("failed to initialize SDP: %w", err)
			}

			fmt.Println("ok SDP structure created")

			// Set adoption_mode: true in config (after sdpinit.Run which writes default config)
			if err := config.SetAdoptionMode(absPath, true); err != nil {
				return fmt.Errorf("failed to set adoption mode: %w", err)
			}
			fmt.Println("ok adoption mode enabled (quality gates disabled)")

			// Commit the .sdp/ and .claude/ structure
			fmt.Println("\nCommitting SDP structure...")
			commitSuccess := true
			if err := commitSDPStructure(); err != nil {
				commitSuccess = false
				fmt.Printf("Warning: failed to commit SDP structure: %v\n", err)
				fmt.Println("  Please commit manually: git add .sdp/ .claude/ && git commit -m 'Initialize SDP'")
			} else {
				fmt.Println("ok SDP structure committed")
			}

			fmt.Println("\nNext steps:")
			fmt.Println("  1. Review the .sdp/ structure")
			fmt.Println("  2. Continue with SDP workflow: sdp plan 'your feature'")
			fmt.Println("  3. When ready for full quality gates: sdp adopt --full")

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
	cmd.Flags().BoolVar(&full, "full", false, "Enable all quality gates (disable adoption mode)")

	return cmd
}

// setFullMode disables adoption mode, enabling all quality gates.
func setFullMode(projectRoot string) error {
	sdpDir := filepath.Join(projectRoot, ".sdp")
	if _, err := os.Stat(sdpDir); os.IsNotExist(err) {
		return fmt.Errorf(".sdp/ directory not found. Run 'sdp adopt' first")
	}

	if err := config.SetAdoptionMode(projectRoot, false); err != nil {
		return fmt.Errorf("failed to disable adoption mode: %w", err)
	}

	fmt.Println("ok adoption mode disabled -- all quality gates enabled")
	fmt.Println("  Gates: coverage, complexity, file size, types, TDD")
	return nil
}

// createSDPDirectory creates the .sdp/ directory structure with essential config files.
// This is separate from sdpinit.Run() which only creates .claude/.
func createSDPDirectory(projectPath string) error {
	sdpDir := filepath.Join(projectPath, ".sdp")

	// Create .sdp/ subdirectories
	dirs := []string{
		filepath.Join(sdpDir, "log"),
		filepath.Join(sdpDir, "evidence"),
		filepath.Join(sdpDir, "checkpoints"),
		filepath.Join(sdpDir, "metrics"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}

	// Create .sdp/config.yml if it doesn't exist
	configPath := filepath.Join(sdpDir, "config.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultCfg := config.DefaultConfig()
		defaultCfg.AdoptionMode = true
		out, marshalErr := yaml.Marshal(defaultCfg)
		if marshalErr != nil {
			return fmt.Errorf("marshal default config: %w", marshalErr)
		}
		if err := os.WriteFile(configPath, out, 0644); err != nil {
			return fmt.Errorf("create config.yml: %w", err)
		}
	}

	// Create .sdp/guard-rules.yml if it doesn't exist
	guardPath := filepath.Join(sdpDir, "guard-rules.yml")
	if _, err := os.Stat(guardPath); os.IsNotExist(err) {
		guardContent := `# SDP Guard Rules
# Controls which files can be edited per workstream
version: "1.0.0"
`
		if err := os.WriteFile(guardPath, []byte(guardContent), 0644); err != nil {
			return fmt.Errorf("create guard-rules.yml: %w", err)
		}
	}

	return nil
}

// commitSDPStructure commits the .sdp/ and .claude/ structure to git
func commitSDPStructure() error {
	// Add .sdp/ and .claude/ directories
	addCmd := exec.Command("git", "add", ".sdp/", ".claude/")
	if output, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add to git: %s: %w", string(output), err)
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
