# 007: sdp guard activate fails with NameError in workstream.py

**Source:** /issue (2026-01-31)  
**Status:** Fixed  
**Priority:** P1 (HIGH)  
**Route:** /bugfix (applied)

## Symptom

`sdp guard activate <ws_id>` fails immediately with:

```
NameError: name 'Command' is not defined
```

**Traceback:**
```
File ".../src/sdp/cli/main.py", line 38, in <module>
    from sdp.cli.workstream import workstream
File ".../src/sdp/cli/workstream.py", line 18, in <module>
    validate_tier: Command | None = None
                   ^^^^^^^
NameError: name 'Command' is not defined
```

Any `sdp` CLI invocation fails because `main.py` imports `workstream`, which raises at import time.

## Root Cause

**Files:** `src/sdp/cli/workstream.py`, `src/sdp/cli/main.py`

`Command` and `Group` are imported only under `TYPE_CHECKING` but used in runtime type annotations:

```python
if TYPE_CHECKING:
    from click import Command, Group

validate_tier: Command | None = None  # ← NameError at runtime
guard: Group | None = None           # ← Same pattern in main.py
```

Type annotations are evaluated at runtime. Since these names are not imported at runtime, the assignments raise `NameError`.

## Impact Chain

1. **Blocking:** All `sdp` CLI commands fail (guard, build, ws, etc.)
2. **Workflow:** Cannot run `sdp guard activate` before editing
3. **Post-build:** Separate concern — post-build reports project-wide mypy issues; those are outside WS 00-020-03 scope (tests/coverage complete)

## Severity Classification

| Criterion   | Assessment                          |
|------------|--------------------------------------|
| Production | No — dev workflow only               |
| Scope      | All CLI commands blocked             |
| Workaround | None — SDP unusable without fix      |
| **Verdict**| **P1 HIGH** → Route to `/bugfix`     |

## Fix

Use `click.Command` and `click.Group` in annotations (click is already imported):

```python
validate_tier: click.Command | None = None   # workstream.py
guard: click.Group | None = None             # main.py
```

**Applied:** Replaced `Command`/`Group` with `click.Command`/`click.Group` in both files; removed `TYPE_CHECKING` block.

## Related

- WS 00-020-03: Tests and coverage complete; blocked by this bug
- Post-build mypy: Project-wide issues; separate WS/backlog item
