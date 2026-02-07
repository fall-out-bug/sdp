package synthesis

import (
	"errors"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/src/sdp/synthesis"
)

// MockAgent for testing
type mockAgent struct {
	id         string
	available  bool
	proposal   *synthesis.Proposal
	shouldFail bool
}

func (m *mockAgent) ID() string { return m.id }
func (m *mockAgent) Available() bool { return m.available }
func (m *mockAgent) Consult(task synthesis.Task, timeout time.Duration) (*synthesis.Proposal, error) {
	if m.shouldFail {
		return nil, errors.New("agent failed")
	}
	return m.proposal, nil
}

// TestSupervisor_New_CreatesSupervisor tests AC1: Supervisor creation
func TestSupervisor_New_CreatesSupervisor(t *testing.T) {
	// Arrange
	engine := synthesis.DefaultRuleEngine()

	// Act
	supervisor := synthesis.NewSupervisor(engine, 5)

	// Assert
	if supervisor == nil {
		t.Fatal("Expected supervisor to be created, got nil")
	}
}

// TestSupervisor_ConsultAgents_GathersProposals tests AC2: Agent coordination
func TestSupervisor_ConsultAgents_GathersProposals(t *testing.T) {
	// Arrange
	engine := synthesis.DefaultRuleEngine()
	supervisor := synthesis.NewSupervisor(engine, 5)

	agent1 := &mockAgent{
		id:        "agent-1",
		available: true,
		proposal:  synthesis.NewProposal("agent-1", map[string]string{"approach": "TDD"}, 0.9, ""),
	}
	agent2 := &mockAgent{
		id:        "agent-2",
		available: true,
		proposal:  synthesis.NewProposal("agent-2", map[string]string{"approach": "TDD"}, 0.85, ""),
	}

	supervisor.RegisterAgent(agent1)
	supervisor.RegisterAgent(agent2)

	task := map[string]interface{}{"description": "Implement feature X"}

	// Act
	proposals, err := supervisor.ConsultAgents(task)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(proposals) != 2 {
		t.Errorf("Expected 2 proposals, got %d", len(proposals))
	}
}

// TestSupervisor_ConsultAgents_SkipsUnavailableAgents tests AC2: Skip unavailable agents
func TestSupervisor_ConsultAgents_SkipsUnavailableAgents(t *testing.T) {
	// Arrange
	engine := synthesis.DefaultRuleEngine()
	supervisor := synthesis.NewSupervisor(engine, 5)

	agent1 := &mockAgent{
		id:        "agent-1",
		available: true,
		proposal:  synthesis.NewProposal("agent-1", nil, 0.9, ""),
	}
	agent2 := &mockAgent{
		id:        "agent-2",
		available: false, // Unavailable
		proposal:  nil,
	}

	supervisor.RegisterAgent(agent1)
	supervisor.RegisterAgent(agent2)

	task := map[string]interface{}{"description": "Test task"}

	// Act
	proposals, err := supervisor.ConsultAgents(task)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(proposals) != 1 {
		t.Errorf("Expected 1 proposal (agent-2 unavailable), got %d", len(proposals))
	}
}

// TestSupervisor_MakeDecision_SynthesizesProposals tests AC3: Synthesizer integration
func TestSupervisor_MakeDecision_SynthesizesProposals(t *testing.T) {
	// Arrange
	engine := synthesis.DefaultRuleEngine()
	supervisor := synthesis.NewSupervisor(engine, 5)

	solution := map[string]string{"approach": "TDD"}
	agent1 := &mockAgent{
		id:        "agent-1",
		available: true,
		proposal:  synthesis.NewProposal("agent-1", solution, 0.9, "TDD is best"),
	}
	agent2 := &mockAgent{
		id:        "agent-2",
		available: true,
		proposal:  synthesis.NewProposal("agent-2", solution, 0.85, "Use TDD"),
	}

	supervisor.RegisterAgent(agent1)
	supervisor.RegisterAgent(agent2)

	task := map[string]interface{}{"description": "Test task"}

	// Act
	decision, err := supervisor.MakeDecision(task)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if decision.Status != "approved" {
		t.Errorf("Expected status 'approved', got %s", decision.Status)
	}

	if decision.Solution == nil {
		t.Error("Expected solution to be set")
	}
}
