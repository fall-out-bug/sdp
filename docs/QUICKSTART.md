# SDP Quick Start

Get from zero to your first feature in 5 minutes.

> **Product Contract:** This quickstart follows the [Local Mode (Default Path)](PRODUCT_CONTRACT.md#path-1-local-mode-default) defined in the SDP Product Contract. For the full product definition, see [PRODUCT_CONTRACT.md](PRODUCT_CONTRACT.md).

## 1. Install

**Full project** (prompts + hooks + optional CLI): default install

```bash
# Into your project (auto-detects Claude Code, Cursor, OpenCode)
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
@oneshot <feature-id>      # Autonomous: build all workstreams
@review <feature-id>       # Multi-agent quality review
@deploy <feature-id>      # Merge to main
```

**Discovery utility (corpus audit):**

```
@strataudit "Audit the strategy corpus in docs/strategy and docs/architecture"
```

**Evidence-pack from existing artifacts:**

```
@strataudit "Build an evidence pack from the existing .strataudit artifacts"
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

- [PRODUCT_CONTRACT.md](PRODUCT_CONTRACT.md) — Product definition (paths, stages, harness support)
- [PROTOCOL.md](PROTOCOL.md) — Full specification
- [MANIFESTO.md](MANIFESTO.md) — Vision and evidence
- [reference/](reference/) — Commands, specs, glossary
