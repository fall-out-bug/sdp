# Spec-Driven Protocol v0.5.0

Workstream-driven development для AI-агентов.

---

## Навигация

```
Ты здесь?                          →  Иди сюда
─────────────────────────────────────────────────────
Нужно понять что делать            →  Phase 1: Analyze
Нужно спланировать WS              →  Phase 2: Plan
Нужно выполнить WS                 →  Phase 3: Execute
Нужно проверить результат          →  Phase 4: Review
Нужно принять архитектурное решение →  ADR Template
Нужны примеры кода hw_checker      →  HW_CHECKER_PATTERNS.md
Непонятно какие правила            →  Guardrails
─────────────────────────────────────────────────────
Multi-agent координация           →  Unified Workflow
Agent spawning/messaging          →  Agent Coordination
Telegram notifications            →  Notification System
Beads task tracking               →  Beads Integration
Feature development               →  @feature skill
```

---

## Workstream Flow

```
┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│  ANALYZE   │───→│    PLAN    │───→│  EXECUTE   │───→│   REVIEW   │
│  (Sonnet)  │    │  (Sonnet)  │    │   (Auto)   │    │  (Sonnet)  │
└────────────┘    └────────────┘    └────────────┘    └────────────┘
     │                  │                  │                  │
     ▼                  ▼                  ▼                  ▼
 Карта WS          План WS            Код            APPROVED/FIX
```

**Промпты:** `@sdp/prompts/structured/phase-{1,2,3,4}-*.md`

---

## Unified Workflow (AI-Comm + Beads)

**Начиная с v0.4.0**: SDP интегрирует AI-Comm архитектуру для multi-agent координации с Beads для task tracking.

### Компоненты Unified Workflow

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

### 1. Agent Coordination

**Agent Spawning:**
```python
from sdp.unified.agent.spawner import AgentSpawner, AgentConfig

spawner = AgentSpawner()
config = AgentConfig(
    name="builder",
    prompt="You are a build agent...",
)
agent_id = spawner.spawn_agent(config)
```

**Inter-Agent Messaging:**
```python
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
message = Message(
    sender="orchestrator",
    content="Execute WS-060-01",
    recipient=agent_id,
)
result = router.send_message(message)
```

**Role Management:**
```python
from sdp.unified.agent.role_loader import RoleLoader
from sdp.unified.agent.role_state import RoleStateManager

# Load role from .agents/{role}.md
loader = RoleLoader()
role = loader.load_role("planner")

# Activate role
state_mgr = RoleStateManager()
state_mgr.activate_role("planner")

# Check active roles
active = state_mgr.list_active()  # ["planner", "builder"]
```

### 2. Notification System

**Configuration:**
```bash
# .env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
```

**Sending Notifications:**
```python
from sdp.unified.notifications.telegram import TelegramConfig, TelegramNotifier
from sdp.unified.notifications.provider import Notification, NotificationType

# Setup
config = TelegramConfig(
    bot_token=os.getenv("TELEGRAM_BOT_TOKEN"),
    chat_id=os.getenv("TELEGRAM_CHAT_ID"),
)
notifier = TelegramNotifier(config=config)

# Send notification
notification = Notification(
    type=NotificationType.SUCCESS,
    message="Feature F24 completed successfully",
)
notifier.send(notification)
```

**Notification Types:**
- `INFO` - ℹ️ Informational messages
- `SUCCESS` - ✅ Successful operations
- `WARNING` - ⚠️ Warnings
- `ERROR` - 🚨 Errors and failures

**Mock Provider (для тестов):**
```python
from sdp.unified.notifications.mock import MockNotificationProvider

mock = MockNotificationProvider()
mock.send(notification)
assert mock.count() == 1
```

### 3. Beads Integration

