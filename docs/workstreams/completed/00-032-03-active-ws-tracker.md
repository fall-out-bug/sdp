---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-01
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '`@build` автоматически ставит WS в IN_PROGRESS через Beads'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: При завершении WS статус меняется на CLOSED
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Попытка активировать второй WS возвращает ошибку
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Unit тесты (≥80%)
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`sdp ws current` показывает текущий активный WS'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-03
---

## 00-032-03: Active WS Tracker

### 🎯 Goal

**What must WORK after completing this WS:**
- Автоматический трекинг статуса WS через Beads
- Только один WS может быть IN_PROGRESS
- Статус синхронизирован между local state и Beads

**Acceptance Criteria:**
- [ ] AC1: `@build` автоматически ставит WS в IN_PROGRESS через Beads
- [ ] AC2: При завершении WS статус меняется на CLOSED
- [ ] AC3: Попытка активировать второй WS возвращает ошибку
- [ ] AC4: `sdp ws current` показывает текущий активный WS
- [ ] AC5: Unit тесты (≥80%)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Guard state только локальный. Beads не знает какой WS активен.

**Solution**: WorkstreamTracker синхронизирует статус с Beads.

### Dependencies

- **00-032-01**: Guard Skill Foundation

### Input Files

- `src/sdp/guard/state.py` (will be created by WS-02, but can work independently)
- `src/sdp/beads/client.py`
- `src/sdp/beads/models.py` (BeadsStatus enum)

### Steps

1. **Create WorkstreamTracker**

   ```python
   # src/sdp/guard/tracker.py
   from datetime import datetime
   from pathlib import Path
   import json
   
   from sdp.beads import BeadsClient, BeadsStatus
   from sdp.errors import SDPError, ErrorCategory
   
   class WorkstreamInProgressError(SDPError):
       """Another WS is already in progress."""
       
       def __init__(self, current_ws: str, requested_ws: str):
           super().__init__(
               category=ErrorCategory.BUSINESS_LOGIC,
               message=f"WS {current_ws} already in progress",
               remediation=(
                   f"1. Complete current WS: sdp guard complete {current_ws}\n"
                   f"2. Or abort: sdp guard abort {current_ws}\n"
                   f"3. Then activate: sdp guard activate {requested_ws}"
               ),
               context={"current_ws": current_ws, "requested_ws": requested_ws}
           )
   
   class WorkstreamTracker:
       """Track active workstream with Beads sync."""
       
       def __init__(
           self, 
           client: BeadsClient,
           state_file: Path = Path(".sdp/state.json")
       ):
           self._client = client
           self._state_file = state_file
       
       def get_active(self) -> str | None:
           """Get currently active WS ID."""
           if not self._state_file.exists():
               return None
           with open(self._state_file) as f:
               state = json.load(f)
           return state.get("active_ws")
       
       def activate(self, ws_id: str) -> None:
           """Activate WS and update Beads status."""
           # Check no other WS is active
           current = self.get_active()
           if current and current != ws_id:
               raise WorkstreamInProgressError(current, ws_id)
           
           # Update Beads status
           self._client.update_task_status(ws_id, BeadsStatus.IN_PROGRESS)
           
           # Get scope from WS metadata
           ws = self._client.get_task(ws_id)
           scope = ws.sdp_metadata.get("scope_files", []) if ws else []
           
           # Save local state
           self._save_state({
               "active_ws": ws_id,
               "started_at": datetime.utcnow().isoformat(),
               "scope_files": scope,
           })
       
       def complete(self, ws_id: str) -> None:
           """Mark WS as complete."""
           current = self.get_active()
           if current != ws_id:
               raise ValueError(f"WS {ws_id} is not active (active: {current})")
           
           self._client.update_task_status(ws_id, BeadsStatus.CLOSED)
           self._clear_state()
       
       def abort(self, ws_id: str) -> None:
           """Abort WS without completing."""
           current = self.get_active()
           if current != ws_id:
               raise ValueError(f"WS {ws_id} is not active")
           
           # Return to OPEN status
           self._client.update_task_status(ws_id, BeadsStatus.OPEN)
           self._clear_state()
       
       def _save_state(self, state: dict) -> None:
           self._state_file.parent.mkdir(exist_ok=True)
           with open(self._state_file, "w") as f:
               json.dump(state, f, indent=2)
       
       def _clear_state(self) -> None:
           self._save_state({"active_ws": None})
   ```

2. **Add CLI command**

   ```python
   # Add to src/sdp/cli/ws.py or create new
   @app.command("current")
   def current_ws() -> None:
       """Show currently active workstream."""
       from sdp.guard.tracker import WorkstreamTracker
       from sdp.beads import create_beads_client
       
       tracker = WorkstreamTracker(create_beads_client())
       ws_id = tracker.get_active()
       
       if not ws_id:
           typer.echo("No active workstream")
           return
       
       typer.echo(f"Active WS: {ws_id}")
   ```

3. **Write tests**

   ```python
   # tests/unit/test_ws_tracker.py
   import pytest
   from unittest.mock import Mock, patch
   from pathlib import Path
   import tempfile
   
   from sdp.guard.tracker import WorkstreamTracker, WorkstreamInProgressError
   from sdp.beads.models import BeadsStatus
   
   class TestWorkstreamTracker:
       def test_activate_updates_beads_status(self):
           """AC1: Activate sets IN_PROGRESS in Beads."""
           client = Mock()
           client.get_task.return_value = Mock(sdp_metadata={})
           
           with tempfile.NamedTemporaryFile(suffix=".json", delete=False) as f:
               tracker = WorkstreamTracker(client, Path(f.name))
               tracker.activate("00-032-01")
           
           client.update_task_status.assert_called_with(
               "00-032-01", BeadsStatus.IN_PROGRESS
           )
       
       def test_complete_sets_closed_status(self):
           """AC2: Complete sets CLOSED in Beads."""
           client = Mock()
           client.get_task.return_value = Mock(sdp_metadata={})
           
           with tempfile.NamedTemporaryFile(suffix=".json", delete=False) as f:
               tracker = WorkstreamTracker(client, Path(f.name))
               tracker.activate("00-032-01")
               tracker.complete("00-032-01")
           
           client.update_task_status.assert_called_with(
               "00-032-01", BeadsStatus.CLOSED
           )
       
       def test_cannot_activate_second_ws(self):
           """AC3: Second activation raises error."""
           client = Mock()
           client.get_task.return_value = Mock(sdp_metadata={})
           
           with tempfile.NamedTemporaryFile(suffix=".json", delete=False) as f:
               tracker = WorkstreamTracker(client, Path(f.name))
               tracker.activate("00-032-01")
               
               with pytest.raises(WorkstreamInProgressError):
                   tracker.activate("00-032-02")
   ```

### Output Files

- `src/sdp/guard/tracker.py`
- `src/sdp/cli/ws.py` (current command)
- `tests/unit/test_ws_tracker.py`

### Completion Criteria

```bash
# Module imports
python -c "from sdp.guard.tracker import WorkstreamTracker"

# Tests pass
pytest tests/unit/test_ws_tracker.py -v

# Coverage
pytest tests/unit/test_ws_tracker.py --cov=src/sdp/guard/tracker --cov-fail-under=80
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
