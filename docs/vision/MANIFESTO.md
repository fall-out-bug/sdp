# SDP v2: Verified Build for AI-Generated Code

> *"The point isn't that AI can write code. It's that you can't tell when it wrote it wrong."*

---

## The Problem

AI writes code. That's solved.

The unsolved problem: **can you trust it?**

A landing page that's slightly wrong? Rebuild it in 5 minutes. A payment system that's slightly wrong? Chargebacks, lawsuits, regulatory fines.

Models get better at generating. They don't get better at guaranteeing. The gap between "it looks right" and "it IS right" — that's where the damage lives.

**SDP exists to close that gap.**

Not with another model reviewing the first. Not with a protocol spec nobody reads. With something simpler and harder to replicate:

> **Decompose into small units. Verify each one. Record everything.**

That's it. That's the entire thesis.

---

## Why Decomposition Is the Insight

Every AI model generates better code in a 200-line focused module than in a 2000-line sprawling file. This isn't an opinion — it's a property of context windows, attention mechanisms, and how autoregressive models lose coherence over long outputs.

SDP's core move: **break features into atomic units before generating code.** Each unit has:
- A clear goal and acceptance criteria
- A bounded scope (files to touch)
- Dependencies on other units
- Verification gates to pass

This is where trust comes from. Not from a second model reviewing the first (that helps ~15%, at 2.5x cost). From the **structure itself** — small, typed, tested units that are individually verifiable.

### The Verification Stack

Ordered by ROI, not by impressiveness:

| Layer | What It Catches | Cost | Status |
|-------|----------------|------|--------|
| **Decomposition** (<200 LOC units) | Context overload, spaghetti generation | Free | Core feature |
| **Type checking** (mypy --strict, tsc) | Type errors — the #1 class of AI bugs | Free | Core feature |
| **Static analysis** (semgrep + custom rules) | Security patterns, anti-patterns | Free | Core feature |
| **Property-based testing** (Hypothesis, fast-check) | Edge cases, invariant violations | 1.2x tokens | High priority |
| **Cross-model review** (different provider) | "Different perspective" bugs | 2.5x tokens | High-risk code only |
| **Human review** (the PR) | Everything above missed | Human time | The output |

The first three layers are free and catch 35-50% more bugs than no verification. That's the foundation. Cross-model review is a premium layer for auth, payments, and data deletion — not the default for every workstream.

---

## What SDP Is

**A build system for AI-generated code with invisible structure and a built-in audit trail.**

One opinion: **AI code should be generated in small verified units, not in large unverified blobs.**

Two modes. One philosophy. Different depths of engagement.

### Mode 1: `ship` — Trust with Verification

For people who want to hand off to AI but need verified results.

```bash
sdp ship "Add OAuth2 login with Google and GitHub"
```

What happens inside: decomposition, specs, TDD, static analysis, audit trail.
What the user sees: a progress bar and a PR with a green checkmark.

**The structure is invisible. The artifacts exist for three consumers:**
1. The next AI session (inter-session memory — context, not docs)
2. The human investigator (forensics when something breaks)
3. The compliance system (SOC2/HIPAA audit trail)

### Mode 2: `mentor` — Guide the AI

For people who want their expertise reflected in the output. The full circle:

```bash
sdp idea "Add OAuth2 login"    # Interactive requirements (you shape the spec)
sdp design idea-oauth           # Architecture decisions (you choose the approach)
sdp build 00-001-01             # TDD execution (you watch each unit)
sdp review F01                  # Multi-agent review (you approve the verdict)
```

Same decomposition. Same verification. Same artifacts. But the human is in the loop at every stage — mentoring the AI, not just consuming its output.

### The Insight

Both modes produce the same thing: **small verified units with an audit trail.** The difference is who makes the decisions — the AI (ship) or the human (mentor). The philosophy doesn't change. The depth of engagement does.

> *"Ship" users trust the framework. "Mentor" users trust themselves. Both get verified code.*

---

## What SDP Is NOT

- **Not a coding methodology.** SDP doesn't enforce TDD, clean architecture, or file size limits as doctrine. These are verification heuristics that the framework tunes based on results.
- **Not a project manager.** No sprints, no velocity, no roadmaps. SDP verifies output.
- **Not for disposable code.** If your code costs less to rewrite than to verify, don't use SDP.
- **Not a straitjacket.** Risk-proportional verification: `payments/` gets the full audit trail, `components/` gets code + tests. The framework infers the risk, the user doesn't configure it.

