# SDP — Codex Setup

This repository ships Codex-oriented guidance and compatibility files.

## Start Here

1. Read [`AGENTS.md`](../AGENTS.md)
2. Read [`.codex/AGENTS.md`](AGENTS.md) in this directory for Codex-specific guidance
3. Use [`prompts/commands.yml`](../prompts/commands.yml) as the canonical command map

## Canonical Prompt Sources

- `prompts/commands/`
- `prompts/skills/`
- `prompts/agents/`

Tool folders such as `.codex/`, `.cursor/`, and `.opencode/` are harness entry
points and lightweight adapters. Prompt logic belongs in `prompts/`.

## Typical Flow

```text
@feature "description"
@build 00-XXX-YY
@review FXXX
```

If your Codex runtime cannot spawn subagents, use the manual workflows in
[`docs/reference/FALLBACK_MODE.md`](../docs/reference/FALLBACK_MODE.md).
