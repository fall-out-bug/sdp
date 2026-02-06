# Business Requirements: Long-Term Memory for SDP

**Feature Code:** LTM-001
**Status:** Draft
**Version:** 1.0
**Date:** 2026-02-06

---

## Executive Summary

SDP currently lacks institutional memory. Developers repeat mistakes, forget past decisions, and cannot learn from historical patterns. This feature implements a comprehensive long-term memory system that captures decisions, usage patterns, session analytics, and project history to enable continuous improvement.

**Business Value:**
- Reduce repeated mistakes by 40%
- Speed up decision-making with historical context
- Enable data-driven process improvements
- Create searchable project intelligence

---

## Stakeholders

| Stakeholder | Role | Interests |
|-------------|------|-----------|
| **Solo Developers** | Primary users | Quick access to past decisions, avoid repeating mistakes, learn from personal patterns |
| **Development Teams** | Collaborative users | Shared decision history, team patterns, coordinated workflow, knowledge transfer |
| **Project Maintainers** | Long-term owners | Historical context, project evolution tracking, onboarding new team members |
| **AI Agents (Claude)** | Decision makers | Access to historical decisions, pattern recognition, better recommendations |
| **SDP Plugin Maintainers** | Platform owners | Aggregate usage patterns, product improvements, telemetry insights |

---

## Problem Statement

### Current Pain Points

**P1 - Repeated Mistakes (Critical):**
- Developers forget why previous approaches failed
- Same errors recur across different workstreams
- No institutional memory of "what didn't work"
- Example: "We tried PostgreSQL indexing for feature X, it didn't scale"

**P1 - Lost Decisions (Critical):**
- Decisions made during @feature or @design sessions are not captured
- Rationale for trade-offs is forgotten
- New team members cannot understand "why" behind architecture
- Example: "Why did we choose PostgreSQL over MongoDB?"

**P2 - Invisible Patterns (High):**
- No visibility into usage patterns (which commands, which workflows)
- Cannot identify bottlenecks in development process
- Don't know which features/workflows are successful
- Example: "80% of our workstreams fail on quality gates"

**P2 - No Session Intelligence (High):**
- Each SDP session is isolated
- Cannot analyze session effectiveness
- Don't know time distribution across activities
- Example: "We spend 60% of time debugging, not building"

**P3 - Fragmented History (Medium):**
- Decisions tracked in multiple places (git commits, markdown docs, code comments)
- No unified view of project evolution
- Difficult to trace "what changed when"
- Example: "When did we switch from X to Y?"

### Impact on Development

**Without Long-Term Memory:**
- 2-3 hours/week wasted re-learning past decisions
- 30% higher bug recurrence rate
- New developer onboarding takes 2x longer
- Cannot measure process improvement effectiveness

**With Long-Term Memory:**
- Instant context: "We tried X in WS-123, failed because Y"
- Pattern analysis: "Quality gates catch 85% of issues before review"
- Continuous improvement: "Session duration decreased 40% after TDD training"

---

## User Stories

### US-001: Поиск принятых решений (Decision Search)

**As a** developer working on a feature
**I want** to search past decisions by keyword, technology, or workstream
**So that** I don't repeat mistakes and can reference previous reasoning

**Acceptance Criteria:**

**Given:** Разработчик работает над функцией и сталкивается с выбором технологии
```
Developer: "Should I use PostgreSQL or MongoDB for this feature?"
```

**When:** Он ищет прошлые решения по ключевым словам
```bash
sdp memory search --query "PostgreSQL vs MongoDB"
# или
sdp memory search --tag "database" --feature "F01"
```

**Then:** Система показывает контекст прошлых решений
```markdown
Found 3 related decisions:

1. [F01-WS-005] 2026-01-15
   Decision: Use PostgreSQL for user data
   Rationale: ACID transactions required, relational data model
   Outcome: Successful, 99.9% uptime

2. [F01-WS-012] 2026-01-20
   Decision: Rejected MongoDB for analytics
   Rationale: Aggregation pipeline too complex, migrated to ClickHouse
   Outcome: Performance improved 5x

3. [F02-WS-003] 2026-02-01
   Decision: Hybrid approach (PostgreSQL + Redis)
   Rationale: PostgreSQL for persistence, Redis for caching
   Outcome: 50% latency reduction

💡 Pattern detected: PostgreSQL preferred for transactional data
```

