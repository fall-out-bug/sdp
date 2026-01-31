---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-01
feature: F032
github_issue: null
project_id: 0
size: MEDIUM
status: completed
traceability:
- ac_description: '`sdp guard check <file>` возвращает exit code 0 (allowed) или 1
    (blocked)'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`sdp guard activate <ws_id>` сохраняет WS в `.sdp/state.json`'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: '`sdp guard status` показывает текущий WS и scope'
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Unit + integration тесты (≥80%)
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`sdp guard deactivate` очищает активный WS'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-02
---

## 00-032-02: Pre-Edit Validator CLI

### 🎯 Goal

**What must WORK after completing this WS:**
- CLI команда `sdp guard check <file>` проверяет разрешение на edit
- CLI команда `sdp guard activate <ws_id>` устанавливает активный WS
- State сохраняется в `.sdp/state.json`

**Acceptance Criteria:**
- [ ] AC1: `sdp guard check <file>` возвращает exit code 0 (allowed) или 1 (blocked)
- [ ] AC2: `sdp guard activate <ws_id>` сохраняет WS в `.sdp/state.json`
- [ ] AC3: `sdp guard status` показывает текущий WS и scope
- [ ] AC4: `sdp guard deactivate` очищает активный WS
- [ ] AC5: Unit + integration тесты (≥80%)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: GuardSkill создан (WS-01), но нет CLI interface для использования.

**Solution**: Typer CLI commands для guard operations.

### Dependencies

- **00-032-01**: Guard Skill Foundation (provides GuardSkill class)

### Input Files

- `src/sdp/guard/skill.py` (from WS-01)
- `src/sdp/cli/main.py` (CLI entry point)

### Steps

1. **Create StateManager**

   ```python
   # src/sdp/guard/state.py
   import json
   from pathlib import Path
   from dataclasses import dataclass, asdict
   from datetime import datetime
   
   @dataclass
   class GuardState:
       active_ws: str | None = None
       activated_at: str | None = None
       scope_files: list[str] | None = None
   
   class StateManager:
       """Manages guard state persistence."""
       
       STATE_FILE = Path(".sdp/state.json")
       
       @classmethod
       def load(cls) -> GuardState:
           """Load state from file."""
           if not cls.STATE_FILE.exists():
               return GuardState()
           
           with open(cls.STATE_FILE) as f:
               data = json.load(f)
           return GuardState(**data)
       
       @classmethod
       def save(cls, state: GuardState) -> None:
           """Save state to file."""
           cls.STATE_FILE.parent.mkdir(exist_ok=True)
           with open(cls.STATE_FILE, "w") as f:
               json.dump(asdict(state), f, indent=2)
       
       @classmethod
       def clear(cls) -> None:
           """Clear state."""
           cls.save(GuardState())
   ```

2. **Create CLI commands**

   ```python
   # src/sdp/cli/guard.py
   import typer
   from pathlib import Path
   from datetime import datetime
   
   from sdp.guard.skill import GuardSkill
   from sdp.guard.state import StateManager, GuardState
   from sdp.beads import create_beads_client
   
   app = typer.Typer(help="Pre-edit guard commands")
   
   @app.command("activate")
   def activate(ws_id: str) -> None:
       """Activate workstream for editing."""
       client = create_beads_client()
       guard = GuardSkill(client)
       
       try:
           guard.activate(ws_id)
       except ValueError as e:
           typer.echo(f"❌ {e}")
           raise typer.Exit(1)
       
       ws = client.get_task(ws_id)
       scope = ws.sdp_metadata.get("scope_files", [])
       
       state = GuardState(
           active_ws=ws_id,
           activated_at=datetime.utcnow().isoformat(),
           scope_files=scope,
       )
       StateManager.save(state)
       
       typer.echo(f"✅ Activated WS: {ws_id}")
       if scope:
           typer.echo(f"   Scope: {len(scope)} files")
       else:
           typer.echo(f"   Scope: unrestricted")
   
   @app.command("check")
   def check_file(file_path: Path) -> None:
       """Check if file edit is allowed."""
       state = StateManager.load()
       
       if not state.active_ws:
           typer.echo("❌ No active WS. Run: sdp guard activate <ws_id>")
           raise typer.Exit(1)
       
       client = create_beads_client()
       guard = GuardSkill(client)
       guard._active_ws = state.active_ws
       
       result = guard.check_edit(str(file_path))
       
       if not result.allowed:
           typer.echo(f"❌ {result.reason}")
           if result.scope_files:
               typer.echo(f"   Allowed files:")
               for f in result.scope_files[:10]:
                   typer.echo(f"     - {f}")
           raise typer.Exit(1)
       
       typer.echo(f"✅ Edit allowed: {file_path}")
   
   @app.command("status")
   def status() -> None:
       """Show current guard status."""
       state = StateManager.load()
       
       if not state.active_ws:
           typer.echo("No active workstream")
           return
       
       typer.echo(f"Active WS: {state.active_ws}")
       typer.echo(f"Activated: {state.activated_at}")
       if state.scope_files:
           typer.echo(f"Scope: {len(state.scope_files)} files")
       else:
           typer.echo("Scope: unrestricted")
   
   @app.command("deactivate")
   def deactivate() -> None:
       """Deactivate current workstream."""
       StateManager.clear()
       typer.echo("✅ Guard deactivated")
   ```

3. **Register in main CLI**

   ```python
   # src/sdp/cli/main.py
   from sdp.cli.guard import app as guard_app
   
   app.add_typer(guard_app, name="guard")
   ```

4. **Write tests**

   ```python
   # tests/unit/test_cli_guard.py
   import pytest
   from typer.testing import CliRunner
   from pathlib import Path
   from sdp.cli.main import app
   from sdp.guard.state import StateManager
   
   runner = CliRunner()
   
   class TestGuardCLI:
       def setup_method(self):
           StateManager.clear()
       
       def test_check_without_active_ws_fails(self):
           """AC1: Check without active WS returns exit 1."""
           result = runner.invoke(app, ["guard", "check", "any.py"])
           assert result.exit_code == 1
           assert "No active WS" in result.output
       
       def test_activate_saves_state(self, mock_beads):
           """AC2: Activate saves state to file."""
           result = runner.invoke(app, ["guard", "activate", "00-032-01"])
           assert result.exit_code == 0
           
           state = StateManager.load()
           assert state.active_ws == "00-032-01"
       
       def test_status_shows_active_ws(self, mock_beads):
           """AC3: Status shows current WS."""
           runner.invoke(app, ["guard", "activate", "00-032-01"])
           result = runner.invoke(app, ["guard", "status"])
           
           assert "00-032-01" in result.output
       
       def test_deactivate_clears_state(self, mock_beads):
           """AC4: Deactivate clears active WS."""
           runner.invoke(app, ["guard", "activate", "00-032-01"])
           runner.invoke(app, ["guard", "deactivate"])
           
           state = StateManager.load()
           assert state.active_ws is None
   ```

### Output Files

- `src/sdp/guard/state.py`
- `src/sdp/cli/guard.py`
- `tests/unit/test_cli_guard.py`
- `tests/integration/test_guard_flow.py`

### Completion Criteria

```bash
# CLI works
sdp guard --help
sdp guard activate --help

# Tests pass
pytest tests/unit/test_cli_guard.py -v

# Coverage
pytest tests/unit/test_cli_guard.py --cov=src/sdp/cli/guard --cov-fail-under=80
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
