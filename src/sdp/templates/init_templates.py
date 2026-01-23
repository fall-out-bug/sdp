"""Templates for sdp init command.

This module contains the template strings used by the `sdp init` command
to bootstrap SDP in a new project.
"""

PROJECT_MAP_TEMPLATE = """# Project Map: {project_name}

**Purpose:** Comprehensive map of project decisions. Read by agents before starting WS.

**Update:** After each significant WS - add entry about decisions made.

---

## How to Use

### For Agents (before starting WS):

1. **MUST READ** before planning/executing WS
2. Verify your WS doesn't conflict with existing decisions
3. If architectural decision needed → create ADR

### Structure:

```
PROJECT_MAP.md (this file)
    ↓ links to
docs/architecture/decisions/ (ADR - Architecture Decision Records)
```

---

## Key Decisions (Quick Reference)

| Area | Decision | ADR | Date |
|------|----------|-----|------|
| **Architecture** | TBD | - | - |
| **Tech Stack** | TBD | - | - |

---

## Current State

### Production Services
- (Add your services here)

### Domains
- (Add your domain areas here)

---

## Patterns & Conventions

### Naming
- Files: `snake_case.py`
- Classes: `PascalCase`
- Functions: `snake_case`
- Constants: `UPPER_SNAKE_CASE`

### Import Order
1. stdlib
2. third-party
3. local (project imports)
4. Relative imports

### Type Hints
- Use modern Python 3.10+ syntax: `list[str]`, `dict[str, int]`, `str | None`
- Always `-> None` for void functions
- Avoid `Any` without explicit justification

### Testing
- Unit tests: `tests/unit/`
- Integration: `tests/integration/`
- Coverage ≥ 80% for new code

---

## Active Constraints

### Code Quality
- Files < 200 LOC
- Complexity < 10 (CC)
- No nesting > 3 levels

### Security
- No hardcoded secrets
- Use environment variables
- Input validation always

---

## Tech Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| [Add component] | [Add technology] | [Version] | [Purpose] |

---

## References

- [SDP Protocol](https://github.com/your-org/sdp)
- [Architecture Decisions](docs/architecture/decisions/)

---

**Last Updated:** {date}
**Maintainers:** [Add maintainers]
"""

INDEX_TEMPLATE = """# Workstream Index

Track all workstreams for this project.

## Legend

- **Status:**
  - `backlog`: Planned but not started
  - `active`: Currently executing
  - `completed`: Finished successfully
  - `blocked`: Waiting on dependency

- **Size:**
  - `SMALL`: < 500 LOC
  - `MEDIUM`: 500-1500 LOC
  - `LARGE`: > 1500 LOC (should be split)

---

## Active Workstreams

Currently executing (status: active)

| WS ID | Title | Feature | Status | Size | Started |
|-------|-------|---------|--------|------|---------|
| - | - | - | - | - | - |

---

## Backlog

Planned workstreams (status: backlog)

| WS ID | Title | Feature | Size | Dependencies |
|-------|-------|---------|------|--------------|
| - | - | - | - | - |

---

## Completed Workstreams

Successfully finished (status: completed)

| WS ID | Title | Feature | Completed | Notes |
|-------|-------|---------|-----------|-------|
| - | - | - | - | - |

---

## Features Overview

High-level feature tracking

| Feature ID | Name | WS Count | Status | Description |
|------------|------|----------|--------|-------------|
| - | - | - | - | - |

---

**Total Workstreams:** 0  
**Completed:** 0  
**In Progress:** 0  
**Backlog:** 0
"""

WS_TEMPLATE = """---
ws_id: WS-{ID}
feature: F{XX}
status: backlog
size: SMALL
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
---

## WS-{ID}: {Title}

### 🎯 Goal

**What must WORK after this WS is complete:**
- [Specific functionality or improvement]
- [Measurable outcome - how to verify goal achieved]

**Acceptance Criteria:**
- [ ] [Verifiable condition 1]
- [ ] [Verifiable condition 2]
- [ ] [Verifiable condition 3]

**⚠️ Rule:** WS NOT complete until Goal achieved (all AC ✅).

---

### Context

[Why this task is needed, current state]

### Dependencies

[WS-XX / Independent]

### Input Files

- `path/to/file.py` — what's there

### Steps

1. [Atomic action]
2. [Next action]
3. ...

### Expected Result

- [What created/modified]
- [File structure]

### Scope Estimate

- **Files:** ~N created + ~M modified
- **Lines:** ~N (SMALL/MEDIUM/LARGE)

### Completion Criteria

```bash
# Tests
pytest tests/unit/test_module.py -v

# Coverage ≥ 80%
pytest --cov=module --cov-fail-under=80

# Regression
pytest tests/unit/ -m fast -v

# Linters
ruff check module/
mypy module/
```

### Constraints

- DON'T: [what not to touch]
- DON'T CHANGE: [what to leave]
"""
