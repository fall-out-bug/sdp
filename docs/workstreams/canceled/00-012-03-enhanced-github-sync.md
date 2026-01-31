---
ws_id: 00-012-03
project_id: 00
feature: F012
status: superseded
superseded_by: 00-032-01
supersede_reason: "F012 daemon/agent framework superseded by F032 Guard + Beads integration"
size: MEDIUM
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
---

## 00-012-03: Enhanced GitHub Sync

### 🎯 Goal

**What must WORK after this WS is complete:**
- `SyncService` detects conflicts (WS vs GitHub status mismatch)
- Conflict resolution: WS file wins (source of truth)
- GitHub status synced to WS frontmatter on conflict
- New `sync_backlog()` method for incremental sync
- `--dry-run` flag previews changes without applying
- Integration with existing `status_mapper.py` and `project_board_sync.py`
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `SyncService` detects conflicts (WS vs GitHub status mismatch)
- [ ] AC2: Conflict resolution: WS file wins (source of truth)
- [ ] AC3: GitHub status synced to WS frontmatter on conflict
- [ ] AC4: New `sync_backlog()` method for incremental sync
- [ ] AC5: `--dry-run` flag previews changes without applying
- [ ] AC6: Integration with existing `status_mapper.py` and `project_board_sync.py`
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

Existing `SyncService` supports unidirectional sync (WS → GitHub). We need bidirectional sync with conflict detection and smart merge for GitHub Agent Orchestrator.

---

### Dependencies

None (can be done in parallel with WS-00-012-01)

---

### Steps

1. Create `src/sdp/github/conflict_resolver.py` with smart merge logic
2. Create `src/sdp/github/sync_enhanced.py` as enhanced sync wrapper
3. Modify `src/sdp/github/sync_service.py` to add `sync_backlog()` and `--dry-run` support
4. Modify `src/sdp/github/cli.py` to add `--dry-run` flag to sync commands
5. Create unit tests for conflict resolver and enhanced sync
6. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/github/
├── conflict_resolver.py  (~200 LOC)
└── sync_enhanced.py       (~150 LOC)

tests/unit/github/
├── test_conflict_resolver.py  (~150 LOC)
└── test_sync_enhanced.py      (~100 LOC)
```

**Modified Files:**
- `src/sdp/github/sync_service.py` - Add sync_backlog(), --dry-run support (~150 LOC)
- `src/sdp/github/cli.py` - Add --dry-run flag (~50 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/github/test_conflict_resolver.py -v
pytest --cov=sdp.github.sync_enhanced --cov-fail-under=80

# Type checking
mypy src/sdp/github/ --strict

# Manual test
sdp-github sync --dry-run --all
```

---

### Constraints

- NOT breaking existing sync functionality
- EXTENDING existing `SyncService` class (not rewriting)
- REUSING `StatusMapper.detect_conflict()` from `status_mapper.py`
- REUSING `ProjectBoardSync` from `project_board_sync.py`
- FOLLOWING existing patterns in `sync_service.py`

---

### Scope Estimate

- **Files:** 6 created/modified
- **Lines:** ~850 LOC
- **Size:** MEDIUM (500-1500 LOC)
