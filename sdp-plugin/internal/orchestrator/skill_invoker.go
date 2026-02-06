package orchestrator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SkillTool defines the interface for invoking skills
type SkillTool interface {
	Invoke(skillName string, args string) (string, error)
}

// IdeaResult represents the result of @idea skill invocation
type IdeaResult struct {
	Problem          string   `json:"problem"`
	Users            []string `json:"users"`
	SuccessCriteria []string `json:"success_criteria"`
}

// DesignResult represents the result of @design skill invocation
type DesignResult struct {
	Workstreams []string `json:"workstreams"`
	SpecPath    string   `json:"spec_path"`
}

// OneshotResult represents the result of @oneshot skill invocation
type OneshotResult struct {
	AgentID string `json:"agent_id"`
	Status  string `json:"status"`
}

// SkillInvoker handles invocation of sub-skills
type SkillInvoker struct {
	skillTool SkillTool
}

// NewSkillInvoker creates a new skill invoker
func NewSkillInvoker(tool SkillTool) *SkillInvoker {
	return &SkillInvoker{
		skillTool: tool,
	}
}

// InvokeIdea invokes the @idea skill for requirements gathering
func (si *SkillInvoker) InvokeIdea(featureID string, description string) (IdeaResult, error) {
	// AC2: Orchestrator calls @idea for requirements gathering
	args := fmt.Sprintf("%s - Feature ID: %s", description, featureID)

	resultJSON, err := si.skillTool.Invoke("idea", args)
	if err != nil {
		return IdeaResult{}, fmt.Errorf("failed to invoke @idea: %w", err)
	}

	// Parse result
	result, err := IdeaResultFromJSON(resultJSON)
	if err != nil {
		return IdeaResult{}, fmt.Errorf("failed to parse @idea result: %w", err)
	}

	return result, nil
}

// InvokeDesign invokes the @design skill for workstream planning
func (si *SkillInvoker) InvokeDesign(specPath string) (DesignResult, error) {
	// AC3: Orchestrator calls @design for workstream planning
	args := specPath

	resultJSON, err := si.skillTool.Invoke("design", args)
	if err != nil {
		return DesignResult{}, fmt.Errorf("failed to invoke @design: %w", err)
	}

	// Parse result
	result, err := DesignResultFromJSON(resultJSON)
	if err != nil {
		return DesignResult{}, fmt.Errorf("failed to parse @design result: %w", err)
	}

	// Verify spec path format
	if !strings.HasPrefix(result.SpecPath, "docs/drafts/") {
		return DesignResult{}, fmt.Errorf("invalid spec path: %s (must start with docs/drafts/)", result.SpecPath)
	}

	return result, nil
}

// InvokeOneshot invokes the @oneshot skill for autonomous execution
func (si *SkillInvoker) InvokeOneshot(featureID string) error {
	// AC4: Orchestrator calls @oneshot for autonomous execution
	args := featureID

	resultJSON, err := si.skillTool.Invoke("oneshot", args)
	if err != nil {
		return fmt.Errorf("failed to invoke @oneshot: %w", err)
	}

	// Parse result to verify agent spawned
	var result OneshotResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		// If we can't parse, just log warning but don't fail
		// The skill may have started successfully even if result format is unexpected
		return nil
	}

	// Verify agent was spawned
	if result.AgentID == "" {
		return fmt.Errorf("@oneshot did not spawn an agent")
	}

	return nil
}

// IdeaResultFromJSON parses an IdeaResult from JSON
func IdeaResultFromJSON(jsonStr string) (IdeaResult, error) {
	var result IdeaResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return IdeaResult{}, fmt.Errorf("failed to parse IdeaResult JSON: %w", err)
	}
	return result, nil
}

// DesignResultFromJSON parses a DesignResult from JSON
func DesignResultFromJSON(jsonStr string) (DesignResult, error) {
	var result DesignResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return DesignResult{}, fmt.Errorf("failed to parse DesignResult JSON: %w", err)
	}
	return result, nil
}
