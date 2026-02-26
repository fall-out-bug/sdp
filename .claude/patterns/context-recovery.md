# Context Recovery Pattern

How to recover original task after context compaction.

---

## Problem

After context compaction, agents tend to:
1. Continue the side task they were doing (tests, debugging)
2. Forget the PRIMARY TASK (roadmap execution, feature implementation)
3. Drift away from original goal

## Solution

**POST-COMPACTION CHECKLIST:**

### Step 1: Check Active Work
```bash
bd list --status=in_progress
bd ready
```

### Step 2: Identify PRIMARY vs SIDE Task

| PRIMARY TASK (continue this) | SIDE TASK (abandon this) |
|------------------------------|--------------------------|
| Roadmap execution | Fixing tests |
| Feature implementation | Improving coverage |
| Review process | Debugging edge case |
| PR creation | Refactoring |

### Step 3: Resume PRIMARY
- Close or de-prioritize side tasks
- Return to feature/roadmap work
- Check checkpoints: `ls .sdp/checkpoints/`

---

## Detection Rules

**You are doing a SIDE TASK if:**
- The work is not in your active workstream
- It started as "I'll just fix this quickly"
- The summary says "was improving X" but no feature ID

**You should return to PRIMARY if:**
- There's a feature ID in progress
- There's a checkpoint file
- `bd ready` shows work available

---

## Example

**Bad (drift):**
```
Summary: "Was improving test coverage for graph package"
Action: Continue writing tests for other packages
```

**Good (recovery):**
```
Summary: "Was improving test coverage for graph package"
Check: bd list --status=in_progress → shows workstream
Action: Resume workstream, tests were just side task
```

---

## Implementation

Add to skill CRITICAL RULES:
```markdown
5. **POST-COMPACTION RECOVERY** - After context compaction, check PRIMARY TASK first. Never drift to side tasks.
```

Add POST-COMPACTION PROTOCOL section with checklist.

---

## See Also

- `@build` - Execute workstream
- `@oneshot` - Execute all workstreams
- `.claude/patterns/tdd.md` - TDD pattern
