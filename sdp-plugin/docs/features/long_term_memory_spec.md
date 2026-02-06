# Long-Term Memory System for SDP
## Functional Specification

**Version:** 1.0
**Date:** 2026-02-06
**Status:** Draft
**Author:** Systems Analyst Agent

---

## Executive Summary

This specification defines a comprehensive long-term memory system for SDP that builds upon existing telemetry and decision tracking infrastructure. The system will provide institutional memory across three domains: **decision tracking**, **session analytics**, and **project history**.

### Key Design Principles

1. **Build on Existing Infrastructure** - Extend `internal/decision/` and `internal/telemetry/` systems
2. **Git-Backed Storage** - All data version-controlled alongside code (like Beads)
3. **Searchable Memory** - Full-text search across decisions, sessions, and events
4. **Privacy-First** - No PII, local-only storage, opt-in telemetry
5. **AI-Ready** - Structured data for Claude Code context retrieval

---

## Functional Requirements

### FR-001: Decision Capture Enhancement

**Priority:** P0 (Existing system - requires extension)

The system MUST extend the existing `Decision` struct to support:

1. **Session Linking** - Associate decisions with Claude Code sessions
2. **Outcome Tracking** - Post-decision validation (did it work?)
3. **Related Decisions** - Link decisions that influence each other
4. **Decision Reversal** - Track when decisions are reversed and why

**Rationale:** Current `Decision` struct (in `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/decision/decision.go`) captures the "what" and "why" but lacks linkage to sessions and outcomes.

---

### FR-002: Session Management

**Priority:** P0 (New feature)

The system MUST introduce a `Session` abstraction that:

1. **Session Identification** - Unique session IDs (UUID v4)
2. **Session Lifecycle** - Track session start, end, duration
3. **Session Context** - Feature/workstream scope, agent roles involved
4. **Session Artifacts** - Code changes, decisions, telemetry events linked to session
5. **Session Handoff** - Context for resuming interrupted work

**Rationale:** Currently no session tracking exists. Sessions bridge the gap between discrete events and cohesive work units.

---

### FR-003: Session Analytics

**Priority:** P1 (New feature)

The system MUST provide analytics on session patterns:

1. **Productivity Metrics** - Sessions per day, average duration, completion rate
2. **Agent Performance** - Which agents/skills are used most frequently
3. **Time Analysis** - Peak development hours, session length distribution
4. **Success Patterns** - Correlate session outcomes with decisions, telemetry
5. **Trend Detection** - Identify changes in productivity over time

**Rationale:** Users want to understand their development patterns ("When am I most productive?", "Which workflows work best?").

---

### FR-004: Project History Timeline

**Priority:** P1 (New feature)

The system MUST maintain a searchable project timeline:

1. **Chronological Events** - Decisions, sessions, telemetry in temporal order
2. **Feature Evolution** - Track features from conception to deployment
3. **Workstream History** - Complete audit trail of all workstreams
4. **Milestones** - Automatic detection of significant events (first commit, deployment, etc.)
5. **Time Travel** - View project state at any point in history

**Rationale:** New team members need onboarding ("How did we arrive at this architecture?"). Existing developers need context ("Why did we choose this approach?").

---

### FR-005: Search and Query Interface

**Priority:** P0 (Core capability)

The system MUST provide powerful search:

1. **Full-Text Search** - Search across all memory types (decisions, sessions, events)
2. **Faceted Filters** - Filter by date range, feature, agent, decision type
3. **Semantic Search** - Find similar decisions ("authentication approach" → JWT decision)
4. **Time-Based Queries** - "What happened last week?", "Decisions since deployment"
5. **CLI Integration** - `sdp memory search`, `sdp memory timeline`, `sdp memory decisions`

**Rationale:** Without effective search, institutional memory is inaccessible.

---

### FR-006: Memory Export and Sharing

**Priority:** P2 (Nice-to-have)

The system SHOULD support exporting memory:

1. **Markdown Reports** - Human-readable decision logs, session summaries
2. **JSON Export** - Machine-readable format for analysis
3. **Git Integration** - Auto-commit memory to repo (configurable)
4. **Team Sharing** - Attach memory to PRs, issues, documentation

**Rationale:** Teams need to share decisions with stakeholders who don't use SDP CLI.

---

## Data Models

