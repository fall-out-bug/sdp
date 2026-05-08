---
name: ship
description: Deployment orchestration. Creates PR to main (after @oneshot) or merges for release.
version: 1.1.0
changes:
  - "1.1.0: Add PI review and worktree merge cleanup checks"
  - "1.0.0: Initial release as @ship (renamed from @deploy)"
---

# @ship - Deployment Orchestration

Create PR to main (after @oneshot) or merge for release.

---

## EXECUTE THIS NOW

When user invokes `@ship F{XX}`:

### Mode 1: PR to Master (default)

**Pre-flight:** Check `.sdp/review_verdict.json` — verdict must be APPROVED and compact. Verify `git branch --show-current` is feature branch. `bd list --status open` — no P0/P1 for the feature/review round. Run quality gates (AGENTS.md). For prompt/agent/skill/eval/model-call changes, confirm PI review has no P0/P1; provider degradation must be explicitly recorded, not silently treated as PASS.

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
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"ship"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than invent placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @ship <target>:
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

## Pre-Ship

`bd list --status open --json | jq '[.[]|select(.priority<=1)]|length'` — must be 0.

---

## Git Safety

Before ANY git: verify `pwd`, `git branch --show-current`.

---

## Troubleshooting

| Issue | Fix |
|-------|-----|
| Not APPROVED | Run @review first |
| P0/P1 open | Fix before ship |
| CI failing | Quality gates locally |
| Push rejected | Pull and retry |
| `gh pr merge --delete-branch` cannot delete local branch | Merge may still have succeeded; verify PR state, then remove the feature worktree and delete the local branch from another worktree |
| Review verdict is huge or contains full provider prompts | Replace with compact schema-valid verdict before ship; do not commit raw `.sdp/runs/pi-review/*` telemetry by default |

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

- `@review` — Must be APPROVED before ship
- `@oneshot` — Autonomous execution
