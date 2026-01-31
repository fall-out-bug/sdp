# Changelog

All notable changes to the Spec-Driven Protocol (SDP).

> **📝 Meta-note:** Versions documented as they are released. Development is AI-assisted.

## [0.7.0] - 2026-01-31

### Added - Feature F034: A+ Quality Initiative

- **00-034-01:** Split Large Files (Phase 1: Core)
- **00-034-02:** Split Large Files (Phase 2: Beads/Unified)
- **00-034-03:** Increase Test Coverage to 80%+ (achieved 85%)
- **00-034-04:** Documentation Consistency
- **00-034-05:** Extract Domain Layer (Clean Architecture)
- **00-034-06:** Add `sdp status` Command
- **00-034-07:** Add Skill Discovery (@help)
- **00-034-08:** Remove Legacy Code (~600 LOC)

**Что нового:**
- Coverage 68% → 85.28%
- Clean Architecture: domain layer extracted, no beads→core violations
- `sdp status` — project and guard status
- Skill discovery via `@help` / `sdp skills`
- Skills optimized (~64% reduction), repo restructured
- Lint clean (F401, F821, E501 fixed)

**Использование:**

```bash
# Project status
sdp status

# Skill discovery
sdp skills list
```

---

## [0.5.2] - 2026-01-31

### Added - Feature F025: pip-audit Security Scanning

- **00-025-01:** pip-audit + Dependabot — dependency vulnerability scanning in CI/CD

**Что нового:**
- pip-audit runs on every PR/push (blocks merge on vulnerabilities)
- PR comments include CVE details, severity, fix versions
- Dependabot weekly PRs for Python + GitHub Actions
- SECURITY.md policy, docs/internals/development.md updated

**Использование:**

```bash
# Run vulnerability scan locally
poetry run pip-audit

# Generate JSON report
poetry run pip-audit --format json --desc -o audit-report.json
```

---

## [0.5.1] - 2026-01-31

### Added - Feature F020: Fast Feedback (Hooks Extraction & Project-Agnostic)

- **00-020-01:** Git hooks extracted to Python — pre-commit, pre-push, post-build, pre-deploy
- **00-020-02:** Hooks project-agnostic — auto-detect project root, `quality-gate.toml` config
- **00-020-04:** Fix `except Exception: pass` in common.py (Issue 006)

**Что нового:**
- Shell hooks → testable Python modules (`src/sdp/hooks/`)
- `find_project_root()`, `find_workstream_dir()` — auto-detection
- Configuration via `quality-gate.toml` [workstreams.dir]
- Hooks coverage 92%, mypy --strict

**Использование:**

```bash
# Hooks run on any SDP project (dogfooding)
python -m sdp.hooks.pre_commit
```

---

## [0.6.1] - 2026-01-31

### Added - Feature F031: Migrate Core Exceptions to SDPError

- **00-031-01:** Core exceptions inherit from SDPError with ErrorCategory, remediation, docs_url

**Что нового:**
- WorkstreamParseError, CircularDependencyError, MissingDependencyError → SDPError
- ModelMappingError, ContractViolationError, HumanEscalationError → SDPError
- format_terminal() and format_json() on SDPError base
- Actionable error messages with remediation steps

**Использование:**

```bash
# Exceptions now include structured output
from sdp.core.workstream import WorkstreamParseError
error = WorkstreamParseError("Invalid ws_id")
print(error.format_terminal())  # Terminal output with remediation
print(error.format_json())      # JSON for CI/CD
```

---

## [0.6.0] - 2026-01-31

### Added - Feature F030: Test Coverage Expansion

- **00-030-01:** GitHub Integration Tests — client 85%, retry 93%, sync 52%
- **00-030-02:** Adapter Tests — base 86%, claude 86%, opencode 93%
- **00-030-03:** Core Functionality Tests — workstream 84%, builder 81%, model 80%

**Что нового:**
- 72 новых unit-тестов для GitHub, adapters, core
- Все тесты используют mocks (без реальных API-вызовов)
- mypy --strict для всех тестовых файлов

