# Observability Bridge Design

**Version:** 1.0
**Status:** Draft
**Author:** SDP Team
**Date:** 2026-02-11

## Abstract

This document specifies the architecture for bridging SDP's evidence layer (what happened during build) with production observability tools (what's happening in production). The bridge enables rapid incident response by connecting deploy markers to evidence chains, allowing engineers to trace production issues back to AI-generated code and human modifications.

---

## Table of Contents

1. [Overview](#overview)
2. [Deploy Markers](#deploy-markers)
3. [OpenTelemetry Span Attributes](#opentelemetry-span-attributes)
4. [Diff-Level Provenance](#diff-level-provenance)
5. [Integration Architecture](#integration-architecture)
6. [Integration Points](#integration-points)
7. [Data Flow Examples](#data-flow-examples)
8. [Privacy & Security](#privacy--security)
9. [Implementation Roadmap](#implementation-roadmap)

---

## Overview

### The Problem

When production incidents occur, engineers need to answer three questions rapidly:

1. **What changed?** - Which code paths are affected?
2. **Who wrote it?** - AI or human? Which model?
3. **How was it tested?** - What quality gates passed?

Today, observability tools answer "what's breaking" but not "what caused it." The evidence layer knows what happened during build, but that knowledge stays disconnected from production.

### The Solution

The **Observability Bridge** connects evidence to production through:

1. **Deploy Markers** - Git tags with SDP metadata embedded
2. **OTel Attributes** - Span enrichment with AI provenance
3. **Diff Attribution** - Per-line AI/human code labeling
4. **Tool Integrations** - Honeycomb, Datadog, Grafana

### Key Decision: Push vs Pull

**Decision:** Push model (SDP exports to observability tools)

**Rationale:**

- **Latency:** Pull during incident response adds critical delay
- **Reliability:** Observability tools may be rate-limited or unavailable
- **Simplicity:** Single export at deploy time vs continuous querying
- **Offline:** Evidence works even if observability tools are down

**Alternative (Pull):** Observability tools query SDP API during investigation
- Rejected: Adds 5-10s latency to incident response
- Rejected: Requires SDP API to be highly available

---

## Deploy Markers

### Purpose

Attach SDP metadata to deploy events so observability tools can query evidence context.

### Format

Deploy markers are **annotated git tags** following naming convention:

```
deploy/<environment>/<timestamp>
```

Examples:
- `deploy/production/20260211-120000`
- `deploy/staging/20260211-113000`
- `deploy/dev/20260211-100000`

### Tag Annotations (Structured)

```bash
git tag -a -m "{
  \"sdp_version\": \"0.9.0\",
  \"environment\": \"production\",
  \"commit_sha\": \"abc123def456\",
  \"deployed_at\": \"2026-02-11T12:00:00Z\",
  \"evidence_log\": \".sdp/log/events.jsonl\",
  \"workstreams\": [\"00-059-01\", \"00-059-02\"],
  \"feature_id\": \"F059\",
  \"deployer\": \"ci-system\",
  \"verification\": {
    \"coverage\": 87.5,
    \"gates_passed\": [\"test\", \"typecheck\", \"lint\"]
  }
}" deploy/production/20260211-120000
```

### Tag Message Schema

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "SDP Deploy Marker",
  "type": "object",
  "required": [
    "sdp_version",
    "environment",
    "commit_sha",
    "deployed_at",
    "evidence_log"
  ],
  "properties": {
    "sdp_version": {
      "type": "string",
      "description": "SDP protocol version"
    },
    "environment": {
      "type": "string",
      "enum": ["production", "staging", "dev"],
      "description": "Deployment environment"
    },
    "commit_sha": {
      "type": "string",
      "pattern": "^[a-f0-9]{40}$",
      "description": "Full Git SHA being deployed"
    },
    "deployed_at": {
      "type": "string",
      "format": "date-time",
      "description": "Deployment timestamp (ISO 8601)"
    },
    "evidence_log": {
      "type": "string",
      "description": "Path to evidence log relative to repo root"
    },
    "workstreams": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Workstream IDs included in this deploy"
    },
    "feature_id": {
      "type": "string",
      "description": "Feature ID if applicable"
    },
    "deployer": {
      "type": "string",
      "description": "Human or system that deployed"
    },
    "verification": {
      "type": "object",
      "properties": {
        "coverage": { "type": "number" },
        "gates_passed": {
          "type": "array",
          "items": { "type": "string" }
        }
      }
    }
  }
}
```

### Tag Creation Workflow

```bash
# 1. Deploy (via CI/CD or @deploy skill)
git checkout main
git pull
# Deployment happens...

# 2. Create deploy marker
TIMESTAMP=$(date -u +"%Y%m%d-%H%M%S")
git tag -a -m "$(cat deploy-marker.json)" "deploy/production/${TIMESTAMP}"

# 3. Push tag to remote
git push origin "deploy/production/${TIMESTAMP}"

# 4. Notify observability tools (see Integration Points)
# Honeycomb, Datadog, Grafana receive deploy event
```

### Tag Query (Incident Response)

```bash
# Find deploy marker for current production
git fetch --tags
TAG=$(git tag -l "deploy/production/*" --sort=-taggerdate | head -n 1)

# Extract metadata
MSG=$(git tag -l "$TAG" --format='%(contents)')
echo $MSG | jq .

# Get evidence log path
EVIDENCE_LOG=$(echo $MSG | jq -r '.evidence_log')

# Trace events for this deploy
sdp log trace --since $(echo $MSG | jq -r '.deployed_at')
```

### Alternatives Considered

#### Alternative 1: OTel Deploy Events
- Emit OTel deploy event directly from CI/CD
- **Rejected:** Requires OTel collector in CI pipeline, adds complexity

#### Alternative 2: Git Notes
- Store deploy metadata in `refs/notes/deploy`
- **Rejected:** Less discoverable than tags, harder to query

#### Alternative 3: External Deploy Tracker
- Store deploy metadata in database (e.g., DeployTracker)
- **Rejected:** Adds external dependency, breaks when DB is down

---

## OpenTelemetry Span Attributes

### Purpose

Enrich OTel spans with AI-generated code provenance so engineers can filter traces by AI vs human code, model ID, and verification status.

### Proposed Semantic Convention

**Namespace:** `sdp.` (SDP-specific attributes)

**Justification:**
- SDP is not an official OTel semantic convention (yet)
- Using vendor namespace prevents conflicts with future standards
- See `docs/design/OTEL-AI-PROVENANCE.md` for standardization proposal

### Core Attributes

| Attribute | Type | Description | Example |
|-----------|------|-------------|---------|
| `sdp.ai_generated` | boolean | Was this code generated by AI? | `true`, `false` |
| `sdp.model` | string | Model identifier if AI-generated | `claude-sonnet-4-20250514` |
| `sdp.evidence_id` | string | Link to evidence event UUID | `evt-123e4567-e89b-12d3` |
| `sdp.verification_status` | string | Quality gate result | `passed`, `failed`, `skipped` |
| `sdp.workstream_id` | string | Workstream that generated code | `00-059-01` |
| `sdp.feature_id` | string | Feature ID | `F059` |
| `sdp.human_modified` | boolean | Was AI code modified by human? | `true`, `false` |

### Full Schema

```yaml
# OTel Span Attributes for AI Provenance
sdp:
  # Core provenance
  ai_generated:
    type: boolean
    description: "True if code path was generated by AI"
    required: true
    example: true

  model:
    type: string
    description: "Model identifier (if ai_generated=true)"
    condition: "required if ai_generated=true"
    example: "claude-sonnet-4-20250514"

  evidence_id:
    type: string
    description: "UUID of generation event in evidence log"
    condition: "required if ai_generated=true"
    example: "evt-123e4567-e89b-12d3-a456-426614174000"

  # Verification context
  verification_status:
    type: string
    enum: [passed, failed, skipped]
    description: "Quality gate result for this code"
    example: "passed"

  verification_coverage:
    type: double
    description: "Test coverage percentage (0-100)"
    example: 87.5

  # Workstream context
  workstream_id:
    type: string
    description: "Workstream ID that generated this code"
    pattern: "^\\d{2}-\\d{3}-\\d{2}$"
    example: "00-059-01"

  feature_id:
    type: string
    description: "Feature ID (if applicable)"
    pattern: "^F\\d+$"
    example: "F059"

  # Human modification
  human_modified:
    type: boolean
    description: "Was AI-generated code modified by human?"
    example: true

  human_author:
    type: string
    description: "Username who modified (if human_modified=true, no email for PII compliance)"
    condition: "optional if human_modified=true"
    example: "@jane_dev"
```

### Span Examples

#### Example 1: Pure AI-Generated Code

```python
# User authentication endpoint (00-059-01, generated by Claude Sonnet 4)

from opentelemetry import trace

tracer = trace.get_tracer(__name__)

@tracer.start_as_current_span("login")
def login(username: str, password: str) -> Token:
    # Span attributes set automatically by SDP instrumentation
    span = trace.get_current_span()
    span.set_attribute("sdp.ai_generated", True)
    span.set_attribute("sdp.model", "claude-sonnet-4-20250514")
    span.set_attribute("sdp.evidence_id", "evt-123e4567-e89b-12d3")
    span.set_attribute("sdp.workstream_id", "00-059-01")
    span.set_attribute("sdp.verification_status", "passed")
    span.set_attribute("sdp.verification_coverage", 92.0)
    span.set_attribute("sdp.human_modified", False)

    # Implementation...
```

**Resulting Span:**
```json
{
  "trace_id": "abc123...",
  "span_id": "def456...",
  "name": "login",
  "attributes": {
    "sdp.ai_generated": true,
    "sdp.model": "claude-sonnet-4-20250514",
    "sdp.evidence_id": "evt-123e4567-e89b-12d3",
    "sdp.workstream_id": "00-059-01",
    "sdp.verification_status": "passed",
    "sdp.verification_coverage": 92.0,
    "sdp.human_modified": false
  }
}
```

#### Example 2: Human-Written Code

```python
# Critical security validation (human-written, no AI involved)

@tracer.start_as_current_span("validate_password")
def validate_password(password: str) -> bool:
    span = trace.get_current_span()
    span.set_attribute("sdp.ai_generated", False)
    # No model, evidence_id, or workstream_id

    # Human-only code path...
```

**Resulting Span:**
```json
{
  "trace_id": "abc123...",
  "span_id": "xyz789...",
  "name": "validate_password",
  "attributes": {
    "sdp.ai_generated": false
  }
}
```

#### Example 3: AI Code Modified by Human

```python
# Payment processing (AI-generated, then human-optimized)

@tracer.start_as_current_span("process_payment")
def process_payment(amount: int, currency: str) -> PaymentResult:
    span = trace.get_current_span()
    span.set_attribute("sdp.ai_generated", True)
    span.set_attribute("sdp.model", "claude-sonnet-4-20250514")
    span.set_attribute("sdp.evidence_id", "evt-234e5678-f89b-23d4")
    span.set_attribute("sdp.workstream_id", "00-060-05")
    span.set_attribute("sdp.human_modified", True)
    span.set_attribute("sdp.human_author", "@alice_security")
    # Note: verification_status reflects current state (after human changes)

    # Human-optimized AI code...
```

### Attribute Propagation

SDP attributes propagate through span links:

```
[HTTP: POST /login]
  └── [Controller: login]
      └── [Service: authenticate]
          └── [Repository: find_user]
              └── (human-written DB query, sdp.ai_generated=false)

All child spans inherit parent's SDP attributes unless explicitly overridden.
```

### Implementation: Automatic Attribute Injection

**Approach:** OTel Processor that inspects git blame and evidence log

```go
// Pseudo-code for OTel SpanProcessor
type SDPProvenanceProcessor struct {
    evidenceLog string
    repo        *git.Repository
}

func (p *SDPProvenanceProcessor) OnStart(span, parent SpanContext) {
    // 1. Get calling code location (runtime.Caller)
    file, line := getCallerFrame()

    // 2. Git blame to get commit and author
    blame, _ := p.repo.Blame(file, line)

    // 3. Search evidence log for generation event
    evidence := p.findEvidence(blame.Commit)

    // 4. Set span attributes
    if evidence != nil {
        span.SetAttribute("sdp.ai_generated", true)
        span.SetAttribute("sdp.model", evidence.ModelID)
        span.SetAttribute("sdp.evidence_id", evidence.ID)
        // ...
    } else {
        span.SetAttribute("sdp.ai_generated", false)
    }
}
```

**Complexity:** O(log n) per span using git blame cache
**Performance:** <1ms overhead per span (acceptable for most applications)

---

## Diff-Level Provenance

### Purpose

Attribute each line of code to AI or human authorship, enabling:
- Precise blame for incidents
- Coverage metrics by author type
- "What changed since last deploy?" with AI/human breakdown

### Algorithm

#### High-Level Approach

1. **Parse Diff** - Extract added/modified lines from deploy diff
2. **Git Blame** - Get commit SHA for each line
3. **Evidence Lookup** - Search evidence log for generation event
4. **Classification** - AI (if evidence found) or Human (if no evidence)
5. **Annotation** - Emit provenance map as JSON

#### Detailed Algorithm

```python
def compute_diff_provenance(deploy_commit: str, prev_commit: str) -> Dict[str, LineProvenance]:
    """
    Compute AI vs human attribution for each line in a deploy diff.

    Returns:
        Mapping of filepath:line -> attribution data
    """
    provenance = {}

    # Step 1: Get diff between commits
    diff = git_diff(prev_commit, deploy_commit)

    for file in diff.changed_files:
        for line in file.added_or_modified_lines:
            line_number = line.number

            # Step 2: Git blame for this specific line
            blame = git_blame(deploy_commit, file.path, line_number)

            # Step 3: Search evidence log for generation event
            evidence = find_generation_event(blame.commit_sha)

            # Step 4: Classify
            if evidence:
                # AI-generated
                provenance[f"{file.path}:{line_number}"] = LineProvenance(
                    author_type="ai",
                    model=evidence.model,
                    evidence_id=evidence.id,
                    workstream_id=evidence.ws_id,
                    generated_at=evidence.timestamp
                )
            else:
                # Human-written (no evidence event found)
                # Note: Strip email from author for PII compliance (use username only)
                provenance[f"{file.path}:{line_number}"] = LineProvenance(
                    author_type="human",
                    commit_sha=blame.commit_sha,
                    author=blame.author.split('<')[0].strip(),  # "Jane Developer <jane@example.com>" -> "Jane Developer"
                    authored_at=blame.date
                )

    return provenance

def find_generation_event(commit_sha: str) -> Optional[GenerationEvent]:
    """
    Search evidence log for generation event matching this commit.
    """
    events = load_evidence_log()

    for event in events:
        if event.type == "generation" and event.commit_sha == commit_sha:
            return GenerationEvent(
                id=event.id,
                model=event.data.model,
                ws_id=event.ws_id,
                timestamp=event.timestamp
            )

    return None
```

### Output Format

```json
{
  "version": "1.0",
  "deploy_commit": "abc123def456",
  "previous_commit": "def456abc123",
  "generated_at": "2026-02-11T12:00:00Z",
  "summary": {
    "total_lines_changed": 150,
    "ai_generated_lines": 120,
    "human_lines": 30,
    "ai_percentage": 80.0
  },
  "files": {
    "src/auth/login.go": {
      "total_lines": 50,
      "ai_lines": [
        {
          "line": 15,
          "author_type": "ai",
          "model": "claude-sonnet-4-20250514",
          "evidence_id": "evt-123e4567-e89b-12d3",
          "workstream_id": "00-059-01",
          "generated_at": "2026-02-11T10:00:00Z"
        },
        {
          "line": 16,
          "author_type": "ai",
          "model": "claude-sonnet-4-20250514",
          "evidence_id": "evt-123e4567-e89b-12d3",
          "workstream_id": "00-059-01",
          "generated_at": "2026-02-11T10:00:00Z"
        }
      ],
      "human_lines": [
        {
          "line": 17,
          "author_type": "human",
          "commit_sha": "789abc012def",
          "author": "@jane_dev",
          "authored_at": "2026-02-11T11:00:00Z"
        }
      ]
    }
  }
}
```

### Usage: Incident Response

```bash
# 1. Get diff provenance for current deploy
sdp provenance diff deploy/production/20260211-120000 > provenance.json

# 2. Filter by files involved in incident
jq '.files["src/auth/login.go"]' provenance.json

# 3. Check if crashing line is AI-generated
# Incident: line 16 of login.go throws nil pointer error
jq '.files["src/auth/login.go"].ai_lines[] | select(.line == 16)' provenance.json

# Output: Shows AI provenance, links to evidence event
# {
#   "line": 16,
#   "author_type": "ai",
#   "model": "claude-sonnet-4-20250514",
#   "evidence_id": "evt-123e4567-e89b-12d3",
#   ...
# }

# 4. Pull full evidence for that generation
sdp log trace --evt evt-123e4567-e89b-12d3
```

### Usage: Coverage Reports

```bash
# Generate coverage breakdown by author type
sdp provenance coverage --deploy deploy/production/20260211-120000

# Output:
# Coverage by Author Type:
# - AI-generated: 87.5% (105/120 lines)
# - Human-written: 95.0% (28.5/30 lines)
# - Overall: 88.7% (133.5/150 lines)
#
# Low Coverage Files (AI-generated):
# - src/payment/processor.go: 65.0% (13/20 lines)
#   Workstream: 00-060-05
#   Model: claude-sonnet-4-20250514
```

### Performance Considerations

| Operation | Complexity | Time (1000 files) |
|-----------|-----------|-------------------|
| Git diff | O(n) | 500ms |
| Git blame (batch) | O(n log m) | 2s |
| Evidence search | O(k) | 100ms |
| **Total** | - | **~2.6s** |

Optimizations:
- Cache git blame results by commit SHA
- Index evidence log by commit_sha (in-memory lookup)
- Parallelize file processing

---

## Integration Architecture

### Data Flow Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          Development Phase                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  [AI Agent]               [Evidence Layer]              [Git Repo]       │
│      |                           |                          |            │
│      | @build 00-059-01          | Write event              | Commit     │
│      |-------------------------->| .sdp/log/events.jsonl  |----------->│
│      |                           |                          |            │
│      |                           | Hash chain               |            │
│      |                           | prev_hash linkage        |            │
└─────────────────────────────────────────────────────────────────────────┘
                                        |
                                        v
┌─────────────────────────────────────────────────────────────────────────┐
│                          Deployment Phase                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  [CI/CD]              [Deploy Marker]                 [Git Remote]      │
│     |                         |                            |             │
│     | Deploy prod             | Create tag                | Push tag    │
│     |------------------------>| deploy/prod/<timestamp>  |------------->│
│     |                         |--------------------------->│             │
│     |                         |                            |             │
│     | Export provenance       | Compute diff              |             │
│     |------------------------>| Git blame + evidence      │             │
│     |                         |                            │             │
└─────────────────────────────────────────────────────────────────────────┘
                                        |
                                        v
┌─────────────────────────────────────────────────────────────────────────┐
│                     Observability Ingestion                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  [SDP Exporter]     [OTel Collector]                [Observability]     │
│       |                    |                              |             │
│       | OTel traces        | Batch process                | Honeycomb   │
│       | with sdp.* attrs   |----------------------------->│ Datadog     │
│       |------------------->|                              │ Grafana     │
│       |                    |                              |             │
│       | Deploy event       | Forward deploy marker       |             │
│       |------------------->|----------------------------->│             │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
                                        |
                                        v
┌─────────────────────────────────────────────────────────────────────────┐
│                      Incident Response                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  [SRE]                  [Observability UI]              [Evidence]       │
│    |                            |                           |            │
│    | Incident alert             | Filter: sdp.ai_generated | Query      │
│    |--------------------------->|-------------------------->| evidence   │
│    |                            |                           |            │
│    | View span                  | Show: model,          | Trace      │
│    | with error                 |   evidence_id            | chain      │
│    |--------------------------->|-------------------------->|            │
│    |                            |                           |            │
│    | Click evidence_id          | Jump to evidence log     │            │
│    |--------------------------->|-------------------------->|            │
│    |                            |                           |            │
│    | Root cause found!          |                           |            │
│    |--------------------------->|                           |            │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### Component Diagram

```mermaid
graph TD
    A[AI Agent: Claude Code] -->|@build| B[Evidence Layer]
    B -->|Write event| C[.sdp/log/events.jsonl]
    C -->|Hash chain| C

    D[CI/CD: GitHub Actions] -->|Deploy| E[Git Remote]
    E -->|Create tag| F[Deploy Marker]

    F -->|Git blame| G[Diff Provenance Calculator]
    G -->|Lookup| C
    G -->|Compute| H[Provenance Map JSON]

    I[SDP Exporter] -->|Enrich spans| J[OTel Collector]
    J -->|OTLP| K[Honeycomb]
    J -->|OTLP| L[Datadog]
    J -->|OTLP| M[Grafana Tempo]

    F -->|Deploy event| J

    N[SRE Dashboard] -->|Filter: sdp.ai_generated=true| K
    K -->|Show evidence_id| N
    N -->|Click link| O[Evidence Log UI]
    O -->|Trace events| C

    style A fill:#e1f5fe
    style B fill:#fff9c4
    style C fill:#fff9c4
    style F fill:#f3e5f5
    style G fill:#f3e5f5
    style H fill:#f3e5f5
    style K fill:#e8f5e9
    style L fill:#e8f5e9
    style M fill:#e8f5e9
    style N fill:#fce4ec
    style O fill:#fce4ec
```

### Sequence Diagram: Incident Response

```mermaid
sequenceDiagram
    participant User as End User
    participant App as Application
    participant Honeycomb as Honeycomb
    participant SRE as SRE
    participant Evidence as Evidence Log

    User->>App: HTTP POST /login
    App->>App: nil pointer error!
    App->>Honeycomb: Error span (sdp.ai_generated=true, sdp.evidence_id=evt-123)
    Honeycomb->>SRE: Alert: High error rate on /login

    SRE->>Honeycomb: Filter: sdp.ai_generated=true, error=true
    Honeycomb-->>SRE: Show span with sdp.evidence_id=evt-123e4567-e89b-12d3

    SRE->>Evidence: GET /evidence/evt-123e4567-e89b-12d3
    Evidence-->>SRE: Generation event:
                    - Model: claude-sonnet-4-20250514
                    - Workstream: 00-059-01
                    - Prompt hash: abc123...
                    - Files: src/auth/login.go:15-20

    SRE->>Evidence: Trace evidence chain
    Evidence-->>SRE: Full history:
                    - Plan: @feature "Add OAuth login"
                    - Generation: AI generated login handler
                    - Verification: Tests passed (92% coverage)
                    - Approval: @review F059 approved

    SRE->>SRE: Root cause identified:
              AI code at line 16 doesn't handle nil user.
              Fix: Add nil check (human patch).

    SRE->>App: Hotfix deployed (rollback + patch)
```

---

## Integration Points

### Honeycomb

#### Overview

Honeycomb is a columnar observability tool optimized for high-cardinality data like SDP provenance attributes.

#### Configuration

```yaml
# honeycomb-sdp.yaml
honeycomb_api_key: ${HONEYCOMB_API_KEY}
dataset: production-traces

# SDP-specific config
sdp_provenance_enabled: true
sdp_evidence_endpoint: https://evidence.example.com
sdp_deploy_markers:
  - git_tag: deploy/production/*
```

#### Span Queries

**Find all errors in AI-generated code:**
```sql
SELECT sdp.model, COUNT(*) as error_count
FROM production-traces
WHERE error = true
  AND sdp.ai_generated = true
GROUP BY sdp.model
ORDER BY error_count DESC
```

**Find low-verification AI code with high latency:**
```sql
SELECT sdp.workstream_id, sdp.verification_coverage, AVG(duration_ms)
FROM production-traces
WHERE sdp.ai_generated = true
  AND sdp.verification_coverage < 80
GROUP BY sdp.workstream_id, sdp.verification_coverage
HAVING AVG(duration_ms) > 1000
```

**Compare error rates: AI vs Human**
```sql
SELECT
  sdp.ai_generated,
  COUNT(*) as total_spans,
  COUNT_IF(error = true) as error_spans,
  (error_spans * 100.0 / total_spans) as error_rate
FROM production-traces
WHERE duration_ms > 100
GROUP BY sdp.ai_generated
```

#### Dashboard Configuration

**Panel 1: AI vs Human Error Rate**
```yaml
type: comparison
query_a: SELECT COUNT_IF(error=true)*100.0/COUNT(*) FROM traces WHERE sdp.ai_generated=true
query_b: SELECT COUNT_IF(error=true)*100.0/COUNT(*) FROM traces WHERE sdp.ai_generated=false
```

**Panel 2: Top Error-Prone Models**
```yaml
type: table
query: SELECT sdp.model, COUNT_IF(error=true) as errors FROM traces WHERE sdp.ai_generated=true GROUP BY sdp.model ORDER BY errors DESC LIMIT 10
```

**Panel 3: Deploy Timeline**
```yaml
type: timeline
query: SELECT deploy_tag, COUNT(*) FROM traces GROUP BY deploy_tag ORDER BY time
```

#### Burn Alerts

```yaml
# Alert: High error rate in AI-generated code
name: AI Code Quality Degraded
description: "Error rate in AI-generated code exceeds 5%"
query: |
  SELECT COUNT_IF(error=true)*100.0/COUNT(*)
  FROM traces
  WHERE sdp.ai_generated = true
threshold: 5
window: 5m
```

#### Integration via OTLP

```go
// Honeycomb OTLP endpoint
endpoint := "https://api.honeycomb.io:443"
headers := map[string]string{
    "x-honeycomb-team": os.Getenv("HONEYCOMB_API_KEY"),
    "x-honeycomb-dataset": "production-traces",
}

exporter, _ := otlptracehttp.New(ctx,
    otlptracehttp.WithEndpoint(endpoint),
    otlptracehttp.WithHeaders(headers),
)

tp := trace.NewTracerProvider(
    trace.WithBatcher(exporter),
    trace.WithResource(resource.NewWithAttributes(
        semconv.SchemaURL,
        semconv.ServiceName("myapp"),
        attribute.String("sdp.version", "0.9.0"),
    )),
)

// Add SDP provenance processor
processor := NewSDPProvenanceProcessor(".sdp/log/events.jsonl")
tp.RegisterSpanProcessor(processor)
```

### Datadog

#### Overview

Datadog integrates with SDP via trace tags and deployment tracking.

#### Configuration

```yaml
# datadog-sdp.yaml
datadog_api_key: ${DATADOG_API_KEY}
datadog_app_key: ${DATADOG_APP_KEY}
site: datadoghq.com

# SDP integration
sdp_provenance:
  enabled: true
  trace_tags:
    - sdp.ai_generated
    - sdp.model
    - sdp.evidence_id
    - sdp.workstream_id

deploy_tracking:
  enabled: true
  git_tags: deploy/production/*
```

#### Trace Tags (Datadog convention)

Datadog uses tags instead of span attributes. Convert SDP attributes to tags:

```
# SDP span attribute -> Datadog tag
sdp.ai_generated:true        -> sdp.ai_generated:true
sdp.model:claude-sonnet-4 -> sdp.model:claude-sonnet-4
sdp.evidence_id:evt-123     -> sdp.evidence_id:evt-123
sdp.workstream_id:00-059-01 -> sdp.workstream_id:00-059-01
```

#### APM Queries

**AI code error rate:**
```
error:true AND sdp.ai_generated:true
```

**Specific model traces:**
```
sdp.model:claude-sonnet-4-20250514
```

**Evidence link:**
```
sdp.evidence_id:evt-123e4567-e89b-12d3
```

#### Deployment Tracking

```python
from datadog import api

# Create deploy event from deploy marker
tag = "deploy/production/20260211-120000"
msg = git_tag_message(tag)

api.Event.create(
    title=f"SDP Deploy: {tag}",
    text=msg,
    tags=[
        "environment:production",
        f"sdp.version:{msg['sdp_version']}",
        f"commit:{msg['commit_sha']}",
        f"workstreams:{','.join(msg['workstreams'])}",
    ],
    alert_type="info"
)
```

#### Dashboard Widgets

**Widget 1: Error Rate by Model**
```json
{
  "viz": "timeseries",
  "requests": [
    {
      "q": "sum:trace.errors{env:prod,sdp.ai_generated:true} by {sdp.model}.as_count()"
    }
  ]
}
```

**Widget 2: Workstream Breakdown**
```json
{
  "viz": "piechart",
  "requests": [
    {
      "q": "sum:trace.requests{env:prod} by {sdp.workstream_id}.as_count()"
    }
  ]
}
```

#### Monitors

```python
# High error rate in AI-generated code
api.Monitor.create(
    type="trace-analytics alert",
    query="error-rate{env:prod,sdp.ai_generated:true} > 5",
    name="SDP: AI Code Error Rate High",
    message="Error rate in AI-generated code exceeds 5%. Check evidence log.",
    tags=["sdp", "production"],
    options={
        "thresholds": {"critical": 5.0},
        "notify_audit": True
    }
)
```

### Grafana

#### Overview

Grafana provides unified dashboards combining traces (Tempo), metrics (Prometheus), and logs (Loki).

#### Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Grafana     │────▶│   Tempo      │────▶│ Evidence Log │
│  Dashboards  │     │  (Traces)    │     │   (Provenance)│
└──────────────┘     └──────────────┘     └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐
│  Prometheus  │    │    Loki      │
│  (Metrics)   │    │   (Logs)     │
└──────────────┘     └──────────────┘
```

#### Tempo Trace Query

**Find AI-generated error traces:**
```
{
  "query": "error=true AND sdp.ai_generated=true",
  "minDuration": "100ms",
  "maxDuration": "5s"
}
```

**Jump to evidence:**
```
Tempo span with sdp.evidence_id=evt-123
  -> External link to evidence log viewer
  -> https://evidence.example.com/evt/123e4567-e89b-12d3
```

#### Dashboard Panels

**Panel 1: Error Rate (AI vs Human)**
```json
{
  "title": "Error Rate by Author Type",
  "type": "graph",
  "targets": [
    {
      "expr": "sum(rate(traces_span_errors{env=\"prod\",sdp_ai_generated=\"true\"}[5m]))",
      "legendFormat": "AI-generated"
    },
    {
      "expr": "sum(rate(traces_span_errors{env=\"prod\",sdp_ai_generated=\"false\"}[5m]))",
      "legendFormat": "Human-written"
    }
  ]
}
```

**Panel 2: Model Performance**
```json
{
  "title": "Latency by Model (p95)",
  "type": "graph",
  "targets": [
    {
      "expr": "histogram_quantile(0.95, sum(rate(traces_span_latency_bucket{sdp_ai_generated=\"true\"}[5m])) by (le, sdp_model))",
      "legendFormat": "{{sdp_model}}"
    }
  ]
}
```

**Panel 3: Deploy Timeline**
```json
{
  "title": "Deploys & Error Spikes",
  "type": "graph",
  "targets": [
    {
      "expr": "sum(rate(traces_span_errors{env=\"prod\"}[5m]))",
      "legendFormat": "Error Rate"
    },
    {
      "expr": "sum(increments(sdp_deploy_total{env=\"prod\"}[1h]))",
      "legendFormat": "Deploys"
    }
  ],
  "annotations": {
    "list": [
      {
        "datasource": "SDP Evidence",
        "query": "deploy",
        "enable": true
      }
    ]
  }
}
```

#### Loki Log Queries

**Find logs from AI-generated code:**
```
{app="myapp", sdp_ai_generated="true"} |= "error"
```

**Correlate with traces:**
```
{trace_id="${__trace.id}"}
```

#### Data Source Configuration

```yaml
# Tempo data source with SDP attributes
apiVersion: 1

datasources:
  - name: Tempo
    type: tempo
    url: http://tempo:4318
    jsonData:
      tracesToLogs:
        datasourceUid: Loki
        filterByTraceID: true
        mapTagNamesEnabled: true
        mappedTags:
          - key: sdp.evidence_id
            value: sdp_evidence_id
      tracesToMetrics:
        datasourceUid: Prometheus
      nodeGraph:
        enabled: true
      search:
        filters:
          - name "AI-generated"
            tag: sdp.ai_generated
            value: true
```

---

## Data Flow Examples

### Example 1: Full Deployment Flow

```bash
# Step 1: Developer generates code (AI)
@build 00-059-01
# Evidence event written to .sdp/log/events.jsonl
# {
#   "id": "evt-123e4567-e89b-12d3",
#   "type": "generation",
#   "ws_id": "00-059-01",
#   "commit_sha": "abc123def456",
#   "data": {
#     "model": "claude-sonnet-4-20250514",
#     "files_changed": ["src/auth/login.go"]
#   }
# }

# Step 2: CI/CD deploys to production
git push origin main

# GitHub Actions workflow:
# - Checkout code
# - Run tests (sdp acceptance run)
# - Create deploy marker
TIMESTAMP=$(date -u +"%Y%m%d-%H%M%S")
git tag -a -m "$(cat deploy-marker.json)" "deploy/production/${TIMESTAMP}"
git push origin "deploy/production/${TIMESTAMP}"

# Step 3: SDP exporter computes diff provenance
sdp provenance export --tag "deploy/production/${TIMESTAMP}" > provenance.json

# Step 4: SDP exporter sends deploy event to observability tools
# Honeycomb
curl -X POST https://api.honeycomb.io/1/markers/deploy \
  -H "X-Honeycomb-Team: ${HC_API_KEY}" \
  -d @deploy-marker.json

# Datadog
curl -X POST https://api.datadoghq.com/api/v1/events \
  -H "DD-API-KEY: ${DD_API_KEY}" \
  -d @deploy-event.json

# Step 5: Application starts, emits OTel traces with SDP attributes
# (Automatic via SDPProvenanceProcessor)
```

### Example 2: Incident Investigation

```bash
# Alert: High error rate on /login endpoint (10% errors)

# Step 1: SRE opens Honeycomb, filters for errors
sdp.ai_generated:true AND error:true AND name:login

# Step 2: Sees span with:
# - sdp.evidence_id: evt-123e4567-e89b-12d3
# - sdp.model: claude-sonnet-4-20250514
# - sdp.workstream_id: 00-059-01

# Step 3: Clicks evidence_id (links to evidence log viewer)
# Opens: https://evidence.example.com/evt/123e4567-e89b-12d3

# Step 4: Views full evidence chain
sdp log trace --evt evt-123e4567-e89b-12d3

# Output:
# Event 1: Plan
#   Type: plan
#   WS: 00-059-01
#   Summary: "Add OAuth login flow"
#
# Event 2: Generation
#   Type: generation
#   WS: 00-059-01
#   Model: claude-sonnet-4-20250514
#   Files: src/auth/login.go (lines 15-25)
#
# Event 3: Verification
#   Type: verification
#   WS: 00-059-01
#   Coverage: 92%
#   Gates: test (passed), typecheck (passed), lint (passed)
#
# Event 4: Approval
#   Type: approval
#   WS: 00-059-01
#   Approved by: Jane Developer

# Step 5: SRE identifies root cause
# Line 16 of login.go (AI-generated) doesn't handle nil user from database
# Verification tests didn't cover nil return from DB

# Step 6: Human hotfix
git checkout -b hotfix/login-nil-check
# Add nil check on line 16
git commit -m "hotfix: handle nil user in login (P0)"
git push origin hotfix/login-nil-check

# Deploy hotfix via expedited flow
```

### Example 3: Quarterly Audit

```bash
# Requirement: "What % of production code is AI-generated?"

# Step 1: Generate provenance report for current deploy
sdp provenance report --deploy deploy/production/20260211-120000

# Output:
# SDP Provenance Report
# Deploy: 20260211-120000 (abc123def456)
# Generated: 2026-02-11T12:00:00Z
#
# Summary:
# - Total files: 250
# - Total lines: 45,000
# - AI-generated lines: 28,500 (63.3%)
# - Human-written lines: 16,500 (36.7%)
#
# By Model:
# - claude-sonnet-4-20250514: 20,000 lines (44.4%)
# - claude-opus-4-20250514: 8,500 lines (18.9%)
#
# By Feature:
# - F059 (Observability): 12,000 lines (26.7%)
# - F058 (Worktrees): 8,500 lines (18.9%)
# - F054 (Evidence): 8,000 lines (17.8%)
#
# Coverage by Author Type:
# - AI-generated: 86.5% average coverage
# - Human-written: 93.2% average coverage
#
# Low Coverage Files (AI-generated):
# - src/payment/processor.go: 65% (13/20 lines)
#   Workstream: 00-060-05
#   Model: claude-sonnet-4-20250514

# Step 2: Export to CSV for compliance
sdp provenance report --deploy deploy/production/20260211-120000 --format csv > audit-2026q1.csv
```

---

## Privacy & Security

### Data Classification

| Data | Classification | Storage | Access |
|------|---------------|---------|--------|
| Deploy markers (metadata) | Public | Git repo | All repo readers |
| Evidence log (event metadata) | Internal | Git repo | All repo readers |
| Prompt hashes | Internal | Evidence log | All repo readers |
| Trace attributes (sdp.*) | Internal | Observability tools | SRE team |
| Provenance maps | Internal | CI/CD artifacts | SRE team |

### No PII in SDP Attributes

**Guarantee:** SDP span attributes contain NO personally identifiable information.

**Safe attributes:**
- `sdp.ai_generated` (boolean)
- `sdp.model` (model name, no user data)
- `sdp.evidence_id` (UUID)
- `sdp.workstream_id` (WS ID format)
- `sdp.verification_status` (enum)

**Excluded attributes:**
- Git author email (use username-only `sdp.human_author` format)
- File paths (use relative paths only)
- Branch names (may contain ticket numbers with user info)

### Evidence Log Privacy

Evidence log is **already privacy-conscious** (see `docs/compliance/COMPLIANCE.md`):

- No raw prompts stored (only hashes)
- No user emails (usernames allowed)
- No code content (only file paths)
- Hash chain provides integrity detection (not tamper-proof)

### Observability Tool Access Control

**Recommendation:** Configure RBAC in observability tools:

- **SRE team**: Full access to traces + SDP attributes
- **Developers**: Read access to own workstream's traces
- **Auditors**: Read access to provenance reports only

**Honeycomb example:**
```yaml
# Team: SRE (full access)
permissions:
  - dataset: production-traces
    filters: []
    actions: [read, write, query]

# Team: Backend Developers (filtered access)
permissions:
  - dataset: production-traces
    filters: ["sdp.workstream_id:00-059-*"]
    actions: [read, query]

# Team: Auditors (provenance only)
permissions:
  - dataset: provenance-reports
    filters: []
    actions: [read]
```

### Compliance Alignment

| Regulation | SDP Compliance |
|------------|----------------|
| **SOC2** | Evidence log provides audit trail for code changes and approvals. Deploy markers support change management controls. |
| **HIPAA** | No PHI in SDP attributes or evidence log by design. Customer must ensure no PHI in traces (e.g., request payloads). |
| **DORA** | Evidence + deploy markers document ICT changes, testing, and deployment for incident response. |
| **EU AI Act** | Provenance attributes support transparency obligations for AI-generated code. |

---

## Implementation Roadmap

### Phase 1: Deploy Markers (P0)

**Workstream:** 00-059-03

**Deliverables:**
- [ ] Deploy marker schema (`schema/deploy-marker.schema.json`)
- [ ] Tag creation script (`hooks/create-deploy-marker.sh`)
- [ ] Tag query CLI (`sdp deploy show`)
- [ ] CI/CD integration examples (GitHub Actions, GitLab CI)

**Timeline:** 1 week

**Acceptance Criteria:**
- Deploy marker created on every production deploy
- Marker contains all required fields (commit_sha, evidence_log, workstreams)
- Tag can be queried via `sdp deploy show`

### Phase 2: OTel Attributes (P1)

**Workstream:** 00-059-02

**Deliverables:**
- [ ] OTel semantic convention draft (`docs/design/OTEL-AI-PROVENANCE.md`)
- [ ] SDPProvenanceProcessor implementation (Go)
- [ ] SDK examples (Python, Go, TypeScript)
- [ ] Unit tests for attribute injection

**Timeline:** 1.5 weeks

**Acceptance Criteria:**
- Spans include `sdp.ai_generated`, `sdp.model`, `sdp.evidence_id`
- Processor performance <1ms overhead per span
- Works with OTel SDK v1.20+

### Phase 3: Diff Provenance (P1)

**Workstream:** 00-059-04

**Deliverables:**
- [ ] Provenance calculator (`sdp provenance diff`)
- [ ] Git blame + evidence lookup engine
- [ ] JSON output schema
- [ ] Coverage report generator

**Timeline:** 2 weeks

**Acceptance Criteria:**
- Compute provenance for 1000-file diff in <5s
- Attribute accuracy: 99%+ (verified manually on sample)
- Output format matches schema

### Phase 4: Honeycomb Integration (P2)

**Workstream:** 00-059-05

**Deliverables:**
- [ ] Honeycomb exporter (`sdp export honeycomb`)
- [ ] Deploy event sender
- [ ] Query examples library
- [ ] Dashboard templates

**Timeline:** 1 week

**Acceptance Criteria:**
- Deploy events appear in Honeycomb within 30s of tag push
- Spans with SDP attributes queryable
- Dashboard templates work out-of-the-box

### Phase 5: Datadog Integration (P2)

**Workstream:** 00-059-06

**Deliverables:**
- [ ] Datadog event sender
- [ ] Trace tag mapper (attributes → tags)
- [ ] Monitor templates
- [ ] Dashboard widgets

**Timeline:** 1 week

**Acceptance Criteria:**
- Deploy events appear in Datadog stream
- Trace tags include SDP attributes
- Monitors trigger on AI code quality degradation

### Phase 6: Grafana Integration (P2)

**Workstream:** 00-059-07

**Deliverables:**
- [ ] Tempo data source config
- [ ] Dashboard JSON files
- [ ] Loki log queries
- [ - Evidence log external link plugin

**Timeline:** 1 week

**Acceptance Criteria:**
- Traces searchable by SDP attributes
- Dashboards display AI vs human metrics
- Evidence links open in log viewer

---

## Appendix

### A. Deploy Marker Examples

**Example 1: Production Deploy**
```json
{
  "sdp_version": "0.9.0",
  "environment": "production",
  "commit_sha": "abc123def4567890abcdef1234567890abcdef12",
  "deployed_at": "2026-02-11T12:00:00Z",
  "evidence_log": ".sdp/log/events.jsonl",
  "workstreams": ["00-059-01", "00-059-02", "00-059-03"],
  "feature_id": "F059",
  "deployer": "github-actions",
  "verification": {
    "coverage": 87.5,
    "gates_passed": ["test", "typecheck", "lint", "security"]
  }
}
```

**Example 2: Staging Deploy**
```json
{
  "sdp_version": "0.9.0",
  "environment": "staging",
  "commit_sha": "def4567890abcdef1234567890abcdef12345678",
  "deployed_at": "2026-02-11T11:30:00Z",
  "evidence_log": ".sdp/log/events.jsonl",
  "workstreams": ["00-060-01"],
  "feature_id": "F060",
  "deployer": "jane-developer",
  "verification": {
    "coverage": 82.0,
    "gates_passed": ["test", "typecheck"]
  }
}
```

### B. Evidence Event Examples

**Generation Event**
```json
{
  "id": "evt-123e4567-e89b-12d3-a456-426614174000",
  "type": "generation",
  "timestamp": "2026-02-11T10:00:00Z",
  "ws_id": "00-059-01",
  "commit_sha": "abc123def4567890abcdef1234567890abcdef12",
  "prev_hash": "sha256-of-previous-event",
  "data": {
    "model": "claude-sonnet-4-20250514",
    "model_version": "20250514",
    "prompt_hash": "a1b2c3d4e5f6...",
    "files_changed": ["src/auth/login.go"]
  }
}
```

### C. Open Questions

1. **Should we include full git commit history in provenance?**
   - **Decision:** No, only direct blame. History is available via `git log`.

2. **Should evidence IDs be traceable back to specific prompts?**
   - **Decision:** No, privacy. Only store prompt hash, not content.

3. **Should we support partial attribution (AI + human on same line)?**
   - **Decision:** No, too complex. Line-level granularity is sufficient.

4. **Should deploy markers be signed?**
   - **Decision:** P0: No. P3: Yes, for compliance-grade non-repudiation.

### D. References

- **Evidence Schema:** [schema/evidence.schema.json](../schema/evidence.schema.json)
- **Evidence Compliance:** [docs/compliance/COMPLIANCE.md](../compliance/COMPLIANCE.md)
- **OTel Spec:** https://opentelemetry.io/docs/reference/specification/
- **Honeycomb Docs:** https://docs.honeycomb.io/
- **Datadog APM:** https://docs.datadoghq.com/tracing/
- **Grafana Tempo:** https://grafana.com/docs/tempo/latest/

---

**Document Version:** 1.0
**Last Updated:** 2026-02-11
**Status:** Ready for implementation (P2 workstreams)
