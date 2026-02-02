# Claude Plugin Distribution

> **Feature ID:** F041
> **Status:** Draft
> **Created:** 2026-02-02

## Problem

Current SDP implementation has significant limitations:

1. **Python Dependency Barrier**
   - Requires Python 3.10+, pip, virtual environment setup
   - Dependencies: click, pydantic, pytest, and 15+ other packages
   - Installation friction for non-Python developers

2. **Language-Specific Assumptions**
   - Workstream validation assumes Python project structure
   - Test coverage commands hardcoded for pytest
   - Pre-commit hooks reference Python-specific tooling

3. **Distribution Friction**
   - Manual pip install or git clone
   - No integration with Claude's plugin ecosystem
   - Users must manage Python environment separately from Claude Code/CLI

## Users

1. **Non-Python Development Teams**
   - Java teams using Maven/Gradle
   - Go teams using go.mod
   - Polyglot teams working across multiple languages
   - **Pain point:** Must install Python stack just for protocol enforcement

2. **Solo Developers on Small Projects**
   - Want protocol enforcement without heavy dependency
   - Prefer lightweight setup via Claude Plugin Marketplace
   - **Pain point:** Python setup feels like overkill for simple projects

3. **DevOps Engineers**
   - Need language-agnostic protocol enforcement
   - Want Git integration without runtime dependencies
   - **Pain point:** Managing Python runtime in CI/CD for protocol checks

## Success Criteria

- [ ] Zero Python dependency for basic SDP protocol usage
- [ ] Installation via Claude Plugin Marketplace (single click)
- [ ] Support for Python, Java, Go projects with language-agnostic validation
- [ ] Backward compatibility with existing Python SDP users
- [ ] Prompts-only workflow works without any binary
- [ ] Optional binary provides CLI convenience

## Goals

### Primary Goal
Transform SDP into a **Claude Plugin** that:
- Distributes via Claude Plugin Marketplace
- Works with any programming language
- Requires zero runtime dependencies for basic usage
- Maintains all protocol enforcement capabilities

### Secondary Goals
1. **Prompts-as-Infrastructure**
   - Extract prompts as standalone delivery mechanism
   - System/agent/skill prompts enforce protocol via AI
   - No code parsing needed - AI validates structure/quality

2. **Optional Binary Tools**
   - Compiled binary (Go/Rust) for convenience features
   - CLI commands (sdp init, sdp doctor)
   - Git hooks installation (pre-commit, pre-push)
   - Protocol enforcement helpers

3. **Split Architecture**
   - **SDP Protocol**: Prompts-only distribution (core product)
   - **SDP Tools**: Optional binary for automation (power users)

### Non-Goals
- Language-specific code parsing (use AI validation instead)
- Enforcing Python as primary runtime (language-agnostic)
- Breaking change for existing users (Python SDP becomes reference)
- Requiring binary for basic usage (prompts-only must work)

## Technical Approach

### Architecture: Split Protocol + Tools

```
SDP Protocol (Core Product)
├── .claude/system/         # Protocol enforcement via system prompts
├── .claude/agents/         # Agent role definitions
├── .claude/skills/         # Skill prompts with protocol validation
└── docs/                   # Protocol documentation (docs/, schemas/)
    Distribution: Claude Plugin Marketplace

SDP Tools (Optional Enhancement)
├── sdp binary              # Compiled CLI (Go/Rust)
│   ├── init               # Project initialization
│   ├── doctor             # Health checks
│   └── hooks              # Git hooks management
└── Language integrations  # Optional language-specific helpers
    ├── Java (Maven/Gradle examples)
    ├── Go (go.mod examples)
    └── Python (pytest examples)
    Distribution: Optional download, not required
```

### Validation Strategy: AI-Based via Prompts

Instead of code parsing, **prompts instruct Claude to validate**:

```markdown
# .claude/system/protocol.md excerpt

When user creates a workstream, you MUST:
1. Validate workstream ID format: PP-FFF-SS
2. Check file size: <200 LOC per file
3. Verify test coverage: ≥80%
4. Ensure type hints present

Validation is performed by:
- Reading file contents
- Checking rules via AI analysis
- Reporting violations with actionable feedback
```

**Examples of Language-Agnostic Validation:**
- File size: Count lines (works for any language)
- Naming: Regex patterns (language-agnostic)
- Structure: Directory existence checks (universal)
- Test coverage: AI detects test files vs production code ratio

### Implementation Phases

#### Phase 1: Protocol Extraction (F041-WS-01)
- Extract prompts as standalone product
- Create Claude Plugin manifest
- Validate prompts work without Python
- **Deliverable**: Installable Claude Plugin

#### Phase 2: Language-Agnostic Validation (F041-WS-02)
- Convert validation logic to prompt-based instructions
- Remove Python-specific assumptions from prompts
- Add language-agnostic examples (Java, Go)
- **Deliverable**: Prompts that work with any language

#### Phase 3: Optional Binary Tools (F041-WS-03)
- Design binary interface (Go or Rust)
- Implement CLI: init, doctor, hooks
- Create package for distribution (brew, apt, etc.)
- **Deliverable**: Optional binary download

#### Phase 4: Marketplace Release (F041-WS-04)
- Submit to Claude Plugin Marketplace
- Create installation documentation
- Gather feedback from early adopters
- **Deliverable**: Public plugin with installation guide

### Migration Path for Existing Users

**Current Python SDP Users:**
1. Continue using Python SDP (supported as reference implementation)
2. Optional: Migrate to plugin (prompts compatible)
3. Binary tools available but not required

**New Users (Non-Python):**
1. Install Claude Plugin from Marketplace (zero setup)
2. Prompts work immediately in Claude Code/CLI
3. Optional: Download binary for CLI convenience

### Backward Compatibility

- **Prompts**: Compatible with existing Python SDP
- **Workstreams**: No format changes (PP-FFF-SS)
- **Documentation**: Existing docs remain valid
- **Migration**: Optional, not forced

## Open Questions

1. **Claude Plugin Format**: What's the exact manifest format for Claude Plugin Marketplace?
2. **Prompt Distribution**: How to package .claude/ structure as installable plugin?
3. **Binary Technology Stack**: Go vs Rust for optional binary?
4. **Validation Granularity**: How much can AI validate vs needing static analysis?
