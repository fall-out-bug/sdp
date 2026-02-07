---
name: prototype
description: Rapid prototyping shortcut for experienced vibecoders
tools: Read, Write, Bash, Glob, Grep, AskUserQuestion
version: 1.0.0
---

# @prototype - Rapid Prototyping Shortcut

Ultra-fast feature planning for experienced vibecoders who want to:
1. **15-minute interview** (5 questions max)
2. **Auto-generate monolithic workstreams** (1-3 big chunks)
3. **Launch agents immediately** (@oneshot with relaxed quality gates)
4. **Track tech debt** for later cleanup

> **Mode:** Working prototype over clean architecture
> **Speed:** 🚀 Maximum velocity, minimum bureaucracy
> **Debt:** All violations tracked as tech debt issues

## When to Use

- ✅ Experienced developers who know the codebase
- ✅ Need working prototype FAST (same day)
- ✅ Technical debt acceptable initially
- ✅ Willing to fix issues later (tech debt bomb)
- ❌ Production features requiring high quality
- ❌ Team collaboration (this is solo mode)

## Quick Reference

```bash
# Basic usage
@prototype "Add user authentication"

# With explicit feature ID
@prototype "Payment processing" --feature=F060

# Skip interview entirely
@prototype "Dashboard widgets" --skip-interview

# Specify workstream count
@prototype "API refactor" --workstreams=2
```

## Workflow

### Step 1: Ultra-Fast Interview (5 Questions)

Ask **only 5 critical questions** instead of @idea's 12-27:

```python
AskUserQuestion(
    questions=[
        {
            "question": "What problem does this feature solve?",
            "header": "Problem",
            "options": [
                {"label": "User pain point", "description": "Fixes existing user friction"},
                {"label": "New capability", "description": "Enables new user workflows"},
                {"label": "Technical debt", "description": "Improves code quality/performance"}
            ]
        },
        {
            "question": "What components are involved?",
            "header": "Scope",
            "options": [
                {"label": "Backend only", "description": "API, database, business logic"},
                {"label": "Frontend only", "description": "UI, user interactions"},
                {"label": "Full stack", "description": "Backend + frontend + integration"}
            ],
            "multiSelect": False
        },
        {
            "question": "Any external dependencies?",
            "header": "Deps",
            "options": [
                {"label": "None", "description": "Pure internal implementation"},
                {"label": "APIs", "description": "Third-party APIs or services"},
                {"label": "Database", "description": "New tables or collections"}
            ],
            "multiSelect": True
        },
        {
            "question": "Any blockers or risks?",
            "header": "Risks",
            "options": [
                {"label": "None known", "description": "Straightforward implementation"},
                {"label": "Technical uncertainty", "description": "New tech or unclear approach"},
                {"label": "Dependencies", "description": "Waiting on other teams or systems"}
            ],
            "multiSelect": True
        },
        {
            "question": "Define success in one sentence:",
            "header": "Success",
            "options": [
                {"label": "User can do X", "description": "Enable specific user workflow"},
                {"label": "Performance gain", "description": "Faster, more efficient"},
                {"label": "Bug fix", "description": "Resolve specific issue"}
            ],
            "multiSelect": False
        }
    ]
)
```

**Time Estimate:** 5-10 minutes (vs 30-60 minutes for @idea)

### Step 2: Auto-Generate Monolithic Workstreams

Based on interview answers, generate **1-3 big workstreams**:

```
Backend Only → 1 WS: "Backend Implementation"
Frontend Only → 1 WS: "Frontend Implementation"
Full Stack → 3 WS: "Backend", "Frontend", "Integration"
```

**Workstream Template:**

```markdown
# WS: 00-FFF-01 - [Component] Implementation

## Status
**Status:** pending
**Type:** prototype
**Priority:** P0 (blocks prototype)
**Mode:** relaxed-quality-gates

## Goal
[One-sentence goal from interview answer]

## Scope
### In Scope
- [Component] implementation
- Core functionality only
- Happy path (no error handling unless critical)

### Out of Scope
- Comprehensive error handling
- Edge cases
- Optimization
- Documentation (code comments only)

## Acceptance Criteria
1. [ ] Core feature works end-to-end
2. [ ] No obvious crashes
3. [ ] Basic manual testing successful

## Quality Gates (PROTOTYPE MODE)
- ❌ TDD: NOT ENFORCED (tests after code OK)
- ❌ Coverage: NO REQUIREMENT
- ❌ File size: NO LIMIT (< 200 LOC waived)
- ❌ Architecture: CLEAN ARCHITECTURE NOT REQUIRED
- ✅ Security: Basic checks only (no XSS, SQL injection)

## Technical Debt
All violations tracked in Beads as tech debt issues:
- Files > 200 LOC
- Missing test coverage
- Architecture violations
- Missing error handling

## Output
- Working code (prototype quality)
- Tech debt issues created
- NOT production-ready

## Next Step
After prototype: Fix tech debt or refactor as production feature
```

**Auto-Generation Logic:**

```python
def generate_prototype_workstreams(interview_answers):
    components = []

    if interview_answers["scope"] == "Backend only":
        components = ["Backend"]
    elif interview_answers["scope"] == "Frontend only":
        components = ["Frontend"]
    else:  # Full stack
        components = ["Backend", "Frontend", "Integration"]

    workstreams = []
    for i, component in enumerate(components, start=1):
        ws_id = f"00-{feature_id}-{i:02d}"
        workstreams.append({
            "id": ws_id,
            "title": f"{component} Implementation",
            "goal": interview_answers["success"],
            "mode": "prototype",
            "quality_gates": "relaxed"
        })

    return workstreams
```

### Step 3: Create Workstream Files

Write workstream files to `docs/workstreams/backlog/`:

```bash
# Example: Feature F060, 3 workstreams
docs/workstreams/backlog/00-060-01-backend.md
docs/workstreams/backlog/00-060-02-frontend.md
docs/workstreams/backlog/00-060-03-integration.md
```

### Step 4: Launch @oneshot with Relaxed Quality Gates

Execute all workstreams with **prototype mode settings**:

```python
Skill(
    skill="oneshot",
    args={
        "feature_id": feature_id,
        "mode": "prototype",  # KEY FLAG
        "quality_gates": "relaxed"
    }
)
```

**@oneshot Behavior in Prototype Mode:**

| Setting | Normal Mode | Prototype Mode |
|---------|------------|----------------|
| TDD Discipline | Red→Green→Refactor | Tests after code (optional) |
| Coverage | ≥80% required | No requirement |
| File Size | <200 LOC | No limit |
| Architecture | Clean layers | Monolithic OK |
| Documentation | Full | Comments only |
| Security | Full review | Basic checks |
| Tech Debt | Fail gates | Create issues, continue |

### Step 5: Tech Debt Bomb (Post-Execution)

After all workstreams complete, generate tech debt summary:

```bash
⚠️ PROTOTYPE MODE COMPLETE

Feature: {feature_id}
Workstreams: {count} completed
Duration: {time}

Tech Debt Created: {X} issues
├─ Architecture violations: {count}
├─ Missing tests: {count}
├─ Files >200 LOC: {count}
└─ Security warnings: {count}

⚠️ NOT PRODUCTION-READY

Next Steps:
1. Manual testing of prototype
2. Fix tech debt issues
3. Add test coverage
4. Refactor for production
5. Re-run @review with full quality gates
```

**Beads Integration (if enabled):**

```bash
# Create tech debt issue for each violation
bd create \
  --title="TechDebt: Prototype F{feature_id} - {violation}" \
  --type=tech-debt \
  --priority=3 \
  --notes="Generated by @prototype in relaxed-quality-gates mode"
```

## Command Flags

```bash
@prototype <feature_description> [options]

Options:
  --feature=FFF     # Explicit feature ID (default: auto-generate)
  --workstreams=N  # Number of workstreams (1-3, default: auto-detect)
  --skip-interview # Skip questions, use defaults
  --no-debt        # Don't create tech debt issues
  --immediate      # Launch @oneshot without confirmation
```

## Examples

### Example 1: Backend API Feature

```bash
User: @prototype "Add telemetry API endpoint"

Claude:
→ Step 1: Interview (5 questions)
   Q1: What problem? → "New capability"
   Q2: Components? → "Backend only"
   Q3: Dependencies? → "Database"
   Q4: Risks? → "None known"
   Q5: Success? → "API returns 200"

→ Step 2: Generate 1 workstream
   docs/workstreams/backlog/00-001-01-backend.md

→ Step 3: Confirm launch
   Launch @oneshot F001 in prototype mode? (Y/n)

→ Step 4: Execute
   ✅ 00-001-01: Backend Implementation (45m)
   ⚠️ Created 3 tech debt issues

→ Step 5: Complete
   ⚠️ PROTOTYPE MODE COMPLETE
   📊 Tech Debt Created: 3 issues
   🚀 Next: Manual testing → Fix debt → Production
```

### Example 2: Full-Stack Feature

```bash
User: @prototype "User dashboard with widgets"

Claude:
→ Step 1: Interview (5 questions)
   Q2: Components? → "Full stack"
   (其他问题...)

→ Step 2: Generate 3 workstreams
   00-002-01: Backend Implementation
   00-002-02: Frontend Implementation
   00-002-03: Integration

→ Step 3-6: Execute all 3 in parallel
   ✅ All workstreams complete (2h 15m total)
   ⚠️ Created 12 tech debt issues
```

### Example 3: Skip Interview, Go Direct

```bash
User: @prototype "Refactor user service" --skip-interview --immediate

Claude:
→ Skipping interview, using defaults
→ Generated 1 workstream: 00-003-01
→ Launching @oneshot F003 immediately...
→ (No confirmation, runs directly)
```

## Quality Gate Overrides

In prototype mode, these quality gates are **DISABLED**:

| Gate | Normal | Prototype | Rationale |
|------|--------|-----------|-----------|
| **TDD** | Red→Green→Refactor | Tests after code | Speed over discipline |
| **Coverage** | ≥80% required | No requirement | Prototype doesn't need full coverage |
| **File Size** | <200 LOC | No limit | Monolithic code OK for prototype |
| **Architecture** | Clean layers | Monolithic OK | Speed over separation |
| **Documentation** | Full docs | Comments only | Self-documenting code sufficient |
| **Error Handling** | Comprehensive | Minimal | Happy path only |

**STRICTLY ENFORCED (Non-Negotiable):**
- ✅ **Code MUST compile and run** (no build failures)
- ✅ **No crashes during execution** (stable runtime)
- ✅ **100% functional** (feature works end-to-end)
- ✅ **Basic security** (no SQL injection, XSS, path traversal)
- ✅ **Git commits succeed** (all code committed)
- ✅ **Manual testing passed** (user verified it works)

**⚠️ CRITICAL: Despite relaxed gates, the result MUST be 100% working.**

This is **prototype speed, NOT prototype quality**. The difference:
- **Fast:** Skip TDD discipline, skip coverage requirements, skip architecture rules
- **NOT broken:** Code still compiles, runs, and works correctly

**Failed prototypes are NOT acceptable.** If it doesn't work, it's not a prototype—it's broken code.

## Tech Debt Tracking

All violations are automatically tracked:

```python
violations = []

# During @oneshot execution
for file in modified_files:
    if count_loc(file) > 200:
        violations.append({
            "type": "architecture",
            "severity": "tech_debt",
            "file": file,
            "loc": count_loc(file),
            "issue": "File exceeds 200 LOC guideline"
        })

if coverage < 0.8:
    violations.append({
        "type": "testing",
        "severity": "tech_debt",
        "package": package_name,
        "coverage": coverage,
        "issue": f"Coverage {coverage:.1%} below 80% threshold"
    })

# Create Beads issues for each violation
if beads_enabled:
    for v in violations:
        bd_create(
            title=f"TechDebt: {v['type']} - {v['issue']}",
            type="tech-debt",
            priority=3,  # Low priority, fix later
            notes=f"Prototype mode violation in {feature_id}"
        )
```

## Output Artifacts

**Created:**
```
docs/drafts/prototype-{feature_id}.md     # Interview summary
docs/workstreams/backlog/00-FFF-*.md      # 1-3 workstreams
```

**Not Created:**
- No detailed spec (use @idea for that)
- No architecture diagrams (use @design for that)
- No comprehensive documentation (prototype only)

## Comparison: @prototype vs @feature vs @idea

| Aspect | @prototype | @feature | @idea |
|--------|------------|----------|-------|
| **Interview** | 5 questions (15 min) | 12-27 questions (30-60 min) | 12-27 questions |
| **Workstreams** | Auto-generated 1-3 | Designed 5-30 | Designed 5-30 |
| **Quality Gates** | Relaxed (prototype) | Strict (production) | Strict (production) |
| **Tech Debt** | Tracked & OK | Tracked & fail | Tracked & fail |
| **Speed** | Same day | Days to weeks | Days to weeks |
| **Use Case** | Working prototype | Production feature | Requirements only |
| **Follow-up** | Fix debt or refactor | Production-ready | Need @design |

## When NOT to Use @prototype

❌ **Production features** (use @feature instead)
❌ **Team projects** (no collaboration, no review)
❌ **Security-critical code** (skip gates is dangerous)
❌ **Long-term maintenance** (tech debt accumulates)
❌ **Learning the codebase** (use @idea for exploration)

---

## ⚠️ CRITICAL DISTINCTION: Speed vs Quality

**What @prototype DOES:**
- ✅ Skip TDD discipline (write tests after code)
- ✅ Skip coverage requirements (no 80% threshold)
- ✅ Skip architecture rules (monolithic files OK)
- ✅ Skip documentation (comments only)

**What @prototype DOES NOT DO:**
- ❌ Accept broken code (MUST compile and run)
- ❌ Accept crashes (MUST be stable)
- ❌ Accept non-functional features (MUST work end-to-end)
- ❌ Skip manual testing (MUST verify it works)

**The Promise:**
> "Despite relaxed quality gates, the result will be 100% working. It may have technical debt, ugly code, or missing tests—but it WILL function correctly."

**The Trade-off:**
- **Speed:** Same-day prototype vs weeks of proper development
- **Debt:** 5-15 hours of tech debt cleanup vs proper development from start
- **Risk:** Prototype might need refactor vs production-ready code

**When to accept this trade-off:**
- ✅ Proof of concept (does this idea work?)
- ✅ Internal tool (temporary usage)
- ✅ MVP for user testing (get feedback fast)
- ✅ Time-sensitive opportunity (beat competitors to market)

**When NOT to accept this trade-off:**
- ❌ Production system (users depend on it)
- ❌ Security-critical (handling payments, personal data)
- ❌ Complex domain (prototype will be throwaway work)

## Follow-up After Prototype

After @prototype completes, you have three paths:

### Path 1: Fix Tech Debt (Make it Production)

```bash
# List tech debt issues
bd list --type=tech-debt --feature=F{feature_id}

# Fix issues systematically
for issue in $(bd list --type=tech-debt); do
    @build $issue  # Fix each debt item
done

# Re-run review with full gates
@review F{feature_id}
```

### Path 2: Refactor from Scratch

```bash
# Prototype worked, now redo properly
@feature "{feature_description}" --based-on=F{feature_id}
# Creates proper workstreams with full quality gates
```

### Path 3: Discard Prototype

```bash
# Prototype didn't work, try different approach
@prototype "{different_description}"
# Or delete and start over
```

## Beads Integration

**If Beads enabled** (`bd --version` works):

```python
# After interview, create feature issue
feature_id = bd create(
    title=feature_description,
    type="prototype",
    priority=1,  # Medium priority
    notes="Prototype mode: relaxed quality gates"
)

# Create workstream issues
for ws in workstreams:
    ws_id = bd create(
        title=f"WS: {ws['title']}",
        type="task",
        priority=2,
        parent=feature_id
    )
    # Track mapping
    write_sdp_beads_mapping(ws_id, f"00-{feature_id}-{i:02d}")
```

**During @oneshot:**
- Mark workstream as in_progress
- Create tech debt issues as violations detected
- Mark workstream as complete
- Update parent feature status

## Implementation Notes

**For @oneshot modification:**

Add `--mode=prototype` flag to @oneshot:

```python
# In @oneshot skill
def execute_oneshot(feature_id, mode="production"):
    if mode == "prototype":
        # Relaxed quality gates
        quality_gates = {
            "tdd": "optional",
            "coverage": "none",
            "file_size": "no_limit",
            "architecture": "relaxed"
        }
        tech_debt_tracking = True
    else:
        # Normal strict gates
        quality_gates = {
            "tdd": "required",
            "coverage": "80%",
            "file_size": "200",
            "architecture": "strict"
        }
        tech_debt_tracking = False
```

## Version

**1.0.0** - Initial release
- Ultra-fast interview (5 questions)
- Auto-generate 1-3 monolithic workstreams
- Relaxed quality gates for prototype mode
- Tech debt tracking
- Immediate @oneshot launch

---

**See Also:**
- `.claude/skills/idea/SKILL.md` — Full requirements gathering
- `.claude/skills/feature/SKILL.md` — Full feature planning
- `.claude/skills/oneshot/SKILL.md` — Autonomous execution
- `.claude/skills/build/SKILL.md` — Workstream execution with TDD
