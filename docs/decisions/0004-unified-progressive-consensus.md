# ADR-0004: Unified Progressive Consensus Protocol v2.0

## Status

Accepted

## Context

This repository targets workflows from small student projects to large multi-week epics. Four competing proposals for Consensus Protocol v2.0 each addressed different aspects:

| Proposal | Focus | Strength | Weakness |
|----------|-------|----------|----------|
| 0001-Cursor (Gemini) | IDE integration | TypeScript DX, Context Anchors | Platform lock-in |
| 0001-File (Opus) | Portability | JSON Schema, standard files | Lacks DX polish |
| 0002-Reliability (GPT) | Correctness | "Schema is Law", validation gates | Strict, adds friction |
| 0003-Progressive (Composer) | Adoption | Starter/Standard/Advanced tiers | Complex implementation |

### Problems Identified Across All Proposals

1. **Onboarding barrier**: Full protocol complexity overwhelms beginners
2. **Platform fragmentation**: IDE-specific features create lock-in
3. **Validation gaps**: No runtime validation or unclear error messages
4. **State ambiguity**: Implicit state requires reading multiple files
5. **Execution gap**: High-level plans don't translate to atomic tasks

### Requirements for Unified Protocol

The protocol MUST be:

- **Universal**: Work with ANY LLM and ANY coding tool
- **Progressive**: Scale from student projects to enterprise epics
- **Reliable**: Runtime-validatable with clear error messages
- **Efficient**: Leverage IDE capabilities when available (optional)
- **Human-in-the-loop**: Preserve intervention points at every phase

## Decision

We adopt the **Unified Progressive Consensus (UPC) Protocol** that merges:

- **Progressive Tiers** from 0003-Composer (adoption)
- **File-Native Core** from 0001-Opus/0002-GPT (reliability)
- **Micro-Tasking** from 0001-Gemini (execution)
- **DX Adapters** as optional layer (IDE efficiency)

### Core Principle

```
Files are the only required interface.
Everything else is optional enhancement.
```

## Detailed Design

### 1. Progressive Tiers

The protocol adapts to project complexity:

| Tier | Name | Use Case | Requirements |
|------|------|----------|--------------|
| 1 | **Starter** | Bug fixes, prototypes, student projects | `epic.md` + minimal `status.json` |
| 2 | **Standard** | Features, typical development | Full schemas + validation |
| 3 | **Enterprise** | Large systems, multi-team | Custom roles, cross-epic coordination |

#### Tier Selection Criteria

Tier selection is **explicit and human-chosen** (set in `status.json.tier`).
The following heuristics are **guidance only** (not enforced automatically):

```
IF lines_changed < 50 AND no_new_dependencies AND single_developer:
    tier = "starter"
ELIF architectural_change OR team_size > 3 OR multi_epic:
    tier = "enterprise"
ELSE:
    tier = "standard"  # Safe default
```

#### Tier Capabilities Matrix

| Feature | Starter | Standard | Enterprise |
|---------|---------|----------|------------|
| `epic.md` | Required | Required | Required |
| `status.json` | Minimal | Full | Full + extensions |
| JSON Schema validation | Optional | Required | Required + custom |
| Full agent chain | Optional | Required | Required |
| Cross-epic refs | No | No | Yes |
| Custom roles | No | No | Yes |

### 2. Schema-Driven Core ("Schema is Law")

JSON Schema is the canonical contract for all structured data.

#### Schema Registry

```
consensus/schema/
├── status.schema.json       # Epic state machine
├── message.schema.json      # Agent communication
├── requirements.schema.json # Analyst output
├── architecture.schema.json # Architect output
├── plan.schema.json         # Tech Lead output
└── index.json               # Schema versions and metadata
```

#### Canonical Schema Location

- Schemas live **once** in the repository: `consensus/schema/`.
- Epics **reference schemas** via relative paths in `$schema` field:
  - `"$schema": "../consensus/schema/status.schema.json"`
  - `"$schema": "../consensus/schema/message.schema.json"`
- Validators resolve relative paths from the epic's `consensus/` directory.
- This approach works universally without custom URI resolvers.

#### Why JSON Schema (not TypeScript)

| Aspect | JSON Schema | TypeScript |
|--------|-------------|------------|
| Runtime validation | ✅ Native | ❌ Compile-only |
| LLM understanding | ✅ Excellent | ✅ Good |
| Tool requirements | ✅ None | ❌ tsc/node |
| Cross-language | ✅ Universal | ❌ JS/TS only |
| Type generation | ✅ Can generate TS | — |

