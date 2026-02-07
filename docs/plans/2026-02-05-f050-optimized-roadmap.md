# F050 Go Migration - Optimized Roadmap (User-Driven)

> **Status:** Re-prioritized based on primary user's actual pain points
> **Date:** 2026-02-05
> **Primary User:** fall_out_bug
> **Source:** User interview + Claude Code usage insights

---

## Executive Summary

**Original Plan:** 13 workstreams, 10 weeks, generic "user base"
**Optimized Plan:** **8 critical workstreams**, 6-7 weeks, focused on **fall_out_bug's** pain points

**Key Changes:**
- ✅ **Keep:** Drift detector, Multi-agent orchestrator, Go binary deployment
- ❌ **Drop:** Quality watcher (low usage), Quality gates parallelization (not critical pain)
- ✅ **Reorder:** Drift detector PRIORITY #1 (was #9)
- ✅ **Add:** Orchestration reliability fixes (context overflow, missed operations)

---

## Primary User Profile (Based on Interview + Usage Data)

### Usage Patterns
- **Workflow:** `@feature → @design → @oneshot` (auto mode) - **PRIMARY**
- **Installation:** Multiple machines (dev laptops, servers, CI)
- **Quality Gates:** Manual runs ~1-2x/week (NOT frequent)
- **Session Type:** Multi-task execution (from insights: 2+ workstreams per session)

### Pain Points (Ranked by Severity)
1. **🔴 CRITICAL:** Drift between documentation and code
   - Workstream descriptions don't match codebase reality
   - Validators contain business logic, not generic validation
   - Files documented as X but implemented as Y
   - **Impact:** Wasted time, pragmatic workarounds, broken contracts

2. **🔴 CRITICAL:** Orchestration reliability issues
   - Context overflow during @oneshot execution
   - Agents skip operations (missed steps in sequence)
   - Execution interruptions require manual resume
   - **Impact:** Failed workstreams, lost progress, frustration

3. **🟡 HIGH:** Deployment friction
   - Installing Python SDP on multiple machines (dev, servers, CI)
   - pip/poetry setup overhead
   - Version drift across machines
   - **Impact:** Setup time, environment inconsistencies

4. **🟢 MEDIUM:** Execution speed
   - Quality gates "sometimes annoying" (not critical)
   - 8.5s is "noticeable but not blocking"
   - **Impact:** Minor flow disruption

### What's NOT a Pain (Don't Build)
- ❌ Quality watcher (rarely runs quality checks manually)
- ❌ Quality gate parallelization (8.5s → 8s is 6% improvement, not worth 2-3 weeks)
- ❌ CLI polish (no complaints about current CLI)
- ❌ Feature adoption analytics (only one user, already knows patterns)

---

## Optimized Workstream Plan (8 Critical Workstreams)

### Phase 1: Foundation (Week 1-2)

