# SDP Documentation Site Map

**Version:** v0.5.0
**Last Updated:** 2026-01-29
**Total Documents:** 80+ files

Complete navigation guide for all SDP documentation.

---

## 📑 Quick Navigation

- [Getting Started](#getting-started) - New user guides
- [Core Documentation](#core-documentation) - Primary references
- [Guides by Role](#guides-by-role) - Role-specific paths
- [Guides by Topic](#guides-by-topic) - Topic-specific deep dives
- [API & Integration](#api--integration) - Technical integrations
- [Reference](#reference) - Glossaries, schemas, templates
- [Architecture](#architecture) - Design documentation
- [Process & Workflow](#process--workflow) - Methodologies
- [Runbooks](#runbooks) - Step-by-step procedures

---

## Getting Started

**For first-time users**

| Document | Location | Purpose | Time |
|----------|----------|---------|------|
| **START_HERE.md** | `/` | Welcome page with learning paths | 5 min |
| **README.md** | `/` | Project overview and features | 5 min |
| **Tutorial** | `docs/beginner/TUTORIAL.md` | Hands-on introduction | 15 min |
| **Glossary** | `docs/reference/GLOSSARY.md` | 150+ term reference | Ongoing |

**Start here:** [START_HERE.md](../START_HERE.md) → [Tutorial](beginner/TUTORIAL.md) → [Protocol](../PROTOCOL.md)

**Progressive Learning Path:** See `docs/beginner/` for structured learning:
1. [00-quick-start.md](beginner/00-quick-start.md) - Get SDP working in 5 min
2. [01-first-feature.md](beginner/01-first-feature.md) - Build your first feature
3. [02-common-tasks.md](beginner/02-common-tasks.md) - Everyday SDP patterns
4. [03-troubleshooting.md](beginner/03-troubleshooting.md) - Solve problems

---

## Core Documentation

**Primary references for SDP**

### Protocol Specifications

| Document | Location | Description |
|----------|----------|-------------|
| **PROTOCOL.md** | `/` | Complete SDP specification (English) |
| **PROTOCOL_RU.md** | `/` | Полная спецификация SDP (Русский) |
| **CODE_PATTERNS.md** | `/` | Implementation patterns and anti-patterns |
| **PRINCIPLES.md** | `docs/PRINCIPLES.md` | SOLID, DRY, KISS, YAGNI, TDD |
| **CLAUDE.md** | `/` | Claude Code integration guide |
| **MODELS.md** | `/` | AI model recommendations |

### Product Documentation

| Document | Location | Description |
|----------|----------|-------------|
| **PRODUCT_VISION.md** | `/` | Project manifesto and long-term goals |
| **PROJECT_MAP.md** | `docs/PROJECT_MAP.md` | Project structure and organization |
| **PROJECT_INDEX.md** | `docs/workstreams/INDEX.md` | Workstream reference index |

---

## Guides by Role

**Documentation targeted to specific roles**

### Team Leads & Managers

| Document | Location | Purpose |
|----------|----------|---------|
| **Overview for Leads** | `docs/overview-for-leads.md` | Executive summary of SDP |
| **Tutorial** | `docs/TUTORIAL.md` | Understand workflow |
| **Project Vision** | `PRODUCT_VISION.md` | Strategic alignment |
| **Project Map** | `docs/PROJECT_MAP.md` | Organization structure |

**Path:** Overview → Tutorial → Protocol (as needed)

---

### Engineers & Developers

| Document | Location | Purpose |
|----------|----------|---------|
| **Tutorial** | `docs/TUTORIAL.md` | Learn by doing |
| **Glossary** | `docs/GLOSSARY.md` | Terminology reference |
| **Code Patterns** | `CODE_PATTERNS.md` | Implementation patterns |
| **Principles** | `docs/PRINCIPLES.md` | Quality standards |
| **Protocol** | `PROTOCOL.md` | Complete reference |

**Path:** Tutorial → Glossary → Code Patterns → Principles → Protocol

---

### DevOps & SRE

| Document | Location | Purpose |
|----------|----------|---------|
| **GitHub Integration** | `docs/github-integration/` | CI/CD setup |
| **GitHub Setup** | `docs/github-integration/SETUP.md` | Configuration guide |
| **GitHub Usage** | `docs/github-integration/USAGE.md` | Daily operations |
| **Deployment** | `PROTOCOL.md#deployment` | Deployment workflow |

**Path:** GitHub Setup → GitHub Usage → Runbooks

---

### QA & Test Engineers

| Document | Location | Purpose |
|----------|----------|---------|
| **Quality Gates** | `PROTOCOL.md#quality-gates` | Quality standards |
| **Verification Protocol** | `docs/verification-protocol.md` | Testing procedures |
| **Test Runbook** | `docs/runbooks/test-runbook.md` | Testing workflow |
| **Two-Stage Review** | `docs/two-stage-review.md` | Review process |

**Path:** Quality Gates → Verification Protocol → Test Runbook

---

### AI Tool Users

| Document | Location | Purpose |
|----------|----------|---------|
| **Claude Code Guide** | `docs/guides/CLAUDE_CODE.md` | Claude Code integration |
| **Cursor Guide** | `docs/guides/CURSOR.md` | Cursor integration |
| **Cursor Advanced** | `docs/guides/CURSOR_ADVANCED.md` | Advanced Cursor features |
| **Multi-IDE Parity** | `docs/multi-ide-parity.md` | Tool comparisons |
| **Claude Integration** | `CLAUDE.md` | Claude Code skills |

**Path:** Choose your IDE → Tool-specific guide → Multi-IDE parity

---

## Guides by Topic

**Deep dives into specific topics**

### Architecture & Design

| Document | Location | Description |
|----------|----------|-------------|
| **Clean Architecture** | `docs/concepts/clean-architecture/README.md` | Layered architecture guide |
| **Artifacts** | `docs/concepts/artifacts/README.md` | Artifact management |
| **Concepts Index** | `docs/concepts/README.md` | Architecture concepts |

**Related:** [CODE_PATTERNS.md](../CODE_PATTERNS.md), [PRINCIPLES.md](PRINCIPLES.md)

---

### Multi-Agent System

| Document | Location | Description |
|----------|----------|-------------|
| **Agent Roles** | `.claude/agents/README.md` | Agent type reference |
| **Agent Coordination** | `PROTOCOL.md#unified-workflow` | Coordination architecture |
| **Orchestrator** | `PROTOCOL.md#autonomous-execution` | Orchestration workflow |

**See also:** Unified Orchestrator, Message Router, Beads Integration

---

### Task Tracking (Beads)

| Document | Location | Description |
|----------|----------|-------------|
| **Beads Integration** | `docs/beads-integration/` | Integration overview |
| **Beads CLI Usage** | `PROTOCOL.md#beads-integration` | CLI reference |
| **Beads Models** | `src/sdp/beads/models.py` | Python API |

---

### Quality & Validation

| Document | Location | Description |
|----------|----------|-------------|
| **Quality Gate Schema** | `docs/quality-gate-schema.md` | Quality configuration |
| **Verification Protocol** | `docs/verification-protocol.md` | Testing procedures |
| **Completion Protocol** | `docs/completion-protocol.md` | Definition of done |
| **Two-Stage Review** | `docs/two-stage-review.md` | Review process |

---

### GitHub Integration

| Document | Location | Description |
|----------|----------|-------------|
| **Setup** | `docs/github-integration/SETUP.md` | Initial configuration |
| **Usage** | `docs/github-integration/USAGE.md` | Daily operations |
| **Troubleshooting** | `docs/github-integration/TROUBLESHOOTING.md` | Common issues |
| **E2E Validation** | `docs/github-integration/E2E_VALIDATION.md` | Testing setup |
| **README** | `docs/github-integration/README.md` | Overview |

---

### Workflow & Processes

| Document | Location | Description |
|----------|----------|-------------|
| **SDP Update Workflow** | `docs/workflows/sdp-update-workflow.md` | Updating SDP itself |
| **Git Workflow** | `docs/guides/GIT_WORKFLOW.md` | Git best practices |
| **Migration Guide** | `docs/migration/ws-naming-migration.md` | Migrating identifiers |

---

## API & Integration

**Technical integration documentation**

### CLI Commands

| Command | Skill | Location | Description |
|---------|-------|----------|-------------|
| `@feature` | Feature | `.claude/skills/feature/SKILL.md` | Unified feature development |
| `@idea` | Idea | `.claude/skills/idea/SKILL.md` | Requirements gathering |
| `@design` | Design | `.claude/skills/design/SKILL.md` | Workstream planning |
| `@build` | Build | `.claude/skills/build/SKILL.md` | Execute workstream |
| `@review` | Review | `.claude/skills/review/SKILL.md` | Quality review |
| `@deploy` | Deploy | `.claude/skills/deploy/SKILL.md` | Production deployment |
| `@oneshot` | Oneshot | `.claude/skills/oneshot/SKILL.md` | Autonomous execution |
| `/debug` | Debug | `.claude/skills/debug/SKILL.md` | Systematic debugging |
| `@issue` | Issue | `.claude/skills/issue/SKILL.md` | Bug routing |
| `@hotfix` | Hotfix | `.claude/skills/hotfix/SKILL.md` | Emergency fix |
| `@bugfix` | Bugfix | `.claude/skills/bugfix/SKILL.md` | Quality fix |
| `/tdd` | TDD | `.claude/skills/tdd/SKILL.md` | TDD enforcement (internal) |

**See also:** [CLAUDE.md](../CLAUDE.md#available-skills)

---

### Python API

| Module | Location | Description |
|--------|----------|-------------|
| **Beads Client** | `src/sdp/beads/` | Task tracking API |
| **Schema Validation** | `src/sdp/schema/` | Intent validation |
| **TDD Runner** | `src/sdp/tdd/` | TDD cycle execution |
| **Feature Management** | `src/sdp/feature/` | Product vision API |
| **Design Graph** | `src/sdp/design/` | Dependency management |
| **Agent System** | `src/sdp/unified/` | Multi-agent coordination |

---

## Reference

**Lookup materials and references**

### Glossaries & References

| Document | Location | Description |
|----------|----------|-------------|
| **Glossary** | `docs/GLOSSARY.md` | 150+ SDP terms defined |
| **Model Mapping** | `docs/model-mapping.md` | AI model recommendations |
| **Quality Gate Schema** | `docs/quality-gate-schema.md` | Quality configuration format |

---

### Templates & Schemas

| Template | Location | Purpose |
|----------|----------|---------|
| **Workstream Template** | `docs/workstreams/TEMPLATE.md` | Workstream file format |
| **Intent Schema** | `docs/schema/intent-schema.json` | Intent validation |
| **Quality Gate Schema** | `docs/quality-gate-schema.md` | Quality configuration |

---

### ADRs (Architecture Decision Records)

| ADR | Location | Topic |
|-----|----------|-------|
| **0001** | `docs/adr/0001-file-native-consensus-protocol.md` | File-native consensus |
| **0002** | `docs/adr/0002-reliability-first-file-native-protocol.md` | Reliability improvements |
| **0003** | `docs/adr/0003-progressive-universal-protocol.md` | Progressive protocol |
| **0004** | `docs/adr/0004-unified-progressive-consensus.md` | Unified progressive |
| **README** | `docs/adr/README.md` | ADR index |

---

## Architecture

**Design and architecture documentation**

### Design Documents

| Document | Location | Description |
|----------|----------|-------------|
| **SDP Analysis** | `docs/plans/2025-01-26-sdp-analysis-design.md` | System analysis |
| **Developer Dashboard** | `docs/plans/2025-01-26-developer-dashboard-design.md` | Dashboard design |
| **AI-Human Comm** | `docs/plans/2025-01-26-ai-human-comm-design.md` | Communication design |
| **Hybrid Workflow** | `docs/plans/2025-01-28-hybrid-sdp-unified-workflow-design.md` | Workflow integration |
| **Beads Integration** | `docs/plans/2025-01-28-beads-sdp-integration-design.md` | Beads integration |
| **SDP Improvement** | `docs/plans/2026-01-29-sdp-improvement-design.md` | Improvement roadmap |

---

## Process & Workflow

**Methodologies and processes**

### Development Processes

| Document | Location | Description |
|----------|----------|-------------|
| **Feature Development** | `PROTOCOL.md#feature-development-flow` | End-to-end feature workflow |
| **TDD Cycle** | `CODE_PATTERNS.md#decomposing-large-tasks` | Red → Green → Refactor |
| **Code Review** | `docs/two-stage-review.md` | Two-stage review process |
| **Debugging** | `docs/runbooks/debug-runbook.md` | Systematic debugging |

---

### Quality Processes

| Document | Location | Description |
|----------|----------|-------------|
| **Quality Gates** | `PROTOCOL.md#quality-gates` | Mandatory checks |
| **Verification** | `docs/verification-protocol.md` | Verification procedures |
| **Completion** | `docs/completion-protocol.md` | Definition of done |

---

## Runbooks

**Step-by-step procedures**

| Runbook | Location | Purpose |
|---------|----------|---------|
| **Debug** | `docs/runbooks/debug-runbook.md` | Debug failing tests |
| **Test** | `docs/runbooks/test-runbook.md` | Run test suite |
| **Oneshot** | `docs/runbooks/oneshot-runbook.md` | Autonomous execution |
| **Git Hooks** | `docs/runbooks/git-hooks-installation.md` | Install git hooks |

**Common workflow:** Debug → Test → Verify → Complete

---

## Workstreams

**Workstream documentation and templates**

### Workstream Index

| Document | Location | Description |
|----------|----------|-------------|
| **Workstream Index** | `docs/workstreams/INDEX.md` | Complete workstream list |
| **Template** | `docs/workstreams/TEMPLATE.md` | Workstream file format |
| **Backlog** | `docs/workstreams/backlog/` | Planned workstreams |
| **In Progress** | `docs/workstreams/in_progress/` | Active workstreams |
| **Completed** | `docs/workstreams/completed/` | Finished workstreams |

---

### Example Workstreams

| Workstream | Location | Status |
|------------|----------|--------|
| **Two-Stage Review** | `docs/workstreams/completed/00-003-01-two-stage-review.md` | ✅ Completed |
| **Systematic Debugging** | `docs/workstreams/completed/00-003-02-systematic-debugging.md` | ✅ Completed |
| **Core Parser** | `docs/workstreams/completed/00-006-01-core-workstream-parser.md` | ✅ Completed |
| **Git Hooks** | `docs/workstreams/completed/00-007-02-git-hooks.md` | ✅ Completed |
| **@idea/@design** | `docs/workstreams/completed/00-007-11-add-idea-design.md` | ✅ Completed |

---

## Drafts & Plans

**In-progress and planning documents**

### Drafts

| Document | Location | Topic |
|----------|----------|-------|
| **Idea F013** | `docs/drafts/idea-f013-ai-comm.md` | AI communication |
| **Idea F014** | `docs/drafts/idea-f014-workflow-efficiency.md` | Workflow improvements |
| **Beads F015** | `docs/drafts/beads-f015-ai-comm.md` | Beads integration |
| **GitHub Orchestrator** | `docs/drafts/idea-github-agent-orchestrator.md` | GitHub agents |

---

### Plans

| Plan | Location | Description |
|------|----------|-------------|
| **F012 Implementation** | `docs/plans/2025-01-26-F012-implementation-plan.md` | Dashboard implementation |
| **SDP Improvement** | `docs/plans/2026-01-29-sdp-improvement-design.md` | Phase 1 improvements |

---

## Directory Structure

```
sdp/
├── START_HERE.md              # 👈 Start here!
├── README.md                  # Project overview
├── PROTOCOL.md                # Full specification
├── PROTOCOL_RU.md             # Спецификация (Русский)
├── CODE_PATTERNS.md           # Implementation patterns
├── CLAUDE.md                  # Claude Code integration
├── MODELS.md                  # AI model recommendations
├── PRODUCT_VISION.md          # Project manifesto
│
├── docs/                      # All documentation
│   ├── GLOSSARY.md            # 150+ term reference
│   ├── SITEMAP.md             # 👈 This file!
│   ├── TUTORIAL.md            # 15-minute tutorial
│   ├── PRINCIPLES.md          # Engineering principles
│   │
│   ├── adr/                   # Architecture Decision Records
│   │   ├── 0001-*.md
│   │   ├── 0002-*.md
│   │   └── README.md
│   │
│   ├── concepts/              # Architecture concepts
│   │   ├── clean-architecture/
│   │   ├── artifacts/
│   │   └── README.md
│   │
│   ├── guides/                # Role and tool guides
│   │   ├── CLAUDE_CODE.md
│   │   ├── CURSOR.md
│   │   ├── CURSOR_ADVANCED.md
│   │   └── GIT_WORKFLOW.md
│   │
│   ├── github-integration/    # GitHub integration docs
│   │   ├── SETUP.md
│   │   ├── USAGE.md
│   │   ├── TROUBLESHOOTING.md
│   │   ├── E2E_VALIDATION.md
│   │   └── README.md
│   │
│   ├── runbooks/              # Step-by-step procedures
│   │   ├── debug-runbook.md
│   │   ├── test-runbook.md
│   │   ├── oneshot-runbook.md
│   │   └── git-hooks-installation.md
│   │
│   ├── workstreams/           # Workstream documentation
│   │   ├── INDEX.md
│   │   ├── TEMPLATE.md
│   │   ├── backlog/
│   │   ├── in_progress/
│   │   └── completed/
│   │
│   ├── drafts/                # In-progress documents
│   ├── plans/                 # Design documents
│   ├── schema/                # JSON schemas
│   ├── intent/                # Machine-readable intents
│   └── beads-integration/     # Beads documentation
│
├── .claude/                   # Claude Code configuration
│   ├── skills/                # Skill definitions (@ and / commands)
│   │   ├── feature/
│   │   ├── idea/
│   │   ├── design/
│   │   ├── build/
│   │   ├── review/
│   │   ├── deploy/
│   │   ├── oneshot/
│   │   ├── debug/
│   │   ├── issue/
│   │   ├── hotfix/
│   │   ├── bugfix/
│   │   └── tdd/
│   │
│   ├── agents/                # Multi-agent configurations
│   │   ├── planner.md
│   │   ├── builder.md
│   │   ├── reviewer.md
│   │   ├── deployer.md
│   │   ├── orchestrator.md
│   │   └── README.md
│   │
│   └── settings.json          # Claude Code settings
│
├── src/sdp/                   # Python source code
│   ├── beads/                 # Beads integration
│   ├── core/                  # Core functionality
│   ├── schema/                # Schema validation
│   ├── tdd/                   # TDD runner
│   ├── feature/               # Feature management
│   ├── design/                # Design graph
│   └── unified/               # Multi-agent system
│
├── prompts/                   # Command prompts
│   └── commands/              # Skill instructions
│
├── hooks/                     # Git hooks
│   ├── pre-build.sh
│   └── post-build.sh
│
└── tests/                     # Test suite
    ├── unit/
    ├── integration/
    └── e2e/
```

---

## Document Statistics

| Category | Count |
|----------|-------|
| **Core Documents** | 7 |
| **Guides** | 12 |
| **Runbooks** | 4 |
| **ADRs** | 5 |
| **Plans** | 7 |
| **Drafts** | 12 |
| **Workstreams** | 50+ |
| **Skills** | 13 |
| **Agents** | 5 |
| **Total** | 115+ files |

---

## Quick Find

**Looking for something specific?**

### By Topic

- **Getting started** → [START_HERE.md](../START_HERE.md), [Tutorial](TUTORIAL.md)
- **Terminology** → [Glossary](GLOSSARY.md)
- **How SDP works** → [Protocol](../PROTOCOL.md)
- **Code patterns** → [Code Patterns](../CODE_PATTERNS.md)
- **Quality standards** → [Principles](PRINCIPLES.md), [Quality Gates](../PROTOCOL.md#quality-gates)
- **AI-IDE integration** → [Guides](guides/), [CLAUDE.md](../CLAUDE.md)
- **Task tracking** → [Beads Integration](beads-integration/), [Protocol](../PROTOCOL.md#beads-integration)
- **Multi-agent system** → [Agent Roles](../.claude/agents/README.md), [Protocol](../PROTOCOL.md#unified-workflow)
- **GitHub setup** → [GitHub Integration](github-integration/)
- **Debugging** → [Debug Runbook](runbooks/debug-runbook.md), [/debug](../PROTOCOL.md#debugging)
- **Deployment** → [Protocol](../PROTOCOL.md#deployment), [GitHub Usage](github-integration/USAGE.md)

### By Document Type

- **Specifications** → [Protocol](../PROTOCOL.md), [Quality Gate Schema](quality-gate-schema.md)
- **Tutorials** → [Tutorial](TUTORIAL.md), [Guides](guides/)
- **References** → [Glossary](GLOSSARY.md), [Model Mapping](model-mapping.md)
- **Procedures** → [Runbooks](runbooks/)
- **Decisions** → [ADRs](adr/)
- **Templates** → [Workstream Template](workstreams/TEMPLATE.md)

### By Role

- **Team Lead** → [Overview for Leads](overview-for-leads.md), [Tutorial](TUTORIAL.md)
- **Engineer** → [Tutorial](TUTORIAL.md), [Code Patterns](../CODE_PATTERNS.md), [Principles](PRINCIPLES.md)
- **DevOps** → [GitHub Integration](github-integration/), [Runbooks](runbooks/)
- **QA** → [Verification Protocol](verification-protocol.md), [Test Runbook](runbooks/test-runbook.md)

---

## Contributing to Documentation

**Adding or updating docs:**

1. **Check existing structure** - Use this sitemap to find the right location
2. **Follow naming conventions** - Use `kebab-case.md` for files
3. **Update this sitemap** - Add new documents to appropriate section
4. **Cross-reference** - Link to related documents
5. **Maintain consistency** - Follow existing formatting and style

**Documentation locations:**
- Guides → `docs/guides/`
- Runbooks → `docs/runbooks/`
- Reference → `docs/` (root level)
- Plans → `docs/plans/`
- Drafts → `docs/drafts/`

---

## Maintenance

**Last full audit:** 2026-01-29
**Next scheduled audit:** 2026-02-28

**Keep this sitemap current:**
- Add new documents when created
- Remove or mark deprecated documents
- Update document counts and statistics
- Verify all links resolve correctly

---

**Version:** 1.0
**Last Updated:** 2026-01-29
**Maintained by:** SDP Documentation Team

---

**See Also:**
- [START_HERE.md](../START_HERE.md) - New user guide
- [GLOSSARY.md](GLOSSARY.md) - Terminology reference
- [PROTOCOL.md](../PROTOCOL.md) - Complete specification
