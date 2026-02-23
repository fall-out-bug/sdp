---
name: feature
description: Feature planning orchestrator (discovery -> idea -> ux -> design -> workstreams)
version: 8.0.0
depends_on: "@discovery v1"
changes:
  - v8: Full product discovery flow with @discovery, @ux, impact analysis
  - Added --quick (skip @discovery), --infra (skip @ux)
  - Step 3.5: Impact analysis after @design
---

# @feature - Feature Planning Orchestrator

**Orchestrate product discovery, requirements gathering, UX research, and workstream design.**

---

## EXECUTE THIS NOW

When user invokes `@feature "Add user authentication"` (or with `--quick` / `--infra`):

### Step 0: Product Discovery (@discovery) — unless `--quick`

Invoke the discovery skill for roadmap pre-check and product research:

```
@discovery "Add user authentication"
```

**What @discovery does:**
- Phase 1: Roadmap pre-check (doc search, overlap report)
- Phase 2: Signal check (2 questions + web search) → route to Obvious / Competitive / Novel
- Phase 3: Product research (track-dependent: OBVIOUS skips; COMPETITIVE single pass; NOVEL max 3 iterations)
- Phase 4: Feature brief in `docs/drafts/discovery-{slug}.md`

**Gate:** If user resolves overlap as "extend" or "supersede", handle before proceeding.

**`--quick` flag:** Skip @discovery entirely. Proceed to Step 1 with current behavior (quick interview → @idea → @design).

---

### Step 1: Quick Interview (3-5 questions)

If @discovery ran: use its output. If `--quick`: ask these questions:

- **Problem**: What problem does this feature solve?
  - User pain point / New capability / Technical debt
- **Users**: Who are the primary users?
  - End users / Internal / Developers
- **Success**: What defines success?
  - Adoption / Efficiency / Quality

**Gate:** If description is vague (< 200 words, unclear scope), ask for clarification before proceeding.

---

### Step 2: Requirements Gathering (@idea)

Invoke the idea skill. Pass discovery output when available:

```
@idea "Add user authentication" --spec docs/drafts/discovery-{slug}.md
```

If `--quick` or no discovery brief:
```
@idea "Add user authentication"
```

**What @idea does:**
- Deep interviewing with the user (or skips cycles pre-answered by @discovery)
- Explores technical approach
- Identifies tradeoffs and concerns
- Generates spec in `docs/drafts/idea-{feature_name}.md` or `docs/intent/{task_id}.json`

---

### Step 2.5: UX Research (@ux) — unless `--infra`

**Auto-trigger heuristic:** Run @ux when @idea output contains user-facing keywords (`ui`, `user`, `interface`, `dashboard`, `form`, `flow`, `UX`, `screen`, `page`, `button`) and lacks infra signals (`K8s`, `CRD`, `reconciler`, `stream`, `JetStream`, `CLI-only`).

```
@ux {feature-id}
```

**What @ux does:**
- 6-question listening session (mental model elicitation)
- Autonomous codebase research (patterns, accessibility, error handling gaps)
- Output: `docs/ux/{feature}.md` with typed schema

**`--infra` flag:** Skip @ux.

---

### Step 3: Workstream Design (@design)

Invoke the design skill for workstream planning:

```
@design {task_id}
```

**What @design does:**
- Loads requirements from @idea (and @ux when present)
- Explores codebase structure
- Asks architecture questions
- Creates `docs/workstreams/backlog/00-FFF-SS.md` files
- Converts UX friction_points and ux_risks into acceptance criteria when docs/ux/ exists

---

### Step 3.5: Impact Analysis

After @design creates workstream files:

1. Read Scope Files from all new workstreams.
2. For each scope file, run:
   ```bash
   grep -rl "<scope_file>" docs/workstreams/backlog/*.md
   ```
3. Also run:
   ```bash
   rg "<domain term1>|<domain term2>" docs/ -l
   sdp drift detect
   ```
4. Categorize matches:
   - **[FILE CONFLICT]** Same file scoped by multiple workstreams → recommend `depends_on`
   - **[DATA BOUNDARY]** New feature modifies type used by another → recommend schema-first or extend
   - **[DEPENDENCY CHAIN]** New feature inserts into existing F00X → F00Y path → show updated graph
   - **[PRIORITY SHIFT]** New feature P0 but depends on P2 blocking other P0 → recommend reprioritize

