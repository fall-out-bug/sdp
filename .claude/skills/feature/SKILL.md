---
name: feature
description: Feature planning orchestrator (idea → design → workstreams)
tools: Read, Write, Edit, Bash, AskUserQuestion, Skill
version: 6.0.0
changes:
  - Refactored to orchestrate @idea and @design skills
  - Removed duplicate agent spawning logic
  - Reduced from 227 LOC to ~100 LOC
---

# @feature - Feature Planning Orchestrator

**Orchestrate requirements gathering and workstream design.**

## Mental Model

```
@feature (Planning Orchestrator)
    │
    ├─► @idea (Requirements)
    │     ├─► Deep interviewing (AskUserQuestion)
    │     ├─► User stories
    │     └─► Success metrics
    │
    └─► @design (Workstream Planning)
          ├─► Codebase exploration
          ├─► Architecture decisions
          └─► Workstream files (00-FFF-SS.md)
```

## When to Use

- Starting new feature from scratch
- Need to gather requirements (@idea phase)
- Need to design workstreams (@design phase)
- Want interactive planning (questions, tradeoffs)

## Verbosity Tiers

```bash
@feature "Add OAuth2" --quiet     # Exit status only: ✅
@feature "Add OAuth2"             # Summary: ✅ Feature F051: OAuth2 Authentication (4 workstreams, 2.1h est.)
@feature "Add OAuth2" --verbose   # Step-by-step progress
@feature "Add OAuth2" --debug     # Internal state + API calls
```

**Examples:**

```bash
# Quiet mode
@feature "Add OAuth2" --quiet
# Output: ✅

# Default mode
@feature "Add OAuth2"
# Output: ✅ Feature F051: OAuth2 Authentication (4 workstreams, 2.1h est.)

# Verbose mode
@feature "Add OAuth2" --verbose
# Output:
# → @idea phase: Requirements gathering
#   → 18 questions across 6 cycles (15m)
# → @design phase: Architecture design
#   → 5 discovery blocks (45m)
# → Workstreams created: 00-051-01 through 00-051-04
# ✅ COMPLETE

# Debug mode
@feature "Add OAuth2" --debug
# Output:
# [DEBUG] Feature description: "Add OAuth2"
# [DEBUG] Calling Skill(@idea)...
# → @idea phase: Requirements gathering
#   [DEBUG] Question cycle 1: 3 questions
#   [DEBUG] Question cycle 2: 3 questions
#   [DEBUG] Question cycle 3: 3 questions
#   → 18 questions across 6 cycles (15m)
# [DEBUG] Creating Beads task for feature...
# [DEBUG] Calling Skill(@design)...
# → @design phase: Architecture design
#   [DEBUG] Discovery block 1: 3 questions
#   [DEBUG] Discovery block 2: 3 questions
#   [DEBUG] Discovery block 3: 3 questions
#   [DEBUG] Discovery block 4: 3 questions
#   [DEBUG] Discovery block 5: 3 questions
#   → 5 discovery blocks (45m)
# [DEBUG] Creating workstream files...
# [DEBUG] Created: docs/workstreams/backlog/00-051-01.md
# [DEBUG] Created: docs/workstreams/backlog/00-051-02.md
# [DEBUG] Created: docs/workstreams/backlog/00-051-03.md
# [DEBUG] Created: docs/workstreams/backlog/00-051-04.md
# → Workstreams created: 00-051-01 through 00-051-04
# ✅ COMPLETE
```

## What @feature Does

**@feature is an ORCHESTRATOR, not a duplicate:**

| Aspect | Old @feature | New @feature |
|--------|--------------|--------------|
| **Role** | Spawns 4 agents directly | Calls Skill('@idea') + Skill('@design') |
| **Logic** | Duplicates interview/discovery | Delegates to specialized skills |
| **Lines of Code** | ~227 LOC | ~100 LOC |
| **Maintainability** | Changes in 2 places | Single source of truth |

## Workflow

### Step 1: Quick Interview (3-5 questions)

Before spawning agents, ask quick questions to understand scope:

