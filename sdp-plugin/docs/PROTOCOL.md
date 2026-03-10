# SDP: Spec-Driven Protocol

**Workstream-driven development** for AI agents with multi-language support.

**Plugin Version:** Language-agnostic (Python, Java, Go)

> **🎯 Documentation Navigation:** See [NAVIGATION.md](NAVIGATION.md) for the complete documentation index with decision trees and progressive disclosure.

---

## Quick Start

```bash
# Install plugin (no Python required)
git clone https://github.com/fall-out-bug/sdp.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/

# Create feature (interactive)
@feature "Add user authentication"

# Plan workstreams
@design feature-auth

# Execute workstream
@build 00-001-01

# Review quality
@review feature-auth

# Deploy to production
@deploy feature-auth
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

Every workstream must pass quality gates. **Commands are language-specific:**

### Test Coverage ≥80%

| Language | Command |
|----------|---------|
| **Python** | `pytest tests/unit/ --cov=src/ --cov-fail-under=80` |
| **Java** | `mvn verify` (JaCoCo report) or `gradle test jacocoTestReport` |
| **Go** | `go test -coverprofile=coverage.out && go tool cover -func=coverage.out` |

### Type Checking

| Language | Command |
|----------|---------|
| **Python** | `mypy src/ --strict` |
| **Java** | `javac -Xlint:all` (compiler checks) |
| **Go** | `go vet ./...` |

### Linting

| Language | Command |
|----------|---------|
| **Python** | `ruff check src/` |
| **Java** | `mvn checkstyle:check` or `gradle checkstyleMain` |
| **Go** | `golangci-lint run ./...` or `staticcheck ./...` |

### File Size <200 LOC

| Language | Command |
|----------|---------|
| **All** | Count lines per file (language-agnostic) |

**Forbidden:**
- ❌ Bare exceptions (`except:`, `catch(Exception e)`, `recover()` without check)
- ❌ Files > 200 LOC
- ❌ Coverage < 80%
- ❌ Time-based estimates
- ❌ TODO without followup WS

---

## Language-Specific Workflows

### Python Projects

```bash
# Prerequisites
pip install pytest pytest-cov mypy ruff

# Workflow
@feature "Add REST API"
@design feature-rest-api
@build 00-001-01
# Runs: pytest tests/ -v
# Quality: pytest --cov=src/ --cov-fail-under=80, mypy src/ --strict, ruff check src/
```

### Java Projects

```bash
# Prerequisites
# Maven: mvn verify (runs JaCoCo)
# Gradle: gradle test jacocoTestReport

# Workflow
@feature "Add REST API"
@design feature-rest-api
@build 00-001-01
# Runs: mvn test
# Quality: JaCoCo coverage ≥80%, javac -Xlint:all, checkstyle
```

### Go Projects

```bash
# Prerequisites
# Go 1.21+ with go tool cover, go vet, golangci-lint

# Workflow
@feature "Add REST API"
@design feature-rest-api
@build 00-001-01
# Runs: go test ./...
# Quality: go tool cover -func=coverage.out (≥80%), go vet ./..., golangci-lint run ./...

# Style: prefer modern stdlib idioms such as slices.SortFunc and strings.Cut
```

---

## Project Type Detection

SDP automatically detects project type:

1. **Python**: `pyproject.toml`, `setup.py`, or `requirements.txt` present
2. **Java**: `pom.xml` or `build.gradle` present
3. **Go**: `go.mod` present
4. **Node.js**: `package.json` present
5. **Rust**: `Cargo.toml` present

If multiple build files exist, SDP prompts user to specify.

---

## Multi-Agent Coordination

SDP integrates multi-agent coordination with role-based message routing.

### Agent Spawning

```python
from sdp.unified.agent.spawner import AgentSpawner, AgentConfig

# Spawn agents
spawner = AgentSpawner()
builder = spawner.spawn_agent(AgentConfig(
    name="builder",
    prompt="Execute workstreams with TDD...",
))
```

### Message Routing

```python
from sdp.unified.agent.router import SendMessageRouter, Message

