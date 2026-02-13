package quality

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (c *Checker) checkPythonCoverage(result *CoverageResult) (*CoverageResult, error) {
	result.ProjectType = "Python"

	// Check if .coverage file exists
	covFile := filepath.Join(c.projectPath, ".coverage")
	if _, err := os.Stat(covFile); os.IsNotExist(err) {
		// Try running pytest with coverage (with 30s timeout)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "pytest", "--cov", "--cov-report=term-missing")
		cmd.Dir = c.projectPath
		output, _ := cmd.CombinedOutput()

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
		// Parse .coveragerc or coverage.json if exists
		jsonFile := filepath.Join(c.projectPath, "coverage.json")
		if data, err := os.ReadFile(jsonFile); err == nil {
			// Parse JSON coverage report
			content := string(data)
			if strings.Contains(content, "percent_covered") {
				// Simple parse - look for "percent_covered": NUMBER
				lines := strings.Split(content, "\n")
				for _, line := range lines {
					if strings.Contains(line, "percent_covered") {
						parts := strings.Split(line, ":")
						if len(parts) == 2 {
							covStr := strings.TrimSpace(strings.Trim(strings.TrimSuffix(strings.TrimSuffix(parts[1], ","), "}"), "\""))
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

	// Default: assume no coverage run yet
	result.Coverage = 0.0
	result.Passed = false
	result.Report = "No coverage data found. Run tests with coverage enabled."

	return result, nil
}

func (c *Checker) checkGoCoverage(result *CoverageResult) (*CoverageResult, error) {
	result.ProjectType = "Go"

	// Run go test with coverage (with 30s timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-cover", "-coverprofile=coverage.out")
	cmd.Dir = c.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Test execution failed, but might still have coverage
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

func (c *Checker) checkJavaCoverage(result *CoverageResult) (*CoverageResult, error) {
	result.ProjectType = "Java"

	// Run mvn test with jacoco (with 30s timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "mvn", "test")
	cmd.Dir = c.projectPath
	_ = cmd.Run()

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
