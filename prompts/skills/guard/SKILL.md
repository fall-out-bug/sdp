---
name: guard
description: Pre-edit gate enforcing WS scope (INTERNAL)
tools:
  - Read
  - Shell
version: 2.0.0
---

# @guard - Pre-Edit Gate (INTERNAL)

**INTERNAL SKILL** — Called automatically before file edits.

## Purpose

1. Enforce that all edits happen within active WS scope
2. Block edits outside scope

## Check Flow

1. Is there an active WS? → No → BLOCK
2. Is file in WS scope? → No → BLOCK
3. Allow edit

## CLI Integration

```bash
# Activate WS (called by @build)
sdp guard activate 00-032-01

# Check file (called before edit)
sdp guard check internal/package/file.go

# Show current status
sdp guard status

# Deactivate when done
sdp guard deactivate
```

## Example Output

```bash
$ sdp guard activate 00-032-01
✓ Activated guard for WS 00-032-01
Scope files:
  - internal/package/*.go
  - internal/package/*_test.go

$ sdp guard check internal/package/file.go
✓ ALLOWED: File within WS scope

$ sdp guard check internal/other/parser.go
✗ BLOCKED: File not in scope
  Active WS: 00-032-01
  Scope: internal/package/*.go
```
