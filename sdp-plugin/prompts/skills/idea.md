---
name: idea
description: Interactive requirements gathering using Beads with deep interviewing
tools: Read, Write, Bash, Glob, Grep, AskUserQuestion
version: 3.0.0
---

# @idea - Requirements Gathering

Deep interviewing to capture comprehensive feature requirements. Creates Beads task with PRODUCT_VISION alignment.

## When to Use

- Starting new feature
- Unclear requirements
- Need comprehensive spec with tradeoffs explored

## Invocation

```bash
@idea "feature description"
@idea "feature description" --spec path/to/SPEC.md  # with existing spec
```

## Workflow

### Step 1: Read Context

```bash
Read(PRODUCT_VISION.md)  # Align with project goals
Glob("docs/specs/**/*")   # Similar features
```

### Step 2: Vision Interview (AskUserQuestion)

Ask product-level questions first:

```markdown
AskUserQuestion({
  "questions": [
    {"question": "What is the core mission of this feature?", "header": "Mission",
     "options": [
       {"label": "Solve pain point", "description": "Addresses user frustration"},
       {"label": "Enable capability", "description": "Unlocks new possibility"},
       {"label": "Improve efficiency", "description": "Faster/cheaper process"}
     ]},
    {"question": "How does this align with PRODUCT_VISION.md?", "header": "Alignment",
     "options": [
       {"label": "Core to mission", "description": "Directly supports vision"},
       {"label": "Enables mission", "description": "Supporting capability"},
       {"label": "New direction", "description": "May need vision update"}
     ]}
  ]
})
```

### Step 3: Deep Dive Interview

Continue with technical questions until complete:

**Topics to cover:**
- Problem & users
- Technical approach (storage, failure modes)
- UI/UX specifics
- Security & performance
- Testing strategy
- Edge cases & tradeoffs

**DON'T ask obvious questions.** Instead ask about:
- Ambiguities and hidden assumptions
- Tradeoffs between approaches
- Failure modes and edge cases

### Step 4: Create Beads Task

```python
from sdp.beads import create_beads_client, BeadsTaskCreate, BeadsPriority

client = create_beads_client()

task = client.create_task(BeadsTaskCreate(
    title=feature_title,
    description=f"""## Context & Problem
{problem_answer}

## Goals & Non-Goals
{goals_answer}

## Technical Approach
{technical_answer}

## Concerns & Tradeoffs
{concerns_answer}
""",
    priority=BeadsPriority.HIGH,
    sdp_metadata={
        "feature_type": "idea",
        "mission": mission_answer,
        "product_vision_alignment": alignment_answer,
    }
))

print(f"✅ Created: {task.id}")
```

### Step 5: Create Intent File

```python
import json

with open(f"docs/intent/{task.id}.json", "w") as f:
    json.dump({
        "task_id": task.id,
        "title": task.title,
        "mission": mission_answer,
        "alignment": alignment_answer,
        "interview_answers": all_answers,
    }, f, indent=2)
```

## Output

**Primary:** Beads task ID (e.g., `bd-0001`)

**Secondary:**
- `docs/intent/{task_id}.json` — Machine-readable intent
- Optional: `docs/drafts/beads-{task_id}.md` — Human-readable export

## Next Steps

```bash
@design bd-0001      # Decompose into workstreams
bd show bd-0001      # View task details
bd ready             # Check ready tasks
```

## Key Principles

1. **Start broad, go deep** — foundational questions first
2. **Product vision first** — align before technical details
3. **No obvious questions** — explore tradeoffs, not yes/no
4. **Continue until complete** — no ambiguities remain

## Example Session

```bash
@idea "Add user authentication"

# Vision interview...
# Technical interview...

✅ Created Beads task: bd-0001
   Title: Add user authentication
   Priority: HIGH
   Intent: docs/intent/bd-0001.json

# Next:
@design bd-0001
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@idea "feature"` | Create task with requirements |
| `@think "complex"` | Deep analysis before @idea |
| `bd show {id}` | View task details |
| `@design {id}` | Decompose into workstreams |

---

**Version:** 3.0.0  
**See Also:** `@design`, `@build`, `@think`
