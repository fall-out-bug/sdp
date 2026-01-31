# SDP A+ Improvement Plan

> **Status:** Research complete
> **Date:** 2026-01-29
> **Goal:** Elevate SDP from B+ to A+ grade through comprehensive architectural, documentation, and quality improvements

---

## Table of Contents

1. [Overview](#overview)
2. [1. Protocol Clarity & Completeness](#1-protocol-clarity--completeness)
3. [2. Developer Experience (DX)](#2-developer-experience-dx)
4. [3. Validation & Enforcement](#3-validation--enforcement)
5. [4. Documentation Structure](#4-documentation-structure)
6. [5. Skill System Design](#5-skill-system-design)
7. [6. TDD Implementation](#6-tdd-implementation)
8. [7. Real-World Integration](#7-real-world-integration)
9. [8. Testing & Quality](#8-testing--quality)
10. [9. Artifact Quality](#9-artifact-quality)
11. [10. Tooling & Automation](#10-tooling--automation)
12. [Implementation Plan](#implementation-plan)
13. [Success Metrics](#success-metrics)

---

## Overview

### Goals

1. **Eliminate Ambiguity** — Make protocol impossible to misunderstand through canonical definitions and consistent terminology
2. **Improve Developer Experience** — Reduce onboarding friction from 2 hours to <30 minutes with interactive setup and health checks
3. **Strengthen Enforcement** — Bridge gap between documented rules and actual validation (coverage ≥80% must really mean ≥80%)
4. **Reorganize Documentation** — Role-based information architecture (newcomer/contributor/maintainer views)
5. **Enable Composability** — Refactor skill system into dependency-injected pipeline stages
6. **Enforce TDD Discipline** — Unbreakable Red→Green→Refactor enforcement via contract immutability + git timestamps
7. **Production Readiness** — CI/CD pipeline templates for team-wide quality enforcement
8. **Meta-Testing** — Dogfood SDP on itself (test the tools that enforce quality on others)
9. **Standardize Artifacts** — Flexible core schema with optional extensions for workstreams and specs
10. **Automate Tooling** — Self-service diagnostics, enhanced error messages, automation over manual processes

### Key Decisions

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| **Protocol Clarity** | Create canonical GLOSSARY.md + standardize workstream IDs | Eliminates confusion from 15+ ID formats |
| **DX** | Interactive setup wizard + health check command | Reduces first-action friction (Nir Eyal's habit formation) |
| **Validation** | Comprehensive enforcement suite (Option A) | Type-safe contracts align with SDP's principles |
| **Documentation** | Role-based reorganization (Option B) | Progressive disclosure matches how users learn |
| **Skill System** | Dependency injection + execution pipeline | Testability + composability without full rewrite |
| **TDD** | Hybrid enforcement (contract + git timestamps) | Unbreakable yet painless when followed correctly |
| **CI/CD** | GitHub Actions templates + Strangler adoption | Production-ready without big-bang migration |
| **Meta-Testing** | Extract shell hooks to Python + meta-quality gates | Dogfooding principle |
| **Artifacts** | Flexible core + optional extensions | Balances consistency with flexibility |
| **Tooling** | `sdp doctor` + enhanced error messages + setup wizard | Self-service over support burden |

---

## 1. Protocol Clarity & Completeness

> **Experts:** Dianna H. Ohnsorg (Technical Documentation), Daniel Kahneman (Cognitive Ease), Alberto Savoia (Actionable Specificity)

### Current State

**Critical Issues Found:**
- **15+ workstream ID formats**: `WS-001-01`, `sdp-118.22`, `PP-FFF-WW`, `00-012-01`, `beads-sdp-XXX`
- **Command prefix confusion**: `@feature` vs `/feature` used interchangeably across docs
- **Mixed language templates**: TEMPLATE.md has Russian headers without English translations
- **Undefined concepts**: "oneshot_ready", "AI-Comm", "veto triggers" mentioned but never defined
- **Grey area quality gates**: "Files < 200 LOC" — includes tests? comments? docstrings?

### Solution: Protocol Definition Overhaul

**Phase 1: Canonical Glossary (Week 1)**

Create `docs/GLOSSARY.md` with 150+ terms:

```markdown
## Workstream
**Definition:** Atomic unit of work in SDP, sized 500-1500 LOC (MEDIUM) or <500 LOC (SMALL)
**Aliases:** WS, task, sub-task (deprecated)
**Identifier Format:** PP-FFF-SS (e.g., 00-012-01)
  - PP: Project ID (00-99)
  - FF: Feature ID (001-999)
  - SS: Sequence number (01-99)
**Related:** Feature (parent), Release (ancestor)
**See:** /docs/workstreams/TEMPLATE.md
```

**Phase 2: Standardize Identifiers (Week 2)**

- **Choose ONE format**: `PP-FFF-SS` (from TEMPLATE.md)
- **Deprecate all others**: `WS-XXX-01`, `sdp-118.22`, `WS-XXX.YY`
- **Migration guide**: `docs/migration/workstream-id-standardization.md`
- **Automated migration script**: `scripts/migrate_workstream_ids.py --dry-run`

**Phase 3: Command Prefix Convention (Week 2)**

- **Skills**: `@feature`, `@idea`, `@design`, `@build`, `@oneshot`, `@review`, `@deploy`
- **Utilities**: `/debug`, `/issue`, `/tdd`
- **Update ALL documentation** to match this convention

**Phase 4: Translate TEMPLATE.md (Week 3)**

Replace Russian headers with English:
- `Цель` → `Goal`
- `Зависимость` → `Dependencies`
- `Входные файлы` → `Input Files`
- `Шаги` → `Steps`
- `Код` → `Code`
- `Ожидаемый результат` → `Expected Result`
- `Критерий завершения` → `Completion Criteria`
- `Ограничения` → `Constraints`

**Phase 5: Binary Rule Clarifications (Week 4)**

```markdown
## Quality Gates: Binary Criteria

### File Size (< 200 LOC)
- **Measure:** `wc -l file.py` (includes blank lines, docstrings, comments)
- **Scope:** All .py files in src/ (excludes tests/)
- **Check:** `find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200'`
- **Exception:** Generated code must be in generated/ subdirectory

### Coverage (≥ 80%)
- **Measure:** pytest --cov (branch coverage)
- **Scope:** All touched files in src/
- **Check:** `pytest --cov=src/ --cov-branch --cov-fail-under=80`
- **Exception:** None

### No TODO Rule
- **Forbidden:** `TODO:`, `FIXME:`, `HACK:`, `XXX:`
- **Allowed:** `TODO(author):` with WS-ID reference
- **Example:** `TODO(@jane): Refactor after WS-123 completes`
- **Enforcement:** `grep -rn "TODO\|FIXME\|HACK\|XXX" src/` must return empty
```

**Phase 6: Protocol Layers (Week 5)**

Add progressive disclosure to `PROTOCOL.md`:

```markdown
## Protocol Layers

### Layer 1: Essentials (First 30 minutes)
- Commands: @feature, @design, @build, @review
- Concepts: Workstream, Feature, TDD
- Quality gates: Coverage ≥80%, mypy --strict

### Layer 2: Collaboration (First day)
- Commands: @oneshot, @deploy, /debug
- Concepts: Beads, dependencies, parallel execution
- Advanced: Multi-agent coordination, checkpoints

### Layer 3: Optimization (First week)
- Commands: @issue, @hotfix, @bugfix
- Concepts: Execution modes, audit logging
- Advanced: Custom agents, notifications
```

### Success Metrics

- **New developer onboarding time**: <30 minutes to first successful `@build` (currently ~2 hours)
- **Protocol ambiguity questions**: Zero in community chat (currently ~5/week)
- **Workstream ID inconsistencies**: Zero across all files (currently 15+ formats)
- **Documentation search time**: <10 seconds to find any definition (currently undefined)

---

## 2. Developer Experience (DX)

> **Experts:** Nir Eyal (UX/Product), Dan Abramov (React/State), Kelsey Hightower (DevOps)

### Current Friction Points

1. **Installation ambiguity** — 3 methods (pipx, pip from source, submodule) with no guidance on WHEN to use which
2. **Scattered docs** — Core concepts across 4+ files; no clear "Start here" path
3. **Optional dependency hell** — Beads, GitHub CLI, Telegram all optional but failure messages unclear
4. **No validation** — Users don't know if setup works until mid-workflow errors
5. **Hard error recovery** — Hook failures lack remediation guidance
6. **User-specific paths** — `.claude/settings.json` has `/home/fall_out_bug/projects/sdp/...`
7. **Version confusion** — README says "Python 3.14+" but pyproject.toml says "^3.10"
8. **No quick feedback** — Cannot verify setup without running full workflow

### Solution: Interactive Setup + Health Check

**Phase 1: `sdp doctor` Command (Week 1)**

```bash
$ sdp doctor
✅ Python 3.14.2 (required: 3.10+)
✅ Poetry 1.8.0 (required: 1.7+)
✅ Git hooks installed (pre-commit, pre-push, post-build)
⚠️  Beads CLI not found (optional, required for task tracking)
   💡 Install: brew install go && go install github.com/steveyegge/beads/cmd/bd@latest
   📖 Or skip: export BEADS_USE_MOCK=true
✅ GitHub CLI 2.50.0 (optional, required for @deploy)
⚠️  Telegram bot not configured (optional, required for notifications)
   💡 Run: sdp config telegram --setup

Overall: 5/7 checks passed
```

**Phase 2: `sdp init` Wizard (Week 2)**

```bash
$ sdp init
🚀 SDP Project Initialization

Step 1/5: Project metadata
? Project name: my-awesome-project
? Project ID (00-99): 01
? Description: A project that does awesome things

Step 2/5: Optional dependencies
? Install Beads CLI for task tracking? [y/N]: y
✓ Beads CLI detected
? Configure GitHub integration? [y/N]: y
✓ GitHub CLI detected
? Configure Telegram notifications? [y/N]: N
ⓘ Skipping Telegram (optional)

Step 3/5: Quality gate configuration
? Coverage threshold [80]:
? Max file size (LOC) [200]:
? Type checking mode [strict|mypy|none]: strict

Step 4/5: Git hooks installation
✓ Installing pre-commit hook...
✓ Installing pre-push hook...
✓ Installing post-build hook...
✓ Installing pre-deploy hook...

Step 5/5: Validation
Running `sdp doctor`...
✅ All checks passed!

✅ Setup complete! Next steps:
  1. Run: @feature "Add user authentication"
  2. Or: cp docs/workstreams/TEMPLATE.md docs/workstreams/backlog/01-001-01.md
```

**Phase 3: Enhanced Error Messages (Week 3)**

Before:
```
Error: Beads not found
```

After:
```
❌ Beads CLI not found

⚠️  Beads is optional for SDP, but required for task tracking

💡 Install:
   brew install go
   go install github.com/steveyegge/beads/cmd/bd@latest

📖 Or skip:
   export BEADS_USE_MOCK=true

🆘 Help: https://sdp.dev/docs/beads.html
📋 Chat: https://discord.gg/sdp-dev
```

**Phase 4: Interactive Tutorial (Week 4)**

```bash
$ sdp tutorial
🎓 SDP Interactive Tutorial

This tutorial will guide you through your first feature using SDP.
Estimated time: 15 minutes

Step 1: Create a feature (2 min)
We'll create a simple "Hello World" feature.

Run: @feature "Hello World feature"

✓ Created docs/drafts/idea-hello-world.md
✓ Created docs/intent/hello-world.json

Step 2: Design workstreams (3 min)
Now let's break this feature into workstreams.

Run: @design idea-hello-world

✓ Created 2 workstreams:
  - 01-001-01: Create hello() function
  - 01-001-02: Add CLI command

Step 3: Implement first workstream (5 min)
Let's implement the hello() function using TDD.

Run: @build 01-001-01

[TodoWrite tracking shown here...]

✓ Tests pass (coverage 100%)
✓ Type checks pass (mypy --strict)
✓ Workstream complete!

Step 4: Review and deploy (5 min)
Run: @review hello-world
Run: @deploy hello-world

✅ Tutorial complete! You've built your first feature with SDP.

Next: Try @feature with your own idea, or see docs/START_HERE.md
```

### Success Metrics

- **Time to first `@build`**: <15 minutes from clone (currently ~2 hours)
- **Setup-related issues**: <5% of total support requests (currently ~30%)
- **Tutorial completion rate**: >80% (currently untracked)
- **`sdp doctor` usage**: Track frequency (indicates self-service troubleshooting)

---

## 3. Validation & Enforcement

> **Experts:** Theo Browne (API Design), Troy Hunt (Security), Kelsey Hightower (DevOps)

### Current Gaps

| Documented Rule | Actual Enforcement | Gap |
|-----------------|-------------------|-----|
| Coverage ≥80% | Runs only if test exists | **Soft fail** |
| Files <200 LOC | Warns only | **Soft fail** |
| CC <10 | No tool configured | **Not enforced** |
| No `except: pass` | Warns only | **Soft fail** |
| Clean architecture | Hardcoded to hw_checker | **Not portable** |
| CI/CD quality gates | None | **Missing** |

### Solution: Comprehensive Enforcement Suite

**Phase 1: Declarative Quality Gates (Week 1)**

Create `quality-gate.toml`:

```toml
[coverage]
min_percent = 80
fail_under = true
branch_coverage = true
exclude = ["tests/", "generated/"]

[complexity]
max_cc = 10
tool = "radon"
fail_over = true

[file_size]
max_lines = 200
fail_over = true
exclude = ["*.pb.py", "*_pb2.py"]

[type_hints]
mode = "strict"
fail_under = true

[error_handling]
forbidden_patterns = ["except:", "except Exception", "except BaseException"]
allow_commented = true  # "except: pass  # noqa: E722" allowed
fail_under = true

[architecture]
enforce_layer_boundaries = true
domain_imports_allowed = false
```

**Phase 2: Strengthen Hooks (Week 2)**

Update all hooks to:

1. **Read from `quality-gate.toml`** instead of hardcoded values
2. **Fail instead of warn** on violations
3. **Provide actionable error messages**

Example transformation:

```bash
# Before (hooks/post-build.sh:82-84)
if [ "$lines" -gt 200 ]; then
    echo "⚠️  Warning: $file exceeds 200 lines ($lines)"
fi

# After
if [ "$lines" -gt 200 ]; then
    echo "❌ File size violation: $file has $lines lines (max: $(config_get file_size.max_lines))"
    echo "💡 Refactor into smaller modules or move to generated/ if auto-generated"
    exit 1
fi
```

**Phase 3: Add Complexity Checking (Week 2)**

Add to `pyproject.toml`:

```toml
[tool.poetry.dependencies]
radon = "^6.0"
```

Update `hooks/post-build.sh`:

```bash
# Check cyclomatic complexity
echo "🔍 Checking complexity..."
radon cc src/ -a -s --json > /tmp/radon.json
complexity=$(jq '[.[] | .average] | add / length' /tmp/radon.json)

if (( $(echo "$complexity > $(config_get complexity.max_cc)" | bc -l) )); then
    echo "❌ Complexity violation: average CC is $complexity (max: $(config_get complexity.max_cc))"
    echo "💡 Run 'radon cc src/ -a' to see detailed report"
    exit 1
fi
```

**Phase 4: GitHub Actions CI/CD (Week 3)**

Create `.github/workflows/quality-gate.yml`:

```yaml
name: SDP Quality Gates

on:
  pull_request:
    branches: [main, dev]
  push:
    branches: [main, dev]

jobs:
  quality-gate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.14'

      - name: Install dependencies
        run: |
          pip install poetry
          poetry install --with dev

      - name: Run quality gates
        run: |
          poetry run pytest --cov=src/ --cov-branch --cov-fail-under=80
          poetry run mypy src/ --strict
          poetry run ruff check src/
          poetry run radon cc src/ -a --json | jq -e 'all(.[] | .average <= 10)'

      - name: Enforce file size
        run: |
          find src/ -name "*.py" -exec wc -l {} + | \
          awk '$1 > 200 {exit 1}' || \
          (echo "❌ Files exceed 200 LOC" && exit 1)

      - name: Check for bare except
        run: |
          ! grep -rn "except:\|except Exception\|except BaseException" src/
```

**Phase 5: Make Architecture Checks Portable (Week 4)**

Replace hardcoded hw_checker checks with regex-based validation:

```python
# src/sdp/hooks/architecture_checker.py
import re
from pathlib import Path
from typing import Dict, List

class ArchitectureChecker:
    def __init__(self, config_path: Path):
        self.rules = self._load_rules(config_path)

    def check_layer_boundaries(self, files: List[str]) -> List[str]:
        """Check that files respect clean architecture boundaries."""
        violations = []

        for file in files:
            module = self._parse_module_path(file)
            layer = self._detect_layer(module)

            imports = self._extract_imports(file)
            for imp in imports:
                imp_layer = self._detect_layer(imp)
                if not self._is_allowed_import(layer, imp_layer):
                    violations.append(
                        f"{file}: Cannot import {imp} ({layer} → {imp_layer})"
                    )

        return violations

    def _is_allowed_import(self, from_layer: str, to_layer: str) -> bool:
        """Check if import direction follows dependency rule."""
        # Domain ← Application ← Infrastructure ← Presentation
        hierarchy = {
            "domain": 0,
            "application": 1,
            "infrastructure": 2,
            "presentation": 3,
        }

        return hierarchy.get(from_layer, 99) <= hierarchy.get(to_layer, 99)
```

**Phase 6: Warm-Up Period (Week 5)**

Run in "warn-only" mode for 1 sprint before flipping to enforcement:

```bash
# In quality-gate.toml
[general]
enforcement_mode = "warn"  # "warn" | "enforce"
```

### Success Metrics

- **Quality gate pass rate**: <5% failures in CI (indicating gates are realistic)
- **Pre-commit hook runtime**: <10 seconds (acceptable friction)
- **Code quality improvements**: CC reduction 10%, coverage increase 5%
- **Team adoption**: All PRs pass quality gates before merge

---

## 4. Documentation Structure

> **Experts:** Dianna Hutts Quammen (Technical Storytelling), Kathy Sierra (Creating Passionate Users), Edward Tufte (Information Design)

### Current Issues

- **No central index** — 100+ markdown files with no navigation hub
- **Scattered information** — Similar content across multiple files
- **124 TODO/FIXME markers** — Incomplete sections
- **Inconsistent naming** — TUTORIAL.md in docs/ but guides/ folder exists
- **Buried workstream INDEX** — Hard to find reference

### Solution: Role-Based Information Architecture

**Phase 1: Create Documentation Hub (Week 1)**

Create `docs/START_HERE.md` and `docs/SITEMAP.md`:

```markdown
# START_HERE.md

# Welcome to SDP

**New to SDP?** Start with the [15-minute tutorial](beginner/00-quick-start.md)

**Need command reference?** See [reference/commands/](reference/commands/)

**Contributing or maintaining?** Read [internals/contributing.md](internals/contributing.md)

---

## Quick Links

- **Installation**: `pipx install sdp` or see [beginner/00-quick-start.md](beginner/00-quick-start.md#installation)
- **First feature**: [beginner/01-first-feature.md](beginner/01-first-feature.md)
- **Quality gates**: [reference/quality-gates.md](reference/quality-gates.md)
- **Troubleshooting**: [beginner/03-troubleshooting.md](beginner/03-troubleshooting.md)

---

## Table of Contents

See [SITEMAP.md](SITEMAP.md) for complete documentation index.
```

**Phase 2: Reorganize by Role (Week 2-3)**

```
docs/
├── START_HERE.md          # Newcomer landing page
├── SITEMAP.md             # Full index with role filters
│
├── beginner/              # Progressive learning path
│   ├── 00-quick-start.md     # 15-minute tutorial (move TUTORIAL.md)
│   ├── 01-first-feature.md   # Walk through @feature → @build
│   ├── 02-common-patterns.md # Code patterns, principles
│   └── 03-troubleshooting.md # FAQ, runbooks
│
├── reference/             # Lookup documentation
│   ├── commands/              # All @ commands
│   ├── concepts/              # Clean architecture, SOLID, TDD
│   ├── quality-gates.md       # Coverage, mypy, ruff
│   └── workflow.md            # Complete workflow diagram
│
├── internals/             # Maintainer documentation
│   ├── architecture/         # ADRs, design docs
│   ├── agents/               # Agent system, roles
│   ├── extending/            # Adding commands, skills
│   └── contributing.md       # Move CONTRIBUTING.md here
│
└── runbooks/              # Operational procedures (keep as-is)
```

**Phase 3: Resolve Documentation Debt (Week 4)**

Create workstream **WS-DOC-01: Resolve Documentation TODOs**:

- Audit all 124 TODO/FIXME markers in docs/
- Either complete the section or remove the TODO
- Estimate: 40 TODOs to complete, 84 to remove

**Phase 4: Add Link Validation (Week 5)**

Add to pre-commit hook:

```bash
# Check markdown links
markdown-link-check docs/ \
  --config .markdown-link-check.json \
  || (echo "❌ Broken documentation links" && exit 1)
```

### Success Metrics

- **Time to find information**: <10 seconds for any concept (currently untracked)
- **Documentation TODOs**: Zero in production docs (currently 124)
- **Broken links**: Zero (currently ~15 detected)
- **Newcomer satisfaction**: Survey after onboarding (target: >4/5)

---

## 5. Skill System Design

> **Experts:** Martin Fowler (Refactoring), Robert C. Martin (SRP), Eric Evans (Ubiquitous Language)

### Current Architecture Issues

| Issue | Impact |
|-------|--------|
| **Tight coupling** | `@build` directly calls `/tdd` Skill tool |
| **Error recovery** | Exception propagation loses context |
| **Testability** | Hard to mock BeadsClient |
| **Composability** | Skills call other skills (implicit dependencies) |
| **Checkpointing** | Manual JSON file writes in `@oneshot` |

### Solution: Dependency Injection + Execution Pipeline

**Phase 1: Extract BuildSkill (Week 1-2)**

```python
# src/sdp/skills/build.py
from dataclasses import dataclass
from sdp.beads.client import BeadsClient
from sdp.tdd.runner import TDDRunner
from sdp.validation.validator import WorkstreamValidator

@dataclass
class BuildSkill:
    """Orchestrates workstream execution with TDD discipline."""

    beads_client: BeadsClient
    tdd_runner: TDDRunner
    validator: WorkstreamValidator

    def execute(self, workstream_id: str) -> ExecutionResult:
        """Execute workstream through TDD pipeline."""
        # Validate prerequisites
        validation = self.validator.validate(workstream_id)
        if not validation.passed:
            return ExecutionResult.failed(validation.errors)

        # Execute TDD phases
        tdd_result = self.tdd_runner.run(workstream_id)

        # Verify quality gates
        quality_result = self.validator.verify_quality(workstream_id)

        return ExecutionResult(
            success=tdd_result.success and quality_result.passed,
            artifacts=tdd_result.artifacts,
            quality_report=quality_result,
        )

@dataclass
class ExecutionContext:
    """Shared context across pipeline stages."""

    workstream_id: str
    state: SkillState
    artifacts: List[Artifact]
    error: Optional[Exception]

class SkillState(Enum):
    IDLE = "idle"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"
```

**Phase 2: Extract IdeaSkill, DesignSkill (Week 3-4)**

```python
# src/sdp/skills/idea.py
@dataclass
class IdeaSkill:
    """Gathers requirements through structured interviewing."""

    interviewer: IdeaInterviewer
    ambiguity_detector: AmbiguityDetector
    beads_client: BeadsClient

    def execute(self, feature_description: str) -> InterviewResult:
        """Run ambiguity detection and interview process."""
        # Detect ambiguities
        ambiguities = self.ambiguity_detector.analyze(feature_description)

        # Conduct interview (interactive via AskUserQuestion)
        interview_data = self.interviewer.conduct(feature_description, ambiguities)

        # Generate artifacts
        intent = self._generate_intent_json(interview_data)
        draft = self._generate_draft_markdown(interview_data)

        # Store in Beads
        self.beads_client.create_task(feature_description, intent, draft)

        return InterviewResult(intent=intent, draft=draft)

# src/sdp/skills/design.py
@dataclass
class DesignSkill:
    """Decomposes features into workstreams with dependency graph."""

    decomposer: FeatureDecomposer
    graph_builder: DependencyGraphBuilder
    beads_client: BeadsClient

    def execute(self, feature_id: str) -> DesignResult:
        """Break feature into atomic workstreams."""
        # Load feature spec
        feature = self.beads_client.get_feature(feature_id)

        # Decompose into workstreams
        workstreams = self.decomposer.decompose(feature)

        # Build dependency graph
        graph = self.graph_builder.build(workstreams)

        # Validate graph (no cycles, reasonable depth)
        self._validate_graph(graph)

        # Create workstream files
        for ws in workstreams:
            self._create_workstream_file(ws)

        return DesignResult(workstreams=workstreams, graph=graph)
```

**Phase 3: Pipeline Executor (Week 5)**

```python
# src/sdp/skills/pipeline.py
from typing import List, Callable

class PipelineExecutor:
    """Executes composable skill pipelines."""

    def __init__(self, context: ExecutionContext):
        self.context = context
        self.stages: List[PipelineStage] = []

    def add_stage(self, stage: PipelineStage) -> "PipelineExecutor":
        """Add stage to pipeline (fluent API)."""
        self.stages.append(stage)
        return self

    def execute(self) -> ExecutionResult:
        """Execute all stages in order."""
        for stage in self.stages:
            self.context.state = SkillState.RUNNING

            try:
                result = stage.execute(self.context)
                if hasattr(stage, "on_success"):
                    stage.on_success(self.context, result)
            except Exception as e:
                self.context.state = SkillState.FAILED
                self.context.error = e
                if hasattr(stage, "on_failure"):
                    stage.on_failure(self.context, e)
                return ExecutionResult.failed(e)

        self.context.state = SkillState.COMPLETED
        return ExecutionResult.success(self.context.artifacts)

# Example: @oneshot as pipeline
def create_oneshot_pipeline(feature_id: str) -> PipelineExecutor:
    """Create multi-agent execution pipeline."""
    context = ExecutionContext(workstream_id=feature_id, state=SkillState.IDLE)

    return (
        PipelineExecutor(context)
        .add_stage(DesignSkill(decomposer, graph_builder, beads))
        .add_stage(BuildSkill(beads, tdd_runner, validator))
        .add_stage(BuildSkill(beads, tdd_runner, validator))  # WS 2
        .add_stage(ReviewSkill(quality_checker))
        .add_stage(DeploySkill(deployer))
    )
```

**Phase 4: Hooks & Event Emission (Week 6)**

```python
# src/sdp/skills/events.py
class EventEmitter:
    """Emits events for audit logging and monitoring."""

    def __init__(self, transport: EventTransport):
        self.transport = transport

    def emit(self, event: Event):
        """Emit event to transport (file, webhook, stdout)."""
        self.transport.send(event)

# Usage in skills
@dataclass
class BuildSkill:
    event_emitter: EventEmitter

    def execute(self, workstream_id: str) -> ExecutionResult:
        self.event_emitter.emit(Event(
            type="build_started",
            workstream_id=workstream_id,
            timestamp=datetime.now(),
        ))

        try:
            result = self._execute_build(workstream_id)
            self.event_emitter.emit(Event(type="build_completed", ...))
            return result
        except Exception as e:
            self.event_emitter.emit(Event(type="build_failed", error=str(e)))
            raise
```

### Success Metrics

- **Skill testability**: >80% unit test coverage for skill classes
- **Pipeline flexibility**: Can create custom pipelines by composing stages
- **Error recovery**: All failures captured in ExecutionContext
- **Checkpointing**: Automatic state tracking, no manual JSON writes

---

## 6. TDD Implementation

> **Experts:** Kent C. Dodds (Testing Behavior), Theo Browne (Type-Safe Contracts), Sam Newman (Bounded Contexts)

### Current Gaps

1. **TDDRunner exists but is not used** — `mock_tdd_success` parameter bypasses enforcement
2. **No timestamp verification** — Can't prove tests were written before implementation
3. **Trust-based enforcement** — Relies on developer/agent honesty
4. **Contract tests not immutable** — Can modify tests during implementation phase

### Solution: Hybrid TDD Enforcement

**Phase 1: Contract Immutability (Week 1-2)**

```python
# src/sdp/tdd/contract.py
from pathlib import Path

class TestContract:
    """Makes test files read-only during implementation phase."""

    def __init__(self, workstream_id: str):
        self.workstream_id = workstream_id
        self.marker_file = Path(f".tdd-contract-{workstream_id}")

    def lock_tests(self, test_files: List[Path]):
        """Lock test files (make read-only) during implementation."""
        self.marker_file.write_text("locked")  # Signal locked state

        for test_file in test_files:
            # Make read-only (Unix: chmod 444)
            test_file.chmod(0o444)

            # Add contract marker
            content = test_file.read_text()
            header = "# TDD-CONTRACT: DO NOT MODIFY DURING BUILD PHASE\n"
            test_file.write_text(header + content)

    def unlock_tests(self):
        """Unlock tests for refactoring phase."""
        self.marker_file.write_text("unlocked")

        # Find all contract files
        for test_file in Path("tests/").rglob("test_*.py"):
            if "# TDD-CONTRACT" in test_file.read_text():
                test_file.chmod(0o644)  # Make writable
                # Remove marker
                content = test_file.read_text()
                content = content.replace("# TDD-CONTRACT: DO NOT MODIFY DURING BUILD PHASE\n", "")
                test_file.write_text(content)

# Integration with TDDRunner
class TDDRunner:
    def run_green_phase(self, ws_id: str):
        """Run implementation phase with locked tests."""
        contract = TestContract(ws_id)

        # Find test files
        test_files = list(Path("tests/").rglob(f"test_*{ws_id}*.py"))

        # Lock them
        contract.lock_tests(test_files)

        try:
            # Run implementation
            self._run_implementation(ws_id)
        finally:
            # Unlock after completion
            contract.unlock_tests()
```

**Phase 2: Git Timestamp Verification (Week 3)**

```bash
# hooks/post-build.sh
verify_tdd_timestamps() {
    local ws_id=$1
    local violations=0

    # Find implementation files for this workstream
    for impl_file in $(find src/ -name "*.py" -newer tests/); do
        # Check if corresponding test file exists and is older
        local test_file="tests/test_$(basename $impl_file)"

        if [ -f "$test_file" ]; then
            local impl_time=$(git log --diff-order=A --follow --format=%ct -1 "$impl_file")
            local test_time=$(git log --diff-order=A --follow --format=%ct -1 "$test_file")

            if [ "$impl_time" -lt "$test_time" ]; then
                echo "❌ TDD violation: $impl_file created BEFORE $test_file"
                violations=$((violations + 1))
            fi
        fi
    done

    if [ $violations -gt 0 ]; then
        echo "❌ Found $violations TDD violations"
        echo "💡 Tests must be written BEFORE implementation"
        exit 1
    fi
}
```

**Phase 3: Pre-commit Phase Enforcement (Week 4)**

```python
# hooks/pre-commit.py
import sys
from pathlib import Path

def check_staged_files_phase(ws_id: str):
    """Ensure only phase-appropriate files are staged."""
    staged_files = get_staged_files()

    # Check current phase
    contract_file = Path(f".tdd-contract-{ws_id}")
    if not contract_file.exists():
        return  # No contract, no enforcement

    phase = contract_file.read_text().strip()  # "red" | "green" | "refactor"

    if phase == "red":
        # Only test files allowed
        for f in staged_files:
            if not f.startswith("tests/"):
                print(f"❌ Cannot commit {f} during RED phase")
                print("💡 Only test files can be committed in RED phase")
                sys.exit(1)

    elif phase == "green":
        # Implementation files allowed, tests locked
        for f in staged_files:
            if f.startswith("tests/") and is_contract_test(f):
                print(f"❌ Cannot modify {f} during GREEN phase (contract locked)")
                print("💡 Tests are immutable during implementation")
                sys.exit(1)

    elif phase == "refactor":
        # Both allowed, but verify tests still pass
        pass
```

**Phase 4: Enhanced TDDRunner (Week 5)**

```python
# src/sdp/tdd/runner.py
class TDDRunner:
    def __init__(self, ws_id: str):
        self.ws_id = ws_id
        self.contract = TestContract(ws_id)
        self.phase = "red"

    def run_red_phase(self):
        """Write failing tests."""
        print("🔴 RED PHASE: Writing failing tests")

        # Set phase
        self.phase = "red"
        self._set_contract_phase("red")

        # Create test file
        test_file = self._create_test_file()

        # Run tests (should fail)
        result = subprocess.run(["pytest", test_file, "--tb=short"])
        if result.returncode == 0:
            raise TDDViolation("Tests must fail initially in RED phase")

        print("✅ RED phase complete: Tests fail as expected")
        return test_file

    def run_green_phase(self, test_file: Path):
        """Write minimal implementation."""
        print("🟢 GREEN PHASE: Implementing to pass tests")

        # Set phase
        self.phase = "green"
        self._set_contract_phase("green")

        # Lock tests
        self.contract.lock_tests([test_file])

        try:
            # Write implementation
            impl_file = self._write_implementation()

            # Run tests (should pass)
            result = subprocess.run(["pytest", test_file, "--tb=short"])
            if result.returncode != 0:
                raise TDDViolation("Tests must pass in GREEN phase")

            print("✅ GREEN phase complete: Tests pass")
            return impl_file
        finally:
            self.contract.unlock_tests()

    def run_refactor_phase(self, impl_file: Path, test_file: Path):
        """Refactor implementation."""
        print("🔵 REFACTOR PHASE: Improving code quality")

        # Set phase
        self.phase = "refactor"
        self._set_contract_phase("refactor")

        # Unlock tests for refactoring
        # (tests can change if behavior changes)

        # Refactor
        self._refactor(impl_file)

        # Verify tests still pass
        result = subprocess.run(["pytest", test_file, "--tb=short"])
        if result.returncode != 0:
            raise TDDViolation("Tests must still pass after refactoring")

        # Run quality gates
        self._run_quality_gates(impl_file)

        print("✅ REFACTOR phase complete")
```

**Phase 5: Post-Build Verification (Week 6)**

```bash
# hooks/post-build.sh
verify_tdd_compliance() {
    local ws_id=$1

    echo "🔍 Verifying TDD compliance..."

    # Check 1: Git timestamps (tests before impl)
    verify_tdd_timestamps "$ws_id"

    # Check 2: Contract markers present
    if ! grep -q "# TDD-CONTRACT" tests/test_*${ws_id}*.py 2>/dev/null; then
        echo "⚠️  No TDD contract markers found (were tests created with /test?)"
    fi

    # Check 3: No test modifications after lock
    # (Git history check)

    echo "✅ TDD compliance verified"
}
```

### Success Metrics

- **TDD violation rate**: <2% of builds (enforcement works)
- **Test timestamp compliance**: 100% of new code follows test-first
- **Contract lock effectiveness**: Zero successful bypass attempts
- **Developer feedback**: Survey on TDD friction (target: "worth it" >80%)

---

## 7. Real-World Integration

> **Experts:** Kelsey Hightower (DevOps), Martin Fowler (Refactoring), Jez Humble (CI/CD)

### Current Gaps

- **No CI/CD pipeline examples** — Projects must figure this out themselves
- **Legacy code migration undocumented** — "Requires refactoring" but no how-to
- **No team collaboration patterns** — Beyond basic GitHub issues
- **No incremental adoption guide** — All-or-nothing approach
- **Time pressure handling missing** — Hotfix exists, but no triage/prioritization

### Solution: CI/CD Templates + Strangler Pattern

**Phase 1: GitHub Actions Template (Week 1)**

Create `.github/workflows/sdp-quality-gate.yml`:

```yaml
name: SDP Quality Gate

on:
  pull_request:
    branches: [main, dev]
  push:
    branches: [main, dev]

jobs:
  quality-gate:
    runs-on: ubuntu-latest
    timeout-minutes: 10

    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # For git history checks

      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.14'

      - name: Install SDP
        run: |
          pipx install sdp
          sdp doctor

      - name: Run SDP Quality Gates
        run: |
          # Coverage (branch ≥80%)
          pytest --cov=src/ --cov-branch --cov-fail-under=80

          # Type checking
          mypy src/ --strict

          # Linting
          ruff check src/

          # Complexity (radon CC <10)
          radon cc src/ -a --json | \
            jq -e 'all(.[] | .average <= 10)' || \
            (echo "❌ Complexity exceeds CC=10" && exit 1)

          # File size (<200 LOC)
          find src/ -name "*.py" -exec wc -l {} + | \
            awk '$1 > 200 {exit 1}' || \
            (echo "❌ Files exceed 200 LOC" && exit 1)

          # No bare except
          ! grep -rn "except:\|except Exception" src/ || \
            (echo "❌ Bare except clauses found" && exit 1)

          # Clean architecture
          python -m sdp.hooks.architecture_checker || \
            (echo "❌ Layer boundary violations" && exit 1)

      - name: Verify TDD Compliance
        run: |
          # Check test timestamps (tests before impl)
          hooks/post-build.sh verify_tdd_timestamps

      - name: Comment PR with Results
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: '✅ All SDP quality gates passed!\n\nCoverage: 85%\nComplexity: CC=8.2\nFile size: Max 187 LOC'
            })
```

**Phase 2: Strangler Pattern Adoption Guide (Week 2)**

Create `docs/guides/strangler-adoption.md`:

```markdown
# Incremental SDP Adoption (Strangler Pattern)

## Strategy

Don't rewrite everything. Use SDP for new features, gradually migrate legacy code.

## Phase 1: New Features Only (Week 1-4)

- All new features use `@feature` → `@design` → `@build`
- Legacy features continue existing process
- Team compares quality/velocity metrics

## Phase 2: High-Churn Modules (Month 2-3)

1. Identify high-churn legacy module (frequent changes)
2. Create facade layer around it
3. Build new features alongside facade using SDP
4. Migrate functionality piece-by-piece
5. Remove legacy code when facade fully replaces it

## Phase 3: Low-Risk Modules (Month 4-6)

- Apply same pattern to stable modules
- Lower urgency, but systematic migration

## Example: Strangling a User Service

### Before (Legacy)
```
src/
└── services/
    └── user_service.py  # 2000 LOC, no tests
```

### Step 1: Create Facade
```python
# src/services/user_facade.py (SDP workstream)
class UserServiceFacade:
    """Facade around legacy user service."""

    def __init__(self):
        self._legacy = LegacyUserService()  # Wrap legacy

    def get_user(self, user_id: int) -> User:
        # Delegate to legacy
        return self._legacy.get_user(user_id)

    def create_user(self, data: dict) -> User:
        # New behavior with TDD
        return self._legacy.create_user(data)
```

### Step 2: Migrate One Method
```python
# New implementation (SDP workstream)
class UserService:
    def __init__(self, repo: UserRepository):
        self._repo = repo

    def create_user(self, data: dict) -> User:
        # New TDD implementation
        user = User(**data)
        self._repo.save(user)
        return user

# Facade updates
class UserServiceFacade:
    def create_user(self, data: dict) -> User:
        # Use new implementation
        return self._new_service.create_user(data)

    def get_user(self, user_id: int) -> User:
        # Still delegate to legacy
        return self._legacy.get_user(user_id)
```

### Step 3: Repeat Until Legacy Gone
```python
# Final state (all methods migrated)
class UserServiceFacade:
    def __init__(self):
        self._service = UserService(repo)  # All new

    def get_user(self, user_id: int) -> User:
        return self._service.get_user(user_id)

    def create_user(self, data: dict) -> User:
        return self._service.create_user(data)

    # Delete LegacyUserService
```
```

**Phase 3: Team Coordination Playbook (Week 3)**

Create `docs/guides/team-collaboration.md`:

```markdown
# SDP for Teams

## Code Review Queues

Use GitHub Projects to track workstream review status:

- **To Review**: PRs awaiting review
- **In Review**: PRs under active review
- **Approved**: Ready to merge

## WIP Limits

Limit work-in-progress per developer:

- Max 2 workstreams in parallel
- Max 1 workstream per feature (reduce context switching)

## Distributed Team Handoffs

When handing off workstreams across timezones:

1. **Update workstream status**: Move to `docs/workstreams/in_progress/`
2. **Add handoff note**:
   ```markdown
   ## Handoff Notes

   - Status: Tests pass, implementation complete
   - Blocked by: Waiting for API design from @alice
   - Next steps: Complete refactoring when design approved
   - Context: See Slack #dev-api thread
   ```
3. **Sync Beads task**: Assign to new owner
4. **GitHub PR**: Mention in PR description

## Async Communication

- Use GitHub issues for questions (mention workstream ID)
- Use Beads task comments for progress updates
- Use PR comments for code review discussions
- Avoid real-time meetings for routine workstreams

## Conflict Resolution

When workstreams have merge conflicts:

1. Rebase local branch: `git pull --rebase`
2. Resolve conflicts
3. Run tests: `pytest`
4. Push: `git push`
5. Update workstream file with "Resolved merge conflicts"
```

**Phase 4: Migration Dashboard (Week 4)**

Create `scripts/migration_progress.py`:

```python
#!/usr/bin/env python3
"""Track legacy → SDP migration progress."""

import json
from pathlib import Path

def scan_repository():
    """Scan repo for SDP adoption metrics."""

    # Count modules using SDP (have workstreams)
    sdp_modules = set()
    for ws_file in Path("docs/workstreams/completed/").rglob("*.md"):
        # Extract module path from workstream
        ws_content = ws_file.read_text()
        if "## Input Files" in ws_content:
            for line in ws_content.split("## Input Files")[1].split("\n")[1:10]:
                if line.strip().startswith("src/"):
                    module = line.strip().split("/")[1]
                    sdp_modules.add(module)

    # Count all modules
    all_modules = set()
    for py_file in Path("src/").rglob("*.py"):
        module = py_file.relative_to("src/").parts[0]
        all_modules.add(module)

    # Calculate metrics
    total = len(all_modules)
    migrated = len(sdp_modules)
    percentage = (migrated / total * 100) if total > 0 else 0

    # Generate report
    report = {
        "total_modules": total,
        "sdp_modules": migrated,
        "legacy_modules": total - migrated,
        "migration_percentage": round(percentage, 1),
        "sdp_module_list": sorted(sdp_modules),
        "legacy_module_list": sorted(all_modules - sdp_modules),
    }

    # Print report
    print(json.dumps(report, indent=2))

    # Generate visualization
    print(f"\n📊 Migration Progress: {percentage:.1f}%")
    print(f"{'█' * int(percentage / 5)}{'░' * (20 - int(percentage / 5))}")
    print(f"{migrated}/{total} modules using SDP")

    return report

if __name__ == "__main__":
    scan_repository()
```

**Phase 5: Hotfix & Triage Process (Week 5)**

Create `docs/runbooks/hotfix-workflow.md`:

```markdown
# Emergency Hotfix Workflow

## When to Use

- **P0 (Critical)**: Production down, data loss, security breach
- **Use @hotfix skill** — fast-tracked, <2 hour resolution

## Workflow

```bash
# 1. Create hotfix from main
@hotfix "Fix database connection leak"

# 2. Skip design phase (goes straight to build)
# 3. Minimal tests (RED phase only critical path)
# 4. Deploy immediately after GREEN
# 5. Create retrospective workstream later
```

## Post-Hotfix

After production is stable, create retrospective workstream:

```bash
# Add tests for edge cases
# Refactor hotfix code (REFACTOR phase)
# Add monitoring/alerting
```

## Triage Process

For bugs that aren't P0:

```bash
# Classify severity
@issue "Login fails on Firefox"
# → Routes to @bugfix (P1/P2) or backlog (P3)
```

**Severity Levels:**
- **P0**: @hotfix — production down
- **P1**: @bugfix — critical functionality broken
- **P2**: @bugfix — important but workaround exists
- **P3**: backlog — nice to have
```

### Success Metrics

- **CI/CD adoption**: >80% of projects using SDP use GitHub Actions template
- **Incremental migration**: >50% of modules migrated within 6 months (no big bang)
- **Team velocity**: Measure features/month before and after SDP adoption
- **Hotfix frequency**: Track @hotfix usage (should decrease over time)

---

## 8. Testing & Quality

> **Experts:** Kent C. Dodds (Testing Behavior), Martin Fowler (Refactoring), Sam Newman (Bounded Contexts)

### Critical Gap: Hypocrisy

SDP enforces quality gates on projects using it, but **doesn't apply those same standards to its own tools**.

**Current State:**
- 104 Python modules in `/src/sdp/` (~608 classes/functions)
- README claims "91% coverage" but that's for hw_checker, not SDP
- **All 16 shell scripts** are untested (hooks/)
- **All GitHub integration** is untested (18 modules)
- **No integration tests** for complete workflows

### Solution: Meta-Testing + Shell Hook Extraction

**Phase 1: Extract Shell Hooks to Python (Week 1-4)**

```python
# src/sdp/hooks/pre_commit.py (extracted from hooks/pre-commit.sh)
import subprocess
from pathlib import Path
from typing import List

def check_time_estimates(files: List[Path]) -> CheckResult:
    """Check files for time-based estimates."""
    violations = []

    for file in files:
        content = file.read_text()

        # Forbidden patterns
        forbidden = [
            r"\d+ hours?",
            r"\d+ days?",
            r"\d+ weeks?",
            "soon",
            "ASAP",
            "TMTR",
        ]

        for pattern in forbidden:
            if re.search(pattern, content, re.IGNORECASE):
                violations.append(f"{file}: Found '{pattern}'")

    if violations:
        return CheckResult.failed(
            "Time estimates found",
            violations,
            remediation="Remove time-based estimates, use relative sizing (SMALL/MEDIUM/LARGE)"
        )

    return CheckResult.passed()

# Shell wrapper becomes thin
# hooks/pre-commit.sh
#!/bin/bash
python -m sdp.hooks.pre_commit "$@"
```

**Phase 2: Meta-Quality Infrastructure (Week 5)**

Create `tests/meta/` directory:

```python
# tests/meta/test_sdp_quality_gates.py
"""Test that SDP follows its own quality gates."""

import pytest
from pathlib import Path

def test_sdp_coverage_meets_threshold():
    """SDP must have ≥80% coverage."""
    result = subprocess.run(
        ["pytest", "--cov=src/sdp", "--cov-report=json"],
        capture_output=True
    )

    # Parse coverage report
    coverage = json.loads(Path("coverage.json").read_text())
    total_coverage = coverage["totals"]["percent_covered"]

    assert total_coverage >= 80, f"SDP coverage {total_coverage}% < 80%"

def test_sdp_file_size_within_limits():
    """SDP source files must be <200 LOC."""
    violations = []

    for py_file in Path("src/sdp").rglob("*.py"):
        lines = len(py_file.read_text().split("\n"))
        if lines > 200:
            violations.append(f"{py_file}: {lines} lines")

    assert not violations, f"Files exceed 200 LOC:\n" + "\n".join(violations)

def test_sdp_no_bare_except():
    """SDP must not have bare except clauses."""
    violations = []

    for py_file in Path("src/sdp").rglob("*.py"):
        content = py_file.read_text()
        if re.search(r"except:\s*$", content, re.MULTILINE):
            violations.append(f"{py_file}: Bare except found")

    assert not violations, f"Bare except clauses found:\n" + "\n".join(violations)

def test_sdp_has_full_type_hints():
    """SDP must have full type hints."""
    # Check that functions have type annotations
    violations = []

    for py_file in Path("src/sdp").rglob("*.py"):
        tree = ast.parse(py_file.read_text())

        for node in ast.walk(tree):
            if isinstance(node, ast.FunctionDef):
                # Skip tests and private methods
                if node.name.startswith("_") or node.name.startswith("test_"):
                    continue

                # Check for type hints
                if not node.returns or not any(arg.annotation for arg in node.args.args):
                    violations.append(f"{py_file}:{node.name} missing type hints")

    assert not violations, f"Functions missing type hints:\n" + "\n".join(violations)

def test_all_hooks_are_tested():
    """All shell hooks must have corresponding Python tests."""
    hook_files = list(Path("hooks/").glob("*.sh"))
    test_files = list(Path("tests/unit/hooks/").glob("test_*.py"))

    hooks_tested = set()
    for test_file in test_files:
        # Extract hook name from test
        content = test_file.read_text()
        matches = re.findall(r'test_(\w+)', content)
        hooks_tested.update(matches)

    missing_tests = []
    for hook_file in hook_files:
        hook_name = hook_file.stem  # e.g., "pre-commit"
        if hook_name.replace("-", "_") not in hooks_tested:
            missing_tests.append(hook_file)

    assert not missing_tests, f"Hooks without tests:\n" + "\n".join(map(str, missing_tests))
```

**Phase 3: GitHub Integration Testing (Week 6)**

```python
# tests/integration/test_github.py
"""Test GitHub integration with mocked API."""

import pytest
from unittest.mock import Mock, patch
from sdp.github.sync import GitHubSync

@pytest.fixture
def mock_github():
    """Mock PyGithub client."""
    with patch("sdp.github.sync.Github") as mock:
        repo = Mock()
        mock.return_value.get_repo.return_value = repo
        yield repo

def test_sync_creates_github_issue(mock_github):
    """Test that workstream sync creates GitHub issue."""
    sync = GitHubSync(token="test", repo="test/test")

    workstream = Workstream(
        id="00-001-01",
        title="Add user authentication",
        status="pending"
    )

    sync.create_issue(workstream)

    # Verify issue created
    mock_github.create_issue.assert_called_once_with(
        title="[WS-00-001-01] Add user authentication",
        body=sync._format_issue_body(workstream),
        labels=["workstream", "pending"]
    )

def test_sync_updates_milestone(mock_github):
    """Test that feature completion updates milestone."""
    sync = GitHubSync(token="test", repo="test/test")

    sync.complete_feature("F01")

    # Verify milestone closed
    mock_gighthouse.get_milestone.assert_called_once_with(number=1)
    mock_gighthouse.edit.assert_called_once_with(state="closed")
```

**Phase 4: Integration Test Suite (Week 7-8)**

```python
# tests/integration/test_full_workflow.py
"""End-to-end tests for complete SDP workflows."""

import subprocess
import tempfile
from pathlib import Path

@pytest.fixture
def sdp_project(tmp_path):
    """Create a temporary SDP project."""
    # Initialize project
    subprocess.run(["sdp", "init", "--dir", str(tmp_path)], check=True)

    # Create example workstream
    ws_file = tmp_path / "docs/workstreams/backlog/00-001-01.md"
    ws_file.write_text("""
# Workstream 00-001-01: Add hello_world function

## Goal
Implement a simple hello_world function.

## Acceptance Criteria
- [ ] Function returns "Hello, World!"
- [ ] Has type hints
- [ ] Has tests

## Context
First workstream for testing.

## Steps
1. Write test (RED)
2. Implement function (GREEN)
3. Add type hints (REFACTOR)

## Dependencies
None

## Size
SMALL
""")

    return tmp_path

def test_build_workflow(sdp_project):
    """Test complete @build workflow."""
    # Run build
    result = subprocess.run(
        ["@build", "00-001-01"],
        cwd=sdp_project,
        capture_output=True
    )

    assert result.returncode == 0, f"Build failed: {result.stderr}"

    # Verify artifacts created
    assert (sdp_project / "src/hello.py").exists()
    assert (sdp_project / "tests/test_hello.py").exists()

    # Verify tests pass
    test_result = subprocess.run(
        ["pytest", "tests/test_hello.py"],
        cwd=sdp_project
    )
    assert test_result.returncode == 0

def test_quality_gates_enforced(sdp_project):
    """Test that quality gates block bad code."""
    # Try to create file >200 LOC
    bad_file = sdp_project / "src/bad.py"
    bad_file.write_text("\n".join(["# Line " + str(i) for i in range(250)]))

    # Run post-build hook
    result = subprocess.run(
        ["hooks/post-build.sh", "src/bad.py"],
        cwd=sdp_project,
        capture_output=True
    )

    assert result.returncode != 0
    assert "exceeds 200 LOC" in result.stdout.decode()
```

### Success Metrics

- **Coverage for `/src/sdp/`**: ≥80% (actual, not hw_checker's)
- **Hooks tested**: 100% (16/16 shell hooks have Python tests)
- **GitHub integration tested**: >80% of modules covered
- **Integration tests**: ≥5 end-to-end workflow tests
- **README badge**: Reflects actual SDP coverage (not hw_checker)

---

## 9. Artifact Quality

> **Experts:** Martin Fowler (Refactoring), Dieter Rams (Minimal Design), Steve Krug (Don't Make Me Think)

### Current Issues

1. **Template has Russian content** — `Цель`, `Зависимость`, etc. without English
2. **Format inconsistency** — Two workstream formats coexist
3. **Execution reports sparse** — Only 7/~70 completed WS have reports
4. **Variable draft quality** — F014 is 233 lines, others are sparse

### Solution: Flexible Core + Optional Extensions

**Phase 1: Standardize Template to English (Week 1)**

Replace `docs/workstreams/TEMPLATE.md` Russian headers:

```markdown
# Workstream PP-FFF-SS: [Title]

**Status**: pending | in_progress | completed
**Feature**: FFF
**Size**: SMALL | MEDIUM | LARGE
**Complexity**: 1-10

---

## Goal
[What does this workstream achieve?]

## Acceptance Criteria
- [ ] [Criterion 1]
- [ ] [Criterion 2]
- [ ] [Criterion 3]

## Context
[Why is this needed? What problem does it solve?]

## Dependencies
**Prerequisites**: [What must exist before starting?]
**Blocked by**: [WS-XXX-YY, WS-XXX-ZZ]
**Blocks**: [WS-XXX-AA]

## Steps
1. [Step 1]
2. [Step 2]
3. [Step 3]

### Step 1 Detail
**Expected Result**: [What should happen?]
**Validation**: [How to verify?]

## Size Estimation
- **Estimated LOC**: [500-1500 for MEDIUM]
- **Estimated Files**: [3-5 for MEDIUM]
- **Estimated Complexity**: [CC target]

## Constraints
- [ ] Files <200 LOC
- [ ] Coverage ≥80%
- [ ] Type hints required
- [ ] No external dependencies (if applicable)

## Completion Criteria
- [ ] All AC met
- [ ] Tests pass
- [ ] Quality gates pass
- [ ] Code reviewed
- [ ] Workstream file moved to completed/

## Execution Report
[Fill after completion - AUTO-GENERATED by @build]

### Time Taken
[Actual duration]

### Actual LOC
[Final LOC]

### Actual Coverage
[Final coverage %]

### Deviations
[What changed from plan?]

### Lessons Learned
[What would you do differently?]
```

**Phase 2: Define Core vs Optional Schema (Week 2)**

Create `docs/schema/workstream.schema.json`:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": [
    "ws_id",
    "feature",
    "status",
    "size",
    "goal",
    "acceptance_criteria",
    "context",
    "steps"
  ],
  "properties": {
    "ws_id": {
      "type": "string",
      "pattern": "^\\d{2}-\\d{3}-\\d{2}$"
    },
    "feature": {
      "type": "string",
      "pattern": "^F\\d{2}$"
    },
    "status": {
      "type": "string",
      "enum": ["pending", "in_progress", "completed"]
    },
    "size": {
      "type": "string",
      "enum": ["SMALL", "MEDIUM", "LARGE"]
    },
    "goal": {
      "type": "string",
      "minLength": 10
    },
    "acceptance_criteria": {
      "type": "array",
      "items": {"type": "string"},
      "minItems": 1
    },
    "context": {
      "type": "string"
    },
    "steps": {
      "type": "array",
      "items": {"type": "string"},
      "minItems": 1
    },
    "dependencies": {
      "type": "object",
      "properties": {
        "prerequisites": {"type": "array"},
        "blocked_by": {"type": "array", "items": {"type": "string"}},
        "blocks": {"type": "array", "items": {"type": "string"}}
      }
    },
    "execution_graph": {
      "type": "object"  // Optional
    },
    "code_blocks": {
      "type": "array"  // Optional
    },
    "completion_criteria": {
      "type": "array"  // Optional
    },
    "constraints": {
      "type": "array"  // Optional
    },
    "execution_report": {
      "type": "object"  // Optional
    }
  }
}
```

**Phase 3: Auto-Generate Execution Reports (Week 3)**

Update `@build` skill to append report automatically:

```python
# src/sdp/skills/build.py
class BuildSkill:
    def complete(self, ws_id: str):
        """Mark workstream complete and append execution report."""
        # Generate report
        report = self._generate_execution_report(ws_id)

        # Append to workstream file
        ws_file = Path(f"docs/workstreams/in_progress/{ws_id}.md")
        content = ws_file.read_text()

        # Add execution report section
        content += f"\n\n## Execution Report\n{report}"

        # Save to completed/
        completed_file = Path(f"docs/workstreams/completed/{ws_id}.md")
        completed_file.write_text(content)

        # Remove from in_progress/
        ws_file.unlink()

    def _generate_execution_report(self, ws_id: str) -> str:
        """Generate execution report from metadata."""
        stats = self._gather_statistics(ws_id)

        return f"""
### Time Taken
{stats['duration']} (started {stats['start_time']})

### Actual LOC
{stats['loc']} files, {stats['total_lines']} total lines

### Actual Coverage
{stats['coverage']}% (target: 80%)

### Quality Gates
- Mypy: {'✅ PASSED' if stats['mypy'] else '❌ FAILED'}
- Ruff: {'✅ PASSED' if stats['ruff'] else '❌ FAILED'}
- File size: {'✅ PASSED' if stats['file_size'] else '❌ FAILED'}

### Deviations
{stats['deviations'] or 'None'}

### Lessons Learned
{stats['lessons'] or 'N/A'}
"""
```

**Phase 4: Enhanced Validation (Week 4)**

Update `hooks/validate-artifacts.sh`:

```bash
validate_workstream() {
    local ws_file=$1

    echo "🔍 Validating workstream: $ws_file"

    # Check required sections
    required_sections=("Goal" "Acceptance Criteria" "Context" "Steps")
    for section in "${required_sections[@]}"; do
        if ! grep -q "^## $section" "$ws_file"; then
            echo "❌ Missing required section: $section"
            return 1
        fi
    done

    # Warn on optional sections
    optional_sections=("Dependencies" "Constraints" "Completion Criteria")
    for section in "${optional_sections[@]}"; do
        if ! grep -q "^## $section" "$ws_file"; then
            echo "⚠️  Missing optional section: $section"
        fi
    done

    # Check AC format (must be testable)
    ac_count=$(grep -c "^- \[ \]" "$ws_file" || true)
    if [ "$ac_count" -lt 1 ]; then
        echo "❌ No acceptance criteria found"
        return 1
    fi

    # Check completed workstreams have execution report
    if [[ "$ws_file" =~ completed/ ]]; then
        if ! grep -q "^## Execution Report" "$ws_file"; then
            echo "⚠️  Completed workstream missing execution report"
        fi
    fi

    echo "✅ Workstream validation passed"
}
```

**Phase 5: Artifact Quality Dashboard (Week 5)**

Create `scripts/artifact_quality.py`:

```python
#!/usr/bin/env python3
"""Analyze artifact quality across repository."""

from pathlib import Path
import re

def analyze_workstreams():
    """Analyze workstream completion and quality."""

    backlog = list(Path("docs/workstreams/backlog/").glob("*.md"))
    in_progress = list(Path("docs/workstreams/in_progress/").glob("*.md"))
    completed = list(Path("docs/workstreams/completed/").glob("*.md"))

    total = len(backlog) + len(in_progress) + len(completed)

    # Check completion rate
    completion_rate = len(completed) / total * 100 if total > 0 else 0

    # Check execution reports
    with_reports = 0
    for ws in completed:
        if "## Execution Report" in ws.read_text():
            with_reports += 1

    report_rate = with_reports / len(completed) * 100 if completed else 0

    # Check required sections
    missing_sections = []
    for ws in backlog + in_progress:
        content = ws.read_text()
        required = ["Goal", "Acceptance Criteria", "Context", "Steps"]
        for section in required:
            if f"## {section}" not in content:
                missing_sections.append(f"{ws}: Missing {section}")

    # Print report
    print(f"📊 Workstream Quality Report")
    print(f"")
    print(f"Total: {total}")
    print(f"  - Backlog: {len(backlog)}")
    print(f"  - In Progress: {len(in_progress)}")
    print(f"  - Completed: {len(completed)}")
    print(f"")
    print(f"Completion Rate: {completion_rate:.1f}%")
    print(f"Report Rate: {report_rate:.1f}% ({with_reports}/{len(completed)} completed)")
    print(f"")

    if missing_sections:
        print(f"⚠️  Missing Required Sections:")
        for item in missing_sections[:10]:
            print(f"  - {item}")
        if len(missing_sections) > 10:
            print(f"  ... and {len(missing_sections) - 10} more")

    return {
        "total": total,
        "completed": len(completed),
        "completion_rate": completion_rate,
        "report_rate": report_rate,
    }

if __name__ == "__main__":
    analyze_workstreams()
```

### Success Metrics

- **Template compliance**: >95% of workstreams follow template
- **Execution report rate**: >80% of completed WS have reports
- **Required sections**: 100% compliance (enforced by validation)
- **Quality dashboard**: Automated reporting runs in CI/CD

---

## 10. Tooling & Automation

> **Experts:** Corey Haines (Pragmatic Automation), Kelsey Hightower (Declarative Config), Troy Hunt (Fail Fast)

### Current Gaps

1. **No setup wizard** — Manual configuration across multiple files
2. **Error messages poor** — "Error: X failed" with no remediation
3. **No health check** — Can't verify setup without running full workflow
4. **Manual artifact tracking** — No dashboards or metrics
5. **Incomplete tooling** — Missing commands for common tasks

### Solution: Self-Service Diagnostics + Enhanced Tooling

**Phase 1: `sdp doctor` Command (Week 1)**

```python
# src/sdp/cli/doctor.py
import sys
import subprocess
from pathlib import Path

class SDPDoctor:
    """Diagnose SDP environment issues."""

    def __init__(self):
        self.checks = []
        self.failures = []
        self.warnings = []

    def run_all_checks(self):
        """Run all diagnostic checks."""
        print("🔍 SDP Environment Check")
        print("=" * 50)

        self.check_python_version()
        self.check_poetry()
        self.check_git_hooks()
        self.check_beads()
        self.check_github_cli()
        self.check_config_files()

        self.print_summary()

    def check_python_version(self):
        """Check Python version."""
        version = sys.version_info
        required = (3, 10)

        if version >= required:
            self.checks.append(("Python", f"✅ {version.major}.{version.minor}.{version.micro}", True))
        else:
            msg = f"❌ Python {version.major}.{version.minor} (required: 3.10+)"
            self.checks.append(("Python", msg, False))
            self.failures.append(("Python", "Upgrade to Python 3.10 or higher", f"pyenv install 3.14 && pyenv global 3.14"))

    def check_beads(self):
        """Check Beads CLI (optional)."""
        result = subprocess.run(["which", "bd"], capture_output=True)

        if result.returncode == 0:
            version = subprocess.run(["bd", "--version"], capture_output=True)
            self.checks.append(("Beads", f"✅ {version.stdout.decode().strip()}", True))
        else:
            msg = "⚠️  Beads CLI not found (optional, required for task tracking)"
            self.checks.append(("Beads", msg, True))
            self.warnings.append(("Beads", "Install for task tracking", "brew install go && go install github.com/steveyegge/beads/cmd/bd@latest"))

    def print_summary(self):
        """Print check summary."""
        passed = sum(1 for _, _, status in self.checks if status)
        total = len(self.checks)

        print("")
        print(f"Summary: {passed}/{total} checks passed")

        if self.failures:
            print("")
            print("❌ Failures (must fix):")
            for name, desc, fix in self.failures:
                print(f"  {name}: {desc}")
                print(f"    💡 Fix: {fix}")

        if self.warnings:
            print("")
            print("⚠️  Warnings (optional):")
            for name, desc, fix in self.warnings:
                print(f"  {name}: {desc}")
                print(f"    💡 Fix: {fix}")

        if self.failures:
            print("")
            print("Run the following to fix:")
            for _, _, fix in self.failures:
                print(f"  {fix}")
            sys.exit(1)

# CLI entry point
def main():
    doctor = SDPDoctor()
    doctor.run_all_checks()

if __name__ == "__main__":
    main()
```

**Phase 2: Enhanced Error Messages (Week 2)**

Create `src/sdp/errors.py`:

```python
from dataclasses import dataclass
from typing import Optional, List

@dataclass
class SDPError:
    """Structured error with remediation."""

    title: str
    message: str
    remediation: str
    docs_url: Optional[str] = None
    chat_url: Optional[str] = None

    def format(self) -> str:
        """Format error for terminal output."""
        output = [f"❌ {self.title}", ""]
        output.append(self.message)
        output.append("")
        output.append(f"💡 {self.remediation}")

        if self.docs_url:
            output.append("")
            output.append(f"📖 Docs: {self.docs_url}")

        if self.chat_url:
            output.append("")
            output.append(f"🆘 Chat: {self.chat_url}")

        return "\n".join(output)

# Predefined errors
class SDPErrors:
    BEADS_NOT_FOUND = SDPError(
        title="Beads CLI not found",
        message="Beads is optional for SDP, but required for task tracking.",
        remediation="Install: brew install go && go install github.com/steveyegge/beads/cmd/bd@latest",
        docs_url="https://sdp.dev/docs/beads.html",
    )

    COVERAGE_TOO_LOW = SDPError(
        title="Coverage below threshold",
        message="Test coverage is 72%, but 80% is required.",
        remediation="Add tests for uncovered lines:\n  pytest --cov=src/ --cov-report=term-missing",
        docs_url="https://sdp.dev/docs/quality-gates.html#coverage",
    )

    FILE_TOO_LARGE = SDPError(
        title="File exceeds size limit",
        message="src/module.py has 245 lines (max: 200).",
        remediation="Refactor into smaller modules:\n  - Extract functions to separate files\n  - Move generated code to generated/",
        docs_url="https://sdp.dev/docs/quality-gates.html#file-size",
    )

    TYPE_HINTS_MISSING = SDPError(
        title="Type hints missing",
        message="Function 'process_data' lacks type annotations.",
        remediation="Add type hints:\n  def process_data(data: List[str]) -> Dict[str, int]:",
        docs_url="https://sdp.dev/docs/code-patterns.html#type-hints",
    )

    BARE_EXCEPT_FOUND = SDPError(
        title="Bare except clause found",
        message="Found 'except:' in src/module.py:42",
        remediation="Replace with specific exception:\n  except ValueError as e:\n  # or: except Exception as e: # noqa: E722",
        docs_url="https://sdp.dev/docs/forbidden-patterns.html#bare-except",
    )

# Usage in hooks
def check_coverage(file: str, coverage: float):
    if coverage < 80:
        error = SDPErrors.COVERAGE_TOO_LOW
        error.message = f"Coverage for {file} is {coverage}%, but 80% is required."
        print(error.format())
        sys.exit(1)
```

**Phase 3: Interactive Setup Wizard (Week 3)**

```python
# src/sdp/cli/init.py
from pathlib import Path
import json

class SetupWizard:
    """Interactive SDP setup wizard."""

    def __init__(self):
        self.config = {}

    def run(self):
        """Run setup wizard."""
        print("🚀 SDP Project Initialization")
        print("=" * 50)

        self.step_project_metadata()
        self.step_optional_dependencies()
        self.step_quality_gates()
        self.step_git_hooks()
        self.step_validation()

        self.save_config()
        self.print_next_steps()

    def step_project_metadata(self):
        """Step 1: Gather project metadata."""
        print("\n📋 Step 1/5: Project metadata")

        self.config["project_name"] = input("Project name: ")
        self.config["project_id"] = input("Project ID (00-99): ")
        self.config["description"] = input("Description: ")

    def step_optional_dependencies(self):
        """Step 2: Configure optional dependencies."""
        print("\n🔌 Step 2/5: Optional dependencies")

        # Beads
        has_beads = subprocess.run(["which", "bd"], capture_output=True)
        if has_beads.returncode == 0:
            use_beads = input("Install Beads CLI for task tracking? [Y/n]: ")
            self.config["beads_enabled"] = use_beads.lower() != "n"

        # GitHub CLI
        has_gh = subprocess.run(["which", "gh"], capture_output=True)
        if has_gh.returncode == 0:
            use_gh = input("Configure GitHub integration? [Y/n]: ")
            self.config["github_enabled"] = use_gh.lower() != "n"

        # Telegram
        use_telegram = input("Configure Telegram notifications? [y/N]: ")
        self.config["telegram_enabled"] = use_telegram.lower() == "y"

    def step_quality_gates(self):
        """Step 3: Configure quality gates."""
        print("\n🎯 Step 3/5: Quality gate configuration")

        self.config["coverage_threshold"] = input("Coverage threshold [80]: ") or "80"
        self.config["max_file_size"] = input("Max file size (LOC) [200]: ") or "200"
        self.config["type_checking"] = input("Type checking (strict/mypy/none) [strict]: ") or "strict"

    def step_git_hooks(self):
        """Step 4: Install git hooks."""
        print("\n🪝 Step 4/5: Git hooks installation")

        hooks = ["pre-commit", "pre-push", "post-build", "pre-deploy"]
        for hook in hooks:
            print(f"  Installing {hook} hook...")
            hook_src = Path(f"hooks/{hook}.sh")
            hook_dst = Path(f".git/hooks/{hook}")
            hook_dst.symlink_to(hook_src.absolute())

        print("✅ Git hooks installed")

    def step_validation(self):
        """Step 5: Validate setup."""
        print("\n✅ Step 5/5: Validation")
        print("Running `sdp doctor`...")

        result = subprocess.run(["sdp", "doctor"], capture_output=True)
        if result.returncode == 0:
            print(result.stdout.decode())
        else:
            print("❌ Setup validation failed")
            print(result.stderr.decode())
            sys.exit(1)

    def save_config(self):
        """Save configuration."""
        config_file = Path("sdp.config.json")
        config_file.write_text(json.dumps(self.config, indent=2))

        print(f"\n✅ Configuration saved to {config_file}")

    def print_next_steps(self):
        """Print next steps."""
        print("\n✅ Setup complete! Next steps:")
        print(f"  1. Run: @feature 'Add user authentication'")
        print(f"  2. Or: cp docs/workstreams/TEMPLATE.md docs/workstreams/backlog/{self.config['project_id']}-001-01.md")
        print(f"  3. Read: docs/START_HERE.md")

# CLI entry point
def main():
    wizard = SetupWizard()
    wizard.run()
```

**Phase 4: Artifact Dashboard (Week 4)**

```python
# src/sdp/cli/dashboard.py
class ArtifactDashboard:
    """Display artifact quality metrics."""

    def show(self):
        """Show dashboard."""
        workstream_stats = self._analyze_workstreams()
        intent_stats = self._analyze_intents()
        quality_metrics = self._analyze_quality()

        self._print_dashboard(workstream_stats, intent_stats, quality_metrics)

    def _analyze_workstreams(self):
        """Analyze workstream artifacts."""
        # Reuse artifact_quality.py logic
        import subprocess
        result = subprocess.run(
            ["python", "scripts/artifact_quality.py"],
            capture_output=True
        )
        # Parse output
        return {...}

    def _print_dashboard(self, ws_stats, intent_stats, quality):
        """Print formatted dashboard."""
        print("📊 SDP Artifact Dashboard")
        print("=" * 60)
        print("")
        print("Workstreams:")
        print(f"  Total: {ws_stats['total']}")
        print(f"  Completed: {ws_stats['completed']} ({ws_stats['completion_rate']:.1f}%)")
        print(f"  With Reports: {ws_stats['report_rate']:.1f}%")
        print("")
        print("Intents:")
        print(f"  Total: {intent_stats['total']}")
        print(f"  Valid Schema: {intent_stats['valid_rate']:.1f}%")
        print("")
        print("Quality Metrics:")
        print(f"  Coverage: {quality['coverage']:.1f}%")
        print(f"  Avg Complexity: {quality['cc']:.1f}")
        print(f"  Files >200 LOC: {quality['large_files']}")
        print("")
        print("=" * 60)

# CLI entry point
def main():
    dashboard = ArtifactDashboard()
    dashboard.show()
```

### Success Metrics

- **`sdp doctor` runtime**: <5 seconds
- **Setup wizard completion**: >90% of new users complete without errors
- **Error remediation success**: >80% of errors resolved without support
- **Dashboard usage**: Track how often users run dashboard (indicates engagement)

---

## Implementation Plan

### Phase 1: Foundation (Weeks 1-4) — **P0**

**Goal:** Eliminate critical blockers and establish core infrastructure.

- [ ] **Week 1: Protocol & Tooling**
  - Create GLOSSARY.md (150+ terms)
  - Implement `sdp doctor` command
  - Add radon for complexity checking
  - Create `quality-gate.toml` schema

- [ ] **Week 2: Documentation & Setup**
  - Create docs/START_HERE.md and docs/SITEMAP.md
  - Implement `sdp init` wizard
  - Standardize command prefixes (`@` vs `/`)
  - Translate TEMPLATE.md to English

- [ ] **Week 3: Validation & Enforcement**
  - Make all hooks fail instead of warn
  - Create `.github/workflows/sdp-quality-gate.yml`
  - Make architecture checks portable
  - Add link validation to pre-commit

- [ ] **Week 4: Artifacts & Quality**
  - Define flexible core schema for workstreams
  - Enhance `validate-artifacts.sh`
  - Auto-generate execution reports
  - Create artifact quality dashboard

### Phase 2: Deep Improvements (Weeks 5-12) — **P1**

**Goal:** Strengthen core systems and improve developer experience.

- [ ] **Week 5-6: TDD Enforcement**
  - Implement contract immutability (chmod + markers)
  - Add git timestamp verification
  - Add pre-commit phase enforcement
  - Enhance TDDRunner with state tracking

- [ ] **Week 7-8: Meta-Testing**
  - Extract pre-commit.sh to Python (testable)
  - Create `tests/meta/` directory
  - Add GitHub integration tests
  - Create integration test suite

- [ ] **Week 9-10: Skill System**
  - Extract BuildSkill to class with DI
  - Extract IdeaSkill, DesignSkill
  - Create PipelineExecutor
  - Add event emission

- [ ] **Week 11-12: Real-World Integration**
  - Create Strangler pattern adoption guide
  - Create team coordination playbook
  - Implement migration dashboard
  - Add hotfix/triage process

### Phase 3: Polish & Optimization (Weeks 13-16) — **P2**

**Goal:** Optimize experience and prepare for v1.0.

- [ ] **Week 13: Documentation Reorganization**
  - Reorganize docs/ by role (beginner/reference/internals)
  - Resolve all 124 documentation TODOs
  - Add cross-reference validation
  - Create decision tree diagrams

- [ ] **Week 14: Enhanced Error Messages**
  - Create SDPError framework
  - Add remediation to all error paths
  - Add "Get help" links to errors
  - Document error patterns

- [ ] **Week 15: Performance & Optimization**
  - Benchmark hook runtimes
  - Optimize slow checks
  - Add parallel execution where possible
  - Create performance dashboard

- [ ] **Week 16: v1.0 Preparation**
  - Update all version references
  - Create migration guides (v0.5 → v1.0)
  - Write comprehensive CHANGELOG
  - Create release notes

### Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| **Onboarding time** | ~2 hours | <30 minutes |
| **Protocol questions/week** | ~5 | 0 |
| **SDP code coverage** | Unknown (hw_checker's 91%) | ≥80% actual |
| **Quality gate pass rate** | Unknown | >95% |
| **Workstream report rate** | ~10% | >80% |
| **CI/CD adoption** | 0 | >80% of projects |
| **Setup completion rate** | Unknown | >90% |
| **Error resolution without support** | Unknown | >80% |

---

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Breaking changes (ID format)** | High | Provide migration script, deprecation period (6 months) |
| **Hook performance** | Medium | Benchmark, parallelize, caching |
| **Template standardization** | Medium | Keep flexible core, warn before enforce |
| **Documentation reorganization** | Low | Use redirects, maintain old paths temporarily |
| **TDD enforcement friction** | Medium | Make painless when followed correctly |
| **Skill system complexity** | Low | Incremental extraction, backward compatibility |

---

## Conclusion

This plan transforms SDP from a B+ framework with excellent ideas but execution gaps into an A+ production-ready protocol that:

1. **Cannot be misunderstood** (canonical glossary, consistent terminology)
2. **Is easy to adopt** (<30 min onboarding, interactive setup)
3. **Enforces quality effectively** (CI/CD, strong validation)
4. **Has composable architecture** (DI pipeline, testable skills)
5. **Dogfoods its principles** (meta-testing, covered hooks)
6. **Supports real teams** (CI/CD templates, Strangler adoption)
7. **Has clear documentation** (role-based, progressive disclosure)
8. **Enforces TDD unbreakably** (contract + timestamp hybrid)
9. **Standardizes artifacts** (flexible core schema)
10. **Provides self-service tooling** (doctor, wizard, dashboard)

**Estimated timeline**: 16 weeks (4 months) for full implementation

**Recommended approach**: Start with Phase 1 (Foundation, 4 weeks), measure impact, then continue to Phase 2 based on feedback.

---

**Next Steps:**

1. Review this plan with stakeholders
2. Prioritize workstreams for Phase 1
3. Create `WS-IMPROVE-001` through `WS-IMPROVE-040` for implementation
4. Begin with `@idea "SDP A+ Improvement Plan"` to generate detailed specs

**Status:** Ready for implementation planning

**Date:** 2026-01-29
**Version:** 1.0
