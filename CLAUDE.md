# Claude Code Integration Guide

Quick reference for using this Spec-Driven Protocol (SDP) repository with Claude Code.

> **📝 Meta-note:** This guide was written with AI assistance (Claude Sonnet 4.5). The workflow is based on real development experience.

## TL;DR

> **🎯 New to SDP?** Start with [docs/NAVIGATION.md](docs/NAVIGATION.md) - Single entry point for all documentation.

### Quick Start (4 Commands)

```bash
@vision "AI-powered task manager"     # Strategic planning
@reality --quick                      # What's actually there?
@feature "Add user authentication"    # Plan feature
@build 00-001-01                      # Execute workstream
@review F01                           # Quality check
```

> **⚠️ Migration Notice:** Upgrading from a previous version? See [Breaking Changes Guide](docs/migrations/breaking-changes.md).

**📋 Decision Trees:** See [docs/NAVIGATION.md](docs/NAVIGATION.md#decision-trees) to choose the right workflow.

**⚠️ Workstream ID Format:** Use `PP-FFF-SS` (e.g., `00-001-01`), NOT legacy `WS-FFF-SS`

## Decision Tree: @vision → @reality → @feature → @oneshot

### Four-Level Planning Model

**SDP has four orchestrators for different planning levels:**

```
Strategic Level                 Analysis Level                 Feature Level                Execution Level
┌──────────────────┐           ┌──────────────────┐         ┌──────────────────┐         ┌──────────────────┐
│     @vision      │           │    @reality      │         │    @feature      │         │    @oneshot      │
│  (strategic)     │           │  (codebase anal) │         │   (planning)     │         │   (execution)    │
└──────────────────┘           └──────────────────┘         └──────────────────┘         └──────────────────┘
         │                              │                             │                             │
         ▼                              ▼                             ▼                             ▼
  7 Expert Agents             8 Expert Agents              @idea + @design            @build (all WS)
  (product analysis)           (codebase analysis)          (requirements + WS)         (implement)
         │                              │                             │                             │
         ▼                              ▼                             ▼                             ▼
  Product Artifacts            Reality Report               workstreams                 @review + @deploy
  (VISION, PRD, ROADMAP)      (health, gaps, debt)         (00-FFF-SS.md)              (quality + merge)
```

### Level Comparison

| Aspect | @vision | @reality | @feature | @oneshot |
|--------|---------|----------|----------|----------|
| **Purpose** | Strategic product planning | Codebase analysis | Feature planning (requirements + workstreams) | Execute workstreams |
| **Input** | Product idea ("AI task manager") | Project directory | Feature description ("Add OAuth") | Feature ID (F01) or WS list |
| **Output** | PRODUCT_VISION.md, PRD.md, ROADMAP.md | Reality report (health, gaps, debt) | Workstream files (00-FFF-SS.md) | Implemented code + deployed feature |
| **Duration** | Quarterly/annual review | Per project or quarterly | Per feature | Per feature |
| **Agents** | 7 experts (product, market, technical, UX, business, growth, risk) | 8 experts (architecture, quality, testing, security, performance, docs, debt, standards) | @idea + @design | @build + @review + @deploy |
| **When to Use** | New project, major pivot, quarterly planning | New to project, before @feature, after @vision, quarterly review | Starting new feature from scratch | Workstreams exist, ready to implement |
| **Human Interaction** | Medium (3-5 interview cycles) | Minimal (auto-scans project) | Heavy (AskUserQuestion, ExitPlanMode) | Minimal (only critical blockers) |

### When to Use Each Level

**Use @vision when:**
- ✅ Starting a new project or product
- ✅ Quarterly strategic review
- ✅ Major pivot or direction change
- ✅ Need comprehensive product analysis
- ✅ Want expert analysis across 7 dimensions (product, market, technical, UX, business, growth, risk)

**Use @reality when:**
- ✅ New to project (what's actually here?)
- ✅ Before @feature (what can we build on?)
- ✅ After @vision (how do docs match code?)
- ✅ Quarterly review (track tech debt and quality trends)
- ✅ Debugging mysteries (why doesn't this work?)
- ✅ Want 8-expert codebase analysis (architecture, quality, testing, security, performance, docs, debt, standards)

**Use @feature when:**
- ✅ You have a feature idea but no workstreams
- ✅ You need to explore requirements (@idea)
- ✅ You need to design architecture (@design)
- ✅ You want interactive planning (questions, tradeoffs)
- ✅ Product vision already exists

**Use @oneshot when:**
- ✅ Workstreams already exist (from @feature or @design)
- ✅ You want autonomous execution (no human interaction)
- ✅ You have 5-30 workstreams to execute
- ✅ You want checkpoint/resume capability

### Typical Full Flow

```bash
# Step 1: Strategic planning (quarterly or new project)
@vision "AI-powered task manager for remote teams"
# → 7 expert agents analyze product, market, technical, UX, business, growth, risk
# → Generates: PRODUCT_VISION.md, PRD.md, ROADMAP.md
# → Result: Clear strategic direction

# Step 2: Feature planning (per feature)
@feature "Add OAuth2 authentication"
# → @idea gathers requirements
# → @design creates workstreams
# → Result: 00-001-01.md, 00-001-02.md, ... in docs/workstreams/backlog/

# Step 3: Execution (autonomous)
@oneshot F01
# → @build executes all workstreams
# → @review checks quality
# → @deploy merges to main
# → Result: Feature shipped
```

### Skip @vision if:
- Product vision already exists (PRODUCT_VISION.md present)
- Working on existing product (not new project)
- Incremental feature (not major pivot)

### Skip @feature if:
- Workstreams already exist (from previous @design)
- You created workstreams manually
- You just want to execute existing WS

### Skip @oneshot if:
- Only 1-2 workstreams (use @build directly)
- You want manual control over each WS
- You're learning the system (use @build to understand workflow)

## Decision Tree: @feature vs @oneshot (Legacy)

> **Note:** This section preserved for historical context. See above for updated three-level model.

### Mental Model

**SDP has two independent orchestrators for different phases:**

```
Planning Phase                  Execution Phase
┌──────────────────┐           ┌──────────────────┐
│   @feature       │           │   @oneshot       │
│  (planning)      │           │  (execution)     │
└──────────────────┘           └──────────────────┘
         │                              │
         ▼                              ▼
    @idea ──────────────────────▶ @build (all WS)
    (gather requirements)              (implement)
         │                              │
         ▼                              ▼
    @design ──────────────────────▶ @review
    (create workstreams)              (quality check)
         │                              │
         ▼                              ▼
   workstreams                     @deploy
   (00-FFF-SS.md)                  (merge to main)
```

### Key Differences

| Aspect | @feature | @oneshot |
|--------|----------|----------|
| **Purpose** | Planning: gather requirements + design workstreams | Execution: implement all workstreams |
| **Input** | Feature description ("Add X") | Feature ID (F01) or list of WS |
| **Output** | Workstream files (00-FFF-SS.md) | Implemented code + deployed feature |
| **Phases** | Discovery (@idea) → Design (@design) | Build (@build) → Review (@review) → Deploy (@deploy) |
| **When to Use** | Starting new feature from scratch | Workstreams exist, ready to implement |
| **Human Interaction** | Heavy (AskUserQuestion, ExitPlanMode) | Minimal (only for critical blockers) |

### When to Use Which

**Use @feature when:**
- ✅ You have a feature idea but no workstreams
- ✅ You need to explore requirements (@idea)
- ✅ You need to design architecture (@design)
- ✅ You want interactive planning (questions, tradeoffs)

**Use @oneshot when:**
- ✅ Workstreams already exist (from @feature or @design)
- ✅ You want autonomous execution (no human interaction)
- ✅ You have 5-30 workstreams to execute
- ✅ You want checkpoint/resume capability

**Typical Flow:**
```bash
# Day 1: Planning phase (interactive)
@feature "Add payment processing"
# → @idea gathers requirements
# → @design creates workstreams
# → Result: 00-050-01.md, 00-050-02.md, ... in docs/workstreams/backlog/

# Day 2-5: Execution phase (autonomous)
@oneshot F050
# → @build executes all workstreams
# → @review checks quality
# → @deploy merges to main
# → Result: Feature shipped
```

**Skip @feature if:**
- Workstreams already exist (from previous @design)
- You created workstreams manually
- You just want to execute existing WS

**Skip @oneshot if:**
- Only 1-2 workstreams (use @build directly)
- You want manual control over each WS
- You're learning the system (use @build to understand workflow)

## Available Skills

| Skill | Purpose | Phase | Example |
|-------|---------|-------|---------|
| `@vision` | **Strategic product planning** (7 expert agents) | Strategic | `@vision "AI-powered task manager"` |
| `@reality` | **Codebase analysis** (8 expert agents) | Analysis | `@reality --quick` or `@reality --focus=security` |
| `@feature` | **Planning orchestrator** (interactive) | Planning | `@feature "Add payment processing"` |
| `@idea` | **Requirements gathering** (AskUserQuestion) | Planning | `@idea "Add payment processing"` |
| `@design` | **Workstream design** (EnterPlanMode) | Planning | `@design idea-payments` |
| `@oneshot` | **Execution orchestrator** (autonomous) | Execution | `@oneshot F01` or `@oneshot F01 --background` |
| `@build` | Execute single workstream (TDD) | Execution | `@build 00-001-01` |
| `@review` | Multi-agent quality review | Execution | `@review F01` |
| `@deploy` | Merge feature branch to main | Execution | `@deploy F01` |
| `/debug` | **Systematic debugging** (scientific method) | Debug | `/debug "Test fails unexpectedly"` |
| `@issue` | Debug and route bugs | Debug | `@issue "Login fails on Firefox"` |
| `@hotfix` | Emergency fix (P0) | Debug | `@hotfix "Critical API outage"` |
| `@bugfix` | Quality fix (P1/P2) | Debug | `@bugfix "Incorrect totals"` |

**Internal skills** (not called directly by users):
| Skill | Purpose | Called By |
|-------|---------|----------|
| `/tdd` | TDD cycle enforcement (Red→Green→Refactor) | `@build` (automatic) |

Skills are defined in `.claude/skills/{name}/SKILL.md`

**Claude Code Integration Highlights:**
- `@vision` — Strategic planning: 7 expert agents (product, market, technical, UX, business, growth, risk) → PRODUCT_VISION.md, PRD.md, ROADMAP.md
- `@reality` — Codebase analysis: 8 expert agents (architecture, quality, testing, security, performance, docs, debt, standards) → Reality report + vision gap analysis
- `@feature` — Planning orchestrator: @idea (requirements) → @design (workstreams)
- `@idea` — Deep interviewing via AskUserQuestion (no obvious questions, explores tradeoffs)
- `@design` — EnterPlanMode for codebase exploration + AskUserQuestion for architecture decisions
- `@oneshot` — Execution orchestrator: @build (all WS) → @review → @deploy
- `@build` — TodoWrite real-time progress tracking through TDD cycle
- `@review` — Multi-agent quality check (QA + Security + DevOps + SRE + TechLead + Documentation)

## Quick Reference

### First Time Setup

1. **Read core docs:**
   - [README.md](README.md) — Overview and quick start
   - [PROTOCOL.md](docs/PROTOCOL.md) — Full SDP specification
   - [RULES_COMMON.md](RULES_COMMON.md) — Common rules

2. **Understand key concepts:**
   - **Workstream (WS)**: Atomic task, one-shot execution
   - **Feature**: 5-30 workstreams
   - **Release**: 10-30 features

3. **Review quality gates:**
   - Files < 200 LOC
   - Coverage ≥80%
   - No `except: pass`
   - Full type hints

4. **Install Beads CLI** (for task tracking):
   ```bash
   # macOS
   brew tap beads-dev/tap
   brew install beads

   # Linux
   curl -sSL https://raw.githubusercontent.com/beads-dev/beads/main/install.sh | bash

   # Verify
   bd --version
   ```

### Typical Workflow

**Option A: Full Strategic Planning (recommended for new projects)**

```bash
# 1. Strategic phase: @vision (quarterly or new project)
@vision "AI-powered task manager for remote teams"
# → 7 expert agents analyze across all dimensions
# → Generates: PRODUCT_VISION.md, PRD.md, ROADMAP.md
# Result: Clear strategic direction and feature roadmap

# 2. Analysis phase: @reality (bridge vision to reality)
@reality --quick
# → 8 expert agents analyze codebase
# → Generates: Reality report (health, gaps, tech debt)
# → Compares: Vision vs Reality gap analysis
# Result: Clear understanding of current state

# 3. Planning phase: @feature (per feature)
@feature "User can reset password via email"
# → @idea gathers requirements (AskUserQuestion)
# → @design creates workstreams (ExitPlanMode)
# Result: 00-050-01.md, 00-050-02.md, ... in docs/workstreams/backlog/

# 4. Execution phase: @oneshot (autonomous)
@oneshot F050
# → @build executes all workstreams
# → @review checks quality
# → @deploy merges to main
# Result: Feature shipped
```

**Option B: Interactive Feature Planning (skip @vision)**

```bash
# 1. Planning phase: @feature (combines @idea + @design)
@feature "User can reset password via email"
# → @idea gathers requirements (AskUserQuestion)
# → @design creates workstreams (ExitPlanMode)
# Result: 00-050-01.md, 00-050-02.md, ... in docs/workstreams/backlog/

# 2. Execution phase: @oneshot (autonomous)
@oneshot F050
# → @build executes all workstreams
# → @review checks quality
# → @deploy merges to main
# Result: Feature shipped
```

**Option B: Manual Planning (skip @feature)**

```bash
# 1. Create workstreams manually
# Edit docs/workstreams/backlog/00-050-01.md, etc.

# 2. Execute manually one by one
@build 00-050-01
@build 00-050-02
# ...

# 3. Review and deploy
@review F050
@deploy F050
```

**Option C: Hybrid (plan interactively, execute manually)**

```bash
# 1. Use @feature for planning
@feature "Add payment processing"
# → Creates workstreams

# 2. Execute manually (for learning or debugging)
@build 00-050-01  # Execute first WS manually
@build 00-050-02  # Execute second WS manually
# ...

# 3. Use @oneshot for remaining WS
@oneshot F050  # Continues from checkpoint
```

### Progress Tracking

When using `@build`, Claude Code automatically tracks progress using TodoWrite:

```markdown
User: @build 00-060-01

Claude:
→ Creating todo list...
  ✓ [in_progress] Pre-build validation
  • [pending] Write failing test (Red)
  • [pending] Implement minimum code (Green)
  • [pending] Refactor implementation
  • [pending] Verify Acceptance Criteria
  • [pending] Run quality gates
  • [pending] Append execution report
  • [pending] Git commit

→ Reading WS file...
  ✓ [completed] Pre-build validation
  ✓ [in_progress] Write failing test (Red)
  • [pending] Implement minimum code (Green)
  ...

→ Test created, running pytest... FAILED (expected)
  ✓ [completed] Write failing test (Red)
  ✓ [in_progress] Implement minimum code (Green)
  ...

→ Implementation done, running pytest... PASSED
  ✓ [completed] Implement minimum code (Green)
  ✓ [in_progress] Refactor implementation
  ...

[All steps complete]
  ✓ All tasks completed
```

This provides real-time visibility into WS execution progress.

### Autonomous Execution with @oneshot

For features with multiple workstreams, use `@oneshot` for autonomous execution:

```markdown
User: @oneshot F01

Claude Code:
→ Spawning orchestrator agent via Task tool...
→ Agent ID: abc123xyz (save for resume)

Orchestrator Agent:
→ Reading feature specification and workstreams...
→ Found 4 workstreams to execute

→ Creating todo list...
  ✓ [in_progress] Wait for PR approval
 • [pending] Execute 00-001-01: Domain entities
 • [pending] Execute 00-001-02: Application services
 • [pending] Execute 00-001-03: Infrastructure layer
 • [pending] Execute 00-001-04: API endpoints
  • [pending] Run final review
  • [pending] Generate UAT guide

→ Creating PR for approval...
→ Waiting for approval...

[User approves PR in GitHub]

→ PR approved!
  ✓ [completed] Wait for PR approval
 ✓ [in_progress] Execute 00-001-01: Domain entities

→ Executing 00-001-01...
 (internal @build TodoWrite tracking for this WS)
→ 00-001-01 complete (45m, 85% coverage)
 ✓ [completed] Execute 00-001-01: Domain entities
 ✓ [in_progress] Execute 00-001-02: Application services

→ Executing 00-001-02...
→ 00-001-02 complete (1h 10m, 82% coverage)
 ✓ [completed] Execute 00-001-02: Application services
 ✓ [in_progress] Execute 00-001-03: Infrastructure layer

... (continues for all WS)

→ All workstreams complete
  ✓ [in_progress] Run final review

→ Running /review F01...
→ Review verdict: APPROVED
  ✓ [completed] Run final review
  ✓ [in_progress] Generate UAT guide

→ Generating UAT guide...
  ✓ [completed] Generate UAT guide

→ All tasks completed! ✅

Main Claude:
## ✅ Feature F01 Execution Complete

**Agent ID:** abc123xyz (for resume)
**Duration:** 3h 45m
**Workstreams:** 4/4 completed
**Coverage:** avg 86%

### Next Steps
1. Human UAT (5-10 min)
2. `@deploy F01` if UAT passes
```

**Background execution** for large features:

```bash
User: @oneshot F01 --background

Claude Code:
→ Starting orchestrator agent in background...
→ Task ID: xyz789
→ Output file: /tmp/agent_xyz789.log

You can continue working. I'll notify when complete.
Check progress: Read("/tmp/agent_xyz789.log")

[5 minutes later]
✅ Background task xyz789 completed!
Feature F01 is done and ready for UAT.
```

**Resume** from interruption:

```bash
# If execution interrupted
User: @oneshot F01 --resume abc123xyz

Claude Code:
→ Resuming agent abc123xyz...
→ Agent continues from last checkpoint (WS-001-03)
```

### File Structure Reference

```
project/
├── PRODUCT_VISION.md      # Project manifesto (generated by @feature)
├── docs/
│   ├── schema/            # Intent JSON schema
│   ├── intent/            # Machine-readable intent files
│   ├── drafts/            # @idea outputs here
│   ├── workstreams/
│   │   ├── backlog/       # @design outputs here
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
│   │   ├── feature/       # Unified entry point
│   │   ├── idea/          # Requirements gathering
│   │   ├── design/        # Workstream planning
│   │   ├── build/         # WS execution
│   │   ├── tdd/           # TDD discipline
│   │   ├── debug/         # Systematic debugging
│   │   └── oneshot/       # Autonomous execution
│   ├── agents/            # Multi-agent mode (advanced)
│   └── settings.json      # Claude Code settings
└── hooks/                 # Git hooks for validation
```

## Key Principles (Quick)

- **SOLID, DRY, KISS, YAGNI** — see [docs/PRINCIPLES.md](docs/PRINCIPLES.md)
- **Clean Architecture** — Domain ← App ← Infra ← Presentation
- **TDD** — Tests first (Red → Green → Refactor)
- **AI-Readiness** — Small files, low complexity, typed

## Validation

### Pre-build Check
```bash
hooks/pre-build.sh WS-001-01
```

### Post-build Check
```bash
hooks/post-build.sh WS-001-01 project.module
```

### Manual Validation
```bash
python scripts/validate.py docs/workstreams/backlog/
```

## Quality Gates (Enforced)

| Gate | Requirement |
|------|-------------|
| **AI-Readiness** | Files < 200 LOC, CC < 10, type hints |
| **Clean Architecture** | No layer violations |
| **Error Handling** | No `except: pass` |
| **Test Coverage** | ≥80% |
| **No TODOs** | All tasks completed or new WS |

## Forbidden Patterns

❌ `except: pass` or bare exceptions  
❌ Time-based estimates  
❌ Layer violations  
❌ Files > 200 LOC  
❌ TODO without followup WS  
❌ Coverage < 80%

## Required Patterns

✅ Type hints everywhere  
✅ Tests first (TDD)  
✅ Explicit error handling  
✅ Clean architecture boundaries  
✅ Conventional commits

## Troubleshooting

### Skill not found
Check `.claude/skills/{name}/SKILL.md` exists

### Validation fails
Run `hooks/pre-build.sh {WS-ID}` to see specific issues

### Workstream blocked
Check dependencies in `docs/workstreams/backlog/{WS-ID}.md`

### Coverage too low
Run `pytest --cov --cov-report=term-missing` to identify gaps

### Legacy Workstream ID Format

**Problem:** Workstreams using old `WS-FFF-SS` format instead of `PP-FFF-SS`

**Solution:** Use the migration script

```bash
# Preview changes (safe)
python scripts/migrate_workstream_ids.py --dry-run

# Migrate SDP workstreams
python scripts/migrate_workstream_ids.py --project-id 00

# Migrate other projects
python scripts/migrate_workstream_ids.py --project-id 02 --path ../hw_checker
```

**What it does:**
- Updates `ws_id` in frontmatter (`WS-001-01` → `00-001-01`)
- Adds `project_id` field
- Renames files to match new format
- Updates cross-WS dependencies
- Validates all changes

**See also:** `docs/migration/ws-naming-migration.md`

## Advanced: Multi-Agent Mode

For complex features, use multi-agent orchestration:

```bash
@orchestrator F01  # Coordinates all agents
```

Agents defined in `.claude/agents/`:
- `planner.md` — Breaks features into workstreams
- `builder.md` — Executes workstreams
- `reviewer.md` — Quality checks
- `deployer.md` — Production deployment
- `orchestrator.md` — Coordinates workflow

## Configuration

See `.claude/settings.json` for:
- Custom Git hooks
- Validation scripts
- Tool integrations

## Resources

| Resource | Purpose |
|----------|---------|
| [PROTOCOL.md](docs/PROTOCOL.md) | Full specification |
| [docs/PRINCIPLES.md](docs/PRINCIPLES.md) | Core principles |
| [docs/SLOS.md](docs/SLOS.md) | SLOs/SLIs for CLI tool |
| [CODE_PATTERNS.md](docs/reference/CODE_PATTERNS.md) | Code patterns |
| [MODELS.md](docs/reference/MODELS.md) | Model recommendations |
| [prompts/commands/](prompts/commands/) | Skill instructions |

---

**Version:** SDP 0.3.0  
**Claude Code Version:** 0.3+  
**Mode:** Skill-based, one-shot execution

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

## Reality-First Development

**Principle:** Always verify actual code before following documentation.

Based on analysis of 827 sessions, the #1 friction point is **documentation-code mismatch**. Workstream descriptions often don't match actual implementation.

### Quick Reality Check (90 seconds)

Before modifying any file based on documentation:

```bash
/reality-check <filename>
```

**Example:**
```markdown
User: Add validation to User model in models.py

Claude: Let me reality-check first...
[Reads src/sdp/quality/models.py]
⚠️ Reality Check: models.py contains dataclasses, NOT validation logic
Recommendation: Create separate validators.py instead
```

### Full Workstream Verification (5-10 minutes)

Before executing workstreams:

```bash
/verify-workstream 00-001-01
```

**Validates:**
- All scope_files exist
- Functions/classes in docs actually present in code
- File purpose matches documentation
- Architectural layers correct

**Output:**
```markdown
## Documentation vs Reality Analysis

| File | Docs | Reality | Status |
|------|------|---------|--------|
| validators.py | Generic validation | Business logic | ❌ Mismatch |
| models.py | Validation models | Dataclasses | ❌ Wrong layer |

**Recommendation:** PAUSE - Update workstream to reflect reality
```

### Integration with Workflow

Add to step 3 of Typical Workflow:

```bash
# 3. Verify workstream (Reality-First)
/verify-workstream 00-001-01

# 4. Execute workstream
@build 00-001-01
```

**Success Metrics:**
- Prevents "wrong_approach" friction (13% of sessions)
- Reduces pragmatic adaptation overhead
- Maintains architectural integrity

