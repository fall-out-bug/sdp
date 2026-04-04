# SDP Quick Start

Get from zero to your first feature in 5 minutes.

Use this doc when your goal is to adopt SDP in your own repo, not to work on `sdp_lab`.

## 0. Choose Your Starting Point

- **Greenfield:** new repo or empty service. Install SDP, run `sdp init --auto`, then start the feature flow.
- **Brownfield:** existing codebase. Install SDP, prefer `sdp init --guided` so you can inspect defaults, then run `sdp doctor` before trusting the flow.

SDP installs prompts, hooks, and optional CLI helpers. You still configure your model provider and API keys in your IDE or agent runtime.

## 1. Install

**Full project** (prompts + hooks + optional CLI): default install

```bash
# Into your project (auto-detects Claude Code, Cursor, OpenCode, Codex)
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh
```

**Binary only** (CLI to ~/.local/bin):

```bash
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only
```

**Or: submodule**

```bash
git submodule add https://github.com/fall-out-bug/sdp.git sdp
```

Use the GitHub URL as the canonical submodule source. A local relative URL like `../sdp` is only a private convenience clone and will break reproducibility for other machines and CI.

Auto-setup today is first-class for:

- `Claude Code`
- `Cursor`
- `OpenCode` / `Windsurf`
- `Codex`

If auto-detect misses your tool, rerun with `SDP_IDE=claude|cursor|opencode|codex`.

Skills load from `sdp/.claude/skills/`, `sdp/.cursor/skills/`, `sdp/.opencode/`, or `.codex/skills/sdp/`.

## 2. Initialize

```bash
sdp init --auto    # Safe defaults, non-interactive
# or
sdp init --guided  # Interactive wizard
```

Creates `.sdp/config.yml`, guard rules, and refreshes the IDE integration already present in the project. If no IDE integration exists yet, `sdp init` falls back to `.claude/`.

*If you get "unknown flag: --auto", upgrade the CLI: `curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only`*

Then verify the environment:

```bash
sdp doctor
```

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
@oneshot <feature-id>      # Autonomous: build all workstreams
@review <feature-id>       # Multi-agent quality review
@deploy <feature-id>      # Merge to main
```

Or step-by-step:

```
@build 00-001-01   # Single workstream with TDD
@build 00-001-02
@review <feature-id>
@deploy <feature-id>
```

## 4. Verify

```bash
sdp verify 00-001-01   # Check workstream completion
sdp status             # Project state
sdp next               # Recommended next action
sdp log show           # Evidence log
```

For a guided dry run of this flow:

```bash
sdp demo
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
@feature "X"  →  @oneshot <feature-id>  →  @review <feature-id>  →  @deploy <feature-id>
     │                  │                │                │
     ▼                  ▼                ▼                ▼
  Workstreams      Execute WS       APPROVED?         Merge PR
```

**Done = @review APPROVED + @deploy completed.**

## Next

- [PROTOCOL.md](PROTOCOL.md) — Full specification
- [MANIFESTO.md](MANIFESTO.md) — Vision and evidence
- [reference/](reference/) — Commands, specs, glossary
