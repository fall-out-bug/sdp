# Changelog

All notable changes to the Spec-Driven Protocol (SDP).

> **📝 Meta-note:** Versions documented as they are released. Development is AI-assisted.

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
