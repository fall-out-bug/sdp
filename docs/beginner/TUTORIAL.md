# SDP: 15-Minute Tutorial

**Spec-Driven Protocol** - workstream-driven development для AI-агентов с multi-agent координацией.

---

## What is SDP? (2 min)

SDP = **Workstream-driven development** + **AI-Comm** + **Beads**

```
Idea → @feature → Workstreams → @build → Review → Deploy
  ↓        ↓           ↓          ↓       ↓        ↓
Beads  Agents    Task tracking   TDD    Quality  Production
```

**Key Concepts:**

| Concept | Description |
|---------|-------------|
| **Workstream (WS)** | Атомарная задача (500-1500 LOC) |
| **Feature** | Крупная фича (5-30 workstreams) |
| **Agent** | AI-агент для выполнения задач |
| **Beads** | Task tracking система |
| **@feature** | Unified entry point для фич |
| **@build** | Execute workstream with TDD |

---

## Quick Start (5 min)

### Installation

```bash
# Clone repository
git clone <repo>
cd sdp

# Install dependencies
pip install -e .

# Or with uv
uv pip install -e .
```

### Your First Feature

```bash
# 1. Create feature (interactive)
@feature "Add user authentication"

# Claude will ask:
# - Technical approach (JWT vs sessions?)
# - UI/UX requirements
# - Database schema
# - Testing strategy

# → Creates: docs/intent/sdp-XXX.json
# → Creates: docs/drafts/beads-sdp-XXX.md
```

```bash
# 2. Plan workstreams
@design beads-sdp-XXX

# Claude explores codebase and creates workstreams:
# - WS-XXX.01: Domain model (450 LOC)
# - WS-XXX.02: Database schema (300 LOC)
# - WS-XXX.03: Repository layer (500 LOC)
# - WS-XXX.04: Service layer (600 LOC)
# - WS-XXX.05: API endpoints (400 LOC)

# → Creates: docs/workstreams/beads-sdp-XXX.md
```

```bash
# 3. Execute workstream
@build WS-XXX.01

# Claude follows TDD:
# 1. Write failing test (Red)
# 2. Implement minimum code (Green)
# 3. Refactor (Refactor)

# → Shows real-time progress with TodoWrite
# → Runs tests, mypy, ruff
# → Commits when complete
```

```bash
# 4. Or execute all workstreams autonomously
@oneshot sdp-XXX

# Executes all WS in dependency order:
# WS-XXX.01 → WS-XXX.02 → WS-XXX.03 → WS-XXX.04 → WS-XXX.05

# → Checkpoint save/restore
# → Background execution support
# → Progress notifications via Telegram
```

---

## Core Concepts (5 min)

### 1. Agent Coordination

SDP spawns specialized agents for different tasks:

```python
from sdp.unified.agent.spawner import AgentSpawner, AgentConfig

# Spawn agents
spawner = AgentSpawner()
planner = spawner.spawn_agent(AgentConfig(
    name="planner",
    prompt="You break features into workstreams...",
))

builder = spawner.spawn_agent(AgentConfig(
    name="builder",
    prompt="You execute workstreams with TDD...",
))

# Send messages between agents
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="orchestrator",
    content="Plan feature F24",
    recipient=planner,
))
```

**Agent Roles:**
- `planner` - Breaks features into workstreams
- `builder` - Executes workstreams with TDD
- `reviewer` - Quality checks
- `deployer` - Production deployment

### 2. Beads Integration

Track tasks and dependencies with Beads:

```python
from sdp.beads import create_beads_client
from sdp.beads.models import BeadsTaskCreate, BeadsStatus

# Create client
client = create_beads_client(use_mock=True)

# Create feature task
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

# Add dependency (ws2 blocked by ws1)
client.add_dependency(ws2.id, ws1.id, dep_type="blocks")

# Update status
client.update_task_status(ws1.id, BeadsStatus.CLOSED)

# Get ready tasks (ws2 becomes ready after ws1)
ready = client.get_ready_tasks()  # [ws2.id]
```

### 3. Telegram Notifications

Get notified about critical events:

```bash
# Setup .env
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
    message="Feature F24 completed successfully",
))
```

**Notification Types:**
- ℹ️ `INFO` - Informational messages
- ✅ `SUCCESS` - Successful operations
- ⚠️ `WARNING` - Warnings
- 🚨 `ERROR` - Errors and failures

### 4. Quality Gates

Every workstream must pass:

```bash
# Test coverage ≥80%
pytest tests/unit/test_module.py \
  --cov=src/sdp/module \
  --cov-fail-under=80

# Type checking
mypy src/sdp/module/ --strict

# Linting
ruff check src/sdp/module/

# All files <200 LOC
find src/sdp/module -name "*.py" -exec wc -l {} +
```

**Quality Checklist:**
- ✅ Tests first (TDD)
- ✅ ≥80% coverage
- ✅ mypy --strict
- ✅ ruff clean
- ✅ Files <200 LOC
- ✅ No `except: pass`
- ✅ Type hints everywhere

---

## Common Workflows (3 min)

### Create Feature from Idea

```bash
# 1. Start with @feature (deep interviewing)
@feature "Add password reset"

# Claude asks:
# - Email service (SendGrid vs AWS SES)?
# - Token storage (database vs Redis)?
# - Token expiry (1 hour vs 24 hours)?
# - Rate limiting (yes/no)?

# → Creates comprehensive spec
```

### Execute Single Workstream

```bash
# 2. Plan and execute
@design beads-password-reset
@build WS-XXX.01

# Claude shows progress:
#   [in_progress] Write failing test (Red)
#   [pending] Implement minimum code (Green)
#   [pending] Refactor implementation
#   ...
```

### Autonomous Feature Execution

```bash
# 3. Let Claude execute everything
@oneshot beads-password-reset

# Spawns orchestrator agent:
# - Executes all WS in order
# - Saves checkpoints after each WS
# - Sends Telegram notifications
# - Resumes from interruption
```

### Review and Deploy

```bash
# 4. Quality review
@review beads-password-reset

# Validates:
# - All quality gates passed
# - Tests ≥80% coverage
# - No tech debt
# - Returns: APPROVED / CHANGES_REQUESTED

# 5. Deploy to production
@deploy beads-password-reset

# - Generates deployment configs
# - Creates PR with changelog
# - Runs smoke tests
```

---

## Next Steps

**Learn More:**
- `PROTOCOL.md` - Full specification (English)
- `PROTOCOL_RU.md` - Полная спецификация (Русский)
- `README.md` - Project overview
- `docs/drafts/beads-sdp-118.md` - Unified workflow implementation

**Try It:**
```bash
# Create your first feature
@feature "Add user comments"

# Or explore existing features
@design beads-sdp-118
```

**Get Help:**
- `@feature` - Interactive requirements gathering
- `@design` - Interactive workstream planning
- `/debug <issue>` - Systematic debugging

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