### 1. Extended Decision Struct

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/decision/decision.go`

```go
// Existing fields (preserved)
type Decision struct {
    Timestamp     time.Time   `json:"timestamp"`
    Type          string      `json:"type"` // "vision", "technical", "tradeoff", "explicit"
    FeatureID     string      `json:"feature_id"`
    WorkstreamID  string      `json:"ws_id,omitempty"`
    Question      string      `json:"question"`
    Decision      string      `json:"decision"`
    Rationale     string      `json:"rationale"`
    Alternatives  []string    `json:"alternatives,omitempty"`
    Outcome       string      `json:"outcome"` // Renamed from ExpectedOutcome
    DecisionMaker string      `json:"decision_maker"`
    Tags          []string    `json:"tags,omitempty"`

    // NEW: Session linkage
    SessionID     string      `json:"session_id,omitempty"` // Link to session

    // NEW: Decision relationships
    RelatedDecisions []string `json:"related_decisions,omitempty"` // IDs of related decisions
    ReversesDecision string  `json:"reverses_decision,omitempty"` // ID if this reverses a previous decision
    ReversedBy       []string `json:"reversed_by,omitempty"`       // IDs of decisions that reversed this one

    // NEW: Outcome validation
    ActualOutcome    string    `json:"actual_outcome,omitempty"`    // What actually happened
    OutcomeValidated bool      `json:"outcome_validated"`          // Has outcome been verified?
    OutcomeValidatedAt time.Time `json:"outcome_validated_at,omitempty"` // When outcome was confirmed

    // NEW: Decision status
    Status        string     `json:"status"` // "proposed", "approved", "implemented", "validated", "reversed"
    ApprovedBy    []string   `json:"approved_by,omitempty"` // Who approved (user IDs, agent IDs)
    ImplementedAt *time.Time `json:"implemented_at,omitempty"`
}
```

**Changes:**
- Added `SessionID` to link decisions to sessions
- Added decision relationship graph (`RelatedDecisions`, `ReversesDecision`, `ReversedBy`)
- Added outcome validation fields (`ActualOutcome`, `OutcomeValidated`)
- Added decision lifecycle status (`Status`, `ApprovedBy`, `ImplementedAt`)

---

### 2. Session Struct (New)

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/session/session.go` (new)

```go
package session

import "time"

// Session represents a Claude Code work session
type Session struct {
    // Identity
    ID        string    `json:"id"`        // UUID v4
    StartedAt time.Time `json:"started_at"`
    EndedAt   *time.Time `json:"ended_at,omitempty"`

    // Context
    ProjectID   string   `json:"project_id"`   // SDP project ID (e.g., "00")
    FeatureID   string   `json:"feature_id,omitempty"`   // Feature being developed
    WorkstreamIDs []string `json:"workstream_ids,omitempty"` // Workstreams in this session

    // Participants
    Agents      []string `json:"agents"`       // Agent roles involved (e.g., ["planner", "builder"])
    SkillsUsed  []string `json:"skills_used"`  // Skills invoked (e.g., ["@build", "@review"])

    // Artifacts (references, not full data)
    Decisions      []string `json:"decisions,omitempty"`      // Decision IDs made in session
    Events         []string `json:"events,omitempty"`         // Telemetry event IDs
    FilesModified  []string `json:"files_modified,omitempty"` // File paths changed
    Commits        []string `json:"commits,omitempty"`        // Git commit SHAs

    // Outcome
    Status         string  `json:"status"` // "active", "completed", "interrupted", "failed"
    CompletionRate float64 `json:"completion_rate"` // 0.0 to 1.0 (percent of planned work completed)

    // Human context
    Summary       string `json:"summary"`       // Human-readable summary
    NextSteps     string `json:"next_steps,omitempty"` // Context for resuming
    InterruptedReason string `json:"interrupted_reason,omitempty"` // Why session ended prematurely

    // Metadata
    Tags          []string `json:"tags,omitempty"`
}
```

---

### 3. SessionEvent Struct (New)

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/session/event.go` (new)

```go
package session

import "time"

// SessionEvent represents a significant event within a session
type SessionEvent struct {
    ID        string    `json:"id"`        // UUID v4
    SessionID string    `json:"session_id"`
    Timestamp time.Time `json:"timestamp"`

    // Event type
    Type      string `json:"type"` // "skill_invoked", "decision_made", "file_created", "test_run", etc.
    Category  string `json:"category"` // "agent", "user", "system"

    // Event data
    Data      map[string]interface{} `json:"data"`

    // Relationships
    CorrelationID string `json:"correlation_id,omitempty"` // Link related events (e.g., request/response)
}
```

---

### 4. ProjectTimeline Struct (New)

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/timeline/timeline.go` (new)

```go
package timeline

import "time"

// TimelineEvent represents any event in project history
type TimelineEvent struct {
    ID        string    `json:"id"`
    Timestamp time.Time `json:"timestamp"`

    // Event type
    Type      string `json:"type"` // "decision", "session_start", "session_end", "commit", "deployment", etc.

    // References
    DecisionID  string `json:"decision_id,omitempty"`
    SessionID   string `json:"session_id,omitempty"`
    FeatureID   string `json:"feature_id,omitempty"`
    WorkstreamID string `json:"workstream_id,omitempty"`

    // Event summary
    Summary     string `json:"summary"`
    Details     string `json:"details,omitempty"`

    // Importance
    Significance string `json:"significance"` // "low", "medium", "high", "milestone"
}

// Timeline is a chronological sequence of events
type Timeline struct {
    Events      []TimelineEvent `json:"events"`
    ProjectID   string          `json:"project_id"`
    StartDate   time.Time       `json:"start_date"`
    EndDate     *time.Time      `json:"end_date,omitempty"`
}
```

---

## Storage Design

### Architecture Decision: Git-Backed JSONL Files

**Recommendation:** Store all memory data in JSONL files, version-controlled with Git.

