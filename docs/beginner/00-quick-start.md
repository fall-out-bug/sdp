# SDP Quick Start Tutorial

**Time:** 15 minutes
**Goal:** Create your first feature using SDP workflow
**Prerequisites:** Python 3.10+, Poetry, Git

---

## Step 1: Initialize Your Project (2 min)

Start by setting up SDP in your project:

```bash
# Navigate to your project directory
cd /path/to/your/project

# Run the interactive setup wizard
sdp init
```

**You'll be prompted for:**
- Project name (default: directory name)
- Description (default: "SDP project")
- Author (default: "Your Name")

**What gets created:**
```
docs/
├── workstreams/
│   ├── INDEX.md          # Workstream tracker
│   ├── TEMPLATE.md       # Workstream template
│   └── backlog/          # Planned workstreams
├── PROJECT_MAP.md        # Project decisions
└── drafts/               # Feature specifications
quality-gate.toml         # Quality rules
.env.template             # Environment variables
```

**Verify setup:**
```bash
sdp doctor
# Should show: ✓ All critical checks passed
```

---

## Step 2: Create Your First Feature (5 min)

Start with **@feature** - the unified entry point:

```bash
@feature "Add user authentication"
```

**What happens:**
1. **Vision Interview** - Claude asks about mission, users, success metrics
2. **Technical Interview** - Claude explores architecture, tradeoffs
3. **Documents Created:**
   - `PRODUCT_VISION.md` - Project manifesto
   - `docs/drafts/idea-user-auth.md` - Full specification
   - `docs/intent/user-auth.json` - Machine-readable intent

**Example conversation:**
```
Claude: What problem does this feature solve?
You: Users need to log in securely with email/password

Claude: Who are the primary users?
You: Web application users

Claude: Technical approach - sessions or JWT?
You: JWT tokens for stateless auth
```

**Power User Tip:** Skip interviews with flags:
```bash
@feature "Add user authentication" --no-interview
```

---

## Step 3: Design Workstreams (3 min)

Now break your feature into workstreams:

```bash
@design idea-user-auth
```

**What happens:**
1. Claude explores your codebase
2. Identifies integration points
3. Proposes workstream decomposition
4. Requests approval

**Example output:**
```
Feature F01: User Authentication
├── WS-F01-01: Domain entities (User, Credential)
├── WS-F01-02: Application services (AuthService)
├── WS-F01-03: Infrastructure layer (JWT, hashing)
└── WS-F01-04: API endpoints (/login, /register)

Approve this plan? (y/n)
```

**Files created:**
- `docs/workstreams/backlog/WS-F01-01.md`
- `docs/workstreams/backlog/WS-F01-02.md`
- `docs/workstreams/backlog/WS-F01-03.md`
- `docs/workstreams/backlog/WS-F01-04.md`

---

## Step 4: Execute Workstreams (3 min)

Implement each workstream with **@build**:

```bash
@build WS-F01-01
```

**What happens (TDD Cycle):**

### Phase 1: Red (Write Failing Test)
```
→ Writing failing test: tests/unit/test_user.py
→ Running pytest... FAILED (expected)
✓ Red phase complete
```

### Phase 2: Green (Implement Minimum Code)
```
→ Implementing User entity
→ Running pytest... PASSED
✓ Green phase complete
```

### Phase 3: Refactor (Improve Code)
```
→ Refactoring: extract validation logic
→ Running pytest... PASSED
✓ Refactor phase complete
```

### Phase 4: Verify
```
→ Checking acceptance criteria
✓ User entity created
✓ Email validation works
✓ Password hashing implemented

→ Running quality gates
✓ Coverage: 87%
✓ File size: 45 lines (<200)
✓ Type hints: 100%
✓ No architecture violations

✓ WS-F01-01 complete!
```

**Repeat for each workstream:**
```bash
@build WS-F01-02
@build WS-F01-03
@build WS-F01-04
```

**Or use autonomous mode:**
```bash
@oneshot F01
# Executes all workstreams automatically
```

---

## Step 5: Review Quality (1 min)

Verify your feature meets standards:

```bash
@review F01
```

**Checks performed:**
- ✓ All workstreams completed
- ✓ Acceptance criteria met
- ✓ Test coverage ≥80%
- ✓ Files <200 LOC
- ✓ Type hints complete
- ✓ No architecture violations
- ✓ No TODOs without followup

**Output:**
```
Review Result: APPROVED ✓

Workstreams: 4/4 completed
Coverage: 86% (target: 80%)
Quality Gates: PASSED

Recommendation: Ready for deployment
```

---

## Step 6: Deploy to Production (1 min)

Deploy your feature:

```bash
@deploy F01
```

**What happens:**
1. Merges feature branch to main
2. Creates git tag (e.g., v1.0.0)
3. Runs final validation
4. Generates release notes

**Output:**
```
✓ Merged dev → main
✓ Created tag v1.0.0
✓ Generated release notes

Feature F01 deployed successfully!
```

---

## Troubleshooting

### "sdp: command not found"
```bash
# Install SDP
pip install sdp

# Or with Poetry
poetry add sdp --group dev
```

### "sdp doctor shows failures"
```bash
# Install missing dependencies
curl -sSL https://install.python-poetry.org | python3 -

# Initialize git repo
git init

# Re-run setup
sdp init --force
```

### "@feature can't find my codebase"
```bash
# Ensure you're in project root
pwd  # Should show project directory

# Check for docs/ directory
ls docs/  # Should exist after sdp init
```

### "@build test fails"
```bash
# Run tests manually to see errors
pytest tests/unit/test_module.py -v

# Check coverage
pytest --cov=module --cov-report=term-missing

# Debug with
@debug "Test failure in test_user.py"
```

### "Workstream blocked by dependencies"
```bash
# Check INDEX.md for dependencies
cat docs/workstreams/INDEX.md

# Complete blocking workstreams first
@build WS-F01-01  # Must complete before WS-F01-02
```

---

## Next Steps

**Explore advanced workflows:**
- `@idea` - Deep requirements gathering
- `@design` - Interactive planning
- `@oneshot` - Autonomous multi-workstream execution
- `@debug` - Systematic debugging

**Read documentation:**
- [PROTOCOL.md](../PROTOCOL.md) - Full SDP specification
- [PRINCIPLES.md](../PRINCIPLES.md) - Core principles
- [CLAUDE_CODE.md](../guides/CLAUDE_CODE.md) - Claude Code integration

**Join the community:**
- GitHub: [github.com/your-org/sdp](https://github.com/your-org/sdp)
- Discord: [discord.gg/sdp](https://discord.gg/sdp)

---

## Checklist

Complete these steps to finish the tutorial:

- [ ] Run `sdp init` successfully
- [ ] Run `sdp doctor` - all checks pass
- [ ] Create feature with `@feature`
- [ ] Design workstreams with `@design`
- [ ] Execute workstream with `@build`
- [ ] Review feature with `@review`
- [ ] Deploy feature with `@deploy`

**Time taken:** _____ minutes

**Feedback?** Open an issue on GitHub!

---

**Version:** SDP 0.3.0
**Last Updated:** 2025-01-29
