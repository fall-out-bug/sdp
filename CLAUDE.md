# Claude Code Integration Guide

Quick reference for using SDP CLI v0.9.4 with Claude Code.

## Quick Start

```bash
@vision "AI-powered task manager"    # Strategic planning
@reality --quick                     # Codebase analysis
@feature "Add user authentication"   # Plan feature
@build 00-001-01                     # Execute workstream
@review F01                          # Quality check
```

**Workstream ID Format:** `PP-FFF-SS` (e.g., `00-001-01`)

---

## Protocol Flow

The correct workflow is:

```
@oneshot F001  →  @review F001  →  @deploy F001
    │                 │                │
    ▼                 ▼                ▼
 Execute WS      APPROVED?         Merge PR
                  │
                  ├─ YES → proceed to @deploy
                  └─ NO → fix loop (max 3 iterations)
```

**"Done" = @review APPROVED + @deploy completed, NOT just "PR merged".**

---

## Decision Tree

```
New project?
|-- Yes --> @vision (strategic) --> @reality (analysis)
+-- No --> Working on existing project?
    |-- Yes --> What's the state?
    |   |-- Don't know --> @reality --quick
    |   +-- Know state --> @feature "add feature"
    +-- No --> Workstreams exist?
        |-- Yes --> @oneshot F050
        +-- No --> @feature "plan feature"
```

### Four-Level Planning Model

| Level | Orchestrator | Purpose | Output |
|-------|-------------|---------|--------|
| **Strategic** | @vision (7 agents) | Product planning | VISION, PRD, ROADMAP |
| **Analysis** | @reality (8 agents) | Codebase analysis | Reality report |
| **Feature** | @feature (roadmap pre-check + @idea + @ux + @design) | Requirements + WS | Workstreams |
| **Execution** | @oneshot (@build) | Parallel execution | Implemented code |

### When to Use Each Level

**@vision** — New project, major pivot, quarterly strategic review

**@reality** — New to project, before @feature, track tech debt, quarterly review

**@feature** — Feature idea but no workstreams, need interactive planning (full discovery flow)

**@ux** — UX research for user-facing features (standalone or auto-triggered by @feature)

**@oneshot** — Workstreams exist, want autonomous execution with checkpoint/resume

**@build** — Execute a single workstream (use instead of @oneshot for 1-2 WS)

---

## Available Skills

### Core Skills

| Skill | Purpose | Phase |
|-------|---------|-------|
| `@vision` | Strategic product planning (7 expert agents) | Strategic |
| `@reality` | Codebase analysis (8 expert agents) | Analysis |
| `@feature` | Planning orchestrator (roadmap pre-check + idea + ux + design) | Planning |
| `@idea` | Requirements gathering (AskUserQuestion) | Planning |
| `@ux` | UX research (mental model elicitation) | Planning |
| `@design` | Workstream design (EnterPlanMode) | Planning |
| `@oneshot` | Execution orchestrator (autonomous) | Execution |
| `@build` | Execute single workstream (TDD) | Execution |
| `@review` | Multi-agent quality review | Execution |
| `@deploy` | Merge feature branch to main | Execution |

### Debug Skills

| Skill | Purpose | Phase |
|-------|---------|-------|
| `@debug` | Systematic debugging (scientific method) | Debug |
| `@issue` | Debug and route bugs | Debug |
| `@hotfix` | Emergency fix (P0) | Debug |
| `@bugfix` | Quality fix (P1/P2) | Debug |

### Utility Skills

| Skill | Purpose |
|-------|---------|
| `@init` | Initialize SDP in current project |
| `@help` | Interactive skill discovery |
| `@prototype` | Rapid prototyping shortcut |
| `@vision --update` | PRD/ diagram regeneration |
| `@test` | Contract test generation |
| `@reality-check` | Quick documentation vs code validation |
| `@verify-workstream` | Validate workstream against codebase |
| `@protocol-consistency` | Audit consistency across docs/CLI/CI |
| `@guard` | Pre-edit gate enforcing WS scope |
| `@tdd` | TDD enforcement (called by @build) |

### Beads Integration

| Skill | Purpose |
|-------|---------|
| `@beads` | Beads task tracker integration |

