package quality

import (
	"os"
	"path/filepath"
	"strings"
)

type Type int

const (
	Python Type = iota
	Go
	Java
)

type Checker struct {
	projectPath string
	projectType Type
	strictMode  bool
}

type CoverageResult struct {
	ProjectType  string
	Coverage     float64
	Threshold    float64
	Passed       bool
	Report       string
	FilesBelow   []FileCoverage
}

type FileCoverage struct {
	File     string
	Coverage float64
}

type ComplexityResult struct {
	AverageCC    float64
	MaxCC        int
	Threshold    int
	Passed       bool
	ComplexFiles []FileComplexity
}

type FileComplexity struct {
	File            string
	AverageCC       float64
	MaxCC           int
	ExceedsThreshold bool
}

type FileSizeResult struct {
	TotalFiles int
	Violators  []FileViolation
	Warnings   []FileViolation
	Threshold  int
	Passed     bool
	AverageLOC int
	Strict     bool
}

type FileViolation struct {
	File    string
	LOC     int
}

type TypeResult struct {
	ProjectType  string
	Passed       bool
	Errors       []TypeError
	Warnings     []TypeError
}

type TypeError struct {
	File    string
	Line    int
	Message string
}

func NewChecker(projectPath string) (*Checker, error) {
	c := &Checker{
		projectPath: projectPath,
	}

	// Detect project type
	pt, err := c.detectProjectType()
	if err != nil {
		return nil, err
	}
	c.projectType = pt

	return c, nil
}

func (c *Checker) detectProjectType() (Type, error) {
	// Check for Python
	if _, err := os.Stat(filepath.Join(c.projectPath, "pyproject.toml")); err == nil {
		return Python, nil
	}
	if _, err := os.Stat(filepath.Join(c.projectPath, "setup.py")); err == nil {
		return Python, nil
	}
	if _, err := os.Stat(filepath.Join(c.projectPath, "requirements.txt")); err == nil {
		return Python, nil
	}

	// Check for Go
	if _, err := os.Stat(filepath.Join(c.projectPath, "go.mod")); err == nil {
		return Go, nil
	}

	// Check for Java
	if _, err := os.Stat(filepath.Join(c.projectPath, "pom.xml")); err == nil {
		return Java, nil
	}

	// Default: check by file extensions
	files, _ := os.ReadDir(c.projectPath)
	pythonCount := 0
	goCount := 0
	javaCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := strings.ToLower(file.Name())
		if strings.HasSuffix(name, ".py") {
			pythonCount++
		} else if strings.HasSuffix(name, ".go") {
			goCount++
		} else if strings.HasSuffix(name, ".java") {
			javaCount++
		}
	}

	// Return type with most files
	if pythonCount >= goCount && pythonCount >= javaCount && pythonCount > 0 {
		return Python, nil
	} else if goCount >= pythonCount && goCount >= javaCount && goCount > 0 {
		return Go, nil
	} else if javaCount > 0 {
		return Java, nil
	}

	// Fallback to Python
	return Python, nil
}

// SetStrictMode enables or disables strict quality gate checking
// In strict mode, file size violations are errors
// In pragmatic mode (default), file size violations are warnings
func (c *Checker) SetStrictMode(strict bool) {
	c.strictMode = strict
}

// IsStrictMode returns whether strict mode is enabled
func (c *Checker) IsStrictMode() bool {
	return c.strictMode
}
