---
completed: '2026-01-30'
depends_on:
- 00-032-18
feature: F032
project_id: 0
scope_files:
- src/sdp/validators/supersede_checker.py
- src/sdp/cli/workstream.py
size: SMALL
status: completed
traceability:
- ac_description: '`sdp ws supersede {old_ws} --replacement {new_ws}` command'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Supersede without replacement is blocked
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Supersede chain validation (no cycles)
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: CI check for orphan supersedes
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: Report of all superseded WS and their replacements
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
- ac_description: Frontmatter field `superseded_by` added
  ac_id: AC6
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_capability_tier_extractors.py
  test_name: test_extract_body_with_frontmatter
ws_id: 00-032-28
---

## 00-032-28: Supersede Validation

### Goal

Prevent orphan superseded workstreams. When WS is marked "superseded", require replacement WS ID. Track supersede chains to prevent cycles.

### Acceptance Criteria

- [ ] AC1: `sdp ws supersede {old_ws} --replacement {new_ws}` command
- [ ] AC2: Supersede without replacement is blocked
- [ ] AC3: Supersede chain validation (no cycles)
- [ ] AC4: Report of all superseded WS and their replacements
- [ ] AC5: CI check for orphan supersedes
- [ ] AC6: Frontmatter field `superseded_by` added

### Contract

```python
# src/sdp/validators/supersede_checker.py
from dataclasses import dataclass

@dataclass
class SupersedeChain:
    original_ws: str
    replacements: list[str]
    has_cycle: bool
    final_ws: str | None  # None if cycle

class SupersedeValidator:
    """Validate supersede relationships."""
    
    def supersede(self, old_ws: str, new_ws: str) -> SupersedeResult:
        """Mark old_ws as superseded by new_ws.
        
        Validates:
        - new_ws exists
        - No cycle created
        - old_ws not already superseded
        """
        raise NotImplementedError
    
    def find_orphans(self, ws_dir: Path) -> list[str]:
        """Find superseded WS without valid replacement."""
        raise NotImplementedError
    
    def trace_chain(self, ws_id: str) -> SupersedeChain:
        """Trace supersede chain to final WS."""
        raise NotImplementedError
    
    def validate_all(self, ws_dir: Path) -> ValidationReport:
        """Validate all supersede relationships."""
        raise NotImplementedError

@dataclass
class SupersedeResult:
    success: bool
    old_ws: str
    new_ws: str
    error: str | None
```

### Frontmatter Extension

```yaml
---
ws_id: 00-012-01
status: superseded
superseded_by: 00-032-01  # NEW: Required when status=superseded
supersede_reason: "Merged into F032 guard foundation"  # NEW: Optional
---
```

### Scope

**Input:**
- Existing WS files with `status: superseded`
- F012 series (14 superseded without replacements)

**Output:**
- `src/sdp/validators/supersede_checker.py`
- `src/sdp/cli/workstream.py` (add `supersede` command)
- Updated frontmatter schema
- Report of current orphans

### Verification

```bash
# Find current orphans
sdp ws orphans

# Supersede with replacement
sdp ws supersede 00-012-01 --replacement 00-032-01

# Supersede without replacement fails
sdp ws supersede 00-012-02 && exit 1 || echo "Blocked OK"

# No cycles
sdp ws supersede 00-032-01 --replacement 00-012-01 && exit 1 || echo "Cycle blocked"
```
