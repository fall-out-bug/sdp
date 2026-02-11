---
name: hotfix
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from main, immediate deploy.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /hotfix - Emergency Production Fixes

Fast-track critical bug fixes for production.

## When to Use

- P0 CRITICAL issues only
- Production down or severely degraded
- All/most users affected
- Data loss/corruption risk

## Invocation

```bash
/hotfix "description" --issue-id=001
```

## Workflow

1. **Create branch** — `git checkout -b hotfix/{issue-id}-{slug}` from main
2. **Minimal fix** — No refactoring, fix bug only
3. **Fast testing** — Smoke + critical path (no full suite)
4. **Commit** — `fix(scope): description (issue NNN)`
5. **MERGE, TAG, PUSH** — See critical section below
6. **Backport** — Merge to dev and feature branches
7. **Close issue** — Update status in issue file

## CRITICAL: Completion Requirements

**You MUST execute these commands yourself. Do NOT give instructions to user.**

```bash
# 1. Merge to main and tag
git checkout main
git merge hotfix/{branch} --no-edit
git tag -a v{VERSION} -m "Hotfix: {description}"
git push origin main --tags

# 2. Backport to dev
git checkout dev
git merge main --no-edit
git push origin dev

# 3. Verify
git status  # MUST show "up to date with origin"
```

**Work is NOT complete until all `git push` commands succeed.**

## Key Rules

- **Minimal changes** — No refactoring!
- **No new features** — Fix bug only
- **Fast testing** — Smoke + critical path
- **SLA target: Immediate** — Emergency response
- **Backport mandatory** — To dev and feature branches

## Output

- Hotfix merged to main with tag
- Backported to dev
- All changes pushed to origin
- Issue marked closed
