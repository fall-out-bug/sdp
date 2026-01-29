# Spec-Driven Protocol (SDP)

Workstream-driven development framework for AI agents with structured Developer-in-the-Loop workflow.

## Overview

SDP enables AI agents (Claude Code, Cursor, OpenCode) to execute software development through atomic **workstreams** — small, self-contained tasks that can be completed in one shot.

**Key concepts:**
- **Workstream (WS):** Atomic task, executable in one shot (`WS-001-01` or `bd-0001.1` with Beads)
- **Feature:** Group of related workstreams (`F001`, 5-30 WS)
- **Release:** Product milestone (10-30 features)

**Scope metrics (no time estimates!):**
- **SMALL:** < 500 LOC
- **MEDIUM:** 500-1500 LOC
- **LARGE:** > 1500 LOC → split into 2+ WS

## Quick Start

```bash
# Install
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
poetry install

# Verify
sdp --version  # Output: sdp version 0.5.0
sdp --help
```

### As Submodule (Recommended)

```bash
# Add to your project
git submodule add git@github.com:fall-out-bug/sdp.git sdp

# Update to latest
git submodule update --remote sdp
```

## Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `@feature` | **Unified feature development** (NEW) | `@feature "Add user auth"` |
| `@idea` | Interactive requirements gathering | `@idea "Add user auth"` |
| `@design` | Plan workstreams | `@design idea-user-auth` |
| `@build` | Execute single workstream | `@build WS-001-01` |
| `@debug` | **Systematic debugging** (NEW) | `@debug "Test fails"` |
| `@oneshot` | Autonomous feature execution | `@oneshot F001` |
| `@review` | Quality check | `@review F001` |
| `@deploy` | Production deployment | `@deploy F001` |
| `@issue` | Debug and route bugs | `@issue "Login fails"` |
| `@hotfix` | Emergency fix (P0) | `@hotfix "Critical bug"` |
| `@bugfix` | Quality fix (P1/P2) | `@bugfix "Incorrect totals"` |

**Internal skills** (called automatically):
- `@tdd` — TDD cycle enforcement (used by `@build`)
- `@think` — Deep structured thinking (used by @idea, @design)

## Project Structure

```
sdp/
├── PRODUCT_VISION.md  # Project manifesto (NEW)
├── src/sdp/          # Source code
│   ├── beads/        # Beads integration (NEW) - Hash-based task IDs, multi-agent coordination
│   ├── core/         # Workstream parser, feature decomposition
│   ├── schema/       # Intent validation (NEW)
│   ├── tdd/          # TDD cycle runner (NEW)
│   ├── feature/      # Product vision management (NEW)
│   ├── design/       # Dependency graph (NEW)
│   ├── github/       # GitHub integration
│   ├── prd/          # PRD command
│   └── validators/   # Quality checks
├── prompts/          # Command prompts
│   └── commands/     # /idea, /design, /build, /review, etc.
├── .claude/skills/    # AI agent skill definitions
│   ├── feature/      # Unified entry point (NEW)
│   ├── idea/         # Requirements gathering
│   ├── design/       # Workstream planning
│   ├── build/        # WS execution
│   ├── tdd/          # TDD discipline (NEW)
│   ├── debug/        # Systematic debugging (NEW)
│   └── oneshot/      # Autonomous execution
├── docs/             # Documentation
│   ├── beads-integration/  # Beads integration docs (NEW)
│   ├── schema/       # Intent JSON schema (NEW)
│   └── intent/       # Machine-readable intent (NEW)
├── tests/            # Test suite
└── hooks/            # Git hooks
```

## Quality Gates

| Gate | Requirement | Check |
|------|-------------|-------|
| AI-Readiness | Files < 200 LOC, CC < 10 | `ruff check --select=C901` |
| Test Coverage | ≥ 80% | `pytest --cov-fail-under=80` |
| Type Checking | mypy strict | `mypy --strict` |
| Clean Architecture | No layer violations | Manual review |
| Error Handling | No `except: pass` | `grep -r "except:"` |

## Features (v0.5.0)

### Completed Features
- ✅ **F003:** Two-Stage Review (5 WS)
- ✅ **F004:** Platform Adapters (4 WS) - Claude Code, Cursor, OpenCode
- ✅ **F005:** Extension System (3 WS)
- ✅ **F006:** Core SDP (6 WS)
- ✅ **F007:** Oneshot & Hooks (10 WS)
- ✅ **F008:** Contract-Driven WS Tiers (9 WS)
- ✅ **F010:** SDP Infrastructure (5 WS)
- ✅ **F011:** PRD Command (6 WS)
- ✅ **F012:** Beads Integration (16 WS) - NEW
- ✅ **F013:** AI-Human Communication Enhancement (9 WS) - NEW

**Total:** 67/67 workstreams completed (100%)

### New in v0.5.0

