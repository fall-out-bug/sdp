---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-17
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: При смене статуса в Beads, local state обновляется
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: При локальном изменении, Beads синхронизируется
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: '`sdp sync` показывает расхождения'
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: ''
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: Conflict resolution с выбором источника
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-18
---

## 00-032-18: Beads Status Sync

### 🎯 Goal

**What must WORK after completing this WS:**
- Автоматическая синхронизация статусов между Beads и local state
- Conflict detection и resolution

**Acceptance Criteria:**
- [ ] AC1: При смене статуса в Beads, local state обновляется
- [ ] AC2: При локальном изменении, Beads синхронизируется
- [ ] AC3: `sdp sync` показывает расхождения
- [ ] AC4: Conflict resolution с выбором источника

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Local state и Beads могут рассинхронизироваться.

**Solution**: Sync service с conflict detection.

### Dependencies

- **00-032-17**: Beads Scope Management

### Steps

1. **Create SyncService**

   ```python
   # src/sdp/beads/sync_service.py
   from dataclasses import dataclass
   from enum import Enum
   
   from sdp.beads import BeadsClient, BeadsStatus
   from sdp.guard.state import StateManager, GuardState
   
   class SyncSource(Enum):
       LOCAL = "local"
       BEADS = "beads"
   
   @dataclass
   class SyncConflict:
       ws_id: str
       local_status: str | None
       beads_status: str
       field: str
   
   @dataclass
   class SyncResult:
       synced: bool
       conflicts: list[SyncConflict]
       changes: list[str]
   
   class BeadsSyncService:
       """Sync local state with Beads."""
       
       def __init__(self, client: BeadsClient):
           self._client = client
       
       def check_sync(self) -> SyncResult:
           """Check if local state matches Beads."""
           local_state = StateManager.load()
           conflicts = []
           
           if not local_state.active_ws:
               return SyncResult(synced=True, conflicts=[], changes=[])
           
           beads_task = self._client.get_task(local_state.active_ws)
           if not beads_task:
               conflicts.append(SyncConflict(
                   ws_id=local_state.active_ws,
                   local_status="active",
                   beads_status="not_found",
                   field="existence"
               ))
               return SyncResult(synced=False, conflicts=conflicts, changes=[])
           
           # Check status match
           if beads_task.status != BeadsStatus.IN_PROGRESS:
               conflicts.append(SyncConflict(
                   ws_id=local_state.active_ws,
                   local_status="active",
                   beads_status=beads_task.status.value,
                   field="status"
               ))
           
           return SyncResult(
               synced=len(conflicts) == 0,
               conflicts=conflicts,
               changes=[]
           )
       
       def sync(self, source: SyncSource = SyncSource.BEADS) -> SyncResult:
           """Sync state from specified source."""
           check = self.check_sync()
           
           if check.synced:
               return check
           
           changes = []
           
           for conflict in check.conflicts:
               if source == SyncSource.BEADS:
                   # Update local from Beads
                   if conflict.beads_status == "not_found":
                       StateManager.clear()
                       changes.append(f"Cleared local (WS {conflict.ws_id} not in Beads)")
                   elif conflict.beads_status != "in_progress":
                       StateManager.clear()
                       changes.append(f"Cleared local (Beads status: {conflict.beads_status})")
               else:
                   # Update Beads from local
                   self._client.update_task_status(
                       conflict.ws_id,
                       BeadsStatus.IN_PROGRESS
                   )
                   changes.append(f"Updated Beads to IN_PROGRESS")
           
           return SyncResult(synced=True, conflicts=[], changes=changes)
   ```

2. **Add CLI command**

   ```python
   # src/sdp/cli/sync.py
   import typer
   from sdp.beads import create_beads_client
   from sdp.beads.sync_service import BeadsSyncService, SyncSource
   
   app = typer.Typer(help="Beads sync commands")
   
   @app.command("check")
   def check_sync() -> None:
       """Check sync status."""
       service = BeadsSyncService(create_beads_client())
       result = service.check_sync()
       
       if result.synced:
           typer.echo("✅ Local and Beads are in sync")
           return
       
       typer.echo("❌ Sync conflicts detected:")
       for c in result.conflicts:
           typer.echo(f"  - {c.ws_id}: local={c.local_status}, beads={c.beads_status}")
   
   @app.command("run")
   def run_sync(
       source: str = typer.Option("beads", help="Source of truth: beads or local")
   ) -> None:
       """Sync local state with Beads."""
       sync_source = SyncSource(source)
       service = BeadsSyncService(create_beads_client())
       result = service.sync(sync_source)
       
       if result.changes:
           for change in result.changes:
               typer.echo(f"  ✅ {change}")
       else:
           typer.echo("✅ Already in sync")
   ```

### Output Files

- `src/sdp/beads/sync_service.py`
- `src/sdp/cli/sync.py`
- `tests/unit/test_sync_service.py`

### Completion Criteria

```bash
# CLI works
sdp sync check
sdp sync run --source beads

# Tests pass
pytest tests/unit/test_sync_service.py -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______
