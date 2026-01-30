---
ws_id: 00-012-14
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

## 00-012-14: Developer Dashboard App

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp dashboard` command launches Textual TUI app
- Tab-based layout: Workstreams, Tests, Activity tabs
- Hotkeys: `1`/`w`=workstreams, `2`/`t`=tests, `3`/`a`=activity, `q`=quit, `r`=refresh
- Workstreams tab: tree view by status, filter by feature/project
- Tests tab: shows test results + coverage, runs watch mode updates
- Activity tab: scrolling log of events (git hooks, daemon, WS changes)
- App uses all Dashboard Core components (00-012-08)
- Live updates via StateBus (workstreams poll every 1s, tests on file change)
- Graceful degradation when daemon not running
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp dashboard` command launches Textual TUI app
- [ ] AC2: Tab-based layout: Workstreams, Tests, Activity tabs
- [ ] AC3: Hotkeys: `1`/`w`=workstreams, `2`/`t`=tests, `3`/`a`=activity, `q`=quit, `r`=refresh
- [ ] AC4: Workstreams tab: tree view by status, filter by feature/project
- [ ] AC5: Tests tab: shows test results + coverage, live updates
- [ ] AC6: Activity tab: scrolling log of events
- [ ] AC7: App uses all Dashboard Core components
- [ ] AC8: Live updates via StateBus
- [ ] AC9: Graceful degradation when daemon not running
- [ ] AC10: Coverage ≥ 80%
- [ ] AC11: mypy --strict passes

---

### Context

Developer Dashboard is the main TUI interface that brings together all DX improvements: workstream status, test watch mode, and activity monitoring. It's the "control center" for SDP development.

---

### Dependencies

00-012-08 (Dashboard Core), 00-012-11 (Status Command), 00-012-12 (Test Watch Mode), 00-012-13 (Auto-State Management)

---

### Steps

1. Create `src/sdp/dashboard/dashboard_app.py` with main Textual app
2. Create `src/sdp/dashboard/tabs/` directory
3. Create `src/sdp/dashboard/tabs/workstreams_tab.py`
4. Create `src/sdp/dashboard/tabs/tests_tab.py`
5. Create `src/sdp/dashboard/tabs/activity_tab.py`
6. Add `@main.command() def dashboard()` to `src/sdp/cli.py`
7. Create unit tests for app logic (using Textual's testing helpers)
8. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/dashboard/
├── dashboard_app.py        # ~250 LOC (main app, tab container)
└── tabs/
    ├── __init__.py
    ├── workstreams_tab.py  # ~150 LOC (tree view, filters)
    ├── tests_tab.py        # ~100 LOC (test panel, watch integration)
    └── activity_tab.py     # ~100 LOC (log viewer)

tests/unit/dashboard/
└── test_dashboard_app.py   # ~150 LOC
```

**Modified Files:**
- `src/sdp/cli.py` - Add dashboard command (~20 LOC)

**UI Mockup:**
```
┌─────────────────────────────────────────────────────────────┐
│  SDP Dashboard                                      [F012]  │
├─────────────────────────────────────────────────────────────┤
│                                                                  │
│  [Workstreams] [Tests] [Activity]                               │
│  ───────────────────────────────────────────────────────────  │
│                                                                  │
│  ┌─ Workstreams ─────────────────────────────────────────┐    │
│  │ 📁 Ideas (2)                                            │    │
│  │   ├─ idea-user-auth [draft]                           │    │
│  │   └─ idea-github-agent [needs_review]                 │    │
│  │                                                        │    │
│  │ 📁 Backlog (11)                                         │    │
│  │   ├─ 00-012-11: Workstream Status [SMALL]             │    │
│  │   ├─ 00-012-12: Test Watch Mode [MEDIUM]              │    │
│  │   └─ 00-012-13: Auto-State Mgmt [MEDIUM]              │    │
│  │                                                        │    │
│  │ 📁 In Progress (1)                                     │    │
│  │   └─ 00-012-14: Developer Dashboard [assignee: @user] │    │
│  │                                                        │    │
│  │ 📁 Completed (40)                                       │    │
│  │   └─ 00-012-08: Dashboard Core ✅                      │    │
│  └────────────────────────────────────────────────────────┘    │
│                                                                  │
│  Press: [w]orkstreams [t]ests [a]ctivity [q]uit                │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│  SDP Dashboard                                      [F012]  │
├─────────────────────────────────────────────────────────────┤
│                                                                  │
│  [Workstreams] [Tests] [Activity]                               │
│  ───────────────────────────────────────────────────────────  │
│                                                                  │
│  ┌─ Tests ─────────────────────────────────────────────────┐   │
│  │ 📊 Test Status                                            │   │
│  │ ✅ PASSED: 142  |  ❌ FAILED: 3  |  ⏭️ SKIPPED: 2        │   │
│  │ 📊 Coverage: 87% ━━━━━━━━━━━━━━━━━━━━                    │   │
│  │ 🕐 Last run: 2 seconds ago                               │   │
│  │                                                          │   │
│  │ Failed Tests:                                            │   │
│  │ ❌ test_dashboard_app.py::test_tab_switch                │   │
│  │    AssertionError: Expected 'tests', got 'workstreams'   │   │
│  │ ❌ test_workstreams_tab.py::test_filter                  │   │
│  │    ValueError: Invalid feature ID                       │   │
│  │ ❌ test_activity_tab.py::test_log_scrolling              │   │
│  │    Timeout: Log did not update                           │   │
│  │                                                          │   │
│  │ [Run All Tests]  [Run Failed Only]                       │   │
│  └────────────────────────────────────────────────────────┘   │
│                                                                  │
│  Watching for file changes...                                   │
│  Press: [w]orkstreams [t]ests [a]ctivity [q]uit                │
└─────────────────────────────────────────────────────────────┘
```

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/dashboard/test_dashboard_app.py -v
pytest --cov=sdp.dashboard.dashboard_app --cov-fail-under=80

# Type checking
mypy src/sdp/dashboard/dashboard_app.py --strict

# Manual test
sdp dashboard
# Expected: TUI launches, shows workstreams tab
# Press 't': switch to tests tab
# Press 'a': switch to activity tab
# Press 'q': exit

# Test live updates (terminal 2)
echo "def test_x(): assert False" >> tests/unit/dashboard/test_state.py
# Expected: tests tab updates within 2s showing failed test
```

---

### Constraints

- USING all widgets from 00-012-08 (WorkstreamTree, TestPanel, ActivityLog)
- USING `WorkstreamReader`, `TestRunner`, `AgentReader` from 00-012-08
- USING Textual's `App` class with tab navigation
- STARTING WorkstreamReader poll task (every 1s)
- STARTING TestRunner watch task (file changes)
- NOT implementing daemon integration (AgentReader returns None if not running)
- FOLLOWING Textual best practices (async/await, proper cleanup)
- ERROR HANDLING: Widget errors show in-place, don't crash app

---

### Scope Estimate

- **Files:** 7 created/modified
- **Lines:** ~770 LOC
- **Size:** MEDIUM (500-1500 LOC)
