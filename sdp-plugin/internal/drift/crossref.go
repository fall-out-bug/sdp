package drift

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CrossRefValidator validates cross-references between docs (AC3)
type CrossRefValidator struct {
	projectRoot string
}

// NewCrossRefValidator creates a new cross-reference validator
func NewCrossRefValidator(projectRoot string) *CrossRefValidator {
	return &CrossRefValidator{projectRoot: projectRoot}
}

// Validate validates all cross-references in docs
func (v *CrossRefValidator) Validate() ([]DriftTypeReport, error) {
	docsDir := filepath.Join(v.projectRoot, "docs")

	if _, err := os.Stat(docsDir); os.IsNotExist(err) {
		return []DriftTypeReport{}, nil
	}

	var issues []EnhancedDriftIssue

	// Walk docs directory
	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Extract and check links
		links := v.extractLinks(string(content))
		for _, link := range links {
			if !v.checkLink(filepath.Dir(path), link) {
				issues = append(issues, EnhancedDriftIssue{
					File:     filepath.Base(path),
					Message:  "Broken link to " + link,
					Severity: SeverityWarning,
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(issues) > 0 {
		return []DriftTypeReport{{
			Type:     DriftTypeDocsDocs,
			Severity: SeverityWarning,
			Issues:   issues,
		}}, nil
	}

	return []DriftTypeReport{}, nil
}

// extractLinks extracts internal markdown links from content
func (v *CrossRefValidator) extractLinks(content string) []string {
	// Match [text](./path.md) or [text](../path.md)
	re := regexp.MustCompile(`\[([^\]]+)\]\((\.{1,2}/[^)]+\.md)\)`)
	matches := re.FindAllStringSubmatch(content, -1)

	var links []string
	for _, match := range matches {
		if len(match) > 2 {
			links = append(links, match[2])
		}
	}
	return links
}

// checkLink checks if a link target exists
func (v *CrossRefValidator) checkLink(baseDir, link string) bool {
	// Resolve relative path
	target := filepath.Join(baseDir, link)
	if _, err := os.Stat(target); err == nil {
		return true
	}
	return false
}