**Rationale:**
1. **Consistency with Beads** - Beads uses Git-backed storage (`.beads/` directory)
2. **Consistency with Existing** - Telemetry uses `~/.sdp/telemetry.jsonl`, decisions use `docs/decisions/decisions.jsonl`
3. **Built-in Versioning** - Git provides history, branching, merging
4. **Simplicity** - No database server, no migration scripts
5. **Code-Reviewed Memory** - Decisions reviewed alongside code in PRs
6. **Privacy** - All data local, no external services

---

### Directory Structure

```
project-root/
├── .sdp-memory/                    # Git-tracked memory directory
│   ├── sessions/                   # Session data
│   │   ├── sessions.jsonl          # Append-only session log
│   │   ├── 2026-02/                # Monthly rotation (optional)
│   │   │   ├── sessions.jsonl
│   │   │   └── events.jsonl
│   │   └── index.json              # Session index (for fast lookup)
│   │
│   ├── decisions/                  # Extended decision tracking
│   │   ├── decisions.jsonl         # Main log (existing, will extend)
│   │   ├── decisions.json          # Generated export (human-readable)
│   │   └── index.json              # Decision index
│   │
│   ├── timeline/                   # Project timeline
│   │   ├── timeline.jsonl          # Chronological event log
│   │   └── milestones.json         # Detected milestones
│   │
│   └── analytics/                  # Computed analytics
│       ├── session_stats.json      # Aggregated statistics
│       ├── productivity_metrics.json
│       └── decision_patterns.json
│
├── docs/decisions/                 # Existing decision location (will be migrated)
│   └── decisions.jsonl             # Legacy support (symlink to .sdp-memory/decisions/)
│
└── .gitignore                      # Ensure .sdp-memory/ is tracked
```

**Migration Strategy:**
- Existing `docs/decisions/decisions.jsonl` will be symlinked to `.sdp-memory/decisions/decisions.jsonl`
- Telemetry (`~/.sdp/telemetry.jsonl`) remains user-local (not project-specific)
- New `.sdp-memory/` directory committed to Git (configurable via `.sdp/memory.gitignore`)

---

### File Formats

#### 1. Sessions Log (JSONL)

**File:** `.sdp-memory/sessions/sessions.jsonl`

```jsonl
{"id":"550e8400-e29b-41d4-a716-446655440000","started_at":"2026-02-06T10:00:00Z","ended_at":"2026-02-06T12:30:00Z","project_id":"00","feature_id":"F01","workstream_ids":["00-001-01","00-001-02"],"agents":["planner","builder"],"skills_used":["@design","@build"],"decisions":["d1","d2"],"events":["e1","e2","e3"],"files_modified":["src/auth.go","src/auth_test.go"],"commits":["abc123","def456"],"status":"completed","completion_rate":1.0,"summary":"Implemented JWT authentication for F01","next_steps":"Deploy to staging","tags":["authentication","jwt"]}
```

#### 2. Session Events Log (JSONL)

**File:** `.sdp-memory/sessions/2026-02/events.jsonl`

```jsonl
{"id":"evt-1","session_id":"550e8400-e29b-41d4-a716-446655440000","timestamp":"2026-02-06T10:05:00Z","type":"skill_invoked","category":"agent","data":{"skill":"@design","feature":"F01"},"correlation_id":"corr-1"}
{"id":"evt-2","session_id":"550e8400-e29b-41d4-a716-446655440000","timestamp":"2026-02-06T10:15:00Z","type":"decision_made","category":"user","data":{"decision_id":"d1","question":"JWT vs sessions?","decision":"JWT"},"correlation_id":"corr-1"}
{"id":"evt-3","session_id":"550e8400-e29b-41d4-a716-446655440000","timestamp":"2026-02-06T11:00:00Z","type":"workstream_started","category":"system","data":{"workstream_id":"00-001-01","title":"Domain entities"},"correlation_id":"corr-2"}
```

#### 3. Extended Decisions Log (JSONL)

**File:** `.sdp-memory/decisions/decisions.jsonl`

```jsonl
{"timestamp":"2026-02-06T10:15:00Z","type":"technical","feature_id":"F01","workstream_id":"00-001-01","question":"Should we use JWT or sessions for authentication?","decision":"Use JWT tokens","rationale":"Stateless, scales horizontally, works well with microservices","alternatives":["Sessions","OAuth2","API keys"],"outcome":"Expected: Faster auth, reduced database load","decision_maker":"user","tags":["authentication","security"],"session_id":"550e8400-e29b-41d4-a716-446655440000","related_decisions":[],"reverses_decision":"","reversed_by":[],"actual_outcome":"Successfully implemented, 40% faster auth","outcome_validated":true,"outcome_validated_at":"2026-02-10T15:00:00Z","status":"validated","approved_by":["user"],"implemented_at":"2026-02-07T10:00:00Z"}
```

#### 4. Timeline Log (JSONL)

**File:** `.sdp-memory/timeline/timeline.jsonl`

```jsonl
{"id":"tl-1","timestamp":"2026-02-06T10:00:00Z","type":"session_start","session_id":"550e8400-e29b-41d4-a716-446655440000","feature_id":"F01","summary":"Started session for F01 authentication feature","details":"Agent: planner","significance":"low"}
{"id":"tl-2","timestamp":"2026-02-06T10:15:00Z","type":"decision","decision_id":"d1","session_id":"550e8400-e29b-41d4-a716-446655440000","feature_id":"F01","workstream_id":"00-001-01","summary":"Chose JWT over sessions for auth","details":"Rationale: Stateless, scalable","significance":"high"}
{"id":"tl-3","timestamp":"2026-02-07T10:00:00Z","type":"implementation","decision_id":"d1","workstream_id":"00-001-01","summary":"Implemented JWT authentication","details":"Files: src/auth.go, src/auth_test.go","significance":"medium"}
```

