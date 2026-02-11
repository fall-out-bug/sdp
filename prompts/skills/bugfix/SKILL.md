---
name: bugfix
description: Quality bug fixes (P1/P2). Full TDD cycle, branch from feature/develop, no production deploy.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /bugfix - Quality Bug Fixes

Standard bug fixes with full quality cycle.

## When to Use

- P1 (HIGH) or P2 (MEDIUM) issues
- Feature broken but not production
- Reproducible errors
- Can wait for proper testing

## Invocation

```bash
/bugfix "description" --feature=F23 --issue-id=002
```

## Workflow

1. **Read issue** — Load issue file from `docs/issues/`
2. **Create branch** — `git checkout -b bugfix/{issue-id}-{slug}` from dev
3. **TDD cycle** — Write failing test → implement fix → refactor
4. **Quality gates** — pytest, coverage ≥80%, mypy --strict, ruff
5. **Commit** — `fix(scope): description (issue NNN)`
6. **Mark issue closed** — Update status in issue file
7. **MERGE AND PUSH** — See critical section below

## CRITICAL: Completion Requirements

**You MUST execute these commands yourself. Do NOT give instructions to user.**

```bash
# 1. Merge to dev
git checkout dev
git merge bugfix/{branch-name} --no-edit

# 2. Push to remote (MANDATORY)
git pull --rebase || true
git push

# 3. Verify
git status  # MUST show "up to date with origin"
```

**Work is NOT complete until `git push` succeeds.**

## Key Difference from Hotfix

| Aspect | Hotfix | Bugfix |
|--------|--------|--------|
| Severity | P0 | P1/P2 |
| Branch from | main | develop/feature |
| Testing | Fast | Full |
| Deploy | Production | Staging |

## Output

- Bug fixed in dev branch
- Tests added with ≥80% coverage
- Issue marked closed
- Changes pushed to origin
