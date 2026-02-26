# SDP: Structured Development Protocol

**Protocol + evidence layer for AI agent workflows.**

SDP gives your AI agents a structured process (Discovery → Delivery → Evidence) and produces proof of what they actually did. It works with Claude Code, Cursor, OpenCode, or anything that can read a markdown file.

> **Read the [Manifesto](docs/MANIFESTO.md)** for the full story — what exists, what's coming, and why evidence matters.

## Demo

![SDP Demo: OpenCode Beads Viz Plugin](demo/demo.gif)

*Install → @feature → @oneshot → @review → @deploy → bd viz. Regenerate: `vhs demo/demo.tape`*

## Quick Start

**Install modes:**

| Mode | Command | What you get |
|------|---------|--------------|
| **Full project** (default) | `curl -sSL .../install.sh \| sh` | Prompts, hooks, config + optional CLI |
| **Binary only** | `curl -sSL .../install.sh \| sh -s -- --binary-only` | CLI binary only |

```bash
# Full: install into your project (auto-detects IDE)
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh

# Binary only: CLI to ~/.local/bin
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only

# Or: force specific IDE integration
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | SDP_IDE=claude sh
```

**Manual install:**

```bash
git submodule add https://github.com/fall-out-bug/sdp.git sdp
```

Skills load automatically from `sdp/.claude/skills/` (Claude Code) or `sdp/.cursor/skills/` (Cursor).

## What SDP Does

### 1. Structures agent work into phases

```
Intent → Plan → Execute → Verify → Review → Publish
```

Each phase has a contract. Skip a phase and the state machine blocks the next one.

### 2. Produces evidence of what agents did

