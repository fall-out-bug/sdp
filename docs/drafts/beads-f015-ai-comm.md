## Context & Problem

F015 addresses multiple communication gaps:
1. **AI не понимает 'зачем'** — AI фокусируется на механике, а не на целях продукта
2. **Потеря контекста** — Решения забываются, нет audit trail  
3. **Разрозненные навыки** — Каждый навык работает изолировано без общего контекста

## Goals & Non-Goals

**Goals:**
- Give AI agents product vision context (зачем, не только как)
- Create decision audit trail for all architectural decisions
- Auto-load PRODUCT_VISION.md in all skills
- Enable new team members to understand decision history

**Non-Goals:**
- Force all decisions to follow rigid process
- Replace human deliberation with auto-logging
- Block work on misaligned features (warning only)

## Primary Users

1. **AI агенты** — Skills (@idea, @design, @build) получают контекст
2. **Разработчики** — Читают decision log, понимают 'зачем'
3. **Product менеджеры** — Видят alignment с PRODUCT_VISION.md
4. **Новые члены команды** — Onboarding через историю решений

## Product Alignment

**Mission:** Improve AI-human communication through vision context and decision logging
**Vision Alignment:** Directly supports PRODUCT_VISION.md mission of "AI agents and humans collaborate on building reliable, maintainable software"

## Technical Approach

**PRODUCT_VISION.md Integration:**
- Auto-load in all skills at startup
- Validation in @idea (check alignment before creating task)
- Explicit misalignment flag when vision conflicts with request

**Decision Logging:**
- Auto-log AskUserQuestion responses as mini-decisions
- Explicit commit for major architectural decisions
- Hybrid: Auto-log for questions, explicit for big decisions

**Storage Format:**
- JSONL: `.sdp/decisions.jsonl` (append-only, machine-readable)
- Markdown: `docs/DECISIONS.md` (human-readable)
- Beads metadata: Store decision IDs in task metadata

## Success Metrics

1. **Alignment Score:** % features aligned-checked with PRODUCT_VISION.md
2. **Decision Coverage:** % architectural decisions with audit trail
3. **AI Context Hits:** # times AI used vision in reasoning

## Technical Implementation

### Components

1. **VisionLoader** — Load PRODUCT_VISION.md, validate schema
2. **DecisionLogger** — Append to JSONL + MD formats
3. **SkillIntegration** — Modify @idea, @design, @build to use vision
4. **AlignmentChecker** — Check feature against vision, flag misalignments

### Data Structures

**Decision Entry (JSONL):**
```json
{
  "timestamp": "2026-01-28T10:00:00Z",
  "type": "ask_user_question" | "explicit_commit",
  "feature_id": "bd-0001",
  "question": "What is the primary problem?",
  "answer": "User pain point",
  "context": {...}
}
```

**VISION.md Schema:**
```yaml
mission: str
users: List[str]
success_metrics: List[Dict[str, str]]
strategic_tradeoffs: Dict[str, str]
non_goals: List[str]
```

## Concerns & Risks

1. **Overhead:** Auto-logging might slow down skills
   - *Mitigation:* Async logging, batch writes

2. **Vision staleness:** PRODUCT_VISION.md might become outdated
   - *Mitigation:* Validation warnings, "last updated" checks

3. **Decision noise:** Auto-logging every answer creates clutter
   - *Mitigation:* Significance filtering, manual commit for major decisions

## Tradeoffs

**Auto-load vs Manual:**
- Auto-load: Always fresh, zero friction
- Manual: Explicit intent, less magic
- **Decision:** Auto-load with opt-out flag

**JSONL vs DB:**
- JSONL: Simple, no dependencies, append-only
- DB: Queryable, transactions, complex
- **Decision:** JSONL (follows audit log pattern from F014)

## Open Questions

1. Should vision updates invalidate existing decisions? (Probably yes)
2. How to handle conflicting decisions? (Show both, ask for resolution)
3. Decision expiration? (Keep all, but mark "superseded by")

## Implementation Workstreams (Draft)

1. **WS-015-01:** VisionLoader + PRODUCT_VISION.md schema validation
2. **WS-015-02:** DecisionLogger (JSONL + MD)
3. **WS-015-03:** @idea skill integration (alignment checking)
4. **WS-015-04:** @design, @build skills integration (vision context)
5. **WS-015-05:** Decision log viewer CLI (`sdp decisions`)
6. **WS-015-06:** Tests + documentation

---

## Interview Answers

{
  "mission": "\u0423\u043b\u0443\u0447\u0448\u0438\u0442\u044c \u043a\u043e\u043c\u043c\u0443\u043d\u0438\u043a\u0430\u0446\u0438\u044e",
  "problem": "\u0412\u0441\u0451 \u0432\u044b\u0448\u0435\u043f\u0435\u0440\u0435\u0447\u0438\u0441\u043b\u0435\u043d\u043d\u043e\u0435 (AI \u043d\u0435 \u043f\u043e\u043d\u0438\u043c\u0430\u0435\u0442 '\u0437\u0430\u0447\u0435\u043c', \u043f\u043e\u0442\u0435\u0440\u044f \u043a\u043e\u043d\u0442\u0435\u043a\u0441\u0442\u0430, \u0440\u0430\u0437\u0440\u043e\u0437\u043d\u0435\u043d\u043d\u044b\u0435 \u043d\u0430\u0432\u044b\u043a\u0438)",
  "users": [
    "Product \u043c\u0435\u043d\u0435\u0434\u0436\u0435\u0440\u044b",
    "\u0420\u0430\u0437\u0440\u0430\u0431\u043e\u0442\u0447\u0438\u043a\u0438",
    "AI \u0430\u0433\u0435\u043d\u0442\u044b",
    "\u041d\u043e\u0432\u044b\u0435 \u0447\u043b\u0435\u043d\u044b \u043a\u043e\u043c\u0430\u043d\u0434\u044b"
  ],
  "technical_approach": "\u0413\u0438\u0431\u0440\u0438\u0434\u043d\u044b\u0439 \u043f\u043e\u0434\u0445\u043e\u0434 (\u0430\u0432\u0442\u043e\u0437\u0430\u0433\u0440\u0443\u0437\u043a\u0430 + validation + \u044f\u0432\u043d\u043e\u0435 \u043d\u0435\u0441\u043e\u043e\u0442\u0432\u0435\u0442\u0441\u0442\u0432\u0438\u0435)",
  "decision_logging": "\u0413\u0438\u0431\u0440\u0438\u0434\u043d\u044b\u0439 (auto-log \u0434\u043b\u044f \u0432\u043e\u043f\u0440\u043e\u0441\u043e\u0432, explicit \u0434\u043b\u044f \u0431\u043e\u043b\u044c\u0448\u0438\u0445 \u0440\u0435\u0448\u0435\u043d\u0438\u0439)",
  "storage": "\u041a\u043e\u043c\u0431\u0438\u043d\u0430\u0446\u0438\u044f JSONL + MD",
  "metrics": [
    "Alignment score",
    "Decision coverage",
    "AI context hits"
  ]
}
