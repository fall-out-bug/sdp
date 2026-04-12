# SDP: Spec-Driven Protocol

**Workstream-driven development** for AI agents with an explicit leaf execution boundary.

---

## Multi-Level Architecture

SDP is designed as a multi-level product. Each level builds on the previous, but works independently.

```
┌─────────────────────────────────────────────────────────────────┐
│  L4: Collaboration (Notifications, Cross-Review)               │
├─────────────────────────────────────────────────────────────────┤
│  L3: Orchestration (Distributed Agents, k8s) — Future          │
├─────────────────────────────────────────────────────────────────┤
│  L2: Go Tools (Evidence Log, Guard, Checkpoints)               │
├─────────────────────────────────────────────────────────────────┤
│  L1: Adapters (Claude Code, Cursor, Windsurf invocation)       │
├─────────────────────────────────────────────────────────────────┤
│  L0: Protocol (THIS DOCUMENT)                                  │
│  ├── Workstream format + Quality gates + TDD                   │
│  ├── Skills (@build, @review, @oneshot, etc.)                  │
│  ├── Agent roles (implementer, reviewer, etc.)                 │
│  └── Beads integration (bd create/close/sync)                  │
└─────────────────────────────────────────────────────────────────┘
```

**Key Principle:** L0 works with ANY AI (Opus, GLM, Codex) in ANY tool (Claude Code, Cursor, Windsurf).

### Level Descriptions

| Level | What It Provides | Required? |
|-------|------------------|-----------|
| **L0** | Protocol, skills, agents, beads | Yes (foundation) |
| **L1** | Tool-specific invocation adapters | Optional |
| **L2** | Go CLI: evidence, guard, checkpoints | Optional |
| **L3** | Distributed orchestration | Future |
| **L4** | AI-Human collaboration features | Future |

### Skills in L0

Skills are LLM-agnostic descriptions of workflows:

```
@build 00-001-01    # Execute executable leaf workstream with TDD
@review <feature-id>        # Multi-agent quality review
@oneshot <feature-id>       # Autonomous feature execution
@deploy <feature-id>        # Create PR and merge
```

Each skill describes WHAT to do. L1 adapters provide HOW to invoke (Task tool, agent panel, etc.).

### Beads in L0

Task tracking works without Go tools:

```bash
bd create --title="Fix bug" --priority=1
bd close sdp-xxx
bd sync
```

Skills target workstream IDs. Runtime binds executable leaf workstreams to their
live Beads issues.

---

## Quick Start

```bash
# Install (WS-067-11: corrected path to sdp-plugin)
go install github.com/fall-out-bug/sdp/sdp-plugin/cmd/sdp@latest

# Create feature (interactive)
@feature "Add user authentication"

# Plan workstreams
@design idea-auth

# Execute executable leaf workstream
@build 00-001-01

# Or execute all autonomously
@oneshot <feature-id>

# Review quality
@review <feature-id>

# Deploy to production
@deploy <feature-id>
```

---

## Core Concepts

### Hierarchy

| Level | Scope | Size | Example |
|-------|-------|------|---------|
| **Release** | Product milestone | 10-30 Features | R1: Submissions E2E |
| **Feature** | Major feature | 5-30 Workstreams | F24: Unified Workflow |
| **Aggregate Workstream** | Non-executable container or roll-up | 2+ child leaves | WS-060: API Contract Roll-up |
| **Leaf Workstream** | Atomic executable task | SMALL/MEDIUM/LARGE | WS-061: Domain Model |

Hard rule:

- only `leaf` workstreams are executable
- aggregate workstreams exist for decomposition and roll-up, not direct dispatch
- maximum nesting depth is one aggregate layer

### Workstream Size

- **SMALL**: < 500 LOC, < 1500 tokens
- **MEDIUM**: 500-1500 LOC, 1500-5000 tokens
- **LARGE**: > 1500 LOC → split into 2+ WS

⚠️ **NO TIME-BASED ESTIMATES** - Use scope metrics (LOC/tokens) only.

