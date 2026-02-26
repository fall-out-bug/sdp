# SDP Quick Start

Get from zero to your first feature in 5 minutes.

## 1. Install

```bash
# Into your project (auto-detects Claude Code, Cursor, OpenCode)
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh

# Or: submodule
git submodule add https://github.com/fall-out-bug/sdp.git sdp
```

Skills load from `sdp/.claude/skills/`, `sdp/.cursor/skills/`, or `sdp/.opencode/`.

## 2. Initialize

```bash
sdp init --auto    # Safe defaults, non-interactive
# or
sdp init --guided  # Interactive wizard
```

Creates `.sdp/config.yml`, guard rules, and IDE integration.

*If you get "unknown flag: --auto", upgrade the CLI: `curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only`*

## 3. Create a Feature

**Discovery (planning):**

```
@feature "OpenCode plugin for beads visualization"
```

This runs:
- `@idea` — Requirements gathering
- `@design` — Workstream decomposition into `docs/workstreams/backlog/00-XXX-YY.md`

**Delivery (implementation):**

```
@oneshot F001      # Autonomous: build all workstreams
@review F001       # Multi-agent quality review
@deploy F001       # Merge to main
```

Or step-by-step:

```
@build 00-001-01   # Single workstream with TDD
@build 00-001-02
@review F001
@deploy F001
```

## 4. Verify

```bash
sdp verify 00-001-01   # Check workstream completion
sdp status             # Project state
sdp log show           # Evidence log
```

## 5. Optional: Beads

Task tracking for multi-session work:

```bash
brew tap beads-dev/tap && brew install beads
bd ready               # Find available tasks
bd create --title="..." # Create task
bd close <id>          # Close task
```

## Flow Summary

```
@feature "X"  →  @oneshot F001  →  @review F001  →  @deploy F001
     │                  │                │                │
     ▼                  ▼                ▼                ▼
  Workstreams      Execute WS       APPROVED?         Merge PR
```

**Done = @review APPROVED + @deploy completed.**

## Next

- [PROTOCOL.md](PROTOCOL.md) — Full specification
- [MANIFESTO.md](MANIFESTO.md) — Vision and evidence
- [reference/](reference/) — Commands, specs, glossary
