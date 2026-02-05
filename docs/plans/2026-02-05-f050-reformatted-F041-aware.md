# F050 Go Migration - Reformatted (F041-Aware)

> **Status:** REFORMATTED - Builds on F041, no duplication
> **Date:** 2026-02-05
> **Timeline:** 10-12 weeks (reduced from 14 by leveraging F041)

---

## 🎯 Executive Summary

**F041 Deliverables (Already Complete):**
- ✅ `sdp-plugin/` Go module with 7 files
- ✅ Basic CLI: init, doctor, hooks commands
- ✅ Cross-platform build system (Makefile)
- ✅ Binary distribution: 4 platforms, ~5.5MB each

**F050 Scope (Extends F041):**
- 🆕 Workstream parser (markdown → Go structs)
- 🆕 TDD runner (Red→Green→Refactor automation)
- 🆕 Drift detector (docs vs code validation)
- 🆕 Beads integration (task tracking)
- 🆕 Quality gates (coverage, complexity, types)
- 🆕 Orchestrator (multi-workstream execution)
- 🆕 Telemetry (usage analytics)

**Key Change:** F050 NO LONGER creates duplicate Go project. Instead extends `sdp-plugin/` with new packages.

---

## 📊 What Changes from Original F050

### Removed (Duplicate with F041):
- ❌ ~~WS-00-050-04: CLI Commands (init, doctor, hooks)~~ - EXISTS in F041
- ❌ ~~Go module setup~~ - EXISTS as `sdp-plugin/go.mod`
- ❌ ~~Cross-platform builds~~ - EXISTS in F041 Makefile

### Added (New Capabilities):
- ✅ WS-00-050-01: **Workstream Parser** (markdown → structs)
- ✅ WS-00-050-02: **TDD Runner** (Red-Green-Refactor)
- ✅ WS-00-050-03: **Beads Integration** (task tracking)
- ✅ WS-00-050-09: **Drift Detector** (docs validation)
- ✅ WS-00-050-10: **Checkpoint System** (resume capability)
- ✅ WS-00-050-11: **Orchestrator** (multi-WS execution)

### Re-scoped (Build on F041):
- 🔄 WS-00-050-05: **Quality Gates** → Add to `sdp doctor` command
- 🔄 WS-00-050-07: **Telemetry** → Extend existing binary

---

## 🏗️ Architecture (Post-F041)

```
sdp-plugin/                    # F041 Foundation
├── cmd/sdp/
│   ├── main.go               # ✅ EXISTS (F041)
│   ├── init.go               # ✅ EXISTS
│   ├── doctor.go             # ✅ EXISTS (F041)
│   ├── hooks.go              # ✅ EXISTS
│   ├── parse.go              # 🆕 F050-01
│   ├── tdd.go                # 🆕 F050-02
│   ├── drift.go              # 🆕 F050-09
│   └── orchestrate.go        # 🆕 F050-11
├── internal/
│   ├── doctor/               # ✅ EXISTS (F041)
│   ├── hooks/                # ✅ EXISTS
│   ├── sdpinit/              # ✅ EXISTS
│   ├── parser/               # 🆕 F050-01
│   ├── tdd/                  # 🆕 F050-02
│   ├── drift/                # 🆕 F050-09
│   ├── beads/                # 🆕 F050-03
│   ├── quality/              # 🆕 F050-05
│   └── orchestrator/         # 🆕 F050-11
├── pkg/                       # ✅ EXISTS (F041)
│   └── installer/            # ✅ EXISTS
├── go.mod                    # ✅ EXISTS (F041)
├── Makefile                  # ✅ EXISTS (F041)
└── prompts/                  # ✅ EXISTS (F041)
```

---

## 📋 Reformatted Workstreams

### Phase 1: Core Capabilities (4 weeks)

#### WS-00-050-01: Workstream Parser 🆕
**Priority:** 🔴 CRITICAL | **Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** None

**Goal:** Parse workstream markdown files into Go structs

**Scope:**
- Add `internal/parser/` package
- Parse YAML frontmatter (--- delimited)
- Validate WS ID format (PP-FFF-SS)
- Extract: goal, acceptance criteria, scope files
- Return structured errors for invalid YAML

