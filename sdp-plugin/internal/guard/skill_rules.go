package guard

import (
	"fmt"
	"os"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
)

// applyGuardRules applies loaded guard rules to staged files (AC1)
func (s *Skill) applyGuardRules(files []string, rules *config.GuardRules) []Finding {
	var findings []Finding

	for _, file := range files {
		absPath, err := ResolvePath(file)
		if err != nil {
			continue
		}

		// Read file content
		content, err := os.ReadFile(absPath)
		if err != nil {
			continue
		}

		// Check each enabled rule
		for _, rule := range rules.Rules {
			if !rule.Enabled {
				continue
			}

			switch rule.ID {
			case "max-file-loc":
				findings = append(findings, s.checkMaxFileLOC(file, content, rule)...)
			case "coverage-threshold":
				// Coverage is checked at project level, not per file
				continue
			case "max-cyclomatic-complexity":
				// Complexity requires parsing - skip for now
				continue
			case "no-commented-code":
				findings = append(findings, s.checkCommentedCode(file, content, rule)...)
			case "no-orphaned-todos":
				findings = append(findings, s.checkOrphanedTODOs(file, content, rule)...)
			}
		}
	}

	return findings
}

// checkMaxFileLOC checks if file exceeds max lines of code (AC1, AC6)
func (s *Skill) checkMaxFileLOC(file string, content []byte, rule config.GuardRule) []Finding {
	var findings []Finding

	maxLines := 200 // default
	if maxLinesVal, ok := rule.Config["max_lines"]; ok {
		switch v := maxLinesVal.(type) {
		case int:
			maxLines = v
		case float64:
			maxLines = int(v)
		}
	}

	lines := strings.Split(string(content), "\n")
	loc := len(lines)

	if loc > maxLines {
		severity := SeverityWarning
		if rule.Severity == "error" {
			severity = SeverityError
		}

		findings = append(findings, Finding{
			Severity: severity,
			Rule:     rule.ID,
			File:     file,
			Message:  fmt.Sprintf("File exceeds maximum size: %d LOC (threshold: %d)", loc, maxLines),
		})
	}

	return findings
}

// checkCommentedCode checks for commented-out code (AC1)
func (s *Skill) checkCommentedCode(file string, content []byte, rule config.GuardRule) []Finding {
	var findings []Finding
	lines := strings.Split(string(content), "\n")

	commentPrefix := "//"
	if strings.HasSuffix(file, ".py") {
		commentPrefix = "#"
	}

	consecutiveComments := 0
	maxComments := 3 // Threshold for detecting commented code blocks

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, commentPrefix) && len(trimmed) > 10 {
			// Skip actual comments that look like documentation
			if strings.Contains(strings.ToLower(trimmed), "todo") ||
				strings.Contains(strings.ToLower(trimmed), "fixme") ||
				strings.Contains(strings.ToLower(trimmed), "note") {
				consecutiveComments = 0
				continue
			}

			consecutiveComments++
			if consecutiveComments >= maxComments {
				severity := SeverityWarning
				if rule.Severity == "error" {
					severity = SeverityError
				}

				findings = append(findings, Finding{
					Severity: severity,
					Rule:     rule.ID,
					File:     file,
					Line:     i + 1,
					Message:  "Possible commented-out code detected",
				})
				break
			}
		} else {
			consecutiveComments = 0
		}
	}

	return findings
}

// checkOrphanedTODOs checks for TODOs without workstream ID (AC1)
func (s *Skill) checkOrphanedTODOs(file string, content []byte, rule config.GuardRule) []Finding {
	var findings []Finding
	lines := strings.Split(string(content), "\n")

	todoPattern := "TODO"

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, todoPattern) {
			// Check if line contains workstream ID pattern
			hasWSID := false
			for _, wsPattern := range []string{"WS-", "ws-"} {
				if strings.Contains(trimmed, wsPattern) {
					// Verify format WS-XXX-YY
					idx := strings.Index(trimmed, wsPattern)
					if idx+len(wsPattern)+6 <= len(trimmed) {
						potentialID := trimmed[idx : idx+len(wsPattern)+6]
						// Check if it matches pattern like WS-063-03
						if len(potentialID) >= 7 && strings.Count(potentialID, "-") >= 2 {
							hasWSID = true
							break
						}
					}
				}
			}

			if !hasWSID {
				severity := SeverityWarning
				if rule.Severity == "error" {
					severity = SeverityError
				}

				findings = append(findings, Finding{
					Severity: severity,
					Rule:     rule.ID,
					File:     file,
					Line:     i + 1,
					Message:  "TODO without workstream ID (format: TODO(WS-XXX-YY))",
				})
			}
		}
	}

	return findings
}
