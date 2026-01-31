---
ws_id: 00-012-11
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

## 00-012-11: Workstream Status Command

### 🎯 Goal

**What must WORK after this WS is complete:**
- `sdp status` command shows all workstreams grouped by status
- Output format: table with WS-ID, Title, Status, Assignee, Feature
- `--filter` flag filters by status (backlog|in_progress|completed)
- `--feature` flag filters by feature ID
- `--watch` flag enables live refresh (every 2s)
- Color-coded output (yellow=backlog, blue=in_progress, green=completed)
- Uses `WorkstreamReader` from Dashboard Core (00-012-08)
- Coverage ≥ 80%
- mypy --strict passes

**Acceptance Criteria:**
- [ ] AC1: `sdp status` command shows all workstreams grouped by status
- [ ] AC2: Output format: table with WS-ID, Title, Status, Assignee, Feature
- [ ] AC3: `--filter` flag filters by status
- [ ] AC4: `--feature` flag filters by feature ID
- [ ] AC5: `--watch` flag enables live refresh (every 2s)
- [ ] AC6: Color-coded output (yellow=backlog, blue=in_progress, green=completed)
- [ ] AC7: Uses `WorkstreamReader` from Dashboard Core
- [ ] AC8: Coverage ≥ 80%
- [ ] AC9: mypy --strict passes

---

### Context

Developers need a quick way to see all workstreams and their status without navigating directories. `sdp status` provides a unified view with filtering options.

---

### Dependencies

00-012-08 (Dashboard Core)

---

### Steps

1. Create `src/sdp/status/` directory
2. Create `src/sdp/status/formatter.py` for table formatting with colors
3. Create `src/sdp/status/command.py` for status command logic
4. Add `@main.command() def status()` to `src/sdp/cli.py`
5. Create unit tests for formatter and command
6. Run TDD cycle: Red → Green → Refactor

---

### Expected Result

**Created Files:**
```
src/sdp/status/
├── __init__.py
├── formatter.py       # ~150 LOC (table formatting, colors)
└── command.py         # ~100 LOC (status logic, filters, watch)

tests/unit/status/
├── test_formatter.py  # ~80 LOC
└── test_command.py    # ~100 LOC
```

**Modified Files:**
- `src/sdp/cli.py` - Add status command (~30 LOC)

**Example Output:**
```
$ sdp status

Workstreams by Status

📁 Backlog (10)
┌─────────────┬──────────────────────────┬────────┬──────────┬─────────┐
│ WS-ID       │ Title                     │ Status │ Assignee │ Feature │
├─────────────┼──────────────────────────┼────────┼──────────┼─────────┤
│ 00-012-01   │ Daemon Service Framework │ backlog │ -        │ F012    │
│ 00-012-02   │ Task Queue Management    │ backlog │ -        │ F012    │
└─────────────┴──────────────────────────┴────────┴──────────┴─────────┘

📁 In Progress (2)
┌─────────────┬──────────────────┬─────────────┬──────────┬─────────┐
│ WS-ID       │ Title             │ Status      │ Assignee │ Feature │
├─────────────┼──────────────────┼─────────────┼──────────┼─────────┤
│ 00-011-06   │ PRD Command       │ in_progress │ @user    │ F011    │
└─────────────┴──────────────────┴─────────────┴──────────┴─────────┘

📁 Completed (40)
...

$ sdp status --filter backlog --feature F012
📁 Backlog (10, F012 only)
...
```

---

### Completion Criteria

```bash
# Unit tests
pytest tests/unit/status/ -v
pytest --cov=sdp.status --cov-fail-under=80

# Type checking
mypy src/sdp/status/ --strict

# Manual tests
sdp status
sdp status --filter backlog
sdp status --feature F012
sdp status --watch  # Refresh every 2s, Ctrl+C to exit
```

---

### Constraints

- USING `WorkstreamReader` from 00-012-08
- USING `rich` library for table formatting and colors (already in SDP)
- NOT implementing TUI (that's 00-012-14)
- FOLLOWING existing CLI patterns from `cli.py`

---

### Scope Estimate

- **Files:** 5 created/modified
- **Lines:** ~460 LOC
- **Size:** SMALL (< 500 LOC)
