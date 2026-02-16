# Prompt Source Layout

This repository uses a single canonical prompt source:

- Skills: `prompts/skills/*/SKILL.md`
- Agents: `prompts/agents/*.md`

Compatibility adapters are provided as symlinks:

- `.claude/skills` -> `../prompts/skills`
- `.claude/agents` -> `../prompts/agents`
- `.cursor/skills` -> `../prompts/skills`
- `.cursor/agents` -> `../prompts/agents`
- `.opencode/skills` -> `../prompts/skills`
- `.opencode/agents` -> `../prompts/agents`
- `.codex/skills/sdp` -> `../../prompts/skills`
- `.codex/agents` -> `../prompts/agents`

Edit only `prompts/*` to avoid prompt drift across tools.
