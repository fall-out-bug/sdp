package synthesis

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// ErrNoProposals is returned when no agents could provide proposals
var ErrNoProposals = errors.New("no proposals received from agents")

// Agent represents a specialist agent that can provide proposals
type Agent interface {
	// ID returns the agent identifier
	ID() string

	// Available returns true if the agent is available for consultation
	Available() bool

	// Consult requests a proposal from the agent for the given task
	Consult(task Task, timeout time.Duration) (*Proposal, error)
}

// Task represents a task to be solved by agents
type Task interface{}

// Decision represents the result of supervisor's decision making
type Decision struct {
	Status    string      // "approved" or "escalated"
	Solution  interface{} // The chosen solution
	Rule      string      // The synthesis rule that was applied
	Proposals []*Proposal // All proposals received
	Reason    string      // Reason for escalation (if applicable)
}

// Supervisor coordinates specialist agents and applies synthesizer
type Supervisor struct {
	engine    *RuleEngine
	agents    map[string]Agent
	maxAgents int
	timeout   time.Duration
}

// NewSupervisor creates a new supervisor
func NewSupervisor(engine *RuleEngine, maxAgents int) *Supervisor {
	return &Supervisor{
		engine:    engine,
		agents:    make(map[string]Agent),
		maxAgents: maxAgents,
		timeout:   30 * time.Second, // Default timeout
	}
}

// RegisterAgent adds an agent to the supervisor
func (s *Supervisor) RegisterAgent(agent Agent) {
	if agent == nil {
		return
	}
	s.agents[agent.ID()] = agent
}

// SetTimeout sets the consultation timeout
func (s *Supervisor) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// ConsultAgents gathers proposals from available specialist agents
func (s *Supervisor) ConsultAgents(task Task) ([]*Proposal, error) {
	proposals := make([]*Proposal, 0)

	for _, agent := range s.agents {
		if !agent.Available() {
			log.Printf("[Supervisor] Agent %s is not available, skipping", agent.ID())
			continue
		}

		// Consult agent with timeout
		proposal, err := agent.Consult(task, s.timeout)
		if err != nil {
			log.Printf("[Supervisor] Agent %s consultation failed: %v", agent.ID(), err)
			continue
		}

		if proposal != nil {
			proposals = append(proposals, proposal)
			log.Printf("[Supervisor] Received proposal from %s (confidence: %.2f)", agent.ID(), proposal.Confidence)
		}
	}

	if len(proposals) == 0 {
		return nil, ErrNoProposals
	}

	log.Printf("[Supervisor] Collected %d proposals from agents", len(proposals))
	return proposals, nil
}

// MakeDecision consults agents and synthesizes proposals into a decision
func (s *Supervisor) MakeDecision(task Task) (*Decision, error) {
	// Step 1: Consult agents
	proposals, err := s.ConsultAgents(task)
	if err != nil {
		return nil, fmt.Errorf("failed to consult agents: %w", err)
	}

	// Step 2: Synthesize proposals
	result, err := s.engine.Execute(proposals)
	if err != nil {
		// Escalate to human
		log.Printf("[Supervisor] Synthesis failed, escalating to human: %v", err)
		return &Decision{
			Status:    "escalated",
			Proposals: proposals,
			Reason:    fmt.Sprintf("Synthesis failed: %v", err),
		}, nil
	}

	// Return approved decision
	log.Printf("[Supervisor] Decision approved using rule '%s'", result.Rule)
	return &Decision{
		Status:    "approved",
		Solution:  result.Solution,
		Rule:      result.Rule,
		Proposals: proposals,
	}, nil
}

// GetAgentStatus returns status of all registered agents
func (s *Supervisor) GetAgentStatus() map[string]bool {
	status := make(map[string]bool)
	for id, agent := range s.agents {
		status[id] = agent.Available()
	}
	return status
}
