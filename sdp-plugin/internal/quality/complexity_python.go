package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (c *Checker) checkPythonComplexity(result *ComplexityResult) (*ComplexityResult, error) {
	// Try using radon (Python complexity tool)
	cmd := exec.Command("radon", "cc", ".", "-a", "-s")
	cmd.Dir = c.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		// radon not available, do basic check
		return c.basicPythonComplexity(result)
	}

	// Parse radon output
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	totalCC := 0.0
	fileCount := 0
	maxCC := 0
	complexFiles := []FileComplexity{}

	for _, line := range lines {
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			// Format: FILE - LINE (CLASS/FUNCTION): TYPE complexity
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				lastPart := parts[len(parts)-1]
				if strings.Contains(lastPart, ":") {
					ccStr := strings.Split(lastPart, ":")[0]
					cc, err := strconv.Atoi(ccStr)
					if err == nil {
						totalCC += float64(cc)
						fileCount++
						if cc > maxCC {
							maxCC = cc
						}

						if cc > result.Threshold {
							// Extract filename
							fileName := parts[0]
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
		}
	}

	if fileCount > 0 {
		result.AverageCC = totalCC / float64(fileCount)
	}
	result.MaxCC = maxCC
	result.Passed = maxCC <= result.Threshold
	result.ComplexFiles = complexFiles

	return result, nil
}

func (c *Checker) basicPythonComplexity(result *ComplexityResult) (*ComplexityResult, error) {
	// Basic check: count lines and functions per file
	err := filepath.Walk(c.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		if !strings.HasSuffix(path, ".py") {
			return nil
		}

		// Skip test files and virtual env
		if strings.Contains(path, "test") || strings.Contains(path, "venv") || strings.Contains(path, ".venv") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		lines := strings.Split(string(content), "\n")
		loc := len(lines)

		// Very rough estimate: CC ≈ LOC/10 for simple code
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
