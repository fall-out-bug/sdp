package main

import (
	"fmt"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/parser"
	"github.com/spf13/cobra"
)

func parseCmd() *cobra.Command {
	var validateFlag bool

	cmd := &cobra.Command{
		Use:   "parse <ws-id>",
		Short: "Parse and display workstream information",
		Long: `Parse a workstream markdown file and display its contents.

Args:
  ws-id    Workstream ID (e.g., 00-050-01)

Examples:
  sdp parse 00-050-01
  sdp parse --validate docs/workstreams/backlog/00-050-01.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if validateFlag {
				return validateRun(cmd, args)
			}

			return parseRun(cmd, args)
		},
	}

	cmd.Flags().BoolVar(&validateFlag, "validate", false, "Validate workstream file")

	return cmd
}

func parseRun(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("workstream ID required")
	}

	wsID := args[0]

	wsID = filepath.Clean(wsID)
	if wsID != args[0] {
		return fmt.Errorf("invalid workstream ID: path traversal detected")
	}

	wsPath, err := findWorkstreamFile(wsID)
	if err != nil {
		return err
	}

	ws, err := parser.ParseWorkstream(wsPath)
	if err != nil {
		return fmt.Errorf("failed to parse workstream: %w", err)
	}

	displayWorkstream(ws)

	if evidence.Enabled() {
		scopeFiles := append([]string{}, ws.Scope.Implementation...)
		scopeFiles = append(scopeFiles, ws.Scope.Tests...)
		if err := evidence.EmitSync(evidence.PlanEventWithFeature(ws.ID, ws.Feature, scopeFiles)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: evidence emit: %v\n", err)
		}
	}

	return nil
}
