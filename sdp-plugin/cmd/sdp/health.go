package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/spf13/cobra"
)

// healthCmd returns the health command
func healthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check SDP system health",
		Long: `Check SDP system health status.

Verifies:
  - .sdp directory structure
  - Memory database accessibility
  - Git configuration
  - Essential directories

Exit codes:
  0 - All checks passed
  1 - Some checks failed`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("failed to find project root: %w", err)
			}

			checks := []struct {
				name string
				fn   func(string) (bool, string)
			}{
				{"SDP directory", checkSDPDir},
				{"Memory database", checkMemoryDB},
				{"Git repository", checkGitRepo},
				{"Workstreams directory", checkWorkstreamsDir},
				{"Config file", checkConfigFile},
			}

			allPassed := true
			fmt.Println("SDP Health Check")
			fmt.Println("================")

			for _, check := range checks {
				passed, message := check.fn(root)
				icon := "✓"
				if !passed {
					icon = "✗"
					allPassed = false
				}
				fmt.Printf("%s %s: %s\n", icon, check.name, message)
			}

			fmt.Println()
			if allPassed {
				fmt.Println("All health checks passed!")
				return nil
			}
			return fmt.Errorf("some health checks failed")
		},
	}
}

func checkSDPDir(root string) (bool, string) {
	sdpDir := filepath.Join(root, ".sdp")
	info, err := os.Stat(sdpDir)
	if err != nil {
		return false, "not found"
	}
	if !info.IsDir() {
		return false, "not a directory"
	}
	return true, "exists"
}

func checkMemoryDB(root string) (bool, string) {
	dbPath := filepath.Join(root, ".sdp", "memory.db")
	info, err := os.Stat(dbPath)
	if err != nil {
		return false, "not found (run 'sdp memory index' to create)"
	}
	return true, fmt.Sprintf("exists (%d bytes)", info.Size())
}

func checkGitRepo(root string) (bool, string) {
	gitDir := filepath.Join(root, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false, "not a git repository"
	}
	if info.IsDir() {
		return true, "valid git repository"
	}
	// Could be a git worktree file
	return true, "git worktree"
}

func checkWorkstreamsDir(root string) (bool, string) {
	wsDir := filepath.Join(root, "docs", "workstreams", "backlog")
	info, err := os.Stat(wsDir)
	if err != nil {
		return false, "not found"
	}
	if !info.IsDir() {
		return false, "not a directory"
	}

	// Count workstream files
	files, err := filepath.Glob(filepath.Join(wsDir, "*.md"))
	if err != nil {
		return true, "exists"
	}
	return true, fmt.Sprintf("exists (%d workstreams)", len(files))
}

func checkConfigFile(root string) (bool, string) {
	configPath := filepath.Join(root, ".sdp", "config.yml")
	_, err := os.Stat(configPath)
	if err != nil {
		return false, "not found (optional)"
	}
	return true, "exists"
}
