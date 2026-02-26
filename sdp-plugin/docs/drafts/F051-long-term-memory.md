# SDP Long-term Memory System

**Vision:** Transform SDP from development tool into organization's collective brain - capturing not just what was built, but **why decisions were made**, what patterns emerge over time, and what knowledge can be reused.

**Status:** Requirements gathering complete
**Created:** 2026-02-06
**Priority:** P1 (High-value feature)

---

## Executive Summary

SDP currently lacks institutional memory. Developers:
- Repeat mistakes (no "we already tried X")
- Forget past decisions (no "we decided Y")
- Lose lessons learned (no failure patterns)
- Can't analyze productivity (no session metrics)

This feature implements comprehensive long-term memory:
1. **Decision Tracking** - Searchable "we decided" log
2. **Session Analytics** - Development activity patterns
3. **Project History** - Failures, successes, lessons
4. **Decision Support** - "We tried X before" warnings

**Success Definition:** Developer asks "Have we faced this before?" and gets answer in <10 seconds with context, rationale, outcomes.

---

## User Stories (English)

### US-001: Decision Search
**As:** Developer working on a feature
**I want:** Search previous decisions by keywords
**So that:** Don't repeat past mistakes

**Scenario:**
- **Given** Developer faces an authentication problem
- **When** They enter `sdp memory search "authentication"`
- **Then** System shows:
  - Decision F01-2024-01-15: "Use JWT" (Worked: success)
  - Decision F02-2024-03-20: "Use sessions" (Failed: session management complexity)
  - Pattern: "JWT decisions work 80% of the time"

**Value:** "We already tried approach X in {workstream}, abandoned due to {reason}"

---

### US-002: Decision Logging
**As:** SDP user
**I want:** Decisions to be automatically saved
**So that:** Don't forget documenting why

**Triggers:**
- `@feature "Add payments"` -> Log vision decisions
- `@design idea-payments` -> Log technical decisions
- `@build 00-001-01` -> Log implementation decisions
- `@review F01` -> Log review findings

**Auto-captured:**
- Timestamp
- Feature/Workstream ID
- Decision maker (user vs agent)
- Related files
- Tags (from content)

**Value:** Complete decision history without effort

---

### US-003: Session Analytics
**As:** Developer
**I want:** See statistics of my work
**So that:** Optimize productivity

**Metrics:**
- Commands per session
- Session duration
- Completion rate (started vs finished)
- Most productive hours
- Agent performance (which agents fail most)

**CLI:**
```bash
sdp memory stats --last 30d
# Output:
# Sessions: 45 (avg 2.3h)
# Completion: 78%
# Peak hours: 10am-12pm
# Top agent: @build (32%)
```

**Value:** Data for process improvement

---

### US-004: "We Already Tried" (Decision Warnings)
**As:** Developer
**I want:** SDP to warn about past failures
**So that:** Don't make the same mistakes again

**Scenario:**
- **Given** Developer proposes "Use MongoDB"
- **When** SDP finds Decision F03-2024-05-10: "MongoDB failed - scaling issues"
- **Then** SDP shows warning:
  ```
  Similar decision found:
  F03-2024-05-10: "Use MongoDB"
  Outcome: Failed - scaling issues at 10M users
  Rationale: Chose for flexibility, regret due to transactions
  Continue anyway? (y/n)
  ```

**Value:** Preventing repeated mistakes

---

### US-005: Lessons Learned
**As:** Development team
**I want:** Automatically extract lessons from completed workstreams
**So that:** Don't lose knowledge

**Auto-extraction:**
- Failed workstreams -> "What went wrong?"
- Successful workstreams -> "What worked?"
- Reversed decisions -> "Why did we change our mind?"

**CLI:**
```bash
sdp memory lessons --feature F01
# Output:
# Lessons from F01 (Payment Processing):
# Worked: Use Stripe SDK (saved 2 weeks)
# Failed: Custom checkout (abandoned after 3 days)
# Risk: Webhook reliability (add retries)
```

**Value:** Systematic learning

---

### US-006: Project Timeline
**As:** New developer on the team
**I want:** See project development chronology
**So that:** Onboard quickly

**Timeline view:**
```
Jan 15: F01 decided "Use JWT" (technical)
Jan 20: WS-001-01 completed (user auth)
Feb 01: F01 decision reversed "JWT -> sessions" (regret: complexity)
Feb 10: WS-001-02 completed (session management)
```

**Features:**
- Filter by feature, type, agent
- View decisions + sessions + commits
- Export to markdown (for docs)

**Value:** Quick onboarding (1 day instead of 1 week)

---

### US-007: Patterns and Metrics
**As:** Tech Lead
**I want:** See decision-making patterns
**So that:** Make better decisions

