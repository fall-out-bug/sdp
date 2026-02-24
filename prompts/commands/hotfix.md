---
description: Emergency P0 production hotfix workflow with minimal-change deployment.
agent: fixer
---

# /hotfix — Emergency Production Fixes

When calling `/hotfix "description" --issue-id=001`:

1. **Create branch** — `git checkout -b hotfix/{id}-{slug}` from master
2. **Minimal fix** — No refactoring, fix bug only
3. **Fast testing** — Smoke + critical path (no full suite)
4. **Commit** — `fix(scope): description (issue NNN)`
5. **MERGE, TAG, PUSH** — Execute yourself!
6. **Backport** — Merge to dev and feature branches
7. **Close issue** — Update status in issue file

## CRITICAL: You MUST Complete

```bash
# Merge to master and tag
git checkout master
git merge hotfix/{branch} --no-edit
git tag -a v{VERSION} -m "Hotfix: {description}"
git push origin master --tags
```

**Work is NOT complete until all `git push` commands succeed.**

## Quick Reference

**Input:** P0 CRITICAL issue  
**Output:** Production fix + pushed to origin

**Key Rules:**
- Minimal changes only
- No refactoring
- No new features
- Fast testing
- Backport mandatory
