# SDP Go Migration + Telemetry Enhancement Roadmap

> **Status:** Ready for Implementation
> **Date:** 2026-02-05
> **Timeline:** 10 weeks
> **Goal:** Migrate to Go + implement telemetry-driven improvements

---

## Overview

This roadmap combines two major initiatives:

1. **Go Migration** — Single binary deployment, leverage Beads CLI
2. **Telemetry Enhancements** — Insights from deep-thinking analysis

**Integration Strategy:** Implement Go migration FIRST (infrastructure), then layer telemetry enhancements on top.

---

## Strategic Alignment

### Deep-Thinking Recommendations → Go Implementation

| Deep-Thinking Recommendation | Go Workstream | Priority |
|------------------------------|---------------|----------|
| **Auto-capture execution telemetry** | 00-GO-04: Telemetry Collector | HIGH |
| **Documentation-code synchronization** | 00-GO-06: Drift Detector | HIGH |
| **Real-time quality gates** | 00-GO-05: Quality Watcher | HIGH |
| **Adaptive parallelism for @oneshot** | 00-GO-08: Multi-Agent Orchestrator | MEDIUM |
| **Telemetry-driven protocol evolution** | 00-GO-09: Insights Analyzer | MEDIUM |

### Key Insight from Deep-Thinking

> "User encounters friction when workstream documentation doesn't match codebase reality. Pragmatic adaptation pattern: 3.6 commits/session, executes quickly but discovers mismatches mid-workstream."

**Solution:** Drift detection (00-GO-06) catches mismatches BEFORE implementation, not after.

---

## Workstream Decomposition

### Legend

- **Size:** SMALL (<500 LOC), MEDIUM (500-1500 LOC), LARGE (1500-3000 LOC)
- **Priority:** P0 (critical), P1 (high), P2 (medium)
- **Dependencies:** `00-GO-01` must complete before `00-GO-02`

---

## Phase 1: Foundation (Week 1-3)

### 00-GO-01: Go Project Setup + Core Parser
**Size:** MEDIUM (600 LOC)
**Priority:** P0
**Depends on:** None

**Goal:** Setup Go project structure, implement workstream YAML parser

**Scope Files:**
- `cmd/sdp/main.go`
- `cmd/sdp/root.go`
- `internal/core/workstream.go`
- `internal/core/validator.go`
- `internal/core/schema.go`
- `go.mod`, `go.sum`

**Acceptance Criteria:**
- AC1: `go build ./cmd/sdp` produces executable binary
- AC2: `sdp --version` displays version string
- AC3: Parse workstream YAML with frontmatter successfully
- AC4: Validate workstream ID format (PP-FFF-SS)
- AC5: Validate capability tier (T0/T1/T2/T3)

**Implementation Notes:**
- Use `gopkg.in/yaml.v3` for YAML parsing
- Regex for WS ID validation: `^\d{2}-\d{3}-\d{2}$`
- Test with real workstream from `docs/workstreams/backlog/`

---

### 00-GO-02: TDD Runner Implementation
**Size:** MEDIUM (400 LOC)
**Priority:** P0
**Depends on:** 00-GO-01

**Goal:** Implement Red-Green-Refactor cycle with pytest integration

**Scope Files:**
- `internal/tdd/runner.go`
- `internal/tdd/phase.go`
- `cmd/sdp/commands/build.go`

**Acceptance Criteria:**
- AC1: Red phase runs test, expects failure
- AC2: Green phase runs test, expects success
- AC3: Refactor phase ensures tests still pass
- AC4: `sdp build 00-001-01` executes TDD cycle end-to-end
- AC5: Output includes phase duration and next phase recommendation

**Implementation Notes:**
- Subprocess calls to `pytest` (not rewriting in Go)
- Capture stdout/stderr for error messages
- Context support for cancellation

---

### 00-GO-03: Beads CLI Wrapper (Thin Layer)
**Size:** SMALL (150 LOC)
**Priority:** P0
**Depends on:** 00-GO-01

**Goal:** Thin wrapper functions for Beads CLI commands

**Scope Files:**
- `internal/beads/wrapper.go`
- `internal/beads/fallback.go`