**Task Tracking:**
```python
from sdp.beads import create_beads_client
from sdp.beads.models import BeadsTaskCreate, BeadsStatus

# Create client (mock for CI, real for dev)
client = create_beads_client(use_mock=True)

# Create feature task
feature = client.create_task(BeadsTaskCreate(
    title="User Authentication",
    description="Add OAuth2 login flow",
    priority=BeadsPriority.HIGH,
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

# Get ready tasks (ws2 becomes ready after ws1 completes)
ready = client.get_ready_tasks()  # [ws2.id]
```

**Checkpoint System:**
```python
from sdp.unified.orchestrator.checkpoint import CheckpointFileManager
from sdp.unified.orchestrator.agent_extension import CheckpointExtension

# Save checkpoint
checkpoint_mgr = CheckpointFileManager()
extension = CheckpointExtension(agent=orchestrator)
checkpoint_mgr.save(
    feature_id="sdp-118",
    agent_id=agent.id,
    completed_ws=["sdp-118.1", "sdp-118.2"],
    checkpoint_ext=extension,
)

# Resume from checkpoint
checkpoint = checkpoint_mgr.load("sdp-118")
if checkpoint:
    resumed = checkpoint_mgr.resume(agent, checkpoint)
```

### 4. Feature Development Flow

**Unified Entry Point (@feature skill):**
```bash
# 1. Gather requirements (interactive)
@feature "Add user authentication"
# → Deep interviewing via AskUserQuestion
# → Creates docs/intent/sdp-XXX.json
# → Creates docs/drafts/beads-sdp-XXX.md

# 2. Plan workstreams (interactive)
@design beads-sdp-XXX
# → EnterPlanMode for codebase exploration
# → Interactive planning via AskUserQuestion
# → Creates WS-XXX.01, WS-XXX.02, ...
# → Generates execution graph

# 3. Execute workstreams
@build WS-XXX.01
# → TodoWrite progress tracking
# → TDD cycle (Red → Green → Refactor)

# Or autonomous execution:
@oneshot sdp-XXX
# → Executes all WS in dependency order
# → Background execution support
# → Checkpoint save/restore

# 4. Quality review
@review sdp-XXX
# → Validates all quality gates
# → Returns APPROVED/CHANGES_REQUESTED

# 5. Deploy
@deploy sdp-XXX
# → Generates deployment configs
# → Creates PR with changelog
```

### 5. Quality Gates (Unified)

**Все прежние gates + новые:**

```bash
# Agent tests (309+ tests)
pytest tests/unified/ -v

# Beads integration
pytest tests/unified/test_e2e/test_beads_client.py

# Telegram E2E (requires credentials)
export TELEGRAM_BOT_TOKEN="..."
export TELEGRAM_CHAT_ID="..."
pytest tests/unified/test_e2e/test_telegram_e2e.py::TestRealTelegramIntegration
```

### 6. Examples

**Multi-Agent Feature Execution:**
```python
# 1. Orchestrator spawns specialized agents
spawner = AgentSpawner()
planner_id = spawner.spawn_agent(AgentConfig(name="planner", ...))
builder_id = spawner.spawn_agent(AgentConfig(name="builder", ...))

# 2. Send messages
router.send_message(Message(
    sender="orchestrator",
    content="Plan feature F24",
    recipient=planner_id,
))

# 3. Receive notifications
notifier.send(Notification(
    type=NotificationType.INFO,
    message="Planner completed: 5 workstreams created",
))

# 4. Track in Beads
client = create_beads_client(use_mock=True)
feature = client.create_task(BeadsTaskCreate(title="F24", ...))
# ... decompose into WS, execute, etc.
```

**Bug Report Workflow:**
```python
from sdp.unified.agent.bug_report import BugReportFlow, BugSeverity

# Create bug report
bug_flow = BugReportFlow()
bug = bug_flow.create_report(
    title="Login fails on Firefox",
    description="OAuth2 token not stored",
    severity=BugSeverity.P1,
    workstream_id="WS-060-01",
)

# Check blocking
if "WS-060-01" in bug_flow.get_blocking_workstreams():
    notifier.send(Notification(
        type=NotificationType.ERROR,
        message="WS-060-01 blocked by P1 bug",
    ))

# Mark resolved
bug_flow.update_status(bug.id, BugStatus.RESOLVED)
```