**Pattern examples:**
- "Security decisions take 3x longer to reach"
- "Tech stack changes cause 60% of delays"
- "Features with 20+ decisions have 50% bug rate"

**Metrics:**
- Decision velocity (decisions/week)
- Outcome tracking (worked vs failed)
- Tag frequency (what we discuss most)

**Value:** Data-driven improvements

---

### US-008: Export and Reports
**As:** Project Manager
**I want:** Export project data
**So that:** Show to stakeholders

**Export formats:**
- Markdown (for documentation)
- JSON (for analysis)
- CSV (for spreadsheets)

**Report types:**
- Monthly decision summary
- Quarterly lessons learned
- Project health report

**Value:** Business transparency

---

## Functional Requirements

### FR-001: Enhanced Decision Data Model
**Extend existing Decision struct:**
```go
type Decision struct {
    // Existing fields
    Timestamp     time.Time
    Type          string // vision, technical, tradeoff, explicit
    FeatureID     string
    WorkstreamID  string
    Question      string
    Decision      string
    Rationale     string
    Alternatives  []string
    Outcome       string
    DecisionMaker string
    Tags          []string

    // New fields
    SessionID         string    `json:"session_id"`
    RelatedDecisions  []string  `json:"related_decisions"`
    ReversesDecision  string    `json:"reverses_decision"` // ID if this reverses past decision
    ActualOutcome     string    `json:"actual_outcome"`    // worked/failed/mixed
    OutcomeCapturedAt time.Time `json:"outcome_captured_at"`
    Confidence        float64   `json:"confidence"`        // 0-1
}
```

**Acceptance:**
- Backward compatible with existing decisions
- Supports decision reversal tracking
- Links to sessions

---

### FR-002: Session Management
**New Session struct:**
```go
type Session struct {
    ID          string    `json:"id"`           // UUID
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Agents      []string  `json:"agents"`       // Which agents used
    Skills      []string  `json:"skills"`       // Which skills invoked
    Files       []string  `json:"files"`        // Files touched
    Commits     []string  `json:"commits"`      // Git commits
    Completed   int       `json:"completed"`    // Tasks completed
    Total       int       `json:"total"`        // Total tasks
}
```

**Acceptance:**
- Auto-start on first command
- Auto-end on `bd sync` or timeout
- Resume after interruption
- Link to decisions made

---

### FR-003: Search & Query
**Full-text search:**
```bash
sdp memory search "authentication" --type technical --last 30d
# Returns ranked results with relevance score
```

**Filters:**
- By type (vision/technical/tradeoff)
- By feature/workstream
- By date range
- By tags
- By outcome (worked/failed)

**CLI commands:**
- `sdp memory search <query>` - Search decisions
- `sdp memory related <id>` - Find related decisions
- `sdp memory timeline` - Project timeline
- `sdp memory patterns` - Pattern analysis

**Acceptance:**
- Search <1sec for 1000 decisions
- Support complex queries
- Export results

---

### FR-004: Analytics & Metrics
**Session analytics:**
```bash
sdp memory stats --sessions
# Output:
# Sessions: 45 (avg 2.3h)
# Completion: 78%
# Peak hours: 10am-12pm
# Most used: @build (32%), @review (18%)
```

**Decision analytics:**
```bash
sdp memory stats --decisions
# Output:
# Total: 127 decisions
# By type: technical (45%), tradeoff (30%)
# With outcomes: 62%
# Success rate: 78%
```

**Pattern detection:**
```bash
sdp memory patterns
# Output:
# Pattern: Security decisions take 3x longer (avg 45min)
# Pattern: Tech stack changes cause 60% of delays
# Pattern: Features with 20+ decisions -> 50% bug rate
```

**Acceptance:**
- Compute statistics in <5sec
- Identify actionable patterns
- Generate reports

---

### FR-005: Storage & Sync
**Git-backed JSONL storage:**
```
.sdp-memory/
├── decisions.jsonl    # Decision log
├── sessions.jsonl     # Session history
├── lessons.jsonl      # Extracted lessons
├── patterns.jsonl     # Detected patterns
└── analytics/         # Cached statistics
```

**Auto-commit:**
- Opt-in git integration
- Commit on decision/session creation
- Sync with remote (bd sync pattern)

**Acceptance:**
- Version-controlled memory
- No database dependency
- Privacy-first (local-only)

---

## Technical Architecture

### Storage Strategy
**Decision:** Git-backed JSONL files (like Beads)

**Rationale:**
- Version-controlled (review decisions in PRs)
- No database server (simple deployment)
- Privacy-first (local-only, no cloud)
- Code-reviewed (decisions alongside code)
- Scalable to 10K decisions (~2MB)