---

## The User Experience

### Ship Mode (default): One Command

```bash
$ sdp ship "Add OAuth2 login with Google and GitHub"

Planning... (3s)

  Workstreams:
  1. Backend OAuth service (passport.js, Google + GitHub providers)
     Scope: src/auth/oauth.ts, src/auth/providers/*.ts
     Gates: types, tests, semgrep

  2. Frontend login component (React, existing AuthContext)
     Scope: src/components/Login.tsx, src/hooks/useOAuth.ts
     Gates: types, tests

  3. Integration tests (OAuth flow e2e)
     Scope: tests/integration/oauth.test.ts
     Gates: tests

  Estimated: ~3 min | 3 units | ~$0.15

Proceed? [Y/n/edit]
```

**Plan by default.** Like `terraform plan`. You see the decomposition before anything executes. You can edit it — "skip unit 3, I'll write those myself" or "use next-auth, not passport."

Then execution:

```
Building...
  [1/3] Backend OAuth ████████████████ done (91% coverage, types ✓, semgrep ✓)
  [2/3] Frontend login ████████████████ done (87% coverage, types ✓)
  [3/3] Integration    ████████████████ done

Verification passed. PR created: github.com/org/repo/pull/42
```

**Streaming progress is non-negotiable.** The progress bar IS the product. Nobody stares at a spinner for 3 minutes.

### Mentor Mode: Full Circle

For engineers who want to shape the output — the same pipeline, but human-in-the-loop at each stage:

```bash
sdp idea "Add OAuth2 login"       # AI interviews you: "Sessions or JWT? PKCE?"
sdp design idea-oauth              # AI proposes architecture, you choose
sdp build 00-001-01               # You watch TDD: red → green → refactor
sdp build 00-001-02               # Unit by unit, your expertise guides the AI
sdp review F01                    # Multi-agent review, you approve the verdict
```

**Same decomposition. Same verification. Same artifacts.** The difference: your decisions are encoded into the specs. The AI doesn't guess "sessions or JWT" — it asks you.

**When to use mentor mode:**
- New domain (you need to teach the AI your constraints)
- High-stakes system (you want sign-off at every stage)
- Learning SDP (you want to see the internals)
- Compliance requires human approval at each gate

### Flags for Control

```bash
sdp ship "Add auth"                  # Plan → approve → build (default)
sdp ship "Add auth" --auto-approve   # Skip plan approval (I trust the decomposition)
sdp ship "Add auth" --cross-review   # Add cross-model review (premium, for high-risk code)
sdp ship "Add auth" --edit           # Open plan in editor before building
sdp ship "Add auth" --dry-run        # Show plan only, don't build
```

### Per-Unit Rollback

When unit 3 fails but units 1-2 succeeded, regenerate just unit 3:

```bash
sdp ship --retry 3    # Regenerate only the failed unit
```

The checkpoint infrastructure makes this free. No wasted work.

---

## The Wedge: GitHub Action

Before `sdp ship` exists, the wedge that gets SDP into repos:

```yaml
# .github/workflows/sdp-verify.yml
name: SDP Verify
on: [pull_request]
jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: sdp-dev/verify-action@v1
        with:
          gates: [types, semgrep, tests]
```

One YAML file. One CI check. Every AI-generated PR gets verified. Teams add it once and forget about it.

**This is the door.** `sdp ship` is the house. You enter through the door.

---

## Who This Is For

### Code That Costs More to Fix Than to Write

| Segment | Why They Need SDP | Pain Today |
|---------|-------------------|-----------|
| **Fintech** | Payment bugs = chargebacks + fines | AI-generated code ships without structural verification |
| **Healthcare** | Data bugs = HIPAA violations | No audit trail for AI code provenance |
| **Infrastructure** | Config bugs = outages | Terraform has `plan`; AI coding tools don't |
| **Enterprise SaaS** | Multi-tenant bugs = data leaks | AI doesn't verify isolation boundaries |
| **ML/Data** | Pipeline bugs = silent data corruption | No property-based testing in AI workflows |

### Early Traction

Enterprise interest is not hypothetical:

- **A top-3 bank** — contracted for AI development workflow integration
- **A major airline** — contracted for AI-assisted development tooling
- **The largest online marketplace** — evaluating for engineering teams