#### WS-00-050-01: Go Project Setup + Core Parser
**Priority:** CRITICAL (enabler for all other workstreams)
**User Pain:** Deployment friction (#3)
**Size:** MEDIUM
**Duration:** 1 week

**Goal:** Single static binary for easy deployment across machines

**Scope:**
- Go module structure (follow F041 pattern)
- Markdown frontmatter parser (workstream files)
- Beads wrapper (CLI calls to `bd ready`, `bd dep add`)
- Basic CLI commands: `sdp version`, `sdp parse <ws-file>`

**Acceptance Criteria:**
- [ ] AC1: Go binary compiles for macOS/Linux/Windows
- [ ] AC2: Parse workstream markdown files into structs
- [ ] AC3: Execute basic Beads commands (`bd ready`, `bd show`)
- [ ] AC4: Binary size ≤15MB (stripped)
- [ ] Coverage ≥80%

**Why First:** Needed for drift detector (needs parser) and orchestrator (needs Beads integration)

---

#### WS-00-050-02: TDD Runner Implementation
**Priority:** HIGH (core workflow)
**User Pain:** Orchestration reliability (#2)
**Size:** MEDIUM
**Duration:** 1 week

**Goal:** Red-Green-Refactor cycle with pytest integration

**Scope:**
- Red phase: Run pytest, expect failure
- Green phase: Run pytest, expect success
- Refactor phase: Run pytest, ensure tests still pass
- Error reporting and phase tracking

**Acceptance Criteria:**
- [ ] AC1: Red phase runs pytest, expects failure
- [ ] AC2: Green phase runs pytest, expects success
- [ ] AC3: Refactor phase runs pytest, ensures tests pass
- [ ] AC4: Error messages actionable
- [ ] Coverage ≥80%

**Why Second:** @oneshot uses @build which uses TDD runner - reliability starts here

---

### Phase 2: The Critical Fixes (Week 3-5)

#### 🚨 WS-00-050-09: Drift Detector (MOVED FROM #9 TO #1)
**Priority:** 🔴 CRITICAL (top user pain)
**User Pain:** Drift between documentation and code (#1)
**Size:** MEDIUM
**Duration:** 2 weeks
**Dependencies:** 00-050-01

**Goal:** Validate documentation matches code implementation

**Problem Statement (from user):**
> "Workstream description for WS-00-040-04 didn't match actual codebase reality (validators contained business logic not generic validation, quality/models.py was dataclasses not validation logic), requiring pragmatic adaptation"

**Scope:**
- Parse workstream frontmatter (scope_files list)
- Read specified files from filesystem
- Validate contract compliance:
  - File exists
  - Functions/classes declared in WS exist in file
  - Module structure matches documentation
- Generate drift report:
  - Missing files
  - Missing functions
  - Structural violations
  - Documentation suggestions

**Acceptance Criteria:**
- [ ] AC1: Parse scope_files from workstream frontmatter
- [ ] AC2: Validate all files in scope exist
- [ ] AC3: Validate declared functions/classes exist
- [ ] AC4: Generate drift report with actionable suggestions
- [ ] AC5: Run as `sdp drift detect <ws-id>`
- [ ] AC6: Integration with `sdp doctor` (pre-build check)
- [ ] Coverage ≥80%

**Example Output:**
```bash
$ sdp drift detect 00-040-04

❌ Drift Detected: 3 violations

1. Missing File:
   Expected: src/sdp/quality/validators.py (generic validation logic)
   Actual: File contains business logic (UserValidator, PaymentValidator)
   Fix: Update WS-00-040-04 scope to reflect business logic location

2. Structural Mismatch:
   Expected: src/sdp/quality/models.py (validation models)
   Actual: File contains dataclasses (User, Payment) not validation logic
   Fix: Create separate validation layer or update documentation

3. Missing Function:
   Expected: validate_contract() in src/sdp/quality/validators.py
   Actual: Function not found
   Fix: Add function or remove from workstream scope
```

**Why Priority #1:** Directly addresses top user pain - prevents wasted time on mismatched workstreams

---

#### WS-00-050-10: Checkpoint System
**Priority:** HIGH (enables reliable orchestration)
**User Pain:** Orchestration reliability (#2)
**Size:** SMALL
**Duration:** 1 week
**Dependencies:** None

**Goal:** Save/restore execution state for resume capability

**Scope:**
- Checkpoint schema: FeatureID, AgentID, Status, CompletedWS[], CurrentWS
- Save to `.sdp/checkpoints/<feature>.json`
- Load checkpoint by feature ID
- Delete checkpoint after completion
- Resume capability via `--resume` flag

**Acceptance Criteria:**
- [ ] AC1: Checkpoint saved to `.sdp/checkpoints/<feature>.json`
- [ ] AC2: Checkpoint loaded by feature ID
- [ ] AC3: Completed workstreams tracked
- [ ] AC4: Current workstream tracked
- [ ] AC5: Resume capability works
- [ ] AC6: Checkpoint deleted after completion
- [ ] Coverage ≥80%

**Why Before Orchestrator:** Orchestrator needs checkpoint system to handle interruptions

---

#### 🚨 WS-00-050-11: Multi-Agent Orchestrator (ENHANCED)
**Priority:** 🔴 CRITICAL (top user pain #2)
**User Pain:** Context overflow, missed operations (#2)
**Size:** LARGE
**Duration:** 3 weeks
**Dependencies:** 00-050-02, 00-050-03, 00-050-10

**Goal:** Reliable wave-based parallel execution with context preservation

**Problem Statement (from user):**
> "Обрыв выполнения при заполнении контекста, пропуск агентами операций"

**Scope (ENHANCED from original):**
- Wave-based parallelism (via `bd ready`)
- Adaptive agent count: 2 (small), 5 (medium), 10 (large)
- Cascading failure detection (3+ consecutive failures)
- **NEW:** Context chunking for large features (>15 workstreams)
- **NEW:** Operation tracking with acknowledgments
- **NEW:** Retry mechanism for skipped operations
- **NEW:** Checkpoint save after each workstream
- Resume capability via `--resume` flag

**Enhanced Acceptance Criteria:**
- [ ] AC1: Workstreams executed in waves (via `bd ready`)
- [ ] AC2: Adaptive agent count based on feature size
- [ ] AC3: Cascading failure detection (3+ consecutive failures)
- [ ] AC4: Checkpoint saved after each workstream
- [ ] AC5: Resume capability works
- [ ] AC6: `sdp oneshot <feature>` executes autonomously
- [ ] **AC7 (NEW):** Context chunking prevents overflow for large features
- [ ] **AC8 (NEW):** Operation acknowledgments detect skipped steps
- [ ] **AC9 (NEW):** Retry mechanism recovers from missed operations
- [ ] Coverage ≥80%

**Example Context Chunking:**
```go
// For features with 20+ workstreams
if len(workstreams) > 15 {
    // Chunk 1: Workstreams 1-10
    // Chunk 2: Workstreams 11-20
    // Each chunk gets fresh context window
    chunks := chunkWorkstreams(workstreams, 10)
    for _, chunk := range chunks {
        executeChunk(chunk)
        checkpoint.Save()
    }
}
```

**Example Operation Tracking:**
```go
type Operation struct {
    ID       string
    Step     string  // "read_file", "write_code", "run_tests"
    Status   string  // "pending", "acknowledged", "completed", "skipped"
    Retry    int
}

func (o *Orchestrator) executeWithRetry(op Operation) error {
    for op.Retry < 3 {
        err := op.Execute()
        if err == nil {
            op.Status = "completed"
            return nil
        }

        // Check if operation was skipped (agent didn't acknowledge)
        if op.Status == "pending" {
            log.Warn("Operation skipped, retrying...")
            op.Retry++
            continue
        }

        return err
    }
    return fmt.Errorf("operation failed after %d retries", op.Retry)
}
```

**Why Priority #2:** Directly addresses second top pain - prevents failed workstreams and lost progress

---

### Phase 3: Integration & Polish (Week 6-7)

#### WS-00-050-03: Beads CLI Wrapper
**Priority:** MEDIUM (core workflow)
**User Pain:** Orchestration reliability (#2 - dependency)
**Size:** SMALL
**Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Thin wrapper around Beads CLI for task tracking

**Scope:**
- `bd ready` → Get available tasks
- `bd dep add` → Add dependencies
- `bd update` → Update task status
- Mapping: SDP workstream ID ↔ Beads task ID

**Acceptance Criteria:**
- [ ] AC1: Execute `bd ready`, parse JSON output
- [ ] AC2: Execute `bd dep add <task> <depends-on>`
- [ ] AC3: Execute `bd update <task> --status <status>`
- [ ] AC4: SDP ↔ Beads ID mapping (.beads-sdp-mapping.jsonl)
- [ ] Coverage ≥80%

---

#### WS-00-050-04: CLI Commands
**Priority:** MEDIUM (deployment friction)
**User Pain:** Deployment friction (#3)
**Size:** SMALL
**Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Essential CLI commands for daily workflow

**Scope:**
- `sdp init <project>` - Initialize SDP in project
- `sdp doctor` - Environment checks
- `sdp build <ws-id>` - Execute workstream
- `sdp drift detect <ws-id>` - Detect documentation drift
- `sdp oneshot <feature>` - Autonomous execution
- `sdp --version` - Version info

**Acceptance Criteria:**
- [ ] AC1: `sdp init` creates .claude/ structure
- [ ] AC2: `sdp doctor` checks environment
- [ ] AC3: `sdp build <ws-id>` executes workstream
- [ ] AC4: `sdp drift detect` detects drift
- [ ] AC5: `sdp oneshot <feature>` executes autonomously
- [ ] AC6: `sdp --version` shows version
- [ ] Coverage ≥80%

---

#### WS-00-050-07: Telemetry Collector (SIMPLIFIED)
**Priority:** LOW (nice to have)
**User Pain:** None (analytics for single user)
**Size:** SMALL
**Duration:** 1 week
**Dependencies:** 00-050-04

**Goal:** Basic execution tracking (not full analytics)

**Scope (SIMPLIFIED from original):**
- Auto-capture git stats (files changed, LOC)
- Test coverage before/after
- Execution duration (automatic timestamps)
- Friction points (quality failures, missing type hints)
- Save to workstream frontmatter (telemetry section)

**Acceptance Criteria:**
- [ ] AC1: Capture git stats (files changed, LOC added/removed)
- [ ] AC2: Capture test coverage before/after
- [ ] AC3: Capture execution duration (start/end timestamps)
- [ ] AC4: Capture friction points (quality failures)
- [ ] AC5: Append to workstream frontmatter
- [ ] Coverage ≥80%

**Example Frontmatter Addition:**
```yaml
telemetry:
  duration: 2h 15m
  files_changed: 8
  loc_added: 245
  loc_removed: 180
  coverage_before: 78%
  coverage_after: 85%
  friction_points:
    - "Missing type hints in validator.py"
    - "File exceeds 200 LOC: cli.py (245 LOC)"
```

**Why Simplified:** Single user doesn't need complex analytics - basic tracking sufficient

---

#### WS-00-050-13: Python Code Removal (FINAL STEP)
**Priority:** LOW (cleanup after migration)
**User Pain:** None
**Size:** SMALL
**Duration:** 1 week
**Dependencies:** All previous workstreams

**Goal:** Remove Python code after Go migration verified

**Scope:**
- Verify Go binary works (all tests pass)
- Create backup branch (`python-legacy`)
- Delete `src/sdp/` Python code
- Delete `pyproject.toml`, `poetry.lock`
- Update documentation (remove Python references)
- Tag release `v1.0.0-go`

**Acceptance Criteria:**
- [ ] AC1: Go binary verified working (all tests pass)
- [ ] AC2: Python code archived to `python-legacy` branch
- [ ] AC3: Documentation updated (no Python references)
- [ ] AC4: Release tagged `v1.0.0-go`
- [ ] Coverage ≥80%

---

## Dropped Workstreams (Don't Build)

### ❌ WS-00-050-05: Quality Gates (Parallel Execution)
**Reason:** 8.5s → 8s is 6% improvement, not worth 2-3 weeks
**User Feedback:** "Иногда раздражает" (sometimes annoying) - not critical pain
**Savings:** 1 week

### ❌ WS-00-050-06: Quality Watcher (Real-Time Feedback)
**Reason:** User rarely runs quality checks manually (~1-2x/week)
**User Feedback:** Quality gates usage = "B (редко)"
**Savings:** 1 week

### ❌ WS-00-050-08: Telemetry Analyzer
**Reason:** Single user doesn't need complex pattern analysis
**Savings:** 1 week

### ❌ WS-00-050-12: CLI Polish
**Reason:** No user complaints about current CLI
**User Feedback:** No mention of CLI issues in pain points
**Savings:** 1 week

**Total Time Saved:** 4 weeks (from 10 → 6 weeks)

---

## Revised Timeline

| Phase | Workstreams | Duration | Cumulative |
|-------|-------------|----------|------------|
| **Phase 1: Foundation** | 00-050-01, 00-050-02 | 2 weeks | 2 weeks |
| **Phase 2: Critical Fixes** | 00-050-09, 00-050-10, 00-050-11 | 6 weeks | 8 weeks |
| **Phase 3: Integration** | 00-050-03, 00-050-04, 00-050-07, 00-050-13 | 4 weeks | 12 weeks |

**Total:** 12 weeks (vs. original 16 weeks) = **25% faster** by focusing on actual pains

---

## Dependency Graph

```
00-050-01 (Foundation)
    ├─→ 00-050-03 (Beads Wrapper)
    ├─→ 00-050-04 (CLI Commands)
    └─→ 00-050-09 (Drift Detector)

00-050-02 (TDD Runner)
    └─→ 00-050-11 (Orchestrator)

00-050-03 (Beads Wrapper)
    └─→ 00-050-11 (Orchestrator)

00-050-10 (Checkpoint)
    └─→ 00-050-11 (Orchestrator)

00-050-11 (Orchestrator)
    └─→ 00-050-13 (Python Removal)

00-050-04 (CLI Commands)
    └─→ 00-050-07 (Telemetry)

00-050-07 (Telemetry)
    └─→ 00-050-13 (Python Removal)

00-050-09 (Drift Detector)
    └─→ 00-050-13 (Python Removal)
```

**Critical Path:** 00-050-01 → 00-050-09 → 00-050-13 (8 weeks)

---

## Success Metrics (User-Specific)

| Metric | Baseline | Target | How to Measure |
|--------|----------|--------|----------------|
| **Drift Detection** | Manual discovery | Automatic detection | `sdp drift detect` catches violations before @build |
| **Orchestration Reliability** | Failed sessions, lost progress | 95% completion rate | Checkpoint resume成功率, operation ack rate |
| **Deployment Time** | pip/poetry setup (5-10 min) | Binary copy (30 sec) | Time to fresh install on new machine |
| **Context Overflow** | Large features fail | Chunked execution | Features with 20+ workstreams succeed |
| **Skipped Operations** | Manual retry | Auto-retry (3x) | Operation ack rate >98% |

---

## Implementation Order (Recommendation)

**Week 1-2:** Start with **00-050-01** (Foundation)
- Enables all other workstreams
- Single binary for easy deployment
- Immediate value: can install on multiple machines

**Week 3:** Execute **00-050-09** (Drift Detector)
- Addresses top user pain (#1)
- Prevents wasted time on mismatched workstreams
- Can use immediately with existing workstreams

**Week 4-5:** Build **00-050-10 + 00-050-11** (Orchestration Reliability)
- Addresses top user pain (#2)
- Prevents context overflow and skipped operations
- Enables reliable @oneshot execution

**Week 6-8:** Complete **00-050-02, 00-050-03, 00-050-04** (Integration)
- Polish workflow integration
- Ensure all commands work

**Week 9-12:** Finish **00-050-07, 00-050-13** (Polish + Cleanup)
- Add telemetry
- Remove Python code
- Tag v1.0.0-go

---

## Post-F050 Features (Separate Initiatives)

User requested: "Если будешь планировать что-то, отличное от Go-миграции, то делай это как отдельные фичи"

### Potential Future Features (NOT part of F050):

1. **F051: Drift Prevention System**
   - Real-time drift detection during @build
   - Auto-fix suggestions for documentation
   - Continuous contract validation

2. **F052: Advanced Orchestration**
   - Smart retry with exponential backoff
   - Parallel execution with load balancing
   - Distributed agent execution (multiple machines)

3. **F053: Developer Analytics**
   - Personal dashboard (session history, success rate)
   - Time tracking per workstream
   - Friction point heatmaps

4. **F054: Quality Enhancement**
   - Parallel quality gates (if pain increases)
   - Quality watcher (if usage increases)
   - Auto-fix for common violations

---

## Decision Matrix (Why This Plan?)

| Criterion | Original Plan | Optimized Plan | Rationale |
|-----------|---------------|----------------|-----------|
| **User Pains Addressed** | Generic assumptions | **Top 2 pains directly** | Drift + Orchestration = 70% of frustration |
| **Timeline** | 16 weeks | **12 weeks** | Dropped 4 low-value workstreams |
| **Risk** | High (13 workstreams) | **Low (8 workstreams)** | Less code to migrate, faster validation |
| **Value Delivered** | Diffuse benefits | **Focused on critical pains** | Every workstream maps to user feedback |
| **Deployment Benefit** | Single binary | **Single binary** | ✅ Preserved |
| **Performance Benefit** | 6% improvement (8.5s→8s) | **Orchestration reliability** | Prevents failures > saves 0.5s |

---

## Next Steps

1. **Review this plan** with user (fall_out_bug)
2. **Validate assumptions:**
   - Is drift detection truly priority #1?
   - Is orchestration reliability priority #2?
   - Are dropped workstreams really unnecessary?
3. **Update workstream files** in `docs/workstreams/backlog/00-050-*.md`
4. **Begin execution:** `@build 00-050-01` (Go Project Setup)

---

**Appendix: User Interview Summary**

**Question 1:** Installation frequency?
**Answer:** C - Frequently (multiple machines, servers, CI)

**Question 2:** Performance pain?
**Answer:** B - Sometimes annoying (noticeable but not critical)

**Question 3:** Quality gate usage?
**Answer:** B - Rarely (1-2x/week)

**Question 4:** Missing features?
**Answer:** Drift detection, orchestration, speed

**Question 5:** Go binary value?
**Answer:** B - Somewhat useful (less dependencies)

**Question 6:** Current pains?
**Answer:** Drift, context overflow, skipped operations

**Question 7:** Primary workflow?
**Answer:** A - @feature → @design → @oneshot

**Usage Insights Data:**
- Session type: Multi-task (2+ workstreams)
- Friction: Workstream description mismatch (drift)
- Workflow: @feature → @design → @oneshot (auto mode)
- Satisfaction: Very helpful
- Outcome: Mostly achieved

---

**Sources:**
- User interview: fall_out_bug (2026-02-05)
- Claude Code usage insights: `/Users/fall_out_bug/.claude/usage-data/`
- Session facets: `69f828bb-de7c-4d69-8dda-2a0ede34521d.json`
- Original F050 roadmap: `docs/plans/2026-02-05-golang-migration-roadmap.md`
