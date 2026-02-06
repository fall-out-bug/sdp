# F024: Unified Workflow - Hybrid SDP Implementation

> **Feature ID:** F024  
> **Status:** In Progress  
> **Created:** 2026-01-28  
> **Workstreams:** 26 (8 completed, 18 open)

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
1. **Progressive Disclosure**: vision → requirements → planning → execution
2. **Team Coordination**: 100+ role management, approval gates
3. **Checkpoint/Resume**: Save state for long-running features
4. **Notifications**: Telegram integration for team updates
5. **Agent Orchestration**: Multi-agent system for autonomous execution

## Success Criteria

- [ ] @feature skill integrates @idea/@design/@oneshot in one workflow
- [ ] TeamManager manages 100+ roles with activation/deactivation
- [ ] Checkpoint save/resume works for orchestrator state
- [ ] Approval gates enforce quality checkpoints
- [ ] Telegram notifications sent for key events
- [ ] Progressive menu UI for user interaction
- [ ] All components have ≥80% test coverage
- [ ] E2E tests validate complete workflow

## Completed Workstreams (1-8)

### ✓ WS-001: Checkpoint database schema
**Status:** CLOSED  
**Implementation:** `internal/checkpoint/checkpoint.go`, `models.go`
- Checkpoint struct with agent_id, status, completed_ws
- JSON serialization for persistence
- CheckpointSaver interface for save/load/resume

### ✓ WS-002: CheckpointRepository implementation  
**Status:** CLOSED  
**Implementation:** `internal/checkpoint/repository.go`
- File-based repository (.oneshot/{feature}-checkpoint.json)
- Thread-safe save/load operations
- Checkpoint merge on resume

### ✓ WS-003: OrchestratorAgent core logic
**Status:** CLOSED  
**Implementation:** `internal/orchestrator/orchestrator.go`
- Dependency graph building (topological sort)
- WorkstreamExecutor interface for @build integration
- Error handling with retries
- Circular dependency detection

### ✓ WS-004: TeamManager role registry
**Status:** CLOSED  
**Implementation:** `internal/orchestrator/team_manager.go`
- Role struct (name, description, permissions, status)
- RoleRegistry with active/dormant role switching
- 100+ role support

### ✓ WS-005: Team lifecycle management
**Status:** CLOSED  
**Implementation:** `internal/orchestrator/team_lifecycle.go`
- AddRole(), RemoveRole(), ActivateRole(), DeactivateRole()
- Role status validation
- Team coordination APIs

### ✓ WS-006: ApprovalGateManager implementation
**Status:** CLOSED  
**Implementation:** `internal/orchestrator/approval.go`
- ApprovalGate struct (name, required_approvals)
- RequestApproval(), Approve(), Reject()
- Gate enforcement before workstream execution

### ✓ WS-007: SkipFlagParser integration
**Status:** CLOSED  
**Implementation:** `internal/orchestrator/skip_parser.go`
- Parse --skip flag with workstream IDs
- Validate skip list against dependencies
- Prevent breaking dependency chains

### ✓ WS-008: Checkpoint save/resume logic
**Status:** CLOSED  
**Implementation:** Orchestrator with CheckpointSaver integration
- Save checkpoint after each workstream
- Resume from interrupted execution
- Skip completed workstreams on resume

## Pending Workstreams (9-26)

### ○ WS-009: @feature skill orchestrator
**Dependencies:** WS-008 ✅  
**Status:** OPEN (next to implement)

**Goal:** Integrate OrchestratorAgent into @feature skill

**What to build:**
- Update `prompts/skills/feature.md` to use orchestrator
- Add workflow: vision → requirements → planning → orchestrator execution
- Progressive menu UI for user choices
- Team coordinator calls (role management, approvals)

**Acceptance Criteria:**
- AC1: @feature calls orchestrator after @design phase
- AC2: Checkpoint saved after each workstream
- AC3: Team approvals requested at gates
- AC4: Telegram notifications on milestones

**Scope Files:**
- `prompts/skills/feature.md` (update with orchestrator flow)
- `internal/orchestrator/feature_coordinator.go` (NEW - coordination logic)
- `internal/orchestrator/feature_coordinator_test.go` (NEW)

**Definition of Done:**
- @feature orchestrates multi-workstream execution
- Progress tracking with timestamps
- Human-in-the-loop at approval gates
- Checkpoint/resume functional

### ○ WS-010: Progressive menu UI
**Dependencies:** WS-009  
**Status:** OPEN

**Goal:** Interactive menu for user choices during @feature

**What to build:**
- AskUserQuestion integration for menu options
- Menu states: vision → technical → execute → review
- Skip to specific phase flags
- Progress display with checkpoints

**Acceptance Criteria:**
- AC1: User can skip vision phase (--vision-only)
- AC2: User can start from existing spec (--spec PATH)
- AC3: Progress shown: "[HH:MM] Executing WS-XXX..."
- AC4: Menu choices logged as decisions

### ○ WS-011: @idea/@design/@oneshot invocation
**Dependencies:** WS-010  
**Status:** OPEN

