# @feature Quick Start Tutorial

**Learn to build features with SDP in 15 minutes**

> **Time:** 15 minutes
> **Difficulty:** Beginner
> **Prerequisites:** Claude Code installed

---

## Table of Contents

1. [Installation](#installation)
2. [Your First Feature](#your-first-feature)
3. [Understanding the Workflow](#understanding-the-workflow)
4. [Executing Workstreams](#executing-workstreams)
5. [Review and Deploy](#review-and-deploy)
6. [Next Steps](#next-steps)

---

## Installation

### Step 1: Install SDP Plugin (1 minute)

```bash
# Clone plugin
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp

# Copy skills to Claude Code
cp -r ~/.claude/sdp/prompts/* .claude/

# Verify installation
ls .claude/skills/
# You should see: feature.md, design.md, build.md, etc.
```

### Step 2: Initialize Project (1 minute)

```bash
# Navigate to your project
cd my-project

# Create SDP directories
mkdir -p docs/{drafts,intent,workstreams/{backlog,in_progress,completed}}

# Verify structure
tree docs/
# docs/
# ├── drafts/
# ├── intent/
# └── workstreams/
#     ├── backlog/
#     ├── in_progress/
#     └── completed/
```

### Step 3: (Optional) Install Beads CLI (1 minute)

```bash
# macOS
brew tap beads-dev/tap
brew install beads

# Linux
curl -sSL https://raw.githubusercontent.com/beads-dev/beads/main/install.sh | bash

# Verify
bd --version
# Output: beads 0.49.3
```

**✅ Installation complete!** Let's build your first feature.

---

## Your First Feature

We'll build a simple **user login** feature together.

### Step 1: Start @feature (30 seconds)

```bash
@feature "Add user login"
```

**What happens:** Claude asks you 3-5 questions about your vision.

---

### Step 2: Vision Interview (3 minutes)

Claude guides you through a structured interview:

```
🎯 Mission: What problem does user login solve?

You: Users need to authenticate to access their personalized
     content and protect their data.

Claude: 👥 Users: Who are the primary users?

You: End users who need to log in, and administrators who
     manage accounts.

Claude: 📊 Success Metrics: How will you measure success?

You: Login success rate >95%, account creation <2 minutes,
     zero security breaches.
```

**What happens:**
- Claude captures your vision
- Generates `PRODUCT_VISION.md`
- Logs strategic tradeoffs

**Output:** `PRODUCT_VISION.md` at project root

```markdown
# PRODUCT_VISION.md

> **Last Updated:** 2026-02-06

## Mission
Enable user authentication for personalized content and data protection.

## Users
1. End Users - Log in to access features
2. Administrators - Manage user accounts

## Success Metrics
- Login success rate >95%
- Account creation <2 minutes
- Zero security breaches
```

---

### Step 3: Technical Interview (4 minutes)

Claude asks technical questions:

```
🔧 Authentication Method: Session-based, JWT, or OAuth?

You: JWT tokens with refresh token rotation.

💾 Storage: Where should we store user credentials?

You: PostgreSQL with bcrypt password hashing.

⚠️  Failure Modes: What if the database is down?

You: Return 503 Service Unavailable, queue for retry.

🔗 Integration Points: Existing user management?

You: No, building from scratch.

🔒 Security: Password hashing requirements?

You: bcrypt with cost factor 12.

📊 Scalability: Expected concurrent users?

You: 1000 concurrent users at peak.

✅ Compliance: Any regulatory requirements?

You: GDPR compliance required.
```

**What happens:**
- Claude captures technical decisions
- Logs each decision with `sdp decisions log`
- Validates tradeoffs

**Output:** Decisions logged to Beads (if installed)

---

### Step 4: Generate Intent (30 seconds)

Claude creates machine-readable intent:

**Output:** `docs/intent/user-login.json`

```json
{
  "problem": "Users need to authenticate to access personalized content",
  "users": ["end_users", "administrators"],
  "success_criteria": [
    "login_success_rate > 95%",
    "account_creation < 2 minutes",
    "zero_security_breaches"
  ],
  "technical_approach": {
    "auth_method": "jwt",
    "storage": "postgresql",
    "security": "bcrypt_cost_12"
  },
  "constraints": ["GDPR_compliant"]
}
```

---

### Step 5: Requirements Draft (1 minute)

Claude creates detailed specification:

**Output:** `docs/drafts/idea-user-login.md`

```markdown
# User Login Feature

> **Feature ID:** F001
> **Status:** Draft

## Problem
Users need to authenticate to access personalized content and data protection.

## Users
1. End Users - Log in to access features
2. Administrators - Manage user accounts

## Success Criteria
- [ ] Login success rate >95%
- [ ] Account creation <2 minutes
- [ ] Zero security breaches

## Goals
1. JWT-based authentication
2. Secure password storage (bcrypt)
3. Refresh token rotation
4. Session management

## Technical Approach
- Auth: JWT tokens with refresh rotation
- Storage: PostgreSQL + bcrypt
- Security: Rate limiting, account lockout

## Architecture
```
┌─────────┐    ┌─────────┐    ┌─────────┐
│ Client  │───→│ AuthAPI │───→│Postgres │
└─────────┘    └─────────┘    └─────────┘
```
```

---

### Step 6: Workstream Planning (2 minutes)

Claude breaks feature into workstreams:

```bash
@design idea-user-login
```

**What happens:**
- Claude analyzes your codebase
- Identifies dependencies
- Creates workstream breakdown

**Output:** Workstreams in `docs/workstreams/backlog/`

```
F001: User Login (8 workstreams)
├── WS-001: Domain Models (SMALL)
├── WS-002: Auth Service (SMALL)
├── WS-003: JWT Management (SMALL)
├── WS-004: Password Hashing (SMALL)
├── WS-005: API Endpoints (SMALL)
├── WS-006: Session Management (SMALL)
├── WS-007: Integration Tests (MEDIUM)
└── WS-008: Documentation (SMALL)
```

**Example Workstream:** `docs/workstreams/backlog/WS-001_domain_models.md`

```markdown
# WS-001: Domain Models

> **Feature:** F001: User Login
> **Size:** SMALL

## Goal
Define User and Session entities with Clean Architecture.

## Acceptance Criteria
- User entity: email, password_hash, created_at
- Session entity: user_id, token, expires_at
- Repository interfaces defined
- Unit tests with ≥80% coverage

## Scope Files
- src/auth/domain/user.go
- src/auth/domain/session.go
```

---

### Step 7: Orchestrator Execution (3 minutes)

Claude executes workstreams autonomously:

```bash
@oneshot F001
```

**What happens:**
```
[10:00] Starting feature execution: F001
[10:00] Loading workstreams...
[10:00] Execution order: [WS-001 WS-002 WS-003 WS-004 WS-005 WS-006 WS-007 WS-008]

[10:01] Executing WS-001: Domain Models (1/8)...
[10:05] → WS-001 complete (4m, 92% coverage)

[10:05] Executing WS-002: Auth Service (2/8)...
[10:09] → WS-002 complete (4m, 88% coverage)

[10:09] Executing WS-003: JWT Management (3/8)...
[10:12] → WS-003 complete (3m, 85% coverage)

[10:12] Executing WS-004: Password Hashing (4/8)...
[10:15] → WS-004 complete (3m, 91% coverage)

[10:15] Executing WS-005: API Endpoints (5/8)...
[10:18] → WS-005 complete (3m, 87% coverage)

[10:18] Executing WS-006: Session Management (6/8)...
[10:21] → WS-006 complete (3m, 89% coverage)

[10:21] Executing WS-007: Integration Tests (7/8)...
[10:27] → WS-007 complete (6m, 94% coverage)

[10:27] Executing WS-008: Documentation (8/8)...
[10:30] → WS-008 complete (3m, docs complete)

[10:30] ✅ Feature complete: 8/8 workstreams, 30m total, 90% avg coverage
```

**What you get:**
- ✅ Working code with tests
- ✅ Documentation
- ✅ ≥80% coverage
- ✅ All quality gates passed

---

## Understanding the Workflow

### The 7 Phases

```
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Phase 1  │→ │ Phase 2  │→ │ Phase 3  │→ │ Phase 4  │
│  Vision  │  │ PRODUCT_ │  │Technical │  │  Intent  │
│Interview │  │ VISION.md│  │Interview │  │   JSON   │
└──────────┘  └──────────┘  └──────────┘  └──────────┘
                                               ↓
                                        ┌──────────┐
                                        │ Phase 5  │
                                        │  Draft   │
                                        └──────────┘
                                               ↓
                                        ┌──────────┐
                                        │ Phase 6  │
                                        │ @design  │
                                        └──────────┘
                                               ↓
                                        ┌──────────┐
                                        │ Phase 7  │
                                        │@oneshot  │
                                        │(optional)│
                                        └──────────┘
```

### Phase Duration

| Phase | Time | Output |
|-------|------|--------|
| 1. Vision Interview | 3 min | User requirements |
| 2. Generate Vision | 30 sec | `PRODUCT_VISION.md` |
| 3. Technical Interview | 4 min | Technical decisions |
| 4. Generate Intent | 30 sec | `docs/intent/*.json` |
| 5. Requirements Draft | 1 min | `docs/drafts/*.md` |
| 6. Workstream Planning | 2 min | Workstream breakdown |
| 7. Execution | 20-30 min | Working code |
| **Total** | **30-40 min** | **Complete feature** |

---

## Executing Workstreams

### Option 1: Autonomous (Recommended)

```bash
@oneshot F001
```

**Pros:**
- Hands-free execution
- Automatic checkpointing
- Resume from interruptions

**Cons:**
- Less control over individual workstreams

### Option 2: Manual

```bash
@build WS-001
@build WS-002
@build WS-003
# ... etc
```

**Pros:**
- Full control
- Review between workstreams

**Cons:**
- More manual effort
- Slower overall

### Option 3: Hybrid

```bash
# Execute first 3 manually
@build WS-001
@build WS-002
@build WS-003

# Then autonomous for rest
@oneshot F001 --resume F001
```

**Best of both worlds!**

---

## Checkpoints & Resume

Orchestrator saves checkpoints automatically:

**Checkpoint:** `.oneshot/F001-checkpoint.json`

```json
{
  "id": "F001",
  "status": "in_progress",
  "completed_workstreams": ["WS-001", "WS-002", "WS-003"],
  "current_workstream": "WS-004",
  "created_at": "2026-02-06T10:00:00Z",
  "updated_at": "2026-02-06T10:12:00Z"
}
```

**Resume after interruption:**

```bash
# Power outage? No problem!
@oneshot F001 --resume F001
# Continues from WS-004
```

---

## Review and Deploy

### Quality Check (Automatic)

Every workstream passes quality gates:

```bash
# Run during @build
go test ./...                    # Tests pass
go tool cover -func=coverage.out # Coverage ≥80%
go vet ./...                     # No warnings
```

**If quality gate fails:**
```
❌ Coverage: 72% (required: ≥80%)
📝 Missing tests: src/auth/user.go:45-52

→ Auto-fix: Adding tests...
→ ✅ Coverage: 84%
```

### Review Feature

```bash
@review F001
```

**Output:**
```
✅ APPROVED

Feature: F001: User Login
Workstreams: 8/8 complete
Coverage: 90% avg
Quality Gates: PASS

Recommendations:
- Consider adding rate limiting (future WS)
- Document API endpoints in Swagger
```

### Deploy to Production

```bash
@deploy F001
```

**What happens:**
1. Runs final tests
2. Creates git tag
3. Pushes to remote
4. Triggers deployment (if configured)

---

## Progressive Menu System

Skip phases with flags:

```bash
# Vision only (stop before technical interview)
@feature "Add analytics" --vision-only

# From existing spec
@feature --spec docs/drafts/idea-payments.md

# No interview (use defaults)
@feature "Add notifications" --no-interview

# Execute immediately
@feature "Add reporting" --execute
```

---

## Common Patterns

### Pattern 1: Quick Feature

```bash
# 5-minute feature (skip interviews)
@feature "Add logging" --no-interview --execute
```

### Pattern 2: Complex Feature

```bash
# Full workflow with manual review
@feature "Add payments"
# ... phases 1-6 ...
@build WS-001  # Review first workstream
@oneshot F001  # Execute rest
```

### Pattern 3: Iterative Development

```bash
# Start with vision
@feature "Add search" --vision-only

# Later, add technical details
@feature --spec docs/drafts/idea-search.md

# Execute when ready
@oneshot F003
```

---

## Troubleshooting

### Issue: Too many questions

**Solution:** Use `--no-interview` flag

```bash
@feature "Add cache" --no-interview
```

### Issue: Want to start from existing spec

**Solution:** Use `--spec` flag

```bash
@feature --spec docs/drafts/idea-auth.md
```

### Issue: Orchestrator failed

**Solution:** Resume from checkpoint

```bash
@oneshot F001 --resume F001
```

### Issue: Need to stop mid-execution

**Solution:** Just press Ctrl+C, checkpoint is saved

```bash
@oneshot F001
[10:15] Executing WS-004...
^C
# Checkpoint saved, resume later
```

---

## Next Steps

### Learn More

- [Full Protocol Documentation](PROTOCOL.md)
- [Workflow Decision Guide](workflow-decision.md)
- [Quality Gates Reference](reference/quality-gates.md)

### Advanced Features

- Multi-agent coordination
- Role-based routing
- Approval gates
- Telegram notifications

### Contribute

- GitHub: https://github.com/ai-masters/sdp-plugin
- Issues: https://github.com/ai-masters/sdp-plugin/issues

---

## Checklist

Use this checklist for your first feature:

- [ ] Install SDP plugin
- [ ] Initialize project structure
- [ ] Run `@feature "Add user login"`
- [ ] Answer vision interview questions (3 min)
- [ ] Review `PRODUCT_VISION.md`
- [ ] Answer technical interview questions (4 min)
- [ ] Review `docs/drafts/idea-user-login.md`
- [ ] Run `@design idea-user-login`
- [ ] Review workstream breakdown
- [ ] Run `@oneshot F001`
- [ ] Review generated code
- [ ] Run `@review F001`
- [ ] Run `@deploy F001` (optional)

**🎉 Congratulations!** You've built your first feature with SDP.

---

## Time Breakdown

| Step | Time | Cumulative |
|------|------|------------|
| Installation | 3 min | 3 min |
| Vision Interview | 3 min | 6 min |
| Technical Interview | 4 min | 10 min |
| Workstream Planning | 2 min | 12 min |
| Execution (8 workstreams) | 30 min | 42 min |
| Review & Deploy | 3 min | **45 min** |

**First feature:** ~45 minutes
**Subsequent features:** ~15-20 minutes (familiarity)

---

**✅ Tutorial complete!** You're ready to build features with SDP.
