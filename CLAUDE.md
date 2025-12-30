# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the **Consensus Workflow** - a file-based multi-agent coordination framework for software development. It enables autonomous agents (Analyst, Architect, Tech Lead, Developer, QA, DevOps, SRE, Security, etc.) to collaborate through structured JSON messages and shared artifacts.

The framework also includes optional education roles (Lecture Assistant, Seminar Assistant) for preparing course materials.

**Note:** This repository is developed with AI assistance. See [MODELS.md](MODELS.md) for model recommendations.

## Core Architecture

### Consensus Protocol (v1.2)

Agents communicate through file-based protocol:
- **Artifacts**: Deliverables in `docs/specs/{epic}/consensus/artifacts/`
- **Messages**: JSON files in `docs/specs/{epic}/consensus/messages/inbox/{agent}/`
- **Decision Logs**: Markdown in `docs/specs/{epic}/consensus/decision_log/`

### Directory Structure

```
docs/specs/{epic}/
├── epic.md                    # Epic definition
├── architecture.md            # Architecture docs
├── implementation.md          # Implementation plan
├── testing.md                 # Testing strategy
├── deployment.md              # Deployment plan
└── consensus/
    ├── artifacts/             # Agent deliverables (JSON)
    ├── messages/inbox/{agent}/ # Agent communication
    └── decision_log/          # Decision history (MD)
```

### Agent Roles

1. **Analyst** - Requirements, scope, success metrics
2. **Architect** - System design, Clean Architecture, layer boundaries
3. **Tech Lead** - Implementation planning, code review orchestration
4. **Developer** - Implementation with TDD, incremental reviews
5. **QA** - Code quality, functional testing
6. **DevOps** - CI/CD, deployment
7. **SRE** - Observability, SLOs
8. **Security** - Threat modeling, auth reviews
9. **Data/ML Quality** - Data specs, ML artifacts
10. **Documentation Steward** - Docs curation
11. **Prompt Engineer** - LLM prompts, RAG

## Critical Protocol Rules

### Language and Format
- **ALL messages MUST be in English**
- JSON format with compact keys: `d`, `st`, `r`, `epic`, `sm`, `nx`, `artifacts`
- Message naming: `{YYYY-MM-DD}-{subject}.json`

### Inbox Rules
- **READ**: Only your own inbox
- **WRITE**: Only to OTHER agents' inboxes

### Quality Gates (Non-Negotiable)
- No silent fallbacks - errors must be explicit
- No layer violations - dependencies point inward
- Continuous code review - after each workstream
- Duplication prevention - search before implementing

### Veto Protocol (Cannot Override)
- Architecture violations (architect)
- Security issues (security)
- Missing rollback plan (devops)
- Code review violations (tech_lead, qa)

## Engineering Principles

### Clean Architecture (Enforced)
```
Presentation → Infrastructure → Application → Domain
```
Dependencies MUST point inward.

### Forbidden Patterns
- `except: pass` (silent failures)
- Default values masking exceptions
- Catch-all error handlers hiding errors

### Required
- Explicit error logging
- Proper exception raising
- Test coverage ≥80% in touched areas

## Model Selection

See [MODELS.md](MODELS.md) for guidance. General principle:
- **Strategic roles** (Analyst, Architect, Security): Use most capable model
- **Implementation roles** (Developer, QA, DevOps): Faster models work well

## Quick Start

1. Create epic: `docs/specs/epic_XX/epic.md`
2. Run Analyst → `requirements.json`
3. Run Architect → `architecture.json` (may veto)
4. Run Tech Lead → `implementation.md`
5. Run Developer → code + tests (continuous review)
6. Run QA → `test_results.md`
7. Run DevOps → deployment

See [modes/](modes/) for workflow guides and [examples/](examples/) for tutorials.

## Integration

- [Claude Code Guide](docs/guides/CLAUDE_CODE.md)
- [Cursor Guide](docs/guides/CURSOR.md)

## Key Files

- `PROTOCOL.md` - Full consensus protocol v2.0
- `RULES_COMMON.md` - Shared rules for all agents
- `MODELS.md` - Model selection guide
- `prompts/` - Agent prompt templates
- `docs/specs/epic_sample/` - Complete example

## Self-Verification Checklist

Before completing work:
- [ ] Clean Architecture boundaries respected
- [ ] Engineering principles followed (DRY, SOLID)
- [ ] No layer violations
- [ ] No fallbacks hiding errors
- [ ] Documentation updated
- [ ] All messages in English
