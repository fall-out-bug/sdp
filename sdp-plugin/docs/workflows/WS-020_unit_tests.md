# WS-020: Unit Tests for Core Components

> **Workstream ID:** WS-020  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-015, WS-018

## Goal

Achieve ≥80% test coverage for orchestrator components.

## Acceptance Criteria

### AC1: Orchestrator Coverage
- [ ] internal/orchestrator/orchestrator.go ≥80%
- [ ] Test dependency graph
- [ ] Test topological sort

### AC2: TeamManager Coverage
- [ ] internal/orchestrator/team_manager.go ≥80%
- [ ] Test role operations
- [ ] Test role validation

### AC3: ApprovalGateManager Coverage
- [ ] internal/orchestrator/approval.go ≥80%
- [ ] Test gate enforcement
- [ ] Test approval workflow

### AC4: Checkpoint Coverage
- [ ] internal/checkpoint ≥80%
- [ ] Test save/load/resume
- [ ] Test JSON serialization

## Scope

Write tests for all core components.

## Estimated Scope

- ~400 LOC (test code)
- Duration: 3 hours
