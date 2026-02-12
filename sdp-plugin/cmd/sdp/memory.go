package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/memory"
	"github.com/spf13/cobra"
)

func memoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Long-term memory management",
		Long: `Manage the SDP long-term memory system.

The memory system indexes project artifacts for fast search and
provides hybrid search capabilities (full-text + semantic).

Examples:
  sdp memory index              # Index all docs/ artifacts
  sdp memory search "API"       # Search for "API" in artifacts
  sdp memory stats              # Show memory statistics`,
	}

	cmd.AddCommand(memoryIndexCmd())
	cmd.AddCommand(memorySearchCmd())
	cmd.AddCommand(memoryStatsCmd())

	return cmd
}

func memoryIndexCmd() *cobra.Command {
	var docsDir string
	var dbPath string

	cmd := &cobra.Command{
		Use:   "index",
		Short: "Index project artifacts",
		Long: `Index all markdown files in the docs/ directory.

Stores indexed artifacts in SQLite database (.sdp/memory.db) with
full-text search capabilities.

Supports incremental updates - only re-indexes changed files
based on SHA256 hash comparison.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Default paths
			if docsDir == "" {
				docsDir = "docs"
			}
			if dbPath == "" {
				dbPath = ".sdp/memory.db"
			}

			// Ensure paths are absolute
			if !filepath.IsAbs(docsDir) {
				cwd, _ := os.Getwd()
				docsDir = filepath.Join(cwd, docsDir)
			}
			if !filepath.IsAbs(dbPath) {
				cwd, _ := os.Getwd()
				dbPath = filepath.Join(cwd, dbPath)
			}

			// Check if docs directory exists
			if _, err := os.Stat(docsDir); os.IsNotExist(err) {
				return fmt.Errorf("docs directory not found: %s", docsDir)
			}

			// Create indexer
			indexer, err := memory.NewIndexer(docsDir, dbPath)
			if err != nil {
				return fmt.Errorf("failed to create indexer: %w", err)
			}
			defer indexer.Close()

			// Index directory
			stats, err := indexer.IndexDirectory()
			if err != nil {
				return fmt.Errorf("indexing failed: %w", err)
			}

			// Print results
			fmt.Println("Indexing complete:")
			fmt.Printf("  Total files: %d\n", stats.TotalFiles)
			fmt.Printf("  New: %d\n", stats.Indexed)
			fmt.Printf("  Updated: %d\n", stats.Updated)
			fmt.Printf("  Skipped (unchanged): %d\n", stats.Skipped)
			if stats.Errors > 0 {
				fmt.Printf("  Errors: %d\n", stats.Errors)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&docsDir, "docs", "", "Docs directory (default: docs)")
	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: .sdp/memory.db)")

	return cmd
}

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

func memoryStatsCmd() *cobra.Command {
	var dbPath string

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show memory statistics",
		Long:  `Display statistics about the indexed artifacts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// Get all artifacts
			artifacts, err := store.ListAll()
			if err != nil {
				return fmt.Errorf("failed to list artifacts: %w", err)
			}

			// Calculate stats
			byType := make(map[string]int)
			byFeature := make(map[string]int)
			for _, a := range artifacts {
				byType[a.Type]++
				if a.FeatureID != "" {
					byFeature[a.FeatureID]++
				}
			}

			fmt.Println("Memory Statistics")
			fmt.Println("================")
			fmt.Printf("Total artifacts: %d\n\n", len(artifacts))

			fmt.Println("By Type:")
			for t, c := range byType {
				fmt.Printf("  %s: %d\n", t, c)
			}

			if len(byFeature) > 0 {
				fmt.Println("\nBy Feature:")
				for f, c := range byFeature {
					fmt.Printf("  %s: %d\n", f, c)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbPath, "db", "", "Database path (default: .sdp/memory.db)")

	return cmd
}
