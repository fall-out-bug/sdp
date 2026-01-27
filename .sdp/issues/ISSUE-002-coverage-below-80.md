# ISSUE-002: Coverage Below 80%

**Severity:** 🟡 MEDIUM
**Status:** Open

## Problem

F012 overall coverage: 71% (target: ≥80%)

Modules below threshold:
- `orchestrator.py`: 66% (missing 34%)
- `daemon.py`: 58% (missing 22%)
- `dashboard_app.py`: 45% (missing 35%)
- `dashboard/widgets/`: 0-25% (missing 75-100%)
- `webhook/server.py`: 17% (missing 83%)
- `test_watch/runner.py`: 0% (missing 100%)
- `workspace/index_updater.py`: 35%
- `workspace/state_updater.py`: 35%

## Acceptance Criteria

- [ ] All F012 modules ≥ 80% coverage
- [ ] Overall F012 coverage ≥ 80%
- [ ] All tests pass

## Steps to Fix

1. **orchestrator.py** (66% → 80%):
   - Add tests for parallel execution paths
   - Add tests for error handling in concurrent execution
   - Add tests for state persistence

2. **daemon.py** (58% → 80%):
   - Add tests for daemon lifecycle (start/stop)
   - Add tests for signal handling (SIGTERM/SIGINT)
   - Add tests for PID file management

3. **dashboard** (0-45% → 80%):
   - Add widget tests (workstream_tree, test_panel, activity_log)
   - Add tab tests (workstreams_tab, tests_tab, activity_tab)
   - Add dashboard_app integration tests

4. **webhook/server.py** (17% → 80%):
   - Add tests for HTTP request handling
   - Add tests for signature validation
   - Add tests for error responses

5. **test_watch/runner.py** (0% → 80%):
   - Add tests for test execution
   - Add tests for output parsing
