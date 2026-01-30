---
ws_id: 00-012-07
project_id: 00
feature: F012
status: superseded
size: MEDIUM
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
---

## 00-012-07: GitHub Project Fields Integration

### 🎯 Goal

**What must WORK after this WS is complete:**
- Sync WS status to GitHub Project "Status" field (single select)
- Sync WS size to GitHub Project "Size" field (text)
- Sync WS feature to GitHub Project "Feature" field (text)
- Auto-create custom fields if missing
- Field mapping configurable via `.sdp/github_fields.toml`
- Bidirectional sync (GitHub changes → WS frontmatter)
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: Sync WS status to GitHub Project "Status" field (single select)
- [ ] AC2: Sync WS size to GitHub Project "Size" field (text)
- [ ] AC3: Sync WS feature to GitHub Project "Feature" field (text)
- [ ] AC4: Auto-create custom fields if missing
- [ ] AC5: Field mapping configurable via `.sdp/github_fields.toml`
- [ ] AC6: Bidirectional sync (GitHub changes → WS frontmatter)
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

GitHub Project Fields (not labels) provide richer status tracking. This sync extends existing project board sync with custom field support.

---

### Dependencies

00-012-03 (Enhanced GitHub Sync)

---

### Steps

1. Create `src/sdp/github/fields_client.py` with custom fields API wrapper
2. Create `src/sdp/github/fields_config.py` for config loading
3. Create `src/sdp/github/fields_sync.py` for field sync logic
4. Modify `src/sdp/github/sync_service.py` to integrate fields sync
5. Modify `src/sdp/github/projects_client.py` to add field creation methods
6. Create unit tests for fields sync
7. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/github/
├── fields_client.py   (~200 LOC)
├── fields_sync.py     (~250 LOC)
└── fields_config.py   (~100 LOC)

tests/unit/github/
└── test_fields_sync.py (~100 LOC)
```

**Modified Files:**
- `src/sdp/github/sync_service.py` - Integrate fields sync (~50 LOC)
- `src/sdp/github/projects_client.py` - Add field creation methods (~100 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/github/test_fields_sync.py -v
pytest --cov=sdp.github.fields_sync --cov-fail-under=80

# Type checking
mypy src/sdp/github/ --strict

# Manual test
sdp-github sync --fields --all
```

---

### Constraints

- EXTENDING `ProjectBoardSync` from `project_board_sync.py`
- EXTENDING `ProjectsClient` from `projects_client.py`
- REUSING `StatusMapper` patterns from `status_mapper.py`
- FOLLOWING existing GitHub integration patterns

---

### Scope Estimate

- **Files:** 6 created/modified
- **Lines:** ~800 LOC
- **Size:** MEDIUM (500-1500 LOC)
