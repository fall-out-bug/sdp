# 00-052-01: @vision Skill Structure

> **Beads ID:** sdp-vxvp
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 1A - Strategic Skills (@vision)
> **Size:** MEDIUM
> **Duration:** 2-3 days
> **Dependencies:**
> - 00-052-00 (Backup & worktree setup)

## Goal

Create `.claude/skills/vision/SKILL.md` for strategic product planning with expert agent analysis.

## Acceptance Criteria

- **AC1:** `.claude/skills/vision/` directory created
- **AC2:** `SKILL.md` created with interview + deep-thinking + artifact generation workflow
- **AC3:** Skill specifies 7 expert agents (product, market, technical, UX, business, growth, risk)
- **AC4:** Output artifacts defined: PRODUCT_VISION.md, docs/prd/PRD.md, docs/roadmap/ROADMAP.md
- **AC5:** Skill integrates with AskUserQuestion for progressive interview
- **AC6:** Skill integrates with Task tool for spawning expert agents

## Files

**Create:**
- `.claude/skills/vision/SKILL.md`

**Create (directory):**
- `.claude/skills/vision/`
- `src/sdp/vision/`
- `tests/sdp/vision/`

## Steps

### Step 1: Create Directory Structure

```bash
mkdir -p .claude/skills/vision
mkdir -p src/sdp/vision
mkdir -p tests/sdp/vision
```

### Step 2: Write SKILL.md

Create `.claude/skills/vision/SKILL.md`:

```markdown
---
name: vision
description: Strategic product planning - vision, PRD, roadmap from expert analysis
tools: Read, Write, Edit, AskUserQuestion, Task, Skill
version: 1.0.0
---

# @vision - Strategic Product Planning

**Transform project ideas into product vision, PRD, and roadmap.**

## When to Use
- **Initial project setup** - "Что мы строим?"
- **Quarterly review** - `@vision --review` - update vision based on progress
- **Major pivot** - "Меняется направление?"
- **New market entry** - "Выходим на новый рынок?"

## Workflow

### Step 1: Quick Interview (3-5 questions)

Use AskUserQuestion to understand:
- What problem are you solving?
- Who are your target users?
- What defines success in 1 year?

**Example:**
```python
AskUserQuestion(
    questions=[
        {
            "question": "What problem are you solving?",
            "header": "Problem",
            "options": [
                {"label": "User pain point", "description": "Fixes existing frustration"},
                {"label": "New opportunity", "description": "Enables new capabilities"},
                {"label": "Technical debt", "description": "Improves foundation"}
            ]
        },
        {
            "question": "Who are your target users?",
            "header": "Users",
            "options": [
                {"label": "Developers", "description": "Tools, libraries, platforms"},
                {"label": "Business", "description": "SaaS, enterprise software"},
                {"label": "Consumers", "description": "End-user applications"}
            ]
        },
        {
            "question": "What defines success in 1 year?",
            "header": "Success",
            "options": [
                {"label": "Adoption", "description": "1000+ active users"},
                {"label": "Revenue", "description": "$10K+ MRR"},
                {"label": "Impact", "description": "Open source community"}
            ],
            "multiSelect": True
        }
    ]
)
```

### Step 2: Deep-Thinking Analysis (7 Expert Agents)

Spawn parallel expert agents via Task tool:

```python
experts = [
    Task("Product expert", prompt="Analyze product-market fit for: {project}"),
    Task("Market expert", prompt="Analyze competitive landscape for: {project}"),
    Task("Technical expert", prompt="Analyze technical feasibility for: {project}"),
    Task("UX expert", prompt="Analyze user experience for: {project}"),
    Task("Business expert", prompt="Analyze business model for: {project}"),
    Task("Growth expert", prompt="Analyze growth strategy for: {project}"),
    Task("Risk expert", prompt="Analyze risks and mitigation for: {project}")
]

# Wait for all experts and synthesize
outputs = {e.description: e.result for e in experts}
synthesis = synthesize_expert_outputs(outputs)
```

### Step 3: Generate Artifacts

**PRODUCT_VISION.md** (project root):
```markdown
# Product Vision

## Why
{Problem statement}

## What
{Product description}

## Who
{Target users}

## Goals (1 year)
- [ ] {Goal 1}
- [ ] {Goal 2}

## Success Metrics
- Adoption: {metric}
- Quality: {metric}
- Growth: {metric}

## Non-Goals
{What we're NOT building}
```

**docs/prd/PRD.md**:
```markdown
# Product Requirements Document

## Requirements
### Functional
- FR1: {requirement}
- FR2: {requirement}

### Non-Functional
- NFR1: Performance <200ms
- NFR2: 99.9% uptime

## Features (Prioritized)
### P0 (Must Have)
- Feature 1: {description}
- Feature 2: {description}

### P1 (Should Have)
- Feature 3: {description}
```

**docs/roadmap/ROADMAP.md**:
```markdown
# Product Roadmap

## Q1 2026: Foundation
- MVP: {core features}

## Q2 2026: Growth
- Feature expansion: {features}

## Q3 2026: Scale
- Performance: {improvements}

## Q4 2026: Maturity
- Platform: {capabilities}
```

### Step 4: Extract Features

Use `src/sdp/vision/extractor.py` to extract features from PRD.

For each P0/P1 feature, create draft in `docs/drafts/feature-{slug}.md`.

## Outputs
- PRODUCT_VISION.md (project root)
- docs/prd/PRD.md
- docs/roadmap/ROADMAP.md
- docs/drafts/feature-*.md (5-10 drafts)

## Example

```bash
@vision "AI-powered task manager"

→ Interview (3-5 questions)
→ Deep-thinking (7 expert agents)
→ Artifacts generated
→ 8 feature drafts created in docs/drafts/
```

## See Also
- `.claude/skills/idea/SKILL.md` - Feature-level requirements
- `.claude/skills/reality/SKILL.md` - Reality check for completed projects
```

### Step 3: Verify Skill Format

```bash
ls -la .claude/skills/vision/
cat .claude/skills/vision/SKILL.md | head -20
```

Expected: Directory exists, SKILL.md has frontmatter with name/description/tools/version

### Step 4: Update Skills Registry

Add to `.claude/skills/index.json` (if exists):

```json
{
  "vision": {
    "name": "vision",
    "description": "Strategic product planning",
    "version": "1.0.0",
    "path": ".claude/skills/vision/SKILL.md"
  }
}
```

## Quality Gates

- SKILL.md valid frontmatter (name, description, tools, version)
- Workflow clearly specified (interview → deep-thinking → artifacts)
- Expert agents documented (7 agents)
- Output artifacts specified
- Integration with AskUserQuestion and Task tools documented

## Success Metrics

- Skill can be invoked via Skill tool
- Expert agents can be spawned via Task tool
- Artifacts generated in correct locations
- Feature extraction works (next workstream)
