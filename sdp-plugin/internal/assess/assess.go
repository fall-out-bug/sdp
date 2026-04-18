package assess

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Recommendation represents a single recommendation for the project
type Recommendation struct {
	Category string
	Title    string
	Message  string
	Priority string
}

// Assessment represents the complete assessment of a project
type Assessment struct {
	Language        string
	Framework       []string
	Structure       []string
	IsMonorepo      bool
	HasTests        bool
	HasCI           bool
	Recommendations []Recommendation
}

// Assess performs a read-only scan of the repository
func Assess(projectPath string) (*Assessment, error) {
	result := &Assessment{
		Framework:       []string{},
		Structure:       []string{},
		Recommendations: []Recommendation{},
	}

	// Detect language
	lang, err := detectLanguage(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect language: %w", err)
	}
	result.Language = lang

	// Detect framework
	frameworks := detectFramework(projectPath, lang)
	result.Framework = frameworks

	// Detect structure
	structure := detectStructure(projectPath)
	result.Structure = structure

	// Detect if monorepo
	result.IsMonorepo = detectMonorepo(projectPath)

	// Detect tests
	result.HasTests = detectTests(projectPath)

	// Detect CI
	result.HasCI = detectCI(projectPath)

	// Generate recommendations
	generateRecommendations(result)

	return result, nil
}

// detectLanguage identifies the primary programming language
func detectLanguage(projectPath string) (string, error) {
	detectors := []struct {
		files    []string
		language string
	}{
		{[]string{"go.mod"}, "Go"},
		{[]string{"package.json"}, "Node.js/TypeScript"},
		{[]string{"requirements.txt", "pyproject.toml", "setup.py", "Pipfile"}, "Python"},
		{[]string{"Cargo.toml"}, "Rust"},
		{[]string{"pom.xml", "build.gradle"}, "Java"},
		{[]string{"Gemfile"}, "Ruby"},
		{[]string{"composer.json"}, "PHP"},
		{[]string{"*.csproj", "*.sln"}, "C#"},
	}

	for _, detector := range detectors {
		for _, file := range detector.files {
			if strings.Contains(file, "*") {
				matches, err := filepath.Glob(filepath.Join(projectPath, file))
				if err != nil {
					// If glob fails, skip this pattern
					continue
				}
				if len(matches) > 0 {
					return detector.language, nil
				}
			} else {
				if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
					return detector.language, nil
				}
			}
		}
	}

	return "Unknown", nil
}

// detectFramework identifies frameworks based on dependencies
func detectFramework(projectPath, language string) []string {
	frameworks := []string{}

	switch language {
	case "Go":
		frameworks = detectGoFrameworks(projectPath)
	case "Node.js/TypeScript":
		frameworks = detectNodeFrameworks(projectPath)
	case "Python":
		frameworks = detectPythonFrameworks(projectPath)
	}

	if len(frameworks) == 0 {
		frameworks = append(frameworks, "None detected")
	}

	return frameworks
}

// detectGoFrameworks identifies Go frameworks
func detectGoFrameworks(projectPath string) []string {
	frameworks := []string{}
	goModPath := filepath.Join(projectPath, "go.mod")
	if content, err := os.ReadFile(goModPath); err == nil {
		contentStr := string(content)
		if strings.Contains(contentStr, "github.com/gin-gonic/gin") {
			frameworks = append(frameworks, "Gin")
		}
		if strings.Contains(contentStr, "github.com/gorilla/mux") {
			frameworks = append(frameworks, "Gorilla Mux")
		}
		if strings.Contains(contentStr, "net/http") {
			frameworks = append(frameworks, "net/http (stdlib)")
		}
	}
	return frameworks
}

// detectNodeFrameworks identifies Node.js/TypeScript frameworks
func detectNodeFrameworks(projectPath string) []string {
	frameworks := []string{}
	packageJsonPath := filepath.Join(projectPath, "package.json")
	if content, err := os.ReadFile(packageJsonPath); err == nil {
		contentStr := string(content)
		if strings.Contains(contentStr, "\"react\"") {
			frameworks = append(frameworks, "React")
		}
		if strings.Contains(contentStr, "\"vue\"") {
			frameworks = append(frameworks, "Vue")
		}
		if strings.Contains(contentStr, "\"next\"") {
			frameworks = append(frameworks, "Next.js")
		}
		if strings.Contains(contentStr, "\"express\"") {
			frameworks = append(frameworks, "Express")
		}
		if strings.Contains(contentStr, "\"@angular\"") {
			frameworks = append(frameworks, "Angular")
		}
	}
	return frameworks
}