---

## Workstream Flow

```
┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│  ANALYZE   │───→│    PLAN    │───→│  EXECUTE   │───→│   REVIEW   │
│  (Sonnet)  │    │  (Sonnet)  │    │   (Auto)   │    │  (Sonnet)  │
└────────────┘    └────────────┘    └────────────┘    └────────────┘
     │                  │                  │                  │
     ▼                  ▼                  ▼                  ▼
  Map WS           Plan WS            Code           APPROVED/FIX
```

---

## Quality Gates

Every executable leaf workstream must pass:

```bash
# Test coverage ≥80%
pytest tests/unit/ --cov=src/ --cov-fail-under=80

# Type checking
mypy src/ --strict

# Linting
ruff check src/

# All files <200 LOC
find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200'
```

**Forbidden:**
- ❌ `except: pass` or bare exceptions
- ❌ Files > 200 LOC
- ❌ Coverage < 80%
- ❌ Time-based estimates
- ❌ TODO without followup WS

---

## Unified Workflow (AI-Comm + Beads)

SDP v0.5+ integrates multi-agent coordination with task tracking.

### Components

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

### Agent Coordination

```python
from sdp.unified.agent.spawner import AgentSpawner, AgentConfig

# Spawn agents
spawner = AgentSpawner()
builder = spawner.spawn_agent(AgentConfig(
    name="builder",
    prompt="Execute executable leaf workstreams with TDD...",
))

# Send messages
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="orchestrator",
    content="Execute 00-060-01",
    recipient=builder,
))
```

### Beads Integration

```python
from sdp.beads import create_beads_client
from sdp.beads.models import BeadsTaskCreate, BeadsStatus

# Create client
client = create_beads_client(use_mock=True)

# Create execution issues for leaf workstreams.
# Workstream topology still lives in the SDP workstream files, not in Beads.
leaf_issue_1 = client.create_task(BeadsTaskCreate(
    title="WS-061: Domain model",
    description="Primary execution issue for executable leaf workstream",
))
leaf_issue_2 = client.create_task(BeadsTaskCreate(
    title="WS-062: Database schema",
    description="Primary execution issue for executable leaf workstream",
))

# Add dependency
client.add_dependency(leaf_issue_2.id, leaf_issue_1.id, dep_type="blocks")

# Update status
client.update_task_status(leaf_issue_1.id, BeadsStatus.CLOSED)

# Get ready tasks
ready = client.get_ready_tasks()  # [ws2.id]
```

### Telegram Notifications

```bash
# .env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
```

```python
from sdp.unified.notifications.telegram import TelegramConfig, TelegramNotifier
from sdp.unified.notifications.provider import Notification, NotificationType

config = TelegramConfig(
    bot_token=os.getenv("TELEGRAM_BOT_TOKEN"),
    chat_id=os.getenv("TELEGRAM_CHAT_ID"),
)
notifier = TelegramNotifier(config=config)

# Send notification
notifier.send(Notification(
    type=NotificationType.SUCCESS,
    message="Feature completed successfully",
))
```

---

## Feature Development Flow

### 1. Requirements (@feature skill)

```bash
@feature "Add user authentication"
```

Claude asks deep questions:
- Technical approach (JWT vs sessions?)
- UI/UX requirements
- Database schema
- Testing strategy
- Security concerns

→ Creates: `docs/intent/sdp-XXX.json`
→ Creates: `docs/drafts/beads-sdp-XXX.md`

### 2. Planning (@design skill)

```bash
@design beads-sdp-XXX
```

Claude explores codebase and creates workstreams:
- 00-XXX-01: Domain model (450 LOC)
- 00-XXX-02: Database schema (300 LOC)
- 00-XXX-03: Repository layer (500 LOC)
- 00-XXX-04: Service layer (600 LOC)
- 00-XXX-05: API endpoints (400 LOC)

→ Creates: `docs/workstreams/beads-sdp-XXX.md`

### 3. Contract Tests (@test skill)

