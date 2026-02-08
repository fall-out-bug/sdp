# SDP: The Accountability Layer for AI-Generated Code

> *"Когда тебя спросят 'что это за хуйня' — ты не можешь ответить 'что-то навайбкодилось, хз как и когда'."*

---

## The Accountability Gap

AI writes code. Somebody ships it. Something breaks. A courtroom, an auditor, or an incident review asks:

**"What happened? When? Who decided? What was verified?"**

Today, the answer is: "I... looked at it. It seemed fine."

That answer will get increasingly expensive. Not because of any one party's liability — the legal frameworks are still being written. But because **someone will need the evidence.** The developer. The model provider. The platform vendor. The regulator. Someone.

**SDP doesn't decide who's responsible. SDP records what happened.**

It's a neutral evidence layer. It captures: which model generated the code, what spec it was built against, what verification ran, what the output was, who approved it. The chain of facts. Not opinions about liability.

> **SDP replaces "it seemed fine" with a record of what actually happened.**

---

## What SDP Is

**An open protocol for accountable AI code generation.**

Three guarantees:

1. **Every piece of AI-generated code has provenance.** Which model. Which version. What parameters. What spec. When. Who initiated.

2. **Every verification is recorded with actual output.** Not "tests passed" — the actual `pytest` output with line counts. Not "types checked" — the actual `mypy` report.

3. **Every decision is traceable.** When something breaks, you trace from the git commit back through the chain: generation → verification → approval. Not "someone approved this PR." **Here's what was generated, here's what was checked, here's who signed off.**

> This is the **forensic system for AI-generated code** — a neutral record of what happened, for whoever needs it: the developer, the auditor, or the incident responder.

---

## The Philosophy

### Decomposition Is Permanent

Not because models are weak. Because **humans can't verify large blobs.**

Code review effectiveness drops to near-random above 400 lines of diff (Microsoft, Google studies). Even if GPT-6 generates flawless 5000-line files — YOU can't verify a 5000-line diff. The bottleneck is human cognition, not AI capability. This doesn't change.

Decomposition compensates for two permanent problems:
- **Specification sparsity** — "Add OAuth" is 15 words. The implementation is 5000 decisions. Smaller units = fewer wrong guesses.
- **Verification capacity** — You can review 200 lines carefully. You can only skim 2000.

### The Verification Stack

Ordered by ROI, not by impressiveness:

| Layer | What It Catches | Cost |
|-------|----------------|------|
| **Decomposition** (<200 LOC units) | Context overload, spaghetti | Free |
| **Type checking** (mypy/tsc) | Type errors — #1 AI bug class | Free |
| **Static analysis** (semgrep) | Security patterns, anti-patterns | Free |
| **Model provenance** | "Which model wrote this and when?" | Free |
| **Property-based testing** | Edge cases, invariants | 1.2x tokens |
| **Cross-model review** | Correlated blind spots | 2.5x tokens |
| **Human review** (the PR) | Everything above missed | Human time |

The first four layers are free and form the **base protocol**. Everything above is optional and additive.

---

## The UX: plan / apply / incident

Borrowed from Terraform — because it works.

### `sdp plan`

Show what you're about to do. Don't do it yet.

```bash
$ sdp plan "Add OAuth2 login with Google and GitHub"

  3 units planned:

  1. Backend OAuth service
     Scope: src/auth/oauth.ts, src/auth/providers/*.ts
     Gates: types, tests, semgrep, provenance
     Risk: high (auth module)

  2. Frontend login component
     Scope: src/components/Login.tsx, src/hooks/useOAuth.ts
     Gates: types, tests, provenance
     Risk: medium

  3. Integration tests
     Scope: tests/integration/oauth.test.ts
     Gates: tests
     Risk: low

  Estimated: ~3 min | 3 units | ~$0.15

  [apply] [edit] [cancel]
```

### `sdp apply`

Execute the plan. Each unit: generate → verify → record.

```
Applying...
  [1/3] Backend OAuth ████████████████ done
        types ✓ | tests ✓ (91% cov) | semgrep ✓ | provenance recorded
  [2/3] Frontend login ████████████████ done
        types ✓ | tests ✓ (87% cov) | provenance recorded
  [3/3] Integration    ████████████████ done
        tests ✓ | provenance recorded

All gates passed. PR created: github.com/org/repo/pull/42
Evidence chain: .sdp/evidence/2026-02-08-oauth-f01.json
```

### `sdp incident`

Production broke. Trace back.

