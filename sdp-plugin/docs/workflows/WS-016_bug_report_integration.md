# WS-016: Bug Report Flow Integration

> **Workstream ID:** WS-016  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-013

## Goal

Integrate bug detection and reporting into orchestrator workflow.

## Acceptance Criteria

### AC1: Bug Detection
- [ ] Bugs detected during execution
- [ ] Bug severity classified (P0/P1/P2/P3)
- [ ] Stack traces captured

### AC2: Auto-Create Issues
- [ ] Bugs create Beads issues
- [ ] Issue includes: description, severity, stack trace
- [ ] Issue linked to parent feature

### AC3: Block on P0
- [ ] P0 bugs block workstream execution
- [ ] P1/P2 logged but don't block
- [ ] P3 logged as warnings

## Scope Files

**internal/orchestrator/bug_handler.go** (NEW)

## Estimated Scope

- ~100 LOC
- Duration: 1 hour
