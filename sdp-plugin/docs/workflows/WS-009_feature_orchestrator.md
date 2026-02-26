# WS-009: @feature Skill Orchestrator

> **Workstream ID:** WS-009  
> **Feature:** Unified Workflow  
> **Status:** Ready to Implement  
> **Dependencies:** WS-008 ✅ (Checkpoint save/resume logic - COMPLETE)

## Goal

Integrate OrchestratorAgent into @feature skill to enable autonomous multi-workstream execution with progressive disclosure, team coordination, and checkpoint/resume capabilities.

## Acceptance Criteria

### AC1: Orchestrator Integration
- @feature calls orchestrator after @design completes
- Orchestrator executes workstreams in dependency order
- Workstreams executed via @build skill (Beads integration)

### AC2: Checkpoint Management
- Checkpoint saved to .oneshot/{feature}-checkpoint.json
- Checkpoint includes: agent_id, status, completed_ws, execution_order
- On resume, skip completed workstreams

### AC3: Team Coordination
- ApprovalGateManager.RequestApproval() called at gates
- Execution pauses until approval granted
- TeamManager roles activated/deactivated as needed

### AC4: Progress Tracking
- Real-time updates: "[HH:MM] Executing WS-XXX..."
- Timestamps for each workstream (start, complete)
- Summary on completion: "X/Y workstreams, Zm total, avg coverage"

### AC5: Error Handling
- Workstream failures logged with context
- Auto-retry for HIGH/MEDIUM issues (max 2 retries)
- CRITICAL errors escalate to human

## Scope Files

### Files to Modify
**prompts/skills/feature.md** - Add Phase 7: Orchestrator Execution

### Files to Create
**internal/orchestrator/feature_coordinator.go** - Main coordination logic
**internal/orchestrator/feature_coordinator_test.go** - Tests

## Implementation Steps

1. Update @feature skill (add Phase 7)
2. Implement FeatureCoordinator struct
3. Add Execute() method
4. Add Resume() method
5. Add tests (≥80% coverage)

## Definition of Done
- All 5 AC met
- Code committed
- Tests passing
- No TODOs

## Estimated Scope
- ~300 LOC implementation
- ~150 LOC tests
- Duration: 2-3 hours
- Size: SMALL
