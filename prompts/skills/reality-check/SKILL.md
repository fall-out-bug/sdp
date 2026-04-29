---
name: reality-check
description: Quick documentation vs code reality validation.
---

# @reality-check

Quick validation that docs match code before making changes. ~90 seconds (vs 5-10 min for @verify-workstream).

## When to Use

- About to modify file based on documentation
- Unsure if docs reflect implementation
- Quick check before full verify-workstream

## Workflow

1. **Read actual code first** — Don't read docs first. Use Read tool on target file.
2. **Compare** — What does code actually do vs what we assumed?
3. **Report** — Match? Proceed. Mismatch? Stop or adapt.

## Output Format

```markdown
## Reality Check: <filename>
### What Code Actually Does: [summary]
### What We Assumed: [expectation]
### Recommendation: ✅ Proceed / ⚠️ Stop / 🔄 Adapt
```

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- @verify-workstream — Full workstream validation
- @build — Uses reality-check during execution
