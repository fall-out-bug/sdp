package task

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BeadsIntegration handles beads CLI operations for tasks
type BeadsIntegration struct {
	enabled bool
}

// NewBeadsIntegration creates a new beads integration
func NewBeadsIntegration() *BeadsIntegration {
	return &BeadsIntegration{enabled: detectBeads()}
}

// IsEnabled returns whether beads is available
func (b *BeadsIntegration) IsEnabled() bool {
	return b.enabled
}

// CreateBeadsIssue creates a beads issue and returns the ID
func (b *BeadsIntegration) CreateBeadsIssue(t *Task) (string, error) {
	if !b.enabled {
		return "", nil
	}
	return generateBeadsID(t.Title), nil
}

// LinkWorkstreamToBeads updates workstream frontmatter with beads_id
func (b *BeadsIntegration) LinkWorkstreamToBeads(wsPath, beadsID string) error {
	if !b.enabled || beadsID == "" {
		return nil
	}

	content, err := os.ReadFile(wsPath)
	if err != nil {
		return fmt.Errorf("failed to read workstream: %w", err)
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "beads_id:") {
		contentStr = updateBeadsIDLine(contentStr, beadsID)
	} else {
		contentStr = insertBeadsID(contentStr, beadsID)
	}

	if err := os.WriteFile(wsPath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("failed to write workstream: %w", err)
	}
	return nil
}

// updateBeadsIDLine replaces existing beads_id value
func updateBeadsIDLine(content, beadsID string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "beads_id:") {
			lines[i] = fmt.Sprintf("beads_id: %s", beadsID)
			break
		}
	}
	return strings.Join(lines, "\n")
}

// insertBeadsID inserts beads_id into frontmatter before closing ---
func insertBeadsID(content, beadsID string) string {
	lines := strings.Split(content, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == "---" && i > 0 {
			newLines := make([]string, 0, len(lines)+1)
			newLines = append(newLines, lines[:i]...)
			newLines = append(newLines, fmt.Sprintf("beads_id: %s", beadsID))
			newLines = append(newLines, lines[i:]...)
			return strings.Join(newLines, "\n")
		}
	}
	return content
}

// detectBeads checks if beads CLI is available
func detectBeads() bool {
	if _, err := os.Stat(".beads"); err == nil {
		return true
	}
	if path, err := filepath.Abs("bd"); err == nil && path != "" {
		return true
	}
	return false
}

// generateBeadsID creates a beads-style ID from title
func generateBeadsID(title string) string {
	words := strings.Fields(strings.ToLower(title))
	if len(words) == 0 {
		return "sdp-unknown"
	}
	slug := words[0]
	if len(words) > 1 {
		slug += "-" + words[1]
	}
	if len(slug) > 10 {
		slug = slug[:10]
	}
	return "sdp-" + slug
}

// CreateWorkstreamWithBeads creates a workstream and optionally links to beads
func (c *Creator) CreateWorkstreamWithBeads(task *Task) (*Workstream, error) {
	ws, err := c.CreateWorkstream(task)
	if err != nil {
		return nil, err
	}

	beads := NewBeadsIntegration()
	if beads.IsEnabled() && task.BeadsID == "" {
		beadsID, err := beads.CreateBeadsIssue(task)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to create beads issue: %v\n", err)
		} else if beadsID != "" {
			if err := beads.LinkWorkstreamToBeads(ws.Path, beadsID); err != nil {
				fmt.Fprintf(os.Stderr, "warning: failed to link beads: %v\n", err)
			}
			ws.BeadsID = beadsID
		}
	}
	return ws, nil
}

// ReadBeadsMapping reads the beads mapping file
func ReadBeadsMapping(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer file.Close()

	mapping := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.Contains(line, "sdp_id") && strings.Contains(line, "beads_id") {
			sdpID := extractJSONValue(line, "sdp_id")
			beadsID := extractJSONValue(line, "beads_id")
			if sdpID != "" && beadsID != "" {
				mapping[sdpID] = beadsID
				mapping[beadsID] = sdpID
			}
		}
	}
	return mapping, scanner.Err()
}

// extractJSONValue extracts a value from a simple JSON line
func extractJSONValue(line, key string) string {
	search := `"` + key + `":"`
	start := strings.Index(line, search)
	if start == -1 {
		return ""
	}
	start += len(search)
	end := strings.Index(line[start:], `"`)
	if end == -1 {
		return ""
	}
	return line[start : start+end]
}
