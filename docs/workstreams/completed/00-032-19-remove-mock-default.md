---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-18
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '`create_beads_client()` по умолчанию использует RealBeadsClient'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: В CI Beads устанавливается в setup step
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Unit tests используют mock через fixture
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
- ac_description: Документация обновлена
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-19
---

## 00-032-19: Remove Mock Default

### 🎯 Goal

**What must WORK after completing this WS:**
- Real Beads is default (`BEADS_USE_MOCK=false`)
- Mock используется только в тестах
- CI настроен с real Beads

**Acceptance Criteria:**
- [ ] AC1: `create_beads_client()` по умолчанию использует RealBeadsClient
- [ ] AC2: В CI Beads устанавливается в setup step
- [ ] AC3: Unit tests используют mock через fixture
- [ ] AC4: Документация обновлена

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Mock default → агенты игнорируют Beads.

**Solution**: Real default, mock только для tests.

### Dependencies

- **00-032-18**: Beads Status Sync

### Steps

1. **Update default in factory**

   ```python
   # src/sdp/beads/__init__.py
   import os
   import shutil
   
   def create_beads_client(use_mock: bool | None = None) -> BeadsClient:
       """Create Beads client.
       
       Default behavior:
       - If bd CLI installed → RealBeadsClient
       - If bd not installed → MockBeadsClient with warning
       
       Override with BEADS_USE_MOCK=true for tests.
       """
       if use_mock is None:
           env_mock = os.getenv("BEADS_USE_MOCK")
           if env_mock is not None:
               use_mock = env_mock.lower() == "true"
           else:
               # Auto-detect: use real if bd installed
               use_mock = shutil.which("bd") is None
               if use_mock:
                   import warnings
                   warnings.warn(
                       "Beads CLI (bd) not found. Using mock client. "
                       "Install: go install github.com/steveyegge/beads/cmd/bd@latest"
                   )
       
       if use_mock:
           return MockBeadsClient()
       
       return RealBeadsClient()
   ```

2. **Update CI workflow**

   ```yaml
   # .github/workflows/ci-critical.yml (add Beads setup)
   steps:
     - name: Setup Go
       uses: actions/setup-go@v5
       with:
         go-version: '1.21'
     
     - name: Install Beads CLI
       run: |
         go install github.com/steveyegge/beads/cmd/bd@latest
         echo "$(go env GOPATH)/bin" >> $GITHUB_PATH
     
     - name: Verify Beads
       run: bd version
   ```

3. **Update test fixtures**

   ```python
   # tests/conftest.py
   import pytest
   import os
   
   @pytest.fixture(autouse=True)
   def use_mock_beads(monkeypatch):
       """Force mock Beads in all tests."""
       monkeypatch.setenv("BEADS_USE_MOCK", "true")
   
   @pytest.fixture
   def real_beads(monkeypatch):
       """Use real Beads for specific tests."""
       monkeypatch.setenv("BEADS_USE_MOCK", "false")
   ```

4. **Update documentation**

   ```markdown
   # docs/setup/beads-installation.md (update)
   
   ## Default Behavior
   
   As of SDP v0.6.0, real Beads is the default:
   
   - If `bd` CLI installed → Uses real Beads
   - If `bd` not installed → Falls back to mock with warning
   
   ## Force Mock Mode
   
   For development without Go:
   
   ```bash
   export BEADS_USE_MOCK=true
   ```
   
   ## CI Configuration
   
   CI workflows install Beads automatically. No action needed.
   ```

### Output Files

- `src/sdp/beads/__init__.py` (updated)
- `.github/workflows/ci-critical.yml` (updated)
- `tests/conftest.py` (updated)
- `docs/setup/beads-installation.md` (updated)

### Completion Criteria

```bash
# With bd installed, uses real
python -c "from sdp.beads import create_beads_client; print(type(create_beads_client()))"
# Expected: RealBeadsClient

# Force mock
BEADS_USE_MOCK=true python -c "from sdp.beads import create_beads_client; print(type(create_beads_client()))"
# Expected: MockBeadsClient

# Tests still pass (use mock)
pytest tests/ -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______
