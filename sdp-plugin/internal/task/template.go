package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateIssue creates an issue file (Beads-free fallback)
func (c *Creator) CreateIssue(task *Task) (*Issue, error) {
	if task.Title == "" {
		return nil, fmt.Errorf("task title is required")
	}

	// Ensure directory exists
	if err := os.MkdirAll(c.config.IssuesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create issues dir: %w", err)
	}

	// Get next issue number
	seq := c.nextIssueSequence()
	issueID := fmt.Sprintf("ISSUE-%04d", seq)

	// Generate content
	content := c.generateIssueContent(issueID, task)
	path := filepath.Join(c.config.IssuesDir, issueID+".md")

	// Write file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write issue: %w", err)
	}

	// Update index
	if err := c.updateIndex(issueID, task, path); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to update index: %v\n", err)
	}

	return &Issue{
		IssueID:   issueID,
		Path:      path,
		CreatedAt: time.Now(),
	}, nil
}

// generateWorkstreamContent generates workstream markdown content
func (c *Creator) generateWorkstreamContent(wsID string, task *Task) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("ws_id: %s\n", wsID))
	sb.WriteString(fmt.Sprintf("feature_id: %s\n", task.FeatureID))
	sb.WriteString(fmt.Sprintf("type: %s\n", task.Type))
	sb.WriteString(fmt.Sprintf("priority: %d\n", task.Priority))
	sb.WriteString(fmt.Sprintf("title: %q\n", task.Title))
	sb.WriteString("status: backlog\n")

	if len(task.DependsOn) > 0 {
		sb.WriteString("depends_on:\n")
		for _, dep := range task.DependsOn {
			sb.WriteString(fmt.Sprintf("  - %s\n", dep))
		}
	} else {
		sb.WriteString("depends_on: []\n")
	}

	sb.WriteString("blocks: []\n")
	sb.WriteString(fmt.Sprintf("project_id: %s\n", c.config.ProjectID))

	if task.BeadsID != "" {
		sb.WriteString(fmt.Sprintf("beads_id: %s\n", task.BeadsID))
	}

	if task.BranchBase != "" {
		sb.WriteString(fmt.Sprintf("branch_base: %s\n", task.BranchBase))
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## Goal\n\n")
	sb.WriteString(task.Goal + "\n\n")

	if task.Context != "" {
		sb.WriteString("## Context\n\n")
		sb.WriteString(task.Context + "\n\n")
	}

	sb.WriteString("## Acceptance Criteria\n\n")
	sb.WriteString("- [ ] AC1: TBD\n\n")

	if len(task.ScopeFiles) > 0 {
		sb.WriteString("## Scope Files\n\n")
		sb.WriteString("```yaml\n")
		sb.WriteString("scope_files:\n")
		for _, f := range task.ScopeFiles {
			sb.WriteString(fmt.Sprintf("  - %s\n", f))
		}
		sb.WriteString("```\n\n")
	}

	sb.WriteString("## Notes\n\n")
	if task.BeadsID != "" {
		sb.WriteString(fmt.Sprintf("- Beads: %s\n", task.BeadsID))
	}

	return sb.String()
}

// generateIssueContent generates issue markdown content
func (c *Creator) generateIssueContent(issueID string, task *Task) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("issue_id: %s\n", issueID))
	sb.WriteString(fmt.Sprintf("title: %q\n", task.Title))
	sb.WriteString("status: open\n")
	sb.WriteString(fmt.Sprintf("priority: %d\n", task.Priority))
	sb.WriteString(fmt.Sprintf("type: %s\n", task.Type))
	sb.WriteString(fmt.Sprintf("created_at: %s\n", time.Now().Format(time.RFC3339)))
	sb.WriteString("beads_id: null\n")

	if task.FeatureID != "" {
		sb.WriteString(fmt.Sprintf("feature_id: %s\n", task.FeatureID))
	} else {
		sb.WriteString("feature_id: null\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## Symptom\n\n")
	sb.WriteString(task.Context + "\n\n")

	sb.WriteString("## Classification\n\n")
	sb.WriteString(fmt.Sprintf("- Severity: P%d\n", task.Priority))
	sb.WriteString(fmt.Sprintf("- Route: /%s\n", task.Type))

	if len(task.ScopeFiles) > 0 {
		sb.WriteString("\n## Scope Files\n\n")
		for _, f := range task.ScopeFiles {
			sb.WriteString(fmt.Sprintf("- %s\n", f))
		}
	}

	return sb.String()
}

// updateIndex appends entry to issues index
func (c *Creator) updateIndex(issueID string, task *Task, path string) error {
	if err := os.MkdirAll(filepath.Dir(c.config.IndexFile), 0755); err != nil {
		return err
	}

	relPath := filepath.Join("docs/issues", issueID+".md")

	entry := issueIndexEntry{
		IssueID:  issueID,
		Title:    task.Title,
		Status:   "open",
		Priority: int(task.Priority),
		File:     relPath,
	}

	f, err := os.OpenFile(c.config.IndexFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(entry)
}
