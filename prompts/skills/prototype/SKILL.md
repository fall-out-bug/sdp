---
name: prototype
description: Rapid prototyping shortcut for experienced vibecoders
---

# @prototype

Ultra-fast feature planning: 5-question interview → 1-3 workstreams → @oneshot with relaxed gates.

## When to Use

Experienced devs, need prototype fast, tech debt acceptable. Not for production or security-critical.

## Gate Overrides

| Gate | Normal | Prototype |
|------|--------|-----------|
| TDD | Required | Optional |
| Coverage | ≥80% | None |
| Architecture | Clean | Monolithic OK |

Non-negotiable: code compiles, runs, no crashes, basic security.

## Workflow

1. AskUserQuestion: problem, scope, dependencies, risks, success
2. Generate 1-3 workstreams
3. Launch @oneshot

## Output

`docs/drafts/prototype-{id}.md`, `docs/workstreams/backlog/00-FFF-*.md`

## Write Plan (F101)

Before creating prototype scaffolding, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (prototype doc, workstream files, event log).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"prototype"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @prototype <feature>:
  CREATE: docs/drafts/prototype-{id}.md — Prototype specification
  CREATE: docs/workstreams/backlog/00-FFF-*.md — Workstream files (1-3)
  MODIFY: .sdp/log/events.jsonl — Write plan event log

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

---

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- @feature — Full planning
- @oneshot — Execution
