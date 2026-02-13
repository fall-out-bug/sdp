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
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Create markdown
			var md string
			md += "# Architectural Decisions\n\n"
			md += fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02"))
			md += fmt.Sprintf("**Total:** %d decisions\n\n", len(decisions))

			for i, d := range decisions {
				md += fmt.Sprintf("## %d. %s\n\n", i+1, d.Decision)
				md += fmt.Sprintf("**Date:** %s\n", d.Timestamp.Format("2006-01-02 15:04:05"))
				md += fmt.Sprintf("**Type:** %s\n", d.Type)
				md += fmt.Sprintf("**Maker:** %s\n", d.DecisionMaker)

				if d.FeatureID != "" {
					md += fmt.Sprintf("**Feature:** %s\n", d.FeatureID)
				}
				if d.WorkstreamID != "" {
					md += fmt.Sprintf("**Workstream:** %s\n", d.WorkstreamID)
				}

				md += "\n### Question\n\n"
				md += d.Question + "\n\n"

				md += "### Decision\n\n"
				md += d.Decision + "\n\n"

				md += "### Rationale\n\n"
				md += d.Rationale + "\n\n"

				if len(d.Alternatives) > 0 {
					md += "### Alternatives Considered\n\n"
					for _, alt := range d.Alternatives {
						md += "- " + alt + "\n"
					}
					md += "\n"
				}

				md += "---\n\n"
			}

			// Write to file
			if err := os.WriteFile(outputPath, []byte(md), 0644); err != nil {
				return err
			}

			fmt.Printf("Exported %d decisions to %s\n", len(decisions), outputPath)

			return nil
		},
	}
}
