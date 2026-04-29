---
name: hotfix
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from main, immediate deploy.
---

# @hotfix

Emergency production fixes. Minimal changes, fast testing, merge to main with tag.

## When to Use

- P0 CRITICAL only
- Production down or severely degraded
- Data loss/corruption risk

## Workflow

1. **Branch** — `git checkout main && git pull && git checkout -b hotfix/{id}-{slug}`
2. **Minimal fix** — No refactoring, fix bug only
3. **Smoke test** — Critical path verification
4. **Merge** — `git checkout main && git merge hotfix/{branch} --no-edit`
5. **Tag** — `git tag -a v{VERSION} -m "Hotfix: {description}"`
6. **Push** — `git push origin main --tags`

## Write Plan (F101)

Before modifying any file, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (fix files and version file).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"hotfix"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @hotfix <target>:
  CREATE: path/to/new/file — <reason>
  MODIFY: path/to/existing/file — <reason>
  DELETE: path/to/removed/file — <reason>

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

## Output

Hotfix merged, tagged, pushed. Issue closed.

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- @bugfix — P1/P2 quality fixes
- @issue — Classification
