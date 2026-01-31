# Beads + SDP Integration Design

> **Status:** Research complete
> **Date:** 2026-01-28
> **Goal:** Evaluate integration between SDP (workstream framework) and Beads (git-backed issue tracker), and design unified workflow for task management

---

## Table of Contents

1. [Overview](#overview)
2. [Executive Summary](#executive-summary)
3. [What is Beads?](#what-is-beads)
4. [Integration Analysis](#1-architecture-separation)
5. [Implementation Plan](#implementation-plan)
6. [Recommendations](#recommendations)

---

## Overview

### Context

- **SDP** (Spec-Driven Protocol): Workstream-driven AI development framework with TDD discipline, quality gates, and skills-based execution
- **Beads**: Git-backed distributed issue tracker optimized for AI agents, written in Go with SQLite storage and JSONL persistence
- **Location**: `/Users/fall_out_bug/projects/vibe_coding/beads`

### Goals

1. Understand Beads architecture and capabilities
2. Evaluate integration potential with SDP
3. Design separation of concerns between both systems
4. Propose implementation path with minimal disruption

---

## Executive Summary

### Key Finding: HIGH Integration Potential 🚀

Beads and SDP are **highly complementary**:

- **Beads excels at**: Issue tracking, dependency graphs, multi-agent coordination, git-backed storage
- **SDP excels at**: Workstream execution, TDD discipline, quality gates, AI skills

**Natural synergy**: Beads manages **what** needs to be done; SDP manages **how** to do it.

### Recommended Approach: **Option B - SDP as Source of Truth, Beads as Optional Cache**

```
SDP Markdown Files (source of truth)
         ↓
    SDP Skills (@idea, @build, @review)
         ↓
   Optional Beads Sync (read-only projection)
         ↓
    Beads CLI for queries/visualization
```

**Why this approach:**
- ✅ Preserves SDP's markdown-first philosophy
- ✅ No hard dependency on Beads (graceful degradation)
- ✅ Incremental adoption (opt-in sync)
- ✅ Git resolves conflicts (human-readable markdown)

### Quick Facts About Beads

| Metric | Value |
|--------|-------|
| **Language** | Go (1.24+) |
| **Storage** | SQLite + JSONL (git-backed) |
| **ID Format** | Hash-based (`bd-a3f8`) prevents conflicts |
| **Architecture** | CLI → SQLite → JSONL → Git |
| **Python LOC** | 10,991 (SDP integration layer) |
| **Documentation** | 186 Markdown files |
| **Maturity** | Production-ready with MCP servers |

---

## What is Beads?

### Core Features

**1. Git-Backed Issue Tracking**
- Issues stored as JSONL in `.beads/` directory
- Version control friendly (line-by-line history)
- No merge conflicts (hash-based IDs)

**2. Graph-Based Dependencies**
- Parent-child relationships: `bd-a3f8`, `bd-a3f8.1`, `bd-a3f8.1.1`
- Dependency types: `blocks`, `parent-child`, `related`, `discovered-from`
- Cycle detection and topological sort

**3. Multi-Agent Optimized**
- Hash-based IDs prevent concurrent creation conflicts
- Daemon mode for background auto-sync (30s debounce)
- JSON output for machine readability
- "Ready" detection (tasks with no open blockers)

**4. Status Management**
- States: `open`, `in_progress`, `blocked`, `deferred`, `closed`, `tombstone`, `pinned`, `hooked`
- Priority levels: 0-4 (0=critical, 4=backlog)
- Automatic status transitions

**5. Integration Support**
- Jira sync (bidirectional)
- GitHub sync
- MCP server for AI agents
- Community tools: `beads_viewer` with robot insights

### Project Scale

```
beads/
├── cmd/bd/              # CLI commands (Go)
├── internal/
│   ├── types/           # Core data types
│   └── storage/         # SQLite + JSONL
├── src/sdp/             # SDP integration (Python)
└── docs/                # 186 Markdown files
```

**Already has SDP integration layer!** → `src/sdp/` directory exists in Beads repo.

---

## Integration Analysis

### 1. Architecture Separation

> **Experts:** Sam Newman, Martin Fowler, Robert C. Martin

**Decision: Option A - Beads = Issue Storage + Graph, SDP = Execution + Quality**

| System | Responsibilities |
|--------|-----------------|
| **Beads** | Issue storage, dependency graph, git operations, status transitions, graph queries |
| **SDP** | TDD workflow, quality gates, WS tier validation, platform adapters, PRD generation |

**Rationale:**
- Clear bounded contexts (Newman's principle)
- Single responsibility per system
- Loose coupling via stable CLI API

**Integration Point:**
```bash
# SDP calls Beads CLI
beads create-issue --json "{'id': 'WS-001-01', 'deps': []}"
beads get-ready-tasks --graph .beads/graph.jsonl
beads update-status --id WS-001-01 --status completed
```

---

### 2. Data Model Mapping

> **Experts:** Martin Fowler, Eric Evans, Martin Kleppmann

**Decision: Option C - Bidirectional Sync with Mapping Layer**

**Field Mapping:**
```
SDP Workstream                    Beads Issue
────────────────────────────────────────────────
ws_id: "00-012-06"           ↔    id: "bd-a3f8"
status: backlog/active/...   ↔    status: open/in_progress/...
dependencies: ["00-012-05"]  ↔    dependencies: [{type: "blocks", id: "bd-a1b2"}]
size: SMALL/MEDIUM/LARGE     ↔    priority: 2/1/0 + custom label
acceptance_criteria          ↔    acceptance_criteria (string)
feature: F012                ↔    parent (via parent-child dep)
```

**ID Mapping Table:**
```json
// .beads-sdp-mapping.jsonl
{"sdp_id": "00-012-06", "beads_id": "bd-a3f8", "updated_at": "2026-01-28T15:30:00Z"}
```

**Conflict Resolution:**
- SDP = Authoritative for content (title, goal, acceptance criteria)
- Beads = Authoritative for execution metadata (status changes, comments)

---

### 3. Workflow Integration

> **Experts:** Kelsey Hightower, Sam Newman, Theo Browne

**Decision: Option B - SDP as Source of Truth with Optional Beads Sync**

**User Experience:**
```bash
# Standard SDP workflow (no Beads needed)
@idea "Add user auth"        # Creates draft in docs/drafts/
@design idea-auth            # Creates workstreams in docs/workstreams/
@build WS-001-01            # Updates markdown, runs TDD

# Beads visibility (if installed)
bd list --project=00 --feature=001
# → Shows WS-001-01 as "completed" with SDP metadata
```

**Hook Integration:**
```bash
# hooks/post-build.sh (existing)
run_quality_gates "$WS_ID"
append_execution_report "$WS_ID"

# NEW: Optional Beads sync
if command -v bd &> /dev/null; then
  bd sync-workstream "$WS_FILE" || log_warning "Beads sync failed"
fi
```

**Progressive Disclosure:**
- Default: SDP works standalone
- Opt-in: `export BEADS_INTEGRATION=1` enables sync
- Power users: Use `bd` CLI for advanced queries

---

### 4. State Synchronization

> **Experts:** Martin Kleppmann, Markus Winand, Sam Newman

**Decision: Option D - Hybrid Consensus Protocol**

**Architecture:**
```
Primary Source: WS Markdown Files (human-readable, Git-tracked)
         ↓
Append-Only Log: JSONL export (`.beads-sdp/artifacts/tasks.jsonl`)
         ↓
Reconciliation: status.json maps WS IDs → JSONL line numbers
```

**Conflict Resolution:**
- Git conflict in WS file → Human resolves (standard git merge)
- JSONL diverges → Append-only, no conflicts
- status.json mismatch → Regenerate from WS + JSONL

**Why not pure SQLite?**
- Preserves SDP's markdown-first philosophy
- Git is the conflict resolver (human-readable)
- JSONL is append-only (no merge conflicts)

---

### 5. CLI/UX Integration

> **Experts:** Nir Eyal, Jakob Nielsen, Martin Fowler

**Decision: Option B - Dual Interface with Smart Routing**

**Two Interfaces:**
```bash
# AI Workflows (@skills) - Interactive, exploratory
@idea "Add auth"        # Deep interviewing
@design idea-auth       # Architecture decisions
@issue "Bug report"     # Debugging

# Operations (sdp CLI) - Scriptable, deterministic
sdp validate WS-001-01  # Quality gates
sdp metrics tier        # Metrics
sdp github sync         # GitHub operations
```

**Smart Routing:**
```bash
# Terminal users can trigger AI workflows
sdp workflow idea "Add auth"  # Wraps @idea for terminal

# AI users can run operations
sdp validate WS-001-01        # Works in Claude Code too
```

**Clear Decision Tree:**
```
Need human decisions? → @skills
Need automation? → sdp CLI
Terminal user? → sdp workflow <cmd>
```

---

### 6. Git Storage Strategy

> **Experts:** Markus Winand, Kelsey Hightower, Sam Newman

**Decision: Hybrid - Keep Separate, Link via References**

**Storage Layout:**
```
project/
├── .beads/                    # Beads data (JSONL + SQLite)
│   └── issues.jsonl
├── docs/workstreams/          # SDP workstreams (markdown)
│   ├── backlog/
│   ├── in_progress/
│   └── completed/
├── .beads-sdp-mapping.jsonl   # ID mapping table
└── .beads-sdp/
    └── artifacts/
        └── tasks.jsonl        # Append-only audit log
```

**Why Separate:**
- Beads needs SQLite for performance (SDP doesn't)
- SDP needs markdown for humans (Beads doesn't)
- Each system optimized for its use case
- Git conflict resolution works independently

---

### 7. Multi-Agent Coordination

> **Experts:** Martin Kleppmann, Kelsey Hightower, Sam Newman

**Decision: Option B - Hybrid Coordination (Lock + Optimistic Concurrency)**

**Approach:**
- **Critical Path**: Lock-based queue management (enqueue/dequeue)
- **Workstream Files**: Version stamps in YAML frontmatter
- **Conflict Detection**: Compare versions before commit
- **Graceful Degradation**: If locks fail, use optimistic mode

**Implementation:**
```python
# Workstream file with version stamp
---
ws_id: 00-012-06
version: 5  # Increment on every write
last_modified: "2026-01-28T15:30:00Z"
feature: F012
status: backlog
---

# Optimistic update logic
def update_workstream(ws_id: str, update_fn: Callable) -> bool:
    for attempt in range(3):
        ws = read_workstream(ws_id)
        old_version = ws.version
        updated_ws = update_fn(ws)
        updated_ws.version += 1
        if write_if_version_match(ws_id, updated_ws, old_version):
            return True
    raise ConcurrentModificationError(f"Failed after 3 attempts: {ws_id}")
```

**Beads Hash IDs (Optional Enhancement):**
- Keep PP-FFF-SS for human readability
- Add hash suffix for concurrent creation: `PP-FFF-SS-{hash}`
- Human tools display only `PP-FFF-SS`, internal tools use full ID

---

### 8. Migration Path

> **Experts:** Martin Fowler, Kelsey Hightower, Sam Newman

**Decision: Option B - Bidirectional Synchronization Adapter (Strangler Fig Pattern)**

**Three-Phase Migration:**

**Phase 1: Export-Only (Week 1-2)**
```python
# src/sdp/beads_exporter.py
def workstream_to_beads(ws: Workstream) -> dict:
    hash_input = f"{ws.ws_id}:{ws.title}"
    beads_id = f"bd-{hashlib.sha256(hash_input.encode()).hexdigest()[:4]}"
    return {
        "id": beads_id,
        "title": f"{ws.ws_id}: {ws.title}",
        "status": _map_status(ws.status),
        "external_ref": f"PP-FFF-SS:{ws.ws_id}",
        ...
    }

# Export all workstreams to .beads/issues.jsonl
export_to_beads("docs/workstreams/", ".beads/issues.jsonl")
```

**Phase 2: Bidirectional Sync (Week 3-4)**
- Read JSONL changes, update markdown frontmatter
- Read markdown changes, regenerate JSONL
- Conflict resolution: markdown wins (source of truth)

**Phase 3: Cutover (Week 5-6)**
- Make JSONL source of truth (optional)
- Archive markdown to `docs/workstreams.legacy/`
- Update @build skill to use Beads CLI directly

**Rollback-Safe:**
- Each phase independently releasable and reversible
- Can stop at Phase 1 (export-only) if that's sufficient
- Git preserves all history

---

### 9. Quality Gate Integration

> **Experts:** Kelsey Hightower, Theo Browne, Troy Hunt

**Decision: Option A - Unidirectional (SDP → Beads)**

**Implementation:**
```bash
# hooks/post-build.sh (after line 204)
if [ "$BEADS_INTEGRATION" = "1" ]; then
    echo ""
    echo "Check 11: Beads sync"
    if command -v bd &> /dev/null; then
        if bd done "$WS_ID" 2>&1; then
            echo "✓ Beads status updated to done"
        else
            echo "⚠️ Beads sync failed (non-critical)"
        fi
    else
        echo "  Skipped (bd command not found)"
    fi
fi
```

**Why Unidirectional:**
- ✅ Quality gates block bad code from being committed (SDP enforcement)
- ✅ Beads `done` status only updated when gates pass (double verification)
- ✅ Minimal coupling (Beads integration is optional)
- ✅ Fail-fast with explicit errors

**Configuration:**
```bash
export BEADS_INTEGRATION=1  # Enable Beads sync
```

---

### 10. Documentation Strategy

> **Experts:** Diátaxis Framework, Martin Fowler, Kelsey Hightower

**Decision: Unified Documentation with Cross-References**

**Documentation Structure:**
```
docs/
├── integration/
│   ├── beads.md           # Beads integration guide
│   ├── workflows.md       # Unified workflows (@skills + sdp)
│   └── migration.md       # Migration guide (3 phases)
├── CLAUDE.md              # Updated with Beads section
└── README.md              # Overview with both systems
```

**Content Sections:**
1. **Conceptual** - Why integrate Beads + SDP?
2. **Tutorials** - Getting started with integrated workflow
3. **How-To Guides** - Common tasks (sync, migrate, troubleshoot)
4. **Reference** - CLI commands, configuration options
5. **Explanation** - Architecture, trade-offs, design decisions

**Quick Start:**
```bash
# 1. Install Beads
go install github.com/steveyegge/beads/cmd/bd@latest

# 2. Enable integration
export BEADS_INTEGRATION=1

# 3. Use SDP as usual
@idea "Add auth"
@build WS-001-01  # Automatically syncs to Beads
```

---

## Implementation Plan

### Phase 1: Foundation (Week 1-2)

**Goals:**
- Basic Beads integration without disrupting existing SDP workflows
- Export-only sync (SDP → Beads)

**Tasks:**
- [ ] Create `src/sdp/beads/exporter.py` - WS to Beads JSONL converter
- [ ] Add `sdp beads export` CLI command
- [ ] Implement ID mapping table (`.beads-sdp-mapping.jsonl`)
- [ ] Add pre-commit hook to regenerate JSONL on WS changes
- [ ] Write integration tests for exporter
- [ ] Document export-only workflow

**Deliverables:**
- Exporter script with field mapping
- ID mapping table format
- Pre-commit hook integration
- Documentation (Phase 1)

---

### Phase 2: Bidirectional Sync (Week 3-4)

**Goals:**
- Two-way sync between SDP and Beads
- Conflict detection and resolution

**Tasks:**
- [ ] Create `src/sdp/beads/sync_service.py`
- [ ] Implement `sync_workstream_to_beads()` method
- [ ] Implement `sync_beads_to_workstream()` method
- [ ] Add conflict detection (version stamps, timestamps)
- [ ] Implement "markdown wins" conflict resolution
- [ ] Add `sdp beads sync` CLI command
- [ ] Extend post-build hook for auto-sync
- [ ] Write tests for sync logic

**Deliverables:**
- Sync service with bidirectional support
- CLI commands for manual sync
- Auto-sync in post-build hook
- Documentation (Phase 2)

---

### Phase 3: Enhanced Features (Week 5-6)

**Goals:**
- Optional cutover to Beads-first mode
- Advanced query capabilities

**Tasks:**
- [ ] Add `sdp beads query` command (uses Beads CLI)
- [ ] Implement "ready task detection" via Beads
- [ ] Add optional Beads-first mode (JSONL source of truth)
- [ ] Create migration script for Beads-first
- [ ] Update @build skill to use Beads when available
- [ ] Add telemetry (sync success rate, performance)
- [ ] Performance optimization (caching, async sync)

**Deliverables:**
- Query commands using Beads
- Optional Beads-first mode
- Migration scripts
- Documentation (Phase 3)

---

### Phase 4: Polish & Documentation (Week 7-8)

**Goals:**
- Production-ready integration
- Comprehensive docs and examples

**Tasks:**
- [ ] End-to-end testing (multi-agent scenarios)
- [ ] Performance benchmarking (sync latency, overhead)
- [ ] Troubleshooting guide (common issues, fixes)
- [ ] Video tutorial: "SDP + Beads in 10 minutes"
- [ ] Example projects (with Beads integration)
- [ ] Update README.md with Beads section
- [ ] Update CLAUDE.md with Beads workflows
- [ ] Release notes for v0.5.0

**Deliverables:**
- Production-ready integration
- Complete documentation suite
- Video tutorial
- Example projects

---

## Recommendations

### For SDP Maintainers

**1. Adopt Option B (SDP Source of Truth, Beads Optional Cache)**
- Preserves SDP's markdown-first philosophy
- No hard dependency on Beads
- Incremental adoption path

**2. Start with Phase 1 (Export-Only)**
- Low risk, high value
- Can stop here if sufficient
- Validates integration assumption

**3. Implement Quality Gate Integration Early**
- Unidirectional SDP → Beads sync
- Add `BEADS_INTEGRATION` flag to post-build hook
- Blocks "done" status until gates pass

**4. Keep Systems Loosely Coupled**
- CLI API integration (not Python imports)
- Each system can evolve independently
- Graceful degradation if Beads missing

---

### For Beads Users

**1. Try SDP for Spec-Driven Development**
- Use @idea for requirements gathering
- Use @design for workstream planning
- Use @build for TDD execution

**2. Leverage Beads for Query & Visualization**
- Use `bd list` for task overview
- Use `bd ready` for next tasks
- Use `bd insights` for hotspots

**3. Provide Feedback**
- Does sync work reliably?
- Is performance acceptable?
- What's missing?

---

### For Teams Considering Both

**Key Questions:**

| Question | Answer |
|----------|--------|
| **Is Beads required?** | No - SDP works standalone |
| **Does sync block builds?** | No - Beads failures are non-critical |
| **Can we migrate gradually?** | Yes - 3-phase migration |
| **What if Beads changes?** | SDP continues working (markdown source of truth) |
| **Performance impact?** | Minimal - async sync, optional caching |

---

## Success Metrics

| Metric | Baseline | Target | How to Measure |
|--------|----------|--------|----------------|
| **Sync success rate** | N/A | ≥ 99% | Post-build hook telemetry |
| **Sync latency** | N/A | < 5s | Time from build complete to JSONL updated |
| **Build overhead** | 0s | < 1s | Beads sync time in post-build hook |
| **Adoption rate** | 0% | 50% of projects | `BEADS_INTEGRATION=1` in projects |
| **User satisfaction** | N/A | ≥ 4.5/5 | Survey after 3 months |

---

## Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| **Sync diverges** | High | Medium | ID mapping table, checksum validation, manual repair command |
| **Beads API changes** | High | Low | Versioned CLI API, adapter pattern, automated tests |
| **Performance degradation** | Medium | Low | Async sync, caching, optional integration |
| **User confusion** | Medium | High | Clear documentation, decision trees, examples |
| **Multi-agent conflicts** | High | Medium | Optimistic concurrency, version stamps, hash IDs (Phase 3) |

---

## Open Questions

1. **Should Beads be optional or recommended?**
   - Current design: Optional (export flag)
   - Future: Consider recommending for multi-agent teams

2. **Should we support Beads-first mode?**
   - Current design: SDP-first (markdown source of truth)
   - Phase 3: Optional Beads-first for power users

3. **How to handle schema drift?**
   - Current design: Ignore unknown fields
   - Future: Version mapping schema, validation

4. **Should we use Beads hash IDs?**
   - Current design: Keep PP-FFF-SS format
   - Phase 3: Optional hash suffix for concurrent creation

---

## Next Steps

1. **Review this design** with SDP and Beads maintainers
2. **Decide on Phase 1 scope** (export-only or bidirectional?)
3. **Create proof-of-concept** for exporter script
4. **Test with real project** (migrate existing SDP workstreams)
5. **Measure performance** (sync latency, overhead)
6. **Gather feedback** from early adopters
7. **Iterate on design** based on learnings

---

## Conclusion

Beads and SDP are **highly complementary** systems that can work together seamlessly:

- **Beads** = Issue tracking, dependency graphs, multi-agent coordination
- **SDP** = Workstream execution, TDD discipline, quality gates

The recommended integration approach (Option B: SDP source of truth, Beads optional cache) provides:

- ✅ Minimal disruption to existing workflows
- ✅ Incremental adoption path (3 phases)
- ✅ Graceful degradation (works without Beads)
- ✅ Clear separation of concerns
- ✅ Production-ready migration strategy

**Ready to implement?** Start with Phase 1 (export-only) and validate assumptions before committing to bidirectional sync.

**Want to discuss further?** Open a GitHub issue or discussion to refine the design.

---

**Sources:**
- [Beads GitHub Repository](https://github.com/steveyegge/beads)
- [Beads ARCHITECTURE.md](https://github.com/steveyegge/beads/blob/main/docs/ARCHITECTURE.md)
- [Beads FAQ Documentation](https://github.com/steveyegge/beads/blob/main/docs/FAQ.md)
- [Introducing Beads - Steve Yegge](https://steve-yegge.medium.com/introducing-beads-a-coding-agent-memory-system-637d7d92514a)
- [Beads: A Git-Friendly Issue Tracker for AI Agents](https://betterstack.com/community/guides/ai/beads-issue-tracker-ai-agents/)
- [SDP PROTOCOL.md](../PROTOCOL.md)
- [SDP CLAUDE.md](../CLAUDE.md)

---

**Version:** 1.0
**Authors:** SDP + Beads Expert Analysis Team
**Status:** Ready for Review
