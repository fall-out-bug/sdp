---
name: deploy
description: Deployment orchestration. Creates PR to dev or merges dev to main for release.
version: 3.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @deploy - Deployment Orchestration

Create PR to dev (after @oneshot) or merge dev to main for release.

---

## EXECUTE THIS NOW

When user invokes `@deploy F020`:

### Mode 1: PR to Dev (default)

Used after `@oneshot F020` completes with review passed.

```
feature/F020-xxx -> dev (via PR)
```

**Steps:**

1. **Pre-flight Checks**
   ```bash
   # CRITICAL: Check review verdict
   if [ -f .sdp/review_verdict.json ]; then
     verdict=$(jq -r '.verdict' .sdp/review_verdict.json)
     if [ "$verdict" != "APPROVED" ]; then
       echo "ERROR: Review not approved. Run @review first."
       exit 1
     fi
   else
     echo "ERROR: No review verdict found. Run @review first."
     exit 1
   fi

   # Verify on feature branch
   git branch --show-current  # Should be feature/F020-xxx

   # Verify no blocking beads issues (P0/P1 open)
   bd list --status open --json | jq '[.[] | select(.priority <= 1)] | length'
   # Must be 0

   # Verify tests pass
   go test ./... -q
   ```

   **Gate:** If no APPROVED review, blocking findings, or tests fail -> STOP.

2. **Push and Create PR**
   ```bash
   # Push feature branch
   git push origin feature/F020-xxx

   # Create PR to master
   gh pr create \
       --base master \
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
   "
   ```

3. **Report**
   ```
   PR Created: https://github.com/owner/repo/pull/123
   Base: master
   Head: feature/F020-xxx
   CI: Running...

   Next steps:
   1. Wait for CI to pass
   2. Human UAT (5-10 min)
   3. Merge PR when ready
   4. Run @deploy F020 --release for production
   ```

### Mode 2: Release to Main (`--release`)

Used after PR merged to master and human UAT complete.

```
master -> release (with version bump)
```

**Steps:**

1. **Pre-flight Checks**
   ```bash
   # Verify on master branch
   git branch --show-current  # Should be master

   # Verify master is up to date
   git pull origin master

   # Verify tests pass
   go test ./... -q
   ```

2. **Version Resolution**

   Read current version from `go.mod` or version file. Bump based on:
   - `patch` (default): 0.5.0 -> 0.5.1
   - `minor`: 0.5.0 -> 0.6.0
   - `major`: 0.5.0 -> 1.0.0

3. **Generate Artifacts**
   ```bash
   # Update CHANGELOG.md
   # Create docs/releases/v{X.Y.Z}.md
   ```

4. **Commit Artifacts**
   ```bash
   git add CHANGELOG.md docs/releases/
   git commit -m "chore(release): v{X.Y.Z}"
   ```

5. **Merge to Main**
   ```bash
   git checkout main
   git pull origin main
   git merge dev --no-ff -m "Release v{X.Y.Z}: F020 Feature Title"
   ```

6. **Tag + Push**
   ```bash
   git tag -a v{X.Y.Z} -m "Release v{X.Y.Z}"
   git push origin main
   git push origin v{X.Y.Z}
   git checkout dev
   ```

7. **Report**
   ```
   Released: v{X.Y.Z}
   Tag: v{X.Y.Z}
   Commit: abc123
   Features: F020

   CHANGELOG: docs/releases/v{X.Y.Z}.md
   ```

---

## Quick Reference

| Mode | Command | Action |
|------|---------|--------|
| PR | `@deploy F020` | Create PR: feature -> master |
| Release | `@deploy F020 --release` | Merge: master -> main |

---

## Pre-Deploy Check

Before any deployment, verify no blocking beads issues (P0/P1 open):

```bash
bd list --status open --json | jq '[.[] | select(.priority <= 1)] | length'
# Must be 0 before deploy
```

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

# Step 3: Only then proceed with deployment
```

**NOTE:** Deployment typically merges to main, which is allowed for @deploy.

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| PR creation fails | Check branch exists and is pushed |
| CI failing | Run `go test ./...` locally |
| Blocking issues | `bd list --status open` then fix P0/P1 |
| Merge conflict | Resolve in feature branch first |

---

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| Tests fail | Pre-flight failed | Fix tests first |
| Not APPROVED | Review pending | Run @review first |
| Merge conflict | Diverged branches | Resolve manually |
| Push rejected | Remote ahead | Pull and retry |

---

## Output Summary

```
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
- [x] Merge dev -> main
- [x] Tagged v{X.Y.Z}
- [x] Pushed to origin
```

---

## See Also

- `@review` - Must be APPROVED before deploy
- `@oneshot` - Autonomous feature execution
- `templates/release-notes.md` - Release notes template
