package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/memory"
	"github.com/spf13/cobra"
)

func memorySearchCmd() *cobra.Command {
	var dbPath string
	var limit int

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search indexed artifacts",
		Long: `Search indexed artifacts using full-text search.

Returns matching artifacts ranked by relevance.
Supports FTS5 query syntax for advanced searches.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			// Default paths
			if dbPath == "" {
				dbPath = ".sdp/memory.db"
			}

			// Ensure path is absolute
			if !filepath.IsAbs(dbPath) {
				cwd, _ := os.Getwd()
				dbPath = filepath.Join(cwd, dbPath)
			}

			// Check if database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				return fmt.Errorf("memory database not found. Run 'sdp memory index' first")
			}

			// Open store
			store, err := memory.NewStore(dbPath)
			if err != nil {
				return fmt.Errorf("failed to open store: %w", err)
			}
			defer store.Close()

			// Search
			results, err := store.Search(query)
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}

			// Print results
			if len(results) == 0 {
				fmt.Println("No results found")
				return nil
			}

			count := len(results)
			if count > limit {
				results = results[:limit]
			}

			fmt.Printf("Found %d results (showing %d):\n\n", count, len(results))
			for i, r := range results {
				fmt.Printf("%d. %s\n", i+1, r.Title)
				fmt.Printf("   Path: %s\n", r.Path)
				if r.FeatureID != "" {
					fmt.Printf("   Feature: %s\n", r.FeatureID)
				}
				if len(r.Tags) > 0 {
					fmt.Printf("   Tags: %v\n", r.Tags)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: .sdp/memory.db)")
	cmd.Flags().IntVarP(&limit, "limit", "n", 10, "Maximum results to show")

	return cmd
}
