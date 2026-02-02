---
name: guard
description: Pre-edit gate enforcing WS scope (INTERNAL)
tools: Read, Shell
---

# @guard - Pre-Edit Gate (INTERNAL)

**INTERNAL SKILL** — Called automatically before file edits.

## Purpose

Enforce that all edits happen within active WS scope.

## Check Flow

1. Is there an active WS? → No → BLOCK
2. Is file in WS scope? → No → BLOCK  
3. Allow edit

## CLI Integration

```bash
# Activate WS (called by @build)
sdp guard activate 00-032-01

# Check file (called before edit)
sdp guard check src/sdp/guard/skill.py

# Show current status
sdp guard status

# Deactivate when done
sdp guard deactivate
```

## Implementation

The guard system consists of:
- `GuardSkill` - Core logic for checking file permissions
- `GuardResult` - Result model with allow/deny decision
- CLI commands - User-facing commands for activation/checking

## Usage in @build Skill

```python
# At start of @build
guard = GuardSkill(beads_client)
guard.activate(ws_id)

# Before each file edit
result = guard.check_edit(file_path)
if not result.allowed:
    raise PermissionError(result.reason)
```

## Example Output

```bash
$ sdp guard activate 00-032-01
✓ Activated guard for WS 00-032-01
Scope files:
  - src/sdp/guard/skill.py
  - src/sdp/guard/state.py
  - tests/unit/test_guard.py

$ sdp guard check src/sdp/guard/skill.py
✓ ALLOWED: File within WS scope

$ sdp guard check src/sdp/core/parser.py
✗ BLOCKED: File not in scope
  Active WS: 00-032-01
  Scope: src/sdp/guard/*.py, tests/unit/test_guard.py

$ sdp guard status
Active WS: 00-032-01
Scope files: 3
Status: ENFORCING
```
