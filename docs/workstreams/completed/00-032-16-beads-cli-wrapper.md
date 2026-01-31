---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-15
feature: F032
github_issue: null
project_id: 0
size: MEDIUM
status: completed
traceability:
- ac_description: '`RealBeadsClient` implements all `BeadsClient` methods'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`create_task`, `get_task`, `update_task_status` работают через
    CLI'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Error handling для CLI failures
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Unit тесты с mocked subprocess (≥80%)
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: JSON parsing из CLI output
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-16
---

## 00-032-16: Beads CLI Wrapper

### 🎯 Goal

**What must WORK after completing this WS:**
- `RealBeadsClient` реализует `BeadsClient` interface
- Subprocess calls к `bd` CLI
- Полная совместимость с MockBeadsClient

**Acceptance Criteria:**
- [ ] AC1: `RealBeadsClient` implements all `BeadsClient` methods
- [ ] AC2: `create_task`, `get_task`, `update_task_status` работают через CLI
- [ ] AC3: Error handling для CLI failures
- [ ] AC4: JSON parsing из CLI output
- [ ] AC5: Unit тесты с mocked subprocess (≥80%)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Только MockBeadsClient существует. Нет интеграции с real Beads.

**Solution**: Python wrapper для `bd` CLI.

### Dependencies

- **00-032-15**: Beads Go Install Check

### Steps

1. **Create RealBeadsClient**

   ```python
   # src/sdp/beads/real_client.py
   import subprocess
   import json
   import shutil
   from typing import Any
   
   from sdp.beads.base import BeadsClient
   from sdp.beads.exceptions import BeadsClientError, BeadsNotInstalledError
   from sdp.beads.models import (
       BeadsTask,
       BeadsTaskCreate,
       BeadsStatus,
       BeadsDependency,
   )
   
   class RealBeadsClient(BeadsClient):
       """Real Beads CLI client implementation."""
       
       def __init__(self, bd_path: str | None = None):
           self._bd = bd_path or shutil.which("bd")
           if not self._bd:
               raise BeadsNotInstalledError(
                   "Beads CLI (bd) not found. "
                   "Install: go install github.com/steveyegge/beads/cmd/bd@latest"
               )
       
       def _run(self, *args: str, check: bool = True) -> dict[str, Any]:
           """Run bd command and parse JSON output."""
           cmd = [self._bd, *args, "--json"]
           
           try:
               result = subprocess.run(
                   cmd,
                   capture_output=True,
                   text=True,
                   timeout=30
               )
           except subprocess.TimeoutExpired:
               raise BeadsClientError(f"Command timed out: {' '.join(cmd)}")
           except Exception as e:
               raise BeadsClientError(f"Failed to run bd: {e}")
           
           if check and result.returncode != 0:
               raise BeadsClientError(
                   f"bd command failed (exit {result.returncode}): {result.stderr}"
               )
           
           if not result.stdout.strip():
               return {}
           
           try:
               return json.loads(result.stdout)
           except json.JSONDecodeError as e:
               raise BeadsClientError(f"Invalid JSON from bd: {e}")
       
       def create_task(self, params: BeadsTaskCreate) -> BeadsTask:
           """Create a new task via CLI."""
           args = [
               "create",
               "--title", params.title,
           ]
           
           if params.description:
               args.extend(["--description", params.description])
           
           if params.priority:
               args.extend(["--priority", params.priority.value])
           
           if params.parent_id:
               args.extend(["--parent", params.parent_id])
           
           result = self._run(*args)
           return self._parse_task(result)
       
       def get_task(self, task_id: str) -> BeadsTask | None:
           """Get task by ID."""
           try:
               result = self._run("show", task_id)
               return self._parse_task(result)
           except BeadsClientError:
               return None
       
       def update_task_status(self, task_id: str, status: BeadsStatus) -> None:
           """Update task status."""
           status_cmd = {
               BeadsStatus.OPEN: "reopen",
               BeadsStatus.IN_PROGRESS: "start",
               BeadsStatus.BLOCKED: "block",
               BeadsStatus.CLOSED: "close",
           }
           
           cmd = status_cmd.get(status)
           if not cmd:
               raise ValueError(f"Unsupported status: {status}")
           
           self._run(cmd, task_id)
       
       def get_ready_tasks(self) -> list[str]:
           """Get tasks ready to work on."""
           result = self._run("ready")
           tasks = result.get("tasks", [])
           return [t["id"] for t in tasks]
       
       def add_dependency(
           self,
           from_id: str,
           to_id: str,
           dep_type: str = "blocks"
       ) -> None:
           """Add dependency between tasks."""
           self._run("dep", "add", from_id, to_id, "--type", dep_type)
       
       def list_tasks(
           self,
           status: BeadsStatus | None = None,
           parent_id: str | None = None
       ) -> list[BeadsTask]:
           """List tasks with optional filters."""
           args = ["list"]
           
           if status:
               args.extend(["--status", status.value])
           
           if parent_id:
               args.extend(["--parent", parent_id])
           
           result = self._run(*args)
           tasks = result.get("tasks", [])
           return [self._parse_task(t) for t in tasks]
       
       def _parse_task(self, data: dict[str, Any]) -> BeadsTask:
           """Parse task from JSON data."""
           return BeadsTask(
               id=data["id"],
               title=data["title"],
               description=data.get("description", ""),
               status=BeadsStatus(data.get("status", "open")),
               priority=data.get("priority"),
               parent_id=data.get("parent_id"),
               dependencies=[
                   BeadsDependency(
                       task_id=d["task_id"],
                       type=d.get("type", "blocks")
                   )
                   for d in data.get("dependencies", [])
               ],
               sdp_metadata=data.get("metadata", {}),
           )
   ```

