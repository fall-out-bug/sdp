# Workstream Breakdown: SDP Long-term Memory System

**Feature ID:** F01 (Memory System)
**Project ID:** 01
**Tech Lead:** Technical Decomposition Agent
**Date:** 2025-02-06

## Executive Summary

Breakdown of Long-term Memory feature for SDP, enabling decision tracking, session analytics, and project history. The feature leverages existing infrastructure (decision logging, telemetry tracking) and adds Git-backed storage, search capabilities, and AI-powered decision support.

**Total Effort:** 14 workstreams, 7.5-9 weeks
**Risk Level:** Medium (Git integration complexity, AI accuracy for decision similarity)

---

## Phase 1: Core Decision Tracking (MVP) - Week 1

### 01-001-01: Enhance Decision Data Model
**Size:** S (2-3 days)
**Dependencies:** None
**Priority:** P0 (Blocking)

**Acceptance Criteria:**
- AC1: Add `RelatedWorkstreams []string` field to Decision struct
- AC2: Add `SessionID string` field (links to telemetry session)
- AC3: Add `OutcomeStatus string` field ("success", "failure", "mixed", "pending")
- AC4: Add `RetrospectiveNotes string` field (post-implementation lessons)
- AC5: Backward compatible - existing decisions load without errors
- AC6: JSON serialization tests pass

**Implementation Notes:**
- Update `internal/decision/decision.go`
- Run `go generate` if using wire for dependency injection
- Add migration logic for existing decisions in `docs/decisions/decisions.jsonl`

**Files:**
- `internal/decision/decision.go` - Update Decision struct
- `internal/decision/decision_test.go` - Add tests for new fields

---

### 01-001-02: Git-backed Decision Storage
**Size:** M (3-5 days)
**Dependencies:** 01-001-01
**Priority:** P0 (Blocking)

**Acceptance Criteria:**
- AC1: Store decisions in `.beads/decisions/` directory (JSONL format, one file per feature)
- AC2: Auto-commit decisions to git with structured commit messages: `docs(decisions): Log {type} decision for {feature-id}`
- AC3: Sync with remote using `bd sync` pattern (reuses Beads sync logic)
- AC4: Import existing decisions from `docs/decisions/decisions.jsonl`
- AC5: Handle merge conflicts (last-write-wins with conflict markers)
- AC6: Unit tests for storage layer (mock git operations)

**Implementation Notes:**
- Create `internal/decision/git_storage.go`
- Reuse patterns from `internal/beads/` for git integration
- Add `--dry-run` flag for testing without actual git commits
- Use `go-git` library if available, else shell out to git binary

**Files:**
- `internal/decision/git_storage.go` - Git-backed storage implementation
- `internal/decision/git_storage_test.go` - Mock git tests
- `internal/decision/logger.go` - Update to use GitStorage
- `cmd/sdp/decisions.go` - Add `sdp decisions sync` command

---

### 01-001-03: Decision Search CLI (Enhanced)
**Size:** M (3-5 days)
**Dependencies:** 01-001-02
**Priority:** P0 (Blocking)

**Acceptance Criteria:**
- AC1: `sdp decisions search "query"` supports full-text search (question, decision, rationale)
- AC2: Filter flags: `--tags`, `--type`, `--status`, `--feature-id`, `--ws-id`
- AC3: Date range filtering: `--after "2025-01-01"`, `--before "2025-02-01"`
- AC4: Show related workstreams in search results
- AC5: Export results: `--format json|markdown|csv`
- AC6: Performance: search < 500ms for 1000 decisions

**Implementation Notes:**
- Current search is basic substring matching - upgrade to full-text search
- Consider `bleve` or `sqlite3 FTS5` for indexing
- Index should update on `git pull` (sync decisions)
- Add search query parser for complex queries: `type:technical status:failure "database"`

**Files:**
- `internal/decision/search.go` - Search engine implementation
- `internal/decision/search_test.go` - Search tests
- `internal/decision/index.go` - Full-text index management
- `cmd/sdp/decisions.go` - Enhance search command

---

### 01-001-04: Decision Metrics CLI
**Size:** S (2-3 days)
**Dependencies:** 01-001-03
**Priority:** P1 (High)

**Acceptance Criteria:**
- AC1: `sdp decisions metrics` shows decision statistics
- AC2: Group by type, feature, outcome status
- AC3: Show decision trends over time (decisions/week)
- AC4: Calculate "decision debt" (pending outcomes)
- AC5: Export metrics as JSON for dashboards