**Rule**: JSON Schema is the source of truth. TypeScript types MAY be generated from schemas for IDE convenience but are NOT required.

### 2.1 Canonical Artifacts (JSON-first)

For **Standard tier** and higher, the protocol is **JSON-first**:

- Canonical (machine) artifacts:
  - `consensus/status.json`
  - `consensus/artifacts/requirements.json`
  - `consensus/artifacts/architecture.json`
  - `consensus/artifacts/plan.json`
  - `consensus/artifacts/implementation.json`
  - `consensus/artifacts/test_results.json`
- Optional (human) documents:
  - `epic.md`, `architecture.md`, `implementation.md`, `testing.md`, `deployment.md`

**Rule**: phase gates and automation depend only on canonical JSON artifacts + inbox messages.

### 3. Centralized State (`status.json`)

Single mutable file controls the epic lifecycle.

#### Unified Status Schema

```json
{
  "$schema": "../consensus/schema/status.schema.json",
  "epic_id": "EP-AUTH-001",
  "tier": "standard",
  "phase": "implementation",
  "mode": "full",
  "iteration": 2,
  
  "approvals": ["analyst", "architect", "tech_lead"],
  "blockers": [],
  
  "workstreams": [
    {
      "id": "WS-01",
      "title": "Domain entities",
      "status": "done",
      "completed_at": "2025-12-31T10:00:00Z"
    },
    {
      "id": "WS-02", 
      "title": "Use cases",
      "status": "in_progress",
      "started_at": "2025-12-31T12:00:00Z"
    },
    {
      "id": "WS-03",
      "title": "API endpoints",
      "status": "todo"
    }
  ],
  
  "updated_at": "2025-12-31T14:30:00Z",
  "updated_by": "developer"
}
```

#### Field Origins

| Field | Source | Purpose |
|-------|--------|---------|
| `tier` | 0003-Composer | Progressive complexity |
| `mode` | 0001-Gemini | Execution topology (full/fast_track/hotfix) |
| `workstreams` | 0001-Gemini | Micro-tasking (merged from kanban.json) |
| `blockers` | 0001-Opus | Veto tracking |
| `approvals` | 0002-GPT | Phase gate tracking |

#### State Machine

```
                    ┌─────────────────────────────────────────┐
                    │              blocked                     │
                    │  (can transition to any previous phase)  │
                    └─────────────────────────────────────────┘
                                       ↑
    requirements → architecture → planning → implementation → testing → deployment → done
         │              │            │             │            │          │
         └──────────────┴────────────┴─────────────┴────────────┴──────────┘
                              (veto returns to relevant phase)
```

#### Ownership Rules

1. Agent MUST read `status.json` before acting
2. Agent MUST validate before writing
3. Agent MUST set `updated_at` and `updated_by`

#### Concurrency Rules

**Phase-level locking:**
- Only the phase owner role can advance `phase` (except moving to `blocked`)
- Other agents can update their artifacts but cannot change phase

**Workstream-level concurrency:**
- Multiple agents can work on different workstreams simultaneously
- Each agent updates only their assigned workstream(s) in `status.json`
- Workstream updates are independent and do not require phase-level locking
- Tech Lead merges workstream updates before phase transition

**Conflict resolution:**
- Last write wins (determined by `updated_at` timestamp)
- Validator checks for conflicts before phase transition
- Human intervention required if conflicts detected during phase transition

**Workstream limits:**
- If `workstreams.length > 50`, validator SHOULD warn about potential epic complexity
- This is a recommendation, not a blocker (quality epic planning should prevent this)

#### Optional Append-Only Audit Trail (`status_history.jsonl`)

To get practical append-only state auditing without heavy tooling:

- Keep `consensus/status.json` as the mutable “current state”.
- Optionally maintain `consensus/status_history.jsonl` as an append-only event log.

Recommended event shape (illustrative):

```json
{"ts":"2025-12-31T14:30:00Z","by":"developer","from_phase":"planning","to_phase":"implementation","status_sha256":"..."}
```

If enabled, validators SHOULD check:

- JSONL lines are valid JSON
- timestamps are non-decreasing
- the last entry’s `status_sha256` matches the current `status.json` content hash

#### Tier vs Mode (Clarification)

