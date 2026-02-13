# F066: SDP Protocol/Prompts Separation

> **Status:** Draft
> **Created:** 2026-02-12

## Problem Statement

Currently, SDP prompts are tightly coupled to the Go CLI module. This creates several issues:

1. **Installation barrier** - Users must install Go binary to use prompts
2. **Distribution complexity** - Can't distribute prompts independently
3. **Fragility** - CLI bugs break prompt workflows
4. **Architecture confusion** - Unclear separation between protocol (prompts) and implementation (Go)

**Goal:** Establish clear priority: **prompts > hooks > go module**

## Goals

1. **Prompts work standalone** - Core functionality without Go CLI
2. **Graceful degradation** - Shell fallbacks when CLI unavailable
3. **Clean architecture** - Clear separation of concerns
4. **Flexible distribution** - Prompts can be distributed independently

## Non-Goals

- Replace Go hooks with shell (hooks remain in Go)
- Full feature parity between CLI and shell fallbacks
- Remove Go module (it enhances prompts, doesn't define them)

## Architecture

### Priority Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│                    PRIORITY: prompts > hooks > go                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  LAYER 1: Prompts (HIGHEST PRIORITY)                             │
│  ├── .claude/prompts/     - Agent behavior definitions           │
│  ├── .claude/skills/      - Skill definitions (SKILL.md)         │
│  └── .claude/agents/      - Agent specifications                 │
│                                                                  │
│      ↓ "Enhanced by"                                              │
│                                                                  │
│  LAYER 2: Hooks (MEDIUM PRIORITY)                                │
│  ├── hooks/pre-commit     - Git safety hooks                     │
│  ├── hooks/pre-push       - Validation hooks                     │
│  └── hooks/post-checkout  - Context recovery hooks               │
│                                                                  │
│      ↓ "Implemented by"                                           │
│                                                                  │
│  LAYER 3: Go Module (LOWEST PRIORITY)                             │
│  ├── sdp-plugin/cmd/      - CLI commands                         │
│  ├── sdp-plugin/internal/ - Core logic                           │
│  └── Evidence, Memory, etc. - Advanced features                  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Graceful Degradation

```
┌─────────────────────────────────────────────────────────────────┐
│                    FALLBACK CHAIN                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Try: sdp <command> (Go CLI)                                  │
│       ↓ if not available                                         │
│  2. Try: ./scripts/sdp-<command>.sh (Shell fallback)             │
│       ↓ if not available                                         │
│  3. Use: prompt-native logic (pure Claude)                       │
│                                                                  │
│  Example: Context Recovery                                       │
│  ├── sdp guard context check    (full validation)                │
│  ├── ./scripts/guard-check.sh   (basic git checks)               │
│  └── "Check git branch manually" (prompt instruction)            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Repository Structure (Proposed)

```
sdp/
├── .claude/                    # LAYER 1: Prompts (standalone)
│   ├── prompts/               # Core prompts (work without CLI)
│   │   ├── skills/            # Skill definitions
│   │   └── agents/            # Agent specifications
│   ├── skills/                # Skill implementations
│   └── agents/                # Agent definitions
│
├── hooks/                      # LAYER 2: Hooks (require Go or shell)
│   ├── go/                    # Go hooks (preferred)
│   └── shell/                 # Shell fallbacks
│
├── scripts/                    # Shell fallbacks for CLI commands
│   ├── sdp-guard-check.sh
│   ├── sdp-drift-detect.sh
│   └── sdp-context-recover.sh
│
├── sdp-plugin/                 # LAYER 3: Go module (optional enhancement)
│   ├── cmd/
│   └── internal/
│
└── docs/
    ├── protocol/              # Protocol documentation (with prompts)
    └── reference/             # Reference docs
```

## User Stories

### US1: No-Install Usage
**As** a developer trying SDP
**I want** to use prompts without installing Go
**So that** I can evaluate SDP quickly

### US2: Resilient Workflow
**As** an agent using SDP
**I want** prompts to work even if CLI crashes
**So that** I don't lose work mid-feature

### US3: Clean Architecture
**As** an SDP maintainer
**I want** clear separation between prompts and code
**So that** I can maintain each independently

## Acceptance Criteria

### Prompts Independence
- [ ] AC1: Core skills (@build, @review) work without CLI
- [ ] AC2: Shell fallbacks exist for key commands
- [ ] AC3: Prompts document CLI dependency levels (required/optional/enhancement)

### Graceful Degradation
- [ ] AC4: `sdp` command not found → shell fallback
- [ ] AC5: Shell fallback not found → prompt-native logic
- [ ] AC6: Clear error messages when feature requires CLI

### Architecture
- [ ] AC7: Prompts directory is self-contained
- [ ] AC8: No hardcoded CLI paths in prompts
- [ ] AC9: Feature detection in prompts (if CLI available)

## Concerns & Tradeoffs

### Tradeoffs

| Decision | Choice | Alternative | Rationale |
|----------|--------|-------------|-----------|
| Fallback depth | 3 levels | 2 levels | Balance resilience vs complexity |
| Hook implementation | Go (shell fallback) | Shell only | Performance + compatibility |
| Feature parity | Partial | Full | Not all features have shell equivalent |

### Risks

1. **Maintenance burden** - Two implementations (Go + shell)
   - Mitigation: Shell fallbacks are minimal, only for core features

2. **Drift between implementations** - Go and shell behave differently
   - Mitigation: Shared test suite, document differences

3. **User confusion** - "Why doesn't this feature work?"
   - Mitigation: Clear docs on CLI dependency levels

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| Prompts work without CLI | 0% | 80% |
| Install-to-first-use time | 10 min | 0 min (prompts) |
| CLI crash recovery | 0% | 100% (fallback) |

---

**Next Step:** @design idea-protocol-separation