**Acceptance Criteria:**
- AC1: `ReadyTasks()` calls `bd ready --json`, parses output
- AC2: `UpdateStatus()` calls `bd update <id> --status <status>`
- AC3: `CreateTask()` calls `bd create` with metadata
- AC4: Fallback to JSON if Beads not installed
- AC5: All functions return Go errors (not panics)

**Implementation Notes:**
- JSON parsing from `bd --json` output
- Only 5-6 wrapper functions (not full client)
- Fallback mechanism for environments without Beads

**Dependency:** Beads CLI must be installed (or use fallback)

---

### 00-GO-04: CLI Commands (init, doctor, build)
**Size:** MEDIUM (500 LOC)
**Priority:** P0
**Depends on:** 00-GO-01, 00-GO-02, 00-GO-03

**Goal:** Implement core CLI commands

**Scope Files:**
- `cmd/sdp/commands/init.go` (enhance existing)
- `cmd/sdp/commands/doctor.go` (enhance existing)
- `cmd/sdp/commands/build.go`

**Acceptance Criteria:**
- AC1: `sdp init` creates `.claude/` directory structure
- AC2: `sdp doctor` checks Python, Beads, git availability
- AC3: `sdp build <ws-id>` executes TDD cycle + updates Beads
- AC4: All commands support `--verbose` and `--quiet` flags
- AC5: Help text displays for all commands

**Implementation Notes:**
- Reuse existing Go code from `sdp-plugin/cmd/sdp/`
- Add cobra commands for build (new)
- Integrate TDD runner + Beads wrapper

**Integration Test:**
```bash
sdp init test-project
cd test-project
sdp doctor  # Should pass
# Create workstream
sdp build 00-001-01  # Should execute TDD
```

---

## Phase 2: Quality Automation (Week 4-5)

### 00-GO-05: Quality Gates (Parallel Execution)
**Size:** MEDIUM (700 LOC)
**Priority:** P1
**Depends on:** 00-GO-04

**Goal:** Automated quality gates with parallel execution (goroutines)

**Scope Files:**
- `internal/quality/checker.go`
- `internal/quality/gate.go`
- `cmd/sdp/commands/quality.go`

