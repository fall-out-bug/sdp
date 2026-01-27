# ISSUE-001: WS-00-012-03 Missing Implementation

**Severity:** 🔴 HIGH
**Workstream:** 00-012-03 (Enhanced GitHub Sync)
**Status:** Open

## Problem

WS-00-012-03 specified creating:
- `src/sdp/github/conflict_resolver.py` (~200 LOC) - **MISSING**
- `src/sdp/github/sync_enhanced.py` (~150 LOC) - **MISSING**

## Acceptance Criteria

- [ ] AC1: `SyncService` detects conflicts (WS vs GitHub status mismatch)
- [ ] AC2: Conflict resolution: WS file wins (source of truth)
- [ ] AC3: GitHub status synced to WS frontmatter on conflict
- [ ] AC4: New `sync_backlog()` method for incremental sync
- [ ] AC5: `--dry-run` flag previews changes without applying
- [ ] AC6: Integration with existing `status_mapper.py` and `project_board_sync.py`
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

## Solution Options

1. **Create missing modules** per original spec
2. **Update WS spec** if enhanced sync not needed (simplify requirements)

## Decision

WS-00-012-03 depends on existing sync infrastructure. The `SyncService` already handles WS→GitHub sync. The "enhanced" bidirectional sync may be out of scope for F012 core goals.

**Recommended:** Mark this WS as simplified/deferred, document current sync capabilities.

## Steps to Fix

1. Update WS-00-012-03 spec to reflect actual implementation
2. Document that basic unidirectional sync already exists
3. Mark bidirectional sync as future work