**Metrics to Track:**
- Total decisions by type
- Decisions per feature
- Outcome status distribution
- Average time to outcome
- Top decision tags

**Files:**
- `internal/decision/metrics.go` - Metrics calculation
- `internal/decision/metrics_test.go` - Metrics tests
- `cmd/sdp/decisions.go` - Add metrics command

---

## Phase 2: Session Analytics - Weeks 2-3

### 01-002-01: Session Tracking Enhancement
**Size:** M (4-5 days)
**Dependencies:** None (parallel with Phase 1)
**Priority:** P1 (High)

**Acceptance Criteria:**
- AC1: Track session start/end events in telemetry
- AC2: Assign unique SessionID to each SDP session
- AC3: Link workstreams to sessions (WS executed in Session X)
- AC4: Capture session context: project, user, time, commands run
- AC5: Session persistence: `~/.sdp/sessions.jsonl`

**Implementation Notes:**
- Extend `internal/telemetry/tracker.go` with session tracking
- Add `SessionStart` and `SessionEnd` event types
- SessionID format: `S-{timestamp}-{random}` (e.g., `S-20250206-abc123`)
- Auto-start session on first command, auto-end after 30min inactivity

**Files:**
- `internal/telemetry/types.go` - Add SessionEventType
- `internal/telemetry/session.go` - Session tracking logic
- `internal/telemetry/session_test.go` - Session tests

---

### 01-002-02: Session Statistics Analysis
**Size:** L (5-7 days)
**Dependencies:** 01-002-01
**Priority:** P1 (High)

**Acceptance Criteria:**
- AC1: `sdp sessions stats` shows session analytics
- AC2: Metrics: total sessions, avg duration, commands/session, success rate
- AC3: Time-series analysis: sessions/day, active hours heatmap
- AC4: Productivity score: (completed WS) / (session duration)
- AC5: Identify patterns: most productive times, common failure modes
- AC6: Export session data: `sdp sessions export --format json`

**Session Metrics:**
- Total sessions (all-time, this week, this month)
- Average session duration
- Commands per session
- Workstreams completed per session
- Session success rate (all commands succeeded?)
- Most active hours (0-23 heatmap)
- Productivity trends (7-day moving average)

**Files:**
- `internal/telemetry/analyzer.go` - Add session analysis
- `internal/telemetry/analyzer_test.go` - Session stats tests
- `cmd/sdp/sessions.go` - New sessions command

---

### 01-002-03: Weekly/Monthly Reports
**Size:** M (3-4 days)
**Dependencies:** 01-002-02
**Priority:** P1 (High)

**Acceptance Criteria:**
- AC1: `sdp sessions report --period week|month` generates summary
- AC2: Report includes: decisions made, WS completed, metrics summary
- AC3: Markdown output for documentation
- AC4: Email integration (optional, P2)
- AC5: Report templates in `templates/reports/`

**Report Sections:**
- Overview: Sessions, decisions, workstreams
- Top Decisions: Most impactful decisions this period
- Productivity: Velocity, success rate, trends
- Lessons Learned: Failed decisions, retrospective notes
- Next Period: Planned workstreams

**Files:**
- `internal/telemetry/report.go` - Report generation
- `internal/telemetry/report_test.go` - Report tests
- `templates/reports/weekly.md.tmpl` - Weekly report template
- `templates/reports/monthly.md.tmpl` - Monthly report template

---

## Phase 3: Decision Support - Weeks 4-5

### 01-003-01: Decision Similarity Engine
**Size:** L (5-7 days)
**Dependencies:** 01-001-03, 01-002-01
**Priority:** P2 (Medium)

**Acceptance Criteria:**
- AC1: Calculate similarity between decisions (embedding-based)
- AC2: `sdp decisions similar --to "WS-ID"` returns related decisions
- AC3: Similarity score > 0.7 threshold
- AC4: Use embeddings: OpenAI API or local model (sentence-transformers)
- AC5: Cache embeddings in `.beads/decisions/embeddings.jsonl`

**Algorithm:**
1. Extract decision text (question + decision + rationale)
2. Generate embedding (384-dim sentence embedding)
3. Cosine similarity search
4. Return top N similar decisions with scores

