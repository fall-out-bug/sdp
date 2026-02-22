# SDP: Structured Development Protocol

## The Question

You have AI agents writing code. Maybe one agent, maybe a swarm. They produce PRs. Some PRs are good. Some are terrible. Most are somewhere in between.

Now ask: **"Show me proof that this PR was planned, tested, reviewed by a separate agent, and stayed within its declared scope."**

No tool in the OpenCode ecosystem answers this. The ecosystem has orchestration (Vibe Kanban, Swarm Tools), policy enforcement (Cupcake), session management (micode, oh-my-opencode), parallel execution — dozens of solutions for running agents. Zero solutions for *proving what agents did*.

That gap is what SDP is about.

## What Exists Today

SDP is a protocol and evidence layer for AI agent workflows. One developer, daily use, 18 releases, actively forming.

### The Protocol

A specification — prompts, JSON schemas, shell hooks — that structures agent work into phases:

```
Intent → Plan → Execute → Verify → Review → Publish
```

Each phase has a contract. Skip a phase and the state machine blocks the next one. The protocol is language-agnostic: it works with OpenCode, Claude Code, Cursor, or anything that can read a markdown file.

**Status:** Implemented. Used daily. Works with OpenCode and Claude Code.

### The Evidence Envelope

Every agent run produces a JSON document with 9 required sections:

| Section | What It Proves |
|---------|---------------|
| **intent** | What the agent was asked to do |
| **plan** | How it decided to approach the task |
| **execution** | What it actually did (files changed, commands run) |
| **verification** | That tests pass and the change works |
| **review** | That a separate agent reviewed the work |
| **risk_notes** | What could go wrong |
| **boundary** | Declared scope vs. observed scope vs. compliance |
| **provenance** | SHA-256 hash chain linking this run to previous runs |
| **trace** | Timing data, phase transitions, gap detection |

This is not logging. Logging says "agent ran at 14:32." An evidence envelope says "the agent declared intent to fix bug #42, planned to modify auth.go and auth_test.go, executed changes to those files, ran tests that passed with 94% coverage, was reviewed by a reviewer agent who approved with 2 suggestions that were addressed, and the boundary check confirms no files outside declared scope were modified." Machine-readable. Validatable. Hashable.

**Status:** Implemented. JSON Schema defined. Validation works. Hash chain provenance works.

### The PR Gate

A CLI command that blocks PR merge unless evidence is complete and valid:

```bash
sdp-evidence validate --evidence .sdp/evidence/run-123.json
```

One command in CI. If the evidence envelope is incomplete, invalid, or has a broken hash chain — merge is blocked.

**Status:** Implemented. Used in our own CI. Not yet published as a standalone release.

### The K8s Bridge

An adapter controller that connects kubeopencode (the K8s operator for running agents) with the evidence layer. When kubeopencode finishes an agent task, the adapter projects the result into a strict evidence envelope.

```
Issue → kubeopencode Task CRD → agent runs → adapter projects evidence → PR (if valid)
```

**Status:** In development. The reconcile loop works. Needs hardening. Depends on kubeopencode upstream.

## Where SDP Fits in the Ecosystem

SDP is not a platform. It's a layer. It composes with the tools you already use.

| You Need | Use | SDP Adds |
|----------|-----|----------|
| Orchestration | Vibe Kanban, Swarm Tools | Evidence envelope for each orchestrated task |
| Policy | Cupcake | Evidence that policy was checked (not just enforced) |
| K8s agent execution | kubeopencode | Evidence projection from Task CRD status |
| Issue tracking | Beads | Evidence-gated state machine transitions |
| CI/CD | GitHub Actions, any CI | `sdp-evidence validate` as a PR gate |

We don't rebuild what the ecosystem already does well. We add the one thing it doesn't have: proof.

## What's Coming

| What | When | How |
|------|------|-----|
| Evidence JSON Schema published | Soon | In this repo (`schema/`) |
| `sdp-evidence` CLI as standalone binary | Soon | `go install` or binary release |
| awesome-opencode listing | After first release | Protocol + evidence CLI |
| kubeopencode upstream contributions | In progress | PRs to kubeopencode |
| OpenCode plugin for local evidence collection | Planned | Evidence during local agent runs |

## What We're Exploring

Some ideas in the lab. No promises, no timelines. Just interesting problems:

- **Multi-persona evaluation** — 5 expert agents (architect, SRE, security, DX, product) reviewing code adversarially, with dissent tracking
- **Self-improvement loops** — agents detecting their own failure patterns and adjusting
- **Cross-project federation** — one evidence layer serving multiple repos
- **Telemetry-driven backlog** — LLM analysis of run patterns that auto-generates improvement issues

These live in [`sdp_lab`](https://github.com/fall-out-bug/sdp_lab) — our research playground. It's private for now, but we're happy to invite people who want to play with these ideas together. Just open an issue and say hi.

## The Honest Numbers

| | |
|---|---|
| **Developers** | 1 |
| **Daily use** | Yes, on our own work |
| **Releases** | 18 (v0.9.6) |
| **Stars** | 16 |
| **Evidence tooling** | Working, not yet published separately |
| **Core code** | ~14K LOC Go |
| **Tests** | ~10K LOC |

SDP is forming. The evidence concept is real and working. The protocol is real and working. Everything else is actively taking shape.

If you care about knowing what your agents actually did — not just that they ran, but that they followed a process and left proof — that's what we're building.

---

*"AI agents can implement features, but without evidence it's just vibes."*
