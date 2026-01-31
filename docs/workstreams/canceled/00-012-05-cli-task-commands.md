---
ws_id: 00-012-05
project_id: 00
feature: F012
status: superseded
superseded_by: 00-032-01
supersede_reason: "F012 daemon/agent framework superseded by F032 Guard + Beads integration"
size: SMALL
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
---

## 00-012-05: CLI Task Commands

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp task enqueue WS-XXX-YY` adds task to queue
- `sdp task execute WS-XXX-YY` runs task immediately
- `sdp task list` shows pending/running/completed tasks
- `sdp task cancel <task_id>` cancels pending task
- CLI follows existing Click patterns
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp task enqueue WS-XXX-YY` adds task to queue
- [ ] AC2: `sdp task execute WS-XXX-YY` runs task immediately
- [ ] AC3: `sdp task list` shows pending/running/completed tasks
- [ ] AC4: `sdp task cancel <task_id>` cancels pending task
- [ ] AC5: CLI follows existing Click patterns
- [ ] AC6: Coverage ≥ 80%
- [ ] AC7: mypy --strict passes

---

### Context

Users need CLI commands for manual task management. This provides the interface to the queue and executor created in previous WS.

---

### Dependencies

00-012-02 (Task Queue Management), 00-012-04 (Agent Executor Interface)

---

### Steps

1. Create `src/sdp/cli_tasks.py` with task CLI group
2. Modify `src/sdp/cli.py` to import and register task commands
3. Create unit tests for task CLI
4. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/
└── cli_tasks.py       (~250 LOC)

tests/unit/cli/
└── test_tasks.py      (~100 LOC)
```

**Modified Files:**
- `src/sdp/cli.py` - Import and register task commands (~20 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/cli/test_tasks.py -v
pytest --cov=sdp.cli_tasks --cov-fail-under=80

# Type checking
mypy src/sdp/cli_tasks.py --strict

# Manual test
sdp task enqueue WS-001-01
sdp task list
sdp task execute WS-001-01 --dry-run
sdp task cancel <task_id>
```

---

### Constraints

- NOT implementing TUI yet
- USING `TaskQueue` from WS-00-012-02
- USING `AgentExecutor` from WS-00-012-04
- FOLLOWING existing CLI patterns in `cli.py`

---

### Scope Estimate

- **Files:** 3 created/modified
- **Lines:** ~350 LOC
- **Size:** SMALL (< 500 LOC)
