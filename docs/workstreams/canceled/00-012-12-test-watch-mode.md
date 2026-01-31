---
ws_id: 00-012-12
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

## 00-012-12: Test Watch Mode

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp test --watch` command starts file watcher on src/ and tests/
- File changes trigger pytest run on affected tests only
- Output shows: PASS/FAIL counts, coverage %, failing test names
- `--pattern` flag filters which tests to run (e.g., `--pattern test_unit*`)
- `--coverage` flag enables coverage report (default: on)
- Watch mode handles new test files (auto-discovery)
- Uses `TestRunner` from Dashboard Core (00-012-08)
- KeyboardInterrupt (Ctrl+C) exits gracefully
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp test --watch` command starts file watcher on src/ and tests/
- [ ] AC2: File changes trigger pytest run on affected tests only
- [ ] AC3: Output shows: PASS/FAIL counts, coverage %, failing test names
- [ ] AC4: `--pattern` flag filters which tests to run
- [ ] AC5: `--coverage` flag enables coverage report
- [ ] AC6: Watch mode handles new test files (auto-discovery)
- [ ] AC7: Uses `TestRunner` from Dashboard Core
- [ ] AC8: KeyboardInterrupt exits gracefully
- [ ] AC9: Coverage ≥ 80%
- [ ] AC10: mypy --strict passes

---

### Context

Fast feedback is critical for TDD. Watch mode provides sub-second test feedback when files change, eliminating manual test runs. This is a key DX improvement from the SDP analysis.

---

### Dependencies

00-012-08 (Dashboard Core)

---

### Steps

1. Create `src/sdp/test_watch/` directory
2. Create `src/sdp/test_watch/watcher.py` for file watching logic
3. Create `src/sdp/test_watch/runner.py` for test execution with output formatting
4. Create `src/sdp/test_watch/affected.py` for determining affected tests
5. Add `@main.command() def test()` to `src/sdp/cli.py` with `--watch`, `--pattern`, `--coverage` flags
6. Create unit tests for watcher and affected test logic
7. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/test_watch/
├── __init__.py
├── watcher.py         # ~150 LOC (watchdog integration, debounce)
├── runner.py          # ~200 LOC (pytest execution, output parsing)
└── affected.py        # ~100 LOC (file -> test mapping)

tests/unit/test_watch/
├── test_watcher.py    # ~100 LOC
└── test_affected.py   # ~80 LOC
```

**Modified Files:**
- `src/sdp/cli.py` - Add test command (~50 LOC)

**Example Output:**
```
$ sdp test --watch

🔍 Watching for changes...
[Press Ctrl+C to exit]

✓ src/sdp/dashboard/state.py changed
→ Running affected tests...

✅ PASSED 42 | ❌ FAILED 2 | ⏭️ SKIPPED 1 | 📊 Coverage 87%

Failed tests:
  - tests/unit/dashboard/test_state.py::test_state_bus_publish
  - tests/unit/dashboard/test_widgets.py::test_workstream_tree_update

Waiting for changes...
```

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/test_watch/ -v
pytest --cov=sdp.test_watch --cov-fail-under=80

# Type checking
mypy src/sdp/test_watch/ --strict

# Manual tests
# Terminal 1
sdp test --watch

# Terminal 2 (edit a test file)
echo "def test_new(): assert True" >> tests/unit/dashboard/test_state.py

# Expected: Terminal 1 shows test run within 2s
```

---

### Constraints

- USING `TestRunner` from 00-012-08
- USING `watchdog` for file watching (already added in 00-012-08)
- USING `pytest` with `--tb=short` for clean output
- NOT implementing TUI (console output only)
- DEBOUNCE file changes (wait 500ms after last change before running tests)
- FOLLOWING existing test patterns from `hooks/post-build.sh`

---

### Scope Estimate

- **Files:** 7 created/modified
- **Lines:** ~780 LOC
- **Size:** MEDIUM (500-1500 LOC)