**Value:** "Мы уже пробовали подход X в {workstream}, отказались из-за {reason}"

---

### US-002: Логирование решений (Decision Logging)

**As a** developer or AI agent making architectural decisions
**I want** to automatically log decisions with context and rationale
**So that** future developers understand the "why" behind choices

**Acceptance Criteria:**

**Given:** Разработчик выполняет @feature или @design
```
Claude: "Should we implement authentication ourselves or use a library?"
```

**When:** Decision is made during workflow
```bash
sdp memory log \
  --question "Authentication implementation approach" \
  --decision "Use auth0 library, not custom implementation" \
  --rationale "Security audits show custom auth has 3x more vulnerabilities" \
  --alternatives "custom,authlib,auth0" \
  --type "technical" \
  --workstream "00-001-05"
```

**Then:** Решение сохраняется в структурированном виде
```json
{
  "timestamp": "2026-02-06T10:30:00Z",
  "question": "Authentication implementation approach",
  "decision": "Use auth0 library, not custom implementation",
  "rationale": "Security audits show custom auth has 3x more vulnerabilities",
  "alternatives": ["custom", "authlib", "auth0"],
  "type": "technical",
  "workstream_id": "00-001-05",
  "feature_id": "F01",
  "decision_maker": "claude",
  "outcome": "pending"
}
```

**Automatic Triggers:**
- After @feature completion (vision decisions)
- After @design completion (architectural decisions)
- When @build fails quality gates (what broke)
- When @review finds issues (what to avoid)

**Value:** "Мы решили" — searchable decision log with full context

---

### US-003: Анализ паттернов использования (Usage Pattern Analysis)

**As a** development team lead
**I want** to see usage patterns and statistics over time
**So that** I can identify bottlenecks and optimize workflow

**Acceptance Criteria:**

**Given:** Команда хочет понять, как используется SDP
```bash
sdp memory analyze --period "30d"
```

**When:** Запуск анализа на основе телеметрии
```bash
sdp memory analyze --period "30d" --by "command,feature,outcome"
```

**Then:** Система показывает статистику паттернов
```markdown
📊 30-Day Usage Analysis
==========================

🎯 Command Usage:
  @build: 45 times (52%) - avg 12min per workstream
  @review: 20 times (23%) - 85% pass rate
  @design: 15 times (17%) - avg 3 features planned
  @feature: 6 times (7%) - avg 45min per feature

📈 Feature Success Rate:
  F01: 92% (12/13 workstreams completed)
  F02: 75% (9/12 workstreams completed) ⚠️
  F03: 100% (5/5 workstreams completed)

⏱️ Time Distribution:
  Building: 65% (9h 15m)
  Reviewing: 20% (2h 50m)
  Planning: 15% (2h 8m)

❌ Failure Patterns:
  1. Test coverage <80%: 8 occurrences (62% of failures)
  2. Type checking errors: 3 occurrences (23%)
  3. Architecture violations: 2 occurrences (15%)

💡 Insights:
  - Features planned with @design have 40% higher success rate
  - Test coverage is #1 bottleneck (consider TDD training)
  - F02 has abnormal failure rate (investigate feature complexity)
```

**Value:** Data-driven process optimization

---

### US-004: История сессий (Session History)

**As a** developer
**I want** to see detailed history of my SDP sessions
**So that** I can understand what I worked on and when

**Acceptance Criteria:**

**Given:** Разработчик хочет просмотреть историю сессий
```bash
sdp memory sessions --period "7d"
```

**When:** Запрос истории сессий
```bash
sdp memory sessions --period "7d" --detail
```