router = SendMessageRouter()
router.send_message(Message(
    sender="orchestrator",
    content="Execute 00-060-01",
    recipient=builder,
))
```

---

## Workstream Format

Workstreams use PP-FFF-SS format:
- **PP**: Project ID (00 = SDP core, 01-99 = custom)
- **FFF**: Feature ID (001-999)
- **SS**: Sequence (01-99)

Example: `00-001-01` = Project 00, Feature 001, Workstream 01

---

## Quality Gate Enforcement

### AI-Based Validation

SDP plugin uses AI validators instead of static analysis tools:

1. **Coverage Validator**: Analyzes test coverage by reading code
2. **Architecture Validator**: Checks Clean Architecture layer separation
3. **Error Validator**: Finds unsafe exception handling
4. **Complexity Validator**: Identifies overly complex code

**Gate:** All validators must PASS for workstream approval.

### Tool-Based Validation (Optional)

If language tools are available, SDP uses them:

```bash
# Python: pytest, mypy, ruff
# Java: Maven/Gradle, JaCoCo, checkstyle
# Go: go test, go tool cover, go vet
```

---

## Documentation

- [Tutorial](TUTORIAL.md) - Full workflow guide
- [Python Examples](examples/python/) - Python-specific guides
- [Java Examples](examples/java/) - Java-specific guides
- [Go Examples](examples/go/) - Go-specific guides
- [Migration Guide](MIGRATION.md) - From Python SDP to plugin

---

## @feature Unified Workflow

The @feature skill provides a **progressive disclosure workflow** for feature development, from vision to execution.

### Overview

@feature guides users through 7 phases:

```
Phase 1: Vision Interview       → AskUserQuestion (3-5 questions)
Phase 2: Generate Vision       → PRODUCT_VISION.md
Phase 3: Technical Interview   → AskUserQuestion (5-8 questions)
Phase 4: Generate Intent       → docs/intent/{slug}.json
Phase 5: Requirements Draft    → docs/drafts/idea-{slug}.md
Phase 6: Workstream Planning   → @design skill
Phase 7: Orchestrator Execution→ @oneshot skill (optional)
```

### Phase 1: Vision Interview

**Goal:** Understand product vision, users, and success metrics.

**Questions asked:**
1. **Mission**: What problem do we solve?
2. **Users**: Who are we building for?
3. **Success Metrics**: How do we measure success?

**Example:**
```markdown
User: @feature "Add user authentication"

Claude: Let me ask a few questions about your vision:

🎯 Mission: What problem does user authentication solve for your product?

User: We need secure login to protect user data and enable personalized experiences.

Claude: 👥 Users: Who are the primary users of this authentication system?

User: End users logging in, plus administrators managing accounts.

Claude: 📊 Success Metrics: How will you measure the success of this feature?

User: Login success rate >95%, account creation time <2 minutes, zero security breaches.
```

### Phase 2: Generate PRODUCT_VISION.md

**Goal:** Create product manifesto at project root.

**Output: `PRODUCT_VISION.md`**
```markdown
# PRODUCT_VISION.md

> **Last Updated:** 2026-02-06
> **Version:** 1.0

## Mission

Enable secure user authentication with personalized experiences while protecting user data.

## Users

1. **End Users** - Log in to access personalized features
2. **Administrators** - Manage user accounts and security

## Success Metrics

- [ ] Login success rate >95%
- [ ] Account creation time <2 minutes
- [ ] Zero security breaches
- [ ] Support for OAuth 2.0 providers

## Strategic Tradeoffs

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| Security | Multi-factor auth required | Balance security vs. UX |
| Storage | Encrypted passwords only | Compliance with GDPR |
| Session | JWT tokens with refresh | Scalability vs. complexity |

## Non-Goals

- Not supporting legacy password hashing
- Not implementing SSO (future feature)
```

**Decision Logging:**
```bash
sdp decisions log --type="tradeoff" \
  --question="Security Approach" \
  --decision="Multi-factor auth required" \
  --rationale="Balance security vs. UX" \
  --feature-id="F001" \
  --maker="user"
```

### Phase 3: Technical Interview

**Goal:** Understand technical approach, tradeoffs, and integration points.

**Questions asked:**
1. **Authentication Method**: Session-based vs. JWT vs. OAuth?
2. **Storage**: Where to store user credentials?
3. **Failure Modes**: What happens if login fails? Database down?
4. **Integration Points**: Existing user management? Email service?
5. **Security Requirements**: Password hashing? 2FA? Rate limiting?
6. **Scalability**: Expected concurrent users?
7. **Compliance**: GDPR? CCPA? SOC2?

**Example:**
```markdown
Claude: 🔧 Authentication Method: Would you prefer session-based auth, JWT tokens, or OAuth 2.0?

