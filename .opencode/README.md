# OpenCode Integration

This directory contains SDP integration for OpenCode.

## Prompt Surface

- Skills: `prompts/skills/`
- Commands: `prompts/commands/`
- Agents: `prompts/agents/`
- Canonical command map: `prompts/commands.yml`

## Hook Surface

OpenCode scope enforcement lives in:

- `.opencode/hooks/pre-tool-use.json`
- `.opencode/hooks/README.md`

The current hook implementation uses `sdp-omc-guard`, which is a stronger
wrapper around SDP guard semantics for edit and write operations.

## Usage

```text
@vision "product"
@feature "add feature"
@build 00-XXX-YY
@review FXXX
@operate "deploy task"
```

## Fallback Mode

If your OpenCode runtime cannot spawn subagents, follow the manual checklists in
[`docs/reference/FALLBACK_MODE.md`](../docs/reference/FALLBACK_MODE.md).
