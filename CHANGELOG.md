# Changelog

All notable changes to the Spec-Driven Protocol (SDP).

## [0.9.8] - 2026-02-26

### Skills Sync, Beads Integration

- **@build** — Post-build `bd close` for beads in WS frontmatter; batch syntax `/build 00-XXX-YY..ZZ`
- **@design** — Pre-draft check, bead verification, default-in-scope
- **@review** — Handoff block when CHANGES_REQUESTED
- **CLAUDE.md** — Sync with sdp_dev/AGENTS.md (placement, "продолжай" convention)

### Evidence + Checkpoint Commit, Guard Strict

- **@build** — Step 3b: commit evidence + checkpoint after `sdp-orchestrate --advance`
- **.gitignore** — `!.sdp/evidence/`, `!.sdp/checkpoints/` (evidence must be committed with PR)
- **sdp guard** — `deactivate` without `|| true`; `activate` in pre-build without `|| true` (blocking)
- **go-ci** — `sdp contract validate` without `|| true` (blocking)

### Installer

- **install.sh** — Fail when no sha256 tool (security); do not skip checksum verification
- **install-project.sh** — Clear retry command on CLI install failure
- **README, QUICKSTART** — Document binary-only vs full project install modes

### Release, Demo, Init

- **Release workflow** — In sdp repo: protocol-e2e with GLM_API_KEY, then sdp CLI release
- **sdp init** — Executable-relative paths for global install; `--auto` in protocol-e2e Phase 2b
- **Phase 5** — Fail when GLM_API_KEY not set (opencode LLM required)
- **demo/** — demo.tape, demo.gif, run-demo.sh; examples/beads-viz-demo

---

### Coverage Context + Protocol-Only Docs Boundary

**New:**
- **CheckCoverage(ctx)** — Coverage checker accepts `context.Context` for cancellation support.
- **Intent schema docs** — `docs/intent/README.md` documents intent specification format.
- **Schema path consistency** — Schemas served from `schema/` at SDP root.
- **ws-verdict, review-verdict** — JSON schemas for build and review outputs.

**Protocol boundary:**
- Removed non-protocol docs (workstreams, roadmap, plans, reviews, decisions, specs) — migrated to sdp_dev.
- Kept minimal `docs/reference/` (PRINCIPLES, GLOSSARY, build/design/review specs).
- Added `docs/README.md` index for protocol documentation.

**Fixes:**
- **intent.schema.json** — Restored valid JSON (was corrupted).
- **Windows build** — Added `lock_windows.go` (no-op flock); evidence layer requires UNIX.
- **Lint** — gofmt, prealloc, gocognit, errcheck nolints.
- **TestVerifyCmd** — Replaced hanging integration test with `TestVerifyCmdConstructed` unit test.

**Audit remediation:**
- Verifier interface abstraction (CoverageChecker, PathValidator, CommandRunner).
- Parser frontmatter fix (Index for both `---` delimiters).
- Configurable coverage timeouts, EmitSync docstrings, Writer hash-chain atomicity.
- Guard rules path validation, ParseDependencies safe fallback, ctx propagation.

**Install:**
- `SDP_REF=v0.9.8` for testing. OpenCode: `SDP_IDE=opencode` or `SDP_IDE=all`.

---

## [0.9.7] - 2026-02-24

### Phase 0 + protocol E2E: Skills, Schema, Constraints, Full Protocol

**New:**
- **Coding Workflow Predicate (in-toto v1)** — `schema/coding-workflow-predicate.schema.json` for attestations. Predicate type: `https://sdp.dev/attestation/coding-workflow/v1`. See [docs/attestation/coding-workflow-v1.md](docs/attestation/coding-workflow-v1.md).
- **@feature --auto** — Generate workstreams directly from roadmap. Skip full discovery for features already in ROADMAP.md.
- **@design** — Workstream file format with Scope Files, beads mapping, INDEX.md update. Required sections documented.
- **PreToolUse constraint enforcement** — `sdp-guard --check-constraints` integration. Reads `.sdp/agent-constraints.yaml` for phase-specific rules (scope, force-push, destructive git).

**Phase 0:**
- Removed phantom guard CLI refs (context, branch, complete, finding)
- Slim @oneshot skill; outer loop via sdp-orchestrate
- Language-agnostic skills: quality gates per AGENTS.md, `master` not `dev`
- @build scope fix; stripped evidence boilerplate
- Deleted `help.md`, `init.md`; compressed deploy, review, implementer

**Install:**
- `SDP_REF` env var for branch/ref (e.g. `SDP_REF=v0.9.7` for testing)
- OpenCode/Windsurf: `SDP_IDE=opencode` or `SDP_IDE=all`

---

## [0.9.4] - 2026-02-18

### Patch Release

**Fixes:**
- Fixed installer portability so `curl ... | sh` works and installs the CLI binary by default.
- Fixed source build path by forcing `CGO_ENABLED=0` in `sdp-plugin/Makefile` to avoid Xcode dependency on macOS.
- Restored `sdp init --guided` compatibility as an alias to interactive mode.

---

## [0.9.3] - 2026-02-17

### Patch Release

**Fixes:**
- Fixed TempDir cleanup in prototype tests (CI stability)
- Fixed Go Release workflow (GPG signing, goreleaser v2 config)
- Added pre-release dry-run check on every PR

---

## [0.9.2] - 2026-02-17

### Patch Release

**New:**
- IDE selection in installer: `SDP_IDE=claude|cursor|opencode`
- OpenCode/Windsurf integration via `.opencode/` directory

**Fixes:**
- Fixed nil pointer panic in `NewWizard` (flaky TestInitCommand)

---

## [0.9.1] - 2026-02-16

### Patch Release

**Improvements:**
- **One-liner installer:** `curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash`
- **OpenCode integration:** `.opencode/` directory with skills, agents, commands
- **Cross-platform sync:** Commands available for Claude Code, Cursor, OpenCode, Windsurf

**Fixes:**
- Fixed Go version mismatch in `sdp-verify-dogfood.yml` (1.24 → 1.26)

**Dependencies:**
- Bump `github.com/spf13/cobra` from 1.8.0 to 1.10.2
- Bump `actions/upload-artifact` from 4 to 6

---

## [0.9.0] - 2026-02-16

### M1 Milestone - UX Excellence & Intelligent Assistance

**Theme:** Enhanced Developer Experience with Smart Recovery and Guidance

This release focuses on UX improvements, intelligent next-step recommendations, structured error handling, and guided onboarding.

### Highlights

- **Next-Step Engine:** Intelligent recommendations with confidence scoring
- **Error Taxonomy:** 38 structured error codes with recovery hints
- **Guided Onboarding:** Interactive wizard and headless mode
- **Self-Healing Doctor:** Automatic environment repair
- **Enhanced Evidence:** Full skills instrumentation

### Statistics

- **Features completed:** 8 (F054, F063, F064, F067, F068, F070, F075, F076)
- **Workstreams:** 57
- **Test coverage:** 80%+ (all packages)
- **M1 Status:** ✅ COMPLETE

---

## New Features

### F068: UX Foundation & First-Run Experience

Guided setup and improved user experience.

**Commands:**
- `sdp init --guided` - Step-by-step setup wizard with auto-fix
- `sdp init --auto` - Safe defaults for quick start
- `sdp demo` - Interactive feature walkthrough
- `sdp status --text/--json` - Quick status for scripts

**Components:**
- Preflight checks with auto-repair
- Quickstart templates
- Improved help text with user intent grouping

### F069: Next-Step Engine

Intelligent recommendation system for development workflow.

**Commands:**
- `sdp next` - Get next recommended action
- `sdp next --json` - Machine-readable output
- `sdp next --alternatives` - Show alternative actions

**Components:**
- Rule-based evaluation with deterministic tie-break
- Confidence scoring (0.0-1.0)
- Categories: execution, recovery, planning, information, setup
- Interactive loop: accept/refine/reject
- Quality metrics: acceptance rate, correction rate

### F070: Failure & Recovery UX

Structured error handling with recovery guidance.

**Commands:**
- `sdp diagnose` - Show error classes and codes
- `sdp diagnose [CODE]` - Get recovery steps for error
- `sdp diagnose --json` - Machine-readable output

**Error Taxonomy:**
| Class | Prefix | Description |
|-------|--------|-------------|
| Environment | ENV | Missing tools, permissions, filesystem |
| Protocol | PROTO | Invalid IDs, malformed files |
| Dependency | DEP | Blocked workstreams, cycles |
| Validation | VAL | Coverage, quality gates |
| Runtime | RUNTIME | External failures, timeouts |

**Components:**
- 38 error codes with messages and recovery hints
- Recovery playbooks with fast/deep path steps
- Diagnostics reports (JSON/text)

### F075: Self-Healing Doctor

Automatic environment repair.

**Commands:**
- `sdp doctor --repair` - Auto-fix detected issues
- `sdp doctor --deep` - Comprehensive diagnostics

**Repair Actions:**
- Install missing hooks
- Fix permissions
- Repair corrupted config
- Sync stale state

### F076: Guided Onboarding Wizard

Interactive project initialization.

**Commands:**
- `sdp init --guided` - Interactive wizard
- `sdp init --headless` - CI/CD mode (no prompts)
- `sdp init --project-type go|node|python|mixed` - Project type selection

**Components:**
- Safe defaults per project type
- Skills selection
- Evidence layer configuration
- JSON output for tooling

---

## Updated Features

### F056: Full Skills Instrumentation (Completed)

Evidence tracking for all skills.

**Events:**
- SkillEvent for @design, @idea, @build
- PlanEvent with approval_data
- DecisionEvent, LessonEvent

**Coverage:** 84.8%

### F065: Agent Git Safety Protocol (Completed)

Git safety with session validation.

**Components:**
- Session hash verification
- Safe git wrapper
- Branch protection
- Context recovery

**Coverage:** 83-90%

### F024: Unified Workflow (Completed)

Unified workflow orchestration.

**Coverage:** Full integration with all skills

---

## CLI Commands (New in 0.9.0)

| Command | Purpose |
|---------|---------|
| `sdp next` | Next recommended action |
| `sdp diagnose [CODE]` | Error lookup and recovery |
| `sdp demo` | Feature walkthrough |
| `sdp doctor --repair` | Auto-fix environment |
| `sdp init --guided` | Interactive setup |
| `sdp init --headless` | CI/CD setup |

---

## [0.8.0] - 2026-02-16

### Major Release - Multi-Agent Architecture + Go Implementation

**Theme:** From Python CLI to Intelligent Orchestration System

This release transforms SDP into a multi-agent orchestration system with autonomous execution, strategic planning, codebase analysis, long-term memory, and comprehensive evidence tracking.

### Highlights

- **Multi-agent orchestration:** 19+ specialized agents
- **Strategic planning:** @vision skill with 7 expert agents
- **Codebase analysis:** @reality skill with 8 expert agents
- **Long-term memory:** SQLite + FTS5 for context recovery
- **Evidence layer:** Hash-chained event log with CLI tools
- **Guard system:** Pre-edit scope enforcement
- **Parallel execution:** ~5x speedup
- **Go CLI:** Full Go implementation

### Statistics

- **Features completed:** 16 (F014, F024, F051-F067)
- **Workstreams:** 120+
- **Test coverage:** 68% → 80%+

---

## Features

### F014: Workflow Efficiency

Workflow optimization and efficiency improvements.

### F024: Unified Workflow

Unified workflow implementation with 18 workstreams covering end-to-end development process.

**Components:**
- Orchestrator with dependency graph and topological sort
- TeamManager for 100+ role management
- ApprovalGateManager for quality checkpoints
- Checkpoint save/resume for long-running features
- NotificationGateway for team updates
- FeatureCoordinator for @feature integration

**Commands:**
- `sdp orchestrate <feature-id>` - Execute all workstreams for a feature
- `sdp orchestrate resume <checkpoint-id>` - Resume from checkpoint

**Packages:**
- `internal/orchestrator/` - 83.4% coverage
- `internal/checkpoint/` - 84.4% coverage
- `internal/notification/` - 82.9% coverage

### F051: Long-term Memory System

Project memory for avoiding duplicated work.

**Commands:**
- `sdp memory index` - Index project artifacts into SQLite + FTS5
- `sdp memory search <query>` - Full-text search
- `sdp memory stats` - Show index statistics
- `sdp drift detect [ws_id]` - Detect code↔docs drift

### F052: Multi-Agent SDP + @vision + @reality

**@vision Skill:**
- 7 expert agents: product, market, technical, UX, business, growth, risk
- Generates VISION.md, PRD.md, ROADMAP.md

**@reality Skill:**
- 8 expert agents: architecture, quality, testing, security, performance, docs, debt, standards
- Generates reality report with tech debt tracking

**Parallel Execution:**
- Kahn's algorithm for dependency-aware parallelization
- Circuit breaker for fault tolerance
- Atomic checkpoint/resume

**Two-Stage Review:**
- Stage 1: Spec compliance
- Stage 2: Code quality (coverage >= 80%)

### F054: SDP Evidence Layer

Hash-chained event log for audit trail.

**Commands:**
- `sdp log show` - Show recent events with filters
- `sdp log trace` - Trace evidence chain
- `sdp log export` - Export as CSV/JSON
- `sdp log stats` - Show statistics

**Architecture:**
```
.sdp/log/events.jsonl  # Hash-chained event log
```

### F055: Compliance Design Doc

- Compliance documentation
- Threat model (THREAT-MODEL.md)
- GDPR/SOC2 compliance reference

### F056: Full Skills Instrumentation

Instrumentation for @review, @design, @idea and remaining skills with evidence tracking.

### F057: CLI plan/apply/log

**Commands:**
- `sdp plan "feature"` - Decompose feature into workstreams
- `sdp apply --ws <id>` - Execute workstreams
- `sdp log show/trace/export/stats` - Evidence operations

### F058: CI/CD GitHub Action

- SDP Verify Action for GitHub Actions
- PR evidence comments
- Release automation

### F059: Observability Bridge Design

- OpenTelemetry semantic conventions
- Observability integration design

### F060: Shared Contracts for Parallel Features

- Cross-feature boundary detection
- Interface contract generation
- Contract-first build workflow

### F061: Data Collection & AI Failure Benchmark

- Metrics collection
- AI failure taxonomy
- Benchmark report generator

### F063: Guardian Hooks and Guard Rails

Pre-edit scope enforcement for workstreams.

**Commands:**
- `sdp guard activate <ws-id>` - Enforce edit scope
- `sdp guard check <file>` - Verify file is in scope
- `sdp guard status` - Show guard status

### F064: Unified Task Resolver

Unified task ID resolution for workstreams, beads, and issues.

- ID resolution (workstream, beads, issue)
- @review artifact creation
- /issue skill backend

### F065: Agent Git Safety Protocol

- Git safety modules with structured logging
- Branch protection
- Safe git operations

### F067: Repository Hardening

**Quality Gates:**
- 80% test coverage threshold in CI
- LOC compliance (all files < 200 lines)
- Go 1.24 across all workflows

**Repository Hygiene:**
- Removed tracked auto-generated files
- Evidence log policy
- Auto-cleanup patterns

---

## CLI Commands (New in 0.8.0)

| Command | Purpose |
|---------|---------|
| `sdp doctor` | Health check |
| `sdp status` | Project state |
| `sdp init` | Initialize SDP |
| `sdp plan` | Decompose feature |
| `sdp apply` | Execute workstreams |
| `sdp guard *` | Scope enforcement |
| `sdp session *` | Session management |
| `sdp log *` | Evidence operations |
| `sdp memory *` | Long-term memory |
| `sdp drift *` | Drift detection |
| `sdp telemetry *` | Telemetry management |
| `sdp skill *` | Skill management |
| `sdp metrics *` | Metrics reporting |

---

## Changed

- **Python → Go:** Full CLI rewrite
- **Test coverage:** 68% → 80%+
- **File organization:** All files < 200 LOC
- **Documentation:** Complete rewrite

---

## Fixed

- Data race in circuit breaker tests
- Security: Checkpoint permissions 0644 → 0600
- Go version consistency (1.24)
- Context recovery in Repair()

---

## [0.7.0] - 2026-01-31

### F034: A+ Quality Initiative

- Split large files
- Test coverage to 85%+
- Domain layer extraction (Clean Architecture)
- `sdp status` command
- Skill discovery via `@help`
- Legacy code removal (~600 LOC)

---

## [0.5.2] - 2026-01-31

### F025: pip-audit Security Scanning

---

## [0.5.1] - 2026-01-31

### F020: Fast Feedback (Git Hooks)

---

## [0.4.0] - 2026-01-27

### F003-F011: Core Features

- F003: Two-stage review
- F004: Platform adapters
- F005: Extension system
- F007: Oneshot + hooks
- F008: Contract-driven tiers
- F010: SDP infrastructure
- F011: PRD command

---

## Earlier Versions

See git history for versions prior to 0.4.0.
