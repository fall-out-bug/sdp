---
ws_id: 00-034-02
feature: F034
status: completed
complexity: MEDIUM
project_id: "00"
depends_on:
  - 00-034-01
---

# Workstream: Split Large Files (Phase 2: Beads/Unified)

**ID:** 00-034-02  
**Feature:** F034 (A+ Quality Initiative)  
**Status:** READY  
**Owner:** AI Agent  
**Complexity:** MEDIUM (~700 LOC refactoring)

---

## Goal

Разбить оставшиеся 15 файлов >200 LOC в `src/sdp/beads/`, `src/sdp/unified/`, `src/sdp/github/` на модули <200 LOC.

---

## Context

Phase 2 фокусируется на модулях интеграции и multi-agent:

**Файлы для рефакторинга:**

| File | LOC | Action |
|------|-----|--------|
| `unified/checkpoint/schema.py` | 373 | Split → models, serializers |
| `beads/skills_oneshot.py` | 369 | Split → executor, checkpoint_handler |
| `beads/sync.py` | 344 | Split → scanner, syncer |
| `beads/skills_design.py` | 312 | Split → analyzer, generator |
| `unified/orchestrator/orchestrator.py` | 298 | Split → coordinator, scheduler |
| `github/issue_sync.py` | 278 | Split → fetcher, updater |
| `github/project_board_sync.py` | 265 | Split → board_api, sync_logic |
| `beads/skills_build.py` | 256 | Split → tdd_runner, reporter |
| `unified/gates/quality_gate.py` | 245 | Split → checkers, aggregator |
| `beads/client.py` | 238 | Split → http_client, cache |
| `github/client.py` | 232 | Split → rest_client, graphql_client |
| `unified/feature/feature_executor.py` | 228 | Split → planner, runner |
| `beads/scope_manager.py` | 221 | Split → analyzer, validator |
| `unified/agent/spawner.py` | 218 | Split → factory, registry |
| `github/sync_service.py` | 205 | Split → issue_sync, pr_sync |

---

## Scope

### In Scope
- ✅ Split 15 files listed above
- ✅ Update all imports
- ✅ Maintain backward compatibility via re-exports

### Out of Scope
- ❌ Files already split in Phase 1
- ❌ New functionality
- ❌ API changes

---

## Dependencies

**Depends On:**
- [x] 00-034-01: Split Large Files (Phase 1: Core)

**Blocks:**
- 00-034-03: Increase Test Coverage

---

## Acceptance Criteria

- [ ] All 15 files split to <200 LOC each
- [ ] `find src/sdp -name "*.py" -exec wc -l {} + | awk '$1 > 200'` returns empty
- [ ] All tests pass: `pytest tests/ -x`
- [ ] mypy --strict passes
- [ ] Backward compatibility maintained

---

## Implementation Plan

### Task 1: Split `unified/checkpoint/schema.py` (373 LOC)

**Target structure:**
```
unified/checkpoint/
├── __init__.py
├── models.py           # Checkpoint, CheckpointState dataclasses (~100 LOC)
├── serializers.py      # to_json(), from_json() (~80 LOC)
├── storage.py          # save_checkpoint(), load_checkpoint() (~100 LOC)
└── schema.py           # DEPRECATED: re-exports (~20 LOC)
```

### Task 2: Split `beads/skills_oneshot.py` (369 LOC)

**Target structure:**
```
beads/skills/
├── __init__.py
├── oneshot/
│   ├── __init__.py
│   ├── executor.py         # Main execution logic (~150 LOC)
│   ├── checkpoint.py       # Checkpoint handling (~100 LOC)
│   └── progress.py         # Progress tracking (~80 LOC)
```

### Task 3: Split `beads/sync.py` (344 LOC)

**Target structure:**
```
beads/
├── sync/
│   ├── __init__.py
│   ├── scanner.py          # Find beads files (~100 LOC)
│   ├── syncer.py           # Sync logic (~120 LOC)
│   └── diff.py             # Diff calculation (~80 LOC)
```

### Task 4-15: Similar splits for remaining files

Apply same pattern: identify logical groupings, extract to submodules, maintain re-exports.

---

## DO / DON'T

### Refactoring

**✅ DO:**
- Follow same patterns as Phase 1
- Keep `beads/` isolated from `core/` (prepare for domain extraction)
- Group related functionality together

**❌ DON'T:**
- Introduce new dependencies between modules
- Change Beads CLI interface
- Modify notification APIs

---

## Files to Modify/Create

**Create:** ~45 new files (3 per split average)

**Modify:**
- [ ] All 15 original files → re-exports only
- [ ] Import statements across codebase

---

## Test Plan

### Unit Tests
- [ ] `tests/unit/test_beads*.py` — all pass
- [ ] `tests/unified/` — all pass
- [ ] `tests/unit/test_github*.py` — all pass

### Integration Tests
- [ ] `bd ready` command works
- [ ] `bd sync` command works
- [ ] Telegram notifications work

---

**Version:** 1.0  
**Created:** 2026-01-31
