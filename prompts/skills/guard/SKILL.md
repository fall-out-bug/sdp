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

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |
