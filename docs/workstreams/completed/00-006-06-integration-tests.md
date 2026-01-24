---
ws_id: 00-190-06
project_id: 00
feature: F006
status: backlog
size: SMALL
github_issue: 1620
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-190-06: Integration Tests for F190 Core

### 🎯 Goal

**What must WORK after this WS is complete:**
- Integration test suite validates end-to-end workflows
- Real file I/O tested (not just mocks)
- Multi-WS feature loading validated

**Acceptance Criteria:**
- [ ] AC1: Integration test file created: `tests/integration/test_core_integration.py`
- [ ] AC2: Test workflow: Parse WS → Load Feature → Calculate execution order
- [ ] AC3: Test workflow: Parse PROJECT_MAP.md from actual hw_checker docs
- [ ] AC4: Test workflow: CLI commands work end-to-end
- [ ] AC5: All integration tests pass

---

### Context

**Current Gap:**
- Only unit tests exist (61 tests total)
- No integration tests validating real file I/O
- No tests for multi-module workflows
- CLI commands not tested end-to-end

**Why This Matters:**
- Unit tests use mocks/fixtures (temporary files)
- Real-world usage involves actual WS files in repos
- Integration tests catch issues unit tests miss:
  - File encoding issues
  - Path resolution bugs
  - Multi-file dependency resolution
  - CLI argument parsing

---

### Dependencies

00--01, 00--02, 00--03, 00--04 (all completed)

---

### Test Scenarios

**Scenario 1: Parse Real WS Files**
- Parse actual 00--01.md from backlog/
- Validate frontmatter, goal, ACs parsed correctly
- Verify file_path populated

**Scenario 2: Load Feature from Directory**
- Load all WS-190-XX files from backlog/
- Validate 4 workstreams loaded
- Check execution order: 00--03, 00--01, 00--02, 00--04
- Verify dependency graph built correctly

**Scenario 3: Parse Real PROJECT_MAP.md**
- Parse `tools/hw_checker/docs/PROJECT_MAP.md`
- Validate decisions table loaded
- Query constraints by category
- Verify tech stack parsed

**Scenario 4: CLI End-to-End**
- Test: `sdp core parse-ws <real_file>`
- Test: `sdp core parse-project-map <real_file>`
- Validate output format (human-readable)
- Check exit codes (0 on success, 1 on error)

**Scenario 5: Error Handling**
- Test: Parse WS file with missing frontmatter
- Test: Parse WS file with circular dependencies
- Test: Load feature from empty directory
- Validate error messages are clear

---

### Steps

1. Create `tests/integration/test_core_integration.py`
2. Implement Scenario 1: Real WS parsing
3. Implement Scenario 2: Feature loading
4. Implement Scenario 3: PROJECT_MAP parsing
5. Implement Scenario 4: CLI end-to-end
6. Implement Scenario 5: Error handling
7. Run: `pytest tests/integration/test_core_integration.py -v`
8. Update pytest markers if needed

---

### Completion Criteria

```bash
# All should pass:
pytest tests/integration/test_core_integration.py -v
pytest tests/integration/test_core_integration.py --cov=src/sdp/core --cov-report=term-missing

# CLI tests:
sdp core parse-ws tools/hw_checker/docs/workstreams/backlog/00--01-core-workstream-parser.md
sdp core parse-project-map tools/hw_checker/docs/PROJECT_MAP.md
```

---

### Expected Test Count

- Scenario 1: 3 tests
- Scenario 2: 4 tests
- Scenario 3: 3 tests
- Scenario 4: 4 tests
- Scenario 5: 3 tests
- **Total: ~17 integration tests**

---

### Constraints

- MUST use real files from repo (not fixtures)
- MUST test CLI commands with subprocess/click.testing
- MUST validate exit codes
- NO mocks for file I/O (that's what unit tests do)
