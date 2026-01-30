---
ws_id: 00-012-13
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

## 00-012-13: Auto-State Management

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp ws move <ws-id> --to <status>` command moves WS files between directories
- Auto-updates `status:` field in YAML frontmatter when moving
- Auto-updates `docs/workstreams/INDEX.md` after move
- `--start` flag moves backlog → in_progress, sets `started:` timestamp
- `--complete` flag moves in_progress → completed, sets `completed:` timestamp
- `@build` skill auto-moves files (backlog → in_progress → completed)
- Validates move (can't move to completed if AC not met)
- Uses `WorkstreamReader` from Dashboard Core (00-012-08)
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp ws move <ws-id> --to <status>` command moves WS files
- [ ] AC2: Auto-updates `status:` field in YAML frontmatter
- [ ] AC3: Auto-updates `docs/workstreams/INDEX.md` after move
- [ ] AC4: `--start` flag moves backlog → in_progress, sets timestamp
- [ ] AC5: `--complete` flag moves in_progress → completed, sets timestamp
- [ ] AC6: `@build` skill auto-moves files (via integration)
- [ ] AC7: Validates move (can't complete if AC not 100% met)
- [ ] AC8: Uses `WorkstreamReader` from Dashboard Core
- [ ] AC9: Coverage ≥ 80%
- [ ] AC10: mypy --strict passes

---

### Context

Currently, workstream files must be manually moved between backlog/in_progress/completed directories, and INDEX.md must be manually updated. This is error-prone and creates friction. Auto-state management eliminates this overhead.

---

### Dependencies

00-012-08 (Dashboard Core)

---

### Steps

1. Create `src/sdp/workspace/` directory
2. Create `src/sdp/workspace/mover.py` for file move logic
3. Create `src/sdp/workspace/index_updater.py` for INDEX.md updates
4. Create `src/sdp/workspace/validator.py` for move validation (AC check)
5. Create `src/sdp/workspace/state_updater.py` for YAML frontmatter updates
6. Add `@main.group() def ws()` to `src/sdp/cli.py` with `move` subcommand
7. Create unit tests for all components
8. Update `@build` skill to call auto-move (separate task)
9. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/workspace/
├── __init__.py
├── mover.py           # ~200 LOC (file moves, directory management)
├── index_updater.py   # ~150 LOC (INDEX.md parsing/updating)
├── validator.py       # ~100 LOC (move validation, AC check)
└── state_updater.py   # ~150 LOC (YAML frontmatter updates)

tests/unit/workspace/
├── test_mover.py      # ~150 LOC
├── test_index_updater.py  # ~100 LOC
├── test_validator.py  # ~100 LOC
└── test_state_updater.py  # ~80 LOC
```

**Modified Files:**
- `src/sdp/cli.py` - Add ws command group (~40 LOC)

**Example Usage:**
```bash
# Start workstream
sdp ws move 00-012-08 --start
# Moves: backlog/00-012-08.md → in_progress/00-012-08.md
# Updates: status: in_progress, started: 2026-01-26
# Updates: INDEX.md

# Complete workstream
sdp ws move 00-012-08 --complete
# Validates: All AC checked? Yes → proceed
# Moves: in_progress/00-012-08.md → completed/00-012-08.md
# Updates: status: completed, completed: 2026-01-26
# Updates: INDEX.md

# Direct move (manual control)
sdp ws move 00-012-08 --to backlog
# Moves file, updates status, updates INDEX
```

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/workspace/ -v
pytest --cov=sdp.workspace --cov-fail-under=80

# Type checking
mypy src/sdp/workspace/ --strict

# Manual tests
sdp ws move 00-012-08 --start
sdp ws move 00-012-08 --complete
sdp ws move 00-012-08 --to backlog

# Verify INDEX.md updated
grep "00-012-08" docs/workstreams/INDEX.md
```

---

### Constraints

- USING `WorkstreamReader` from 00-012-08
- PARSING YAML frontmatter with `yaml.safe_load()`
- PRESERVING YAML ordering and comments when possible
- VALIDATING completion: all checkbox AC must be checked (`- [x]`)
- NOT moving files if validation fails (show clear error)
- FOLLOWING existing file structure patterns
- INDEX.md format: must match existing structure

---

### Scope Estimate

- **Files:** 10 created/modified
- **Lines:** ~1,220 LOC
- **Size:** MEDIUM (500-1500 LOC)
