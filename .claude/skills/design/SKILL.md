---
name: design
description: System design with progressive disclosure
tools: Read, Write, Bash, Glob, Grep, AskUserQuestion
version: 4.0.0
---

# @design - System Design with Progressive Disclosure

Multi-agent system design (Arch + Security + SRE) with progressive discovery blocks.

## When to Use

- After @idea requirements gathering
- Need architecture decisions
- Creating workstream breakdown

## Invocation

```bash
@design <task_id>
@design <task_id> --quiet    # Minimal design blocks
@design "feature description"  # Skip @idea, design directly
```

## Progressive Discovery Workflow

### Overview

**Discovery Blocks:** 3-5 focused blocks (not one big questionnaire)

**Block Structure:**
- Each block: 3 questions
- After each block: trigger point (Continue / Skip block / Done)
- User can skip blocks not relevant to feature

### Discovery Blocks

**Block 1: Data & Storage (3 questions)**
- Data models?
- Storage requirements?
- Persistence strategy?

**Block 2: API & Integration (3 questions)**
- API endpoints?
- External integrations?
- Authentication/authorization?

**Block 3: Architecture (3 questions)**
- Component structure?
- Layer boundaries?
- Error handling strategy?

**Block 4: Security (3 questions)**
- Input validation?
- Sensitive data handling?
- Rate limiting?

**Block 5: Operations (3 questions)**
- Monitoring?
- Deployment?
- Rollback strategy?

### After Each Block: Trigger Point

```markdown
AskUserQuestion({
  "questions": [
    {"question": "Block complete. Continue to next block?",
     "header": "Discovery",
     "options": [
       {"label": "Continue (Recommended)", "description": "Next discovery block"},
       {"label": "Skip block", "description": "Skip remaining blocks"},
       {"label": "Done", "description": "Generate workstreams with current info"}
     ]}
  ]
})
```

## Integration with @idea

```python
# Uses requirements from @idea
idea_result = load_idea_result(task_id)

# Skip already covered topics
skip_topics = idea_result.covered_topics

# Focus on design-specific questions
design_blocks = filter_blocks(skip_topics)
```

## --quiet Mode

Minimal blocks (2 blocks, 6 questions):
1. Data & Storage
2. Core Architecture

## Output

**Primary:** Workstream files in `docs/workstreams/backlog/`

**Secondary:**
- `docs/drafts/<task_id>-design.md` - Design document

## Next Steps

```bash
@build <ws_id>      # Execute workstream
@oneshot <feature>  # Execute all workstreams
```

---

**Version:** 4.0.0 (Progressive Disclosure)
**See Also:** `@idea`, `@build`, `@oneshot`
