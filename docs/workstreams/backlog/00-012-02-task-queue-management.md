---
ws_id: 00-012-02
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

## 00-012-02: Task Queue Management

### 🎯 Goal

**What must WORK after this WS is complete:**
- `TaskQueue` class with `enqueue()`, `dequeue()`, `peek()` methods
- Priority-based scheduling (backlog < active < blocked)
- Retry logic with exponential backoff (max 2 attempts)
- Thread-safe queue operations (Lock-based)
- Queue state persisted to `.sdp/queue_state.json`
- CLI command `sdp queue status` shows pending tasks
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `TaskQueue` class with `enqueue()`, `dequeue()`, `peek()` methods
- [ ] AC2: Priority-based scheduling (backlog < active < blocked)
- [ ] AC3: Retry logic with exponential backoff (max 2 attempts)
- [ ] AC4: Thread-safe queue operations (Lock-based)
- [ ] AC5: Queue state persisted to `.sdp/queue_state.json`
- [ ] AC6: CLI command `sdp queue status` shows pending tasks
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

Agents need a task queue to manage async workstream execution with priority scheduling and retry logic. This queue must be thread-safe for concurrent agent access.

---

### Dependencies

00-012-01 (Daemon Service Framework)

---

### Steps

1. Create `src/sdp/queue/__init__.py` with exports
2. Create `src/sdp/queue/task.py` with Task dataclass (ws_id, priority, retries)
3. Create `src/sdp/queue/priority.py` with Priority enum and comparison logic
4. Create `src/sdp/queue/task_queue.py` with queue data structure
5. Create `src/sdp/queue/state.py` for queue state persistence
6. Add `@main.group() def queue()` to `src/sdp/cli.py` with `status` subcommand
7. Create unit tests for all queue modules
8. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/queue/
├── __init__.py
├── task_queue.py     (~300 LOC)
├── task.py           (~100 LOC)
├── priority.py       (~80 LOC)
└── state.py          (~120 LOC)

tests/unit/queue/
├── __init__.py
├── test_task_queue.py  (~150 LOC)
└── test_priority.py    (~100 LOC)
```

**Modified Files:**
- `src/sdp/cli.py` - Add queue command group (~50 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/queue/ -v
pytest --cov=sdp.queue --cov-fail-under=80

# Type checking
mypy src/sdp/queue/ --strict

# Integration test
sdp queue status  # Should show empty queue
```

---

### Constraints

- NOT integrating with GitHub yet
- NOT implementing agent execution yet
- NOT implementing TUI yet
- REUSING `WSMetadata` from `ws_parser.py` for task payloads
- FOLLOWING existing patterns in `core/workstream.py` for data structures

---

### Scope Estimate

- **Files:** 10 created/modified
- **Lines:** ~900 LOC
- **Size:** MEDIUM (500-1500 LOC)
