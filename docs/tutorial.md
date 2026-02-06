# SDP Tutorial

Welcome to the SDP (Spec-Driven Protocol) tutorial! This guide will walk you through the complete workflow from project setup to deployment.

## Prerequisites

Before starting, ensure you have:

- Git initialized
- Claude Code installed
- SDP prompts installed in `.claude/`
- (Optional) Go binary `sdp` for convenience features

## Quick Start (15 min)

New to SDP? Start with our [beginner tutorial](beginner/00-quick-start.md) for a hands-on introduction.

## Tutorial Sections

### 1. Project Initialization

Learn how to set up SDP in your project using the interactive wizard.

```bash
sdp init
```

**Topics covered:**
- Interactive metadata collection
- Directory structure creation
- Quality gate configuration
- Git hooks installation
- Health check validation

**Documentation:** [.claude/skills/init/SKILL.md](../.claude/skills/init/SKILL.md)

---

### 2. Feature Development

Discover the progressive disclosure workflow for developing features.

```bash
@feature "Your feature name"
```

**Workflow stages:**
1. **Vision Interview** - Define mission, users, success metrics
2. **Technical Interview** - Explore architecture, tradeoffs
3. **Specification** - Generate machine-readable intent
4. **Draft Creation** - Create comprehensive requirements doc

**Outputs:**
- `PRODUCT_VISION.md` - Project manifesto
- `docs/drafts/idea-{slug}.md` - Feature specification
- `docs/intent/{slug}.json` - Machine-readable intent

**Documentation:** [.claude/skills/feature/SKILL.md](../.claude/skills/feature/SKILL.md)

---

### 3. Requirements Gathering

Deep dive into feature requirements with interactive interviewing.

```bash
@idea "Your feature idea"
```

**Interview topics:**
- Problem statement
- User personas
- Success criteria
- Technical approach
- Integration points
- Failure modes
- Security considerations

**Documentation:** [.claude/skills/idea/SKILL.md](../.claude/skills/idea/SKILL.md)

---

### 4. Workstream Design

Break features into atomic, executable workstreams.

```bash
@design idea-{slug}
```

**Process:**
1. Codebase exploration
2. Dependency analysis
3. Workstream decomposition
4. Approval request

**Output:** Workstream markdown files in `docs/workstreams/backlog/`

**Documentation:** [.claude/skills/design/SKILL.md](../.claude/skills/design/SKILL.md)

---

### 5. Workstream Execution

Implement workstreams using Test-Driven Development.

```bash
@build WS-XXX-YY
```

**TDD Cycle:**

#### Phase 1: Red
- Write failing test
- Verify test fails
- Commit test

#### Phase 2: Green
- Implement minimum code
- Verify test passes
- Commit implementation

#### Phase 3: Refactor
- Improve code quality
- Verify tests still pass
- Commit refactoring

#### Phase 4: Verify
- Check acceptance criteria
- Run quality gates
- Append execution report

**Quality Gates:**
- Coverage ≥ 80%
- Files < 200 LOC
- Complexity < 10
- Type hints 100%
- No architecture violations

**Documentation:** [.claude/skills/build/SKILL.md](../.claude/skills/build/SKILL.md)

---

### 6. Autonomous Execution

Execute entire features automatically with multi-agent orchestration.

```bash
@oneshot F01
```

**Features:**
- Parallel workstream execution
- Dependency tracking
- Checkpoint/resume capability
- Background execution support
- PR-less execution mode

**Documentation:** [.claude/skills/oneshot/SKILL.md](../.claude/skills/oneshot/SKILL.md)

---

### 7. Quality Review

Validate feature quality before deployment.

```bash
@review F01
```

**Review checks:**
- All workstreams completed
- Acceptance criteria met
- Quality gates passed
- No TODOs without followup
- Documentation complete

**Documentation:** [.claude/skills/review/SKILL.md](../.claude/skills/review/SKILL.md)

---

### 8. Deployment

Deploy features to production.

```bash
@deploy F01
```

**Deployment steps:**
1. Merge feature to main
2. Create version tag
3. Generate release notes
4. Run final validation

**Documentation:** [.claude/skills/deploy/SKILL.md](../.claude/skills/deploy/SKILL.md)

---

### 9. Debugging

Systematic debugging using scientific method.

```bash
@debug "Issue description"
```

**Process:**
1. Observe symptoms
2. Form hypotheses
3. Design experiments
4. Analyze results
5. Propose solution

**Documentation:** [.claude/skills/debug/SKILL.md](../.claude/skills/debug/SKILL.md)

