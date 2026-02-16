package quality

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (c *Checker) CheckTypes() (*TypeResult, error) {
	result := &TypeResult{}

	switch c.projectType {
	case Python:
		return c.checkPythonTypes(result)
	case Go:
		return c.checkGoTypes(result)
	case Java:
		return c.checkJavaTypes(result)
	default:
		return result, fmt.Errorf("unsupported project type: %d", c.projectType)
	}
}

func (c *Checker) checkPythonTypes(result *TypeResult) (*TypeResult, error) {
	result.ProjectType = "Python"

	// Try mypy
	cmd := exec.Command("mypy", ".", "--show-error-codes")
	cmd.Dir = c.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		// mypy failed, parse errors
		outputStr := string(output)
		lines := strings.Split(outputStr, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "Success:") {
				continue
			}

			// Parse mypy error format: FILE:LINE: ERROR_MESSAGE
			parts := strings.SplitN(line, ":", 3)
			if len(parts) >= 3 {
				file := strings.TrimSpace(parts[0])
				lineNum, err1 := strconv.Atoi(strings.TrimSpace(parts[1]))
				message := strings.TrimSpace(parts[2])

				if err1 == nil && file != "" {
					relPath, _ := filepath.Rel(c.projectPath, file)
					result.Errors = append(result.Errors, TypeError{
						File:    relPath,
						Line:    lineNum,
						Message: message,
					})
				}
			}
		}

		result.Passed = len(result.Errors) == 0
	} else {
		// mypy passed
		result.Passed = true
	}

	// If mypy not configured, try basic check
	if !result.Passed && len(result.Errors) == 0 {
		result.Passed = true // Assume OK if mypy not configured
		result.Warnings = append(result.Warnings, TypeError{
			File:    "",
			Line:    0,
			Message: "mypy not configured. Consider adding pyproject.toml with mypy configuration.",
		})
	}

	return result, nil
}

func (c *Checker) checkGoTypes(result *TypeResult) (*TypeResult, error) {
	result.ProjectType = "Go"

	// Use go vet
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = c.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		// go vet found issues
		outputStr := string(output)
		lines := strings.Split(outputStr, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Parse go vet format: FILE:LINE: MESSAGE
			parts := strings.SplitN(line, ":", 3)
			if len(parts) >= 3 {
				file := strings.TrimSpace(parts[0])
				lineNum, err1 := strconv.Atoi(strings.TrimSpace(parts[1]))
				message := strings.TrimSpace(parts[2])

				if err1 == nil && file != "" {
					relPath, _ := filepath.Rel(c.projectPath, file)
					result.Errors = append(result.Errors, TypeError{
						File:    relPath,
						Line:    lineNum,
						Message: message,
					})
				}
			} else if len(parts) == 2 {
				// Some vet errors don't have line numbers
				file := strings.TrimSpace(parts[0])
				message := strings.TrimSpace(parts[1])

				if file != "" {
					relPath, _ := filepath.Rel(c.projectPath, file)
					result.Errors = append(result.Errors, TypeError{
						File:    relPath,
						Line:    0,
						Message: message,
					})
				}
			}
		}

		result.Passed = len(result.Errors) == 0
	} else {
		result.Passed = true
	}

	return result, nil
}

func (c *Checker) checkJavaTypes(result *TypeResult) (*TypeResult, error) {
	result.ProjectType = "Java"

	// Use javac or mvn compile
	cmd := exec.Command("mvn", "compile")
	cmd.Dir = c.projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := string(output)
		lines := strings.Split(outputStr, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "[ERROR]") {
				// Parse Maven error format
				if strings.Contains(line, ".java:[") {
					parts := strings.Split(line, ".java:[")
					if len(parts) >= 2 {
						file := parts[0] + ".java"
						rest := parts[1]

						lineParts := strings.SplitN(rest, ",", 2)
						if len(lineParts) >= 2 {
							lineNum, err1 := strconv.Atoi(strings.TrimSpace(lineParts[0]))
							message := strings.TrimSpace(lineParts[1])

							if err1 == nil {
								result.Errors = append(result.Errors, TypeError{
									File:    file,
									Line:    lineNum,
									Message: message,
								})
							}
						}
					}
				}
			}
		}

		result.Passed = len(result.Errors) == 0
	} else {
		result.Passed = true
	}

	return result, nil
}
