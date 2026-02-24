package orchestrate

import (
	"regexp"
	"strings"
)

var (
	reScopeFile   = regexp.MustCompile(`^-\s+` + "`" + `([^` + "`" + `]+)` + "`")
	reAcceptance  = regexp.MustCompile(`^-\s+\[[ x]\]\s+(.+)`)
	reDependsOn   = regexp.MustCompile(`(?m)^depends_on:\s*\[(.*?)\]`)
)

func parseWorkstreamSections(content string) (acceptance []string, scopeFiles []string) {
	lines := strings.Split(content, "\n")
	var inScopeFiles, inAcceptance bool
	for _, line := range lines {
		if strings.TrimSpace(line) == "## Scope Files" {
			inScopeFiles = true
			inAcceptance = false
			continue
		}
		if strings.TrimSpace(line) == "## Acceptance Criteria" {
			inAcceptance = true
			inScopeFiles = false
			continue
		}
		if strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "## Scope") && !strings.HasPrefix(line, "## Acceptance") {
			inScopeFiles = false
			inAcceptance = false
			continue
		}
		if inAcceptance {
			if m := reAcceptance.FindStringSubmatch(line); len(m) > 1 {
				acceptance = append(acceptance, strings.TrimSpace(m[1]))
			}
		}
		if inScopeFiles {
			if m := reScopeFile.FindStringSubmatch(line); len(m) > 1 {
				scopeFiles = append(scopeFiles, strings.TrimSpace(m[1]))
			}
		}
	}
	return acceptance, scopeFiles
}

func parseDependsOn(content string) []string {
	var deps []string
	if m := reDependsOn.FindStringSubmatch(content); len(m) > 1 {
		for _, s := range strings.Split(m[1], ",") {
			id := strings.Trim(strings.Trim(s, `"`), " ")
			if id != "" {
				deps = append(deps, id)
			}
		}
	}
	return deps
}

func parseQualityGates(agentsContent string) string {
	idx := strings.Index(agentsContent, "## Quality Gates")
	if idx < 0 {
		return ""
	}
	rest := agentsContent[idx:]
	end := strings.Index(rest, "\n## ")
	if end > 0 {
		rest = rest[:end]
	}
	return strings.TrimSpace(rest)
}
