package beads

import (
	"fmt"
	"os"
	"os/exec"
)

// runBeadsCommand executes a Beads CLI command
func (c *Client) runBeadsCommand(args ...string) (string, error) {
	cmd := exec.Command("bd", args...)
	output, err := cmd.Output()
	if err != nil {
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
