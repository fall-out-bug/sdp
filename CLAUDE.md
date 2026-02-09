# Claude Code Integration Guide

Quick reference for using SDP v0.9.0 with Claude Code.

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
| **Feature** | @feature (@idea + @design) | Requirements + WS | Workstreams |
| **Execution** | @oneshot (@build) | Parallel execution | Implemented code |

### When to Use Each Level

**@vision** — New project, major pivot, quarterly strategic review

**@reality** — New to project, before @feature, track tech debt, quarterly review

**@feature** — Feature idea but no workstreams, need interactive planning

**@oneshot** — Workstreams exist, want autonomous execution with checkpoint/resume

**@build** — Execute a single workstream (use instead of @oneshot for 1-2 WS)

---

## Available Skills

| Skill | Purpose | Phase |
|-------|---------|-------|
| `@vision` | Strategic product planning (7 expert agents) | Strategic |
| `@reality` | Codebase analysis (8 expert agents) | Analysis |
| `@feature` | Planning orchestrator (interactive) | Planning |
| `@idea` | Requirements gathering (AskUserQuestion) | Planning |
| `@design` | Workstream design (EnterPlanMode) | Planning |
| `@oneshot` | Execution orchestrator (autonomous) | Execution |
| `@build` | Execute single workstream (TDD) | Execution |
| `@review` | Multi-agent quality review | Execution |
| `@deploy` | Merge feature branch to main | Execution |
| `@debug` | Systematic debugging (scientific method) | Debug |
| `@issue` | Debug and route bugs | Debug |
| `@hotfix` | Emergency fix (P0) | Debug |
| `@bugfix` | Quality fix (P1/P2) | Debug |

**Internal:** `/tdd` (TDD enforcement, called by `@build`)

Skills defined in `.claude/skills/{name}/SKILL.md`

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
│   ├── cmd/               # CLI commands
│   └── internal/          # Core logic
├── src/sdp/               # Go source modules
│   ├── agents/            # Code analysis, contracts
│   ├── graph/             # Dependency graph, dispatcher
│   ├── monitoring/        # Metrics, SLO tracking
│   ├── synthesis/         # Agent synthesis engine
│   ├── reality/           # Codebase scanners
│   └── vision/            # Vision extractor
├── tests/                 # Go test suite
├── .claude/
│   ├── skills/            # Skill definitions
│   └── agents/            # Multi-agent definitions
├── docs/
│   ├── PROTOCOL.md        # Core specification
│   ├── reference/         # Command and API reference
│   ├── vision/            # Strategic vision docs
│   ├── drafts/            # @idea output
│   └── workstreams/       # Backlog + completed WS
├── hooks/                 # Git hooks and validators
├── templates/             # Workstream templates
├── PRODUCT_VISION.md      # Product vision v3.0
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

## Validation

```bash
hooks/pre-build.sh 00-001-01     # Pre-build check
hooks/post-build.sh 00-001-01    # Post-build check
```

---

## Evidence Layer

Build and verify flows emit events to `.sdp/log/events.jsonl` (hash-chained). Use:

| Command | Purpose |
|---------|---------|
| `sdp acceptance run` | Run acceptance gate from `.sdp/config.yml` (e.g. `go test -run TestSmoke`) |
| `sdp log trace [commit]` | Show evidence timeline; `--ws`, `--json`, `--verify` for chain check |
| `sdp collision check` | Detect scope overlaps between active workstreams |

Config: `.sdp/config.yml` with `version`, `acceptance.command`, `evidence.enabled`, `evidence.log_path`. @build emits plan/generation/verification events when evidence is enabled.

---

## Parallel Execution

The @oneshot orchestrator uses a parallel dispatcher (`src/sdp/graph/`) with:

1. **Build Graph** — Parse WS files, extract dependencies, build DAG
2. **Topological Sort** — Kahn's algorithm for valid execution order
3. **Parallel Dispatch** — Execute independent WS concurrently (3-5 agents)
4. **Circuit Breaker** — Fault tolerance with retry logic
5. **Checkpoint** — Atomic save/restore for resume after interruption

Speedup: ~5x for 5-10 workstreams.

---

## Multi-Agent Synthesis

When agents disagree, the Synthesizer resolves conflicts:

1. **Unanimous** — All agents agree
2. **Domain Expertise** — Highest confidence wins
3. **Quality Gate** — Best quality score
4. **Merge** — Combine best parts
5. **Escalate** — Ask human

See [docs/reference/agent-catalog.md](docs/reference/agent-catalog.md) for agent documentation.

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

## Reality-First Development

Always verify actual code before following documentation.

**Before modifying any file:**
```bash
@reality-check <filename>
```

**Before executing workstreams:**
```bash
@verify-workstream 00-001-01
```

Validates: scope files exist, functions/classes match docs, architectural layers correct.

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

---

**Version:** 0.9.0
