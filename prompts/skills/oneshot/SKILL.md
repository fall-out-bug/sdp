---
name: oneshot
description: Autonomous feature execution via sdp orchestrate outer loop
cli: sdp-orchestrate
version: 9.0.0
changes:
  - F016: Outer loop — sdp-orchestrate drives phases; LLM only for @build and @review
  - Slim prompt: 3 rules, positive framing
  - PR and CI handled by CLI
---

# oneshot

Outer loop: `sdp-orchestrate` (or `sdp orchestrate` if available) drives phases. You execute @build and @review inline. Ensure `sdp-orchestrate` is on PATH (see AGENTS.md build instructions).

## Rules

1. **Get next action** — Run `sdp-orchestrate --feature F{XX} --next-action`. Parse the JSON output.
2. **Execute phase** — For `build`: run @build {ws_id}, commit, then `sdp-orchestrate --feature F{XX} --advance --result $(git rev-parse HEAD)`. For `review`: run @review F{XX}, fix P0/P1 until approved (max 3 iterations), then `sdp-orchestrate --feature F{XX} --advance`.
3. **Advance** — After each phase, run `sdp-orchestrate --feature F{XX} --advance`. PR and CI run automatically. When action is `done`, output only: `CI GREEN - @oneshot complete`.

## Post-compaction

If context was compacted, read `.sdp/checkpoints/F{XX}.json` and `git checkout $(jq -r .branch .sdp/checkpoints/F{XX}.json)`. Resume from step 1.

## Claude Code

Use Task tool to spawn @build and @review subagents. Each subagent gets a fresh context window. Stop hook (F015) blocks premature exit when CI phase is incomplete.
