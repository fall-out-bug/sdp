---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-16
feature: F032
github_issue: null
project_id: 0
size: MEDIUM
status: completed
traceability:
- ac_description: WS metadata содержит `scope_files` array
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`sdp ws scope <ws_id>` показывает scope файлов'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: '`sdp ws scope <ws_id> add <file>` добавляет файл'
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Guard использует scope из Beads metadata
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`sdp ws scope <ws_id> remove <file>` удаляет файл'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-17
---

## 00-032-17: Beads Scope Management

### 🎯 Goal

**What must WORK after completing this WS:**
- `@design` записывает `scope_files` в WS metadata
- `@guard` читает scope из Beads
- CLI для управления scope

**Acceptance Criteria:**
- [ ] AC1: WS metadata содержит `scope_files` array
- [ ] AC2: `sdp ws scope <ws_id>` показывает scope файлов
- [ ] AC3: `sdp ws scope <ws_id> add <file>` добавляет файл
- [ ] AC4: `sdp ws scope <ws_id> remove <file>` удаляет файл
- [ ] AC5: Guard использует scope из Beads metadata

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Scope хранится локально. Нет sync с Beads.

**Solution**: Хранить scope в Beads task metadata.

### Dependencies

- **00-032-16**: Beads CLI Wrapper

### Steps

1. **Create ScopeManager**

   ```python
   # src/sdp/beads/scope_manager.py
   from sdp.beads import BeadsClient
   
   class ScopeManager:
       """Manage workstream file scope in Beads metadata."""
       
       def __init__(self, client: BeadsClient):
           self._client = client
       
       def get_scope(self, ws_id: str) -> list[str]:
           """Get scope files for workstream."""
           task = self._client.get_task(ws_id)
           if not task:
               raise ValueError(f"WS not found: {ws_id}")
           
           return task.sdp_metadata.get("scope_files", [])
       
       def set_scope(self, ws_id: str, files: list[str]) -> None:
           """Set scope files for workstream."""
           task = self._client.get_task(ws_id)
           if not task:
               raise ValueError(f"WS not found: {ws_id}")
           
           metadata = task.sdp_metadata.copy()
           metadata["scope_files"] = files
           
           # Update via CLI
           self._client.update_metadata(ws_id, metadata)
       
       def add_file(self, ws_id: str, file_path: str) -> None:
           """Add file to scope."""
           scope = self.get_scope(ws_id)
           if file_path not in scope:
               scope.append(file_path)
               self.set_scope(ws_id, scope)
       
       def remove_file(self, ws_id: str, file_path: str) -> None:
           """Remove file from scope."""
           scope = self.get_scope(ws_id)
           if file_path in scope:
               scope.remove(file_path)
               self.set_scope(ws_id, scope)
       
       def is_in_scope(self, ws_id: str, file_path: str) -> bool:
           """Check if file is in scope."""
           scope = self.get_scope(ws_id)
           
           # Empty scope = all files allowed
           if not scope:
               return True
           
           return file_path in scope
   ```

2. **Add CLI commands**

   ```python
   # src/sdp/cli/ws.py (add scope subcommand)
   import typer
   from sdp.beads import create_beads_client
   from sdp.beads.scope_manager import ScopeManager
   
   scope_app = typer.Typer(help="Manage WS file scope")
   
   @scope_app.command("show")
   def show_scope(ws_id: str) -> None:
       """Show scope files for workstream."""
       manager = ScopeManager(create_beads_client())
       
       try:
           scope = manager.get_scope(ws_id)
       except ValueError as e:
           typer.echo(f"❌ {e}")
           raise typer.Exit(1)
       
       if not scope:
           typer.echo(f"Scope for {ws_id}: unrestricted (all files allowed)")
           return
       
       typer.echo(f"Scope for {ws_id}: {len(scope)} files")
       for f in scope:
           typer.echo(f"  - {f}")
   
   @scope_app.command("add")
   def add_to_scope(ws_id: str, file_path: str) -> None:
       """Add file to workstream scope."""
       manager = ScopeManager(create_beads_client())
       manager.add_file(ws_id, file_path)
       typer.echo(f"✅ Added {file_path} to {ws_id} scope")
   
   @scope_app.command("remove")
   def remove_from_scope(ws_id: str, file_path: str) -> None:
       """Remove file from workstream scope."""
       manager = ScopeManager(create_beads_client())
       manager.remove_file(ws_id, file_path)
       typer.echo(f"✅ Removed {file_path} from {ws_id} scope")
   
   # Register
   app.add_typer(scope_app, name="scope")
   ```

3. **Update Guard to use ScopeManager**

   ```python
   # src/sdp/guard/skill.py (update)
   from sdp.beads.scope_manager import ScopeManager
   
   class GuardSkill:
       def __init__(self, beads_client: BeadsClient):
           self._client = beads_client
           self._scope_manager = ScopeManager(beads_client)
           # ...
       
       def check_edit(self, file_path: str) -> GuardResult:
           if not self._active_ws:
               return GuardResult(allowed=False, ...)
           
           # Use ScopeManager
           if not self._scope_manager.is_in_scope(self._active_ws, file_path):
               scope = self._scope_manager.get_scope(self._active_ws)
               return GuardResult(
                   allowed=False,
                   ws_id=self._active_ws,
                   reason=f"File {file_path} not in WS scope",
                   scope_files=scope
               )
           
           return GuardResult(allowed=True, ...)
   ```

### Output Files

- `src/sdp/beads/scope_manager.py`
- `src/sdp/cli/ws.py` (scope commands)
- `tests/unit/test_scope_manager.py`

### Completion Criteria

```bash
# CLI works
sdp ws scope --help
sdp ws scope show bd-0001

# Tests pass
pytest tests/unit/test_scope_manager.py -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
