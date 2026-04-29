---
name: verify-workstream
description: Validate workstream documentation against codebase reality.
---

# @verify-workstream

Before @build or @oneshot, validate docs match codebase.

## Workflow

1. **Read WS** — Parse frontmatter: goal, scope_files, acceptance_criteria
2. **Find files** — Glob to locate scope files. Gate: all must exist
3. **Read implementation** — Parse structure, identify patterns
4. **Compare** — Table: Documentation | Reality | Status
5. **Recommend** — PAUSE (high mismatch) / PROCEED / PROCEED WITH ADAPTATIONS

## Output

Verification complete. Severity. Recommendation. Comparison table.

## Integration

@build invokes this before execution.

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- @reality-check — Quick single-file check
- @build — Auto-runs verification