**Goal:** Call sub-skills from orchestrator

**What to build:**
- Skill tool integration in orchestrator
- Call @idea for requirements gathering
- Call @design for workstream planning
- Call @oneshot for autonomous execution

**Acceptance Criteria:**
- AC1: Orchestrator can invoke Skill tool
- AC2: Sub-skill results captured and saved
- AC3: Error handling if sub-skill fails

### ○ WS-012: AgentSpawner via Task tool
**Dependencies:** WS-011  
**Status:** OPEN

**Goal:** Spawn specialist agents (QA, Security, DevOps, etc.)

**What to build:**
- Task tool integration for agent spawning
- Agent type registry (planner, builder, reviewer, deployer)
- Agent result aggregation

**Acceptance Criteria:**
- AC1: Can spawn planner agent for architecture
- AC2: Can spawn builder agent for implementation
- AC3: Can spawn reviewer agent for quality checks
- AC4: Agent results saved to checkpoint

### ○ WS-013: SendMessage router
**Dependencies:** WS-012  
**Status:** OPEN

**Goal:** Route messages between agents and orchestrator

**What to build:**
- Message bus for agent communication
- Agent → Orchestrator status updates
- Orchestrator → Agent commands

**Acceptance Criteria:**
- AC1: Agents send progress updates
- AC2: Orchestrator can send pause/resume commands
- AC3: Message logging for debugging

### ○ WS-014: RoleLoader and prompt management
**Dependencies:** WS-012  
**Status:** OPEN

**Goal:** Load role prompts from prompts/roles/

**What to build:**
- RolePromptLoader interface
- Load role definitions from .md files
- Role validation (name, description, permissions)

**Acceptance Criteria:**
- AC1: Load 100+ roles from filesystem
- AC2: Validate role schema
- AC3: Cache role prompts for performance

### ○ WS-015: Dormant/active role switching
**Dependencies:** WS-014  
**Status:** OPEN

**Goal:** Dynamic role activation/deactivation

**What to build:**
- ActivateRole(roleName, agentID) API
- DeactivateRole(agentID) API
- Role availability checking

**Acceptance Criteria:**
- AC1: Can activate dormant role
- AC2: Can release active role
- AC3: Prevent role conflicts (same role to 2 agents)

### ○ WS-016: Bug report flow integration
**Dependencies:** WS-013  
**Status:** OPEN

**Goal:** Integrate bug reports into orchestrator workflow

**What to build:**
- Bug detection during execution
- Auto-create issues in tracking system
- Block execution on P0 bugs

**Acceptance Criteria:**
- AC1: Bugs create Beads issues
- AC2: P0 bugs block workstream
- AC3: Bug reports include stack traces

### ○ WS-017: NotificationProvider interface
**Dependencies:** WS-013  
**Status:** OPEN

**Goal:** Define notification provider interface

**What to build:**
- NotificationProvider interface (Send(), Notify())
- Event types: workstream_complete, approval_needed, bug_found
- Provider registration

**Acceptance Criteria:**
- AC1: Interface defined with 3 methods
- AC2: Support multiple providers (Telegram, Email, Mock)
- AC3: Async notification sending

### ○ WS-018: NotificationRouter implementation
**Dependencies:** WS-017  
**Status:** OPEN

**Goal:** Route notifications to appropriate providers

**What to build:**
- NotificationRouter with provider registry
- Event type → provider mapping
- Fallback providers

**Acceptance Criteria:**
- AC1: Routes workstream events to Telegram
- AC2: Routes approval events to all providers
- AC3: Mock provider for testing

### ○ WS-019: TelegramNotifier + Mock provider
**Dependencies:** WS-018  
**Status:** OPEN

**Goal:** Telegram bot integration for notifications

**What to build:**
- Telegram bot client
- Send message to channel/group
- Mock provider for testing

**Acceptance Criteria:**
- AC1: Sends formatted messages to Telegram
- AC2: Handles send errors gracefully
- AC3: Mock provider for unit tests

### ○ WS-020: Unit tests for core components
**Dependencies:** WS-015, WS-018  
**Status:** OPEN

**Goal:** ≥80% coverage for orchestrator components

**What to test:**
- Orchestrator dependency graph
- TeamManager role operations
- ApprovalGateManager enforcement
- CheckpointRepository persistence

**Acceptance Criteria:**
- AC1: All packages ≥80% coverage
- AC2: No code without tests
- AC3: All edge cases covered

### ○ WS-021: Integration tests for agent coordination
**Dependencies:** WS-019  
**Status:** OPEN

**Goal:** Test multi-agent workflows

**What to build:**
- Test orchestrator with 2+ agents
- Test message routing
- Test checkpoint/resume with agents

**Acceptance Criteria:**
- AC1: Orchestrator coordinates 3 agents
- AC2: Messages routed correctly
- AC3: Checkpoint saves agent state

### ○ WS-022: E2E tests with real Beads
**Dependencies:** WS-020  
**Status:** OPEN

