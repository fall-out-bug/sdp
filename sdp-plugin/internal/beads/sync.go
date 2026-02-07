package beads

import (
	"fmt"
	"os/exec"
)

// Sync runs "bd sync" to synchronize Beads state
func (c *Client) Sync() error {
	if !c.beadsInstalled {
		// Beads not installed, skip sync
		return nil
	}

	cmd := exec.Command("bd", "sync")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bd sync failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}