**Дополнительная документация:**
- `src/sdp/unified/agent/README.md` - Agent system details
- `src/sdp/unified/notifications/README.md` - Notification system
- `src/sdp/beads/README.md` - Beads integration
- `docs/drafts/beads-sdp-118.md` - Unified workflow implementation

---

## Терминология

| Термин | Scope | Размер | Пример |
|--------|-------|--------|--------|
| **Release** | Продуктовая веха | 10-30 Features | R1: Submissions E2E |
| **Feature** | Крупная фича | 5-30 Workstreams | F24: Obsidian Vault |
| **Workstream** | Атомарная задача | SMALL/MEDIUM/LARGE | WS-140: Vault Domain |

**Scope метрики для Workstream:**
- **SMALL**: < 500 LOC, < 1500 tokens
- **MEDIUM**: 500-1500 LOC, 1500-5000 tokens  
- **LARGE**: > 1500 LOC → разбить на 2+ WS

### ⚠️ Важно: NO TIME-BASED ESTIMATES

**ЗАПРЕЩЕНО использовать время для оценки:**
- ❌ "Это займёт 2 часа"
- ❌ "Нужно 3 дня"
- ❌ "Не успеваю за неделю"
- ❌ "Времени нет"
- ❌ "Это долго"

**ИСПОЛЬЗУЙ scope метрики:**
- ✅ "Это MEDIUM workstream (1000 LOC, 3000 tokens)"
- ✅ "Scope превышен, нужно разбить на 2 WS"
- ✅ "По scope это SMALL задача"

#### ✅ Разрешённые упоминания времени (исключения)

Время **разрешено** только в следующих случаях (и **не является оценкой scope**):

- **Telemetry / измерения**: elapsed time, timestamps в логах, метрики выполнения (например, `"elapsed": "1h 23m"`).
- **SLA / операционные цели**: hotfix/bugfix target windows (например, “P0 hotfix: <2h”, “P1/P2 bugfix: <24h”).
- **Human Verification (UAT)**: ориентиры для человека (“Smoke test: 30 sec”, “Scenarios: 5–10 min”).

Во всех остальных контекстах **время запрещено** — используем только LOC/tokens и sizing (SMALL/MEDIUM/LARGE).

**Почему НЕ время:**
1. AI agents работают с разной скоростью (Sonnet ≠ Haiku ≠ GPT)
2. Scope объективен (LOC, tokens), время субъективно
3. Время создаёт ложное давление ("не успеваю" → спешка → баги)
4. One-shot execution: агент выполняет WS за один проход, независимо от "времени"

### Иерархия (Product)

```
PORTAL_VISION.md (продукт)
    ↓
RELEASE_PLAN.md (релизы)
    ↓
Feature (F01-F99) — крупные фичи
    ↓
Workstream (WS-001-WS-999) — атомарные задачи
```

### Устаревшие термины

- ~~Epic (EP)~~ → **Feature (F)** (с 2026-01-07)
- ~~Sprint~~ → не используется

---

## Workstream Naming Convention (PP-FFF-SS)

### Format

```
PP-FFF-SS
├─ PP: Project ID (2 digits, 00-99)
├─ FFF: Feature ID (3 digits, 000-999)
└─ SS: Workstream Sequence (2 digits, 00-99)
```

### Project ID Registry

| ID | Project | Description |
|----|---------|-------------|
| 00 | **SDP Protocol** | Universal meta-protocol (uses itself) |
| 01 | *Reserved* | Available for future use |
| 02 | hw_checker | Homework validation system |
| 03 | mlsd | ML System Design course |
| 04 | bdde | Big Data course |
| 05 | msu_ai_masters | Meta-repo configuration |

**Principle:** PP = who owns the workstream. All projects (02-05) use SDP (00) as their tool.

### Examples

