---
name: feature
description: Feature planning orchestrator (idea -> design -> workstreams)
version: 7.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @feature - Feature Planning Orchestrator

**Orchestrate requirements gathering and workstream design.**

---

## EXECUTE THIS NOW

When user invokes `@feature "Add user authentication"`:

### Step 1: Quick Interview (3-5 questions)

Ask the user quick questions to understand scope:

- **Problem**: What problem does this feature solve?
  - User pain point / New capability / Technical debt
- **Users**: Who are the primary users?
  - End users / Internal / Developers
- **Success**: What defines success?
  - Adoption / Efficiency / Quality

**Gate:** If description is vague (< 200 words, unclear scope), ask for clarification before proceeding.

### Step 2: Requirements Gathering (@idea)

Invoke the idea skill for deep requirements gathering:

```
@idea "Add user authentication"
```

**What @idea does:**
- Deep interviewing with the user
- Explores technical approach
- Identifies tradeoffs and concerns
- Generates comprehensive spec in `docs/drafts/idea-{feature_name}.md`

**Output:**
- `docs/drafts/idea-{feature_name}.md` with requirements
- User stories, acceptance criteria
- Success metrics, stakeholders

### Step 3: Workstream Design (@design)

Invoke the design skill for workstream planning:

```
@design idea-user-authentication
```

**What @design does:**
- Explores codebase structure
- Asks architecture questions
- Designs workstream decomposition
- Requests user approval
- Creates `docs/workstreams/backlog/00-FFF-SS.md` files

**Output:**
- Workstream files (e.g., `00-050-01.md`, `00-050-02.md`)
- Dependency graph
- Architecture decisions

### Step 4: Verify Outputs

```bash
# Check that @idea created spec
ls docs/drafts/idea-{feature_name}.md

# Check that @design created workstreams
ls docs/workstreams/backlog/00-FFF-*.md

# Count workstreams
ws_count=$(ls docs/workstreams/backlog/00-FFF-*.md | wc -l)
echo "Created $ws_count workstreams"
```

---

## Mental Model

```
@feature (Planning Orchestrator)
    |
    +-> @idea (Requirements)
    |     +-> Deep interviewing
    |     +-> User stories
    |     +-> Success metrics
    |
    +-> @design (Workstream Planning)
          +-> Codebase exploration
          +-> Architecture decisions
          +-> Workstream files (00-FFF-SS.md)
```

---

## When to Use

- Starting new feature from scratch
- Need to gather requirements (@idea phase)
- Need to design workstreams (@design phase)
- Want interactive planning (questions, tradeoffs)

---

## Output

**Success:**
```
Feature planning complete
Requirements: docs/drafts/idea-{feature_name}.md
Workstreams: N created in docs/workstreams/backlog/00-FFF-*.md
Next step: @oneshot F{FF} or @build 00-FFF-01
```

**Example:**
```
User: @feature "Add payment processing"

Step 1: Quick Interview (3 questions)
Step 2: @idea "Add payment processing"
  Interviewing requirements...
  Created: docs/drafts/idea-payment-processing.md
Step 3: @design idea-payment-processing
  Exploring codebase...
  Designing workstreams...
  Created: 00-050-01.md, 00-050-02.md, 00-050-03.md
Step 4: Verification
  3 workstreams created

Feature F050 planning complete
docs/drafts/idea-payment-processing.md
docs/workstreams/backlog/00-050-*.md (3 files)

Next: @oneshot F050 or @build 00-050-01
```

---

## Beads Integration

**Detect Beads:**
```bash
if bd --version &>/dev/null && [ -d .beads ]; then
  BEADS_ENABLED=true
else
  BEADS_ENABLED=false
fi
```

**Beads operations:**
- @idea creates feature task if enabled
- @design creates workstream tasks if enabled
- @feature itself does NOT create Beads tasks (delegates)

---

## Key Differences from @oneshot

| Aspect | @feature | @oneshot |
|--------|----------|----------|
| **Phase** | Planning | Execution |
| **Input** | Feature description | Feature ID or workstreams |
| **Output** | Workstream files | Implemented code |
| **Skills used** | @idea, @design | @build, @review, @deploy |
| **Human interaction** | Heavy (interviewing) | Minimal (only blockers) |
| **When to use** | Starting new feature | Workstreams exist |

---

## Skip @feature If...

**Use @idea directly when:**
- You already have workstreams
- Only need requirements gathering
- Skip workstream design

**Use @design directly when:**
- You have requirements (idea file)
- Only need workstream planning
- Requirements already gathered

**Use @oneshot when:**
- Workstreams already exist
- Ready to implement
- Want autonomous execution

---

## See Also

- `@idea` - Requirements gathering
- `@design` - Workstream planning
- `@oneshot` - Execution orchestrator
- `CLAUDE.md` - Decision tree: @feature vs @oneshot