User: JWT tokens with refresh token rotation.

Claude: 💾 Storage: Where should we store user credentials?

User: PostgreSQL database with encrypted passwords.

Claude: ⚠️  Failure Modes: What should happen if the database is down during login?

User: Return 503 Service Unavailable, queue login attempts for retry.

[... 5 more questions ...]

Claude: ✅ Technical approach documented. Logging decisions...
```

### Phase 4: Generate Intent JSON

**Goal:** Create machine-readable intent file.

**Output: `docs/intent/user-authentication.json`**
```json
{
  "problem": "Users need secure authentication to access personalized features",
  "users": ["end_users", "administrators"],
  "success_criteria": [
    "login_success_rate > 95%",
    "account_creation_time < 2 minutes",
    "zero_security_breaches"
  ],
  "technical_approach": {
    "auth_method": "jwt",
    "storage": "postgresql",
    "security": "bcrypt_password_hashing"
  },
  "constraints": [
    "GDPR compliant",
    "OAuth 2.0 support"
  ]
}
```

**Validation:**
```bash
# Python
from sdp.schema.validator import IntentValidator
validator = IntentValidator()
validator.validate(intent_dict)  # Raises ValidationError if invalid

# Go
validator := NewIntentValidator()
if err := validator.Validate(intent); err != nil {
    log.Fatal(err)
}
```

### Phase 5: Requirements Draft

**Goal:** Create human-readable specification.

**Output: `docs/drafts/idea-user-authentication.md`**
```markdown
# User Authentication

> **Feature ID:** F001
> **Status:** Draft
> **Created:** 2026-02-06

## Problem

Users need secure authentication to access personalized features while protecting their data.

## Users

1. **End Users** - Need to log in and manage their accounts
2. **Administrators** - Need to manage user accounts and security

## Success Criteria

- [ ] Login success rate >95%
- [ ] Account creation time <2 minutes
- [ ] Zero security breaches
- [ ] OAuth 2.0 support

## Goals

1. Implement JWT-based authentication
2. Support email/password and OAuth login
3. Secure password storage with bcrypt
4. Refresh token rotation
5. Session management

## Non-Goals

- Single sign-on (SSO) - future feature
- Legacy password migration - out of scope
- Biometric auth - not in MVP

## Technical Approach

**Authentication:** JWT tokens with refresh token rotation
**Storage:** PostgreSQL with encrypted passwords (bcrypt)
**Security:** Rate limiting, account lockout, 2FA support
**Compliance:** GDPR compliant data handling

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───→│  Auth API   │───→│  Postgres   │
└─────────────┘    └─────────────┘    └─────────────┘
                          │
                          ▼
                   ┌─────────────┐
                   │ Email Svc   │
                   └─────────────┘
```

## Workstreams (Preliminary)

1. Domain models (User, Session, Token)
2. Authentication service
3. JWT token management
4. OAuth integration
5. API endpoints
6. Frontend integration
7. Testing (unit, integration, E2E)
8. Documentation
```

### Phase 6: Workstream Planning

**Goal:** Break feature into workstreams via @design skill.

**Call @design:**
```bash
@design idea-user-authentication
```

**@design analyzes:**
- Codebase structure
- Existing dependencies
- Team skills
- Risk factors

**Output: Workstreams in `docs/workstreams/backlog/`**
```markdown
# WS-001: Domain Models

> **Feature:** F001: User Authentication
> **Size:** SMALL

## Goal

Define User, Session, and Token entities with Clean Architecture.

## Acceptance Criteria

- User entity with email, password_hash, created_at
- Session entity with user_id, token, expires_at
- Token entity with jti, user_id, expires_at
- Repository interfaces defined
- Unit tests with ≥80% coverage

## Scope Files

**src/auth/domain/user.go**
**src/auth/domain/session.go**
**src/auth/domain/token.go**

## Definition of Done

- All 4 AC met
- Tests passing
- No TODOs
- Files <200 LOC
```

