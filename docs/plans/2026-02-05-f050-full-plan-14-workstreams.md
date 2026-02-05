# F050 Go Migration - Full Execution Plan (14 Workstreams)

> **Status:** All workstreams preserved, correctly prioritized
> **Date:** 2026-02-05
> **Timeline:** 14 weeks across 3 phases
> **Source:** 827 sessions analysis + user interview + Go migration requirements

---

## 📊 Execution Strategy

### Principle: Critical First, Then Essential, Then Polish

**Not:** "What do I use now?" (Python-centric thinking)
**But:** "What does Go version need?" (migration requirements)

**Rationale:** All workstreams will be needed in Go version. Question is: when to build them?

---

## Phase 1: Critical (MUST have for Go version)

**Timeline:** 6 weeks | **Goal:** Functional Go SDP with drift detection + reliable orchestration

### WS-00-050-01: Go Project Setup + Core Parser
**Priority:** 🔴 CRITICAL (enabler for all)
**Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** None

**Goal:** Single static binary foundation

**Scope:**
- Go module structure
- Markdown frontmatter parser (workstream files)
- Basic CLI: version, parse, help
- Binary size ≤15MB (stripped)
- Cross-platform builds (macOS/Linux/Windows)

**Acceptance Criteria:**
- [ ] AC1: Go binary compiles for all platforms
- [ ] AC2: Parse workstream markdown into structs
- [ ] AC3: Binary size ≤15MB
- [ ] AC4: Cross-platform build succeeds
- [ ] Coverage ≥80%

---

### WS-00-050-02: TDD Runner Implementation
**Priority:** 🔴 CRITICAL (core workflow)
**Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Red-Green-Refactor cycle with pytest

**Scope:**
- Red phase: Run pytest, expect failure
- Green phase: Run pytest, expect success
- Refactor phase: Ensure tests still pass
- Error reporting and phase tracking

**Acceptance Criteria:**
- [ ] AC1: Red phase runs pytest, expects failure
- [ ] AC2: Green phase runs pytest, expects success
- [ ] AC3: Refactor phase validates tests pass
- [ ] AC4: Phase transitions tracked
- [ ] Coverage ≥80%

---

### WS-00-050-09: Drift Detector ⚡ PRIORITY #1
**Priority:** 🔴 CRITICAL (top pain from 827 sessions)
**Size:** MEDIUM | **Duration:** 2 weeks
**Dependencies:** 00-050-01

**Goal:** Validate documentation matches implementation

**Problem Statement (from usage data):**
> "Documentation-code mismatch causes 4,903 friction events across 827 sessions. Workstream descriptions don't match actual codebase reality."

**Scope:**
- Parse scope_files from workstream frontmatter
- Validate files exist
- Check functions/classes present
- Compare architecture (layers, purposes)
- Generate drift report with recommendations

**Acceptance Criteria:**
- [ ] AC1: Parse scope_files from WS frontmatter
- [ ] AC2: Validate all files exist
- [ ] AC3: Validate declared functions/classes exist
- [ ] AC4: Compare file purpose with documentation
- [ ] AC5: Generate drift report (discrepancies listed)
- [ ] AC6: Integrate with `sdp doctor` (pre-build check)
- [ ] AC7: Run as `sdp drift detect <ws-id>`
- [ ] Coverage ≥80%

**Example Output:**
```bash
$ sdp drift detect 00-040-04

❌ Drift Detected: 3 violations

1. Missing File:
   Expected: src/sdp/quality/validators.py (generic validation)
   Actual: Contains business logic (UserValidator, PaymentValidator)
   Fix: Update WS scope or extract generic validation

2. Structural Mismatch:
   Expected: src/sdp/quality/models.py (validation models)
   Actual: Contains dataclasses (User, Payment)
   Fix: Create separate validation layer

3. Missing Function:
   Expected: validate_contract() in validators.py
   Actual: Function not found
   Fix: Implement or remove from scope
```

---

### WS-00-050-14: Command Auto-Retry 🆕 NEW!
**Priority:** 🔴 CRITICAL (4,904 command failures!)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** None

**Goal:** Auto-retry failed Bash commands with exponential backoff

**Problem Statement (from usage data):**
> "4,904 command failures from 36,868 Bash commands = 13% failure rate. Causes massive friction in workflows."

**Scope:**
- Detect command failures (exit code != 0)
- Retry with exponential backoff (1s, 2s, 4s)
- Max 3 retries before escalating
- Log retry attempts to telemetry
- Integration with TDD runner (auto-retry pytest)
- Integration with orchestrator (auto-retry subprocess calls)

