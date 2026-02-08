# SDP Product Vision

> *AI generates code. You ship it. You own the consequences.*

---

## The World We See

By 2028, AI will write most production code. This is not a prediction — it's a trajectory visible in every metric: GitHub Copilot adoption, Cursor growth, enterprise AI budgets.

The question is not "will AI write code?" The question is: **who is responsible when that code fails?**

Not OpenAI. Not Anthropic. Not the tool vendor. **The person who shipped it.**

Today, a developer using Cursor generates 500 lines of payment logic. The code looks right. The tests pass. It ships. Three months later, a currency conversion bug causes $200K in incorrect charges. Who is accountable?

The developer. Always the developer. The one who said "looks good, ship it."

**SDP exists to give that developer evidence — not hope — that the code is correct.**

---

## The Philosophy: Four Pillars

### 1. Specs Before Code

Not requirements documents. Not JIRA tickets. **Checkpoints.**

Like `terraform plan` — the AI proposes what it's about to do, the human blesses it. Five lines: goal, acceptance criteria, scope. Takes three seconds. Catches "wait, you're putting auth logic on the frontend?" BEFORE 2000 lines are generated.

The spec is not bureaucracy. It's **provenance** — the answer to "why was it built this way?" that you'll need in three months.

### 2. Atomic Units of Work

Every feature decomposes into small verified units. Not because models are weak — because **humans can't verify large blobs.**

The research is clear: code review effectiveness drops to near-random above 400 lines of diff (Microsoft, Google studies). The cognitive bottleneck is human, not AI. This is permanent.

Decomposition also reduces **specification sparsity** — the gap between what you said ("Add OAuth") and what the AI must decide (token storage, refresh rotation, PKCE, error handling). Smaller units = more specific specs = fewer wrong guesses.

The unit size is tunable. The principle is not.

### 3. Invisible Rails and Artifacts

Every step produces an artifact: spec, plan, code, tests, review verdict, audit entry. But the user never touches them directly.

Three consumers of artifacts:
- **The next AI session** — inter-session memory. Without it, every session starts cold.
- **The human investigator** — forensics when something breaks. "What was the spec? What was verified? What was the test output?"
- **The compliance system** — SOC2/HIPAA audit trail. "Prove all AI-generated code was independently verified."

Artifacts should be **risk-proportional**: payments module gets the full trail, UI components get code + tests.

### 4. Verification Requires Structure

You CAN run mypy on a 5000-line blob. That's syntactic verification — it catches type errors.

You CANNOT semantically verify a blob — "does this code do what it's supposed to do?" — without knowing what it's supposed to do. That requires a spec. A spec requires decomposition.

**Decomposition for generation and decomposition for verification are the same act.** Specification and constraint are two views of one thing. They're inseparable.

> SDP is not a linter (post-hoc verification) and not a build system (generation control). It's both — because decomposition IS both.

---

## Two Modes, One Philosophy

### `sdp ship` — Autopilot with Verification

You describe what you want. SDP decomposes, generates, verifies, records, ships. You see a progress bar and a PR.

**For:** Teams that trust the framework. Rapid iteration. "I want verified code, fast."

### `sdp drive` — You're at the Wheel

You describe what you want. SDP decomposes — but stops at every fork and asks you. Your expertise shapes every spec. Your decisions are recorded.

**For:** Engineers who own the outcome. New domains. High-stakes systems. Compliance environments.

### The Accountability Thesis

In both modes, the human is accountable. SDP doesn't remove responsibility — it gives you **evidence to stand behind your decisions.** The artifacts prove what was specified, what was verified, and what passed.

The difference between "I think it works" and "here's proof it works" is the difference between a blog post and a court filing.

---

## The Five Bets

1. **AI will write most code by 2028.** Not controversial.
2. **Trust will be the bottleneck.** Controversial today. Obvious after the first major AI-code incident.
3. **Decomposition is a permanent advantage.** Human verification capacity is finite. This doesn't change with better models.
4. **Static analysis beats cross-model review on ROI.** Free tools that catch 35-50% of bugs > expensive tools that catch 15% more.
5. **Data is the moat.** The team with the best dataset on "what AI gets wrong" wins. Not the team with the best architecture.

If bet 3 is wrong — if humans can somehow verify 5000-line blobs — SDP is unnecessary.

We bet they can't.

---

## Who This Is For

**Code that costs more to fix than to write.**

| Segment | The Pain |
|---------|----------|
| Fintech | Payment bugs = chargebacks + regulatory fines |
| Healthcare | Data bugs = HIPAA violations |
| Infrastructure | Config bugs = outages at scale |
| Enterprise SaaS | Multi-tenant bugs = data leaks |
| ML/Data | Pipeline bugs = silent corruption |

Enterprise interest is validated:
- A top-3 bank — contracted
- A major airline — contracted
- The largest online marketplace — evaluating

**Not for:** Landing pages, MVPs, prototypes, disposable code. If rewriting is cheaper than verifying — don't verify. Just ship.

---

## The 12-Month Window

Cursor and GitHub will ship "verified AI code" features within 18 months. They have distribution, model access, and enterprise demand.

SDP's advantage: **speed and focus.** Ship the wedge now. Build the dataset now. Become the tool enterprises require — or the engine Cursor embeds.

The clock is running.

---

*SDP Vision v1.0 — February 2026*
