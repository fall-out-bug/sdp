---
ws_id: 00-032-29
title: Add AC→Test Mappings for F032 Workstreams
status: completed
feature: F032
project_id: 00
size: M
estimated_loc: 200
dependencies: []
completed: "2026-01-30"
review_source: docs/reports/2026-01-30-F032-review.md
---

# 00-032-29: Add AC→Test Mappings for F032 Workstreams

## Goal

Achieve 100% AC→Test traceability for all F032 workstreams using `sdp trace` tooling.

## Context

F032 review found 0% AC mapped. The `sdp trace check` command now works (after Issue 002 bugfix), but no mappings exist.

## Acceptance Criteria

- [x] AC1: Run `sdp trace auto --apply` on all F032 workstreams
- [x] AC2: Manual mappings added where auto-detect fails
- [x] AC3: `sdp trace check 00-032-XX` passes for all 28 workstreams
- [x] AC4: Traceability report shows 100% coverage

## Technical Approach

1. **Auto-detect first:**
   ```bash
   for ws in docs/workstreams/completed/00-032-*.md; do
     sdp trace auto "$ws" --apply
   done
   ```

2. **Review unmapped ACs** — add test docstrings or explicit mappings

3. **Verify:**
   ```bash
   sdp trace check 00-032-01
   # ... repeat for all
   ```

## Out of Scope

- Writing new tests (see 00-032-30)
- Fixing failing tests (see Issue 003)

## Execution Report

**Date:** 2026-01-30

**Changes:**
1. Extended `TraceabilityService.add_mapping` with markdown frontmatter fallback (when Beads has no task)
2. Ran `sdp trace auto --apply` on all 28 F032 workstreams
3. Manual mappings for 4 unmapped AC6s:
   - 00-032-06 AC6 → test_validate_skill_valid
   - 00-032-25 AC6 → test_validate_skill_too_long
   - 00-032-26 AC6 → test_sync_updates_beads_from_local
   - 00-032-28 AC6 → test_extract_body_with_frontmatter

**Result:** 28/28 workstreams pass `sdp trace check` (100% traceability)
