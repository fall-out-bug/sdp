# SDP Commands Reference

Complete reference for all SDP CLI commands and skills.

---

## Table of Contents

- [Feature Commands](#feature-commands)
- [Utility Commands](#utility-commands)
- [Internal Commands](#internal-commands)
- [Command Options](#command-options)

---

## Feature Commands

### @feature

**Purpose:** Unified entry point for feature development

**Usage:**
```bash
@feature "Feature description"
```

**What it does:**
1. Interactive requirements gathering
2. Creates PRODUCT_VISION.md
3. Generates workstream plan
4. Sets up Beads task tracking

**Output:**
- `docs/drafts/idea-{slug}.md` - Requirements document
- Beads task with feature ID

**Example:**
```bash
@feature "Add user authentication"
```

**See:** [.claude/skills/feature/SKILL.md](../../.claude/skills/feature/SKILL.md)

---

### @design

**Purpose:** Plan workstreams from requirements

**Usage:**
```bash
@design idea-{slug}
```

**What it does:**
1. Reads requirements document
2. Decomposes into workstreams
3. Creates dependency graph
4. Generates workstream files

**Output:**
- `docs/workstreams/backlog/WS-*.md` - Workstream files
- Dependency graph visualization

**Example:**
```bash
@design idea-user-auth
```

**See:** [.claude/skills/design/SKILL.md](../../.claude/skills/design/SKILL.md)

---

### @build

**Purpose:** Execute a single workstream

**Usage:**
```bash
@build WS-{ID}
```

**What it does:**
1. Pre-build validation
2. TDD cycle (Red → Green → Refactor)
3. Quality gate checks
4. Git commit
5. Beads status update

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

**See:** [.claude/skills/build/SKILL.md](../../.claude/skills/build/SKILL.md)

---

### @review

**Purpose:** Quality check for completed feature

**Usage:**
```bash
@review F{ID}
```

**What it checks:**
- All workstreams completed
- Tests passing
- Coverage ≥80%
- Type hints complete
- No TODO markers

**Example:**
```bash
@review F001
```

**See:** [.claude/skills/review/SKILL.md](../../.claude/skills/review/SKILL.md)

---

### @deploy

**Purpose:** Deploy feature to production

**Usage:**
```bash
@deploy F{ID}
```

**What it does:**
1. Final quality verification
2. Creates git tag
3. Merges to main branch
4. Triggers deployment pipeline

**Pre-deployment Checklist:**
- [ ] All workstreams completed
- [ ] All tests passing
- [ ] Coverage ≥80%
- [ ] Code review approved
- [ ] Documentation updated

**Example:**
```bash
@deploy F001
```

**See:** [.claude/skills/deploy/SKILL.md](../../.claude/skills/deploy/SKILL.md)

---

## Utility Commands

### @oneshot

**Purpose:** Autonomous feature execution

**Usage:**
```bash
@oneshot F{ID}
```

**What it does:**
1. Spawns orchestrator agent
2. Executes all workstreams
3. Handles dependencies
4. Save/restore checkpoints
5. Progress notifications

**Features:**
- Checkpoint save/restore
- Background execution
- Telegram notifications
- Resume after interruption

**Example:**
```bash
@oneshot F001
```

**See:** [.claude/skills/oneshot/SKILL.md](../../.claude/skills/oneshot/SKILL.md)

---

### /debug

**Purpose:** Systematic debugging

**Usage:**
```bash
/debug "Problem description"
```

**Process:**
1. Observe problem
2. Form hypothesis
3. Design experiment
4. Run experiment
5. Update hypothesis

**Example:**
```bash
/debug "Test fails unexpectedly"
```

**See:** [.claude/skills/debug/SKILL.md](../../.claude/skills/debug/SKILL.md)

---

### @issue

**Purpose:** Bug routing and classification

**Usage:**
```bash
@issue "Bug description"
```

**Routes to:**
- `@hotfix` for P0 (critical) issues
- `@bugfix` for P1/P2 (quality) issues
- Backlog for P3 (minor) issues

**Example:**
```bash
@issue "Login fails on Firefox"
```

**See:** [.claude/skills/issue/SKILL.md](../../.claude/skills/issue/SKILL.md)

---

### @hotfix

**Purpose:** Emergency fix for P0 issues

**Usage:**
```bash
@hotfix "Critical security vulnerability"
```

**Characteristics:**
- Branch from main
- Minimal changes only
- Deploy < 2 hours
- Skip full process

**Example:**
```bash
@hotfix "Production database connection fails"
```

**See:** [.claude/skills/hotfix/SKILL.md](../../.claude/skills/hotfix/SKILL.md)

---

### @bugfix

**Purpose:** Quality fix for P1/P2 issues

**Usage:**
```bash
@bugfix "Incorrect totals calculation"
```

**Process:**
1. Branch from feature/develop
2. Full TDD cycle
3. Quality gates enforced
4. No production deploy

**Example:**
```bash
@bugfix "User profile image not loading"
```

**See:** [.claude/skills/bugfix/SKILL.md](../../.claude/skills/bugfix/SKILL.md)

---

## Internal Commands

### /tdd

**Purpose:** TDD cycle enforcement (internal)

**Usage:**
Automatic (called by @build)

**Process:**
1. Red - Write failing test
2. Green - Write minimal code
3. Refactor - Improve code

**See:** [.claude/skills/tdd/SKILL.md](../../.claude/skills/tdd/SKILL.md)

---

### @idea

**Purpose:** Requirements gathering (internal)

**Usage:**
```bash
@idea "Feature concept"
```

**Process:**
- Deep questioning via AskUserQuestion
- Explores tradeoffs
- Generates comprehensive spec

**Called by:** @feature command

**See:** [.claude/skills/idea/SKILL.md](../../.claude/skills/idea/SKILL.md)

---

## Command Options

### Verbosity Levels

Most commands support verbosity:

```bash
@build WS-001-01 --verbose
@review F001 --quiet
```

### Background Execution

For long-running commands:

```bash
@oneshot F001 --background
```

### Checkpoint Resume

Resume interrupted execution:

```bash
@oneshot F001 --resume <agent-id>
```

---

## Quick Reference

| Command | Purpose | Time |
|---------|---------|------|
| `@feature` | Create feature | 10-15 min |
| `@design` | Plan workstreams | 5-10 min |
| `@build` | Execute workstream | 30-90 min |
| `@review` | Quality check | 5-10 min |
| `@deploy` | Deploy feature | 5-10 min |
| `@oneshot` | Autonomous execution | 2-6 hours |
| `/debug` | Debug issue | 15-60 min |
| `@issue` | Report bug | 2-5 min |
| `@hotfix` | Emergency fix | < 2 hours |
| `@bugfix` | Quality fix | 1-4 hours |

---

## Command Flow

### Standard Feature Development

```
@feature → @design → @build → @build → ... → @review → @deploy
```

### Autonomous Execution

```
@feature → @oneshot → @review → @deploy
```

### Bug Fix Flow

```
@issue → @hotfix OR @bugfix → @review
```

---

## See Also

- [skills.md](skills.md) - Skill system details
- [quality-gates.md](quality-gates.md) - Quality standards
- [beginner/02-common-tasks.md](../beginner/02-common-tasks.md) - Common workflows

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
