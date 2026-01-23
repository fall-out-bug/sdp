# Spec-Driven Protocol (SDP)

Workstream-driven development protocol for AI agents with Developer-in-the-Loop structured workflow.

## Terminology

| Term | Scope | Size | Example |
|------|-------|------|---------|
| **Release** | Product milestone | 10-30 Features | R1: Submissions E2E |
| **Feature** | Large feature | 5-30 Workstreams | F24: Obsidian Vault |
| **Workstream** | Atomic task | SMALL/MEDIUM/LARGE | WS-140: Vault Domain |

**Scope metrics:**
- **SMALL**: < 500 LOC, < 1500 tokens
- **MEDIUM**: 500-1500 LOC, 1500-5000 tokens
- **LARGE**: > 1500 LOC → split into 2+ WS

**⚠️ NO time-based estimates** (hours/days/weeks). Use scope metrics only.

**Hierarchy:**
```
Product:      PORTAL_VISION → RELEASE_PLAN → Feature (F) → Workstream (WS)
Architecture: L1 (System) → L2 (Domain) → L3 (Component) → L4 (Workstream)
```

**Deprecated:** ~~Epic (EP)~~ → Feature (F) since 2026-01-07

## Installation

### From Source

```bash
# Navigate to sdp directory
cd sdp

# Install in development mode using Poetry
poetry install

# Or using pip
pip install -e .
```

### Verify Installation

After installation, verify the package is working:

```bash
# Check version
sdp --version
# Output: sdp version 0.3.0

# Show help
sdp --help

# Test core commands
sdp core parse-ws <path-to-ws-file.md>
sdp core parse-project-map <path-to-PROJECT_MAP.md>

# GitHub integration commands
sdp-github sync-all --ws-dir <path-to-workstreams>
sdp-github sync-ws <path-to-ws-file.md>
```

### Git Hooks

SDP uses cross-platform Git hooks for quality checks that work across all IDEs (Claude Code, Cursor, OpenCode).

**Installation:**
```bash
# From repository root
bash sdp/hooks/install-hooks.sh
```

**Available Hooks:**
- `pre-commit` - Quality checks (time estimates, code quality, Python code, Clean Architecture, WS format, breaking changes, test quality)
- `post-commit` - GitHub issue sync (if GITHUB_TOKEN set)
- `pre-push` - Regression tests and coverage checks

**Uninstallation:**
```bash
rm .git/hooks/pre-commit .git/hooks/post-commit .git/hooks/pre-push
```

**Environment Variables (for GitHub integration):**
```bash
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_REPO="owner/repo"
```

### Package Structure

The package provides:

- **Core modules** (`sdp.core`):
  - `parse_workstream()` - Parse workstream markdown files
  - `parse_project_map()` - Parse PROJECT_MAP.md files
  - `Feature` - Feature decomposition with dependency management
  - `ProjectMap` - Project-level decisions and constraints

- **GitHub integration** (`sdp.github`):
  - Workstream to GitHub issue syncing
  - Project board automation
  - Milestone management

- **CLI commands**:
  - `sdp` - Main CLI with core operations
  - `sdp-github` - GitHub integration commands

## Quick Start

### Option 1: Slash Commands (Recommended)

**Feature Development:**
```bash
# 1. Gather requirements
/idea "LMS integration for courses"

# 2. Design all workstreams
/design idea-lms-integration

# 3a. Execute manually (step-by-step)
/test WS-060-01    # Generate test contract
/build WS-060-01   # Implement (makes tests GREEN)
/test WS-060-02
/build WS-060-02

# 3b. Execute autonomously (hands-off)
/oneshot F60

# 4. Review feature
/codereview F60

# 5. Deploy (after human UAT)
/deploy F60
```

**Issue Management:**
```bash
# 1. Analyze issue
/issue "API returns 500 on submissions"
# → Routes to /hotfix (P0) or /bugfix (P1/P2)

# 2. Debug (if root cause unknown)
/debug "API returns 500 on /submissions endpoint"
# → Systematic debugging: Symptom → Hypothesis → Test → Root Cause → Impact

# 3a. Emergency fix (P0 CRITICAL)
/hotfix "fix duplicate ID race" --issue-id=001
# → Deploy to production (SLA target: < 2h)

# 3b. Feature bug fix (P1/P2)
/bugfix "handle large repos" --feature=F23 --issue-id=002
# → Full TDD cycle → merge to feature/develop
```

See: `sdp/prompts/commands/` for command specs.

**Contract-Driven Workflow:**

The `/test` command implements contract-driven development:

```bash
/test WS-060-01  # Generate test contract (T0 tier)
/build WS-060-01 # Implement to make tests GREEN
```

**Principles:**
- Tests = contract (single source of truth)
- Tests created before implementation
- Tests executable (fail with NotImplementedError)
- Tests NOT changed during /build
- Capability tiers: T0 (contract), T1 (basic), T2 (refactor), T3 (fill-in)

