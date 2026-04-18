---
name: hotfix
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from master, immediate deploy.
---

# @hotfix

Emergency production fixes. Minimal changes, fast testing, merge to master with tag.

## When to Use

- P0 CRITICAL only
- Production down or severely degraded
- Data loss/corruption risk

## Workflow

1. **Branch** — `git checkout master && git pull && git checkout -b hotfix/{id}-{slug}`
2. **Minimal fix** — No refactoring, fix bug only
3. **Smoke test** — Critical path verification
4. **Merge** — `git checkout master && git merge hotfix/{branch} --no-edit`
5. **Tag** — `git tag -a v{VERSION} -m "Hotfix: {description}"`
6. **Push** — `git push origin master --tags`

## Write Plan (F101)

Before modifying any file, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (minimal fix files, tag, merge operations).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl`:
   ```json
   {"ts":"<ISO-8601>","type":"write_plan","skill":"hotfix","ws_id":"<ws-id>","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]}
   ```

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

## See Also

- @bugfix — P1/P2 quality fixes
- @issue — Classification