**Skills defined in:** `.claude/skills/{name}/SKILL.md`

---

## Typical Workflow

### Full Flow (new project)

```bash
# 1. Strategic planning
@vision "AI-powered task manager for remote teams"

# 2. Codebase analysis
@reality --quick

# 3. Feature planning (per feature)
@feature "User can reset password via email"

# 4. Autonomous execution
@oneshot F050
```

### Quick Flow (existing project)

```bash
# 1. Plan feature
@feature "Add payment processing"

# 2. Execute all workstreams
@oneshot F050
```

### Manual Flow (learning or debugging)

```bash
@build 00-050-01   # Execute one at a time
@build 00-050-02
@review F050       # Review when done
@deploy F050       # Deploy
```

---

## First Time Setup

1. **Read core docs:**
   - [README.md](README.md)
   - [PROTOCOL.md](docs/PROTOCOL.md)

2. **Key concepts:**
   - **Workstream (WS)**: Atomic task, one-shot execution
   - **Feature**: 5-30 workstreams
   - **Release**: 10-30 features

3. **Install Beads CLI** (task tracking):
   ```bash
   brew tap beads-dev/tap && brew install beads
   bd --version
   ```

---

## Project Structure

```
sdp/
├── sdp-plugin/            # Go implementation (CLI + agents)
│   ├── cmd/sdp/           # CLI commands
│   └── internal/          # Core logic
├── .claude/
│   ├── skills/            # Skill definitions
│   └── agents/            # Multi-agent definitions
├── docs/
│   ├── PROTOCOL.md        # Core specification
│   ├── reference/         # Command and API reference
│   ├── vision/            # Strategic vision docs
│   ├── decisions/         # Architecture decisions
│   ├── drafts/            # @idea output
│   └── workstreams/       # Backlog + completed WS
├── hooks/                 # Git hooks and validators
├── templates/             # Workstream templates
├── PRODUCT_VISION.md      # Product vision
└── go.mod                 # Go module
```

---

## Quality Gates

| Gate | Requirement |
|------|-------------|
| **File Size** | < 200 LOC |
| **Test Coverage** | >= 80% |
| **Type Hints** | Full strict typing |
| **Clean Architecture** | No layer violations |
| **Error Handling** | Explicit, no bare exceptions |
| **TODOs** | All resolved or tracked in WS |

### Forbidden Patterns
- Files > 200 LOC
- Time-based estimates
- Layer violations
- Coverage < 80%
- TODO without followup WS

### Required Patterns
- Type hints everywhere
- Tests first (TDD)
- Explicit error handling
- Clean architecture boundaries
- Conventional commits

---

## Key Principles

- **SOLID, DRY, KISS, YAGNI** — see [docs/reference/PRINCIPLES.md](docs/reference/PRINCIPLES.md)
- **Clean Architecture** — Domain <- App <- Infra <- Presentation
- **TDD** — Tests first (Red -> Green -> Refactor)
- **AI-Readiness** — Small files, low complexity, typed

---

## CLI Reference

The SDP CLI provides terminal commands for planning, executing, and tracking workstreams.

### Core Commands

| Command | Purpose |
|---------|---------|
| `sdp doctor` | Health check (hooks, config, deps) |
| `sdp status` | Show project state |
| `sdp init` | Initialize SDP in a new project |
| `sdp parse <ws-file>` | Parse workstream file |

### Guard Commands

| Command | Purpose |
|---------|---------|
| `sdp guard activate <ws-id>` | Enforce edit scope for workstream |
| `sdp guard check <file>` | Verify file is in scope |
| `sdp guard status` | Show guard status |
| `sdp guard deactivate` | Clear edit scope |
| `sdp guard finding list` | List guard findings |
| `sdp guard finding resolve <id>` | Resolve a finding |

### Session Commands

| Command | Purpose |
|---------|---------|
| `sdp session show` | Show current session |
| `sdp session clear` | Clear session |

### Log Commands

| Command | Purpose |
|---------|---------|
| `sdp log show` | Show recent events with filters |
| `sdp log trace` | Trace evidence chain |
| `sdp log export` | Export events as CSV/JSON |
| `sdp log stats` | Show event statistics |

