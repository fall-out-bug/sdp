# F067 Repository Hardening - Migration Guide

**Date:** 2026-02-13
**Feature:** F067 Repository Hardening
**Status:** Complete

---

## Overview

F067 hardens the SDP repository with consistent tooling, quality gates, and documentation. This guide explains what changed and how to migrate.

---

## What Changed

### 1. Go Toolchain Alignment (WS-04)

**Before:** Mixed Go versions (1.25.6, 1.26, 1.21, 1.24)
**After:** Single Go 1.24 across all surfaces

| Location | Before | After |
|----------|--------|-------|
| `go.mod` (root) | 1.25.6 | 1.24 |
| `sdp-plugin/go.mod` | 1.26 | 1.24 |
| `.github/workflows/*.yml` | 1.26 | 1.24 |
| `sdp-plugin/Dockerfile` | 1.21-alpine | 1.24-alpine |
| NEW: `.go-version` | N/A | 1.24 |

**Migration:**
```bash
# Install correct Go version
brew install go@1.24  # or download from go.dev

# Verify
cat .go-version
go version
```

### 2. Quality Gate Alignment (WS-05)

**Before:** CI enforced 60% coverage, docs said 80%
**After:** 80% everywhere, guard checks are blocking

| Gate | Before | After |
|------|--------|-------|
| CI coverage threshold | 60% | 80% |
| Guard check in CI | Non-blocking (`|| true`) | Blocking |
| Config source | 3 files | 1 file (`.sdp/guard-rules.yml`) |

**Removed Files:**
- `quality-gate.toml`
- `ci-gates.toml`

**Migration:**
```bash
# Local coverage check
cd sdp-plugin && go test -cover ./...

# If below 80%, add tests before pushing
```

### 3. Prompt Source of Truth (WS-02)

**Before:** Prompts duplicated in `prompts/` and `sdp-plugin/prompts/`
**After:** Single canonical source in `prompts/`

**Removed:**
- `sdp-plugin/prompts/` (entire directory)

**Symlinks:**
- `.claude/skills` → `prompts/skills/`
- `.claude/agents` → `prompts/agents/`

**Migration:**
```bash
# Always edit canonical files
vim prompts/skills/build/SKILL.md  # CORRECT
# NOT: vim sdp-plugin/prompts/skills/build/SKILL.md  # REMOVED

# Check for drift
./hooks/check-prompt-drift.sh
```

### 4. Repository Hygiene (WS-07)

**Before:** Tracked build artifacts (`coverage_quality.out`)
**After:** Clean source tree, evidence policy documented

**Removed from tracking:**
- `coverage_quality.out`

**Added to .gitignore:**
- `.sdp/memory.db`
- `.sdp/checkpoints/`
- `bin/`, `dist/`

**Migration:**
```bash
# Clean local artifacts
git pull
rm -f coverage_quality.out  # if exists locally
```

### 5. Git Structural Hygiene (WS-11)

**Before:** Self-referential submodule, Python artifacts
**After:** Clean structure, Go-first

**Removed:**
- `.sdp/.sdp` (self-referential submodule)
- `poetry.lock`

**Fixed:**
- `.gitmodules` now empty (no submodules)
- `.cursor/worktrees.json` uses `go mod download` instead of `poetry install`

### 6. Adapter Consistency (WS-03)

**Before:** OpenCode referenced `prompts/commands/` (non-existent)
**After:** All adapters reference `prompts/skills/*/SKILL.md`

**Updated:**
- `.opencode/opencode.json` - Fixed all skill paths
- `.cursor/worktrees.json` - Uses Go instead of Poetry

### 7. Installation (WS-06)

**New:** Install script and improved docs

```bash
# New curl install
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/scripts/install.sh | bash

# Correct go install path
go install github.com/fall-out-bug/sdp/sdp-plugin/cmd/sdp@latest
```

---

## Breaking Changes

### For Contributors

1. **Go version:** Must use 1.24 (enforced in CI)
2. **Coverage:** Must be ≥80% (enforced in CI)
3. **Prompt edits:** Only edit `prompts/` directory
4. **Guard checks:** Now block CI on error

### For CI

1. **Guard check:** No longer has `|| true` - will fail on violations
2. **Coverage gate:** Threshold is 80%, not 60%
3. **Adapter validation:** CI checks all adapter paths exist

---

## Rollback Instructions

### If Go version causes issues:
```bash
# Temporarily revert in go.mod
# But note: CI will fail until fixed
```

### If coverage threshold too strict:
1. Update `.sdp/guard-rules.yml` to lower threshold
2. Update `.github/workflows/go-ci.yml` to match
3. Run config consistency check in CI

### If guard blocking needed:
```bash
# Local only - add || true temporarily
./sdp guard check --staged || true

# DO NOT commit the || true
```

---

## Contributor Checklist

Before submitting PR:

- [ ] Go version is 1.24 (`go version`)
- [ ] Tests pass (`cd sdp-plugin && go test ./...`)
- [ ] Coverage ≥80% (`go test -cover ./...`)
- [ ] Guard checks pass (`./sdp guard check --staged`)
- [ ] Prompt edits in `prompts/` only
- [ ] No `.out` files staged
- [ ] Version strings consistent

---

## Common Failure Modes

| Failure | Cause | Fix |
|---------|-------|-----|
| Go version mismatch | Wrong Go installed | Install Go 1.24 |
| Coverage < 80% | Missing tests | Add tests |
| Guard check fails | LOC/type/TODO violations | Fix violations |
| Adapter path error | Wrong skill path | Use `prompts/skills/*/SKILL.md` |
| Duplicate prompt | Created outside `prompts/` | Move to canonical location |

---

## Files Changed Summary

| Category | Files |
|----------|-------|
| **Added** | `.go-version`, `DEVELOPMENT.md`, `scripts/install.sh`, `prompts/agents/VERSIONS.md`, `docs/decisions/ADR-001-dual-module-structure.md`, `hooks/check-prompt-drift.sh` |
| **Removed** | `poetry.lock`, `quality-gate.toml`, `ci-gates.toml`, `coverage_quality.out`, `sdp-plugin/prompts/`, `.sdp/.sdp` |
| **Modified** | `.github/workflows/go-ci.yml`, `.gitignore`, `.gitmodules`, `CONTRIBUTING.md`, `README.md`, `docs/PROTOCOL.md`, `docs/reference/quality-gates.md`, `docs/reference/EVIDENCE-COVERAGE.md`, `.opencode/opencode.json`, `.cursor/worktrees.json` |

---

**Version:** 1.0.0
**Last Updated:** 2026-02-13
