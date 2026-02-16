package sdpinit

import (
	"encoding/json"
	"fmt"
	"os"
)

// Exit codes for headless mode.
const (
	ExitSuccess          = 0
	ExitError            = 1
	ExitValidationFailed = 2
)

// HeadlessOutput is the JSON output format for headless mode.
type HeadlessOutput struct {
	Success     bool             `json:"success"`
	Error       string           `json:"error,omitempty"`
	ErrorCode   int              `json:"error_code,omitempty"`
	ProjectType string           `json:"project_type"`
	Created     []string         `json:"created"`
	Preflight   *PreflightResult `json:"preflight,omitempty"`
	Config      *ConfigSummary   `json:"config,omitempty"`
}

// ConfigSummary provides a summary of the configuration applied.
type ConfigSummary struct {
	Skills          []string `json:"skills"`
	EvidenceEnabled bool     `json:"evidence_enabled"`
	BeadsEnabled    bool     `json:"beads_enabled"`
}

// HeadlessRunner manages headless initialization.
type HeadlessRunner struct {
	config    Config
	preflight *PreflightResult
	output    *HeadlessOutput
}

// NewHeadlessRunner creates a new headless runner.
func NewHeadlessRunner(cfg Config) *HeadlessRunner {
	return &HeadlessRunner{
		config: cfg,
		output: &HeadlessOutput{
			Created: []string{},
		},
	}
}

// Run executes headless initialization.
func (h *HeadlessRunner) Run() (*HeadlessOutput, error) {
	// Run preflight checks
	h.preflight = RunPreflight()
	h.output.Preflight = h.preflight

	// Determine project type (config overrides detection)
	if h.config.ProjectType == "" {
		h.config.ProjectType = h.preflight.ProjectType
	}
	// Set output to final project type (may differ from detected)
	h.output.ProjectType = h.config.ProjectType

	// Validate configuration
	if err := h.validate(); err != nil {
		h.output.Success = false
		h.output.Error = err.Error()
		h.output.ErrorCode = ExitValidationFailed
		return h.output, err
	}

	// Get defaults
	defaults := MergeDefaults(h.config.ProjectType, &h.config)

	// Track what would be created
	h.trackCreatedFiles()

	// Set config summary
	h.output.Config = &ConfigSummary{
		Skills:          defaults.Skills,
		EvidenceEnabled: defaults.EvidenceEnabled,
		BeadsEnabled:    !h.config.SkipBeads,
	}

	// If dry-run, return without creating
	if h.config.DryRun {
		h.output.Success = true
		return h.output, nil
	}

	// Run actual initialization
	if err := Run(h.config); err != nil {
		h.output.Success = false
		h.output.Error = err.Error()
		h.output.ErrorCode = ExitError
		return h.output, err
	}

	h.output.Success = true
	return h.output, nil
}

func (h *HeadlessRunner) validate() error {
	// Validate project type
	validTypes := map[string]bool{
		"go": true, "node": true, "python": true,
		"mixed": true, "unknown": true,
	}

	if !validTypes[h.config.ProjectType] {
		return fmt.Errorf("invalid project type: %s", h.config.ProjectType)
	}

	// Check for critical conflicts
	for _, conflict := range h.preflight.Conflicts {
		if conflict == ".claude/settings.json" && !h.config.Force {
			return fmt.Errorf("conflict: %s exists (use --force to overwrite)", conflict)
		}
	}

	return nil
}

func (h *HeadlessRunner) trackCreatedFiles() {
	h.output.Created = []string{
		".claude/",
		".claude/skills/",
		".claude/agents/",
		".claude/validators/",
		".claude/settings.json",
	}
}

// RunHeadless is a convenience function for headless initialization.
func RunHeadless(cfg Config) (*HeadlessOutput, error) {
	runner := NewHeadlessRunner(cfg)
	return runner.Run()
}

// OutputJSON writes the headless output as JSON to stdout.
func (o *HeadlessOutput) OutputJSON() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(o)
}

// GetExitCode returns the appropriate exit code for the output.
func (o *HeadlessOutput) GetExitCode() int {
	if o.Success {
		return ExitSuccess
	}
	if o.ErrorCode == 0 {
		return ExitError
	}
	return o.ErrorCode
}