// detectPythonFrameworks identifies Python frameworks
func detectPythonFrameworks(projectPath string) []string {
	frameworks := []string{}
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")

	var content []byte
	var err error

	if content, err = os.ReadFile(requirementsPath); err != nil {
		content, err = os.ReadFile(pyprojectPath)
		if err != nil {
			// Neither file could be read
			return frameworks
		}
	}

	if len(content) > 0 {
		contentStr := string(content)
		if strings.Contains(contentStr, "django") {
			frameworks = append(frameworks, "Django")
		}
		if strings.Contains(contentStr, "flask") {
			frameworks = append(frameworks, "Flask")
		}
		if strings.Contains(contentStr, "fastapi") {
			frameworks = append(frameworks, "FastAPI")
		}
	}
	return frameworks
}

// detectStructure identifies project structure patterns
func detectStructure(projectPath string) []string {
	structures := []string{}

	dirs := []string{
		"src", "cmd", "internal", "pkg", "lib", "app",
		"components", "pages", "services", "handlers", "models", "utils",
		"tests", "test", "__tests__", "spec",
	}

	for _, dir := range dirs {
		if info, err := os.Stat(filepath.Join(projectPath, dir)); err == nil && info.IsDir() {
			structures = append(structures, dir)
		}
	}

	return structures
}

// detectMonorepo checks if this is a monorepo
func detectMonorepo(projectPath string) bool {
	indicators := []string{
		"packages", "apps", "services", "workspaces",
		".gitmodules", "pnpm-workspace.yaml", "lerna.json",
	}

	for _, indicator := range indicators {
		if info, err := os.Stat(filepath.Join(projectPath, indicator)); err == nil {
			if info.IsDir() || strings.HasPrefix(indicator, ".") {
				return true
			}
		}
	}

	return false
}

// detectTests checks if the project has tests.
// Walks the directory tree recursively since filepath.Glob does not support "**".
func detectTests(projectPath string) bool {
	testDirs := []string{
		"tests", "test", "__tests__", "spec",
	}

	for _, dir := range testDirs {
		if info, err := os.Stat(filepath.Join(projectPath, dir)); err == nil && info.IsDir() {
			return true
		}
	}

	suffixes := []string{"_test.go", "_test.py", ".test.ts", ".test.js", ".spec.ts", ".spec.js"}
	found := false
	err := filepath.WalkDir(projectPath, func(path string, d os.DirEntry, err error) error {
		if err != nil || found {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		for _, suf := range suffixes {
			if strings.HasSuffix(name, suf) {
				found = true
				return filepath.SkipAll
			}
		}
		return nil
	})
	if err != nil {
		// If walk fails, return false (no tests found)
		return false
	}
	return found
}

// detectCI checks if the project has CI configuration
func detectCI(projectPath string) bool {
	ciIndicators := []string{
		".github", ".gitlab-ci.yml", ".circleci",
		".travis.yml", "jenkins.yml", "azure-pipelines.yml",
	}

	for _, ci := range ciIndicators {
		if _, err := os.Stat(filepath.Join(projectPath, ci)); err == nil {
			return true
		}
	}

	return false
}

// generateRecommendations creates recommendations based on assessment
func generateRecommendations(result *Assessment) {
	// SDP setup recommendation
	result.Recommendations = append(result.Recommendations, Recommendation{
		Category: "sdp",
		Title:    "Initialize SDP",
		Message:  "Run 'sdp init' to initialize Spec-Driven Protocol for this project",
		Priority: "medium",
	})

	// Test recommendations
	if !result.HasTests {
		result.Recommendations = append(result.Recommendations, Recommendation{
			Category: "testing",
			Title:    "Add Tests",
			Message:  "No tests detected. Consider adding tests for reliability.",
			Priority: "high",
		})
	}

	// CI recommendations
	if !result.HasCI {
		result.Recommendations = append(result.Recommendations, Recommendation{
			Category: "ci",
			Title:    "Setup CI",
			Message:  "No CI detected. Consider setting up GitHub Actions or similar.",
			Priority: "medium",
		})
	}

	// Language-specific recommendations
	switch result.Language {
	case "Go":
		result.Recommendations = append(result.Recommendations, Recommendation{
			Category: "tooling",
			Title:    "Use go.mod",
			Message:  "Ensure dependencies are managed via go.mod",
			Priority: "low",
		})
	case "Node.js/TypeScript":
		result.Recommendations = append(result.Recommendations, Recommendation{
			Category: "tooling",
			Title:    "Use package.json",
			Message:  "Ensure dependencies are managed via package.json",
			Priority: "low",
		})
	case "Python":
		result.Recommendations = append(result.Recommendations, Recommendation{
			Category: "tooling",
			Title:    "Use Virtual Environment",
			Message:  "Consider using venv or pyenv for dependency isolation",
			Priority: "medium",
		})
	}

	// Monorepo recommendations
	if result.IsMonorepo {
		result.Recommendations = append(result.Recommendations, Recommendation{
			Category: "structure",
			Title:    "Monorepo Detected",
			Message:  "This appears to be a monorepo. SDP can work with monorepos - consider using workspaces.",
			Priority: "low",
		})
	}
}