**Beads Integration (F012):**
- Hash-based task IDs (no conflicts in multi-agent workflows)
- Parent-child dependencies with automatic unblocking
- Multi-agent coordination with `bd ready`
- Status tracking (OPEN → IN_PROGRESS → CLOSED/BLOCKED)
- Checkpoint/resume capability
- Migration tools from markdown workstreams

**AI-Human Communication (F013):**
- `@feature` unified skill for progressive disclosure
- `PRODUCT_VISION.md` for project alignment
- `@think` skill for deep structured thinking
- Intent schema validation (machine-readable requirements)
- Enhanced workstreams with execution graphs
- EnterPlanMode for interactive planning
- TDD cycle runner for Red-Green-Refactor enforcement
- Systematic debugging with evidence-based root cause analysis

## Documentation

### Core Documentation
- [PROTOCOL.md](PROTOCOL.md) - Full specification (English)
- [PROTOCOL_RU.md](PROTOCOL_RU.md) - Полная спецификация (Русский)
- [PRODUCT_VISION.md](PRODUCT_VISION.md) - Project manifesto and direction (NEW)
- [CLAUDE.md](CLAUDE.md) - Claude Code integration guide
- [docs/workstreams/INDEX.md](docs/workstreams/INDEX.md) - Workstream index
- [CHANGELOG.md](CHANGELOG.md) - Release notes

### Beads Integration (NEW)
- [docs/beads-integration/COMPLETE.md](docs/beads-integration/COMPLETE.md) - Integration summary
- [docs/beads-integration/INSTALLATION.md](docs/beads-integration/INSTALLATION.md) - Beads CLI setup

### Multi-IDE Support
- [docs/multi-ide-parity.md](docs/multi-ide-parity.md) - Multi-IDE support

## Workflow Examples

### Traditional Markdown Workflow

```bash
@idea "Add user auth"
# → docs/drafts/idea-add-auth.md

@design idea-add-auth
# → docs/workstreams/backlog/WS-001-01.md, WS-001-02.md, ...

@build WS-001-01
# → Executes workstream, moves to completed/
```

### Beads-First Workflow (NEW - Recommended)

```bash
@idea "Add user auth"
# → bd-0001 (Beads task) + docs/intent/bd-0001.json

@design bd-0001
# → bd-0001.1, bd-0001.2, bd-0001.3 (sub-tasks with dependencies)
# + Execution graph for correct ordering

@oneshot bd-0001 --agents 3
# → Executes all workstreams in parallel with checkpoint/resume
```

**Benefits of Beads workflow:**
- No ID conflicts (hash-based vs manual PP-FFF-SS)
- Multi-agent ready (execute WS in parallel)
- Built-in dependency tracking
- `bd ready` shows what to work on next
- Checkpoint/resume for fault tolerance
- Automatic unblocking (complete WS1 → WS2 becomes ready)

### Enhanced Workflow with @feature (NEW)

```bash
@feature "Add user auth"
# → Progressive disclosure through idea → design → oneshot
# → Aligns with PRODUCT_VISION.md
# → Uses @think for complex decisions
# → Creates intent files for automation
```

## Beads Integration (NEW)

SDP now integrates with [Beads](https://github.com/steveyegge/beads) - a git-backed issue tracker for AI agents.

**Key benefits:**
- **Hash-based IDs:** `bd-a3f8` format prevents conflicts in multi-agent workflows
- **Parent-child dependencies:** Hierarchical task structure with automatic unblocking
- **Ready detection:** `bd ready` command shows tasks with no open blockers
- **Status tracking:** Automatic state transitions (OPEN → IN_PROGRESS → CLOSED/BLOCKED)
- **Multi-agent coordination:** Multiple agents can execute tasks in parallel safely

**Setup:**
```bash
# Install Beads (optional - can use mock mode for dev)
brew install go
go install github.com/steveyegge/beads/cmd/bd@latest

# Use mock mode for development
export BEADS_USE_MOCK=true

# Use real Beads for production
export BEADS_USE_MOCK=false
```

**Migration:**
```bash
# Migrate existing markdown workstreams to Beads
sdp beads migrate docs/workstreams/backlog/

# View migration status
sdp beads status
```

## Inspiration

SDP is inspired by and incorporates ideas from:

- [superpowers](https://github.com/obra/superpowers) - AI agent capabilities framework
- [vibe-kanban](https://github.com/BloopAI/vibe-kanban) - Task management for AI agents
- [sub-agents.directory](https://github.com/ayush-that/sub-agents.directory) - Multi-agent orchestration
- [Beads](https://github.com/steveyegge/beads) - Git-backed issue tracker for AI agents (NEW)

## License

MIT License - see LICENSE file for details

## Support

- GitHub Issues: https://github.com/fall-out-bug/sdp/issues
- Documentation: https://github.com/fall-out-bug/sdp/tree/main/docs

---

**Version:** 0.5.0 | **Status:** ✅ Active | **Last Updated:** 2026-01-28