```bash
$ sdp incident abc1234  # git commit hash

  Commit abc1234 — AI-generated via SDP
  ├── Model: claude-sonnet-4-20250514, temp=0.3
  ├── Spec: "Backend OAuth service with Google + GitHub"
  ├── Acceptance criteria: 4/4 passed at build time
  ├── Verification: types ✓, tests ✓ (91%), semgrep ✓
  ├── Approved by: @developer at 2026-02-08T14:32:00Z
  └── Evidence: .sdp/evidence/2026-02-08-oauth-f01.json

  Related units in same feature:
  ├── unit-2: Frontend login (commit def5678)
  └── unit-3: Integration tests (commit ghi9012)
```

When someone asks "what happened?" — you have the answer. Every detail. Timestamped. Machine-readable.

### Modes

```bash
sdp plan "Add auth"                   # Show plan, wait for approval
sdp plan "Add auth" --auto-apply      # Plan + apply immediately (ship mode)
sdp plan "Add auth" --interactive     # Stop at every fork, you decide (drive mode)
sdp apply                             # Apply last plan
sdp apply --retry 3                   # Retry failed unit 3 only
sdp incident <commit>                 # Forensic trace
```

**`--auto-apply` is "ship".** Trust the framework, get a PR.
**`--interactive` is "drive".** You decide at every fork: "Sessions or JWT?" → your call.
**Default is `plan`.** Show the intent. Build trust. Like `terraform plan`.

---

## The Strategy: Protocol First

### Why Protocol, Not CLI

The old plan: build a CLI tool, then maybe extract a protocol.
The new plan: **protocol first, then tools on top.**

Why:
1. **What already exists is a protocol.** SDP has 19 skills, a synthesis engine, dependency graph, checkpoint system — all running as a Claude Code plugin. That's a protocol implementation, not a CLI.
2. **Open protocol gets embedded everywhere.** If Cursor, Replit, Windsurf, OpenCode can all speak SDP protocol — the standard wins regardless of which tool wins.
3. **Enterprise wants control, not tools.** Banks don't adopt your CLI. They adopt your spec and implement it with their infra. The protocol is what they buy. The CLI is what indie devs use.

### The Stack

```
┌─────────────────────────────────────────┐
│  Tools (what users touch)               │
│  Claude Code plugin | CLI | IDE plugins │
├─────────────────────────────────────────┤
│  Protocol (the standard)                │
│  plan → apply → verify → evidence       │
│  Model provenance | Audit chain         │
├─────────────────────────────────────────┤
│  Engine (the implementation)            │
│  Decomposition | Verification | Forensics│
└─────────────────────────────────────────┘
```

**Protocol** = open, embeddable by anyone.
**Engine** = reference implementation (open-source verification, proprietary orchestration).
**Tools** = surfaces for different users (plugin for devs, CLI for automation, dashboards for enterprise).

---

## Who This Is For

**Code that costs more to fix than to write.**

| Segment | What They Need |
|---------|---------------|
| **Fintech** | Forensic proof that payment code was verified before merge |
| **Healthcare** | HIPAA-compliant audit trail for AI-generated code provenance |
| **Infrastructure** | `terraform plan` for AI code — show before apply |
| **Enterprise SaaS** | Policy: "all AI PRs must have SDP evidence chain" |
| **Regulated industries** | SOC2/DORA/ISO 27001 compliance for AI code |

Enterprise interest validated: top-3 bank (contracted), major airline (contracted), largest marketplace (evaluating).

**Not for:** Landing pages, MVPs, prototypes. If rewriting is cheaper than verifying — just ship.

---

## The Moat

Code is reproducible. Any team can build a verification CLI.

What they can't build:

1. **Decomposition heuristics.** Learned from thousands of verified builds. "OAuth → 3 units, not 5." Gets better with every run.

2. **AI failure taxonomy.** A dataset of what AI actually gets wrong, organized by model, language, domain. This doesn't exist anywhere.

3. **The evidence standard.** If SDP becomes how enterprises prove AI code was verified — that's a standard, not a product. Standards have network effects.

---

## The Five Bets

1. **AI will write most code by 2028.** Not controversial.
2. **"What happened?" will be mandatory.** Courtrooms, audits, incident reviews will demand records of AI code generation.
3. **Provenance is permanent.** Regardless of whether code is generated in 50 small units or one big prompt — the record of what model, what spec, what verification, who approved never stops being valuable.
4. **Forensics > Verification.** Everyone claims "verified AI code." Nobody has a forensic chain from production incident back to generation.
5. **Protocol > Product.** The team that sets the evidence standard wins.

Decomposition is a powerful technique that improves generation quality today. It may or may not be permanent. **Provenance and evidence are permanent regardless.**

---

*SDP Manifesto v3.0 — February 2026*
*Protocol-first. Accountability-first. Forensics-first.*