```bash
@test 00-XXX-01
```

Generate contract tests that define **immutable interfaces**:

- **Function signatures** - Stable API contracts
- **Input/output contracts** - Data format specifications
- **Error conditions** - Expected failure modes
- **Invariants** - Business rules that must hold

**Workflow:**
1. Analyze interface requirements from spec
2. Design test contracts (signatures, I/O, errors, invariants)
3. Create contract test file: `tests/contract/test_{component}.py`
4. Get stakeholder approval
5. **Lock contracts** - once approved, they CANNOT be modified during /build

**⚠️ Contract Immutability:**
- ✅ `/build` CAN implement code to pass contracts
- ❌ `/build` CANNOT modify contract test files
- ❌ `/build` CANNOT change function signatures
- ❌ `/build` CANNOT relax error conditions

**If interface change is needed:**
1. Stop `/build`
2. Create new workstream: "Update contract for {Component}"
3. Run `/test` with revised contracts
4. Get explicit approval
5. Resume `/build`

Creates: `tests/contract/test_{component}.py`

### 4. Implementation (@build skill)

```bash
@build 00-XXX-01
```

Claude follows TDD:
1. **Red** - Write failing test
2. **Green** - Implement minimum code
3. **Refactor** - Improve design

→ Shows real-time progress
→ Runs tests, mypy, ruff
→ Commits when complete

**⚠️ Contract Test Enforcement:**
- Guard prevents editing contract test files during `/build`
- Interface changes require new `/test` cycle

### 5. Autonomous Execution (@oneshot skill)

```bash
@oneshot sdp-XXX
```

Orchestrator agent:
- Executes all WS in dependency order
- Saves checkpoints after each WS
- Sends Telegram notifications
- Resumes from interruption

### 6. Quality Review (@review skill)

```bash
@review sdp-XXX
```

Validates:
- ✅ All quality gates passed
- ✅ Tests ≥80% coverage
- ✅ No tech debt
- ✅ Clean architecture

→ Returns: APPROVED / CHANGES_REQUESTED

### 7. Deployment (@deploy skill)

```bash
@deploy sdp-XXX
```

Generates:
- `docker-compose.yml`
- `.github/workflows/deploy.yml`
- `CHANGELOG.md` entry
- Git tag: `v{version}`

---

## Guardrails

### YAGNI (You Aren't Gonna Need It)

- Implement requirements **only**
- No "nice to have" features
- No "we might need this later"
- Delete unused code immediately

### KISS (Keep It Simple, Stupid)

- Prefer simple solutions
- Avoid over-engineering
- No premature abstraction
- One-liner > function > class

### DRY (Don't Repeat Yourself)

- Extract duplicated code
- Create reusable utilities
- But avoid premature abstraction

### SOLID Principles

- **S**ingle Responsibility - One reason to change
- **O**pen/Closed - Open for extension, closed for modification
- **L**iskov Substitution - Subtypes must be substitutable
- **I**nterface Segregation - No fat interfaces
- **D**ependency Inversion - Depend on abstractions

---

## Workstream Naming Convention

**Format:** `PP-FFF-SS`

- **PP** - Product/Project (01-99)
- **FFF** - Feature number (001-999)
- **SS** - Workstream sequence (01-99)

**Examples:**
- `00-001-01` - First workstream of SDP feature 001
- `02-150-01` - First workstream of hw_checker feature 150

**Legacy terms (no longer used):**
- ~~`WS-FFF-SS`~~ — replaced by `PP-FFF-SS`
- ~~`Epic`~~ — replaced by **Feature**
- ~~`Sprint`~~ — not used

**Migration Features:**
- ✅ `--dry-run` mode for safe preview
- ✅ Updates frontmatter (`ws_id` and `project_id`)
- ✅ Renames files to match new format
- ✅ Updates cross-WS dependencies
- ✅ Comprehensive validation and error reporting
- ✅ Full test coverage (≥80%)

---

## Clean Architecture