See: `tools/hw_checker/docs/workstreams/completed/2026-01/WS-410-01-contract-driven-ws-spec.md` for full spec.

### Sub-Agents

SDP provides specialized sub-agents for different phases of development:

| Agent | Description | When to Use | Tools |
|-------|-------------|-------------|-------|
| **planner** | Codebase analysis and WS decomposition | `/design`, creating implementation plans | Read, Glob, Grep, Bash |
| **builder** | TDD execution for single WS | `/build WS-XXX-YY` | Read, Write, Edit, Bash, Glob, Grep |
| **reviewer** | 17-point quality checklist review | `/codereview F{XX}` | Read, Bash, Glob, Grep |
| **deployer** | DevOps, CI/CD, git merge/tag | `/deploy F{XX}` | Read, Write, Edit, Bash, Glob, Grep |
| **orchestrator** | Autonomous feature execution | `/oneshot F{XX}` | Read, Write, Edit, Bash, Glob, Grep |

**All agents delegate to master prompts** from `sdp/prompts/commands/` for single source of truth.

**Multi-IDE Support:**
- **Claude Code:** `.claude/agents/` + `.claude/skills/`
- **Cursor:** `.cursor/agents/` (same format as Claude Code)
- **OpenCode:** `.opencode/opencode.json` (JSON config format)

All agents reference same master prompts, ensuring consistent behavior across IDEs.

### Multi-IDE Parity Documentation

SDP supports three AI coding tools with full parity:

- **[Multi-IDE Parity Documentation](docs/multi-ide-parity.md)** - Complete parity matrix and architecture
- **[Oneshot Runbook](docs/runbooks/oneshot-runbook.md)** - Autonomous feature execution
- **[Debug Runbook](docs/runbooks/debug-runbook.md)** - Systematic debugging workflow
- **[Test Runbook](docs/runbooks/test-runbook.md)** - Contract-driven test generation
- **[Git Hooks Installation](docs/runbooks/git-hooks-installation.md)** - Cross-platform hooks

**Key Features:**
- ✅ Same master prompts across all IDEs
- ✅ Consistent workflows and behavior
- ✅ Universal Git hooks (work in all IDEs)
- ✅ OpenCode-specific format constraints documented
- ✅ Parity matrix for feature comparison

**Issue Management:**
```bash
# 1. Analyze issue
/issue "API returns 500 on submissions"
# → Routes to /hotfix (P0) or /bugfix (P1/P2)

# 2a. Emergency fix (P0 CRITICAL)
/hotfix "fix duplicate ID race" --issue-id=001
# → Deploy to production (SLA target: < 2h)

# 2b. Feature bug fix (P1/P2)
/bugfix "handle large repos" --feature=F23 --issue-id=002
# → Full TDD cycle → merge to feature/develop
```

See: `sdp/prompts/commands/` for command specs.

### Option 2: 4-Phase Workflow (DEPRECATED)

> ⚠️ **DEPRECATED:** Use slash commands instead. See `sdp/prompts/structured/DEPRECATED.md` for migration guide.

Legacy 4-phase workflow is superseded by `/design`, `/build`, `/codereview`, `/deploy` commands which provide better UX and single source of truth.

## File Organization

```
sdp/
├── PROTOCOL.md                 # Full protocol specification
├── HW_CHECKER_PATTERNS.md      # Code patterns reference
├── README.md                   # This file
├── prompts/
│   ├── commands/               # Slash commands (/idea, /design, /build, /oneshot, /issue, /hotfix, /bugfix, /codereview, /deploy)
│   └── structured/             # 4-phase workflow (phase-1 through phase-4)
├── templates/                  # Templates (uat-guide, idea-draft, release-notes)
├── hooks/                      # Git hooks (pre-commit, post-commit, pre-push, install-hooks.sh)
├── schema/                     # JSON Schema for validation
└── archive/                    # Archived v1.2 materials
```

## Core Concepts

### Workstream-Driven Development

- **Workstream (WS):** Small, self-contained task executable in one shot by AI agent
- **One-shot execution:** Agent completes WS without iterative human feedback
- **AI-Readiness:** Files < 200 LOC, CC < 10, full type hints, coverage ≥80%

### Structured vs Multi-Agent Mode

**Structured Mode (recommended for most work):**
- 4 phases with human checkpoints
- Single developer + AI agents
- Faster, less coordination overhead

**Multi-Agent Mode (for complex epics):**
- Individual agent prompts (Analyst, Architect, Tech Lead, Developer, QA, DevOps)
- Parallel execution possible
- Use for large system changes

## Key Features