**Goal:** End-to-end tests with Beads task tracking

**What to build:**
- Test @feature → bd create → bd close
- Test bd dep add for dependencies
- Test bd sync with git

**Acceptance Criteria:**
- AC1: @feature creates Beads tasks
- AC2: Workstreams update Beads status
- AC3: Git commits include Beads metadata

### ○ WS-023: E2E tests with real Telegram
**Dependencies:** WS-019  
**Status:** OPEN

**Goal:** Test Telegram notifications end-to-end

**What to build:**
- Send test notification via Telegram
- Verify message format
- Test error handling

**Acceptance Criteria:**
- AC1: Test message received in Telegram
- AC2: Message formatted correctly
- AC3: Errors logged, don't crash

### ○ WS-024: Update PROTOCOL.md with unified workflow
**Dependencies:** WS-009  
**Status:** OPEN

**Goal:** Document @feature workflow in PROTOCOL.md

**What to build:**
- Add @feature section to PROTOCOL.md
- Document progressive disclosure phases
- Document orchestrator integration

**Acceptance Criteria:**
- AC1: PROTOCOL.md describes @feature workflow
- AC2: Examples for each phase
- AC3: Beads integration documented

### ○ WS-025: Create 15-minute tutorial
**Dependencies:** WS-024  
**Status:** OPEN

**Goal:** Quick start guide for @feature

**What to build:**
- Tutorial: docs/TUTORIAL_FEATURE.md
- Step-by-step example: @feature "Add login"
- Screenshots/diagrams

**Acceptance Criteria:**
- AC1: Tutorial takes 15 min to complete
- AC2: Working example produced
- AC3: Screenshots for key steps

### ○ WS-026: English translation + role setup guide
**Dependencies:** WS-025  
**Status:** OPEN

**Goal:** Translate docs to English, setup guide

**What to build:**
- Translate TUTORIAL_FEATURE.md to English
- Create role setup guide
- Document 100+ roles

**Acceptance Criteria:**
- AC1: All docs in English
- AC2: Role setup guide complete
- AC3: 10 example roles documented

## Technical Approach

**Architecture:**
```
@feature (skill)
  ├── Phase 1: Vision Interview (AskUserQuestion)
  ├── Phase 2: Generate PRODUCT_VISION.md
  ├── Phase 3: Technical Interview (AskUserQuestion)
  ├── Phase 4: Generate intent.json
  ├── Phase 5: Create requirements draft (docs/drafts/idea-{slug}.md)
  └── Phase 6: @design (creates workstreams)
      └── Orchestrator (executes workstreams)
          ├── TeamManager (role coordination)
          ├── ApprovalGateManager (quality gates)
          ├── CheckpointSaver (persistence)
          └── NotificationRouter (Telegram)
```

**Technologies:**
- Go orchestrator (already built in WS-003 through WS-008)
- Skill system (@feature, @idea, @design, @oneshot)
- Beads for task tracking
- Telegram Bot API for notifications

## Non-Goals

- ❌ Time-based estimates (use scope: LOC/tokens only)
- ❌ Automatic workstream creation (still manual via @design)
- ❌ Real-time collaboration (async via Beads/Telegram)
- ❌ Multi-user race conditions (single orchestrator at a time)

## Strategic Tradeoffs

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| Orchestrator in Go | Use Go for performance | Faster than Python, better concurrency |
| Beads integration | Use existing Beads | Leverage git-backed task tracking |
| Telegram notifications | Single channel | Simple, async, no UI needed |
| Checkpoint format | JSON file | Human-readable, git-friendly |
| Role prompts | Markdown files | Easy to edit, version control |

## Success Metrics

- [ ] @feature reduces multi-workstream coordination time by 50%
- [ ] 100+ roles supported with <100ms activation time
- [ ] Checkpoint/resume works for 20+ workstream features
- [ ] Telegram notifications sent within 5 seconds of events
- [ ] Zero data loss (checkpoint on every workstream)
- [ ] E2E tests validate complete workflow (vision → deployment)

## Risks

- **Risk 1:** Orchestrator complexity - mitigated by WS-003 through WS-008 building foundation
- **Risk 2:** Role management scale - mitigated by TeamManager (WS-004, WS-005)
- **Risk 3:** Telegram API rate limits - mitigated by async sending, error handling
- **Risk 4:** Checkpoint corruption - mitigated by JSON validation, git backups

## Next Steps

1. **Immediate:** Implement WS-009 (@feature skill orchestrator) - READY
2. **After WS-009:** WS-010 (Progressive menu UI)
3. **After WS-010:** WS-011 (@idea/@design/@oneshot invocation)
4. **Continue sequentially** through WS-026 following dependency chain

---

**Total Estimated Scope:**
- **Workstreams:** 26 (8 completed, 18 open)
- **Lines of Code:** ~3,000-5,000 LOC (excluding tests)
- **Test Coverage:** ≥80% target
- **Duration:** 3-4 weeks (sequential execution)

**Current Progress:** 31% complete (8/26 workstreams done)
