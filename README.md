# Consensus Workflow

A file-based multi-agent coordination framework for software development. Enables autonomous agents (Analyst, Architect, Tech Lead, Developer, QA, DevOps, etc.) to collaborate through structured JSON messages and shared artifacts.

## Quick Start

See [QUICKSTART.md](QUICKSTART.md) for a 10-minute tutorial.

## Integration Guides

- [Claude Code Integration](docs/guides/CLAUDE_CODE.md)
- [Cursor IDE Integration](docs/guides/CURSOR.md)
- [Model Recommendations](MODELS.md)

## Development Roles (Epic Workflow)

Every epic must traverse the following sequence. Agents must stay within their lane unless another role explicitly requests collaboration:

1. **Analyst** – clarifies scope, success metrics, and user stories, then records requirements.
2. **Architect** – derives system design, validates Clean Architecture boundaries, produces diagrams.
3. **Tech Lead** – translates architecture into an executable plan (tasks, sequencing, testing, deployment).
4. **Developer** – implements via TDD, keeps documentation updated, reports blockers.
5. **QA** – verifies acceptance criteria, coverage, and regression safety.
6. **DevOps** – maintains build/deploy pipelines, container images, secrets, and runtime environments.
7. **SRE** – defines telemetry, alerting, reliability signals, and runbooks.
8. **Security** – conducts threat modeling, secret audits, auth/authorization reviews.
9. **Data & ML Quality** – keeps data specs, validation logic, and ML artifacts accurate.
10. **Documentation Steward** – curates specs, READMEs, indexes, and consensus logs.
11. **Prompt & Context Engineer** – manages MCP/LLM prompts, RAG corpora, and safety rails.

## Education Roles (Optional)

Additional roles for preparing educational content:

12. **Lecture Assistant** – helps prepare lectures: structure, slides, code examples, talking points.
13. **Seminar Assistant** – helps prepare practical sessions: assignments, tests, grading criteria, solutions.

## Artifacts

Development artifacts are stored under `docs/specs/epic_xx/consensus/`.

See [docs/examples/sample_epic/](docs/examples/sample_epic/) for a complete worked example.

## Usage

Use the prompts in `prompts/` when spinning up an agent session to guarantee consistent expectations and outputs.

### Full Prompts (1500-2500 tokens)
For complex epics, new team members, or when full context needed:
```bash
cat prompts/analyst_prompt.md
```

### Quick Prompts (200-400 tokens)
For routine tasks, experienced users, or context-limited sessions:
```bash
cat prompts/quick/analyst_quick.md
```

### Token Comparison

| Role | Full | Quick | Savings |
|------|------|-------|---------|
| analyst | ~1500 | ~250 | 83% |
| architect | ~2500 | ~300 | 88% |
| tech_lead | ~2000 | ~350 | 82% |
| developer | ~2000 | ~350 | 82% |
| qa | ~1800 | ~300 | 83% |
| devops | ~1500 | ~250 | 83% |

**Quick prompts reference `RULES_COMMON.md` for shared rules.**

### Meta-Orchestrator (Simple Epics)
For epics with ≤3 workstreams, use single-session orchestrator:
```bash
cat prompts/meta_orchestrator_prompt.md
# One chat instead of 6, ~60% faster
```

## Key Documentation

- [PROTOCOL.md](PROTOCOL.md) – Full consensus protocol specification (v1.2)
- [RULES_COMMON.md](RULES_COMMON.md) – Shared rules for all agents
- [USER_GUIDE.md](USER_GUIDE.md) – Operational workflow and best practices
- [MODELS.md](MODELS.md) – Model selection guide
- [CONTROL_PLAYBOOK.md](CONTROL_PLAYBOOK.md) – Recovery workflow for complex situations
- [consensus_architecture.json](consensus_architecture.json) – Technical architecture specification

## Core Principles

- **Clean Architecture** – Layer boundaries are sacred (Domain → Application → Infrastructure → Presentation)
- **English Only** – All messages and artifacts must be in English
- **No Silent Failures** – Errors must be explicit, logged, or raised
- **Continuous Code Review** – Review after each workstream, not just at epic completion
- **Veto Power** – Architecture violations and security issues cannot be overridden

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
