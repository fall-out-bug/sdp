package hooks

import (
	"fmt"
	"os"
	"path/filepath"
)

const hookTemplate = `#!/bin/sh
# SDP Git Hook
# Run: sdp validate pre-commit
# Or customize with your own checks

# Check if .claude/ exists
if [ -d ".claude" ]; then
    echo "SDP: Running quality checks..."
    # Add your validation commands here
    # Example: claude "@review" or run tests
fi
`

func Install() error {
	gitDir := ".git/hooks"

	// Check if .git exists
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf(".git directory not found. Run 'git init' first")
	}

	// Create hooks directory if missing
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return fmt.Errorf("create hooks dir: %w", err)
	}

	hooks := map[string]string{
		"pre-commit": hookTemplate,
		"pre-push":   hookTemplate,
	}

	for name, content := range hooks {
		path := filepath.Join(gitDir, name)

		// Check if hook already exists
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("⚠ %s already exists, skipping\n", name)
			continue
		}

		if err := os.WriteFile(path, []byte(content), 0755); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
		fmt.Printf("✓ Installed %s\n", name)
	}

	fmt.Println("\nGit hooks installed!")
	fmt.Println("Customize hooks in .git/hooks/ if needed")

	return nil
}

func Uninstall() error {
	gitDir := ".git/hooks"

	hooks := []string{"pre-commit", "pre-push"}

	for _, name := range hooks {
		path := filepath.Join(gitDir, name)

		if err := os.Remove(path); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("remove %s: %w", name, err)
			}
		} else {
			fmt.Printf("✓ Removed %s\n", name)
		}
	}

	return nil
}