---

### Storage APIs

#### SessionLogger Interface

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/session/logger.go` (new)

```go
package session

import "context"

// Logger manages session persistence
type Logger struct {
    sessionsPath string
    eventsPath   string
}

// NewLogger creates a new session logger
func NewLogger(projectRoot string) (*Logger, error)

// StartSession creates and starts a new session
func (l *Logger) StartSession(ctx context.Context, sess Session) error

// EndSession marks a session as complete
func (l *Logger) EndSession(ctx context.Context, sessionID string, endTime time.Time, status string) error

// RecordEvent appends an event to the session log
func (l *Logger) RecordEvent(ctx context.Context, event SessionEvent) error

// LoadSession retrieves a session by ID
func (l *Logger) LoadSession(ctx context.Context, sessionID string) (*Session, error)

// LoadActiveSessions retrieves all active (incomplete) sessions
func (l *Logger) LoadActiveSessions(ctx context.Context) ([]Session, error)

// ListSessions returns sessions matching filters
func (l *Logger) ListSessions(ctx context.Context, filters SessionFilters) ([]Session, error)

// SessionFilters defines query filters
type SessionFilters struct {
    ProjectID     string
    FeatureID     string
    Status        string
    StartedAfter  *time.Time
    StartedBefore *time.Time
    Agent         string // Filter by agent role
    Skill         string // Filter by skill used
}
```

#### TimelineBuilder Interface

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/timeline/builder.go` (new)

