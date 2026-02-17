# SDP: Spec-Driven Protocol

**A structured protocol for AI-assisted development.**

SDP turns your AI assistant into a predictable process: Discovery → Delivery → Evidence.

## What it is

SDP is a set of prompts (skills) that load into Claude Code, Cursor, or other AI tools. Skills define:

- **Discovery**: How to gather requirements and plan features
- **Delivery**: How to write code with TDD and quality gates
- **Evidence**: How to record decisions in an audit log

**Everything works through skills.** CLI and Beads are optional add-ons.

## Quick Start

**One-liner installer (recommended):**

```bash
# Install for all IDEs
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash

# Or specify your IDE
SDP_IDE=claude curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash
SDP_IDE=cursor curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash
```

**Supported IDEs:** `claude`, `cursor`, `opencode`, `all` (default)

**Manual install:**

```bash
git submodule add https://github.com/fall-out-bug/sdp.git sdp
echo "sdp/.git" >> .gitignore
```

Skills load automatically from `sdp/.claude/skills/` (Claude Code) or `sdp/.cursor/skills/` (Cursor).

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
│   ├── .claude/              # Claude Code integration
│   ├── docs/                 # Documentation
│   └── CLAUDE.md             # Quick reference
└── docs/workstreams/         # Your workstreams
```

---

# Optional Components

## Go CLI (experimental)

CLI provides helper commands. **Not required for protocol to work.**

```bash
# Install
cd sdp/sdp-plugin && go build -o sdp ./cmd/sdp

# Commands
sdp doctor              # Health check
sdp status              # Project state
sdp guard activate WS   # Edit scope enforcement
sdp log show            # Evidence log
```

## Beads (experimental)

Task tracking for multi-session work. **Not required.**

```bash
brew tap beads-dev/tap && brew install beads
bd ready                # Find available tasks
bd create --title="..." # Create task
bd close <id>           # Close task
```

## Evidence Layer (experimental)

Audit log in `.sdp/log/events.jsonl` with hash-chain.

```bash
sdp log show            # Show events
sdp log trace           # Trace by commit/workstream
```

---

## Documentation

| File | Content |
|------|---------|
| [CLAUDE.md](CLAUDE.md) | Quick reference for Claude Code |
| [docs/PROTOCOL.md](docs/PROTOCOL.md) | Full specification |
| [docs/vision/ROADMAP.md](docs/vision/ROADMAP.md) | Roadmap and milestones |
| [CHANGELOG.md](CHANGELOG.md) | Version history |

## License

MIT

---

**GitHub:** [fall-out-bug/sdp](https://github.com/fall-out-bug/sdp)
