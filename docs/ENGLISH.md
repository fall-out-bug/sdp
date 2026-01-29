# SDP: Spec-Driven Protocol (English Summary)

**Workstream-driven development** for AI agents with multi-agent coordination.

---

## Quick Start

```bash
# Install
pip install -e .

# Create feature (interactive)
@feature "Add user authentication"

# Plan workstreams
@design beads-auth

# Execute workstream
@build WS-AUTH-01

# Or execute all autonomously
@oneshot beads-auth

# Review quality
@review beads-auth

# Deploy to production
@deploy beads-auth
```

---

## Core Concepts

### Hierarchy

| Level | Scope | Size | Example |
|-------|-------|------|---------|
| **Release** | Product milestone | 10-30 Features | R1: Submissions E2E |
| **Feature** | Major feature | 5-30 Workstreams | F24: Unified Workflow |
| **Workstream** | Atomic task | SMALL/MEDIUM/LARGE | WS-060: Domain Model |

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

Every workstream must pass:

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
    prompt="Execute workstreams with TDD...",
))

# Send messages
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="orchestrator",
    content="Execute WS-060-01",
    recipient=builder,
))
```

### Beads Integration

```python
from sdp.beads import create_beads_client
from sdp.beads.models import BeadsTaskCreate, BeadsStatus

# Create client
client = create_beads_client(use_mock=True)

# Create feature
feature = client.create_task(BeadsTaskCreate(
    title="User Authentication",
    description="Add OAuth2 login flow",
))

# Decompose into workstreams
ws1 = client.create_task(BeadsTaskCreate(
    title="Domain model",
    parent_id=feature.id,
))
ws2 = client.create_task(BeadsTaskCreate(
    title="Database schema",
    parent_id=feature.id,
))

# Add dependency
client.add_dependency(ws2.id, ws1.id, dep_type="blocks")

# Update status
client.update_task_status(ws1.id, BeadsStatus.CLOSED)

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
- WS-XXX.01: Domain model (450 LOC)
- WS-XXX.02: Database schema (300 LOC)
- WS-XXX.03: Repository layer (500 LOC)
- WS-XXX.04: Service layer (600 LOC)
- WS-XXX.05: API endpoints (400 LOC)

→ Creates: `docs/workstreams/beads-sdp-XXX.md`

### 3. Execution (@build skill)

```bash
@build WS-XXX.01
```

Claude follows TDD:
1. **Red** - Write failing test
2. **Green** - Implement minimum code
3. **Refactor** - Improve design

→ Shows real-time progress
→ Runs tests, mypy, ruff
→ Commits when complete

### 4. Autonomous Execution (@oneshot skill)

```bash
@oneshot sdp-XXX
```

Orchestrator agent:
- Executes all WS in dependency order
- Saves checkpoints after each WS
- Sends Telegram notifications
- Resumes from interruption

### 5. Quality Review (@review skill)

```bash
@review sdp-XXX
```

Validates:
- ✅ All quality gates passed
- ✅ Tests ≥80% coverage
- ✅ No tech debt
- ✅ Clean architecture

→ Returns: APPROVED / CHANGES_REQUESTED

### 6. Deployment (@deploy skill)

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
- `WS-001-01` - First workstream of first feature
- `sdp-118.22` - Workstream 22 of SDP-118

**Legacy:**
- ~~Epic~~ → **Feature**
- ~~Sprint~~ → Not used

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
@build WS-XXX-01          # Execute workstream
@oneshot beads-XXX        # Autonomous execution
@review beads-XXX         # Quality review
@deploy beads-XXX         # Production deployment

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

## Documentation

- `PROTOCOL.md` - Full specification (Russian)
- `docs/TUTORIAL.md` - 15-minute quick start
- `.claude/agents/README.md` - Agent roles guide
- `README.md` - Project overview

---

## Version

**SDP v0.5.0** - Unified Workflow (AI-Comm + Beads)

Updated: 2026-01-29

---

**See Also:**
- Russian version: `PROTOCOL.md` (полная спецификация)
- Tutorial: `docs/TUTORIAL.md` (15-минутное введение)
- Agent Roles: `.claude/agents/README.md` (справочник по ролям)
