---
completed: '2026-01-30'
depends_on:
- 00-032-21
feature: F032
project_id: 0
scope_files:
- src/sdp/validators/ws_completion.py
- src/sdp/cli/workstream.py
size: MEDIUM
status: completed
traceability:
- ac_description: '`sdp ws verify {ws_id}` checks all output files exist'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`sdp ws verify` runs completion criteria commands'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: '`sdp ws verify` checks test coverage for module'
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Verification result stored in WS frontmatter
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`sdp ws verify` validates AC checkboxes match reality'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
- ac_description: Integration with Beads status sync
  ac_id: AC6
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_sync_service.py
  test_name: test_sync_updates_beads_from_local
ws_id: 00-032-26
---

## 00-032-26: Automated WS Completion Verification

### Goal

Verify WS is actually complete before marking done. Check output files exist, tests pass, coverage met. Prevent agents from self-reporting "✅ DONE" without evidence.

### Acceptance Criteria

- [ ] AC1: `sdp ws verify {ws_id}` checks all output files exist
- [ ] AC2: `sdp ws verify` runs completion criteria commands
- [ ] AC3: `sdp ws verify` checks test coverage for module
- [ ] AC4: `sdp ws verify` validates AC checkboxes match reality
- [ ] AC5: Verification result stored in WS frontmatter
- [ ] AC6: Integration with Beads status sync

### Contract

```python
# src/sdp/validators/ws_completion.py
from dataclasses import dataclass
from pathlib import Path

@dataclass
class VerificationResult:
    ws_id: str
    passed: bool
    checks: list[CheckResult]
    coverage_actual: float | None
    missing_files: list[str]
    failed_commands: list[str]

@dataclass  
class CheckResult:
    name: str
    passed: bool
    message: str
    evidence: str | None  # Command output or file path

class WSCompletionVerifier:
    """Verify workstream completion with evidence."""
    
    def verify(self, ws_id: str) -> VerificationResult:
        """Run all verification checks.
        
        Checks:
        1. All scope_files output exist
        2. All Verification commands pass
        3. Coverage meets threshold
        4. AC checkboxes accurate
        """
        raise NotImplementedError
    
    def verify_output_files(self, ws: Workstream) -> list[CheckResult]:
        """Check all output files in scope exist."""
        raise NotImplementedError
    
    def verify_commands(self, ws: Workstream) -> list[CheckResult]:
        """Run verification commands and check exit codes."""
        raise NotImplementedError
    
    def verify_coverage(self, ws: Workstream) -> CheckResult:
        """Check test coverage meets threshold."""
        raise NotImplementedError
```

### Scope

**Input:**
- WS file with `scope_files`, `Verification` section
- Codebase for file existence checks
- pytest for coverage verification

**Output:**
- `src/sdp/validators/ws_completion.py`
- `src/sdp/cli/workstream.py` (add `verify` command)
- `tests/unit/validators/test_ws_completion.py`

### Verification

```bash
# Verify command works
sdp ws verify 00-032-26

# Incomplete WS fails
# (create incomplete WS, should fail)
sdp ws verify 00-032-99-fake && exit 1 || echo "Failed correctly"

# Coverage check
pytest tests/unit/validators/test_ws_completion.py --cov-fail-under=80
```