```go
package timeline

import "context"

// Builder constructs project timelines
type Builder struct {
    timelinePath string
    decisionLogger *decision.Logger
    sessionLogger  *session.Logger
    telemetryPath  string
}

// NewBuilder creates a new timeline builder
func NewBuilder(projectRoot string) (*Builder, error)

// BuildTimeline constructs a timeline from all sources
func (b *Builder) BuildTimeline(ctx context.Context, startTime, endTime *time.Time) (*Timeline, error)

// DetectMilestones identifies significant events
func (b *Builder) DetectMilestones(ctx context.Context) ([]Milestone, error)

// GetTimeline retrieves a pre-built timeline
func (b *Builder) GetTimeline(ctx context.Context) (*Timeline, error)

// Milestone represents a significant project event
type Milestone struct {
    ID          string    `json:"id"`
    Date        time.Time `json:"date"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Type        string    `json:"type"` // "first_commit", "first_deployment", "decision_reversal", etc.
    Significance string   `json:"significance"`
}
```

---

## Use Cases

### UC-001: Developer Searches for Similar Decisions

**Actor:** Developer
**Goal:** Find previous decisions about authentication to inform current work

**Flow:**
1. User invokes: `sdp memory search "authentication"`
2. System searches across:
   - Decision text (question, decision, rationale)
   - Session summaries
   - Timeline events
3. System returns ranked results:
   ```
   Found 3 decisions matching "authentication":

   1. JWT vs Sessions [2026-02-06]
      Type: technical
      Feature: F01
      Status: validated
      Decision: Use JWT tokens
      Rationale: Stateless, scales horizontally
      Actual Outcome: 40% faster auth performance

   2. OAuth2 Flow Choice [2026-01-15]
      Type: technical
      Feature: E05
      Status: reversed
      Decision: Use OAuth2 Authorization Code flow
      Rationale: Most secure for third-party integrations
      Reversed By: "Use custom SSO" (2026-01-20)
      Reason: Too complex for our use case

   3. Session Storage Backend [2025-12-10]
      Type: technical
      Feature: A02
      Status: implemented
      Decision: Redis for session storage
      Rationale: Fast in-memory lookups
   ```
4. User selects decision #1 to view details
5. System shows full decision + related workstreams + session context

**API Used:** `MemoryEngine.Search(query, filters)`

---

### UC-002: Understanding "Why Did We Choose This Approach?"

**Actor:** New team member (onboarding)
**Goal:** Understand rationale behind current architecture

**Flow:**
1. User views file: `src/auth/jwt.go`
2. User invokes: `sdp memory why src/auth/jwt.go`
3. System searches for decisions linked to this file:
   ```
   File: src/auth/jwt.go
   Linked to 2 decisions:

   1. Use JWT for authentication [2026-02-06]
      Session: 550e8400-e29b-41d4-a716-446655440000
      Rationale: Stateless, scales horizontally
      Alternatives: Sessions, OAuth2, API keys
      Outcome: 40% faster auth performance
      Related Workstreams: 00-001-01

   2. Add refresh token rotation [2026-02-10]
      Session: 660e8400-e29b-41d4-a716-446655440001
      Rationale: Improve security without UX impact
      Status: validated
   ```
4. User clicks through to view full session context

**API Used:** `MemoryEngine.GetDecisionsForFile(filePath)`

---

### UC-003: Productivity Analysis - "When Am I Most Productive?"

**Actor:** Developer
**Goal:** Optimize work schedule based on performance data

**Flow:**
1. User invokes: `sdp memory analyze productivity`
2. System aggregates session data:
   ```
   Productivity Analysis (Last 30 days)
   ======================================

   Session Summary:
     Total Sessions: 45
     Completed: 42 (93%)
     Interrupted: 3 (7%)

   Time Patterns:
     Most Productive Hours: 9am - 12pm (avg completion rate: 98%)
     Least Productive: 3pm - 5pm (avg completion rate: 75%)
     Avg Session Duration: 2h 15m

   Skill Usage:
     @build: 35 times (78% of sessions)
     @review: 28 times (62% of sessions)
     @design: 12 times (27% of sessions)

   Agent Performance:
     builder: 30 sessions (highest completion rate: 99%)
     reviewer: 28 sessions (avg quality gate pass: 94%)
     planner: 12 sessions (avg session duration: 45m)

   Trends:
     Productivity increased 15% over last month
     @review usage up 20% (improving quality)
   ```
3. User adjusts schedule based on insights

**API Used:** `AnalyticsEngine.GenerateProductivityReport(timeRange)`

---

### UC-004: Resume Interrupted Session

**Actor:** Developer
**Goal:** Resume work after interruption (crash, context switch)

**Flow:**
1. User invokes: `sdp memory resume`
2. System finds most recent interrupted session:
   ```
   Resume Session: 550e8400-e29b-41d4-a716-446655440000
   =================================================

   Feature: F01 (User Authentication)
   Workstreams: 00-001-02 (in progress), 00-001-03 (pending)

   Last Activity: 2026-02-06 12:30 (interrupted: system crash)

   Summary:
   Implemented JWT authentication layer, completed workstream 00-001-01.
   Currently working on 00-001-02 (Application services).

   Next Steps:
   1. Complete AuthService.Login() implementation
   2. Add unit tests for token validation
   3. Review quality gate results (failed: coverage 72%, need 80%)

   Context:
     Last Decision: "Use middleware for JWT validation" (d3)
     Files Modified: src/auth/service.go, src/auth/middleware.go
     Current Branch: feature/f01-auth
   ```
3. User confirms resume
4. System restores session context (reload variables, set active workstream)

**API Used:** `SessionLogger.LoadActiveSessions()`, `SessionLogger.Resume(sessionID)`

---

### UC-005: View Project Timeline for Onboarding

**Actor:** New team member
**Goal:** Understand project evolution

**Flow:**
1. User invokes: `sdp memory timeline --last 30d`
2. System renders chronological timeline:
   ```
   Project Timeline (Last 30 Days)
   ================================

   Feb 6, 2026 (Wednesday)
   ───────────────────────
   10:00  Started session for F01 authentication
   10:15  DECISION: Use JWT tokens (d1) ★ HIGH
   11:00  Started workstream 00-001-01 (Domain entities)
   12:30  Completed session (100% completion)

   Feb 7, 2026 (Thursday)
   ──────────────────────
   09:00  Started session: Implement JWT layer
   10:00  IMPLEMENTED: JWT authentication (d1)
   11:30  COMMIT: abc123def - "Add JWT middleware"
   14:00  QUALITY GATE: Passed (coverage 85%)

   Feb 10, 2026 (Sunday)
   ────────────────────
   15:00  OUTCOME VALIDATED: JWT decision (d1)
          Actual: 40% faster auth performance
   ```

3. User drills down into events for details

**API Used:** `TimelineBuilder.BuildTimeline(startTime, endTime)`

---

### UC-006: Detect Decision Pattern - "We Keep Reversing These"

**Actor:** Tech Lead
**Goal:** Identify problematic decision patterns

**Flow:**
1. User invokes: `sdp memory patterns --type reversals`
2. System analyzes decision graph:
   ```
   Decision Reversal Patterns
   ===========================

   High Reversal Rate Areas:
   ────────────────────────

   1. Authentication Strategy (3 reversals in 90 days)
      Feb 1: OAuth2 → Custom SSO
      Jan 20: Custom SSO → Sessions
      Jan 15: Sessions → OAuth2
      Pattern: Requirements keep changing
      Recommendation: Freeze auth requirements for 1 sprint

   2. Database Choice (2 reversals in 60 days)
      Dec 10: PostgreSQL → MongoDB
      Nov 15: MySQL → PostgreSQL
      Pattern: Unclear scalability needs
      Recommendation: Run load tests before deciding

   3. Testing Framework (1 reversal)
      Nov 1: pytest → unittest
      Reason: pytest incompatible with CI
      Status: Stable (no reversal since)
   ```
3. User shares insights with team to improve decision-making

**API Used:** `AnalyticsEngine.AnalyzeDecisionPatterns()`

---

## APIs and Interfaces

### 1. Memory Engine (Core Interface)

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/memory/engine.go` (new)

