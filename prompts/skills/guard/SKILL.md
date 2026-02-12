---
name: guard
description: Pre-edit gate enforcing WS scope (INTERNAL)
tools: Read, Shell
version: 2.0.0
---

# @guard - Pre-Edit Gate (INTERNAL)

**INTERNAL SKILL** — Called automatically before file edits.

## Purpose

1. Enforce that all edits happen within active WS scope
2. Track and display review findings
3. Block progress if P0/P1 findings unresolved

## Check Flow

1. Is there an active WS? → No → BLOCK
2. Is file in WS scope? → No → BLOCK
3. Are there blocking findings? → Yes → WARN (but allow)
4. Allow edit

## CLI Integration

```bash
# Activate WS (called by @build)
sdp guard activate 00-032-01

# Check file (called before edit)
sdp guard check src/sdp/guard/skill.py

# Show current status (includes findings)
sdp guard status

# Deactivate when done
sdp guard deactivate
```

## Review Findings Integration

### Register Findings (called by @review)

```bash
# Register a finding from review
sdp guard finding add \
  --feature=F051 \
  --area=SRE \
  --title="Missing logging in memory.Store" \
  --priority=1 \
  --beads=sdp-abc123

# List all findings
sdp guard finding list
sdp guard finding list --all  # include resolved

# Resolve a finding
sdp guard finding resolve finding-123 --by="Fixed in commit abc123"

# Clear resolved findings
sdp guard finding clear
```

### Finding Priorities

| Priority | Name | Behavior |
|----------|------|----------|
| P0 | Critical | BLOCK - Must resolve immediately |
| P1 | High | WARN - Should resolve before merge |
| P2 | Medium | Track - Resolve when possible |
| P3 | Low | Track - Optional |

### Status Output with Findings

```bash
$ sdp guard status
Guard Status: ACTIVE
Active WS: 00-051-03
Scope files:
  - sdp-plugin/internal/memory/store.go
  - sdp-plugin/internal/memory/search.go

Review Findings: 2 open (1 blocking), 1 resolved

⚠️  BLOCKING FINDINGS (must resolve before merge):
  [SRE] P1 Add context.Context support
    → Beads: sdp-abc123
```

## Integration with @review

**Review agents MUST register findings:**

```bash
# After creating beads issue
bd create --title="SRE: Add logging" --type=task --priority=1

# Register in guard (enables blocking check)
sdp guard finding add \
  --feature=$FEATURE_ID \
  --area=SRE \
  --title="Add logging" \
  --priority=1 \
  --beads=$(bd list --search="Add logging" --format=id)
```

**@deploy checks for blocking findings:**

```bash
# Before merge, check for blockers
sdp guard status
if blocking > 0; then
  echo "Cannot deploy: unresolved P0/P1 findings"
  exit 1
fi
```

## Implementation

The guard system consists of:
- `GuardSkill` - Core logic for checking file permissions
- `GuardState` - State with scope files + review findings
- `ReviewFinding` - Finding model with priority/status
- CLI commands - User-facing commands for activation/checking/findings

## Usage in @build Skill

```python
# At start of @build
guard = GuardSkill(beads_client)
guard.activate(ws_id)

# Before each file edit
result = guard.check_edit(file_path)
if not result.allowed:
    raise PermissionError(result.reason)

# Check for blocking findings
if state.HasBlockingFindings():
    print("⚠️ Warning: There are unresolved P0/P1 findings")
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

$ sdp guard finding add --feature=F051 --area=SRE --title="Missing logging" --priority=1
✓ Registered finding: finding-1739123456
  Feature: F051
  Area: SRE
  Priority: P1

⚠️  BLOCKING: P0/P1 finding requires resolution before merge
```

## Version

**2.0.0** - Review Findings Integration
- Added `sdp guard finding` commands
- Track review findings with priority
- Block progress on P0/P1 findings
- Integration with @review skill

**1.0.0** - Initial implementation
- WS scope enforcement
- Context validation
- Branch safety checks
