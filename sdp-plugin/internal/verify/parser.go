package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Parser handles workstream file parsing
type Parser struct {
	wsDir string
}

// NewParser creates a new workstream parser
func NewParser(wsDir string) *Parser {
	return &Parser{
		wsDir: wsDir,
	}
}

// FindWSFile locates a workstream file by ID
func (p *Parser) FindWSFile(wsID string) (string, error) {
	// Check backlog first
	backlogPath := filepath.Join(p.wsDir, "backlog")
	if files, err := filepath.Glob(filepath.Join(backlogPath, wsID+"*.md")); err == nil && len(files) > 0 {
		return files[0], nil
	}

	// Check in_progress
	inProgressPath := filepath.Join(p.wsDir, "in_progress")
	if files, err := filepath.Glob(filepath.Join(inProgressPath, wsID+"*.md")); err == nil && len(files) > 0 {
		return files[0], nil
	}

	// Check completed
	completedPath := filepath.Join(p.wsDir, "completed")
	if files, err := filepath.Glob(filepath.Join(completedPath, wsID+"*.md")); err == nil && len(files) > 0 {
		return files[0], nil
	}

	return "", fmt.Errorf("workstream file not found: %s", wsID)
}

// ParseWSFile parses a workstream markdown file
func (p *Parser) ParseWSFile(wsPath string) (*WorkstreamData, error) {
	// Read file
	content, err := os.ReadFile(wsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract frontmatter (between --- markers)
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---") {
		return nil, fmt.Errorf("no frontmatter found in %s", wsPath)
	}
	afterOpen := strings.TrimPrefix(contentStr, "---")
	frontmatter, _, ok := strings.Cut(afterOpen, "---")
	if !ok {
		return nil, fmt.Errorf("no frontmatter found in %s", wsPath)
	}

	// Parse frontmatter fields
	data := &WorkstreamData{
		ScopeFiles:           []string{},
		VerificationCommands: []string{},
		CoverageThreshold:    80.0, // Default threshold
	}

	lines := strings.Split(frontmatter, "\n")
	inList := false
	currentList := &data.ScopeFiles

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and markers
		if line == "" || line == "---" {
			continue
		}

		// Parse key-value pairs
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "-") {
			inList = false
			key, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}

			key = strings.TrimSpace(key)
			value = strings.TrimSpace(strings.Trim(value, `"'`))

			switch key {
			case "ws_id":
				data.WSID = value
			case "title":
				data.Title = value
			case "status":
				data.Status = value
			case "coverage_threshold":
				if _, err := fmt.Sscanf(value, "%f", &data.CoverageThreshold); err != nil {
					// Use default value if parsing fails
					data.CoverageThreshold = 80.0
				}
			case "scope_files":
				inList = true
				currentList = &data.ScopeFiles
			case "verification_commands":
				inList = true
				currentList = &data.VerificationCommands
			}
		} else if inList {
			// Parse list items
			if value, ok := strings.CutPrefix(line, "-"); ok {
				value = strings.TrimSpace(value)
				value = strings.Trim(value, `"'`)
				*currentList = append(*currentList, value)
			}
		}
	}

	return data, nil
}
