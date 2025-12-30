# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2024-12-31

### Added
- **Three-mode architecture**: Solo, Structured, Multi-Agent modes for different project sizes
- **Solo mode** (`modes/solo/`)
  - README with workflow guide
  - CLAUDE.md.example template
  - Prompt examples for common tasks
- **Structured mode** (`modes/structured/`)
  - Four-phase workflow (Analyze → Design → Implement → Review)
  - Phase-specific prompts with templates
  - CLAUDE.md.example for structured projects
- **Multi-Agent mode** (`modes/multi-agent/`)
  - README explaining full consensus protocol
  - References existing prompts/
- **Concepts documentation** (`concepts/`)
  - Roles in software development
  - Artifacts guide (specs, designs, ADRs, reports)
  - ADR templates and examples
  - Clean Architecture explanation
- **Worked examples** (`examples/`)
  - Solo bug fix example
  - Structured feature example (password reset)
  - Multi-Agent epic example (notification system)
- **CONCEPTS.md** - Overview of core concepts

### Changed
- **README.md** - Complete rewrite with mode selection guide
- **CLAUDE.md** - Simplified, removed fabricated model claims
- **MODELS.md** - Rewritten with principle-based recommendations (removed unverified benchmarks)
- **docs/guides/CLAUDE_CODE.md** - Simplified, clarified Claude-only limitation
- **docs/guides/CURSOR.md** - Simplified, removed fabricated settings

### Removed
- Fabricated SWE-bench percentages from documentation
- Non-existent model names (Gemini 3 Flash, Kimi K2 Thinking, etc.)
- Fabricated Cursor settings and multi-provider claims
- Over-complicated configuration examples

### Fixed
- Documentation now reflects actual tool capabilities
- Model recommendations based on principles, not fabricated benchmarks

## [1.2.0] - Previous

### Added
- Continuous code review after each workstream
- Duplication prevention rules
- Early quality gates
- Cross-epic code review

---

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
