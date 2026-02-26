# SDP Roadmap — Protocol & Evidence Layer

> **For users:** Where SDP is going. Protocol-first, standards-based.  
> **For contributors:** See [sdp_dev roadmap](https://github.com/fall-out-bug/sdp_dev/blob/master/docs/roadmap/ROADMAP.md) for internal phases and implementation details.

---

## What SDP Is

SDP is a **protocol + evidence layer** for AI coding agents. It structures work (Intent → Plan → Execute → Verify → Review → Publish) and produces **proof** of what agents actually did. Works with Claude Code, Cursor, OpenCode, or anything that can read markdown.

| Layer | What |
|-------|------|
| **Protocol** | Prompts, JSON schemas, hooks — language-agnostic |
| **Evidence** | JSON envelope with 9 sections (intent, plan, execution, verification, review, etc.) |
| **Gate** | `sdp-evidence validate` — blocks merge if evidence is incomplete or invalid |

---

## Current State

| Artifact | Status |
|----------|--------|
| Protocol (prompts, skills) | Published, used daily |
| Evidence JSON Schema | Published (`schema/evidence-envelope.schema.json`) |
| in-toto predicate spec | Published (`https://sdp.dev/attestation/coding-workflow/v1`) |
| `sdp` CLI | Published (init, build, verify, guard, etc.) |
| `sdp-evidence` CLI | Working, not yet standalone release |
| K8s bridge | In development |

---

## Roadmap (User-Facing)

### Near Term

| What | When |
|------|------|
| `sdp-evidence` as standalone binary | Soon — `go install` or release |
| awesome-opencode listing | After first release |
| Multi-repo protocol | Protocol defines path→scope, commit per-repo; repo names = project config |

### Standards Migration

| What | Why |
|------|-----|
| in-toto attestation format | Replace custom envelope with DSSE + Sigstore |
| OPA/Rego policies | Replace markdown policies with executable rules |
| CI auto-attestation | Agent cannot "forget" evidence — CI creates it from facts |

### Longer Term

| What |
|------|
| OpenCode plugin for local evidence collection |
| K8s pipeline with enforcement built-in |
| Cross-project federation — one evidence layer, multiple repos |

---

## Multi-Repo Protocol

The protocol is **repo-agnostic**. It defines:

- Scope can span multiple git roots
- Path prefix → repo boundary
- Commit workflow is per-repo (commit in child before parent if submodule)

Concrete repo names (sdp, sdp_dev, opencode, …) live in project config, not in the protocol.

---

## Key References

| Topic | Link |
|-------|------|
| Vision & evidence | [MANIFESTO.md](MANIFESTO.md) |
| Full protocol spec | [PROTOCOL.md](PROTOCOL.md) |
| in-toto predicate | [attestation/coding-workflow-v1.md](attestation/coding-workflow-v1.md) |
| Quick start | [QUICKSTART.md](QUICKSTART.md) |
| Internal roadmap | [sdp_dev ROADMAP](https://github.com/fall-out-bug/sdp_dev/blob/master/docs/roadmap/ROADMAP.md) |

---

*"AI agents can implement features, but without evidence it's just vibes."*