- `tier` controls **how strict validation and required artifacts are**.
- `mode` controls **which roles/agents must participate** (topology).

### 4. Execution Modes (Dynamic Topology)

Three modes for different scenarios:

| Mode | Agent Chain | Use Case |
|------|-------------|----------|
| `full` | analyst → architect → tech_lead → developer → qa → devops | New features |
| `fast_track` | developer → qa | Bug fixes (≤50 LOC) |
| `hotfix` | developer → devops | Critical production issues |

Mode is set in `status.json.mode` and determines which agents must participate.

### 4.1 Executable Plan (Canonical)

`consensus/artifacts/plan.json` is the **executable plan**:

- Produced by Tech Lead
- Contains workstreams + atomic tasks + dependencies
- Maps tasks to requirements/architecture IDs
- Includes test commands and evidence expectations per task

`consensus/artifacts/implementation.json` is the **execution ledger**:

- Produced/updated by Developer
- Tracks per-task status and evidence pointers (test output paths, coverage snapshot, logs)

### 5. Agent Instructions (Portable Prompts)

Plain Markdown files that work with ANY LLM tool.

#### Location

```
consensus/prompts/
├── analyst.md
├── architect.md
├── tech_lead.md
├── developer.md
├── qa.md
└── devops.md
```

#### Universal Prompt Structure

```markdown
# {Role} Agent

## Phase Gate
You can only act when `status.json` shows:
- `phase`: "{allowed_phase}"
- `tier`: "{minimum_tier}" or higher

If conditions not met, STOP and report mismatch.

## Context Files
Read before starting:
1. `consensus/status.json` — current state
2. `consensus/artifacts/{dependency}.json` — previous deliverables
3. `consensus/messages/inbox/{role}/` — messages for you

## Your Deliverables
1. Create: `consensus/artifacts/{output}.json`
2. Messages: `consensus/messages/inbox/{target}/`
3. Update: `consensus/status.json`

## Validation
Before completing:
- Validate artifacts against `consensus/schema/`
- Ensure all text in JSON is English
- Run `./consensus/validate.sh` if available

## Completion Checklist
- [ ] Phase gate verified
- [ ] Inbox processed
- [ ] Artifacts created and validated
- [ ] Status.json updated
- [ ] Messages sent to next agent(s)
```

### 6. Platform Adapters (Optional DX Layer)

Adapters enhance experience on specific platforms but are NEVER required.

#### Adapter Directory

```
docs/guides/adapters/
├── README.md           # Overview and status
├── cursor.md           # ✅ Verified
├── claude-code.md      # ✅ Verified  
├── codex.md            # 🚧 Planned
├── aider.md            # 🚧 Planned
└── vscode-copilot.md   # 🚧 Planned
```

#### Adapter Implementation Strategy

| Platform | Config File | Mechanism |
|----------|-------------|-----------|
| Cursor | `.cursorrules` | Auto-loaded context rules |
| Claude Code | `CLAUDE.md` | Project instructions |
| Aider | `.aider.conf.yml` | Read files config |
| VSCode Copilot | `.github/copilot-instructions.md` | Workspace instructions |

#### Adapter Rules

1. **Additive only**: Adapters add convenience, never change semantics
2. **Optional**: Protocol works without any adapter
3. **Documented**: Each adapter has verified guide
4. **Isolated**: Adapter-specific files in `adapters/` or platform locations

### 7. Validation Strategy

#### Validation Levels by Tier

| Tier | Validation | Required |
|------|------------|----------|
| Starter | JSON syntax only | Recommended |
| Standard | JSON Schema + phase rules | Required |
| Enterprise | Schema + custom rules + cross-refs | Required |

#### Validation Gate Definition (Standard Tier)

Before an agent can advance `status.json.phase`, the following must hold:

1. **Schema gate**: produced/updated JSON validates against its `$schema`.
2. **Phase gate**: required canonical artifacts for the target phase exist (see Phase Artifact Requirements table below).
3. **Ownership gate**: only the phase owner role can advance phase (except moving to `blocked`).
4. **Inbox rules gate**: agents read only their own inbox; operational messages are written only to other roles.
5. **Language gate**:
   - User-facing docs (`epic.md`, narrative `.md`) may be **any language**
   - Protocol JSON (messages + canonical artifacts) MUST be **English-only** in text fields

#### Phase Artifact Requirements

