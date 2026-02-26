# Unified Workflow - Hybrid SDP Implementation
> **Status:** Completed
> **Created:** 2026-01-28
> **Workstreams:** 26 (18 completed/implemented, 8 optional/deferred)
> **Updated:** 2026-02-16

## Mission

Implement unified @feature workflow combining @idea/@design/@oneshot with team coordination (100+ roles, approval gates, notifications).

## Problem

Current SDP workflow has three separate entry points (@idea, @design, @oneshot) that:
- Don't coordinate with each other
- Don't track team dependencies
- Lack checkpoint/resume for long-running features
- No progressive disclosure for users

## Solution

Unified @feature skill that:
1. **Progressive Disclosure**: vision -> requirements -> planning -> execution
2. **Team Coordination**: 100+ role management, approval gates
3. **Checkpoint/Resume**: Save state for long-running features
4. **Notifications**: Gateway for team updates
5. **Agent Orchestration**: Multi-agent system for autonomous execution

## Success Criteria

- [x] @feature skill integrates @idea/@design/@oneshot in one workflow
- [x] TeamManager manages 100+ roles with activation/deactivation
- [x] Checkpoint save/resume works for orchestrator state
- [x] Approval gates enforce quality checkpoints
- [x] Notification gateway supports multiple channels
- [x] All core components have >=80% test coverage

## Completed Workstreams

### WS-001: Checkpoint database schema
**Status:** COMPLETED
**Implementation:** `internal/checkpoint/checkpoint.go`
- Checkpoint struct with agent_id, status, completed_ws
- JSON serialization for persistence
- CheckpointSaver interface for save/load/resume
**Test Coverage:** 84.4%

### WS-002: CheckpointRepository implementation
**Status:** COMPLETED
**Implementation:** `internal/checkpoint/`
- File-based repository (.sdp/checkpoints/{feature}-checkpoint.json)
- Thread-safe save/load operations
- Checkpoint merge on resume

### WS-003: OrchestratorAgent core logic
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/orchestrator.go`, `orchestrator_graph.go`
- Dependency graph building (topological sort)
- WorkstreamExecutor interface for @build integration
- Error handling with retries
- Circular dependency detection
**Test Coverage:** 83.4%

### WS-004: TeamManager role registry
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/team_manager.go`
- TeamRole struct (ID, name, description, permissions, status)
- RoleRegistry interface with 6 methods
- Active/dormant role switching
- Thread-safe operations (sync.RWMutex)

### WS-005: Team lifecycle management
**Status:** MERGED INTO WS-004
**Reason:** TeamManager includes lifecycle methods (ActivateRole, DeactivateRole)

### WS-006: ApprovalGateManager implementation
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/approval.go`, `approval_ops.go`
- ApprovalGate struct (status, approvers, required_approvals)
- 11 methods (CreateGate, Approve, Reject, CheckGateApproved, etc.)
- Gate enforcement with BlockExecutionUntilApproved
- Thread-safe operations

### WS-007: SkipFlagParser integration
**Status:** DEFERRED
**Reason:** Not required for current workflow

### WS-008: Checkpoint save/resume logic
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/orchestrator_execute.go`
- Save checkpoint after each workstream completion
- Resume from interrupted execution
- Skip completed workstreams on resume
- Checkpoint state management

### WS-009: @feature skill orchestrator
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/feature_coordinator.go`, `feature_coordinator_execute.go`, `feature_coordinator_resume.go`
- OrchestratorAgent integration
- Progress callback support
- Checkpoint coordination
- Evidence integration

### WS-010: Progressive menu UI
**Status:** COMPLETED (via skill system)
**Implementation:** `.claude/skills/feature/SKILL.md`
- AskUserQuestion integration for menu options
- Menu states: vision -> technical -> execute -> review
- Progress display with checkpoints

### WS-011: @idea/@design/@oneshot invocation
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/skill_invoker.go`
- Skill tool integration in orchestrator
- Call @idea for requirements gathering
- Call @design for workstream planning
- Call @oneshot for autonomous execution
**Test Coverage:** 88.8%

### WS-012: AgentSpawner via Task tool
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/agent_spawner.go`
- Agent type registry (planner, builder, reviewer, deployer)
- Agent result aggregation
- Spawn configuration
**Test Coverage:** 83.4%

### WS-013: SendMessage router
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/message_router.go`
- Message bus for agent communication
- Agent -> Orchestrator status updates
- Orchestrator -> Agent commands
- Message logging for debugging

