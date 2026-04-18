---
name: bugfix
description: Quality bug fixes (P1/P2). Full TDD cycle, branch from master via feature/, no production deploy.
---

# @bugfix

Quality bug fixes with full TDD cycle. Branch from master via feature/.

## When to Use

- P1 (HIGH) or P2 (MEDIUM) issues
- Feature broken but not production
- Reproducible errors

## Workflow

1. **Read issue** — `bd show <id>` or load from `docs/issues/`
2. **Branch** — `git checkout master && git pull && git checkout -b fix/{id}-{slug}`
3. **TDD** — Red: failing test → Green: minimal fix → Refactor
4. **Quality gates** — Run quality gates (see Quality Gates in AGENTS.md)
5. **Commit** — `git commit -m "fix(scope): description"`
6. **Push** — `git push -u origin fix/{branch}` then `gh pr create --base master`

## Write Plan (F101)

Before modifying any file, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (test files, fix files, branch operations).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl`:
   ```json
   {"ts":"<ISO-8601>","type":"write_plan","skill":"bugfix","ws_id":"<ws-id>","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]}
   ```

**Output format:**
```
WRITE PLAN for @bugfix <target>:
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

Bug fixed, tests added, issue closed, changes pushed.

## See Also

- @hotfix — P0 emergency
- @issue — Classification
- @debug — Root cause analysis
