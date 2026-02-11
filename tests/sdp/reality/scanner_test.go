package reality

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/reality"
)

func TestDetectLanguage(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Test Go detection
	goDir := filepath.Join(tmpDir, "go_project")
	os.MkdirAll(goDir, 0755)
	os.WriteFile(filepath.Join(goDir, "go.mod"), []byte("module example.com"), 0644)

	lang, framework := reality.DetectLanguage(goDir)
	if lang != "go" {
		t.Errorf("Expected 'go', got '%s'", lang)
	}
	if framework != "" {
		t.Errorf("Expected empty framework, got '%s'", framework)
	}

	// Test Python detection
	pyDir := filepath.Join(tmpDir, "python_project")
	os.MkdirAll(pyDir, 0755)
	os.WriteFile(filepath.Join(pyDir, "pyproject.toml"), []byte("[tool.poetry]"), 0644)

	lang, framework = reality.DetectLanguage(pyDir)
	if lang != "python" {
		t.Errorf("Expected 'python', got '%s'", lang)
	}
	if framework != "" {
		t.Errorf("Expected empty framework, got '%s'", framework)
	}

	// Test Java detection
	javaDir := filepath.Join(tmpDir, "java_project")
	os.MkdirAll(javaDir, 0755)
	os.WriteFile(filepath.Join(javaDir, "pom.xml"), []byte("<project></project>"), 0644)

	lang, framework = reality.DetectLanguage(javaDir)
	if lang != "java" {
		t.Errorf("Expected 'java', got '%s'", lang)
	}
	if framework != "spring" {
		t.Errorf("Expected 'spring', got '%s'", framework)
	}

	// Test Node.js detection
	nodeDir := filepath.Join(tmpDir, "node_project")
	os.MkdirAll(nodeDir, 0755)
	os.WriteFile(filepath.Join(nodeDir, "package.json"), []byte(`{"name": "test"}`), 0644)

	lang, framework = reality.DetectLanguage(nodeDir)
	if lang != "nodejs" {
		t.Errorf("Expected 'nodejs', got '%s'", lang)
	}
	if framework != "" {
		t.Errorf("Expected empty framework, got '%s'", framework)
	}
}

func TestDetectFramework(t *testing.T) {
	// Test Gin (Go)
	tmpDir := t.TempDir()
	ginDir := filepath.Join(tmpDir, "gin_project")
	os.MkdirAll(ginDir, 0755)
	os.WriteFile(filepath.Join(ginDir, "go.mod"), []byte("module example.com\n\ngithub.com/gin-gonic/gin v1.9.1"), 0644)

	_, framework := reality.DetectLanguage(ginDir)
	if framework != "gin" {
		t.Errorf("Expected 'gin', got '%s'", framework)
	}

	// Test Django (Python)
	djangoDir := filepath.Join(tmpDir, "django_project")
	os.MkdirAll(djangoDir, 0755)
	os.WriteFile(filepath.Join(djangoDir, "pyproject.toml"), []byte("[tool.poetry]\ndependencies = ['django']"), 0644)

	_, framework = reality.DetectLanguage(djangoDir)
	if framework != "django" {
		t.Errorf("Expected 'django', got '%s'", framework)
	}
}

func TestScanProject(t *testing.T) {
	// Create test project structure
	tmpDir := t.TempDir()

	// Create some test files
	os.MkdirAll(filepath.Join(tmpDir, "src"), 0755)
	os.MkdirAll(filepath.Join(tmpDir, "tests"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "src", "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module example.com"), 0644)

	scanner := reality.NewProjectScanner(tmpDir)
	result, err := scanner.Scan()

	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Language != "go" {
		t.Errorf("Expected language 'go', got '%s'", result.Language)
	}

	if result.TotalFiles != 2 { // go.mod, src/main.go (dirs don't count as files)
		t.Errorf("Expected 2 files, got %d", result.TotalFiles)
	}

	if result.LinesOfCode < 3 { // Should count lines in main.go
		t.Errorf("Expected LOC >= 3, got %d", result.LinesOfCode)
	}
}

func TestCountLinesOfCode(t *testing.T) {
	// Create test file with known line count
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	content := "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n"
	os.WriteFile(testFile, []byte(content), 0644)

	scanner := reality.NewProjectScanner(tmpDir)
	result, err := scanner.Scan()

	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should count 4 non-empty, non-comment lines
	if result.LinesOfCode != 4 {
		t.Errorf("Expected 4 lines, got %d", result.LinesOfCode)
	}
}