**Beyond 10K:** Add SQLite FTS5 index

---

### Search Strategy
**Phase 1 (MVP):** Linear scan with `strings.Contains`
- Works for 1000 decisions
- Simple, no dependencies

**Phase 2:** SQLite FTS5 full-text search
- <100ms for 10K decisions
- Ranked results
- Complex queries

**Phase 3:** Embeddings + vector search
- Semantic similarity
- "Decisions like this one"
- Requires ML model

---

### Session Tracking
**Integration points:**
- `@feature`, `@design` -> Create session
- `@build`, `@review` -> Update session
- `bd sync` -> End session
- Crash/interruption -> Resume on restart

**Data capture:**
- Commands executed
- Agents/skills used
- Files modified
- Time spent
- Completion rate

---

## Workstream Breakdown

### Phase 1: Foundation (MVP) - 2 weeks

**01-001-01: Enhanced Decision Data Model** (S, 2-3 days)
- AC1: Add SessionID, RelatedDecisions, ReversesDecision to Decision struct
- AC2: Add ActualOutcome, OutcomeCapturedAt, Confidence fields
- AC3: Backward compatible with existing decisions
- AC4: Unit tests for new fields

**01-001-02: Git-backed Decision Storage** (M, 3-5 days)
- AC1: Create `.sdp-memory/decisions/` directory
- AC2: Store decisions in JSONL format
- AC3: Auto-commit to git (opt-in)
- AC4: Sync with remote (bd sync pattern)
- AC5: Migrate existing decisions

**01-001-03: Decision Search CLI** (M, 3-5 days)
- AC1: `sdp memory search <query>` command
- AC2: Filter by type, date, tags, outcome
- AC3: Ranked results (relevance score)
- AC4: Export to JSON/markdown
- AC5: Search <1sec for 1000 decisions

**01-001-04: Decision Metrics CLI** (S, 2-3 days)
- AC1: `sdp memory stats --decisions` command
- AC2: Show decision count by type
- AC3: Show outcome success rate
- AC4: Tag frequency analysis

---

### Phase 2: Session Analytics - 3 weeks

**01-002-01: Session Management** (M, 4-5 days)
- AC1: Implement Session struct
- AC2: Auto-start on first command
- AC3: Auto-end on bd sync/timeout
- AC4: Resume after interruption
- AC5: Link sessions to decisions

**01-002-02: Session Statistics** (L, 5-7 days)
- AC1: `sdp memory stats --sessions` command
- AC2: Compute session duration, completion rate
- AC3: Agent/skill usage analysis
- AC4: Peak productivity hours detection
- AC5: Time-based trends

**01-002-03: Session Reports** (M, 3-4 days)
- AC1: Weekly/monthly session summaries
- AC2: Productivity metrics (velocity trend)
- AC3: Export to PDF/markdown
- AC4: Email reports (opt-in)

---

### Phase 3: Decision Support - 2 weeks

**01-003-01: Decision Similarity** (L, 5-7 days)
- AC1: Find similar decisions (TF-IDF)
- AC2: `sdp memory similar <id>` command
- AC3: Suggest related decisions
- AC4: Show "decisions like this"

**01-003-02: "We Already Tried" Detection** (XL, 7-10 days)
- AC1: Detect similar decisions in real-time
- AC2: Warning before repeating past mistakes
- AC3: Show past outcomes
- AC4: Integration with @feature, @design

**01-003-03: Decision Recommendations** (L, 5-7 days)
- AC1: Suggest decisions based on history
- AC2: "When you chose X before, it failed"
- AC3: Confidence scores
- AC4: Learn from outcomes

---

### Phase 4: Project History - 2 weeks

**01-004-01: Project Timeline** (M, 3-4 days)
- AC1: `sdp memory timeline` command
- AC2: Unified chronological view
- AC3: Filter by feature, type
- AC4: Export timeline

**01-004-02: Lessons Learned** (M, 3-4 days)
- AC1: Auto-extract lessons from completed workstreams
- AC2: Categorize (worked/failed/risk)
- AC3: `sdp memory lessons` command
- AC4: Link to source decisions

**01-004-03: Project Dashboard** (M, 3-4 days)
- AC1: Web dashboard (existing ui/dashboard)
- AC2: Visual timeline
- AC3: Decision graph
- AC4: Real-time metrics

**01-004-04: Retrospective Workflow** (S, 2-3 days)
- AC1: `sdp memory retrospective <feature>` command
- AC2: Generate retrospective document
- AC3: Include decisions, lessons, metrics
- AC4: Export to markdown

---

## Success Metrics

### Adoption Metrics
- **Decision logging rate:** >=80% of workstreams have logged decisions
- **Memory search usage:** >=5 searches/day per active project
- **Session history views:** >=10 views/week

