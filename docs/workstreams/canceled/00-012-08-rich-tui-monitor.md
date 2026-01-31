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

## 00-012-08: Rich TUI Monitor

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp monitor` command launches TUI
- TUI shows: active agents, task queue, completed tasks, errors
- Live updates every 1s (refresh from queue state)
- Color-coded status (green=active, red=error, yellow=queued)
- Keyboard controls: `q`=quit, `r`=refresh, `s`=sort
- TUI handles terminal resize gracefully
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp monitor` command launches TUI
- [ ] AC2: TUI shows: active agents, task queue, completed tasks, errors
- [ ] AC3: Live updates every 1s (refresh from queue state)
- [ ] AC4: Color-coded status (green=active, red=error, yellow=queued)
- [ ] AC5: Keyboard controls: `q`=quit, `r`=refresh, `s`=sort
- [ ] AC6: TUI handles terminal resize gracefully
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

Users need real-time visibility into agent execution. TUI provides kanban-style monitoring with live updates.

---

### Dependencies

00-012-06 (Multi-Agent Orchestration)

---

### Steps

1. Create `src/sdp/monitor/__init__.py` with exports
2. Create `src/sdp/monitor/state_reader.py` to read queue/orchestrator state
3. Create `src/sdp/monitor/widgets.py` with UI components (progress bars, tables)
4. Create `src/sdp/monitor/tui.py` with TUI main loop
5. Modify `pyproject.toml` to add `rich` dependency for TUI
6. Add `@main.command() def monitor()` to `src/sdp/cli.py`
7. Create unit tests for TUI
8. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/monitor/
├── __init__.py
├── tui.py            (~400 LOC)
├── widgets.py        (~300 LOC)
└── state_reader.py   (~150 LOC)

tests/unit/monitor/
└── test_tui.py       (~150 LOC)
```

**Modified Files:**
- `pyproject.toml` - Add rich dependency (~5 LOC)
- `src/sdp/cli.py` - Add monitor command (~10 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/monitor/test_tui.py -v
pytest --cov=sdp.monitor --cov-fail-under=80

# Type checking
mypy src/sdp/monitor/ --strict

# Manual test
sdp monitor
```

---

### Constraints

- READING state from `TaskQueue` (WS-00-012-02)
- READING state from `Orchestrator` (WS-00-012-06)
- USING `rich` library for TUI (standard Python TUI library)
- FOLLOWING patterns from existing CLI commands

---

### Scope Estimate

- **Files:** 7 created/modified
- **Lines:** ~1000 LOC
- **Size:** MEDIUM (500-1500 LOC)
