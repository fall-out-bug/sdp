---
name: design
description: Decompose feature into workstreams with scope
tools: Read, Write, Shell, AskUserQuestion
---

# @design - Feature Decomposition

Analyze requirements and create workstreams with dependencies and scope.

## Quick Reference

| Step | Action | Gate |
|------|--------|------|
| 1 | Read feature | Requirements clear |
| 2 | Explore codebase | Context gathered |
| 3 | Ask architecture | Decisions made |
| 4 | Create workstreams | All WS have AC + scope |
| 5 | Verify deps | No cycles |

## Workflow

### Step 1: Read Feature

```bash
bd show {feature-id}
```

Or for markdown:

```bash
Read("docs/drafts/{feature}.md")
```

### Step 2: Explore Codebase

```bash
Glob("src/**/*.py")
Grep("relevant patterns")
```

Understand:
- Existing architecture
- Integration points
- Dependencies

### Step 3: Architecture Questions

Use AskUserQuestion for:
- Complexity level (simple/medium/large)
- Layers needed (domain/repo/service/api)
- Database changes
- External integrations

### Step 4: Create Workstreams

For each WS:

```yaml
ws_id: 00-032-01
title: Domain entities
size: MEDIUM
depends_on: []
scope_files:
  - src/domain/entities.py
  - src/domain/value_objects.py
  - tests/unit/test_entities.py
acceptance_criteria:
  - AC1: Entity created with required fields
  - AC2: Value objects immutable
```

**Key:** Include `scope_files` for guard enforcement.

### Step 5: Verify Dependencies

```bash
sdp ws graph {feature-id}
```

Check for cycles. Ensure topological order possible.

### Step 6: Beads Registration (when Beads enabled)

When project uses Beads (bd installed, `.beads/` exists):

```bash
poetry run sdp beads migrate docs/workstreams/backlog/ --real
```

Creates Beads tasks and `.beads-sdp-mapping.jsonl` entries. Agents can then use `bd ready` and `@build` integrates with Beads.

**When Beads NOT enabled:** Skip. Workstreams remain in markdown only.

## Quality Gates

See [Quality Gates Reference](../../docs/reference/quality-gates.md)

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| Cycle detected | Circular deps | Break cycle |
| Missing scope | No files listed | Add scope_files |
| Too large | WS >500 LOC | Split WS |

## See Also

- [Full Design Spec](../../docs/reference/design-spec.md)
- [Sizing Guide](../../docs/reference/ws-sizing.md)
