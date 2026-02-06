package drift

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// Detector detects documentation-code drift
type Detector struct {
	projectRoot string
}

// NewDetector creates a new drift detector
func NewDetector(projectRoot string) *Detector {
	return &Detector{
		projectRoot: projectRoot,
	}
}

// DetectDrift detects drift between documentation and actual code
func (d *Detector) DetectDrift(wsPath string) (*DriftReport, error) {
	// Parse workstream
	ws, err := parser.ParseWorkstream(wsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse workstream: %w", err)
	}

	report := &DriftReport{
		WorkstreamID: ws.ID,
		Timestamp:    time.Now(),
		Issues:       []DriftIssue{},
	}

	// Combine implementation and test files
	allFiles := append(ws.Scope.Implementation, ws.Scope.Tests...)

	// Check each file in scope
	for _, filePath := range allFiles {
		// Make path absolute
		fullPath := filePath
		if !filepath.IsAbs(filePath) {
			fullPath = filepath.Join(d.projectRoot, filePath)
		}

		// Check file existence
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			report.AddIssue(DriftIssue{
				File:          filePath,
				Status:        StatusError,
				Expected:      "File exists",
				Actual:        "File not found",
				Recommendation: fmt.Sprintf("Create file: %s", filePath),
			})
			continue
		}

		// Check for entities (functions, classes)
		entities := extractEntities(fullPath)
		if len(entities) == 0 {
			report.AddIssue(DriftIssue{
				File:          filePath,
				Status:        StatusWarning,
				Expected:      "Contains functions/classes",
				Actual:        "No entities found",
				Recommendation: fmt.Sprintf("Add implementation to %s", filePath),
			})
			continue
		}

		// File is OK
		report.AddIssue(DriftIssue{
			File:     filePath,
			Status:   StatusOK,
			Expected: "File exists with content",
			Actual:   fmt.Sprintf("Found %d entities", len(entities)),
		})
	}

	// Generate verdict
	report.Verdict = report.GenerateVerdict()

	return report, nil
}

// extractEntities extracts function and class names from a source file
func extractEntities(filePath string) []string {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	contentStr := string(content)
	entities := []string{}

	// Go patterns
	if strings.HasSuffix(filePath, ".go") {
		// Functions: func Name(
		funcPattern := regexp.MustCompile(`func\s+(\w+)\s*\(`)
		matches := funcPattern.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matches {
			if len(match) > 1 {
				entities = append(entities, "func "+match[1])
			}
		}

		// Types: type Name struct/interface
		typePattern := regexp.MustCompile(`type\s+(\w+)\s+(struct|interface)`)
		matches = typePattern.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matches {
			if len(match) > 1 {
				entities = append(entities, "type "+match[1])
			}
		}
	}

	// Python patterns
	if strings.HasSuffix(filePath, ".py") {
		// Functions: def name(
		funcPattern := regexp.MustCompile(`def\s+(\w+)\s*\(`)
		matches := funcPattern.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matches {
			if len(match) > 1 {
				entities = append(entities, "def "+match[1])
			}
		}

		// Classes: class Name:
		classPattern := regexp.MustCompile(`class\s+(\w+)\s*:`)
		matches = classPattern.FindAllStringSubmatch(contentStr, -1)
		for _, match := range matches {
			if len(match) > 1 {
				entities = append(entities, "class "+match[1])
			}
		}
	}

	return entities
}
