---
assignee: null
completed: '2026-01-30'
depends_on: []
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '`sdp doctor` показывает статус Beads CLI'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Если `bd` отсутствует, показывает инструкции
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: '`docs/setup/beads-installation.md` создан'
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
- ac_description: Health check интегрирован в doctor module
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-15
---

## 00-032-15: Beads Go Installation Check

### 🎯 Goal

**What must WORK after completing this WS:**
- `sdp doctor` проверяет наличие `bd` command
- Инструкции по установке если отсутствует
- Health check в CI pipeline

**Acceptance Criteria:**
- [ ] AC1: `sdp doctor` показывает статус Beads CLI
- [ ] AC2: Если `bd` отсутствует, показывает инструкции
- [ ] AC3: `docs/setup/beads-installation.md` создан
- [ ] AC4: Health check интегрирован в doctor module

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Beads требует Go и `bd` CLI. Нет проверки при старте.

**Solution**: Health check в `sdp doctor`.

### Dependencies

None (independent)

### Steps

1. **Create Beads health check**

   ```python
   # src/sdp/health_checks/beads.py
   import subprocess
   import shutil
   from dataclasses import dataclass
   
   from sdp.health_checks.base import HealthCheck, HealthStatus
   
   @dataclass
   class BeadsHealthResult:
       installed: bool
       version: str | None
       go_installed: bool
       path: str | None
   
   class BeadsHealthCheck(HealthCheck):
       """Check Beads CLI availability."""
       
       name = "beads"
       description = "Beads task management CLI"
       
       def check(self) -> HealthStatus:
           result = self._check_beads()
           
           if result.installed:
               return HealthStatus(
                   healthy=True,
                   message=f"Beads CLI v{result.version} at {result.path}",
               )
           
           # Build remediation message
           remediation = []
           if not result.go_installed:
               remediation.append("1. Install Go: brew install go")
           remediation.append("2. Install Beads: go install github.com/steveyegge/beads/cmd/bd@latest")
           remediation.append("3. Add to PATH: export PATH=$PATH:$(go env GOPATH)/bin")
           
           return HealthStatus(
               healthy=False,
               message="Beads CLI not found",
               remediation="\n".join(remediation),
               docs_url="docs/setup/beads-installation.md"
           )
       
       def _check_beads(self) -> BeadsHealthResult:
           # Check Go
           go_installed = shutil.which("go") is not None
           
           # Check bd
           bd_path = shutil.which("bd")
           if not bd_path:
               return BeadsHealthResult(
                   installed=False,
                   version=None,
                   go_installed=go_installed,
                   path=None
               )
           
           # Get version
           try:
               result = subprocess.run(
                   ["bd", "version"],
                   capture_output=True,
                   text=True
               )
               version = result.stdout.strip()
           except Exception:
               version = "unknown"
           
           return BeadsHealthResult(
               installed=True,
               version=version,
               go_installed=go_installed,
               path=bd_path
           )
   ```

2. **Register in doctor**

   ```python
   # src/sdp/cli/doctor.py (update)
   from sdp.health_checks.beads import BeadsHealthCheck
   
   HEALTH_CHECKS = [
       # ... existing checks ...
       BeadsHealthCheck(),
   ]
   ```

3. **Create installation guide**

   ```markdown
   # docs/setup/beads-installation.md
   # Beads CLI Installation
   
   ## Prerequisites
   
   - Go 1.21+ (for building Beads)
   
   ## Installation
   
   ### macOS
   
   ```bash
   # Install Go
   brew install go
   
   # Install Beads
   go install github.com/steveyegge/beads/cmd/bd@latest
   
   # Add to PATH (add to ~/.zshrc)
   export PATH=$PATH:$(go env GOPATH)/bin
   
   # Verify
   bd version
   ```
   
   ### Linux
   
   ```bash
   # Install Go (Ubuntu/Debian)
   sudo apt install golang-go
   
   # Install Beads
   go install github.com/steveyegge/beads/cmd/bd@latest
   
   # Add to PATH (add to ~/.bashrc)
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
   
   ## Verify Installation
   
   ```bash
   # Check with sdp doctor
   sdp doctor
   
   # Expected output:
   # ✅ beads: Beads CLI v1.0.0 at /Users/you/go/bin/bd
   ```
   
   ## Troubleshooting
   
   ### "bd: command not found"
   
   1. Check Go installation: `go version`
   2. Check GOPATH: `go env GOPATH`
   3. Add to PATH: `export PATH=$PATH:$(go env GOPATH)/bin`
   
   ### Permission denied
   
   ```bash
   chmod +x $(go env GOPATH)/bin/bd
   ```
   ```

### Output Files

- `src/sdp/health_checks/beads.py`
- `docs/setup/beads-installation.md`
- `tests/unit/test_beads_health_check.py`

### Completion Criteria

```bash
# Health check works
sdp doctor

# Module imports
python -c "from sdp.health_checks.beads import BeadsHealthCheck"

# Tests pass
pytest tests/unit/test_beads_health_check.py -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______
