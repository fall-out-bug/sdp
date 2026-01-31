# Hybrid SDP: Unified Workflow Design

> **Status:** Design Complete
> **Date:** 2026-01-28
> **Goal:** Unify @idea/@design/@oneshot with team coordination into seamless hybrid intelligence workflow

---

## Table of Contents

1. [High-Level Architecture](#high-level-architecture)
2. [Component Architecture](#component-architecture)
3. [Data Flow & Execution Lifecycle](#data-flow--execution-lifecycle)
4. [Error Handling & Edge Cases](#error-handling--edge-cases)
5. [Testing Strategy](#testing-strategy)
6. [Success Metrics](#success-metrics)

---

## High-Level Architecture

### Unified SDP Workflow: Hybrid Intelligence Mode

**Единый pipeline** объединяющий @idea/@design/@oneshot с team coordination:

```
1. @feature "Feature X"
   ↓
   Orchestrator presents progressive menu
   ↓
2. @idea phase (optional: --skip-requirements)
   ↓ Analyst agent (active)
   ↓ Interactive interviewing (AskUserQuestion)
   ↓ Requirements spec + Intent JSON
   ↓ RequirementsGate: human approval
   ↓
3. @design phase (optional: --skip-architecture)
   ↓ Architect agent (active)
   ↓ Workstream decomposition (EnterPlanMode)
   ↓ Beads tasks: parent feature + N workstreams
   ↓ ArchitectureGate: human approval
   ↓
4. @oneshot phase
   ↓ Orchestrator agent создаёт Team "feature-X"
   ↓ Team spawns: developer × 3, qa × 2, devops
   ↓ Additional 95+ roles dormant
   ↓
5. Parallel execution (Beads get_ready_tasks)
   ↓ Developer agents выполняют workstreams
   ↓ SendMessage coordination между team members
   ↓ Agent concerns → Beads bug reports
   ↓
6. @review phase
   ↓ QA agents run full test suite
   ↓ Code quality verification
   ↓
7. Final UAT (optional: --skip-uat)
   ↓ Human tests on stage
   ↓
8. @deploy production
```

### Ключевые принципы

**1. Все роли всегда в team**
- Team создаётся один раз при `@oneshot`
- 100+ ролей могут быть зарегистрированы
- Базовые 5 ролей активны сразу
- Дополнительные "спят" пока не получат SendMessage

**2. Role activation via SendMessage**
```
Architect → SendMessage("security", "review auth design")
  ↓
Security agent "просыпается" (dormant → active)
  ↓
Выполняет review
  ↓
SendMessage("architect", "review complete, no concerns")
  ↓
Security снова "спит" (active → dormant)
```

**3. Approval gates с skip-флагами**
```bash
# Полный режим (все gates)
@feature "Add auth"

# Pure vibe mode (без gates)
@feature "Add auth" --skip-requirements --skip-architecture --skip-uat

# Hybrid (только architecture gate)
@feature "Add auth" --skip-requirements
```

**4. Beads как единый state store**
- Tasks: parent feature + child workstreams
- Dependencies: автоматическое unblocking
- Bug reports: от всех агентов (architect, security, sre, etc.)
- Progress tracking: готовые WS, blocked, pending

### Масштабируемость ролей

**Базовые роли** (5 штук, всегда активны):
- Analyst, Architect, Developer, QA, DevOps

**Плагинные роли** (95+ штук, dormant):
- Security, SRE, Data/ML, Documentation
- Lecture, Seminar, Prompt Engineer
- Performance, Accessibility, Privacy
- Custom: reviewer, codeowner, integrator, ...

**Добавление новой роли:**
```yaml
# .claude/roles/custom_reviewer.yml
name: reviewer
description: "Code reviewer for critical paths"
prompt_path: prompts/custom_reviewer.md
active_by_default: false
wakeup_messages: ["review_required"]
```

---

## Component Architecture

### Key Components

**1. Orchestrator Agent** (`@oneshot` core)
- **Ответственность**: Создаёт team, координирует выполнение workstreams
- **Input**: Feature ID from Beads (parent task)
- **Output**: Execution report + checkpoint JSON
- **Ключевые методы**:
  - `create_team(feature_id)` — создаёт Team с 100+ ролей
  - `dispatch_workstreams()` — распределяет WS по developer agents
  - `monitor_progress()` — отслеживает completion через Beads
  - `handle_bugs()` — агрегирует bug reports от всех agents
  - `create_checkpoint()` — fault tolerance checkpoints

**2. Team Manager** (Team tool wrapper)
- **Ответственность**: Жизненный цикл team, роль активации
- **State**: `~/.claude/teams/{feature_id}/config.json`
- **Ключевые методы**:
  - `register_role(role_config)` — регистрирует роль (все 100+ при init)
  - `activate_role(role_name)` — переводит роль из dormant → active
  - `send_message(from, to, content)` — меж-агентская коммуникация
  - `get_dormant_roles()` — возвращает спящие роли

**3. Beads Integration Layer**
- **Ответственность**: Dependency tracking, bug reports, state management
- **Key components**:
  - `BeadsClient` — CLI wrapper
  - `BeadsTask` model — task/dependency model
  - `BugReporter` — создаёт bug tasks от agents
  - `DependencyResolver` — topological sort для execution order

**4. Agent Runtime** (Task tool + SendMessage)
- **Ответственность**: Execution agents, меж-агентская коммуникация
- **Key components**:
  - `AgentSpawner` — создаёт agents через Task tool
  - `MessageRouter` — маршрутизирует SendMessage
  - `RoleLoader` — загружает role prompts из `prompts/{role}.md`
  - `StateTracker` — отслеживает active/dormant

**5. Approval Gate Manager**
- **Ответственность**: Decision gates с skip-флагами
- **State**: Транзакционное (git commit per decision)
- **Key components**:
  - `RequirementsGate` — approves/rejects @idea output
  - `ArchitectureGate` — approves/rejects @design output
  - `UATGate` — final UAT before deploy
  - `SkipFlagParser` — парсит `--skip-*` флаги

**6. Notification Router** (multi-provider)
- **Ответственность**: Маршрутизирует уведомления в разные каналы
- **Providers**:
  - `TelegramNotifier` (базовый, реализован)
  - `SlackNotifier` (future)
  - `DiscordNotifier` (future)
  - `EmailNotifier` (future)
- **Key methods**:
  - `notify(message)` — отправляет во все enabled providers
  - `add_provider(provider)` — регистрирует новый provider

### Component Interaction Diagram

```
┌──────────────────┐
│  @feature skill  │
│  (User invokes)  │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐      ┌──────────────────┐
│ Orchestrator     │─────→│ Team Manager     │
│ Agent            │      │ (100+ roles)     │
└────────┬─────────┘      └────────┬─────────┘
         │                         │
         │ Beads dependencies      │ SendMessage
         ▼                         ▼
┌──────────────────┐      ┌──────────────────┐
│ Beads Integration│      │ Agent Runtime    │
│ (bugs, state)    │      │ (execution)      │
└────────┬─────────┘      └────────┬─────────┘
         │                         │
         │ Bug reports             │ Agent status
         ▼                         │
┌──────────────────┐              │
│ Notification     │◄─────────────┘
│ Router           │
│ (Telegram/Slack/)│
└──────────────────┘
```

---

## Data Flow & Execution Lifecycle

### Entry Point: @feature (Progressive Disclosure)

```
User: @feature "Add user authentication"
  ↓
┌─────────────────────────────────────┐
│ @feature skill orchestrator         │
│ ├─ Progressive disclosure           │
│ └─ Step-by-step menu               │
└─────────────────────────────────────┘
  ↓
Interactive menu:
"Let's build 'Add user authentication' together!

Step 1 of 5: Requirements Gathering
  Run @idea to clarify requirements?
  [1] Yes - Run @idea with interactive interview
  [2] Skip - I have specs already
  [3] Custom - Skip both @idea and @design
"
```

### Phase 1: Requirements Gathering (@feature → @idea)

```
@feature invokes @idea skill
  ↓
┌─────────────────────────────────────┐
│ Analyst agent (active in team)      │
│ ├─ Load prompt: prompts/analyst.md  │
│ ├─ Interview user via AskUserQuestion│
│ ├─ Explore tradeoffs (multi-choice) │
│ └─ Deep dive on concerns            │
└─────────────────────────────────────┘
  ↓
Output artifacts:
├─ docs/drafts/idea-user-auth.md
│  ├─ Problem statement
│  ├─ User stories
│  ├─ Technical approach
│  └─ Acceptance criteria
│
└─ docs/intent/f0xx-user-auth.json
   ├─ mission: "Enable secure auth..."
   ├─ alignment: [PRODUCT_VISION.md refs]
   └─ interview_answers: {...}
  ↓
@feature: Requirements approval
├─ Present artifacts to human
├─ Ask: "Requirements look good? Proceed to @design?"
└─ Human: [1] Yes
```

### Phase 2: Architecture Design (@feature → @design)

```
@feature invokes @design skill
  ↓
┌─────────────────────────────────────┐
│ Architect agent (active in team)    │
│ ├─ Read: docs/drafts/idea-*.md      │
│ ├─ Explore codebase (EnterPlanMode) │
│ ├─ Decompose into workstreams       │
│ └─ Create Beads tasks               │
└─────────────────────────────────────┘
  ↓
Output artifacts:
├─ docs/workstreams/backlog/
│  ├─ WS-0xx-01: Auth domain entities
│  ├─ WS-0xx-02: Auth repository
│  └─ ...
│
└─ Beads tasks:
   ├─ bd-f0x: parent feature (status: open)
   ├─ bd-f0x-01: child WS (depends: [])
   └─ bd-f0x-02: child WS (depends: [bd-f0x-01])
  ↓
@feature: Architecture approval
├─ Present workstreams to human
├─ Ask: "Workstreams look good? Ready to execute?"
└─ User: [1] Yes
```

### Phase 3: Team Creation & Execution (@feature → @oneshot)

```
@feature invokes @oneshot skill
  ↓
@feature: "Ready to execute! Choose mode:

  [1] Standard mode
      - Requirements approval ✓
      - Architecture approval ✓
      - Final UAT required

  [2] Pure vibe mode
      - Skip all approvals
      - Agents work autonomously
      - No UAT required

  [3] Custom mode
      - Mix & match approvals

Your choice: "
  ↓
User selects: [1] Standard mode
  ↓
Orchestrator agent:
├─ Load feature spec from Beads
├─ Generate unique team ID
├─ Create checkpoint
└─ Create team with 100+ roles
```

### Team Creation Flow

```
TeamManager.create_team("feature-f0x")
  ↓
┌─────────────────────────────────────┐
│ Team config:                        │
│ ~/.claude/teams/feature-f0x.json   │
│ {                                   │
│   "team_name": "feature-f0x",       │
│   "created_at": "2026-01-28...",    │
│   "roles": {                        │
│     "analyst": {"status": "idle"},  │
│     "architect": {"status": "idle"},│
│     "developer_1": {"status": "active"},│
│     "developer_2": {"status": "active"},│
│     "developer_3": {"status": "active"},│
│     "qa_1": {"status": "active"},   │
│     "qa_2": {"status": "active"},   │
│     "devops": {"status": "idle"},   │
│     "security": {"status": "dormant"},│
│     "sre": {"status": "dormant"},   │
│     ... (97 more dormant roles)     │
│   }                                 │
│ }                                   │
└─────────────────────────────────────┘
  ↓
Checkpoint saved:
.oneshot/bd-f0x-checkpoint.json
```

### Phase 4: Parallel Execution (Main Loop)

```
Orchestrator: while True:
  ├─ ready = beads.get_ready_tasks()
  ├─ feature_ws = filter_by_parent(ready, "bd-f0x")
  │
  ├─ if not feature_ws:
  │   └─ break  # Все workstreams завершены
  │
  ├─ Dispatch to developers:
  │  for ws in feature_ws:
  │    developer = get_idle_developer()
  │    developer.execute(ws)
  │
  └─ Wait for completion
     └─ Update checkpoint
```

### Developer Agent Execution

```
Developer agent executing WS:
┌─────────────────────────────────────┐
│ Developer_3 (active)                │
│ ├─ Receive workstream WS-0xx-03     │
│ ├─ Read: docs/workstreams/...       │
│ ├─ Create worktree (git worktree)   │
│ ├─ Execute TDD cycle:               │
│ │   ├─ Red: Write failing test      │
│ │   ├─ Green: Write implementation  │
│ │   └─ Refactor: Clean up           │
│ ├─ Run quality gates                │
│ ├─ Move WS: backlog → completed     │
│ └─ Close worktree                   │
└─────────────────────────────────────┘
  ↓
Discover concern: "Auth service has no rate limiting"
  ↓
SendMessage("architect", "Concern: WS-0xx-03 missing rate limiting")
  ↓
┌─────────────────────────────────────┐
│ Architect agent (idle → active)     │
│ ├─ Receive message from developer   │
│ ├─ Review code                      │
│ ├─ Decision: Not critical for MVP   │
│ └─ Create Beads bug                 │
└─────────────────────────────────────┘
  ↓
Beads bug created:
bd-f0x-bug-001: "Add rate limiting to auth service"
severity: P2
reporter: architect
feature: bd-f0x
  ↓
Architect returns to idle
```

### Agent Coordination Examples

**Security review (dormant agent activation):**

```
Developer: "Implementing OAuth2 flow"
  ↓
SendMessage("security", "Review OAuth2 implementation for F0xx")
  ↓
┌─────────────────────────────────────┐
│ Security agent (dormant → active)   │
│ ├─ Receive wakeup message           │
│ ├─ Load: prompts/security_prompt.md │
│ ├─ Review code + docs               │
│ ├─ Check: PKCE flow, token storage  │
│ ├─ Decision: CONCERN found          │
│ └─ Create Beads bug (P1)            │
└─────────────────────────────────────┘
  ↓
Security: SendMessage("developer", "Bug created: bd-f0x-bug-002")
  ↓
Security returns to dormant
```

**SRE observability (dormant agent activation):**

```
DevOps: "Deploying auth service"
  ↓
SendMessage("sre", "Add observability for auth service")
  ↓
┌─────────────────────────────────────┐
│ SRE agent (dormant → active)        │
│ ├─ Review deployment                │
│ ├─ Decision: Needs metrics          │
│ ├─ Create Beads bug (P2)            │
│ └─ SendMessage("devops", "Metrics added")│
└─────────────────────────────────────┘
  ↓
SRE returns to dormant
```

### Phase 5: Completion & UAT

```
Orchestrator: All workstreams complete
  ↓
┌─────────────────────────────────────┐
│ QA agents (active)                  │
│ ├─ Run full test suite              │
│ ├─ Verify acceptance criteria       │
│ ├─ Check coverage ≥80%              │
│ └─ Aggregate bug reports            │
└─────────────────────────────────────┘
  ↓
Execution report:
{
  "feature": "bd-f0x",
  "workstreams": 5,
  "completed": 5,
  "bugs_found": 3,
  "bugs": [
    "bd-f0x-bug-001: rate limiting (P2)",
    "bd-f0x-bug-002: token storage (P1)",
    "bd-f0x-bug-003: auth metrics (P2)"
  ],
  "coverage": 87%
}
  ↓
Final UAT Gate:
├─ If --skip-uat: AUTO-APPROVE
├─ Else: present report + stage URL
│   ├─ Human tests on stage
│   ├─ Approve → @deploy
│   └─ Reject → create hotfix WS
└─ Commit: "UAT passed for F0xx"
```

### Checkpoint Structure (Fault Tolerance)

```json
{
  "feature": "bd-f0x",
  "team_id": "feature-f0x",
  "status": "in_progress",
  "started_at": "2026-01-28T10:00:00Z",
  "last_update": "2026-01-28T12:30:00Z",
  "completed_ws": ["bd-f0x-01", "bd-f0x-02"],
  "pending_ws": ["bd-f0x-03", "bd-f0x-04", "bd-f0x-05"],
  "bugs_created": 3,
  "agents_active": 7,
  "resume_count": 0
}
```

---

## Error Handling & Edge Cases

### Error Recovery Strategy

**1. Workstream Execution Failures**

```
Developer agent executing WS-0xx-03:
├─ Try: Implement auth service
├─ Exception: Tests failing after 3 refactor attempts
├─ Decision: Workstream blocked
└─ Actions:
   ├─ Create Beads bug: "WS-0xx-03 blocked: test failures"
   ├─ Move WS: backlog → blocked
   ├─ Mark checkpoint: "ws_blocked"
   └─ Continue: Execute other ready workstreams
```

**Recovery flow:**
```
Orchestrator detects blocked WS
  ↓
Checkpoint updated
  ↓
Options presented to human:
┌─────────────────────────────────────┐
│ "WS-0xx-03 blocked. Choose action:  │
│                                      │
│  [1] Skip - Mark as WONTFIX         │
│  [2] Retry - Resume from checkpoint │
│  [3] Debug - Run /debug skill       │
│  [4] Manual - I'll fix it myself    │
│  [5] Abort - Cancel entire feature  │
│ "                                    │
└─────────────────────────────────────┘
```

**2. Agent Communication Failures**

```
Developer → SendMessage("security", "Review auth")
  ↓
[ERROR] Security agent not responding
  ↓
Error handling:
├─ Log to team state: "security agent unreachable"
├─ Retry: 3 attempts with exponential backoff
├─ Fallback: Create Beads bug with @security mention
└─ Continue: Developer proceeds without review
```

**3. Orchestrator Crashes**

```
Orchestrator agent crashes mid-execution
  ↓
Checkpoints already saved
  ↓
Recovery flow:
@oneshot bd-f0x --resume
  ↓
┌─────────────────────────────────────┐
│ Orchestrator (resumed)              │
│ ├─ Load checkpoint                  │
│ ├─ Detect: agent_id changed         │
│ ├─ Verify: Beads state matches     │
│ ├─ Continue from WS-0xx-03          │
│ └─ Update: resume_count += 1       │
└─────────────────────────────────────┘
```

**4. Team State Corruption**

```
Team config corrupted: ~/.claude/teams/feature-f0x.json
  ↓
Error detected on TeamManager.load()
  ↓
Recovery strategy:
├─ Backup: Copy corrupted file to .backup
├─ Reconstruct: From Beads state
├─ Validate: Checkpoint consistency
└─ Fallback: Create new team
```

**5. Beads Client Failures**

```
Orchestrator: beads.get_ready_tasks()
  ↓
[ERROR] Beads CLI not responding
  ↓
Fallback strategy:
├─ Detect: 3 consecutive failures
├─ Switch: Use MockBeadsClient (in-memory state)
├─ Log: "Beads CLI down, using mock mode"
└─ Continue: Execution proceeds with mock state
```

**Mock fallback with sync:**
```python
class HybridBeadsClient:
    def __init__(self):
        self.cli = CLIBeadsClient()
        self.mock = MockBeadsClient()
        self.mode = "cli"

    def get_ready_tasks(self):
        try:
            return self.cli.get_ready_tasks()
        except BeadsError:
            logger.warning("CLI failed, switching to mock")
            self.mode = "mock"
            return self.mock.get_ready_tasks()

    def sync_on_recovery(self):
        """Sync mock state back to CLI when it recovers"""
        if self.mode == "mock":
            try:
                self.cli.list_tasks()
                for task in self.mock.tasks:
                    self.cli.create_task(task)
                self.mode = "cli"
            except BeadsError:
                pass
```

**6. Human Decision Timeouts**

```
RequirementsGate waiting for human approval
  ↓
Timeout: 24 hours (configurable)
  ↓
Action:
├─ Send reminder: "F0xx: Requirements pending approval"
├─ Auto-pause: Mark feature as "awaiting_human"
└─ Resume on next @feature invocation
```

**7. Git Conflict Scenarios**

```
Developer_1: WS-0xx-01 (modifies auth.py)
Developer_2: WS-0xx-02 (modifies auth.py)
  ↓
Both complete, attempt to merge
  ↓
Git conflict detected
  ↓
Resolution:
├─ Detect: Via pre-commit hooks
├─ Strategy: Sequentialize conflicting WS
├─ Fallback: Notify human to resolve
```

**8. Approval Gate Rejections**

```
RequirementsGate: Human rejects @idea output
  ↓
@feature menu:
┌─────────────────────────────────────┐
│ "Requirements rejected. Options:     │
│                                      │
│  [1] Revise - Analyst tries again   │
│  [2] Modify - I'll edit manually    │
│  [3] Skip - Proceed to @design      │
│  [4] Abort - Cancel feature         │
│ "                                    │
└─────────────────────────────────────┘
```

### Notification System (Multi-Channel)

**Notification providers** (плагинная архитектура):

```
NotificationSystem
├─ TelegramNotifier (базовый)
├─ SlackNotifier (future)
├─ DiscordNotifier (future)
├─ EmailNotifier (future)
└─ WebhookNotifier (custom)
```

**Provider interface:**
```python
class NotificationProvider(ABC):
    @abstractmethod
    async def send(self, message: NotificationMessage) -> bool:
        """Send notification, return success"""
        pass

    @abstractmethod
    def supports_interactive(self) -> bool:
        """Can this provider show buttons?"""
        pass
```

**Telegram provider** (первая реализация):

```python
class TelegramNotifier(NotificationProvider):
    """Telegram notification provider"""

    def __init__(self, bot_token, chat_ids):
        self.bot_token = bot_token
        self.chat_ids = chat_ids

    async def send(self, message: NotificationMessage) -> bool:
        """Send to Telegram"""
        # Implementation...

    def supports_interactive(self) -> bool:
        return True  # Inline keyboards
```

**Alert categories:**
- **Info**: Feature started, workstream complete, requirements approved
- **Warning**: Bug created (P2/P3), workstream blocked
- **Critical**: Execution failed, bug created (P0/P1), agent unreachable

**Configuration:**
```yaml
# config/notifications.yml
notifications:
  enabled_providers:
    - telegram
    # - slack  # Future

  telegram:
    bot_token: ${TELEGRAM_BOT_TOKEN}
    features_chat_id: "@sdp_features"
    alerts_chat_id: "@sdp_alerts"
```

**Adding new provider:**
```python
# Discord example
class DiscordNotifier(NotificationProvider):
    async def send(self, message: NotificationMessage) -> bool:
        # Discord webhook implementation
        pass
```

### Error Recovery Matrix

| Error | Detection | Recovery | Auto-resume |
|-------|-----------|----------|-------------|
| WS test failure | Post-build hooks | Block WS + bug report | No |
| Agent unreachable | SendMessage timeout | Fallback bug report | Yes |
| Orchestrator crash | Checkpoint missing | `--resume` flag | Yes |
| Team config corrupt | JSON parse error | Reconstruct from Beads | Yes |
| Beads CLI down | 3 consecutive errors | MockBeadsClient | Yes |
| Human timeout | 24h timer | Pause + reminder | Manual |
| Git conflict | Merge fails | Sequentialize | No |
| Approval rejection | Human input | Menu of options | Manual |

---

## Testing Strategy

### Testing Pyramid

```
         ┌──────────────┐
         │   E2E Tests  │  @feature workflow (rare)
         │    < 5%      │  Full @idea → @design → @oneshot
         ├──────────────┤
         │ Integration  │  Agent coordination (common)
         │    ~ 25%     │  Team messaging, Beads integration
         ├──────────────┤
         │   Unit Tests │  Individual components (heavy)
         │    ~ 70%     │  Orchestrator, TeamManager, Notifier
         └──────────────┘
```

### 1. Unit Tests (70% coverage)

```python
class TestOrchestrator:
    def test_create_team_with_all_roles(self):
        """Orchestrator creates team with 100+ roles"""
        orchestrator = Orchestrator(mock_beads_client)
        team = orchestrator.create_team("feature-f01")

        assert len(team.roles) >= 100
        assert team.roles["developer_1"]["status"] == "active"
        assert team.roles["security"]["status"] == "dormant"

    def test_checkpoint_save_and_resume(self):
        """Checkpoint persists state correctly"""
        orchestrator = Orchestrator(mock_beads_client)

        checkpoint = orchestrator.create_checkpoint(
            feature="bd-f01",
            completed=["ws-1", "ws-2"],
            pending=["ws-3", "ws-4"]
        )

        orchestrator.save_checkpoint(checkpoint)

        # Load in new orchestrator
        new_orchestrator = Orchestrator(mock_beads_client)
        loaded = new_orchestrator.load_checkpoint("bd-f01")

        assert loaded.feature == "bd-f01"
        assert loaded.completed_ws == ["ws-1", "ws-2"]
```

### 2. Integration Tests (25% coverage)

```python
class TestTeamCoordination:
    @pytest.mark.asyncio
    async def test_dormant_agent_wakeup(self):
        """Security agent wakes up on message"""
        team = TeamManager.create("feature-test")
        assert team.roles["security"]["status"] == "dormant"

        await team.send_message(
            from_role="developer",
            to_role="security",
            content="Review auth implementation"
        )

        assert team.roles["security"]["status"] == "active"

    @pytest.mark.asyncio
    async def test_parallel_workstream_execution(self):
        """Multiple developers execute in parallel"""
        orchestrator = Orchestrator(BeadsClient())
        team = orchestrator.create_team("feature-test", developer_count=3)

        ws1 = create_workstream("ws-1", dependencies=[])
        ws2 = create_workstream("ws-2", dependencies=[])
        ws3 = create_workstream("ws-3", dependencies=[])

        results = await orchestrator.execute_parallel([ws1, ws2, ws3])

        assert all(r.success for r in results)
```

### 3. E2E Tests (5% coverage)

```python
class TestFeatureWorkflow:
    @pytest.mark.e2e
    def test_idea_to_oneshot_flow(self):
        """Full @idea → @design → @oneshot pipeline"""

        # Step 1: @idea
        idea_output = run_idea_skill("Add user notifications")
        assert idea_output["requirements_file"] exists

        # Step 2: Requirements gate
        gate = RequirementsGate(idea_output)
        decision = gate.ask_human()
        assert decision == "approve"

        # Step 3: @design
        design_output = run_design_skill(idea_output)
        assert design_output["workstreams_count"] >= 3

        # Step 4: @oneshot (with mock team)
        with mock_team_execution():
            result = run_oneshot_skill(design_output["feature_id"])
            assert result["status"] == "completed"
```

### 4. Real Integration E2E Tests

**Beads Integration:**

```python
class TestBeadsIntegration:
    @pytest.fixture(scope="class")
    def real_beads_client(self, tmp_path_factory):
        """Create real Beads database for testing"""
        temp_dir = tmp_path_factory.mktemp("beads_e2e")

        subprocess.run(
            ["bd", "init", "--repo-id", "test-e2e"],
            cwd=temp_dir,
            capture_output=True
        )

        client = CLIBeadsClient(beads_dir=temp_dir)
        yield client

        # Cleanup
        shutil.rmtree(temp_dir)

    def test_create_parent_and_child_tasks(self, real_beads_client):
        """Create feature task with workstream children"""

        parent = real_beads_client.create_task(BeadsTask(
            title="F099: Test Feature",
            issue_type="feature"
        ))

        ws1 = create_task(real_beads_client, "WS-1", parent)
        ws2 = create_task(real_beads_client, "WS-2", parent, deps=[ws1])

        # Verify get_ready_tasks()
        ready = real_beads_client.get_ready_tasks()
        assert ws1 in ready
        assert ws2 not in ready  # Blocked

        real_beads_client.mark_done(ws1)

        # Now ws-2 should be ready
        ready = real_beads_client.get_ready_tasks()
        assert ws2 in ready
```

**Telegram Alerts:**

```python
@pytest.mark.e2e
@pytest.mark.telegram
class TestTelegramAlerts:
    @pytest.fixture(scope="class")
    def telegram_notifier(self):
        """Create real Telegram notifier for test"""
        token = os.getenv("TELEGRAM_TEST_BOT_TOKEN")
        if not token:
            pytest.skip("TELEGRAM_TEST_BOT_TOKEN not set")

        return TelegramNotifier(
            bot_token=token,
            features_chat_id=os.getenv("TELEGRAM_TEST_CHAT_ID")
        )

    @pytest.mark.asyncio
    async def test_send_feature_start_notification(self, telegram_notifier):
        """Send feature start alert to Telegram"""

        await telegram_notifier.feature_started(
            feature_id="F099-E2E",
            name="E2E Test Feature",
            user="@test_runner"
        )

        # Verify via bot updates API
        updates = await telegram_notifier.get_updates(timeout=5)
        assert any("F099-E2E" in u.message for u in updates)
```

### Test Execution

```bash
# Unit tests (fast)
pytest tests/unit/ -v

# Integration tests (medium)
pytest tests/integration/ -v

# E2E tests (slow)
pytest tests/e2e/ -v -m e2e

# Coverage report
pytest --cov=src/sdp --cov-report=html
```

---

## Success Metrics

### Developer Experience Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Time to first feature** | < 15 min | Onboarding tutorial completion |
| **Cognitive load** | < 50 rule occurrences | Count rule violations in docs |
| **Manual state tracking** | 0 manual file moves | Automated WS state transitions |
| **CLI command complexity** | < 3 args per command | Average args count |

### Workflow Efficiency Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Idea → Deploy cycle time** | < 4 hours (pure vibe) | End-to-end duration |
| **Approval gate overhead** | < 5 min per gate | Human interaction time |
| **Agent parallelization** | ≥ 2 agents avg | Concurrent active agents |
| **Checkpoint resume success** | > 95% | Resume completion rate |

### Quality Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Test coverage** | ≥ 80% | pytest --cov |
| **Bugs per feature** | < 3 bugs | Beads bug count |
| **Critical bugs (P0/P1)** | 0 | Bug severity distribution |
| **Code review pass rate** | > 90% | First-time approval rate |

### Reliability Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Orchestrator crash recovery** | > 98% | Successful resumes |
| **Agent communication success** | > 99% | SendMessage success rate |
| **Beads client uptime** | > 99.5% | CLI availability |
| **Notification delivery** | > 95% | Telegram/Slack success |

### Feature Success Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Feature completion rate** | > 85% | Deployed / Started |
| **Workstream blocking rate** | < 10% | Blocked / Total WS |
| **Human intervention frequency** | < 3 per feature | Manual decisions |
| **UAT pass rate** | > 90% | First-time UAT approval |

### Notification Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Notification latency** | < 5 sec | Send → Deliver time |
| **False positive rate** | < 5% | Non-critical alerts |
| **Interactive button response** | > 80% | Button click rate |

---

## Implementation Roadmap

### Phase 1: Core Infrastructure (Week 1-2)
- [ ] Implement Orchestrator agent
- [ ] Implement Team Manager with 100+ roles
- [ ] Implement Approval Gate Manager
- [ ] Implement Checkpoint/Resume system

### Phase 2: @feature Skill (Week 3)
- [ ] Implement @feature skill with progressive menu
- [ ] Integrate @idea/@design/@oneshot invocation
- [ ] Implement skip-flags logic

### Phase 3: Agent Runtime (Week 4-5)
- [ ] Implement Agent Spawner via Task tool
- [ ] Implement SendMessage router
- [ ] Implement dormant/active role switching
- [ ] Implement bug report flow

### Phase 4: Notification System (Week 6)
- [ ] Implement Notification Router
- [ ] Implement TelegramNotifier
- [ ] Implement alert categories
- [ ] Configuration system

### Phase 5: Testing & Validation (Week 7-8)
- [ ] Unit tests (70% coverage)
- [ ] Integration tests (25% coverage)
- [ ] E2E tests with real Beads
- [ ] E2E tests with real Telegram

### Phase 6: Documentation & Onboarding (Week 9)
- [ ] Update PROTOCOL.md
- [ ] Create "First 15 Minutes" tutorial
- [ ] Create role configuration guide
- [ ] Create notification setup guide

---

**Status:** ✅ Design Complete
**Next Steps:** Implementation planning with @design skill
