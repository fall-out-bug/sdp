// Package nextstep provides deterministic next-step recommendations for SDP workflows.
// It defines the contract for "what should happen next" based on project and execution state.
package nextstep

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ContractVersion is the current version of the recommendation contract.
// This should be incremented when breaking changes are made to the schema.
const ContractVersion = "1.0.0"

// Common errors for recommendation validation.
var (
	ErrEmptyCommand    = errors.New("command cannot be empty")
	ErrEmptyReason     = errors.New("reason cannot be empty")
	ErrInvalidCategory = errors.New("invalid or empty category")
	ErrConfidenceRange = errors.New("confidence must be between 0 and 1")
)

// RecommendationCategory defines the type of recommendation being made.
type RecommendationCategory string

const (
	// CategoryExecution recommends executing a workstream or command.
	CategoryExecution RecommendationCategory = "execution"
	// CategoryRecovery recommends recovering from an error or failure.
	CategoryRecovery RecommendationCategory = "recovery"
	// CategoryPlanning recommends planning activities (design, idea).
	CategoryPlanning RecommendationCategory = "planning"
	// CategoryInformation recommends informational commands (status, help).
	CategoryInformation RecommendationCategory = "information"
	// CategorySetup recommends setup or initialization activities.
	CategorySetup RecommendationCategory = "setup"
)

// Recommendation represents a next-step recommendation output.
// This is the primary contract for all consumer surfaces (CLI, agents, docs).
type Recommendation struct {
	ActionID string `json:"action_id,omitempty"`
	// Command is the recommended SDP command to execute.
	Command string `json:"command"`
	// Reason explains why this command is recommended.
	Reason string `json:"reason"`
	// Confidence indicates how strongly this recommendation is made (0.0-1.0).
	Confidence float64 `json:"confidence"`
	// Category classifies the type of recommendation.
	Category RecommendationCategory `json:"category"`
	// Version is the contract version for this recommendation.
	Version string `json:"version"`
	// Alternatives provides fallback options if the primary is not suitable.
	Alternatives         []Alternative  `json:"alternatives,omitempty"`
	RequiredContext      map[string]any `json:"required_context,omitempty"`
	OptionalContext      map[string]any `json:"optional_context,omitempty"`
	PolicyExpectations   []string       `json:"policy_expectations,omitempty"`
	EvidenceExpectations []string       `json:"evidence_expectations,omitempty"`
	// Metadata contains additional context-specific information.
	Metadata map[string]any `json:"metadata,omitempty"`
}

type InstructionPayload = Recommendation

// Alternative represents a secondary recommendation option.
type Alternative struct {
	Command string `json:"command"`
	Reason  string `json:"reason"`
}

// Validate checks if the recommendation is valid.
func (r *Recommendation) Validate() error {
	if r.Command == "" {
		return ErrEmptyCommand
	}
	if r.Reason == "" {
		return ErrEmptyReason
	}
	if r.Category == "" {
		return ErrInvalidCategory
	}
	if r.Confidence < 0 || r.Confidence > 1 {
		return ErrConfidenceRange
	}
	return nil
}

// String returns a human-readable representation of the recommendation.
func (r *Recommendation) String() string {
	return fmt.Sprintf("[%s] %s (%.0f%% confidence): %s",
		r.Category, r.Command, r.Confidence*100, r.Reason)
}

func (r *Recommendation) ToJSON() ([]byte, error) {
	r.enrich()
	return json.Marshal(r)
}

func FromJSON(data []byte) (*Recommendation, error) {
	var rec Recommendation
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *Recommendation) enrich() {
	if r.ActionID == "" {
		r.ActionID = actionIDFromCommand(r.Command)
	}
	if len(r.RequiredContext) == 0 {
		r.RequiredContext = inferRequiredContext(r.Command, r.Metadata)
	}
	if len(r.OptionalContext) == 0 {
		r.OptionalContext = inferOptionalContext(r.Metadata, r.RequiredContext)
	}
	if len(r.PolicyExpectations) == 0 {
		r.PolicyExpectations = inferPolicyExpectations(r.Command, r.Category)
	}
	if len(r.EvidenceExpectations) == 0 {
		r.EvidenceExpectations = inferEvidenceExpectations(r.Command, r.Category)
	}
}

func actionIDFromCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	if strings.HasPrefix(parts[0], "@") {
		return "skill." + strings.TrimPrefix(parts[0], "@")
	}
	if len(parts) > 1 {
		return parts[0] + "." + parts[1]
	}
	return parts[0]
}

func inferRequiredContext(command string, metadata map[string]any) map[string]any {
	if len(metadata) == 0 {
		return nil
	}

	required := map[string]any{}
	if strings.Contains(command, "--ws") {
		for _, key := range []string{"workstream_id", "failed_workstream", "active_workstream", "blocker"} {
			if value, ok := metadata[key]; ok {
				required[key] = value
			}
		}
	}
	if strings.Contains(command, "sdp review") || strings.Contains(command, "sdp deploy") {
		if value, ok := metadata["feature_id"]; ok {
			required["feature_id"] = value
		}
	}
	if len(required) == 0 {
		return nil
	}
	return required
}

func inferOptionalContext(metadata map[string]any, required map[string]any) map[string]any {
	if len(metadata) == 0 {
		return nil
	}
	optional := map[string]any{}
	for key, value := range metadata {
		if _, ok := required[key]; ok {
			continue
		}
		optional[key] = value
	}
	if len(optional) == 0 {
		return nil
	}
	return optional
}

func inferPolicyExpectations(command string, category RecommendationCategory) []string {
	switch {
	case strings.HasPrefix(command, "sdp apply"), strings.HasPrefix(command, "sdp build"):
		return []string{"respect scope and workstream boundaries"}
	case strings.HasPrefix(command, "sdp review"):
		return []string{"review must pass before deployment"}
	case strings.HasPrefix(command, "sdp deploy"):
		return []string{"deploy only after explicit review approval"}
	case category == CategoryRecovery:
		return []string{"confirm failure context before retrying"}
	default:
		return nil
	}
}

func inferEvidenceExpectations(command string, category RecommendationCategory) []string {
	switch {
	case strings.HasPrefix(command, "sdp apply"), strings.HasPrefix(command, "sdp build"):
		return []string{"record execution evidence for the targeted workstream"}
	case strings.HasPrefix(command, "sdp review"):
		return []string{"record review verdict before deployment"}
	case strings.HasPrefix(command, "sdp deploy"):
		return []string{"record deployment approval in the evidence log"}
	case category == CategoryRecovery:
		return []string{"capture failure context before resuming execution"}
	default:
		return nil
	}
}