| WS ID | Project | Feature | Description |
|-------|---------|---------|-------------|
| 00-500-01 | SDP | F500 | Sync SDP content |
| 00-410-01 | SDP | F410 | Contract-driven WS spec |
| 02-150-01 | hw_checker | F150 | Config fixes |
| 02-201-01 | hw_checker | F201 | Multi-IDE parity |
| 03-100-01 | mlsd | F100 | Question domain |
| 04-050-01 | bdde | F050 | Data pipeline |

### Cross-Project Dependencies

Projects can depend on SDP workstreams:

```yaml
# In hw_checker (02-150-03.md):
---
depends_on:
  - 00-100-05  # SDP Protocol WS-100-05
---
```

**Rule:** Projects (02-05) may depend on SDP (00), but SDP does not depend on specific projects.

### Migration from Legacy Format

| Old Format | New Format | Example |
|------------|------------|---------|
| `WS-FFF-SS` | `PP-FFF-SS` | WS-193-01 → 00-193-01 |
| `WS-FFF-SS` | `PP-FFF-SS` | WS-150-01 → 02-150-01 |

The SDP parser supports both formats for backward compatibility. Legacy `WS-FFF-SS` format is automatically interpreted as Project 00 (SDP).

### Automated Migration

See `sdp/docs/migration/ws-naming-migration.md` for detailed migration guide.

---

## Guardrails

### AI-Readiness (БЛОКИРУЮЩИЕ)

| Правило | Порог | Проверка |
|---------|-------|----------|
| File size | < 200 LOC | `wc -l` |
| Complexity | CC < 10 | `ruff --select=C901` |
| Type hints | 100% public | Visual |
| Nesting | ≤ 3 levels | Visual |

### Clean Architecture (БЛОКИРУЮЩИЕ)

```
Domain      →  НЕ импортирует ничего из других слоёв
Application →  НЕ импортирует infrastructure напрямую
```

```bash
# Проверка
grep -r "from hw_checker.infrastructure" hw_checker/domain/ hw_checker/application/
# Должно быть пусто
```

### Error Handling (БЛОКИРУЮЩИЕ)

```python
# ЗАПРЕЩЕНО
except:
    pass

except Exception:
    return None

# ОБЯЗАТЕЛЬНО
except SpecificError as e:
    log.error("operation.failed", error=str(e), exc_info=True)
    raise
```

### Security (для DinD)

- [ ] Нет `privileged: true`
- [ ] Нет `/var/run/docker.sock` mounts
- [ ] Resource limits заданы
- [ ] Нет string interpolation в shell commands

---

## Quality Gates

### Gate 1: Analyze → Plan
- [ ] Карта WS сформирована
- [ ] Зависимости указаны
- [ ] AI-Readiness оценён для каждого WS

### Gate 2: Plan → Execute
- [ ] **WS не существует** в INDEX (проверено)
- [ ] **Scope оценён**, не превышает MEDIUM
- [ ] Все пути файлов указаны
- [ ] Код готов к copy-paste
- [ ] Критерии завершения включают: tests + coverage + regression
- [ ] Ограничения явные
- [ ] **НЕТ временных оценок** (часов/дней)

### Gate 3: Execute → Review
- [ ] Все шаги выполнены
- [ ] Критерии завершения пройдены
- [ ] **Coverage ≥ 80%** для изменённых файлов
- [ ] **Regression passed** (fast tests)
- [ ] **Нет TODO/Later** в коде
- [ ] Отчёт сформирован

### Gate 4: Review → Done
- [ ] AI-Readiness: ✅
- [ ] Clean Architecture: ✅
- [ ] Error Handling: ✅
- [ ] Tests & Coverage: ✅ (≥80%)
- [ ] Regression: ✅ (all fast tests)
- [ ] Review записан **в конец WS файла** (не отдельный файл)

### Gate 5: Done → Deploy (Human UAT)

**UAT (User Acceptance Testing)** — проверка человеком перед деплоем:

