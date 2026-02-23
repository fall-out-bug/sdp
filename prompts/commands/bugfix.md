---
description: Quality bug-fix workflow for P1/P2 issues with full TDD cycle.
agent: builder
---

# /bugfix — Quality Bug Fixes

When calling `/bugfix issue NNN`:

1. **Read issue** — Load `docs/issues/{NNN}-*.md`
2. **Create branch** — `git checkout -b bugfix/{NNN}-{slug}` from dev
3. **TDD cycle** — Write failing test → implement fix → refactor
4. **Quality gates** — pytest, coverage ≥80%, mypy --strict, ruff
5. **Commit** — `fix(scope): description (issue NNN)`
6. **Mark issue closed** — Update status in issue file
7. **MERGE AND PUSH** — Execute yourself, not instructions!

## CRITICAL: You MUST Complete

```bash
git checkout dev
git merge bugfix/{branch} --no-edit
git push
git status  # MUST show "up to date with origin"
```

**Work is NOT complete until `git push` succeeds.**

## Quick Reference

**Input:** P1/P2 issue  
**Output:** Bug fixed + tests + pushed to origin

| Aspect | Hotfix | Bugfix |
|--------|--------|--------|
| Severity | P0 | P1/P2 |
| Branch from | main | dev |
| Testing | Fast | Full |
