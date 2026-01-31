---
completed: '2026-01-30'
depends_on:
- 00-032-26
feature: F032
project_id: 0
scope_files:
- hooks/post-ws-complete.sh
- src/sdp/hooks/ws_complete.py
size: SMALL
status: completed
traceability:
- ac_description: '`hooks/post-ws-complete.sh` calls `sdp ws verify`'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Hook blocks status change if verification fails
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Hook integrates with Beads status update
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Bypass flag for manual override (requires reason)
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: Clear error messages with remediation steps
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-27
---

## 00-032-27: Post-WS-Complete Enforcement Hook

### Goal

Create hook that runs after agent marks WS complete. Verifies completion is real before allowing status change. Blocks fake "✅ DONE" claims.

### Acceptance Criteria

- [ ] AC1: `hooks/post-ws-complete.sh` calls `sdp ws verify`
- [ ] AC2: Hook blocks status change if verification fails
- [ ] AC3: Hook integrates with Beads status update
- [ ] AC4: Clear error messages with remediation steps
- [ ] AC5: Bypass flag for manual override (requires reason)

### Contract

```python
# src/sdp/hooks/ws_complete.py
from sdp.validators.ws_completion import WSCompletionVerifier

class PostWSCompleteHook:
    """Hook to verify WS completion before status change."""
    
    def __init__(self, verifier: WSCompletionVerifier):
        raise NotImplementedError
    
    def run(self, ws_id: str, bypass: bool = False, reason: str = "") -> HookResult:
        """Run verification and return result.
        
        Args:
            ws_id: Workstream to verify
            bypass: Skip verification (requires reason)
            reason: Bypass justification (logged)
            
        Returns:
            HookResult with pass/fail and messages
        """
        raise NotImplementedError

@dataclass
class HookResult:
    passed: bool
    ws_id: str
    messages: list[str]
    bypass_used: bool
    bypass_reason: str | None
```

```bash
# hooks/post-ws-complete.sh
#!/bin/bash
# Post-WS-Complete Hook - Verifies completion before status change

WS_ID=$1
BYPASS=${2:-false}
REASON=${3:-""}

# Run Python verification
python -m sdp.hooks.ws_complete "$WS_ID" --bypass="$BYPASS" --reason="$REASON"

exit $?
```

### Scope

**Input:**
- WS ID from agent status change
- `sdp ws verify` output

**Output:**
- `hooks/post-ws-complete.sh`
- `src/sdp/hooks/__init__.py`
- `src/sdp/hooks/ws_complete.py`
- Updated Beads integration

### Verification

```bash
# Hook exists and executable
test -x hooks/post-ws-complete.sh

# Incomplete WS blocked
hooks/post-ws-complete.sh 00-032-99-incomplete && exit 1 || echo "Blocked OK"

# Complete WS passes
hooks/post-ws-complete.sh 00-032-27

# Bypass works with reason
hooks/post-ws-complete.sh 00-032-99 true "Manual verification done"
```
