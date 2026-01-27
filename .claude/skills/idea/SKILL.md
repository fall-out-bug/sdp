---
name: idea
description: Interactive requirements gathering through deep interviewing using AskUserQuestion. Creates feature draft with goals, scope, and open questions.
tools: Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion
---

# /idea - Requirements Gathering

Deep, interactive interviewing to capture comprehensive feature requirements using AskUserQuestion tool.

## When to Use

- Starting new feature
- Unclear requirements
- Need comprehensive requirements document
- Want to explore tradeoffs and edge cases

## Workflow

**IMPORTANT:** Use AskUserQuestion for deep, continuous interviewing until requirements are complete.

### When to Call /think

If requirements are **ambiguous, complex, or have conflicting needs**, call `/think` first:

```python
Skill("think")
# Returns structured analysis before interviewing
```

Use /think for:
- Multiple user types with conflicting needs
- Technical approach unclear
- Success metrics debatable
- Significant unknowns

### Step 1: Read Product Vision (If Exists)

Check for PRODUCT_VISION.md to align with project goals:

```bash
Read(PRODUCT_VISION.md)  # if exists
```

### Step 2: Vision Interview (NEW)

**Ask product vision questions using AskUserQuestion:**

```markdown
AskUserQuestion({
  "questions": [{
    "question": "What is the core mission of this feature?",
    "header": "Mission",
    "options": [
      {"label": "Solve specific pain point", "description": "Addresses clear user frustration"},
      {"label": "Enable new capability", "description": "Unlocks something previously impossible"},
      {"label": "Improve efficiency", "description": "Makes existing process faster/cheaper"}
    ]
  }, {
    "question": "How does this feature align with PRODUCT_VISION.md?",
    "header": "Alignment",
    "options": [
      {"label": "Directly supports mission", "description": "Core to product vision"},
      {"label": "Enables mission", "description": "Supporting capability"},
      {"label": "Extends mission", "description": "Natural evolution"},
      {"label": "New direction", "description": "May need vision update"}
    ]
  }]
})
```

### Step 3: Technical Interview

**Ask foundational questions:**

```markdown
AskUserQuestion({
  "questions": [{
    "question": "What is the primary problem this feature solves?",
    "header": "Problem",
    "options": [
      {"label": "User pain point", "description": "Addresses frustration or inefficiency"},
      {"label": "Business requirement", "description": "Enables new revenue or reduces cost"},
      {"label": "Technical debt", "description": "Improves maintainability or performance"},
      {"label": "Competitive parity", "description": "Matches competitor capabilities"}
    ]
  }, {
    "question": "Who are the primary users of this feature?",
    "header": "Users",
    "options": [
      {"label": "End users", "description": "Direct product users"},
      {"label": "Administrators", "description": "System managers and ops teams"},
      {"label": "Developers", "description": "Engineering team integration"},
      {"label": "API consumers", "description": "External integrations"}
    ],
    "multiSelect": true
  }]
})
```

### Step 4: Deep Dive Interview

Continue with progressively detailed questions about:
- Technical Implementation (storage, failure modes, auth)
- UI & UX (location, discoverability)
- Concerns & Tradeoffs (performance, complexity, security)

### Step 5: Create Draft

Write comprehensive spec to `docs/drafts/idea-{slug}.md`

### Step 6: Generate intent.json (NEW)

Create machine-readable intent at `docs/intent/{slug}.json`:

```json
{
  "problem": "Detailed problem statement (50+ chars)",
  "users": ["end_users", "admins"],
  "success_criteria": [
    {"criterion": "Success rate", "measurement": ">95%"},
    {"criterion": "Latency", "measurement": "<500ms p95"}
  ],
  "tradeoffs": {
    "security": "prioritize",
    "performance": "accept",
    "complexity": "accept",
    "time_to_market": "prioritize"
  },
  "technical_approach": {
    "architecture": "monolith",
    "storage": "relational_db",
    "failure_mode": "graceful_degradation",
    "auth_method": "jwt"
  }
}
```

### Step 7: Validate Intent (NEW)

Validate against schema:

```python
from sdp.schema.validator import IntentValidator

validator = IntentValidator()
validator.validate_file("docs/intent/{slug}.json")
```

## Output

- `docs/drafts/idea-{slug}.md` — Comprehensive spec
- `docs/intent/{slug}.json` — Machine-readable intent (validated)

## Next Step

`/design idea-{slug}` — Decompose into workstreams with full context
