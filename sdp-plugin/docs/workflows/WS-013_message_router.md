# WS-013: SendMessage Router

> **Workstream ID:** WS-013  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-012

## Goal

Implement message bus for agent-orchestrator communication.

## Acceptance Criteria

### AC1: Agent → Orchestrator
- [ ] Agents send progress updates
- [ ] Status changes captured
- [ ] Errors reported immediately

### AC2: Orchestrator → Agent
- [ ] Orchestrator sends pause/resume commands
- [ ] Agent receives commands
- [ ] Command execution confirmed

### AC3: Message Logging
- [ ] All messages logged for debugging
- [ ] Timestamp on each message
- [ ] Message format validated

## Scope Files

**internal/orchestrator/message_router.go** (NEW)
**internal/orchestrator/message_router_test.go** (NEW)

## Estimated Scope

- ~150 LOC
- Duration: 1.5 hours