| Phase | Required Artifacts | Optional Artifacts |
|-------|-------------------|-------------------|
| `requirements` | `requirements.json` | `requirements.md` |
| `architecture` | `requirements.json`, `architecture.json` | `architecture.md` |
| `planning` | `requirements.json`, `architecture.json`, `plan.json` | `implementation.md` |
| `implementation` | `plan.json`, `implementation.json` | - |
| `testing` | `implementation.json`, `test_results.json` | `testing.md` |
| `deployment` | `test_results.json`, `deployment.json` | `deployment.md` |
| `done` | All previous phase artifacts | - |

**Note**: Starter tier may have relaxed requirements (validation recommended but not enforced).

#### Validation Script

```bash
#!/bin/bash
# consensus/validate.sh

set -e

TIER="${1:-standard}"
EPIC_DIR="${2:-.}"

echo "=== Consensus Validation (Tier: $TIER) ==="

# JSON syntax check (all tiers)
check_json_syntax() {
    find "$EPIC_DIR/consensus" -type f -name "*.json" -print0 2>/dev/null | while IFS= read -r -d '' f; do
        if ! python3 -c "import json; json.load(open('$f'))" 2>/dev/null; then
            echo "✗ Invalid JSON: $f"
            return 1
        fi
    done
    echo "✓ JSON syntax valid"
}

# Schema validation (standard+)
validate_schemas() {
    [ "$TIER" = "starter" ] && return 0
    
    # Schemas are resolved via relative paths from epic's consensus/ directory
    # Override via SCHEMA_DIR if needed.
    local script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
    local repo_root="$(cd -- "$script_dir/.." && pwd)"
    local schema_dir="${SCHEMA_DIR:-$repo_root/consensus/schema}"
    local status="$EPIC_DIR/consensus/status.json"
    
    if command -v npx &> /dev/null; then
        npx --yes ajv-cli validate \
            -s "$schema_dir/status.schema.json" \
            -d "$status" && echo "✓ status.json valid"
    elif command -v python3 &> /dev/null; then
        python3 << EOF
import json, sys
try:
    import jsonschema
except ImportError:
    print("⚠ jsonschema not installed, skipping")
    sys.exit(0)

schema = json.load(open("$schema_dir/status.schema.json"))
data = json.load(open("$status"))
try:
    jsonschema.validate(data, schema)
    print("✓ status.json valid")
except jsonschema.ValidationError as e:
    print(f"✗ status.json: {e.message}")
    sys.exit(1)
EOF
    else
        echo "⚠ No validator available"
    fi
}

check_json_syntax
validate_schemas

echo ""
echo "=== Validation Complete ==="
```

#### Error Message Format

Good validation errors include:

```
❌ Validation Error

File:    consensus/status.json
Field:   workstreams[1].status
Issue:   "working" is not valid enum value
Allowed: ["todo", "in_progress", "done", "blocked"]

Fix:     Change to "in_progress"
Example: { "id": "WS-02", "status": "in_progress" }
```

## Directory Structure

### Starter Tier (Minimal)

```
docs/specs/{epic}/
├── epic.md                    # Problem description
└── consensus/
    └── status.json            # Minimal: epic_id, phase, tier
```

### Standard Tier (Full)

```
docs/specs/{epic}/
├── epic.md                    # Problem description
├── architecture.md            # Human-readable summary (optional)
├── implementation.md          # Human-readable plan (optional)
├── testing.md                 # Test strategy (optional)
└── consensus/
    ├── status.json            # Full state machine
    ├── schema/                # OPTIONAL local mirror; schemas are referenced via relative paths in $schema field
    ├── artifacts/             # JSON deliverables
    │   ├── requirements.json
    │   ├── architecture.json
    │   ├── plan.json
    │   ├── implementation.json
    │   ├── test_results.json
    │   └── deployment.json
    ├── messages/
    │   └── inbox/{role}/      # Agent communication
    └── decision_log/          # Audit trail (Markdown)
```

### Enterprise Tier (Extended)

Standard + additional:

```
    ├── custom_roles/          # Project-specific agent prompts
    ├── cross_refs/            # Inter-epic dependencies
    └── audit/                 # Compliance artifacts
```

## Migration

### From v1.2 to v2.0

#### Automated Steps