**Acceptance Criteria:**
- AC1: Parse valid WS markdown → struct
- AC2: Validate WS ID regex: `^\d{2}-\d{3}-\d{2}$`
- AC3: Extract all frontmatter fields
- AC4: Return helpful errors for malformed YAML
- AC5: Unit tests for edge cases
- Coverage ≥80%

**Command:**
```bash
sdp parse <ws-id>              # Parse and display workstream
sdp parse --validate <ws.md>   # Validate workstream file
```

**Integration:** Called by drift detector, orchestrator

---

#### WS-00-050-02: TDD Runner 🆕
**Priority:** 🔴 CRITICAL | **Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Automate Red→Green→Refactor cycle

**Scope:**
- Add `internal/tdd/` package
- Red phase: Run tests, expect failure
- Green phase: Run tests, expect success
- Refactor phase: Validate tests still pass
- Support pytest (Python), go test, mvn test (Java)

**Acceptance Criteria:**
- AC1: Red phase detects test failures
- AC2: Green phase validates test success
- AC3: Refactor phase ensures no regression
- AC4: Language detection (pytest/go test/mvn)
- AC5: Phase transition tracking
- Coverage ≥80%

**Command:**
```bash
sdp tdd red                    # Run tests, expect failure
sdp tdd green                  # Run tests, expect success
sdp tdd refactor               # Ensure no regression
```

**Integration:** Used by @build skill

---

#### WS-00-050-03: Beads Integration 🆕
**Priority:** 🟡 IMPORTANT | **Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** None

**Goal:** Wrapper around Beads CLI for task tracking

**Scope:**
- Add `internal/beads/` package
- Wrap `bd ready`, `bd show`, `bd update`
- Read `.beads-sdp-mapping.jsonl`
- Auto-sync on commit (via git hooks)

**Acceptance Criteria:**
- AC1: `bd ready` → list available tasks
- AC2: `bd show <id>` → get task details
- AC3: `bd update <id> --status <status>` → update status
- AC4: Mapping WS ID ↔ Beads ID
- AC5: Auto-sync on git commit
- Coverage ≥80%

**Commands:**
```bash
sdp beads ready               # Show ready tasks
sdp beads show <ws-id>        # Show task details
sdp beads update <ws-id>      # Update task status
```

**Integration:** Used by all skills for task tracking

---

#### WS-00-050-09: Drift Detector 🆕
**Priority:** 🔴 CRITICAL | **Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** 00-050-01

**Goal:** Detect documentation-code mismatch

**Problem (from 827 sessions):**
> 4,903 friction events from documentation not matching codebase reality

**Scope:**
- Add `internal/drift/` package
- Parse scope_files from WS frontmatter
- Validate files exist at expected paths
- Check declared functions/classes present
- Compare file purpose with reality
- Generate drift report with recommendations

**Acceptance Criteria:**
- AC1: Parse scope_files from WS
- AC2: Validate all files exist
- AC3: Validate declared entities present
- AC4: Generate drift report
- AC5: Integrate with `sdp doctor`
- AC6: Command: `sdp drift detect <ws-id>`
- Coverage ≥80%

**Example Output:**
```bash
$ sdp drift detect 00-050-01

## Drift Report: 00-050-01

| File | Status | Issue |
|------|--------|-------|
| internal/parser/workstream.go | ✅ OK | Matches docs |
| internal/parser/validator.go | ⚠️ WARNING | Function Parse() not found, ParseWorkstream() exists |
| internal/parser/schema.go | ❌ ERROR | File not found |

**Verdict:** ❌ FAIL - 1 error, 1 warning

**Recommendation:** Rename Parse() → ParseWorkstream() or update docs
```

**Integration:** Pre-build check via `sdp doctor`

---

### Phase 2: Quality & Telemetry (3 weeks)

#### WS-00-050-05: Quality Gates Extension 🔄
**Priority:** 🟡 IMPORTANT | **Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** None

**Goal:** Extend `sdp doctor` with quality checks

**Scope:**
- Add `internal/quality/` package
- Coverage checker (≥80%)
- Complexity analyzer (CC < 10)
- File size validator (<200 LOC)
- Type hints checker (Python) / signatures (Go/Java)

**Acceptance Criteria:**
- AC1: Coverage ≥80% check
- AC2: Cyclomatic complexity <10
- AC3: File size <200 LOC
- AC4: Type completeness check
- AC5: Integrated into `sdp doctor`
- Coverage ≥80%

