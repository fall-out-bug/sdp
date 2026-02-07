# Agent Roles Setup Guide

**SDP Agent System** - multi-agent coordination for feature development.

---

## What are Agent Roles?

Agent roles define specialized AI agents with specific capabilities and responsibilities. Each role has:

- **Name** - Role identifier (e.g., `planner`, `builder`)
- **Purpose** - What the role does
- **Capabilities** - Specific skills and tasks
- **Prompt** - System prompt for Claude agent

```
┌─────────────────────────────────────────────────────┐
│              Orchestrator Agent                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │ Planner  │  │ Builder  │  │ Reviewer │         │
│  └──────────┘  └──────────┘  └──────────┘         │
│  ┌──────────┐                                       │
│  │ Deployer │                                       │
│  └──────────┘                                       │
└─────────────────────────────────────────────────────┘
```

---

## Built-in Roles

### 1. Planner Agent (`planner.md`)

**Purpose:** Break features into workstreams

**Capabilities:**
- Analyze feature requirements
- Design workstream decomposition
- Create dependency graphs
- Estimate workstream size (SMALL/MEDIUM/LARGE)

**When to use:**
```bash
@design beads-XXX
```

**Example output:**
```
WS-XXX.01: Domain model (450 LOC, MEDIUM)
WS-XXX.02: Database schema (300 LOC, MEDIUM)
WS-XXX.03: Repository layer (500 LOC, MEDIUM)
```

### 2. Builder Agent (`builder.md`)

**Purpose:** Execute workstreams with TDD

**Capabilities:**
- Test-Driven Development (Red → Green → Refactor)
- Write clean, testable code
- Follow quality gates (coverage, mypy, ruff)
- Commit work when complete

**When to use:**
```bash
@build WS-XXX.01
```

**Workflow:**
```python
# 1. Red: Write failing test
def test_feature():
    assert feature_not_implemented()

# 2. Green: Implement minimum code
def feature():
    return "working"

# 3. Refactor: Improve design
def feature_refactored():
    return "clean code"
```

### 3. Reviewer Agent (`reviewer.md`)

**Purpose:** Quality validation of features

**Capabilities:**
- Validate quality gates
- Check test coverage ≥80%
- Verify mypy --strict compliance
- Review for tech debt
- Return verdict: APPROVED / CHANGES_REQUESTED

**When to use:**
```bash
@review beads-XXX
```

**Quality checklist:**
- ✅ Tests first (TDD)
- ✅ Coverage ≥80%
- ✅ mypy --strict
- ✅ ruff clean
- ✅ Files <200 LOC
- ✅ No `except: pass`
- ✅ Type hints

### 4. Deployer Agent (`deployer.md`)

**Purpose:** Production deployment

**Capabilities:**
- Generate deployment configs (docker-compose, CI/CD)
- Create PR with changelog
- Run smoke tests
- Merge to main with tagging

**When to use:**
```bash
@deploy beads-XXX
```

**Artifacts:**
- `docker-compose.yml`
- `.github/workflows/deploy.yml`
- `CHANGELOG.md` entry
- Git tag: `v{version}`

### 5. Orchestrator Agent (`orchestrator.md`)

**Purpose:** Coordinate all agents

**Capabilities:**
- Spawn specialized agents
- Route messages between agents
- Manage agent lifecycle
- Handle checkpoints
- Send notifications

**When to use:**
```bash
@oneshot beads-XXX
```

**Workflow:**
```
1. Spawn planner → Create workstreams
2. Spawn builder → Execute each WS
3. Spawn reviewer → Validate quality
4. Send notifications → Progress updates
5. Save checkpoints → Resume support
```

---

## Creating Custom Roles

### Role File Format

Create file: `.claude/agents/{role-name}.md`

```markdown
# {Role Name}

{One-line description of what this role does}

## Purpose

{Detailed explanation of role's purpose}

## Capabilities

- **{Capability 1}**: {Description}
- **{Capability 2}**: {Description}
- **{Capability 3}**: {Description}

## When to Use

{When this agent should be spawned}

## Workflow

{Step-by-step process}

## Examples

{Code examples if applicable}
```

### Example: Security Reviewer

Create `.claude/agents/security-reviewer.md`:

```markdown
# Security Reviewer

Reviews code for security vulnerabilities and best practices.

## Purpose

Identify security issues before code reaches production.

## Capabilities

- **SQL Injection Detection**: Find unsafe query patterns
- **XSS Prevention**: Check output encoding
- **Authentication**: Verify auth logic
- **Authorization**: Check permission checks

## When to Use

```bash
@security-review beads-XXX
```

## Workflow

1. Review all database queries
2. Check authentication flows
3. Verify authorization logic
4. Test for common vulnerabilities (OWASP Top 10)
5. Generate security report

## Examples

```python
# ❌ Bad: SQL injection risk
query = f"SELECT * FROM users WHERE id={user_id}"

