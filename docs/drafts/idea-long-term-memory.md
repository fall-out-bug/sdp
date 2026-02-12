# F051: SDP Long-term Memory System

> **Status:** Draft
> **Created:** 2026-02-12
> **Beads:** sdp-6fx

## Problem Statement

Throughout a project, decisions and changes accumulate across many artifacts and files. Agents often:
- Duplicate already-completed research
- Propose previously rejected approaches
- Are unaware of existing architectural decisions

**Goal:** Provide quick access to decision history, their rationale, and current project state.

## Goals

1. **Avoid duplicated work** - agent should know what's already been researched
2. **Preserve decisions** - remember what was accepted and what was rejected
3. **Detect drift** - discover discrepancies between code and documentation

## Non-Goals

- Global memory across projects (project-level only)
- Real-time synchronization between agents (event sourcing is sufficient)
- ML training on history (out of scope)

## Technical Approach

### Integration with Existing Subsystems

| Subsystem | Role | Integration |
|-----------|------|-------------|
| evidence.jsonl | Event log with hash chain | Source of truth for events |
| Beads issues | Persistent tasks | Classification and status |
| .sdp/session.json | Session state (F065) | Context recovery |

### Hybrid Search

1. **Full-text** - keyword search across artifacts
2. **Semantic** - embeddings for meaning-based retrieval
3. **Graph** - relationships between decisions, features, files

### Drift Detection (all types)

| Type | Description | Trigger |
|------|-------------|---------|
| Code↔Docs | Code diverges from documentation | `sdp drift detect` |
| Decisions↔Code | Implementation doesn't match decisions | ADR validation |
| Docs↔Docs | Contradictions between documents | Cross-reference check |

### Parallel Tracks

- **Memory/Search** - indexing, storage, retrieval
- **Agent Coordination** - event sourcing, role switching

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Long-term Memory System                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │   Memory     │    │    Search    │    │    Drift     │       │
│  │   Store      │    │    Engine    │    │   Detector   │       │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘       │
│         │                   │                   │                │
│         └───────────────────┼───────────────────┘                │
│                             │                                    │
│                    ┌────────▼────────┐                          │
│                    │   Event Store   │                          │
│                    │ (evidence.jsonl)│                          │
│                    └────────┬────────┘                          │
│                             │                                    │
│  ┌──────────────┐    ┌──────▼───────┐    ┌──────────────┐       │
│  │    Beads     │    │   Session    │    │   Agent      │       │
│  │   Issues     │    │   Manager    │    │ Coordinator  │       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## User Stories

### US1: Context Recovery After Compaction
**As** an agent after session compaction
**I want** to quickly restore project context
**So that** I can continue work without re-researching

### US2: Decision Discovery
**As** an agent planning a feature
**I want** to find related architectural decisions
**So that** I don't propose rejected approaches

### US3: Drift Detection
**As** a developer before merge
**I want** to detect code-documentation discrepancies
**So that** documentation stays up-to-date

### US4: Agent Coordination
**As** an orchestrator agent
**I want** to coordinate multiple agents' work
**So that** conflicts and duplication are avoided

## Acceptance Criteria

### Memory & Search
- [ ] AC1: Indexed artifacts include all .md files in docs/
- [ ] AC2: Semantic search returns relevant results (precision > 80%)
- [ ] AC3: Query response time < 500ms for project with 1000+ artifacts

### Drift Detection
- [ ] AC4: Code↔Docs drift detected with file:line references
- [ ] AC5: Decision drift detected via ADR validation
- [ ] AC6: Drift report generated on demand

### Agent Coordination
- [ ] AC7: Event-sourced coordination via evidence.jsonl
- [ ] AC8: Dynamic role switching based on task type
- [ ] AC9: Continuous requirement verification

## Concerns & Tradeoffs

### Tradeoffs

| Decision | Choice | Alternative | Rationale |
|----------|--------|-------------|-----------|
| Storage | evidence.jsonl extension | New database | Less complexity, existing hash chain |
| Search | SQLite FTS5 + embeddings | Elasticsearch | No external dependency |
| Coordination | Event sourcing | Distributed state | Single source of truth |

### Risks

1. **Embeddings cost** - Semantic search requires API calls
   - Mitigation: Cache embeddings, batch processing

2. **Index size growth** - Large projects may have many artifacts
   - Mitigation: Incremental indexing, compression

3. **Drift false positives** - Not all changes are drift
   - Mitigation: Configurable thresholds, whitelist

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| Context recovery time | 5-10 min | < 30 sec |
| Duplicate research | ~30% | < 10% |
| Drift detection recall | N/A | > 90% |
| Agent coordination latency | N/A | < 100ms |

## Related Beads Tasks

- sdp-6fx.1: Add @deploy to @oneshot workflow ✅
- sdp-6fx.2: Add checkpoint to @feature before agents
- sdp-6fx.3: Clarify @feature vs @oneshot in documentation
- sdp-6fx.4: Refactor @feature to orchestrate @idea/@design
- sdp-6fx.5: Implement event-sourced agent coordination
- sdp-6fx.6: Implement dynamic agent role switching
- sdp-6fx.7: Implement continuous requirement verification
- sdp-6fx.8: Enhance checkpoint with snapshots and rollback
- sdp-6fx.9: Implement parallel workstream execution
- sdp-6fx.10: Create notification gateway system
- sdp-6fx.11: Implement drift detection system

---

**Next Step:** @design idea-long-term-memory
