# WS-015: Dormant/Active Role Switching

> **Workstream ID:** WS-015  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-014

## Goal

Dynamic role activation/deactivation for agent coordination.

## Acceptance Criteria

### AC1: Activate Role
- [ ] ActivateRole(roleName, agentID) works
- [ ] Role marked as active
- [ ] Agent granted permissions

### AC2: Deactivate Role
- [ ] DeactivateRole(agentID) releases role
- [ ] Role marked as dormant
- [ ] Permissions revoked

### AC3: Conflict Prevention
- [ ] Same role cannot be active for 2 agents
- [ ] Queue for role requests
- [ ] Timeout for role release

## Scope Files

**internal/orchestrator/role_switching.go** (NEW)

## Estimated Scope

- ~150 LOC
- Duration: 1.5 hours
