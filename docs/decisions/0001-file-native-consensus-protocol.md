# ADR-0001: File-Native Consensus Protocol v2.0

## Status

Proposed

## Context

The current Consensus Protocol (v1.2) relies on manual orchestration and implicit state management.

### Problems Identified

1. **Fragility**: Agents rely on textual descriptions in `PROTOCOL.md` or loose JSON schemas, leading to structure errors
2. **State Dispersion**: State is implicit in artifacts and messages, requiring reading multiple files to determine current phase or blockers
3. **Context Overload**: Large epics flood context windows with raw messages, increasing cost and latency
4. **Process Rigidity**: Simple bug fixes require the same full chain (Analyst → Architect → Tech Lead → Developer → QA → DevOps) as complex features
5. **Execution Gap**: Tech Lead plans (`implementation.md`) are often too high-level for reliable single-pass implementation

### Design Constraints

The protocol MUST be:

- **Model-agnostic**: Work with any LLM (Claude, GPT, Gemini, Mistral, Llama, Qwen, etc.)
- **Tool-agnostic**: Work with any AI coding tool (Cursor, Claude Code, Codex, Aider, Antigravity, VSCode Copilot, etc.)
- **File-based**: Use only standard file operations as the communication primitive
- **Validatable**: Support runtime validation without IDE-specific features
- **Human-in-the-loop**: Preserve ability for human intervention at any point

### Rejected Alternatives

| Alternative | Reason for Rejection |
|-------------|---------------------|
| External orchestration (LangGraph, CrewAI) | Treats agent process as black box, breaks human-in-the-loop workflow |
| IDE-native features (Cursor-specific directives) | Creates platform lock-in, not portable |
| Database-backed state | Adds infrastructure complexity, breaks file-based simplicity |
| TypeScript-only schemas | Requires compilation, no runtime validation |

## Decision

We will upgrade to a **File-Native Protocol** where files are the only API contract.

### Core Principle

```
Any tool that can READ and WRITE files can participate in the protocol.
No other capability is required.
```

### 1. JSON Schema as Contract ("Schema is Law")

Instead of TypeScript interfaces (compile-time only) or textual descriptions (ambiguous), we use JSON Schema for runtime-validatable contracts.

**Directory Structure:**

```
consensus/schema/
├── status.schema.json      # Epic state machine
├── message.schema.json     # Agent-to-agent communication
├── artifact.schema.json    # Base artifact structure
└── index.json              # Schema registry
```

**Why JSON Schema:**

- Runtime validators exist for every language (ajv, jsonschema, etc.)
- LLMs understand JSON Schema natively
- Can generate TypeScript/Pydantic from schemas (optional bonus)
- Human-readable AND machine-validatable

**Example — Status Schema:**

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "consensus://schema/status",
  "type": "object",
  "required": ["epic_id", "phase", "iteration", "mode", "updated_at", "updated_by"],
  "properties": {
    "epic_id": {
      "type": "string",
      "pattern": "^EP-[A-Z0-9-]+$"
    },
    "phase": {
      "enum": ["requirements", "architecture", "planning", "implementation", "testing", "deployment", "done", "blocked"]
    },
    "iteration": {
      "type": "integer",
      "minimum": 1
    },
    "mode": {
      "enum": ["full", "fast_track", "hotfix"],
      "description": "Execution topology"
    },
    "blockers": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["type", "by", "reason"],
        "properties": {
          "type": { "enum": ["veto", "question", "dependency"] },
          "by": { "type": "string" },
          "reason": { "type": "string" },
          "created_at": { "type": "string", "format": "date-time" }
        }
      }
    },
    "approvals": {
      "type": "array",
      "items": { "type": "string" }
    },
    "workstreams": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "title", "status"],
        "properties": {
          "id": { "type": "string" },
          "title": { "type": "string" },
          "status": { "enum": ["todo", "in_progress", "done", "blocked"] },
          "started_at": { "type": "string", "format": "date-time" },
          "completed_at": { "type": "string", "format": "date-time" }
        }
      }
    },
    "updated_at": { "type": "string", "format": "date-time" },
    "updated_by": { "type": "string" }
  }
}
```

### 2. Centralized State File (`status.json`)

A single mutable file acts as the state machine for the epic.

**Location:** `docs/specs/{epic}/consensus/status.json`

**Ownership Rules:**

- Only ONE agent updates at a time (determined by `phase`)
- Agent MUST read before acting
- Agent MUST validate before writing
- Agent MUST include `updated_at` and `updated_by`

**State Transitions:**

```
requirements → architecture → planning → implementation → testing → deployment → done
     ↑              ↑            ↑             ↑            ↑          ↑
     └──────────────┴────────────┴─────────────┴────────────┴──────────┘
                              (blocked → any previous phase)
