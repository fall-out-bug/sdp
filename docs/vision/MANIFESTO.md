# SDP v2: The Trust Layer for AI-Generated Software

> *"The point of SDP isn't plan + build. It's verified build with decomposition and adversarial review. No other tool does this automatically."*

---

## The Problem

In 2026, AI writes code. That's solved.

The unsolved problem: **can you trust it?**

- A landing page built by AI that's slightly wrong? Rebuild it. Cost: $0.
- A payment system built by AI that's slightly wrong? Chargebacks, lawsuits, regulatory fines. Cost: $100K+.

As AI writes more code, the **trust gap** widens. Models get better at generating. They don't get better at guaranteeing. The gap between "it looks right" and "it IS right" is where the damage happens.

**SDP exists to close that gap.**

---

## What SDP Is (v2)

SDP is a **trust standard for AI-generated software.**

Not a coding framework. Not a project management tool. Not a skill taxonomy.

A trust standard. Three words:

> **Decompose. Verify. Audit.**

### Decompose

AI generates better code in small, focused units than in large, sprawling ones. SDP decomposes features into workstreams — atomic units with clear acceptance criteria and dependency order. The model works within its reliable context window, not beyond it.

### Verify

Every generated unit passes through:
1. **Tests** — written before implementation (TDD), catching interface mismatches
2. **Type checking** — `mypy --strict` or language equivalent. Highest ROI quality gate for AI code.
3. **Cross-model review** — a *different* model checks the work. Correlated failures between same-model implementer and reviewer are the #1 blind spot in AI coding. Cross-model review breaks the correlation.

### Audit

Every change is traceable:
- Which model generated it
- Against which spec
- With what acceptance criteria
- Reviewed by which verifier
- At what cost (tokens, time, dollars)

---

## What SDP Is NOT

- **Not a coding framework.** SDP doesn't tell agents *how* to code. TDD, clean architecture, file size limits — these are opinions, not protocol. Agents apply their own engineering judgment.
- **Not a project manager.** SDP doesn't manage sprints, estimate timelines, or track velocity. It verifies output.
- **Not a prompt library.** SDP doesn't provide better prompts. It provides a verification layer that works regardless of prompt quality.
- **Not for landing pages.** If your code is disposable, you don't need trust guarantees. SDP is for code that costs more to fix than to write.

---

## The Architecture

### Two Layers, Cleanly Separated

```
┌─────────────────────────────────────────┐
│           SDP Protocol (the spec)        │
│                                          │
│  • Workstream format                     │
│  • Agent contract interface              │
│  • Verification gate interface           │
│  • Checkpoint/audit format               │
│  • Coordination semantics                │
│                                          │
│  Anyone can implement this.              │
│  This is the HTTP of agent coordination. │
└─────────────────────────────────────────┘
                    │
                    │ implements
                    ▼
┌─────────────────────────────────────────┐
│         SDP Framework (one impl)         │
│                                          │
│  • Go binary (dispatcher, verifier)      │
│  • @ship command (one-command UX)        │
│  • Git hooks + CI integration            │
│  • Cross-model review engine             │
│  • Agent plugin SDK                      │
│                                          │
│  The reference implementation.           │
│  Others can build their own.             │
└─────────────────────────────────────────┘
```

### The Protocol (sdp-spec)

A specification document. Not code. Defines:

| Concept | What It Specifies |
|---------|-------------------|
| **Workstream** | Format for atomic task units: goal, acceptance criteria, scope, dependencies |
| **Agent Contract** | How agents declare capabilities, quality guarantees, cost profile |
| **Verification Gate** | Interface for quality checks (tests pass, types check, review approved) |
| **Checkpoint** | Format for saving/restoring execution state across sessions |
| **Audit Entry** | Format for tracing who did what, when, with which model, at what cost |
| **Coordination** | How multiple agents divide work, avoid conflicts, merge results |

This is a *standard*. Like OpenAPI defines how APIs describe themselves, SDP Protocol defines how AI agents describe, coordinate, and verify their work.

### The Framework (sdp)

A Go binary that implements the protocol. Ships with:

**One command for users:**
```
sdp ship "Add OAuth2 login with Google and GitHub"
```

