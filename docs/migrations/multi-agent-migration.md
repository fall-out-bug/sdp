# Multi-Agent SDP Migration Guide

Migrate from old SDP (single-agent) to new SDP (multi-agent with progressive disclosure).

## Overview

SDP v4.0 introduces a **four-level planning model** with multi-agent synthesis, progressive disclosure, and verbosity tiers.

**Key Changes:**
- **Four planning levels**: @vision → @reality → @feature → @oneshot
- **Progressive disclosure**: 3-question cycles with trigger points (reduced from 25+ to 12-27 questions)
- **Verbosity tiers**: --quiet, --verbose, --debug flags for all skills
- **Multi-agent synthesis**: Parallel expert agents for complex decisions
- **Beads integration**: Git-backed issue tracking

## Before/After Comparison

### Old Workflow (v3.x)
```
@feature "Add OAuth"
→ 15-25 questions asked upfront
→ No structure, no breaks
→ User fatigue
→ Single agent does everything
```

### New Workflow (v4.x)
```
# Strategic level (new)
@vision "product idea"      # 7 expert agents (product, market, technical, UX, business, growth, risk)
@reality --quick           # 8 expert agents (architecture, quality, testing, security, performance, docs, debt, standards)

# Feature level (updated)
@idea "Add OAuth"          # 3-question cycles with trigger points
@design task-id            # Discovery blocks (skip irrelevant)
@oneshot F001              # Parallel execution with checkpoints

# Verbosity control (new)
@build 00-001-01 --quiet   # Exit status only
@build 00-001-01 --verbose # Step-by-step progress
@build 00-001-01 --debug   # Internal state + API calls
```

## Four-Level Planning Model

| Level | Orchestrator | Purpose | Output | Duration |
|-------|-------------|---------|--------|----------|
| **Strategic** | @vision (7 agents) | Product planning | VISION, PRD, ROADMAP | Quarterly/new project |
| **Analysis** | @reality (8 agents) | Codebase analysis | Reality report | Per project/quarterly |
| **Feature** | @feature (@idea + @design) | Requirements + WS | Workstreams | Per feature |
| **Execution** | @oneshot (@build) | Parallel execution | Implemented code | Per feature |

### When to Use Each Level

**Use @vision when:**
- Starting new project or product
- Quarterly strategic review
- Major pivot or direction change
- Need comprehensive product analysis (7 dimensions)

