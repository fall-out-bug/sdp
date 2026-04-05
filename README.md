# SDP: Structured Development Protocol

**Prompt bundle plus optional local CLI for stage-based AI work with evidence.**

SDP installs prompt and agent surfaces into supported IDE integrations and can install the `sdp` CLI for local setup, planning, execution, and inspection. The public install flow supports `Claude Code`, `Cursor`, `OpenCode` / `Windsurf`, and `Codex`.

> [Manifesto](docs/MANIFESTO.md) — why SDP exists and where the project is headed.

## Quick Start

```bash
# Install prompts, hooks, and optional CLI
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh

# Or install only the CLI binary
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | sh -s -- --binary-only

# Or vendor SDP as a submodule
git submodule add https://github.com/fall-out-bug/sdp.git sdp
```

Installer auto-detects `Claude Code`, `Cursor`, `OpenCode` / `Windsurf`, and `Codex`. If detection misses your tool, rerun with `SDP_IDE=claude|cursor|opencode|codex`.

If you embed SDP as a submodule inside another repo, keep the public GitHub URL in `.gitmodules`. Do not point the submodule at a local sibling path such as `../sdp`.

**Recommended first success:**

```bash
sdp init --auto
sdp doctor
sdp demo
```

After that, use the local CLI flow:

```bash
sdp plan "Add auth"
sdp apply --dry-run
sdp apply
```

→ [5-minute guide](docs/QUICKSTART.md)

## Two Modes

| Mode | Use when | What it needs |
|------|----------|---------------|
| Local Mode | You want the fastest first success in one repo | `sdp` CLI plus one supported IDE integration |
| Operator Mode | You already run queue-backed work with workstreams and operators | Prompt surfaces plus `Beads` for the live queue |

**Local Mode is the recommended starting point.** It works without Beads and is the current public onboarding path.

**Operator Mode is advanced.** Use it only if you already want a board-backed queue and multi-session execution. Beads is required for that operating model.

## Current Workflow Surfaces

| Stage | Local CLI | Prompt surface | Notes |
|-------|-----------|----------------|-------|
| Bootstrap | `sdp init`, `sdp doctor`, `sdp demo` | Installed into `.claude/`, `.cursor/`, `.opencode/`, or `.codex/` | `sdp init` refreshes existing integrations and falls back to `.claude/` only when none exists yet |
| Plan | `sdp plan` | `/feature`, `/idea`, `/design` | CLI is the clearest first-run path |
| Execute | `sdp apply`, `sdp build` | `/build`, `/oneshot` | `sdp build` executes one workstream; `sdp apply` runs ready workstreams |
| Verify and inspect | `sdp verify`, `sdp status`, `sdp next`, `sdp log show` | `/review` | `status` and `next` are current inspection surfaces |
| Record approval | `sdp deploy` | `/deploy` | `sdp deploy` records an approval event after merge; it does not merge branches or deploy infrastructure |

## What SDP Installs

1. `.sdp/` project config, guard rules, and evidence paths.
2. Prompt and agent adapters for the supported IDE integration already present in your repo.
3. Optional CLI helpers for setup, planning, execution, verification, and diagnostics.

Canonical prompt sources live in `prompts/`. Tool-specific directories such as `.claude/`, `.cursor/`, `.opencode/`, and `.codex/` are adapters around that source tree.

## Optional Components

- **CLI:** `sdp init`, `sdp doctor`, `sdp plan`, `sdp apply`, `sdp status`, `sdp next`, `sdp log`, `sdp demo`
- **Beads:** `brew tap beads-dev/tap && brew install beads` for board-backed, multi-session work
- **Platform note:** some evidence helpers rely on `flock`, so macOS/Linux is the tested path

## Docs

| File | Content |
|------|---------|
| [QUICKSTART.md](docs/QUICKSTART.md) | Recommended first-success path |
| [CLI_REFERENCE.md](docs/CLI_REFERENCE.md) | Current `sdp` command surfaces |
| [PROTOCOL.md](docs/PROTOCOL.md) | Current protocol overview |
| [reference/README.md](docs/reference/README.md) | Reference index and legacy-doc status |
| [.codex/INSTALL.md](.codex/INSTALL.md) | Codex-specific install notes |
| [MANIFESTO.md](docs/MANIFESTO.md) | Vision and rationale |
| [ROADMAP.md](docs/ROADMAP.md) | Product direction |

## License

MIT