What happens inside (invisible to user):
1. Scans codebase (what's here, what tech stack, what patterns)
2. Decomposes into workstreams with dependency order
3. For each workstream: generates tests → implements → type-checks
4. Cross-model review (sends to a different provider for verification)
5. Produces a PR with full audit trail

**Two integration paths:**

1. **Git hooks + CI** — SDP as quality gate in the pipeline. Any code pushed goes through verification. Developers don't install SDP; the *repo* has SDP.

2. **Agent SDK** — Libraries that agent platforms embed. Cursor calls `sdp.Decompose(task)` internally. Replit calls `sdp.Verify(code)` before deploying. SDP is the engine inside, invisible to users.

---

## The User Experience

### For Developers (Direct Use)

```bash
# Install
go install github.com/user/sdp@latest

# The only command you need
sdp ship "Add user authentication with OAuth2"

# Watch it work
Decomposing... 3 workstreams identified
  [1/3] Backend auth service ████████████ done (87% coverage)
  [2/3] Frontend login flow   ████████████ done (92% coverage)
  [3/3] Integration tests     ████████████ done
Cross-model review... approved ✓
PR created: github.com/org/repo/pull/42
```

That's it. No skills to learn. No workstream IDs to manage. No guards to activate. No questions to answer.

Want more control? It's there:
```bash
# See what it planned
sdp ship "Add auth" --explain

# Review workstreams before execution
sdp ship "Add auth" --approve-plan

# Force a specific decomposition
sdp ship "Add auth" --workstreams=2

# Skip cross-model review (fast mode, less trust)
sdp ship "Add auth" --skip-review
```

### For CI/CD (Invisible)

```yaml
# .github/workflows/sdp-verify.yml
on: [pull_request]
jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: sdp-dev/verify-action@v1
        with:
          gates: [tests, types, review]
          review-model: "gpt-4o"  # Cross-model review
```

Every PR gets verified. Developers never think about SDP.

### For Agent Platforms (SDK)

```go
// Inside Cursor/Replit/Claude Code
import "github.com/user/sdp/sdk"

func handleUserRequest(prompt string) {
    // Decompose into verified units
    plan := sdp.Decompose(prompt, codebase)
    
    for _, ws := range plan.Workstreams {
        code := agent.Generate(ws.Spec)
        result := sdp.Verify(code, ws.AcceptanceCriteria)
        if !result.Passed {
            code = agent.Fix(code, result.Issues)
        }
    }
    
    // Cross-model review
    review := sdp.CrossModelReview(allChanges, "gpt-4o")
    
    // Audit trail
    sdp.RecordAudit(plan, review, tokenCost)
}
```

The platform uses SDP. The user never knows.

---

## Who This Is For

### Primary: Deeptech & Enterprise

Code that costs more to fix than to write:
- **Fintech** — Payment processing, regulatory compliance
- **Healthcare** — Patient data, HIPAA compliance
- **Infrastructure** — Cloud provisioning, network config
- **ML/Data** — Pipeline correctness, model serving
- **Enterprise SaaS** — Multi-tenant, data isolation

### Secondary: AI Agent Platforms

Tools that want "quality built in":
- **IDE agents** (Cursor, Windsurf, Continue) — embed verification
- **Autonomous agents** (Devin, SWE-Agent) — coordination protocol
- **CI/CD platforms** (GitHub Actions, GitLab CI) — verification gates

### Not For:

- Landing pages and portfolio sites
- Disposable prototypes
- Vibe coding experiments
- Projects where "works on my machine" is good enough

---

## The Moat

SDP's moat is not code. Code is reproducible.

SDP's moat is **accumulated trust knowledge**:
- Which verification gates catch which failure modes
- How to decompose features for optimal AI generation quality
- What cross-model review patterns break correlated failures
- How coordination protocols prevent multi-agent conflicts

This knowledge deepens with every verified build. Every failure mode caught. Every cross-model disagreement resolved. It's the engineering knowledge equivalent of Google's search index — built through usage, not engineering alone.

---

## The Roadmap

### Phase 1: @ship (The One-Command UX)

Make the default experience zero-ceremony:
- `sdp ship "description"` — one command, invisible internals
- Quality gates run silently (tests, types, file size)
- Cross-model review for every generation
- Audit trail recorded automatically

### Phase 2: CI/CD Integration (The Invisible Gate)

SDP as a GitHub Action / GitLab CI step:
- Every PR verified automatically
- Cross-model review on every AI-generated diff
- Trust scores per repository, per developer, per model
- Dashboard for engineering leads

### Phase 3: Agent SDK (The Invisible Engine)

SDP as a library agent platforms embed:
- `sdp.Decompose()` — break features into verified units
- `sdp.Verify()` — run all gates on generated code
- `sdp.CrossReview()` — send to different model for check
- `sdp.Audit()` — record everything

### Phase 4: Protocol Standard (The HTTP Moment)

SDP Protocol v1.0 as an open specification:
- Workstream format standard
- Agent contract standard
- Verification gate interface
- Published independently of any implementation
- Multiple implementations by different teams/companies

---

## The Bet

We bet that:

1. **AI will write most code by 2028.** Not controversial.
2. **Trust will be the bottleneck, not generation.** Controversial today. Obvious in 2028.
3. **Verification requires structure.** You can't verify a 5000-line monolithic AI generation. You can verify 10 focused units of 200 lines each.
4. **Cross-model review breaks correlated failures.** Same model reviewing its own code is theater. Different model catches different bugs.
5. **The invisible standard wins.** HTTP won because users never think about it. SDP wins when developers never think about it.

If any of these bets are wrong, SDP is unnecessary. If all five are right, SDP is essential infrastructure.

---

*SDP v2 — February 2026*
