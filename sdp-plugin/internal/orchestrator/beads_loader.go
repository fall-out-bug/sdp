package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// BeadsLoader loads workstreams from Beads and workstream files
type BeadsLoader struct {
	workstreamDir string
	mappingPath   string
}

// NewBeadsLoader creates a new Beads-based workstream loader
func NewBeadsLoader(workstreamDir, mappingPath string) *BeadsLoader {
	return &BeadsLoader{
		workstreamDir: workstreamDir,
		mappingPath:   mappingPath,
	}
}

// LoadWorkstreams loads all workstreams for a feature
func (b *BeadsLoader) LoadWorkstreams(featureID string) ([]WorkstreamNode, error) {
	// Get all workstream files (pattern: PP-FFF-SS.md)
	matches, err := filepath.Glob(filepath.Join(b.workstreamDir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("failed to glob workstream files: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("%w: no workstreams found in %s", ErrFeatureNotFound, b.workstreamDir)
	}

	// Pre-allocate slice with known capacity
	workstreams := make([]WorkstreamNode, 0, len(matches))

	for _, match := range matches {
		// Parse workstream file
		ws, err := parser.ParseWorkstream(match)
		if err != nil {
			return nil, fmt.Errorf("failed to parse workstream %s: %w", match, err)
		}

		// Skip if not for this feature
		if ws.Feature != featureID {
			continue
		}

		// Extract dependencies from workstream file
		deps := b.extractDependencies(ws)

		node := WorkstreamNode{
			ID:           ws.ID,
			Feature:      ws.Feature,
			Status:       ws.Status,
			Dependencies: deps,
		}

		workstreams = append(workstreams, node)
	}

	if len(workstreams) == 0 {
		return nil, fmt.Errorf("%w: feature %s has no workstreams", ErrFeatureNotFound, featureID)
	}

	return workstreams, nil
}

// extractDependencies extracts dependencies from a workstream
func (b *BeadsLoader) extractDependencies(ws *parser.Workstream) []string {
	// Read the workstream file to find dependencies
	wsPath := filepath.Join(b.workstreamDir, ws.ID+".md")

	content, err := os.ReadFile(wsPath)
	if err != nil {
		return []string{}
	}

	// Parse dependencies from the content
	// Look for "Dependencies:" in the frontmatter or content
	lines := strings.Split(string(content), "\n")

	inDepsSection := false
	var deps []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check for Dependencies section (### or ** or plain)
		if strings.HasPrefix(trimmed, "Dependencies:") ||
			strings.HasPrefix(trimmed, "**Dependencies:**") ||
			strings.HasPrefix(trimmed, "### Dependencies") ||
			strings.HasPrefix(trimmed, "###Dependencies") {
			inDepsSection = true
			continue
		}

		// Check for next section (any heading marker)
		if inDepsSection && (strings.HasPrefix(trimmed, "##") || strings.HasPrefix(trimmed, "**")) {
			// Only exit if it's a different section
			if !strings.Contains(strings.ToLower(trimmed), "dependencies") {
				inDepsSection = false
			}
			continue
		}

		// Extract dependency IDs (must start with "-")
		if inDepsSection && strings.HasPrefix(trimmed, "-") {
			depID := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			depID = strings.TrimPrefix(depID, "`")
			depID = strings.TrimSuffix(depID, "`")
			depID = strings.TrimSpace(depID)
			if depID != "" {
				deps = append(deps, depID)
			}
		}
	}

	return deps
}

// CLIExecutor executes workstreams via the SDP CLI
type CLIExecutor struct {
	sdpCommand string
}

// NewCLIExecutor creates a new CLI executor
func NewCLIExecutor(sdpCommand string) *CLIExecutor {
	if sdpCommand == "" {
		sdpCommand = "sdp"
	}
	return &CLIExecutor{sdpCommand: sdpCommand}
}

// Execute executes a workstream by calling the SDP CLI
func (e *CLIExecutor) Execute(wsID string) error {
	// Call: sdp build <ws-id>
	cmd := exec.Command(e.sdpCommand, "build", wsID)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("sdp build %s failed: %w\nOutput: %s", wsID, err, string(output))
	}

	return nil
}
