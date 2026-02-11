package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/fall-out-bug/sdp/internal/decision"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

const (
	maxFieldLength = 10 * 1024 // 10KB max per field
)

// decisionsCmd returns the decisions command
func decisionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decisions",
		Short: "Manage decision audit trail",
		Long: `View and search decision log for architectural and product decisions.

Decisions are automatically logged during @feature skill interviews
and can be queried using this command.`,
	}

	cmd.AddCommand(decisionsListCmd())
	cmd.AddCommand(decisionsSearchCmd())
	cmd.AddCommand(decisionsExportCmd())
	cmd.AddCommand(decisionsLogCmd())

	return cmd
}

// decisionsListCmd lists all decisions
func decisionsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all decisions",
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
				fmt.Println("No decisions found yet")
				fmt.Println("\nDecisions will be automatically logged during @feature skill interviews.")
				return nil
			}

			fmt.Printf("Found %d decision(s):\n\n", len(decisions))

			for i, d := range decisions {
				fmt.Printf("%d. [%s] %s\n", i+1, d.Timestamp.Format("2006-01-02"), d.Decision)
				fmt.Printf("   Type: %s\n", d.Type)
				fmt.Printf("   Question: %s\n", d.Question)
				fmt.Printf("   Decision: %s\n", d.Decision)
				fmt.Printf("   Rationale: %s\n", d.Rationale)
				if d.FeatureID != "" {
					fmt.Printf("   Feature: %s\n", d.FeatureID)
				}
				if d.WorkstreamID != "" {
					fmt.Printf("   Workstream: %s\n", d.WorkstreamID)
				}
				fmt.Println()
			}

			return nil
		},
	}
}

// decisionsSearchCmd searches decisions
func decisionsSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search decisions by keyword",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

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

			var found []decision.Decision

			for _, d := range decisions {
				if strings.Contains(d.Question, query) ||
					strings.Contains(d.Decision, query) ||
					strings.Contains(d.Rationale, query) {
					found = append(found, d)
				}
			}

			if len(found) == 0 {
				fmt.Printf("No decisions found matching '%s'\n", query)
				return nil
			}

			fmt.Printf("Found %d decision(s) matching '%s':\n\n", len(found), query)

			for i, d := range found {
				fmt.Printf("%d. [%s] %s\n", i+1, d.Timestamp.Format("2006-01-02"), d.Decision)
				fmt.Printf("   %s\n", d.Question)
				fmt.Println()
			}

			return nil
		},
	}
}

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

// decisionsLogCmd logs a new decision
func decisionsLogCmd() *cobra.Command {
	var decisionType, featureID, workstreamID, question, decisionStr, rationale, alternatives, outcome, maker, reverses string
	var tags []string

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Log a new decision",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required fields
			if question == "" || decisionStr == "" {
				return fmt.Errorf("required flags: --question, --decision")
			}

			// Validate field lengths
			if err := validateFieldLength("question", question, maxFieldLength); err != nil {
				return err
			}
			if err := validateFieldLength("decision", decisionStr, maxFieldLength); err != nil {
				return err
			}
			if rationale != "" {
				if err := validateFieldLength("rationale", rationale, maxFieldLength); err != nil {
					return err
				}
			}
			if alternatives != "" {
				if err := validateFieldLength("alternatives", alternatives, maxFieldLength); err != nil {
					return err
				}
			}

			// Strip control characters
			question = stripControlChars(question)
			decisionStr = stripControlChars(decisionStr)
			rationale = stripControlChars(rationale)
			alternatives = stripControlChars(alternatives)

			// Validate decision type
			if decisionType != "" {
				validTypes := []string{
					decision.DecisionTypeVision,
					decision.DecisionTypeTechnical,
					decision.DecisionTypeTradeoff,
					decision.DecisionTypeExplicit,
				}
				valid := false
				for _, t := range validTypes {
					if decisionType == t {
						valid = true
						break
					}
				}
				if !valid {
					return fmt.Errorf("invalid decision type %q, must be one of: vision, technical, tradeoff, explicit", decisionType)
				}
			}

			root, err := findProjectRoot()
			if err != nil {
				return err
			}

			logger, err := decision.NewLogger(root)
			if err != nil {
				return err
			}

			// Parse alternatives (comma-separated)
			var altList []string
			if alternatives != "" {
				altList = strings.Split(alternatives, ",")
				for i := range altList {
					altList[i] = strings.TrimSpace(altList[i])
				}
			}

			// Parse tags (comma-separated)
			var tagList []string
			if len(tags) > 0 {
				tagList = tags
			}

			// Default maker to "user" if not specified
			if maker == "" {
				maker = "user"
			}

			// Create decision (bridge in logger.Log emits to evidence log)
			d := decision.Decision{
				Question:      question,
				Decision:      decisionStr,
				Rationale:     rationale,
				Type:          decisionType,
				FeatureID:     featureID,
				WorkstreamID:  workstreamID,
				Alternatives:  altList,
				Outcome:       outcome,
				DecisionMaker: maker,
				Tags:          tagList,
				Reverses:      reverses,
			}

			if err := logger.Log(d); err != nil {
				return err
			}
			// AC1: Emit decision event to evidence log (bridge in CLI to avoid decision->evidence cycle)
			var rev *string
			if reverses != "" {
				rev = &reverses
			}
			if err := evidence.EmitSync(evidence.DecisionEvent(workstreamID, question, decisionStr, rationale, altList, 0, tagList, rev)); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "warning: evidence emit: %v\n", err)
			}

			fmt.Printf("✓ Logged decision: %s\n", d.Decision)
			return nil
		},
	}

	cmd.Flags().StringVar(&decisionType, "type", decision.DecisionTypeExplicit, "Decision type")
	cmd.Flags().StringVar(&featureID, "feature-id", "", "Feature ID")
	cmd.Flags().StringVar(&workstreamID, "workstream-id", "", "Workstream ID")
	cmd.Flags().StringVar(&question, "question", "", "Question or problem")
	cmd.Flags().StringVar(&decisionStr, "decision", "", "Decision made")
	cmd.Flags().StringVar(&rationale, "rationale", "", "Rationale for decision")
	cmd.Flags().StringVar(&alternatives, "alternatives", "", "Alternatives considered (comma-separated)")
	cmd.Flags().StringVar(&outcome, "outcome", "", "Expected outcome")
	cmd.Flags().StringVar(&maker, "maker", "", "Decision maker (user/claude/system)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags for categorization")
	cmd.Flags().StringVar(&reverses, "reverses", "", "ID of previous decision being overturned (AC7)")

	return cmd
}

// findProjectRoot finds the project root by looking for .git directory
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root
		}
		dir = parent
	}

	return "", fmt.Errorf("not in a git repository")
}

// validateFieldLength validates field length
func validateFieldLength(fieldName, value string, maxLen int) error {
	if len(value) > maxLen {
		return fmt.Errorf("%s exceeds maximum length of %d bytes (got %d)", fieldName, maxLen, len(value))
	}
	return nil
}

// stripControlChars removes control characters except newline/tab
func stripControlChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			return -1 // Remove control char
		}
		return r
	}, s)
}
