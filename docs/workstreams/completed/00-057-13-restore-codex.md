---
ws_id: "00-057-13"
feature_id: "F057"
title: "Restore .codex/ directory (artifact lost in git clean)"
status: completed
priority: "P1"
depends_on: []
blocks: []
project_id: "00"
completed_at: "2026-02-10"
---

# 00-057-13: Restore .codex/ directory (artifact lost in git clean)

## Goal

Recreate the `.codex/` directory structure that was the deliverable of task sdp-5ch (00-192-03: Codex Adapter) and was deleted by `git clean -fd`.

## Context

Codex Adapter (completed) produced: `.codex/INSTALL.md`, `.codex/skills/`. That directory was untracked and was removed during repo cleanup. This WS restores the artifact only (no Python adapter — repo is Go now).

## Acceptance Criteria

- [x] AC1: `.codex/INSTALL.md` exists with setup instructions for Codex users
- [x] AC2: `.codex/skills/` exists (project-level skills; can reference or copy from `.claude/skills/`)
- [x] AC3: Structure matches Codex expectation: INSTALL.md readable by Codex, skills/ present
- [x] AC4: `.codex/` is committed (tracked) so it is not lost by future `git clean`

## Scope Files

**Implementation:**
- .codex/INSTALL.md (new)
- .codex/skills/ (new — copy or symlink from .claude/skills for project skills)

**Verification:**
- `test -f .codex/INSTALL.md && test -d .codex/skills`

## Execution Report

**Executed by:** @build (00-057-13)  
**Date:** 2026-02-10

### Goal Status

- [x] AC1: `.codex/INSTALL.md` — ✅ created
- [x] AC2: `.codex/skills/` + README.md — ✅ created
- [x] AC3: Structure matches Codex — ✅
- [x] AC4: `.codex/` committed — ✅

**Goal Achieved:** ✅ YES

### Files Changed

| File | Action |
|------|--------|
| .codex/INSTALL.md | created |
| .codex/skills/README.md | created |
| docs/workstreams/in_progress/.gitkeep | created |
| .contracts/.gitkeep | created |

### Commit

feat(00-057-13): Restore .codex/ and in_progress/.contracts/ (artifact lost in git clean)

## Notes

- Original spec: docs/workstreams/completed/00-004-03-codex-adapter.md
- Beads task sdp-5ch was closed; this WS restores the lost artifact.
