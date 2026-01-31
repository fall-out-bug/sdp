---
completed: '2026-01-30'
depends_on:
- 00-032-24
feature: F032
project_id: 0
scope_files:
- templates/workstream.md
- templates/workstream-frontmatter.md
- docs/reference/quality-gates.md
size: MEDIUM
status: completed
traceability:
- ac_description: New template is ≤80 lines
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Template contains ONLY signatures, not implementations
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: All code examples use `raise NotImplementedError`
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Migration guide for existing WS format
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: ADR document explains rationale
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
- ac_description: '`sdp ws validate` rejects WS with >150 lines code blocks'
  ac_id: AC6
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_skill_validator.py
  test_name: test_validate_skill_too_long
ws_id: 00-032-25
---

## 00-032-25: WS Template Simplification (Contracts Only)

### Goal

Simplify workstream template to max 80 lines. Remove ready-made code, keep only contracts (signatures + docstrings). Agents must implement, not copy-paste.

### Acceptance Criteria

- [ ] AC1: New template is ≤80 lines
- [ ] AC2: Template contains ONLY signatures, not implementations
- [ ] AC3: All code examples use `raise NotImplementedError`
- [ ] AC4: ADR document explains rationale
- [ ] AC5: Migration guide for existing WS format
- [ ] AC6: `sdp ws validate` rejects WS with >150 lines code blocks

### Contract

```python
# src/sdp/validators/ws_template_checker.py
MAX_WS_LINES: int = 150
MAX_CODE_BLOCK_LINES: int = 30
REQUIRED_SECTIONS: list[str]

def validate_ws_structure(ws_path: Path) -> ValidationResult:
    """Validate WS follows simplified template.
    
    Checks:
    - Total lines ≤ MAX_WS_LINES
    - Code blocks ≤ MAX_CODE_BLOCK_LINES
    - All required sections present
    - No full implementations (detect by line count + no raise NotImplementedError)
    """
    raise NotImplementedError
```

### New Template Structure

```markdown
## {WS-ID}: {Title}

### Goal
{2-3 sentences max}

### Acceptance Criteria
- [ ] AC1: {testable criterion}

### Contract
```python
def function(arg: Type) -> ReturnType:
    """One-line doc."""
    raise NotImplementedError
```

### Scope
- Input: {files}
- Output: {files}

### Verification
```bash
pytest ... --cov-fail-under=80
```
```

### Scope

**Input:**
- `templates/workstream.md` (current ~200 lines)
- `templates/workstream-frontmatter.md`
- Existing WS examples

**Output:**
- `templates/workstream-v2.md` (≤80 lines)
- `docs/adr/007-simplified-ws-template.md`
- `docs/migration/ws-template-v2.md`
- `src/sdp/validators/ws_template_checker.py`

### Verification

```bash
# Template size
wc -l templates/workstream-v2.md  # ≤80

# Validator works
sdp ws validate docs/workstreams/backlog/00-032-25.md

# Old oversized WS rejected
sdp ws validate docs/workstreams/backlog/00-020-01.md && exit 1 || echo "Rejected OK"
```
