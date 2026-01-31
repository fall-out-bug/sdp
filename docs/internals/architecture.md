# SDP Architecture

Detailed architecture documentation for SDP framework.

---

## Table of Contents

- [System Architecture](#system-architecture)
- [Component Design](#component-design)
- [Data Flow](#data-flow)
- [Design Patterns](#design-patterns)
- [Technology Stack](#technology-stack)

---

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      AI-IDE Interface                         │
│                   (Claude Code / Cursor)                     │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Skill System                            │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │ @feature│ │ @design │ │ @build  │ │ @review │          │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘          │
│       └──────────┬─────────────┘─────────────┘              │
└───────────────────┼───────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│                    Core Layer                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Workstream   │  │   Feature    │  │    Design    │     │
│  │   Parser     │  │   Manager    │  │     Graph    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└───────────────────────────┬─────────────────────────────────┘
                            │
                    ┌───────┴───────┐
                    │               │
                    ▼               ▼
┌──────────────────────┐  ┌──────────────────────────────────┐
│   Quality Gates      │  │     Multi-Agent System           │
│  ┌────────────────┐  │  │  ┌──────────┐  ┌─────────────┐  │
│  │ TDD Runner    │  │  │  │Planner   │  │ Orchestrator│  │
│  │ Coverage Check│  │  │  │Builder   │  │  Message    │  │
│  │ Mypy Strict   │  │  │  │Reviewer  │  │   Router    │  │
│  └────────────────┘  │  │  └──────────┘  └─────────────┘  │
└──────────────────────┘  └──────────────────────────────────┘
                            │
                    ┌───────┴───────┐
                    │               │
                    ▼               ▼
┌──────────────────────┐  ┌──────────────────────────────────┐
│  Integrations         │  │      Task Tracking               │
│  ┌────────────────┐  │  │  ┌──────────┐  ┌─────────────┐  │
│  │  Beads CLI    │  │  │  │  Beads   │  │   GitHub    │  │
│  │  GitHub API   │  │  │  │  Tasks   │  │   Issues    │  │
│  │  Telegram     │  │  │  │  Sync    │  │  Projects   │  │
│  └────────────────┘  │  │  └──────────┘  └─────────────┘  │
└──────────────────────┘  └──────────────────────────────────┘
```

---

## Component Design

### Workstream Model

**Purpose:** Atomic unit of work (500-1500 LOC)

**Structure:**
```python
@dataclass
class Workstream:
    ws_id: WorkstreamID          # PP-FFF-SS format
    title: str                   # Human-readable title
    feature: str                 # Parent feature ID
    status: WorkstreamStatus     # backlog|active|completed|blocked
    size: WorkstreamSize         # SMALL|MEDIUM|LARGE
    acceptance_criteria: list[AcceptanceCriterion]
    dependencies: list[str]      # WS IDs
```

**Lifecycle:**
```
backlog → active → completed
    ↓         ↓
  blocked ←─────┘
```

---

### Feature Model

**Purpose:** Group 5-30 related workstreams

**Structure:**
```python
@dataclass
class Feature:
    feature_id: str              # F{NNN} format
    title: str
    description: str
    workstreams: list[Workstream]
    status: FeatureStatus        # planned|active|completed
```

**Hierarchy:**
```
Release (product milestone)
  └─ Feature (5-30 workstreams)
      └─ Workstream (atomic task)
```

---

### Design Graph

**Purpose:** Manage workstream dependencies

**Algorithm:**
- Topological sort for execution order
- Detect circular dependencies
- Identify blocking workstreams

**Example:**
```
WS-001-01 (independent)
  ↓
WS-001-02 (depends on WS-001-01)
  ↓
WS-001-03 (depends on WS-001-02)
```

---

### Quality Gate Pipeline

**Checks:**
1. **TDD** - Tests written before code
2. **Coverage** - ≥80% line coverage
3. **Type Hints** - mypy --strict
4. **Linting** - ruff clean
5. **File Size** - <200 LOC
6. **Error Handling** - No bare exceptions

**Validation Points:**
- Pre-build (before starting)
- Post-build (after completion)
- Pre-deploy (before production)

---

## Data Flow

### Feature Development Flow

```
1. @feature "Add user auth"
   ├─ @idea (interview user)
   ├─ Generate PRODUCT_VISION.md
   └─ Create Beads task

2. @design idea-user-auth
   ├─ Read requirements
   ├─ Decompose into workstreams
   ├─ Create dependency graph
   └─ Generate WS files

3. @build WS-001-01 (for each WS)
   ├─ Pre-build validation
   ├─ TDD cycle (Red→Green→Refactor)
   ├─ Quality gate checks
   ├─ Git commit
   └─ Beads status update

4. @review F001
   ├─ Verify all WS completed
   ├─ Run quality checks
   └─ Generate report

5. @deploy F001
   ├─ Final verification
   ├─ Create git tag
   └─ Merge to main
```

---

### Autonomous Execution Flow

```
@oneshot F001
  │
  ├─ Spawn orchestrator agent
  │   ├─ Read feature spec
  │   ├─ Load workstreams
  │   ├─ Execute in dependency order
  │   │   ├─ @build WS-001-01
  │   │   ├─ @build WS-001-02
  │   │   └─ @build WS-001-03
  │   ├─ Save checkpoints
  │   └─ Send progress notifications
  │
  └─ Return agent ID (for resume)
```

---

## Design Patterns

### Repository Pattern

**Purpose:** Abstract data access

**Example:**
```python
class WorkstreamRepository:
    """Repository for workstream data."""

    def load(self, ws_id: str) -> Workstream:
        """Load workstream from file."""
        pass

    def save(self, ws: Workstream) -> None:
        """Save workstream to file."""
        pass
```

---

### Factory Pattern

**Purpose:** Create workstream instances

**Example:**
```python
class WorkstreamFactory:
    """Factory for creating workstreams."""

    @staticmethod
    def from_markdown(file_path: Path) -> Workstream:
        """Create workstream from markdown file."""
        pass
```

---

### State Machine

**Purpose:** Manage workstream status

**States:**
- `backlog` → `active` → `completed`
- `backlog` → `blocked` → `active`

**Transitions:**
```python
def transition_to(self, new_status: WorkstreamStatus) -> None:
    """Transition to new status if valid."""
    valid_transitions = {
        WorkstreamStatus.BACKLOG: [
            WorkstreamStatus.ACTIVE,
            WorkstreamStatus.BLOCKED
        ],
        # ...
    }
    # Validate and transition
```

---

### Observer Pattern

**Purpose:** Notify on workstream completion

**Example:**
```python
class WorkstreamObserver:
    """Observer for workstream events."""

    def on_complete(self, ws: Workstream) -> None:
        """Handle workstream completion."""
        # Update Beads
        # Send Telegram notification
        # Update GitHub issue
```

---

## Technology Stack

### Core Technologies

- **Python 3.14+** - Primary language
- **Click** - CLI framework
- **Pydantic** - Data validation
- **Pytest** - Testing framework
- **Mypy** - Type checking
- **Ruff** - Linting

### Dependencies

**Required:**
- `click` - CLI commands
- `pydantic` - Data validation
- `pyyaml` - YAML parsing
- `requests` - HTTP client

**Development:**
- `pytest` - Testing
- `pytest-cov` - Coverage
- `mypy` - Type checking
- `ruff` - Linting

**Optional:**
- `beads` - Task tracking
- `github` - GitHub integration
- `python-telegram-bot` - Notifications

---

## Architecture Principles

### Clean Architecture

**Layers:**
1. **Domain** - Core business logic
2. **Application** - Use cases
3. **Infrastructure** - External integrations
4. **Presentation** - CLI/API

**Rule:** Dependencies point inward

```
Presentation → Application → Domain
                     ↑          ↑
Infrastructure ──────────┘
```

---

### SOLID Principles

1. **SRP** - Single responsibility per class
2. **OCP** - Open for extension, closed for modification
3. **LSP** - Substitutable interfaces
4. **ISP** - Interface segregation
5. **DIP** - Depend on abstractions

---

### DRY (Don't Repeat Yourself)

- Reusable components
- Shared utilities
- Template-based generation

---

### KISS (Keep It Simple, Stupid)

- Simple over complex
- Clear over clever
- Explicit over implicit

---

## Performance

### Targets

- **Test Suite:** < 30 seconds
- **Validation:** < 5 seconds/file
- **Skill Overhead:** < 2 seconds
- **Message Latency:** < 100ms

### Optimization

- Lazy loading of modules
- Parallel test execution
- Incremental validation
- Cached dependency graphs

---

## Security

### Considerations

- No hardcoded secrets
- Environment variables for config
- Git hooks for validation
- Quality gates for security

### Enforced Checks

- `[security]` section in quality-gate.toml
- No SQL injection patterns
- No eval() usage
- HTTPS URLs required

---

## See Also

- [CODE_PATTERNS.md](../../CODE_PATTERNS.md) - Implementation patterns
- [PRINCIPLES.md](../PRINCIPLES.md) - Engineering principles
- [extending.md](extending.md) - How to extend SDP
- [contributing.md](contributing.md) - Contribution guide

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
