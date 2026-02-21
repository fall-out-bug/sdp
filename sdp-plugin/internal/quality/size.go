package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
)

func (c *Checker) CheckFileSize() (*FileSizeResult, error) {
	// Load threshold from guard rules (AC6)
	threshold := 200 // default
	projectRoot, rootErr := config.FindProjectRoot()
	if rootErr == nil {
		guardRules, rulesErr := config.LoadGuardRules(projectRoot)
		if rulesErr == nil {
			// Find max-file-loc rule and get its threshold
			for _, rule := range guardRules.Rules {
				if rule.Enabled && rule.ID == "max-file-loc" {
					if maxLines, ok := rule.Config["max_lines"]; ok {
						switch v := maxLines.(type) {
						case int:
							threshold = v
						case float64:
							threshold = int(v)
						}
					}
					break
				}
			}
		}
	}

	result := &FileSizeResult{
		Threshold: threshold,
		Strict:    c.strictMode,
	}

	skipPrefixes := []string{"vendor/", "node_modules/", ".git/", "target/", "__pycache__/", ".venv/", "venv/", ".tmp/", "/sdp/"}
	if projectRoot, rootErr := config.FindProjectRoot(); rootErr == nil {
		if cfg, cfgErr := config.Load(projectRoot); cfgErr == nil && len(cfg.Quality.SizeExclude) > 0 {
			skipPrefixes = append(skipPrefixes, cfg.Quality.SizeExclude...)
		}
	}

	var totalLOC int
	var totalFiles int
	var sumLOC int

	// Walk the project directory
	err := filepath.Walk(c.projectPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return walkErr
		}

		// Skip certain directories (base + config size_exclude)
		for _, prefix := range skipPrefixes {
			if strings.Contains(path, prefix) {
				return nil
			}
		}

		// Check file extensions based on project type
		var shouldCheck bool
		switch c.projectType {
		case Python:
			shouldCheck = strings.HasSuffix(path, ".py")
		case Go:
			shouldCheck = strings.HasSuffix(path, ".go")
		case Java:
			shouldCheck = strings.HasSuffix(path, ".java")
		}

		if !shouldCheck {
			return nil
		}

		// Skip test files for size check (they can be longer)
		if strings.Contains(filepath.Base(path), "test") ||
			strings.Contains(filepath.Base(path), "_test") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		loc := len(lines)

		totalFiles++
		totalLOC += loc
		sumLOC += loc

		if loc > result.Threshold {
			// Make path relative to project path
			relPath, _ := filepath.Rel(c.projectPath, path)
			violation := FileViolation{
				File: relPath,
				LOC:  loc,
			}

			// In strict mode, violations are errors
			// In pragmatic mode, violations are warnings
			if c.strictMode {
				result.Violators = append(result.Violators, violation)
			} else {
				result.Warnings = append(result.Warnings, violation)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking project: %w", err)
	}

	if totalFiles > 0 {
		result.AverageLOC = totalLOC / totalFiles
	}

	result.TotalFiles = totalFiles

	// In strict mode, fail on violations
	// In pragmatic mode, always pass (warnings are OK)
	if c.strictMode {
		result.Passed = len(result.Violators) == 0
	} else {
		result.Passed = true
	}

	return result, nil
}