**Then:** Система показывает chronologically ordered sessions
```markdown
📅 Session History (Last 7 Days)
=================================

2026-02-06 (14:30-16:45) - 2h 15m
  Feature: F01 - User Authentication
  Workstreams: 00-001-05, 00-001-06
  Outcome: 2 completed, 0 failed
  Decisions Made: 2
    • Use auth0 for authentication (see: `sdp memory show D001`)
    • Rejected custom session management (see: `sdp memory show D002`)

2026-02-05 (10:00-11:30) - 1h 30m
  Feature: F02 - Payment Processing
  Workstreams: 00-002-03
  Outcome: 0 completed, 1 failed (quality gate)
  Failure Reason: Test coverage 72% (required ≥80%)
  Lesson Learned: Write tests first (see: `sdp memory lessons --tag "tdd"`)

2026-02-04 (09:00-12:00) - 3h 0m
  Feature: F03 - API Gateway
  Workstreams: 00-003-01, 00-003-02, 00-003-03
  Outcome: 3 completed, 0 failed
  Decisions Made: 1
    • Use gRPC for inter-service communication

📈 Summary:
  Total Sessions: 3
  Total Time: 6h 45m
  Success Rate: 83% (5/6 workstreams)
  Decisions Logged: 3
```

**Value:** Full traceability of development activity

---

### US-005: Извлеченные уроки (Extracted Lessons)

**As a** developer starting a new workstream
**I want** to see lessons learned from similar past work
**So that** I can avoid repeating mistakes

**Acceptance Criteria:**

**Given:** Разработчик начинает новый workstream
```bash
@build 00-001-07
```

**When:** SDP обнаруживает похожие прошлые workstreams
```bash
sdp memory lessons --workstream "00-001-07" --similar
```

**Then:** Система показывает извлеченные уроки
```markdown
💡 Lessons Learned for Similar Workstreams
==========================================

Based on 3 similar workstreams (database-related, backend):

🚫 Anti-Patterns to Avoid:
  1. Missing database migrations in tests
     Occurred in: 00-001-02, 00-002-05
     Impact: 2h debugging, 1 failed deployment
     Fix: Use test database fixtures (see `docs/testing-database.md`)

  2. Not indexing foreign keys
     Occurred in: 00-001-03
     Impact: 10x query performance degradation
     Fix: Add db_index=True in model definition

  3. Hardcoding database URLs
     Occurred in: 00-002-04
     Impact: Security issue, config management problems
     Fix: Use environment variables (see `docs/config.md`)

✅ Proven Patterns:
  1. Use pytest fixtures for test data
     Success rate: 100% (3/3 workstreams)
     Example: 00-001-05 (tests/conftest.py)

  2. Repository pattern for data access
     Success rate: 100% (2/2 workstreams)
     Example: 00-001-06 (src/repositories/)

📊 Risk Assessment:
  This workstream: MEDIUM complexity
  Recommended: Allocate 30% more time for testing
  Common pitfalls: Test data setup, migration handling

🔗 Related Decisions:
  • D001: Use PostgreSQL for all persistent data (2026-01-15)
  • D003: Always version database migrations (2026-01-20)
```

**Automatic Extraction:**
- When workstream fails: capture what went wrong
- When workstream succeeds: capture what worked well
- When review finds issues: capture quality patterns
- When similar workstreams complete: identify common patterns

**Value:** "Мы отказались" — avoid repeating failures

---

### US-006: Аналитика проекта (Project Analytics)

**As a** project maintainer
**I want** to see project-level analytics and trends
**So that** I can understand project health and evolution

**Acceptance Criteria:**

**Given:** Мейнтейнер хочет проанализировать проект
```bash
sdp memory analytics --project "myproject"
```

**When:** Запуск аналитики на основе всей истории
```bash
sdp memory analytics --project "myproject" --period "90d" --trends
```

