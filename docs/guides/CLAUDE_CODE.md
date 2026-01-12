# Claude Code Integration

This guide explains how to use Consensus Workflow with [Claude Code](https://claude.ai/code).

## What is Claude Code?

Claude Code is Anthropic's official CLI for Claude. It provides:
- Interactive terminal-based AI coding assistant
- Automatic codebase understanding
- File editing and command execution capabilities
- Project context via CLAUDE.md files

**Important**: Claude Code works **only with Claude models** (Anthropic). It does not support other providers like Google AI, OpenAI, or local models.

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

Claude Code automatically reads `CLAUDE.md` from your project root for context and instructions.

### Best Practices (from official docs)

- **Keep it short**: 60-300 lines recommended
- **Focus on universal rules**: Instructions should apply to most tasks
- **Use clear headings**: Structure with # and ## for organization
- **Include**:
  - Code style guidelines
  - Testing instructions
  - Repository conventions
  - Common commands

### Example CLAUDE.md for Consensus Workflow

```markdown
# Project Instructions

## Consensus Protocol
This project uses Consensus Workflow for multi-agent coordination.
See PROTOCOL.md for full specification.

## Key Rules
- ALL output must be in English
- JSON messages use compact keys: d, st, r, feature, sm, nx, artifacts
- No silent fallbacks (except: pass forbidden)
- Clean Architecture: dependencies point inward

## Artifacts Location
- docs/specs/{epic}/consensus/artifacts/
- docs/specs/{epic}/consensus/messages/inbox/{agent}/
- docs/specs/{epic}/consensus/decision_log/

## Testing
Run: pytest tests/ --cov=src
Target: ≥80% coverage in touched areas

## Agent Prompts
Full prompts: prompts/{role}_prompt.md
Quick prompts: prompts/quick/{role}_quick.md
```

## Model Selection

Switch models using the `/model` command in Claude Code:

```
/model opus    # Claude Opus 4.5 - best reasoning
/model sonnet  # Claude Sonnet 4.5 - balanced
/model haiku   # Claude Haiku 4.5 - fastest
```

### Recommended for Consensus Roles

| Role | Recommended Model | Why |
|------|-------------------|-----|
| Analyst | Opus | Complex requirements analysis |
| Architect | Opus | Architecture decisions, veto logic |
| Tech Lead | Sonnet | Implementation planning |
| Developer | Sonnet or Haiku | Code implementation |
| QA | Sonnet or Haiku | Test verification |
| DevOps | Haiku | Configuration generation |

## Workflow

### Using Agent Prompts

```bash
# Start Claude Code in your project
claude

# In Claude Code session:
> Read prompts/analyst_prompt.md and analyze docs/specs/epic_XX/epic.md
> Create requirements.json following the protocol

> Now read prompts/architect_prompt.md
> Review requirements.json and create architecture.json
> Veto if Clean Architecture is violated
```

### Sequential Agent Approach

Run each agent role, switching context:

1. **Analyst** (use /model opus)
   ```
   Read prompts/analyst_prompt.md
   Analyze epic.md, create requirements.json
   ```

2. **Architect** (use /model opus)
   ```
   Read prompts/architect_prompt.md
   Review requirements, create architecture.json
   VETO if layer violations detected
   ```

3. **Tech Lead** (use /model sonnet)
   ```
   Read prompts/tech_lead_prompt.md
   Create implementation.md with workstreams
   ```

4. **Developer** (use /model sonnet)
   ```
   Read prompts/developer_prompt.md
   Implement workstream with TDD
   ```

5. **QA** (use /model haiku)
   ```
   Read prompts/qa_prompt.md
   Run tests, verify coverage
   ```

6. **DevOps** (use /model haiku)
   ```
   Read prompts/devops_prompt.md
   Create deployment configuration
   ```

## Tips

1. **Keep CLAUDE.md short**: Long instructions may be ignored
2. **Use /model command**: Switch to appropriate model for each role
3. **Clear context**: Use `/clear` between major agent switches
4. **Verify JSON format**: Ask Claude to show JSON before saving
5. **Use quick prompts**: For routine tasks, reference prompts/quick/

## Limitations

- **Single provider**: Only Claude models (no Gemini, GPT, etc.)
- **No parallel agents**: Run agents sequentially
- **Context limits**: Very long conversations may lose context

## Resources

- [Claude Code Documentation](https://docs.anthropic.com/claude/docs/claude-code)
- [PROTOCOL.md](../../PROTOCOL.md) - Consensus protocol specification
- [RULES_COMMON.md](../../RULES_COMMON.md) - Shared agent rules
- [prompts/](../../prompts/) - Agent prompt templates

---

**Note**: Claude Code features may change. Check official Anthropic documentation for current capabilities.