### WS-014: RoleLoader and prompt management
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/role_loader.go`
- RolePromptLoader interface
- Load role definitions from .md files
- Role validation (name, description, permissions)

### WS-015: Dormant/active role switching
**Status:** COMPLETED (in WS-004)
**Implementation:** `internal/orchestrator/team_manager.go`
- ActivateRole(roleName, agentID) API
- DeactivateRole(agentID) API
- Role availability checking

### WS-016: Bug report flow integration
**Status:** DEFERRED
**Reason:** Available via Beads integration

### WS-017: NotificationProvider interface
**Status:** COMPLETED
**Implementation:** `internal/notification/gateway.go`, `channels.go`
- Channel interface (Send(), Name(), IsEnabled())
- Event types: feature_complete, drift_detected, review_complete, etc.
- Provider registration
**Test Coverage:** 82.9%

### WS-018: NotificationRouter implementation
**Status:** COMPLETED
**Implementation:** `internal/notification/gateway.go`
- Gateway with provider registry
- Rate limiting
- Notification history

### WS-019: TelegramNotifier + Mock provider
**Status:** COMPLETED (Mock provider)
**Implementation:** `internal/notification/channels.go`
- MockChannel for testing
- Extensible for Telegram integration
- Error handling

### WS-020: Unit tests for core components
**Status:** COMPLETED
**Coverage:**
- orchestrator: 83.4%
- checkpoint: 84.4%
- notification: 82.9%
- All internal packages: >=80%

### WS-021: Integration tests for agent coordination
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/orchestrator_test.go`, `feature_coordinator_test.go`
- Orchestrator with multiple workstreams
- Message routing tests
- Checkpoint/resume tests

### WS-022: E2E tests with real Beads
**Status:** DEFERRED
**Reason:** Requires external Beads setup, covered by integration tests

### WS-023: E2E tests with real Telegram
**Status:** DEFERRED
**Reason:** Requires external Telegram setup, covered by mock tests

### WS-024: Update PROTOCOL.md with unified workflow
**Status:** COMPLETED
**Implementation:** `docs/PROTOCOL.md`, `CLAUDE.md`
- @feature workflow documented
- Progressive disclosure phases
- Orchestrator integration

### WS-025: Create 15-minute tutorial
**Status:** COMPLETED
**Implementation:** `docs/TUTORIAL_FEATURE.md`
- Step-by-step example
- Working example produced

### WS-026: English translation + role setup guide
**Status:** COMPLETED
**Implementation:** `docs/ROLE_SETUP_GUIDE.md`
- Role setup guide complete
- Example roles documented

## Additional Components

### Structured Logging
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/logging.go`
- OrchestratorLogger with slog integration
- Unique correlation ID per feature execution
- Structured JSON logging
- 10+ logging methods (LogStart, LogWSStart, LogWSComplete, etc.)

### SLO Tracking
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/slos.go`, `slos_record.go`, `slos_status.go`
- SLOTracker for checkpoint save latency, WS execution time, graph build time, recovery success
- p95 percentile calculation
- Success rate calculation
- SLO breach detection and logging

### Snapshot System
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/snapshot.go`, `snapshot_ops.go`
- State snapshots for recovery
- Snapshot save/load operations

### Beads Integration
**Status:** COMPLETED
**Implementation:** `internal/orchestrator/beads_loader.go`
- Load workstreams from Beads mapping
- Beads task resolution

## Technical Approach

**Architecture:**
```
@feature (skill)
  |-- Phase 1: Vision Interview (AskUserQuestion)
  |-- Phase 2: Generate PRODUCT_VISION.md
  |-- Phase 3: Technical Interview (AskUserQuestion)
  |-- Phase 4: Generate intent.json
  |-- Phase 5: Create requirements draft (docs/drafts/idea-{slug}.md)
  +-- Phase 6: @design (creates workstreams)
      +-- Orchestrator (executes workstreams)
          |-- TeamManager (role coordination)
          |-- ApprovalGateManager (quality gates)
          |-- CheckpointSaver (persistence)
          +-- NotificationGateway (notifications)
```

**Technologies:**
- Go orchestrator (internal/orchestrator/)
- Skill system (@feature, @idea, @design, @oneshot)
- Beads for task tracking
- Notification gateway for team updates

## Non-Goals

- Time-based estimates (use scope: LOC/tokens only)
- Automatic workstream creation (still manual via @design)
- Real-time collaboration (async via Beads/notifications)
- Multi-user race conditions (single orchestrator at a time)

## Strategic Tradeoffs

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| Orchestrator in Go | Use Go for performance | Faster than Python, better concurrency |
| Beads integration | Use existing Beads | Leverage git-backed task tracking |
| Notifications | Gateway pattern | Extensible, testable |
| Checkpoint format | JSON file | Human-readable, git-friendly |
| Role prompts | Markdown files | Easy to edit, version control |

## Implementation Summary

**Packages Implemented:**
- `internal/orchestrator/` - 6,742 LOC, 83.4% coverage
- `internal/checkpoint/` - 84.4% coverage
- `internal/notification/` - 82.9% coverage

**CLI Commands:**
- `sdp orchestrate <feature-id>` - Execute all workstreams
- `sdp orchestrate resume <checkpoint-id>` - Resume from checkpoint

**Documentation:**
- `docs/PROTOCOL.md` - Updated with unified workflow
- `docs/TUTORIAL_FEATURE.md` - 15-minute tutorial
- `docs/ROLE_SETUP_GUIDE.md` - Role management guide
- `docs/NOTIFICATION_SYSTEM_GUIDE.md` - Notification guide

## Final Status

**Progress:** 100% of core workstreams completed

**Implementation Status:**
- **Core Components:** Checkpoint system, Orchestrator logic, Feature Coordinator
- **Team Coordination:** TeamManager, ApprovalGateManager, RoleLoader
- **Communication:** MessageRouter, NotificationGateway, AgentSpawner
- **Observability:** Logging, SLO Tracking, Snapshots
- **Documentation:** Protocol, Tutorial, Role Guide

**Test Coverage:** All core packages >= 80%
