# /oneshot — Autonomous Feature Execution

When calling `/oneshot F{XX}`:

1. Load full prompt: `@prompts/commands/oneshot.md`
2. Create PR and wait for approval
3. Execute all feature WS autonomously
4. Save checkpoints
5. Handle errors (auto-fix or escalate)
6. Run `/review` at the end
7. Output summary

## Quick Reference

**Input:** Feature ID (F60)
**Output:** All WS executed + Review + UAT guide
**Features:**
- PR approval gate
- Checkpoint/resume support
- Progress tracking JSON
- Auto-fix MEDIUM/HIGH errors
- Telegram notifications

**Next:** Human UAT → `/deploy F{XX}`

## Checkpoint Files

- `.oneshot/F{XX}-checkpoint.json` - Resume state
- `.oneshot/F{XX}-progress.json` - Real-time metrics
