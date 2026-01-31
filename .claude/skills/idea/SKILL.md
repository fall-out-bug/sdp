---
name: idea
description: Interactive requirements gathering using Beads for task storage. Creates Beads task with comprehensive requirements from deep interviewing.
tools: Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion
version: 2.1.0-beads-ai-comm
---

# @idea - Requirements Gathering (Beads + AI-Comm Integration)

Deep, interactive interviewing to capture comprehensive feature requirements using AskUserQuestion tool. **Outputs Beads task with PRODUCT_VISION alignment.**

## When to Use

- Starting new feature (Beads-first workflow)
- Unclear requirements
- Need comprehensive requirements document
- Want to explore tradeoffs and edge cases

## Beads vs Markdown Workflow

**This skill creates Beads tasks** (hash-based IDs, multi-agent ready) with enhanced metadata.

For traditional markdown workflow, use `prompts/commands/idea.md` instead.

## Invocation

```bash
@idea "feature description"
# or with existing spec
@idea "feature description" --spec path/to/SPEC.md
```

**Environment Variables:**
- `BEADS_USE_MOCK=true` - Use mock Beads (default for projects without Beads)
- `BEADS_USE_MOCK=false` - Use real Beads CLI (requires Go + bd installed)

**SDP repo:** Beads is always enabled. Use real Beads (bd installed, .beads/ exists).

## Workflow

**IMPORTANT:** Use AskUserQuestion for deep, continuous interviewing until requirements are complete.

### When to Call /think

If requirements are **ambiguous, complex, or have conflicting needs**, call `@think` first:

```python
Skill("think")
# Returns structured analysis before interviewing
```

Use @think for:
- Multiple user types with conflicting needs
- Technical approach unclear
- Success metrics debatable
- Significant unknowns

### Step 0: Initialize Beads Client

```python
from sdp.beads import create_beads_client, BeadsTaskCreate, BeadsPriority
import os

# Detect which client to use
use_mock = os.getenv("BEADS_USE_MOCK", "true").lower() == "true"
client = create_beads_client(use_mock=use_mock)

print(f"Using {'Mock' if use_mock else 'Real'} Beads client")
```

### Step 1: Read Product Vision (If Exists)

Check for PRODUCT_VISION.md to align with project goals:

```bash
Read(PRODUCT_VISION.md)  # if exists
```

### Step 2: Read Context (If Exists)

If user provides existing spec or similar features exist:

```bash
# Read existing spec
@path/to/SPEC.md

# Search for similar features
Glob("docs/specs/**/*")
Grep("similar feature keywords")
```

### Step 3: Vision Interview (NEW from ai-comm)

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

### Step 4: Initial Interview

**Ask foundational questions using AskUserQuestion:**

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
    ],
    "multiSelect": false
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

### Step 5: Deep Dive Interview

**Continue with progressively detailed questions. Be VERY in-depth:**

**Technical Implementation:**
```markdown
AskUserQuestion({
  "questions": [{
    "question": "How should data be persisted for this feature?",
    "header": "Data Storage",
    "options": [
      {"label": "Relational DB (Recommended)", "description": "ACID guarantees, complex queries. Adds DB dependency."},
      {"label": "NoSQL", "description": "Flexible schema, horizontal scaling. Eventual consistency."},
      {"label": "In-memory cache", "description": "Fast, ephemeral. Data loss on restart."},
      {"label": "File system", "description": "Simple, no DB. Limited query capabilities."}
    ],
    "multiSelect": false
  }, {
    "question": "What happens if this feature's service is unavailable?",
    "header": "Failure Mode",
    "options": [
      {"label": "Graceful degradation", "description": "Fallback to limited functionality"},
      {"label": "Fail closed", "description": "Block operation entirely"},
      {"label": "Queue and retry", "description": "Async processing with eventual completion"},
      {"label": "Best-effort", "description": "Continue with warnings"}
    ],
    "multiSelect": false
  }]
})
```