**Использование:**

```bash
# Запуск тестов F030
uv run pytest tests/unit/adapters/ tests/unit/core/ tests/unit/github/ -v

# Проверка типов
uv run mypy tests/unit/adapters/ tests/unit/core/ tests/unit/github/ --strict
```

## [0.4.0] - 2026-01-27

### Added - Feature F011: PRD Command (6 workstreams)
- PRD command with project type profiles (cli, service, library)
- Auto-generated architecture diagrams from code annotations
- `@prd:` annotation parser for documentation updates
- Line limits validator for PRD documents
- Diagram generator (Mermaid format)
- CodeReview hook integration for PRD validation

### Added - Feature F003: Two-Stage Review (5 workstreams)
- Stage 1: Spec Compliance (goal achievement, AC coverage, specification alignment)
- Stage 2: Code Quality (coverage >= 80%, mypy strict, AI-readiness)
- Stage 2 only runs if Stage 1 passes — no polishing incorrect code
- Updated `/codereview` skill with two-stage workflow

### Added - Feature F004: Platform Adapters (4 workstreams)
- `PlatformAdapter` interface for unified API
- `detect_platform()` for auto IDE detection
- Claude Code adapter (`.claude/` support)
- Cursor adapter (`.codex/` support)
- OpenCode adapter (`.opencode/` support)

### Added - Feature F005: Extension System (3 workstreams)
- `sdp.local/` and `~/.sdp/extensions/{name}/` support
- `extension.yaml` manifest format
- Extension auto-discovery and loading
- Hooks, patterns, skills, integrations components

### Added - Feature F007: Oneshot & Hooks (10 workstreams)
- `/oneshot` command for autonomous feature execution
- Git hooks: pre-commit, post-commit, pre-push
- Quality gates enforcement
- Cursor agents integration
- `/debug` and `/test` commands
- `/idea` and `/design` skills

### Added - Feature F008: Contract-Driven WS Tiers (9 workstreams)
- Starter, Standard, Advanced tiers
- Capability tier validator
- Model mapping registry
- Tier auto-promotion
- Escalation metrics

### Added - Feature F010: SDP Infrastructure (5 workstreams)
- Submodule support
- PP-FFF-SS naming convention
- Content synchronization

### Changed
- Workstream ID format changed to PP-FFF-WW
- Enhanced GitHub bidirectional sync service
- Improved documentation structure

### Statistics
- **Total Workstreams:** 58
- **Completed:** 48 (83%)
- **Features:** 8 (F003, F004, F005, F006, F007, F008, F010, F011)

## [0.4.0-rc] - 2026-01-25

### Added - Feature F003: Two-Stage Review (5 workstreams)
- Peer review skill with 17-point quality checklist
- Systematic debugging skill with 4-phase process
- Fix verification tests skill for completed workstreams
- Debug command implementation with breakpoint-style workflow
- Enhanced test coverage metrics (F191 rule enforcement)

### Added - Feature F004: Platform Adapters (4 workstreams)
- Claude Code adapter implementation
- Cursor agent adapter interface
- OpenCode multi-IDE support
- Unified platform adapter interface

### Added - Feature F005: Extension System (3 workstreams)
- Extension interface and base classes
- hw_checker extension implementation
- GitHub integration extension

### Added - Feature F006: Core SDP (6 workstreams)
- Workstream parser with PP-FFF-WW format support
- Feature decomposition from requirements
- Project map generation and maintenance
- PIP package structure for distribution
- File size reduction utilities
- Integration test suite

### Added - Feature F007: Oneshot & Hooks (11 workstreams)
- Oneshot autonomous execution with orchestrator agent
- Git hooks integration (pre-commit, post-commit)
- Debug command with 4-phase systematic debugging
- Test command coverage validation
- Documentation generation workflow
- Test artifact cleanup utilities
- EP30 misclassification fix
- Debug title correction
- Idea/design skill integration

