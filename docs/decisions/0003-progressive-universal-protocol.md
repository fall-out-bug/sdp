# ADR-0003: Progressive Universal Protocol v2.0

## Status

Proposed

## Context

The repository targets workflows from small student projects to large multi-week epics. Current proposals (ADR-0001 Cursor-native, ADR-0001 File-native, ADR-0002 Reliability-first) each solve different aspects but miss key requirements:

**Problems:**
- **Onboarding barrier**: Beginners face steep learning curve with full protocol complexity
- **Error diagnosis**: Validation errors don't guide users to fixes
- **Model compatibility**: No clear guidance for different LLM capabilities
- **IDE lock-in risk**: Some proposals favor specific tools
- **Context overload**: Full protocol too heavy for simple tasks

**Requirements:**
- Must work with ANY IDE (Cursor, VS Code, JetBrains, terminal-based tools)
- Must work with ANY model (Claude, GPT, Gemini, local models)
- Must scale from student projects to enterprise epics
- Must provide clear error messages and guidance
- Must support progressive complexity (start simple, add features as needed)

## Decision

We adopt a **Progressive Universal Protocol** with three tiers:

### Tier 1: Starter Mode (Zero-Config)
- Minimal setup: just `epic.md` + auto-generated `status.json`
- Simplified schemas (only required fields)
- Human-friendly error messages with examples
- Works with any model, any IDE
- Perfect for students and small projects

### Tier 2: Standard Mode (Validated)
- Full JSON Schema validation
- Complete `status.json` state machine
- All artifacts in JSON format
- Platform-agnostic validation scripts
- Recommended for most projects

### Tier 3: Advanced Mode (Extended)
- Custom roles and workflows
- Cross-epic coordination
- Advanced validation rules
- IDE-specific adapters (optional)
- For large teams and complex projects

## Core Principles

### 1. Universal Compatibility
```
Any tool that can READ and WRITE files can participate.
No compilation, no special APIs, no vendor lock-in.
```

### 2. Progressive Disclosure
- Start with minimal structure
- Add complexity only when needed
- Clear migration path between tiers

### 3. Schema-First Validation
- JSON Schema as canonical contract
- Runtime validation (no compilation required)
- Human-readable error messages with fix suggestions

### 4. Smart Defaults
- Auto-detect project type and suggest mode
- Pre-filled templates for common scenarios
- Sensible defaults for all optional fields

### 5. Error Guidance
Every validation error includes:
- What's wrong (clear description)
- Where it's wrong (file + field path)
- How to fix (example or template link)
- Why it matters (context)

## Directory Structure

### Starter Mode (Minimal)
```
docs/specs/{epic}/
├── epic.md                    # Simple description
└── consensus/
    ├── status.json            # Auto-generated, minimal fields
    └── artifacts/             # Optional, simple JSON
```

### Standard Mode (Full)
```
docs/specs/{epic}/
├── epic.md
├── architecture.md            # Optional human-readable
├── implementation.md         # Optional human-readable
└── consensus/
    ├── status.json           # Full state machine
    ├── schema/               # JSON Schemas (or symlink to global)
    │   ├── status.schema.json
    │   ├── message.schema.json
    │   └── artifact.schema.json
    ├── artifacts/            # JSON-first
    │   ├── requirements.json
    │   ├── architecture.json
    │   └── plan.json
    ├── messages/inbox/{role}/ # Validated JSON
    └── decision_log/         # Markdown audit trail
```

## Validation Strategy

### Starter Mode Validation
```bash
# Simple checks, friendly messages
consensus validate --mode starter --epic EP-001

# Output:
✓ Epic file exists
✓ Status file is valid JSON
⚠ Missing 'phase' field - add: "phase": "requirements"
💡 Example: see examples/starter/status.json.example
```

### Standard Mode Validation
```bash
# Full schema validation with diagnostics
consensus validate --mode standard --epic EP-001

# Output includes:
- JSON Schema validation results
- Missing required artifacts for current phase
- Phase transition rules check
- Cross-reference integrity
- Actionable error messages
```

## Mode Selection

### Auto-Detection
```python
def suggest_mode(epic_size, project_type, team_size):
    if epic_size == "bug_fix" and lines_changed < 50:
        return "starter"
    elif team_size == 1 and epic_size == "small_feature":
        return "starter"
    elif project_type == "student_project":
        return "starter"
    elif epic_size == "large_feature" or team_size > 3:
        return "standard"
    else:
        return "standard"  # Safe default
```

