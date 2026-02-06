package decision

import "time"

// Decision represents an architectural or product decision
type Decision struct {
	// Timestamp of when decision was made
	Timestamp time.Time `json:"timestamp"`

	// Type of decision
	Type string `json:"type"` // "vision", "technical", "tradeoff", "explicit"

	// Feature ID associated with this decision
	FeatureID string `json:"feature_id"`

	// Workstream ID if linked to specific workstream
	WorkstreamID string `json:"ws_id,omitempty"`

	// Question or problem that prompted the decision
	Question string `json:"question"`

	// Decision made (what was chosen)
	Decision string `json:"decision"`

	// Rationale behind the decision
	Rationale string `json:"rationale"`

	// Alternatives considered
	Alternatives []string `json:"alternatives,omitempty"`

	// Outcome or impact
	Outcome string `json:"outcome"`

	// Person/AI agent who made the decision
	DecisionMaker string `json:"decision_maker"` // "user", "claude", "system"

	// Tags for categorization
	Tags []string `json:"tags,omitempty"`
}

// DecisionLog is the append-only log of decisions
type DecisionLog struct {
	Decisions []Decision `json:"decisions"`
}

// DecisionType constants
const (
	DecisionTypeVision    = "vision"
	DecisionTypeTechnical = "technical"
	DecisionTypeTradeoff   = "tradeoff"
	DecisionTypeExplicit  = "explicit"
)