Every run creates a strict evidence envelope — a JSON document proving intent, plan, execution, verification, review, boundary compliance, and provenance (SHA-256 hash chain). [Details in the Manifesto](docs/MANIFESTO.md#the-evidence-envelope).

**Schemas:** `schema/evidence-envelope.schema.json` (legacy), `schema/coding-workflow-predicate.schema.json` (in-toto v1). See [docs/attestation/coding-workflow-v1.md](docs/attestation/coding-workflow-v1.md).

### 3. Gates PRs on evidence

```bash
sdp-evidence validate --evidence .sdp/evidence/run-123.json
```

One command in CI. Incomplete evidence = blocked merge.

## Core Workflow

### Discovery (planning)

```
@vision "AI task manager"      → VISION.md, PRD.md, ROADMAP.md
@reality --quick               → Codebase analysis
@feature "Add authentication"  → Workstreams for feature
```

### Delivery (implementation)

```
@oneshot F001                  → Autonomous execution of all workstreams
@review F001                   → Multi-agent quality review
@deploy F001                   → Merge to main
```

### Manual mode

```
@build 00-001-01               → Single workstream with TDD
@build 00-001-02
@review F001
@deploy F001
```

### Debug

```
@debug "Test fails"            → Systematic debugging
@hotfix "API down"             → Emergency fix (P0)
@bugfix "Wrong totals"         → Quality fix (P1/P2)
```

## Skills

| Skill | Purpose | Phase |
|-------|---------|-------|
| `@vision` | Strategic planning (7 agents) | Discovery |
| `@reality` | Codebase analysis (8 agents) | Discovery |
| `@feature` | Feature planning (@idea + @design) | Discovery |
| `@idea` | Requirements gathering | Discovery |
| `@design` | Workstream decomposition | Discovery |
| `@oneshot` | Autonomous execution | Delivery |
| `@build` | TDD for single workstream | Delivery |
| `@review` | Quality review (6 agents) | Delivery |
| `@deploy` | Deploy to main | Delivery |
| `@debug` | Systematic debugging | Debug |
| `@hotfix` | Emergency fix | Debug |
| `@bugfix` | Quality fix | Debug |

## Protocol Flow

```
@oneshot F001  →  @review F001  →  @deploy F001
     │                 │                │
     ▼                 ▼                ▼
Execute WS       APPROVED?         Merge PR
                    │
                    ├─ YES → proceed
                    └─ NO → fix loop
```

**Done = @review APPROVED + @deploy completed**, not just "PR merged".

## Quality Gates

| Gate | Requirement |
|------|-------------|
| TDD | Tests first |
| Coverage | >= 80% |
| File size | < 200 LOC |
| Architecture | No layer violations |

## Workstream ID

Format: `PP-FFF-SS`

- `PP` — Project (00 = SDP itself)
- `FFF` — Feature number
- `SS` — Step number

Example: `00-024-03` = SDP, feature 24, step 3

## Project Structure

```
your-project/
├── sdp/                      # SDP submodule
│   ├── prompts/skills/       # Skills (source of truth)
│   ├── prompts/agents/       # Agent definitions
│   ├── schema/               # JSON schemas (evidence, protocol)
│   ├── hooks/                # Git hooks and validators
│   ├── .claude/              # Claude Code integration
│   ├── .cursor/              # Cursor integration
│   ├── .opencode/            # OpenCode integration
│   └── docs/                 # Documentation
└── docs/workstreams/         # Your workstreams
```

---

## Optional Components

### Go CLI

CLI provides helper commands. **Not required for the protocol to work.**

```bash
cd sdp/sdp-plugin && CGO_ENABLED=0 go build -o sdp ./cmd/sdp

sdp doctor              # Health check
sdp status              # Project state
sdp guard activate WS   # Edit scope enforcement
sdp log show            # Evidence log
```

### Beads

Task tracking for multi-session work. **Not required.**

```bash
brew tap beads-dev/tap && brew install beads
bd ready                # Find available tasks
bd create --title="..." # Create task
bd close <id>           # Close task
```

### Evidence Layer

Audit log in `.sdp/log/events.jsonl` with hash-chain provenance.

```bash
sdp log show            # Show events
sdp log trace           # Trace by commit/workstream
```

---

## Platform

**Evidence file lock requires UNIX (macOS/Linux).** The evidence layer uses `flock` for inter-process safety. Windows is not supported. SDP targets macOS/Linux CI.

---

## Where SDP Fits in the Ecosystem

SDP composes with the tools you already use. It doesn't replace them — it adds evidence.

| You Need | Use | SDP Adds |
|----------|-----|----------|
| Orchestration | Vibe Kanban, Swarm Tools | Evidence envelope for each task |
| Policy | Cupcake | Evidence that policy was checked |
| K8s agents | kubeopencode | Evidence projection from CRD status |
| Issue tracking | Beads | Evidence-gated state transitions |
| CI/CD | GitHub Actions, any CI | `sdp-evidence validate` as a PR gate |

## Research Lab

We're exploring multi-persona adversarial review, self-improvement loops, cross-project federation, and telemetry-driven backlog generation in [`sdp_lab`](https://github.com/fall-out-bug/sdp_lab). It's private for now — open an issue if you'd like to play with us.

## Documentation

| File | Content |
|------|---------|
| [demo/README.md](demo/README.md) | Terminal demo (vhs), examples |
| [docs/QUICKSTART.md](docs/QUICKSTART.md) | 5-minute getting started |
| [docs/MANIFESTO.md](docs/MANIFESTO.md) | Why SDP exists, what's real, what's next |
| [docs/PROTOCOL.md](docs/PROTOCOL.md) | Full specification |
| [docs/attestation/coding-workflow-v1.md](docs/attestation/coding-workflow-v1.md) | in-toto predicate spec (evidence v2) |
| [docs/reference/](docs/reference/) | Principles, glossary, specs |
| [CLAUDE.md](CLAUDE.md) | Quick reference for Claude Code |
| [CHANGELOG.md](CHANGELOG.md) | Version history |

## License

MIT

---

*"AI agents can implement features, but without evidence it's just vibes."*