All via independent contracts. Signal: **enterprises with compliance requirements already want this.** The pain is real. The budget exists. The question is execution speed.

### Who Pays

Not developers. Developers adopt free tools.

**Compliance budgets.** The pitch: "Auditable proof that all AI-generated code was independently verified." That's a SOC2 line item.

**Engineering lead budgets.** The pitch: "SDP-verified code has N% fewer defects. Here's the data."

**The enterprise path is validated.** The next step is packaging what works into a repeatable product — not proving demand exists.

---

## The Moat

Code is reproducible. Cursor can build `sdp ship` in a sprint.

What Cursor can't build in a sprint:

1. **Decomposition heuristics.** How to break "Add OAuth2 login" into the right 3 units (not 5, not 1). Learned from thousands of verified builds. Each build refines the heuristics.

2. **Verification gate calibration.** Which semgrep rules catch AI-specific anti-patterns. Which property-based test strategies expose AI-specific edge cases. Built from data on what AI actually gets wrong.

3. **Trust dataset.** "SDP-verified code has X% fewer defects than unverified AI code." The dataset that proves this is the moat. Without it, SDP is a claim. With it, SDP is evidence.

**The moat deepens with usage, not engineering.** Every verified build adds data. Every caught bug refines a heuristic. Every false positive tunes a gate. This is a flywheel, not a feature list.

---

## The Roadmap

### Phase 1: The Wedge (Weeks 1-4)

**Ship the GitHub Action.**
- `sdp-dev/verify-action@v1` — types + static analysis + test verification
- Free for open source, paid for private repos
- Target: 100 repos in month 1

**Ship `sdp ship` with plan-by-default.**
- Decomposition → approval → verified execution → PR
- Streaming progress, per-unit rollback
- Static analysis pipeline (mypy + semgrep) as default gates

**Start collecting data.**
- 10 features with SDP vs 10 without
- Measure: bugs found in review, time to correct implementation, token spend

### Phase 2: The Evidence (Months 2-4)

**Publish the dataset.**
- "SDP-verified code: X% fewer defects, Y% faster to merge, Z% less rework"
- Open dataset, reproducible methodology
- Blog posts, case studies

**Add property-based testing.**
- AI-generated Hypothesis/fast-check strategies
- Second-highest ROI verification layer

**Add compliance reporting.**
- SOC2-ready audit trails
- Per-PR verification certificates
- This is where the money is

### Phase 3: The Platform (Months 4-8)

**Agent SDK.**
- `sdp.Decompose()`, `sdp.Verify()`, `sdp.Audit()`
- Pitch to platforms WITH the evidence dataset: "Our data shows X% defect reduction. Want to embed this?"

**Cross-model review for premium tier.**
- High-risk code only (auth, payments, data)
- Additional cost, additional confidence
- Justified by data, not by thesis

### Phase 4: The Standard (If We Win)

**Extract the protocol — but only if the framework wins.**
- If 1000+ repos use SDP, platforms will want interoperability
- The protocol is what DHH called "Rack" — the thinnest interface for composability
- It emerges from success, not from planning

---

## The Five Bets

1. **AI will write most code by 2028.** Not controversial.
2. **Trust will be the bottleneck.** Controversial today. Obvious after the first major AI-code incident.
3. **Decomposition is a permanent advantage.** Small typed tested units > large unverified blobs, regardless of model capability.
4. **Static analysis beats cross-model review on ROI.** Free tools that catch 35-50% of bugs beat expensive tools that catch 15% more.
5. **Data is the moat.** The team with the best dataset on "what AI gets wrong and how to catch it" wins. Not the team with the best architecture.

If bet 3 is wrong — if future models produce flawless 5000-line generations — SDP is unnecessary.

We bet they won't.

---

## The 12-Month Window

Cursor and GitHub will ship "verified AI code" features. Probably within 18 months. They have distribution, model access, and enterprise customers demanding it.

SDP's advantage: **speed and focus.** Ship the wedge now. Build the dataset now. Become the tool that Cursor embeds — or the tool that enterprises require alongside Cursor.

The clock is running.

---

*SDP v2.2 — February 2026*
*Updated after unified philosophy debate (PG, DHH, Levels, Karpathy, Masad, Willison).*
*Two-mode model (ship/mentor) based on enterprise feedback from T-Bank, S7, Avito.*
