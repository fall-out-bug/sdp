package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/memory"
	"github.com/spf13/cobra"
)

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
