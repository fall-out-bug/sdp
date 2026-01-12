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

**Advanced diff features:**
- **Inline diff view**: Click file in Source Control → see changes inline
- **Compare with previous version**: Right-click file → "Compare with..."
- **Diff between branches**: `Ctrl+Shift+P` → "Git: Compare References"
- **Unified diff view**: See all changes in one view
- **Side-by-side diff**: Split view for easier comparison

**During `/build` workflow:**
1. After implementing code, open Git UI
2. Review diff before committing
3. Check for:
   - Type hints added
   - Tests included
   - No `except: pass`
   - File size < 200 LOC
4. Stage and commit if all checks pass

### Code Actions (Quick Fixes)

**Use Code Actions to auto-fix linting errors:**

Cursor provides quick fixes for common issues. Use them during quality gates:

**How to use:**
1. Open file with linting errors
2. Hover over error (yellow/red squiggle)
3. Click lightbulb icon 💡 or press `Ctrl+.`
4. Select "Fix all auto-fixable problems"

**Common fixes:**
- **Import organization**: Auto-sort imports (stdlib → third-party → local)
- **Type hints**: Add missing type hints
- **Unused imports**: Remove unused imports
- **Formatting**: Auto-format code
- **Missing docstrings**: Add docstring templates

**Integration with SDP:**

During `/build` post-build validation:
```bash
# Run linter
ruff check src/module/

# If errors found:
# 1. Open file in Cursor
# 2. Use Code Actions (Ctrl+.)
# 3. Fix all auto-fixable
# 4. Re-run validation
```

**Keyboard shortcuts:**
- `Ctrl+.` — Show Code Actions
- `Ctrl+Shift+.` — Quick Fix (next error)
- `F8` — Go to next error
- `Shift+F8` — Go to previous error

**Best practices:**
- Fix linting errors immediately during `/build`
- Use Code Actions before manual fixes
- Verify fixes don't break tests
- Commit fixes separately: `fix({scope}): WS-XXX-YY - fix linting errors`

### Chat Context Management

**Strategically manage chat context for better AI performance:**

#### Pin Important Files

**Always pin these files:**
- `PROJECT_CONVENTIONS.md` — Project-specific rules (always pinned)
- `docs/workstreams/INDEX.md` — For `/design` and `/build`
- `PROTOCOL.md` — For reference during any command

**How to pin:**
1. @-mention file in chat
2. Click pin icon 📌 next to file name
3. File stays in context across messages

**Command-specific pins:**

**For `/idea`:**
- `PROTOCOL.md`
- `templates/idea-draft.md`

**For `/design`:**
- `docs/PROJECT_MAP.md` (architecture decisions)
- `docs/workstreams/INDEX.md` (check duplicates)
- `PROJECT_CONVENTIONS.md`

**For `/build`:**
- `docs/workstreams/backlog/WS-{ID}-*.md` (WS plan)
- `PROJECT_CONVENTIONS.md`
- `CODE_PATTERNS.md`

**For `/review`:**
- `docs/workstreams/INDEX.md`
- `PROJECT_CONVENTIONS.md`
- `CODE_PATTERNS.md`

#### Clear Context Strategically

**When to clear chat:**
- ✅ Between different features
- ✅ After completing `/review`
- ✅ When switching commands (e.g., `/build` → `/review`)
- ✅ When context becomes too long (> 50 messages)

**When NOT to clear:**
- ❌ During single `/build` execution
- ❌ When debugging related issues
- ❌ During `/oneshot` autonomous execution

**How to clear:**
- Click "Clear Chat" button
- Or start new chat: `Ctrl+Shift+N`

#### Context Budget Management

**Monitor token usage:**
- Large files consume more tokens
- Use @file selectively
- Pin only essential files
- Clear old context regularly

**Tips:**
- Use file paths instead of full content when possible
- Reference specific sections: `@file.md#section`
- Clear context between major phases

### Workspace Settings

**Configure Cursor for optimal SDP workflow:**

Create `.vscode/settings.json` in your project:

```json
{
  // Python settings
  "python.linting.enabled": true,
  "python.linting.ruffEnabled": true,
  "python.linting.mypyEnabled": true,
  "python.formatting.provider": "ruff",
  "python.analysis.typeCheckingMode": "strict",
  
  // File size limits (SDP requirement)
  "files.maxMemoryForLargeFilesMB": 50,
  
  // Exclude patterns for indexing
  "files.exclude": {
    "**/__pycache__": true,
    "**/*.pyc": true,
    "**/.pytest_cache": true,
    "**/.mypy_cache": true,
    "**/node_modules": true
  },
  
  // Search exclude patterns
  "search.exclude": {
    "**/node_modules": true,
    "**/__pycache__": true,
    "**/.pytest_cache": true,
    "**/.mypy_cache": true,
    "**/venv": true,
    "**/.venv": true
  },
  
  // Editor settings
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": "explicit",
    "source.fixAll": "explicit"
  },
  
  // File associations
  "files.associations": {
    "*.md": "markdown",
    "*.mdc": "markdown"
  },
  
  // Git settings
  "git.enableSmartCommit": true,
  "git.confirmSync": false,
  "git.autofetch": true,
  
  // Terminal settings
  "terminal.integrated.defaultProfile.linux": "bash",
  "terminal.integrated.cwd": "${workspaceFolder}",
  
  // Cursor-specific
  "cursor.chat.maxContextLength": 50000,
  "cursor.chat.model": "claude-sonnet-4.5"
}
```

**SDP-specific settings:**

```json
{
  // Enforce file size limits
  "[python]": {
    "editor.rulers": [200],
    "editor.wordWrap": "wordWrapColumn",
    "editor.wordWrapColumn": 200
  },
  
  // Test discovery
  "python.testing.pytestEnabled": true,
  "python.testing.pytestArgs": [
    "tests",
    "--cov=src",
    "--cov-report=term-missing"
  ],
  
  // Type checking
  "python.analysis.typeCheckingMode": "strict",
  "python.analysis.diagnosticMode": "workspace"
}
```

**Installation:**
1. Copy settings to `.vscode/settings.json`
2. Cursor will auto-load settings
3. Adjust model preference in Cursor Settings → Models

**Benefits:**
- Auto-format on save
- Auto-fix linting errors
- Type checking enabled
- Test discovery configured
- File size warnings

## Tips

1. **Use command autocomplete**: Type `/` to see all available commands
2. **Use @file references**: Always include recommended files for context
3. **Use Composer for multi-file**: Edit related files simultaneously
4. **Use Terminal integration**: Run validation in Cursor terminal
5. **Use Git UI**: Visual diff and branch management
6. **Use Code Actions**: Auto-fix linting errors (`Ctrl+.`)
7. **Pin important files**: PROJECT_CONVENTIONS.md, INDEX.md always pinned
8. **Clear context strategically**: Between features, after `/review`
9. **Configure workspace settings**: Use `.vscode/settings.json` template
10. **Let hooks validate**: Don't bypass Git hooks
11. **Follow conventional commits**: `feat(scope): WS-XXX-YY - description`

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
