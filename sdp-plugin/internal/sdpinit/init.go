package sdpinit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Config holds initialization configuration options.
type Config struct {
	ProjectType string
	SkipBeads   bool
	// Skills to enable (empty = use defaults)
	Skills []string
	// NoEvidence disables evidence logging
	NoEvidence bool
	// ProjectName for display purposes
	ProjectName string
	// Interactive mode flag
	Interactive bool
	// Headless mode for CI/CD
	Headless bool
	// Output format (text, json)
	Output string
	// Force overwrites existing files
	Force bool
	// DryRun previews changes without writing
	DryRun bool
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

	// In headless mode, don't print text output
	if !cfg.Headless {
		fmt.Println("✓ SDP initialized in .claude/")
		fmt.Printf("  Project type: %s\n", cfg.ProjectType)
		fmt.Println("\nNext steps:")
		fmt.Println("  1. Review .claude/settings.json")
		fmt.Println("  2. Start using Claude Code with SDP prompts")
	}

	return nil
}

func copyPrompts(destDir string) error {
	promptsDir, err := resolvePromptsDir()
	if err != nil {
		return err
	}

	promptsAbs, err := filepath.Abs(promptsDir)
	if err != nil {
		return err
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

		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			if isSymlink(pathJoinClean(destDir, relPath)) {
				return filepath.SkipDir
			}
			return os.MkdirAll(destPath, 0755)
		}

		same, sameErr := sameResolvedPath(path, destPath)
		if sameErr == nil && same {
			return nil
		}

		destAbs, absErr := filepath.Abs(destPath)
		if absErr == nil {
			if strings.HasPrefix(destAbs, promptsAbs+string(os.PathSeparator)) {
				return nil
			}
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
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

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

func sameResolvedPath(a, b string) (bool, error) {
	ra, err := filepath.EvalSymlinks(a)
	if err != nil {
		return false, err
	}
	rb, err := filepath.EvalSymlinks(b)
	if err != nil {
		return false, err
	}
	raAbs, err := filepath.Abs(ra)
	if err != nil {
		return false, err
	}
	rbAbs, err := filepath.Abs(rb)
	if err != nil {
		return false, err
	}
	return raAbs == rbAbs, nil
}

func pathJoinClean(parts ...string) string {
	return filepath.Clean(filepath.Join(parts...))
}

func resolvePromptsDir() (string, error) {
	if envDir, err := envPromptsSourceDir(); err != nil {
		return "", err
	} else if envDir != "" {
		return envDir, nil
	}

	for _, dir := range localPromptsCandidates() {
		if isValidPromptsDir(dir) {
			return dir, nil
		}
	}

	downloadedDir, err := downloadPromptsToCache()
	if err != nil {
		return "", fmt.Errorf("prompts not found locally and remote fetch failed: %w", err)
	}
	return downloadedDir, nil
}

func createSettings(claudeDir string, cfg Config) error {
	// Get defaults for the project type
	defaults := MergeDefaults(cfg.ProjectType, &cfg)

	// Build skills list
	skills := defaults.Skills
	if len(cfg.Skills) > 0 {
		skills = cfg.Skills
	}

	// Build settings JSON
	settings := fmt.Sprintf(`{
  "skills": %s,
  "projectType": "%s",
  "evidence": {
    "enabled": %t
  },
  "sdpVersion": "1.0.0"
}`, formatStringsAsJSON(skills), cfg.ProjectType, defaults.EvidenceEnabled)

	return os.WriteFile(
		filepath.Join(claudeDir, "settings.json"),
		[]byte(settings),
		0600,
	)
}

// formatStringsAsJSON formats a string slice as a JSON array.
func formatStringsAsJSON(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	result := "[\n"
	for i, item := range items {
		result += fmt.Sprintf("    \"%s\"", item)
		if i < len(items)-1 {
			result += ","
		}
		result += "\n"
	}
	result += "  ]"
	return result
}
