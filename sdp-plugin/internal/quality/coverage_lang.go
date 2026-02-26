package quality

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/security"
)

func (c *Checker) coverageTimeout(cfgKey string, envKey string, fallback time.Duration) time.Duration {
	cfg, err := config.Load(c.projectPath)
	if err != nil {
		return config.TimeoutFromEnv(envKey, fallback)
	}
	switch cfgKey {
	case "python":
		return config.TimeoutFromConfigOrEnv(cfg.Timeouts.CoveragePython, envKey, fallback)
	case "go":
		return config.TimeoutFromConfigOrEnv(cfg.Timeouts.CoverageGo, envKey, fallback)
	case "list":
		return config.TimeoutFromConfigOrEnv(cfg.Timeouts.CoverageList, envKey, fallback)
	case "java":
		return config.TimeoutFromConfigOrEnv(cfg.Timeouts.CoverageJava, envKey, fallback)
	default:
		return fallback
	}
}

func (c *Checker) checkPythonCoverage(ctx context.Context, result *CoverageResult) (*CoverageResult, error) {
	result.ProjectType = "Python"

	// Check if .coverage file exists
	covFile := filepath.Join(c.projectPath, ".coverage")
	if _, err := os.Stat(covFile); os.IsNotExist(err) {
		// Try running pytest with coverage
		if ctx == nil {
			ctx = context.TODO()
		}
		timeout := c.coverageTimeout("python", "SDP_TIMEOUT_COVERAGE_PYTHON", 30*time.Second)
		runCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		cmd, errCmd := security.SafeCommand(runCtx, "pytest", "--cov", "--cov-report=term-missing")
		if errCmd != nil {
			result.Coverage = 0.0
			result.Passed = false
			result.Report = fmt.Sprintf("pytest not allowed or invalid args: %v", errCmd)
			return result, nil
		}
		cmd.Dir = c.projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Coverage = 0.0
			result.Passed = false
			result.Report = fmt.Sprintf("pytest failed: %v\n%s", err, truncateOutput(output, 500))
			return result, nil
		}

		// Parse output for coverage percentage
		outputStr := string(output)
		if strings.Contains(outputStr, "%") {
			lines := strings.Split(outputStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "TOTAL") && strings.Contains(line, "%") {
					fields := strings.Fields(line)
					for _, field := range fields {
						if strings.HasSuffix(field, "%") {
							covStr := strings.TrimSuffix(field, "%")
							cov, err := strconv.ParseFloat(covStr, 64)
							if err == nil {
								result.Coverage = cov
								result.Passed = cov >= result.Threshold
								return result, nil
							}
						}
					}
				}
			}
		}
	}

	// Try reading .coverage file
	if _, err := os.Stat(covFile); err == nil {
		jsonFile := filepath.Join(c.projectPath, "coverage.json")
		if data, err := os.ReadFile(jsonFile); err == nil {
			if cov := parseCoverageJSON(data); cov >= 0 {
				result.Coverage = cov
				result.Passed = cov >= result.Threshold
				return result, nil
			}
		}
	}

	// Default: assume no coverage run yet
	result.Coverage = 0.0
	result.Passed = false
	result.Report = "No coverage data found. Run tests with coverage enabled."

	return result, nil
}