```
src/
├── domain/          # Business logic (no framework deps)
│   ├── entities/    # Core business objects
│   └── value_objects/  # Immutable values
├── application/     # Use cases (orchestration)
│   └── services/    # Application services
├── infrastructure/  # External concerns (DB, API)
│   ├── persistence/ # Database access
│   └── api/         # Controllers, views
└── presentation/    # UI layer (optional)
```

**Rules:**
- Domain ← No dependencies on other layers
- Application ← Can use Domain
- Infrastructure ← Can use Domain, Application
- Presentation ← Can use all layers

**Forbidden:**
```python
# ❌ Layer violation
from src.infrastructure.persistence import Database

class UserEntity:
    def save(self):
        db = Database()  # Domain shouldn't know about DB
```

```python
# ✅ Clean separation
class UserEntity:
    def __init__(self, name: str, email: str):
        self.name = name
        self.email = email
```

---

## Error Handling

**Forbidden:**
```python
# ❌ Bare except
try:
    risky_operation()
except:
    pass  # SWALLOWS ALL ERRORS
```

**Required:**
```python
# ✅ Explicit error handling
try:
    risky_operation()
except SpecificError as e:
    logger.error(f"Failed: {e}")
    raise  # Re-raise or handle
```

---

## Quick Reference

### Commands

```bash
# Development
@feature "title"           # Gather requirements
@design beads-XXX          # Plan workstreams
@build 00-XXX-01          # Execute executable leaf workstream
@oneshot FXXX             # Autonomous execution
@review FXXX              # Quality review
@deploy FXXX              # Production deployment

# Debugging
/debug "<issue>"           # Systematic debugging

# Issue routing
/issue "<bug report>"      # Classify and route bugs
@hotfix "<P0 issue>"       # Emergency fix <2h
@bugfix "<P1/P2 issue>"    # Quality fix <24h
```

### Quality Checks

```bash
# AI-Readiness
find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200'
ruff check src/ --select=C901  # Complexity

# Clean Architecture
grep -r "from.*infrastructure" src/domain/

# Error handling
grep -rn "except:" src/
grep -rn "except Exception" src/ | grep -v "exc_info"

# Test coverage
pytest tests/ --cov=src/ --cov-fail-under=80

# Full test suite
pytest -x --tb=short
pytest --cov=src/ --cov-report=term-missing
```

---

## Feature Branch Rule

**CRITICAL:** Features MUST be implemented in feature branches.

### Allowed Branches

| Branch Type | Purpose | Example |
|-------------|---------|---------|
| `feature/<id>` | Feature implementation | `feature/auth-login` |
| `bugfix/issue-id` | Bug fixes | `bugfix/sdp-1234` |
| `hotfix/issue-id` | Emergency fixes | `hotfix/sdp-1234` |

### Protected Branches

| Branch | Allowed Operations |
|--------|-------------------|
| `main` | Merge only (via PR) |
| `dev` | Merge only (via PR) |

### Enforcement

- Guard rejects commits to protected branches when `feature_id` is active
- `@build` verifies feature branch before starting work
- Pre-commit hooks block direct commits to `dev`/`main` for feature work

### Commands

```bash
# Check if current branch is valid for feature
sdp guard branch check --feature=<feature-id>

# Validate branch naming convention
sdp guard branch validate feature/<feature-id>
```

### Error Recovery

If you're on `dev` or `main` when you should be on a feature branch:

```bash
# Create feature branch
git checkout -b feature/<feature-id>

# Or switch to existing branch
git checkout feature/<feature-id>
```

---

## Documentation

- `.claude/agents/README.md` - Agent roles guide
- `README.md` - Project overview

---

## Version

**SDP v0.9.8** — Multi-Agent Architecture

Updated: 2026-02-26

---

**See Also:**
- Agent Roles: `.claude/agents/README.md`
- Reference: `docs/reference/`
- Schema Registry: `docs/reference/schema-registry.md`
- Integration Contracts Guide: `docs/reference/integration-contracts.md`
