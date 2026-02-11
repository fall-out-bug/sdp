package parser

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ParseWorkstream parses a workstream markdown file with YAML frontmatter
func ParseWorkstream(wsPath string) (*Workstream, error) {
	// Read file
	content, err := os.ReadFile(wsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Check if file is empty
	if len(content) == 0 {
		return nil, fmt.Errorf("file is empty: %s", wsPath)
	}

	// Validate file size
	if err := ValidateFileSize(content); err != nil {
		return nil, fmt.Errorf("file size validation failed: %w", err)
	}

	// Extract frontmatter (delimited by ---)
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format: expected --- delimiters")
	}

	// Parse YAML frontmatter with security checks
	var fm frontmatter
	if err := SafeYAMLUnmarshal(parts[1], &fm); err != nil {
		return nil, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	// Validate required fields
	if fm.WSID == "" {
		return nil, fmt.Errorf("missing required field: ws_id")
	}
	feature := fm.getFeature()
	if feature == "" {
		return nil, fmt.Errorf("missing required field: feature or feature_id")
	}

	// Extract main content
	mainContent := string(parts[2])

	// Validate content length
	if err := ValidateContentLength(mainContent); err != nil {
		return nil, fmt.Errorf("content validation failed: %w", err)
	}

	// Parse goal section
	goal := extractSection(mainContent, "Goal")

	// Parse acceptance criteria
	acceptance := extractListItems(mainContent, "Acceptance Criteria")

	// Parse scope files
	scopeFiles := extractScopeFiles(mainContent)

	ws := &Workstream{
		ID:         fm.WSID,
		Parent:     fm.Parent,
		Feature:    fm.getFeature(),
		Status:     fm.Status,
		Size:       fm.Size,
		ProjectID:  fm.ProjectID,
		Goal:       goal,
		Acceptance: acceptance,
		Scope:      scopeFiles,
	}

	return ws, nil
}

// extractSection extracts a section by heading name
func extractSection(content, sectionName string) string {
	// Look for ## SectionName or ### SectionName
	headingPattern := regexp.MustCompile(`##+ ` + regexp.QuoteMeta(sectionName))

	// Find the heading
	headingIndex := headingPattern.FindStringIndex(content)
	if headingIndex == nil {
		return ""
	}

	// Start after the heading
	start := headingIndex[1]

	// Find the next heading or end of content
	nextHeading := regexp.MustCompile(`\n##+ `).FindStringIndex(content[start:])
	if nextHeading == nil {
		// No next heading, return rest of content
		section := content[start:]
		// Skip first line (the heading itself)
		if idx := strings.Index(section, "\n"); idx != -1 {
			section = section[idx+1:]
		}
		return strings.TrimSpace(section)
	}

	// Extract content up to next heading
	end := start + nextHeading[0]
	section := content[start:end]

	// Skip first line (the heading itself)
	if idx := strings.Index(section, "\n"); idx != -1 {
		section = section[idx+1:]
	}

	return strings.TrimSpace(section)
}

// extractListItems extracts list items from a section
func extractListItems(content, sectionName string) []string {
	section := extractSection(content, sectionName)
	if section == "" {
		return []string{}
	}

	// Extract markdown list items
	listPattern := regexp.MustCompile(`^[\s]*-\s*\[[ x ]\]\s*(.+)$`)
	lines := strings.Split(section, "\n")

	var items []string
	for _, line := range lines {
		matches := listPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			items = append(items, strings.TrimSpace(matches[1]))
		}
	}

	return items
}

// extractScopeFiles extracts scope files from the workstream
func extractScopeFiles(content string) Scope {
	section := extractSection(content, "Scope Files")
	if section == "" {
		return Scope{Implementation: []string{}, Tests: []string{}}
	}

	lines := strings.Split(section, "\n")
	var implementation []string
	var tests []string

	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Detect section headers
		if strings.HasPrefix(line, "**Implementation:**") || strings.HasPrefix(line, "Implementation:") {
			currentSection = "implementation"
			continue
		}
		if strings.HasPrefix(line, "**Tests:**") || strings.HasPrefix(line, "Tests:") {
			currentSection = "tests"
			continue
		}

		// Extract file paths (marked with -)
		if strings.HasPrefix(line, "-") {
			filePath := strings.TrimPrefix(line, "-")
			filePath = strings.TrimSpace(filePath)
			filePath = strings.TrimPrefix(filePath, "`")
			filePath = strings.TrimSuffix(filePath, "`")

			// Skip (NEW), (ENHANCE) markers
			if strings.Contains(filePath, "(NEW)") || strings.Contains(filePath, "(ENHANCE)") {
				filePath = regexp.MustCompile(`\s*\(.*?\)`).ReplaceAllString(filePath, "")
			}
			filePath = strings.TrimSpace(filePath)

			if currentSection == "implementation" {
				implementation = append(implementation, filePath)
			} else if currentSection == "tests" {
				tests = append(tests, filePath)
			}
		}
	}

	return Scope{
		Implementation: implementation,
		Tests:          tests,
	}
}
