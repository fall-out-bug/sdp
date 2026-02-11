package hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const sdpManagedMarker = "# SDP-MANAGED-HOOK"

const hookTemplate = `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook
# This hook is managed by SDP. Do not edit manually.

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks: https://github.com/fall-out-bug/sdp"
    exit 0
fi

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
		"pre-commit":    hookTemplate,
		"pre-push":      hookTemplate,
		"post-merge":    hookTemplate,
		"post-checkout": hookTemplate,
	}

	for name, content := range hooks {
		path := filepath.Join(gitDir, name)

		// Check if hook already exists
		existingContent, err := os.ReadFile(path)
		if err == nil {
			// Hook exists - check if it's SDP-managed
			if strings.Contains(string(existingContent), sdpManagedMarker) {
				// Update existing SDP-managed hook
				if err := os.WriteFile(path, []byte(content), 0755); err != nil {
					return fmt.Errorf("update %s: %w", name, err)
				}
				fmt.Printf("✓ Updated %s\n", name)
				continue
			}
			// Non-SDP hook - skip it
			fmt.Printf("⚠ %s already exists (not SDP-managed), skipping\n", name)
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

	hooks := []string{"pre-commit", "pre-push", "post-merge", "post-checkout"}

	for _, name := range hooks {
		path := filepath.Join(gitDir, name)

		// Check if hook exists
		content, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				// Hook doesn't exist, skip it
				continue
			}
			return fmt.Errorf("read %s: %w", name, err)
		}

		// Check if it's an SDP-managed hook
		if !strings.Contains(string(content), sdpManagedMarker) {
			fmt.Printf("⚠ %s exists but not SDP-managed, skipping\n", name)
			continue
		}

		// Remove the SDP-managed hook
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove %s: %w", name, err)
		}
		fmt.Printf("✓ Removed %s\n", name)
	}

	return nil
}
