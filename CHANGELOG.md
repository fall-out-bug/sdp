# Changelog

All notable changes to the Spec-Driven Protocol (SDP).

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
- `sdp guard finding list` - List findings
- `sdp guard finding resolve <id>` - Resolve finding

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
