package main

import (
	"fmt"
	"strings"

	"github.com/fall-out-bug/sdp/internal/decision"
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