```python
AskUserQuestion(
    questions=[{
        "question": "What problem does this feature solve?",
        "header": "Problem",
        "options": [
            {"label": "User pain point", "description": "Fixes existing user friction"},
            {"label": "New capability", "description": "Enables new user workflows"},
            {"label": "Technical debt", "description": "Improves code quality/performance"}
        ],
        "multiSelect": false
    }, {
        "question": "Who are the primary users?",
        "header": "Users",
        "options": [
            {"label": "End users", "description": "External customers or users"},
            {"label": "Internal", "description": "Internal tools/operations"},
            {"label": "Developers", "description": "Developer experience/API"}
        ],
        "multiSelect": true
    }, {
        "question": "What defines success?",
        "header": "Success",
        "options": [
            {"label": "Adoption", "description": "% of users using the feature"},
            {"label": "Efficiency", "description": "Time savings, automation"},
            {"label": "Quality", "description": "Bug reduction, reliability"}
        ],
        "multiSelect": true
    }]
)
```

**Gate:** If description is vague (< 200 words, unclear scope), ask for clarification before proceeding.

### Step 2: Requirements Gathering (@idea)

```python
# Delegate to @idea skill for deep requirements gathering
Skill(
    skill="idea",
    args=feature_description  # e.g., "Add payment processing"
)
```

**What @idea does:**
- Deep interviewing via AskUserQuestion
- Explores technical approach
- Identifies tradeoffs and concerns
- Generates comprehensive spec in `docs/drafts/idea-{feature_name}.md`

**Output:**
- `docs/drafts/idea-{feature_name}.md` with requirements
- User stories, acceptance criteria
- Success metrics, stakeholders

### Step 3: Workstream Design (@design)

```python
# Delegate to @design skill for workstream planning
Skill(
    skill="design",
    args=idea_file  # e.g., "idea-payment-processing"
)
```

**What @design does:**
- EnterPlanMode for codebase exploration
- Asks architecture questions
- Designs workstream decomposition
- Requests approval via ExitPlanMode
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

## Output

**Success:**
```
✅ Feature planning complete
📄 Requirements: docs/drafts/idea-{feature_name}.md
📊 Workstreams: N created in docs/workstreams/backlog/00-FFF-*.md
🎯 Next step: @oneshot F{FF} or @build 00-FFF-01
```

**Example:**
```bash
User: @feature "Add payment processing"

Claude:
→ Step 1: Quick Interview (3 questions)
→ Step 2: @idea "Add payment processing"
   → Interviewing requirements...
   → Created: docs/drafts/idea-payment-processing.md
→ Step 3: @design idea-payment-processing
   → Exploring codebase...
   → Designing workstreams...
   → Created: 00-050-01.md, 00-050-02.md, 00-050-03.md
→ Step 4: Verification
   → ✅ 3 workstreams created

✅ Feature F050 planning complete
📄 docs/drafts/idea-payment-processing.md
📊 docs/workstreams/backlog/00-050-*.md (3 files)

Next: @oneshot F050 or @build 00-050-01
```

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

## Key Differences from @oneshot

| Aspect | @feature | @oneshot |
|--------|----------|----------|
| **Phase** | Planning | Execution |
| **Input** | Feature description | Feature ID or workstreams |
| **Output** | Workstream files | Implemented code |
| **Skills used** | @idea, @design | @build, @review, @deploy |
| **Human interaction** | Heavy (interviewing) | Minimal (only blockers) |
| **When to use** | Starting new feature | Workstreams exist |

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

## Version

**6.0.0** - Orchestrator refactoring
- Delegates to Skill('@idea') for requirements
- Delegates to Skill('@design') for workstreams
- Reduced from 227 LOC to ~100 LOC
- Removed duplicate agent spawning logic

**See Also:**
- `.claude/skills/idea/SKILL.md` — Requirements gathering
- `.claude/skills/design/SKILL.md` — Workstream planning
- `.claude/skills/oneshot/SKILL.md` — Execution orchestrator
- `CLAUDE.md` — Decision tree: @feature vs @oneshot