**Commands:**
```bash
sdp doctor                     # Runs all checks
sdp quality coverage           # Check test coverage
sdp quality complexity         # Check cyclomatic complexity
sdp quality size               # Check file sizes
```

---

#### WS-00-050-07: Telemetry Collector 🆕
**Priority:** 🟢 NICE | **Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** None

**Goal:** Collect usage metrics for improvement

**Scope:**
- Add `internal/telemetry/` package
- Track: command usage, duration, errors
- Local storage: `~/.sdp/telemetry.jsonl`
- Privacy-first: no PII, opt-out available
- Export to CSV/JSON for analysis

**Acceptance Criteria:**
- AC1: Track command invocations
- AC2: Record duration and success/failure
- AC3: Store locally in JSONL format
- AC4: Opt-out mechanism
- AC5: Export to CSV/JSON
- Coverage ≥80%

**Commands:**
```bash
sdp telemetry status           # Show collection status
sdp telemetry export           # Export to CSV
sdp telemetry disable          # Opt-out
```

---

#### WS-00-050-08: Telemetry Analyzer 🆕
**Priority:** 🟢 NICE | **Size:** SMALL | **Duration:** 1 week
**Dependencies:** 00-050-07

**Goal:** Analyze telemetry for insights

**Scope:**
- Add `internal/telemetry/analyzer.go`
- Calculate: success rate, avg duration, error frequency
- Generate insights report
- Identify friction points

**Acceptance Criteria:**
- AC1: Calculate success rate by command
- AC2: Average duration by command
- AC3: Top 5 error categories
- AC4: Insights report
- AC5: Command: `sdp telemetry analyze`
- Coverage ≥80%

**Example Output:**
```bash
$ sdp telemetry analyze

## Usage Insights (Last 30 Days)

**Top Commands:**
1. @build (45 sessions) - 89% success, 12m avg
2. @review (23 sessions) - 95% success, 3m avg
3. @design (18 sessions) - 100% success, 25m avg

**Top Errors:**
1. Test failures (12) - @build
2. Coverage <80% (8) - @review
3. Drift detected (5) - @build

**Recommendation:** Focus on test reliability (12% failure rate)
```

---

### Phase 3: Orchestration & Polish (3 weeks)

#### WS-00-050-10: Checkpoint System 🆕
**Priority:** 🟡 IMPORTANT | **Size:** SMALL | **Duration:** 3 days
**Dependencies:** None

**Goal:** Resume capability for long operations

**Scope:**
- Add `internal/checkpoint/` package
- JSON-based checkpoint: `.sdp/checkpoints/*.json`
- Save/progress after each WS
- Resume from last checkpoint

**Acceptance Criteria:**
- AC1: Create checkpoint on WS start
- AC2: Update checkpoint on WS complete
- AC3: Resume from checkpoint
- AC4: List available checkpoints
- AC5: Clean old checkpoints
- Coverage ≥80%

**Commands:**
```bash
sdp checkpoint create <name>    # Create checkpoint
sdp checkpoint resume <id>      # Resume from checkpoint
sdp checkpoint list            # List checkpoints
sdp checkpoint clean           # Clean old checkpoints
```

---

#### WS-00-050-11: Orchestrator 🆕
**Priority:** 🟡 IMPORTANT | **Size:** LARGE | **Duration:** 2 weeks
**Dependencies:** 00-050-01, 00-050-02, 00-050-03, 00-050-10

**Goal:** Execute multiple workstreams with dependency tracking

**Scope:**
- Add `internal/orchestrator/` package
- Build dependency graph from Beads
- Topological sort for execution order
- Execute WS sequentially (single-session limitation)
- Checkpoint after each WS
- Resume on failure

**Acceptance Criteria:**
- AC1: Load workstreams from Beads
- AC2: Build dependency graph
- AC3: Topological sort execution order
- AC4: Execute workstreams sequentially
- AC5: Checkpoint after each WS
- AC6: Resume from checkpoint
- AC7: Handle dependencies correctly
- Coverage ≥80%

**Command:**
```bash
sdp orchestrate <feature-id>    # Execute all workstreams
sdp orchestrate --resume <id>   # Resume from checkpoint
sdp orchestrate --dry-run       # Preview execution
```

