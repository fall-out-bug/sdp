# SDP Product Contract

> **Version:** 1.0.0  
> **Status:** Stable  
> **Last Updated:** 2026-04-18

## Overview

This document is the **single source of truth** for SDP's product definition. It defines the two user paths, the stage model, control surfaces, and harness support policy.

All other documentation (QUICKSTART.md, PROTOCOL.md) references this contract. If you find contradictions, this document takes precedence.

---

## User Paths

### Path 1: Local Mode (Default)

**Target:** Individual developers working locally.

**Entry Point:** `@feature` skill.

**Flow:**

```
┌─────────────────────────────────────────────────────────────────────┐
│                    LOCAL MODE — DEFAULT PATH                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. BOOTSTRAP                                                       │
│     ├─ Install SDP: curl install.sh | sh                            │
│     ├─ Run: sdp init --auto                                         │
│     └─ Result: .sdp/config.yml + project structure                 │
│                                                                     │
│  2. INTAKE                                                          │
│     ├─ Run: @feature "Add user authentication"                      │
│     ├─ Agent asks: technical approach, UI/UX, testing, security     │
│     └─ Result: docs/intent/sdp-XXX.json + docs/drafts/beads-sdp-XXX.md │
│                                                                     │
│  3. SHAPING                                                         │
│     ├─ Run: @design beads-sdp-XXX                                   │
│     ├─ Agent explores codebase, creates workstreams                │
│     └─ Result: docs/workstreams/beads-sdp-XXX.md (5-30 leaf WS)    │
│                                                                     │
│  4. EXECUTION                                                       │
│     ├─ Option A (Autonomous): @oneshot <feature-id>                │
│     │  ├─ Orchestrator executes all ready leaf WS                   │
│     │  ├─ Saves checkpoints after each WS                           │
│     │  └─ Resumes from interruption                                 │
│     │                                                                │
│     └─ Option B (Manual): @build 00-XXX-01 (repeat per WS)         │
│        ├─ Execute single leaf workstream with TDD                  │
│        └─ Commit when complete                                      │
│                                                                     │
│  5. FINDINGS                                                        │
│     ├─ Run: @review <feature-id>                                    │
│     ├─ Multi-agent quality review (6 agents)                        │
│     └─ Result: APPROVED / CHANGES_REQUESTED                         │
│                                                                     │
│  6. DELIVERY                                                        │
│     ├─ Run: @deploy <feature-id>                                    │
│     ├─ Generate: PR, changelog, git tag                             │
│     └─ Merge to main                                                │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Control Surfaces:**
- **Primary:** Claude Code CLI (skills: @feature, @build, @review, @deploy)
- **Companion:** Beads CLI (bd ready, bd create, bd close)

**Data Storage:** Local git repo + `.sdp/` directory.

---

### Path 2: Operator Mode (Advanced)

**Target:** Platform teams running SDP in CI/CD.

**Entry Point:** `sdp-evidence` CLI binary.

**Flow:**

```
┌─────────────────────────────────────────────────────────────────────┐
│                  OPERATOR MODE — ADVANCED PATH                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  1. BOOTSTRAP                                                       │
│     ├─ Install evidence CLI: go install github.com/fall-out-bug/sdp/sdp-evidence/cmd/sdp-evidence@latest │
│     ├─ Configure CI: .github/workflows/sdp-evidence.yml             │
│     └─ Result: Evidence gate in PR pipeline                         │
│                                                                     │
│  2. INTAKE                                                          │
│     ├─ Same as Local Mode (@feature → @design)                      │
│     └─ Or: Import workstreams from external tools                   │
│                                                                     │
│  3. SHAPING                                                         │
│     ├─ Same as Local Mode                                           │
│     └─ Or: Use Strataudit for corpus analysis                       │
│                                                                     │
│  4. EXECUTION                                                       │
│     ├─ Agents run in CI (OpenCode, Claude Code, etc.)               │
│     ├─ Evidence envelope emitted: .sdp/evidence/<run-id>.json       │
│     └─ Hash-chain provenance enforced                               │
│                                                                     │
│  5. FINDINGS                                                        │
│     ├─ Run: sdp-evidence validate .sdp/evidence/<run-id>.json       │
│     ├─ Validate: completeness, schema, hash-chain                  │
│     └─ Result: VALID / INVALID + details                            │
│                                                                     │
│  6. DELIVERY                                                        │
│     ├─ PR gate: Evidence validation must pass                       │
│     ├─ Run: sdp-evidence gate check                                 │
│     └─ Block merge if evidence missing or invalid                   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Control Surfaces:**
- **Primary:** `sdp-evidence` CLI (validate, gate check)
- **Companion:** CI/CD system (GitHub Actions, GitLab CI, etc.)

**Data Storage:** Evidence envelopes in `.sdp/evidence/` + git for provenance.

---

## Stage Model

SDP follows a six-stage model. Both paths use the same stages, but differ in control surfaces.

