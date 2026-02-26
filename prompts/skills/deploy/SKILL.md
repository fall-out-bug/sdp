---
name: deploy
description: Deployment orchestration. Creates PR to master (after @oneshot) or merges for release.
version: 4.0.0
changes:
  - "4.0.0: Compress to ~150 lines (P2 remediation)"
---

# @deploy - Deployment Orchestration

Create PR to master (after @oneshot) or merge for release.

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

## Quick Reference

| Mode | Action |
|------|--------|
| PR | feature -> master via gh pr create |
| Release | Version bump + tag on master |

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

## See Also

- `@review` — Must be APPROVED before deploy
- `@oneshot` — Autonomous execution