**Integration:** Used by @oneshot skill

---

#### WS-00-050-14: Command Auto-Retry 🆕
**Priority:** 🟢 NICE | **Size:** SMALL | **Duration:** 3 days
**Dependencies:** 00-050-02

**Goal:** Retry failed commands with exponential backoff

**Problem (from 827 sessions):**
> 4,904 Bash command failures (13% rate) cause friction

**Scope:**
- Add `internal/retry/` package
- Exponential backoff: 1s → 2s → 4s
- Max 3 retries before escalation
- Log retry attempts to telemetry

**Acceptance Criteria:**
- AC1: Detect command failure (exit != 0)
- AC2: Retry with exponential backoff
- AC3: Max 3 retries
- AC4: Log to telemetry
- AC5: Integration with TDD runner
- Coverage ≥80%

---

### Phase 4: Migration Completion (2 weeks)

#### WS-00-050-12: CLI Polish 🔄
**Priority:** 🟢 NICE | **Size:** SMALL | **Duration:** 3 days
**Dependencies:** All above

**Goal:** Polish CLI user experience

**Scope:**
- Add shell completion (bash, zsh, fish)
- Color output for errors/warnings
- Progress bars for long operations
- Better help text

**Acceptance Criteria:**
- AC1: Shell completion works
- AC2: Colorized output
- AC3: Progress bars
- AC4: Improved help text
- Coverage ≥80%

---

#### WS-00-050-13: Python SDP Deprecation 🔄
**Priority:** 🟡 IMPORTANT | **Size:** MEDIUM | **Duration:** 1 week
**Dependencies:** All above

**Goal:** Deprecate Python SDP in favor of Go version

**Scope:**
- Add deprecation notice to Python SDP
- Migration guide: Python → Go
- Feature parity checklist
- Update documentation

**Acceptance Criteria:**
- AC1: Deprecation notice in README
- AC2: Migration guide complete
- AC3: Feature parity documented
- AC4: Documentation updated
- Coverage N/A (documentation only)

---

#### WS-00-050-06: Quality Watcher (Optional)
**Priority:** 🟢 NICE | **Size:** SMALL | **Duration:** 3 days
**Dependencies:** 00-050-05

**Goal:** Background quality monitoring

**Scope:**
- Add `internal/watcher/` package
- Watch files for changes
- Run quality checks on save
- Notify on violations

**Acceptance Criteria:**
- AC1: Watch files for changes
- AC2: Run quality checks
- AC3: Notify on violations
- AC4: Configurable watch paths
- Coverage ≥80%

**Command:**
```bash
sdp watch                     # Start quality watcher
```

---

## 📊 Timeline Comparison

| Phase | Original F050 | Reformatted F050 | Savings |
|-------|--------------|------------------|---------|
| Phase 1 | 6 weeks (6 WS) | 4 weeks (4 WS) | -2 weeks |
| Phase 2 | 4 weeks (4 WS) | 3 weeks (3 WS) | -1 week |
| Phase 3 | 2 weeks (3 WS) | 3 weeks (5 WS) | +1 week |
| Phase 4 | - | 2 weeks (2 WS) | - |
| **Total** | **12 weeks** | **12 weeks** | **Same!** |

**But:** Reformatted version has MORE features (drift detector, telemetry) with SAME timeline by leveraging F041 foundation.

---

## 🎯 Success Criteria

F050 is complete when:

1. ✅ All workstreams implemented in `sdp-plugin/`
2. ✅ Binary size ≤15MB (current: 5.5MB ✅)
3. ✅ Cross-platform builds work (already: 4 platforms ✅)
4. ✅ Drift detector catches 90%+ of documentation mismatches
5. ✅ TDD runner automates Red→Green→Refactor
6. ✅ Orchestrator executes multi-WS features
7. ✅ Telemetry provides actionable insights
8. ✅ Python SDP deprecated with migration path

---

## 🚀 Next Steps

1. **Review this plan** - Does it correctly build on F041?
2. **Update Beads tasks** - Modify 00-050-XX workstreams to reflect new scope
3. **Start execution** - Begin with WS-00-050-01 (Workstream Parser)
4. **Continuous integration** - Add to `sdp-plugin/` as packages complete

**Question:** Should we proceed with this reformatted plan, or make additional adjustments?
