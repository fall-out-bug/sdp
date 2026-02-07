package prd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDetector(t *testing.T) {
	projectPath := t.TempDir()
	detector := NewDetector(projectPath)

	if detector == nil {
		t.Fatal("NewDetector returned nil")
	}

	if detector.projectPath != projectPath {
		t.Errorf("projectPath = %s, want %s", detector.projectPath, projectPath)
	}
}

func TestHasGoMod(t *testing.T) {
	projectPath := t.TempDir()
	detector := NewDetector(projectPath)

	// No go.mod
	if detector.hasGoMod() {
		t.Error("hasGoMod should return false when go.mod doesn't exist")
	}

	// Create go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	os.WriteFile(goModPath, []byte("module test\n"), 0644)

	if !detector.hasGoMod() {
		t.Error("hasGoMod should return true when go.mod exists")
	}
}

func TestHasDockerCompose(t *testing.T) {
	projectPath := t.TempDir()
	detector := NewDetector(projectPath)

	// No docker-compose
	if detector.hasDockerCompose() {
		t.Error("hasDockerCompose should return false when no compose files exist")
	}

	// Create docker-compose.yml
	composePath := filepath.Join(projectPath, "docker-compose.yml")
	os.WriteFile(composePath, []byte("version: '3'\n"), 0644)

	if !detector.hasDockerCompose() {
		t.Error("hasDockerCompose should return true when docker-compose.yml exists")
	}
}

func TestHasPythonProject(t *testing.T) {
	projectPath := t.TempDir()
	detector := NewDetector(projectPath)

	// No Python indicators
	if detector.hasPythonProject() {
		t.Error("hasPythonProject should return false when no Python files exist")
	}

	// Create pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	os.WriteFile(pyprojectPath, []byte("[project]\n"), 0644)

	if !detector.hasPythonProject() {
		t.Error("hasPythonProject should return true when pyproject.toml exists")
	}
}

func TestHasJavaProject(t *testing.T) {
	projectPath := t.TempDir()
	detector := NewDetector(projectPath)

	// No Java indicators
	if detector.hasJavaProject() {
		t.Error("hasJavaProject should return false when no Java files exist")
	}

	// Create pom.xml
	pomPath := filepath.Join(projectPath, "pom.xml")
	os.WriteFile(pomPath, []byte("<project></project>\n"), 0644)

	if !detector.hasJavaProject() {
		t.Error("hasJavaProject should return true when pom.xml exists")
	}
}

func TestHasMainPackage(t *testing.T) {
	projectPath := t.TempDir()
	detector := NewDetector(projectPath)

	// No main.go
	if detector.hasMainPackage() {
		t.Error("hasMainPackage should return false when no main.go exists")
	}

	// Create cmd/test/main.go structure
	cmdDir := filepath.Join(projectPath, "cmd", "test")
	os.MkdirAll(cmdDir, 0755)
	mainPath := filepath.Join(cmdDir, "main.go")
	os.WriteFile(mainPath, []byte("package main\n"), 0644)

	if !detector.hasMainPackage() {
		t.Error("hasMainPackage should return true when cmd/*/main.go exists")
	}
}

func TestHasCLIEntryPoints(t *testing.T) {
	projectPath := t.TempDir()
	detector := NewDetector(projectPath)

	// No entry points
	if detector.hasCLIEntryPoints() {
		t.Error("hasCLIEntryPoints should return false when no entry points exist")
	}

	// Create pyproject.toml with [project.scripts]
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	os.WriteFile(pyprojectPath, []byte("[project.scripts]\ncli = 'cli:main'\n"), 0644)

	if !detector.hasCLIEntryPoints() {
		t.Error("hasCLIEntryPoints should return true when [project.scripts] exists")
	}
}

func TestDetectTypeGoCLI(t *testing.T) {
	projectPath := t.TempDir()

	// Create go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	os.WriteFile(goModPath, []byte("module test\n"), 0644)

	// Create cmd/test/main.go
	cmdDir := filepath.Join(projectPath, "cmd", "test")
	os.MkdirAll(cmdDir, 0755)
	mainPath := filepath.Join(cmdDir, "main.go")
	os.WriteFile(mainPath, []byte("package main\n"), 0644)

	detector := NewDetector(projectPath)
	projectType := detector.DetectType()

	if projectType != CLI {
		t.Errorf("DetectType = %s, want %s", projectType, CLI)
	}
}

func TestDetectTypeGoLibrary(t *testing.T) {
	projectPath := t.TempDir()

	// Create go.mod only
	goModPath := filepath.Join(projectPath, "go.mod")
	os.WriteFile(goModPath, []byte("module test\n"), 0644)

	detector := NewDetector(projectPath)
	projectType := detector.DetectType()

	if projectType != Go {
		t.Errorf("DetectType = %s, want %s", projectType, Go)
	}
}

func TestDetectTypePythonService(t *testing.T) {
	projectPath := t.TempDir()

	// Create pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	os.WriteFile(pyprojectPath, []byte("[project]\n"), 0644)

	// Create docker-compose.yml
	composePath := filepath.Join(projectPath, "docker-compose.yml")
	os.WriteFile(composePath, []byte("version: '3'\n"), 0644)

	detector := NewDetector(projectPath)
	projectType := detector.DetectType()

	if projectType != Service {
		t.Errorf("DetectType = %s, want %s", projectType, Service)
	}
}

func TestDetectTypeJavaProject(t *testing.T) {
	projectPath := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(projectPath, "pom.xml")
	os.WriteFile(pomPath, []byte("<project></project>\n"), 0644)

	detector := NewDetector(projectPath)
	projectType := detector.DetectType()

	if projectType != Java {
		t.Errorf("DetectType = %s, want %s", projectType, Java)
	}
}

func TestDetectTypeUnknown(t *testing.T) {
	projectPath := t.TempDir()
	// Empty directory with no indicators

	detector := NewDetector(projectPath)
	projectType := detector.DetectType()

	if projectType != Library {
		t.Errorf("DetectType = %s, want %s (default)", projectType, Library)
	}
}

func TestProjectTypeString(t *testing.T) {
	tests := []struct {
		name     string
		p        ProjectType
		expected string
	}{
		{"Service", Service, "service"},
		{"CLI", CLI, "cli"},
		{"Library", Library, "library"},
		{"Go", Go, "go"},
		{"Python", Python, "python"},
		{"Java", Java, "java"},
		{"Unknown", Unknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.expected {
				t.Errorf("String() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestProjectTypeValue(t *testing.T) {
	tests := []struct {
		name     string
		p        ProjectType
		expected string
	}{
		{"Service", Service, "service"},
		{"CLI", CLI, "cli"},
		{"Library", Library, "library"},
		{"Go", Go, "go"},
		{"Python", Python, "python"},
		{"Java", Java, "java"},
		{"Unknown", Unknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Value(); got != tt.expected {
				t.Errorf("Value() = %s, want %s", got, tt.expected)
			}
		})
	}
}
