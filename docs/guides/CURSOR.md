# Cursor IDE Integration

Guide for using Spec-Driven Protocol (SDP) with [Cursor IDE](https://cursor.com).

> **📝 Meta-note:** Documentation developed with AI assistance (Claude Sonnet 4.5).

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

## Advanced Features

### Composer (Multi-file Editing)

**Use Composer for related files during `/build`:**

When implementing features that span multiple files, use Composer to edit them simultaneously:

```
@src/domain/user.py @src/application/get_user.py @tests/test_get_user.py
"Implement GetUser use case with TDD: write test first, then implementation"
```

**Benefits:**
- Edit related files in one operation
- Maintain consistency across layers
- Faster iteration

**When to use:**
- Domain entity + Application use case + Tests
- Service + Repository + Tests
- Multiple related refactorings
- Cross-layer changes

### @file References

**Always use @file references for context:**

Each command has recommended @file references documented in `prompts/commands/{command}.md`.

**Example for `/build`:**
```
@docs/workstreams/backlog/WS-001-01.md
@PROJECT_CONVENTIONS.md
@docs/workstreams/INDEX.md
```

**Benefits:**
- Explicit file context
- Better token management
- Clearer AI understanding

### Terminal Integration

**Use Cursor's built-in terminal for validation:**

Instead of external terminal, use Cursor's integrated terminal (`` Ctrl+` `` or View → Terminal):

```bash
# Pre-build validation
hooks/pre-build.sh WS-001-01

# Post-build validation
hooks/post-build.sh WS-001-01 project.module

# Run tests with coverage
pytest tests/unit/test_module.py --cov=src/module --cov-report=term-missing

# Check linting
ruff check src/module/
mypy src/module/ --strict
```

**Benefits:**
- See results inline
- Copy output easily
- Run commands without leaving Cursor

**Keyboard shortcuts:**
- `` Ctrl+` `` — Toggle terminal
- `Ctrl+Shift+` `` — Create new terminal
- `Ctrl+Shift+K` — Clear terminal

### Git UI Integration

**Use Cursor's Git UI for visual workflow:**

#### Visual Diff Before Commit

1. Open Source Control panel (`Ctrl+Shift+G`)
2. Review changes visually
3. Check for:
   - Clean Architecture violations (domain imports)
   - Test coverage (new files have tests)
   - File size (< 200 LOC)

#### Branch Management

1. Click branch name in status bar
2. Create new branch: `feature/{slug}`
3. Switch branches visually
4. See branch history

#### Staging Area

1. Stage files selectively (checkboxes)
2. Review staged changes
3. Write commit message inline
4. Commit with conventional format

**Example commit workflow:**
```
1. Stage src/ files → Commit: "feat(auth): WS-001-01 - implement domain layer"
2. Stage tests/ files → Commit: "test(auth): WS-001-01 - add unit tests"
3. Stage docs/ files → Commit: "docs(auth): WS-001-01 - execution report"
```

#### Visual Diff for Review

During `/review`, use Git UI to:
- Compare changes visually
- See line-by-line differences
- Check Clean Architecture boundaries
- Verify test coverage

**Access:**
- `Ctrl+Shift+G` — Source Control panel
- Click file → View diff
- Compare branches visually

## Tips

1. **Use command autocomplete**: Type `/` to see all available commands
2. **Use @file references**: Always include recommended files for context
3. **Use Composer for multi-file**: Edit related files simultaneously
4. **Use Terminal integration**: Run validation in Cursor terminal
5. **Use Git UI**: Visual diff and branch management
6. **Keep context focused**: Pin important files, clear between features
7. **Let hooks validate**: Don't bypass Git hooks
8. **Follow conventional commits**: `feat(scope): WS-XXX-YY - description`

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
