package resolver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BeadsDetector detects if beads CLI is available
type BeadsDetector struct {
	projectDir string
}

// NewBeadsDetector creates a new beads detector
func NewBeadsDetector() *BeadsDetector {
	return &BeadsDetector{projectDir: "."}
}

// NewBeadsDetectorWithDir creates a detector with a specific project directory
func NewBeadsDetectorWithDir(dir string) *BeadsDetector {
	return &BeadsDetector{projectDir: dir}
}

// IsAvailable checks if beads CLI is installed and accessible
func (d *BeadsDetector) IsAvailable() bool {
	_, err := exec.LookPath("bd")
	return err == nil
}

// HasBeadsDirectory checks if .beads directory exists
func (d *BeadsDetector) HasBeadsDirectory() bool {
	beadsDir := filepath.Join(d.projectDir, ".beads")
	info, err := os.Stat(beadsDir)
	return err == nil && info.IsDir()
}

// IsEnabled returns true if beads is fully configured (CLI + directory)
func (d *BeadsDetector) IsEnabled() bool {
	return d.IsAvailable() && d.HasBeadsDirectory()
}

// ResolveBeadsID resolves a beads ID to its linked workstream
func (r *Resolver) ResolveBeadsID(beadsID string) (*Result, error) {
	return r.resolveBeads(beadsID)
}

// FindWorkstreamByBeadsID searches workstream files for a matching beads_id
func (r *Resolver) FindWorkstreamByBeadsID(beadsID string) (wsID string, path string, err error) {
	result, err := r.resolveBeads(beadsID)
	if err != nil {
		return "", "", err
	}
	return result.WSID, result.Path, nil
}

// CreateBeadsIssue creates a beads issue via the CLI
func CreateBeadsIssue(title, issueType string) (string, error) {
	detector := NewBeadsDetector()
	if !detector.IsEnabled() {
		return "", fmt.Errorf("beads not enabled (CLI or .beads directory missing)")
	}

	cmd := exec.Command("bd", "create", "--title", title, "--type", issueType)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("bd create failed: %w", err)
	}

	// Extract beads ID from output
	// Output format: "Created issue: sdp-abc123"
	beadsID := parseBeadsIDFromOutput(string(output))
	if beadsID == "" {
		return "", fmt.Errorf("failed to parse beads ID from output: %s", string(output))
	}

	return beadsID, nil
}

// UpdateBeadsNotes updates beads issue notes with workstream reference
func UpdateBeadsNotes(beadsID, wsPath string) error {
	detector := NewBeadsDetector()
	if !detector.IsEnabled() {
		return nil // Silently skip if beads not enabled
	}

	notes := fmt.Sprintf("Workstream: %s", wsPath)
	cmd := exec.Command("bd", "update", beadsID, "--notes", notes)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("bd update notes failed: %w", err)
	}

	return nil
}

// parseBeadsIDFromOutput extracts beads ID from bd create output
func parseBeadsIDFromOutput(output string) string {
	// Try to find pattern like "Created issue: sdp-abc123"
	// or just the ID itself at the end
	lines := splitLines(output)
	for _, line := range lines {
		// Look for beads ID pattern
		if len(line) >= 7 && line[:4] == "sdp-" {
			return line
		}
		// Look for "Created issue: <id>" pattern
		if len(line) > 15 && line[:15] == "Created issue: " {
			return line[15:]
		}
	}
	return ""
}

func splitLines(s string) []string {
	result := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			result = append(result, line)
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}
