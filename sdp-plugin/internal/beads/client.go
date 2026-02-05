package beads

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Client wraps the Beads CLI for task tracking
type Client struct {
	mappingPath     string
	beadsInstalled bool
}

// Task represents a Beads task
type Task struct {
	ID       string
	Title    string
	Status   string
	Priority string
}

// mappingEntry represents a line in the mapping file
type mappingEntry struct {
	SdpID     string `json:"sdp_id"`
	BeadsID   string `json:"beads_id"`
	UpdatedAt string `json:"updated_at"`
}

// NewClient creates a new Beads client
func NewClient() (*Client, error) {
	beadsInstalled := isBeadsInstalled()
	mappingPath, err := findMappingFile()
	if err != nil {
		mappingPath = ".beads-sdp-mapping.jsonl"
	}

	return &Client{
		mappingPath:     mappingPath,
		beadsInstalled: beadsInstalled,
	}, nil
}

// Ready returns available tasks
func (c *Client) Ready() ([]Task, error) {
	if !c.beadsInstalled {
		return []Task{}, nil
	}

	output, err := c.runBeadsCommand("ready")
	if err != nil {
		return []Task{}, fmt.Errorf("bd ready failed: %w", err)
	}

	tasks := c.parseTaskList(output)
	if tasks == nil {
		return []Task{}, nil
	}
	return tasks, nil
}

// Show returns details of a specific task
func (c *Client) Show(beadsID string) (*Task, error) {
	if !c.beadsInstalled {
		return nil, fmt.Errorf("beads CLI not installed")
	}

	output, err := c.runBeadsCommand("show", beadsID)
	if err != nil {
		return nil, fmt.Errorf("bd show failed: %w", err)
	}

	task := &Task{ID: beadsID}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Title:") {
			task.Title = strings.TrimSpace(strings.TrimPrefix(line, "Title:"))
		} else if strings.HasPrefix(line, "Status:") {
			task.Status = strings.TrimSpace(strings.TrimPrefix(line, "Status:"))
		} else if strings.HasPrefix(line, "Priority:") {
			task.Priority = strings.TrimSpace(strings.TrimPrefix(line, "Priority:"))
		}
	}

	return task, nil
}

// Update updates the status of a task
func (c *Client) Update(beadsID string, status string) error {
	if !c.beadsInstalled {
		return fmt.Errorf("beads CLI not installed")
	}

	_, err := c.runBeadsCommand("update", beadsID, "--status", status)
	if err != nil {
		return fmt.Errorf("bd update failed: %w", err)
	}

	return nil
}

// MapWSToBeads converts workstream ID to Beads ID
func (c *Client) MapWSToBeads(wsID string) (string, error) {
	entries, err := c.readMapping()
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.SdpID == wsID {
			return entry.BeadsID, nil
		}
	}

	return "", fmt.Errorf("workstream ID not found in mapping: %s", wsID)
}

// MapBeadsToWS converts Beads ID to workstream ID
func (c *Client) MapBeadsToWS(beadsID string) (string, error) {
	entries, err := c.readMapping()
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.BeadsID == beadsID {
			return entry.SdpID, nil
		}
	}

	return "", fmt.Errorf("beads ID not found in mapping: %s", beadsID)
}

// readMapping reads the mapping file
func (c *Client) readMapping() ([]mappingEntry, error) {
	file, err := os.Open(c.mappingPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open mapping file: %w", err)
	}
	defer file.Close()

	var entries []mappingEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry mappingEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read mapping file: %w", err)
	}

	return entries, nil
}

// parseTaskList parses the output of "bd ready"
func (c *Client) parseTaskList(output string) []Task {
	var tasks []Task
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, "sdp-") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) > 0 {
			taskID := parts[0]
			if strings.HasPrefix(taskID, "sdp-") {
				task := Task{
					ID:    taskID,
					Title: strings.Join(parts[1:], " "),
				}
				tasks = append(tasks, task)
			}
		}
	}

	if tasks == nil {
		return []Task{}
	}
	return tasks
}
