package planner

import (
	"fmt"
	"strings"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

// Planner decomposes feature descriptions into workstreams.
type Planner struct {
	BacklogDir     string           // Directory for workstream files
	Description    string           // Feature description to decompose
	Interactive    bool             // AC2: Interactive mode (ask questions)
	AutoApply      bool             // AC3: Auto-apply mode (execute after plan)
	DryRun         bool             // AC7: Show what would be created
	OutputFormat   string           // AC6: Output format (human, json)
	ModelAPI       string           // AC8: AI model API endpoint
	Questions      []string         // Questions for interactive mode
	EvidenceWriter *evidence.Writer // Evidence log writer
}

// DecompositionResult contains the output of feature decomposition.
type DecompositionResult struct {
	Workstreams  []Workstream `json:"workstreams"`
	Dependencies []Dependency `json:"dependencies"`
	Summary      string       `json:"summary"`
	FeatureID    string       `json:"feature_id,omitempty"`
	CreatedAt    string       `json:"created_at"`
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