**Then:** Система показывает проектную аналитику
```markdown
📊 Project Analytics: myproject
================================

📅 Period: Last 90 days (2024-11-08 to 2026-02-06)

📈 Development Velocity:
  Workstreams Completed: 47 (avg 5.2/week)
  Features Completed: 8 (avg 0.9/week)
  Success Rate: 89% (47/53 workstreams)

  Trend: ↗️ Improving (was 75% in previous period)

⏱️ Time to Complete:
  Median Workstream: 12min
  Median Feature: 2h 45m

  By Size:
    SMALL: 8min (avg)
    MEDIUM: 18min (avg)
    LARGE: 45min (avg)

🎯 Quality Metrics:
  Avg Test Coverage: 87% (↑ 5% from previous period)
  Type Safety: 98% (go vet pass rate)
  Architecture Violations: 2 (↓ from 8)

💡 Decision Patterns:
  Total Decisions: 34
    • Technical: 20 (59%)
    • Vision: 8 (24%)
    • Tradeoff: 6 (18%)

  Most Influenced:
    • PostgreSQL chosen 12 times (100% success rate)
    • Custom implementation rejected 8 times (avg 3x effort saved)

🚫 Recurring Issues:
  1. Test coverage below 80%: 6 occurrences
     Resolved by: TDD training (decreased to 1 occurrence)

  2. Missing error handling: 4 occurrences
     Resolved by: Error validator added to @review

📊 Workstream Patterns:
  Most Common Size: MEDIUM (65%)
  Most Common Feature Type: Backend API (40%)

  Completion by Feature:
    F01 (Auth): 92% (12/13)
    F02 (Payments): 75% (9/12) ⚠️
    F03 (API Gateway): 100% (8/8)

💰 ROI Metrics:
  Time Saved by Reusing Decisions: ~6h
  Time Saved by Pattern Recognition: ~4h
  Total Time Saved This Period: ~10h

🔮 Predictions:
  Next period velocity: 6.0 workstreams/week (trending up)
  Recommended focus: Improve F02 success rate (investigate complexity)

📈 Health Score: 85/100 (Good)
  ↑ Velocity: 18/20
  ↑ Quality: 19/20
  → Consistency: 16/20
  ↑ Learning: 17/20
  ↓ Documentation: 15/20 (improve decision docs)
```

**Value:** Project intelligence and continuous improvement

---

### US-007: Контекст при принятии решений (Decision Context)

**As an** AI agent (Claude)
**I want** to access historical decisions and patterns
**So that** I can provide better recommendations

**Acceptance Criteria:**

**Given:** Claude помогает разработчику с @feature
```
User: @feature "Add real-time notifications"
Claude: I'll help you design this feature...
```

**When:** Claude предлагает техническое решение
```bash
# Claude internally queries:
sdp memory query --context "notifications,real-time,websocket" --format json
```

**Then:** Claude получает контекст для建议
```json
{
  "related_decisions": [
    {
      "workstream": "00-001-08",
      "decision": "Use WebSocket for real-time features",
      "rationale": "Bidirectional communication required, low latency",
      "outcome": "successful, 50ms avg latency"
    },
    {
      "workstream": "00-002-10",
      "decision": "Rejected Server-Sent Events (SSE)",
      "rationale": "No bidirectional support, browser compatibility",
      "outcome": "migrated to WebSocket, improved UX"
    }
  ],
  "patterns": {
    "recommended": "WebSocket",
    "anti_patterns": ["polling", "SSE for chat"],
    "success_rate": "100% (2/2 workstreams with WebSocket)"
  },
  "lessons": [
    "Use connection pooling (see 00-001-08)",
    "Implement heartbeat for connection health (see 00-002-11)"
  ]
}
```

**Claude's Response:**
```
Based on project history, I recommend WebSocket for notifications:

✅ Proven in this project:
  • WS-008 (2026-01-20): WebSocket for chat, 50ms latency
  • WS-011 (2026-01-25): WebSocket for notifications, 99.9% uptime

❌ Approaches we rejected:
  • SSE (WS-010): No bidirectional support
  • Polling (WS-009): Too much server load

💡 Key lessons learned:
  1. Use connection pooling (reduced connections by 80%)
  2. Implement heartbeat (detected dead connections 2x faster)

Would you like me to design the WebSocket architecture?
```

**Value:** AI agents with institutional memory

---

