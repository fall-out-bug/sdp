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
	var quick bool
	var deep bool
	var bootstrapSDP bool
	var focus string

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
			modeCount := 0
			if quick {
				modeCount++
			}
			if deep {
				modeCount++
			}
			if bootstrapSDP {
				modeCount++
			}
			if modeCount > 1 {
				return fmt.Errorf("choose only one of --quick, --deep, or --bootstrap-sdp")
			}

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

			opts := realityemitter.Options{
				Mode:  realityemitter.ModeDeep,
				Focus: focus,
			}
			switch {
			case quick:
				opts.Mode = realityemitter.ModeQuick
			case bootstrapSDP:
				opts.Mode = realityemitter.ModeBootstrapSDP
			case deep:
				opts.Mode = realityemitter.ModeDeep
			}

			paths, err := realityemitter.EmitOSSWithOptions(absRoot, opts)
			if err != nil {
				return fmt.Errorf("emit reality OSS artifacts: %w", err)
			}

			ui.SuccessLine("Generated %d reality OSS artifact(s) in %s [%s]", len(paths), absRoot, opts.Mode)
			for _, path := range paths {
				fmt.Printf("  - %s\n", path)
			}

			return nil
		},
	}

	emitOSSCmd.Flags().StringVar(&root, "root", "", "Project root path (default: auto-detect from current directory)")
	emitOSSCmd.Flags().BoolVar(&quick, "quick", false, "Run the fast baseline mode with reduced detail")
	emitOSSCmd.Flags().BoolVar(&deep, "deep", false, "Run the full single-repo baseline mode (default)")
	emitOSSCmd.Flags().BoolVar(&bootstrapSDP, "bootstrap-sdp", false, "Prioritize SDP bootstrap recommendations in the emitted artifacts")
	emitOSSCmd.Flags().StringVar(&focus, "focus", "", "Optional analysis focus: architecture, quality, testing, docs, security")
	cmd.AddCommand(emitOSSCmd)

	return cmd
}