```go
package memory

import (
    "context"
    "time"
)

// Engine provides unified access to all memory systems
type Engine struct {
    decisionLogger *decision.Logger
    sessionLogger  *session.Logger
    timelineBuilder *timeline.Builder
    indexer        *Indexer
}

// NewEngine creates a new memory engine
func NewEngine(projectRoot string) (*Engine, error)

// Search performs full-text search across all memory types
func (e *Engine) Search(ctx context.Context, query string, filters SearchFilters) (*SearchResults, error)

// GetDecision retrieves a decision by ID
func (e *Engine) GetDecision(ctx context.Context, decisionID string) (*Decision, error)

// GetSession retrieves a session by ID
func (e *Engine) GetSession(ctx context.Context, sessionID string) (*Session, error)

// GetTimeline retrieves the project timeline
func (e *Engine) GetTimeline(ctx context.Context, startTime, endTime *time.Time) (*Timeline, error)

// GetDecisionsForFile finds decisions related to a file
func (e *Engine) GetDecisionsForFile(ctx context.Context, filePath string) ([]Decision, error)

// GetRelatedDecisions finds decisions linked to a given decision
func (e *Engine) GetRelatedDecisions(ctx context.Context, decisionID string) ([]Decision, error)

// SearchFilters defines search criteria
type SearchFilters struct {
    DateRange      *DateRange
    DecisionTypes  []string
    Features       []string
    Workstreams    []string
    Agents         []string
    Tags           []string
    MinSignificance string // "low", "medium", "high", "milestone"
}

// SearchResults contains ranked search results
type SearchResults struct {
    Decisions []DecisionResult `json:"decisions"`
    Sessions  []SessionResult  `json:"sessions"`
    Events    []EventResult    `json:"events"`
    Total     int              `json:"total"`
}

// DecisionResult is a scored decision search result
type DecisionResult struct {
    Decision   Decision `json:"decision"`
    Score      float64  `json:"score"`
    Highlights []string `json:"highlights"` // Matching text snippets
}
```

---

### 2. Analytics Engine

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/analytics/engine.go` (new)

```go
package analytics

import (
    "context"
    "time"
)

// Engine generates insights from memory data
type Engine struct {
    memoryEngine *memory.Engine
}

// NewEngine creates a new analytics engine
func NewEngine(memoryEngine *memory.Engine) (*Engine, error)

// GenerateProductivityReport analyzes session patterns
func (e *Engine) GenerateProductivityReport(ctx context.Context, timeRange DateRange) (*ProductivityReport, error)

// AnalyzeDecisionPatterns detects decision-making trends
func (e *Engine) AnalyzeDecisionPatterns(ctx context.Context) (*DecisionPatternsReport, error)

// GetAgentPerformance analyzes agent/skill effectiveness
func (e *Engine) GetAgentPerformance(ctx context.Context, timeRange DateRange) (*AgentPerformanceReport, error)

// ProductivityReport contains productivity insights
type ProductivityReport struct {
    TimeRange         DateRange                `json:"time_range"`
    SessionStats      SessionStatistics        `json:"session_stats"`
    TimePatterns      []TimePattern            `json:"time_patterns"`
    SkillUsage        map[string]int           `json:"skill_usage"`
    AgentPerformance  map[string]AgentMetrics  `json:"agent_performance"`
    Trends            []Trend                  `json:"trends"`
}

// DecisionPatternsReport contains decision pattern analysis
type DecisionPatternsReport struct {
    ReversalHotspots  []ReversalHotspot `json:"reversal_hotspots"`
    DecisionClusters  []DecisionCluster `json:"decision_clusters"` // Related decisions
    ValidationRate    float64           `json:"validation_rate"`   // % of decisions with validated outcomes
    AvgTimeToValidate time.Duration     `json:"avg_time_to_validate"`
}

// ReversalHotspot identifies areas with frequent reversals
type ReversalHotspot struct {
    Topic         string    `json:"topic"`
    ReversalCount int       `json:"reversal_count"`
    TimeSpan      time.Duration `json:"time_span"`
    Timeline      []Decision   `json:"timeline"`
    Recommendation string     `json:"recommendation"`
}
```

---

### 3. CLI Commands

**File:** `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/cmd/sdp/memory.go` (new)

```go
package main

import (
    "github.com/spf13/cobra"
)

var memoryCmd = &cobra.Command{
    Use:   "memory",
    Short: "Access long-term memory (decisions, sessions, history)",
    Long:  `Query and search institutional memory for decisions, sessions, and project history.`,
}

func init() {
    memoryCmd.AddCommand(memorySearchCmd())
    memoryCmd.AddCommand(memoryTimelineCmd())
    memoryCmd.AddCommand(memoryDecisionsCmd())
    memoryCmd.AddCommand(memorySessionsCmd())
    memoryCmd.AddCommand(memoryAnalyzeCmd())
    memoryCmd.AddCommand(memoryResumeCmd())
}

// memorySearchCmd searches across all memory
func memorySearchCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "search <query>",
        Short: "Search decisions, sessions, and timeline",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
        },
    }
}

// memoryTimelineCmd shows project timeline
func memoryTimelineCmd() *cobra.Command {
    var lastDays int

    cmd := &cobra.Command{
        Use:   "timeline",
        Short: "Show project history timeline",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
        },
    }

    cmd.Flags().IntVarP(&lastDays, "last", "l", 30, "Show last N days")
    return cmd
}