**Acceptance Criteria:**
- AC1: Run mypy, ruff, coverage in parallel (goroutines)
- AC2: File size check (<200 LOC per file)
- AC3: `sdp quality check` displays pass/fail per gate
- AC4: Exit code 1 if any gate fails
- AC5: Total execution time <8 seconds (vs Python's 8.5s)

**Implementation Notes:**
- Use `sync.WaitGroup` for parallel execution
- Channels for aggregating results
- Individual gate functions: `runMypy()`, `runRuff()`, `runCoverage()`

**Integration Test:**
```bash
sdp quality check
✅ mypy: PASSED (0 errors)
✅ ruff: PASSED (0 warnings)
✅ coverage: 85% (target: 80%)
✅ file size: all files <200 LOC
```

**Recommendation from Deep-Thinking:**
> "Quality gates should run in parallel to reduce feedback delay from 30-60s to <10s."

---

### 00-GO-06: Quality Watcher (Real-Time Feedback)
**Size:** MEDIUM (500 LOC)
**Priority:** P1
**Depends on:** 00-GO-05

**Goal:** Background file watcher with real-time quality feedback

**Scope Files:**
- `internal/quality/watcher.go`
- `internal/quality/cache.go`

**Acceptance Criteria:**
- AC1: Watcher monitors `.py` files via fsnotify
- AC2: Debounce file changes (500ms delay)
- AC3: Run incremental quality checks on changed files
- AC4: Cache results in `.quality-cache.json`
- AC5: Desktop notification on quality failures

**Implementation Notes:**
- Use `github.com/fsnotify/fsnotify` for cross-platform watching
- Cache for fast lookups (avoid redundant checks)
- Desktop notifier: `terminal-notifier` (macOS), `notify-send` (Linux)

**Integration Test:**
```bash
# Terminal 1
sdp quality watch &
[Watching for file changes...]

# Terminal 2
vim src/module.py  # Edit file

# Terminal 1 (desktop notification)
⚠️ Quality gate FAILED: src/module.py:42: missing type hints
```

**Recommendation from Deep-Thinking:**
> "Real-time hooks during coding (not just at commit time) reduce manual gate running by 80%."

---

## Phase 3: Telemetry System (Week 6-7)

### 00-GO-07: Telemetry Collector
**Size:** MEDIUM (600 LOC)
**Priority:** P1
**Depends on:** 00-GO-04

**Goal:** Auto-capture execution metrics (git stats, coverage, friction)

**Scope Files:**
- `internal/telemetry/collector.go`
- `internal/telemetry/git_stats.go`

**Acceptance Criteria:**
- AC1: Capture files changed, LOC added/deleted from git
- AC2: Capture test coverage before/after implementation
- AC3: Detect friction points (quality failures, missing type hints)
- AC4: Append telemetry to workstream frontmatter automatically
- AC5: `sdp telemetry scan` generates metrics summary

**Implementation Notes:**
- Git stats via `git diff --stat` and `git diff --shortstat`
- Coverage via `pytest --cov-report=json`
- Friction detection via quality cache analysis

**Integration Test:**
```bash
sdp build 00-001-01
# ... TDD cycle executes ...

# Workstream file automatically updated:
---
telemetry:
  files_changed: [src/module.py, tests/test_module.py]
  loc_added: 45
  loc_deleted: 12
  test_coverage_before: 0
  test_coverage_after: 87
  friction_points:
    - "Quality gate failed 2x due to missing type hints"
  outcome: "success"
---
```

**Recommendation from Deep-Thinking:**
> "Auto-capture telemetry eliminates manual documentation (zero friction) and enables pattern analysis."

---

### 00-GO-08: Telemetry Analyzer (Pattern Detection)
**Size:** MEDIUM (500 LOC)
**Priority:** P2
**Depends on:** 00-GO-07

**Goal:** Analyze telemetry patterns, generate insights

**Scope Files:**
- `internal/telemetry/analyzer.go`
- `internal/telemetry/patterns.go`
- `cmd/sdp/commands/telemetry.go`

**Acceptance Criteria:**
- AC1: `sdp telemetry insights --days 30` generates weekly report
- AC2: Identify top 5 friction points with occurrence counts
- AC3: Generate workstream proposals for high-severity patterns
- AC4: Calculate skill usage statistics (@build, @review, @oneshot)
- AC5: Analyze escalation trends (via `bd history` or JSON)

**Implementation Notes:**
- Pattern detection: group telemetry events by friction point
- Severity levels: HIGH (≥7 occurrences), MEDIUM (4-6), LOW (<4)
- Workstream proposal generation: `00-900-01`, `00-900-02`, etc.

**Integration Test:**
```bash
sdp telemetry insights --days 30

# Telemetry Insights: Week of 2026-01-29

## Friction Points (Top 5)
1. Quality gate failed 2x due to missing type hints (7 occurrences)
2. Documentation inconsistent with actual workflow (5 occurrences)

## Suggested Protocol Improvements
1. **00-900-01**: Add pre-build type hint validation to guard system
2. **00-900-02**: Update docs/workflow-decision.md based on real usage

## Skill Usage
- @build: 42 invocations, 95% success rate
- @review: 28 invocations, 89% success rate
```

**Recommendation from Deep-Thinking:**
> "Telemetry-driven insights enable continuous protocol improvement (auto-generate workstream proposals)."

---

### 00-GO-09: Drift Detector (Documentation-Code Sync)
**Size:** MEDIUM (800 LOC)
**Priority:** P1
**Depends on:** 00-GO-04, 00-GO-07

**Goal:** Detect mismatches between Contract section and implementation

**Scope Files:**
- `internal/validators/drift/fingerprint.go`
- `internal/validators/drift/detector.go`

**Acceptance Criteria:**
- AC1: Extract fingerprint from Contract section (pre-build)
- AC2: Parse implementation code to extract actual functions/classes
- AC3: Calculate drift score (Jaccard similarity)
- AC4: Block commit if drift > threshold (0.5) without deviation declaration
- AC5: `sdp drift check <ws-id>` displays drift analysis

**Implementation Notes:**
- Fingerprint: contract SHA256, expected logic type (generic/business), function signatures
- Drift calculation: `1.0 - Jaccard_similarity(expected, actual)`
- Deviation declaration: structured YAML in frontmatter

**Integration Test:**
```bash
# Pre-build: extract fingerprint
sdp build 00-040-04 --fingerprint-only
✅ Fingerprint stored: contract_sha256=abc123...
  Expected logic: generic validator
  Expected functions: [validate(), apply_rules()]

# Post-build: detect drift
sdp build 00-040-04
⚠️ Drift detected: score 0.75 > threshold 0.5
   Expected: generic validator
   Actual: business logic in quality/models.py
   Declare deviation? [y/N]: y
   Deviation reason: Required for edge case handling in payment flow
   ✅ Workstream 00-040-04 completed with deviation
```

**Recommendation from Deep-Thinking:**
> "Documentation-code mismatch is core pain point. Drift detection catches it BEFORE implementation, not after."

---

## Phase 4: Multi-Agent + Polish (Week 8-9)

### 00-GO-10: Checkpoint System (Simplified)
**Size:** SMALL (200 LOC)
**Priority:** P2
**Depends on:** 00-GO-04

**Goal:** Simplified checkpoint system for @oneshot resume

**Scope Files:**
- `internal/checkpoint/repository.go`

**Acceptance Criteria:**
- AC1: Save checkpoint to `.sdp/checkpoints/<feature>.json`
- AC2: Load checkpoint by feature ID
- AC3: Track completed workstreams, current workstream, status
- AC4: `sdp oneshot --resume <feature>` continues from checkpoint
- AC5: Delete checkpoint after successful completion

**Implementation Notes:**
- JSON storage (no SQLite) — 200 LOC vs 800 LOC Python
- Simple schema: Feature, AgentID, Status, CompletedWS[], CurrentWS

**Recommendation from Deep-Thinking:**
> "Use Beads gates for coordination if possible, fallback to JSON for simplicity."

---

### 00-GO-11: Multi-Agent Orchestrator (Adaptive Parallelism)
**Size:** MEDIUM (600 LOC)
**Priority:** P2
**Depends on:** 00-GO-03, 00-GO-10

**Goal:** Simplified multi-agent execution with adaptive parallelism

**Scope Files:**
- `internal/unified/orchestrator/executor.go`
- `internal/unified/orchestrator/wave.go`

**Acceptance Criteria:**
- AC1: Execute workstreams in waves via `bd ready`
- AC2: Adaptive agent count: 2 agents (small), 5 (medium), 10 (large)
- AC3: Cascading failure detection (3+ consecutive failures)
- AC4: Checkpoint save/restore for resume capability
- AC5: `sdp oneshot <feature>` executes all workstreams autonomously

**Implementation Notes:**
- Simplified from Python (1000 LOC → 600 LOC)
- Delegate dependency resolution to `bd ready`
- Wave-based parallelism: execute ready tasks, wait, repeat

**Integration Test:**
```bash
sdp oneshot F01
✅ Wave 1: Executing 2 tasks in parallel
  ✓ 00-001-01 completed (2m 15s)
  ✓ 00-001-04 completed (1m 45s)

✅ Wave 2: Executing 1 task
  ✓ 00-001-02 completed (3m 30s)

✅ Wave 3: Executing 1 task
  ✓ 00-001-03 completed (2m 50s)

✅ Feature F01 completed (10m 20s total)
```

**Recommendation from Deep-Thinking:**
> "Adaptive parallelism reduces execution time for large features (10+ workstreams)."

---

### 00-GO-12: CLI Polish + Cross-Platform Builds
**Size:** SMALL (300 LOC)
**Priority:** P2
**Depends on:** All previous workstreams

**Goal:** Production-ready CLI with polish

**Scope Files:**
- `Makefile`
- `cmd/sdp/commands/completion.go`
- Various cosmetic improvements

**Acceptance Criteria:**
- AC1: Shell completion (bash, zsh, fish)
- AC2: Color output (lipgloss library)
- AC3: Progress bars for long operations
- AC4: Cross-platform builds (Linux, macOS, Windows)
- AC5: Binary size ≤20MB (stripped)

**Implementation Notes:**
- `make build-all` produces 5 binaries
- Stripped binaries: `-ldflags "-s -w"`
- Upload to GitHub releases

**Integration Test:**
```bash
# Build all platforms
make build-all
✅ sdp-linux-amd64 (15MB)
✅ sdp-darwin-amd64 (14MB)
✅ sdp-darwin-arm64 (13MB)
✅ sdp-windows-amd64.exe (16MB)

# Test shell completion
sdp completion bash > /etc/bash_completion.d/sdp
source ~/.bashrc
sdp bu<TAB>  # Autocomplete to "sdp build"
```

---

## Phase 5: Python Cleanup (Week 10)

### 00-GO-13: Python Code Removal
**Size:** LARGE (3000 LOC deleted)
**Priority:** P0
**Depends on:** 00-GO-12 (all Go features verified)

**Goal:** Remove obsolete Python code after successful migration

**Scope Files:**
- Delete: `src/sdp/` (entire directory)
- Delete: `pyproject.toml`, `poetry.lock`, `tests/`
- Keep: `.claude/skills/`, `.claude/agents/`, `docs/`, `prompts/`

**Acceptance Criteria:**
- AC1: All Python code deleted (25K LOC removed)
- AC2: Go binary works flawlessly (all tests pass)
- AC3: Documentation updated (no Python references)
- AC4: `.gitignore` updated (no Python artifacts)
- AC5: Python code archived to `python-legacy` branch

**Implementation Notes:**
- **CRITICAL:** Verify Go works perfectly before deletion
- Create backup branch: `git checkout -b python-legacy && git push origin python-legacy`
- Tag release: `v1.0.0-go`

**Files to DELETE:**
```
src/sdp/adapters/          # Not needed (skills/ handle integration)
src/sdp/beads/             # Replaced by internal/beads/wrapper.go
src/sdp/cli/               # Replaced by cmd/sdp/
src/sdp/core/              # Replaced by internal/core/
src/sdp/design/            # Beads handles dependencies
src/sdp/doctor.py          # Replaced by cmd/sdp/commands/doctor.go
src/sdp/errors/            # Go error handling
src/sdp/extensions/        # Not needed
src/sdp/feature/           # skills/ handle this
src/sdp/health_checks/     # sdp doctor replacement
src/sdp/init_*.py          # sdp init replacement
src/sdp/prd/               # Keep for now (product vision)
src/sdp/quality/           # internal/quality/ replacement
src/sdp/report/            # telemetry handles this
src/sdp/schema/            # internal/core/schema.go replacement
src/sdp/tdd/               # internal/tdd/ replacement
src/sdp/traceability/      # drift detector covers this
src/sdp/unified/           # Partially replaced (checkpoint simplified)
src/sdp/validators/        # internal/validators/ replacement
```

**Files to KEEP:**
```
.claude/skills/            # Core SDP functionality (prompts)
.claude/agents/            # Multi-agent prompts
docs/                      # Documentation
hooks/                     # Git hooks (may need Go rewrite later)
prompts/commands/          # Skill definitions (markdown)
```

**Execution Script:**
```bash
#!/bin/bash
# cleanup-python.sh

# 1. Verify Go works
./sdp doctor || { echo "Go binary broken! Aborting."; exit 1; }
./sdp --version

# 2. Backup Python code
git checkout -b python-legacy
git push origin python-legacy

# 3. Return to main
git checkout main

# 4. Delete Python code
git rm -r src/sdp/
git rm pyproject.toml poetry.lock
git rm tests/

# 5. Update documentation
# Edit README.md, CLAUDE.md (remove Python references)

# 6. Commit
git commit -m "feat: migrate to Go implementation

- Replace Python CLI with Go binary
- Remove 25K LOC Python code
- Leverage Beads CLI for task tracking
- Single static binary for easy deployment

BREAKING CHANGE: Python runtime no longer required"

# 7. Tag release
git tag v1.0.0-go
git push origin main --tags

echo "✅ Python cleanup complete. Go binary is now the only implementation."
```

---

## Dependency Graph

```
00-GO-01 (Setup)
    ├─→ 00-GO-02 (TDD Runner)
    │       ├─→ 00-GO-04 (CLI Commands)
    │               ├─→ 00-GO-05 (Quality Gates)
    │               │       ├─→ 00-GO-06 (Quality Watcher)
    │               │
    │               ├─→ 00-GO-07 (Telemetry Collector)
    │               │       ├─→ 00-GO-08 (Telemetry Analyzer)
    │               │       └─→ 00-GO-09 (Drift Detector)
    │               │
    │               ├─→ 00-GO-10 (Checkpoint System)
    │               │       └─→ 00-GO-11 (Multi-Agent Orchestrator)
    │               │               └─→ 00-GO-12 (CLI Polish)
    │               │                       └─→ 00-GO-13 (Python Cleanup)
    │
    └─→ 00-GO-03 (Beads Wrapper)
            └─→ 00-GO-04 (CLI Commands)
                    (continues as above)
```

**Critical Path:**
```
00-GO-01 → 00-GO-02 → 00-GO-03 → 00-GO-04 → 00-GO-07 → 00-GO-09 → 00-GO-12 → 00-GO-13
```

**Parallel Work:**
- `00-GO-05` (Quality Gates) can parallelize with `00-GO-07` (Telemetry)
- `00-GO-06` (Quality Watcher) depends on `00-GO-05` only
- `00-GO-08` (Telemetry Analyzer) depends on `00-GO-07` only

---

## Timeline

### Week 1-3: Foundation (MVP)
- 00-GO-01: Setup + Parser (MEDIUM)
- 00-GO-02: TDD Runner (MEDIUM)
- 00-GO-03: Beads Wrapper (SMALL)
- 00-GO-04: CLI Commands (MEDIUM)

**Deliverable:** Working `sdp build` command with TDD cycle

### Week 4-5: Quality Automation
- 00-GO-05: Quality Gates (MEDIUM)
- 00-GO-06: Quality Watcher (MEDIUM)

**Deliverable:** Real-time quality feedback + desktop notifications

### Week 6-7: Telemetry System
- 00-GO-07: Telemetry Collector (MEDIUM)
- 00-GO-08: Telemetry Analyzer (MEDIUM)
- 00-GO-09: Drift Detector (MEDIUM)

**Deliverable:** Auto-capture metrics + drift detection + insights

### Week 8-9: Multi-Agent + Polish
- 00-GO-10: Checkpoint System (SMALL)
- 00-GO-11: Multi-Agent Orchestrator (MEDIUM)
- 00-GO-12: CLI Polish (SMALL)

**Deliverable:** Production-ready binary with @oneshot support

### Week 10: Python Cleanup
- 00-GO-13: Python Code Removal (LARGE)

**Deliverable:** Clean codebase, Go-only implementation

---

## Success Metrics

### Feature Parity
- ✅ All Python CLI commands work in Go
- ✅ All skills work with Go binary
- ✅ Beads integration functional
- ✅ @oneshot executes autonomously

### Performance
- ✅ `sdp build` ≤ Python latency (3.5s target)
- ✅ Quality gates <8s (parallel execution)
- ✅ Binary size ≤20MB

### Quality
- ✅ Test coverage ≥80%
- ✅ No race conditions (`go test -race`)
- ✅ Lint passed (`golangci-lint run`)

### Adoption
- ✅ 5 pilot projects successful
- ✅ Zero data loss (checkpoints, telemetry)
- ✅ User satisfaction ≥4/5 stars

---

## Risk Mitigation

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| **Go learning curve** | Medium | Medium | - Start with simple wrappers<br>- Use Go resources (Effective Go)<br>- Pair programming for review |
| **Beads not installed** | High | Low | - Fallback to JSON storage<br>- Auto-install prompt<br>- Clear error messages |
| **Performance regression** | Low | High | - Benchmark critical paths<br>- Optimize hotspots<br>- Profile before/after |
| **Python data loss** | Low | Critical | - Verify Go works perfectly<br>- Backup to python-legacy branch<br>- Tag before deletion |
| **Breaking skills** | Medium | High | - Test all skills with Go<br>- Keep skill format unchanged<br>- Document migration |

---

## Rollback Plan

If Go version has critical issues:

1. **Stop using Go binary**
   ```bash
   rm /usr/local/bin/sdp

   # Restore Python version
   git checkout python-legacy
   pip install -e .
   ```

2. **Report bug**
   - GitHub issue with reproduction
   - Tag: `critical`, `go-migration`

3. **Fix and retry**
   - Debug in development branch
   - Add regression test
   - Retry cutover in 1 week

---

## Next Steps

1. **Review and approve this roadmap** — Ensure alignment with goals
2. **Setup Beads integration** — Run `sdp beads migrate docs/plans/ --real`
3. **Start Phase 1** — Execute 00-GO-01 through 00-GO-04
4. **Weekly reviews** — Check progress against timeline
5. **Adjust if needed** — Reprioritize based on learnings

---

**Document Version:** 1.0
**Status:** Ready for Beads Migration
**Next Action:** Run `sdp beads migrate` to create Beads tasks
