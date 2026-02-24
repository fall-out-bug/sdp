---
name: guard
description: Pre-edit gate enforcing WS scope (INTERNAL)
---

# @guard (INTERNAL)

Pre-edit gate. Called automatically before file edits. Enforce edits within active WS scope.

## Commands

```bash
sdp guard activate <ws-id>   # Set scope
sdp guard check <file>      # Verify file in scope
sdp guard status            # Show current
sdp guard deactivate        # Clear
```

## Flow

1. Active WS? No → BLOCK
2. File in scope? No → BLOCK
3. Allow edit

## Output

ALLOWED or BLOCKED with scope details.
