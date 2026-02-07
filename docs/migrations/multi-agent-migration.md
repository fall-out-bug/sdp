# Multi-Agent SDP Migration Guide

Migrate from old SDP (single-agent) to new SDP (multi-agent with progressive disclosure).

## Before/After Comparison

### Old Workflow (v3.x)
```
@feature "Add OAuth"
→ 15-25 questions asked upfront
→ No structure, no breaks
→ User fatigue
```

### New Workflow (v4.x)
```
@idea "Add OAuth"
→ 3-question cycles with trigger points
→ User controls depth
→ Average 18 questions (vs 25)

@design task-id
→ Discovery blocks
→ Skip irrelevant blocks
→ Faster design phase
```

## Migration Steps

### Step 1: Update Skills
```bash
# Skills are already updated in .claude/skills/
# @idea: v4.0.0 (progressive disclosure)
# @design: v4.0.0 (progressive disclosure)
```

### Step 2: Learn New Commands
```bash
# Strategic level (new)
@vision "product idea"      # 7 expert agents
@reality --quick           # 8 expert agents

# Feature level (updated)
@idea "feature"            # Progressive disclosure
@design task-id            # Discovery blocks

# Execution (same)
@build ws-id               # Execute workstream
@oneshot feature-id        # Execute all workstreams
```

### Step 3: Update Workflows

**Old:**
```bash
@feature "Add OAuth"       # One big questionnaire
@build 00-001-01
```

**New:**
```bash
@idea "Add OAuth"          # 3-question cycles
@design task-id            # Discovery blocks
@oneshot F001              # Parallel execution
```

## Pattern Mapping

| Old Pattern | New Pattern |
|------------|------------|
| `@feature` | `@idea` + `@design` |
| Deep dive upfront | Progressive cycles/blocks |
| Unbounded questions | 12-27 questions target |
| Single agent | Multi-agent synthesis |

## Rollback Procedure

If issues arise:
```bash
# Revert to old skills
git checkout v3.1.0 -- .claude/skills/idea/SKILL.md
git checkout v3.1.0 -- .claude/skills/design/SKILL.md

# Use old commands
@feature "feature description"
```

## Testing Checklist

- [ ] Try @idea with progressive disclosure
- [ ] Try @design with discovery blocks
- [ ] Verify trigger points work
- [ ] Check question count (12-27 range)
- [ ] Test --quiet mode

## Common Patterns

### Pattern 1: Quick Feature
**Old:** `@feature "add button"` (25 questions)
**New:** `@idea "add button" --quiet` (5 questions)

### Pattern 2: Complex Feature
**Old:** `@feature "OAuth integration"` (30+ questions)
**New:**
```bash
@idea "OAuth integration"      # 18 questions average
@design task-id                # Skip irrelevant blocks
@oneshot F001                  # Parallel execution
```

### Pattern 3: New Product
**Old:** Manual planning
**New:**
```bash
@vision "AI task manager"      # 7 expert agents
@reality --quick               # 8 expert agents
@feature "first feature"       # Plan based on analysis
```

---

**Version:** 4.0.0
**Migrating From:** 3.x
**Date:** 2026-02-07
