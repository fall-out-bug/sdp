# SDP: Spec-Driven Protocol

**Workstream-driven framework that turns AI coding tools (Claude Code, Cursor, OpenCode) into a structured software development process.**

[![Go](https://img.shields.io/badge/go-1.25+-blue.svg)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Coverage](https://img.shields.io/badge/coverage-83.2%25-brightgreen.svg)](tests/)
[![Version](https://img.shields.io/badge/version-0.9.0-brightgreen.svg)](CHANGELOG.md)

---

## v0.9.0 Architecture Update

**SDP v0.9.0 is a major architectural transformation** from Python CLI to Go-based multi-agent orchestration system.

### What Changed

- **Current Status:** Active development (26/26 workstreams complete)
- **Architecture:** Multi-agent orchestration with 19 specialized agents
- **Installation:** Go binary (language-agnostic protocol)
- **Quality:** 83.2% test coverage, 4.96x parallel execution speedup

### Migration from Python SDP

**If you're upgrading from Python SDP (v0.7.x or earlier):**

| Feature | Python SDP (v0.7.x) | Go SDP (v0.9.0) |
|---------|---------------------|-----------------|
| **Installation** | `pip install sdp` | `go install github.com/fall-out-bug/sdp@latest` |
| **Runtime** | Python 3.10+ | Go 1.25.6+ |
| **Execution** | Manual @build per WS | Autonomous @oneshot with parallel dispatcher |
| **Agents** | 0 (skill-based) | 19 (strategic planning + codebase analysis) |
| **Planning** | @feature only | @vision → @reality → @feature → @oneshot |

**Your existing workstreams are compatible** - the protocol is language-agnostic!

See [docs/migrations/multi-agent-migration.md](docs/migrations/multi-agent-migration.md) for complete migration guide.

---

## When to Use SDP

SDP is for you if:

- ✅ **You already have an AI-IDE** (Claude Code, Cursor, OpenCode), but lack a structured process
- ✅ **You need to manage complex features** through atomic workstreams with dependencies
- ✅ **You want repeatable quality gates** for AI-generated code (TDD, coverage, type hints)
- ✅ **You prefer progress tracking** with task systems (Beads CLI) over manual to-do lists
- ✅ **You're building solo or small team** projects with 5-500 workstreams

**New to SDP?** See [START_HERE.md](https://github.com/fall-out-bug/sdp/blob/main/docs/beginner/START_HERE.md) for beginner guides!

**Upgrading from an earlier version?** See [Migration Guide](docs/migration/breaking-changes.md) for breaking changes!

---

## Quick Start

### Install via Go Binary (Recommended)

```bash
# Install Go 1.25.6+
# macOS
brew install go
# OR download from golang.org

# Install SDP
go install github.com/fall-out-bug/sdp@latest

# Verify installation
sdp status
```

### Use as Submodule (For project-specific customization)

```bash
# Add to your project
git submodule add git@github.com:fall-out-bug/sdp.git sdp

# Initialize (sets up Git hooks, skills, agents)
cd sdp
./scripts/init.sh
```

---

## Minimal Example

Complete workflow in 30 seconds:

```bash
# 1. Create feature (interactive interviews)
@feature "Add user comments"

# 2. Plan workstreams
@design beads-comments

# 3. Execute one workstream
@build 00-COMMENTS-01

# Or execute all autonomously
@oneshot beads-comments
```

**What happens:**
- AI interviews you about requirements, users, success metrics
- Creates workstreams with dependencies (00-001-01 → 00-001-02 → 00-001-03)
- Executes with TDD (Red → Green → Refactor)
- Validates quality gates (≥80% coverage, mypy --strict, ruff)
- Tracks progress in Beads CLI
- Sends Telegram notifications on completion

---

## Features

### Multi-IDE Support
Works with Claude Code, Cursor, OpenCode through unified skill system. Switch between AI tools without changing your workflow.

### Autonomous Execution
`@oneshot` mode executes entire features autonomously with:
- Checkpoint save/restore (resume after interruption)
- Background execution support
- Progress notifications via Telegram
- Automatic dependency resolution

### Quality Gates Built-in
Every workstream passes:
- **TDD** - Tests first, code second
- **Coverage ≥80%** - Enforced on all files
- **Type hints** - Full mypy --strict compliance
- **Linting** - ruff for code quality
- **File size <200 LOC** - Keep code focused

### Beads Integration
Native task tracking with:
- Hash-based task IDs (bd-0001, bd-0001.1, etc.)
- Dependency DAG (00-001-02 blocked by 00-001-01)
- Ready task detection (`bd ready` shows executable tasks)
- JSONL storage for Git versioning

### Progressive Disclosure
`@feature` skill uses 5-minute interview to understand requirements before planning. No premature design, ask questions as you go.

*See [PROTOCOL.md](docs/PROTOCOL.md) for full feature list.*

---

## Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `sdp doctor` | **Health checks** | `sdp doctor` |
| `@feature` | **Unified feature development** | `@feature "Add user auth"` |
| `@idea` | Interactive requirements gathering | `@idea "Add user auth"` |
| `@design` | Plan workstreams | `@design idea-user-auth` |
| `@build` | Execute single workstream | `@build 00-001-01` |
| `@debug` | **Systematic debugging** | `@debug "Test fails"` |
| `@oneshot` | Autonomous feature execution | `@oneshot F001` |
| `@review` | Quality check | `@review F001` |
| `@deploy` | Production deployment | `@deploy F001` |
| `@issue` | Debug and route bugs | `@issue "Login fails"` |
| `@hotfix` | Emergency fix (P0) | `@hotfix "Critical bug"` |
| `@bugfix` | Quality fix (P1/P2) | `@bugfix "Incorrect totals"` |

### Health Checks

The `sdp doctor` command performs diagnostic checks on your SDP installation:

```bash
sdp doctor                    # Human-readable output
sdp doctor --format json      # Machine-readable JSON
```

**Checks performed:**
- Python version (>= 3.10) - **Critical**
- Poetry installation - **Critical**
- Git hooks configuration - **Critical**
- Beads CLI (optional)
- GitHub CLI (optional)
- Telegram configuration (optional)

**Exit codes:**
- `0` - All critical checks passed
- `1` - One or more critical checks failed

---

## Architecture

### Workstream Hierarchy

```
Release (product milestone)
  └─ Feature (5-30 workstreams)
      └─ Workstream (atomic task, one-shot)
          ├─ SMALL: < 500 LOC
          ├─ MEDIUM: 500-1500 LOC
          └─ LARGE: > 1500 LOC (split into 2+)
```

### Example: Adding User Authentication

```
Feature F24: User Authentication
├─ 00-024-01: Domain model (450 LOC, MEDIUM)
├─ 00-024-02: Database schema (300 LOC, MEDIUM)
├─ 00-024-03: Repository layer (500 LOC, MEDIUM)
├─ 00-024-04: Service layer (600 LOC, MEDIUM)
└─ 00-024-05: API endpoints (400 LOC, MEDIUM)
```

---

## Project Structure

```
sdp/
├── PRODUCT_VISION.md     # Project manifesto
├── src/sdp/              # Source code
│   ├── beads/            # Beads integration - task tracking
│   ├── core/             # Workstream parser, decomposition
│   ├── schema/           # Intent validation
│   ├── tdd/              # TDD cycle runner
│   ├── feature/          # Product vision management
│   ├── design/           # Dependency graph
│   └── unified/          # Multi-agent coordination (NEW)
│       ├── agent/        # Agent spawning, messaging, roles
│       └── notifications/ # Telegram + Console providers
├── prompts/              # Command prompts
├── .claude/skills/       # AI agent skill definitions
├── docs/                  # Documentation
│   ├── TUTORIAL.md        # 15-minute quick start
│   ├── schema/            # JSON schemas
│   ├── intent/            # Machine-readable intents
│   └── drafts/            # Feature specifications
└── tests/                 # Test suite (309 tests, 91% coverage)
```

---

## Documentation

### Beginner Guides (Progressive Learning)
**New to SDP? Start here:**
- [docs/beginner/START_HERE.md](docs/beginner/START_HERE.md) - Welcome page with learning paths
- [docs/beginner/00-quick-start.md](docs/beginner/00-quick-start.md) - 5-minute overview
- [docs/beginner/01-first-feature.md](docs/beginner/01-first-feature.md) - Hands-on tutorial
- [docs/beginner/02-common-tasks.md](docs/beginner/02-common-tasks.md) - Common workflows
- [docs/beginner/03-troubleshooting.md](docs/beginner/03-troubleshooting.md) - Troubleshooting

### Migration Guides
**Upgrading from a previous version?**
- [docs/migrations/breaking-changes.md](docs/migrations/breaking-changes.md) - Complete guide for all breaking changes
- [docs/migration/ws-naming-migration.md](docs/migration/ws-naming-migration.md) - Workstream ID format migration (WS-FFF-SS → PP-FFF-SS)

### Reference Documentation (Lookup)
**Looking up specific commands or config?**
- [docs/reference/README.md](docs/reference/README.md) - Reference overview
- [docs/reference/commands.md](docs/reference/commands.md) - All commands
- [docs/reference/quality-gates.md](docs/reference/quality-gates.md) - Quality standards
- [docs/reference/configuration.md](docs/reference/configuration.md) - Config files
- [docs/reference/skills.md](docs/reference/skills.md) - Skill system
- [docs/SLOS.md](docs/SLOS.md) - SLOs/SLIs for CLI tool

### Internals (Maintainer Docs)
**Extending SDP or contributing?**
- [docs/internals/README.md](docs/internals/README.md) - Internals overview
- [docs/internals/architecture.md](docs/internals/architecture.md) - System architecture
- [docs/internals/extending.md](docs/internals/extending.md) - How to extend
- [docs/internals/development.md](docs/internals/development.md) - Dev setup
- [docs/internals/contributing.md](docs/internals/contributing.md) - Contributing

### Core Documentation
- [PROTOCOL.md](docs/PROTOCOL.md) - Complete specification
- [CODE_PATTERNS.md](docs/reference/CODE_PATTERNS.md) - Implementation patterns
- [CLAUDE.md](CLAUDE.md) - Claude Code integration
- [docs/SITEMAP.md](docs/SITEMAP.md) - Full documentation index

### Key Concepts

| Concept | Description | Link |
|----------|-------------|------|
| **Workstreams** | Atomic tasks, one-shot execution | [PROTOCOL.md#Workstream-Flow](PROTOCOL.md#workstream-flow) |
| **Scope Metrics** | Size by LOC, not time | [PROTOCOL.md#terminology](PROTOCOL.md#terminology) |
| **Quality Gates** | Coverage, mypy, ruff, file size | [PROTOCOL.md#quality-gates](PROTOCOL.md#quality-gates) |
| **Agent System** | Multi-agent coordination | [.claude/agents/README.md](.claude/agents/README.md) |
| **Beads Integration** | Task tracking | [PROTOCOL.md#unified-workflow-ai-comm--beads](PROTOCOL.md#unified-workflow-ai-comm--beads) |

---

## Development Status

**Current Version:** v0.9.0 (Multi-Agent Architecture)

**Status:** 🚀 **Active Development** - 26/26 workstreams complete

**Implemented:**
- ✅ Multi-agent orchestration (19 specialized agents)
- ✅ Four-level planning model (@vision → @reality → @feature → @oneshot)
- ✅ Parallel execution (4.96x speedup with dependency-aware dispatcher)
- ✅ Fault tolerance (circuit breaker + atomic checkpoint/resume)
- ✅ Progressive disclosure (12-27 questions vs unbounded)
- ✅ Two-stage review (implementer → spec reviewer → quality)
- ✅ Test coverage 83.2% (30 tests, graph package 88.3%)

**Roadmap:**
- 🔄 v1.0.0: Claude Plugin marketplace distribution
- 📋 v1.1.0: OpenTelemetry integration, enhanced synthesis rules
- 📋 v1.2.0: Real-time dashboards (Grafana), alerting (PagerDuty)

See [ROADMAP.md](docs/roadmap/ROADMAP.md) for complete planning.

---

## Contributing

Contributions are welcome! This is an active project.

**Areas of contribution:**
- Bug fixes and improvements
- New expert agents (see [docs/reference/agent-catalog.md](docs/reference/agent-catalog.md))
- Documentation and examples
- Performance optimizations

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Topics

ai-agents developer-tools workflow prompt-engineering spec-driven-development claude-code cursor opencode terminal-workflows multi-agent-coordination task-tracking quality-gates

---

**Website:** [Documentation Index](docs/workstreams/INDEX.md)
**GitHub:** [fall-out-bug/sdp](https://github.com/fall-out-bug/sdp)

*Workstream-driven development for AI agents*
