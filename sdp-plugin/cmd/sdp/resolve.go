package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/resolver"
	"github.com/spf13/cobra"
)

// resolveCmd represents the resolve command
func resolveCmd() *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "resolve <id>",
		Short: "Resolve any task identifier to its file path",
		Long: `Resolve workstream IDs, beads IDs, or issue IDs to their file paths.

Supports:
  - Workstream IDs: 00-064-01, 99-F064-0001
  - Beads IDs: sdp-ushh, abc-123
  - Issue IDs: ISSUE-0001

The resolver auto-detects the ID type from its pattern and returns
the file path and metadata.

Examples:
  sdp resolve 00-064-01
  sdp resolve sdp-ushh
  sdp resolve ISSUE-0001 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			r := resolver.NewResolver(
				resolver.WithWorkstreamDir("docs/workstreams/backlog"),
				resolver.WithIssuesDir("docs/issues"),
				resolver.WithIndexFile(".sdp/issues-index.jsonl"),
			)

			result, err := r.Resolve(id)
			if err != nil {
				return fmt.Errorf("failed to resolve %s: %w", id, err)
			}

			if outputJSON {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				return encoder.Encode(result)
			}

			// Human-readable output
			fmt.Printf("Type: %s\n", result.Type)
			fmt.Printf("ID: %s\n", result.ID)
			if result.WSID != "" {
				fmt.Printf("Workstream: %s\n", result.WSID)
			}
			fmt.Printf("Path: %s\n", result.Path)
			if result.Title != "" {
				fmt.Printf("Title: %s\n", result.Title)
			}
			if result.Status != "" {
				fmt.Printf("Status: %s\n", result.Status)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	return cmd
}