**Workstream breakdown example:**
```
F001: User Authentication
├── WS-001: Domain Models (SMALL)
├── WS-002: Auth Service (SMALL)
├── WS-003: JWT Management (SMALL)
├── WS-004: OAuth Integration (MEDIUM)
├── WS-005: API Endpoints (SMALL)
├── WS-006: Frontend Forms (SMALL)
├── WS-007: Integration Tests (MEDIUM)
└── WS-008: Documentation (SMALL)
```

### Phase 7: Orchestrator Execution (Optional)

**Goal:** Execute workstreams autonomously with checkpoint/resume.

**Call @oneshot:**
```bash
@oneshot F001
```

**Orchestrator executes:**
```
[15:23] Starting feature execution: F001
[15:23] Loading workstreams...
[15:23] Building dependency graph...
[15:23] Execution order: [WS-001 WS-002 WS-003 WS-004 WS-005 WS-006 WS-007 WS-008]

[15:24] Executing WS-001: Domain Models (1/8)...
[15:46] → WS-001 complete (22m, 85% coverage)

[15:46] Executing WS-002: Auth Service (2/8)...
[16:08] → WS-002 complete (22m, 89% coverage)

[16:08] Executing WS-003: JWT Management (3/8)...
[16:29] → WS-003 complete (21m, 82% coverage)

[16:29] Executing WS-004: OAuth Integration (4/8)...
[17:15] → WS-004 complete (46m, 81% coverage)

[17:15] Executing WS-005: API Endpoints (5/8)...
[17:35] → WS-005 complete (20m, 87% coverage)

[17:35] Executing WS-006: Frontend Forms (6/8)...
[17:55] → WS-006 complete (20m, 84% coverage)

[17:55] Executing WS-007: Integration Tests (7/8)...
[18:20] → WS-007 complete (25m, 91% coverage)

[18:20] Executing WS-008: Documentation (8/8)...
[18:30] → WS-008 complete (10m, docs complete)

[18:30] Feature execution complete: 8/8 workstreams, 3h 7m total, 86% avg coverage
```

**Checkpoint format:**
```json
{
  "id": "F001",
  "feature_id": "F001",
  "status": "in_progress",
  "completed_workstreams": ["WS-001", "WS-002", "WS-003"],
  "current_workstream": "WS-004",
  "created_at": "2026-02-06T15:23:00Z",
  "updated_at": "2026-02-06T17:15:00Z"
}
```

**Resume from checkpoint:**
```bash
@oneshot F001 --resume F001
# Resumes from WS-004 (OAuth Integration)
```

**Error handling:**
```
[16:29] Executing WS-004: OAuth Integration (4/8)...
[16:45] → WS-004 failed: OAuth provider API changed
[16:45] ⚠️  Retrying (attempt 1/2)...
[17:05] → WS-004 failed: Rate limit exceeded
[17:05] ⚠️  Retrying (attempt 2/2)...
[17:25] → WS-004 complete after retries (40m, 81% coverage)
```

### Progressive Menu System

Users can skip phases or start from existing specs:

**Power User Flags:**
- `--vision-only` - Only Phase 1-2 (stop before technical interview)
- `--no-interview` - Skip AskUserQuestion, use defaults
- `--spec PATH` - Load existing draft from docs/drafts/
- `--execute` - Automatically start orchestrator in Phase 7

**Examples:**
```bash
# Vision only
@feature "Add payments" --vision-only

# From existing spec
@feature --spec docs/drafts/idea-auth.md

# No interview (defaults)
@feature "Add notifications" --no-interview

# Full workflow with execution
@feature "Add analytics" --execute
```

### Beads Integration

@feature integrates with Beads for task tracking:

```bash
# Before Phase 1
bd create feature --title="F001: User Authentication" --description="..."

# Phase 1-6: Log decisions
sdp decisions log --type="technical" --question="Auth method?" --decision="JWT"

# Phase 7: Execute workstreams
for ws in WS-001 WS-002 ...; do
  bd create task --title="$ws" --parent="F001"
  bd update "$ws" --status in_progress
  @build "$ws"
  bd close "$ws" --reason="Complete"
done

# After feature complete
bd close "F001" --reason="Feature complete, all workstreams done"
```

