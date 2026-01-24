---
ws_id: 00-190-05
project_id: 00
feature: F006
status: backlog
size: SMALL
github_issue: 1617
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-190-05: File Size Reduction - project_map.py

### 🎯 Goal

**What must WORK after this WS is complete:**
- `project_map.py` split into smaller, focused modules
- All functionality preserved
- No impact on tests or public API

**Acceptance Criteria:**
- [ ] AC1: `project_map.py` split into ≤3 files, each <300 LOC
- [ ] AC2: All 16 tests in `test_project_map.py` still pass
- [ ] AC3: Public API unchanged (`parse_project_map()`, `create_project_map_template()`, query methods)
- [ ] AC4: Test coverage remains ≥90%
- [ ] AC5: No ruff or mypy errors

---

### Context

**Current Issue:**
- `project_map.py` is 577 lines (significantly over 200 LOC guideline)
- Single file contains parsing, template generation, and query logic
- Violates AI-Readiness constraint

**Root Cause:**
- Multiple responsibilities in one file:
  1. Core dataclasses (ProjectMap, Decision, Constraint, TechStackItem)
  2. Parsing logic (_parse_decisions_table, _parse_constraints, _parse_tech_stack_table)
  3. Template generation (create_project_map_template)
  4. Query methods (get_decision, get_constraint)

**Impact:**
- 🟡 MEDIUM severity - File is too large for efficient AI agent processing
- Affects maintainability and test isolation

---

### Dependencies

00--03 (already completed)

---

### Proposed Solution

**Split into 3 modules:**

1. **`project_map_types.py`** (~100 LOC)
   - Dataclasses: `ProjectMap`, `Decision`, `Constraint`, `TechStackItem`
   - Exception: `ProjectMapParseError`
   - Query methods: `get_decision()`, `get_constraint()`

2. **`project_map_parser.py`** (~250 LOC)
   - Main parser: `parse_project_map()`
   - Helper parsers: `_parse_decisions_table()`, `_parse_constraints()`, `_parse_tech_stack_table()`
   - Section extractors: `_extract_section()`, `_extract_project_name()`

3. **`project_map_template.py`** (~150 LOC)
   - Template generator: `create_project_map_template()`
   - Template string constants

**Public API preservation:**
- All exports remain in `sdp/core/__init__.py`
- Importing from `sdp.core` unchanged
- No breaking changes

---

### Steps

1. Create `project_map_types.py` - Move dataclasses and query methods
2. Create `project_map_parser.py` - Move parsing logic
3. Create `project_map_template.py` - Move template generation
4. Update imports in existing files
5. Run tests: `pytest tests/unit/core/test_project_map.py -v`
6. Verify public API: `python -c "from sdp.core import parse_project_map, ProjectMap"`
7. Run linters: `ruff check src/sdp/core/` && `mypy src/sdp/core/`
8. Remove original `project_map.py`

---

### Completion Criteria

```bash
# All should pass:
pytest tests/unit/core/test_project_map.py -v --cov=src/sdp/core/project_map_types --cov=src/sdp/core/project_map_parser --cov=src/sdp/core/project_map_template --cov-report=term-missing
ruff check src/sdp/core/
mypy src/sdp/core/ --ignore-missing-imports
wc -l src/sdp/core/project_map_*.py  # All <300 LOC
```

---

### Constraints

- NO changes to test files
- NO breaking changes to public API
- NO functional changes (pure refactoring)
- MUST maintain ≥90% test coverage
