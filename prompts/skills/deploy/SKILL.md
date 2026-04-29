---
name: deploy
description: DEPRECATED: Use @ship instead. Deployment orchestration. Creates PR to main (after @oneshot) or merges for release.
version: 5.0.0
deprecated: true
deprecated_in_favor_of: ship
deprecation_version: "5.0.0"
removal_version: "8.0.0"
changes:
  - "5.0.0: DEPRECATED - Renamed to @ship"
  - "4.0.0: Compress to ~150 lines (P2 remediation)"
---

# @deploy - DEPRECATED

⚠️ **DEPRECATED** — Use `@ship` instead.

This skill will be removed in version 8.0.0. Both `@deploy` and `@ship` will work for 3 minor versions.

---

## Deployment Orchestration (Legacy)

Create PR to main (after @oneshot) or merge for release.

---

## EXECUTE THIS NOW

When user invokes `@deploy F{XX}`:

### Mode 1: PR to Master (default)

**Pre-flight:** Check `.sdp/review_verdict.json` — verdict must be APPROVED. Verify `git branch --show-current` is feature branch. `bd list --status open` — no P0/P1. Run quality gates (AGENTS.md).

**Steps:** Push feature branch. Base branch: `git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's|.*/||'` (or `main`). `gh pr create --base {base} --head feature/F{XX}-xxx --title "feat(F{XX}): ..." --body "..."`. Do not hardcode `master`.

**Report:** PR Created: {url}. CI: Running...

### Mode 2: Release (`--release`)

**Pre-flight:** On default branch (main/master). `git pull`. Quality gates pass.

**Steps:** Detect version file (go.mod, package.json, Cargo.toml, etc.). Bump (patch/minor/major). Update CHANGELOG.md, docs/releases/v{X.Y.Z}.md. Commit. Tag v{X.Y.Z}. Push default branch + tag.

**Report:** Released: v{X.Y.Z}. Tag: v{X.Y.Z}.

---

## Write Plan (F101)

Before modifying any file, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (CHANGELOG.md, release docs, version file).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"deploy"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @deploy <target>:
  CREATE: path/to/new/file — <reason>
  MODIFY: path/to/existing/file — <reason>
  DELETE: path/to/removed/file — <reason>

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

## Quick Reference

| Mode | Action |
|------|--------|
| PR | feature -> main via gh pr create |
| Release | Version bump + tag on main |

---

## Pre-Deploy

`bd list --status open --json | jq '[.[]|select(.priority<=1)]|length'` — must be 0.

---

## Git Safety

Before ANY git: verify `pwd`, `git branch --show-current`.

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Not APPROVED | Run @review first |
| P0/P1 open | Fix before deploy |
| CI failing | Quality gates locally |
| Push rejected | Pull and retry |

---

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- `@review` — Must be APPROVED before deploy
- `@oneshot` — Autonomous execution
