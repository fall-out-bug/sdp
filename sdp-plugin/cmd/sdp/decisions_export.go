package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/decision"
	"github.com/spf13/cobra"
)

// decisionsExportCmd exports decisions to markdown
func decisionsExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export [output]",
		Short: "Export decisions to markdown",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := findProjectRoot()
			if err != nil {
				return err
			}

			logger, err := decision.NewLogger(root)
			if err != nil {
				return err
			}

			decisions, err := logger.LoadAll()
			if err != nil {
				return err
			}

			if len(decisions) == 0 {
				fmt.Println("No decisions to export")
				return nil
			}

			// Determine output path
			outputPath := filepath.Join(root, "docs", "decisions", "DECISIONS.md")
			if len(args) > 0 {
				// Validate path is within project root
				userPath := args[0]
				if filepath.IsAbs(userPath) {
					return fmt.Errorf("absolute paths not allowed: %s", userPath)
				}
				// Clean path to resolve any ".." elements
				cleanPath := filepath.Clean(userPath)
				// Ensure path doesn't escape root
				fullPath := filepath.Join(root, cleanPath)
				if !strings.HasPrefix(fullPath, root) {
					return fmt.Errorf("path escapes project root: %s", userPath)
				}
				outputPath = fullPath
			}

			// Create output directory if needed
			outputDir := filepath.Dir(outputPath)
			if err := os.MkdirAll(outputDir, 0o755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Create markdown
			var md strings.Builder
			md.WriteString("# Architectural Decisions\n\n")
			md.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02")))
			md.WriteString(fmt.Sprintf("**Total:** %d decisions\n\n", len(decisions)))

			for i, d := range decisions {
				md.WriteString(fmt.Sprintf("## %d. %s\n\n", i+1, d.Decision))
				md.WriteString(fmt.Sprintf("**Date:** %s\n", d.Timestamp.Format("2006-01-02 15:04:05")))
				md.WriteString(fmt.Sprintf("**Type:** %s\n", d.Type))
				md.WriteString(fmt.Sprintf("**Maker:** %s\n", d.DecisionMaker))

				if d.FeatureID != "" {
					md.WriteString(fmt.Sprintf("**Feature:** %s\n", d.FeatureID))
				}
				if d.WorkstreamID != "" {
					md.WriteString(fmt.Sprintf("**Workstream:** %s\n", d.WorkstreamID))
				}

				md.WriteString("\n### Question\n\n")
				md.WriteString(d.Question)
				md.WriteString("\n\n")

				md.WriteString("### Decision\n\n")
				md.WriteString(d.Decision)
				md.WriteString("\n\n")

				md.WriteString("### Rationale\n\n")
				md.WriteString(d.Rationale)
				md.WriteString("\n\n")

				if len(d.Alternatives) > 0 {
					md.WriteString("### Alternatives Considered\n\n")
					for _, alt := range d.Alternatives {
						md.WriteString("- ")
						md.WriteString(alt)
						md.WriteString("\n")
					}
					md.WriteString("\n")
				}

				md.WriteString("---\n\n")
			}

			// Write to file
			if err := os.WriteFile(outputPath, []byte(md.String()), 0o644); err != nil {
				return err
			}

			fmt.Printf("Exported %d decisions to %s\n", len(decisions), outputPath)

			return nil
		},
	}
}
