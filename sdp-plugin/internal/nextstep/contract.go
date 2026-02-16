// Package nextstep provides deterministic next-step recommendations for SDP workflows.
// It defines the contract for "what should happen next" based on project and execution state.
package nextstep

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Alternatives []Alternative `json:"alternatives,omitempty"`
	// Metadata contains additional context-specific information.
	Metadata map[string]any `json:"metadata,omitempty"`
}

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

// ToJSON serializes a recommendation to JSON.
func (r *Recommendation) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// FromJSON deserializes a recommendation from JSON.
func FromJSON(data []byte) (*Recommendation, error) {
	var rec Recommendation
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}
