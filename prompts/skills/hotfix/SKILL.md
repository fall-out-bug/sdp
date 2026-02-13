---
name: hotfix
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from main, immediate deploy.
version: 2.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @hotfix - Emergency Production Fixes

Fast-track critical bug fixes for production.

---

## EXECUTE THIS NOW

When user invokes `@hotfix "description"` or `@hotfix <issue-id>`:

### Step 1: Create Branch

```bash
git checkout -b hotfix/{issue-id}-{slug} main
```

Branch from main (NOT dev or feature).

### Step 2: Minimal Fix

- No refactoring!
- No new features!
- Fix bug only!

### Step 3: Fast Testing

- Smoke tests only
- Critical path verification
- No full test suite required

### Step 4: Commit

```bash
git add .
git commit -m "fix(scope): description (issue NNN)"
```

### Step 5: Merge, Tag, Push (CRITICAL)

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

### Step 6: Close Issue

Update status in issue file.

---

## When to Use

- P0 CRITICAL issues only
- Production down or severely degraded
- All/most users affected
- Data loss/corruption risk

---

## Key Rules

| Rule | Description |
|------|-------------|
| **Minimal changes** | No refactoring! |
| **No new features** | Fix bug only |
| **Fast testing** | Smoke + critical path |
| **SLA target** | Immediate (emergency) |
| **Backport mandatory** | To dev and feature branches |

---

## Output

- Hotfix merged to main with tag
- Backported to dev
- All changes pushed to origin
- Issue marked closed

---

## See Also

- `@bugfix` - Quality fixes (P1/P2)
- `@issue` - Bug classification and routing
- `@deploy` - Standard deployment