```

**Example:**

```json
{
  "epic_id": "EP-AUTH",
  "phase": "implementation",
  "iteration": 2,
  "mode": "full",
  "blockers": [],
  "approvals": ["analyst", "architect", "tech_lead"],
  "workstreams": [
    { "id": "WS-01", "title": "Domain entities", "status": "done", "completed_at": "2025-12-31T10:00:00Z" },
    { "id": "WS-02", "title": "Use cases", "status": "in_progress", "started_at": "2025-12-31T12:00:00Z" },
    { "id": "WS-03", "title": "API endpoints", "status": "todo" }
  ],
  "updated_at": "2025-12-31T14:30:00Z",
  "updated_by": "developer"
}
```

### 3. Agent Prompts as Portable Files

Instead of IDE-specific mechanisms, agent instructions live in plain Markdown files that any LLM can understand.

**Directory:** `consensus/prompts/{role}.md`

**Standard Structure:**

```markdown
# {Role} Agent

## Phase Gate
You can only act when `status.json` shows phase: "{allowed_phase}"
If phase doesn't match, STOP and report the mismatch.

## Context Files
Read these before starting:
1. `consensus/status.json` — current state
2. `consensus/artifacts/{dependency}.json` — previous agent output
3. `consensus/messages/inbox/{role}/` — messages for you

## Your Deliverables
1. Create/update `consensus/artifacts/{output}.json`
2. Send messages to `consensus/messages/inbox/{target}/`
3. Update `consensus/status.json` with new phase and your approval

## Output Validation
- Validate against `consensus/schema/{artifact}.schema.json`
- All text content MUST be in English

## Completion Checklist
- [ ] Read status.json and verified phase
- [ ] Processed all inbox messages
- [ ] Created required artifacts
- [ ] Validated artifacts against schemas
- [ ] Updated status.json
- [ ] Sent messages to next agent(s)
```

**Why Markdown:**

- Every LLM understands Markdown
- Human-readable and reviewable
- Easy to version control
- No special parsing required
- Works with ANY AI tool

### 4. Workstream Tracking (Merged into status.json)

Instead of separate `kanban.json`, workstreams are tracked directly in `status.json`:

```json
{
  "workstreams": [
    {
      "id": "WS-01",
      "title": "Implement domain entities",
      "status": "done",
      "completed_at": "2025-12-31T10:00:00Z"
    },
    {
      "id": "WS-02",
      "title": "Implement use cases",
      "status": "in_progress",
      "started_at": "2025-12-31T12:00:00Z"
    }
  ]
}
```

**Why merged (not separate kanban.json):**

- Single source of truth
- No synchronization issues between files
- Simpler mental model for agents
- Reduces context needed

### 5. Dynamic Execution Topology

Three modes defined in `status.json.mode`:

| Mode | Flow | Use Case |
|------|------|----------|
| `full` | Analyst → Architect → Tech Lead → Developer → QA → DevOps | New features, architectural changes |
| `fast_track` | Developer → QA | Bug fixes, minor changes (≤50 LOC) |
| `hotfix` | Developer → (optional QA) → DevOps | Critical production issues |

**Mode Selection Criteria:**

| Criterion | `full` | `fast_track` | `hotfix` |
|-----------|--------|--------------|----------|
| New functionality | ✓ | ✗ | ✗ |
| Architectural change | ✓ | ✗ | ✗ |
| Lines changed | any | ≤50 | ≤20 |
| New dependencies | allowed | not allowed | not allowed |
| Production incident | ✗ | ✗ | ✓ |

### 6. Validation Layer

Platform-agnostic validation script that works everywhere:

```bash
#!/bin/bash
# consensus/validate.sh

set -e

SCHEMA_DIR="${1:-consensus/schema}"
EPIC_DIR="${2:-.}"

validate_json() {
  local schema="$1"
  local file="$2"
  
  if [ ! -f "$file" ]; then
    echo "⚠ File not found: $file"
    return 0
  fi

  # Try Node.js ajv
  if command -v npx &> /dev/null; then
    npx --yes ajv-cli validate -s "$schema" -d "$file" 2>/dev/null && echo "✓ $file" && return 0
  fi
  
  # Try Python jsonschema
  if command -v python3 &> /dev/null; then
    python3 -c "
import json, sys
try:
    import jsonschema
except ImportError:
    print('⚠ jsonschema not installed, skipping')
    sys.exit(0)
schema = json.load(open('$schema'))
data = json.load(open('$file'))
try:
    jsonschema.validate(data, schema)
    print('✓ $file')
except jsonschema.ValidationError as e:
    print(f'✗ $file: {e.message}')
    sys.exit(1)
" && return 0
  fi
  
  echo "⚠ No validator available (install ajv-cli or python jsonschema)"
  return 0
}

