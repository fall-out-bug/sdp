# Cursor IDE Integration

Guide for using Spec-Driven Protocol (SDP) with [Cursor IDE](https://cursor.com).

## Quick Start

Cursor automatically reads `.cursorrules` for project-specific rules and context.

Use **slash commands** for SDP workflow:

```
/idea "Add user authentication"
/design idea-user-auth
/build WS-001-01
/review F01
/deploy F01
```

## Setup

1. **Open project in Cursor**
2. **Cursor auto-loads** `.cursorrules` from project root
3. **Commands auto-complete** from `.cursor/commands/*.md`

## Available Slash Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `/idea` | Requirements gathering | `/idea "Add payment processing"` |
| `/design` | Create workstreams | `/design idea-payments` |
| `/build` | Execute workstream | `/build WS-001-01` |
| `/review` | Quality check | `/review F01` |
| `/deploy` | Production deployment | `/deploy F01` |
| `/issue` | Debug and route bugs | `/issue "Login fails on Firefox"` |
| `/hotfix` | Emergency fix (P0) | `/hotfix "Critical API outage"` |
| `/bugfix` | Quality fix (P1/P2) | `/bugfix "Incorrect totals"` |
| `/oneshot` | Autonomous execution | `/oneshot F01` |

Commands are defined in `.cursor/commands/{command}.md`

## Typical Workflow

### 1. Gather Requirements

```
/idea "Users need password reset via email"
```

**Output:** `docs/drafts/idea-password-reset.md`

### 2. Design Workstreams

```
/design idea-password-reset
```

**Output:**
- `docs/workstreams/backlog/WS-001-01-domain.md`
- `docs/workstreams/backlog/WS-001-02-service.md`
- `docs/workstreams/backlog/WS-001-03-api.md`
- etc.

### 3. Execute Workstreams

**Option A: Manual execution**
```
/build WS-001-01
/build WS-001-02
/build WS-001-03
```

**Option B: Autonomous execution**
```
/oneshot F01
```

### 4. Review Quality

```
/review F01
```

Checks:
- ✅ All acceptance criteria met
- ✅ Coverage ≥80%
- ✅ No TODO/FIXME
- ✅ Clean Architecture followed

### 5. Deploy

```
/deploy F01
```

Generates:
- Docker configs
- CI/CD pipelines
- Release notes
- Deployment plan

## Model Selection

Cursor supports multiple AI models. Use Settings → Models to switch.

### Recommended by Command

| Command | Recommended Model | Why |
|---------|------------------|-----|
| `/idea` | Claude Opus/Sonnet | Requirements analysis |
| `/design` | Claude Opus/Sonnet | Workstream decomposition |
| `/build` | Claude Sonnet | Code implementation |
| `/review` | Claude Sonnet | Quality checks |
| `/deploy` | Claude Sonnet/Haiku | Config generation |
| `/oneshot` | Claude Opus | Autonomous orchestration |

See [MODELS.md](../../MODELS.md) for detailed recommendations.

## File Structure

```
project/
├── .cursorrules          # Project rules (auto-loaded)
├── .cursor/
│   ├── commands/         # Slash command definitions
│   └── worktrees.json    # Git worktree config
├── docs/
│   ├── drafts/           # /idea outputs
│   ├── workstreams/
│   │   ├── backlog/      # /design outputs
│   │   ├── in_progress/  # /build working
│   │   └── completed/    # /build done
│   └── specs/            # Feature specs
├── prompts/commands/     # Full command instructions
├── hooks/                # Git hooks (validation)
└── schema/               # JSON validation
```

## Quality Gates (Enforced)

| Gate | Requirement |
|------|-------------|
| **AI-Readiness** | Files < 200 LOC, CC < 10, type hints |
| **Clean Architecture** | No layer violations |
| **Error Handling** | No `except: pass` |
| **Test Coverage** | ≥80% |
| **No TODOs** | All tasks done or new WS |

## Git Hooks

Automatic validation via Git hooks:

### Pre-build
```bash
hooks/pre-build.sh WS-001-01
```

Checks:
- Workstream exists and READY
- Dependencies satisfied
- Previous WS completed

### Post-build
```bash
hooks/post-build.sh WS-001-01 project.module
```

Checks:
- Tests pass (coverage ≥80%)
- No TODO/FIXME
- Type hints complete
- Files < 200 LOC
- Clean Architecture compliance

### Pre-commit
```bash
hooks/pre-commit.sh
```

Ensures:
- Linting passes
- Tests pass
- No secrets
- Conventional commits

## Tips

1. **Use command autocomplete**: Type `/` to see all available commands
2. **Keep context focused**: Only @-mention relevant files
3. **Clear chat between features**: Fresh context for new work
4. **Let hooks validate**: Don't bypass Git hooks
5. **Follow conventional commits**: `feat(scope): WS-XXX-YY - description`

## Troubleshooting

### Command not found
Ensure `.cursor/commands/{command}.md` exists

### Validation fails
Run `hooks/pre-build.sh {WS-ID}` to see specific issues

### Workstream blocked
Check dependencies in `docs/workstreams/backlog/{WS-ID}.md`

### Coverage too low
Run `pytest --cov --cov-report=term-missing`

## Resources

| Resource | Purpose |
|----------|---------|
| [PROTOCOL.md](../../PROTOCOL.md) | Full SDP specification |
| [docs/PRINCIPLES.md](../../docs/PRINCIPLES.md) | Core principles |
| [CODE_PATTERNS.md](../../CODE_PATTERNS.md) | Code patterns |
| [MODELS.md](../../MODELS.md) | Model recommendations |
| [.cursorrules](../../.cursorrules) | Project rules |

---

**Version:** SDP 0.3.0  
**Cursor Compatibility:** Latest  
**Mode:** Slash commands, one-shot execution