✅ **TDD enforcement:** Red-Green-Refactor cycle in every WS  
✅ **SOLID principles:** Explicit checks in review (SRP, OCP, LSP, ISP, DIP)  
✅ **Clean Architecture:** Layer violations detected automatically  
✅ **Strict typing:** mypy --strict validation  
✅ **Code patterns:** 8 patterns + 5 anti-patterns documented  
✅ **Coverage requirement:** ≥80% mandatory, enforced in CI  
✅ **Regression suite:** Fast tests must pass after every WS  

## Guardrails

Non-negotiable rules that block progress if violated:

1. **AI-Readiness:** Files < 200 LOC, CC < 10, type hints everywhere
2. **Clean Architecture:** No layer violations (Domain ← App ← Infra ← Presentation)
3. **Error Handling:** No silent failures (`except: pass` forbidden)
4. **Security:** No privileged containers, resource limits, no string interpolation in commands
5. **No TODOs:** All tasks must be completed or explicitly deferred to new WS

## Quality Gates

Checkpoints before phase transitions:

- **Analyze → Plan:** Workstream map complete, dependencies clear, scope estimated
- **Plan → Execute:** Pre-flight checks passed (no duplicates, scope ≤ MEDIUM)
- **Execute → Review:** Coverage ≥80%, regression passing, no TODO comments
- **Review → Next WS:** All checklists ✅ (AI-Ready, SOLID, tests, architecture)

## Usage Examples

### Example 1: Refactor Large File (AI-Readiness)

```bash
# Phase 1: Generate workstream map
@sdp/prompts/structured/phase-1-analyze.md
# Input: "Refactor orchestrator.py (350 LOC) for AI-Readiness"
# Reads: L1 (SYSTEM_OVERVIEW), L2 (domains/execution/), L4 (workstreams/INDEX.md)
# Output: 4 workstreams (structure, extract commands, orchestrator, wrapper)

# Phase 2: Plan first WS
@sdp/prompts/structured/phase-2-design.md
# Input: WS-001 from map
# Checks: L1-L4 context, duplicates in INDEX.md
# Output: Detailed plan with code snippets, bash checks

# Phase 3: Execute
@sdp/prompts/structured/phase-3-implement.md
# Follows plan, runs tests, checks coverage

# Phase 4: Review
@sdp/prompts/structured/phase-4-review.md
# Validates AI-Readiness, updates L2/L3 docs if needed, appends results to WS-001 file
```

### Example 2: New Feature (Clean Architecture)

```bash
# Create feature spec
vim tools/hw_checker/docs/specs/feature_24_obsidian_vault/feature.md

# Phase 1: Decompose into layers
@sdp/prompts/structured/phase-1-analyze.md
# Reads: PORTAL_VISION, RELEASE_PLAN, feature_24/feature.md
# Output: WS-140 (Domain), WS-141 (Application), WS-142 (Infrastructure), WS-143 (Presentation)

# Creates L2 domain map (if new domain)
vim tools/hw_checker/docs/domains/content/DOMAIN_MAP.md

# Phase 2-4: Execute each WS
# Each WS validated for layer boundaries before proceeding
# Phase 4 updates L2/L3 docs as needed
```

## Validation

```bash
# Validate workstream structure
python sdp/scripts/validate.py tools/hw_checker/docs/workstreams/WS-XXX

# Expected output:
# ✓ JSON Syntax Check
# ✓ Schema Validation
# ✓ Phase Requirements
# ✓ All checks passed
```

## Metrics

Track these for continuous improvement:

| Metric | Target | Command |
|--------|--------|---------|
| WS completion scope | SMALL/MEDIUM/LARGE | Manual tracking (LOC/tokens) |
| Coverage | ≥ 80% | `pytest --cov-fail-under=80` |
| Regression pass rate | 100% | `pytest -m fast` |
| Type checking | 0 errors | `mypy --strict` |
| Code quality | CC < 10 | `ruff check --select=C901` |

## Migration Notes

**From v1.2 (multi-agent) to v2.0 (structured):**
- Old epics (EP01-EP10) remain valid, no migration needed
- New work uses structured mode (Phase 1-4)
- Multi-agent prompts archived in `archive/v1.2/` for reference
- All v1.2 documentation in `archive/migration/`

## Resources

- **Full Protocol:** [PROTOCOL.md](PROTOCOL.md)
- **Code Patterns:** [HW_CHECKER_PATTERNS.md](HW_CHECKER_PATTERNS.md)
- **Common Rules:** [RULES_COMMON.md](RULES_COMMON.md)
- **Prompts:** [prompts/README.md](prompts/README.md)

## Support

For questions or issues:
1. Check `PROTOCOL.md` for detailed specifications
2. Review `HW_CHECKER_PATTERNS.md` for code examples
3. Validate your structure with `scripts/validate.py`
4. Review archived v1.2 materials in `archive/` if needed

---

**Protocol Version:** v0.3.0  
**Last Updated:** 2026-01-11  
**Status:** ✅ Active
