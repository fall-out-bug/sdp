# F052: Multi-Agent SDP + @vision + @reality - Complete Workstreams Summary

> **Total Workstreams:** 25
> **Timeline:** ~13 weeks (3 months)
> **Approach:** Quality Lock-in (Phase 2 first) → Parallel tracks

## Parent Feature

**sdp-db2w** - Feature F052: Multi-Agent SDP + @vision + @reality (P1)

## Critical Path

```
00-052-00 (Backup)
  ↓
00-052-01, 00-052-02 (Phase 1A/B - parallel)
  ↓
00-052-08 (Implementer) + 00-052-09 (Spec Reviewer) - parallel
  ↓
00-052-10 (@build update)
  ↓
00-052-11 (Test two-stage) → READY FOR PARALLEL PHASES
  ↓
Phase 3/4/5 execute in parallel
  ↓
Phase 6 (Documentation)
```

---

## Phase 0: Preparation (1 day)

### 00-052-00: Backup & Worktree Setup
- **Beads ID:** sdp-wqv8
- **Size:** SMALL | **Duration:** 1 day
- **Dependencies:** None (MUST BE FIRST)
- **Status:** ✅ Workstream created
- **Beads Status:** ⚠️ Has cyclic dependency issue (depends on parent, but should be independent)

**File:** `docs/workstreams/backlog/00-052-00-backup-and-worktree.md` ✅ CREATED

---

## Phase 1A: @vision Skill (1 week, parallel with Phase 1B)

### 00-052-01: @vision Skill Structure
- **Beads ID:** sdp-vxvp
- **Size:** MEDIUM | **Duration:** 2-3 days
- **Dependencies:** 00-052-00
- **Acceptance Criteria:**
  - AC1: `.claude/skills/vision/SKILL.md` created
  - AC2: Interview workflow (3-5 questions via AskUserQuestion)
  - AC3: Deep-thinking with 7 expert agents
  - AC4: Artifact generation (PRODUCT_VISION.md, PRD.md, ROADMAP.md)

**File:** `docs/workstreams/backlog/00-052-01-vision-skill-structure.md` ✅ CREATED

### 00-052-02: Vision Extractor Implementation
- **Beads ID:** sdp-wbyc
- **Size:** MEDIUM | **Duration:** 2-3 days
- **Dependencies:** 00-052-01
- **Scope:**
  - `src/sdp/vision/extractor.py` - Extract P0/P1 features from PRD
  - `tests/sdp/vision/test_extractor.py`
- **AC:** Parse PRD, extract features, create drafts in `docs/drafts/`

### 00-052-03: Update CLAUDE.md with @vision
- **Beads ID:** sdp-luo8
- **Size:** SMALL | **Duration:** 1-2 days
- **Dependencies:** 00-052-01
- **Scope:** `CLAUDE.md` - Add @vision to decision tree

---

## Phase 1B: @reality Skill (1 week, parallel with Phase 1A)

### 00-052-04: @reality Skill Structure
- **Beads ID:** sdp-uugq
- **Size:** MEDIUM | **Duration:** 2-3 days
- **Dependencies:** 00-052-00
- **AC:**
  - `.claude/skills/reality/SKILL.md` created
  - Modes: --quick (5-10 min), --deep (30-60 min), --focus=topic
  - 8 expert agents for deep analysis
  - Universal project analysis (works without SDP)

### 00-052-05: Project Scanner Implementation
- **Beads ID:** sdp-bsds
- **Size:** MEDIUM | **Duration:** 2-3 days
- **Dependencies:** 00-052-04
- **Scope:**
  - `src/sdp/reality/scanner.py` - ProjectScanner class
  - `src/sdp/reality/detectors.py` - Language/framework detection
  - `tests/sdp/reality/test_scanner.py`

### 00-052-06: Update CLAUDE.md with @reality
- **Beads ID:** sdp-kxif
- **Size:** SMALL | **Duration:** 1-2 days
- **Dependencies:** 00-052-04
- **Scope:** `CLAUDE.md` - Add @reality examples

### 00-052-07: Test @vision + @reality Integration
- **Beads ID:** sdp-riqq
- **Size:** SMALL | **Duration:** 1-2 days
- **Dependencies:** 00-052-02, 00-052-05
- **AC:** E2E test on sample projects, verify artifact generation

---

## Phase 2: Two-Stage Review (QUALITY LOCK-IN - 2 weeks)

**CRITICAL:** This phase MUST complete before Phases 3-6.

