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
	// Do not assume "---" at position 0: use Index to find both delimiters.
	contentStr := string(content)
	openDelim := strings.Index(contentStr, "---")
	if openDelim == -1 {
		return nil, fmt.Errorf("no frontmatter found in %s", wsPath)
	}
	afterOpen := openDelim + 3 // skip "---"
	closeDelim := strings.Index(contentStr[afterOpen:], "---")
	if closeDelim == -1 {
		return nil, fmt.Errorf("no frontmatter found in %s", wsPath)
	}
	frontmatter := contentStr[afterOpen : afterOpen+closeDelim]

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
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(strings.Trim(parts[1], `"'`))

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
		} else if inList && strings.HasPrefix(line, "-") {
			// Parse list items
			value := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			value = strings.Trim(value, `"'`)
			*currentList = append(*currentList, value)
		}
	}

	return data, nil
}
