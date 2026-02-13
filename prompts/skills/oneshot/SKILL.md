---
name: oneshot
description: Autonomous multi-agent execution with review-fix loop and PR creation
cli: sdp orchestrate
version: 6.0.0
---

# oneshot

> **CLI:** `sdp orchestrate <feature-id>`

Autonomous feature execution with review-fix loop and PR creation.

---

## EXECUTE THIS NOW

When user invokes `@oneshot F067`, run the CLI command:

```bash
sdp orchestrate F067
```

This is not an example. Execute this command.

---

## Quick Start

```bash
sdp orchestrate F067              # Execute all workstreams
sdp orchestrate resume F067       # Resume from checkpoint
sdp orchestrate --retry 3 F067    # Allow 3 retries per WS
```

---

## What Happens

```
Phase 1: Execute Workstreams (CLI handles this)
    └─ Loads docs/workstreams/backlog/00-067-*.md
    └─ Builds dependency graph
    └─ Executes in topological order
    └─ Checkpoints after each WS

Phase 2: Review-Fix Loop
    └─ Run @review after all WS complete
    └─ Fix P0/P1 findings
    └─ Repeat until approved (max 3)

Phase 3: Verify Clean
    └─ sdp guard finding list (0 blocking)

Phase 4: Create PR
    └─ Push to feature branch
    └─ PR to dev (NOT main)
```

---

## Finding Priority

| Priority | Action |
|----------|--------|
| P0 | Fix immediately, commit |
| P1 | Create bugfix with `bd create` |
| P2+ | Track only, don't block |

---

## Resume After Interruption

```bash
# Check checkpoint
cat .sdp/checkpoints/F067-*.json

# Resume execution
sdp orchestrate resume F067
```

---

## Output

```
🚀 Orchestrating feature F067
   Workstream dir: docs/workstreams/backlog
   Checkpoint dir: .sdp/checkpoints

→ [14:30] Executing 00-067-01...
→ [14:35] ✅ 00-067-01 complete
...
✅ Feature F067 completed successfully
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Command not found | `go install ./sdp-plugin/cmd/sdp` |
| Checkpoint corrupted | Delete `.sdp/checkpoints/F067-*.json` |
| WS blocked | Check dependencies in WS frontmatter |

---

## See Also

- `@build` - Single workstream
- `@review` - Quality review
- `.claude/patterns/tdd.md` - TDD pattern

**Implementation:** `sdp-plugin/cmd/sdp/orchestrate.go`
