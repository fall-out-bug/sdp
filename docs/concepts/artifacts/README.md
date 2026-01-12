# Artifacts in Software Development

Artifacts are the outputs that document decisions, enable collaboration, and provide an audit trail.

## Why Artifacts Matter

- **Documentation**: Future you (or teammates) will understand why
- **Communication**: Share context without long meetings
- **Quality**: Formal outputs catch issues early
- **Audit**: Track decisions and their rationale

## Artifact Types by Workflow

| Artifact | Slash Commands | Task Orchestrator |
|----------|----------------|-------------------|
| CLAUDE.md | Required | Required |
| idea-{slug}.md | Required (from /idea) | Required |
| WS-XXX-YY.md | Required (from /design) | Required |
| Execution Report | Required (from /build) | Required |
| Review Report | Required (from /review) | Required |
| UAT Guide | Required (from /review) | Required |
| ADRs | When needed | When needed |
| Checkpoint JSON | No | Required (for /oneshot) |

## Core Artifacts

### 1. CLAUDE.md (Project Instructions)

**Purpose**: Give AI context about your project.

**Contents**:
- Tech stack
- Key commands
- Coding standards
- Project structure
- Important rules

**Example**:
```markdown
# Project: MyApp

## Tech Stack
Python 3.11, FastAPI, PostgreSQL

## Commands
- Test: pytest
- Run: uvicorn main:app --reload

## Rules
- Clean Architecture
- No silent failures
- Tests required for new code
```

### 2. Specification (spec.md)

**Purpose**: Define what to build.

**Contents**:
- Requirements
- User stories
- Acceptance criteria
- Out of scope
- Open questions

**Template**:
```markdown
# Feature: [Name]

## Overview
[1-2 sentences]

## Requirements
- REQ-1: The system must...
- REQ-2: The system must...

## User Stories
- As a [user], I want [action] so that [benefit]

## Acceptance Criteria
- [ ] [testable condition]
- [ ] [testable condition]

## Out of Scope
- [explicitly not included]

## Open Questions
- [things to clarify]
```

### 3. Design Document (design.md)

**Purpose**: Define how to build it.

**Contents**:
- Components
- Data flow
- API specification
- Database changes
- Dependencies

**Template**:
```markdown
# Design: [Feature Name]

## Components
- [Component 1] - [responsibility]
- [Component 2] - [responsibility]

## Data Flow
1. User action
2. System response
3. ...

## API Changes
POST /api/endpoint
  Request: { field: type }
  Response: { field: type }

## Database Changes
- New table: [name]
- New column: [table.column]

## Dependencies
- [new library or service]
```

### 4. Architecture Decision Record (ADR)

**Purpose**: Document why architectural decisions were made.

**Contents**:
- Context (the problem)
- Decision (the choice)
- Consequences (trade-offs)

**Template**:
```markdown
# ADR-NNNN: [Title]

## Status
Proposed | Accepted | Deprecated | Superseded by ADR-XXXX

## Context
[What is the issue? What forces are at play?]

## Decision
[What is the change we're making?]

## Consequences
[What becomes easier? What becomes harder?]
```

See [concepts/adr/](../adr/) for more details.

### 5. Test Report (test-report.md)

**Purpose**: Document quality verification results.

**Contents**:
- Test results summary
- Coverage report
- Issues found
- Recommendations

**Template**:
```markdown
# Test Report: [Feature]

## Summary
- Tests: 42 passed, 0 failed
- Coverage: 87%
- Issues: 2 found

## Coverage by Component
- domain/: 95%
- application/: 88%
- infrastructure/: 75%

## Issues Found
1. [Issue description] - [severity]
2. [Issue description] - [severity]

## Recommendations
- [what to fix before merge]
```

### 6. Decision Log (Multi-Agent only)

**Purpose**: Audit trail of all decisions made during development.

**Contents**:
- Date and agent
- Decision made
- Rationale
- Related artifacts

**Format**:
```markdown
# Decision Log: Epic [Name]

## 2024-01-15 - Architect
**Decision**: Use JWT for authentication
**Rationale**: Stateless, works well with microservices
**ADR**: docs/adr/0005-jwt-auth.md

## 2024-01-16 - Tech Lead
**Decision**: Split into 3 workstreams
**Rationale**: Independent components, parallel work possible
**Artifact**: consensus/artifacts/implementation.md
```

## Artifact Lifecycle

```
Requirement → Spec → Design → Code → Tests → Report
                ↓
               ADR (when decisions made)
```

### When to Create

| Artifact | Create When |
|----------|-------------|
| Spec | Starting any non-trivial feature |
| Design | After spec approved, before coding |
| ADR | Making technology or pattern choice |
| Test Report | After implementation complete |

### When to Update

| Artifact | Update When |
|----------|-------------|
| Spec | Requirements change |
| Design | Implementation deviates from plan |
| ADR | Decision revisited |
| Test Report | After fixing issues |

## Artifact Organization

### Simple Projects
```
project/
├── CLAUDE.md
└── docs/
    ├── specs/
    │   └── feature-name.md
    └── adr/
        └── 0001-decision.md
```

### Larger Projects
```
project/
├── CLAUDE.md
└── docs/
    ├── specs/
    │   ├── feature-1/
    │   │   ├── spec.md
    │   │   └── design.md
    │   └── feature-2/
    │       ├── spec.md
    │       └── design.md
    ├── adr/
    │   ├── 0001-database-choice.md
    │   └── 0002-auth-approach.md
    └── reviews/
        └── feature-1-review.md
```

## Tips

### Keep Artifacts Light
Don't write novels. Focus on:
- What's essential
- What's not obvious
- What you'll forget

### Use Templates
Templates ensure consistency and remind you what to include.

### Update or Delete
Outdated artifacts are worse than no artifacts. Keep them current or remove them.

### Link Related Artifacts
```markdown
## Related
- Spec: docs/specs/auth.md
- Design: docs/specs/auth-design.md
- ADR: docs/adr/0005-jwt.md
```
