package task

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateWorkstream creates a workstream file from a task definition
func (c *Creator) CreateWorkstream(task *Task) (*Workstream, error) {
	if err := c.validateWorkstream(task); err != nil {
		return nil, err
	}

	prefix := c.wsPrefix(task.Type, task.FeatureID)
	seq := c.nextSequence(prefix)
	wsID := c.generateWSID(task.Type, task.FeatureID, seq)

	if err := os.MkdirAll(c.config.WorkstreamDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workstream dir: %w", err)
	}

	content := c.generateWorkstreamContent(wsID, task)
	path := filepath.Join(c.config.WorkstreamDir, wsID+".md")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write workstream: %w", err)
	}

	return &Workstream{
		WSID:      wsID,
		Path:      path,
		FeatureID: task.FeatureID,
		BeadsID:   task.BeadsID,
		CreatedAt: time.Now(),
	}, nil
}

// validateWorkstream validates task for workstream creation
func (c *Creator) validateWorkstream(task *Task) error {
	if task.Title == "" {
		return fmt.Errorf("task title is required")
	}
	if task.FeatureID == "" {
		return fmt.Errorf("feature ID is required for workstream creation")
	}
	return nil
}

// wsPrefix generates the prefix for workstream ID based on type
func (c *Creator) wsPrefix(taskType Type, featureID string) string {
	featureNum := strings.TrimPrefix(featureID, "F")

	switch taskType {
	case TypeBug, TypeHotfix:
		return fmt.Sprintf("99-%s", featureNum)
	default:
		return fmt.Sprintf("00-%s", featureNum)
	}
}

// generateWSID generates a workstream ID
func (c *Creator) generateWSID(taskType Type, featureID string, seq int) string {
	featureNum := strings.TrimPrefix(featureID, "F")

	switch taskType {
	case TypeBug, TypeHotfix:
		return fmt.Sprintf("99-%s-%02d", featureNum, seq)
	default:
		return fmt.Sprintf("00-%s-%02d", featureNum, seq)
	}
}

// issueIndexEntry represents the index file entry (shared with template.go)
type issueIndexEntry struct {
	IssueID  string `json:"issue_id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority int    `json:"priority"`
	File     string `json:"file"`
}

// ReadIndex reads all entries from the issues index
func (c *Creator) ReadIndex() ([]issueIndexEntry, error) {
	f, err := os.Open(c.config.IndexFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func() { _ = f.Close() }() //nolint:errcheck // cleanup

	var entries []issueIndexEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry issueIndexEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	return entries, scanner.Err()
}
