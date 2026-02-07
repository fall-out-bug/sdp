# WS-021: Integration Tests for Agent Coordination

> **Workstream ID:** WS-021  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-019

## Goal

Test multi-agent workflows with orchestrator.

## Acceptance Criteria

### AC1: Multi-Agent Coordination
- [ ] Orchestrator coordinates 3+ agents
- [ ] Agents receive commands
- [ ] Agents send status updates

### AC2: Message Routing
- [ ] Messages routed correctly
- [ ] Async message delivery
- [ ] No lost messages

### AC3: Checkpoint with Agents
- [ ] Checkpoint saves agent state
- [ ] Resume restores agent context
- [ ] Agent continues from checkpoint

## Scope

Integration tests using mock agents.

## Estimated Scope

- ~300 LOC (test code)
- Duration: 2.5 hours
