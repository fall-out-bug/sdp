package beads

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
)

// Client wraps the Beads CLI for task tracking
type Client struct {
	mappingPath    string
	beadsInstalled bool
	projectRoot    string
}

// Task represents a Beads task
type Task struct {
	ID       string
	Title    string
	Status   string
	Priority string
}

// CreateOptions configures bd create invocations.
type CreateOptions struct {
	Type        string
	Priority    string
	Labels      []string
	Description string
	Parent      string
	Silent      bool
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
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		projectRoot, err = os.Getwd()
		if err != nil {
			projectRoot = "."
		}
	}
	mappingPath := filepath.Join(projectRoot, ".beads-sdp-mapping.jsonl")

	return &Client{
		mappingPath:    mappingPath,
		beadsInstalled: beadsInstalled,
		projectRoot:    projectRoot,
	}, nil
}

// Ready returns available tasks
func (c *Client) Ready() ([]Task, error) {
	if !c.beadsInstalled {
		return []Task{}, nil
	}

	output, err := c.runBeadsCommand("ready")
	if err != nil {
		if errors.Is(err, ErrNoBeadsDatabase) {
			return []Task{}, nil
		}
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
	for line := range strings.SplitSeq(output, "\n") {
		if value, ok := strings.CutPrefix(line, "Title:"); ok {
			task.Title = strings.TrimSpace(value)
		} else if value, ok := strings.CutPrefix(line, "Status:"); ok {
			task.Status = strings.TrimSpace(value)
		} else if value, ok := strings.CutPrefix(line, "Priority:"); ok {
			task.Priority = strings.TrimSpace(value)
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

// Create creates a Beads issue and returns its ID.
func (c *Client) Create(title string, opts CreateOptions) (string, error) {
	if !c.beadsInstalled {
		return "", fmt.Errorf("beads CLI not installed")
	}
	if strings.TrimSpace(title) == "" {
		return "", fmt.Errorf("title is required")
	}

	args := []string{"create", "--title", title}
	if opts.Type != "" {
		args = append(args, "--type", opts.Type)
	}
	if opts.Priority != "" {
		args = append(args, "--priority", opts.Priority)
	}
	if len(opts.Labels) > 0 {
		args = append(args, "--labels", strings.Join(opts.Labels, ","))
	}
	if opts.Description != "" {
		args = append(args, "--description", opts.Description)
	}
	if opts.Parent != "" {
		args = append(args, "--parent", opts.Parent)
	}
	if opts.Silent {
		args = append(args, "--silent")
	}

	output, err := c.runBeadsCommand(args...)
	if err != nil {
		return "", fmt.Errorf("bd create failed: %w", err)
	}

	beadsID := parseCreatedTaskID(output)
	if beadsID == "" {
		return "", fmt.Errorf("bd create succeeded but no beads ID found in output: %q", strings.TrimSpace(output))
	}

	return beadsID, nil
}

// Close closes a Beads issue with an optional reason.
func (c *Client) Close(beadsID, reason string) error {
	if !c.beadsInstalled {
		return fmt.Errorf("beads CLI not installed")
	}
	if strings.TrimSpace(beadsID) == "" {
		return fmt.Errorf("beads ID is required")
	}

	args := []string{"close", beadsID}
	if strings.TrimSpace(reason) != "" {
		args = append(args, "--reason", reason)
	}

	if _, err := c.runBeadsCommand(args...); err != nil {
		return fmt.Errorf("bd close failed: %w", err)
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
		if os.IsNotExist(err) {
			return []mappingEntry{}, nil
		}
		return nil, fmt.Errorf("failed to open mapping file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close mapping file: %v\n", cerr)
		}
	}()

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

	for line := range strings.SplitSeq(output, "\n") {
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

var createdTaskIDRe = regexp.MustCompile(`\b[a-z][a-z0-9_-]*-[a-z0-9][a-z0-9_-]*\b`)

func parseCreatedTaskID(output string) string {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if value, ok := strings.CutPrefix(line, "Created issue:"); ok {
			return strings.TrimSpace(value)
		}
		if value, ok := strings.CutPrefix(line, "Created:"); ok {
			return strings.TrimSpace(value)
		}
		if match := createdTaskIDRe.FindString(line); match != "" {
			return match
		}
	}
	return ""
}