echo "=== Consensus Validation ==="
validate_json "$SCHEMA_DIR/status.schema.json" "$EPIC_DIR/consensus/status.json"

echo ""
echo "=== Validation Complete ==="
```

### 7. Platform Adapters (Optional Layer)

Platform-specific integrations are optional and live in separate guides:

```
docs/guides/adapters/
├── README.md           # Overview and status
├── cursor.md           # ✅ Verified
├── claude-code.md      # ✅ Verified
├── codex.md            # 🚧 Planned
├── antigravity.md      # 🚧 Planned
├── aider.md            # 🚧 Planned
└── vscode-copilot.md   # 🚧 Planned
```

Each adapter documents:

- How to configure the tool
- How to invoke agents
- How to run validation
- Platform-specific tips

**Important:** Adapters are convenience layers. The core protocol works without them.

## Directory Structure (Complete)

```
docs/specs/{epic}/
├── epic.md                          # Epic definition (human-readable)
├── architecture.md                  # Architecture documentation
├── implementation.md                # Implementation plan
├── testing.md                       # Test strategy
├── deployment.md                    # Deployment plan
└── consensus/
    ├── status.json                  # 🆕 Centralized state
    ├── schema/                      # 🆕 JSON Schemas (or symlink to global)
    │   ├── status.schema.json
    │   ├── message.schema.json
    │   └── artifact.schema.json
    ├── prompts/                     # 🆕 Agent instructions
    │   ├── analyst.md
    │   ├── architect.md
    │   ├── tech_lead.md
    │   ├── developer.md
    │   ├── qa.md
    │   └── devops.md
    ├── artifacts/                   # Agent deliverables (JSON)
    │   ├── requirements.json
    │   ├── architecture.json
    │   └── ...
    ├── messages/inbox/{role}/       # Agent communication (JSON)
    └── decision_log/                # Decision history (Markdown)
```

## Consequences

### Positive

| Benefit | Description |
|---------|-------------|
| **Universal Compatibility** | Works with ANY LLM tool that can read/write files |
| **Runtime Validation** | JSON Schema catches errors before they propagate |
| **Single Source of Truth** | `status.json` eliminates state ambiguity |
| **Flexible Execution** | Three modes match real-world workflows |
| **Easy Onboarding** | Portable prompts work anywhere |
| **Future-Proof** | No dependency on any specific vendor or tool |
| **Human Readable** | All state is inspectable in standard formats |

### Negative

| Drawback | Mitigation |
|----------|------------|
| Schema maintenance overhead | Schemas change rarely; version them |
| Validation requires runtime | Node.js or Python are ubiquitous |
| Migration effort from v1.2 | Provide migration script |
| Learning curve for JSON Schema | Provide examples and templates |

### Neutral

- Platform adapters are optional extra work for optimized DX
- Prompts need updating when protocol changes

## Migration from v1.2

### Automated Steps

```bash
# 1. Create schema directory (or use global schemas)
mkdir -p docs/specs/{epic}/consensus/schema
cp consensus/schema/*.schema.json docs/specs/{epic}/consensus/schema/

# 2. Generate status.json from existing artifacts
./scripts/migrate-v1-to-v2.sh {epic}

# 3. Create prompts from templates
cp -r consensus/prompts/ docs/specs/{epic}/consensus/prompts/
```

### Manual Steps

1. Update existing agent instructions to read `status.json` first
2. Add validation step to workflow
3. Choose appropriate `mode` for in-progress epics
4. Review and adjust workstreams in status.json

### Backwards Compatibility

- v1.2 message format remains valid
- Existing artifacts don't need changes
- Only adds new files, doesn't modify existing ones

## References

- **Replaces:** Consensus Protocol v1.2 (implicit state, text-based schemas)
- **Related:** [Platform Adapter Guides](../../docs/guides/adapters/)
- **Rejects:** IDE-specific "native" features in core protocol
- **Inspiration:** 
  - Stigmergy (coordination through environment modification)
  - Unix philosophy (files as universal interface)
  - JSON Schema specification (draft-07)

## Open Questions

1. **Schema versioning:** How to handle schema evolution across epics?
2. **Conflict resolution:** What if two agents update status.json simultaneously?
3. **Audit trail:** Should status.json changes be logged separately?

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 2.0-draft | 2025-12-31 | Initial v2.0 proposal, platform-agnostic rewrite |

