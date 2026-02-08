# SDP Product Vision

> *AI generates code. Somebody ships it. Somebody asks "what happened?" SDP has the answer.*

---

## The World in 2028

Most production code is AI-generated. Developers describe intent, AI writes implementation. This is already happening — the question is what happens when it goes wrong.

**The accountability gap:** when AI-generated code fails in production, nobody has a record. No trace of what was specified. No proof of what was verified. No chain from the incident back to the generation. Just a git blame pointing at a PR that got rubber-stamped.

Who's liable? The developer who shipped it? The model provider? The tool vendor? That's a legal question — and it'll be contested for a decade.

But regardless of who's liable, **someone needs the evidence.** SDP is that evidence.

---

## What SDP Is

**An open protocol for accountable AI code generation.**

Three layers:

### 1. The Protocol (open standard)

A specification that any tool can implement:

- **plan** — decompose a feature into verified units, show the plan
- **apply** — generate, verify, and record each unit
- **evidence** — cryptographically linked chain: spec → code → verification → approval
- **incident** — trace from any commit back through the full evidence chain

The protocol is tool-agnostic. Claude Code, Cursor, Replit, any IDE, any CI system can speak it.

### 2. The Engine (reference implementation)

Open-source verification. Proprietary orchestration.

- **Decomposition engine** — NL → atomic units with dependency graph
- **Verification stack** — types, static analysis, tests, property-based testing
- **Evidence chain** — model provenance, verification output, approval records
- **Forensic tools** — `sdp incident` traces from production back to spec

### 3. The Tools (user surfaces)

- **Claude Code plugin** — first implementation (exists today as skill system)
- **CLI** — `sdp plan` / `sdp apply` / `sdp incident` for automation and enterprise
- **GitHub Action** — verification in CI/CD
- **IDE plugins** — Cursor, VS Code, JetBrains
- **Enterprise dashboards** — policy, compliance, team analytics

---

## The Strategy: Protocol First

### Why This Order

**Old plan:** CLI → GitHub Action → SDK → Protocol (if we win)
**New plan:** Protocol → Plugin → CLI → SDK → Enterprise

Why reversed:

1. **What exists IS a protocol implementation.** 19 skills, synthesis engine, dependency graph, TDD pipeline — running as Claude Code plugin. Not a CLI.

2. **Open protocol = network effects.** If every AI coding tool embeds SDP protocol, the standard wins regardless of which tool wins. Open source enables this.

3. **Enterprise adopts specs, not tools.** A top-3 bank doesn't install your CLI on 5000 developer machines. They adopt your protocol specification and implement it with their infra team.

4. **Plugin-first = fastest to real users.** Claude Code users can run SDP today. CLI users can't — because `sdp ship` doesn't exist.

### The Open-Core Split

| Layer | License | Why |
|-------|---------|-----|
| **Protocol spec** | Open (CC-BY) | Standard must be open to get adoption |
| **Verification engine** | Open-source (Apache 2.0) | Community builds integrations, ecosystem grows |
| **Orchestration + evidence** | Proprietary | This is what enterprises pay for |
| **Enterprise features** | Commercial | Dashboards, policy, compliance export |

---

## The Philosophy: Four Pillars

### 1. Evidence Is Non-Negotiable

AI generates code. Someone ships it. When it breaks — and it will — someone needs the record. SDP is the neutral ledger: what was generated, how it was verified, who approved it.

Not a verdict about who's responsible. A record of what happened.

### 2. Decomposition Is Permanent

Not because models are weak — because humans can't verify large blobs. Code review effectiveness drops to near-random above 400 LOC. Better models don't fix human cognition.

Unit size is tunable. The principle is not.

### 3. Forensics Over Verification

Everyone claims "verified AI code." Nobody has a forensic chain from production incident back through spec → generation → verification → approval. SDP does.

Verification answers: "did this pass?" Forensics answers: "what happened, why, and who decided?"

### 4. Evidence Is the Product

Not the CLI. Not the plugin. Not the decomposition. **The evidence chain** — the cryptographically linked record of what was specified, generated, verified, and approved. That's what enterprises buy. That's what courtrooms accept. That's what survives model changes, tool changes, and team changes.

---

## The Five Bets

1. **AI writes most code by 2028.** Trajectory, not prediction.
2. **"What happened?" becomes mandatory.** Courtrooms, audits, incident reviews will demand records.
3. **Provenance is permanent.** Regardless of how code is generated — decomposed or monolithic — the record of what model, what spec, who approved never stops being valuable.
4. **Forensics > Verification.** "Did it pass?" is table stakes. "What happened?" is the moat.
5. **Protocol > Product.** The team that sets the evidence standard wins.

---

## Who This Is For

**Code that costs more to fix than to write.**

| Segment | What They Need |
|---------|---------------|
| Fintech | Forensic proof that payment code was verified |
| Healthcare | HIPAA audit trail for AI code provenance |
| Infrastructure | `terraform plan` for AI code |
| Enterprise SaaS | Policy: "all AI PRs need SDP evidence" |
| Regulated industries | SOC2/DORA/ISO 27001 compliance |

Enterprise traction: top-3 bank (contracted), major airline (contracted), largest marketplace (evaluating).

**Not for:** Landing pages, MVPs, disposable code.

---

## The Moat

1. **AI failure taxonomy.** A dataset of what AI gets wrong, by model/language/domain. Doesn't exist anywhere.
2. **Decomposition heuristics.** Learned from thousands of runs. Gets better with every build.
3. **The evidence standard.** If SDP becomes how enterprises prove AI code was verified — that's a standard with network effects. Not a product. A standard.

---

## Kill Criteria

> If after 500 SDP runs, verification catch rate is below 5% and post-merge defect rate is not measurably different from baseline — kill the product.

Specific. Testable. Honest.

---

*SDP Vision v2.0 — February 2026*
*Protocol-first. Accountability-first. Forensics-first.*
