package beads

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Sync exports Beads state back into the tracked repo snapshot.
func (c *Client) Sync() error {
	if !c.beadsInstalled {
		// Beads not installed, skip sync
		return nil
	}

	cmd := exec.Command("bd", "sync")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	if !strings.Contains(string(output), `unknown command "sync"`) {
		return fmt.Errorf("bd sync failed: %w\nOutput: %s", err, string(output))
	}

	if err := os.MkdirAll(".beads", 0o755); err != nil {
		return fmt.Errorf("create .beads directory: %w", err)
	}

	exportCmd := exec.Command("bd", "export", "-o", ".beads/issues.jsonl")
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