### Memory Commands

| Command | Purpose |
|---------|---------|
| `sdp memory index` | Index project artifacts |
| `sdp memory search <query>` | Search indexed artifacts |
| `sdp memory stats` | Show index statistics |

### Drift Commands

| Command | Purpose |
|---------|---------|
| `sdp drift detect [ws-id]` | Detect code↔docs drift |

### Metrics Commands

| Command | Purpose |
|---------|---------|
| `sdp metrics report` | Show metrics report |
| `sdp metrics classify` | Classify metrics |

### Telemetry Commands

| Command | Purpose |
|---------|---------|
| `sdp telemetry status` | Show telemetry status |
| `sdp telemetry analyze` | Analyze telemetry data |

### Skill Commands

| Command | Purpose |
|---------|---------|
| `sdp skill list` | List available skills |
| `sdp skill show <name>` | Show skill details |
| `sdp skill validate` | Validate skill definitions |

### Plan Options

- **Dry run**: `--dry-run` - Preview without writing files
- **JSON output**: `--output=json` - Machine-readable format

### Log Filters

- **By type**: `--type generation` - Filter event type
- **By workstream**: `--ws 00-054-01` - Filter by workstream ID
- **By date**: `--since 2026-02-01T00:00:00Z` - Filter by ISO date

---

## Evidence Layer

Build and verify flows emit events to `.sdp/log/events.jsonl` (hash-chained).

Config: `.sdp/config.yml` with `version`, `evidence.enabled`, `evidence.log_path`. @build emits plan/generation/verification events when evidence is enabled.

---

## Long-term Memory (F051)

Project memory system for avoiding duplicated work. Integrates with evidence.jsonl and Beads issues.

### Architecture

```
.sdp/
├── memory.db        # SQLite + FTS5 index
├── log/
│   └── events.jsonl # Evidence log (hash-chained)
└── notifications.log # Notification channel log
```

### Commands

| Command | Purpose |
|---------|---------|
| `sdp memory index` | Index all docs/ artifacts into memory.db |
| `sdp memory search <query>` | Full-text search across indexed artifacts |
| `sdp memory stats` | Show index statistics |
| `sdp drift detect [ws_id]` | Detect code↔docs drift |

### Use Cases

1. **Context Recovery:** After session compaction, search memory to restore context
2. **Decision Discovery:** Find related decisions before proposing approaches
3. **Drift Detection:** Detect code-documentation discrepancies

---

## Troubleshooting

### Skill not found
Check `.claude/skills/{name}/SKILL.md` exists

### Workstream blocked
Check dependencies in `docs/workstreams/backlog/{WS-ID}.md`

### Coverage too low
Run test coverage tool with verbose output to identify gaps

---

## Landing the Plane (Session Completion)

**When ending a work session**, complete ALL steps. Work is NOT complete until `git push` succeeds.

1. **File issues** for remaining work
2. **Run quality gates** (if code changed)
3. **Update issue status** — Close finished, update in-progress
4. **PUSH TO REMOTE:**
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** — Clear stashes, prune branches
6. **Verify** — All committed AND pushed
7. **Hand off** — Context for next session

---

## Resources

| Resource | Purpose |
|----------|---------|
| [PROTOCOL.md](docs/PROTOCOL.md) | Full specification |
| [docs/reference/PRINCIPLES.md](docs/reference/PRINCIPLES.md) | Core principles |
| [docs/SLOS.md](docs/SLOS.md) | SLOs/SLIs |
| [docs/reference/CODE_PATTERNS.md](docs/reference/CODE_PATTERNS.md) | Code patterns |
| [docs/reference/MODELS.md](docs/reference/MODELS.md) | Model recommendations |
| [.claude/skills/](.claude/skills/) | Skill definitions |
| [docs/compliance/COMPLIANCE.md](docs/compliance/COMPLIANCE.md) | Enterprise compliance |
| [docs/compliance/THREAT-MODEL.md](docs/compliance/THREAT-MODEL.md) | Threat model |

---

**CLI Version:** 0.9.4 | **Protocol Version:** 0.10.0
