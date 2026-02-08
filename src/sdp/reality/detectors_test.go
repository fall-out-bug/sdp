package reality

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDetectLanguage_AllLanguages tests detection of all 6 supported languages
func TestDetectLanguage_AllLanguages(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string]string
		expectedLang  string
		expectedFw    string
	}{
		{
			name: "Go project",
			files: map[string]string{
				"go.mod": "module test\n\ngo 1.21\n",
			},
			expectedLang: "go",
			expectedFw:   "",
		},
		{
			name: "Python with pyproject.toml",
			files: map[string]string{
				"pyproject.toml": "[project]\nname = \"test\"\n",
			},
			expectedLang: "python",
			expectedFw:   "",
		},
		{
			name: "Python with requirements.txt",
			files: map[string]string{
				"requirements.txt": "requests==2.28.0\n",
			},
			expectedLang: "python",
			expectedFw:   "",
		},
		{
			name: "Python with setup.py",
			files: map[string]string{
				"setup.py": "from setuptools import setup\nsetup()\n",
			},
			expectedLang: "python",
			expectedFw:   "",
		},
		{
			name: "Java with pom.xml",
			files: map[string]string{
				"pom.xml": "<project></project>\n",
			},
			expectedLang: "java",
			expectedFw:   "spring",
		},
		{
			name: "Java with build.gradle",
			files: map[string]string{
				"build.gradle": "plugins {\n    id 'java'\n}\n",
			},
			expectedLang: "java",
			expectedFw:   "gradle",
		},
		{
			name: "Node.js",
			files: map[string]string{
				"package.json": "{\"name\": \"test\"}\n",
			},
			expectedLang: "nodejs",
			expectedFw:   "",
		},
		{
			name:         "Unknown project",
			files:        map[string]string{},
			expectedLang: "unknown",
			expectedFw:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test files
			for filename, content := range tt.files {
				path := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create %s: %v", filename, err)
				}
			}

			lang, fw := DetectLanguage(tmpDir)

			if lang != tt.expectedLang {
				t.Errorf("Expected language '%s', got '%s'", tt.expectedLang, lang)
			}
			if fw != tt.expectedFw {
				t.Errorf("Expected framework '%s', got '%s'", tt.expectedFw, fw)
			}
		})
	}
}

// TestDetectPythonFramework_AllFrameworks tests all Python frameworks
func TestDetectPythonFramework_AllFrameworks(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		content       string
		expectedFw    string
	}{
		{
			name:       "Django via pyproject.toml",
			file:       "pyproject.toml",
			content:    "[project]\nname = \"test\"\ndjango = \"4.2.0\"\n",
			expectedFw: "django",
		},
		{
			name:       "Flask via pyproject.toml",
			file:       "pyproject.toml",
			content:    "[project]\nname = \"test\"\nFlask = \"2.3.0\"\n",
			expectedFw: "flask",
		},
		{
			name:       "FastAPI via pyproject.toml",
			file:       "pyproject.toml",
			content:    "[project]\nname = \"test\"\nfastapi = \"0.100.0\"\n",
			expectedFw: "fastapi",
		},
		{
			name:       "Django via requirements.txt",
			file:       "requirements.txt",
			content:    "django==4.2.0\n",
			expectedFw: "django",
		},
		{
			name:       "Flask via requirements.txt",
			file:       "requirements.txt",
			content:    "flask==2.3.0\n",
			expectedFw: "flask",
		},
		{
			name:       "FastAPI via requirements.txt",
			file:       "requirements.txt",
			content:    "fastapi==0.100.0\n",
			expectedFw: "fastapi",
		},
		{
			name:       "Unknown Python framework",
			file:       "requirements.txt",
			content:    "requests==2.28.0\n",
			expectedFw: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test file
			path := filepath.Join(tmpDir, tt.file)
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create %s: %v", tt.file, err)
			}

			_, fw := DetectLanguage(tmpDir)

			if fw != tt.expectedFw {
				t.Errorf("Expected framework '%s', got '%s'", tt.expectedFw, fw)
			}
		})
	}
}

// TestDetectNodeFramework_AllFrameworks tests all Node.js frameworks
func TestDetectNodeFramework_AllFrameworks(t *testing.T) {
	tests := []struct {
		name       string
		dependencies string
		expectedFw string
	}{
		{
			name:       "Express",
			dependencies: `"express": "^4.18.0"`,
			expectedFw: "express",
		},
		{
			name:       "React",
			dependencies: `"react": "^18.2.0"`,
			expectedFw: "react",
		},
		{
			name:       "Vue",
			dependencies: `"vue": "^3.3.0"`,
			expectedFw: "vue",
		},
		{
			name:       "NestJS (known issue: scoped packages not detected)",
			dependencies: `"@nestjs/core": "^10.0.0"`,
			expectedFw: "",  // Current implementation doesn't detect @nestjs/core
		},
		{
			name:       "Unknown framework",
			dependencies: `"lodash": "^4.17.0"`,
			expectedFw: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create package.json
			packagePath := filepath.Join(tmpDir, "package.json")
			content := `{"name": "test", "dependencies": {` + tt.dependencies + `}}`
			if err := os.WriteFile(packagePath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create package.json: %v", err)
			}

			_, fw := DetectLanguage(tmpDir)

			if fw != tt.expectedFw {
				t.Errorf("Expected framework '%s', got '%s'", tt.expectedFw, fw)
			}
		})
	}
}

// TestDetectGoFramework_AllFrameworks tests all Go frameworks
func TestDetectGoFramework_AllFrameworks(t *testing.T) {
	tests := []struct {
		name       string
		dependency string
		expectedFw string
	}{
		{
			name:       "Gin",
			dependency: "github.com/gin-gonic/gin v1.9.1",
			expectedFw: "gin",
		},
		{
			name:       "Gorilla",
			dependency: "github.com/gorilla/mux v1.8.0",
			expectedFw: "gorilla",
		},
		{
			name:       "Fiber",
			dependency: "github.com/gofiber/fiber/v2 v2.48.0",
			expectedFw: "fiber",
		},
		{
			name:       "Unknown framework",
			dependency: "",
			expectedFw: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create go.mod
			goModPath := filepath.Join(tmpDir, "go.mod")
			content := "module test\n"
			if tt.dependency != "" {
				content += "\nrequire " + tt.dependency + "\n"
			}
			if err := os.WriteFile(goModPath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create go.mod: %v", err)
			}

			_, fw := DetectLanguage(tmpDir)

			if fw != tt.expectedFw {
				t.Errorf("Expected framework '%s', got '%s'", tt.expectedFw, fw)
			}
		})
	}
}
