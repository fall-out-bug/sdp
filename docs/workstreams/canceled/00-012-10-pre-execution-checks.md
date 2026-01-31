---
ws_id: 00-012-10
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

## 00-012-10: Pre-Execution Checks

### 🎯 Goal

**What must WORK after this WS is complete:**
- `PreExecutionChecker` class with `check(ws_id: str)` method
- Check 1: WS file exists and is valid YAML
- Check 2: Dependencies are satisfied (all deps completed)
- Check 3: No circular dependencies
- Check 4: WS size ≤ MEDIUM (< 1500 LOC)
- Checks run before `AgentExecutor.execute()`
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `PreExecutionChecker` class with `check(ws_id: str)` method
- [ ] AC2: Check 1: WS file exists and is valid YAML
- [ ] AC3: Check 2: Dependencies are satisfied (all deps completed)
- [ ] AC4: Check 3: No circular dependencies
- [ ] AC5: Check 4: WS size ≤ MEDIUM (< 1500 LOC)
- [ ] AC6: Checks run before `AgentExecutor.execute()`
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

Pre-execution validation prevents bad WS from reaching agents. This catches issues early and provides clear error messages.

---

### Dependencies

00-012-04 (Agent Executor Interface)

---

### Steps

1. Create `src/sdp/agents/pre_check.py` with pre-execution validation
2. Modify `src/sdp/agents/executor.py` to add pre-check before execution
3. Create unit tests for pre-check
4. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/agents/
└── pre_check.py      (~200 LOC)

tests/unit/agents/
└── test_pre_check.py  (~150 LOC)
```

**Modified Files:**
- `src/sdp/agents/executor.py` - Add pre-check before execution (~50 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/agents/test_pre_check.py -v
pytest --cov=sdp.agents.pre_check --cov-fail-under=80

# Type checking
mypy src/sdp/agents/pre_check.py --strict

# Manual test
sdp agent check WS-001-01
```

---

### Constraints

- USING `ws_parser.parse_ws_file()` from `ws_parser.py`
- USING `dependency_graph` from WS-00-012-06
- FOLLOWING existing validation patterns from `validators/`

---

### Scope Estimate

- **Files:** 3 created/modified
- **Lines:** ~450 LOC
- **Size:** SMALL (< 500 LOC)
