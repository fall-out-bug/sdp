package hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const sdpManagedMarker = "# SDP-MANAGED-HOOK"

// Hooks to install (AC1-AC3): pre-commit, pre-push, post-checkout
// Also include post-merge for backward compatibility
var managedHooks = []string{"pre-commit", "pre-push", "post-checkout", "post-merge"}

// Install installs SDP-managed git hooks from the hooks/ directory.
// AC4: Implements `sdp hooks install` command
// AC5: Hooks work in both main repo and worktrees
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

	// Try to find hooks source directory
	// First check for local hooks/ directory (development)
	hooksSourceDir := "hooks"
	if _, err := os.Stat(hooksSourceDir); os.IsNotExist(err) {
		// Fall back to embedded hooks
		return installEmbeddedHooks(gitDir)
	}

	return installFromDirectory(gitDir, hooksSourceDir)
}

// installFromDirectory installs hooks from a source directory.
func installFromDirectory(gitDir, sourceDir string) error {
	for _, name := range managedHooks {
		sourcePath := filepath.Join(sourceDir, name+".sh")
		targetPath := filepath.Join(gitDir, name)

		// Read source hook
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			fmt.Printf("⚠ Source hook %s not found, skipping\n", sourcePath)
			continue
		}

		// Add SDP marker if not present
		hookContent := string(content)
		if !strings.Contains(hookContent, sdpManagedMarker) {
			hookContent = "#!/bin/bash\n" + sdpManagedMarker + "\n" + hookContent
		}

		// Check if hook already exists
		existingContent, err := os.ReadFile(targetPath)
		if err == nil {
			// Hook exists - check if it's SDP-managed
			if strings.Contains(string(existingContent), sdpManagedMarker) {
				// Update existing SDP-managed hook
				if err := os.WriteFile(targetPath, []byte(hookContent), 0755); err != nil {
					return fmt.Errorf("update %s: %w", name, err)
				}
				fmt.Printf("✓ Updated %s\n", name)
				continue
			}
			// Non-SDP hook - skip it
			fmt.Printf("⚠ %s already exists (not SDP-managed), skipping\n", name)
			continue
		}

		if err := os.WriteFile(targetPath, []byte(hookContent), 0755); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
		fmt.Printf("✓ Installed %s\n", name)
	}

	fmt.Println("\nGit hooks installed!")
	fmt.Println("These hooks provide:")
	fmt.Println("  - Session validation before commits")
	fmt.Println("  - Branch tracking validation before pushes")
	fmt.Println("  - Session updates when switching branches")
	fmt.Println("")
	fmt.Println("Customize hooks in .git/hooks/ if needed")

	return nil
}

// installEmbeddedHooks installs hooks from embedded templates.
// This is used when the hooks/ directory is not available.
func installEmbeddedHooks(gitDir string) error {
	embeddedHooks := map[string]string{
		"pre-commit":    getPreCommitTemplate(),
		"pre-push":      getPrePushTemplate(),
		"post-checkout": getPostCheckoutTemplate(),
		"post-merge":    getPostMergeTemplate(),
	}

	for name, content := range embeddedHooks {
		targetPath := filepath.Join(gitDir, name)

		// Check if hook already exists
		existingContent, err := os.ReadFile(targetPath)
		if err == nil {
			// Hook exists - check if it's SDP-managed
			if strings.Contains(string(existingContent), sdpManagedMarker) {
				// Update existing SDP-managed hook
				if err := os.WriteFile(targetPath, []byte(content), 0755); err != nil {
					return fmt.Errorf("update %s: %w", name, err)
				}
				fmt.Printf("✓ Updated %s\n", name)
				continue
			}
			// Non-SDP hook - skip it
			fmt.Printf("⚠ %s already exists (not SDP-managed), skipping\n", name)
			continue
		}

		if err := os.WriteFile(targetPath, []byte(content), 0755); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
		fmt.Printf("✓ Installed %s\n", name)
	}

	fmt.Println("\nGit hooks installed (embedded version)!")
	fmt.Println("For full functionality, copy hooks from hooks/ directory")

	return nil
}

func Uninstall() error {
	gitDir := ".git/hooks"

	for _, name := range managedHooks {
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