// memoryAnalyzeCmd generates analytics reports
func memoryAnalyzeCmd() *cobra.Command {
    var reportType string

    cmd := &cobra.Command{
        Use:   "analyze <report-type>",
        Short: "Generate analytics report",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
        },
    }

    cmd.Flags().StringVarP(&reportType, "type", "t", "productivity", "Report type: productivity, patterns, agents")
    return cmd
}

// memoryResumeCmd resumes interrupted session
func memoryResumeCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "resume [session-id]",
        Short: "Resume interrupted session",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
        },
    }
}
```

---

## Implementation Phases

### Phase 1: Foundation (P0 - 2 weeks)

**Goal:** Extend existing decision tracking, add session management

**Workstreams:**
1. **WS-MEM-001:** Extend `Decision` struct with new fields
2. **WS-MEM-002:** Create `Session` and `SessionEvent` structs
3. **WS-MEM-003:** Implement `SessionLogger` (JSONL storage)
4. **WS-MEM-004:** Migrate existing decisions to `.sdp-memory/` directory
5. **WS-MEM-005:** Add `sdp memory search` (basic text search)

**Deliverables:**
- Extended `Decision` struct (backward compatible)
- `Session` struct + `SessionLogger` implementation
- CLI: `sdp memory search <query>`
- Decision index for fast lookups

---

### Phase 2: Timeline & Analytics (P1 - 3 weeks)

**Goal:** Build timeline and basic analytics

**Workstreams:**
1. **WS-MEM-006:** Implement `TimelineBuilder`
2. **WS-MEM-007:** Add `sdp memory timeline` command
3. **WS-MEM-008:** Implement session analytics (productivity metrics)
4. **WS-MEM-009:** Add `sdp memory analyze` command
5. **WS-MEM-010:** Create decision pattern analysis

**Deliverables:**
- Timeline view (CLI: `sdp memory timeline`)
- Productivity report (CLI: `sdp memory analyze productivity`)
- Decision reversal pattern detection

---

### Phase 3: Enhanced Search & Resume (P1 - 2 weeks)

**Goal:** Advanced search and session resumption

**Workstreams:**
1. **WS-MEM-011:** Implement full-text search index (SQLite FTS5 or blevesearch)
2. **WS-MEM-012:** Add semantic search (vector embeddings)
3. **WS-MEM-013:** Implement `sdp memory resume` command
4. **WS-MEM-014:** Add file-to-decision linkage (via git blame)

**Deliverables:**
- Fast full-text search across all memory
- Semantic similarity search
- Session resumption workflow

---

### Phase 4: Polish & Integration (P2 - 2 weeks)

**Goal:** Export, visualization, AI integration

**Workstreams:**
1. **WS-MEM-015:** Add memory export to markdown/JSON
2. **WS-MEM-016:** Create dashboard UI (terminal-based or web)
3. **WS-MEM-017:** Integrate with Claude Code (provide memory as context)
4. **WS-MEM-018:** Documentation and tutorials

**Deliverables:**
- Export workflows (markdown reports)
- Optional dashboard UI
- Claude Code integration (auto-inject relevant decisions)
- User documentation

---

## Technical Considerations

### 1. Backward Compatibility

**Challenge:** Existing `Decision` struct is used throughout codebase.

**Solution:**
- All new fields in `Decision` struct are `omitempty` (JSON)
- Existing code continues to work without modifications
- Migration script updates old decisions to new format (adds defaults)
- CLI flag `--legacy` uses old format if needed

---

### 2. Performance

**Challenge:** Searching thousands of decisions/sessions could be slow.

**Solutions:**
- **Index File:** Maintain `.sdp-memory/index.json` for fast lookups
- **FTS Index:** Use SQLite FTS5 or blevesearch for full-text search
- **Lazy Loading:** Only load decision details on demand
- **Pagination:** Limit search results to top 20, allow pagination
- **Caching:** Cache frequently accessed decisions in memory

---

### 3. Concurrency

**Challenge:** Multiple processes may write to memory simultaneously (e.g., two @build runs).

**Solutions:**
- **File Locking:** Use `flock()` on JSONL files during writes
- **Append-Only:** JSONL format supports concurrent appends (with locking)
- **Atomic Writes:** Write to temp file, then rename (for index updates)
- **Mutex Protection:** In-memory caches protected by `sync.Mutex`

---

### 4. Privacy

**Challenge:** Memory may contain sensitive information.

**Solutions:**
- **Opt-In:** Memory collection disabled by default (like telemetry)
- **Local Only:** No data transmission (all stored in `.sdp-memory/`)
- **GitIgnore Config:** `.sdp/memory.gitignore` controls what gets committed
- **PII Filtering:** Strip sensitive data (usernames, paths) before indexing
- **Encryption:** Optional encryption for `.sdp-memory/` (future enhancement)

---

### 5. Storage Growth

**Challenge:** Memory logs grow unbounded over time.

**Solutions:**
- **Rotation:** Rotate monthly (like telemetry) → `.sdp-memory/sessions/2026-02/`
- **Compression:** Compress old logs with gzip
- **Archive:** Move logs > 1 year to `.sdp-memory/archive/`
- **Retention Config:** User-configurable retention policy (default: keep all)
- **Pruning:** CLI command `sdp memory prune --before 2024-01-01`

---

## Success Criteria

### Functional Success
- [ ] Users can search decisions by keyword (accuracy > 90%)
- [ ] Users can view project timeline (last 30 days in < 2 seconds)
- [ ] Users can resume interrupted sessions (restore context in < 5 seconds)
- [ ] Analytics detect decision reversals (precision > 80%)
- [ ] Export to markdown works (all decisions included)

### Non-Functional Success
- [ ] Search latency < 500ms for 10,000 decisions
- [ ] Memory overhead < 100MB for 1 year of data
- [ ] No data loss (concurrent writes tested with 10 parallel agents)
- [ ] Backward compatibility with existing decisions (100%)
- [ ] Test coverage > 80% on all new code

---

## Open Questions

1. **Storage Backend:** Should we use SQLite for indexes (complex) or JSON files (simple)?
   - **Recommendation:** Start with JSON + simple in-memory index, add SQLite later if needed

2. **Search Engine:** Should we use external libraries (blevesearch, sqlite3) or pure Go?
   - **Recommendation:** Pure Go for MVP, add SQLite FTS5 in Phase 3 if performance insufficient

3. **Session ID Format:** UUID v4 vs. timestamp-based vs. sequential?
   - **Recommendation:** UUID v4 (unique, distributed, no coordination needed)

4. **Memory Retention:** Should old memory be auto-deleted?
   - **Recommendation:** No, keep forever (user can manually prune)

5. **Git Integration:** Should memory be auto-committed?
   - **Recommendation:** Opt-in (add `memory.autoCommit` to `.sdp/config.json`)

---

## Appendices

### Appendix A: Example Session Flow

```
User: @feature "Add user authentication"

