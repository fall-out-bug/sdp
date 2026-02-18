# OpenCode Integration

This directory contains SDP integration for OpenCode.

## Setup

SDP skills and agents are available via symlinks:
- `skills/` → All SDP skills (@vision, @build, @review, etc.)
- `agents/` → Agent definitions (orchestrator, reviewer, etc.)

Primary agent cards shown in OpenCode are configured in `.opencode/opencode.json`.
Agent metadata (name/description/prompt body) comes from files in `prompts/agents/*.md` via the `agents/` symlink.
Each agent file must have valid closed YAML frontmatter so descriptions render in UI.

## Usage

Skills work the same as in Claude Code:

```
@vision "your product"     # Strategic planning
@feature "add feature"     # Plan feature
@build 00-001-01           # Execute workstream
@review F01                # Quality check
```

## Commands

If your tool supports slash commands, create command files pointing to skills:

```
/oneshot → skills/oneshot/SKILL.md
/build   → skills/build/SKILL.md
/review  → skills/review/SKILL.md
```