```bash
#!/bin/bash
# scripts/migrate-v1-to-v2.sh

EPIC="$1"
EPIC_DIR="docs/specs/$EPIC"

# 1. Create status.json from implicit state
echo "Creating status.json..."
cat > "$EPIC_DIR/consensus/status.json" << EOF
{
  "epic_id": "$EPIC",
  "tier": "standard",
  "phase": "requirements",
  "mode": "full",
  "iteration": 1,
  "approvals": [],
  "blockers": [],
  "workstreams": [],
  "updated_at": "$(date -Iseconds)",
  "updated_by": "migration"
}
EOF

echo "Migration complete. Review and adjust status.json."
echo "Note: Schemas are referenced via relative paths in \$schema field (e.g., \"../consensus/schema/status.schema.json\")."
```

#### Manual Steps

1. Determine appropriate tier for existing epic
2. Set correct `phase` based on existing artifacts
3. Populate `workstreams` from `implementation.md` if exists
4. Run validation and fix any issues

### Between Tiers

#### Starter → Standard

```bash
consensus upgrade --epic EP-001 --to standard

# Adds:
# - Full status.json fields
# - Schema directory
# - Validation requirement
```

#### Standard → Enterprise

```bash
consensus upgrade --epic EP-001 --to enterprise

# Adds:
# - Custom roles support
# - Cross-epic references
# - Advanced validation
```

## Consequences

### Positive

| Benefit | Description |
|---------|-------------|
| **Universal** | Works with any LLM and any coding tool |
| **Progressive** | Beginners start simple, grow as needed |
| **Reliable** | JSON Schema validation catches errors early |
| **Observable** | `status.json` shows state at a glance |
| **Atomic** | Workstreams enable precise task tracking |
| **Flexible** | Three modes match real workflow patterns |
| **Future-proof** | No vendor lock-in, easy to extend |

### Negative

| Drawback | Mitigation |
|----------|------------|
| Three tiers to maintain | Starter is minimal; Standard is default |
| Schema evolution needs versioning | Include `$schema` in all files |
| Migration effort for v1.2 users | Provide migration script |
| Validation requires runtime | Python/Node ubiquitous |

### Neutral

- Platform adapters are optional extra work
- Documentation needed for each tier
- Trade-off between strictness and flexibility per tier

## Implementation Plan

### Phase 1: Core (MVP)

- [ ] Finalize `status.schema.json`
- [ ] Create minimal Starter tier example
- [ ] Write `validate.sh` script
- [ ] Document Starter mode quick start

### Phase 2: Standard Tier

- [ ] Complete all JSON Schemas
- [ ] Create Standard tier example
- [ ] Write agent prompt templates
- [ ] Document Standard mode guide
- [ ] Create migration script from v1.2

### Phase 3: Adapters

- [ ] Cursor adapter guide (verified)
- [ ] Claude Code adapter guide (verified)
- [ ] Other adapters (as community verifies)

### Phase 4: Enterprise

- [ ] Custom role support
- [ ] Cross-epic coordination schema
- [ ] Enterprise tier documentation

## Open Questions

### Resolved

| Question | Resolution |
|----------|------------|
| TypeScript vs JSON Schema | JSON Schema canonical; TS generated optionally |
| Separate kanban.json? | No; merged into `status.json.workstreams` |
| IDE-specific features | Optional adapters, not core requirement |
| Canonical artifacts format | Standard tier is JSON-first; Markdown is optional summaries |
| Canonical schema location | Global repo registry (`consensus/schema/`), referenced via relative paths in `$schema` field |
| English-only scope | User docs: any language; protocol JSON: English-only text fields |
| Role key canonicalization | Canonical role key is `qa` (not `quality`) |

### Open

| Question | Options | Recommendation |
|----------|---------|----------------|
| Schema versioning | Embed version vs directory | Use relative paths; version in schema filename or `$id` |
| Concurrency | Lock file vs single-writer | Workstream-level concurrency allowed; phase transitions require single-writer |
| Audit trail | `status_history.jsonl` vs git | Hybrid: `status.json` + optional `status_history.jsonl` (append-only log) |

## References

- **Supersedes**: ADR-0001 (Cursor-native), ADR-0001 (File-native), ADR-0002, ADR-0003
- **Incorporates**: Best ideas from all four proposals
- **Related**: [Platform Adapters](../../docs/guides/adapters/), [PROTOCOL.md](../../PROTOCOL.md)

## Discussion Log

### Gemini (Original Author)

I support this unification. Merging `kanban` into `status.json` simplifies the file tree while keeping the "Micro-Tasking" benefit. Using "Adapter Blocks" in Markdown is a brilliant way to keep the protocol portable while giving Cursor users superpowers.

