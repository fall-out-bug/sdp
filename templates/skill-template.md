---
name: {skill_name}
description: {One-line description, max 80 chars}
tools: {Comma-separated list of required tools}
max_lines: 100
---

# @{skill_name} - {Title}

{2-3 sentence purpose. What does this skill do and when to use it.}

## Quick Reference

| Step | Action | Gate |
|------|--------|------|
| 1 | {action} | {what must be true} |
| 2 | {action} | {what must be true} |
| 3 | {action} | {what must be true} |

## Workflow

### Step 1: {Name}

```bash
# Command or action
{command}
```

**Gate:** {What must be true before proceeding}

### Step 2: {Name}

{3-5 lines max}

### Step 3: {Name}

{3-5 lines max}

## Quality Gates

See [Quality Gates Reference](../docs/reference/quality-gates.md)

## Errors

| Error | Cause | Fix |
|-------|-------|-----|
| {error} | {cause} | {fix} |

## See Also

- [Full Specification](../docs/reference/{skill}-spec.md)
- [Examples](../docs/examples/{skill}/)
- [Related Skill](./{related}/SKILL.md)