| Шаг | Описание | Время |
|-----|----------|-------|
| 1 | Quick Smoke Test | 30 сек |
| 2 | Detailed Scenarios (happy path + errors) | 5-10 мин |
| 3 | Red Flags Check | 2 мин |
| 4 | Sign-off | 1 мин |

**UAT Guide создаётся автоматически** после `/codereview APPROVED`:
- Feature-level: `docs/uat/F{XX}-uat-guide.md`
- WS-level: секция "Human Verification (UAT)" в WS файле

**Без Sign-off человека → Deploy блокирован.**

---

## WS Scope Control

**Метрики размера (вместо времени):**

| Размер | Строк кода | Токенов | Действие |
|--------|-----------|---------|----------|
| **SMALL** | < 500 | < 1500 | ✅ Оптимально |
| **MEDIUM** | 500-1500 | 1500-5000 | ✅ Допустимо |
| **LARGE** | > 1500 | > 5000 | ❌ **РАЗБИТЬ** |

**Правило:** Все WS должны быть SMALL или MEDIUM.

**Если scope превышен во время Execute:**
→ STOP, вернуться к Phase 2 для разбиения на WS-XXX-1, WS-XXX-2

---

## Test Coverage Gate

**Минимум:** 80% для изменённых/созданных файлов

```bash
pytest tests/unit/test_module.py -v \
  --cov=hw_checker/module \
  --cov-report=term-missing \
  --cov-fail-under=80
```

**Если coverage < 80% → CHANGES REQUESTED (HIGH)**

---

## Regression Gate

**После каждого WS:**

```bash
# Все fast tests ДОЛЖНЫ проходить
pytest tests/unit/ -m fast -v
```

**Если регресс нарушен → CHANGES REQUESTED (CRITICAL)**

---

## TODO/Later Gate

**СТРОГО ЗАПРЕЩЕНО в коде:**
- `# TODO: ...`
- `# FIXME: ...`
- Комментарии "оставлю на потом", "временное решение"

**Исключение:** `# NOTE:` — только для пояснений

**Если обнаружено → CHANGES REQUESTED (HIGH)**

---

## ⛔ NO TECH DEBT

**Концепция Tech Debt ЗАПРЕЩЕНА в проекте.**

❌ "Это tech debt, сделаем потом"
❌ "Временное решение, вернёмся позже"
❌ "Грязный код, но работает"
❌ "Отложим рефакторинг"

✅ **Правило: всё говно убираем сразу.**

**Если код не соответствует стандартам:**
1. Исправь в текущем WS
2. Если scope превышен → разбей на WS (см. ниже)
3. НЕ оставляй "на потом"

**Философия:** Каждый WS оставляет код в идеальном состоянии. Нет накапливающегося долга.

---

## 🔀 Substreams: Правила разбиения

**Если WS нужно разбить на части:**

### Формат нумерации (СТРОГО)

```
WS-{PARENT_ID}-{SEQ}

Где:
- PARENT_ID = ID родительского WS (3 цифры, с ведущими нулями)
- SEQ = порядковый номер substream (2 цифры: 01, 02, ... 99)
```

**Примеры:**
```
WS-050         ← родительский (разбивается)
├── WS-050-01  ← первый substream
├── WS-050-02  ← второй substream
├── WS-050-03  ← третий substream
├── ...
├── WS-050-10  ← десятый (сортировка корректна!)
└── WS-050-15  ← пятнадцатый
```

**ЗАПРЕЩЁННЫЕ форматы:**
```
❌ WS-050-A, WS-050-B      (буквы)
❌ WS-050-part1            (слова)
❌ WS-050.1, WS-050.2      (точки)
❌ WS-50-1                 (без ведущих нулей в PARENT)
❌ WS-050-1                (однозначный SEQ — всегда 01, 02...)
```

### ОБЯЗАТЕЛЬНО при разбиении:

1. **Создай ВСЕ файлы substreams** в `workstreams/backlog/`:
   ```
   WS-050-01-domain-entities.md
   WS-050-02-application-layer.md
   WS-050-03-infrastructure.md
   ```

