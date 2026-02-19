---
description: Autonomous multi-step feature execution with checkpoints and review loop.
agent: orchestrator
---

# /oneshot — Autonomous Feature Execution

**Note:** This is Cursor-specific command. For Claude Code, use Task tool-based orchestration (see `.claude/skills/oneshot/SKILL.md`).

When calling `/oneshot F{XX}` in Cursor:

1. Load full prompt: `@.claude/skills/oneshot.md`
2. Follow TodoWrite tracking (create todo list at start)
3. Create PR and wait for approval
4. Execute all feature WS autonomously (inline, no Task tool)
5. Save checkpoints with progress
6. Handle errors (auto-fix or escalate)
7. Run `/review` at the end
8. Update TodoWrite: mark all completed
9. Output summary

## Quick Reference

**Input:** Feature ID (F60)
**Output:** All WS executed + Review + UAT guide

**Features:**
- TodoWrite progress tracking (real-time UI updates)
- PR approval gate
- Checkpoint/resume support
- Progress tracking JSON
- Auto-fix MEDIUM/HIGH errors
- Telegram notifications

**Difference from Claude Code:**
- Cursor: Inline execution (no Task tool)
- Claude Code: Task tool orchestrator with isolated agent

**Next:** Human UAT → `/deploy F{XX}`

## Checkpoint Files

- `.oneshot/F{XX}-checkpoint.json` - Resume state (includes agent_id for Claude Code)
- `.oneshot/F{XX}-progress.json` - Real-time metrics
