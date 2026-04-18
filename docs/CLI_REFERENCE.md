# SDP CLI Reference

`sdp <command> --help` is the authoritative source for flags and command semantics. This document summarizes the current surfaces that are already present in the CLI.

## Recommended First-Success Path

```bash
sdp init --auto
sdp doctor
sdp demo
```

Then continue with:

```bash
sdp plan "Add auth"
sdp apply --dry-run
sdp apply
sdp status --text
sdp next
```

## Core Local Workflow

| Command | Current purpose |
|---------|-----------------|
| `sdp init` | Create `.sdp/` config and guard rules, refresh existing supported IDE integrations, and fall back to `.claude/` only when no integration exists yet |
| `sdp doctor` | Check Git, Go, supported IDE integration, and optional drift state |
| `sdp demo` | Run a temporary first-success walkthrough using `init`, `doctor`, and `status --text` |
| `sdp plan <description>` | Decompose a feature description into workstreams from the terminal |
| `sdp apply` | Execute ready workstreams with streaming progress output |
| `sdp build <ws-id>` | Execute one workstream; for the full agent-driven cycle use `/build` or `sdp orchestrate` |
| `sdp verify <ws-id>` | Verify workstream completion against outputs, verification commands, and coverage threshold |
| `sdp status` | Show project state; default output is a TUI, with `--text` and `--json` for scripts |
| `sdp next` | Recommend the next action based on workstream, git, and config state |
| `sdp log show` | Inspect evidence log events |
| `sdp deploy` | Record an approval event after merge; it does not merge branches or deploy infrastructure |

## Common Modes and Flags

| Command | Useful modes |
|---------|--------------|
| `sdp init` | `--auto`, `--dry-run`, `--interactive`, `--headless`, `--force` |
| `sdp doctor` | `--repair`, `--deep`, `--migrate`, `--rollback`, `--drift` |
| `sdp plan` | `--interactive`, `--auto-apply`, `--dry-run`, `--output=json` |
| `sdp apply` | `--ws <id>`, `--retry <n>`, `--dry-run`, `--output=json` |
| `sdp status` | default TUI, `--text`, `--json` |
| `sdp next` | `--json`, `--alternatives` |
| `sdp demo` | `--template`, `--verbose`, `--cleanup=false` |

## Broader Command Tree

The top-level help currently exposes these command groups:

| Area | Commands |
|------|----------|
| Setup and state | `init`, `doctor`, `health`, `status`, `next`, `demo`, `hooks`, `completion` |
| Planning and execution | `parse`, `plan`, `build`, `apply`, `orchestrate`, `verify`, `tdd`, `deploy` |
| Guard and session | `guard`, `session`, `resolve`, `git`, `collision` |
| Evidence and audit | `log`, `decisions`, `checkpoint`, `coordination`, `design`, `idea` |
| Quality and diagnostics | `quality`, `drift`, `diagnose`, `watch`, `contract`, `acceptance` |
| Workflow support | `beads`, `task`, `memory`, `prd`, `prototype`, `skill` |
| Telemetry and metrics | `telemetry`, `metrics` |

## Relationship to Prompt Surfaces

The CLI is optional. Core SDP prompt surfaces are installed into the supported integration directory in your repo:

- `.claude/`
- `.cursor/`
- `.opencode/`
- `.codex/`

Use [reference/skills.md](reference/skills.md) for the current prompt-surface map and [PROTOCOL.md](PROTOCOL.md) for the current protocol overview.
