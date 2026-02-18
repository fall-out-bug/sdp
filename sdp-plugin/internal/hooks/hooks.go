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
var baseHooks = []string{"pre-commit", "pre-push", "post-checkout", "post-merge"}
var provenanceHooks = []string{"commit-msg", "post-commit"}

// InstallOptions controls optional hook groups.
type InstallOptions struct {
	WithProvenance bool
}

func hooksForOptions(opts InstallOptions) []string {
	hooks := append([]string{}, baseHooks...)
	if opts.WithProvenance {
		hooks = append(hooks, provenanceHooks...)
	}
	return hooks
}

func allManagedHooks() []string {
	all := append([]string{}, baseHooks...)
	all = append(all, provenanceHooks...)
	return all
}

// Install installs SDP-managed git hooks from the hooks/ directory.
// AC4: Implements `sdp hooks install` command
// AC5: Hooks work in both main repo and worktrees
func Install() error {
	return InstallWithOptions(InstallOptions{})
}

// InstallWithOptions installs SDP-managed git hooks from the hooks/ directory.
func InstallWithOptions(opts InstallOptions) error {
	gitDir := ".git/hooks"
	hookNames := hooksForOptions(opts)

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
		return installEmbeddedHooks(gitDir, hookNames)
	}

	return installFromDirectory(gitDir, hooksSourceDir, hookNames)
}

// installFromDirectory installs hooks from a source directory.
func installFromDirectory(gitDir, sourceDir string, hookNames []string) error {
	for _, name := range hookNames {
		sourcePath := filepath.Join(sourceDir, name+".sh")
		targetPath := filepath.Join(gitDir, name)

		// Read source hook
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			fmt.Printf("⚠ Source hook %s not found, skipping\n", sourcePath)
			continue
		}

		hookContent := ensureManagedMarker(string(content))

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

func ensureManagedMarker(hookContent string) string {
	if strings.Contains(hookContent, sdpManagedMarker) {
		return hookContent
	}

	lines := strings.Split(hookContent, "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "#!") {
		withMarker := append([]string{lines[0], sdpManagedMarker}, lines[1:]...)
		return strings.Join(withMarker, "\n")
	}

	return "#!/bin/sh\n" + sdpManagedMarker + "\n" + hookContent
}

// installEmbeddedHooks installs hooks from embedded templates.
// This is used when the hooks/ directory is not available.
func installEmbeddedHooks(gitDir string, hookNames []string) error {
	embeddedHooks := map[string]string{
		"pre-commit":    getPreCommitTemplate(),
		"pre-push":      getPrePushTemplate(),
		"post-checkout": getPostCheckoutTemplate(),
		"post-merge":    getPostMergeTemplate(),
		"commit-msg":    getCommitMsgTemplate(),
		"post-commit":   getPostCommitTemplate(),
	}

	for _, name := range hookNames {
		content, ok := embeddedHooks[name]
		if !ok {
			continue
		}
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

	for _, name := range allManagedHooks() {
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
