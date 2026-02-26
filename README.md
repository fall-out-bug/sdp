# SDP: Structured Development Protocol

**Protocol + evidence layer for AI agent workflows.**

SDP gives your AI agents a structured process (Discovery → Delivery → Evidence) and produces proof of what they actually did. Works with Claude Code, Cursor, OpenCode, or anything that can read markdown.

> [Manifesto](docs/MANIFESTO.md) — what exists, what's coming, why evidence matters.

## Quick Start

```bash
# Install (prompts, hooks, optional CLI)
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh

# Or binary only
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only

# Or submodule
git submodule add https://github.com/fall-out-bug/sdp.git sdp
```

Skills load from `sdp/.claude/skills/` (Claude) or `sdp/.cursor/skills/` (Cursor).

**First run:**

```bash
sdp init --auto
@feature "Your feature"
@oneshot F001
@review F001
@deploy F001
```

→ [5-minute guide](docs/QUICKSTART.md)

## What SDP Does

1. **Structures agent work** — Intent → Plan → Execute → Verify → Review → Publish. Each phase has a contract.

2. **Produces evidence** — JSON envelope with intent, plan, execution, verification, provenance (hash chain). [Details](docs/MANIFESTO.md#the-evidence-envelope).

3. **Gates PRs** — `sdp-evidence validate` in CI. Incomplete evidence = blocked merge.

## Core Workflow

| Phase | Command |
|-------|---------|
| Planning | `@vision "AI task manager"` → `@feature "Add auth"` |
| Execution | `@oneshot F001` or `@build 00-001-01` |
| Review | `@review F001` |
| Deploy | `@deploy F001` |
| Debug | `@debug`, `@hotfix`, `@bugfix` |

**Done = @review APPROVED + @deploy completed.**

## Skills

| Skill | Purpose |
|-------|---------|
| `@vision` | Strategic planning |
| `@feature` | Feature planning (→ workstreams) |
| `@oneshot` | Autonomous execution |
| `@build` | Single workstream (TDD) |
| `@review` | Multi-agent quality review |
| `@deploy` | Merge to main |
| `@debug` / `@hotfix` / `@bugfix` | Debug flows |

## Optional

**CLI:** `sdp doctor`, `sdp status`, `sdp guard activate`, `sdp log show`

**Beads:** `brew tap beads-dev/tap && brew install beads` — task tracking for multi-session work.

**Platform:** Evidence layer uses `flock` — macOS/Linux only. Windows not supported.

**Research Lab:** We're exploring multi-persona adversarial review, self-improvement loops, cross-project federation, and telemetry-driven backlog generation in [sdp_lab](https://github.com/fall-out-bug/sdp_lab). Private for now — open an issue if you'd like to play with us.

## Docs

| File | Content |
|------|---------|
| [QUICKSTART.md](docs/QUICKSTART.md) | 5-minute getting started |
| [MANIFESTO.md](docs/MANIFESTO.md) | Vision, evidence, what exists |
| [ROADMAP.md](docs/ROADMAP.md) | Where SDP is going |
| [PROTOCOL.md](docs/PROTOCOL.md) | Full specification |
| [reference/](docs/reference/) | Principles, glossary, specs |

## License

MIT
