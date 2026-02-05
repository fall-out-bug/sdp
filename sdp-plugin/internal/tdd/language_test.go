package tdd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectLanguageByGoMod(t *testing.T) {
	// Create temporary directory with go.mod
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	runner := &Runner{}
	lang, err := runner.DetectLanguage(tmpDir)

	if err != nil {
		t.Errorf("Expected to detect Go project, got error: %v", err)
	}

	if lang != Go {
		t.Errorf("Expected language Go, got %v", lang)
	}
}

func TestDetectLanguageByPyprojectToml(t *testing.T) {
	// Create temporary directory with pyproject.toml
	tmpDir := t.TempDir()
	pyprojectPath := filepath.Join(tmpDir, "pyproject.toml")
	err := os.WriteFile(pyprojectPath, []byte("[tool.poetry]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create pyproject.toml: %v", err)
	}

	runner := &Runner{}
	lang, err := runner.DetectLanguage(tmpDir)

	if err != nil {
		t.Errorf("Expected to detect Python project, got error: %v", err)
	}

	if lang != Python {
		t.Errorf("Expected language Python, got %v", lang)
	}
}

func TestDetectLanguageByPomXml(t *testing.T) {
	// Create temporary directory with pom.xml
	tmpDir := t.TempDir()
	pomPath := filepath.Join(tmpDir, "pom.xml")
	err := os.WriteFile(pomPath, []byte("<project></project>\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create pom.xml: %v", err)
	}

	runner := &Runner{}
	lang, err := runner.DetectLanguage(tmpDir)

	if err != nil {
		t.Errorf("Expected to detect Java project, got error: %v", err)
	}

	if lang != Java {
		t.Errorf("Expected language Java, got %v", lang)
	}
}

func TestDetectLanguageNoProjectFiles(t *testing.T) {
	// Create empty temporary directory
	tmpDir := t.TempDir()

	runner := &Runner{}
	_, err := runner.DetectLanguage(tmpDir)

	if err == nil {
		t.Error("Expected error when no project files found, got nil")
	}
}

func TestDetectLanguagePriority(t *testing.T) {
	// Create directory with both go.mod and pyproject.toml
	// Should prefer Go (first check)
	tmpDir := t.TempDir()

	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	pyprojectPath := filepath.Join(tmpDir, "pyproject.toml")
	err = os.WriteFile(pyprojectPath, []byte("[tool.poetry]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create pyproject.toml: %v", err)
	}

	runner := &Runner{}
	lang, err := runner.DetectLanguage(tmpDir)

	if err != nil {
		t.Errorf("Expected to detect language, got error: %v", err)
	}

	if lang != Go {
		t.Errorf("Expected Go (priority), got %v", lang)
	}
}

func TestGetTestCommand(t *testing.T) {
	tests := []struct {
		language Language
		expected string
	}{
		{Python, "pytest"},
		{Go, "go test"},
		{Java, "mvn test"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			runner := &Runner{language: tt.language}
			cmd := runner.getTestCommand()
			if cmd != tt.expected {
				t.Errorf("Expected command %q, got %q", tt.expected, cmd)
			}
		})
	}
}
