# Spec Driven Protocol (SDP)

Workstream-driven development protocol for AI agents with structured, one-shot execution.

[Русская версия](README_RU.md)

---

## Core Idea

**Workstream** = atomic task that AI completes in one shot, without iterative loops.

```
Feature → Workstreams → One-shot execution → Done
```

## Terminology

| Term | Scope | Size | Example |
|------|-------|------|---------|
| **Release** | Product milestone | 10-30 Features | R1: MVP |
| **Feature** | Large capability | 5-30 Workstreams | F1: User Auth |
| **Workstream** | Atomic task | SMALL/MEDIUM/LARGE | WS-001: Domain entities |

**Scope metrics:**
- **SMALL**: < 500 LOC, < 1500 tokens
- **MEDIUM**: 500-1500 LOC, 1500-5000 tokens
- **LARGE**: > 1500 LOC → **split into 2+ WS**

**NO time-based estimates.** Use LOC/tokens only.

## Workflow

Use slash commands for streamlined execution:

```bash
/idea "User authentication"      # 1. Requirements gathering
/design idea-user-auth           # 2. Create workstreams
/build WS-001-01                 # 3. Implement workstream
/review F01                      # 4. Quality review
/deploy F01                      # 5. Deploy to production
```

## Quick Start

### 1. Gather Requirements

```bash
/idea "Add user authentication with email/password"
```

**Output:** `docs/drafts/idea-user-auth.md`

### 2. Create Workstreams

```bash
/design idea-user-auth
```

**Output:** 
- `docs/workstreams/backlog/WS-001-01-domain.md`
- `docs/workstreams/backlog/WS-001-02-repository.md`
- `docs/workstreams/backlog/WS-001-03-service.md`
- `docs/workstreams/backlog/WS-001-04-api.md`
- `docs/workstreams/backlog/WS-001-05-tests.md`

### 3. Implement Workstreams

```bash
/build WS-001-01  # Domain layer
/build WS-001-02  # Repository
/build WS-001-03  # Service
# ... etc
```

Or use autonomous execution:

```bash
/oneshot F01  # Executes all WS automatically
```

### 4. Review Quality

```bash
/review F01
```

Checks:
- ✅ All acceptance criteria met
- ✅ Coverage ≥ 80%
- ✅ No TODO/FIXME
- ✅ Clean Architecture followed

### 5. Deploy

```bash
/deploy F01
```

Generates:
- Docker configs
- CI/CD pipelines
- Release notes
- Deployment plan

## Commands Reference

| Command | Purpose | When to Use |
|---------|---------|-------------|
| `/idea` | Requirements gathering | Start new feature |
| `/design` | Create workstreams | After requirements clear |
| `/build` | Implement workstream | Execute single WS |
| `/review` | Quality check | After all WS complete |
| `/deploy` | Production deployment | After review APPROVED |
| `/issue` | Debug and route | Analyze bugs |
| `/hotfix` | Emergency fix | P0 production issue |
| `/bugfix` | Quality fix | P1/P2 bugs |
| `/oneshot` | Autonomous execution | Execute all WS hands-free |

## Quality Gates

| Gate | Requirements |
|------|--------------|
| **AI-Readiness** | Files < 200 LOC, CC < 10, type hints |
| **Clean Architecture** | No layer violations |
| **Error Handling** | No `except: pass` |
| **Test Coverage** | ≥ 80% |
| **No TODOs** | All completed or new WS |

## Core Principles

| Principle | Summary |
|-----------|---------|
| **SOLID** | SRP, OCP, LSP, ISP, DIP |
| **DRY** | Don't Repeat Yourself |
| **KISS** | Keep It Simple |
| **YAGNI** | Build only what's needed |
| **TDD** | Tests first (Red → Green → Refactor) |
| **Clean Code** | Readable, maintainable |
| **Clean Architecture** | Dependencies point inward |

See [docs/PRINCIPLES.md](docs/PRINCIPLES.md) for details.

## File Structure

```
sdp/
├── PROTOCOL.md              # Full specification
├── CODE_PATTERNS.md         # Implementation patterns
├── RULES_COMMON.md          # Common rules
├── docs/
│   ├── PRINCIPLES.md        # SOLID, DRY, KISS, YAGNI
│   └── concepts/            # Clean Architecture, Artifacts, Roles
├── prompts/
│   └── commands/            # Slash commands (/idea, /design, etc.)
├── schema/                  # JSON validation
├── scripts/                 # Utilities
├── hooks/                   # Git hooks
└── templates/               # Document templates
```

## Resources

| Resource | Purpose |
|----------|---------|
| [PROTOCOL.md](PROTOCOL.md) | Full specification |
| [docs/PRINCIPLES.md](docs/PRINCIPLES.md) | SOLID, DRY, KISS, YAGNI |
| [docs/concepts/](docs/concepts/) | Architecture concepts |
| [CODE_PATTERNS.md](CODE_PATTERNS.md) | Code patterns |
| [MODELS.md](MODELS.md) | Model recommendations |
| [CLAUDE.md](CLAUDE.md) | Claude Code integration |

## Integration

### For Claude Code

1. Copy files to your project:
```bash
cp -r prompts/ your-project/
cp -r schema/ your-project/
cp -r .claudecode/ your-project/
cp CLAUDE.md your-project/
```

2. Use skills: `@idea`, `@design`, `@build`, etc.

### For Cursor

1. Copy files to your project:
```bash
cp -r prompts/ your-project/
cp -r schema/ your-project/
cp -r .cursor/ your-project/
cp .cursorrules your-project/
```

2. Use slash commands: `/idea`, `/design`, `/build`, etc.

---

**Version:** 0.3.0 | **Status:** Active