2. **Заполни каждый substream** полностью (не stub):
   - Контекст
   - Зависимости (WS-XXX-1 → WS-XXX-2 → ...)
   - Входные файлы
   - Шаги
   - Код
   - Критерии завершения

3. **Обнови INDEX.md** с новыми WS

4. **Удали или пометь родительский WS** как "Разбит → WS-XXX-1, WS-XXX-2"

### ЗАПРЕЩЕНО:

❌ Ссылаться на несуществующие WS ("см. WS-050-02" без создания файла)
❌ Оставлять пустые stubs ("TODO: заполнить")
❌ Разбивать без создания файлов
❌ Partial execution ("сделал часть, остальное в другом WS")
❌ Форматы: `24.1`, `WS-24-1`, `WS-050-1`, `WS-050-part1`
❌ Time estimates: "0.5 дня", "3 дня" — только LOC/tokens
❌ Создавать отдельные `-ANALYSIS.md` файлы (анализ → сразу в WS файлы)

### Пример правильного разбиения:

```markdown
## WS-050: Large Feature → РАЗБИТ

**Статус:** Разбит на substreams
**Причина:** Scope > MEDIUM (2500 LOC)

**Substreams:** (формат: WS-{PARENT}-{SEQ}, SEQ всегда 2 цифры)
| ID | Файл | Scope |
|----|------|-------|
| WS-050-01 | WS-050-01-domain-entities.md | SMALL (400 LOC) |
| WS-050-02 | WS-050-02-application-layer.md | MEDIUM (800 LOC) |
| WS-050-03 | WS-050-03-infrastructure.md | MEDIUM (700 LOC) |
| WS-050-04 | WS-050-04-presentation.md | SMALL (300 LOC) |
| WS-050-05 | WS-050-05-integration-tests.md | SMALL (300 LOC) |

Все файлы созданы в backlog/, добавлены в INDEX.md.
```

### Проверка перед ссылкой на substream

```bash
# ОБЯЗАТЕЛЬНО перед тем как написать "см. WS-050-02":
ls tools/hw_checker/docs/workstreams/backlog/WS-050-02-*.md

# Если "No such file" → СНАЧАЛА создай файл!

# Проверка формата нумерации (должны быть 2 цифры для SEQ):
ls tools/hw_checker/docs/workstreams/backlog/ | grep -E "WS-[0-9]{3}-[0-9]{2}-"
# ✅ WS-050-01-domain.md, WS-050-02-app.md
# ❌ WS-050-1-domain.md (SEQ должен быть 01, не 1)

# Проверка на time estimates (должно быть пусто):
grep -rE "дн[яей]|час[ов]|недел" tools/hw_checker/docs/workstreams/backlog/WS-050*.md
```

---

## ADR Template

Когда принимаешь архитектурное решение, создай:

`docs/architecture/adr/YYYY-MM-DD-{title}.md`

```markdown
# ADR: {Title}

## Status
Proposed / Accepted / Deprecated

## Context
[Какая проблема? Какие ограничения?]

## Decision
[Что решили делать?]

## Alternatives Considered
1. [Альтернатива 1] — почему нет
2. [Альтернатива 2] — почему нет

## Consequences
- [+] Плюс
- [-] Минус
- [!] Риск
```

---

## Workstream Format

```markdown
## WS-{ID}: {Title}

### Контекст
[Почему нужно]

### Зависимость  
[WS-XX / Независимый]

### Входные файлы
- `path/to/file.py` — что там

### Шаги
1. [Атомарное действие]
2. ...

### Код
```python
# Готовый код
```

### Ожидаемый результат
- [Что должно быть]

### Критерий завершения
```bash
pytest ...
ruff check ...
```

### Ограничения
- НЕ делать: ...
```

---

## Иерархия документации (C4-подобная)

