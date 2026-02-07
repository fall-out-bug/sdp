package drift

import (
	"os"
	"path/filepath"
	"strings"
)

// Analyzer analyzes source code for purpose and structure
type Analyzer struct {
	projectRoot string
}

// NewAnalyzer creates a new code analyzer
func NewAnalyzer(projectRoot string) *Analyzer {
	return &Analyzer{
		projectRoot: projectRoot,
	}
}

// AnalyzeFile analyzes a file and returns its purpose
func (a *Analyzer) AnalyzeFile(filePath string) (string, error) {
	// Make path absolute
	fullPath := filePath
	if !filepath.IsAbs(filePath) {
		fullPath = filepath.Join(a.projectRoot, filePath)
	}

	// Read file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	contentStr := string(content)

	// Extract purpose from comments/docstrings
	purpose := extractPurpose(contentStr)

	return purpose, nil
}

// extractPurpose extracts purpose from file comments
func extractPurpose(content string) string {
	lines := strings.Split(content, "\n")
	purposeLines := []string{}

	inDocstring := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Go comment
		if strings.HasPrefix(trimmed, "//") {
			comment := strings.TrimPrefix(trimmed, "//")
			comment = strings.TrimSpace(comment)
			if comment != "" {
				purposeLines = append(purposeLines, comment)
			}
			if len(purposeLines) >= 3 {
				break
			}
		}

		// Python docstring
		if strings.HasPrefix(trimmed, `"""`) || strings.HasPrefix(trimmed, `'''`) {
			if !inDocstring {
				inDocstring = true
				continue
			} else {
				break
			}
		}

		if inDocstring {
			purposeLines = append(purposeLines, trimmed)
		}
	}

	if len(purposeLines) == 0 {
		return "No purpose documentation found"
	}

	return strings.Join(purposeLines, " ")
}

// ComparePurpose compares documented purpose with actual purpose
func (a *Analyzer) ComparePurpose(filePath string, documentedPurpose string) (bool, string) {
	actualPurpose, err := a.AnalyzeFile(filePath)
	if err != nil {
		return false, "Error reading file"
	}

	// Simple comparison: check if documented purpose contains key terms
	if documentedPurpose == "" {
		return true, actualPurpose
	}

	// Extract key terms from documented purpose
	terms := extractKeyTerms(documentedPurpose)

	// Check if any terms appear in actual purpose
	for _, term := range terms {
		if strings.Contains(strings.ToLower(actualPurpose), strings.ToLower(term)) {
			return true, actualPurpose
		}
	}

	// No match found
	return false, actualPurpose
}

// extractKeyTerms extracts meaningful terms from a purpose string
func extractKeyTerms(purpose string) []string {
	// Remove common words and extract key terms
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true,
		"or": true, "but": true, "in": true, "on": true,
		"at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true,
	}

	words := strings.Fields(purpose)
	terms := []string{}

	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if len(word) > 3 && !stopWords[word] {
			terms = append(terms, word)
		}
	}

	return terms
}