func (c *Checker) checkGoCoverage(ctx context.Context, result *CoverageResult) (*CoverageResult, error) {
	result.ProjectType = "Go"

	if ctx == nil {
		ctx = context.TODO()
	}

	// Load coverage_exclude from config
	var excludePrefixes []string
	if projectRoot, err := config.FindProjectRoot(); err == nil {
		if cfg, err := config.Load(projectRoot); err == nil && len(cfg.Quality.CoverageExclude) > 0 {
			excludePrefixes = cfg.Quality.CoverageExclude
		}
	}

	// Build package list, optionally excluding configured paths
	testArgs := []string{"test", "-cover", "-coverprofile=coverage.out"}
	if len(excludePrefixes) > 0 {
		listTimeout := c.coverageTimeout("list", "SDP_TIMEOUT_COVERAGE_LIST", 10*time.Second)
		listCtx, listCancel := context.WithTimeout(ctx, listTimeout)
		defer listCancel()
		listCmd, errList := security.SafeCommand(listCtx, "go", "list", "./...")
		if errList != nil {
			testArgs = append(testArgs, "./...")
		} else {
			listCmd.Dir = c.projectPath
			listOut, listErr := listCmd.Output()
			if listErr != nil {
				testArgs = append(testArgs, "./...")
			} else {
				var pkgs []string
				for _, line := range strings.Split(strings.TrimSpace(string(listOut)), "\n") {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}
					excluded := false
					for _, prefix := range excludePrefixes {
						if strings.Contains(line, prefix) {
							excluded = true
							break
						}
					}
					if !excluded {
						pkgs = append(pkgs, line)
					}
				}
				if len(pkgs) > 0 {
					testArgs = append(testArgs, pkgs...)
				} else {
					testArgs = append(testArgs, "./...")
				}
			}
		}
	} else {
		testArgs = append(testArgs, "./...")
	}

	goTimeout := c.coverageTimeout("go", "SDP_TIMEOUT_COVERAGE_GO", 60*time.Second)
	runCtx, cancel := context.WithTimeout(ctx, goTimeout)
	defer cancel()

	cmd, errCmd := security.SafeCommand(runCtx, "go", testArgs...)
	if errCmd != nil {
		result.Coverage = 0.0
		result.Passed = false
		result.Report = fmt.Sprintf("go test not allowed or invalid args: %v", errCmd)
		return result, nil
	}
	cmd.Dir = c.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Coverage = 0.0
		result.Passed = false
		result.Report = fmt.Sprintf("Test execution failed: %s", string(output))
		return result, nil
	}

	// Parse coverage output
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	totalCoverage := 0.0
	count := 0

	for _, line := range lines {
		if strings.Contains(line, "coverage:") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if strings.HasSuffix(field, "%") && i > 0 {
					covStr := strings.TrimSuffix(field, "%")
					if cov, err := strconv.ParseFloat(covStr, 64); err == nil {
						totalCoverage += cov
						count++
					}
				}
			}
		}
	}

	if count > 0 {
		result.Coverage = totalCoverage / float64(count)
	} else {
		result.Coverage = 0.0
	}

	result.Passed = result.Coverage >= result.Threshold
	return result, nil
}

func (c *Checker) checkJavaCoverage(ctx context.Context, result *CoverageResult) (*CoverageResult, error) {
	result.ProjectType = "Java"

	if ctx == nil {
		ctx = context.TODO()
	}
	javaTimeout := c.coverageTimeout("java", "SDP_TIMEOUT_COVERAGE_JAVA", 30*time.Second)
	runCtx, cancel := context.WithTimeout(ctx, javaTimeout)
	defer cancel()

	cmd, errCmd := security.SafeCommand(runCtx, "mvn", "test")
	if errCmd != nil {
		result.Coverage = 0.0
		result.Passed = false
		result.Report = fmt.Sprintf("mvn not allowed or invalid args: %v", errCmd)
		return result, nil
	}
	cmd.Dir = c.projectPath
	if err := cmd.Run(); err != nil {
		result.Report = fmt.Sprintf("mvn test failed: %v", err)
		result.Coverage = 0.0
		result.Passed = false
		return result, nil
	}

	// Try to find jacoco.csv
	jacocoFile := filepath.Join(c.projectPath, "target/site/jacoco/jacoco.csv")
	if file, err := os.Open(jacocoFile); err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		totalLines := 0
		coveredLines := 0

		// Skip header
		if scanner.Scan() {
			// Header row
		}

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, ",")
			if len(fields) >= 7 {
				// INSTRUCTION_MISSED, INSTRUCTION_COVERED are at indices 4,5
				if missed, err1 := strconv.Atoi(fields[4]); err1 == nil {
					if covered, err2 := strconv.Atoi(fields[5]); err2 == nil {
						totalLines += missed + covered
						coveredLines += covered
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			result.Coverage = 0.0
			result.Passed = false
			result.Report = fmt.Sprintf("jacoco.csv read error: %v", err)
			return result, nil
		}

		if totalLines > 0 {
			result.Coverage = float64(coveredLines) / float64(totalLines) * 100
		} else {
			result.Coverage = 0.0
		}
	} else {
		result.Coverage = 0.0
		result.Report = "No JaCoCo coverage report found. Run 'mvn test' with jacoco plugin."
	}

	result.Passed = result.Coverage >= result.Threshold
	return result, nil
}

func truncateOutput(b []byte, maxLen int) string {
	s := string(b)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// parseCoverageJSON extracts percent_covered from JSON. Returns -1 if not found or invalid.
func parseCoverageJSON(data []byte) float64 {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return -1
	}
	if v, ok := m["percent_covered"]; ok {
		switch x := v.(type) {
		case float64:
			return x
		case int:
			return float64(x)
		}
	}
	if totals, ok := m["totals"].(map[string]interface{}); ok {
		if v, ok := totals["percent_covered"]; ok {
			switch x := v.(type) {
			case float64:
				return x
			case int:
				return float64(x)
			}
		}
	}
	return -1
}
