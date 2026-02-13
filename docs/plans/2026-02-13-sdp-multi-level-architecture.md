# SDP Multi-Level Architecture Design

> **Status:** Research complete
> **Date:** 2026-02-13
> **Goal:** Design SDP as a multi-level product - from pure protocol (L0) to full orchestration platform (L3-L4)

---

## Table of Contents

1. [Overview](#overview)
2. [Level 0: Protocol Definition](#1-level-0-protocol-definition)
3. [Level 1: Hooks & Basic CLI](#2-level-1-hooks--basic-cli)
4. [Level 2: Go Tools - Provenance](#3-level-2-go-tools---provenance)
5. [Level 2: Go Tools - Guard](#4-level-2-go-tools---guard)
6. [Level 3: Agent Orchestration](#5-level-3-agent-orchestration)
7. [Level 4: AI-Human Collaboration](#6-level-4-ai-human-collaboration)
8. [Level Interfaces](#7-level-interfaces)
9. [Migration Path](#8-migration-path)
10. [Implementation Plan](#implementation-plan)

---

## Overview

### Goals

1. **L0 (Protocol)** — LLM-agnostic, tool-agnostic workflow specification that ANY AI can follow
2. **L1 (Hooks)** — Lightweight safety guardrails (git hooks, basic CLI) that help but don't require
3. **L2 (Go Tools)** — Reliability layer: provenance, traceability, guard enforcement
4. **L3 (Orchestration)** — Cross-tool agent coordination, distributed execution
5. **L4 (Collaboration)** — AI-Human shared workspace, cross-review, notification

### Key Decisions

| Aspect | Decision |
|--------|----------|
| L0 Protocol | Protocol + Reference Implementations (L1 adapters per tool) |
| L1 CLI | Minimal Safety-First (safety hooks, optional enforcement) |
| L2 Provenance | Dual-Log Architecture (evidence + state, hash-chained) |
| L2 Guard | Pre-Commit + Post-Commit with Declarative Policies |
| L3 Orchestration | Hierarchical Orchestrator with Agent Backends |
| L4 Collaboration | Hybrid Notification System (multi-channel) |
| Interfaces | Protocol-First Interface Definition with schemas |
| Migration | Single-Level Protocol with Optional Enhancements |

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│  L4: AI-Human Collaboration                                     │
│  ├── Beads Issues (shared workspace)                            │
│  ├── Notification Channels (Webhook, Desktop, Log)              │
│  └── Cross-Review (AI ↔ Human)                                  │
├─────────────────────────────────────────────────────────────────┤
│  L3: Agent Orchestration (Future)                               │
│  ├── Orchestrator (dependency graph, topological sort)          │
│  ├── AgentBackend interface (Local, HTTP, gRPC)                 │
│  └── Distributed execution (k8s-ready)                          │
├─────────────────────────────────────────────────────────────────┤
│  L2: Go Tools (Reliability)                                     │
│  ├── Evidence Log (hash-chained events.jsonl)                   │
│  ├── State Log (checkpoints, resume)                            │
│  ├── Guard (scope enforcement, drift detection)                 │
│  └── Collision Detection (parallel WS safety)                   │
├─────────────────────────────────────────────────────────────────┤
│  L1: Hooks & Basic CLI                                          │
│  ├── Git Hooks (pre-commit, pre-push, post-checkout)            │
│  ├── sdp init / doctor / hooks                                  │
│  └── Claude Code Hooks (PreToolUse, PostToolUse)                │
├─────────────────────────────────────────────────────────────────┤
│  L0: Protocol (Foundation)                                      │
│  ├── Workstream format (PP-FFF-SS, YAML frontmatter)            │
│  ├── Quality gates (coverage, LOC, type hints, clean arch)      │
│  ├── TDD workflow (Red → Green → Refactor)                      │
│  └── Skills/Agents (tool-agnostic definitions)                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 1. Level 0: Protocol Definition

> **Experts:** Martin Fowler (Refactoring), Theo Browne (API Design), Sam Newman (Architecture)

### Decision: Protocol + Reference Implementations

L0 contains ONLY the workflow essence that ALL implementations share. Tool-specific adapters live in L1.

### What Belongs in L0

**Core Concepts:**
- **Workstream**: Atomic task (SMALL/MEDIUM/LARGE)
- **Feature**: 5-30 workstreams
- **Release**: 10-30 features

**Workflow Phases:**
1. Strategic Planning → Product vision, PRD, roadmap
2. Analysis → Codebase health, gaps, tech debt
3. Feature Planning → Requirements → Workstreams
4. Execution → TDD cycle → Quality gates → Review → Deploy

**Quality Gates:**
- Test coverage ≥80%
- Files ≤200 LOC
- Type hints required
- Clean architecture (no layer violations)
- Explicit error handling (no bare exceptions)

**Workstream Format:**
```yaml
---
ws_id: PP-FFF-SS
feature_id: FFFF
title: "Description"
status: backlog|in-progress|complete
priority: 0|1|2|3
size: SMALL|MEDIUM|LARGE
depends_on: [ws-id, ...]
blocks: [ws-id, ...]
---
## Goal
{What this workstream achieves}

## Acceptance Criteria
- [ ] AC1: {testable criterion}

## Scope Files
{Files that may be modified}
```

### What Does NOT Belong in L0

- Skill definitions (@vision, @build, etc.) → Tool-specific L1
- Agent roles (implementer, reviewer) → Tool-specific L1
- Tool-specific APIs (AskUserQuestion, Task) → L1 adapters
- CLI commands (sdp plan, sdp apply) → L2
- Evidence logging, Beads integration → L2

### L1 Adapters per Tool

| Tool | L1 Location | Adapter Pattern |
|------|-------------|-----------------|
| Claude Code | `.claude/` | Skills use AskUserQuestion, Task tools |
| Cursor | `.cursor/` | Skills use Cursor's agent panel |
| Windsurf | `.windsurf/` | Skills use Windsurf's cascade |
| Codex | `.codex/` | Skills use Codex's capabilities |

**Translation Layer:**
- L0 "ask questions" → Claude Code `AskUserQuestion(...)` or Cursor interactive prompts
- L0 "spawn specialist agents" → Claude Code `Task(subagent_type=...)` or Cursor agent panel

---

## 2. Level 1: Hooks & Basic CLI

> **Experts:** Kelsey Hightower (DevOps), Theo Browne (API Design), Martin Fowler (Refactoring)

### Decision: Minimal Safety-First

L1 provides ONLY essential safety guardrails (git safety). L2 adds protocol compliance.

### Essential L1 CLI Commands

| Command | Purpose |
|---------|---------|
| `sdp init` | Bootstrap environment |
| `sdp doctor` | Validate environment setup |
| `sdp hooks install/uninstall` | Manage safety hooks |
| `sdp guard activate/check` | Scope enforcement (bridge to L2) |

### Essential Git Hooks

| Hook | Purpose |
|------|---------|
| **pre-commit** | Block commits to protected branches, validate session |
| **pre-push** | Prevent direct pushes to main/dev |
| **PreToolUse** (Claude) | Block `git reset --hard`, `git clean`, destructive commands |

### Design Philosophy

- **Helpful but not required**: L1 enhances but doesn't gate usage
- **Clear boundaries**: L1 = environment safety, L2+ = protocol compliance
- **Graceful degradation**: CLI → Shell fallback → Prompt-native logic

### What Moves to L2

- Quality gates (coverage, type checking, file size)
- Protocol compliance (workstream validation, drift detection)
- Telemetry and evidence tracking

---

## 3. Level 2: Go Tools - Provenance

> **Experts:** Martin Kleppmann (Distributed Systems), Theo Browne (API Design), Troy Hunt (Security)

### Decision: Dual-Log Architecture

Separate evidence log (audit) from state log (checkpoint) with cross-referencing.

### Evidence Log (`events.jsonl`)

**Purpose:** Audit trail for all workstream activities

**Format:** JSONL with hash-chaining
```json
{"type":"plan","ws_id":"00-067-01","timestamp":"...","prev_hash":"abc123","data":{...}}
```

**Event Types:**
- `plan`: Workstream activation/design
- `generation`: Code generation with model metadata
- `verification`: Quality gate results
- `approval`: Deployment/approval events
- `decision`: Architectural decisions
- `checkpoint`: State snapshots (NEW)

### State Log (`state.jsonl`)

**Purpose:** Execution state for resume/recovery

**Format:** Mutable snapshots, NOT hash-chained
```json
{"feature":"F067","completed_ws":["00-067-01","00-067-02"],"current":"00-067-03"}
```

### L1-L2 Contract

```go
// EvidenceEmitter interface - L2 provides to L1
type EvidenceEmitter interface {
    Emit(ev *Event) error         // Async, non-blocking
    EmitSync(ev *Event) error     // Sync, for critical events
    Checkpoint(wsID string, state interface{}) error
}
```

### Key Improvements Needed

1. **Add checkpoint events** to evidence log (currently invisible)
2. **Add session_id** for linking events across resumes
3. **Add retry tracking** for failed generations
4. **Fix Emit()** to not silently drop errors

---

## 4. Level 2: Go Tools - Guard

> **Experts:** Kelsey Hightower (DevOps), Martin Fowler (Refactoring), Theo Browne (API Design)

### Decision: Pre-Commit + Post-Commit with Declarative Policies

Guard enforces scope at commit time, CI verifies before merge.

### Current Implementation

| Command | Purpose |
|---------|---------|
| `sdp guard activate <ws-id>` | Set active workstream |
| `sdp guard check <file>` | Validate single file edit |
| `sdp guard check --staged` | Validate staged files |
| `sdp guard status` | Show active WS and scope |
| `sdp drift detect <ws-id>` | Detect doc vs code drift |
| `sdp collision check` | Detect scope overlaps |

### Declarative Policy Format

```yaml
# .sdp/guard-policy.yml
apiVersion: sdp.io/v1
kind: GuardPolicy
spec:
  enforcement:
    mode: hybrid  # strict | hybrid | permissive
    scopeViolation: warning
    driftDetection: error

  scope:
    allowNewFiles: true
    requireTests: true

  rules:
    - id: max-file-loc
      enabled: true
      severity: error
      config:
        maxLines: 200
```

### Key Improvements Needed

1. **Move state** from `~/.config/sdp/state.json` to `.sdp/guard-state.json` (project-relative)
2. **Add glob patterns** to scope_files
3. **Auto-install hook** on `sdp init`
4. **Provide GitHub Action** for CI drift detection

---

## 5. Level 3: Agent Orchestration

> **Experts:** Martin Kleppmann (Distributed Systems), Kelsey Hightower (DevOps), Sam Newman (Architecture)

### Decision: Hierarchical Orchestrator with Agent Backends

Keep central orchestrator but abstract execution backend. Gradual migration to distributed.

### Architecture

```
┌─────────────────────────────────────────┐
│           Orchestrator                   │
│  ┌─────────────────────────────────────┐│
│  │ Dependency Graph (DAG)              ││
│  │ Topological Sort                    ││
│  │ Checkpoint Manager                  ││
│  └─────────────────────────────────────┘│
│  ┌─────────────────────────────────────┐│
│  │ BackendRegistry                     ││
│  │ ├── LocalBackend (current)          ││
│  │ ├── HTTPBackend (future)            ││
│  │ └── gRPCBackend (future)            ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
```

### AgentBackend Interface

```go
type AgentBackend interface {
    ID() string
    Kind() string           // "local", "http", "grpc"
    Execute(wsID string) error
    Stream(wsID string) (<-chan Event, error)
    Capabilities() []string
    HealthCheck() error
}
```

### Migration Path

1. **Phase 1**: LocalBackend wraps existing ExecuteFunc (zero behavior change)
2. **Phase 2**: Add HTTPBackend for remote agents
3. **Phase 3**: Add gRPCBackend for streaming
4. **Phase 4**: k8s deployment with backend discovery

### Why Not Pure Event Mesh?

- **Preserves dependency ordering**: SDP's core value is DAG-based execution
- **Backward compatible**: LocalBackend is thin wrapper
- **Simpler fault model**: Orchestrator is SPOF but highly available (k8s replicas)

---

## 6. Level 4: AI-Human Collaboration

> **Experts:** Nir Eyal (UX/Product), Dan Abramov (React/State), Kent C. Dodds (Testing)

### Decision: Hybrid Notification System (Multi-channel)

Build on existing notification infrastructure (`channels.go`) with Beads integration.

### Current State

- **Beads Issues**: Shared workspace with `owner`, `status`, `priority` fields
- **Notification Channels**: Log, Webhook, Desktop (partially implemented)
- **Handoff Protocol**: Session templates for AI→Human transitions

### Key Features to Implement

1. **Claiming Mechanism**: `sdp claim <issue-id>` prevents AI-human conflicts
2. **Enhanced Beads Schema**: Add `claimed_by`, `claimed_at` fields
3. **Notification Triggers**:
   - AI creates issue → Webhook/Desktop notify
   - AI claims issue → Update Beads status
   - Human updates issue → AI detects on next cycle
4. **Conflict Detection**: Check `claimed_by` matches current agent before WS execution
5. **Handoff Enhancement**: Include `claimed_tasks`, `blocked_tasks` in session template

### Notification Flow

```
AI creates issue
    ↓
Webhook/Desktop notification → Human
    ↓
Human reviews, assigns back
    ↓
AI detects assignment (poll/trigger)
    ↓
AI executes workstream
    ↓
AI creates PR → Human review request
```

---

## 7. Level Interfaces

> **Experts:** Theo Browne (API Design), Sam Newman (Architecture), Kent C. Dodds (Testing)

### Decision: Protocol-First Interface Definition with Schemas

Define explicit schemas for all data formats with validation at boundaries.

### Key Interfaces

| Boundary | Format | Validation |
|----------|--------|------------|
| L0 → L1 | WS YAML | JSON Schema, Go structs |
| L1 → L2 | CLI JSON output | Structured output contracts |
| L2 → L3 | Evidence events | Event schema registry |
| Cross-Level | Config YAML | Versioned, validated |

### Skill Frontmatter Contract

```yaml
---
name: oneshot
cli: sdp orchestrate        # CLI command (optional)
mode: hybrid                # cli_only | llm_only | hybrid
llm_subagents: [implementer, spec-reviewer, quality-reviewer]
input_format:               # Expected inputs
  feature_id: string
output_format:              # Expected outputs
  status: success | failure
  workstreams_executed: int
errors:                     # Error conditions
  - feature_not_found
  - dependency_cycle
---
```

### Schema Registry

```
src/sdp/schema/
├── workstream.schema.json
├── evidence.schema.json
├── checkpoint.schema.json
├── config.schema.json
└── session.schema.json
```

---

## 8. Migration Path

> **Experts:** Martin Fowler (Refactoring), Dan Abramov (Progressive Disclosure), Kelsey Hightower (DevOps)

### Decision: Single-Level Protocol with Optional Enhancements

Protocol is universal. Hooks and CLI are additive enhancements, not separate levels.

### Migration Steps

```
L0: Protocol Only
├── Read CLAUDE.md + PROTOCOL.md
├── Follow workstream convention manually
├── Use @feature, @build skills
└── No validation automation

L0 → L1: Add Hooks
├── bash hooks/install-hooks.sh
└── Get: Pre-commit validation (file size, architecture)

L1 → L2: Add CLI
├── cd sdp-plugin && go build
└── Get: Parallel execution, evidence logging, guard
```

### Handling Mixed Levels

- **Protocol is universal** — L0, L1, L2 teams collaborate seamlessly
- **Hooks are local** — Each developer's machine has its own hooks
- **CLI is optional** — Not required for protocol compliance
- **`.sdp/config.json`** declares enabled capabilities

### Backward Compatibility

| From | To | Breaking? | Migration |
|------|-----|-----------|-----------|
| L0 | L1 | No | Just install hooks |
| L1 | L2 | No | Just build CLI |
| L2 → L1 | No | Just stop using CLI |
| Any | L0 | No | Just ignore tooling |

---

## Implementation Plan

### Phase 1: Foundation (Current - v0.10.x)

- [x] L0 Protocol documented in PROTOCOL.md
- [x] L1 Hooks implemented (pre-commit, pre-push, PreToolUse)
- [x] L2 Evidence logging (hash-chained events.jsonl)
- [x] L2 Guard commands (activate, check, status)
- [x] L2 Checkpoint system
- [ ] **TODO**: L0-L1 separation (remove tool-specific refs from PROTOCOL.md)
- [ ] **TODO**: Skill frontmatter standardization

### Phase 2: Reliability (v0.11.x)

- [ ] Dual-Log Architecture (separate state.jsonl)
- [ ] EvidenceEmitter interface for L1-L2 contract
- [ ] Checkpoint events in evidence log
- [ ] Guard policy file (`.sdp/guard-policy.yml`)
- [ ] GitHub Action for CI drift detection
- [ ] Schema registry (`src/sdp/schema/`)

### Phase 3: Orchestration (v0.12.x)

- [ ] AgentBackend interface
- [ ] LocalBackend (wraps existing ExecuteFunc)
- [ ] BackendRegistry
- [ ] Per-backend circuit breakers
- [ ] HTTPBackend for remote agents

### Phase 4: Collaboration (v0.13.x)

- [ ] Enhanced Beads schema (claimed_by, claimed_at)
- [ ] Notification triggers on AI actions
- [ ] Conflict detection (claimed_by check)
- [ ] Webhook channel completion
- [ ] Desktop channel enablement

### Phase 5: Scale (v1.0.x)

- [ ] gRPCBackend for streaming
- [ ] k8s deployment manifests
- [ ] Agent discovery service
- [ ] Multi-region support

---

## Success Metrics

| Metric | Baseline | Target (v1.0) |
|--------|----------|---------------|
| L0 adoption time | 10 min | < 5 min |
| L1 hook effectiveness | Unknown | 90% of violations caught |
| L2 evidence coverage | Partial | 100% of WS executions |
| L3 parallel speedup | 1x | 5x for 5+ WS |
| L4 notification latency | N/A | < 30 sec |
| Mixed-level compatibility | Untested | 100% |

---

## See Also

- [docs/PROTOCOL.md](../PROTOCOL.md) — Current protocol specification
- [docs/vision/2026-02-08-sdp-multi-agent-future-design.md](../vision/2026-02-08-sdp-multi-agent-future-design.md) — Future design doc
- [CLAUDE.md](../../CLAUDE.md) — Integration guide
- [.claude/commands.json](../../.claude/commands.json) — Command mappings

---

**Version:** 1.0.0 (Multi-Level Architecture Design)
**Authors:** Expert Analysis Team (8 specialists)
**Date:** 2026-02-13
