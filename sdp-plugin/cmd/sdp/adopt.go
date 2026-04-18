package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/sdpinit"
	"github.com/spf13/cobra"
)

func adoptCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "adopt",
		Short: "Adopt current changes into SDP",
		Long: `Convert a successful 'sdp try' session into a full SDP setup:
  - Creates .sdp/ structure equivalent to 'sdp init'
  - Commits the .sdp/ structure
  - Preserves all changes from the trial

This is the next step after accepting a trial with 'sdp try --keep'.`,
		Example: `  # Adopt current changes
  sdp adopt

  # Force adopt even if .sdp exists
  sdp adopt --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath := "."

			// Convert to absolute path
			absPath, err := filepath.Abs(projectPath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}

			// Check if .sdp already exists
			sdpPath := filepath.Join(absPath, ".sdp")
			if _, err := os.Stat(sdpPath); err == nil {
				if !force {
					return fmt.Errorf(".sdp directory already exists. Use --force to reinitialize")
				}
				fmt.Println("⚠ Reinitializing existing .sdp directory")
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

			// TODO: Commit the .sdp/ structure
			// This would involve: git add .sdp/ && git commit -m "Initialize SDP"

			fmt.Println("\nNext steps:")
			fmt.Println("  1. Review the .sdp/ structure")
			fmt.Println("  2. Commit the changes: git add .sdp/ && git commit -m 'Initialize SDP'")
			fmt.Println("  3. Continue with SDP workflow: sdp plan 'your feature'")

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Reinitialize even if .sdp exists")

	return cmd
}
