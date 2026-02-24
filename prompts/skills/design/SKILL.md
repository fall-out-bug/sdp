---
name: design
description: System design with progressive disclosure
---

# @design

Multi-agent design (Arch + Security + SRE) with progressive discovery blocks.

## When to Use

After @idea. Need architecture decisions. Creating workstream breakdown.

## Workflow

1. **Load requirements** — `docs/intent/{task_id}.json`, `docs/drafts/idea-*.md`
2. **Progressive discovery** — 3-5 blocks, 3 questions each. Blocks: Data & Storage, API & Integration, Architecture, Security, Operations. After each: Continue / Skip / Done
3. **Generate workstreams** — `docs/workstreams/backlog/00-FFF-SS.md`
4. **Create Beads** — `bd create --title="WS-FFF-01: {title}" --type=task --priority=2`

## Modes

| Mode | Blocks |
|------|--------|
| Default | 3-5 |
| --quiet | 2 (Data + Architecture) |

## Output

Workstream files. `docs/drafts/<task_id>-design.md`.

## See Also

- @idea — Requirements
- @build — Execute workstream
- @oneshot — Execute all
