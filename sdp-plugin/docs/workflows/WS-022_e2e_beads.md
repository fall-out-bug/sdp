# WS-022: E2E Tests with Real Beads

> **Workstream ID:** WS-022  
> **Feature:** F024: Unified Workflow  
> **Dependencies:** WS-020

## Goal

End-to-end tests with Beads task tracking.

## Acceptance Criteria

### AC1: Feature → Beads
- [ ] @feature creates Beads feature issue
- [ ] Workstreams create child tasks

### AC2: Workstream → Beads
- [ ] @build updates Beads status (in_progress → closed)
- [ ] Dependencies tracked

### AC3: Git Integration
- [ ] Commits include Beads metadata
- [ ] bd sync works
- [ ] Beads issues linked to commits

## Scope

E2E tests with real Beads CLI.

## Estimated Scope

- ~250 LOC (test code)
- Duration: 2 hours
