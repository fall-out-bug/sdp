---
name: bugfix
description: Quality bug fixes (P1/P2). Full TDD cycle, branch from master via feature/, no production deploy.
version: 2.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @bugfix - Quality Bug Fixes

Standard bug fixes with full quality cycle.

---

## EXECUTE THIS NOW

When user invokes `@bugfix "description"` or `@bugfix <issue-id>`:

### Step 1: Read Issue

Load issue file from `docs/issues/` or resolve via `bd show <id>`.

### Step 2: Create Branch

```bash
git checkout master && git pull
git checkout -b fix/{issue-id}-{slug}
```

Branch from master via feature/ (NOT main).

### Step 3: TDD Cycle

1. **Red** - Write failing test that reproduces the bug
2. **Green** - Implement minimum fix to pass
3. **Refactor** - Clean up if needed

### Step 4: Quality Gates

```bash
go test ./...
go vet ./...
go build ./...
```

### Step 5: Commit

```bash
git add .
git commit -m "fix(scope): description (issue NNN)"
```

### Step 6: Merge and Push (CRITICAL)

```bash
# 1. Push branch (MANDATORY)
git pull --rebase || true
git push -u origin fix/{branch-name}

# 2. Create PR to master
gh pr create --base master --head fix/{branch-name} --title "fix: description"

# 3. Verify
git status  # MUST show "up to date with origin"
```

**Work is NOT complete until `git push` succeeds.**

---

## When to Use

- P1 (HIGH) or P2 (MEDIUM) issues
- Feature broken but not production
- Reproducible errors
- Can wait for proper testing

---

## Accepts Any Identifier Format

```bash
@bugfix "description" --feature=F23 --issue-id=002
@bugfix 99-F064-01     # Workstream ID (fix format)
@bugfix sdp-xxx        # Beads task ID
@bugfix ISSUE-0001     # Issue ID
```

**Resolution:** Uses `bd show <id>` to find task file.

---

## Key Difference from Hotfix

| Aspect | Hotfix | Bugfix |
|--------|--------|--------|
| Severity | P0 | P1/P2 |
| Branch from | master | master via feature/ |
| Testing | Fast | Full |
| Deploy | Production | Staging |

---

## Output

- Bug fixed in feature branch
- Tests added with >=80% coverage
- Issue marked closed
- Changes pushed to origin

---

## Git Safety

**CRITICAL:** Before ANY git operation, verify context.

**MANDATORY before any git command:**

```bash
# Step 1: Verify context
pwd
git branch --show-current
sdp guard context check

# Step 2: If check fails, recover
sdp guard context go $FEATURE_ID

# Step 3: Only then proceed
git add .
git commit -m "..."
```

---

## See Also

- `@hotfix` - Emergency P0 fixes (production)
- `@issue` - Bug classification and routing
- `@debug` - Systematic debugging
