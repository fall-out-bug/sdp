package planner

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
