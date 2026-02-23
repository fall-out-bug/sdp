---
description: Deployment orchestration from approved feature to release handoff.
agent: builder
---

# /deploy — Deploy Feature

When calling `/deploy {feature} [version_bump]`:

1. Load skill: `.claude/skills/deploy/SKILL.md`
2. Pre-flight: pytest, verify APPROVED
3. Version: bump semver (patch/minor/major)
4. Generate: CHANGELOG, release notes, pyproject.toml
5. **EXECUTE** (do NOT propose):
   - `git commit` artifacts
   - `git merge dev → main`
   - `git tag v{X.Y.Z}`
   - `git push origin main v{X.Y.Z}`
6. Report summary

## Quick Reference

**Input:** APPROVED feature + version bump (default: patch)
**Output:** Production deployment + v{X.Y.Z} tag
**Rule:** Do NOT stop after artifacts — EXECUTE all git operations

## Version Bump

- `@deploy F020` — patch (0.5.0 → 0.5.1)
- `@deploy F020 minor` — minor (0.5.0 → 0.6.0)
- `@deploy F020 major` — major (0.5.0 → 1.0.0)
