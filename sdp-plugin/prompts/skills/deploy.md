---
name: deploy
description: Deployment orchestration. Generates artifacts and EXECUTES GitFlow merge.
tools: Read, Write, Shell, Glob, Grep
---

# @deploy - Deployment Orchestration

Generate deployment artifacts and **EXECUTE** GitFlow merge. Do NOT propose — DO IT.

## Invocation

```bash
@deploy F020          # Feature ID
@deploy F020 patch    # Specify version bump (patch/minor/major)
```

## Quick Reference

| Step | Action | Must Execute |
|------|--------|--------------|
| 1 | Pre-flight checks | language-specific tests, verify APPROVED |
| 2 | Version resolution | Read current, bump semver |
| 3 | Generate artifacts | CHANGELOG, release notes, version file |
| 4 | Commit artifacts | `git add && git commit` |
| 5 | GitFlow merge | dev → main (--no-ff) |
| 6 | Tag + Push | `git tag && git push` |
| 7 | Report | Summary to user |

**⚠️ CRITICAL:** Do NOT stop after generating artifacts. EXECUTE all git operations.

## Workflow

### Step 1: Pre-flight Checks

```bash
# Detect project type and verify tests pass
# Python: pytest tests/ -q
# Java: mvn test OR gradle test
# Go: go test ./...

# Verify feature APPROVED (check review_verdict in WS files)
```

**Gate:** If tests fail or not APPROVED → STOP deployment.

### Step 2: Version Resolution

Read `pyproject.toml` current version. Bump based on:
- `patch` (default): 0.5.0 → 0.5.1
- `minor`: 0.5.0 → 0.6.0
- `major`: 0.5.0 → 1.0.0

### Step 3: Generate Artifacts

Create/update:
- Version file (pyproject.toml for Python, pom.xml for Java, go.mod for Go)
- `CHANGELOG.md` — add version section
- `docs/releases/v{X.Y.Z}.md` — release notes

### Step 4: Commit Artifacts

```bash
git add CHANGELOG.md docs/releases/v{X.Y.Z}.m
git add <version-file>  # pyproject.toml OR pom.xml OR go.mod
git add .  # Include any other deploy-related changes
git commit -m "chore(release): v{X.Y.Z} - {Feature} {Title}"
```

**⚠️ DO NOT ask user permission. Commit immediately.**

### Step 5: GitFlow Merge

```bash
# Push dev first
git push origin dev

# Switch to main, pull, merge
git checkout main
git pull origin main
git merge dev --no-ff -m "Merge dev: v{X.Y.Z} {Feature}"
```

**⚠️ DO NOT ask user permission. Execute merge.**

### Step 6: Tag + Push

```bash
# Create annotated tag
git tag -a v{X.Y.Z} -m "Release v{X.Y.Z}: {Feature} - {Title}"

# Push main + tag
git push origin main
git push origin v{X.Y.Z}

# Return to dev
git checkout dev
```

**⚠️ DO NOT ask user permission. Tag and push immediately.**

### Step 7: Report

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

## See Also

- [@review skill](../review/SKILL.md) — Must be APPROVED before deploy
- [Release Notes Template](../../templates/release-notes.md)
