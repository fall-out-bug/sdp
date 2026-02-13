package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
)

func checkFilePermissions() CheckResult {
	// List of sensitive files to check
	sensitiveFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".sdp", "telemetry.jsonl"),
		".beads/beads.db",
		".oneshot",
	}

	insecureFiles := []string{}
	for _, path := range sensitiveFiles {
		info, err := os.Stat(path)
		if err != nil {
			// File doesn't exist, skip
			continue
		}

		// Check if file or directory
		if info.IsDir() {
			// Check files in directory
			entries, err := os.ReadDir(path)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				fullPath := filepath.Join(path, entry.Name())
				fileInfo, err := os.Stat(fullPath)
				if err != nil {
					continue
				}

				// Check permissions (should be 0600 for files)
				if fileInfo.Mode().Perm()&0077 != 0 {
					insecureFiles = append(insecureFiles, fullPath)
				}
			}
		} else {
			// Check single file permissions
			if info.Mode().Perm()&0077 != 0 {
				insecureFiles = append(insecureFiles, path)
			}
		}
	}

	if len(insecureFiles) > 0 {
		return CheckResult{
			Name:    "File Permissions",
			Status:  "warning",
			Message: fmt.Sprintf("Sensitive files have insecure permissions: %s (run 'chmod 0600 <file>' to fix)", strings.Join(insecureFiles, ", ")),
		}
	}

	return CheckResult{
		Name:    "File Permissions",
		Status:  "ok",
		Message: "All sensitive files have secure permissions",
	}
}

func checkProjectConfig() CheckResult {
	root, err := config.FindProjectRoot()
	if err != nil {
		return CheckResult{
			Name:    ".sdp/config.yml",
			Status:  "warning",
			Message: fmt.Sprintf("Could not find project root: %v", err),
		}
	}
	path := filepath.Join(root, ".sdp", "config.yml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return CheckResult{
			Name:    ".sdp/config.yml",
			Status:  "ok",
			Message: "No project config (using defaults)",
		}
	}
	cfg, err := config.Load(root)
	if err != nil {
		return CheckResult{
			Name:    ".sdp/config.yml",
			Status:  "error",
			Message: fmt.Sprintf("Invalid config: %v", err),
		}
	}
	if err := cfg.Validate(); err != nil {
		return CheckResult{
			Name:    ".sdp/config.yml",
			Status:  "error",
			Message: fmt.Sprintf("Validation failed: %v", err),
		}
	}
	return CheckResult{
		Name:    ".sdp/config.yml",
		Status:  "ok",
		Message: fmt.Sprintf("Valid (version %d)", cfg.Version),
	}
}
