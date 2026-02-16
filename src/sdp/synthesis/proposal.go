package synthesis

import (
	"encoding/json"
	"time"
)

// Proposal represents a solution proposed by an agent
type Proposal struct {
	AgentID    string      `json:"agent_id"`
	Solution   interface{} `json:"solution"`
	Confidence float64     `json:"confidence"`
	Reasoning  string      `json:"reasoning"`
	Timestamp  time.Time   `json:"timestamp"`
}

// NewProposal creates a new proposal with current timestamp
func NewProposal(agentID string, solution interface{}, confidence float64, reasoning string) *Proposal {
	return &Proposal{
		AgentID:    agentID,
		Solution:   solution,
		Confidence: confidence,
		Reasoning:  reasoning,
		Timestamp:  time.Now().UTC(),
	}
}

// Equals checks if two proposals are equal (excluding timestamp)
func (p *Proposal) Equals(other *Proposal) bool {
	if p == nil || other == nil {
		return p == other
	}

	if p.AgentID != other.AgentID {
		return false
	}

	if p.Confidence != other.Confidence {
		return false
	}

	if p.Reasoning != other.Reasoning {
		return false
	}

	// Compare solutions as JSON
	pJSON, err1 := json.Marshal(p.Solution)
	otherJSON, err2 := json.Marshal(other.Solution)
	if err1 != nil || err2 != nil {
		return false
	}

	return string(pJSON) == string(otherJSON)
}

// IsHigherConfidenceThan checks if this proposal has higher confidence than another
func (p *Proposal) IsHigherConfidenceThan(other *Proposal) bool {
	if p == nil || other == nil {
		return false
	}
	return p.Confidence > other.Confidence
}

// MarshalJSON implements custom JSON marshaling
func (p *Proposal) MarshalJSON() ([]byte, error) {
	type Alias Proposal
	return json.Marshal(&struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Timestamp: p.Timestamp.Format(time.RFC3339Nano),
		Alias:     (*Alias)(p),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling
func (p *Proposal) UnmarshalJSON(data []byte) error {
	type Alias Proposal
	aux := &struct {
		Timestamp string `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Timestamp != "" {
		ts, err := time.Parse(time.RFC3339Nano, aux.Timestamp)
		if err != nil {
			return err
		}
		p.Timestamp = ts
	}

	return nil
}
