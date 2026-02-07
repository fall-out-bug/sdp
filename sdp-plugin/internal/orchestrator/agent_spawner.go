package orchestrator

import (
	"encoding/json"
	"fmt"
	"time"
)

// TaskTool defines the interface for spawning agents via Task tool
type TaskTool interface {
	Spawn(agentType string, prompt string) (string, error)
	GetResult(agentID string) (string, error)
	Terminate(agentID string) error
}

// AgentResult represents the result from an agent
type AgentResult struct {
	Status string `json:"status"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// AgentSpawner handles spawning and managing specialist agents
type AgentSpawner struct {
	taskTool TaskTool
}

// NewAgentSpawner creates a new agent spawner
func NewAgentSpawner(tool TaskTool) *AgentSpawner {
	return &AgentSpawner{
		taskTool: tool,
	}
}

// Spawn spawns a new agent of the specified type
func (as *AgentSpawner) Spawn(agentType string, prompt string) (string, error) {
	// Validate agent type
	validTypes := map[string]bool{
		"planner":  true,
		"builder":  true,
		"reviewer": true,
		"tester":   true,
		"security": true,
	}

	if !validTypes[agentType] {
		return "", fmt.Errorf("invalid agent type: %s (valid: planner, builder, reviewer, tester, security)", agentType)
	}

	// Spawn agent via Task tool
	agentID, err := as.taskTool.Spawn(agentType, prompt)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAgentSpawnFailed, err)
	}

	return agentID, nil
}

// GetResult retrieves the result from an agent
func (as *AgentSpawner) GetResult(agentID string) (AgentResult, error) {
	resultJSON, err := as.taskTool.GetResult(agentID)
	if err != nil {
		return AgentResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	// Parse result
	result, err := AgentResultFromJSON(resultJSON)
	if err != nil {
		// If parsing fails, return error result
		return AgentResult{
			Status: "error",
			Output: "",
			Error:  fmt.Sprintf("Failed to parse agent result: %v", err),
		}, nil
	}

	return result, nil
}

// Terminate terminates an agent
func (as *AgentSpawner) Terminate(agentID string) error {
	err := as.taskTool.Terminate(agentID)
	if err != nil {
		return fmt.Errorf("failed to terminate agent: %w", err)
	}
	return nil
}

// SpawnAndWait spawns an agent and waits for completion
func (as *AgentSpawner) SpawnAndWait(agentType string, prompt string, timeout time.Duration) (AgentResult, error) {
	// Spawn agent
	agentID, err := as.Spawn(agentType, prompt)
	if err != nil {
		return AgentResult{}, err
	}

	// Wait for result with timeout
	resultChan := make(chan AgentResult, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := as.GetResult(agentID)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		// Terminate agent after getting result
		_ = as.Terminate(agentID)
		return result, nil
	case err := <-errChan:
		// Terminate agent on error
		_ = as.Terminate(agentID)
		return AgentResult{}, err
	case <-time.After(timeout):
		// Terminate agent on timeout
		_ = as.Terminate(agentID)
		return AgentResult{}, fmt.Errorf("agent timeout after %v", timeout)
	}
}

// SpawnMultiple spawns multiple agents in parallel
func (as *AgentSpawner) SpawnMultiple(agentConfigs []AgentConfig) ([]string, error) {
	agentIDs := make([]string, 0, len(agentConfigs))

	for _, config := range agentConfigs {
		agentID, err := as.Spawn(config.Type, config.Prompt)
		if err != nil {
			// Terminate all spawned agents on failure
			for _, id := range agentIDs {
				_ = as.Terminate(id)
			}
			return nil, fmt.Errorf("failed to spawn agent %s: %w", config.Type, err)
		}
		agentIDs = append(agentIDs, agentID)
	}

	return agentIDs, nil
}

// AgentConfig represents configuration for spawning an agent
type AgentConfig struct {
	Type   string
	Prompt string
}

// AgentResultFromJSON parses an AgentResult from JSON
func AgentResultFromJSON(jsonStr string) (AgentResult, error) {
	var result AgentResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return AgentResult{}, fmt.Errorf("failed to parse AgentResult JSON: %w", err)
	}
	return result, nil
}
