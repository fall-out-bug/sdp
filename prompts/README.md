# Prompt Source Layout

This repository uses a single canonical prompt source:

- Commands: `prompts/commands/*.md`
- Skills: `prompts/skills/*/SKILL.md`
- Agents: `prompts/agents/*.md`

Compatibility adapters are provided as symlinks:

- `.claude/skills` -> `../prompts/skills`
- `.claude/agents` -> `../prompts/agents`
- `.cursor/skills` -> `../prompts/skills`
- `.cursor/agents` -> `../prompts/agents`
- `.opencode/skills` -> `../prompts/skills`
- `.opencode/agents` -> `../prompts/agents`
- `.codex/skills` -> per-skill individual symlinks (e.g. `.codex/skills/build` -> `../../prompts/skills/build`)
- `.codex/agents` -> `../prompts/agents`

Edit only `prompts/*` to avoid prompt drift across tools.

When prompts touch Go code, reference `@go-modern` from this canonical tree instead of duplicating tool-specific Go style rules.
