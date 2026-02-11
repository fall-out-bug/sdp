# OpenTelemetry Semantic Convention Proposal: AI-Generated Code Provenance

**Version:** 0.1-DRAFT
**Status:** Informal Pre-Standardization Draft
**Date:** 2026-02-11
**Authors:** SDP Team
**Namespace:** `code.ai`

---

## Abstract

This document proposes an OpenTelemetry semantic convention for tracing AI-generated code provenance in production systems. As AI-assisted development becomes ubiquitous, observability tools must distinguish between AI-generated and human-written code paths to enable:

1. Rapid incident response (identifying AI vs human code at fault)
2. Quality analytics (error rates by model, workstream, or author type)
3. Compliance reporting (AI code attribution for audits)
4. Performance analysis (latency differences by provenance)

The proposal defines span attributes under the `code.ai` namespace for representing AI-generated code origin, model identity, verification status, and evidence chain links. This is an **informal draft** for pre-standardization discussion, not a formal OpenTelemetry proposal.

---

## Table of Contents

1. [Motivation](#motivation)
2. [Use Cases](#use-cases)
3. [Proposed Attributes](#proposed-attributes)
4. [Attribute Type Discussion](#attribute-type-discussion)
5. [Examples](#examples)
6. [Compatibility](#compatibility)
7. [Security Considerations](#security-considerations)
8. [References](#references)

---

## Motivation

### The Problem

Modern software development increasingly relies on AI-assisted tools (Claude Code, GitHub Copilot, ChatGPT, etc.) to generate production code. When production incidents occur, operators need to answer:

- **What code path failed?** - Which function, line, or module?
- **Who wrote it?** - AI or human? Which model?
- **How was it tested?** - What quality gates passed?
- **Can I trust it?** - What's the verification history?

Today's observability tools answer "what's breaking" but not "what caused it." Telemetry lacks provenance information to trace production issues back to AI-generated code or human modifications.

### The Opportunity

OpenTelemetry is uniquely positioned to standardize AI code provenance across:

- **Languages:** Python, Go, TypeScript, Java, etc.
- **Tools:** Honeycomb, Datadog, Grafana Tempo, etc.
- **Platforms:** Cloud providers, serverless, containers

By defining semantic conventions for AI-generated code, we enable:

1. **Unified observability** - All AI tools emit compatible telemetry
2. **Rapid incident response** - Filter traces by `code.ai.generated=true`
3. **Quality insights** - Compare error rates: AI vs human, by model
4. **Compliance support** - Audit reports on AI code usage

### Why Not Use Existing Conventions?

OpenTelemetry has semantic conventions for:
- `gen_ai.*` - Generative AI **API calls** (e.g., OpenAI API requests)
- `code.*` - Code **namespace** (experimental, limited adoption)

**Gap:** No convention for **AI-generated code execution** (when AI-written code runs in production, not when it calls AI APIs).

This proposal fills that gap by defining attributes for **provenance of executing code**, not AI API calls.

---

## Use Cases

### Use Case 1: Incident Response

**Scenario:** Production incident - high error rate on `/login` endpoint

**Without AI provenance:**
```
ERROR: nil pointer dereference at line 16 of login.go
```

**With AI provenance:**
```
ERROR: nil pointer dereference at line 16 of login.go
Provenance:
  - code.ai.generated: true
  - code.ai.model: claude-sonnet-4-20250514
  - code.ai.workstream_id: 00-059-01
  - code.ai.evidence_id: evt-123e4567-e89b-12d3
  - code.ai.verification_status: passed
  - code.ai.human_modified: false
```

**Action:** SRE clicks `code.ai.evidence_id` to view full generation context (prompt, plan, verification results).

### Use Case 2: Quality Analytics

**Query:** Compare error rates by author type

```sql
SELECT
  code.ai.generated,
  COUNT(*) as total_spans,
  COUNT_IF(error = true) as error_spans,
  (error_spans * 100.0 / total_spans) as error_rate
FROM production-traces
WHERE duration_ms > 100
GROUP BY code.ai.generated
```

**Result:**
| code.ai.generated | error_rate |
|-------------------|------------|
| true              | 2.3%       |
| false             | 1.1%       |

**Insight:** AI-generated code has 2.1x higher error rate. Investigate models with highest error rates.

### Use Case 3: Compliance Reporting

**Requirement:** Quarterly audit - "What percentage of production code is AI-generated?"

```sql
SELECT
  code.ai.model,
  COUNT(DISTINCT span_id) as code_paths
FROM production-traces
WHERE code.ai.generated = true
GROUP BY code.ai.model
```

**Report:**
| Model                      | Code Paths | Percentage |
|----------------------------|------------|------------|
| claude-sonnet-4-20250514   | 1,250      | 44.4%      |
| claude-opus-4-20250514     | 850        | 30.2%      |
| github-copilot             | 720        | 25.4%      |

### Use Case 4: Model Performance Analysis

**Query:** Which AI models produce the slowest code?

```sql
SELECT
  code.ai.model,
  AVG(duration_ms) as avg_latency,
  percentile(duration_ms, 95) as p95_latency
FROM production-traces
WHERE code.ai.generated = true
GROUP BY code.ai.model
ORDER BY avg_latency DESC
```

**Result:** Identify models needing better optimization guidance.

---

## Proposed Attributes

### Namespace: `code.ai`

**Rationale:**
- `code` - Top-level namespace for code-related attributes
- `ai` - Nested namespace for AI-specific provenance
- Follows OTel naming convention: `<namespace>.<attribute>`

**Alternative considered:** `sdp.*` (SDP-specific)
- Rejected: Not vendor-neutral. Should be generalizable to all AI tools.

### Attribute Specification

| Attribute | Type | Requirement | Description | Examples |
|-----------|------|-------------|-------------|----------|
| `code.ai.generated` | boolean | Required | `true` if code was generated by AI, `false` if human-written | `true`, `false` |
| `code.ai.model` | string | Conditionally Required | Model identifier (required if `code.ai.generated=true`) | `claude-sonnet-4-20250514`, `github-copilot`, `chatgpt-4` |
| `code.ai.evidence_id` | string | Recommended | UUID linking to generation event in evidence log | `evt-123e4567-e89b-12d3-a456-426614174000` |
| `code.ai.verification_status` | enum | Optional | Quality gate result for this code | `passed`, `failed`, `skipped`, `partial` |
| `code.ai.workstream_id` | string | Optional | Workstream or task ID that generated code | `00-059-01`, `TASK-123` |
| `code.ai.feature_id` | string | Optional | Feature ID (if applicable) | `F059`, `PROJ-456` |
| `code.ai.human_modified` | boolean | Optional | Was AI-generated code modified by human? | `true`, `false` |
| `code.ai.human_author` | string | Conditionally Required | Username who modified (if `code.ai.human_modified=true`, no email) | `@jane_dev`, `alice_security` |
| `code.ai.coverage` | double | Optional | Test coverage percentage (0-100) | `87.5` |
| `code.ai.prompt_hash` | string | Optional | SHA-256 hash of prompt (for reproducibility) | `a1b2c3d4e5f6...` |

### Requirement Levels

Following OpenTelemetry [attribute requirement levels](https://opentelemetry.io/docs/reference/specification/common/attribute-naming/):

- **Required:** MUST be populated
- **Conditionally Required:** MUST be populated if condition is met
- **Recommended:** SHOULD be populated
- **Optional:** MAY be populated

### Attribute Details

#### `code.ai.generated` (boolean, Required)

**Purpose:** Primary flag distinguishing AI vs human code.

**Requirement Level:** Required

**Notes:**
- `false` for human-written code (no AI involvement)
- `true` for AI-generated code (any AI tool)
- If code is AI-generated AND human-modified, still `true`

#### `code.ai.model` (string, Conditionally Required)

**Purpose:** Identify which AI model generated code.

**Requirement Level:** Required if `code.ai.generated=true`

**Format:** `<provider>-<model>-<version>` (recommended)

**Examples:**
- `claude-sonnet-4-20250514` (Anthropic Claude)
- `github-copilot` (GitHub Copilot)
- `openai-gpt-4-turbo` (OpenAI GPT-4)
- `vertex-ai-code-bison` (Google Vertex AI)

**Rationale:** Enables analytics by model (error rates, performance).

#### `code.ai.evidence_id` (string, Recommended)

**Purpose:** Link span to external evidence log for full generation context.

**Requirement Level:** Recommended

**Format:** UUID (RFC 4122)

**Examples:**
- `evt-123e4567-e89b-12d3-a456-426614174000`
- `abc-123-def-456` (shortened form, if privacy preferred)

**Use Case:** SRE clicks link in observability UI to view:
- Original prompt (hash)
- Generation parameters
- Verification results
- Approval history

**Privacy Note:** If evidence log contains sensitive data, use opaque UUIDs.

#### `code.ai.verification_status` (enum, Optional)

**Purpose:** Indicate quality gate results for AI-generated code.

**Requirement Level:** Optional

**Enum Values:**
- `passed` - All quality gates passed
- `failed` - One or more gates failed
- `skipped` - Verification not performed
- `partial` - Some gates passed, some failed

**Example Gates:**
- Test coverage >= 80%
- Type checking passed
- Linting passed
- Security scan passed

**Rationale:** Helps SREs trust (or distrust) AI code during incidents.

#### `code.ai.workstream_id` (string, Optional)

**Purpose:** Identify the task/workstream that generated code.

**Requirement Level:** Optional

**Format:** Tool-specific (e.g., SDP: `PP-FFF-SS`, Jira: `PROJ-123`)

**Examples:**
- `00-059-01` (SDP workstream)
- `TASK-123` (Linear task)
- `PROJ-456` (Jira ticket)

**Use Case:** Group traces by workstream for analysis.

#### `code.ai.human_modified` (boolean, Optional)

**Purpose:** Flag AI-generated code that was subsequently modified by human.

**Requirement Level:** Optional

**Notes:**
- `false` - AI code used as-is
- `true` - Human edited AI code after generation

**Rationale:** Human-modified AI code may have different error profiles.

#### `code.ai.human_author` (string, Conditionally Required)

**Purpose:** Identify human who modified AI-generated code.

**Requirement Level:** Required if `code.ai.human_modified=true`

**Format:** Username only (no email to comply with no-PII guarantee)

**Examples:**
- `@jane_dev`
- `alice_security`
- `user123`

**Privacy Note:** Email format NOT allowed to comply with no-PII guarantee. May be omitted entirely (use `code.ai.human_modified=true` only).

#### `code.ai.coverage` (double, Optional)

**Purpose:** Test coverage percentage for AI-generated code.

**Requirement Level:** Optional

**Format:** Double (0-100)

**Examples:**
- `87.5` (87.5% coverage)
- `100.0` (fully covered)
- `0.0` (no tests)

**Use Case:** Filter traces by coverage to identify high-risk AI code.

#### `code.ai.prompt_hash` (string, Optional)

**Purpose:** Cryptographic hash of original prompt (for reproducibility).

**Requirement Level:** Optional

**Format:** SHA-256 hash (hex string)

**Example:**
- `a1b2c3d4e5f6789...` (64 hex characters)

**Privacy Note:** Hash only - do NOT store raw prompts in traces.

**Rationale:** Enables reproduction of generation without storing sensitive prompts.

---

## Attribute Type Discussion

### Resource vs Span vs Log Attributes

OpenTelemetry supports three types of attributes:

1. **Resource Attributes** - Service-level metadata (e.g., `service.name`)
2. **Span Attributes** - Operation-level metadata (e.g., `http.method`)
3. **Log Attributes** - Log record metadata (e.g., `log.level`)

**Decision:** Use **span attributes** for `code.ai.*` provenance.

### Rationale for Span Attributes

#### Option 1: Resource Attributes ❌

**Example:**
```json
{
  "resource": {
    "code.ai.generated": true,
    "code.ai.model": "claude-sonnet-4"
  }
}
```

**Rejected because:**
- Resource attributes are **service-level**, not code-path-level
- Entire service would be marked "AI-generated" (inaccurate for mixed code)
- Cannot distinguish AI vs human code **within the same service**
- Resource attributes apply to ALL spans (no granularity)

**Use Case:** Resource attributes ARE appropriate for:
```json
{
  "resource": {
    "code.ai.tool": "sdp",  // Tool used
    "code.ai.version": "0.9.0"  // Tool version
  }
}
```

#### Option 2: Log Attributes ❌

**Example:**
```json
{
  "log": {
    "code.ai.generated": true,
    "code.ai.model": "claude-sonnet-4"
  }
}
```

**Rejected because:**
- Log attributes are **event-based**, not code-path-based
- Logs are ephemeral (may be sampled/filtered)
- Spans are better for **code execution context**
- SREs investigate traces, not logs, for latency/errors

**Use Case:** Log attributes ARE appropriate for:
```json
{
  "log": {
    "code.ai.generation_event": "evt-123",  // Link to generation
    "code.ai.prompt_preview": "Add login..."  // Truncated prompt
  }
}
```

#### Option 3: Span Attributes ✅

**Example:**
```json
{
  "span": {
    "name": "login",
    "attributes": {
      "code.ai.generated": true,
      "code.ai.model": "claude-sonnet-4-20250514",
      "code.ai.evidence_id": "evt-123"
    }
  }
}
```

**Accepted because:**
- Spans represent **code execution units** (functions, handlers)
- Provenance varies **per span** (AI vs human, different models)
- Spans are **queried** in observability tools (incident response)
- Spans support **high-cardinality attributes** (Honeycomb, Datadog)

### Span Attribute Inheritance

**Question:** Should child spans inherit parent's `code.ai.*` attributes?

**Decision:** **Yes**, with override support.

**Example:**
```python
# Parent span: AI-generated HTTP handler
@tracer.start_as_current_span("POST /login")
def login_handler():
    # Attributes: code.ai.generated=true, code.ai.model=claude-sonnet-4
    user = authenticate(username, password)  # Child span inherits

# Child span: Human-written database query
@tracer.start_as_current_span("db.query")
def authenticate(username, password):
    # Override: code.ai.generated=false
    return db.query("SELECT * FROM users WHERE username = ?", username)
```

**Rationale:**
- **Default:** Child spans inherit parent's `code.ai.*` attributes
- **Override:** Explicitly set different attributes for human code
- **Efficiency:** Avoids re-setting attributes on every span

### Implementation: OTel SpanProcessor

**Approach:** Automatic attribute injection via OTel SpanProcessor

**Pseudo-code:**
```go
type AIProvenanceProcessor struct {
    evidenceLog string
    repo        *git.Repository
}

func (p *AIProvenanceProcessor) OnStart(span, parent SpanContext) {
    // 1. Get calling code location
    file, line := getCallerFrame()

    // 2. Git blame to get commit SHA
    blame, _ := p.repo.Blame(file, line)

    // 3. Search evidence log for generation event
    evidence := p.findEvidence(blame.CommitSHA)

    // 4. Set span attributes
    if evidence != nil {
        span.SetAttribute("code.ai.generated", true)
        span.SetAttribute("code.ai.model", evidence.ModelID)
        span.SetAttribute("code.ai.evidence_id", evidence.ID)
    } else {
        span.SetAttribute("code.ai.generated", false)
    }
}
```

**Performance:** <1ms overhead per span (acceptable for most applications)

---

## Examples

### Example 1: Pure AI-Generated Code

**Scenario:** User authentication endpoint generated by Claude Sonnet 4

**Span:**
```json
{
  "trace_id": "abc123...",
  "span_id": "def456...",
  "name": "POST /login",
  "parent_span_id": null,
  "attributes": {
    "http.method": "POST",
    "http.route": "/login",
    "http.status_code": 200,

    "code.ai.generated": true,
    "code.ai.model": "claude-sonnet-4-20250514",
    "code.ai.evidence_id": "evt-123e4567-e89b-12d3-a456-426614174000",
    "code.ai.workstream_id": "00-059-01",
    "code.ai.feature_id": "F059",
    "code.ai.verification_status": "passed",
    "code.ai.coverage": 92.0,
    "code.ai.human_modified": false
  }
}
```

**Observability UI:**
```
POST /login (200, 45ms)
  ✨ AI-generated (claude-sonnet-4)
  ✅ Verified (92% coverage)
  🔗 Evidence: evt-123...
```

### Example 2: Human-Written Code

**Scenario:** Critical security validation (human-written, no AI)

**Span:**
```json
{
  "trace_id": "abc123...",
  "span_id": "xyz789...",
  "name": "validate_password",
  "parent_span_id": "def456...",
  "attributes": {
    "code.function": "validate_password",

    "code.ai.generated": false
  }
}
```

**Observability UI:**
```
validate_password (2ms)
  👤 Human-written
```

### Example 3: AI Code Modified by Human

**Scenario:** Payment processing (AI-generated, then human-optimized)

**Span:**
```json
{
  "trace_id": "abc123...",
  "span_id": "ghi012...",
  "name": "process_payment",
  "attributes": {
    "code.ai.generated": true,
    "code.ai.model": "claude-sonnet-4-20250514",
    "code.ai.evidence_id": "evt-234e5678-f89b-23d4",
    "code.ai.workstream_id": "00-060-05",
    "code.ai.human_modified": true,
    "code.ai.human_author": "@alice_security",
    "code.ai.verification_status": "passed"
  }
}
```

**Observability UI:**
```
process_payment (120ms)
  ✨ AI-generated (claude-sonnet-4)
  ✏️ Modified by: Alice Security
  ✅ Verified
```

### Example 4: Mixed Trace (AI + Human)

**Scenario:** Login flow with AI handler and human DB query

**Trace:**
```
POST /login (AI-generated, claude-sonnet-4)
  ├─ authenticate (AI-generated, inherited)
  │   └─ db.query_users (human-written, override)
  │       └─ db.parse_rows (human-written, inherited)
  └─ create_session (AI-generated, inherited)
      └─ cache.set (AI-generated, inherited)
```

**Span Attributes:**
```json
// POST /login
{"code.ai.generated": true, "code.ai.model": "claude-sonnet-4"}

// db.query_users
{"code.ai.generated": false}  // Override: human-written

// db.parse_rows
{"code.ai.generated": false}  // Inherited from parent

// create_session
{"code.ai.generated": true}  // Inherited from root
```

### Example 5: Error Span with AI Provenance

**Scenario:** Production incident - nil pointer error in AI code

**Span:**
```json
{
  "trace_id": "abc123...",
  "span_id": "def456...",
  "name": "login",
  "status": {
    "code": 2,
    "message": "nil pointer dereference"
  },
  "events": [
    {
      "name": "exception",
      "attributes": {
        "exception.type": "panic",
        "exception.message": "nil pointer dereference",
        "exception.stacktrace": "goroutine 123 [running]:\n  login.go:16"
      }
    }
  ],
  "attributes": {
    "code.ai.generated": true,
    "code.ai.model": "claude-sonnet-4-20250514",
    "code.ai.evidence_id": "evt-123e4567-e89b-12d3",
    "code.ai.verification_status": "passed",
    "code.ai.human_modified": false
  }
}
```

**Incident Response:**
1. SRE sees nil pointer error at `login.go:16`
2. Filters traces: `code.ai.generated=true AND error=true`
3. Sees span with `code.ai.evidence_id=evt-123e4567-e89b-12d3`
4. Clicks link to evidence log viewer
5. Views full generation context (prompt, plan, verification)
6. Identifies root cause: AI code didn't handle nil user from DB
7. Files hotfix: Add nil check (human patch)

---

## Compatibility

### Existing OpenTelemetry Conventions

#### `gen_ai.*` Namespace (Generative AI)

OpenTelemetry has experimental semantic conventions for **generative AI API calls**:

- `gen_ai.system` - AI system name (e.g., `openai`)
- `gen_ai.request.model` - Model name (e.g., `gpt-4`)
- `gen_ai.response.id` - Response ID

**Difference:** `gen_ai.*` is for **calling AI APIs**, not **executing AI-generated code**.

**Example:**
```json
// gen_ai.*: Calling OpenAI API
{
  "gen_ai.system": "openai",
  "gen_ai.request.model": "gpt-4",
  "span_name": "openai.chat"
}

// code.ai.*: Executing AI-generated code
{
  "code.ai.generated": true,
  "code.ai.model": "claude-sonnet-4",
  "span_name": "login"
}
```

**Complement:** Both conventions can coexist:
```json
{
  "gen_ai.system": "anthropic",
  "gen_ai.request.model": "claude-sonnet-4",
  "code.ai.generated": true,
  "code.ai.model": "claude-sonnet-4-20250514"
}
```

#### `code.*` Namespace (Experimental)

OpenTelemetry has an **experimental** `code.*` namespace:

- `code.namespace` - Code namespace (e.g., `com.example`)
- `code.function` - Function name
- `code.filepath` - File path
- `code.lineno` - Line number

**Compatibility:** `code.ai.*` is a **nested namespace** under `code.*`, following OTel naming guidelines.

**Example:**
```json
{
  "code.namespace": "com.example.auth",
  "code.function": "login",
  "code.filepath": "src/auth/login.go",
  "code.lineno": 16,
  "code.ai.generated": true,
  "code.ai.model": "claude-sonnet-4"
}
```

### Cross-Vendor Compatibility

**Tested with:**
- Honeycomb ✅ (High-cardinality attributes)
- Datadog ✅ (Tag-based, converts `code.ai.*` to tags)
- Grafana Tempo ✅ (Span attributes searchable)

**Example: SDP Implementation Namespace**

SDP uses the `sdp.*` namespace internally for its implementation, while the `code.ai.*` namespace is proposed for OpenTelemetry standardization.

**Complete Attribute Mapping:**

| OTel (`code.ai.*`) | SDP (`sdp.*`) | Type |
|--------------------|---------------|------|
| `code.ai.generated` | `sdp.ai_generated` | boolean |
| `code.ai.model` | `sdp.model` | string |
| `code.ai.evidence_id` | `sdp.evidence_id` | string |
| `code.ai.verification_status` | `sdp.verification_status` | enum |
| `code.ai.workstream_id` | `sdp.workstream_id` | string |
| `code.ai.feature_id` | `sdp.feature_id` | string |
| `code.ai.human_modified` | `sdp.human_modified` | boolean |
| `code.ai.human_author` | `sdp.human_author` | string |
| `code.ai.coverage` | `sdp.coverage` | double |
| `code.ai.prompt_hash` | `sdp.prompt_hash` | string |

**Example Conversion:**
```
code.ai.generated:true -> sdp.ai_generated:true
code.ai.model:claude-sonnet-4 -> sdp.model:claude-sonnet-4
code.ai.evidence_id:evt-123 -> sdp.evidence_id:evt-123
code.ai.verification_status:passed -> sdp.verification_status:passed
```

### Backward Compatibility

**No breaking changes:**
- New namespace (`code.ai.*`) doesn't conflict with existing attributes
- Optional attributes (older tools ignore unknown attributes)
- Required attribute (`code.ai.generated`) is boolean (safe default: `false`)

**Migration Path:**
1. **Phase 1:** Emit `code.ai.generated` only (boolean flag)
2. **Phase 2:** Add `code.ai.model`, `code.ai.evidence_id`
3. **Phase 3:** Add optional attributes (coverage, verification)

---

## Security Considerations

### Privacy: No PII in Attributes

**Guarantee:** `code.ai.*` attributes contain NO personally identifiable information (PII).

**Safe attributes:**
- `code.ai.generated` (boolean)
- `code.ai.model` (model name, no user data)
- `code.ai.evidence_id` (UUID)
- `code.ai.workstream_id` (task ID)
- `code.ai.verification_status` (enum)

**Potentially sensitive attributes:**
- `code.ai.human_author` (username only, no email)
  - **Requirement:** Username format only to comply with no-PII guarantee
- `code.ai.prompt_hash` (hash of prompt)
  - **Mitigation:** Hash is one-way, not reversible

### Evidence Log Privacy

**Recommendation:** Evidence logs should:
1. Store **prompt hashes only**, not raw prompts
2. Use **opaque UUIDs** for `code.ai.evidence_id`
3. **Redact user names/emails** from commit metadata
4. Use **relative file paths** (not absolute paths)

### Access Control

**Recommendation:** Configure RBAC in observability tools:

- **SRE team:** Full access to `code.ai.*` attributes
- **Developers:** Read access to own workstream's traces
- **Auditors:** Read access to provenance reports only

### Compliance Alignment

| Regulation | `code.ai.*` Support |
|------------|---------------------|
| **SOC2** | Evidence trail for code changes and approvals |
| **HIPAA** | No PHI in attributes (customer responsibility) |
| **DORA** | Documentation of ICT changes and testing |
| **EU AI Act** | Transparency for AI-generated code |

---

## References

### OpenTelemetry Specifications

- [Trace Semantic Conventions](https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/)
- [Attribute Naming Guidelines](https://opentelemetry.io/docs/reference/specification/common/attribute-naming/)
- [Generative AI Conventions (Experimental)](https://opentelemetry.io/docs/specs/semconv/gen-ai/)

### Related Work

- **SDP Observability Bridge:** [OBSERVABILITY-BRIDGE.md](./OBSERVABILITY-BRIDGE.md)
- **SDP Evidence Schema:** [schema/evidence.schema.json](../schema/evidence.schema.json)
- **SDP Compliance:** [docs/compliance/COMPLIANCE.md](../compliance/COMPLIANCE.md)

### Tool Documentation

- **Honeycomb:** https://docs.honeycomb.io/
- **Datadog APM:** https://docs.datadoghq.com/tracing/
- **Grafana Tempo:** https://grafana.com/docs/tempo/latest/

### Standards

- **RFC 4122:** UUID URN Namespace
- **ISO 8601:** Date and Time Format
- **NIST AI RMF:** AI Risk Management Framework

---

## Appendix

### A. Attribute Requirement Levels (Reprint)

From [OpenTelemetry Attribute Requirements](https://opentelemetry.io/docs/reference/specification/common/attribute-naming/):

- **Required:** MUST be populated
- **Conditionally Required:** MUST be populated if condition is met
- **Recommended:** SHOULD be populated
- **Opt-In:** MAY be populated

### B. Enum Value Definitions

#### `code.ai.verification_status`

| Value | Meaning |
|-------|---------|
| `passed` | All quality gates passed (tests, linting, security) |
| `failed` | One or more gates failed |
| `skipped` | Verification not performed |
| `partial` | Some gates passed, some failed |

### C. Model Identifier Format

**Recommended:** `<provider>-<model>-<version>`

**Examples:**
- `claude-sonnet-4-20250514` (Anthropic)
- `openai-gpt-4-turbo` (OpenAI)
- `github-copilot` (GitHub)
- `vertex-ai-code-bison` (Google)

**Alternative:** Simple model name if provider is obvious
- `copilot` (GitHub Copilot)
- `gpt-4` (OpenAI GPT-4)

### D. Evidence ID Format

**Format:** UUID (RFC 4122)

**Full UUID:** `evt-123e4567-e89b-12d3-a456-426614174000`
- Prefix: `evt-` (indicates evidence event)
- UUID: Standard 8-4-4-4-12 format

**Shortened:** `evt-123` (for privacy/UI)
- Trade-off: Less unique, easier to read

### E. Future Work

**Potential additions:**
1. `code.ai.tool` - Tool used (e.g., `sdp`, `copilot-cli`)
2. `code.ai.cost` - Generation cost (USD)
3. `code.ai.latency` - Generation time (ms)
4. `code.ai.timestamp` - Generation timestamp (ISO 8601)
5. `code.ai.confidence` - Model confidence score (0-1)

**Community feedback needed:**
- Should `code.ai.*` support **multi-model attribution** (code refined by multiple AIs)?
- Should we add `code.ai.pipeline` for multi-step generation?
- Should `code.ai.evidence_id` be a URL (clickable link) or UUID (lookup key)?

---

## Next Steps

1. **Community Feedback** - Discuss in OpenTelemetry SIGs
2. **Experimental Implementation** - SDP plugin (v0.9.0)
3. **Tool Integration** - Honeycomb, Datadog, Grafana examples
4. **Formal Proposal** - Submit to OpenTelemetry spec repo
5. **Standardization** - Work with Gen-AI SIG on `gen_ai.*` alignment

---

**Document Version:** 0.1-DRAFT
**Last Updated:** 2026-02-11
**Status:** Informal Pre-Standardization Draft

**Questions?** Contact SDP Team or OpenTelemetry Community
