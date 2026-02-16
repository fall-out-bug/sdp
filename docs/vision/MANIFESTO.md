# SDP: Agents Are Roles, Not People. Accountability Is Real.

> *"Agents are roles, not people. Accountability is real."* [¹](#footnotes)

---

## Tech Lead Vasily

You created a task in the tracker. Timofey, a friendly analyst who loves talking about his travels, pulls you into a video call. He gently but persistently asks about your requirements, then sends you meeting minutes, creates all the subtasks, assigns the right people.

Next, Anastasia catches you — a stern architect with a steely gaze. She lectures you that your ideas violate the project's architecture. You argue. Eventually, you agree on specifications. She locks everything down, writes specs to the wiki, draws C4 diagrams. You see Timofey and Anastasia arguing in the comments — work is humming. Disagreements resolve fast.

Then Vasily, the tech lead, takes over. Cheerful, covered in tattoos, built like Dwayne Johnson — at least that's what his video avatar looks like. You joke around, polish the design, he grabs everything into sprint. Meanwhile he mentions the team is chipping in for a gift for Kostya, the QA engineer. You don't remember Kostya, but you politely ask what they're getting him. "Claude tokens!" Vasily jokes. Everyone's a vibe-coder these days. But you know the project runs strict TDD — every sneeze is tested. The product is nearly bug-free.

By evening, you have working prototypes. By morning, the full feature. You're not worried — the team delivers like this every time.

You stand up from your chair in the small office. Almost the entire team is remote. You glance at Georgy, the senior DevOps, perpetually tinkering with GPUs in your corporate Kubernetes. Strange — you barely have any ML. Though Alexander sits next to him — the lead ML specialist, lost in thought. He doesn't even say hello. Why are the only people in the office such NPCs? You wish Vasily would come by. You'd grab a beer. You want to see real teammates, not these two.

**Except Timofey, Anastasia, Kostya, and Vasily aren't real people. They're agents.**

And Georgy and Alexander are the ones cleaning Anastasia's hallucinations from your DMs.

Or maybe you're an agent too. Alexander's experiment.

Scary? Not to me. I like it. [¹](#footnotes)

---

## The Problem Under the Story

The multi-agent future is coming. The question isn't whether AI agents will write code, run reviews, and ship features. They will. The question is:

**When something breaks — and it will — what do you have?**

Today: nothing. A git blame pointing at a rubber-stamped PR. "Something vibed itself into existence. No idea how or when."

That answer will get increasingly expensive. Not because of one party's liability — the legal frameworks are still being written [²](#footnotes). But because **someone will need the evidence.** The developer. The model provider. The auditor. The incident responder.

> **SDP doesn't decide who's responsible. SDP records what happened.**

---

## What SDP Is

**An open protocol for recording what AI agents did to your codebase.**

Three guarantees:

1. **Provenance.** Every piece of AI-generated code records: which model, which version, what parameters, what spec, when, who initiated. Whether it was Vasily or Claude — you'll know.

2. **Verification evidence.** Not "tests passed" — the actual pytest output. Not "types checked" — the actual mypy report. Real output, not assertions.

3. **Forensic trace.** When production breaks, trace from the git commit back through the chain: what was generated → what was verified → who approved. The reconstruction tool for the 3 AM page.

---

## The Philosophy

### Provenance Is Permanent

Models will improve. Context windows will grow. Maybe one day AI generates flawless 5000-line files from a paragraph.

But the question "which model wrote this, from what spec, and who approved it?" **never stops being valuable.** Provenance survives every regime change in AI capability.

### Decomposition Helps Today

Breaking features into small verified units improves AI generation quality and makes human review possible. Code review effectiveness drops to near-random above 400 lines of diff [³](#footnotes). This may or may not be permanent — but it works now, and SDP supports it.

The protocol works **with or without decomposition.** Provenance is the invariant. Decomposition is a powerful technique.

### Forensics Over Verification

Everyone claims "verified AI code." Nobody has a forensic chain from a production incident back to the generation spec.

Verification answers: "did this pass?" Forensics answers: "what happened, when, and who decided?"

SDP is the forensic system. The first one for AI-generated code.

### Evidence, Not Opinions

SDP is a neutral ledger. It doesn't embed opinions about who's liable. It records facts: what was generated, how it was verified, who approved it. The courts, the auditors, the incident reviewers — they decide responsibility. SDP gives them the record.

---

## The UX: plan / apply / incident

Borrowed from Terraform — because it works.

**`sdp plan`** — Timofey and Anastasia. Show what you're about to do. Decompose the feature, display the units, estimate the cost. Don't execute yet.

```bash
$ sdp plan "Add OAuth2 login with Google and GitHub"

  3 units planned:
  1. Backend OAuth service    [high risk — auth module]
  2. Frontend login component [medium risk]
  3. Integration tests        [low risk]

  Estimated: ~3 min | ~$0.15
  [apply] [edit] [cancel]
```

**`sdp apply`** — Vasily. Execute the plan. Each unit: generate → verify → record. Streaming progress. Per-unit rollback if something fails.

```
Applying...
  [1/3] Backend OAuth ████████████ done
        model: claude-sonnet-4 | types ✓ | tests ✓ (91%) | provenance ✓
  [2/3] Frontend      ████████████ done
  [3/3] Integration   ████████████ done

All gates passed. Evidence recorded.
PR created: github.com/org/repo/pull/42
```

**`sdp incident`** — Georgy and Alexander, cleaning up at 3 AM. Trace from a commit back through the full chain.

```bash
$ sdp incident abc1234

  Commit abc1234 — AI-generated
  ├── Model: claude-sonnet-4-20250514
  ├── Spec: "Backend OAuth service with Google + GitHub"
  ├── Verification: types ✓, tests ✓ (91%), semgrep ✓
  ├── Approved by: @developer at 2026-02-08T14:32Z
  └── Evidence: .sdp/log/2026-02-08-oauth.json
```

### Modes

```bash
sdp plan "Add auth"                   # Show plan, wait for approval
sdp plan "Add auth" --auto-apply      # Ship mode: plan + apply, no stop
sdp plan "Add auth" --interactive     # Drive mode: you decide at every fork
sdp apply                             # Execute last plan
sdp apply --retry 3                   # Retry failed unit only
sdp incident <commit>                 # Forensic trace
```

---

## The Strategy: Protocol First

What already exists is a Claude Code plugin with 19 agent roles [⁴](#footnotes). That's a protocol implementation — not a CLI.

**The order:** Protocol (open standard) → Plugin → CLI → SDK → Ecosystem.

Why:
- Open protocol gets embedded everywhere. If Cursor, Replit, and others can speak SDP — the standard wins.
- High-assurance teams adopt specs, not tools. They can implement the protocol inside their own infra.
- Plugin-first is the fastest path to real users.

### Public OSS Scope

| Layer | License |
|-------|---------|
| Protocol schema | CC-BY (open) |
| Verification engine | Open-source (MIT) |
| Orchestration + evidence analysis | Open-source (MIT) |
| Public adoption layers (`L0-L2`) | Open-source (MIT) |

---

## Who This Is For

**Code that costs more to fix than to write.**

Fintech. Healthcare. Infrastructure. High-assurance SaaS. Regulated industries.

**Not for:** Landing pages, MVPs, prototypes. If rewriting is cheaper than recording — don't record. Just ship.

---

## The Five Bets

1. **AI will write most code by 2028.**
2. **"What happened?" will be mandatory.** Courtrooms, audits, incident reviews will demand records.
3. **Provenance is permanent.** Regardless of how code is generated — the record never stops being valuable.
4. **Forensics > Verification.** "Did it pass?" is table stakes. "What happened?" is the moat.
5. **Protocol > Product.** The team that sets the evidence standard wins.

---

## Kill Criteria

> After 500 SDP runs: if verification catch rate < 5% and post-merge defect rate equals baseline — kill the product. Specific. Testable. Honest.

---

## Footnotes

¹ From [the original post](https://t.me/data_intensive_boar/25), November 13, 2025. The multi-agent vision that became SDP — written before the "Swarm" hype, before the agent frameworks gold rush.

² Accountability for AI-generated code is being contested across multiple legal frameworks: strict developer liability, shared model-provider liability, product liability on tool vendors, and negligence-based frameworks. SDP is framework-agnostic — it records evidence, not verdicts.

³ Microsoft and Google internal studies on code review effectiveness. Review quality drops to near-random above ~400 lines of diff.

⁴ The SDP codebase includes a Go-based engine with decomposition, dependency graphs, parallel dispatch, circuit breakers, checkpoint recovery, and a synthesis engine — running as a Claude Code skill system since 2025.

---

*SDP Manifesto v4.0 — February 2026*
*Born from a Telegram post about Tech Lead Vasily.*
*Protocol-first. Evidence-first. Human-first.*
