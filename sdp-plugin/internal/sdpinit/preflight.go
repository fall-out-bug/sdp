package sdpinit

import (
	"os"
	"path/filepath"
)

// PreflightResult contains preflight check results
type PreflightResult struct {
	ProjectType string   // "go", "node", "python", "mixed", "unknown"
	HasSDP      bool     // .sdp directory exists
	HasClaude   bool     // .claude directory exists
	HasGit      bool     // .git directory exists
	Conflicts   []string // Existing files that would be overwritten
	Warnings    []string // Non-fatal issues
}

// RunPreflight runs all preflight checks
func RunPreflight() *PreflightResult {
	result := &PreflightResult{
		ProjectType: "unknown",
	}

	// Detect project type
	result.ProjectType = DetectProjectType()

	// Check for existing SDP structure
	result.HasSDP = dirExists(".sdp")
	result.HasClaude = dirExists(".claude")
	result.HasGit = dirExists(".git")

	// Check for conflicts
	result.Conflicts = checkConflicts()

	// Add warnings
	if !result.HasGit {
		result.Warnings = append(result.Warnings, "Not a git repository - version control recommended")
	}

	return result
}

// DetectProjectType determines the primary project type
func DetectProjectType() string {
	detectors := []struct {
		check func() bool
		typ   string
	}{
		{isGoProject, "go"},
		{isNodeProject, "node"},
		{isPythonProject, "python"},
	}

	detected := []string{}
	for _, d := range detectors {
		if d.check() {
			detected = append(detected, d.typ)
		}
	}

	if len(detected) == 0 {
		return "unknown"
	}
	if len(detected) == 1 {
		return detected[0]
	}
	return "mixed"
}

func isGoProject() bool {
	if _, err := os.Stat("go.mod"); err == nil {
		return true
	}
	if _, err := os.Stat("Gopkg.toml"); err == nil {
		return true
	}
	// Check for .go files
	return hasFilesWithExt(".go")
}

func isNodeProject() bool {
	if _, err := os.Stat("package.json"); err == nil {
		return true
	}
	if _, err := os.Stat("yarn.lock"); err == nil {
		return true
	}
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return true
	}
	return hasFilesWithExt(".ts", ".js", ".tsx", ".jsx")
}

func isPythonProject() bool {
	if _, err := os.Stat("setup.py"); err == nil {
		return true
	}
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return true
	}
	if _, err := os.Stat("requirements.txt"); err == nil {
		return true
	}
	return hasFilesWithExt(".py")
}

func hasFilesWithExt(exts ...string) bool {
	entries, err := os.ReadDir(".")
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		for _, ext := range exts {
			if filepath.Ext(entry.Name()) == ext {
				return true
			}
		}
	}
	return false
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func checkConflicts() []string {
	conflicts := []string{}

	// Check if .claude would overwrite anything
	if _, err := os.Stat(".claude/settings.json"); err == nil {
		conflicts = append(conflicts, ".claude/settings.json")
	}

	return conflicts
}
