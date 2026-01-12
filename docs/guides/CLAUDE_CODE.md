# Claude Code Integration

Guide for using Spec-Driven Protocol (SDP) with [Claude Code](https://claude.ai/code).

> **📝 Meta-note:** Documentation developed with AI assistance (Claude Sonnet 4.5).

## Quick Start

Use **skills** (@ commands) for SDP workflow:

```
@idea "Add user authentication"
@design idea-user-auth
@build WS-001-01
@review F01
@deploy F01
```

## What is Claude Code?

Claude Code is Anthropic's official AI coding assistant:
- Interactive terminal-based interface
- Automatic codebase understanding
- File editing and command execution
- Project context via CLAUDE.md files

**Important:** Claude Code works **only with Claude models** (Anthropic).

## Setup

### Installation

```bash
# Via npm
npm install -g @anthropic-ai/claude-code

# Verify
claude --version
```

### Authentication

```bash
# Login (opens browser)
claude login

# Or set API key
export ANTHROPIC_API_KEY="sk-ant-..."
```

## CLAUDE.md File

Claude Code automatically reads `CLAUDE.md` from project root.

See [CLAUDE.md](../../CLAUDE.md) for this project's configuration.

## Available Skills

| Skill | Purpose | Example |
|-------|---------|---------|
| `@idea` | **Interactive requirements** (AskUserQuestion) | `@idea "Add payment processing"` |
| `@design` | **Interactive planning** (EnterPlanMode) | `@design idea-payments` |
| `@build` | Execute workstream (TodoWrite tracking) | `@build WS-001-01` |
| `@review` | Quality check | `@review F01` |
| `@deploy` | Production deployment | `@deploy F01` |
| `@issue` | Debug and route bugs | `@issue "Login fails on Firefox"` |
| `@hotfix` | Emergency fix (P0) | `@hotfix "Critical API outage"` |
| `@bugfix` | Quality fix (P1/P2) | `@bugfix "Incorrect totals"` |
| `@oneshot` | **Autonomous execution** (Task-based) | `@oneshot F01` or `@oneshot F01 --background` |

Skills are defined in `.claude/skills/{name}/SKILL.md`

**Claude Code Integration Highlights:**
- `@idea` — Deep interviewing via AskUserQuestion (explores tradeoffs, no obvious questions)
- `@design` — EnterPlanMode for codebase exploration + AskUserQuestion for architecture decisions
- `@build` — TodoWrite real-time progress tracking through TDD cycle
- `@oneshot` — Task tool spawns isolated orchestrator agent with background execution support

## Typical Workflow

### 1. Gather Requirements (Interactive)

```bash
# Start Claude Code
claude

# In session:
> @idea "Users need password reset via email"
```

**Claude uses AskUserQuestion for deep interviewing:**
- Technical approach (email service, token storage)
- UI/UX (where in app, error messages)
- Security (token expiry, rate limiting)
- Concerns (complexity, failure modes)

**Output:** `docs/drafts/idea-password-reset.md` (comprehensive spec)

### 2. Design Workstreams (Interactive Planning)

```
> @design idea-password-reset
```

**Claude enters Plan Mode:**
- Explores codebase (existing auth, email infrastructure)
- Asks architecture questions via AskUserQuestion (JWT vs sessions, etc.)
- Designs WS decomposition
- Requests approval via ExitPlanMode

**Output:**
- `docs/workstreams/backlog/WS-001-01-domain.md`
- `docs/workstreams/backlog/WS-001-02-service.md`
- `docs/workstreams/backlog/WS-001-03-api.md`
- etc.

### 3. Execute Workstreams

**Option A: Manual execution (with TodoWrite tracking)**
```
> @build WS-001-01
# Claude shows TodoWrite progress:
#   [in_progress] Pre-build validation
#   [pending] Write failing test (Red)
#   [pending] Implement minimum code (Green)
#   ... (updates in real-time)

> @build WS-001-02
> @build WS-001-03
```

**Option B: Autonomous execution via Task tool**
```
> @oneshot F01
# Spawns orchestrator agent with TodoWrite tracking
# Executes all WS with PR approval gate

> @oneshot F01 --background
# Run in background for large features

> @oneshot F01 --resume {agent_id}
# Resume from checkpoint if interrupted
```

### 4. Review Quality

```
> @review F01
```

Checks:
- ✅ All acceptance criteria met
- ✅ Coverage ≥80%
- ✅ No TODO/FIXME
- ✅ Clean Architecture followed

### 5. Deploy

```
> @deploy F01
```

Generates:
- Docker configs
- CI/CD pipelines
- Release notes
- Deployment plan

## Model Selection

Switch models using `/model` command:

```
/model opus    # Claude Opus 4.5 - best reasoning
/model sonnet  # Claude Sonnet 4.5 - balanced
/model haiku   # Claude Haiku 4.5 - fastest
```

### Recommended by Skill

| Skill | Model | Why |
|-------|-------|-----|
| `@idea` | Opus | Requirements analysis |
| `@design` | Opus | Workstream decomposition |
| `@build` | Sonnet | Code implementation |
| `@review` | Sonnet | Quality checks |
| `@deploy` | Sonnet/Haiku | Config generation |
| `@oneshot` | Opus | Autonomous orchestration (Task tool) |

See [MODELS.md](../../MODELS.md) for detailed recommendations.

## File Structure

```
project/
├── CLAUDE.md             # Claude Code config (auto-loaded)
├── .claude/
│   ├── skills/           # Skill definitions
│   ├── agents/           # Multi-agent mode (advanced)
│   └── settings.json     # Settings
├── docs/
│   ├── drafts/           # @idea outputs
│   ├── workstreams/
│   │   ├── backlog/      # @design outputs
│   │   ├── in_progress/  # @build working
│   │   └── completed/    # @build done
│   └── specs/            # Feature specs
├── prompts/commands/     # Full skill instructions
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

### Post-build
```bash
hooks/post-build.sh WS-001-01 project.module
```

### Pre-commit
```bash
hooks/pre-commit.sh
```

See [CURSOR.md](CURSOR.md) for hook details.

## Advanced Features

### Task Tool Integration (@oneshot)

`@oneshot` uses Claude Code's Task tool to spawn an isolated orchestrator agent:

**How it works:**
1. Main Claude spawns Task agent with orchestrator instructions
2. Agent executes all WS autonomously
3. Real-time progress via TodoWrite
4. PR approval gate before execution
5. Checkpoint/resume capability

**Background execution:**
```bash
> @oneshot F01 --background
# Agent runs in background
# Check progress: Read("/tmp/agent_{id}.log")
# Notification when complete
```

**Resume from interruption:**
```bash
> @oneshot F01 --resume {agent_id}
# Agent continues from last checkpoint
```

### AskUserQuestion Integration (@idea, @design)

**@idea** uses AskUserQuestion for deep requirements gathering:
- No obvious questions
- Explores tradeoffs
- Uncovers hidden requirements
- Technical and business concerns

**@design** uses EnterPlanMode + AskUserQuestion:
- Codebase exploration in Plan Mode
- Architecture decisions via AskUserQuestion
- Approval workflow via ExitPlanMode

### TodoWrite Progress Tracking (@build, @oneshot)

**@build** shows real-time progress:
- Pre-build validation
- TDD cycle (Red → Green → Refactor)
- Quality gates
- Execution report

**@oneshot** tracks high-level progress:
- PR approval
- Each WS execution
- Final review
- UAT guide generation

### Multi-Agent Mode (Legacy)

For complex features, use multi-agent orchestration:

```
> @orchestrator F01
```

Agents defined in `.claude/agents/`:
- `planner.md` — Breaks features into workstreams
- `builder.md` — Executes workstreams
- `reviewer.md` — Quality checks
- `deployer.md` — Production deployment
- `orchestrator.md` — Coordinates workflow

**Note:** `@oneshot` is preferred for autonomous execution (uses Task tool).

## Tips

1. **Keep CLAUDE.md concise**: 60-300 lines recommended
2. **Use /model command**: Switch for appropriate complexity
3. **Clear context**: Use `/clear` between major features
4. **Verify before saving**: Ask to show outputs before writing
5. **Follow skill instructions**: Each skill has specific requirements

## Troubleshooting

### Skill not found
Check `.claude/skills/{name}/SKILL.md` exists

### Validation fails
Run `hooks/pre-build.sh {WS-ID}` to see issues

### Workstream blocked
Check dependencies in `docs/workstreams/backlog/{WS-ID}.md`

### Coverage too low
Run `pytest --cov --cov-report=term-missing`

## Resources

| Resource | Purpose |
|----------|---------|
| [Claude Code Docs](https://docs.anthropic.com/claude/docs/claude-code) | Official documentation |
| [PROTOCOL.md](../../PROTOCOL.md) | Full SDP specification |
| [CLAUDE.md](../../CLAUDE.md) | Project configuration |
| [docs/PRINCIPLES.md](../../docs/PRINCIPLES.md) | Core principles |
| [MODELS.md](../../MODELS.md) | Model recommendations |

---

**Version:** SDP 0.3.0  
**Claude Code Compatibility:** 0.3+  
**Mode:** Skill-based, one-shot execution
