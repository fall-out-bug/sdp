# Verbosity Tiers Specification

**Version:** 1.0.0
**Status:** Active
**Last Updated:** 2026-02-08

## Overview

All SDP skills MUST support four verbosity levels to provide users with control over output detail.

## Levels

```
Level 0 (--quiet):   Exit status only (✅/❌)
Level 1 (default):   Summary (1-3 lines with key metrics)
Level 2 (--verbose): Step-by-step progress
Level 3 (--debug):   Internal state + API calls
```

## Per-Skill Examples

### @build Skill

```bash
# Quiet mode
@build 00-050-01 --quiet
# Output: ✅

# Default mode
@build 00-050-01
# Output: ✅ 00-050-01: Workstream Parser (22m, 85%, commit:abc123)

# Verbose mode
@build 00-050-01 --verbose
# Output:
# → Reading WS spec...
# → TDD cycle: Red (3m) → Green (12m) → Refactor (7m)
# → Quality check: PASS (coverage 85%, mypy clean)
# ✅ COMPLETE

# Debug mode
@build 00-050-01 --debug
# Output:
# [DEBUG] Loading workstream: docs/workstreams/backlog/00-050-01.md
# [DEBUG] Scope files: [src/sdp/parser.py, tests/sdp/test_parser.py]
# → Reading WS spec...
# [DEBUG] AC count: 5
# → TDD cycle: Red (3m) → Green (12m) → Refactor (7m)
# [DEBUG] Test run: pytest tests/sdp/parser.py -v
# [DEBUG] Coverage: 85.3%
# → Quality check: PASS (coverage 85%, mypy clean)
# ✅ COMPLETE
```

### @review Skill

```bash
# Quiet mode
@review F01 --quiet
# Output: ✅

# Default mode
@review F01
# Output: ✅ APPROVED (QA:82%, Security:PASS, 5 agents)

# Verbose mode
@review F01 --verbose
# Output:
# → Spawning 6 review agents...
# → QA review: PASS (82% coverage, 145/145 tests)
# → Security review: PASS (no vulnerabilities)
# → DevOps review: PASS (CI/CD validated)
# → SRE review: PASS (SLOs defined)
# → TechLead review: PASS (code quality good)
# → Documentation review: PASS (0% drift)
# ✅ APPROVED

# Debug mode
@review F01 --debug
# Output:
# [DEBUG] Feature: F01
# [DEBUG] Workstreams: 5
# [DEBUG] Spawning agents via Task tool...
# [DEBUG] Agent 1: QA (subagent_type=general-purpose)
# [DEBUG] Agent 2: Security (subagent_type=general-purpose)
# → Spawning 6 review agents...
# [QA agent output...]
# → QA review: PASS (82% coverage, 145/145 tests)
# → Security review: PASS (no vulnerabilities)
# → DevOps review: PASS (CI/CD validated)
# → SRE review: PASS (SLOs defined)
# → TechLead review: PASS (code quality good)
# → Documentation review: PASS (0% drift)
# ✅ APPROVED
```

### @feature Skill

```bash
# Quiet mode
@feature "Add OAuth2" --quiet
# Output: ✅ (4 workstreams created)

# Default mode
@feature "Add OAuth2"
# Output: ✅ Feature F051: OAuth2 Authentication (4 workstreams, 2.1h est.)

# Verbose mode
@feature "Add OAuth2" --verbose
# Output:
# → @idea phase: Requirements gathering (18 questions, 15m)
# → @design phase: Architecture design (5 discovery blocks, 45m)
# → Workstreams created: 00-051-01 through 00-051-04
# ✅ COMPLETE

# Debug mode
@feature "Add OAuth2" --debug
# Output:
# [DEBUG] Feature description: "Add OAuth2"
# [DEBUG] Starting @idea phase...
# [DEBUG] Question cycle 1: 3 questions
# [DEBUG] Question cycle 2: 3 questions
# → @idea phase: Requirements gathering (18 questions, 15m)
# [DEBUG] Starting @design phase...
# [DEBUG] Discovery block 1: 3 questions
# [DEBUG] Discovery block 2: 3 questions
# → @design phase: Architecture design (5 discovery blocks, 45m)
# [DEBUG] Creating workstream files...
# [DEBUG] Created: docs/workstreams/backlog/00-051-01.md
# [DEBUG] Created: docs/workstreams/backlog/00-051-02.md
# [DEBUG] Created: docs/workstreams/backlog/00-051-03.md
# [DEBUG] Created: docs/workstreams/backlog/00-051-04.md
# → Workstreams created: 00-051-01 through 00-051-04
# ✅ COMPLETE
```

### @vision Skill

