---
name: deploy
description: Deployment orchestration. Creates PR to dev or merges dev to main for release.
tools: Read, Write, Shell, Glob, Grep
version: 2.0.0
---

# @deploy - Deployment Orchestration

Create PR to dev (after @oneshot) or merge dev to main for release.

## Invocation

```bash
@deploy F020              # Create PR to dev (after review passed)
@deploy F020 --release    # Merge dev to main with version bump
@deploy F020 --release minor  # Minor version bump
```

## Two Deployment Modes

### Mode 1: PR to Dev (default)

Used after `@oneshot F020` completes with review passed.

```
feature/F020-xxx → dev (via PR)
```

**Steps:**
1. Verify review passed
2. Verify no blocking findings (`sdp guard finding list`)
3. Push branch
4. Create PR to dev
5. CI validates

### Mode 2: Release to Main (`--release`)

Used after PR merged to dev and human UAT complete.

```
dev → main (with version bump)
```

**Steps:**
1. Version resolution
2. Generate artifacts (CHANGELOG, release notes)
3. Commit artifacts to dev
4. Merge dev to main (--no-ff)
5. Tag + Push
6. Report

## Quick Reference

| Mode | Command | Action |
|------|---------|--------|
| PR | `@deploy F020` | Create PR: feature → dev |
| Release | `@deploy F020 --release` | Merge: dev → main |

## Workflow: Mode 1 (PR to Dev)

### Step 1: Pre-flight Checks

```bash
# Verify on feature branch
git branch --show-current  # Should be feature/F020-xxx

# Verify no blocking findings
sdp guard finding list
# Must show: "0 blocking"

# Verify tests pass
go test ./... -q
```

**Gate:** If blocking findings or tests fail → STOP.

### Step 2: Push and Create PR

```bash
# Push feature branch
git push origin feature/F020-xxx

# Create PR to dev
gh pr create \
    --base dev \
    --head feature/F020-xxx \
    --title "feat(F020): Feature Title" \
    --body "## Summary
{summary_from_idea_file}

## Workstreams
{list_of_completed_workstreams}

## Test plan
- [x] All workstreams completed
- [x] Review passed
- [x] Tests pass locally

🤖 Generated with [Claude Code](https://claude.com/claude-code)"
```

### Step 3: Report

```
✅ PR Created: https://github.com/owner/repo/pull/123
Base: dev
Head: feature/F020-xxx
CI: Running...

Next steps:
1. Wait for CI to pass
2. Human UAT (5-10 min)
3. Merge PR when ready
4. Run @deploy F020 --release for production
```

## Workflow: Mode 2 (Release to Main)

### Step 1: Pre-flight Checks

```bash
# Verify on dev branch
git branch --show-current  # Should be dev

# Verify dev is up to date
git pull origin dev

# Verify tests pass
go test ./... -q
```

### Step 2: Version Resolution

Read current version from `go.mod` or version file. Bump based on:
- `patch` (default): 0.5.0 → 0.5.1
- `minor`: 0.5.0 → 0.6.0
- `major`: 0.5.0 → 1.0.0

### Step 3: Generate Artifacts

```bash
# Update CHANGELOG.md
# Create docs/releases/v{X.Y.Z}.md
```

### Step 4: Commit Artifacts

```bash
git add CHANGELOG.md docs/releases/
git commit -m "chore(release): v{X.Y.Z}"
```

### Step 5: Merge to Main

```bash
git checkout main
git pull origin main
git merge dev --no-ff -m "Release v{X.Y.Z}: F020 Feature Title"
```

### Step 6: Tag + Push

```bash
git tag -a v{X.Y.Z} -m "Release v{X.Y.Z}"
git push origin main
git push origin v{X.Y.Z}
git checkout dev
```

### Step 7: Report

```
✅ Released: v{X.Y.Z}
Tag: v{X.Y.Z}
Commit: abc123
Features: F020

CHANGELOG: docs/releases/v{X.Y.Z}.md
```

## Guard Integration

Before any deployment, check for blocking findings:

```bash
sdp guard finding list

# If blocking findings exist:
sdp guard finding resolve finding-xxx --by="Fixed in commit abc123"
sdp guard finding clear
```

## Troubleshooting

| Issue | Solution |
|-------|----------|
| PR creation fails | Check branch exists and is pushed |
| CI failing | Run `go test ./...` locally |
| Blocking findings | `sdp guard finding list` then fix |
| Merge conflict | Resolve in feature branch first |

---

**Version:** 2.0.0 (PR-based deployment)
**See Also:** `@oneshot`, `@review`, `@build`

Output summary:

```markdown
## Deploy Complete: v{X.Y.Z}

**Feature:** {FXX} - {Title}
**Tag:** v{X.Y.Z}
**Branch:** main

### Artifacts Created
- pyproject.toml (version bump)
- CHANGELOG.md (release entry)
- docs/releases/v{X.Y.Z}.md

### Git Operations
- [x] Committed release artifacts
- [x] Merged dev → main
- [x] Tagged v{X.Y.Z}
- [x] Pushed to origin
```

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| Tests fail | Pre-flight failed | Fix tests first |
| Not APPROVED | Review pending | Run @review first |
| Merge conflict | Diverged branches | Resolve manually |
| Push rejected | Remote ahead | Pull and retry |

## Git Safety

**CRITICAL:** Before ANY git operation, verify context.

See [GIT_SAFETY.md](../../.claude/GIT_SAFETY.md) for full guidelines.

**MANDATORY before any git command:**

```bash
# Step 1: Verify context
pwd
git branch --show-current
sdp guard context check

# Step 2: If check fails, recover
sdp guard context go $FEATURE_ID

# Step 3: Only then proceed with deployment
```

**NOTE:** Deployment typically merges to main, which is allowed for @deploy.

## See Also

- [@review skill](../review/SKILL.md) — Must be APPROVED before deploy
- [Release Notes Template](../../templates/release-notes.md)