# ✅ Good: Parameterized query
query = "SELECT * FROM users WHERE id=?"
cursor.execute(query, (user_id,))
```
```

---

## Role Activation

### Using RoleLoader

```python
from sdp.unified.agent.role_loader import RoleLoader
from sdp.unified.agent.role_state import RoleStateManager

# Load role from file
loader = RoleLoader(agents_dir=".claude/agents")
role = loader.load_role("planner")

# Activate role
state_mgr = RoleStateManager()
state_mgr.activate_role("planner")

# Check active roles
active = state_mgr.list_active()  # ["planner"]
```

### Manual Role Switching

```python
# Switch roles during execution
state_mgr.deactivate_role("planner")
state_mgr.activate_role("builder")

assert state_mgr.is_active("builder")
assert not state_mgr.is_active("planner")
```

---

## Best Practices

### 1. Single Responsibility

Each role should have **one clear purpose**.

❌ **Bad:** `fullstack-dev.md` - Does everything
✅ **Good:** `builder.md` - Executes workstreams only

### 2. Clear Capabilities

List specific capabilities, not vague goals.

❌ **Bad:** "Can do development tasks"
✅ **Good:** "Execute TDD cycle (Red → Green → Refactor)"

### 3. When to Use

Explicitly state when to spawn this role.

```markdown
## When to Use

```bash
@build WS-XXX.01
```

Or automatically:
- After @design creates workstreams
- Before @review validates quality
```

### 4. Examples

Provide runnable examples for role-specific tasks.

```python
# Good: Complete example
def test_feature():
    client = create_client()
    response = client.get("/api/users")
    assert response.status_code == 200
    assert len(response.json()) > 0
```

### 5. Role Composition

Roles can use other roles via agent spawning.

```python
# Orchestrator spawns planner
planner_id = spawner.spawn_agent(AgentConfig(
    name="planner",
    prompt="Break feature into workstreams",
))

# Planner returns workstreams
# Orchestrator spawns builder for each WS
builder_id = spawner.spawn_agent(AgentConfig(
    name="builder",
    prompt=f"Execute {workstream}",
))
```

---

## Agent Communication

### Send Message to Agent

```python
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()

# Send message
message = Message(
    sender="orchestrator",
    content="Execute WS-060-01: Domain model",
    recipient=builder_id,
)

result = router.send_message(message)
assert result.success
```

### Receive Messages

Agents automatically receive messages via:
- `recipient` field in Message
- Agent listens for messages with its ID
- Process message and respond

---

## Role Templates

Copy these templates for new roles:

### Template 1: Specialist Agent

```markdown
# {Specialist}

{One-line description}

## Purpose

{What this specialist does}

## Capabilities

- **{Task 1}**: {How it's done}
- **{Task 2}**: {How it's done}

## When to Use

```bash
@{specialist} {target}
```

## Workflow

1. {Step 1}
2. {Step 2}
3. {Step 3}

## Examples

{Code examples}
```

### Template 2: Review Agent

```markdown
# {Reviewer}

Reviews {domain} for {criteria}.

## Purpose

Ensure {quality standard} is met.

## Capabilities

- **Check {aspect 1}**: {Method}
- **Check {aspect 2}**: {Method}
- **Report findings**: {Format}

## When to Use

```bash
@{reviewer} {target}
```

## Checklist

- ✅ {Check 1}
- ✅ {Check 2}
- ✅ {Check 3}

## Examples

{Before/After comparisons}
```

---

## Troubleshooting

### Role Not Loading

**Problem:** Role loads as `None`

**Solution:**
```bash
# Check file exists
ls .claude/agents/{role}.md

# Check file format
head .claude/agents/{role}.md

# Must start with: # {Role Name}
```

### Agent Not Responding

**Problem:** Agent spawned but doesn't process messages

**Solution:**
```python
# Verify agent ID
print(f"Agent ID: {agent_id}")

# Check message routing
result = router.send_message(message)
print(f"Success: {result.success}")
print(f"Error: {result.error}")
```

### Role Not Activating

**Problem:** `state_mgr.is_active(role)` returns False

**Solution:**
```python
# Check role was loaded
role = loader.load_role("my-role")
assert role is not None

# Check activation
state_mgr.activate_role("my-role")
assert state_mgr.is_active("my-role")
```

---

## See Also

- `PROTOCOL.md` - Full SDP specification (English)
- `PROTOCOL_RU.md` - Полная спецификация (Русский)
- `docs/TUTORIAL.md` - 15-minute quick start
- `src/sdp/unified/agent/README.md` - Agent system internals

---

**Version:** SDP v0.6.0
**Updated:** 2026-01-29