### Added - Feature F008: Contract-Driven WS Tiers (9 workstreams)
- Workstream contract specification format
- Capability tier validator (T0-T3)
- Model mapping registry for LLM selection
- Test command workflow with tier routing
- Model-agnostic builder/router implementation
- Model selection optimization (cost/latency tradeoffs)
- Tier auto-promotion based on success metrics
- Escalation metrics tracking and analysis
- Runtime contract validation

### Added - Feature F010: SDP Infrastructure (5 workstreams)
- Project Map with PRD v2.0 format
- Command reference documentation
- Configuration file support (~/.sdp/config.toml)
- Usage examples and interactive workflows
- Error handling with recovery strategies

### Added - Feature F011: PRD Command (6 workstreams)
- PRD command with project type profiles (cli, service, library)
- Line limits validator (PRD section constraints)
- Annotation parser (@prd_flow, @prd_step decorators)
- Diagram generator (Mermaid, PlantUML)
- Codereview hook integration for PRD validation
- hw_checker PRD migration

### Changed
- Enhanced GitHub bidirectional sync service
- Improved project fields integration
- Workstream ID format changed to PP-FFF-WW
- Index tracking with completion percentages
- Pre-deploy hooks adapted for SDP
- Project map parser supports multiple title formats

### Fixed
- Oneshot premature stop bug (explicit backlog count check)
- Project map parsing with `# PROJECT_MAP:` format
- Session quality check hook path resolution
- Pre-edit check validation

### Infrastructure
- 204 unit tests, 16 integration tests
- 88% average test coverage
- Full mypy --strict type checking
- Ruff linting with SDP rules
- Clean Architecture compliance (Domain-App-Infra-Presentation)

## [0.3.1] - 2026-01-12

### Added
- `docs/PRINCIPLES.md` - SOLID, DRY, KISS, YAGNI, TDD principles
- `docs/concepts/` - Clean Architecture, Artifacts, Roles documentation
- `README_RU.md` - Russian translation of README

### Removed
- `archive/` directory - legacy v1.2 materials cleaned up
- `IMPLEMENTATION_SUMMARY.md` - redundant with PROTOCOL.md

### Changed
- Simplified README.md with clearer structure
- Updated CLAUDE.md with links to new docs

## [2.0.0] - 2025-12-31

### Added
- **Unified Progressive Consensus (UPC) Protocol** - Complete rewrite
- **Three Protocol Tiers**: Starter, Standard, Enterprise
- **Three Execution Modes**: full, fast_track, hotfix
- **JSON Schemas** for all artifacts (`consensus/schema/`)
- **Centralized state** (`status.json`) as single source of truth
- **Workstreams** for micro-task tracking (merged from kanban concept)
- **Universal agent prompts** (`consensus/prompts/`)
- **Validation scripts** (`consensus/scripts/validate.py`)
- **Epic initialization** (`consensus/scripts/init.py`)
- **ADR-0004** documenting the unified protocol design

### Changed
- Protocol now **schema-driven** (JSON Schema is law)
- **Phase transitions** are explicit and validated
- **Agent prompts** are now portable Markdown (work with any LLM)
- **Directory structure** simplified and standardized

### Removed
- Legacy `modes/` directory (archived to `archive/v1.2/`)
- Legacy `prompts/` directory (archived to `archive/v1.2/`)
- Legacy `WORKFLOW.md`, `CONCEPTS.md` (archived)
- Implicit state management (replaced by `status.json`)

### Migration
See [PROTOCOL.md](docs/PROTOCOL.md) for migration instructions from v1.2.

## [1.2.0] - 2024-12-29

### Added
- Continuous code review after each workstream
- Duplication prevention rules
- Early quality gates
- Cross-epic code review
- Strict code review at epic completion
- No error-hiding fallbacks rule

## [1.1.0] - 2024-11-XX

### Added
- JSON format for inbox messages (compact keys)
- Per-epic directory structure
- Decision logs in Markdown

## [1.0.0] - 2024-11-XX

### Added
- Initial file-based consensus protocol
- Agent roles and responsibilities
- Veto protocol
- Clean Architecture enforcement
