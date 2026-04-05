# SDP Protocol Overview

This document describes the current public SDP model at a high level. For exact command behavior, use `sdp <command> --help`.

## What Ships Today

SDP currently ships three layers:

1. **Canonical prompt sources** in `prompts/commands/`, `prompts/skills/`, and `prompts/agents/`
2. **Harness adapters** in `.claude/`, `.cursor/`, `.opencode/`, and `.codex/`
3. **Optional Go CLI** in `sdp-plugin/` for setup, planning, execution, verification, and inspection

The CLI is a convenience surface, not the only way to use SDP. Prompt surfaces are installed into supported IDE integrations and can be used independently.

## Two Operating Modes

| Mode | Current use | Required components |
|------|-------------|---------------------|
| Local Mode | First-run onboarding and single-repo use | `sdp` CLI plus one supported IDE integration |
| Operator Mode | Queue-backed, multi-session work | Prompt surfaces plus `Beads` |

**Local Mode** is the recommended public starting point.

**Operator Mode** is advanced and assumes you already want live queue management. If you want board-backed operation, Beads is part of the contract.

## Current Stage Model

| Stage | Main surfaces | Current outcome |
|-------|---------------|-----------------|
| Bootstrap | `sdp init`, `sdp doctor`, `sdp demo` | Create `.sdp/` config and verify environment |
| Plan | `sdp plan`, prompt planning commands such as `/feature` | Create or refine workstreams |
| Execute | `sdp apply`, `sdp build`, prompt execution commands such as `/build` and `/oneshot` | Run ready workstreams |
| Verify and inspect | `sdp verify`, `sdp status`, `sdp next`, `sdp log show`, `/review` | Check completion and inspect state |
| Record approval | `sdp deploy`, `/deploy` | Record post-merge approval or follow harness-specific release flow |

Important distinction:

- `sdp deploy` records an approval event after merge.
- `sdp deploy` does not merge branches.
- `sdp deploy` does not deploy infrastructure.

## Supported Integrations

The current public installer supports:

- `Claude Code`
- `Cursor`
- `OpenCode` / `Windsurf`
- `Codex`

`sdp init` refreshes the supported integration already present in the repo and creates `.claude/` only as a fallback when no supported integration exists yet.

## Current Artifacts

Common runtime artifacts today:

- `.sdp/config.yml`
- `.sdp/guard-rules.yml`
- `.sdp/` evidence and local state files
- workstream documents created by planning flows

Canonical prompt content lives in `prompts/`. Tool-specific directories are adapters around that source tree.

## Authoritative Sources

Use these sources in order:

1. `sdp <command> --help` for CLI behavior
2. [README.md](../README.md) and [QUICKSTART.md](QUICKSTART.md) for onboarding
3. [CLI_REFERENCE.md](CLI_REFERENCE.md) for the current command map
4. [reference/skills.md](reference/skills.md) for prompt-surface layout
5. `prompts/` source files when you need exact prompt definitions

## Legacy Note

Older deep reference documents in `docs/reference/` include historical design material from earlier iterations. If a legacy note disagrees with CLI help or the onboarding docs above, trust the runtime help and current onboarding docs.
