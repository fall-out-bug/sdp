---
ws_id: 00-013-07
feature: F013
dependencies: []
oneshot_ready: true
status: completed
completed: "2026-01-30"
---

## @build Skill

### Goal

Create workstream execution skill with TDD integration.

### Files

**Create:** `.claude/skills/build/SKILL.md` (~100 LOC)

### Steps

**Step 1: Define workflow**
- Read workstream
- Verify prerequisites
- Execute steps with /tdd
- Verify acceptance criteria
- Move to completed

**Step 2: Integrate /tdd**
```python
# For each step:
Skill("tdd")
# Follows Red → Green → Refactor
```

**Step 3: Progress reporting**
```markdown
✅ Step 1/5: Create module skeleton
✅ Step 2/5: Implement core class
```

### Acceptance Criteria

- [ ] SKILL.md defines workflow
- [ ] Calls /tdd internally
- [ ] Progress reporting format
- [ ] Prerequisites verification
