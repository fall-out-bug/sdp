# /hotfix — Emergency Production Fixes

When calling `/hotfix "description" --issue-id=001`:

1. **Create branch** — `git checkout -b hotfix/{id}-{slug}` from main
2. **Minimal fix** — No refactoring, fix bug only
3. **Fast testing** — Smoke + critical path (no full suite)
4. **Commit** — `fix(scope): description (issue NNN)`
5. **MERGE, TAG, PUSH** — Execute yourself!
6. **Backport** — Merge to dev and feature branches
7. **Close issue** — Update status in issue file

## CRITICAL: You MUST Complete

```bash
# Merge to main and tag
git checkout main
git merge hotfix/{branch} --no-edit
git tag -a v{VERSION} -m "Hotfix: {description}"
git push origin main --tags

# Backport to dev
git checkout dev
git merge main --no-edit
git push origin dev
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
