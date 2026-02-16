---
name: vision
description: Strategic product planning - vision, PRD, roadmap from expert analysis
version: 2.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @vision - Strategic Product Planning

**Transform project ideas into product vision, PRD, and roadmap.**

---

## Workflow

When user invokes `@vision "AI task manager"`:

1. Interview user to gather requirements
2. Run parallel expert analysis (7 agents)
3. Generate vision artifacts
4. Extract feature drafts

---

## Step 1: Quick Interview (3-5 questions)

Use AskUserQuestion tool to gather requirements with multiSelect support:

- What problem are you solving?
- Who are your target users?
- What defines success in 1 year?
- What's your MVP?
- Who are your competitors?

## Step 2: Deep-Thinking Analysis (7 Expert Agents)

Run parallel expert analysis:

1. Product expert - Product-market fit analysis
2. Market expert - Competitive landscape analysis
3. Technical expert - Technical feasibility analysis
4. UX expert - User experience analysis
5. Business expert - Business model analysis
6. Growth expert - Growth strategy analysis
7. Risk expert - Risk and mitigation analysis

Synthesize outputs into coherent strategy.

## Step 3: Generate Artifacts

**PRODUCT_VISION.md** (project root):
- Why: Problem statement
- What: Product description
- Who: Target users
- Goals (1 year)
- Success Metrics
- Non-Goals

**docs/prd/PRD.md**:
- Functional Requirements
- Non-Functional Requirements
- Features (Prioritized P0/P1/P2)

**docs/roadmap/ROADMAP.md**:
- Q1: Foundation
- Q2: Growth
- Q3: Scale
- Q4: Maturity

## Step 4: Extract Features

For each P0/P1 feature, create draft in `docs/drafts/feature-{slug}.md`.

---

## When to Use

- **Initial project setup** - "What are we building?"
- **Quarterly review** - `@vision --review` - update vision based on progress
- **Major pivot** - "Is the direction changing?"
- **New market entry** - "Entering a new market?"

---

## Modes

| Mode | Output | Purpose |
|------|--------|---------|
| Default | Summary | Vision: AI-Powered Task Manager (7 experts) |
| `--quiet` | Exit status | Just check if complete |
| `--verbose` | Step-by-step | Full progress output |
| `--debug` | Internal state | Debug mode |

---

## Output

- `PRODUCT_VISION.md` (project root)
- `docs/prd/PRD.md`
- `docs/roadmap/ROADMAP.md`
- `docs/drafts/feature-*.md` (5-10 drafts)

---

## Example

```
@vision "AI-powered task manager"

Interview (3-5 questions)
Deep-thinking (7 expert agents)
Artifacts generated
8 feature drafts created in docs/drafts/
```

---

## See Also

- `@idea` - Feature-level requirements
- `@reality` - Reality check for completed projects
- `@feature` - Feature planning orchestrator
