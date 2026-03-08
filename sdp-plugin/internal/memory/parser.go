package memory

import (
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ParseFile parses a markdown file and extracts metadata
func (i *Indexer) ParseFile(content, filename string) (*Artifact, error) {
	artifact := &Artifact{
		Type:      "doc",
		Tags:      []string{},
		IndexedAt: time.Now(),
	}

	fm, mainContent := extractFrontmatter(content)
	artifact.Content = mainContent

	if fm != nil {
		if title, ok := fm["title"].(string); ok {
			artifact.Title = title
		}
		if fid, ok := fm["feature_id"].(string); ok {
			artifact.FeatureID = fid
		}
		if wsID, ok := fm["ws_id"].(string); ok {
			artifact.WorkstreamID = wsID
		}
		if tags, ok := fm["tags"].([]any); ok {
			for _, t := range tags {
				if tag, ok := t.(string); ok {
					artifact.Tags = append(artifact.Tags, tag)
				}
			}
		} else if tags, ok := fm["tags"].(string); ok {
			artifact.Tags = strings.Split(tags, ",")
		}
	}

	// Extract ws_id from filename if not in frontmatter
	if artifact.WorkstreamID == "" {
		wsIDPattern := regexp.MustCompile(`^\d{2}-\d{3}-\d{2}`)
		if match := wsIDPattern.FindString(filename); match != "" {
			artifact.WorkstreamID = match
			if len(match) >= 6 {
				artifact.FeatureID = "F" + match[3:6]
			}
		}
	}

	if artifact.Title == "" {
		artifact.Title = extractFirstHeading(mainContent)
	}
	if artifact.Title == "" {
		artifact.Title = strings.TrimSuffix(filename, ".md")
	}

	return artifact, nil
}

// extractFrontmatter extracts YAML frontmatter from content
func extractFrontmatter(content string) (map[string]any, string) {
	if !strings.HasPrefix(content, "---\n") {
		return nil, content
	}

	fmContent, mainContent, ok := strings.Cut(content[4:], "\n---\n")
	if !ok {
		return nil, content
	}
	mainContent = strings.TrimSpace(mainContent)

	var fm map[string]any
	if err := yaml.Unmarshal([]byte(fmContent), &fm); err != nil {
		return nil, content
	}

	return fm, mainContent
}

// extractFirstHeading extracts the first heading from markdown content
func extractFirstHeading(content string) string {
	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if heading, ok := strings.CutPrefix(line, "# "); ok {
			return heading
		}
		if heading, ok := strings.CutPrefix(line, "## "); ok {
			return heading
		}
	}
	return ""
}
