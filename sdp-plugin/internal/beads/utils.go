package beads

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrNoBeadsDatabase = errors.New("no beads database found")

// runBeadsCommand executes a Beads CLI command
func (c *Client) runBeadsCommand(args ...string) (string, error) {
	cmd := exec.Command("bd", args...)
	if c.projectRoot != "" {
		cmd.Dir = c.projectRoot
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		if strings.Contains(trimmed, "no beads database found") {
			return "", fmt.Errorf("%w: %s", ErrNoBeadsDatabase, trimmed)
		}
		if trimmed != "" {
			return "", fmt.Errorf("command failed: %s: %w", trimmed, err)
		}
		return "", fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// isBeadsInstalled checks if Beads CLI is available
func isBeadsInstalled() bool {
	_, err := exec.LookPath("bd")
	return err == nil
}

// findMappingFile finds the Beads mapping file
func findMappingFile() (string, error) {
	// Try common locations
	locations := []string{
		".beads-sdp-mapping.jsonl",
		"../.beads-sdp-mapping.jsonl",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			return loc, nil
		}
	}

	return "", fmt.Errorf("mapping file not found")
}

func (c *Client) issuesSnapshotPath() string {
	if c.projectRoot == "" {
		return filepath.Join(".beads", "issues.jsonl")
	}
	return filepath.Join(c.projectRoot, ".beads", "issues.jsonl")
}