### US-008: Экспорт и отчетность (Export & Reporting)

**As a** development team
**I want** to export memory data for reports and presentations
**So that** we can share insights with stakeholders

**Acceptance Criteria:**

**Given:** Команда хочет подготовить отчет о разработке
```bash
sdp memory export --format markdown --period "30d" --output report.md
```

**When:** Экспорт данных памяти
```bash
sdp memory export \
  --format markdown \
  --period "30d" \
  --include "decisions,sessions,analytics,lessons" \
  --output "monthly-report-2026-02.md"
```

**Then:** Система генерирует структурированный отчет
```markdown
# SDP Monthly Report: February 2026

## Executive Summary
- **Workstreams Completed:** 47 (89% success rate)
- **Features Delivered:** 8
- **Total Development Time:** 94 hours
- **Decisions Made:** 34
- **Time Saved from Pattern Reuse:** ~10 hours

## Key Achievements
✅ Delivered F01 (User Authentication) with 92% success rate
✅ Improved test coverage from 82% to 87%
✅ Reduced recurring bugs by 40% (through pattern detection)

## Major Decisions
1. [D001] Use PostgreSQL for all persistent data (2026-01-15)
   - Impact: 12 workstreams using this decision
   - Outcome: 100% success rate, 99.9% uptime

2. [D008] Adopt TDD workflow (2026-01-20)
   - Impact: Test coverage ↑ 5%, bug recurrence ↓ 40%
   - Outcome: ROI of 6h time saved

## Lessons Learned
### What Worked Well
- WebSocket for real-time features (2/2 success rate)
- Repository pattern for data access (100% success)
- Beads integration for task tracking

### What Didn't Work
- Custom authentication (abandoned, switched to auth0)
- MongoDB for analytics (migrated to ClickHouse)
- Monolithic approach (split into microservices)

## Trends
📈 Velocity: 5.2 → 6.0 workstreams/week (↑ 15%)
📈 Quality: 82% → 87% test coverage (↑ 5%)
📈 Success Rate: 75% → 89% (↑ 14%)

## Next Month Focus
1. Improve F02 (Payments) success rate (currently 75%)
2. Reduce workstream completion time (target: <10min median)
3. Document architectural decisions in ADR format

## Appendix
- [Full Decision Log](decisions-2026-02.md)
- [Session History](sessions-2026-02.md)
- [Quality Metrics](quality-2026-02.md)
```

**Other Export Formats:**
- JSON: For data analysis
- CSV: For spreadsheet import
- PDF: For presentations
- HTML: For web dashboards

**Value:** Stakeholder communication and project visibility

---

## Success Metrics (KPIs)

### Adoption Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| **Decision Logging Rate** | ≥80% of workstreams | `(logged decisions / completed workstreams) * 100` |
| **Memory Search Usage** | ≥5 searches/day | Average searches per active user per day |
| **Session History Views** | ≥10 views/week | Number of `sdp memory sessions` commands |
| **Export Usage** | ≥2 exports/month | Number of `sdp memory export` commands |

### Quality Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| **Mistake Recurrence Rate** | ↓40% from baseline | `(repeated mistakes / total mistakes) * 100` |
| **Decision Reuse Rate** | ≥50% | `(workstreams using past decisions / total workstreams) * 100` |
| **Pattern Detection Accuracy** | ≥80% precision | `(correct patterns / total patterns detected) * 100` |
| **Search Result Relevance** | ≥70% user satisfaction | User feedback on search results |

### Time Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| **Time Saved per Workstream** | ≥5min | Avg time difference with vs without memory |
| **Decision-Making Time** | ↓50% | Time to decision with historical context vs without |
| **Session Analysis Time** | <30s | Time to generate session analytics |
| **Onboarding Time** | ↓30% | Time for new developer to reach productivity |

### Project Health Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| **Feature Success Rate** | ≥90% | `(completed features / started features) * 100` |
| **Workstream Success Rate** | ≥85% | `(completed workstreams / started workstreams) * 100` |
| **Test Coverage Trend** | ↑ quarterly | Compare current to previous quarter |
| **Documentation Coverage** | ≥70% | `(documented decisions / total decisions) * 100` |