```bash
# Quiet mode
@vision "AI task manager" --quiet
# Output: ✅ (PRODUCT_VISION.md, PRD.md, ROADMAP.md)

# Default mode
@vision "AI task manager"
# Output: ✅ Vision: AI-Powered Task Manager (7 experts, 45m)

# Verbose mode
@vision "AI task manager" --verbose
# Output:
# → Quick interview: 5 questions (8m)
# → Deep-thinking: 7 parallel experts (30m)
# → Artifacts generated: PRODUCT_VISION.md, PRD.md, ROADMAP.md
# ✅ COMPLETE

# Debug mode
@vision "AI task manager" --debug
# Output:
# [DEBUG] Product idea: "AI task manager"
# [DEBUG] Starting quick interview...
# [DEBUG] Question 1: What problem are you solving?
# [DEBUG] Question 2: Who are your target users?
# → Quick interview: 5 questions (8m)
# [DEBUG] Spawning 7 parallel experts...
# [DEBUG] Expert 1: Product (subagent_type=general-purpose)
# [DEBUG] Expert 2: Market (subagent_type=general-purpose)
# → Deep-thinking: 7 parallel experts (30m)
# [DEBUG] Generating artifacts...
# [DEBUG] Created: PRODUCT_VISION.md
# [DEBUG] Created: docs/prd/PRD.md
# [DEBUG] Created: docs/roadmap/ROADMAP.md
# → Artifacts generated: PRODUCT_VISION.md, PRD.md, ROADMAP.md
# ✅ COMPLETE
```

### @reality Skill

```bash
# Quiet mode
@reality --quick --quiet
# Output: ✅ (Health: 72/100)

# Default mode
@reality --quick
# Output: ✅ Reality Check: 15K LOC, Health 72/100, 3 critical issues

# Verbose mode
@reality --quick --verbose
# Output:
# → Project scan: 15,234 LOC, 127 files (Go)
# → Architecture: Layer violations detected (2)
# → Testing: Coverage ~65% (below 80% threshold)
# → Documentation: 4 drift issues found
# ⚠️ 3 critical issues found
# ✅ COMPLETE

# Debug mode
@reality --quick --debug
# Output:
# [DEBUG] Mode: quick
# [DEBUG] Detecting project type...
# [DEBUG] Found go.mod → Go project
# → Project scan: 15,234 LOC, 127 files (Go)
# [DEBUG] Scanning src/ directory...
# [DEBUG] Analyzing imports...
# [DEBUG] Checking for circular dependencies...
# → Architecture: Layer violations detected (2)
# [DEBUG] Scanning tests/ directory...
# [DEBUG] Found 45 test files
# → Testing: Coverage ~65% (below 80% threshold)
# [DEBUG] Comparing docs to code...
# → Documentation: 4 drift issues found
# ⚠️ 3 critical issues found
# ✅ COMPLETE
```

## Implementation Pattern

### Skill File Header

Each skill MUST document verbosity tiers in its header:

```markdown
## Verbosity Tiers

```bash
@skill "args" --quiet     # Exit status only
@skill "args"             # Summary with metrics (default)
@skill "args" --verbose   # Step-by-step progress
@skill "args" --debug     # Internal state + API calls
```

### Output Format

**Quiet (--quiet):**
- Single symbol: `✅` or `❌`
- No additional text
- For scripts/ci

**Default (no flags):**
- 1-3 lines max
- Key metrics only
- Human-readable

**Verbose (--verbose):**
- Step-by-step progress
- Bullet points or arrows (→)
- Clear section headers
- No [DEBUG] tags

**Debug (--debug):**
- All verbose output
- Plus [DEBUG] tagged internal state
- API calls shown
- File paths logged
- Intermediate values

## Detection Logic

Skills should auto-detect verbosity from invocation:

```python
# In skill execution
def execute_skill(args):
    verbosity = "default"
    if "--quiet" in args or "-q" in args:
        verbosity = "quiet"
    elif "--verbose" in args or "-v" in args:
        verbosity = "verbose"
    elif "--debug" in args or "-d" in args:
        verbosity = "debug"

    return generate_output(verbosity)
```

## Quality Gates

- **Quiet mode:** MUST still exit with correct status code
- **Default mode:** MUST fit in 3 lines or less
- **Verbose mode:** MUST NOT include [DEBUG] tags
- **Debug mode:** MUST include all verbose output plus debug info

## Migration Checklist

For each skill:
- [ ] Add verbosity tiers section to SKILL.md
- [ ] Document output format for each level
- [ ] Provide examples for each level
- [ ] Update workflow to respect verbosity flags
- [ ] Test all 4 levels manually

## Rollout Plan

1. **Phase 1:** Implement for core skills (@build, @review, @feature)
2. **Phase 2:** Implement for planning skills (@idea, @design, @vision, @reality)
3. **Phase 3:** Implement for utility skills (@deploy, @hotfix, @bugfix)
4. **Phase 4:** Implement for remaining skills

## Version History

- **1.0.0** (2026-02-08): Initial specification