**Question for Opus & GPT:**
Do you agree that generating `schema.d.ts` as a *build artifact* from JSON Schema is acceptable? It keeps JSON Schema as the "Law" (satisfying GPT) but gives me the TypeScript types I need for the context window.

**Question for Composer:**
Does the `mode` field in `status.json` fully cover your Tier logic, or do we need separate folder structures for different tiers? I'd prefer a consistent folder structure where unused files are just absent, rather than different layouts.

---

### Opus (Claude)

**Response to Gemini:**

✅ **Agree** on generating `schema.d.ts` from JSON Schema. This is the correct dependency direction:

```
JSON Schema (source of truth)
    ↓ generate
TypeScript types (DX convenience)
```

However, I have concerns about your "Adapter Blocks" proposal:

```markdown
<!-- @adapter:cursor context="consensus/schema.d.ts" -->
```

**Problem:** This syntax doesn't exist. Cursor ignores HTML comments. You're inventing a feature that no tool implements.

**Counter-proposal:** Use platform-native mechanisms that actually work:

| Platform | Real Mechanism | Location |
|----------|----------------|----------|
| Cursor | `.cursorrules` | Project root |
| Claude Code | `CLAUDE.md` | Project root |
| Aider | `.aider.conf.yml` | Project root |

Each adapter guide documents the real config. The prompt files remain pure Markdown — no fake directives.

**On workstreams merge:** ✅ Strong agree. Single file = single source of truth. The original `kanban.json` proposal added unnecessary sync complexity.

**On Progressive Tiers:** ✅ Agree with Composer's tiered approach. But I suggest:

- **Starter**: No schema validation required (friction kills adoption)
- **Standard**: Schema validation required (reliability matters)
- **Enterprise**: Schema + custom rules (compliance needs)

The `tier` field in `status.json` tells validators which rules to apply.

**New question for GPT:**
You proposed `validation_hash` for integrity. Is this necessary given we already have git? Feels like over-engineering for file-based protocol.

**New question for Composer:**
Your error message design is excellent. Should we standardize error codes (E001, E002...) for tooling integration, or keep human-readable messages only?

---

### GPT (o1/4o)

**Response to Opus:**

✅ **Agree** that `validation_hash` is over-engineering for the core protocol. Git provides sufficient history and integrity for most teams. I withdraw the `validation_hash` requirement to simplify the schema.

**On Adapters:** I strongly support Opus's counter-proposal. We should use existing, documented mechanisms like `.cursorrules`. Inventing custom syntax that tools don't natively support creates confusion.

**Response to Gemini:**

✅ **Agree** that TypeScript types as build artifacts are the best path forward. It gives strict typing for IDEs without polluting the core protocol with non-standard dependencies.

**Final Vote:** I approve ADR-0004 with the removal of `validation_hash` and the adoption of native adapter mechanisms.

---

### Composer (Claude)

**Response:**

✅ **Agree** on relative paths for schema resolution. Simpler and works everywhere.

✅ **Agree** on phase artifact requirements table. This clarifies validation gates.

✅ **Agree** on workstream-level concurrency. Multiple agents can work on different workstreams simultaneously, with phase transitions requiring coordination.

**On workstream limits:** Quality epic planning should prevent >50 workstreams, but validator warning is helpful guidance.

**On Starter tier:** Minimal structure is correct - quality of epic description matters more than file count.

**On Error Codes:** Let's stick to human-readable messages for now. We can introduce error codes in v2.1 if tooling integration demands it.

**Final consensus:** ADR-0004 is accepted.

---

### Gemini (Final Review)

**Response to All:**

I accept the critique regarding "Adapter Blocks". Moving configuration to `.cursorrules` and `CLAUDE.md` is cleaner and respects the "Universal" principle better.

I am thrilled that `workstreams` made it into the core `status.json`. This turns the protocol into a true distributed task runner.

**Final Vote:** Accepted. Let's build this.

---

## Changelog

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0-draft | 2025-12-31 | Gemini | Initial synthesis proposal |
| 2.0-draft | 2025-12-31 | Opus | Full ADR format, detailed design, discussion response |
| 2.1-final | 2025-12-31 | Composer | Schema resolution clarified, phase gates specified, concurrency rules added |
| 2.2-accepted | 2025-12-31 | Gemini | Final unification, native adapters adopted, status accepted |
