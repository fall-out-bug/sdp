package resolver

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// resolveWorkstream resolves a workstream ID to its file path
func (r *Resolver) resolveWorkstream(wsID string) (*Result, error) {
	path := filepath.Join(r.workstreamDir, wsID+".md")

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("workstream not found: %s", wsID)
	}

	result := &Result{
		Type: TypeWorkstream,
		ID:   wsID,
		Path: path,
	}

	// Extract frontmatter metadata
	r.extractFrontmatter(path, result)

	return result, nil
}

// resolveBeads resolves a beads ID to its linked workstream
func (r *Resolver) resolveBeads(beadsID string) (*Result, error) {
	// Search workstream files for matching beads_id in frontmatter
	matches, err := filepath.Glob(filepath.Join(r.workstreamDir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("failed to search workstreams: %w", err)
	}

	for _, path := range matches {
		frontmatter, err := r.parseFrontmatter(path)
		if err != nil {
			continue
		}

		if frontmatter["beads_id"] == beadsID {
			result := &Result{
				Type: TypeBeads,
				ID:   beadsID,
				WSID: frontmatter["ws_id"],
				Path: path,
			}

			if title, ok := frontmatter["title"]; ok {
				result.Title = title
			}
			if status, ok := frontmatter["status"]; ok {
				result.Status = status
			}

			return result, nil
		}
	}

	return nil, fmt.Errorf("beads ID not found in any workstream: %s", beadsID)
}

// parseFrontmatter extracts YAML frontmatter from a markdown file
func (r *Resolver) parseFrontmatter(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close() //nolint:errcheck // cleanup
	}()

	scanner := bufio.NewScanner(file)
	frontmatter := make(map[string]string)
	inFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			}
			// End of frontmatter
			break
		}

		if inFrontmatter {
			// Parse key: value
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				value = strings.Trim(value, "\"'")
				frontmatter[key] = value
			}
		}
	}

	return frontmatter, scanner.Err()
}

// extractFrontmatter populates result with frontmatter metadata
func (r *Resolver) extractFrontmatter(path string, result *Result) {
	frontmatter, err := r.parseFrontmatter(path)
	if err != nil {
		return
	}

	if title, ok := frontmatter["title"]; ok {
		result.Title = title
	}
	if status, ok := frontmatter["status"]; ok {
		result.Status = status
	}
	if wsID, ok := frontmatter["ws_id"]; ok && result.WSID == "" {
		result.WSID = wsID
	}
}
