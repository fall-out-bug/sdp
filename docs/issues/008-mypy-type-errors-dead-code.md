# 008: mypy type errors — dead code and missing annotations

**Source:** /issue (2026-01-31)  
**Status:** Fixed  
**Priority:** P2 (MEDIUM)  
**Route:** /bugfix (applied)

## Problem

30 mypy errors across 14 files. Project-wide type checking is broken.

**Measurement:**
```bash
uv run mypy src/sdp --ignore-missing-imports
# Found 30 errors in 14 files
```

## Error Categories

### 1. Dead Code — CLI for Unimplemented API (Critical)

**Files:** `cli/tier_metrics.py`, `cli/tier_promote.py`

Code references attributes/methods that don't exist:

```python
# tier_metrics.py:42 — TierMetrics has NO 'history' attribute
if metrics.history:
    for entry in metrics.history[-5:]:
        ...

# tier_metrics.py:51 — TierMetricsStore has NO 'list_all' method
all_metrics = store.list_all()

# tier_promote.py:91 — TierMetricsStore has NO 'record_promotion' method
store.record_promotion(ws_id, current.value, target.value)
```

**Root cause:** CLI written for future API that was never implemented.

### 2. Missing Generic Type Parameters (14 errors)

**Files:** `beads/base.py`, `beads/cli.py`, `beads/mock.py`, `traceability/models.py`, `validators/supersede_checker.py`, `validators/ws_completion.py`

```python
# Bad
def to_dict(self) -> dict:
    ...

# Good
def to_dict(self) -> dict[str, Any]:
    ...
```

### 3. Missing Type Annotations (3 errors)

**Files:** `traceability/service.py`, `cli/tier_promote.py`, `validators/capability_tier_checks_scope.py`

```python
# Bad
def check_scope(ws):
    ...

# Good
def check_scope(ws: Workstream) -> ValidationCheck:
    ...
```

### 4. Actual Bugs (3 errors)

```python
# beads/cli.py:192 — JSONDecodeError has no 'stderr' attribute
except json.JSONDecodeError as e:
    raise BeadsClientError(f"... {e.stderr}")  # Wrong!

# prd/parser_python.py:237 — Optional passed where str required
FlowStep(flow_name=maybe_none_value)  # Should handle None
```

## Affected Files

| File | Errors | Type |
|------|--------|------|
| `cli/tier_metrics.py` | 3 | Dead code |
| `cli/tier_promote.py` | 5 | Dead code + missing types |
| `beads/cli.py` | 3 | Bug + missing generics |
| `validators/ws_completion.py` | 4 | Missing generics |
| `validators/supersede_checker.py` | 3 | Missing generics |
| `traceability/models.py` | 3 | Missing generics |
| Others | 9 | Mixed |

## Action Plan

**Option A: Remove dead code (faster)**
- Delete `cli/tier_metrics.py`, `cli/tier_promote.py`
- Remove from `cli/main.py` imports
- Fix remaining 22 errors

**Option B: Implement missing API (fuller)**
- Add `history` to `TierMetrics`
- Add `list_all()`, `record_promotion()` to `TierMetricsStore`
- Fix all 30 errors

**Recommended:** Option A — dead code provides no value.

## Severity

- **P2 MEDIUM** — doesn't block runtime, but:
  - Prevents `mypy --strict` in CI
  - Masks real bugs
  - Accumulates debt

## Resolution (2026-01-31)

- **Option A applied:** Removed dead code `cli/tier_metrics.py`, `cli/tier_promote.py`
- Fixed 22 remaining mypy errors across 12 files
- Fixed beads/cli.py JSONDecodeError bug (e.stderr → e.msg)
- Fixed prd/parser_python.py flow_name Optional handling
- Branch: `bugfix/008-mypy-type-errors-dead-code`

**Note:** Project coverage 59% is pre-existing; 80% gate deferred to separate WS.

## Related

- Issue 007: NameError in workstream.py (fixed)
- WS 00-020-03: Triggered discovery of this debt
