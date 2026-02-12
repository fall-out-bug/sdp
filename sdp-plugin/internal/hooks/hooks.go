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
		"pre-commit": getPreCommitTemplate(),
		"pre-push":   getPrePushTemplate(),
		"post-checkout": getPostCheckoutTemplate(),
		"post-merge": getPostMergeTemplate(),
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

// Embedded hook templates for when hooks/ directory is not available

func getPreCommitTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Pre-commit
# Part of F065 - Agent Git Safety Protocol

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
fi

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

CURRENT_BRANCH=$(git branch --show-current)

# Session validation
if [ -f ".sdp/session.json" ] && command -v jq >/dev/null 2>&1; then
    EXPECTED_BRANCH=$(jq -r '.expected_branch' .sdp/session.json 2>/dev/null)
    if [ -n "$EXPECTED_BRANCH" ] && [ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]; then
        echo "ERROR: Branch mismatch! Expected: $EXPECTED_BRANCH, Current: $CURRENT_BRANCH"
        exit 1
    fi
fi

# Protected branch check
if [ -f ".sdp/session.json" ]; then
    case "$CURRENT_BRANCH" in
        main|dev)
            FEATURE_ID=$(jq -r '.feature_id' .sdp/session.json 2>/dev/null)
            if [ -n "$FEATURE_ID" ] && [ "$FEATURE_ID" != "null" ]; then
                echo "ERROR: Cannot commit to $CURRENT_BRANCH for feature $FEATURE_ID"
                exit 1
            fi
            ;;
    esac
fi

exit 0
`
}

func getPrePushTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Pre-push
# Part of F065 - Agent Git Safety Protocol

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
fi

set -e

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

CURRENT_BRANCH=$(git branch --show-current)

# Prevent pushing to protected branches
case "$CURRENT_BRANCH" in
    main|dev)
        echo "ERROR: Direct push to $CURRENT_BRANCH is not allowed!"
        echo "Create a feature branch and use PR workflow."
        exit 1
        ;;
esac

exit 0
`
}

func getPostCheckoutTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Post-checkout
# Part of F065 - Agent Git Safety Protocol

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
fi

# Only update session on branch checkout
if [ "$3" != "1" ]; then
    exit 0
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
cd "$REPO_ROOT"

if [ -f ".sdp/session.json" ] && command -v jq >/dev/null 2>&1; then
    NEW_BRANCH=$(git branch --show-current)
    TEMP_FILE=$(mktemp)
    jq --arg branch "$NEW_BRANCH" '.expected_branch = $branch' .sdp/session.json > "$TEMP_FILE" 2>/dev/null
    if [ $? -eq 0 ]; then
        mv "$TEMP_FILE" .sdp/session.json
        echo "Session updated: now on branch $NEW_BRANCH"
    else
        rm -f "$TEMP_FILE"
    fi
fi

exit 0
`
}

func getPostMergeTemplate() string {
	return `#!/bin/sh
` + sdpManagedMarker + `
# SDP Git Hook - Post-merge
# Runs after a git merge completes

# Check if sdp binary exists
if ! command -v sdp >/dev/null 2>&1; then
    echo "Warning: SDP CLI (sdp) not found in PATH"
    echo "Install SDP to enable quality checks"
    exit 0
fi

# Check if .sdp/ exists
if [ -d ".sdp" ]; then
    echo "SDP: Post-merge checks..."
    # Add post-merge validation here
fi

exit 0
`
}
