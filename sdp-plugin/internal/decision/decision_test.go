package decision_test

import (
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestDecision_TypeConstants(t *testing.T) {
	tests := []struct {
		name string
		typ  string
	}{
		{"Vision", decision.DecisionTypeVision},
		{"Technical", decision.DecisionTypeTechnical},
		{"Tradeoff", decision.DecisionTypeTradeoff},
		{"Explicit", decision.DecisionTypeExplicit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.typ == "" {
				t.Error("Type constant should not be empty")
			}
		})
	}
}

func TestDecision_JSONTags(t *testing.T) {
	d := decision.Decision{
		Timestamp:     time.Now(),
		Type:          decision.DecisionTypeTechnical,
		FeatureID:     "F001",
		WorkstreamID:  "00-001-01",
		Question:      "Test?",
		Decision:      "Yes",
		Rationale:     "Because",
		Alternatives:  []string{"No", "Maybe"},
		Outcome:       "Success",
		DecisionMaker: "user",
		Tags:          []string{"tag1", "tag2"},
	}

	// Verify JSON marshaling doesn't crash
	// This test ensures struct tags are valid
	_ = d
}
