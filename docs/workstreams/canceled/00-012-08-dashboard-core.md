---
ws_id: 00-012-08
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

## 00-012-08: Dashboard Core (Reusable UI Components)

### 🎯 Goal

**What must WORK after this WS is complete:**
- `DashboardState` dataclass holds workstreams, test results, agent activity
- `StateBus` pub/sub for state updates (subscribe/publish)
- `WorkstreamReader` scans and parses workstream YAML files
- `TestRunner` watches files, runs pytest, parses output (watchdog)
- `AgentReader` reads from daemon queue (optional, graceful degradation)
- Reusable Textual widgets: `WorkstreamTree`, `TestPanel`, `ActivityLog`
- All widgets update via StateBus subscription
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `DashboardState` dataclass with workstreams, test_results, agent_activity
- [ ] AC2: `StateBus` with subscribe() and publish() methods
- [ ] AC3: `WorkstreamReader` scans docs/workstreams/{backlog,in_progress,completed}/
- [ ] AC4: `TestRunner` uses watchdog for file watching, runs pytest
- [ ] AC5: `AgentReader` connects to daemon or returns None (graceful degradation)
- [ ] AC6: `WorkstreamTree` Textual widget shows tree by status
- [ ] AC7: `TestPanel` Textual widget shows test results + coverage
- [ ] AC8: `ActivityLog` Textual widget shows scrolling event log
- [ ] AC9: All widgets update via StateBus subscription
- [ ] AC10: Coverage ≥ 80%
- [ ] AC11: mypy --strict passes

---

### Context

Dashboard Core provides reusable UI components and state management for both Developer Dashboard (00-012-14) and Agent Monitor (existing monitoring). This is the foundation for all DX improvements.

---

### Dependencies

None (foundation for other dashboard workstreams)

---

### Steps

1. Create `src/sdp/dashboard/` directory structure
2. Create `src/sdp/dashboard/state.py` with DashboardState and StateBus
3. Create `src/sdp/dashboard/sources/workstream_reader.py`
4. Create `src/sdp/dashboard/sources/test_runner.py`
5. Create `src/sdp/dashboard/sources/agent_reader.py`
6. Create `src/sdp/dashboard/widgets/workstream_tree.py`
7. Create `src/sdp/dashboard/widgets/test_panel.py`
8. Create `src/sdp/dashboard/widgets/activity_log.py`
9. Modify `pyproject.toml` to add `textual` and `watchdog` dependencies
10. Create unit tests for all components
11. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/dashboard/
├── __init__.py
├── state.py              # DashboardState, StateBus (~100 LOC)
├── sources/
│   ├── __init__.py
│   ├── workstream_reader.py   # ~150 LOC
│   ├── test_runner.py         # ~200 LOC
│   └── agent_reader.py        # ~100 LOC
└── widgets/
    ├── __init__.py
    ├── workstream_tree.py     # ~200 LOC
    ├── test_panel.py          # ~150 LOC
    └── activity_log.py        # ~100 LOC

tests/unit/dashboard/
├── test_state.py         # ~50 LOC
├── test_workstream_reader.py   # ~100 LOC
├── test_test_runner.py   # ~100 LOC
└── test_widgets.py       # ~150 LOC
```

**Modified Files:**
- `pyproject.toml` - Add textual, watchdog dependencies (~10 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/dashboard/ -v
pytest --cov=sdp.dashboard --cov-fail-under=80

# Type checking
mypy src/sdp/dashboard/ --strict

# Manual smoke test (python -c)
python -c "
from sdp.dashboard.state import StateBus, DashboardState
from sdp.dashboard.sources.workstream_reader import WorkstreamReader
from sdp.dashboard.widgets.workstream_tree import WorkstreamTree
print('Dashboard Core imports OK')
"
```

---

### Constraints

- NOT implementing full TUI app yet (that's 00-012-14)
- NOT integrating with GitHub yet
- USING Textual for widgets (standard Python TUI framework)
- USING watchdog for file watching
- FOLLOWING existing patterns from `cli.py` for data structures
- ERROR HANDLING: Each source returns cached state on error, don't crash

---

### Scope Estimate

- **Files:** 16 created/modified
- **Lines:** ~1,500 LOC
- **Size:** MEDIUM (500-1500 LOC)
