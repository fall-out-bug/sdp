package synthesis

import (
	"errors"
	"testing"
	"time"
)

// MockAgent implements Agent interface for testing
type MockAgent struct {
	id        string
	available bool
	proposal  *Proposal
	err       error
}

func (m *MockAgent) ID() string {
	return m.id
}

func (m *MockAgent) Available() bool {
	return m.available
}

func (m *MockAgent) Consult(task Task, timeout time.Duration) (*Proposal, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.proposal, nil
}

// TestNewSupervisor verifies supervisor creation
func TestNewSupervisor(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	if s == nil {
		t.Fatal("NewSupervisor returned nil")
	}
	if s.agents == nil {
		t.Error("agents map not initialized")
	}
	if s.timeout != 30*time.Second {
		t.Error("default timeout not set")
	}
}

// TestSupervisor_RegisterAgent verifies agent registration
func TestSupervisor_RegisterAgent(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	agent := &MockAgent{id: "agent-1", available: true}
	s.RegisterAgent(agent)

	if len(s.agents) != 1 {
		t.Error("agent not registered")
	}
}

// TestSupervisor_RegisterAgent_Nil verifies nil handling
func TestSupervisor_RegisterAgent_Nil(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(nil)
	if len(s.agents) != 0 {
		t.Error("nil agent should not be registered")
	}
}

// TestSupervisor_SetTimeout verifies timeout setting
func TestSupervisor_SetTimeout(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	s.SetTimeout(10 * time.Second)
	if s.timeout != 10*time.Second {
		t.Error("timeout not set")
	}
}

// TestSupervisor_ConsultAgents_Success verifies successful consultation
func TestSupervisor_ConsultAgents_Success(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(&MockAgent{
		id:        "agent-1",
		available: true,
		proposal:  NewProposal("agent-1", "solution", 0.9, ""),
	})
	s.RegisterAgent(&MockAgent{
		id:        "agent-2",
		available: true,
		proposal:  NewProposal("agent-2", "solution", 0.8, ""),
	})

	proposals, err := s.ConsultAgents("test-task")
	if err != nil {
		t.Fatalf("ConsultAgents failed: %v", err)
	}
	if len(proposals) != 2 {
		t.Errorf("expected 2 proposals, got %d", len(proposals))
	}
}

// TestSupervisor_ConsultAgents_SkipsUnavailable verifies unavailable agents skipped
func TestSupervisor_ConsultAgents_SkipsUnavailable(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(&MockAgent{
		id:        "agent-1",
		available: false,
	})
	s.RegisterAgent(&MockAgent{
		id:        "agent-2",
		available: true,
		proposal:  NewProposal("agent-2", "solution", 0.8, ""),
	})

	proposals, err := s.ConsultAgents("test-task")
	if err != nil {
		t.Fatalf("ConsultAgents failed: %v", err)
	}
	if len(proposals) != 1 {
		t.Errorf("expected 1 proposal (unavailable skipped), got %d", len(proposals))
	}
}

// TestSupervisor_ConsultAgents_SkipsErrors verifies error handling
func TestSupervisor_ConsultAgents_SkipsErrors(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(&MockAgent{
		id:        "agent-1",
		available: true,
		err:       errors.New("consultation error"),
	})
	s.RegisterAgent(&MockAgent{
		id:        "agent-2",
		available: true,
		proposal:  NewProposal("agent-2", "solution", 0.8, ""),
	})

	proposals, err := s.ConsultAgents("test-task")
	if err != nil {
		t.Fatalf("ConsultAgents failed: %v", err)
	}
	if len(proposals) != 1 {
		t.Errorf("expected 1 proposal (error skipped), got %d", len(proposals))
	}
}

// TestSupervisor_ConsultAgents_NoProposals verifies no proposals error
func TestSupervisor_ConsultAgents_NoProposals(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	// No agents registered
	_, err := s.ConsultAgents("test-task")
	if err == nil {
		t.Error("expected error for no proposals")
	}
	if !errors.Is(err, ErrNoProposals) {
		t.Errorf("expected ErrNoProposals, got %v", err)
	}
}

// TestSupervisor_ConsultAgents_NilProposal verifies nil proposal handling
func TestSupervisor_ConsultAgents_NilProposal(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(&MockAgent{
		id:        "agent-1",
		available: true,
		proposal:  nil, // Returns nil
	})

	_, err := s.ConsultAgents("test-task")
	if err == nil {
		t.Error("expected error when all agents return nil")
	}
}

// TestSupervisor_MakeDecision_Approved verifies approved decision
func TestSupervisor_MakeDecision_Approved(t *testing.T) {
	engine := DefaultRuleEngine()
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(&MockAgent{
		id:        "agent-1",
		available: true,
		proposal:  NewProposal("agent-1", "solution", 0.9, ""),
	})
	s.RegisterAgent(&MockAgent{
		id:        "agent-2",
		available: true,
		proposal:  NewProposal("agent-2", "solution", 0.8, ""),
	})

	decision, err := s.MakeDecision("test-task")
	if err != nil {
		t.Fatalf("MakeDecision failed: %v", err)
	}
	if decision.Status != "approved" {
		t.Errorf("expected approved status, got %s", decision.Status)
	}
}

// TestSupervisor_MakeDecision_Escalated verifies escalated decision
func TestSupervisor_MakeDecision_Escalated(t *testing.T) {
	engine := NewRuleEngine() // Empty engine - will always fail
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(&MockAgent{
		id:        "agent-1",
		available: true,
		proposal:  NewProposal("agent-1", "solution", 0.9, ""),
	})

	decision, err := s.MakeDecision("test-task")
	if err != nil {
		t.Fatalf("MakeDecision failed: %v", err)
	}
	if decision.Status != "escalated" {
		t.Errorf("expected escalated status, got %s", decision.Status)
	}
}

// TestSupervisor_MakeDecision_NoAgents verifies no agents error
func TestSupervisor_MakeDecision_NoAgents(t *testing.T) {
	engine := DefaultRuleEngine()
	s := NewSupervisor(engine, 5)

	_, err := s.MakeDecision("test-task")
	if err == nil {
		t.Error("expected error for no agents")
	}
}

// TestSupervisor_GetAgentStatus verifies agent status
func TestSupervisor_GetAgentStatus(t *testing.T) {
	engine := NewRuleEngine()
	s := NewSupervisor(engine, 5)

	s.RegisterAgent(&MockAgent{id: "agent-1", available: true})
	s.RegisterAgent(&MockAgent{id: "agent-2", available: false})

	status := s.GetAgentStatus()
	if status["agent-1"] != true {
		t.Error("agent-1 should be available")
	}
	if status["agent-2"] != false {
		t.Error("agent-2 should be unavailable")
	}
}

// TestDecision_Structure verifies decision structure
func TestDecision_Structure(t *testing.T) {
	decision := &Decision{
		Status:    "approved",
		Solution:  "test-solution",
		Rule:      "unanimous",
		Proposals: []*Proposal{{AgentID: "agent-1"}},
		Reason:    "test reason",
	}

	if decision.Status != "approved" {
		t.Error("Status not set")
	}
	if decision.Solution != "test-solution" {
		t.Error("Solution not set")
	}
	if decision.Rule != "unanimous" {
		t.Error("Rule not set")
	}
}

// TestTask_Interface verifies Task interface
func TestTask_Interface(t *testing.T) {
	var task Task = "string task"
	if task == nil {
		t.Error("string should be valid Task")
	}

	task = map[string]interface{}{"key": "value"}
	if task == nil {
		t.Error("map should be valid Task")
	}
}