5. Present **Impact Report** — user must acknowledge before @oneshot.
6. For resolved conflicts: update workstream frontmatter (depends_on, related_to, status).

---

### Step 4: Verify Outputs

```bash
# Check discovery brief (if not --quick)
ls docs/drafts/discovery-{slug}.md

# Check @idea spec
ls docs/drafts/idea-{feature_name}.md docs/intent/*.json

# Check @ux output (if user-facing)
ls docs/ux/{feature}.md

# Check workstreams
ls docs/workstreams/backlog/00-FFF-*.md
ws_count=$(ls docs/workstreams/backlog/00-FFF-*.md 2>/dev/null | wc -l)
echo "Created $ws_count workstreams"
```

---

## Mental Model

```
@feature (Planning Orchestrator)
    |
    +-> @discovery (Product Discovery) [unless --quick]
    |     +-> Roadmap pre-check
    |     +-> Signal check → Obvious / Competitive / Novel
    |     +-> Feature brief
    |
    +-> @idea (Requirements)
    |     +-> Deep interviewing
    |     +-> User stories, success metrics
    |
    +-> @ux (UX Research) [unless --infra, user-facing only]
    |     +-> Mental model elicitation
    |     +-> UX Risk Register
    |
    +-> @design (Workstream Planning)
    |     +-> Architecture decisions
    |     +-> Workstream files (00-FFF-SS.md)
    |
    +-> Impact Analysis (Step 3.5)
          +-> FILE CONFLICT / DATA BOUNDARY / DEPENDENCY CHAIN / PRIORITY SHIFT
```

---

## Flags

| Flag | Effect |
|------|--------|
| `--quick` | Skip @discovery; use original flow (quick interview → @idea → @design) |
| `--infra` | Skip @ux (infrastructure feature, no user-facing surface) |

---

## When to Use

- Starting new feature from scratch
- Need full product discovery (@discovery phase)
- Need requirements gathering (@idea phase)
- Need UX research for user-facing features (@ux phase)
- Need workstream design (@design phase)

---

## Output

**Success:**
```
Feature planning complete
Discovery: docs/drafts/discovery-{slug}.md
Requirements: docs/drafts/idea-{feature_name}.md
UX: docs/ux/{feature}.md (if user-facing)
Workstreams: N created in docs/workstreams/backlog/00-FFF-*.md
Next step: @oneshot F{FF} or @build 00-FFF-01
```

---

## Example Session

```
User: @feature "Add payment processing"

Step 0: @discovery "Add payment processing"
  Roadmap pre-check: no overlaps
  Signal check: COMPETITIVE track
  Research: 3 alternatives, build-vs-adopt decision
  Created: docs/drafts/discovery-payment-processing.md

Step 2: @idea "Add payment processing" --spec docs/drafts/discovery-payment-processing.md
  Skipped Vision, Problem (pre-answered by discovery)
  Created: docs/drafts/idea-payment-processing.md

Step 2.5: @ux payment-processing (auto-triggered)
  6 UX questions, codebase scan
  Created: docs/ux/payment-processing.md

Step 3: @design sdp-xxx
  Created: 00-050-01.md, 00-050-02.md, 00-050-03.md

Step 3.5: Impact Analysis
  [LOW] No conflicts found.

Feature F050 planning complete
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
| **Skills used** | @discovery, @idea, @ux, @design | @build, @review, @deploy |
| **Human interaction** | Heavy (interviewing) | Minimal (only blockers) |
| **When to use** | Starting new feature | Workstreams exist |

---

## Skip @feature If...

**Use @discovery directly when:**
- Only need roadmap pre-check or product research
- Stop after discovery brief

**Use @idea directly when:**
- You already have workstreams
- Only need requirements gathering
- Skip workstream design

**Use @design directly when:**
- You have requirements (idea file)
- Only need workstream planning

**Use @oneshot when:**
- Workstreams already exist
- Ready to implement
- Want autonomous execution

---

## See Also

- `@discovery` - Product discovery gate
- `@idea` - Requirements gathering
- `@ux` - UX research
- `@design` - Workstream planning
- `@oneshot` - Execution orchestrator
- `CLAUDE.md` - Decision tree: @feature vs @oneshot
