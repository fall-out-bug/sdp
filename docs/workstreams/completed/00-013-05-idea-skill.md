---
ws_id: 00-013-05
feature: F013
dependencies: []
oneshot_ready: true
status: completed
completed: "2026-01-30"
---

## @idea Skill Update

### Goal

Enhance @idea with product vision questions and schema validation.

### Files

**Modify:** `.claude/skills/idea/SKILL.md` (+70 LOC)

### Steps

**Step 1: Add vision interview**
- Mission question
- PRODUCT_VISION.md alignment question

**Step 2: Add intent.json generation**
```json
{
  "problem": "...",
  "users": [...],
  "success_criteria": [...]
}
```

**Step 3: Add schema validation**
```python
from sdp.schema.validator import IntentValidator
validator.validate_file("docs/intent/{slug}.json")
```

### Acceptance Criteria

- [ ] Vision interview questions added
- [ ] intent.json generation documented
- [ ] Schema validation step added
- [ ] Links to @design as next step