### Manual Override
```bash
# Force specific mode
consensus init --epic EP-001 --mode starter
consensus init --epic EP-001 --mode standard
consensus init --epic EP-001 --mode advanced
```

## Error Messages Design

### Bad (Current)
```
Error: Validation failed
File: artifacts/requirements.json
```

### Good (Proposed)
```
❌ Validation Error: Missing required field

File:    consensus/artifacts/requirements.json
Field:  acceptance_criteria
Issue:  Field is required but missing

Fix:    Add acceptance_criteria array:
        {
          "acceptance_criteria": [
            "User can login with valid credentials",
            "System rejects invalid credentials"
          ]
        }

Example: examples/starter/requirements.json.example
Docs:    docs/guides/starter-mode.md#acceptance-criteria
```

## Migration Path

### Starter → Standard
```bash
# Upgrade existing epic
consensus upgrade --epic EP-001 --to standard

# Automatically:
# - Adds full schema validation
# - Expands status.json with all fields
# - Creates schema/ directory
# - Validates existing artifacts
```

### Standard → Advanced
```bash
# Enable advanced features
consensus upgrade --epic EP-001 --to advanced

# Adds:
# - Custom role support
# - Cross-epic coordination
# - Advanced validation rules
```

## Model Compatibility

### Tier 1 Models (Basic)
- Local models (Llama, Qwen)
- GPT-3.5-turbo
- Claude Haiku
- **Use**: Starter mode, simple tasks

### Tier 2 Models (Standard)
- GPT-4, GPT-4-turbo
- Claude Sonnet
- Gemini Pro
- **Use**: Standard mode, most tasks

### Tier 3 Models (Advanced)
- Claude Opus
- GPT-4 with advanced reasoning
- **Use**: Advanced mode, complex decisions

**Rule**: Protocol works with ANY model tier, but recommendations guide optimal usage.

## IDE Adapters (Optional)

Adapters provide IDE-specific enhancements but are NOT required:

```
docs/guides/adapters/
├── cursor.md          # Cursor-specific tips
├── vscode.md          # VS Code setup
├── jetbrains.md       # IntelliJ/PyCharm
├── terminal.md        # CLI-only workflow
└── universal.md       # Works everywhere (default)
```

**Key**: Core protocol works without any adapter.

## Implementation Plan

### Phase 1: Starter Mode (MVP)
- [ ] Minimal schema (status.json only)
- [ ] Auto-generation script
- [ ] Friendly error messages
- [ ] Starter examples
- [ ] Quick start guide

### Phase 2: Standard Mode
- [ ] Full JSON Schema definitions
- [ ] Validation script (Python + Node.js)
- [ ] Error diagnostics
- [ ] Migration tool (starter → standard)
- [ ] Standard examples

### Phase 3: Advanced Mode
- [ ] Custom role support
- [ ] Cross-epic coordination
- [ ] IDE adapters
- [ ] Advanced validation rules

## Consequences

### Positive
- **Lower barrier to entry**: Students can start immediately
- **Universal compatibility**: Works everywhere, no lock-in
- **Progressive complexity**: Add features as needed
- **Better error handling**: Users know how to fix issues
- **Model flexibility**: Works with any LLM capability level
- **Future-proof**: Easy to extend without breaking changes

### Negative
- **More complexity in implementation**: Three tiers to maintain
- **Migration tooling needed**: Upgrade paths between tiers
- **Documentation overhead**: Need guides for each tier

### Mitigations
- Starter mode is minimal (easy to maintain)
- Migration tools automate upgrades
- Documentation structured by tier (easy to navigate)

## Alternatives Considered

### Single Protocol for All
**Rejected**: Too complex for beginners, too simple for advanced users

### IDE-Specific Protocols
**Rejected**: Creates lock-in, reduces universality

### Model-Specific Optimizations
**Rejected**: Protocol should work with any model; optimizations are optional

## Open Questions

1. **Schema versioning**: How to handle schema evolution across epics?
2. **Backward compatibility**: Should v1.2 epics auto-upgrade or require manual migration?
3. **Validation performance**: For large epics, should validation be incremental?
4. **Error language**: Should error messages support multiple languages or English-only?

## References

- Builds on: ADR-0001 (File-native), ADR-0002 (Reliability-first)
- Rejects: ADR-0001 (Cursor-native) as core, but keeps as optional adapter
- Related: [Universal Setup Guide](../../docs/guides/universal-setup.md) (to be created)

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 2.0-draft | 2025-01-XX | Initial proposal: Progressive Universal Protocol |