**Implementation Notes:**
- Start with keyword-based TF-IDF (simpler, no API needed)
- Upgrade to embeddings if TF-IDF insufficient
- Add `--embeddings-model` flag for model selection

**Files:**
- `internal/decision/similarity.go` - Similarity engine
- `internal/decision/similarity_test.go` - Similarity tests
- `internal/decision/embeddings.go` - Embedding generation

---

### 01-003-02: "We Already Tried" Detection
**Size:** XL (7-10 days)
**Dependencies:** 01-003-01
**Priority:** P2 (Medium)

**Acceptance Criteria:**
- AC1: AI agent checks decisions before starting workstream
- AC2: Warning message if similar decision exists with "failure" outcome
- AC3: Integrate with `@build` skill (pre-build check)
- AC4: Show related decisions: "You tried X in WS-001-02, it failed because Y"
- AC5: Configurable: opt-out via `--no-decision-check`

**Integration Points:**
- Hook into `pre-build.sh` or `.claude/skills/build/SKILL.md`
- Add `sdp decisions check --ws-id` command
- Prompt engineering: "⚠️ Decision found: You tried similar approach in WS-XXX"

**Files:**
- `hooks/pre-build-decision-check.sh` - Pre-build hook
- `internal/decision/checker.go` - Decision conflict detection
- `cmd/sdp/decisions.go` - Add check command

---

### 01-003-03: Decision Recommendations
**Size:** L (5-7 days)
**Dependencies:** 01-003-01
**Priority:** P2 (Medium)

**Acceptance Criteria:**
- AC1: `sdp decisions recommend --context "adding user auth"` suggests relevant decisions
- AC2: Recommend successful decisions (status=success)
- AC3: Show rationale and outcome for context
- AC4: Rank by relevance (similarity + outcome status)
- AC5: Include "what to avoid" (failed decisions)

**Use Cases:**
- Before starting feature: "What decisions should I review?"
- During planning: "What did we decide about X?"
- Onboarding: "Show me all database decisions"

**Files:**
- `internal/decision/recommend.go` - Recommendation engine
- `internal/decision/recommend_test.go` - Recommendation tests
- `cmd/sdp/decisions.go` - Add recommend command

---

## Phase 4: Project History & Lessons - Week 6

### 01-004-01: Project Timeline View
**Size:** M (3-4 days)
**Dependencies:** 01-001-02, 01-002-01
**Priority:** P1 (High)

**Acceptance Criteria:**
- AC1: `sdp history timeline` shows chronological project history
- AC2: Combine decisions, workstreams, sessions into timeline
- AC3: Group by week/month: "Week of Jan 6: 3 decisions, 2 WS completed"
- AC4: Filter by feature, type, status
- AC5: Export as markdown for changelog generation

**Timeline Format:**
```
2025-01-06 to 2025-01-12
├── Decisions: 3 (2 technical, 1 tradeoff)
├── Workstreams: 2 completed (00-001-01, 00-001-02)
└── Sessions: 5 (avg 2.3h each)
```

**Files:**
- `internal/history/timeline.go` - Timeline generation
- `internal/history/timeline_test.go` - Timeline tests
- `cmd/sdp/history.go` - New history command

---

### 01-004-02: Lessons Learned Capture
**Size:** M (3-4 days)
**Dependencies:** 01-001-01, 01-002-02
**Priority:** P1 (High)

**Acceptance Criteria:**
- AC1: `sdp decisions lessons` extracts insights from decisions
- AC2: Auto-generate lessons: "When X, prefer Y (learned from WS-Z)"
- AC3: Manual lesson capture: `sdp decisions add-lesson --ws-id "XXX" --lesson "Prefer PostgreSQL over SQLite for production"`
- AC4: Categorize lessons: performance, security, UX, architecture
- AC5: Store lessons in `.beads/lessons/lessons.jsonl`

**Lesson Format:**
```json
{
  "lesson_id": "L-20250106-001",
  "ws_id": "00-001-02",
  "category": "architecture",
  "lesson": "Use repository pattern for data access",
  "context": "Direct DB coupling caused testing issues",
  "impact": "positive",
  "timestamp": "2025-01-06T10:00:00Z"
}
```

**Files:**
- `internal/decision/lessons.go` - Lesson management
- `internal/decision/lessons_test.go` - Lesson tests
- `cmd/sdp/decisions.go` - Add lessons command

---