**Continue interviewing** until you have comprehensive answers about:
- Technical implementation details
- UI/UX specifics
- Error handling approach
- Performance considerations
- Security implications
- Testing strategy
- Deployment approach
- Monitoring and observability

**DON'T ask obvious questions.** Instead, ask about:
- Ambiguities in requirements
- Hidden assumptions
- Edge cases
- Failure modes
- Integration points
- Tradeoffs between approaches

### Step 6: Create Beads Task with Enhanced Metadata (MERGED)

After interviewing is complete, create Beads task with ai-comm metadata:

```python
# Build comprehensive description from interview
description = f"""## Context & Problem

{problem_answer}

## Goals & Non-Goals

**Goals:**
{goals_answer}

**Non-Goals:**
{nongoals_answer}

## Primary Users

{users_answer}

## Product Alignment (NEW)

**Mission:** {mission_answer}
**Vision Alignment:** {alignment_answer}

## Technical Approach

**Data Storage:** {storage_answer}
**Failure Mode:** {failure_answer}

## UI/UX Design

{ui_answer}

## Concerns & Risks

{concerns_answer}

## Tradeoffs

{tradeoffs_answer}

## Open Questions

{open_questions}
"""

# Determine priority from scope
scope_priority = {
    "Critical path": BeadsPriority.CRITICAL,
    "Important but not urgent": BeadsPriority.HIGH,
    "Nice to have": BeadsPriority.MEDIUM,
    "Backlog": BeadsPriority.BACKLOG
}

# Create Beads task with enhanced metadata
task = client.create_task(BeadsTaskCreate(
    title=feature_title,
    description=description,
    priority=scope_priority.get(scope_answer, BeadsPriority.MEDIUM),
    sdp_metadata={
        "feature_type": "idea",
        "interview_answers": all_answers,
        "created_by": "@idea skill",
        "product_vision_alignment": alignment_answer,  # NEW
        "mission": mission_answer,  # NEW
        "think_analysis": think_result if think_called else None,  # NEW
    }
))

print(f"✅ Created Beads task: {task.id}")
print(f"   Title: {task.title}")
print(f"   Status: {task.status}")
print(f"   Priority: {task.priority}")
```

### Step 7: Optional Markdown Export + Intent File (NEW from ai-comm)

For git history and human readability, export to markdown:

```python
markdown_path = f"docs/drafts/beads-{task.id}.md"

with open(markdown_path, "w") as f:
    f.write(f"""# {task.title}

> **Beads Task ID:** {task.id}
> **Created:** {datetime.utcnow().isoformat()}
> **Priority:** {task.priority.value}

{description}

---

## Next Steps

1. Run `@design {task.id}` to decompose into workstreams
2. Or run `bd show {task.id}` to view in Beads CLI
""")
```

**Create machine-readable intent file (NEW):**
```python
import json

intent_path = f"docs/intent/{task.id}.json"

with open(intent_path, "w") as f:
    json.dump({
        "task_id": task.id,
        "title": task.title,
        "mission": mission_answer,
        "alignment": alignment_answer,
        "priority": task.priority.value,
        "created_at": datetime.utcnow().isoformat(),
        "interview_answers": all_answers,
    }, f, indent=2)
```

## Output

**Primary:** Beads task ID (e.g., `bd-0001`)

**Secondary:**
- Optional markdown export to `docs/drafts/beads-{task_id}.md`
- Intent file at `docs/intent/{task_id}.json` (NEW)

**Beads Task Fields:**
- `id`: Hash-based task ID (auto-generated)
- `title`: Feature title from input
- `description`: Comprehensive requirements from interview
- `status`: OPEN (default)
- `priority`: CRITICAL/HIGH/MEDIUM/LOW/BACKLOG
- `sdp_metadata`: Interview answers, product alignment, think analysis (enhanced)

## Next Steps

After creating idea task:

1. **Decompose into workstreams:**
   ```bash
   @design bd-0001
   ```

