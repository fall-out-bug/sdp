# SDP: Spec-Driven Protocol

**Turn AI coding tools into a structured development process with workstreams, quality gates, and multi-agent orchestration.**

[![Go](https://img.shields.io/badge/go-1.24+-blue.svg)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.9.0-brightgreen.svg)](CHANGELOG.md)

## What is SDP

SDP is a protocol for AI-assisted development. You describe what to build, SDP breaks it into atomic workstreams, and your AI tool (Claude Code, Cursor, or OpenCode) executes them with TDD, quality gates, and dependency tracking.

**Works with:** Claude Code, Cursor, OpenCode
**Written in:** Go (protocol is language-agnostic)
**Tracks tasks with:** [Beads CLI](https://github.com/beads-dev/beads)

---

## Installation

### Option 1: Submodule (recommended)

Add SDP to your project. The AI tools read skills, agents, and hooks directly.

```bash
# Add SDP as a submodule
git submodule add https://github.com/fall-out-bug/sdp.git sdp

# Install git hooks
cd sdp && bash hooks/install-hooks.sh && cd ..
```

Canonical prompts live in `sdp/prompts/skills/` and `sdp/prompts/agents/`.
Tool adapters (`sdp/.claude/`, `sdp/.cursor/`, `sdp/.opencode/`) point to the same source via symlinks.

### Option 2: Build the CLI (optional)

The `sdp` CLI provides health checks, guard enforcement, workstream parsing, and telemetry. It's optional — all core workflows work through AI tool skills alone.

```bash
cd sdp/sdp-plugin
go build -o sdp ./cmd/sdp
mv sdp /usr/local/bin/  # or anywhere on your PATH

# Verify
sdp doctor
sdp status
```

### Option 3: Clone standalone

For exploring or contributing without a parent project:

```bash
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
```

### Install Beads (task tracking)

Beads is the task tracker SDP uses for dependency management and issue tracking.

```bash
# macOS
brew tap beads-dev/tap && brew install beads

# Linux
curl -sSL https://raw.githubusercontent.com/beads-dev/beads/main/install.sh | bash

# Verify
bd --version
```

---

## Workflow

### Full flow (new project)

```bash
# 1. Define product vision (7 expert agents analyze your idea)
@vision "AI-powered task manager"

# 2. Scan codebase (8 expert agents analyze what's actually built)
@reality --quick

# 3. Plan a feature (interactive interview + workstream design)
@feature "Add user authentication"

# 4. Execute all workstreams autonomously
@oneshot F001
```

### Quick flow (feature on existing project)

```bash
@feature "Add password reset"    # Plan
@oneshot F050                    # Execute all workstreams
```

### Manual flow (one workstream at a time)

```bash
@build 00-050-01    # Execute single workstream with TDD
@build 00-050-02
@review F050        # Quality review
@deploy F050        # Merge to main
```

### What happens during execution

1. AI interviews you about requirements, users, success metrics
2. Creates workstreams with dependency graph (`00-001-01` -> `00-001-02` -> `00-001-03`)
3. Executes each with TDD: Red (failing test) -> Green (pass) -> Refactor
4. Validates quality gates: coverage >= 80%, type hints, linting, files < 200 LOC
5. Tracks progress in Beads CLI
6. Parallel dispatch for independent workstreams (up to 5x speedup)

---

## CLI Quick Start

The SDP CLI provides terminal commands for planning, executing, and tracking workstreams.

### Basic Workflow

```bash
# 1. Plan a feature (decompose into workstreams)
sdp plan "Add OAuth2 authentication"

# 2. See what would be created (dry-run)
sdp plan "Add OAuth2" --dry-run

# 3. Get machine-readable output
sdp plan "Add OAuth2" --output=json | jq .

# 4. Execute all ready workstreams
sdp apply

# 5. Execute specific workstream
sdp apply --ws 00-054-01

# 6. Trace evidence chain
sdp log trace
sdp log trace --ws 00-054-01

# 7. Export events for analysis
sdp log export --format=json | jq .
sdp log export --format=csv > events.csv

# 8. Show event statistics
sdp log stats
```

### Plan Modes

- **Drive mode** (default): `sdp plan "Add feature"` - Shows plan, waits for confirmation
- **Interactive mode**: `sdp plan "Add feature" --interactive` - Ask questions to refine
- **Ship mode**: `sdp plan "Add feature" --auto-apply` - Plan then execute immediately
- **Dry run**: `sdp plan "Add feature" --dry-run` - Preview without writing files

### Apply Options

- **Execute all**: `sdp apply` - Run all workstreams (no unresolved blockers)
- **Execute one**: `sdp apply --ws 00-054-01` - Run specific workstream
- **With retry**: `sdp apply --retry 3` - Retry failed workstreams up to N times
- **Dry run**: `sdp apply --dry-run` - Show execution plan without running
- **JSON output**: `sdp apply --output=json` - Machine-readable progress events

### Log Commands

```bash
# Show recent events (last 20, paginated)
sdp log show
sdp log show --page 2

# Filter by type, model, date, workstream
sdp log show --type generation
sdp log show --model claude-sonnet-4
sdp log show --since 2026-02-01T00:00:00Z
sdp log show --ws 00-054-01

# Export events
sdp log export --format=json
sdp log export --format=csv

# Trace evidence chain
sdp log trace                          # All events
sdp log trace abc123def                # By commit SHA
sdp log trace --ws 00-054-01           # By workstream
sdp log trace --ws 00-054-01 --verify  # With chain integrity check

# Statistics
sdp log stats                          # Summary by type, model, date
```

### JSON Output

All commands support JSON output for scripting and CI/CD integration:

```bash
# Plan output
sdp plan "Add feature" --output=json | jq '.workstreams[] | .id'

# Apply progress
sdp apply --output=json | jq '.workstreams[] | select(.status == "failed")'

# Log export
sdp log export --format=json | jq '.events[] | select(.type == "generation")'

# Trace output
sdp log trace --json | jq '.events | length'
```

### Environment Configuration

Set `MODEL_API` to enable automated planning:

```bash
# OpenAI API
export MODEL_API="openai:gpt-4"

# Anthropic API
export MODEL_API="anthropic:claude-sonnet-4-20250514"

# Custom endpoint
export MODEL_API="http://localhost:11434/v1:llama3"
```

Or configure in `.sdp/config.json`:

```json
{
  "version": "0.9.0",
  "model_api": "anthropic:claude-sonnet-4-20250514",
  "evidence": {
    "enabled": true,
    "log_path": ".sdp/log/events.jsonl"
  }
}
```

---

## Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `@vision` | Strategic product planning (7 agents) | `@vision "AI task manager"` |
| `@reality` | Codebase analysis (8 agents) | `@reality --quick` |
| `@feature` | Plan feature (idea + design) | `@feature "Add auth"` |
| `@build` | Execute single workstream (TDD) | `@build 00-001-01` |
| `@oneshot` | Autonomous feature execution | `@oneshot F001` |
| `@review` | Multi-agent quality review | `@review F001` |
| `@deploy` | Merge to main | `@deploy F001` |
| `@debug` | Systematic debugging | `@debug "Test fails"` |
| `@hotfix` | Emergency fix (P0) | `@hotfix "API down"` |
| `@bugfix` | Quality fix (P1/P2) | `@bugfix "Wrong totals"` |

**CLI commands** (requires Go binary):

| Command | Purpose |
|---------|---------|
| `sdp doctor` | Health check (dependencies, hooks, config) |
| `sdp status` | Show active workstream and project state |
| `sdp plan "Add feature"` | Decompose feature into workstreams |
| `sdp apply` | Execute workstreams from terminal |
| `sdp log show` | Show recent events with filters |
| `sdp log trace` | Trace evidence chain by commit/workstream |
| `sdp log export` | Export events as CSV/JSON |
| `sdp log stats` | Show event statistics |
| `sdp guard activate 00-001-01` | Enforce edit scope to a workstream |
| `sdp init` | Initialize SDP in a new project |
| `sdp parse` | Parse and validate workstream files |

---

## How it works

### Workstream hierarchy

```
Release (product milestone)
  Feature (5-30 workstreams)
    Workstream (atomic task, one-shot execution)
```

### Workstream ID: `PP-FFF-SS`

- `PP` — Project (00 = SDP itself)
- `FFF` — Feature number (001-999)
- `SS` — Step number (01-99)

Example: `00-024-03` = SDP project, feature 24 (User Auth), step 3 (Repository layer)

### Quality gates

Every workstream must pass:

| Gate | Requirement |
|------|-------------|
| TDD | Tests first, then implementation |
| Coverage | >= 80% |
| Type hints | Strict typing |
| File size | < 200 LOC per file |
| Architecture | No layer violations |
| Error handling | Explicit, no bare exceptions |

---

## Project structure

```
your-project/
├── sdp/                      # SDP submodule
│   ├── prompts/skills/       # canonical skill prompts (source of truth)
│   ├── prompts/agents/       # canonical agent prompts (source of truth)
│   ├── .claude/skills -> ../prompts/skills   # compatibility symlink
│   ├── .claude/agents -> ../prompts/agents   # compatibility symlink
│   ├── .cursor/              # Cursor IDE integration
│   ├── .opencode/            # OpenCode integration
│   ├── sdp-plugin/           # Go CLI source
│   ├── src/sdp/              # Go engine (graph, synthesis, monitoring)
│   ├── tests/                # Go test suite
│   ├── hooks/                # Git hooks
│   ├── templates/            # Workstream templates
│   ├── docs/
│   │   ├── PROTOCOL.md       # Core specification
│   │   ├── reference/        # Command and API reference
│   │   ├── vision/           # Strategic vision docs
│   │   ├── drafts/           # Feature specs (@idea output)
│   │   ├── adr/              # Architecture decisions
│   │   └── workstreams/      # Backlog and completed
│   ├── CLAUDE.md             # Claude Code guide
│   ├── PRODUCT_VISION.md     # Product vision v3.0
│   └── CHANGELOG.md          # Version history
├── docs/workstreams/         # Your project's workstreams
└── .beads/                   # Task tracking (auto-created)
```

---

## Documentation

| Document | What's in it |
|----------|-------------|
| [PROTOCOL.md](docs/PROTOCOL.md) | Full SDP specification |
| [CLAUDE.md](CLAUDE.md) | Claude Code quick reference |
| [PRODUCT_VISION.md](PRODUCT_VISION.md) | Product vision and architecture |
| [docs/reference/](docs/reference/) | Commands, quality gates, configuration |
| [docs/reference/agent-catalog.md](docs/reference/agent-catalog.md) | All 19+ agent definitions |
| [docs/vision/](docs/vision/) | Strategic vision and roadmap |
| [docs/runbooks/](docs/runbooks/) | Operational runbooks |
| [docs/compliance/COMPLIANCE.md](docs/compliance/COMPLIANCE.md) | Enterprise compliance (evidence, GDPR, SOC2) |
| [docs/compliance/THREAT-MODEL.md](docs/compliance/THREAT-MODEL.md) | Threat model and accepted risks |
| [CHANGELOG.md](CHANGELOG.md) | Version history |

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).

---

**GitHub:** [fall-out-bug/sdp](https://github.com/fall-out-bug/sdp)
