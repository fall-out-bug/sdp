# Changelog

All notable changes to the Spec-Driven Protocol (SDP).

> **📝 Meta-note:** Versions documented as they are released. Development is AI-assisted.

## [0.4.0] - 2026-01-25

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
See [PROTOCOL.md](PROTOCOL.md) for migration instructions from v1.2.

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
