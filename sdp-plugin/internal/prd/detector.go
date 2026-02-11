package prd

import (
	"os"
	"path/filepath"
	"strings"
)

// Detector analyzes project structure to detect project type
type Detector struct {
	projectPath string
}

// NewDetector creates a new project type detector
func NewDetector(projectPath string) *Detector {
	return &Detector{
		projectPath: projectPath,
	}
}

// DetectType determines the project type from file structure
func (d *Detector) DetectType() ProjectType {
	// Check for Go project
	if d.hasGoMod() {
		// Further classify Go projects
		if d.hasDockerCompose() {
			return Service
		}
		if d.hasMainPackage() {
			return CLI
		}
		return Go
	}

	// Check for Python project
	if d.hasPythonProject() {
		if d.hasDockerCompose() {
			return Service
		}
		if d.hasCLIEntryPoints() {
			return CLI
		}
		return Python
	}

	// Check for Java project
	if d.hasJavaProject() {
		if d.hasDockerCompose() {
			return Service
		}
		return Java
	}

	// Check for service with Docker
	if d.hasDockerCompose() {
		return Service
	}

	// Default: library
	return Library
}

// hasGoMod checks if go.mod exists
func (d *Detector) hasGoMod() bool {
	_, err := os.Stat(filepath.Join(d.projectPath, "go.mod"))
	return err == nil
}

// hasDockerCompose checks for docker-compose files
func (d *Detector) hasDockerCompose() bool {
	composeFiles := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
		"compose.yml",
		"compose.yaml",
	}

	for _, file := range composeFiles {
		if _, err := os.Stat(filepath.Join(d.projectPath, file)); err == nil {
			return true
		}
	}

	return false
}

// hasMainPackage checks if Go project has main package
func (d *Detector) hasMainPackage() bool {
	// Look for cmd/*/main.go or main.go
	mainPatterns := []string{
		"cmd/*/main.go",
		"main.go",
	}

	for _, pattern := range mainPatterns {
		matches, err := filepath.Glob(filepath.Join(d.projectPath, pattern))
		if err == nil && len(matches) > 0 {
			return true
		}
	}

	return false
}

// hasPythonProject checks if this is a Python project
func (d *Detector) hasPythonProject() bool {
	pythonIndicators := []string{
		"pyproject.toml",
		"setup.py",
		"requirements.txt",
		"Pipfile",
		"setup.cfg",
	}

	for _, indicator := range pythonIndicators {
		if _, err := os.Stat(filepath.Join(d.projectPath, indicator)); err == nil {
			return true
		}
	}

	return false
}

// hasCLIEntryPoints checks for Python CLI entry points
func (d *Detector) hasCLIEntryPoints() bool {
	// Check pyproject.toml for [project.scripts]
	pyproject := filepath.Join(d.projectPath, "pyproject.toml")
	if content, err := os.ReadFile(pyproject); err == nil {
		contentStr := string(content)
		if strings.Contains(contentStr, "[project.scripts]") ||
			strings.Contains(contentStr, "[tool.poetry.scripts]") {
			return true
		}
	}

	// Check for cli.py or main.py
	cliFiles := []string{"cli.py", "main.py"}
	for _, file := range cliFiles {
		if _, err := os.Stat(filepath.Join(d.projectPath, file)); err == nil {
			return true
		}
	}

	return false
}

// hasJavaProject checks if this is a Java project
func (d *Detector) hasJavaProject() bool {
	javaIndicators := []string{
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
	}

	for _, indicator := range javaIndicators {
		if _, err := os.Stat(filepath.Join(d.projectPath, indicator)); err == nil {
			return true
		}
	}

	return false
}
