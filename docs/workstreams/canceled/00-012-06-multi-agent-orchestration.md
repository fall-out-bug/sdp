---
ws_id: 00-012-06
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

## 00-012-06: Multi-Agent Orchestration

### 🎯 Goal

**What must WORK after this WS is complete:**
- `Orchestrator` class manages agent pool (max 3 concurrent)
- Dependency resolution: execute WS in topological order
- Load balancing: assign WS to least busy agent
- Deadlock detection: circular dependency detection
- Orchestrator state persisted to `.sdp/orchestrator_state.json`
- `sdp orchestrator run --feature F012` executes all WS for feature
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `Orchestrator` class manages agent pool (max 3 concurrent)
- [ ] AC2: Dependency resolution: execute WS in topological order
- [ ] AC3: Load balancing: assign WS to least busy agent
- [ ] AC4: Deadlock detection: circular dependency detection
- [ ] AC5: Orchestrator state persisted to `.sdp/orchestrator_state.json`
- [ ] AC6: `sdp orchestrator run --feature F012` executes all WS for feature
- [ ] AC7: Coverage ≥ 80%
- [ ] AC8: mypy --strict passes

---

### Context

Multiple agents need to work concurrently without conflicts. Orchestrator manages agent pool, resolves dependencies, and balances load.

---

### Dependencies

00-012-04 (Agent Executor Interface)

---

### Steps

1. Create `src/sdp/agents/dependency_graph.py` for topological sort
2. Create `src/sdp/agents/agent_pool.py` for agent pool management
3. Create `src/sdp/agents/orchestrator.py` for main orchestration logic
4. Add `@main.group() def orchestrator()` to `src/sdp/cli.py` with `run` subcommand
5. Create unit tests for orchestrator and dependency graph
6. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/agents/
├── orchestrator.py        (~400 LOC)
├── dependency_graph.py    (~250 LOC)
└── agent_pool.py          (~200 LOC)

tests/unit/agents/
├── test_orchestrator.py   (~200 LOC)
└── test_dependency_graph.py (~50 LOC)
```

**Modified Files:**
- `src/sdp/cli.py` - Add orchestrator command group (~50 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/agents/test_orchestrator.py -v
pytest --cov=sdp.agents.orchestrator --cov-fail-under=80

# Type checking
mypy src/sdp/agents/ --strict

# Integration test
sdp orchestrator run --feature F012 --dry-run
```

---

### Constraints

- NOT implementing TUI yet
- USING `AgentExecutor` from WS-00-012-04
- USING `TaskQueue` from WS-00-012-02
- PARSING WS dependencies from `ws_parser.py` (`WSMetadata.dependencies`)
- FOLLOWING patterns from existing `orchestrator.md` agent spec

---

### Scope Estimate

- **Files:** 7 created/modified
- **Lines:** ~1100 LOC
- **Size:** MEDIUM (500-1500 LOC)
