package beads

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Sync exports Beads state back into the tracked repo snapshot.
func (c *Client) Sync() error {
	if !c.beadsInstalled {
		// Beads not installed, skip sync
		return nil
	}

	cmd := exec.Command("bd", "sync")
	if c.projectRoot != "" {
		cmd.Dir = c.projectRoot
	}
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	if !strings.Contains(string(output), `unknown command "sync"`) {
		return fmt.Errorf("bd sync failed: %w\nOutput: %s", err, string(output))
	}

	snapshotPath := c.issuesSnapshotPath()
	if err := os.MkdirAll(filepath.Dir(snapshotPath), 0o755); err != nil {
		return fmt.Errorf("create .beads directory: %w", err)
	}

	exportCmd := exec.Command("bd", "export", "-o", snapshotPath)
	if c.projectRoot != "" {
		exportCmd.Dir = c.projectRoot
	}
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
