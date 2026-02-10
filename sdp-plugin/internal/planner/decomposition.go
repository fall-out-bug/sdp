package planner

import (
	"fmt"
	"time"
)

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
