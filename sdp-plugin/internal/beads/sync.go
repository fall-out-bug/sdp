package beads

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Sync persists Beads state into the tracked repo snapshot.
func (c *Client) Sync() error {
	if !c.beadsInstalled {
		// Beads not installed, skip sync
		return nil
	}

	projectRoot, err := findBeadsProjectRoot()
	if err != nil {
		return err
	}

	cmd := exec.Command("bd", "sync")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	if !strings.Contains(string(output), `unknown command "sync"`) {
		return fmt.Errorf("bd sync failed: %w\nOutput: %s", err, string(output))
	}

	snapshotPath := filepath.Join(projectRoot, ".beads", "issues.jsonl")
	if err := os.MkdirAll(filepath.Dir(snapshotPath), 0o755); err != nil {
		return fmt.Errorf("create .beads directory: %w", err)
	}

	exportCmd := exec.Command("bd", "export", "-o", snapshotPath)
	exportCmd.Dir = projectRoot
	exportOutput, exportErr := exportCmd.CombinedOutput()
	if exportErr != nil {
		return fmt.Errorf(
			"bd sync unavailable and export fallback failed: %w\nSync output: %s\nExport output: %s",
			exportErr,
			string(output),
			string(exportOutput),
		)
	}

	return nil
}

func findBeadsProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}

	current := cwd
	for {
		for _, marker := range []string{".beads", ".git"} {
			if _, err := os.Stat(filepath.Join(current, marker)); err == nil {
				return current, nil
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			return cwd, nil
		}
		current = parent
	}
}
