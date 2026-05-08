# Cursor Integration

**Primary context source:** [`.cursorrules`](../.cursorrules)

Cursor reads `.cursorrules` automatically at the project root. Keep that file
small and operational: role, decision tree, quality gates, and where to find
canonical prompts.

## Canonical Prompt Source

- Commands: `../prompts/commands/`
- Skills: `../prompts/skills/`
- Agents: `../prompts/agents/`
- Command mapping: `../prompts/commands.yml`

Edit canonical prompt sources only. Do not fork Cursor-only copies of prompt
logic unless the harness format itself requires it.

## Usage

Use `@` commands for the main SDP flows:

```text
@feature "description"
@idea "description"
@design 00-XXX-YY
@build 00-XXX-YY
@review FXXX
@fix "regression description"
@operate "deploy or release task"
```

## Fallback Mode

If your Cursor runtime cannot spawn subagents, use the manual checklists in
[`docs/reference/FALLBACK_MODE.md`](../docs/reference/FALLBACK_MODE.md).

## See Also

- [`AGENTS.md`](../AGENTS.md)
- [`docs/reference/project-map.md`](../docs/reference/project-map.md)
- [`docs/reference/FALLBACK_MODE.md`](../docs/reference/FALLBACK_MODE.md)
- [`prompts/commands.yml`](../prompts/commands.yml)
