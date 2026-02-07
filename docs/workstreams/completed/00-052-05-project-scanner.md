# 00-052-05: Project Scanner Implementation

> **Beads ID:** sdp-bsds
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 1B - Analysis Skills (@reality)
> **Size:** MEDIUM
> **Duration:** 2-3 days
> **Dependencies:**
> - 00-052-04 (@reality Skill Structure)

## Goal

Implement project scanner that detects language, framework, and project structure.

## Acceptance Criteria

- **AC1:** `src/sdp/reality/scanner.go` created with ProjectScanner struct
- **AC2:** `DetectLanguage()` identifies Go, Python, JS, TS, Rust, Java
- **AC3:** `DetectFramework()` identifies common frameworks (React, Django, Gin, etc.)
- **AC4:** `ScanProject()` returns structure overview (directories, packages, files)
- **AC5:** `tests/sdp/reality/scanner_test.go` with ≥80% coverage

## Files

**Create:**
- `src/sdp/reality/scanner.go` - ProjectScanner implementation
- `tests/sdp/reality/scanner_test.go` - Test suite

**Modify:**
- None

## Steps

### Step 1: Implement Project Scanner

Create `src/sdp/reality/scanner.go`:

```go
package reality

import (
	"os"
	"path/filepath"
	"strings"
)

// ProjectType represents the detected project type
type ProjectType struct {
	Language   string   // go, python, javascript, typescript, rust, java
	Frameworks []string // django, react, gin, etc.
	Type       string   // service, web, cli, library
}

// ProjectStats represents project statistics
type ProjectStats struct {
	TotalFiles  int
	TotalLines  int
	LangBreakdown map[string]int // language -> LOC
}

// ProjectScanner scans and analyzes project structure
type ProjectScanner struct {
	projectPath string
}

// NewProjectScanner creates a new scanner
func NewProjectScanner(projectPath string) *ProjectScanner {
	return &ProjectScanner{
		projectPath: projectPath,
	}
}

// DetectLanguage identifies primary programming language
func (s *ProjectScanner) DetectLanguage() (string, error) {
	indicators := map[string][]string{
		"go":        {"go.mod", "*.go"},
		"python":    {"requirements.txt", "setup.py", "pyproject.toml", "*.py"},
		"javascript": {"package.json", "*.js"},
		"typescript": {"tsconfig.json", "*.ts"},
		"rust":      {"Cargo.toml", "*.rs"},
		"java":      {"pom.xml", "build.gradle", "*.java"},
	}

	// Check for indicator files
	for lang, files := range indicators {
		for _, file := range files {
			if strings.Contains(file, "*") {
				// Check for file extension
				matches, _ := filepath.Glob(filepath.Join(s.projectPath, file))
				if len(matches) > 0 {
					return lang, nil
				}
			} else {
				// Check for specific file
				if _, err := os.Stat(filepath.Join(s.projectPath, file)); err == nil {
					return lang, nil
				}
			}
		}
	}

	return "unknown", nil
}

// DetectFramework identifies web frameworks
func (s *ProjectScanner) DetectFramework(language string) ([]string, error) {
	frameworks := []string{}

	lang, err := s.DetectLanguage()
	if err != nil {
		return nil, err
	}

	switch lang {
	case "go":
		if s.fileExists("go.mod") {
			content, _ := os.ReadFile(filepath.Join(s.projectPath, "go.mod"))
			contentStr := string(content)
			if strings.Contains(contentStr, "gin") {
				frameworks = append(frameworks, "gin")
			}
			if strings.Contains(contentStr, "echo") {
				frameworks = append(frameworks, "echo")
			}
		}
	case "python":
		if s.fileExists("requirements.txt") {
			content, _ := os.ReadFile(filepath.Join(s.projectPath, "requirements.txt"))
			contentStr := string(content)
			if strings.Contains(contentStr, "django") {
				frameworks = append(frameworks, "django")
			}
			if strings.Contains(contentStr, "flask") {
				frameworks = append(frameworks, "flask")
			}
			if strings.Contains(contentStr, "fastapi") {
				frameworks = append(frameworks, "fastapi")
			}
		}
	case "javascript", "typescript":
		if s.fileExists("package.json") {
			content, _ := os.ReadFile(filepath.Join(s.projectPath, "package.json"))
			contentStr := string(content)
			if strings.Contains(contentStr, "react") {
				frameworks = append(frameworks, "react")
			}
			if strings.Contains(contentStr, "vue") {
				frameworks = append(frameworks, "vue")
			}
			if strings.Contains(contentStr, "express") {
				frameworks = append(frameworks, "express")
			}
		}
	}

	return frameworks, nil
}

// ScanProject analyzes project structure
func (s *ProjectScanner) ScanProject() (*ProjectType, *ProjectStats, error) {
	lang, err := s.DetectLanguage()
	if err != nil {
		return nil, nil, err
	}

	frameworks, err := s.DetectFramework(lang)
	if err != nil {
		return nil, nil, err
	}

	stats, err := s.calculateStats()
	if err != nil {
		return nil, nil, err
	}

	projectType := &ProjectType{
		Language:   lang,
		Frameworks: frameworks,
		Type:       s.inferType(lang, frameworks),
	}

	return projectType, stats, nil
}

// inferType guesses project type from lang/framework
func (s *ProjectScanner) inferType(lang string, frameworks []string) string {
	for _, fw := range frameworks {
		if fw == "react" || fw == "vue" {
			return "web"
		}
		if fw == "gin" || fw == "express" || fw == "django" {
			return "service"
		}
	}

	if s.fileExists("cmd") || s.fileExists("main.go") || s.fileExists("main.py") {
		return "cli"
	}

	return "library"
}

// calculateStats counts lines of code
func (s *ProjectScanner) calculateStats() (*ProjectStats, error) {
	stats := &ProjectStats{
		LangBreakdown: make(map[string]int),
	}

	filepath.Walk(s.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		// Skip vendor, node_modules, etc.
		if strings.Contains(path, "vendor") || strings.Contains(path, "node_modules") {
			return nil
		}

		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		if ext == "go" || ext == "py" || ext == "js" || ext == "ts" {
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			lines := strings.Count(string(content), "\n")
			stats.TotalFiles++
			stats.TotalLines += lines
			stats.LangBreakdown[ext] += lines
		}

		return nil
	})

	return stats, nil
}

// fileExists checks if file exists in project root
func (s *ProjectScanner) fileExists(name string) bool {
	_, err := os.Stat(filepath.Join(s.projectPath, name))
	return err == nil
}
```

