---
name: oneshot
description: Autonomous feature execution via sdp orchestrate outer loop
cli: sdp-orchestrate
version: 9.0.0
changes:
  - Outer loop — sdp-orchestrate drives phases; LLM only for @build and @review
  - Slim prompt: 3 rules, positive framing
  - PR and CI handled by CLI
---

# oneshot

Outer loop: `sdp-orchestrate` (or `sdp orchestrate` if available) drives phases. You execute @build and @review inline.

**Run orchestrate:** Either `sdp-orchestrate` on PATH, or from project root: `go run ./cmd/sdp-orchestrate`. See AGENTS.md for build/install.

## Rules

0. **Scope** — Do not change workstream scope mid-run. If scope must change, stop and start a new run.
1. **Get next action** — Run `sdp-orchestrate --feature F{XX} --next-action`. Parse the JSON output (schema: `sdp/schema/next-action.schema.json`).
2. **Execute phase and advance** — For `build`: run @build {ws_id}, commit, then `sdp-orchestrate --feature F{XX} --advance --result $(git rev-parse HEAD)`. For `review`: run @review F{XX}, fix P0/P1 until approved (max 3 iterations), then `sdp-orchestrate --feature F{XX} --advance`. **One advance per phase** — run `--advance` exactly once after build, exactly once after review. PR and CI run automatically. When action is `done`, output only: `CI GREEN - @oneshot complete`.

## Write Plan (F101)

Before the orchestration loop begins, emit a write plan covering all files touched across phases:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason. Covers workstream verdicts, evidence files, checkpoints, and any source files modified by @build.
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl`:
   ```json
   {"ts":"<ISO-8601>","type":"write_plan","skill":"oneshot","ws_id":"<ws-id>","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]}
   ```

**Output format:**
```
WRITE PLAN for @oneshot <feature-id>:
  CREATE: path/to/new/file — <reason>
  MODIFY: path/to/existing/file — <reason>
  DELETE: path/to/removed/file — <reason>

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

## Post-compaction

If context was compacted, read `.sdp/checkpoints/F{XX}.json` and `git checkout $(jq -r .branch .sdp/checkpoints/F{XX}.json)`. Resume from step 1.

## Claude Code

Use Task tool to spawn @build and @review subagents. Each subagent gets a fresh context window. Stop hook blocks premature exit when CI phase is incomplete.
