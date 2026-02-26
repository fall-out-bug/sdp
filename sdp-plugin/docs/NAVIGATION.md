# SDP Documentation Navigation

> **🎯 Single Entry Point:** Start here to find what you need.
> **Last Updated:** 2026-02-06

## 🚀 Where to Start? (By User Role)

### New to SDP?
1. **[Quick Start](#quick-start)** - Get started in 5 minutes
2. **[Tutorial](docs/TUTORIAL.md)** - 15-minute hands-on tutorial
3. **[PROTOCOL.md](docs/PROTOCOL.md)** - Full specification (when you need details)

### Experienced User
1. **[Skill Reference](#skill-reference)** - Quick command lookup
2. **[Decision Trees](#decision-trees)** - Choose the right workflow
3. **[Quality Gates](docs/quality-gates.md)** - Code quality standards

### Enterprise/Team
1. **[Beads Workflow](docs/workflow-decision.md)** - Task tracking integration
2. **[SRE SLOs](docs/slos/orchestrator.md)** - Service level objectives
3. **[Security](docs/SECURITY.md)** - Security guidelines

---

## Quick Start

### 3 Essential Commands

```bash
# 1. Plan a feature
@feature "Add user authentication"

# 2. Execute workstream
@build 00-001-01

# 3. Quality check
@review F01
```

That's it! Everything else builds on these three commands.

---

## Table of Contents

### Level 1: Getting Started (L1)
- [Quick Start](#quick-start) - 3 commands to know
- [Tutorial](docs/TUTORIAL.md) - Learn by doing
- [README](README.md) - Project overview

### Level 2: Core Concepts (L2)
- [Workstreams](docs/PROTOCOL.md#workstream) - Atomic units of work
- [Features](docs/PROTOCOL.md#feature) - Collections of workstreams
- [Quality Gates](docs/quality-gates.md) - Code quality standards
- [Skills](#skill-reference) - Available commands

### Level 3: Workflows & Decisions (L3)
- [Decision Trees](#decision-trees) - Choose the right approach
- [Workflow Comparison](docs/workflow-decision.md) - Beads vs Traditional
- [Debugging](#debugging) - Systematic problem solving

### Level 4: Advanced Topics (L4)
- [SRE & Operations](#sre--operations) - Monitoring, SLOs
- [Architecture](#architecture) - System design
- [Extensions](#extensions) - Customization

---

## Decision Trees

### Tree 1: Feature Development Workflow

```
START: I want to build a feature
│
├─ Is it a NEW feature with unknown requirements?
│  └─ YES → @feature "Description"
│     (Progressive disclosure: vision → requirements → planning → execution)
│
├─ Do you have clear requirements but need to plan workstreams?
│  └─ YES → @design idea-name
│     (Interactive planning with EnterPlanMode)
│
├─ Do you have workstreams ready to execute?
│  └─ YES → @build 00-001-01
│     (TDD discipline with progress tracking)
│
└─ Do you want autonomous execution (no human intervention)?
   └─ YES → @oneshot F01
      (Spawns orchestrator agent, runs workstreams autonomously)
```

**See Also:** [Workflow Comparison](docs/workflow-decision.md)

---

### Tree 2: Quality & Review

```
START: I need to check code quality
│
├─ Is it a completed feature?
│  └─ YES → @review F01
│     (Multi-agent review: QA + Security + DevOps + SRE + TechLead + Docs)
│
├─ Did you find a bug in production?
│  └─ YES → @issue "Bug description"
│     (Analyzes severity, routes to @hotfix or @bugfix)
│
├─ Is code failing tests unexpectedly?
│  └─ YES → /debug "Test failure description"
│     (Systematic debugging: scientific method)
│
└─ Do you need to validate quality gates manually?
   └─ YES → Check docs/quality-gates.md
      (Coverage ≥80%, files <200 LOC, type hints, etc.)
```

**See Also:** [Quality Gates Reference](docs/quality-gates.md)

---

### Tree 3: Execution Strategy

```
START: How should I execute workstreams?
│
├─ Need real-time progress visibility?
│  └─ YES → @build 00-001-01
│     (TodoWrite tracking, Red-Green-Refactor TDD cycle)
│
├─ Want hands-off execution?
│  └─ YES → @oneshot F01
│     (Autonomous orchestrator, checkpoint/resume support)
│
├─ Need to coordinate multiple features?
│  └─ YES → @feature "Description"
│     (Manages multi-feature orchestration)
│
└─ Just testing/prototyping?
   └─ YES → Use @build or @oneshot (same, but no Beads tracking)
```

**See Also:** [Beads Integration](docs/workflow-decision.md#beads-first-workflow)

---

## Skill Reference

### Feature Development
| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `@feature` | Unified feature workflow (recommended) | New features, progressive disclosure |
| `@idea` | Interactive requirements gathering | Exploring user requirements |
| `@design` | Interactive workstream planning | Planning from requirements draft |
| `@oneshot` | Autonomous execution | Hands-off multi-workstream execution |

### Execution
| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `@build` | Execute workstream | Manual execution with progress tracking |
| `/tdd` | TDD cycle enforcement | Automatic (used by @build) |

### Quality & Debugging
| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `@review` | Multi-agent quality review | Completed features |
| `/debug` | Systematic debugging | Unexpected failures |
| `@issue` | Bug analysis and routing | Production bugs |
| `@hotfix` | Emergency fix (P0) | Critical production issues |
| `@bugfix` | Quality fix (P1/P2) | Non-emergency bugs |

### Deployment
| Skill | Purpose | When to Use |
|-------|---------|-------------|
| `@deploy` | Production deployment | Ready-to-ship features |

### Internal Skills (not called directly)
| Skill | Purpose | Called By |
|-------|---------|----------|
| `/tdd` | TDD discipline | @build (automatic) |
| `guard` | Pre-edit validation | @build (automatic) |

---

## Documentation Levels (Progressive Disclosure)

### Level 1: Essential (L1)
**Goal:** Get started immediately
- [Quick Start](#quick-start) - 3 commands
- [Tutorial](docs/TUTORIAL.md) - 15-minute intro
- [README](README.md) - Project overview

### Level 2: Core Concepts (L2)
**Goal:** Understand how SDP works
- [PROTOCOL.md](docs/PROTOCOL.md) - Full specification
- [Workstreams](docs/PROTOCOL.md#workstream) - Atomic tasks
- [Quality Gates](docs/quality-gates.md) - Code quality standards

### Level 3: Workflows (L3)
**Goal:** Choose the right approach
- [Decision Trees](#decision-trees) - This page
- [Workflow Comparison](docs/workflow-decision.md) - Beads vs Traditional
- [CLAUDE.md](CLAUDE.md) - Integration guide

### Level 4: Advanced (L4)
**Goal:** Deep dive and customization
- [SRE & Operations](#sre--operations) - Monitoring, SLOs
- [Architecture](#architecture) - System design
- [Security](docs/SECURITY.md) - Security guidelines

---

## SRE & Operations

### Monitoring & Observability
- [Orchestrator SLOs](docs/slos/orchestrator.md) - Service level objectives
- [Structured Logging](docs/workflows/backlog/sdp-zig-structured_logging.md) - Logging implementation
- [Telemetry Guide](docs/TELEMETRY_HOWTO.md) - Telemetry setup

### Runbooks
- [Decision Logging](docs/runbooks/decision-logging.md) - Logging runbook
- [Decision Backup](docs/operations/decision-backup.md) - Backup procedures

---

## Architecture

### Core Components
- [Checkpoint System](docs/PROTOCOL.md#checkpoint-system) - State persistence
- [Orchestrator](docs/F024_unified_workflow_spec.md) - Multi-workstream execution
- [Beads Integration](docs/workflow-decision.md#beads-integration) - Task tracking

### Feature Specifications
- [Unified Workflow](docs/F024_unified_workflow_spec.md) - Multi-workstream orchestration
- [Workstream Specifications](docs/workflows/) - Detailed workstream docs

---

## Extensions

### Role Management
- [Role Setup Guide](docs/ROLE_SETUP_GUIDE.md) - Configure 100+ roles
- [Role Switching Guide](docs/ROLE_SWITCHING_GUIDE.md) - Dynamic role activation

### Notifications
- [Notification System Guide](docs/NOTIFICATION_SYSTEM_GUIDE.md) - Telegram integration
- [Notification Provider](docs/workflows/WS-017_notification_provider.md) - Provider interface

### Testing
- [Testing Guide](docs/TESTING_GUIDE.md) - Testing best practices
- [Bug Report Guide](docs/BUG_REPORT_GUIDE.md) - Bug reporting workflow

---

## Language-Specific Quick Starts

- [Python Quick Start](docs/examples/python/QUICKSTART.md)
- [Java Quick Start](docs/examples/java/QUICKSTART.md)
- [Go Quick Start](docs/examples/go/QUICKSTART.md)

---

## Debugging

### Systematic Debugging Workflow
Use `/debug` for unexpected failures:

```bash
/debug "Test fails unexpectedly"
```

**Debugging Process:**
1. **Observe** - Gather evidence (logs, error messages)
2. **Hypothesize** - Formulate hypotheses
3. **Test** - Verify hypotheses experimentally
4. **Fix** - Apply root cause fix
5. **Verify** - Confirm fix works

**See Also:** [CLAUDE.md - Debugging Section](CLAUDE.md#debugging)

---

## FAQ

### Q: When should I use @build vs @oneshot?
**A:** Use `@build` for real-time progress visibility and hands-on control. Use `@oneshot` for autonomous execution without intervention. See [Decision Tree 3](#tree-3-execution-strategy).

### Q: What's the difference between @feature and @design?
**A:** `@feature` is the **unified workflow** (vision → requirements → planning → execution). `@design` is just the **planning phase** - use it when you already have requirements but need to break down into workstreams.

### Q: Do I need Beads CLI?
**A:** Recommended for team collaboration and multi-session work. Optional for single-developer projects. See [Workflow Comparison](docs/workflow-decision.md).

### Q: How do I get help?
**A:** Use `@help` or `claude help` to see available commands and get interactive guidance.

---

## Need Help?

- **Quick Help:** Run `@help` in Claude Code
- **Report Issues:** See [Bug Report Guide](docs/BUG_REPORT_GUIDE.md)
- **Contributing:** See [README - Contributing Section](README.md#contributing)

---

**This page is maintained as part of the SDP project. Last updated: 2026-02-06**