### ROI Metrics
| Metric | Target | Measurement |
|--------|--------|-------------|
| **Time Saved Monthly** | ≥10 hours/developer | Sum of time savings from decision reuse |
| **Bug Reduction** | ↓30% | Compare bug count before/after memory feature |
| **Development Velocity** | ↑20% | Workstreams completed per week before/after |
| **Developer Satisfaction** | ≥4/5 | Survey rating on usefulness of memory feature |

---

## Dependencies

### Existing Systems

**1. Telemetry System** (`internal/telemetry/`)
- **Current:** Basic event tracking (command_start, command_complete, ws_start, ws_complete, quality_gate_result)
- **Gaps:** No session analytics, no pattern analysis, no decision correlation
- **Needed:** Extended event types for memory feature

**2. Decision Tracking** (`internal/decision/decision.go`)
- **Current:** Data structures defined, but NOT implemented
- **Gaps:** No storage, no search, no logging mechanism
- **Needed:** Full implementation of decision log

**3. Checkpoint System** (`internal/checkpoint/checkpoint.go`)
- **Current:** Feature execution state persistence
- **Gaps:** No session-level tracking, no analytics
- **Needed:** Session-level checkpointing

**4. Beads Integration** (`internal/beads/`)
- **Current:** Task tracking client (Beads CLI installed but not configured)
- **Status:** Available but not used in this project
- **Opportunity:** Optional integration for task history

### New Components Required

**1. Memory Storage** (`internal/memory/`)
- Decision log storage (JSONL format, similar to telemetry)
- Session history storage
- Pattern extraction cache
- Lesson learned index

**2. Memory Analytics** (`internal/memory/analytics.go`)
- Pattern detection algorithms
- Trend analysis
- Correlation engine (decisions → outcomes)
- Lesson extraction

**3. Memory Search** (`internal/memory/search.go`)
- Full-text search over decisions
- Tag-based filtering
- Similarity matching (for "similar workstreams")
- Ranking by relevance

**4. Memory API** (`cmd/sdp/memory.go`)
- CLI commands for all user stories
- Export functionality
- Report generation

**5. Integration Points**
- `@feature` → Auto-log vision decisions
- `@design` → Auto-log architectural decisions
- `@build` → Log outcomes, extract lessons
- `@review` → Log quality patterns
- `@issue` → Correlate bugs with decisions

---

## Data Model

### Decision Log
```json
{
  "id": "D001",
  "timestamp": "2026-02-06T10:30:00Z",
  "type": "technical|vision|tradeoff|explicit",
  "question": "What database to use?",
  "decision": "Use PostgreSQL",
  "rationale": "ACID transactions required",
  "alternatives": ["MongoDB", "MySQL"],
  "outcome": "successful|failed|mixed|pending",
  "impact_score": 1-10,
  "workstream_id": "00-001-05",
  "feature_id": "F01",
  "decision_maker": "claude|user|team",
  "tags": ["database", "postgresql", "backend"],
  "related_decisions": ["D002", "D003"],
  "evidence": ["test results", "performance metrics"]
}
```

### Session Record
```json
{
  "id": "S001",
  "start_time": "2026-02-06T10:00:00Z",
  "end_time": "2026-02-06T12:00:00Z",
  "duration_seconds": 7200,
  "workstreams": ["00-001-05", "00-001-06"],
  "feature_id": "F01",
  "decisions_made": ["D001", "D002"],
  "commands_run": [
    {"command": "@build", "count": 2, "duration_seconds": 900}
  ],
  "outcome": "success|partial|failure",
  "lessons_learned": ["LL001", "LL002"]
}
```

### Lesson Learned
```json
{
  "id": "LL001",
  "type": "anti-pattern|success-pattern|optimization",
  "title": "Always index foreign keys",
  "description": "Missing foreign key indexes caused 10x slowdown",
  "source_workstreams": ["00-001-02", "00-001-05"],
  "occurrence_count": 2,
  "impact": "high|medium|low",
  "fix_recommendation": "Add db_index=True to all foreign keys",
  "tags": ["database", "performance", "sql"],
  "created_at": "2026-02-06T11:00:00Z"
}
```