### Quality Metrics
- **Mistake recurrence:** Down 40% from baseline (measured via surveys)
- **Decision reuse:** >=50% of new decisions reference past decisions
- **Pattern detection accuracy:** >=80% precision (user feedback)

### Time Savings
- **Time saved per workstream:** >=5min (no searching Slack/email)
- **Decision-making time:** Down 50% (past decisions available)
- **Onboarding time:** Down 30% (timeline view)

### ROI
- **Time saved:** >=10 hours/developer/month
- **Bug reduction:** Down 30% (pattern detection)
- **Development velocity:** Up 20% (fewer mistakes)

---

## Implementation Priority

### P0 (MVP - Phase 1)
- Enhanced Decision data model
- Git-backed storage
- Basic search CLI
- Decision metrics

**Effort:** 2 weeks
**Value:** Core foundation

### P1 (High Value - Phase 2-3)
- Session management
- Session analytics
- Decision similarity
- "We already tried" warnings

**Effort:** 5 weeks
**Value:** Major time savings

### P2 (Nice to Have - Phase 4)
- Project timeline
- Lessons learned
- Dashboard
- Retrospective workflow

**Effort:** 2 weeks
**Value:** Polish & reporting

---

## Open Questions

1. **Beads Integration:**
   - Should decisions create Beads tasks automatically?
   - Should decisions link to Beads tasks?
   - **Recommendation:** Optional, Phase 6

2. **Search Algorithm:**
   - MVP: Simple grep (works for 1000 decisions)
   - Phase 2: SQLite FTS5 (better relevance)
   - Phase 3: Embeddings (semantic search)
   - **Recommendation:** Start simple, upgrade based on usage

3. **Outcome Prompting:**
   - When to ask for outcomes? (Immediate? 30 days? 60 days?)
   - **Recommendation:** Prompt 30 days after decision, email reminder

4. **Privacy:**
   - Decisions may contain sensitive info
   - **Recommendation:** Local-only, opt-in cloud sync, respect .gitignore

5. **Cross-Project Search:**
   - Search decisions across multiple SDP projects?
   - **Recommendation:** Per-project with optional global index (Phase 4)

---

## Dependencies

**Existing:**
- Decision struct (`internal/decision/decision.go`)
- Decision logger (`internal/decision/logger.go`)
- CLI commands (`cmd/sdp/decisions.go`)
- Telemetry system (`internal/telemetry/`)
- Beads pattern (git-backed storage)

**New Required:**
- Memory storage (`internal/memory/`)
- Session tracking (`internal/memory/session.go`)
- Search engine (`internal/memory/search.go`)
- Analytics engine (`internal/memory/analytics.go`)
- CLI commands (`cmd/sdp/memory.go`)

---

## Risk Assessment

### Risk 1: Low Adoption
**Mitigation:**
- Make logging automatic (no user effort)
- Integrate into existing workflows (@feature, @design)
- Show immediate value (search before decision)

### Risk 2: Data Quality
**Mitigation:**
- Capture metadata automatically (timestamp, feature, agent)
- Prompt for outcomes (30 days later)
- Validate decision structure (required fields)

### Risk 3: Search Performance
**Mitigation:**
- Start with simple grep (works for 1000 decisions)
- Upgrade to SQLite FTS5 at 10K decisions
- Index in background

### Risk 4: Privacy Concerns
**Mitigation:**
- Local-only by default
- Opt-in cloud sync
- Respect .gitignore
- Anonymize for analytics

---

## Next Steps

1. **Review this specification** with stakeholders
2. **Prioritize phases** based on team needs
3. **Start Phase 1** (Foundation - 2 weeks)
4. **Measure adoption** after Phase 1
5. **Iterate** based on usage patterns

---

## Appendix: Existing Code Inventory

**Already Built:**
- `internal/decision/decision.go` - Decision struct (basic fields)
- `internal/decision/logger.go` - JSONL logging
- `cmd/sdp/decisions.go` - CLI commands (list, search, export)
- `internal/telemetry/` - Command tracking
- `docs/runbooks/decision-logging.md` - Operational docs

**Needs Enhancement:**
- Decision struct (add session link, outcome tracking)
- Search (full-text index)
- Session tracking (new component)
- Analytics (new component)
- "We already tried" warnings (new feature)

**Quick Wins (1 day each):**
- Add pagination to `decisions list`
- Add `--format csv` to `decisions export`
- Add `--outcome` filter to `decisions search`
- Add `decisions stats` command

---

**Document Status:** Requirements complete
**Next:** Technical design -> Implementation
**Estimated Timeline:** 9 weeks (all 4 phases)
**Recommended Start:** Phase 1 (MVP)
