# WS-012: AgentSpawner via Task Tool

> **Workstream ID:** WS-012  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-011

## Goal

Enable orchestrator to spawn specialist agents (planner, builder, reviewer) using Task tool.

## Acceptance Criteria

### AC1: Spawn Planner Agent
- [ ] Can spawn planner agent for architecture
- [ ] Agent ID captured for resume
- [ ] Results aggregated

### AC2: Spawn Builder Agent
- [ ] Can spawn builder agent for implementation
- [ ] Agent executes workstreams
- [ ] Results captured

### AC3: Spawn Reviewer Agent
- [ ] Can spawn reviewer agent for quality checks
- [ ] Multi-agent review (6 agents)
- [ ] Verdict: APPROVED/CHANGES_REQUESTED

### AC4: Agent Result Aggregation
- [ ] Agent results saved to checkpoint
- [ ] Errors logged with context
- [ ] Agent terminated after completion

## Scope Files

**internal/orchestrator/agent_spawner.go** (NEW)
```go
package orchestrator

type AgentSpawner struct {
    taskTool TaskTool
}

func (as *AgentSpawner) Spawn(agentType string, prompt string) (AgentID, error)
func (as *AgentSpawner) GetResult(agentID string) (AgentResult, error)
func (as *AgentSpawner) Terminate(agentID string) error
```

## Definition of Done

- [ ] All 4 AC met
- [ ] Agent spawning functional
- [ ] Results captured
- [ ] Tests ≥80% coverage

## Estimated Scope

- ~200 LOC implementation
- ~100 LOC tests
- Duration: 2 hours
- Size: SMALL
