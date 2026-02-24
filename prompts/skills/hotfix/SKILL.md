---
name: hotfix
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from master, immediate deploy.
---

# @hotfix

Emergency production fixes. Minimal changes, fast testing, merge to master with tag.

## When to Use

- P0 CRITICAL only
- Production down or severely degraded
- Data loss/corruption risk

## Workflow

1. **Branch** — `git checkout master && git pull && git checkout -b hotfix/{id}-{slug}`
2. **Minimal fix** — No refactoring, fix bug only
3. **Smoke test** — Critical path verification
4. **Merge** — `git checkout master && git merge hotfix/{branch} --no-edit`
5. **Tag** — `git tag -a v{VERSION} -m "Hotfix: {description}"`
6. **Push** — `git push origin master --tags`

## Output

Hotfix merged, tagged, pushed. Issue closed.

## See Also

- @bugfix — P1/P2 quality fixes
- @issue — Classification
