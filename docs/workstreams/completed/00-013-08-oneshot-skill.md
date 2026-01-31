---
ws_id: 00-013-08
feature: F013
dependencies: []
oneshot_ready: true
status: completed
completed: "2026-01-30"
---

## @oneshot Skill Update

### Goal

Enhance @oneshot with checkpoint/resume and execution graph reading.

### Files

**Modify:** `.claude/skills/oneshot/SKILL.md` (+70 LOC)

### Steps

**Step 1: Add execution graph reading**
```python
from sdp.design.graph import DependencyGraph
graph = DependencyGraph()
execution_order = graph.topological_sort()
```

**Step 2: Add checkpoint format**
```json
{
  "feature": "F013",
  "completed_ws": [...],
  "current_ws": "WS-013-03",
  "execution_order": [...]
}
```

**Step 3: Add resume capability**
```bash
/oneshot F013 --resume {agent_id}
```

### Acceptance Criteria

- [ ] Reads DependencyGraph
- [ ] Checkpoint after each WS
- [ ] Resume from checkpoint
- [ ] Two-stage review
