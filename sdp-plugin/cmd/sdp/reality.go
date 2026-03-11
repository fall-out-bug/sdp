package main

import (
	"fmt"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/config"
	realityemitter "github.com/fall-out-bug/sdp/internal/reality"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

func realityCmd() *cobra.Command {
	var root string

	cmd := &cobra.Command{
		Use:   "reality",
		Short: "Generate reality/reality-pro baseline artifacts",
		Long: `Generate reality/reality-pro baseline artifacts.

OSS command set:
  emit-oss    Emit the required OSS artifact set into docs/reality and .sdp/reality`,
	}

	emitOSSCmd := &cobra.Command{
		Use:   "emit-oss",
		Short: "Emit the OSS reality artifact set",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectRoot := root
			if projectRoot == "" {
				detectedRoot, err := config.FindProjectRoot()
				if err != nil {
					return fmt.Errorf("find project root: %w", err)
				}
				projectRoot = detectedRoot
			}

			absRoot, err := filepath.Abs(projectRoot)
			if err != nil {
				return fmt.Errorf("resolve project root: %w", err)
			}

			paths, err := realityemitter.EmitOSS(absRoot)
			if err != nil {
				return fmt.Errorf("emit reality OSS artifacts: %w", err)
			}

			ui.SuccessLine("Generated %d reality OSS artifact(s) in %s", len(paths), absRoot)
			for _, path := range paths {
				fmt.Printf("  - %s\n", path)
			}

			return nil
		},
	}

	emitOSSCmd.Flags().StringVar(&root, "root", "", "Project root path (default: auto-detect from current directory)")
	cmd.AddCommand(emitOSSCmd)

	return cmd
}
