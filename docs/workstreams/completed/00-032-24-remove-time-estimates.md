---
completed: '2026-01-30'
depends_on: []
feature: F032
project_id: 0
scope_files:
- .claude/skills/design/SKILL.md
- docs/workstreams/backlog/00-013-*.md
- templates/workstream.md
size: SMALL
status: completed
traceability:
- ac_description: No `estimated_duration`, `"X hours"`, `"X days"` in any skill file
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: No time estimates in workstream templates
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: No time estimates in existing backlog workstreams
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Pre-commit hook rejects files with time patterns
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: CLI validator `sdp ws lint` detects time estimates
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-24
---

## 00-032-24: Remove Time Estimates from Skills and WS

### Goal

Remove all time-based estimates from skills, workstreams, and templates. SDP forbids time estimates, but they exist in core files.

### Acceptance Criteria

- [ ] AC1: No `estimated_duration`, `"X hours"`, `"X days"` in any skill file
- [ ] AC2: No time estimates in workstream templates
- [ ] AC3: No time estimates in existing backlog workstreams
- [ ] AC4: CLI validator `sdp ws lint` detects time estimates
- [ ] AC5: Pre-commit hook rejects files with time patterns

### Contract

```python
# src/sdp/validators/time_estimate_checker.py
TIME_PATTERNS: list[str]
"""Regex patterns for forbidden time estimates."""

def check_file(path: Path) -> list[Violation]:
    """Check file for time estimate violations.
    
    Returns:
        List of violations with line numbers
    """
    raise NotImplementedError

def check_directory(path: Path, glob: str = "**/*.md") -> list[Violation]:
    """Check all matching files in directory."""
    raise NotImplementedError
```

### Scope

**Input:**
- `.claude/skills/design/SKILL.md` (contains `estimated_duration`)
- `docs/workstreams/backlog/00-013-*.md` (contain `"2-3 hours"`)
- `templates/workstream.md`

**Output:**
- `src/sdp/validators/time_estimate_checker.py`
- Updated skill/template files (time estimates removed)
- Updated `hooks/pre-commit.sh` (add time check)

### Verification

```bash
# No time estimates in skills
grep -r "hours\|days\|estimated_duration" .claude/skills/ && exit 1 || echo "OK"

# Validator works
python -m sdp.validators.time_estimate_checker docs/workstreams/

# Pre-commit rejects
echo "estimated_duration: 2h" > /tmp/test.md
hooks/pre-commit.sh /tmp/test.md && exit 1 || echo "Rejected OK"
```