```
L1: System      docs/SYSTEM_OVERVIEW.md
    ↓ Общий контекст системы, границы, основные домены
    
L2: Domain      docs/domains/{domain}/DOMAIN_MAP.md  
    ↓ Структура домена, компоненты, интеграции
    
L3: Component   docs/domains/{domain}/components/{comp}/SPEC.md
    ↓ Детальная спецификация компонента
    
L4: Workstream  docs/workstreams/WS-XXX.md
    ↓ Конкретная задача для выполнения
```

### Navigation Flow

**Phase 1 (Analyze):**
1. Читай L1 (`SYSTEM_OVERVIEW.md`) для общего контекста
2. Выбери релевантный домен, читай L2 (`domains/{domain}/DOMAIN_MAP.md`)
3. Если затрагиваешь компонент, читай L3 (component SPEC)
4. Генерируй L4 (workstream map)

**Phase 2 (Plan):**
1. Читай L4 (`workstreams/INDEX.md`) — проверь дубликаты
2. Читай L1/L2/L3 для контекста конкретного WS
3. Создай детальный план WS

**Phase 3 (Execute):**
1. Работай по плану WS (L4)

**Phase 4 (Review):**
1. Проверь качество кода
2. Если WS изменил domain boundaries → обновить L2
3. Если WS изменил component → обновить L3

### Product vs Architecture Hierarchy

**Product (планирование фичей):**
```
PORTAL_VISION.md → RELEASE_PLAN.md → Feature (F) → Workstream (WS)
```

**Architecture (структура кода/документации):**
```
L1 (System) → L2 (Domain) → L3 (Component) → L4 (Workstream)
```

**Пересечение:**
- Feature F24 → создаёт/модифицирует L2 (content domain)
- Workstream WS-140 → создаёт L3 (vault component)

---

## Quick Reference

```bash
# AI-Readiness check
find hw_checker -name "*.py" -exec wc -l {} + | awk '$1 > 200'
ruff check hw_checker --select=C901

# Clean Architecture check  
grep -r "from hw_checker.infrastructure" hw_checker/domain/ hw_checker/application/

# Error handling check
grep -rn "except:" hw_checker/
grep -rn "except Exception" hw_checker/ | grep -v "exc_info"

# Test coverage (≥80%)
pytest tests/unit/test_module.py -v \
  --cov=hw_checker/module \
  --cov-report=term-missing \
  --cov-fail-under=80

# Regression (fast tests)
pytest tests/unit/ -m fast -v

# TODO/Later check
grep -rn "TODO\|FIXME" hw_checker/ --include="*.py" | grep -v "# NOTE"

# Full test suite
pytest -m fast -x --tb=short
pytest --cov=hw_checker --cov-report=term-missing
```

---

## Observability

### Telegram Notifications

Automated notifications for critical events:

```bash
# Setup
export TELEGRAM_BOT_TOKEN="..."
export TELEGRAM_CHAT_ID="..."

# Events: oneshot_started, oneshot_completed, oneshot_blocked,
#         ws_failed, review_failed, breaking_changes, e2e_failed,
#         deploy_success, hotfix_deployed
```

See: `sdp/notifications/TELEGRAM.md`

### Audit Log

Centralized logging of all workflow events:

```bash
# Configuration
export AUDIT_LOG_FILE="/var/log/consensus-audit.log"

# Format: ISO8601|EVENT_TYPE|USER|GIT_BRANCH|EVENT_DATA
# Example:
# 2026-01-11T00:30:15+03:00|WS_START|user|feature/lms|ws=WS-060-01

# Query
grep "feature=F60" /var/log/consensus-audit.log
grep "WS_FAILED" /var/log/consensus-audit.log
```

See: `sdp/notifications/AUDIT_LOG.md`

### Breaking Changes Detection

Automatic detection and documentation:

```bash
# Runs in pre-commit hook
python scripts/detect_breaking_changes.py --staged

# Generates:
# - BREAKING_CHANGES.md
# - MIGRATION_GUIDE.md (template)
```

See: `tools/hw_checker/scripts/detect_breaking_changes.py`

---