### Pattern Extracted
```json
{
  "id": "P001",
  "pattern_type": "success|anti-pattern|decision",
  "description": "PostgreSQL chosen for transactional data",
  "confidence": 0.95,
  "supporting_instances": 12,
  "success_rate": 1.0,
  "contexts": ["user-auth", "payments", "orders"],
  "first_seen": "2026-01-15",
  "last_seen": "2026-02-05"
}
```

---

## Technical Requirements

### Non-Functional Requirements

**Performance:**
- Search response time: <500ms for 1000 decisions
- Analytics generation: <30s for 90 days of data
- Memory overhead: <100MB for 10K decisions
- Storage growth: <1MB per 100 decisions

**Reliability:**
- Data durability: Write-ahead log, fsync on critical writes
- No data loss: Graceful shutdown, crash recovery
- Backward compatibility: Support old decision formats

**Security:**
- Access control: Owner-only read/write (0600 permissions)
- No PII: Anonymize any user-specific data
- Audit log: Track who accessed/modified decisions

**Usability:**
- CLI ergonomics: Intuitive command structure
- Clear output: Human-readable with option for JSON
- Contextual help: `--help` on all commands

**Maintainability:**
- Code coverage: ≥80%
- Clean architecture: Separate storage, analytics, API layers
- Testable: Unit tests for all core logic

### Storage Requirements

**Location:**
```
.sdp/memory/
├── decisions.jsonl          # Decision log (JSONL format)
├── sessions.jsonl           # Session history
├── lessons.jsonl            # Extracted lessons
├── patterns.jsonl           # Detected patterns
├── analytics/               # Cached analytics
│   ├── daily.json
│   ├── weekly.json
│   └── monthly.json
└── index/                   # Search index
    ├── decisions.index
    └── lessons.index
```

**Format:**
- JSONL for append-only logs (decisions, sessions, lessons)
- JSON for analytics snapshots
- Custom binary format for search index (or SQLite)

**Retention:**
- Active period: 90 days in primary storage
- Archive: Compressed after 90 days
- Purge: Delete after 1 year (configurable)

**Backup:**
- Git integration: Commit to git automatically
- Export functionality: Manual backup on demand
- Sync: Optional cloud backup (user-provided)

---

## Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| **Low adoption** (developers don't log decisions) | High | Medium | Auto-logging from @feature, @design, @build; Make it frictionless |
| **False pattern detection** (wrong conclusions) | Medium | Medium | Confidence scoring, human review of patterns, minimal threshold |
| **Storage bloat** (unbounded growth) | Medium | Low | Automatic rotation, compression, archive old data |
| **Performance degradation** (slow analytics) | Medium | Low | Incremental updates, caching, lazy loading |
| **Privacy concerns** (sensitive decisions) | High | Low | Local-only storage, no remote sync, explicit opt-in |
| **Data corruption** (lost memory) | High | Low | Write-ahead log, git versioning, backup tools |

---

## Phasing & Rollout

### Phase 1: Foundation (Week 1-2)
**Goal:** Basic decision logging and search

**Deliverables:**
- Decision log storage (`decisions.jsonl`)
- CLI commands: `sdp memory log`, `sdp memory search`, `sdp memory list`
- Manual decision logging
- Basic search by keyword/tag

**Acceptance:**
- Can log a decision and search for it
- Storage persists across restarts
- No data loss on crashes

### Phase 2: Integration (Week 3-4)
**Goal:** Automatic logging from SDP workflows

**Deliverables:**
- Auto-logging from `@feature`, `@design`, `@build`, `@review`
- Session tracking (start/end time, workstreams, outcomes)
- CLI command: `sdp memory sessions`

**Acceptance:**
- Decisions auto-logged without manual intervention
- Sessions tracked automatically
- Can view session history

### Phase 3: Analytics (Week 5-6)
**Goal:** Pattern extraction and analytics

