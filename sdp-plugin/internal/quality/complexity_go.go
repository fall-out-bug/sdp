package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
)

func (c *Checker) complexityExcludePrefixes() []string {
	base := []string{".tmp/", "/sdp/", "vendor/", "node_modules/"}
	if projectRoot, err := config.FindProjectRoot(); err == nil {
		if cfg, err := config.Load(projectRoot); err == nil && len(cfg.Quality.ComplexityExclude) > 0 {
			return append(base, cfg.Quality.ComplexityExclude...)
		}
	}
	return base
}

func (c *Checker) shouldExcludeComplexity(fileName string) bool {
	for _, prefix := range c.complexityExcludePrefixes() {
		if strings.Contains(fileName, prefix) {
			return true
		}
	}
	return false
}

func (c *Checker) checkGoComplexity(result *ComplexityResult) (*ComplexityResult, error) {
	// Use gocyclo if available
	cmd := exec.Command("gocyclo", "-over", "10", ".")
	cmd.Dir = c.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		// gocyclo not available, do basic check
		return c.basicGoComplexity(result)
	}

	// Parse gocyclo output
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	totalCC := 0
	fileCount := 0
	maxCC := 0
	complexFiles := []FileComplexity{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) >= 2 {
			fileName := parts[len(parts)-1]
			if c.shouldExcludeComplexity(fileName) {
				continue
			}

			cc, err := strconv.Atoi(parts[0])
			if err == nil {
				totalCC += cc
				fileCount++
				if cc > maxCC {
					maxCC = cc
				}

				if cc > result.Threshold {
					complexFiles = append(complexFiles, FileComplexity{
						File:             fileName,
						AverageCC:        float64(cc),
						MaxCC:            cc,
						ExceedsThreshold: true,
					})
				}
			}
		}
	}

	if fileCount > 0 {
		result.AverageCC = float64(totalCC) / float64(fileCount)
	}
	result.MaxCC = maxCC
	result.Passed = maxCC <= result.Threshold
	result.ComplexFiles = complexFiles

	return result, nil
}

func (c *Checker) basicGoComplexity(result *ComplexityResult) (*ComplexityResult, error) {
	// Basic check similar to Python
	err := filepath.Walk(c.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip generated, test, and excluded paths
		if strings.Contains(path, "generated") || strings.Contains(path, "_test.go") {
			return nil
		}
		if c.shouldExcludeComplexity(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		loc := len(lines)

		estimatedCC := loc / 10
		if estimatedCC > result.Threshold {
			result.ComplexFiles = append(result.ComplexFiles, FileComplexity{
				File:             path,
				AverageCC:        float64(estimatedCC),
				MaxCC:            estimatedCC,
				ExceedsThreshold: true,
			})
		}

		if estimatedCC > result.MaxCC {
			result.MaxCC = estimatedCC
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result.Passed = len(result.ComplexFiles) == 0
	return result, nil
}
