---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-02
- 00-032-03
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: 'E2E тест: happy path (activate → edit allowed → complete)'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: 'E2E тест: edit blocked без активного WS'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: 'E2E тест: edit blocked для файла вне scope'
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Coverage report показывает ≥80% для guard module
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: 'E2E тест: concurrent WS activation blocked'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-04
---

## 00-032-04: Guard Integration Tests

### 🎯 Goal

**What must WORK after completing this WS:**
- E2E тесты полного flow: activate → check → edit → complete
- Все edge cases покрыты интеграционными тестами

**Acceptance Criteria:**
- [ ] AC1: E2E тест: happy path (activate → edit allowed → complete)
- [ ] AC2: E2E тест: edit blocked без активного WS
- [ ] AC3: E2E тест: edit blocked для файла вне scope
- [ ] AC4: E2E тест: concurrent WS activation blocked
- [ ] AC5: Coverage report показывает ≥80% для guard module

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Unit tests проверяют компоненты по отдельности. Нужны E2E тесты.

**Solution**: Интеграционные тесты с реальным (mock) Beads client.

### Dependencies

- **00-032-02**: Pre-Edit Validator CLI
- **00-032-03**: Active WS Tracker

### Steps

1. **Create test fixtures**

   ```python
   # tests/conftest.py (add to existing)
   import pytest
   from sdp.beads import MockBeadsClient
   from sdp.beads.models import BeadsTaskCreate, BeadsPriority
   
   @pytest.fixture
   def beads_client():
       """Fresh mock Beads client for each test."""
       return MockBeadsClient()
   
   @pytest.fixture
   def sample_workstream(beads_client):
       """Create a sample WS for testing."""
       task = beads_client.create_task(BeadsTaskCreate(
           title="Test WS",
           description="Test workstream",
           priority=BeadsPriority.MEDIUM,
           sdp_metadata={
               "scope_files": [
                   "src/sdp/guard/skill.py",
                   "src/sdp/guard/models.py",
               ]
           }
       ))
       return task
   ```

2. **Write E2E tests**

   ```python
   # tests/e2e/test_guard_flow.py
   import pytest
   from typer.testing import CliRunner
   from pathlib import Path
   import tempfile
   import os
   
   from sdp.cli.main import app
   from sdp.guard.state import StateManager
   
   runner = CliRunner()
   
   class TestGuardE2E:
       """End-to-end tests for guard workflow."""
       
       @pytest.fixture(autouse=True)
       def setup(self, tmp_path, monkeypatch):
           """Setup clean environment for each test."""
           # Use temp directory for state
           state_file = tmp_path / ".sdp" / "state.json"
           monkeypatch.setattr(StateManager, "STATE_FILE", state_file)
           
           # Use mock Beads
           monkeypatch.setenv("BEADS_USE_MOCK", "true")
       
       def test_happy_path_activate_edit_complete(self, sample_workstream):
           """AC1: Full workflow succeeds."""
           ws_id = sample_workstream.id
           
           # 1. Activate
           result = runner.invoke(app, ["guard", "activate", ws_id])
           assert result.exit_code == 0
           assert "Activated" in result.output
           
           # 2. Check allowed file
           result = runner.invoke(app, [
               "guard", "check", "src/sdp/guard/skill.py"
           ])
           assert result.exit_code == 0
           assert "allowed" in result.output
           
           # 3. Complete (via guard complete or ws complete)
           result = runner.invoke(app, ["guard", "deactivate"])
           assert result.exit_code == 0
       
       def test_edit_blocked_without_active_ws(self):
           """AC2: Edit blocked without activation."""
           result = runner.invoke(app, [
               "guard", "check", "any/file.py"
           ])
           assert result.exit_code == 1
           assert "No active WS" in result.output
       
       def test_edit_blocked_outside_scope(self, sample_workstream):
           """AC3: Edit blocked for file outside scope."""
           ws_id = sample_workstream.id
           
           # Activate WS with restricted scope
           runner.invoke(app, ["guard", "activate", ws_id])
           
           # Try to edit file outside scope
           result = runner.invoke(app, [
               "guard", "check", "src/other/forbidden.py"
           ])
           assert result.exit_code == 1
           assert "not in WS scope" in result.output
       
       def test_concurrent_activation_blocked(self, sample_workstream, beads_client):
           """AC4: Cannot activate second WS."""
           ws1_id = sample_workstream.id
           
           # Create second WS
           ws2 = beads_client.create_task(BeadsTaskCreate(
               title="Second WS",
               description="Another workstream",
           ))
           
           # Activate first
           runner.invoke(app, ["guard", "activate", ws1_id])
           
           # Try to activate second
           result = runner.invoke(app, ["guard", "activate", ws2.id])
           assert result.exit_code == 1
           assert "already in progress" in result.output
   ```

3. **Generate coverage report**

   ```bash
   # tests/e2e/conftest.py
   # Configure pytest-cov for guard module
   ```

### Output Files

- `tests/e2e/test_guard_flow.py`
- `tests/conftest.py` (updated with fixtures)

### Completion Criteria

```bash
# E2E tests pass
pytest tests/e2e/test_guard_flow.py -v

# Full coverage report
pytest tests/ --cov=src/sdp/guard --cov-report=term-missing --cov-fail-under=80

# All guard module tests
pytest tests/ -k guard -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
