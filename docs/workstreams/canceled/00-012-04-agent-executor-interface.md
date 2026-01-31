---
ws_id: 00-012-04
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

## 00-012-04: Agent Executor Interface

### 🎯 Goal

**What must WORK after this WS is complete:**
- `AgentExecutor` class with `execute(ws_id: str)` method
- Executor calls `/build` skill internally (subprocess or direct API)
- Progress tracking via `TaskUpdate` (real-time status)
- Error handling: retry on failure, escalate after 2 attempts
- Execution metrics stored in `.sdp/execution_metrics.json`
- Timeout protection (max 1h per WS)
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `AgentExecutor` class with `execute(ws_id: str)` method
- [ ] AC2: Executor calls `/build` skill internally (subprocess or direct API)
- [ ] AC3: Progress tracking via `TaskUpdate` (real-time status)
- [ ] AC4: Error handling: retry on failure, escalate after 2 attempts
- [ ] AC5: Execution metrics stored in `.sdp/execution_metrics.json`
- [ ] AC6: Timeout protection (max 1h per WS)
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

Agents need an executor interface for running workstreams autonomously with progress tracking and error handling. This bridges the queue with agent execution.

---

### Dependencies

00-012-02 (Task Queue Management)

---

### Steps

1. Create `src/sdp/agents/__init__.py` with exports
2. Create `src/sdp/agents/executor.py` with main executor logic
3. Create `src/sdp/agents/metrics.py` for execution metrics tracking
4. Create `src/sdp/agents/errors.py` with custom error classes
5. Create unit tests for executor and metrics
6. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/agents/
├── __init__.py
├── executor.py       (~350 LOC)
├── metrics.py         (~200 LOC)
└── errors.py          (~100 LOC)

tests/unit/agents/
├── __init__.py
├── test_executor.py   (~200 LOC)
└── test_metrics.py    (~100 LOC)
```

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/agents/ -v
pytest --cov=sdp.agents --cov-fail-under=80

# Type checking
mypy src/sdp/agents/ --strict

# Integration test (requires mock /build)
sdp agent execute WS-001-01 --dry-run
```

---

### Constraints

- NOT implementing TUI yet
- NOT implementing multi-agent orchestration yet
- USING `TaskQueue` from WS-00-012-02 for task management
- CALLING existing `/build` skill (`.claude/skills/build/SKILL.md`)
- FOLLOWING patterns from `builder_router.py` for model routing
- REUSING `WSMetadata` from `ws_parser.py`

---

### Scope Estimate

- **Files:** 8 created
- **Lines:** ~950 LOC
- **Size:** MEDIUM (500-1500 LOC)
