package reality

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// DetectLanguage detects the programming language and framework
func DetectLanguage(projectPath string) (language string, framework string) {
	// Check for Go
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		language = "go"
		framework = detectGoFramework(projectPath)
		return
	}

	// Check for Python
	if _, err := os.Stat(filepath.Join(projectPath, "pyproject.toml")); err == nil {
		language = "python"
		framework = detectPythonFramework(projectPath)
		return
	}
	if _, err := os.Stat(filepath.Join(projectPath, "requirements.txt")); err == nil {
		language = "python"
		framework = detectPythonFramework(projectPath)
		return
	}
	if _, err := os.Stat(filepath.Join(projectPath, "setup.py")); err == nil {
		language = "python"
		framework = detectPythonFramework(projectPath)
		return
	}

	// Check for Java
	if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		language = "java"
		framework = "spring"
		return
	}
	if _, err := os.Stat(filepath.Join(projectPath, "build.gradle")); err == nil {
		language = "java"
		framework = "gradle"
		return
	}

	// Check for Node.js
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
		language = "nodejs"
		framework = detectNodeFramework(projectPath)
		return
	}

	return "unknown", ""
}

// detectGoFramework detects Go web frameworks
func detectGoFramework(projectPath string) string {
	goModPath := filepath.Join(projectPath, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "gin-gonic/gin") {
			return "gin"
		}
		if strings.Contains(line, "gorilla/mux") {
			return "gorilla"
		}
		if strings.Contains(line, "fiber") {
			return "fiber"
		}
	}

	return ""
}

// detectPythonFramework detects Python frameworks
func detectPythonFramework(projectPath string) string {
	// Check pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if file, err := os.Open(pyprojectPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.ToLower(scanner.Text())
			if strings.Contains(line, "django") {
				return "django"
			}
			if strings.Contains(line, "flask") {
				return "flask"
			}
			if strings.Contains(line, "fastapi") {
				return "fastapi"
			}
		}
	}

	// Check requirements.txt
	reqPath := filepath.Join(projectPath, "requirements.txt")
	if file, err := os.Open(reqPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.ToLower(scanner.Text())
			if strings.HasPrefix(line, "django") {
				return "django"
			}
			if strings.HasPrefix(line, "flask") {
				return "flask"
			}
			if strings.HasPrefix(line, "fastapi") {
				return "fastapi"
			}
		}
	}

	return ""
}

// detectNodeFramework detects Node.js frameworks
func detectNodeFramework(projectPath string) string {
	packagePath := filepath.Join(projectPath, "package.json")
	file, err := os.Open(packagePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		if strings.Contains(line, "\"express\"") {
			return "express"
		}
		if strings.Contains(line, "\"react\"") {
			return "react"
		}
		if strings.Contains(line, "\"vue\"") {
			return "vue"
		}
		if strings.Contains(line, "\"nestjs\"") {
			return "nestjs"
		}
	}

	return ""
}
