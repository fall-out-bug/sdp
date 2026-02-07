# 00-052-04: @reality Skill Structure

> **Beads ID:** sdp-uugq
> **Feature:** F052 - Multi-Agent SDP + @vision + @reality
> **Phase:** 1B - Analysis Skills (@reality)
> **Size:** MEDIUM
> **Duration:** 2-3 days
> **Dependencies:**
> - 00-052-00 (Backup & worktree setup)

## Goal

Create `.claude/skills/reality/SKILL.md` for codebase analysis and health assessment.

## Acceptance Criteria

- **AC1:** `.claude/skills/reality/` directory created
- **AC2:** `SKILL.md` created with 3 modes (--quick, --deep, --focus)
- **AC3:** 8 expert agents specified (architecture, quality, testing, security, performance, docs, debt, standards)
- **AC4:** Universal project analysis (works without SDP)
- **AC5:** Output report format defined

## Files

**Create:**
- `.claude/skills/reality/SKILL.md`

**Create (directory):**
- `.claude/skills/reality/`
- `src/sdp/reality/`
- `tests/sdp/reality/`

## Steps

### Step 1: Create Directory Structure

```bash
mkdir -p .claude/skills/reality
mkdir -p src/sdp/reality
mkdir -p tests/sdp/reality
```

### Step 2: Write SKILL.md

Create `.claude/skills/reality/SKILL.md`:

```markdown
---
name: reality
description: Codebase analysis and architecture validation - what's actually there vs documented
tools: Read, Glob, Grep, Bash, AskUserQuestion, Task
version: 1.0.0
---

# @reality - Codebase Analysis

**Understand what's actually in your codebase.**

## When to Use

- **New to project** - "Что здесь вообще есть?"
- **Before @feature** - "На что можем опираться?"
- **After @vision** - "Как код соответствует видению?"
- **Quarterly review** - Track tech debt and quality trends
- **Debugging mysteries** - "Почему это не работает?"

## Modes

### --quick (5-10 minutes)

Fast analysis for quick health check:

```bash
@reality --quick
```

**Analysis:**
- Language/framework detection
- File structure overview
- Quick metrics (LOC, files, packages)
- Obvious issues (missing tests, TODOs)

**Output:** Terminal summary + JSON report

### --deep (30-60 minutes)

Comprehensive analysis with 8 expert agents:

```bash
@reality --deep
```

**Analysis:**
- Architecture: Layer violations, coupling, patterns
- Quality: Code smells, complexity, duplication
- Testing: Coverage, gaps, flakiness
- Security: Vulnerabilities, secrets, dependencies
- Performance: Bottlenecks, resource usage
- Documentation: Completeness, accuracy
- Technical Debt: Magnitude, prioritization
- Standards: Linting, conventions, best practices

**Output:** Detailed report + recommendations

### --focus=topic (10-15 minutes)

Deep dive on specific aspect:

```bash
@reality --focus=security
@reality --focus=architecture
@reality --focus=testing
@reality --focus=debt
```

## Workflow

### Step 1: Detect Project Type

```python
# Scan for language/framework indicators
project_type = detect_project(language, framework)
# Examples: Go service, React app, Python monolith
```

### Step 2: Universal Analysis (--quick mode)

```python
# Works on ANY project, even without SDP
analysis = {
  "language": detect_language(),
  "frameworks": detect_frameworks(),
  "structure": analyze_structure(),
  "metrics": calculate_metrics(),
  "issues": find_obvious_issues()
}
```

### Step 3: Expert Agent Analysis (--deep mode)

Spawn 8 expert agents in parallel:

```python
experts = [
  Task("Architecture expert", "Analyze layers, patterns, violations"),
  Task("Quality expert", "Find code smells, complexity issues"),
  Task("Testing expert", "Analyze coverage, test quality"),
  Task("Security expert", "Find vulnerabilities, exposed secrets"),
  Task("Performance expert", "Find bottlenecks, slow operations"),
  Task("Documentation expert", "Check completeness, accuracy"),
  Task("Technical debt expert", "Categorize, prioritize debt"),
  Task("Standards expert", "Check linting, conventions")
]
```

### Step 4: Generate Report

**docs/reality/REALITY_REPORT.md:**
```markdown
# Reality Report

## Project Overview
- **Type:** {project_type}
- **Language:** {language}
- **Frameworks:** {frameworks}
- **Size:** {LOC}

## Health Score: {score}/100

### Architecture: {score}/100
- Strengths: {list}
- Issues: {list}
- Recommendations: {list}

### Quality: {score}/100
...

## Vision vs Reality Gap
{Compare PRODUCT_VISION.md vs actual code}

## Technical Debt
{Prioritized list}

## Action Items
1. {High priority fix}
2. {Medium priority fix}
```

## Outputs

- **Terminal:** Summary table
- **JSON:** `.oneshot/reality-report.json` (machine-readable)
- **Markdown:** `docs/reality/REALITY_REPORT.md` (human-readable)

## Example

```bash
@reality --quick

→ Detected: Go service (2,450 LOC across 38 files)
→ Health: 72/100 (Good)
→ Issues: 3 TODOs, 12% test coverage gap
→ Recommendations: Add integration tests, update docs
```

## See Also

- `.claude/skills/vision/SKILL.md` - Strategic planning
- `src/sdp/reality/scanner.go` - Project scanner implementation
```

### Step 3: Verify Skill Format

```bash
ls -la .claude/skills/reality/
cat .claude/skills/reality/SKILL.md | head -20
```

Expected: Directory exists, SKILL.md has frontmatter with name/description/tools/version

## Quality Gates

- SKILL.md valid frontmatter (name, description, tools, version)
- All 3 modes documented (--quick, --deep, --focus)
- 8 expert agents listed
- Universal analysis documented (works without SDP)
- Output format specified

## Success Metrics

- Skill can be invoked via Skill tool
- --quick mode completes in <10 minutes
- --deep mode spawns 8 expert agents
- Works on projects without SDP
