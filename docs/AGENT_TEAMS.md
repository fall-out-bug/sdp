# AGENT_TEAMS Integration Guide

## Enable AGENT_TEAMS

Add to `settings.json`:

```json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  },
  "teammateMode": "auto"
}
```

**Display modes:**
- `auto` - split panes в tmux, иначе in-process
- `in-process` - все в одном терминале (Shift+Up/Down для навигации)
- `tmux` - split panes (требует tmux или iTerm2)

## SDP Agent Teams Configuration

### Discovery Team (@feature)

Когда пользователь вызывает `@feature "Add authentication"`:

```markdown
Create an agent team for feature discovery. Spawn 4 teammates:
- Business Analyst - gather requirements, user stories, KPIs
- Product Manager - define vision, roadmap, RICE prioritization
- Systems Analyst - specify functional requirements, APIs, data models
- Technical Decomposition - break into workstreams, estimate effort

Each teammate works independently. They can message each other directly.
Lead synthesizes their findings into docs/drafts/FXXX.md.

Beads integration: if enabled, teammates create/update tasks.
```

**Teammate spawning:**
```python
# Lead (main session) creates teammates:
Teammate-1 (Business Analyst): .claude/agents/business-analyst.md
Teammate-2 (Product Manager): .claude/agents/product-manager.md
Teammate-3 (Systems Analyst): .claude/agents/systems-analyst.md
Teammate-4 (Technical Decomposition): .claude/agents/technical-decomposition.md
```

### Design Team (@design)

```markdown
Create an agent team for system design. Spawn 3 teammates:
- System Architect - architecture pattern, tech stack, ADRs
- Security - threat model, auth design, compliance
- SRE - SLOs, monitoring strategy, incident response

Each reviews the design from their perspective.
They discuss tradeoffs and converge on a shared architecture.
Lead synthesizes into docs/designs/FXXX.md.
```

### Review Team (@review)

```markdown
Create an agent team for quality review. Spawn 5 teammates:
- QA - test coverage, quality metrics, quality gates
- Security - vulnerabilities, security controls, compliance
- DevOps - CI/CD pipeline, infrastructure, deployment
- SRE - SLOs, monitoring, disaster recovery
- Tech Lead - code quality, architecture decisions

Each reviews from their specialty.
They coordinate findings: any FAIL = overall CHANGES_REQUESTED.
Only if all 5 PASS → APPROVED.
```

### Implementation Team (@oneshot)

```markdown
Create an agent team for feature execution. Spawn N teammates:
- One teammate per workstream (e.g., 13 teammates for F050)

Each teammate executes @build for their workstream:
Teammate-1: @build 00-050-01
Teammate-2: @build 00-050-02
...

Teammates work IN PARALLEL on independent workstreams.
Shared task list tracks dependencies.
When WS completes, teammate claims next unblocked WS.

Lead orchestrates, handles errors, synthesizes final report.
```

## Workflow Comparison

### OLD (Subagents via Task tool):
```
@feature "Add auth"
  ↓
Main session spawns subagents via Task tool
  ↓
Subagents run, return results to main
  ↓
Main synthesizes
```

### NEW (AGENT_TEAMS):
```
@feature "Add auth"
  ↓
Main session creates AGENT_TEAM
  ↓
Spawns 4 teammate SESSIONS (separate Claude Code instances)
  ↓
Each teammate loads .claude/agents/{name}.md
  ↓
Teammates work independently, message each other
  ↓
Shared task list coordinates work
  ↓
Lead synthesizes results
```

## Benefits

1. **True parallelism** - teammates work simultaneously, not sequential
2. **Direct communication** - teammates can challenge each other
3. **Shared task list** - automatic dependency resolution
4. **Independent contexts** - each has own context window
5. **Interactive** - can message any teammate directly

## Configuration File

Create `.claude/teams/sdp-feature/config.json`:

```json
{
  "name": "sdp-feature",
  "description": "SDP feature discovery team",
  "members": [
    {
      "name": "business-analyst",
      "agent": ".claude/agents/business-analyst.md",
      "role": "Discovery",
      "focus": "Requirements, user stories, KPIs"
    },
    {
      "name": "product-manager",
      "agent": ".claude/agents/product-manager.md",
      "role": "Discovery",
      "focus": "Vision, roadmap, prioritization"
    },
    {
      "name": "systems-analyst",
      "agent": ".claude/agents/systems-analyst.md",
      "role": "Analysis",
      "focus": "Functional specs, APIs, data models"
    },
    {
      "name": "technical-decomposition",
      "agent": ".claude/agents/technical-decomposition.md",
      "role": "Analysis",
      "focus": "Workstreams, dependencies, estimates"
    }
  ],
  "tasks": {
    "mode": "shared",
    "autoClaim": true,
    "dependencies": "automatic"
  },
  "beads": {
    "enabled": true,
    "detection": "bd --version && .beads/"
  }
}
```

## Next Steps

1. Update settings.json to enable AGENT_TEAMS
2. Test with simple team first (e.g., code review)
3. Update @feature/@design/@review/@oneshot skills to use AGENT_TEAMS
4. Create team config files
5. Document workflow for users

## Resources

- [Claude Code Agent Teams Documentation](https://code.claude.com/docs/en/agent-teams)
- [Building a C compiler with a team of parallel Claudes](https://www.anthropic.com/engineering/building-c-compiler)
- [Claude Code Best Practices for Agentic Coding](https://www.anthropic.com/engineering/claude-code-best-practices)
