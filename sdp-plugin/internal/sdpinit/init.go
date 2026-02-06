package sdpinit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	ProjectType string
	SkipBeads   bool
}

func Run(cfg Config) error {
	// Create .claude/ directory
	claudeDir := ".claude"
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("create .claude/: %w", err)
	}

	// Create subdirectories
	dirs := []string{
		"skills",
		"agents",
		"validators",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(claudeDir, dir), 0755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}

	// Copy prompts from prompts/ directory
	if err := copyPrompts(claudeDir); err != nil {
		return fmt.Errorf("copy prompts: %w", err)
	}

	// Create settings.json
	if err := createSettings(claudeDir, cfg); err != nil {
		return fmt.Errorf("create settings: %w", err)
	}

	fmt.Println("✓ SDP initialized in .claude/")
	fmt.Printf("  Project type: %s\n", cfg.ProjectType)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review .claude/settings.json")
	fmt.Println("  2. Start using Claude Code with SDP prompts")

	return nil
}

func copyPrompts(destDir string) error {
	promptsDir := "prompts"

	// Check if prompts directory exists
	if _, err := os.Stat(promptsDir); os.IsNotExist(err) {
		return fmt.Errorf("prompts directory not found: %s", promptsDir)
	}

	// Walk the prompts directory and copy to .claude/
	return filepath.Walk(promptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(promptsDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := srcFile.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to close source file %s: %v\n", path, cerr)
			}
		}()

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := dstFile.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to close destination file %s: %v\n", destPath, cerr)
			}
		}()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

func createSettings(claudeDir string, cfg Config) error {
	settings := fmt.Sprintf(`{
  "skills": [
    "feature",
    "idea",
    "design",
    "build",
    "review",
    "deploy",
    "debug",
    "bugfix",
    "hotfix",
    "oneshot"
  ],
  "projectType": "%s",
  "sdpVersion": "1.0.0"
}`, cfg.ProjectType)

	return os.WriteFile(
		filepath.Join(claudeDir, "settings.json"),
		[]byte(settings),
		0600,
	)
}