### 00-052-08: Implementer Agent
- **Beads ID:** sdp-vrij
- **Size:** MEDIUM | **Duration:** 3-4 days
- **Dependencies:** 00-052-00
- **Scope:** `.claude/agents/implementer.md`
- **AC:**
  - TDD cycle spec (Red → Green → Refactor)
  - Self-report format
  - Quality check before commit

### 00-052-09: Spec Compliance Reviewer Agent
- **Beads ID:** sdp-01q4
- **Size:** MEDIUM | **Duration:** 3-4 days
- **Dependencies:** 00-052-00
- **Scope:** `.claude/agents/spec-reviewer.md`
- **AC:**
  - "DO NOT TRUST" pattern
  - Read actual code, don't trust implementer report
  - Spec vs reality comparison

### 00-052-10: Update @build for Two-Stage Review
- **Beads ID:** sdp-syga
- **Size:** LARGE | **Duration:** 4-5 days
- **Dependencies:** 00-052-08, 00-052-09
- **Scope:** `.claude/skills/build/SKILL.md`
- **AC:**
  - Orchestrate: Implementer → Spec Reviewer → Quality Reviewer
  - Max 2 retries per stage
  - All stages pass → mark complete

### 00-052-11: Test Two-Stage Review End-to-End
- **Beads ID:** sdp-n7sp
- **Size:** MEDIUM | **Duration:** 2-3 days
- **Dependencies:** 00-052-10
- **AC:** E2E test on sample workstream, verify all 3 stages execute

---

## Phase 3: Speed Track (3 weeks, parallel with Phase 4-5)

**Dependencies:** Phase 2 complete

### 00-052-12: Parallel Dispatcher for @oneshot
- **Beads ID:** sdp-tsc5
- **Size:** LARGE | **Duration:** 5-7 days
- **Dependencies:** 00-052-11
- **Scope:** `.claude/agents/orchestrator.md`, `.claude/skills/oneshot/SKILL.md`
- **AC:**
  - Build dependency graph from WS files
  - Kahn's algorithm for topological sort
  - Execute ready WS in parallel (3-5 agents)
  - 3x speedup for 5+ WS

### 00-052-13: Circuit Breaker Implementation
- **Beads ID:** sdp-fd04
- **Size:** MEDIUM | **Duration:** 4-5 days
- **Dependencies:** 00-052-11
- **Scope:** `.claude/agents/circuit-breaker.md`
- **AC:**
  - Error categorization: transient/permanent/cascade/crash
  - Retry with exponential backoff (transient)
  - Skip WS (permanent) or stop feature (cascade)

### 00-052-14: Checkpoint Atomic Writes
- **Beads ID:** sdp-8wf8
- **Size:** SMALL | **Duration:** 2-3 days
- **Dependencies:** 00-052-13
- **Scope:** `.claude/agents/orchestrator.md`
- **AC:**
  - Temp file + rename pattern
  - Prevents corruption on crash

### 00-052-15: Test Parallel Execution
- **Beads ID:** sdp-n856
- **Size:** MEDIUM | **Duration:** 3-4 days
- **Dependencies:** 00-052-12, 00-052-14
- **AC:** Test 3-5 builder agents parallel, verify checkpoint updates

---

## Phase 4: Synthesis Track (2 weeks, parallel with Phase 3-5)

**Dependencies:** Phase 2 complete

### 00-052-16: Agent Synthesizer Core
- **Beads ID:** sdp-ychs
- **Size:** MEDIUM | **Duration:** 4-5 days
- **Dependencies:** Phase 2 complete
- **Scope:** `.claude/agents/synthesizer.md`
- **AC:**
  - Synthesis rules (priority order)
  - Unanimous agreement → domain expertise → quality gate → merge → escalate

### 00-052-17: Synthesis Rules Engine
- **Beads ID:** sdp-89ug
- **Size:** MEDIUM | **Duration:** 3-4 days
- **Dependencies:** 00-052-16
- **Scope:** `src/sdp/synthesis/rules.py`, `src/sdp/synthesis/engine.py`
- **AC:**
  - Modular rule system
  - Priority-based execution
  - Tests for each rule

### 00-052-18: Hierarchical Supervisor
- **Beads ID:** sdp-8r2n
- **Size:** MEDIUM | **Duration:** 3-4 days
- **Dependencies:** 00-052-17
- **Scope:** `.claude/agents/supervisor.md`
- **AC:**
  - Coordinate specialist agents
  - Apply synthesizer
  - Escalate to human if needed

---

## Phase 5: UX Track (2 weeks, parallel with Phase 3-4)

**Dependencies:** Phase 2 complete