### Step 2: Write Tests

Create `tests/sdp/reality/scanner_test.go`:

```go
package reality

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/reality"
)

func TestDetectLanguage(t *testing.T) {
	// Test Go project
	goDir := t.TempDir()
	os.WriteFile(filepath.Join(goDir, "go.mod"), []byte("module test\n\ngo 1.21"), 0644)

	scanner := reality.NewProjectScanner(goDir)
	lang, err := scanner.DetectLanguage()
	if err != nil {
		t.Fatalf("DetectLanguage failed: %v", err)
	}

	if lang != "go" {
		t.Errorf("Expected 'go', got '%s'", lang)
	}
}

func TestDetectFramework(t *testing.T) {
	// Test Django project
	pyDir := t.TempDir()
	os.WriteFile(filepath.Join(pyDir, "requirements.txt"), []byte("django==4.0\n"), 0644)
	os.WriteFile(filepath.Join(pyDir, "manage.py"), []byte("# Django"), 0644)

	scanner := reality.NewProjectScanner(pyDir)
	lang, _ := scanner.DetectLanguage()
	frameworks, err := scanner.DetectFramework(lang)
	if err != nil {
		t.Fatalf("DetectFramework failed: %v", err)
	}

	if len(frameworks) == 0 || frameworks[0] != "django" {
		t.Errorf("Expected 'django', got %v", frameworks)
	}
}

func TestScanProject(t *testing.T) {
	// Create test Go project
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n\ngo 1.21"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644)

	scanner := reality.NewProjectScanner(tmpDir)
	projectType, stats, err := scanner.ScanProject()
	if err != nil {
		t.Fatalf("ScanProject failed: %v", err)
	}

	if projectType.Language != "go" {
		t.Errorf("Expected language 'go', got '%s'", projectType.Language)
	}

	if stats.TotalFiles < 1 {
		t.Errorf("Expected at least 1 file, got %d", stats.TotalFiles)
	}
}
```

### Step 3: Verify Coverage

```bash
go test -v -coverprofile=coverage.out ./tests/sdp/reality
go tool cover -func=coverage.out | grep scanner.go
```

Target: ≥80% coverage

## Quality Gates

- Files < 200 LOC each
- Coverage ≥80%
- No `go vet` warnings
- Tests use t.TempDir() for cleanup
- No hardcoded paths

## Success Metrics

- Detects all 6 target languages
- Detects common frameworks (≥5)
- Accurate line counting
- Works on any project directory
