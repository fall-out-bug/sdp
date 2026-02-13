package drift

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ADRValidator validates ADR decisions against code (AC2)
type ADRValidator struct {
	projectRoot string
}

// NewADRValidator creates a new ADR validator
func NewADRValidator(projectRoot string) *ADRValidator {
	return &ADRValidator{projectRoot: projectRoot}
}

// Validate validates all ADRs in the decisions directory
func (v *ADRValidator) Validate() ([]DriftTypeReport, error) {
	decisionsDir := filepath.Join(v.projectRoot, "docs", "decisions")

	// Check if decisions directory exists
	if _, err := os.Stat(decisionsDir); os.IsNotExist(err) {
		return []DriftTypeReport{}, nil
	}

	var reports []DriftTypeReport
	var issues []EnhancedDriftIssue

	// Walk decisions directory
	err := filepath.Walk(decisionsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Parse frontmatter
		status := extractStatus(string(content))
		decisionID := extractDecisionID(string(content))

		// Check for superseded decisions
		if status == "superseded" {
			issues = append(issues, EnhancedDriftIssue{
				File:     filepath.Base(path),
				Message:  "Decision " + decisionID + " is superseded",
				Severity: SeverityWarning,
			})
		}

		// Check for deprecated decisions
		if status == "deprecated" {
			issues = append(issues, EnhancedDriftIssue{
				File:     filepath.Base(path),
				Message:  "Decision " + decisionID + " is deprecated",
				Severity: SeverityWarning,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(issues) > 0 {
		reports = append(reports, DriftTypeReport{
			Type:     DriftTypeDecisionCode,
			Severity: SeverityWarning,
			Issues:   issues,
		})
	}

	return reports, nil
}

// extractKeywords extracts key technical terms from content
func (v *ADRValidator) extractKeywords(content string) []string {
	keywords := []string{}
	content = strings.ToLower(content)

	// Extract technology keywords
	patterns := []string{
		`sqlite`, `postgres`, `mysql`, `mongodb`,
		`redis`, `fts5`, `elasticsearch`,
		`rest`, `graphql`, `grpc`,
		`docker`, `kubernetes`, `aws`, `gcp`,
	}

	for _, p := range patterns {
		if matched, _ := regexp.MatchString(p, content); matched { //nolint:errcheck // pattern is hardcoded
			keywords = append(keywords, p)
		}
	}

	return keywords
}

// extractStatus extracts status from frontmatter
func extractStatus(content string) string {
	re := regexp.MustCompile(`(?m)^status:\s*(\w+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// extractDecisionID extracts decision_id from frontmatter
func extractDecisionID(content string) string {
	re := regexp.MustCompile(`(?m)^decision_id:\s*(\S+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