**Acceptance Criteria:**
- [ ] AC1: Detect command failure (exit code)
- [ ] AC2: Retry with exponential backoff (1s → 2s → 4s)
- [ ] AC3: Max 3 retries before escalation
- [ ] AC4: Log retry attempts (timestamp, attempt#, error)
- [ ] AC5: Integration with TDD runner
- [ ] AC6: Integration with orchestrator
- [ ] Coverage ≥80%

**Retry Logic:**
```go
type RetryConfig struct {
    MaxAttempts int    // 3
    BaseDelay   time.Duration  // 1s
    MaxDelay    time.Duration  // 4s
}

func (r *Retrier) Run(cmd *Command) error {
    for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
        err := cmd.Execute()
        if err == nil {
            return nil  // Success
        }

        // Log retry
        log.Printf("Attempt %d/%d failed: %v", attempt, r.MaxAttempts, err)

        // Exponential backoff
        delay := r.BaseDelay * time.Duration(1<<uint(attempt-1))
        if delay > r.MaxDelay {
            delay = r.MaxDelay
        }
        time.Sleep(delay)
    }
    return fmt.Errorf("failed after %d attempts", r.MaxAttempts)
}
```

---

### WS-00-050-11: Multi-Agent Orchestrator (Enhanced)
**Priority:** 🔴 CRITICAL (second top pain)
**Size:** LARGE | **Duration:** 3 weeks
**Dependencies:** 00-050-02, 00-050-03, 00-050-10

**Goal:** Reliable wave-based parallel execution

**Problem Statement (from user interview):**
> "Обрыв выполнения при заполнении контекста, пропуск агентами операций"

**Enhanced Scope (beyond original plan):**
- Wave-based parallelism (via `bd ready`)
- Adaptive agent count: 2 (small), 5 (medium), 10 (large)
- Cascading failure detection (3+ consecutive failures)
- **NEW:** Context chunking for large features (>15 workstreams)
- **NEW:** Operation tracking with acknowledgments
- **NEW:** Retry mechanism for skipped operations
- **NEW:** Checkpoint save after each workstream
- Resume capability via `--resume` flag

**Acceptance Criteria:**
- [ ] AC1: Workstreams executed in waves (via `bd ready`)
- [ ] AC2: Adaptive agent count based on feature size
- [ ] AC3: Cascading failure detection (3+ consecutive failures)
- [ ] AC4: Checkpoint saved after each workstream
- [ ] AC5: Resume capability works
- [ ] AC6: `sdp oneshot <feature>` executes autonomously
- [ ] **AC7 (NEW):** Context chunking prevents overflow
- [ ] **AC8 (NEW):** Operation acknowledgments detect skips
- [ ] **AC9 (NEW):** Retry mechanism recovers operations
- [ ] Coverage ≥80%

**Context Chunking Example:**
```go
if len(workstreams) > 15 {
    // Split into chunks of 10
    chunks := chunkWorkstreams(workstreams, 10)
    for i, chunk := range chunks {
        log.Printf("Executing chunk %d/%d (%d workstreams)", i+1, len(chunks), len(chunk))
        executeChunk(chunk)  // Fresh context per chunk
        checkpoint.Save()
    }
}
```

---

### WS-00-050-10: Checkpoint System
**Priority:** 🔴 HIGH (orchestration enabler)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** None

**Goal:** Save/restore execution state

**Scope:**
- JSON-based checkpoint storage (`.sdp/checkpoints/<feature>.json`)
- Track: FeatureID, AgentID, Status, CompletedWS[], CurrentWS
- Save after each workstream
- Load on resume
- Delete after completion

**Acceptance Criteria:**
- [ ] AC1: Checkpoint saved to JSON file
- [ ] AC2: Checkpoint loaded by feature ID
- [ ] AC3: Completed workstreams tracked
- [ ] AC4: Current workstream tracked
- [ ] AC5: Resume via `--resume` flag
- [ ] AC6: Deleted after completion
- [ ] Coverage ≥80%

---

**Phase 1 Summary:** 6 workstreams, 6 weeks, **functional Go SDP**

---

## Phase 2: Essential (SHOULD have for complete Go SDP)

**Timeline:** 4 weeks | **Goal:** Full-featured Go SDP with telemetry + integration

### WS-00-050-03: Beads CLI Wrapper
**Priority:** 🟡 HIGH (dependency for orchestrator)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Thin wrapper around Beads CLI

**Scope:**
- Execute `bd ready`, `bd create`, `bd dep add`, `bd update`
- Parse JSON output
- Maintain SDP ↔ Beads mapping (`.beads-sdp-mapping.jsonl`)
- Helper commands for common workflows

**Acceptance Criteria:**
- [ ] AC1: Execute `bd ready`, parse JSON
- [ ] AC2: Execute `bd dep add`, `bd update`
- [ ] AC3: Maintain mapping file
- [ ] AC4: Helper commands work
- [ ] Coverage ≥80%

---

### WS-00-050-04: CLI Commands
**Priority:** 🟡 HIGH (deployment + daily workflow)
**Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Essential CLI commands

**Scope:**
- `sdp init <project>` - Initialize SDP
- `sdp doctor` - Environment checks
- `sdp build <ws-id>` - Execute workstream
- `sdp drift detect <ws-id>` - Detect drift
- `sdp oneshot <feature>` - Autonomous execution
- `sdp --version` - Version info

**Acceptance Criteria:**
- [ ] AC1: All commands implemented
- [ ] AC2: Help text works
- [ ] AC3: Error handling works
- [ ] Coverage ≥80%

---

### WS-00-050-05: Quality Gates (Parallel Execution)
**Priority:** 🟡 MEDIUM (performance optimization)
**Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Parallel quality gate execution for speed

**Scope:**
- Run mypy, ruff, pytest coverage in parallel (goroutines)
- Target: <8s total (vs Python 8.5s serial)
- Aggregate results
- Unified error reporting

**Acceptance Criteria:**
- [ ] AC1: Mypy, ruff, coverage run in parallel
- [ ] AC2: Execution time <8s
- [ ] AC3: Results aggregated correctly
- [ ] AC4: Error reporting unified
- [ ] Coverage ≥80%

**Why Still Needed:** Go version will run quality gates, parallelism improves UX

---

### WS-00-050-07: Telemetry Collector (Simplified)
**Priority:** 🟢 MEDIUM (observability)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** 00-050-04

**Goal:** Basic execution tracking

**Scope (simplified from original):**
- Auto-capture git stats (files, LOC)
- Test coverage before/after
- Execution duration (timestamps)
- Friction points (quality failures)
- Save to workstream frontmatter

**Acceptance Criteria:**
- [ ] AC1: Capture git stats
- [ ] AC2: Capture coverage before/after
- [ ] AC3: Capture execution duration
- [ ] AC4: Capture friction points
- [ ] AC5: Append to WS frontmatter
- [ ] Coverage ≥80%

**Why Simplified:** Single user doesn't need complex analytics pipeline

---

**Phase 2 Summary:** 4 workstreams, 4 weeks, **complete Go SDP**

---

## Phase 3: Polish (NICE to have for v1.0.0-go)

**Timeline:** 2 weeks | **Goal:** Production-ready polish

### WS-00-050-06: Quality Watcher
**Priority:** 🟢 LOW (UX enhancement)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** 00-050-05

**Goal:** Real-time quality feedback

**Scope:**
- Watch .py files via fsnotify (cross-platform)
- Run quality gates on file change
- Desktop notifications (macOS/Linux/Windows)
- Configurable watch paths

**Acceptance Criteria:**
- [ ] AC1: Watcher monitors .py files
- [ ] AC2: Quality gates run on change
- [ ] AC3: Desktop notifications work
- [ ] AC4: Configurable watch paths
- [ ] Coverage ≥80%

**Why Still Needed:** Go fsnotify is more reliable than Python watchdog

---

### WS-00-050-08: Telemetry Analyzer
**Priority:** 🟢 LOW (analytics)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** 00-050-07

**Goal:** Pattern detection and insights

**Scope:**
- Analyze telemetry data
- Detect patterns (friction points, slow workstreams)
- Generate insights report
- Suggest optimizations

**Acceptance Criteria:**
- [ ] AC1: Parse telemetry data
- [ ] AC2: Detect friction patterns
- [ ] AC3: Generate insights
- [ ] AC4: Suggest optimizations
- [ ] Coverage ≥80%

**Why Still Needed:** Even single user benefits from pattern analysis

---

### WS-00-050-12: CLI Polish
**Priority:** 🟢 LOW (UX polish)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** 00-050-04

**Goal:** Production-ready CLI

**Scope:**
- Shell completion (bash, zsh, fish)
- Color output (lipgloss)
- Progress bars (bubbletea)
- Man page generation

**Acceptance Criteria:**
- [ ] AC1: Bash completion works
- [ ] AC2: Zsh completion works
- [ ] AC3: Fish completion works
- [ ] AC4: Color output displays
- [ ] AC5: Progress bars display
- [ ] Coverage ≥80%

**Why Still Needed:** Production CLI needs polish

---

### WS-00-050-13: Python Code Removal
**Priority:** 🟡 CRITICAL (cleanup)
**Size:** SMALL | **Duration:** 1 week
**Dependencies:** All previous workstreams

**Goal:** Remove Python after Go verified

**Scope:**
- Verify Go binary works (all tests pass)
- Create backup branch (`python-legacy`)
- Delete `src/sdp/` Python code
- Delete `pyproject.toml`, `poetry.lock`
- Update documentation (remove Python references)
- Tag release `v1.0.0-go`

**Acceptance Criteria:**
- [ ] AC1: Go binary verified working
- [ ] AC2: Python archived to `python-legacy` branch
- [ ] AC3: Documentation updated
- [ ] AC4: Release tagged
- [ ] Coverage ≥80%

---

**Phase 3 Summary:** 4 workstreams, 2 weeks (overlap), **production-ready**

---

## 📊 Final Summary

### Workstream Count
- **Original Plan:** 13 workstreams
- **New Plan:** 14 workstreams (+1: Command Auto-Retry)
- **Timeline:** 12 weeks (Phase 1+2) + 2 weeks polish = **14 weeks**

### By Phase

| Phase | Workstreams | Duration | Purpose |
|-------|-------------|----------|---------|
| **Phase 1: Critical** | 6 | 6 weeks | Functional Go SDP |
| **Phase 2: Essential** | 4 | 4 weeks | Complete Go SDP |
| **Phase 3: Polish** | 4 | 2 weeks | Production-ready |
| **Total** | **14** | **12-14 weeks** | **Full migration** |

### By Priority

| Priority | Count | Workstreams |
|----------|-------|-------------|
| 🔴 CRITICAL | 6 | 01, 02, 09, 14, 11, 10 |
| 🟡 HIGH | 2 | 03, 04 |
| 🟡 MEDIUM | 2 | 05, 07 |
| 🟢 LOW | 3 | 06, 08, 12 |
| 🟡 CRITICAL | 1 | 13 (cleanup) |

### Execution Order

**Week 1-2: Foundation**
- WS-00-050-01: Go Project Setup
- WS-00-050-02: TDD Runner

**Week 3-4: Critical Fixes**
- WS-00-050-14: Command Auto-Retry
- WS-00-050-10: Checkpoint System

**Week 5-6: Drift + Orchestration**
- WS-00-050-09: Drift Detector ⚡ PRIORITY #1
- WS-00-050-11: Orchestrator (starts, 3 weeks)

**Week 7-8: Integration**
- WS-00-050-03: Beads Wrapper
- WS-00-050-04: CLI Commands

**Week 9-11: Features**
- WS-00-050-05: Quality Gates (parallel)
- WS-00-050-07: Telemetry
- WS-00-050-11: Orchestrator (continues)

**Week 12-14: Polish + Cleanup**
- WS-00-050-06: Quality Watcher
- WS-00-050-08: Telemetry Analyzer
- WS-00-050-12: CLI Polish
- WS-00-050-13: Python Removal

---

## 🎯 Success Metrics (From Usage Data)

| Metric | Baseline | Target | How to Measure |
|--------|----------|--------|----------------|
| **Drift Detection** | Manual discovery (4,903 events) | Automatic detection | `/verify-workstream` violations caught |
| **Command Failures** | 4,904 (13% rate) | <5% (auto-retry) | Retry telemetry, failure rate |
| **Orchestration Reliability** | Failed sessions | 95% completion | Checkpoint resume, operation ack |
| **Quality Gate Speed** | 8.5s serial | <8s parallel | Goroutine execution time |
| **Deployment Time** | pip/poetry (5-10 min) | Binary copy (30 sec) | Fresh install time |
| **Context Overflow** | Large features fail | Chunked execution | 20+ WS features succeed |

---

## 🚀 Ready to Execute

### Next Steps

**Option A: Execute Phase 1 (Critical Path)**
```bash
@build 00-050-01  # Week 1: Go Project Setup
@build 00-050-02  # Week 1: TDD Runner
@build 00-050-14  # Week 2: Command Auto-Retry
@build 00-050-10  # Week 2: Checkpoint System
@build 00-050-09  # Week 3-4: Drift Detector ⚡
@build 00-050-11  # Week 5-7: Orchestrator
```

**Option B: Autonomous Execution**
```bash
@oneshot F050  # Execute all 14 workstreams
```

**Option C: Update Files First**
```bash
# Create/update all 14 workstream files
# Then execute Phase 1
```

---

**What's your choice?**
