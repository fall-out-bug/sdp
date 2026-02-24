---
description: Quality bug-fix workflow for P1/P2 issues with full TDD cycle.
agent: builder
---

# /bugfix — Quality Bug Fixes

When calling `/bugfix issue NNN`:

1. **Read issue** — Load `docs/issues/{NNN}-*.md`
2. **Create branch** — `git checkout -b bugfix/{NNN}-{slug}` from master
3. **TDD cycle** — Write failing test → implement fix → refactor
4. **Quality gates** — run quality gates (see AGENTS.md)
5. **Commit** — `fix(scope): description (issue NNN)`
6. **Mark issue closed** — Update status in issue file
7. **MERGE AND PUSH** — Execute yourself, not instructions!

## CRITICAL: You MUST Complete

```bash
git checkout master
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
| Branch from | master | master |
| Testing | Fast | Full |
