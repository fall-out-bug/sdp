package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checkWorkstreamCircularDeps checks for circular dependencies in workstreams
func checkWorkstreamCircularDeps() DeepCheckResult {
	start := getTime()
	details := make(map[string]interface{})

	wsDir := "docs/workstreams/backlog"
	entries, err := os.ReadDir(wsDir)
	if os.IsNotExist(err) {
		return DeepCheckResult{
			Check:    "Workstream Deps",
			Status:   "ok",
			Duration: since(start),
			Message:  "No workstreams directory",
			Details:  details,
		}
	}
	if err != nil {
		return DeepCheckResult{
			Check:    "Workstream Deps",
			Status:   "warning",
			Duration: since(start),
			Message:  fmt.Sprintf("Failed to read workstreams: %v", err),
			Details:  details,
		}
	}

	// Build dependency graph
	deps := make(map[string][]string)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		wsID := strings.TrimSuffix(entry.Name(), ".md")
		content, err := os.ReadFile(filepath.Join(wsDir, entry.Name()))
		if err != nil {
			continue
		}

		// Simple "Depends on:" parsing
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "**Depends on**:") {
				depList := strings.TrimPrefix(line, "**Depends on**:")
				depList = strings.TrimSpace(depList)
				if depList != "" && depList != "-" {
					for _, dep := range strings.Split(depList, ",") {
						deps[wsID] = append(deps[wsID], strings.TrimSpace(dep))
					}
				}
			}
		}
	}

	// Check for cycles using DFS
	cycles := detectCycles(deps)
	details["workstreams"] = len(deps)
	details["cycles"] = cycles

	if len(cycles) > 0 {
		return DeepCheckResult{
			Check:    "Workstream Deps",
			Status:   "error",
			Duration: since(start),
			Message:  fmt.Sprintf("Circular dependencies: %v", cycles),
			Details:  details,
		}
	}

	return DeepCheckResult{
		Check:    "Workstream Deps",
		Status:   "ok",
		Duration: since(start),
		Message:  fmt.Sprintf("No cycles in %d workstreams", len(deps)),
		Details:  details,
	}
}

// detectCycles finds circular dependencies using DFS
func detectCycles(deps map[string][]string) []string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles []string

	var dfs func(node string, path []string) bool
	dfs = func(node string, path []string) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range deps[node] {
			if !visited[neighbor] {
				if dfs(neighbor, append(path, neighbor)) {
					return true
				}
			} else if recStack[neighbor] {
				cycles = append(cycles, fmt.Sprintf("%s -> %s", strings.Join(path, " -> "), neighbor))
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for node := range deps {
		if !visited[node] {
			dfs(node, []string{node})
		}
	}

	return cycles
}

// checkSkillsSyntax validates skill files have valid structure
func checkSkillsSyntax() DeepCheckResult {
	start := getTime()
	details := make(map[string]interface{})

	skillsDir := ".claude/skills"
	entries, err := os.ReadDir(skillsDir)
	if os.IsNotExist(err) {
		return DeepCheckResult{
			Check:    "Skills Syntax",
			Status:   "warning",
			Duration: since(start),
			Message:  ".claude/skills directory not found",
			Details:  details,
		}
	}
	if err != nil {
		return DeepCheckResult{
			Check:    "Skills Syntax",
			Status:   "error",
			Duration: since(start),
			Message:  fmt.Sprintf("Failed to read skills: %v", err),
			Details:  details,
		}
	}

	invalid := []string{}
	checked := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillFile := filepath.Join(skillsDir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			invalid = append(invalid, entry.Name()+" (missing SKILL.md)")
			continue
		}

		// Read and validate basic structure
		content, err := os.ReadFile(skillFile)
		if err != nil {
			invalid = append(invalid, entry.Name()+" (read error)")
			continue
		}

		// Check for YAML frontmatter
		if !strings.HasPrefix(string(content), "---") {
			invalid = append(invalid, entry.Name()+" (missing frontmatter)")
			continue
		}

		checked++
	}

	details["checked"] = checked
	details["invalid"] = invalid

	if len(invalid) > 0 {
		return DeepCheckResult{
			Check:    "Skills Syntax",
			Status:   "warning",
			Duration: since(start),
			Message:  fmt.Sprintf("Issues: %s", strings.Join(invalid, ", ")),
			Details:  details,
		}
	}

	return DeepCheckResult{
		Check:    "Skills Syntax",
		Status:   "ok",
		Duration: since(start),
		Message:  fmt.Sprintf("All %d skills valid", checked),
		Details:  details,
	}
}