### 01-004-03: Project Dashboard Integration
**Size:** M (3-4 days)
**Dependencies:** 01-004-01, 01-004-02
**Priority:** P2 (Medium)

**Acceptance Criteria:**
- AC1: Add "Memory" tab to existing dashboard (`internal/ui/dashboard/`)
- AC2: Show: recent decisions, lessons learned, session stats
- AC3: Interactive search and filter
- AC4: Link decisions to workstreams

**Implementation Notes:**
- Extend existing TUI dashboard (Bubble Tea)
- Add memory view to tab navigation
- Reuse search logic from `internal/decision/search.go`

**Files:**
- `internal/ui/dashboard/memory_view.go` - Memory tab
- `internal/ui/dashboard/styles.go` - Add memory styles

---

### 01-004-04: Decision Retrospective Workflow
**Size:** S (2-3 days)
**Dependencies:** 01-001-01, 01-004-02
**Priority:** P2 (Medium)

**Acceptance Criteria:**
- AC1: `sdp decisions retrospective --ws-id "XXX"` opens retrospective editor
- AC2: Prompt for: outcome status, lessons learned, what worked/didn't
- AC3: Update decision with retrospective notes
- AC4: Auto-generate lesson if insight provided

**Retrospective Questions:**
1. What was the outcome? (success/failure/mixed)
2. What worked well?
3. What didn't work?
4. What would you do differently?
5. Lesson for future workstreams?

**Files:**
- `internal/decision/retrospective.go` - Retrospective workflow
- `internal/decision/retrospective_test.go` - Retrospective tests
- `cmd/sdp/decisions.go` - Add retrospective command

---

## Cross-Cutting Concerns

### Testing Strategy
- **Unit tests:** All packages, >80% coverage
- **Integration tests:** Git operations with test repo
- **E2E tests:** Full workflow (log → sync → search)
- **Performance tests:** 1000 decisions search < 500ms

### Documentation
- `docs/DECISIONS.md` - How decision logging works
- `docs/SESSIONS.md` - Session tracking guide
- `docs/WORKFLOW_DECISIONS.md` - Decision-first workflow
- `API.md` - Decision/Session API reference

### Migration Path
1. **Phase 1:** Existing decisions auto-imported to `.beads/decisions/`
2. **Phase 2:** Telemetry sessions auto-created from historical data
3. **Phase 3:** AI features opt-in (requires API key)

### Error Handling
- Git failures: Fallback to local storage, warn user
- Search failures: Degrade to substring search
- Embedding failures: Fallback to keyword search

---

## Risk Assessment

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Git merge conflicts on decisions | Medium | High | Auto-resolve with timestamps, manual review flag |
| Search performance degradation | High | Medium | Add pagination, indexing, caching |
| AI similarity accuracy | Low | High | User feedback loop, threshold tuning |
| Beads integration complexity | Medium | Low | Reuse existing sync logic, add feature flag |

---

## Effort Summary

**Phase 1 (MVP):** 1 S + 2 M + 1 S = ~1 week
**Phase 2:** 1 M + 1 L + 1 M = ~2 weeks
**Phase 3:** 1 L + 1 XL + 1 L = ~3 weeks
**Phase 4:** 3 M + 1 S = ~1.5 weeks

**Total:** 14 workstreams, **7.5 weeks**
**Buffer (+20%):** ~9 weeks total

---

## Dependencies on Other Features

None - this is a standalone feature that enhances existing SDP workflows.

---

## Beads Integration Status

**Status:** SKIPPED - Beads is not enabled in this project (`.beads/` directory does not exist)

**Note:** Beads CLI is installed (`bd version 0.49.3`) but the project has not been initialized with Beads task tracking. To enable Beads integration:

```bash
# Initialize Beads in project
bd init

# Then create feature and workstream tasks
bd create --title="Feature F01: SDP Long-term Memory" --type=feature
bd create --title="WS 01-001-01: Enhance Decision Data Model" --type=task --parent=F01
# ... (repeat for all 14 workstreams)
```

---

## Next Steps

1. **Product Manager approval** of workstream breakdown
2. **Start with Phase 1** (01-001-01, 01-001-02, 01-001-03)
3. **Parallel work** on 01-002-01 (session tracking) while Phase 1 completes
4. **Enable Beads** if task tracking is desired

---

**Tech Lead:** Technical Decomposition Agent
**Date:** 2025-02-06
**Status:** Ready for Business Review
