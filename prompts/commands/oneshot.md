---
description: Autonomous feature execution via sdp-orchestrate outer loop.
agent: orchestrator
---

# /oneshot — Autonomous Feature Execution

When calling `/oneshot F{XX}` in Cursor:

1. Load skill: `@.claude/skills/oneshot/SKILL.md`
2. Run `sdp-orchestrate --feature F{XX} --next-action` as the outer loop
3. Execute each phase inline:
   - **build**: @build {ws_id} → commit → `sdp-orchestrate --feature F{XX} --advance --result <commit>`
   - **review**: @review F{XX} → fix P0/P1 → `sdp-orchestrate --feature F{XX} --advance`
4. PR creation and CI loop are handled by the CLI — no agent involvement
5. When done: output only `CI GREEN - @oneshot complete`

**Input:** Feature ID (e.g. F016)
**Output:** All WS executed + CI green. No "Next steps" or handoff lists.

**opencode:** Use `sdp-orchestrate --feature F{XX} --runtime opencode` as the outer loop. opencode lacks Stop hooks — the outer loop CLI replaces them.
