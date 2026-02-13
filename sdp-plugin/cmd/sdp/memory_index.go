package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/memory"
	"github.com/spf13/cobra"
)

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
