# Consensus Workflow

An educational framework for AI-assisted software development. Learn to work effectively with AI coding assistants through structured approaches that scale with project complexity.

## Choose Your Mode

| Mode | Complexity | Time | Use When |
|------|------------|------|----------|
| [**Solo**](modes/solo/) | Low | <2h | Bug fixes, small features |
| [**Structured**](modes/structured/) | Medium | 2-8h | Features needing documentation |
| [**Multi-Agent**](modes/multi-agent/) | High | Days/Weeks | Large projects, teams |

### Quick Decision Guide

```
Task < 2 hours and < 10 files?
├─ Yes → Solo Mode
└─ No  → Need documentation or ADRs?
         ├─ Yes → Structured Mode
         └─ No  → Parallel work or audit trail?
                  ├─ Yes → Multi-Agent Mode
                  └─ No  → Structured Mode
```

## Quick Start

### Solo Mode (Start Here)
```bash
# 1. Add CLAUDE.md to your project
cp modes/solo/CLAUDE.md.example your-project/CLAUDE.md

# 2. Start AI assistant and describe your task
"Fix the bug where users can't login with + in email"

# 3. Iterate until done
```
See [modes/solo/](modes/solo/) for guide and prompts.

### Structured Mode
```bash
# 1. Phase 1 - Analyze requirements
"Create specification for password reset feature"

# 2. Phase 2 - Design solution
"Create technical design and ADR"

# 3. Phase 3 - Implement
"Implement following the design, TDD approach"

# 4. Phase 4 - Review
"Verify implementation meets all requirements"
```
See [modes/structured/](modes/structured/) for phase prompts.

### Multi-Agent Mode
```bash
# Run specialized agents sequentially
# Each agent has dedicated prompts in prompts/

Analyst    → requirements.md
Architect  → architecture.md (may veto)
Tech Lead  → implementation.md
Developer  → code + tests
QA         → test-report.md
DevOps     → deployment.md
```
See [modes/multi-agent/](modes/multi-agent/) for full protocol.

## Learn the Concepts

| Concept | What You'll Learn |
|---------|-------------------|
| [Roles](concepts/roles/) | What Analyst, Architect, Developer do |
| [Artifacts](concepts/artifacts/) | Specs, designs, ADRs, reports |
| [ADR](concepts/adr/) | Documenting architectural decisions |
| [Clean Architecture](concepts/clean-architecture/) | Layer separation and dependencies |

See [CONCEPTS.md](CONCEPTS.md) for overview.

## Examples

| Example | Mode | Scenario |
|---------|------|----------|
| [Bug Fix](examples/solo-bug-fix/) | Solo | Fix email validation bug |
| [Feature](examples/structured-feature/) | Structured | Add password reset |
| [Epic](examples/multi-agent-epic/) | Multi-Agent | Build notification system |

## Integration Guides

- [Claude Code](docs/guides/CLAUDE_CODE.md) - Using with Claude Code CLI
- [Cursor](docs/guides/CURSOR.md) - Using with Cursor IDE
- [Model Selection](MODELS.md) - Choosing models for different tasks

## Core Principles

These apply to all modes:

- **Clean Architecture** - Dependencies point inward
- **No Silent Failures** - Errors must be explicit
- **Test Coverage** - ≥80% for new code
- **English Only** - All outputs in English

## Multi-Agent Protocol (Advanced)

For full multi-agent workflow:

| Resource | Description |
|----------|-------------|
| [PROTOCOL.md](PROTOCOL.md) | Complete protocol specification |
| [RULES_COMMON.md](RULES_COMMON.md) | Shared rules for all agents |
| [prompts/](prompts/) | Full agent prompts |
| [prompts/quick/](prompts/quick/) | Shorter prompts for routine tasks |

### Agent Roles

| Role | Responsibility | Can Veto? |
|------|----------------|-----------|
| Analyst | Requirements, scope | No |
| Architect | System design | Yes |
| Tech Lead | Planning, code review | Yes |
| Developer | Implementation, TDD | No |
| QA | Testing, verification | Yes |
| DevOps | Deployment | Yes |

Additional roles: Security, SRE, Data/ML, Documentation Steward.

## Repository Structure

```
consensus/
├── modes/              # Three workflow modes
│   ├── solo/           # Single agent, simple tasks
│   ├── structured/     # Phased approach with artifacts
│   └── multi-agent/    # Full consensus protocol
├── concepts/           # Educational content
│   ├── roles/          # Development roles explained
│   ├── artifacts/      # Document types
│   ├── adr/            # Architecture Decision Records
│   └── clean-architecture/
├── examples/           # Complete worked examples
├── prompts/            # Agent prompts for multi-agent mode
│   └── quick/          # Shorter versions
└── docs/
    ├── guides/         # Tool integration guides
    └── specs/          # Sample specifications
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