2. **View in Beads CLI:**
   ```bash
   bd show bd-0001
   bd ready  # Check if ready to work on
   ```

3. **Start execution:**
   ```bash
   @build bd-0001.1  # First workstream sub-task
   ```

## Key Principles

**Interviewing Strategy:**
1. **Start broad, go deep** — foundational questions first, then drill into details
2. **Product vision first** — align with PRODUCT_VISION.md before technical details (NEW)
3. **No obvious questions** — don't ask "should we test?" Ask "integration tests or just unit tests?"
4. **Expose tradeoffs** — every option should show pros/cons in description
5. **Continue until complete** — keep asking until no ambiguities remain
6. **Capture decisions** — record why certain approaches were chosen/rejected

**Beads Integration:**
1. **Hash-based IDs** — No conflicts, multi-agent safe
2. **Priority levels** — Map scope to priority (critical → P0, etc.)
3. **Enhanced metadata** — Product alignment, mission, think analysis (NEW)
4. **Intent files** — Machine-readable intent for automation (NEW)
5. **Optional markdown** — Export for git history, but Beads is source of truth

## Migration from Markdown Workflow

**Old workflow:**
```bash
@idea "Add auth"  # → docs/drafts/idea-add-auth.md
@design idea-add-auth  # → docs/workstreams/backlog/WS-*.md
@build WS-001-01
```

**New Beads + ai-comm workflow:**
```bash
@idea "Add auth"  # → bd-0001 (Beads task) + docs/intent/bd-0001.json
@design bd-0001  # → bd-0001.1, bd-0001.2, ... (sub-tasks with execution graphs)
@build bd-0001.1  # → Updates Beads status
```

**Benefits:**
- No ID conflicts (hash-based vs manual PP-FFF-SS)
- Multi-agent ready (sub-tasks can be executed in parallel)
- Built-in dependency tracking
- `bd ready` shows what to work on next
- Product vision alignment tracked in metadata (NEW)
- Machine-readable intent for automation (NEW)

## Troubleshooting

**Beads not found:**
```bash
# Use mock mode
export BEADS_USE_MOCK=true

# Or install Beads
brew install go
go install github.com/steveyegge/beads/cmd/bd@latest
```

**Task not created:**
```bash
# Check Beads status
bd status

# View Beads logs
bd logs
```

**Intent validation failed (NEW):**
```bash
# Validate intent file
sdp schema validate docs/intent/bd-0001.json
```

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@idea "feature"` | Create Beads task with requirements + intent |
| `@think "complex analysis"` | Deep analysis before @idea (NEW) |
| `bd show {id}` | View task details |
| `bd ready` | List ready tasks |
| `bd list --status open` | List all open tasks |
| `@design {id}` | Decompose into workstreams |
| `@build {id}` | Execute workstream |
| `sdp schema validate` | Validate intent file (NEW) |

## Example Session

```bash
# Optional: Analyze complexity first
@think "Should we use JWT or sessions for auth?"

# Start interviewing
@idea "Add user authentication"

# ... (interviewing happens with PRODUCT_VISION alignment) ...

# Output:
✅ Created Beads task: bd-0001
   Title: Add user authentication
   Status: BeadsStatus.OPEN
   Priority: BeadsPriority.HIGH
   Intent: docs/intent/bd-0001.json

# Next step: decompose
@design bd-0001

# Output:
✅ Created 3 workstreams:
   bd-0001.1: Domain entities
   bd-0001.2: Repository layer
   bd-0001.3: Service layer

# Check what's ready
bd ready

# Output:
Ready tasks:
- bd-0001.1 (Domain entities)

# Execute
@build bd-0001.1

# After completion, bd-0001.2 becomes ready automatically!
```

---

**Version:** 2.1.0-beads-ai-comm
**Status:** Beads + AI-Comm Integration
**See Also:** `@design`, `@build`, `@oneshot`, `@think`