2. **Update factory function**

   ```python
   # src/sdp/beads/__init__.py
   import os
   
   from sdp.beads.base import BeadsClient
   from sdp.beads.mock import MockBeadsClient
   from sdp.beads.real_client import RealBeadsClient
   
   def create_beads_client(use_mock: bool | None = None) -> BeadsClient:
       """Create appropriate Beads client.
       
       Args:
           use_mock: Force mock (True) or real (False). 
                    If None, uses BEADS_USE_MOCK env var.
       """
       if use_mock is None:
           use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"
       
       if use_mock:
           return MockBeadsClient()
       
       return RealBeadsClient()
   ```

3. **Write tests**

   ```python
   # tests/unit/test_real_beads_client.py
   import pytest
   from unittest.mock import patch, Mock
   
   from sdp.beads.real_client import RealBeadsClient
   from sdp.beads.models import BeadsTaskCreate, BeadsStatus
   from sdp.beads.exceptions import BeadsClientError
   
   class TestRealBeadsClient:
       @patch("subprocess.run")
       @patch("shutil.which", return_value="/usr/local/bin/bd")
       def test_create_task_calls_cli(self, mock_which, mock_run):
           """AC1: create_task calls bd create."""
           mock_run.return_value = Mock(
               returncode=0,
               stdout='{"id": "bd-0001", "title": "Test"}',
               stderr=""
           )
           
           client = RealBeadsClient()
           task = client.create_task(BeadsTaskCreate(
               title="Test Task",
               description="Description"
           ))
           
           assert task.id == "bd-0001"
           mock_run.assert_called_once()
           call_args = mock_run.call_args[0][0]
           assert "create" in call_args
           assert "--title" in call_args
       
       @patch("subprocess.run")
       @patch("shutil.which", return_value="/usr/local/bin/bd")
       def test_error_handling(self, mock_which, mock_run):
           """AC3: CLI errors raise BeadsClientError."""
           mock_run.return_value = Mock(
               returncode=1,
               stdout="",
               stderr="Error: task not found"
           )
           
           client = RealBeadsClient()
           
           with pytest.raises(BeadsClientError):
               client.get_task("nonexistent")
   ```

### Output Files

- `src/sdp/beads/real_client.py`
- `src/sdp/beads/__init__.py` (updated)
- `tests/unit/test_real_beads_client.py`

### Completion Criteria

```bash
# Module imports
python -c "from sdp.beads.real_client import RealBeadsClient"

# Tests pass
pytest tests/unit/test_real_beads_client.py -v

# Coverage
pytest tests/unit/test_real_beads_client.py --cov=src/sdp/beads/real_client --cov-fail-under=80
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
