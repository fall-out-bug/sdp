package planner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateWorkstreamFiles creates workstream markdown files in the backlog directory.
// AC4: Creates workstream files in docs/workstreams/backlog/
// AC7: In dry-run mode, returns without creating files.
func (p *Planner) CreateWorkstreamFiles(result *DecompositionResult) error {
	// AC7: Dry-run mode - don't create files
	if p.DryRun {
		return nil
	}

	// Ensure backlog directory exists
	if err := os.MkdirAll(p.BacklogDir, 0755); err != nil {
		return fmt.Errorf("create backlog dir: %w", err)
	}

	// Create each workstream file
	for _, ws := range result.Workstreams {
		filename := ws.Filename()
		path := filepath.Join(p.BacklogDir, filename)

		// Check if file already exists
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("workstream file already exists: %s", path)
		}

		// Generate workstream content
		content := p.generateWorkstreamContent(ws, result)

		// Write file
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("write workstream file %s: %w", path, err)
		}
	}

	return nil
}

// generateWorkstreamContent generates the markdown content for a workstream file.
func (p *Planner) generateWorkstreamContent(ws Workstream, result *DecompositionResult) string {
	var sb strings.Builder

	// Frontmatter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("ws_id: %s\n", ws.ID))
	sb.WriteString(fmt.Sprintf("feature: %s\n", result.FeatureID))
	sb.WriteString("status: pending\n")
	if ws.Complexity != "" {
		sb.WriteString(fmt.Sprintf("complexity: %s\n", ws.Complexity))
	}
	sb.WriteString(fmt.Sprintf("project_id: \"%s\"\n", ws.ID[:2]))
	sb.WriteString("---\n\n")

	// Title and header
	sb.WriteString(fmt.Sprintf("# Workstream: %s\n\n", ws.Title))
	sb.WriteString(fmt.Sprintf("**ID:** %s\n", ws.ID))
	sb.WriteString(fmt.Sprintf("**Feature:** %s\n", result.FeatureID))
	sb.WriteString("**Status:** READY\n")
	sb.WriteString("**Owner:** AI Agent\n")
	if ws.Complexity != "" {
		sb.WriteString(fmt.Sprintf("**Complexity:** %s\n", ws.Complexity))
	}
	sb.WriteString("\n---\n\n")

	// Goal section
	sb.WriteString("## Goal\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", ws.Description))
	sb.WriteString("---\n\n")

	// Context section
	sb.WriteString("## Context\n\n")
	sb.WriteString("Generated from feature decomposition.\n\n")
	sb.WriteString("---\n\n")

	// Acceptance Criteria section
	sb.WriteString("## Acceptance Criteria\n\n")
	sb.WriteString("- [ ] AC1: Implementation complete\n")
	sb.WriteString("- [ ] AC2: Tests passing with >= 80% coverage\n")
	sb.WriteString("- [ ] AC3: Code review approved\n")
	sb.WriteString("\n---\n\n")

	// Scope section
	sb.WriteString("## Scope\n\n")
	sb.WriteString("### In Scope\n\n")
	sb.WriteString("- Implementation of core functionality\n")
	sb.WriteString("- Unit and integration tests\n")
	sb.WriteString("- Documentation updates\n\n")
	sb.WriteString("### Out of Scope\n\n")
	sb.WriteString("- Performance optimization (unless specified)\n")
	sb.WriteString("- UI/UX refinements (unless specified)\n")
	sb.WriteString("\n---\n\n")

	// Dependencies section
	sb.WriteString("## Dependencies\n\n")
	deps := []string{}
	for _, dep := range result.Dependencies {
		if dep.From == ws.ID {
			deps = append(deps, fmt.Sprintf("- **%s**: %s", dep.To, dep.Reason))
		}
	}
	if len(deps) > 0 {
		for _, dep := range deps {
			sb.WriteString(fmt.Sprintf("%s\n", dep))
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("No dependencies.\n\n")
	}
	sb.WriteString("---\n\n")

	// Notes section
	sb.WriteString("## Notes\n\n")
	sb.WriteString(fmt.Sprintf("Created: %s\n", result.CreatedAt))
	sb.WriteString(fmt.Sprintf("Source Feature: %s\n", result.Summary))

	return sb.String()
}
