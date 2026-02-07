---
name: idea
description: Interactive requirements gathering with progressive disclosure
tools: Read, Write, Bash, Glob, Grep, AskUserQuestion
version: 4.0.0
---

# @idea - Requirements Gathering with Progressive Disclosure

Deep interviewing to capture comprehensive feature requirements using progressive disclosure (3-question cycles). Creates markdown spec, optionally creates Beads task.

## When to Use

- Starting new feature
- Unclear requirements
- Need comprehensive spec with tradeoffs explored

## Invocation

```bash
@idea "feature description"
@idea "feature description" --quiet    # Minimal questions (core only)
@idea "feature description" --spec path/to/SPEC.md  # with existing spec
```

## Progressive Disclosure Workflow

### Overview

**Question Target:**
- Minimum: 12 questions (bounded exploration)
- Maximum: 27 questions (deep analysis)
- Average: 18-20 questions per feature

**3-Question Cycles:**
1. Ask 3 focused questions
2. Offer trigger point after each cycle
3. User chooses: continue / deep design / skip to @design

### Step 1: Read Context

```bash
Read(PRODUCT_VISION.md)  # Align with project goals
Glob("docs/specs/**/*")   # Similar features
```

### Step 2: Vision Interview (Cycle 1 - 3 Questions)

AskUserQuestion with trigger point:

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
     ]},
    {"question": "Who are the primary users?", "header": "Users",
     "options": [
       {"label": "End users", "description": "Direct consumers"},
       {"label": "Developers", "description": "Internal tooling"},
       {"label": "Admins", "description": "Management/ops"}
     ]}
  ],
  "multiSelect": false
})

# TRIGGER POINT (after 3 questions)
AskUserQuestion({
  "questions": [
    {"question": "Continue exploring requirements? (3-question cycle complete)",
     "header": "Depth",
     "multiSelect": false,
     "options": [
       {"label": "Continue (Recommended)", "description": "More questions to understand requirements"},
       {"label": "Deep design", "description": "Jump to @design with detailed architectural exploration"},
       {"label": "Skip to @design", "description": "Move to workstream decomposition with current info"}
     ]}
  ]
})
```

### Step 3: Technical Interview (Cycles 2-5)

Continue 3-question cycles based on user choice:

**Cycle 2 - Problem & Users (3 questions):**
- What problem does this solve?
- What are the user pain points?
- What happens if we don't build this?

**Cycle 3 - Technical Approach (3 questions):**
- Storage/data requirements?
- Failure modes to handle?
- Integration points?

**Cycle 4 - UI/UX & Quality (3 questions):**
- UI/UX requirements?
- Performance targets?
- Security considerations?

**Cycle 5 - Testing & Edge Cases (3 questions):**
- Testing strategy?
- Edge cases to handle?
- Success metrics?

After each cycle, offer trigger point (Continue / Deep design / Skip).

### Step 4: TMI Detection

If user provides extensive detail upfront (indicators):

```markdown
# TMI Indicators:
- "detailed spec", "full implementation", "complete architecture"
- "here's the full code", "let me describe everything"
- User writes >500 characters in initial prompt

# Response:
"I notice you've provided detailed context. Would you like me to:
- Continue with targeted questions (recommended)
- Skip to @design with your detailed spec
- Use --quiet mode for minimal questions"
```

### Step 5: Create Beads Task

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
        "question_count": len(interview_answers),
    }
))

print(f"✅ Created: {task.id}")
```

## --quiet Mode

Minimal questions (3-5 core questions only):
1. Mission?
2. Users?
3. Core requirement?

Skip deep-dive cycles, move directly to @design.

## Output

**Primary:** Beads task ID (e.g., `bd-0001`)

**Secondary:**
- `docs/intent/{task_id}.json` — Machine-readable intent
- Question count included in metadata

## Next Steps

```bash
@design bd-0001      # Decompose into workstreams
bd show bd-0001      # View task details
bd ready             # Check ready tasks
```

## Key Principles

1. **Progressive disclosure** — 3 questions at a time
2. **User-controlled depth** — trigger points after each cycle
3. **Respect brevity** --quiet mode for experienced users
4. **No obvious questions** — explore tradeoffs, not yes/no
5. **TMI detection** — offer shortcuts when user over-explains

## Example Session

```bash
@idea "Add user authentication"

# Cycle 1: Vision (3 questions)
# [Mission] What is the core mission?
# [Alignment] How does this align with vision?
# [Users] Who are the primary users?

# TRIGGER: Continue? (yes/deep design/skip)
# User selects: Continue

# Cycle 2: Problem (3 questions)
# ...

# TRIGGER: Continue? (yes/deep design/skip)
# User selects: Deep design

# Jump to @design with architectural exploration

✅ Created Beads task: bd-0001
   Title: Add user authentication
   Questions asked: 6
   Priority: HIGH

# Next:
@design bd-0001
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@idea "feature"` | Create task with progressive interview |
| `@idea "feature" --quiet` | Minimal questions (3-5 core only) |
| `@think "complex"` | Deep analysis before @idea |
| `bd show {id}` | View task details |
| `@design {id}` | Decompose into workstreams |

---

**Version:** 4.0.0 (Progressive Disclosure)
**See Also:** `@design`, `@build`, `@think`
