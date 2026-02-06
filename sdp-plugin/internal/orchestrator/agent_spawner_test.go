package orchestrator

import (
	"errors"
	"testing"
	"time"
)

// MockTaskTool mocks the Task tool interface
type MockTaskTool struct {
	spawnedAgents map[string]string // agentID -> agentType
	agentResults  map[string]string // agentID -> result
	shouldFail    bool
}

func (m *MockTaskTool) Spawn(agentType string, prompt string) (string, error) {
	if m.shouldFail {
		return "", ErrAgentSpawnFailed
	}

	if m.spawnedAgents == nil {
		m.spawnedAgents = make(map[string]string)
	}
	if m.agentResults == nil {
		m.agentResults = make(map[string]string)
	}

	// Generate fake agent ID
	agentID := agentType + "-" + time.Now().Format("20060102-150405")
	m.spawnedAgents[agentID] = agentType

	return agentID, nil
}

func (m *MockTaskTool) GetResult(agentID string) (string, error) {
	if m.shouldFail {
		return "", ErrAgentNotFound
	}

	if result, ok := m.agentResults[agentID]; ok {
		return result, nil
	}

	// Return fake result
	return `{"status": "complete", "output": "Task completed successfully"}`, nil
}

func (m *MockTaskTool) Terminate(agentID string) error {
	if m.shouldFail {
		return ErrAgentNotFound
	}

	delete(m.spawnedAgents, agentID)
	return nil
}

func (m *MockTaskTool) SetResult(agentID string, result string) {
	if m.agentResults == nil {
		m.agentResults = make(map[string]string)
	}
	m.agentResults[agentID] = result
}

func (m *MockTaskTool) GetSpawnedAgentType(agentID string) string {
	return m.spawnedAgents[agentID]
}

func TestAgentSpawner_NewAgentSpawner(t *testing.T) {
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	if spawner == nil {
		t.Fatal("Expected non-nil spawner")
	}

	if spawner.taskTool == nil {
		t.Error("Expected taskTool to be initialized")
	}
}

func TestAgentSpawner_SpawnPlanner(t *testing.T) {
	// AC1: Can spawn planner agent for architecture
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	agentID, err := spawner.Spawn("planner", "Design architecture for user auth")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify agent ID format
	if agentID == "" {
		t.Fatal("Expected non-empty agent ID")
	}

	// Verify agent type
	if tool.GetSpawnedAgentType(agentID) != "planner" {
		t.Errorf("Expected agent type 'planner', got '%s'", tool.GetSpawnedAgentType(agentID))
	}
}

func TestAgentSpawner_SpawnBuilder(t *testing.T) {
	// AC2: Can spawn builder agent for implementation
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	agentID, err := spawner.Spawn("builder", "Execute WS-001")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify agent ID
	if agentID == "" {
		t.Fatal("Expected non-empty agent ID")
	}

	// Verify agent type
	if tool.GetSpawnedAgentType(agentID) != "builder" {
		t.Errorf("Expected agent type 'builder', got '%s'", tool.GetSpawnedAgentType(agentID))
	}
}

func TestAgentSpawner_SpawnReviewer(t *testing.T) {
	// AC3: Can spawn reviewer agent for quality checks
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	agentID, err := spawner.Spawn("reviewer", "Review F001 for quality")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify agent type
	if tool.GetSpawnedAgentType(agentID) != "reviewer" {
		t.Errorf("Expected agent type 'reviewer', got '%s'", tool.GetSpawnedAgentType(agentID))
	}
}

func TestAgentSpawner_GetResult(t *testing.T) {
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	// Spawn agent
	agentID, _ := spawner.Spawn("builder", "Execute WS-001")

	// Set mock result
	tool.SetResult(agentID, `{"status": "complete", "output": "WS-001 executed successfully"}`)

	// Get result
	result, err := spawner.GetResult(agentID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify result
	if result.Status != "complete" {
		t.Errorf("Expected status 'complete', got '%s'", result.Status)
	}
}

func TestAgentSpawner_GetResult_ParseError(t *testing.T) {
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	// Spawn agent
	agentID, _ := spawner.Spawn("builder", "Execute WS-001")

	// Set invalid JSON result
	tool.SetResult(agentID, "invalid json")

	// Get result - should not fail, just return error in result
	result, err := spawner.GetResult(agentID)

	// Should return result with error status
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", result.Status)
	}
}

func TestAgentSpawner_Terminate(t *testing.T) {
	// AC4: Agent terminated after completion
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	// Spawn agent
	agentID, _ := spawner.Spawn("builder", "Execute WS-001")

	// Terminate agent
	err := spawner.Terminate(agentID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify agent was terminated
	if tool.GetSpawnedAgentType(agentID) != "" {
		t.Error("Expected agent to be terminated")
	}
}

func TestAgentSpawner_SpawnFailure(t *testing.T) {
	tool := &MockTaskTool{shouldFail: true}
	spawner := NewAgentSpawner(tool)

	_, err := spawner.Spawn("builder", "Execute WS-001")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, ErrAgentSpawnFailed) {
		t.Errorf("Expected ErrAgentSpawnFailed, got %v", err)
	}
}

func TestAgentSpawner_AgentResultFromJSON(t *testing.T) {
	json := `{"status": "complete", "output": "Task completed", "error": ""}`

	result, err := AgentResultFromJSON(json)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Status != "complete" {
		t.Errorf("Expected status 'complete', got '%s'", result.Status)
	}

	if result.Output != "Task completed" {
		t.Errorf("Expected output 'Task completed', got '%s'", result.Output)
	}
}

func TestAgentSpawner_SpawnMultipleAgents(t *testing.T) {
	// Test spawning multiple agents for parallel execution
	tool := &MockTaskTool{}
	spawner := NewAgentSpawner(tool)

	// Spawn multiple agents
	agentIDs := []string{}
	for i := 0; i < 3; i++ {
		agentID, err := spawner.Spawn("builder", "Execute WS-00"+string(rune('1'+i)))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		agentIDs = append(agentIDs, agentID)
	}

	// Verify all agents spawned
	if len(agentIDs) != 3 {
		t.Fatalf("Expected 3 agents, got %d", len(agentIDs))
	}

	// Terminate all agents
	for _, agentID := range agentIDs {
		err := spawner.Terminate(agentID)
		if err != nil {
			t.Errorf("Expected no error terminating %s, got %v", agentID, err)
		}
	}
}
