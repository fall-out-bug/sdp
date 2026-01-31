---
ws_id: 00-012-01
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

## 00-012-01: Daemon Service Framework

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp daemon --watch` command starts daemon process
- Daemon creates PID file at `.sdp/daemon.pid` on startup
- Daemon responds to SIGTERM/SIGINT with graceful shutdown
- Daemon logs to `.sdp/daemon.log` with timestamps
- `--stop` flag kills daemon by PID file
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp daemon --watch` command starts daemon process
- [ ] AC2: Daemon creates PID file at `.sdp/daemon.pid` on startup
- [ ] AC3: Daemon responds to SIGTERM/SIGINT with graceful shutdown
- [ ] AC4: Daemon logs to `.sdp/daemon.log` with timestamps
- [ ] AC5: `--stop` flag kills daemon by PID file
- [ ] AC6: Coverage ≥ 80%
- [ ] AC7: mypy --strict passes

---

### Context

SDP needs a daemon framework that can run as a background service, watch for file changes, and trigger sync operations. This is the foundation for all other daemon features.

---

### Dependencies

None (first WS in F012)

---

### Steps

1. Create `src/sdp/daemon/__init__.py` with exports
2. Create `src/sdp/daemon/daemon.py` with daemon lifecycle management
3. Create `src/sdp/daemon/pid_manager.py` for PID file handling
4. Add `@main.group() def daemon()` to `src/sdp/cli.py` with `--watch`, `--stop`, `--status` flags
5. Create unit tests for daemon and PID manager
6. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/daemon/
├── __init__.py
├── daemon.py        (~150 LOC)
└── pid_manager.py   (~100 LOC)

tests/unit/daemon/
├── __init__.py
├── test_daemon.py    (~80 LOC)
└── test_pid_manager.py  (~70 LOC)
```

**Modified Files:**
- `src/sdp/cli.py` - Add daemon command group (~50 LOC)

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/daemon/ -v
pytest --cov=sdp.daemon --cov-fail-under=80

# Type checking
mypy src/sdp/daemon/ --strict

# Manual smoke test
sdp daemon --watch &
sleep 2
cat .sdp/daemon.pid
sdp daemon --stop

# Verify daemon stopped
! ps aux | grep "sdp daemon"
```

---

### Constraints

- NOT integrating with GitHub yet (just daemon framework)
- NOT implementing task watching (future WS)
- NOT implementing webhooks (future WS)
- FOLLOWING existing Click patterns from `cli.py`
- FOLLOWING logging patterns from `sync_service.py`

---

### Scope Estimate

- **Files:** 7 created/modified
- **Lines:** ~450 LOC
- **Size:** SMALL (< 500 LOC)