### 00-052-19: Progressive Disclosure for @idea
- **Beads ID:** sdp-62ts
- **Size:** MEDIUM | **Duration:** 4-5 days
- **Dependencies:** Phase 2 complete
- **Scope:** `.claude/skills/idea/SKILL.md`
- **AC:**
  - 3-question cycles with trigger points
  - TMI detection (suggest --quiet)
  - 12-27 questions per feature (down from unbounded)

### 00-052-20: Progressive Disclosure for @design
- **Beads ID:** sdp-jr2n
- **Size:** MEDIUM | **Duration:** 3-4 days
- **Dependencies:** 00-052-19
- **Scope:** `.claude/skills/design/SKILL.md`
- **AC:**
  - Progressive discovery blocks
  - Trigger points: continue/deep design/skip

### 00-052-21: Deep-Thinking Integration (@idea/@design)
- **Beads ID:** sdp-heyv
- **Size:** LARGE | **Duration:** 5-7 days
- **Dependencies:** 00-052-20
- **Scope:** `.claude/skills/idea/SKILL.md`, `.claude/skills/design/SKILL.md`
- **AC:**
  - @idea: 10 expert agents (enhanced mode)
  - @design: 3-5 expert agents (lightweight mode)
  - Synthesizer integration

### 00-052-22: Verbosity Tiers Implementation
- **Beads ID:** sdp-ah0j
- **Size:** SMALL | **Duration:** 2-3 days
- **Dependencies:** 00-052-21
- **Scope:** All skills
- **AC:**
  - --quiet (critical errors only)
  - --verbose (full details)
  - --debug (internal logs + timestamps)

---

## Phase 6: Documentation Track (1 week)

**Dependencies:** Phases 1-5 complete

### 00-052-23: Agent Catalog Documentation
- **Beads ID:** sdp-t6qo
- **Size:** MEDIUM | **Duration:** 3-4 days
- **Dependencies:** Phases 1-5
- **Scope:** `docs/reference/agent-catalog.md`
- **AC:**
  - All 19 agents documented
  - Roles, workflows, examples

### 00-052-24: Updated CLAUDE.md
- **Beads ID:** sdp-ybou
- **Size:** SMALL | **Duration:** 2-3 days
- **Dependencies:** Phases 1-5
- **Scope:** `CLAUDE.md`
- **AC:**
  - Multi-agent architecture section
  - @vision/@reality in decision trees

### 00-052-25: Migration Guide
- **Beads ID:** sdp-u9kb
- **Size:** SMALL | **Duration:** 2-3 days
- **Dependencies:** Phases 1-5
- **Scope:** `docs/migrations/multi-agent-migration.md`
- **AC:**
  - Before/After comparison
  - Migration steps
  - Rollback procedure

---

## Beads Dependency Issues

**Current Problems:**
1. **Cyclic dependencies:** sdp-vxvp → sdp-wbyc → sdp-vxvp (should be sdp-wbyc → sdp-vxvp)
2. **Wrong order:** Some tasks have dependencies backwards

**Correct Dependency Structure:**

```
Phase 0:
  sdp-wqv8 (independent)

Phase 1A (@vision):
  sdp-vxvp (@vision skill) → independent
  sdp-wbyc (extractor) → depends on sdp-vxvp
  sdp-luo8 (CLAUDE.md) → depends on sdp-vxvp

Phase 1B (@reality):
  sdp-uugq (@reality skill) → independent
  sdp-bsds (scanner) → depends on sdp-uugq
  sdp-kxif (CLAUDE.md) → depends on sdp-uugq

Phase 2:
  sdp-vrij (Implementer) → independent
  sdp-01q4 (Spec reviewer) → independent
  sdp-syga (@build update) → depends on sdp-vrij + sdp-01q4
  sdp-n7sp (Test) → depends on sdp-syga
```

**Workaround:** Workstream markdown files have correct dependencies. Follow those, not beads dependency graph.

---

## Next Steps

1. ✅ Backup & worktree setup (00-052-00) - Workstream created
2. ✅ @vision skill structure (00-052-01) - Workstream created
3. ⏳ Create remaining workstream markdown files
4. ⏳ Execute Phase 0 (sdp-wqv8)
5. ⏳ Execute Phase 1A + 1B in parallel
6. ⏳ Execute Phase 2 (QUALITY LOCK-IN)
7. ⏳ Execute Phases 3-6 in parallel tracks

---

**Generated:** 2026-02-07
**Implementation Plan:** `docs/plans/2026-02-07-multi-agent-vision-reality-implementation.md`
