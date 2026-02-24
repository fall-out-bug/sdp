---
name: prototype
description: Rapid prototyping shortcut for experienced vibecoders
---

# @prototype

Ultra-fast feature planning: 5-question interview → 1-3 workstreams → @oneshot with relaxed gates.

## When to Use

Experienced devs, need prototype fast, tech debt acceptable. Not for production or security-critical.

## Gate Overrides

| Gate | Normal | Prototype |
|------|--------|-----------|
| TDD | Required | Optional |
| Coverage | ≥80% | None |
| Architecture | Clean | Monolithic OK |

Non-negotiable: code compiles, runs, no crashes, basic security.

## Workflow

1. AskUserQuestion: problem, scope, dependencies, risks, success
2. Generate 1-3 workstreams
3. Launch @oneshot

## Output

`docs/drafts/prototype-{id}.md`, `docs/workstreams/backlog/00-FFF-*.md`

## See Also

- @feature — Full planning
- @oneshot — Execution
