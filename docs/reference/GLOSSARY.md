# SDP Canonical Glossary

**Version:** 1.0
**Last Updated:** 2026-01-29
**Maintainer:** SDP Protocol Team

This glossary provides canonical definitions for all SDP (Spec-Driven Protocol) terms, concepts, and conventions. It resolves naming conflicts and serves as the single source of truth for terminology.

---

## Quick Navigation

- [Core Concepts](#core-concepts)
- [Hierarchy Levels](#hierarchy-levels)
- [Identifier Formats](#identifier-formats)
- [Commands & Skills](#commands--skills)
- [Development Workflow](#development-workflow)
- [Quality Gates](#quality-gates)
- [Architecture](#architecture)
- [Testing](#testing)
- [Multi-Agent](#multi-agent)
- [Principles](#principles)
- [Tools & Integrations](#tools--integrations)

---

## Core Concepts

### SDP (Spec-Driven Protocol)

**Definition:** Workstream-driven development methodology for AI agents with multi-agent coordination.

**Context:** The overarching framework that defines how AI agents should plan, execute, and track development work through workstreams.

**Example:** "SDP v0.5.0 integrates multi-agent coordination with Beads task tracking."

**See Also:** Workstream, Feature, Release, Beads

**Related Terms:** AI-Comm, Beads CLI, Unified Orchestrator

---

### Workstream (WS)

**Definition:** Atomic task unit representing the smallest executable piece of work in SDP. Executed in one-shot fashion with clear acceptance criteria.

**Canonical Format:** `PP-FFF-SS` or `WS-PP-FFF-SS`

**Size Classifications:**
- **SMALL:** < 500 LOC, < 1500 tokens
- **MEDIUM:** 500-1500 LOC, 1500-5000 tokens
- **LARGE:** > 1500 LOC → split into 2+ WS

**Key Characteristics:**
- One-shot execution (completed in single session)
- TDD-driven (Red → Green → Refactor)
- Clear acceptance criteria
- ≤ 200 LOC per file
- ≥ 80% test coverage

**Examples:**
- `WS-001-01`: First workstream of first feature
- `00-012-03`: Third workstream of feature 012 in project 00

**Legacy Names:** ~~Epic~~, ~~Sprint~~, ~~Task~~

**See Also:** Feature, Release, Acceptance Criteria, TDD

**Related Terms:** Beads Task, WS-PP-FFF-SS format

---

### Feature

**Definition:** Major deliverable composed of 5-30 workstreams. Represents a significant piece of functionality that delivers user value.

**Canonical Format:** `PP-FFF` or `F-FFF`

**Scope:** 5-30 workstreams

**Lifecycle:**
1. Requirements gathering (`@idea` or `@feature`)
2. Workstream planning (`@design`)
3. Execution (`@build` or `@oneshot`)
4. Review (`@review`)
5. Deployment (`@deploy`)

**Example:** "F24: Unified Workflow" - Feature implementing unified workflow system

**See Also:** Workstream, Release, Product Vision

**Related Terms:** Beads Task (parent), Intent, Draft

---

### Release

**Definition:** Product milestone grouping 10-30 features. Represents a significant version or delivery to production.

**Canonical Format:** `R1`, `R2`, etc.

**Scope:** 10-30 features

**Example:** "R1: Submissions E2E" - First release with end-to-end submission functionality

**See Also:** Feature, Workstream, Deployment

---

### Beads

**Definition:** Task tracking system integrated with SDP for managing hierarchical tasks and dependencies.

**Context:** Provides CLI and Python client for task management, dependency tracking, and status management.

**Key Capabilities:**
- Create hierarchical tasks (Feature → Workstreams)
- Add dependencies between tasks
- Track task status
- Query ready tasks (unblocked work)

**Examples:**
```python
client.create_task(BeadsTaskCreate(
    title="User Authentication",
    parent_id=feature.id,
))
client.add_dependency(ws2.id, ws1.id, dep_type="blocks")
ready = client.get_ready_tasks()
```

**See Also:** Workstream, Feature, Unified Orchestrator

---

## Hierarchy Levels

### Three-Level Hierarchy

**Definition:** SDP organizes work into three levels: Release → Feature → Workstream.

| Level | Scope | Size | Example |
|-------|-------|------|---------|
| **Release** | Product milestone | 10-30 Features | R1: Submissions E2E |
| **Feature** | Major feature | 5-30 Workstreams | F24: Unified Workflow |
| **Workstream** | Atomic task | SMALL/MEDIUM/LARGE | WS-060: Domain Model |

**Rules:**
- Workstreams belong to exactly one Feature
- Features belong to exactly one Release
- Dependencies flow only between siblings (WS→WS, Feature→Feature)

**See Also:** Workstream, Feature, Release

---

## Identifier Formats

### PP-FFF-SS Format (Canonical)

**Definition:** Standard identifier format for workstreams.

**Structure:**
- **PP** - Product/Project ID (01-99)
- **FFF** - Feature number (001-999)
- **SS** - Workstream sequence (01-99)

**Examples:**
- `01-001-01` - First workstream of first feature in project 01
- `00-012-03` - Third workstream of feature 012 in project 00

**Prefix Variants:**
- `WS-PP-FFF-SS` - Full prefix (canonical)
- `PP-FFF-SS` - Short prefix (acceptable)
- `WS-001-01` - Legacy format (deprecated)

**Validation:**
```python
# Regex pattern
r'^\d{2}-\d{3}-\d{2}$'
```

**See Also:** Workstream, Feature ID, Naming Conflicts

---

### Feature ID (PP-FFF)

**Definition:** Identifier for features.

**Structure:**
- **PP** - Product/Project ID (01-99)
- **FFF** - Feature number (001-999)

**Examples:**
- `01-001` - First feature in project 01
- `00-012` - Feature 012 in project 00

**Prefix Variants:**
- `F-FFF` - Short form (acceptable)
- `PP-FFF` - Full form (canonical)

**See Also:** Feature, Workstream ID

---

### beads-sdp-XXX Format (Legacy)

**Definition:** Legacy identifier format from early SDP versions. Still found in some documentation and file names.

**Status:** Deprecated, but may appear in:
- File names (e.g., `beads-sdp-118.md`)
- Legacy documentation
- Old Beads task titles

**Migration:** Convert to `PP-FFF-SS` format.

**Example:** `beads-sdp-118` → `01-118-XX` (multiple workstreams)

**See Also:** PP-FFF-SS Format, Naming Conflicts

---

### Idea ID

**Definition:** Identifier for requirements gathered via `@idea` command.

**Format:** `idea-{slug}` where slug is derived from the feature title.

**Examples:**
- `idea-user-auth`
- `idea-password-reset`

**Usage:**
- Input to `@design` command
- Prefix for draft files
- Temporary identifier until Feature ID assigned

**Lifecycle:** idea → Feature → Workstreams

**See Also:** Feature, @idea, @design

---

## Naming Conflicts (Resolution)

### Multiple Identifier Systems

**Issue:** SDP documentation uses multiple identifier formats that can be confusing:

| Format | Example | Usage | Status |
|--------|---------|-------|--------|
| `WS-PP-FFF-SS` | `WS-001-01` | Canonical workstream ID | ✅ Preferred |
| `PP-FFF-SS` | `01-001-01` | Short workstream ID | ✅ Acceptable |
| `beads-sdp-XXX` | `beads-sdp-118` | Legacy feature/workstream | ⚠️ Deprecated |
| `F-FFF` | `F24` | Short feature ID | ✅ Acceptable |
| `PP-FFF` | `01-001` | Feature ID | ✅ Preferred |
| `EP-XX` | `EP08` | Epic (legacy) | ❌ Replaced by Feature |

**Resolution Strategy:**
1. **Workstreams:** Always use `PP-FFF-SS` format in new docs
2. **Features:** Use `PP-FFF` format
3. **Releases:** Use `R1`, `R2` format
4. **Ideas:** Use `idea-{slug}` format
5. **Migration:** Gradually update legacy references

**See Also:** PP-FFF-SS Format, Feature ID, beads-sdp-XXX Format

---

## Commands & Skills

### @ Command Prefix

**Definition:** Prefix for user-invocable skills in Claude Code and Cursor.

**Convention:** `@{skill-name}` for skills that require human interaction or approval.

**Examples:**
- `@feature "Add user authentication"`
- `@design idea-user-auth`
- `@build WS-001-01`
- `@review F01`
- `@deploy F01`

**Usage:** Invoked by users to initiate workflows.

**See Also:/ Command Prefix, Skill, Command

---

### / Command Prefix

**Definition:** Prefix for built-in CLI commands and internal skills.

**Convention:** `/{command-name}` for commands that don't require approval.

**Examples:**
- `/debug "Test fails unexpectedly"`
- `/issue "Login fails on Firefox"`
- `/tdd` (internal, not user-invocable)

**Usage:** Invoked by users or called automatically by other skills.

**See Also:** @ Command Prefix, Skill, Command

---

### Skill

**Definition:** Reusable workflow automation defined in `.claude/skills/{name}/SKILL.md`.

**Location:** `.claude/skills/{skill-name}/SKILL.md`

**Types:**
- **User-invocable:** `@feature`, `@idea`, `@design`, `@build`, `@review`, `@deploy`, `@hotfix`, `@bugfix`, `@issue`, `@oneshot`
- **Internal:** `/tdd` (called by `@build`)

**Invocation:**
```bash
@feature "Add payment processing"
@design idea-payments
@build WS-001-01
```

**See Also:** @ Command Prefix, Command, Prompt

---

### @feature (Skill)

**Definition:** Unified entry point for feature development with progressive disclosure.

**Purpose:** Interactive requirements gathering via deep questioning.

**Workflow:**
1. Asks deep questions (technical approach, UI/UX, security, concerns)
2. Creates Intent file: `docs/intent/sdp-XXX.json`
3. Creates Draft: `docs/drafts/beads-sdp-XXX.md`
4. Creates Product Vision: `PRODUCT_VISION.md`

**Example:**
```bash
@feature "Add user authentication"
# Claude asks:
# - Technical approach (JWT vs sessions?)
# - UI/UX requirements
# - Database schema
# - Testing strategy
# - Security concerns
```

**Output:** `docs/intent/sdp-XXX.json`, `docs/drafts/beads-sdp-XXX.md`

**See Also:** @idea, @design, Intent, Draft

---

### @idea (Skill)

**Definition:** Interactive requirements gathering using AskUserQuestion for deep interviewing.

**Purpose:** Explore tradeoffs and requirements without obvious questions.

**Workflow:**
1. Asks clarifying questions via AskUserQuestion
2. Explores technical approaches
3. Creates comprehensive spec in `docs/drafts/`

**Example:**
```bash
@idea "User can reset password via email"
# Claude asks:
# - Technical approach (email service, token storage)
# - UI/UX (where in app, error messages)
# - Security (token expiry, rate limiting)
# - Concerns (complexity, failure modes)
```

**Output:** `docs/drafts/idea-{slug}.md`

**See Also:** @feature, @design, Draft

---

### @design (Skill)

**Definition:** Interactive planning using EnterPlanMode for codebase exploration.

**Purpose:** Decompose features into workstreams with dependency tracking.

**Workflow:**
1. Enters Plan Mode (codebase exploration allowed)
2. Asks architecture questions via AskUserQuestion
3. Designs workstream decomposition
4. Requests approval via ExitPlanMode

**Input:** `idea-{slug}` or feature description

**Output:** `docs/workstreams/backlog/PP-FFF-SS.md`

**Example:**
```bash
@design idea-password-reset
# Claude explores codebase:
# - Existing auth infrastructure
# - Email service availability
# Creates: WS-XXX-01, WS-XXX-02, etc.
```

**See Also:** @idea, @build, Workstream

---

### @build (Skill)

**Definition:** Execute workstream with TodoWrite real-time progress tracking.

**Purpose:** Atomic workstream execution with TDD discipline.

**Workflow:**
1. Pre-build validation
2. Red: Write failing test
3. Green: Implement minimum code
4. Refactor: Improve design
5. Run quality gates
6. Append execution report
7. Git commit

**Progress Tracking:** Real-time TodoWrite updates

**Input:** Workstream ID (`PP-FFF-SS`)

**Output:** Completed workstream with execution report

**Example:**
```bash
@build WS-001-01
# Shows real-time progress:
# ✓ [completed] Pre-build validation
# ✓ [in_progress] Write failing test (Red)
# • [pending] Implement minimum code (Green)
```

**See Also:/tdd, Workstream, TDD

---

### /tdd (Skill)

**Definition:** Internal skill enforcing TDD cycle (Red → Green → Refactor).

**Purpose:** Called automatically by `@build` to ensure TDD discipline.

**Workflow:**
1. **Red:** Write failing test
2. **Green:** Implement minimum code to pass
3. **Refactor:** Improve design while tests pass

**Called By:** `@build` (automatic)

**Not User-Invocable:** Internal implementation detail

**See Also:** @build, TDD, Test-Driven Development

---

### @oneshot (Skill)

**Definition:** Autonomous execution of all workstreams for a feature using Task-based multi-agent orchestration.

**Purpose:** Execute entire feature without human intervention.

**Workflow:**
1. Spawn orchestrator agent via Task tool
2. Execute all WS in dependency order
3. Save checkpoints after each WS
4. Send Telegram notifications
5. Resume from interruption if needed

**Options:**
- `@oneshot F01` - Normal execution
- `@oneshot F01 --background` - Background execution
- `@oneshot F01 --resume {agent-id}` - Resume from interruption

**Example:**
```bash
@oneshot F01
# Spawns orchestrator agent with ID: abc123xyz
# Executes all workstreams
# Provides checkpoint for resume
```

**Output:** Completed feature with UAT guide

**See Also:** Orchestrator Agent, Task Tool, Checkpoint

---

### @review (Skill)

**Definition:** Quality review of completed features or workstreams.

**Purpose:** Validate all quality gates passed.

**Validates:**
- ✅ Tests ≥80% coverage
- ✅ No tech debt
- ✅ Clean architecture
- ✅ All quality gates passed

**Output:** APPROVED or CHANGES_REQUESTED with feedback

**Example:**
```bash
@review F01
# Returns: APPROVED
```

**See Also:** Quality Gates, @deploy

---

### @deploy (Skill)

**Definition:** Production deployment with artifact generation.

**Purpose:** Generate deployment artifacts and create release.

**Generates:**
- `docker-compose.yml`
- `.github/workflows/deploy.yml`
- `CHANGELOG.md` entry
- Git tag: `v{version}`

**Example:**
```bash
@deploy F01
# Generates deployment configuration
# Creates git tag v1.0.0
```

**See Also:** Release, Feature, @review

---

### @hotfix (Skill)

**Definition:** Emergency fix for P0 (critical) issues with <2 hour turnaround.

**Purpose:** Fast-track critical fixes to production.

**Workflow:**
1. Branch from `main`
2. Minimal fix
3. Quick validation
4. Deploy to production
5. Followup with proper workstream if needed

**Use When:** System is down or critical bug blocking users

**Constraint:** < 2 hours from start to deployment

**See Also:** @bugfix, @issue, P0 Severity

---

### @bugfix (Skill)

**Definition:** Quality fix for P1/P2 issues with <24 hour turnaround.

**Purpose:** Address non-critical bugs with proper workflow.

**Workflow:**
1. Create workstream from issue
2. Execute with @build
3. Review with @review
4. Deploy with @deploy

**Use When:** Bug is important but not critical

**Constraint:** < 24 hours

**See Also:** @hotfix, @issue, P1/P2 Severity

---

### @issue (Skill)

**Definition:** Debug and route bugs to appropriate fix workflow.

**Purpose:** Classify bug severity and route to hotfix/bugfix/backlog.

**Workflow:**
1. Analyze bug report
2. Classify severity (P0/P1/P2/P3)
3. Route appropriately:
   - P0 → @hotfix
   - P1/P2 → @bugfix
   - P3 → backlog

**Example:**
```bash
@issue "Login fails on Firefox"
# Analyzes and routes to @bugfix
```

**See Also:** @hotfix, @bugfix, Severity Classification

---

### /debug (Skill)

**Definition:** Systematic debugging using scientific method for evidence-based root cause analysis.

**Purpose:** Debug test failures, unexpected behavior, or production issues.

**Method:**
1. Observe (gather evidence)
2. Hypothesize (formulate theory)
3. Test (verify hypothesis)
4. Confirm (root cause identified)

**Example:**
```bash
/debug "Test fails unexpectedly"
# Systematic investigation
```

**See Also:** @issue, Scientific Method, Root Cause Analysis

---

## Development Workflow

### TDD (Test-Driven Development)

**Definition:** Development methodology where tests are written before implementation.

**Cycle:** Red → Green → Refactor

**Steps:**
1. **Red:** Write failing test that defines desired behavior
2. **Green:** Write minimum code to make test pass
3. **Refactor:** Improve code design while tests pass

**Benefits:**
- Tests document expected behavior
- Confident refactoring
- 100% coverage by design
- Better API design

**Anti-pattern:** Writing tests after implementation

**Example:**
```python
# Step 1: RED - Write failing test
def test_user_can_be_created():
    user = User(email="test@example.com", name="Test")
    assert user.email == "test@example.com"
# Result: NameError: name 'User' is not defined

# Step 2: GREEN - Minimal implementation
@dataclass
class User:
    email: str
    name: str
# Result: Test passes

# Step 3: REFACTOR - Add validation
@dataclass
class User:
    email: str
    name: str
    def __post_init__(self) -> None:
        if "@" not in self.email:
            raise ValueError("Invalid email")
# Result: Tests still pass, code improved
```

**Enforcement:** `/tdd` skill called by `@build`

**See Also:/tdd, @build, Test Coverage

---

### Red-Green-Refactor

**Definition:** The three phases of TDD cycle.

**Red:**
- Write failing test
- Test should fail for expected reason
- Confirms test is valid

**Green:**
- Write minimum code to pass test
- No concern for design quality yet
- Just make it pass

**Refactor:**
- Improve code design
- Remove duplication
- Enhance readability
- Tests must still pass

**See Also:** TDD, @build,/tdd

---

### Intent

**Definition:** Machine-readable specification of feature requirements in JSON format.

**Location:** `docs/intent/sdp-XXX.json`

**Purpose:** Structured representation of user requirements for AI consumption.

**Structure:**
```json
{
  "title": "Feature title",
  "description": "Detailed description",
  "requirements": ["req1", "req2"],
  "acceptance_criteria": ["ac1", "ac2"],
  "constraints": ["constraint1"]
}
```

**Created By:** `@feature` or `@idea` skills

**Used By:** `@design` skill for workstream planning

**See Also:** Draft, Feature, Schema

---

### Draft

**Definition:** Human-readable specification of feature requirements in Markdown format.

**Location:** `docs/drafts/beads-sdp-XXX.md` or `docs/drafts/idea-{slug}.md`

**Purpose:** Detailed requirements document for human review and planning.

**Created By:** `@feature` or `@idea` skills

**Used By:** `@design` skill for workstream planning

**See Also:** Intent, Feature, Specification

---

### Specification (Spec)

**Definition:** Formal feature document with requirements, workstreams, and acceptance criteria.

**Location:** `docs/specs/{feature-id}.md`

**Purpose:** Complete feature specification for implementation.

**Contents:**
- Requirements
- Workstream breakdown
- Acceptance criteria
- Dependencies
- Technical approach

**See Also:** Intent, Draft, Workstream

---

### Product Vision

**Definition:** High-level project manifesto describing long-term goals and direction.

**Location:** `PRODUCT_VISION.md` (root)

**Purpose:** Align all features and workstreams with strategic objectives.

**Created By:** `@feature` skill

**See Also:** Feature, Release, Roadmap

---

### Acceptance Criteria (AC)

**Definition:** Specific, measurable conditions that must be met for a workstream to be considered complete.

**Format:** Checkbox list in workstream file.

**Example:**
```markdown
**Acceptance Criteria:**
- [ ] AC1: User can authenticate with valid credentials
- [ ] AC2: Invalid credentials return error message
- [ ] AC3: Session expires after 30 minutes
- [ ] Coverage ≥ 80%
- [ ] mypy --strict passes
```

**Rule:** Workstream is NOT complete until all AC are checked.

**See Also:** Workstream, Quality Gates, Definition of Done

---

### Execution Report

**Definition:** Section in workstream file documenting actual execution results.

**Location:** Bottom of workstream file (## Execution Report)

**Contents:**
- Executor name/date
- Goal status (AC check)
- Files changed with LOC
- Commit hash

**Example:**
```markdown
## Execution Report
**Executed by:** claude-sonnet-4.5
**Date:** 2026-01-29

### Goal Status
- [x] AC1-AC3 — ✅

**Goal Achieved:** Yes

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| src/module.py | created | 145 |

### Commit
abc123def
```

**See Also:** Workstream, Acceptance Criteria, Commit

---

## Quality Gates

### Quality Gates

**Definition:** Mandatory checks that all workstreams must pass before completion.

**Purpose:** Ensure code quality, maintainability, and AI-readiness.

**Gates:**

| Gate | Requirement | Check Command |
|------|-------------|---------------|
| **AI-Readiness** | Files < 200 LOC, CC < 10, type hints | `find src/ -name "*.py" -exec wc -l {} + \| awk '$1 > 200'` |
| **Clean Architecture** | No layer violations | `grep -r "from.*infrastructure" src/domain/` |
| **Error Handling** | No `except: pass` | `grep -rn "except:" src/` |
| **Test Coverage** | ≥ 80% | `pytest tests/ --cov=src/ --cov-fail-under=80` |
| **Type Checking** | Strict mypy | `mypy src/ --strict` |
| **No TODOs** | All tasks completed | Manual review |

**See Also:** Forbidden Patterns, Required Patterns, @review

---

### AI-Readiness

**Definition:** Code characteristics that make it easy for AI models to understand and modify.

**Criteria:**
- Files < 200 LOC
- Cyclomatic Complexity < 10
- Full type hints
- Clear naming
- Small, focused functions

**Purpose:** Enable reliable AI agent code comprehension and modification.

**Enforcement:** Pre-build and post-build hooks

**See Also:** Quality Gates, Forbidden Patterns

---

### Forbidden Patterns

**Definition:** Code patterns that are explicitly prohibited in SDP.

**List:**
- ❌ `except: pass` or bare exceptions
- ❌ Time-based estimates
- ❌ Layer violations
- ❌ Files > 200 LOC
- ❌ TODO without followup WS
- ❌ Coverage < 80%
- ❌ Implicit type hints (Python 2 style)
- ❌ `Any` type without justification

**Enforcement:** Pre-build validation, post-build verification

**See Also:** Quality Gates, Required Patterns, Clean Architecture

---

### Required Patterns

**Definition:** Code patterns that are mandatory in SDP.

**List:**
- ✅ Type hints everywhere
- ✅ Tests first (TDD)
- ✅ Explicit error handling
- ✅ Clean architecture boundaries
- ✅ Conventional commits
- ✅ Docstrings for public APIs
- ✅ Protocol-based abstractions

**Enforcement:** Post-build validation

**See Also:** Quality Gates, Forbidden Patterns, TDD

---

### Test Coverage

**Definition:** Percentage of codebase executed by test suite.

**Requirement:** ≥ 80% for all workstreams.

**Measurement:**
```bash
pytest tests/ --cov=src/ --cov-fail-under=80
```

**Report:** `--cov-report=term-missing` shows uncovered lines

**Purpose:** Ensure code reliability and enable confident refactoring.

**See Also:** TDD, Quality Gates, pytest

---

### Type Checking

**Definition:** Static type analysis using mypy to ensure type correctness.

**Requirement:** `mypy --strict` must pass with no errors.

**Enforcement:**
```bash
mypy src/ --strict --no-implicit-optional
```

**Purpose:** Catch type errors before runtime, improve IDE support.

**See Also:** Quality Gates, Type Hints, mypy

---

## Architecture

### Clean Architecture

**Definition:** Layered architecture where dependencies point inward toward the domain.

**Layers (outside → inside):**
1. **Presentation** - UI layer (CLI, API, controllers)
2. **Infrastructure** - External concerns (DB, API, frameworks)
3. **Application** - Use cases (orchestration, services)
4. **Domain** - Business logic (entities, value objects)

**Dependency Rule:** Dependencies point INWARD

```
┌─────────────────────────────────────────────────────┐
│                  Presentation                        │
│              (Controllers, Views, API)               │
├─────────────────────────────────────────────────────┤
│                  Infrastructure                      │
│          (Database, External APIs, Frameworks)       │
├─────────────────────────────────────────────────────┤
│                   Application                        │
│              (Use Cases, Services)                   │
├─────────────────────────────────────────────────────┤
│                     Domain                           │
│           (Entities, Business Rules)                 │
└─────────────────────────────────────────────────────┘

        ↑ Dependencies point INWARD (toward Domain)
```

**Layer Rules:**

| Layer | Can Import From | Cannot Import From |
|-------|-----------------|-------------------|
| Domain | Nothing | Application, Infrastructure, Presentation |
| Application | Domain | Infrastructure, Presentation |
| Infrastructure | Domain, Application | Presentation |
| Presentation | Application | - |

**Violation Example:**
```python
# ❌ Layer violation - Domain importing Infrastructure
from src.infrastructure.persistence import Database

class UserEntity:
    def save(self):
        db = Database()  # Domain shouldn't know about DB
```

**Correct Example:**
```python
# ✅ Clean separation
class UserEntity:
    def __init__(self, name: str, email: str):
        self.name = name
        self.email = email
```

**See Also:** Dependency Inversion, Layer Violation, SOLID

---

### Layer Violation

**Definition:** When code in an inner layer imports from an outer layer.

**Detection:**
```bash
grep -r "from.*infrastructure" src/domain/
```

**Examples:**
- Domain entity importing Database class
- Application service importing FastAPI
- Repository imported into Domain layer

**Consequence:** Violates dependency inversion, creates tight coupling

**Fix:** Use Protocol-based abstractions (ports)

**See Also:** Clean Architecture, Dependency Inversion, Protocol

---

### Dependency Inversion

**Definition:** SOLID principle stating that high-level modules should depend on abstractions, not concretions.

**Purpose:** Enable decoupling and testability.

**Example:**
```python
# ❌ BAD: Direct dependency on concrete class
class OrderService:
    def __init__(self):
        self.db = MySQLDatabase()  # Tightly coupled

# ✅ GOOD: Depend on abstraction
class OrderRepository(Protocol):
    def save(self, order: Order) -> None: ...

class OrderService:
    def __init__(self, repository: OrderRepository):  # Abstract
        self.repository = repository
```

**See Also:** SOLID, Clean Architecture, Protocol

---

### Protocol

**Definition:** Python structural typing interface defining behavior without inheritance.

**Purpose:** Define abstractions for dependency injection.

**Example:**
```python
class SubmissionRepository(Protocol):
    def get_by_id(self, submission_id: str) -> Submission | None:
        ...
    def save(self, submission: Submission) -> None:
        ...
```

**Usage:** Enables loose coupling and testability.

**See Also:** Dependency Inversion, Clean Architecture, Type Hints

---

### Repository Pattern

**Definition:** Data access abstraction that isolates domain from database details.

**Components:**
- **Port** (Protocol in Application layer)
- **Adapter** (Implementation in Infrastructure layer)

**Example:**
```python
# Application layer (port)
class SubmissionRepository(Protocol):
    def get_by_id(self, submission_id: str) -> Submission | None: ...
    def save(self, submission: Submission) -> None: ...

# Infrastructure layer (adapter)
class PostgresSubmissionRepository:
    def __init__(self, session: Session) -> None:
        self._session = session

    def get_by_id(self, submission_id: str) -> Submission | None:
        row = self._session.query(SubmissionModel).filter_by(id=submission_id).first()
        return Submission.from_orm(row) if row else None
```

**See Also:** Clean Architecture, Protocol, Adapter

---

### Adapter Pattern

**Definition:** Converts external system interfaces to domain Protocol interfaces.

**Purpose:** Integrate external systems without violating clean architecture.

**Example:**
```python
class GCPStorageAdapter:
    """Adapts GCS to StoragePort interface."""

    def __init__(self, client: storage.Client, bucket_name: str) -> None:
        self._client = client
        self._bucket = self._client.bucket(bucket_name)

    def upload(self, local_path: Path, remote_path: str) -> str:
        blob = self._bucket.blob(remote_path)
        blob.upload_from_filename(str(local_path))
        return blob.public_url
```

**See Also:** Clean Architecture, Protocol, Repository

---

### Factory Pattern

**Definition:** Creation pattern that encapsulates object creation logic.

**Purpose:** Centralize complex object creation, enable polymorphism.

**Example:**
```python
class ExecutorFactory:
    def __init__(self, config: ExecutorConfig, logger: structlog.BoundLogger) -> None:
        self._config = config
        self._logger = logger

    def create(self, executor_type: ExecutorType) -> Executor:
        match executor_type:
            case ExecutorType.DIND:
                return DindExecutor(self._config.dind, self._logger)
            case ExecutorType.K8S:
                return K8sExecutor(self._config.k8s, self._logger)
            case _:
                raise ValueError(f"Unknown executor type: {executor_type}")
```

**See Also:** Design Patterns, Dependency Injection

---

### State Machine

**Definition:** Pattern for managing state transitions in a system.

**Components:**
- **States** (Enum)
- **Context** (Data)
- **Commands** (Transitions)
- **Orchestrator** (Executor)

**Example:**
```python
class CleanupState(Enum):
    INITIAL = auto()
    RUNNING = auto()
    COMPLETED = auto()
    FAILED = auto()

class CleanupOrchestrator:
    def run(self, ctx: CleanupContext) -> CleanupResult:
        state = CleanupState.INITIAL
        while state not in (CleanupState.COMPLETED, CleanupState.FAILED):
            cmd = self._commands.get(state)
            if cmd is None:
                break
            ctx, state = cmd.execute(ctx)
        return CleanupResult(success=state == CleanupState.COMPLETED)
```

**See Also:** Orchestrator, Context, Command Pattern

---

### Orchestrator

**Definition:** Coordinator that manages execution of multiple commands or workstreams.

**Types:**
- **Workstream Orchestrator:** Coordinates WS execution in `@oneshot`
- **Command Orchestrator:** Coordinates state machine transitions

**Purpose:** Coordinate complex workflows with dependencies.

**See Also:** @oneshot, State Machine, Multi-Agent

---

## Testing

### pytest

**Definition:** Python testing framework used by SDP.

**Usage:**
```bash
# Run all tests
pytest tests/ -v

# Run with coverage
pytest tests/ --cov=src/ --cov-fail-under=80

# Run specific file
pytest tests/unit/test_module.py -v

# Show missing coverage
pytest --cov=src/ --cov-report=term-missing
```

**Configuration:** `pytest.ini` or `pyproject.toml`

**See Also:** Test Coverage, TDD, Quality Gates

---

### Unit Test

**Definition:** Test that verifies individual units of code in isolation.

**Characteristics:**
- Fast execution
- No external dependencies (mocked)
- Tests single function/class
- Located in `tests/unit/`

**Example:**
```python
def test_user_can_be_created():
    user = User(email="test@example.com", name="Test")
    assert user.email == "test@example.com"
    assert user.name == "Test"
```

**See Also:** Integration Test, TDD, pytest

---

### Integration Test

**Definition:** Test that verifies interaction between multiple components.

**Characteristics:**
- Slower than unit tests
- Real external dependencies (DB, API)
- Tests component interactions
- Located in `tests/integration/`

**Example:**
```python
def test_user_can_be_saved_to_database():
    user = User(email="test@example.com", name="Test")
    repository = PostgresUserRepository(session)
    repository.save(user)

    retrieved = repository.get_by_id(user.id)
    assert retrieved.email == "test@example.com"
```

**See Also:** Unit Test, pytest, Repository

---

### Test Double

**Definition:** Generic term for any substitute object used in testing.

**Types:**
- **Dummy:** Object passed around but never used
- **Stub:** Provides canned answers to calls
- **Fake:** Working but simplified implementation
- **Mock:** Pre-programmed with expectations
- **Spy:** Records calls for verification

**Purpose:** Isolate code under test from dependencies.

**See Also:** Unit Test, Mocking

---

### Mocking

**Definition:** Technique for replacing real dependencies with test doubles.

**Tools:** `unittest.mock`, `pytest-mock`

**Example:**
```python
from unittest.mock import Mock

def test_send_notification():
    email_service = Mock()
    user = User(email="test@example.com", name="Test")

    user.send_welcome_email(email_service)

    email_service.send.assert_called_once_with(
        to="test@example.com",
        subject="Welcome"
    )
```

**See Also:** Test Double, Unit Test

---

## Multi-Agent

### Multi-Agent System

**Definition:** System where multiple AI agents collaborate to complete complex tasks.

**Components:**
- **Agent Spawner** - Creates agents dynamically
- **Message Router** - Routes messages between agents
- **Role Manager** - Assigns roles and capabilities
- **Notification Router** - Dispatches notifications

**Purpose:** Enable parallel, specialized task execution.

**See Also:** Agent, Orchestrator, Task Tool

---

### Agent

**Definition:** Autonomous AI entity with specific role and capabilities.

**Types:**
- **Planner** - Breaks features into workstreams
- **Builder** - Executes workstreams
- **Reviewer** - Quality checks
- **Deployer** - Production deployment
- **Orchestrator** - Coordinates workflow

**Location:** `.claude/agents/{role}.md`

**Environment Variables:**
- `CLAUDE_CODE_AGENT_ID` - Unique identifier
- `CLAUDE_CODE_AGENT_TYPE` - Role/type
- `CLAUDE_CODE_TEAM_NAME` - Team membership

**See Also:** Multi-Agent System, Orchestrator, Task Tool

---

### Task Tool

**Definition:** Claude Code tool for spawning isolated agents for autonomous execution.

**Purpose:** Enable background execution and resume capability.

**Usage:**
```python
# Spawn orchestrator agent
task = Task("Execute feature F01", subagent_type="orchestrator")
# Returns agent ID for resume
```

**Features:**
- Isolated execution context
- Checkpoint/resume support
- Background execution
- Progress tracking

**See Also:** @oneshot, Agent, Orchestrator, Checkpoint

---

### Checkpoint

**Definition:** Saved state after each workstream completion in autonomous mode.

**Purpose:** Enable resume from interruption.

**Format:** Progress summary with completed workstreams

**Usage:**
```bash
@oneshot F01 --resume abc123xyz
# Resumes from last checkpoint (WS-001-03)
```

**See Also:** @oneshot, Task Tool, Agent

---

### Unified Orchestrator

**Definition:** SDP v0.5+ component integrating multi-agent coordination with Beads task tracking.

**Components:**
```
┌─────────────────────────────────────────────────────────────┐
│                    Unified Orchestrator                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Agent Spawner│──│Message Router│──│ Role Manager │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│         │                  │                  │             │
│         ▼                  ▼                  ▼             │
│  ┌──────────────────────────────────────────────────┐     │
│  │              Notification Router                  │     │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────┐   │     │
│  │  │ Console  │  │ Telegram │  │    Mock      │   │     │
│  │  └──────────┘  └──────────┘  └──────────────┘   │     │
│  └──────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │  Beads CLI  │
                    │ Task Tracker│
                    └─────────────┘
```

**See Also:** Multi-Agent System, Beads, Orchestrator

---

### Message Router

**Definition:** Component that routes messages between agents in multi-agent system.

**Purpose:** Enable agent communication and coordination.

**Types:**
- **Direct Message** - One-to-one communication
- **Broadcast** - One-to-all communication
- **Notification** - Event-based messaging

**See Also:** Multi-Agent System, Unified Orchestrator, Agent

---

### Notification Router

**Definition:** Component that dispatches notifications to various channels.

**Channels:**
- **Console** - Terminal output
- **Telegram** - Bot notifications
- **Mock** - Testing/debugging

**Example:**
```python
notifier.send(Notification(
    type=NotificationType.SUCCESS,
    message="Feature completed successfully",
))
```

**See Also:** Unified Orchestrator, Telegram, Message Router

---

### Telegram Notification

**Definition:** Integration for sending SDP notifications via Telegram bot.

**Configuration:**
```bash
# .env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
```

**Usage:** Automatic notifications for:
- Feature completion
- Workstream failures
- Deployment status

**See Also:** Notification Router, Unified Orchestrator

---

## Principles

### SOLID

**Definition:** Five principles of object-oriented design for maintainable software.

**Acronym:**
- **S** - Single Responsibility Principle (SRP)
- **O** - Open/Closed Principle (OCP)
- **L** - Liskov Substitution Principle (LSP)
- **I** - Interface Segregation Principle (ISP)
- **D** - Dependency Inversion Principle (DIP)

**Purpose:** Create maintainable, scalable software.

**See Also:** [docs/PRINCIPLES.md](docs/PRINCIPLES.md), Clean Architecture

---

### SRP (Single Responsibility Principle)

**Definition:** A class should have only one reason to change.

**Example:**
```python
# ❌ BAD: Multiple responsibilities
class UserService:
    def authenticate(self, email: str, password: str) -> User: ...
    def send_welcome_email(self, user: User) -> None: ...
    def generate_report(self, user: User) -> str: ...

# ✅ GOOD: Separated concerns
class AuthService:
    def authenticate(self, email: str, password: str) -> User: ...

class EmailService:
    def send_welcome(self, user: User) -> None: ...

class ReportGenerator:
    def generate_user_report(self, user: User) -> str: ...
```

**See Also:** SOLID, OCP

---

### OCP (Open/Closed Principle)

**Definition:** Open for extension, closed for modification.

**Example:**
```python
# ❌ BAD: Modifying existing code for new types
class PaymentProcessor:
    def process(self, payment_type: str, amount: float) -> None:
        if payment_type == "credit_card":
            self._process_credit_card(amount)
        elif payment_type == "paypal":
            self._process_paypal(amount)
        elif payment_type == "crypto":  # Adding new type = modifying class
            self._process_crypto(amount)

# ✅ GOOD: Extending via new classes
class PaymentProcessor(Protocol):
    def process(self, amount: float) -> None: ...

class CryptoProcessor:  # New type = new class, no modification
    def process(self, amount: float) -> None: ...
```

**See Also:** SOLID, Strategy Pattern

---

### LSP (Liskov Substitution Principle)

**Definition:** Subtypes must be substitutable for their base types.

**Example:**
```python
# ❌ BAD: Subtype breaks parent contract
class Square(Rectangle):
    def set_width(self, width: int) -> None:
        self.width = width
        self.height = width  # Unexpected side effect

# ✅ GOOD: Proper abstraction
class Shape(Protocol):
    def area(self) -> float: ...
```

**See Also:** SOLID, Protocol

---

### ISP (Interface Segregation Principle)

**Definition:** Clients should not depend on interfaces they don't use.

**Example:**
```python
# ❌ BAD: Fat interface
class IWorker(Protocol):
    def work(self) -> None: ...
    def eat(self) -> None: ...

class Robot:  # Robots don't eat!
    def eat(self) -> None:
        raise NotImplementedError  # Violation!

# ✅ GOOD: Segregated interfaces
class IWorkable(Protocol):
    def work(self) -> None: ...

class Robot:
    def work(self) -> None: ...  # Only implements what it needs
```

**See Also:** SOLID, Protocol

---

### DIP (Dependency Inversion Principle)

**Definition:** Depend on abstractions, not concretions.

**Example:**
```python
# ❌ BAD: Direct dependency on concrete class
class OrderService:
    def __init__(self):
        self.db = MySQLDatabase()  # Tightly coupled

# ✅ GOOD: Depend on abstraction
class OrderRepository(Protocol):
    def save(self, order: Order) -> None: ...

class OrderService:
    def __init__(self, repository: OrderRepository):  # Abstract
        self.repository = repository
```

**See Also:** SOLID, Protocol, Clean Architecture

---

### DRY (Don't Repeat Yourself)

**Definition:** Every piece of knowledge must have a single, unambiguous representation.

**Example:**
```python
# ❌ BAD: Repeated validation
def create_user(email: str) -> User:
    if "@" not in email or "." not in email:
        raise InvalidEmailError(email)

def update_email(user: User, email: str) -> None:
    if "@" not in email or "." not in email:  # Duplicated!
        raise InvalidEmailError(email)

# ✅ GOOD: Single source of truth
class Email:
    def __init__(self, value: str):
        if not self._is_valid(value):
            raise InvalidEmailError(value)
        self.value = value

    @staticmethod
    def _is_valid(value: str) -> bool:
        return "@" in value and "." in value
```

**See Also:** KISS, YAGNI

---

### KISS (Keep It Simple, Stupid)

**Definition:** The simplest solution is usually the best.

**Example:**
```python
# ❌ BAD: Over-engineered
def is_palindrome(s: str) -> bool:
    import re
    cleaned = re.sub(r'[^a-zA-Z0-9]', '', s).lower()
    stack = []
    for char in cleaned:
        stack.append(char)
    reversed_str = ''
    while stack:
        reversed_str += stack.pop()
    return cleaned == reversed_str

# ✅ GOOD: Simple and clear
def is_palindrome(s: str) -> bool:
    cleaned = ''.join(c.lower() for c in s if c.isalnum())
    return cleaned == cleaned[::-1]
```

**See Also:** DRY, YAGNI

---

### YAGNI (You Aren't Gonna Need It)

**Definition:** Don't build features until they're actually needed.

**Example:**
```python
# ❌ BAD: Building for hypothetical future
class Config:
    def __init__(
        self,
        database_url: str,
        cache_url: str | None = None,  # "Might need cache later"
        message_queue_url: str | None = None,  # "Might need queue later"
    ):
        ...

# ✅ GOOD: Build only what's needed NOW
class Config:
    def __init__(self, database_url: str):
        self.database_url = database_url

# Add cache_url when you actually need caching
```

**Guardrails:**
- Implement requirements **only**
- No "nice to have" features
- No "we might need this later"
- Delete unused code immediately

**See Also:** KISS, DRY

---

## Tools & Integrations

### Claude Code

**Definition:** Anthropic's CLI tool for AI-assisted software development.

**Website:** https://claude.com/claude-code

**Key Features:**
- Multi-file editing
- Bash execution
- Tool invocation (Read, Write, Edit, Glob, Grep, Task)
- Skill system (@ and / commands)
- Multi-agent orchestration

**Model Selection:**
```bash
/model opus    # Most capable - for /design, /review, /oneshot
/model sonnet  # Balanced - for /idea, /issue
/model haiku   # Fastest - for /build, /deploy, /hotfix, /bugfix
```

**See Also:** [CLAUDE.md](CLAUDE.md), Skill, @ Command Prefix

---

### Cursor

**Definition:** AI-powered code editor by Cursor Inc.

**Website:** https://cursor.sh

**Key Features:**
- AI chat inline
- Multi-file editing
- Command palette
- Agent system

**Integration:** SDP provides skills and workflows for Cursor.

**See Also:** [docs/guides/CURSOR.md](docs/guides/CURSOR.md), Claude Code

---

### mypy

**Definition:** Static type checker for Python.

**Requirement:** `mypy --strict` must pass

**Usage:**
```bash
mypy src/ --strict --no-implicit-optional
```

**Purpose:** Catch type errors before runtime.

**See Also:** Type Checking, Type Hints, Quality Gates

---

### ruff

**Definition:** Fast Python linter written in Rust.

**Usage:**
```bash
ruff check src/ --select=C901  # Complexity check
```

**Purpose:** Enforce code style and catch complexity issues.

**See Also:** Quality Gates, Linting

---

### structlog

**Definition:** Structured logging library for Python.

**Usage:**
```python
import structlog
log = structlog.get_logger()

log.info("operation.start", submission_id=sid)
log.error("operation.failed", error=str(e), exc_info=True)
```

**Purpose:** Consistent, structured logging across codebase.

**See Also:** Logging, Context

---

### Git Hooks

**Definition:** Scripts that run automatically at specific Git events.

**Location:** `hooks/`

**Types:**
- `pre-build.sh` - Validation before workstream execution
- `post-build.sh` - Validation after workstream completion
- `pre-commit` - Validation before commit
- `pre-push` - Validation before push

**Purpose:** Enforce quality gates automatically.

**See Also:** Quality Gates, Validation

---

### Beads CLI

**Definition:** Command-line interface for Beads task tracking system.

**Usage:**
```bash
bd task create --title "Feature name"
bd task list --status ready
bd task dependency add --blocks ws2.id ws1.id
```

**Purpose:** Manage tasks, dependencies, and status.

**See Also:** Beads, Task Tracking, Unified Orchestrator

---

### Docker Compose

**Definition:** Tool for defining and running multi-container Docker applications.

**Generated By:** `@deploy` skill

**Purpose:** Local development and production deployment.

**See Also:** @deploy, Deployment

---

### GitHub Actions

**Definition:** CI/CD platform integrated with GitHub.

**Generated By:** `@deploy` skill (`.github/workflows/deploy.yml`)

**Purpose:** Automated testing, deployment, and quality checks.

**See Also:** @deploy, CI/CD

---

## Severity Classification

### P0 (Critical)

**Definition:** Critical issue requiring immediate fix (< 2 hours).

**Examples:**
- System down
- Data loss
- Security breach
- Complete feature failure

**Workflow:** `@issue` → `@hotfix` → production

**See Also:** @hotfix, P1, P2, P3

---

### P1 (High)

**Definition:** High-priority bug requiring fix within 24 hours.

**Examples:**
- Major feature broken
- Significant performance degradation
- User workflow blocked

**Workflow:** `@issue` → `@bugfix` → `@review` → `@deploy`

**See Also:** @bugfix, P0, P2

---

### P2 (Medium)

**Definition:** Medium-priority bug requiring fix within 1 week.

**Examples:**
- Minor feature broken
- Annoying but workaround exists
- Edge case failure

**Workflow:** `@issue` → `@bugfix` → `@review` → `@deploy`

**See Also:** @bugfix, P0, P1

---

### P3 (Low)

**Definition:** Low-priority issue that can be deferred.

**Examples:**
- Cosmetic issues
- Nice-to-have improvements
- Rare edge cases

**Workflow:** `@issue` → backlog → future sprint

**See Also:** @issue, Backlog, P0, P1, P2

---

## File Conventions

### Conventional Commits

**Definition:** Standardized commit message format.

**Format:**
```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `refactor` - Code refactoring
- `test` - Adding tests
- `docs` - Documentation
- `chore` - Maintenance tasks

**Example:**
```
feat(auth): add OAuth2 login flow

Implement OAuth2 authentication with Google and GitHub providers.

Closes #123
```

**See Also:** Commit, Git Workflow

---

### `.env` File

**Definition:** Environment configuration file.

**Location:** Project root

**Purpose:** Store sensitive configuration (API keys, tokens)

**Example:**
```bash
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
DATABASE_URL=postgresql://...
```

**Security:** Never commit `.env` files (in `.gitignore`)

**See Also:** Configuration, Secrets

---

### `pyproject.toml`

**Definition:** Python project configuration file.

**Purpose:**
- Dependencies
- Tool configuration (pytest, mypy, ruff)
- Project metadata

**Example:**
```toml
[tool.pytest.ini_options]
testpaths = ["tests"]
addopts = "--cov=src/ --cov-fail-under=80"

[tool.mypy]
strict = true

[tool.ruff]
select = ["E", "F", "C901"]
```

**See Also:** Python Packaging, Tool Configuration

---

### Directory Structure

**Definition:** Standard organization of SDP project files.

**Example:**
```
project/
├── PRODUCT_VISION.md      # Project manifesto
├── docs/
│   ├── schema/            # Intent JSON schema
│   ├── intent/            # Machine-readable intent files
│   ├── drafts/            # @idea outputs
│   ├── workstreams/
│   │   ├── backlog/       # @design outputs
│   │   ├── in_progress/   # @build moves here
│   │   └── completed/     # @build finalizes here
│   └── specs/             # Feature specifications
├── src/sdp/
│   ├── schema/            # Intent validation
│   ├── tdd/               # TDD cycle runner
│   ├── feature/           # Product vision management
│   └── design/            # Dependency graph
├── prompts/commands/      # Skill instructions
├── .claude/
│   ├── skills/            # Skill definitions
│   ├── agents/            # Multi-agent mode
│   └── settings.json      # Claude Code settings
└── hooks/                 # Git hooks
```

**See Also:** Workstream, Intent, Draft

---

## Advanced Concepts

### Veto Protocol

**Definition:** Mechanism for agents to block work that violates rules.

**Used By:** Architect, Security, DevOps, Tech Lead, QA agents

**Format:**
```json
{
  "d": "2025-12-09",
  "st": "veto",
  "r": "architect",
  "epic": "EP08",
  "sm": [
    "Violation: layer_violation",
    "Location: application/foo.py imports infrastructure",
    "Impact: Breaks dependency inversion"
  ],
  "nx": ["Remove direct import", "Use port injection"]
}
```

**Non-overrideable Vetoes:**
- Architecture violations (architect)
- Security issues (security)
- Missing rollback plan (devops)
- Code review violations (tech_lead, qa)

**See Also:** [RULES_COMMON.md](RULES_COMMON.md), Agent, Quality Gates

---

### Handoff Protocol

**Definition:** Process for transferring work between agents or preserving context for next epic.

**Purpose:** Preserve key decisions, technical debt, and lessons learned.

**Contents:**
- Key decisions
- Technical debt identified
- Lessons learned
- Foundation for next work

**See Also:** [RULES_COMMON.md](RULES_COMMON.md), Agent, Epic

---

### Inbox Message

**Definition:** Communication mechanism between agents in multi-agent system.

**Location:** `messages/inbox/{role}/{date}-{subject}.json`

**Format:**
```json
{
  "d": "2025-12-09",
  "st": "status|request|veto|approval|handoff",
  "r": "developer",
  "epic": "EP08",
  "sm": ["summary points"],
  "nx": ["next actions"],
  "artifacts": ["paths"],
  "answers": {"Q1": "A1"}
}
```

**Rules:**
- READ: Only your own inbox
- WRITE: Only to OTHER agents' inboxes
- All messages in English

**See Also:** [RULES_COMMON.md](RULES_COMMON.md), Agent, Multi-Agent

---

### Epic (Legacy)

**Definition:** **DEPRECATED** - Replaced by Feature in modern SDP.

**Status:** Use `PP-FFF` (Feature) instead of `EP-XX` (Epic)

**Migration:** `EP08` → `00-008` or similar

**See Also:** Feature, Release, Naming Conflicts

---

### Sprint (Legacy)

**Definition:** **DEPRECATED** - SDP does not use time-based iterations.

**Reasoning:** SDP uses scope-based estimation (LOC/tokens), not time.

**Replacement:** Workstream (scope-based)

**See Also:** Workstream, YAGNI

---

### Progressive Disclosure

**Definition:** UX pattern where information is revealed incrementally.

**Usage:** `@feature` skill uses progressive disclosure to:
1. Start with high-level vision
2. Reveal details as needed
3. Avoid overwhelming users

**See Also:** @feature, UX

---

### AskUserQuestion

**Definition:** Tool for deep interviewing in Claude Code skills.

**Purpose:** Ask non-obvious questions exploring tradeoffs.

**Used By:** `@idea`, `@design`

**Characteristics:**
- No yes/no questions
- Explores technical approaches
- Uncovers concerns
- Validates assumptions

**See Also:** @idea, @design, Skill

---

### EnterPlanMode

**Definition:** State in which codebase exploration is allowed.

**Used By:** `@design` skill

**Purpose:** Enable agents to read files and explore codebase before planning.

**Flow:**
1. `@design` invoked
2. EnterPlanMode (codebase exploration allowed)
3. Explore codebase
4. AskUserQuestion for architecture decisions
5. ExitPlanMode (require approval)

**See Also:** @design, ExitPlanMode, AskUserQuestion

---

### ExitPlanMode

**Definition:** State transition requesting approval for proposed plan.

**Used By:** `@design` skill

**Purpose:** Require human approval before proceeding with implementation.

**Flow:**
1. Plan created
2. ExitPlanMode (request approval)
3. Human approves
4. Plan accepted, workstreams created

**See Also:** @design, EnterPlanMode, Approval

---

### TodoWrite

**Definition:** Claude Code tool for real-time progress tracking.

**Used By:** `@build` skill

**Purpose:** Show real-time progress during workstream execution.

**Example:**
```
✓ [completed] Pre-build validation
✓ [in_progress] Write failing test (Red)
• [pending] Implement minimum code (Green)
```

**See Also:** @build, Progress Tracking

---

### Background Execution

**Definition:** Execution mode for long-running features that doesn't block user.

**Usage:**
```bash
@oneshot F01 --background
```

**Purpose:** Allow user to continue working while feature executes.

**Features:**
- Non-blocking
- Progress logging to file
- Notification on completion

**See Also:** @oneshot, Task Tool, Agent

---

---

## Index

### By Category

**Core Concepts:** SDP, Workstream, Feature, Release, Beads

**Hierarchy:** Three-Level Hierarchy, Release, Feature, Workstream

**Identifiers:** PP-FFF-SS, Feature ID, beads-sdp-XXX, Idea ID, Naming Conflicts

**Commands:** @feature, @idea, @design, @build, @review, @deploy, @hotfix, @bugfix, @issue, /debug, /tdd, @oneshot

**Workflow:** TDD, Red-Green-Refactor, Intent, Draft, Specification, Product Vision, Acceptance Criteria, Execution Report

**Quality:** Quality Gates, AI-Readiness, Forbidden Patterns, Required Patterns, Test Coverage, Type Checking

**Architecture:** Clean Architecture, Layer Violation, Dependency Inversion, Protocol, Repository, Adapter, Factory, State Machine, Orchestrator

**Testing:** pytest, Unit Test, Integration Test, Test Double, Mocking

**Multi-Agent:** Multi-Agent System, Agent, Task Tool, Checkpoint, Unified Orchestrator, Message Router, Notification Router, Telegram

**Principles:** SOLID, SRP, OCP, LSP, ISP, DIP, DRY, KISS, YAGNI

**Tools:** Claude Code, Cursor, mypy, ruff, structlog, Git Hooks, Beads CLI, Docker Compose, GitHub Actions

**Severity:** P0, P1, P2, P3

**Files:** Conventional Commits, .env, pyproject.toml, Directory Structure

**Advanced:** Veto Protocol, Handoff Protocol, Inbox Message, Epic (Legacy), Sprint (Legacy)

**UX:** Progressive Disclosure, AskUserQuestion, EnterPlanMode, ExitPlanMode, TodoWrite, Background Execution

---

## Contributing

To add or update terms:

1. Read existing glossary to avoid conflicts
2. Add term in appropriate section
3. Include definition, example, cross-references
4. Update index
5. Follow existing format

---

**Version:** 1.0
**Last Updated:** 2026-01-29
**Total Terms:** 150+

---

**See Also:**
- [PROTOCOL.md](PROTOCOL.md) — Full SDP specification
- [CODE_PATTERNS.md](CODE_PATTERNS.md) — Code patterns reference
- [docs/PRINCIPLES.md](docs/PRINCIPLES.md) — Engineering principles
- [CLAUDE.md](CLAUDE.md) — Claude Code integration guide
