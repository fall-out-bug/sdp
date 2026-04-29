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
1. **Get next action** — Run `sdp-orchestrate --feature F{XX} --next-action`. Parse the JSON output (schema: `schema/next-action.schema.json`).
2. **Execute phase and advance** — For `build`: run @build {ws_id}, commit, then `sdp-orchestrate --feature F{XX} --advance --result $(git rev-parse HEAD)`. For `review`: run @review F{XX}, fix P0/P1 until approved (max 3 iterations), then `sdp-orchestrate --feature F{XX} --advance`. **One advance per phase** — run `--advance` exactly once after build, exactly once after review. PR and CI run automatically. When action is `done`, output only: `CI GREEN - @oneshot complete`.

## Post-compaction

If context was compacted, read `.sdp/checkpoints/F{XX}.json` and `git checkout $(jq -r .branch .sdp/checkpoints/F{XX}.json)`. Resume from step 1.

## Write Plan (F101)

Before the orchestration loop begins, emit a write plan for orchestrator-owned artifacts. @build and @review emit their own detailed file plans when invoked:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason. Covers `.sdp/checkpoints`, `.sdp/evidence`, `.sdp/ws-verdicts`, and orchestrator state files. @build and @review handle their own file plans.
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"oneshot"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

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

## Claude Code

Use Task tool to spawn @build and @review subagents. Each subagent gets a fresh context window. Stop hook blocks premature exit when CI phase is incomplete.

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |
