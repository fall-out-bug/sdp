# SDP Skills Reference

Complete reference for SDP skill system and available skills.

---

## Table of Contents

- [Skill System](#skill-system)
- [Feature Skills](#feature-skills)
- [Utility Skills](#utility-skills)
- [Internal Skills](#internal-skills)
- [Skill Development](#skill-development)

---

## Skill System

### What are Skills?

Skills are Claude Code commands that execute specific SDP workflows. They are defined in `.claude/skills/{name}/SKILL.md` files.

### Skill Invocation

```bash
# Using @ prefix (user-facing)
@feature "Add authentication"

# Using / prefix (utilities)
/debug "Test fails"
```

### Skill Locations

**Feature Skills:**
- `.claude/skills/feature/`
- `.claude/skills/idea/`
- `.claude/skills/design/`
- `.claude/skills/build/`
- `.claude/skills/review/`
- `.claude/skills/deploy/`

**Utility Skills:**
- `.claude/skills/oneshot/`
- `.claude/skills/debug/`
- `.claude/skills/issue/`
- `.claude/skills/hotfix/`
- `.claude/skills/bugfix/`

**Internal Skills:**
- `.claude/skills/tdd/`

---

## Feature Skills

### @feature

**Location:** `.claude/skills/feature/SKILL.md`

**Purpose:** Unified entry point for feature development

**Workflow:**
1. Calls @idea for requirements
2. Generates PRODUCT_VISION.md
3. Creates Beads task
4. Outputs workstream plan

**Example:**
```bash
@feature "Add user authentication"
```

**Output:**
- `docs/drafts/idea-user-auth.md`
- Beads task F{ID}

---

### @idea

**Location:** `.claude/skills/idea/SKILL.md`

**Purpose:** Requirements gathering (internal)

**Process:**
- Deep questioning via AskUserQuestion
- Explores technical approaches
- Identifies tradeoffs
- Generates comprehensive spec

**Example:**
```bash
@idea "User authentication with OAuth"
```

**Called By:** @feature skill

---

### @design

**Location:** `.claude/skills/design/SKILL.md`

**Purpose:** Plan workstreams from requirements

**Process:**
1. Read requirements document
2. Enter Plan Mode for exploration
3. Decompose into workstreams
4. Create dependency graph
5. Request approval

**Example:**
```bash
@design idea-user-auth
```

**Output:**
- `docs/workstreams/backlog/WS-*.md`
- Dependency visualization

---

### @build

**Location:** `.claude/skills/build/SKILL.md`

**Purpose:** Execute workstream with TDD

**Process:**
1. Pre-build validation
2. Red: Write failing test
3. Green: Write minimal code
4. Refactor: Improve code
5. Quality gate checks
6. Git commit
7. Beads update

**Quality Gates:**
- Coverage ≥80%
- mypy --strict
- ruff clean
- Files <200 LOC
- No bare exceptions

**Example:**
```bash
@build WS-001-01
```

**Progress Tracking:** Real-time TodoWrite updates

---

### @review

**Location:** `.claude/skills/review/SKILL.md`

**Purpose:** Quality check for completed feature

**Checks:**
- All workstreams completed
- Tests passing
- Coverage ≥80%
- Type hints complete
- No TODO markers
- Clean architecture respected

**Example:**
```bash
@review <feature-id>
```

**Output:** Pass/fail verdict with details

---

### @deploy

**Location:** `.claude/skills/deploy/SKILL.md`

**Purpose:** Deploy feature to production

**Process:**
1. Final quality verification
2. Create git tag
3. Merge to main branch
4. Generate changelog
5. Trigger deployment

**Example:**
```bash
@deploy <feature-id>
```

**Prerequisites:**
- All workstreams completed
- @review passed
- Documentation updated

---

## Utility Skills

### @oneshot

**Location:** `.claude/skills/oneshot/SKILL.md`

**Purpose:** Autonomous feature execution

**Features:**
- Spawns orchestrator agent
- Executes all workstreams
- Handles dependencies
- Checkpoint save/restore
- Background execution
- Progress notifications

**Example:**
```bash
@oneshot <feature-id>

# Background mode
@oneshot <feature-id> --background

# Resume from checkpoint
@oneshot <feature-id> --resume <agent-id>
```

**Output:** Agent ID for resume capability

---

### /debug

**Location:** `.claude/skills/debug/SKILL.md`

**Purpose:** Systematic debugging using scientific method

**Process:**
1. Observe problem
2. Form hypothesis
3. Design experiment
4. Run experiment
5. Update hypothesis

**Example:**
```bash
/debug "Test fails when running full suite"
```

**Method:** Evidence-based root cause analysis

---

### @issue

**Location:** `.claude/skills/issue/SKILL.md`

**Purpose:** Bug routing and classification

**Process:**
1. Analyze bug description
2. Classify severity (P0/P1/P2/P3)
3. Route to appropriate fix command
   - P0 → @hotfix
   - P1/P2 → @bugfix
   - P3 → backlog

**Example:**
```bash
@issue "Login fails on Firefox with error 500"
```

**Severity Classification:**
- **P0** - Critical security, data loss, production down
- **P1** - Major functionality broken
- **P2** - Minor issues, workarounds available
- **P3** - Cosmetic, nice to have

---

### @hotfix

**Location:** `.claude/skills/hotfix/SKILL.md`

**Purpose:** Emergency fix for P0 issues

**Characteristics:**
- Branch from main
- Minimal changes only
- Deploy < 2 hours
- Skip full process
- Direct to production

**Example:**
```bash
@hotfix "Production database connection fails"
```

**Workflow:**
1. Create hotfix branch from main
2. Implement minimal fix
3. Fast verification
4. Deploy immediately
5. Create regular WS for follow-up

---

### @bugfix

**Location:** `.claude/skills/bugfix/SKILL.md`

**Purpose:** Quality fix for P1/P2 issues

**Characteristics:**
- Branch from feature/develop
- Full TDD cycle
- Quality gates enforced
- No production deploy

**Example:**
```bash
@bugfix "User profile image not loading"
```

**Workflow:**
1. Create bugfix branch
2. Full @build process
3. Quality verification
4. Merge to feature branch

---

## Internal Skills

### /tdd

**Location:** `.claude/skills/tdd/SKILL.md`

**Purpose:** TDD cycle enforcement (internal)

**Process:**
1. **Red** - Write failing test
2. **Green** - Write minimal code
3. **Refactor** - Improve code

**Called By:** @build skill (automatic)

**Not for direct user invocation**

---

## Skill Development

### Creating a New Skill

**Directory Structure:**
```
.claude/skills/{skill_name}/
├── SKILL.md          # Skill definition
├── prompt.md         # Optional: Additional prompts
└── examples/         # Optional: Usage examples
```

**SKILL.md Format:**
```markdown
# @skill-name

One-line description.

## Usage
\`\`\`bash
@skill-name "argument"
\`\`\`

## Process
1. Step one
2. Step two
3. Step three

## Output
What the skill produces

## Examples
Common usage examples
```

### Skill Best Practices

**DO ✅:**
- Use clear, descriptive names
- Provide usage examples
- Document all steps
- Handle errors gracefully
- Validate inputs

**DON'T ❌:**
- Don't create overlapping skills
- Don't skip error handling
- Don't omit documentation
- Don't make skills too complex

---

## Skill Invocation Flow

### Standard Feature Flow

```
@feature
  ↓ (calls)
@idea (gathers requirements)
  ↓ (outputs)
@design (plans workstreams)
  ↓ (outputs)
@build (executes each WS)
  ↓ (repeats for all WS)
@review (quality check)
  ↓ (if passed)
@deploy (production)
```

### Autonomous Flow

```
@feature
  ↓
@oneshot (spawns orchestrator)
  ↓
Orchestrator agent executes all workstreams
  ↓
@review + @deploy
```

### Bug Fix Flow

```
@issue
  ↓ (classifies)
@hotfix OR @bugfix
  ↓
@review (quality check)
```

---

## Quick Reference

| Skill | Purpose | Time | User-Facing |
|-------|---------|------|-------------|
| `@feature` | Create feature | 10-15 min | ✅ |
| `@idea` | Requirements | 5-10 min | ❌ Internal |
| `@design` | Plan workstreams | 5-10 min | ✅ |
| `@build` | Execute workstream | 30-90 min | ✅ |
| `@review` | Quality check | 5-10 min | ✅ |
| `@deploy` | Deploy feature | 5-10 min | ✅ |
| `@oneshot` | Autonomous exec | 2-6 hours | ✅ |
| `/debug` | Debug issue | 15-60 min | ✅ |
| `@issue` | Report bug | 2-5 min | ✅ |
| `@hotfix` | Emergency fix | < 2 hours | ✅ |
| `@bugfix` | Quality fix | 1-4 hours | ✅ |
| `/tdd` | TDD enforcement | Auto | ❌ Internal |

---

## See Also

- [commands.md](commands.md) - Command reference
- [quality-gates.md](quality-gates.md) - Quality standards
- [beginner/02-common-tasks.md](../beginner/02-common-tasks.md) - Common workflows

---

**Version:** SDP v0.9.8
**Updated:** 2026-02-26