### Workflow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     @feature Workflow                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Phase 1  │→ │ Phase 2  │→ │ Phase 3  │→ │ Phase 4  │   │
│  │  Vision  │  │ PRODUCT_ │  │Technical │  │  Intent  │   │
│  │Interview │  │ VISION.md│  │Interview │  │   JSON   │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                 ↓            │
│                                          ┌──────────┐       │
│                                          │ Phase 5  │       │
│                                          │  Draft   │       │
│                                          └──────────┘       │
│                                                 ↓            │
│                                          ┌──────────┐       │
│                                          │ Phase 6  │       │
│                                          │ @design  │       │
│                                          └──────────┘       │
│                                                 ↓            │
│                                          ┌──────────┐       │
│                                          │ Phase 7  │       │
│                                          │@oneshot  │       │
│                                          │(optional)│       │
│                                          └──────────┘       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Agent Interaction Diagram

```
┌──────────────┐         ┌──────────────┐
│    User      │         │  Orchestrator│
└──────┬───────┘         └──────┬───────┘
       │                        │
       │ @feature "Add X"       │
       ├───────────────────────>│
       │                        │
       │        Phase 1-6       │
       │                        │
       │  Ask questions         │
       │<───────────────────────┤
       │  Answer                │
       ├───────────────────────>│
       │                        │
       │  Generate docs         │
       │<───────────────────────┤
       │                        │
       │  Execute? (y/n)        │
       │<───────────────────────┤
       │  Yes                   │
       ├───────────────────────>│
       │                        │
       │         ┌──────────────┴──────────────┐
       │         │    @oneshot spawns agents   │
       │         ├─────────────────────────────┤
       │         │                             │
       │    ┌────┴────┐  ┌────┴────┐  ┌────┴─┴──┐
       │    │Builder  │  │Designer │  │Tester   │
       │    └────┬────┘  └────┬────┘  └────┬───┘
       │         │            │            │
       │         └────────────┴────────────┘
       │                      │
       │    Progress updates  │
       │<─────────────────────┤
       │                      │
       │  Feature complete    │
       │<─────────────────────┤
```

### Sequence Diagram

```
User       Orchestrator    @idea        @design       @oneshot
 │             │              │             │             │
 │ @feature    │              │             │             │
 ├────────────>│              │             │             │
 │             │              │             │             │
 │             │ Ask questions│             │             │
 │<────────────┤              │             │             │
 │ Answers     │              │             │             │
 ├────────────>│              │             │             │
 │             │              │             │             │
 │             │ Invoke       │             │             │
 │             ├─────────────>│             │             │
 │             │ Result       │             │             │
 │             │<─────────────┤             │             │
 │             │              │             │             │
 │             │ Invoke                    │             │
 │             ├──────────────────────────>│             │
 │             │ Workstreams               │             │
 │             │<──────────────────────────┤             │
 │             │              │             │             │
 │             │ Invoke                                 │
 │             ├────────────────────────────────────────>│
 │             │              │             │             │
 │ Progress    │              │             │             │
 │<────────────┤              │             │             │
 │             │              │             │             │
 │             │              │             │ Execute WS  │
 │             │              │             ├───────────> │
 │             │              │             │   WS done   │
 │             │              │             │<───────────┤
 │             │              │             │             │
 │             │              │             │ Execute WS  │
 │             │              │             ├───────────> │
 │             │              │             │   WS done   │
 │             │              │             │<───────────┤
 │             │              │             │             │
 │ Complete    │              │             │             │
 │<────────────┤              │             │             │
```

### Best Practices

1. **Start with @feature** - Even for experienced developers, the structured workflow helps
2. **Use checkpoints** - Orchestrator saves state, resume from interruptions
3. **Log decisions** - Use `sdp decisions log` for reproducibility
4. **Test incrementally** - Each workstream should have ≥80% coverage
5. **Keep workstreams small** - Split LARGE workstreams (>1500 LOC)
6. **Review at checkpoints** - Pause after critical workstreams
7. **Document as you go** - Don't leave documentation for last

### Troubleshooting

**Issue:** @feature asks too many questions
**Solution:** Use `--no-interview` flag to skip questions

**Issue:** Want to start from existing spec
**Solution:** Use `--spec docs/drafts/idea-{slug}.md` flag

**Issue:** Orchestrator failed mid-execution
**Solution:** Resume with `@oneshot F001 --resume F001`

**Issue:** Need to skip phases
**Solution:** Use `--vision-only` or jump directly to @design

**Issue:** Workstream blocked by dependency
**Solution:** Check `.oneshot/F001-checkpoint.md` for status
