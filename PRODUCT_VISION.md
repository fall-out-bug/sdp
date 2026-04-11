# SDP — Product Vision

> **Updated:** 2026-04-11
> **Version:** 5.0 (AI PDLC + SDLC)

## Mission

Take a software idea from conception to deployed feature through structured AI-driven phases — with agents, gates, and verifiable evidence at every step.

## What SDP Is

**An AI-managed platform for the full development lifecycle** — Product Development (PDLC) + Software Development (SDLC).

Two first-class phases:
- **Discovery** — research, shaping, scope decision via LLM pipeline → validated spec
- **Delivery** — implementation through agentloop FSM with gate enforcement → deployed code with evidence

- **The Protocol** — prompts, JSON schemas, shell hooks that structure agent work into phases (Intent → Plan → Execute → Verify → Review → Publish). Language-agnostic. Works with OpenCode, Claude Code, Cursor.
- **The Evidence Envelope** — a strict 9-section JSON document (intent, plan, execution, verification, review, risk, boundary, provenance, trace) that every agent run produces. Validated by schema. Hash-chained for tamper detection.
- **The PR Gate** — one CLI command in CI that blocks merge unless evidence is complete and valid.

## What SDP Is Not

- Not an orchestrator (use Vibe Kanban, Swarm Tools)
- Not a policy engine (use Cupcake)
- Not a K8s operator (use kubeopencode — we contribute upstream)
- Not a session manager (use micode, oh-my-opencode)

SDP composes with all of these. It adds evidence to whatever workflow you already have.

## Users

1. **Individual developer** — installs SDP as a submodule, gets structured agent workflow + evidence log
2. **Team using AI agents in CI** — adds `sdp-evidence validate` to PR gates, gets audit trail
3. **Platform team on K8s** — uses kubeopencode + SDP evidence bridge for agent-to-PR pipeline with proof

## Adoption Model

SDP is adoptable without "the whole spaceship":

| Level | What You Get | What You Install |
|-------|-------------|-----------------|
| **Protocol only** | Structured agent workflow via skills | `curl install.sh` — prompts, schemas, hooks |
| **+ Evidence** | Audit log with hash-chain provenance | CLI: `sdp log`, `sdp-evidence validate` |
| **+ K8s Bridge** | Agent runs on kubeopencode produce evidence | Adapter controller (in development) |

Each level is independently valuable. Upgrade path is additive.

## Ecosystem Position

```
┌─────────────────────────────────────────────────┐
│              OpenCode Ecosystem                  │
│                                                  │
│  Orchestration: Vibe Kanban, Swarm Tools         │
│  Policy: Cupcake                                 │
│  K8s: kubeopencode                               │
│  Issues: Beads                                   │
│  Sessions: micode, oh-my-opencode                │
│                                                  │
│  Evidence: SDP  ◄── this is us                   │
│                                                  │
└─────────────────────────────────────────────────┘
```

## What's Shipped

- Protocol spec with 12 skills (v0.9.8, 18 releases)
- Evidence log with hash-chain provenance
- CLI: `sdp doctor`, `sdp guard`, `sdp log`, `sdp status`
- Multi-agent review (6 agents), strategic planning (7 agents), codebase analysis (8 agents)
- Install script with auto-detect for Claude Code, Cursor, OpenCode
- 1,004 commits, 16 stars, MIT license

## What's Next

| Priority | What | Status |
|----------|------|--------|
| **P0** | Publish evidence JSON Schema in `schema/` | Ready |
| **P0** | Release `sdp-evidence` CLI as standalone binary | In progress |
| **P1** | awesome-opencode listing | After first evidence CLI release |
| **P1** | kubeopencode upstream PRs (retry budget, evidence hooks) | In progress |
| **P2** | OpenCode plugin for local evidence collection | Planned |
| **P2** | Adapter controller hardening for K8s bridge | In development |

## Research (sdp_dev)

Exploratory work, no promises:

- Multi-persona adversarial review with dissent tracking
- Agent self-improvement loops (failure pattern → adjustment)
- Cross-project federation via NATS
- Telemetry-driven backlog generation

Private repo [`sdp_lab`](https://github.com/fall-out-bug/sdp_lab) — open an issue to request access.

## Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| Protocol releases | 18 (v0.9.8) | v1.0 |
| Evidence CLI published | No | Yes |
| awesome-opencode listing | No | Yes |
| External users of evidence CLI | 0 | >= 1 |
| kubeopencode upstream PRs | In progress | >= 1 merged |

## Principles

- **Evidence is the moat.** Everything else is ecosystem.
- **Compose, don't replace.** Use Vibe Kanban for orchestration, Cupcake for policy, kubeopencode for K8s.
- **Ship what exists.** Don't announce what isn't built.
- **One developer, honest scope.** Scale the community, not the claims.
- **MIT forever.** Public SDP layers stay open source.