**Deliverables:**
- Pattern detection algorithms
- Lesson extraction from failures
- CLI commands: `sdp memory analyze`, `sdp memory lessons`
- Trend analysis (velocity, quality, decisions)

**Acceptance:**
- Patterns automatically detected and shown
- Lessons extracted from failed workstreams
- Analytics reports generated

### Phase 4: AI Integration (Week 7-8)
**Goal:** AI agents use memory for recommendations

**Deliverables:**
- Memory query API for AI agents
- Integration with Claude Code skills
- Context-aware recommendations

**Acceptance:**
- Claude references past decisions in responses
- Recommendations based on project history
- Lessons shown before similar workstreams

### Phase 5: Polish & Export (Week 9-10)
**Goal:** Reporting, dashboards, documentation

**Deliverables:**
- Export functionality (markdown, JSON, PDF)
- Interactive dashboards (optional web UI)
- Documentation and tutorials
- Performance optimization

**Acceptance:**
- Can export memory data for reports
- Performance targets met (<500ms search)
- Documentation complete

---

## Open Questions

1. **Beads Integration:** Should we integrate with Beads for task history?
   - **Pros:** Rich task context, automated correlation
   - **Cons:** Beads not configured in this project, adds dependency
   - **Decision:** Optional integration (Phase 6)

2. **Storage Backend:** JSONL vs SQLite vs custom binary?
   - **JSONL:** Simple, human-readable, git-friendly
   - **SQLite:** Fast queries, ACID, but adds dependency
   - **Decision:** JSONL for logs, SQLite for index (Phase 3)

3. **Privacy:** Should decisions be encrypted?
   - **Context:** May contain sensitive architectural decisions
   - **Decision:** No encryption (local-only), but add option in Phase 6

4. **AI Model:** Should we use ML for pattern detection?
   - **Context:** Can improve accuracy of pattern detection
   - **Decision:** Start with rule-based, add ML in Phase 7 (future)

5. **Multi-Project:** Should memory be per-project or global?
   - **Context:** Developer works on multiple projects
   - **Decision:** Per-project by default, optional global search

---

## Next Steps

1. **Review this document** with stakeholders (user, maintainers)
2. **Prioritize user stories** based on value and effort
3. **Create technical specification** for Phase 1
4. **Set up data model** and storage layer
5. **Implement Phase 1** (Foundation)
6. **Gather feedback** from early users
7. **Iterate** based on usage patterns

---

## Appendix

### Example Workflow

**Scenario:** Developer starting new workstream

```bash
# 1. Check for similar past work
$ sdp memory lessons --workstream "00-004-01" --similar

💡 Based on 3 similar workstreams:
🚫 Avoid: Missing database migrations (caused 2h debugging)
✅ Use: Repository pattern (100% success rate)

# 2. Search for relevant decisions
$ sdp memory search --query "database migration"

Found 2 decisions:
1. [D005] Use Alembic for migrations (2026-01-20)
   Outcome: Successful, 0 downtime deployments

# 3. Start workstream with context
$ @build 00-004-01

Claude: I see we've done similar database work before.
Based on D005, I'll use Alembic for migrations...
(implementation proceeds)

# 4. Workstream completes, outcome logged
$ sdp memory sessions --last

2026-02-06 (10:00-11:30) - 1h 30m
  Workstream: 00-004-01
  Outcome: ✅ Success
  Decisions: Used Alembic (see D008)
  Lessons: N/A (no issues)
```

### Glossary

- **Decision:** Architectural or technical choice made during development
- **Lesson:** Extracted insight from past work (success or failure)
- **Pattern:** Recurrent approach across multiple workstreams
- **Session:** Single development session (start → end time)
- **Workstream:** Atomic unit of work (SDP concept)
- **Feature:** Collection of workstreams (SDP concept)
- **Memory:** Long-term storage of decisions, sessions, lessons, patterns

---

**Document Status:** Draft for Review
**Author:** Business Analyst Agent (Claude Sonnet 4.5)
**Review Date:** 2026-02-06
**Version:** 1.0
