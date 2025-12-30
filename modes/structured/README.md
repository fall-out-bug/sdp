# Structured Mode

One agent working through explicit phases. Creates documentation and artifacts along the way.

## When to Use

- Task takes 2-8 hours
- Involves multiple components
- Needs documentation for team/future reference
- Making architectural decisions worth recording
- Want quality without full multi-agent overhead

## How It Works

```
Phase 1: Analyze    →  spec.md (what to build)
    ↓
Phase 2: Design     →  design.md + ADRs (how to build)
    ↓
Phase 3: Implement  →  code + tests
    ↓
Phase 4: Review     →  verification
```

**Key insight**: Same AI, but you explicitly guide it through phases. This creates checkpoints and documentation.

## Setup

1. Create project structure:
```
your-project/
├── CLAUDE.md              # Project instructions
└── docs/
    ├── specs/             # Feature specifications
    ├── adr/               # Architecture Decision Records
    └── reviews/           # Code review reports
```

2. Use phase prompts from [prompts/](prompts/)

## The Four Phases

### Phase 1: Analyze

**Goal**: Understand and document what needs to be built.

**Input**: Feature request, bug report, or user story

**Output**: `docs/specs/{feature-name}.md`

**Prompt template**: [prompts/phase-1-analyze.md](prompts/phase-1-analyze.md)

```
Analyze this feature request and create a specification:

[paste feature request]

Save specification to docs/specs/[feature-name].md
Include: requirements, user stories, acceptance criteria, out of scope.
```

### Phase 2: Design

**Goal**: Plan the technical implementation.

**Input**: Specification from Phase 1

**Output**:
- `docs/specs/{feature-name}-design.md`
- `docs/adr/NNNN-{decision}.md` (if architectural decisions made)

**Prompt template**: [prompts/phase-2-design.md](prompts/phase-2-design.md)

```
Read docs/specs/[feature].md and create technical design.

Include: components, data flow, API changes, database changes.
If making architectural decisions, create ADR in docs/adr/
Save design to docs/specs/[feature]-design.md
```

### Phase 3: Implement

**Goal**: Write the code and tests.

**Input**: Specification + Design from Phases 1-2

**Output**: Code in `src/`, tests in `tests/`

**Prompt template**: [prompts/phase-3-implement.md](prompts/phase-3-implement.md)

```
Implement based on:
- Spec: docs/specs/[feature].md
- Design: docs/specs/[feature]-design.md

Follow TDD: write tests first.
Work layer by layer: domain → application → infrastructure → presentation.
```

### Phase 4: Review

**Goal**: Verify quality and completeness.

**Input**: All artifacts from Phases 1-3

**Output**: `docs/reviews/{feature-name}-review.md` (if issues found)

**Prompt template**: [prompts/phase-4-review.md](prompts/phase-4-review.md)

```
Review implementation of [feature]:

1. Check acceptance criteria from docs/specs/[feature].md
2. Verify design from docs/specs/[feature]-design.md was followed
3. Run tests, check coverage
4. Look for code smells, security issues

Report issues with file:line references.
```

## Example Workflow

Adding user notifications feature:

```bash
# Phase 1: Analyze
> Read the feature request in JIRA-123 and create specification.
  Save to docs/specs/user-notifications.md

# Review spec, approve or refine
> Looks good. Proceed.

# Phase 2: Design
> Read docs/specs/user-notifications.md and create design.
  We should decide on push vs email - create ADR.
  Save design to docs/specs/user-notifications-design.md

# Review design and ADR
> I agree with the decision. Proceed.

# Phase 3: Implement
> Implement user notifications based on spec and design.
  Start with domain layer.

> Now implement application layer.

> Now implement infrastructure layer with email provider.

> Now add API endpoints in presentation layer.

# Phase 4: Review
> Review the implementation:
  - Check all acceptance criteria
  - Verify design was followed
  - Run tests
  - Check for issues
```

## CLAUDE.md for Structured Mode

Extends Solo mode with phase awareness:

```markdown
# Project: MyApp

## Tech Stack
Python 3.11, FastAPI, PostgreSQL, pytest

## Structured Development

### Artifacts Location
- Specifications: docs/specs/{feature}.md
- Designs: docs/specs/{feature}-design.md
- ADRs: docs/adr/NNNN-{title}.md
- Reviews: docs/reviews/{feature}-review.md

### Phase Checklist
Before implementing:
- [ ] Specification approved
- [ ] Design document created
- [ ] ADRs for architectural decisions

After implementing:
- [ ] All acceptance criteria met
- [ ] Tests pass with adequate coverage
- [ ] No code smells or security issues

## Rules
- Clean Architecture: dependencies point inward
- No silent failures
- Test coverage ≥80% for new code
```

See [CLAUDE.md.example](CLAUDE.md.example) for full template.

## When to Upgrade to Multi-Agent

Switch to [Multi-Agent Mode](../multi-agent/) when:
- Task grows beyond one person/session
- Need parallel development (frontend + backend)
- Formal audit trail required
- Team wants explicit role separation

## Tips

### Don't Skip Phases
Each phase catches issues early:
- Analyze catches requirement gaps
- Design catches architectural problems
- Review catches implementation bugs

Skipping phases creates technical debt.

### Keep Artifacts Updated
If implementation diverges from design, update the design doc. Artifacts should reflect reality.

### Use ADRs Liberally
Any decision with alternatives is worth an ADR:
- "Why PostgreSQL over MongoDB?"
- "Why JWT over sessions?"
- "Why this library over that one?"

Future you will thank present you.

### Phase Boundaries are Flexible
You don't need to finish all of Phase 1 before starting Phase 2. Iterate:
- Rough spec → rough design → prototype → refine spec → refine design → implement
