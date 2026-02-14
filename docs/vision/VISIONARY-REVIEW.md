# Visionary Review of SDP v2.1 Manifesto

> **Date:** 2026-02-08
> **Reviewed by:** Andrej Karpathy, Amjad Masad, Pieter Levels, Simon Willison
> **Status:** APPROVED with corrections. Ship the wedge.

---

## One-Line Verdicts

| Expert | Verdict |
|--------|---------|
| **Karpathy** | "Decomposition compensates for specification sparsity, not model weakness. That's a stronger argument — even perfect models benefit." |
| **Masad** | "Adopt with developers first, compliance second. Free GitHub Action → devs use it → managers notice → compliance asks for audit trails." |
| **Levels** | "Ship it yesterday. The manifesto is 2x too long. The tool is the manifesto." |
| **Willison** | "Ship a semgrep ruleset for AI-generated code in week 1. The Action follows. The CLI follows. The data collection starts immediately." |

---

## Key Corrections

### 1. Reframe WHY Decomposition Works (Karpathy)

**Before:** "Models lose coherence over long outputs — property of context windows and attention."

**After:** "Decomposition compensates for **specification sparsity**. When you say 'Add OAuth2 login,' that has hundreds of implicit decisions. The model makes them all simultaneously. Some will be wrong. Decomposition reduces the decision surface per step."

**Why this matters:** If framed as "model weakness," someone ships GPT-6 with long-output training and the thesis evaporates. If framed as "humans can't specify 2000 lines of behavior in one prompt," that's permanently true. Stronger argument.

### 2. Plan Default Should Be Context-Dependent (Karpathy)

Not just plan-by-default or auto-approve-by-default. Smart default:
- Project has `.sdp/config` with `risk: high`? → Show plan
- Changed files touch auth/payments/infra? → Show plan
- Fresh project, no production traffic? → Auto-approve
- Let the tool infer the right mode

But if forced to pick one: **plan-by-default is correct for SDP's positioning.** Selling to fintech/healthcare = trust brand = show the plan.

### 3. Sell to Developers FIRST, Compliance SECOND (Masad)

**Before:** "Target: compliance budgets, not dev tools budgets."

**After:** The manifesto has this backwards.

**The correct sequence:**
1. Free GitHub Action → developers adopt
2. Developers love it → engineering managers notice
3. Managers ask "can we get audit trails?" → compliance sale

**Compliance buyers don't buy from indie projects.** They need SOC2 certs, enterprise support, legal entities to sue. That's Phase 3, not Phase 1.

### 4. The Decomposition Is the Weakest Link (Karpathy)

**The hole:** `sdp ship` decomposes → verifies each unit. But who verifies the decomposition?

What if the AI puts token validation in the frontend instead of giving it a backend unit? The human skims the plan and hits Y.

**Fix needed:** Heuristics that flag suspicious decompositions ("this unit touches auth AND database — consider splitting"). The data flywheel should feed into decomposition quality, not just gate calibration.

### 5. Ship Sequence Is Wrong (Willison + Levels)

**Manifesto says:** Phase 1 = GitHub Action + `sdp ship` (both in 4 weeks).

**Correct sequence:**
- **Week 1:** Semgrep ruleset for AI-generated code (zero install: `semgrep --config p/sdp-ai-python`)
- **Week 2:** Wrap in GitHub Action with semgrep gate
- **Week 3:** Add type checking gate
- **Week 4:** Add test verification gate
- **Month 2:** `sdp ship` with plan-by-default

Action first. CLI second. The Action is the wedge. The CLI is the product.

### 6. The "Data as Moat" Needs Honest Bootstrap (Willison)

**The dataset structure:**
- Input prompt + model used + generated code
- Decomposition applied (units, boundaries)
- Gate results per unit
- Human review outcome
- **Post-merge defects** (track reverts/hotfixes within 48h of SDP-verified merge)

**Honest admission needed:** Phase 1 data comes from own usage + 10 design partners. Not from Action user base. The chicken-and-egg is real.

### 7. Add "Solo Devs Shipping Payments/Auth" to Target (Levels)

**Before:** Target is enterprise (fintech, healthcare, infra).

**After:** Add a row: "Solo devs with payment/auth code." They're the viral path. Enterprises don't tweet about CI pipelines. Indie hackers tweet about tools that saved their Stripe integration.

### 8. SDK Conversation in Month 2, Not Month 8 (Masad)

**Before:** Phase 3 (months 4-8) = Agent SDK.

**After:** Start the SDK conversation with IDE vendors (Cursor, Windsurf, Continue) in month 2 — with early data from the GitHub Action. The CLI is the demo. The SDK is the business.

### 9. Consider Snapshot Testing (Willison)

Add between static analysis and property-based testing:
- **Snapshot/regression testing**: "Did the output change unexpectedly?"
- Cheap, catches subtle AI behavior changes
- Especially valuable when AI regenerates existing code

### 10. The Manifesto Is Still 2x Too Long (Levels)

The `sdp ship` UX section + the GitHub Action wedge = the manifesto.
Everything else is appendix material. The "Five Bets," the Phase 4 protocol extraction, the moat section — move to appendix.

**"A manifesto that nobody reads because the tool is so good they don't need to — that's the goal."**

---

## Revised Ship Sequence

| Week | Ship | Purpose |
|------|------|---------|
| **1** | Semgrep ruleset for AI Python/TS code | Zero friction, starts data collection |
| **2** | GitHub Action (semgrep gate) | The wedge |
| **3** | Add type checking + test gates to Action | Full verification stack |
| **4** | Publish first dataset (own builds) | Evidence begins |
| **5-8** | `sdp ship` with plan-by-default | The product |
| **8+** | Start SDK conversations with IDE vendors | The business |

---

## Would They Use It?

| Expert | Would Use | For What |
|--------|-----------|----------|
| **Karpathy** | GitHub Action: yes. `sdp ship`: for boring/commodity code only. Not for creative work. | "Safety net for vibe coding" |
| **Masad** | Would embed SDK in Replit Agent, not use CLI directly. | "The decomposition logic — if it's visibly smarter than what I'd do manually" |
| **Levels** | GitHub Action: no (PHP + no tests). `sdp ship` for Stripe/auth code: maybe. | "Use SDP for the 10% of code that would ruin you if wrong" |
| **Willison** | Semgrep ruleset: yes, immediately. Action: yes. `sdp ship`: after it's proven. | "Ship the semgrep ruleset and I'll write a blog post about it" |

---

## The Consensus

**All four agree:**
1. The thesis is sound. Decomposition as the core insight is right.
2. The manifesto is dramatically better than v1.
3. **Stop writing manifestos and ship the semgrep ruleset.**
4. The window is real (9-18 months depending on who moves first).
5. Data > architecture > marketing narratives.

**The one sentence that matters:**

> "A manifesto that nobody reads because the tool is so good they don't need to — that's the goal." — Levels + Willison