**Use @reality when:**
- New to project (what's actually here?)
- Before @feature (what can we build on?)
- Quarterly review (track tech debt)
- Want 8-expert codebase analysis

**Use @feature when:**
- You have a feature idea but no workstreams
- Need interactive planning (@idea → @design)
- Want progressive discovery blocks

**Use @oneshot when:**
- Workstreams exist, ready to execute
- Want autonomous execution (no human interaction)
- Have 5-30 workstreams to execute

## Migration Steps

### Step 1: Update Skills
```bash
# Skills are already updated in .claude/skills/
# @vision: v1.0.0 (7 expert agents)
# @reality: v1.0.0 (8 expert agents)
# @idea: v4.0.0 (progressive disclosure)
# @design: v4.0.0 (progressive disclosure)
# @build: v1.0.0 (verbosity tiers)
# @review: v7.0.0 (multi-agent + contract validation)
```

### Step 2: Learn New Commands
```bash
# Strategic level (new)
@vision "product idea"           # 7 expert agents
@reality --quick                 # 8 expert agents, 5-10 min
@reality --deep                  # 8 expert agents, 30-60 min
@reality --focus=security        # Single expert deep dive

# Feature level (updated)
@feature "Add OAuth"             # Orchestrator: @idea → @design
@idea "Add OAuth"                # Progressive disclosure
@design task-id                  # Discovery blocks

# Execution level (same)
@build ws-id                     # Execute workstream with TDD
@oneshot feature-id              # Execute all workstreams in parallel
@review feature-id               # Multi-agent quality review
@deploy feature-id               # Deploy to production

# Verbosity control (new)
<command> --quiet                # Exit status only (✅/❌)
<command> --verbose              # Step-by-step progress
<command> --debug                # Internal state + API calls
```

### Step 3: Update Workflows

**Old Workflow:**
```bash
@feature "Add OAuth"       # One big questionnaire (25+ questions)
@build 00-001-01
@build 00-001-02
...
@review F001
@deploy F001
```

**New Workflow (Option A - Full Strategic Planning):**
```bash
# Step 1: Strategic phase (quarterly or new project)
@vision "AI-powered task manager for remote teams"
# → 7 expert agents analyze across all dimensions
# → Generates: PRODUCT_VISION.md, PRD.md, ROADMAP.md

# Step 2: Analysis phase (bridge vision to reality)
@reality --quick
# → 8 expert agents analyze codebase
# → Generates: Reality report (health, gaps, tech debt)

# Step 3: Planning phase (per feature)
@feature "User can reset password via email"
# → @idea gathers requirements (18 questions average vs 25+)
# → @design creates workstreams (5 discovery blocks)

# Step 4: Execution phase (autonomous)
@oneshot F050
# → @build executes all workstreams in parallel (4.96x speedup)
# → @review checks quality (6 agents)
# → @deploy merges to main
```

**New Workflow (Option B - Interactive Feature Planning):**
```bash
# Skip @vision if product vision already exists
@feature "User can reset password via email"
# → @idea gathers requirements (3-question cycles)
# → @design creates workstreams (discovery blocks)
# → Result: 00-050-01.md, 00-050-02.md, ... in docs/workstreams/backlog/

@oneshot F050
# → Autonomous execution
# → Result: Feature shipped
```

**New Workflow (Option C - Manual Planning):**
```bash
# Create workstreams manually
# Edit docs/workstreams/backlog/00-050-01.md, etc.

@build 00-050-01
@build 00-050-02
...

@review F050
@deploy F050
```

### Step 4: Initialize Beads (Optional)

```bash
# Install Beads CLI
brew tap beads-dev/tap
brew install beads

# Initialize in project
bd init

# Check ready tasks
bd ready

# Track work
bd create --title="Implement feature X" --type=feature
bd update sdp-xxx --status in_progress
bd close sdp-xxx
bd sync
```

## Pattern Mapping

| Old Pattern | New Pattern | Benefit |
|------------|------------|---------|
| `@feature` | `@idea` + `@design` | Separated concerns, progressive disclosure |
| Deep dive upfront | 3-question cycles with triggers | Reduced fatigue (18 vs 25 questions) |
| Unbounded questions | 12-27 questions target | Predictable duration |
| Single agent | Multi-agent synthesis | Better quality, parallel experts |
| Sequential @build | `@oneshot` parallel execution | 4.96x speedup |
| No output control | Verbosity tiers (--quiet/--verbose/--debug) | User control |
| Manual task tracking | Beads integration | Multi-session persistence |

## Verbosity Tiers

All skills now support four verbosity levels:

```bash
# Level 0 (--quiet): Exit status only
@build 00-050-01 --quiet
# Output: ✅

# Level 1 (default): Summary with metrics
@build 00-050-01
# Output: ✅ 00-050-01: Workstream Parser (22m, 85%, commit:abc123)

# Level 2 (--verbose): Step-by-step progress
@build 00-050-01 --verbose
# Output:
# → Activating guard...
# → Reading WS spec...
# → TDD cycle: Red → Green → Refactor
# → Quality check: PASS
# ✅ COMPLETE

# Level 3 (--debug): Internal state + API calls
@build 00-050-01 --debug
# Output:
# [DEBUG] WS ID: 00-050-01
# [DEBUG] Beads ID: sdp-abc123
# → Activating guard...
# [DEBUG] Guard activated: /tmp/guard-00-050-01.json
# → Reading WS spec...
# [DEBUG] Scope files: [src/sdp/parser.py, tests/sdp/test_parser.py]
# ...
# ✅ COMPLETE
```

**Supported skills:** @build, @review, @feature, @vision, @reality, @deploy, @hotfix, @bugfix

## Progressive Disclosure

### @idea (Requirements Gathering)

**Old pattern:**
- 15-25 questions asked upfront
- No breaks, user fatigue
- Unbounded depth

**New pattern:**
- 3-question cycles with trigger points
- User controls depth
- Target: 12-27 questions (average 18)

**Trigger points after each cycle:**
- Continue with targeted questions (recommended)
- Deep design (skip to @design with detailed spec)
- Skip to @design immediately
- Use --quiet mode (5 questions core only)

**Example:**
```python
# Cycle 1 (3 questions)
→ What problem does this feature solve?
→ Who are your target users?
→ What defines success?

# Trigger: Continue / Deep design / Skip / --quiet?

# Cycle 2 (3 questions)
→ What are the core requirements?
→ What are the edge cases?
→ What are the constraints?

# Trigger: Continue / Deep design / Skip?

# Continue until requirements clear (typically 4-6 cycles)
```

### @design (Workstream Design)

**Old pattern:**
- Single exploration phase
- No structure
- Easy to get lost

**New pattern:**
- 3-5 discovery blocks
- 3 questions per block
- Skip irrelevant blocks

**Discovery blocks:**
1. Architecture (component boundaries, patterns)
2. Technology stack (language, framework, libraries)
3. Data modeling (schema, relationships)
4. Integration points (APIs, services)
5. Quality attributes (performance, security)

**Example:**
```python
# Block 1: Architecture (3 questions)
→ What are the main components?
→ How do components interact?
→ What's the component boundary?

# Trigger: Continue / Skip block / Done?

# Block 2: Technology (3 questions)
→ Language preference?
→ Framework selection?
→ Library dependencies?

# Trigger: Continue / Skip block / Done?

# Continue until architecture clear (typically 3-5 blocks)
```

## Multi-Agent Synthesis

### @vision (7 Expert Agents)

**Parallel experts:**
1. Product expert - Business value, user needs
2. Market expert - Competitive landscape, positioning
3. Technical expert - Feasibility, architecture
4. UX expert - User experience, workflows
5. Business expert - Revenue model, monetization
6. Growth expert - Acquisition, retention, virality
7. Risk expert - Technical, market, execution risks

**Synthesis:**
- Unanimous agreement (all 7 agree)
- Domain expertise (highest confidence wins)
- Quality gate (best quality score)
- Merge (combine best parts)
- Escalate (ask human)

**Output:**
- PRODUCT_VISION.md - Project manifesto
- PRD.md - Product requirements document
- ROADMAP.md - Feature roadmap

### @reality (8 Expert Agents)

**Parallel experts:**
1. Architecture - Layers, patterns, violations
2. Code Quality - LOC, complexity, duplication
3. Testing - Coverage, quality, infrastructure
4. Security - Secrets, OWASP, dependencies
5. Performance - Bottlenecks, memory, scalability
6. Documentation - Coverage, drift, quality
7. Technical Debt - TODO/FIXME, smells, design debt
8. Standards - Conventions, compliance, best practices

**Output:**
- Health score (0-100)
- Critical issues (fix now)
- Quick wins (fix today)
- Trends (if --review mode)

### @review (7 Agents + Contract Validation)

**Parallel agents:**
1. QA - Coverage, test quality, metrics
2. Security - Threats, vulnerabilities, compliance
3. DevOps - CI/CD, infrastructure, deployment
4. SRE - SLOs, monitoring, incidents
5. TechLead - Code quality, architecture
6. Documentation - Drift (3 levels: Vision→Specs→Code)
7. Contract Validator - API contract validation

**Verdict:**
- APPROVED (all 7 PASS)
- CHANGES_REQUESTED (any FAIL)

**No middle ground.**

## Beads Integration

### Task Tracking

```bash
# Create issue
bd create --title="Implement feature X" --type=feature --priority=2

# Start work
bd ready                    # Find unblocked tasks
bd update sdp-xxx --status in_progress

# Complete work
bd close sdp-xxx --reason "Feature implemented"

# Sync to git
bd sync                     # Commit changes to .beads/

# Dependencies
bd dep add sdp-yyy sdp-xxx  # sdp-yyy depends on sdp-xxx
bd blocked                  # Show blocked tasks
```

### @build Integration

```bash
# @build automatically updates Beads if enabled
@build 00-001-01
# → bd update sdp-xxx --status in_progress
# → TDD cycle...
# → bd close sdp-xxx --reason "WS completed"
# → bd sync
```

## Rollback Procedure

If issues arise:
```bash
# Revert to old skills
git checkout v3.1.0 -- .claude/skills/idea/SKILL.md
git checkout v3.1.0 -- .claude/skills/design/SKILL.md
git checkout v3.1.0 -- .claude/skills/feature/SKILL.md

# Use old commands
@feature "feature description"
```

## Testing Checklist

### @vision
- [ ] Try @vision with product idea
- [ ] Verify 7 expert agents spawned in parallel
- [ ] Check artifacts generated (VISION, PRD, ROADMAP)
- [ ] Test verbosity tiers (--quiet, --verbose, --debug)

### @reality
- [ ] Try @reality --quick (5-10 min)
- [ ] Try @reality --deep (30-60 min)
- [ ] Try @reality --focus=security
- [ ] Check health score calculated
- [ ] Test verbosity tiers

### @idea
- [ ] Try @idea with progressive disclosure
- [ ] Verify 3-question cycles
- [ ] Check trigger points work
- [ ] Verify question count (12-27 range)
- [ ] Test --quiet mode (5 questions)

### @design
- [ ] Try @design with discovery blocks
- [ ] Verify 3 questions per block
- [ ] Check skip block works
- [ ] Test verbosity tiers

### @oneshot
- [ ] Try @oneshot with 5+ workstreams
- [ ] Verify parallel execution (4.96x speedup)
- [ ] Check checkpoint/resume works
- [ ] Test background execution

### @build
- [ ] Try @build with ws-id
- [ ] Try @build with beads-id
- [ ] Verify TDD cycle (Red→Green→Refactor)
- [ ] Check verbosity tiers

### @review
- [ ] Try @review with feature-id
- [ ] Verify 6 agents spawned in parallel
- [ ] Check contract validation
- [ ] Verify verdict (APPROVED/CHANGES_REQUESTED)
- [ ] Test verbosity tiers

## Common Patterns

### Pattern 1: Quick Feature
**Old:** `@feature "add button"` (25 questions)
**New:** `@idea "add button" --quiet` (5 questions)

### Pattern 2: Complex Feature
**Old:** `@feature "OAuth integration"` (30+ questions)
**New:**
```bash
@idea "OAuth integration"      # 18 questions average
@design task-id                # Skip irrelevant blocks
@oneshot F001                  # Parallel execution (4.96x speedup)
```

### Pattern 3: New Product
**Old:** Manual planning
**New:**
```bash
@vision "AI task manager"      # 7 expert agents
@reality --quick               # 8 expert agents
@feature "first feature"       # Plan based on analysis
@oneshot F001                  # Execute
```

### Pattern 4: Quarterly Review
**New:**
```bash
@vision --review               # Update vision based on progress
@reality --review              # Track tech debt trends
# Compare: Vision vs Reality gap
# Plan: Next quarter features
```

### Pattern 5: Debugging Issues
**New:**
```bash
@reality --focus=security      # Deep dive on security
@issue "Login fails"           # Classify and route
# → /hotfix (P0), /bugfix (P1/P2), or backlog
```

## Performance Improvements

| Metric | Old | New | Improvement |
|--------|-----|-----|-------------|
| Questions per feature | 25+ | 18 (avg) | 28% reduction |
| @oneshot speed | 1x | 4.96x | 396% faster |
| Expert agents | 1 | 7-8 | Better quality |
| Output control | None | 4 levels | User control |
| Task tracking | Manual | Beads | Multi-session |

## Breaking Changes

1. **@feature refactored** - Now orchestrates @idea + @design (not a direct replacement)
2. **@idea separated** - Progressive disclosure (not deep dive upfront)
3. **@design updated** - Discovery blocks (not single exploration)
4. **@oneshot new** - Parallel execution with checkpoints (new orchestrator)
5. **@review expanded** - 7 agents + contract validation (was fewer)

## Resources

- [Verbosity Tiers Spec](../reference/verbosity-tiers.md)
- [Deep-Thinking Integration](../design/deep-thinking-integration.md)
- [Progressive Disclosure](../feature/progressive-disclosure.md)
- [Beads Workflow](https://github.com/beads-dev/beads)

---

**Version:** 4.0.0
**Migrating From:** 3.x
**Date:** 2026-02-08
**Authors:** SDP Team + Claude Sonnet 4.5
