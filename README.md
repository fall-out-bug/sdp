# Spec-Driven Protocol (SDP)

Workstream-driven development framework for AI agents with structured Developer-in-the-Loop workflow.

## Overview

SDP enables AI agents (Claude Code, Cursor, OpenCode) to execute software development through atomic **workstreams** — small, self-contained tasks that can be completed in one shot.

**Key concepts:**
- **Workstream (WS):** Atomic task, executable in one shot (`WS-001-01`)
- **Feature:** Group of related workstreams (`F001`, 5-30 WS)
- **Release:** Product milestone (10-30 features)

**Scope metrics (no time estimates!):**
- **SMALL:** < 500 LOC
- **MEDIUM:** 500-1500 LOC
- **LARGE:** > 1500 LOC → split into 2+ WS

## Quick Start

```bash
# Install
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
poetry install

# Verify
sdp --version  # Output: sdp version 0.4.0
sdp --help
```

### As Submodule (Recommended)

```bash
# Add to your project
git submodule add git@github.com:fall-out-bug/sdp.git sdp

# Update to latest
git submodule update --remote sdp
```

## Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `/idea` | Interactive requirements gathering | `/idea "Add user auth"` |
| `/design` | Plan workstreams | `/design idea-user-auth` |
| `/build` | Execute single workstream | `/build WS-001-01` |
| `/oneshot` | Autonomous feature execution | `/oneshot F001` |
| `/review` | Quality check | `/review F001` |
| `/deploy` | Production deployment | `/deploy F001` |
| `/issue` | Debug and route bugs | `/issue "Login fails"` |
| `/hotfix` | Emergency fix (P0) | `/hotfix "Critical bug"` |
| `/bugfix` | Quality fix (P1/P2) | `/bugfix "Incorrect totals"` |

## Project Structure

```
sdp/
├── src/sdp/          # Source code
│   ├── core/         # Workstream parser, feature decomposition
│   ├── github/       # GitHub integration
│   ├── prd/          # PRD command
│   └── validators/   # Quality checks
├── prompts/          # Command prompts
│   └── commands/     # /idea, /design, /build, /review, etc.
├── docs/             # Documentation
├── tests/            # Test suite
└── hooks/            # Git hooks
```

## Quality Gates

| Gate | Requirement | Check |
|------|-------------|-------|
| AI-Readiness | Files < 200 LOC, CC < 10 | `ruff check --select=C901` |
| Test Coverage | ≥ 80% | `pytest --cov-fail-under=80` |
| Type Checking | mypy strict | `mypy --strict` |
| Clean Architecture | No layer violations | Manual review |
| Error Handling | No `except: pass` | `grep -r "except:"` |

## Features (v0.4.0)

- ✅ F003: Two-Stage Review (5 WS)
- ✅ F004: Platform Adapters (4 WS) - Claude Code, Cursor, OpenCode
- ✅ F005: Extension System (3 WS)
- ✅ F006: Core SDP (6 WS)
- ✅ F007: Oneshot & Hooks (10 WS)
- ✅ F008: Contract-Driven WS Tiers (9 WS)
- ✅ F010: SDP Infrastructure (5 WS)
- ✅ F011: PRD Command (6 WS)

**Total:** 48/58 workstreams completed (83%)

## Documentation

- [PROTOCOL.md](PROTOCOL.md) - Full protocol specification (Russian)
- [docs/multi-ide-parity.md](docs/multi-ide-parity.md) - Multi-IDE support
- [docs/workstreams/INDEX.md](docs/workstreams/INDEX.md) - Workstream index
- [CHANGELOG.md](CHANGELOG.md) - Release notes

## Inspiration

SDP is inspired by and incorporates ideas from:

- [superpowers](https://github.com/obra/superpowers) - AI agent capabilities framework
- [vibe-kanban](https://github.com/BloopAI/vibe-kanban) - Task management for AI agents
- [sub-agents.directory](https://github.com/ayush-that/sub-agents.directory) - Multi-agent orchestration

## License

MIT License - see LICENSE file for details

## Support

- GitHub Issues: https://github.com/fall-out-bug/sdp/issues
- Documentation: https://github.com/fall-out-bug/sdp/tree/main/docs

---

**Version:** 0.4.0 | **Status:** ✅ Active | **Last Updated:** 2026-01-25