| Stage | Purpose | Local Mode Trigger | Operator Mode Trigger |
|-------|---------|-------------------|----------------------|
| **Bootstrap** | Initialize project | `sdp init --auto` | `sdp-evidence init` |
| **Intake** | Gather requirements | `@feature` skill | Same (local dev) |
| **Shaping** | Plan workstreams | `@design` skill | Same (local dev) |
| **Execution** | Implement work | `@oneshot` or `@build` | CI agents + evidence emit |
| **Findings** | Quality validation | `@review` skill | `sdp-evidence validate` |
| **Delivery** | Ship to production | `@deploy` skill | PR gate + merge |

**Key Principle:** Stages are sequential but can loop back (e.g., Findings → Execution if review fails).

---

## Control Surfaces

### Primary Control Surface

**Claude Code CLI** (recommended for Local Mode):

```bash
@feature "Add X"      # Intake + Shaping
@oneshot <feature-id> # Execution
@review <feature-id>  # Findings
@deploy <feature-id>  # Delivery
```

### Companion Control Surface

**Beads CLI** (task tracking):

```bash
bd ready              # Find ready tasks
bd create --title="X" # Create task
bd close <id>         # Close task
```

### Operator Control Surface

**sdp-evidence CLI** (Operator Mode):

```bash
sdp-evidence validate .sdp/evidence/<run-id>.json
sdp-evidence gate check
```

### Board Visibility

**Evidence Dashboard** (planned):

- View evidence envelopes across all features
- Filter by stage, agent, harness
- Trace hash-chain provenance

### Quickstart Commands (planned)

**CLI-based first experience:**

```bash
sdp assess [project-path]    # Read-only project assessment
sdp try "task description"   # Try a task on temporary branch
sdp adopt                    # Adopt successful trial into SDP
```

These commands provide a lightweight entry point before full SDP setup. Currently in planning phase.

---

## Harness Support Policy

SDP is designed to work across multiple AI harnesses. Support levels:

### Recommended

| Harness | Version | Notes |
|---------|---------|-------|
| **Claude Code** | Latest | Primary target. Full skill support. |
| **OpenCode** | Latest | Full skill support via `.opencode/` adapters. |

### Supported

| Harness | Version | Notes |
|---------|---------|-------|
| **Cursor** | Latest | Skills via `.cursor/` adapters. |
| **Windsurf** | Latest | Skills via `.windsurf/` adapters. |

**"Supported" means:** Skills load and execute, but some advanced features may be limited.

### Compatible

| Harness | Version | Notes |
|---------|---------|-------|
| **Copilot** | Latest | Protocol-compatible (use prompts/schemas manually). |
| **Zed** | Latest | Protocol-compatible (use prompts/schemas manually). |

**"Compatible" means:** You can use SDP prompts/schemas, but no skill auto-loading.

---

## Default Path Diagram

This diagram is reusable across all SDP documentation:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   BOOTSTRAP  │────▶│    INTAKE    │────▶│   SHAPING    │────▶│  EXECUTION   │
│              │     │              │     │              │     │              │
│ sdp init     │     │ @feature     │     │ @design      │     │ @oneshot     │
│              │     │              │     │              │     │ @build       │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
                                                                 │
                                                                 ▼
                                                          ┌──────────────┐
                                                          │  FINDINGS    │◀─────┐
                                                          │              │      │
                                                          │ @review      │      │
                                                          │              │      │
                                                          └──────────────┘      │
                                                                 │             │
                                                                 ▼             │
                                                          ┌──────────────┐      │
                                                          │   DELIVERY   │──────┘ (if CHANGES_REQUESTED)
                                                          │              │
                                                          │ @deploy      │
                                                          │              │
                                                          └──────────────┘
```

---

## Contract Validity

This contract is valid for:

- **SDP Protocol:** v0.10.0+
- **SDP CLI:** v0.9.8+
- **Evidence CLI:** v1.0.0+ (when released)

### Versioning Policy

- **Major version bump:** Breaking change to stages, paths, or harness support.
- **Minor version bump:** New skills, new harness support, UX improvements.
- **Patch version bump:** Bug fixes, documentation updates.

---

## Related Documents

| Document | Purpose | Link |
|----------|---------|------|
| **QUICKSTART.md** | Get started in 5 minutes | [Quick Start](QUICKSTART.md) |
| **PROTOCOL.md** | Full protocol specification | [Protocol](PROTOCOL.md) |
| **PRODUCT_VISION.md** | Product vision and positioning | [Vision](../PRODUCT_VISION.md) |
| **CLAUDE.md** | Claude Code integration guide | [Claude Guide](../CLAUDE.md) |

---

## Change Log

| Date | Version | Change |
|------|---------|--------|
| 2026-04-18 | 1.0.0 | Initial product contract (F097) |

---

## Getting Help

- **Documentation:** [docs/](../)
- **Issues:** [GitHub Issues](https://github.com/fall-out-bug/sdp/issues)
- **Community:** [OpenCode Ecosystem](https://github.com/kubeopencode)

---

**End of Contract**