Claude (Systems Analyst):
→ Creating session: sess-123
→ Recording event: skill_invoked (@feature)
→ Interviewing user...

[User answers questions]

Claude (Systems Analyst):
→ Decision made: "Use JWT tokens"
→ Logging decision: d1
→ Linking decision to session: sess-123
→ Recording event: decision_made (d1)

Claude (Technical Architect):
→ @design invoked
→ Recording event: skill_invoked (@design)
→ Creating workstreams...

[User approves design]

Claude (Builder):
→ @build 00-001-01 invoked
→ Recording event: workstream_started (00-001-01)
→ Implementing...

[Build completes]

Claude (Builder):
→ Recording event: workstream_completed (00-001-01)
→ Updating session status: in_progress

[User invokes @build 00-001-02]

[...]

[Session ends - crash or user exit]

Claude (Session Manager):
→ Session interrupted: sess-123
→ Saving context for resumption
→ User can run: sdp memory resume
```

---

### Appendix B: Decision Lifecycle State Machine

```
[proposed]
    |
    v
[approved] <-- approved_by added
    |
    v
[implemented] <-- implemented_at set
    |
    +-> [validated] <-- outcome_validated = true
    |
    +-> [reversed] <-- reversed_by updated
```

**Transitions:**
- `proposed` → `approved`: User or agent approves decision
- `approved` → `implemented`: Decision code merged to main
- `implemented` → `validated`: Outcome confirmed (after monitoring)
- `implemented` → `reversed`: New decision reverses this one
- `validated` → `reversed`: Even validated decisions can be reversed

---

### Appendix C: File Schema Reference

#### Decision Schema (v2)
```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Decision",
  "type": "object",
  "required": ["timestamp", "type", "question", "decision", "rationale", "status"],
  "properties": {
    "timestamp": {"type": "string", "format": "date-time"},
    "type": {"enum": ["vision", "technical", "tradeoff", "explicit"]},
    "feature_id": {"type": "string"},
    "ws_id": {"type": "string"},
    "question": {"type": "string"},
    "decision": {"type": "string"},
    "rationale": {"type": "string"},
    "alternatives": {"type": "array", "items": {"type": "string"}},
    "outcome": {"type": "string"},
    "decision_maker": {"type": "string"},
    "tags": {"type": "array", "items": {"type": "string"}},
    "session_id": {"type": "string"},
    "related_decisions": {"type": "array", "items": {"type": "string"}},
    "reverses_decision": {"type": "string"},
    "reversed_by": {"type": "array", "items": {"type": "string"}},
    "actual_outcome": {"type": "string"},
    "outcome_validated": {"type": "boolean"},
    "outcome_validated_at": {"type": "string", "format": "date-time"},
    "status": {"enum": ["proposed", "approved", "implemented", "validated", "reversed"]},
    "approved_by": {"type": "array", "items": {"type": "string"}},
    "implemented_at": {"type": "string", "format": "date-time"}
  }
}
```

---

## Related Documents

- **Existing Systems:**
  - `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/decision/decision.go` - Decision struct
  - `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/decision/logger.go` - Decision logger
  - `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/telemetry/collector.go` - Telemetry collector
  - `/Users/fall_out_bug/projects/vibe_coding/sdp/sdp-plugin/internal/beads/client.go` - Beads client (task tracking)

- **Documentation:**
  - `docs/PROTOCOL.md` - SDP workflow specification
  - `docs/PRIVACY.md` - Privacy policy (telemetry)
  - `docs/TELEMETRY_HOWTO.md` - Telemetry usage guide
  - `CLAUDE.md` - Project instructions

---

**End of Specification**
