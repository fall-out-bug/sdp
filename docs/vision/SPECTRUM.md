# SDP: The Spectrum

> Wearing a tie, but t-shirts on Fridays.

---

## Two Poles

On one end — a jam session. A personal AI assistant that listens and does. You're alone, you're in the zone, you owe nothing to anyone. That's OpenClaw, that's Copilot in free mode, that's vibe coding in its purest form. Territory of lawlessness. And it's beautiful — for a Saturday prototype, for an experiment, for the joy of it.

On the other end — a symphony orchestra. A conductor. A score. Every instrument knows its part. The performance is recorded. You can audit any measure, any note. That's SDP. Law and order. And it's necessary — when someone is accountable for the code, when the auditor arrives, when production goes down at 3 AM.

```
Jam session                              Orchestra
(OpenClaw, vibe coding)                  (SDP)
─────────────────────────────────────────────────────
Personal                                 Team
"Just do it"                             "Why, what, how, prove it"
One agent, any channel                   Many agents, one protocol
No traces                                Every step recorded
Freedom                                  Accountability
Instant result                           Verifiable result
```

## Both Are Legitimate

This is not "bad approach vs. good approach." It's a spectrum. Context determines position:

| Situation | Pole |
|-----------|------|
| Saturday prototype | Jam session |
| MVP for investors | Closer to orchestra |
| Fintech in production | Orchestra |
| Hackathon | Jam session |
| SOC2 audit | Orchestra |
| Personal pet project | Jam session |
| Code you're accountable for | Orchestra |

The question is never "what's better." The question is always "what's needed right now."

---

## SDP Doesn't Kill the Vibe

SDP makes the vibe **accountable.**

Discovery — it's the same conversation with AI. You say "I want OAuth." AI asks "why JWT, not sessions?" You answer. That's the vibe. But now that conversation is **recorded** — and six months later, when someone asks "why JWT?", the answer isn't in someone's head, it's in the evidence log.

Delivery — it's the same code generation. AI writes, tests run, coverage counts. That's the vibe. But now every generation is **recorded** — which model, which spec, which tests, what result. And when something breaks, `sdp trace` shows the full chain in 5 seconds.

**The protocol doesn't replace creativity. The protocol records it.**

---

## Progressive Disclosure: From T-shirt to Tie

The key design principle: the protocol plugs in layer by layer, not as a wall.

### Level 0 — T-shirt (Friday)

You're solo. Small project. You need results, not process.

```
@feature "Add dark mode" → @oneshot F01 → done
```

SDP works, evidence is recorded, but you don't think about it. Protocol runs in the background. You're in flow.

### Level 1 — Smart Casual (Team)

Two or three of you. Coordination needed. Who's doing what. Any conflicts?

```
@feature "Add OAuth" → workstreams with dependencies
@oneshot F01 → parallel execution
contracts → "Auth and Payments both touch User — here's the contract"
beads → task tracking across sessions
```

The protocol helps you avoid stepping on each other's toes. Evidence helps understand what happened while you were away.

### Level 2 — Tie (Enterprise)

Audit. Regulation. Accountability.

```
sdp trace abc1234 → full chain from decision to deploy
compliance export → SOC2/HIPAA/EU AI Act
signed evidence → cryptographic confirmation
risk-proportional verification → auth/ full trail, components/ light
```

Same protocol. Same tools. But now evidence isn't just a log — it's a legal artifact.

### One Protocol, Three Levels

```
T-shirt         Smart Casual        Tie
(solo)          (team)              (enterprise)
    │               │                   │
    ▼               ▼                   ▼
 @feature        contracts          compliance export
 @oneshot        beads              signed evidence
 evidence        scope collision    audit trail
 (background)    (coordination)     (governance)
```

**You don't switch protocols. You raise the level.**

Project grew? Enable beads. Auditor arrived? Enable compliance export. Evidence is already there — it's been recording since day one, back when you were in a t-shirt.

---

## Our Vibe

SDP is not anti-vibe-coding. It's vibe coding that grew up.

When you're solo on Saturday — create. When you're on a team — coordinate. When you answer to regulators — prove it. The protocol is the same. Depth adjusts to need.

**Wearing a tie, but t-shirts on Fridays.**

---

*SDP Spectrum v1.0 — February 14, 2026*
*From jam session to symphony. One protocol, three volume levels.*
