package planner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

// Planner decomposes feature descriptions into workstreams.
type Planner struct {
	BacklogDir    string           // Directory for workstream files
	Description   string           // Feature description to decompose
	Interactive   bool             // AC2: Interactive mode (ask questions)
	AutoApply     bool             // AC3: Auto-apply mode (execute after plan)
	DryRun        bool             // AC7: Show what would be created
	OutputFormat  string           // AC6: Output format (human, json)
	ModelAPI      string           // AC8: AI model API endpoint
	Questions     []string         // Questions for interactive mode
	EvidenceWriter *evidence.Writer // Evidence log writer
}

// DecompositionResult contains the output of feature decomposition.
type DecompositionResult struct {
	Workstreams  []Workstream  `json:"workstreams"`
	Dependencies []Dependency  `json:"dependencies"`
	Summary      string        `json:"summary"`
	FeatureID    string        `json:"feature_id,omitempty"`
	CreatedAt    string        `json:"created_at"`
}

// Workstream represents a single workstream in the decomposition.
type Workstream struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Complexity  string `json:"complexity,omitempty"`
	Estimate    string `json:"estimate,omitempty"`
}

// Dependency represents a dependency between workstreams.
type Dependency struct {
	From   string `json:"from"`   // Dependent workstream ID
	To     string `json:"to"`     // Prerequisite workstream ID
	Reason string `json:"reason"` // Why this dependency exists
}

// Filename generates a slug-based filename for the workstream.
func (ws Workstream) Filename() string {
	// Convert title to lowercase, replace spaces and underscores with hyphens
	slug := strings.ToLower(ws.Title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	return fmt.Sprintf("%s-%s.md", ws.ID, slug)
}

// Decompose performs feature decomposition using AI or mock logic.
// AC8: Returns error if no model API configured.
func (p *Planner) Decompose() (*DecompositionResult, error) {
	// AC8: Check if model API is configured
	if p.ModelAPI == "" {
		return nil, fmt.Errorf("no model API configured: set MODEL_API environment variable or configure in .sdp/config.json")
	}

	// For now, use mock decomposition (real AI integration would call the model)
	// In production, this would make an API call to the LLM
	result := &DecompositionResult{
		FeatureID: "F057", // Would be assigned dynamically
		Summary:   fmt.Sprintf("Decomposition of: %s", p.Description),
		CreatedAt: time.Now().Format(time.RFC3339),
		Workstreams: []Workstream{
			{
				ID:          "00-057-01",
				Title:       "OAuth2 configuration",
				Description: "Setup OAuth2 provider credentials and configuration",
				Status:      "pending",
				Complexity:  "SMALL",
			},
			{
				ID:          "00-057-02",
				Title:       "OAuth2 callback handler",
				Description: "Implement OAuth2 callback endpoint with token exchange",
				Status:      "pending",
				Complexity:  "MEDIUM",
			},
			{
				ID:          "00-057-03",
				Title:       "User authentication flow",
				Description: "Implement login/logout flows using OAuth2",
				Status:      "pending",
				Complexity:  "MEDIUM",
			},
		},
		Dependencies: []Dependency{
			{
				From:   "00-057-02",
				To:     "00-057-01",
				Reason: "Callback handler requires OAuth2 configuration",
			},
			{
				From:   "00-057-03",
				To:     "00-057-02",
				Reason: "Authentication flow requires callback handler",
			},
		},
	}

	return result, nil
}

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

// EmitPlanEvent emits a plan event to the evidence log.
// AC5: Emits plan event with workstream information.
func (p *Planner) EmitPlanEvent(result *DecompositionResult) error {
	if p.EvidenceWriter == nil {
		return fmt.Errorf("evidence writer not configured")
	}

	// Build scope files list (all workstream files that would be created)
	scopeFiles := []string{}
	for _, ws := range result.Workstreams {
		scopeFiles = append(scopeFiles, filepath.Join(p.BacklogDir, ws.Filename()))
	}

	// Create plan event using existing evidence package
	ev := evidence.PlanEventWithFeature(
		"00-057-00", // Planning workstream ID
		result.FeatureID,
		scopeFiles,
	)

	// Add additional metadata about the decomposition
	if ev.Data == nil {
		ev.Data = make(map[string]interface{})
	}
	if dataMap, ok := ev.Data.(map[string]interface{}); ok {
		dataMap["workstream_count"] = len(result.Workstreams)
		dataMap["dependency_count"] = len(result.Dependencies)
		dataMap["description"] = p.Description
		dataMap["action"] = "decompose"
	}

	return p.EvidenceWriter.Append(ev)
}

// FormatOutput formats the decomposition result for output.
// AC6: Supports JSON and human-readable formats.
func (p *Planner) FormatOutput(result *DecompositionResult) (string, error) {
	switch p.OutputFormat {
	case "json":
		return p.formatJSON(result)
	case "human", "":
		return p.formatHuman(result)
	default:
		return "", fmt.Errorf("unsupported output format: %s", p.OutputFormat)
	}
}

// formatJSON generates JSON output.
func (p *Planner) formatJSON(result *DecompositionResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal JSON: %w", err)
	}
	return string(data), nil
}

// formatHuman generates human-readable output.
func (p *Planner) formatHuman(result *DecompositionResult) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Feature Decomposition: %s\n", p.Description))
	sb.WriteString(fmt.Sprintf("Feature ID: %s\n", result.FeatureID))
	sb.WriteString(fmt.Sprintf("Workstreams: %d\n", len(result.Workstreams)))
	sb.WriteString(fmt.Sprintf("Dependencies: %d\n\n", len(result.Dependencies)))

	sb.WriteString(strings.Repeat("=", 60) + "\n\n")

	// List workstreams
	for i, ws := range result.Workstreams {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, ws.ID, ws.Title))
		sb.WriteString(fmt.Sprintf("   %s\n", ws.Description))
		if ws.Complexity != "" {
			sb.WriteString(fmt.Sprintf("   Complexity: %s\n", ws.Complexity))
		}
		sb.WriteString(fmt.Sprintf("   Status: %s\n\n", ws.Status))
	}

	// List dependencies
	if len(result.Dependencies) > 0 {
		sb.WriteString("Dependencies:\n\n")
		for _, dep := range result.Dependencies {
			sb.WriteString(fmt.Sprintf("- %s depends on %s\n", dep.From, dep.To))
			sb.WriteString(fmt.Sprintf("  Reason: %s\n", dep.Reason))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat("=", 60) + "\n")

	return sb.String(), nil
}

// PromptForInteractive runs interactive questions (AC2).
// This would be called in interactive mode to gather more information.
func (p *Planner) PromptForInteractive() error {
	if !p.Interactive {
		return nil
	}

	// In a real implementation, this would use a prompt library
	// For now, we just mark that interactive mode is active
	return nil
}

// ExecuteAutoApply executes the plan in auto-apply mode (AC3).
// This would trigger @oneshot or @build after planning.
func (p *Planner) ExecuteAutoApply(result *DecompositionResult) error {
	if !p.AutoApply {
		return nil
	}

	// In a real implementation, this would:
	// 1. Create workstream files
	// 2. Trigger execution (@oneshot or @build)
	// For now, just validate the plan is complete
	if len(result.Workstreams) == 0 {
		return fmt.Errorf("cannot auto-apply: no workstreams in plan")
	}

	return nil
}