---

### 10. Health Checks

Verify SDP installation and configuration.

```bash
sdp doctor
```

**Checks performed:**
- Python version (≥3.10)
- Poetry installation
- Git hooks configuration
- Beads CLI (optional)
- GitHub CLI (optional)
- Telegram config (optional)

**Documentation:** [src/sdp/doctor.py](../src/sdp/doctor.py)

---

## Common Workflows

### Creating a New Feature

```bash
# 1. Initialize project (first time only)
sdp init

# 2. Create feature specification
@feature "Add user authentication"

# 3. Design workstreams
@design idea-user-auth

# 4. Execute each workstream
@build WS-F01-01
@build WS-F01-02
@build WS-F01-03
@build WS-F01-04

# 5. Review quality
@review F01

# 6. Deploy to production
@deploy F01
```

### Autonomous Feature Development

```bash
# Single command for entire feature
@oneshot F01

# Resume if interrupted
@oneshot F01 --resume <agent-id>

# Background execution
@oneshot F01 --background
```

### Bug Fixing

```bash
# Classify and route bug
@issue "Login fails on Firefox"

# For P0 emergencies
@hotfix "Critical API outage"

# For P1/P2 quality issues
@bugfix "Incorrect totals"
```

---

## Best Practices

### Workstream Design

- Keep workstreams SMALL (< 500 LOC)
- One workstream = one atomic task
- Define clear acceptance criteria
- List dependencies explicitly

### Code Quality

- Write tests first (TDD)
- Keep files under 200 LOC
- Maintain 80%+ coverage
- Use type hints everywhere
- Follow clean architecture

### Git Workflow

- Commit early, commit often
- Use conventional commits
- Never skip `@build` verification
- Always run `@review` before `@deploy`

---

## Troubleshooting

### Common Issues

**Issue:** `sdp: command not found`
```bash
pip install sdp
# or
poetry add sdp --group dev
```

**Issue:** `sdp doctor` shows failures
```bash
# Install Poetry
curl -sSL https://install.python-poetry.org | python3 -

# Initialize git
git init

# Re-run setup
sdp init --force
```

**Issue:** Tests failing in `@build`
```bash
# Run tests manually
pytest tests/unit/test_module.py -v

# Check coverage
pytest --cov=module --cov-report=term-missing

# Debug systematically
@debug "Test failure in test_module.py"
```

**Issue:** Workstream blocked
```bash
# Check dependencies
cat docs/workstreams/INDEX.md

# Complete blocking workstreams first
@build WS-XXX-01  # Blocks WS-XXX-02
```

---

## Advanced Topics

### Multi-Agent Orchestration

Learn about autonomous agent coordination for complex features.

**Documentation:** [.claude/agents/](../.claude/agents/)

### Extension Development

Create custom extensions for SDP.

**Documentation:** [src/sdp/extensions/](../src/sdp/extensions/)

### Quality Gates

Configure custom quality rules for your project.

**Documentation:** [quality-gate.toml](../quality-gate.toml)

---

## Additional Resources

### Core Documentation

- [PROTOCOL.md](PROTOCOL.md) - Full SDP specification
- [PRINCIPLES.md](PRINCIPLES.md) - Core principles
- [CODE_PATTERNS.md](CODE_PATTERNS.md) - Code patterns
- [MODELS.md](MODELS.md) - Model recommendations

### Guides

- [CLAUDE_CODE.md](guides/CLAUDE_CODE.md) - Claude Code integration
- [CURSOR.md](guides/CURSOR.md) - Cursor IDE integration
- [GIT_WORKFLOW.md](guides/GIT_WORKFLOW.md) - Git workflow

### Runbooks

- [debug-runbook.md](runbooks/debug-runbook.md) - Debugging procedures
- [test-runbook.md](runbooks/test-runbook.md) - Testing procedures
- [oneshot-runbook.md](runbooks/oneshot-runbook.md) - Autonomous execution

---

## Getting Help

### Community

- GitHub: [github.com/your-org/sdp](https://github.com/your-org/sdp)
- Discord: [discord.gg/sdp](https://discord.gg/sdp)
- Documentation: [docs.sdp.dev](https://docs.sdp.dev)

### Reporting Issues

Found a bug? Have a feature request?

1. Search existing issues
2. Create new issue with template
3. Provide minimal reproduction
4. Include SDP version and environment

---

## Version

**SDP Version:** 0.3.0
**Last Updated:** 2025-01-29
**Tutorial Version:** 1.0
