package beadscli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CreateOptions configures bd create in the root module.
type CreateOptions struct {
	Title    string
	Priority string
	Labels   []string
}

// CreateAndSync creates a Beads issue and persists the repo snapshot.
func CreateAndSync(opts CreateOptions) error {
	if strings.TrimSpace(opts.Title) == "" {
		return fmt.Errorf("title is required")
	}

	args := []string{"create", "--title", opts.Title}
	if opts.Priority != "" {
		args = append(args, "--priority", opts.Priority)
	}
	if len(opts.Labels) > 0 {
		args = append(args, "--labels", strings.Join(opts.Labels, ","))
	}

	cmd := exec.Command("bd", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bd create failed: %w\nOutput: %s", err, string(output))
	}

	if err := Sync(); err != nil {
		return fmt.Errorf("bd create succeeded but sync failed: %w", err)
	}

	return nil
}

// Sync persists the Beads snapshot to .beads/issues.jsonl.
func Sync() error {
	cmd := exec.Command("bd", "sync")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	if !strings.Contains(string(output), `unknown command "sync"`) {
		return fmt.Errorf("bd sync failed: %w\nOutput: %s", err, string(output))
	}

	snapshotPath := filepath.Join(".beads", "issues.jsonl")
	if err := os.MkdirAll(filepath.Dir(snapshotPath), 0o755); err != nil {
		return fmt.Errorf("create .beads directory: %w", err)
	}

	exportCmd := exec.Command("bd", "export", "-o", snapshotPath)
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
